package mysql

import (
	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/plan"
	"github.com/vatebur/dbinstall/internal/provider"
	"github.com/vatebur/dbinstall/internal/providers/mysqlfamily"
	"github.com/vatebur/dbinstall/internal/spec"
)

type Provider struct{}

func (Provider) Descriptor() provider.Descriptor {
	return provider.Descriptor{
		Name: "mysql", DisplayName: "Oracle MySQL", Family: "mysql",
		Version: "v1alpha1", SupportedSeries: []string{"8.4"},
		Capabilities:     []provider.Capability{provider.ArchiveInstall, provider.PackageInstall, provider.MultiInstance, provider.StaticTuning},
		FoundationStatus: "planning-only",
	}
}

func (Provider) Validate(document spec.DatabaseInstance, machine host.Info) error {
	return mysqlfamily.Validate(document, machine, "8.4")
}

func (Provider) Plan(document spec.DatabaseInstance, _ host.Info) ([]plan.Step, []plan.Warning, error) {
	steps, warnings := mysqlfamily.Plan(document, "Oracle MySQL")
	return steps, warnings, nil
}
