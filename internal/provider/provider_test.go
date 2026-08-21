package provider

import (
	"testing"

	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/spec"
)

type testProvider struct{ name string }

func (p testProvider) Descriptor() Descriptor                        { return Descriptor{Name: p.name} }
func (testProvider) Validate(spec.DatabaseInstance, host.Info) error { return nil }

func TestRegistryRejectsDuplicateAndSorts(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(testProvider{"z"}); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(testProvider{"a"}); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(testProvider{"a"}); err == nil {
		t.Fatal("expected duplicate registration to fail")
	}
	descriptors := registry.Descriptors()
	if descriptors[0].Name != "a" || descriptors[1].Name != "z" {
		t.Fatalf("descriptors are not sorted: %#v", descriptors)
	}
}
