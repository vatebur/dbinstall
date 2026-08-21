package state

import (
	"context"
	"testing"
	"time"
)

func TestStoreLifecycle(t *testing.T) {
	directory := t.TempDir()
	store, err := Open(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	now := time.Unix(100, 0).UTC()
	instance := Instance{Name: "orders", Provider: "mysql", Version: "8.4.6", Status: "planned", SpecDigest: "sha256:test", UpdatedAt: now}
	if err := store.UpsertInstance(context.Background(), instance, map[string]string{"safe": "value"}); err != nil {
		t.Fatal(err)
	}
	instances, err := store.Instances(context.Background(), "orders")
	if err != nil {
		t.Fatal(err)
	}
	if len(instances) != 1 || instances[0].Provider != "mysql" {
		t.Fatalf("instances = %#v", instances)
	}
	operation := Operation{ID: "op1", InstanceName: "orders", Kind: "apply", Status: "running", SpecDigest: "sha256:test", StartedAt: now}
	if err := store.StartOperation(context.Background(), operation); err != nil {
		t.Fatal(err)
	}
	if err := store.FinishOperation(context.Background(), "op1", "failed", "test", now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
}
