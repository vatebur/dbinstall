package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/vatebur/dbinstall/internal/plan"
)

type fakeAction struct {
	step       plan.Step
	satisfied  bool
	applyErr   error
	verifyErr  error
	rolledBack bool
}

func (a *fakeAction) Descriptor() plan.Step               { return a.step }
func (a *fakeAction) Check(context.Context) (bool, error) { return a.satisfied, nil }
func (a *fakeAction) Apply(context.Context) error         { return a.applyErr }
func (a *fakeAction) Verify(context.Context) error        { return a.verifyErr }
func (a *fakeAction) Rollback(context.Context) error      { a.rolledBack = true; return nil }

func TestExecutorSkipsSatisfiedAndRollsBackVerifiedFailure(t *testing.T) {
	first := &fakeAction{step: plan.Step{ID: "first", Implemented: true}, satisfied: true}
	second := &fakeAction{step: plan.Step{ID: "second", DependsOn: []string{"first"}, Implemented: true}}
	third := &fakeAction{step: plan.Step{ID: "third", DependsOn: []string{"second"}, Implemented: true}, verifyErr: errors.New("unhealthy")}
	result, err := (Executor{}).Run(context.Background(), []Action{first, second, third})
	if err == nil {
		t.Fatal("expected verification failure")
	}
	if first.rolledBack {
		t.Fatal("satisfied action must not be rolled back")
	}
	if !second.rolledBack || !third.rolledBack {
		t.Fatal("applied actions were not rolled back")
	}
	if len(result.RolledBack) != 2 || result.RolledBack[0] != "third" {
		t.Fatalf("rollback order = %#v", result.RolledBack)
	}
}

func TestExecutorRejectsForwardDependency(t *testing.T) {
	action := &fakeAction{step: plan.Step{ID: "first", DependsOn: []string{"later"}, Implemented: true}}
	if _, err := (Executor{}).Run(context.Background(), []Action{action}); err == nil {
		t.Fatal("expected graph validation failure")
	}
}
