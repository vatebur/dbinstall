package provider

import (
	"fmt"
	"sort"
	"sync"

	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/spec"
)

type Capability string

const (
	ArchiveInstall Capability = "archive-install"
	PackageInstall Capability = "package-install"
	MultiInstance  Capability = "multi-instance"
	StaticTuning   Capability = "static-tuning"
	Observation    Capability = "runtime-observation"
)

type Descriptor struct {
	Name             string       `json:"name"`
	DisplayName      string       `json:"displayName"`
	Family           string       `json:"family"`
	Version          string       `json:"providerVersion"`
	Capabilities     []Capability `json:"capabilities"`
	SupportedSeries  []string     `json:"supportedSeries"`
	FoundationStatus string       `json:"foundationStatus"`
}

type Provider interface {
	Descriptor() Descriptor
	Validate(spec.DatabaseInstance, host.Info) error
}

type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]Provider)}
}

func (r *Registry) Register(value Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	name := value.Descriptor().Name
	if name == "" {
		return fmt.Errorf("provider name is empty")
	}
	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %q already registered", name)
	}
	r.providers[name] = value
	return nil
}

func (r *Registry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %q is not registered", name)
	}
	return value, nil
}

func (r *Registry) Descriptors() []Descriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	values := make([]Descriptor, 0, len(r.providers))
	for _, value := range r.providers {
		values = append(values, value.Descriptor())
	}
	sort.Slice(values, func(i, j int) bool { return values[i].Name < values[j].Name })
	return values
}
