package planner

import (
	"testing"
	"time"

	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/providers/builtin"
	"github.com/vatebur/dbinstall/internal/spec"
)

func TestMySQLFoundationPlanIsSafeAndNotExecutable(t *testing.T) {
	registry, err := builtin.Registry()
	if err != nil {
		t.Fatal(err)
	}
	document := spec.DatabaseInstance{
		APIVersion: spec.APIVersion, Kind: spec.Kind, Metadata: spec.Metadata{Name: "mysql"},
		Spec: spec.InstanceSpec{
			Provider: "mysql", Version: "8.4.6", Installation: spec.Installation{Method: "archive", Source: "online"},
			Instance: spec.Runtime{Port: 3306, BindAddress: "0.0.0.0", Paths: spec.Paths{Data: "/data", Logs: "/logs", Temp: "/temp"}},
			Workload: spec.Workload{Profile: "oltp", Durability: "strict"},
		},
	}
	machine := host.Info{OSFamily: "debian", Certified: true}
	result, err := (Planner{Registry: registry, Now: func() time.Time { return time.Unix(0, 0) }}).Build(document, machine)
	if err != nil {
		t.Fatal(err)
	}
	if result.Executable {
		t.Fatal("foundation plan must not be executable")
	}
	if len(result.Steps) < 5 {
		t.Fatalf("expected a meaningful lifecycle, got %#v", result.Steps)
	}
	if len(result.Warnings) == 0 || result.Warnings[0].Code != "network.all-interfaces" {
		t.Fatalf("missing all-interface warning: %#v", result.Warnings)
	}
}
