package planner

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/plan"
	"github.com/vatebur/dbinstall/internal/provider"
	"github.com/vatebur/dbinstall/internal/spec"
)

type Clock func() time.Time

type Planner struct {
	Registry *provider.Registry
	Now      Clock
}

func (p Planner) Build(document spec.DatabaseInstance, machine host.Info) (plan.Document, error) {
	value, err := p.Registry.Get(document.Spec.Provider)
	if err != nil {
		return plan.Document{}, err
	}
	if err := value.Validate(document, machine); err != nil {
		return plan.Document{}, fmt.Errorf("provider validation: %w", err)
	}
	steps, warnings, err := value.Plan(document, machine)
	if err != nil {
		return plan.Document{}, fmt.Errorf("provider plan: %w", err)
	}
	warnings = append(commonWarnings(document, machine), warnings...)
	digest, err := specDigest(document)
	if err != nil {
		return plan.Document{}, err
	}
	now := p.Now
	if now == nil {
		now = time.Now
	}
	return plan.Document{
		APIVersion: "dbinstall.vatebur.dev/plan/v1alpha1", OperationID: operationID(),
		GeneratedAt: now().UTC(), InstanceName: document.Metadata.Name, SpecDigest: digest,
		Provider: document.Spec.Provider, Version: document.Spec.Version, Host: machine,
		Executable: allImplemented(steps), Warnings: warnings, Steps: steps,
	}, nil
}

func commonWarnings(document spec.DatabaseInstance, machine host.Info) []plan.Warning {
	var warnings []plan.Warning
	if document.Spec.Instance.BindAddress == "0.0.0.0" {
		warnings = append(warnings, plan.Warning{Code: "network.all-interfaces", Risk: plan.RiskHigh, Message: "database will listen on every IPv4 interface; verify firewall or security-group policy before apply"})
	}
	if !machine.Certified {
		warnings = append(warnings, plan.Warning{Code: "host.not-certified", Risk: plan.RiskHigh, Message: machine.Certification})
	}
	if document.Spec.SystemTuning {
		warnings = append(warnings, plan.Warning{Code: "host.system-tuning", Risk: plan.RiskMedium, Message: "system tuning was requested and requires explicit apply-time authorization"})
	}
	if len(document.Spec.ExtraConfig) > 0 {
		warnings = append(warnings, plan.Warning{Code: "config.unverified", Risk: plan.RiskMedium, Message: "extraConfig contains provider-unverified database parameters"})
	}
	return warnings
}

func specDigest(document spec.DatabaseInstance) (string, error) {
	data, err := json.Marshal(document)
	if err != nil {
		return "", fmt.Errorf("serialize normalized specification: %w", err)
	}
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func operationID() string {
	var value [12]byte
	if _, err := rand.Read(value[:]); err != nil {
		panic(fmt.Sprintf("secure random source unavailable: %v", err))
	}
	return hex.EncodeToString(value[:])
}

func allImplemented(steps []plan.Step) bool {
	if len(steps) == 0 {
		return false
	}
	for _, step := range steps {
		if !step.Implemented {
			return false
		}
	}
	return true
}
