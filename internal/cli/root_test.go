package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestProvidersJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := NewRootCommand(&stdout, &stderr)
	command.SetArgs([]string{"providers", "--output", "json"})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, provider := range []string{"greatsql", "mysql", "postgresql"} {
		if !strings.Contains(stdout.String(), `"name": "`+provider+`"`) {
			t.Fatalf("missing provider %s in %s", provider, stdout.String())
		}
	}
}
