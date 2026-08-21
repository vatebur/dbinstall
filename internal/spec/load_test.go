package spec

import (
	"errors"
	"strings"
	"testing"
)

func TestLoadDefaultsAndRejectsUnknownFields(t *testing.T) {
	valid := `apiVersion: dbinstall.vatebur.dev/v1alpha1
kind: DatabaseInstance
metadata:
  name: orders
spec:
  provider: mysql
  version: 8.4.6
  installation: {}
  instance:
    paths: {}
  workload: {}
  credentials: {}
`
	document, err := Load([]byte(valid))
	if err != nil {
		t.Fatal(err)
	}
	if document.Spec.Instance.BindAddress != "0.0.0.0" {
		t.Fatalf("bind address = %q", document.Spec.Instance.BindAddress)
	}
	if document.Spec.Instance.Port != 3306 {
		t.Fatalf("port = %d", document.Spec.Instance.Port)
	}
	if document.Spec.Workload.Profile != "development" {
		t.Fatalf("profile = %q", document.Spec.Workload.Profile)
	}

	_, err = Load([]byte(valid + "unknown: true\n"))
	if err == nil || !strings.Contains(err.Error(), "field unknown not found") {
		t.Fatalf("expected unknown field error, got %v", err)
	}
}

func TestValidateReturnsAllRelevantProblems(t *testing.T) {
	document := DatabaseInstance{
		APIVersion: APIVersion,
		Kind:       Kind,
		Metadata:   Metadata{Name: "Bad_Name"},
		Spec: InstanceSpec{
			Provider:     "mysql",
			Version:      "latest",
			Installation: Installation{Method: "archive", Source: "offline", ArtifactPath: "relative"},
			Instance: Runtime{Port: 80, BindAddress: "all", Paths: Paths{
				Data: "/data", Logs: "/data", Temp: "/tmp/dbinstall",
			}},
			Workload: Workload{Profile: "unknown", Durability: "strict"},
		},
	}
	err := Validate(document)
	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if len(validation.Problems) < 6 {
		t.Fatalf("expected several problems, got %#v", validation.Problems)
	}
}

func TestParseBytes(t *testing.T) {
	value, err := ParseBytes("8GiB")
	if err != nil {
		t.Fatal(err)
	}
	if value != 8<<30 {
		t.Fatalf("value = %d", value)
	}
	if _, err := ParseBytes("8GB"); err == nil {
		t.Fatal("expected SI quantity to be rejected")
	}
}
