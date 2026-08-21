package spec

const (
	APIVersion = "dbinstall.vatebur.dev/v1alpha1"
	Kind       = "DatabaseInstance"
)

type DatabaseInstance struct {
	APIVersion string       `yaml:"apiVersion" json:"apiVersion"`
	Kind       string       `yaml:"kind" json:"kind"`
	Metadata   Metadata     `yaml:"metadata" json:"metadata"`
	Spec       InstanceSpec `yaml:"spec" json:"spec"`
}

type Metadata struct {
	Name string `yaml:"name" json:"name"`
}

type InstanceSpec struct {
	Provider         string                 `yaml:"provider" json:"provider"`
	Version          string                 `yaml:"version" json:"version"`
	Installation     Installation           `yaml:"installation" json:"installation"`
	Instance         Runtime                `yaml:"instance" json:"instance"`
	Workload         Workload               `yaml:"workload" json:"workload"`
	Credentials      Credentials            `yaml:"credentials" json:"credentials"`
	ConfigOverrides  map[string]any         `yaml:"configOverrides,omitempty" json:"configOverrides,omitempty"`
	ExtraConfig      map[string]any         `yaml:"extraConfig,omitempty" json:"extraConfig,omitempty"`
	SystemTuning     bool                   `yaml:"systemTuning,omitempty" json:"systemTuning,omitempty"`
	ProviderSettings map[string]interface{} `yaml:"providerSettings,omitempty" json:"providerSettings,omitempty"`
}

type Installation struct {
	Method       string `yaml:"method" json:"method"`
	Source       string `yaml:"source" json:"source"`
	ArtifactPath string `yaml:"artifactPath,omitempty" json:"artifactPath,omitempty"`
}

type Runtime struct {
	Port        int    `yaml:"port" json:"port"`
	BindAddress string `yaml:"bindAddress" json:"bindAddress"`
	Paths       Paths  `yaml:"paths" json:"paths"`
}

type Paths struct {
	Data string `yaml:"data" json:"data"`
	Logs string `yaml:"logs" json:"logs"`
	Temp string `yaml:"temp" json:"temp"`
}

type Workload struct {
	Profile             string `yaml:"profile" json:"profile"`
	MemoryLimit         string `yaml:"memoryLimit,omitempty" json:"memoryLimit,omitempty"`
	ExpectedConnections int    `yaml:"expectedConnections,omitempty" json:"expectedConnections,omitempty"`
	Durability          string `yaml:"durability,omitempty" json:"durability,omitempty"`
}

type Credentials struct {
	AdminPassword string `yaml:"adminPassword,omitempty" json:"adminPassword,omitempty"`
}
