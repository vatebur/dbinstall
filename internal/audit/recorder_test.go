package audit

import (
	"os"
	"strings"
	"testing"

	"github.com/vatebur/dbinstall/internal/plan"
)

func TestWritePlanPublishesJSON(t *testing.T) {
	path, err := WritePlan(t.TempDir(), plan.Document{APIVersion: "test", OperationID: "abc", InstanceName: "orders"})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"operationID": "abc"`) {
		t.Fatalf("unexpected plan: %s", data)
	}
}
