package postgresql

import (
	"fmt"

	"github.com/vatebur/dbinstall/internal/host"
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
