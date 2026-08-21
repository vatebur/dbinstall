package spec

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	namePattern    = regexp.MustCompile(`^[a-z0-9](?:[-a-z0-9]{0,61}[a-z0-9])?$`)
	versionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`)
	quantity       = regexp.MustCompile(`^([1-9][0-9]*)(KiB|MiB|GiB|TiB)$`)
)

type ValidationError struct {
	Problems []string `json:"problems"`
}

func (e *ValidationError) Error() string {
	return "invalid specification: " + strings.Join(e.Problems, "; ")
}

func Validate(document DatabaseInstance) error {
	var problems []string
	add := func(ok bool, problem string) {
		if !ok {
			problems = append(problems, problem)
		}
	}

	add(document.APIVersion == APIVersion, fmt.Sprintf("apiVersion must be %q", APIVersion))
	add(document.Kind == Kind, fmt.Sprintf("kind must be %q", Kind))
	add(namePattern.MatchString(document.Metadata.Name), "metadata.name must be a lowercase DNS label")
	add(oneOf(document.Spec.Provider, "mysql", "greatsql", "postgresql"), "spec.provider is unsupported")
	add(versionPattern.MatchString(document.Spec.Version), "spec.version must be a complete x.y.z version")
	add(oneOf(document.Spec.Installation.Method, "archive", "package"), "installation.method must be archive or package")
	add(oneOf(document.Spec.Installation.Source, "online", "offline"), "installation.source must be online or offline")
	if document.Spec.Installation.Source == "offline" {
		add(filepath.IsAbs(document.Spec.Installation.ArtifactPath), "offline installation requires an absolute artifactPath")
	}
	add(document.Spec.Instance.Port >= 1024 && document.Spec.Instance.Port <= 65535, "instance.port must be between 1024 and 65535")
	add(net.ParseIP(document.Spec.Instance.BindAddress) != nil, "instance.bindAddress must be an IP address")

	paths := []string{document.Spec.Instance.Paths.Data, document.Spec.Instance.Paths.Logs, document.Spec.Instance.Paths.Temp}
	seen := make(map[string]bool, len(paths))
	for _, path := range paths {
		add(filepath.IsAbs(path), "instance paths must be absolute")
		clean := filepath.Clean(path)
		add(clean != "/", "instance paths cannot be the filesystem root")
		add(!seen[clean], "instance data, log, and temp paths must be distinct")
		seen[clean] = true
	}
	add(oneOf(document.Spec.Workload.Profile, "development", "oltp", "olap", "custom"), "workload.profile is unsupported")
	add(oneOf(document.Spec.Workload.Durability, "strict", "balanced", "relaxed"), "workload.durability is unsupported")
	add(document.Spec.Workload.ExpectedConnections >= 0, "workload.expectedConnections cannot be negative")
	if document.Spec.Workload.MemoryLimit != "" {
		_, err := ParseBytes(document.Spec.Workload.MemoryLimit)
		add(err == nil, "workload.memoryLimit must use KiB, MiB, GiB, or TiB")
	}
	if reference := document.Spec.Credentials.AdminPassword; reference != "" {
		add(validSecretReference(reference), "credentials.adminPassword must use env://NAME or an absolute file:///path")
	}
	for key := range document.Spec.ExtraConfig {
		add(!reservedConfigKey(key), fmt.Sprintf("extraConfig.%s is owned by the instance model", key))
	}

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}
	return nil
}

func ParseBytes(value string) (uint64, error) {
	match := quantity.FindStringSubmatch(value)
	if match == nil {
		return 0, errors.New("invalid IEC quantity")
	}
	n, _ := strconv.ParseUint(match[1], 10, 64)
	shift := map[string]uint{"KiB": 10, "MiB": 20, "GiB": 30, "TiB": 40}[match[2]]
	if n > ^uint64(0)>>shift {
		return 0, errors.New("quantity overflows uint64")
	}
	return n << shift, nil
}

func oneOf(value string, choices ...string) bool {
	for _, choice := range choices {
		if value == choice {
			return true
		}
	}
	return false
}

func validSecretReference(reference string) bool {
	if strings.HasPrefix(reference, "env://") {
		name := strings.TrimPrefix(reference, "env://")
		return regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(name)
	}
	if strings.HasPrefix(reference, "file://") {
		return filepath.IsAbs(strings.TrimPrefix(reference, "file://"))
	}
	return false
}

func reservedConfigKey(key string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(key, "-", "_"))
	reserved := map[string]bool{
		"port": true, "bind_address": true, "datadir": true, "data_directory": true,
		"socket": true, "pid_file": true, "user": true, "log_error": true,
	}
	return reserved[normalized]
}
