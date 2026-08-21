package postgresql

import (
	"fmt"

	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/plan"
	"github.com/vatebur/dbinstall/internal/provider"
	"github.com/vatebur/dbinstall/internal/spec"
)

type Provider struct{}

func (Provider) Descriptor() provider.Descriptor {
	return provider.Descriptor{
		Name: "postgresql", DisplayName: "PostgreSQL", Family: "postgresql",
		Version: "v1alpha1", SupportedSeries: []string{"pending-stable-selection"},
		Capabilities:     []provider.Capability{provider.ArchiveInstall, provider.PackageInstall},
		FoundationStatus: "descriptor-only",
	}
}

func (Provider) Validate(_ spec.DatabaseInstance, machine host.Info) error {
	if machine.OSFamily != "rhel" && machine.OSFamily != "debian" {
		return fmt.Errorf("unsupported OS family %q", machine.OSFamily)
	}
	return nil
}

func (Provider) Plan(_ spec.DatabaseInstance, _ host.Info) ([]plan.Step, []plan.Warning, error) {
	steps := []plan.Step{
		{ID: "preflight.host", Description: "verify host support, paths, ports, privileges, and conflicts", Risk: plan.RiskInfo, Privilege: plan.PrivilegeUser, Implemented: true, Timeout: "30s"},
		{ID: "provider.pending", Description: "PostgreSQL installation steps require stable-series artifact selection", DependsOn: []string{"preflight.host"}, Risk: plan.RiskInfo, Privilege: plan.PrivilegeUser, Implemented: false, Timeout: "0s"},
	}
	warnings := []plan.Warning{{Code: "provider.descriptor-only", Risk: plan.RiskHigh, Message: "PostgreSQL is a descriptor-only foundation provider and cannot be applied"}}
	return steps, warnings, nil
}
