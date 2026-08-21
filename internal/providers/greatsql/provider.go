package greatsql

import (
	"github.com/vatebur/dbinstall/internal/host"
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
		return mysqlfamily.Validate(document, machine, "unsupported")
	}
	return nil
}
