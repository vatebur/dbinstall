package executor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vatebur/dbinstall/internal/plan"
)

type Action interface {
	Descriptor() plan.Step
	Check(context.Context) (satisfied bool, err error)
	Apply(context.Context) error
	Verify(context.Context) error
	Rollback(context.Context) error
}

type Event struct {
	StepID  string    `json:"stepID"`
	Phase   string    `json:"phase"`
	Time    time.Time `json:"time"`
	Message string    `json:"message,omitempty"`
}

type Sink interface{ Record(Event) error }

type SinkFunc func(Event) error

func (f SinkFunc) Record(event Event) error { return f(event) }

type Result struct {
	Applied    []string `json:"applied"`
	Skipped    []string `json:"skipped"`
	RolledBack []string `json:"rolledBack"`
}

type ExecutionError struct {
	StepID          string
	Phase           string
	Cause           error
	RollbackFailure error
}

func (e *ExecutionError) Error() string {
	if e.RollbackFailure != nil {
		return fmt.Sprintf("step %s failed during %s: %v; rollback was incomplete: %v", e.StepID, e.Phase, e.Cause, e.RollbackFailure)
	}
	return fmt.Sprintf("step %s failed during %s: %v", e.StepID, e.Phase, e.Cause)
}

func (e *ExecutionError) Unwrap() error { return e.Cause }

type Executor struct {
	Sink Sink
	Now  func() time.Time
}

func (e Executor) Run(ctx context.Context, actions []Action) (Result, error) {
	if err := validateGraph(actions); err != nil {
		return Result{}, err
	}
	now := e.Now
	if now == nil {
		now = time.Now
	}
	sink := e.Sink
	if sink == nil {
		sink = SinkFunc(func(Event) error { return nil })
	}
	var result Result
	var applied []Action
	for _, action := range actions {
		descriptor := action.Descriptor()
		if !descriptor.Implemented {
			return result, &ExecutionError{StepID: descriptor.ID, Phase: "preflight", Cause: errors.New("step is not implemented")}
		}
		if err := sink.Record(Event{StepID: descriptor.ID, Phase: "check", Time: now().UTC()}); err != nil {
			return result, &ExecutionError{StepID: descriptor.ID, Phase: "journal", Cause: err}
		}
		satisfied, err := action.Check(ctx)
		if err != nil {
			return result, &ExecutionError{StepID: descriptor.ID, Phase: "check", Cause: err}
		}
		if satisfied {
			result.Skipped = append(result.Skipped, descriptor.ID)
			continue
		}
		if err := action.Apply(ctx); err != nil {
			rollbackErr := rollback(ctx, applied, &result, sink, now)
			return result, &ExecutionError{StepID: descriptor.ID, Phase: "apply", Cause: err, RollbackFailure: rollbackErr}
		}
		applied = append(applied, action)
		result.Applied = append(result.Applied, descriptor.ID)
		if err := action.Verify(ctx); err != nil {
			rollbackErr := rollback(ctx, applied, &result, sink, now)
			return result, &ExecutionError{StepID: descriptor.ID, Phase: "verify", Cause: err, RollbackFailure: rollbackErr}
		}
		if err := sink.Record(Event{StepID: descriptor.ID, Phase: "verified", Time: now().UTC()}); err != nil {
			rollbackErr := rollback(ctx, applied, &result, sink, now)
			return result, &ExecutionError{StepID: descriptor.ID, Phase: "journal", Cause: err, RollbackFailure: rollbackErr}
		}
	}
	return result, nil
}

func validateGraph(actions []Action) error {
	seen := make(map[string]bool, len(actions))
	for _, action := range actions {
		step := action.Descriptor()
		if step.ID == "" {
			return errors.New("step ID is empty")
		}
		if seen[step.ID] {
			return fmt.Errorf("duplicate step ID %q", step.ID)
		}
		for _, dependency := range step.DependsOn {
			if !seen[dependency] {
				return fmt.Errorf("step %q dependency %q is missing or ordered after it", step.ID, dependency)
			}
		}
		seen[step.ID] = true
	}
	return nil
}

func rollback(ctx context.Context, actions []Action, result *Result, sink Sink, now func() time.Time) error {
	var failures []error
	for index := len(actions) - 1; index >= 0; index-- {
		step := actions[index].Descriptor()
		if err := actions[index].Rollback(ctx); err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", step.ID, err))
			continue
		}
		result.RolledBack = append(result.RolledBack, step.ID)
		if err := sink.Record(Event{StepID: step.ID, Phase: "rolled-back", Time: now().UTC()}); err != nil {
			failures = append(failures, fmt.Errorf("journal %s: %w", step.ID, err))
		}
	}
	return errors.Join(failures...)
}
