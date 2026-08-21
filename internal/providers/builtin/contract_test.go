package builtin

import (
	"testing"

	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/spec"
)

func TestBuiltInProviderContracts(t *testing.T) {
	registry, err := Registry()
	if err != nil {
		t.Fatal(err)
	}
	machine := host.Info{OSFamily: "debian", Architecture: "amd64", InitSystem: "systemd", Certified: true}
	versions := map[string]string{"mysql": "8.4.6", "greatsql": "8.0.32", "postgresql": "17.6"}
	methods := map[string]string{"mysql": "archive", "greatsql": "archive", "postgresql": "package"}
	for _, descriptor := range registry.Descriptors() {
		t.Run(descriptor.Name, func(t *testing.T) {
			if descriptor.DisplayName == "" || descriptor.Family == "" || descriptor.Version == "" {
				t.Fatalf("incomplete descriptor: %#v", descriptor)
			}
			value, err := registry.Get(descriptor.Name)
			if err != nil {
				t.Fatal(err)
			}
			document := spec.DatabaseInstance{
				APIVersion: spec.APIVersion, Kind: spec.Kind, Metadata: spec.Metadata{Name: descriptor.Name},
				Spec: spec.InstanceSpec{
					Provider: descriptor.Name, Version: versions[descriptor.Name],
					Installation: spec.Installation{Method: methods[descriptor.Name], Source: "online"},
					Instance:     spec.Runtime{Port: 3306, BindAddress: "0.0.0.0", Paths: spec.Paths{Data: "/data/" + descriptor.Name, Logs: "/logs/" + descriptor.Name, Temp: "/tmp/" + descriptor.Name}},
					Workload:     spec.Workload{Profile: "development", Durability: "strict"},
				},
			}
			if err := value.Validate(document, machine); err != nil {
				t.Fatal(err)
			}
			steps, _, err := value.Plan(document, machine)
			if err != nil {
				t.Fatal(err)
			}
			seen := map[string]bool{}
			for _, step := range steps {
				if step.ID == "" || seen[step.ID] {
					t.Fatalf("invalid step ID %q", step.ID)
				}
				for _, dependency := range step.DependsOn {
					if !seen[dependency] {
						t.Fatalf("step %s has missing or forward dependency %s", step.ID, dependency)
					}
				}
				seen[step.ID] = true
			}
		})
	}
}
