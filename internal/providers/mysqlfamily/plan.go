package mysqlfamily

import (
	"fmt"
	"path/filepath"

	"github.com/vatebur/dbinstall/internal/plan"
	"github.com/vatebur/dbinstall/internal/spec"
)

func Plan(document spec.DatabaseInstance, displayName string) ([]plan.Step, []plan.Warning) {
	name := document.Metadata.Name
	config := filepath.Join("/etc/dbinstall/instances", name, "my.cnf")
	service := "dbinstall-" + name + ".service"
	steps := []plan.Step{
		{ID: "preflight.host", Description: "verify host support, paths, ports, privileges, and conflicts", Risk: plan.RiskInfo, Privilege: plan.PrivilegeUser, Implemented: true, Timeout: "30s"},
		{ID: "artifact.resolve", Description: fmt.Sprintf("resolve and verify the pinned %s distribution", displayName), DependsOn: []string{"preflight.host"}, Risk: plan.RiskLow, Privilege: plan.PrivilegeUser, Implemented: false, Timeout: "15m"},
		{ID: "distribution.install", Description: "install the immutable database distribution", DependsOn: []string{"artifact.resolve"}, Risk: plan.RiskMedium, Privilege: plan.PrivilegeRoot, Implemented: false, Timeout: "10m", OwnedResources: []plan.Resource{{Kind: "distribution", Name: document.Spec.Provider + "-" + document.Spec.Version}}},
		{ID: "instance.directories", Description: "create dedicated data, log, temp, and configuration directories", DependsOn: []string{"distribution.install"}, Risk: plan.RiskMedium, Privilege: plan.PrivilegeRoot, Implemented: false, Timeout: "1m", OwnedResources: []plan.Resource{{Kind: "directory", Name: document.Spec.Instance.Paths.Data}, {Kind: "directory", Name: document.Spec.Instance.Paths.Logs}, {Kind: "directory", Name: document.Spec.Instance.Paths.Temp}}},
		{ID: "instance.configure", Description: "render, validate, and atomically publish instance configuration", DependsOn: []string{"instance.directories"}, Risk: plan.RiskMedium, Privilege: plan.PrivilegeRoot, Implemented: false, Timeout: "2m", OwnedResources: []plan.Resource{{Kind: "config", Name: config}}},
		{ID: "instance.initialize", Description: "initialize a new empty data directory and local-only administrative account", DependsOn: []string{"instance.configure"}, Risk: plan.RiskHigh, Privilege: plan.PrivilegeRoot, Implemented: false, Timeout: "10m"},
		{ID: "service.install", Description: "install and enable the dedicated systemd unit", DependsOn: []string{"instance.initialize"}, Risk: plan.RiskMedium, Privilege: plan.PrivilegeRoot, Implemented: false, Timeout: "2m", OwnedResources: []plan.Resource{{Kind: "systemd-unit", Name: service}}},
		{ID: "instance.verify", Description: "verify startup, local administration, port, TLS, and basic read/write", DependsOn: []string{"service.install"}, Risk: plan.RiskInfo, Privilege: plan.PrivilegeRoot, Implemented: false, Timeout: "5m"},
	}
	return steps, nil
}
