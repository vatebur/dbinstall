package greatsql

import (
	"fmt"

	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/plan"
	"github.com/vatebur/dbinstall/internal/provider"
	"github.com/vatebur/dbinstall/internal/providers/mysqlfamily"
	"github.com/vatebur/dbinstall/internal/spec"
)

type Provider struct{}

func (Provider) Descriptor() provider.Descriptor {
	return provider.Descriptor{
		Name: "greatsql", DisplayName: "GreatSQL", Family: "mysql",
		Version: "v1alpha1", SupportedSeries: []string{"pending-artifact-selection"},
		Capabilities:     []provider.Capability{provider.ArchiveInstall, provider.MultiInstance, provider.StaticTuning},
		FoundationStatus: "descriptor-only",
	}
}

func (Provider) Validate(document spec.DatabaseInstance, machine host.Info) error {
	if machine.OSFamily != "rhel" && machine.OSFamily != "debian" {
		return fmt.Errorf("unsupported OS family %q", machine.OSFamily)
	}
	if document.Spec.Installation.Method != "archive" {
		return fmt.Errorf("GreatSQL foundation provider supports archive planning only")
	}
	return nil
}

func (Provider) Plan(document spec.DatabaseInstance, _ host.Info) ([]plan.Step, []plan.Warning, error) {
	steps, warnings := mysqlfamily.Plan(document, "GreatSQL")
	return steps, warnings, nil
}
