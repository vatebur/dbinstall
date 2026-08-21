package secrets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalResolver(t *testing.T) {
	t.Setenv("DBINSTALL_TEST_SECRET", "from-env")
	value, err := (LocalResolver{}).Resolve("env://DBINSTALL_TEST_SECRET")
	if err != nil || string(value) != "from-env" {
		t.Fatalf("environment secret = %q, %v", value, err)
	}

	path := filepath.Join(t.TempDir(), "secret")
	if err := os.WriteFile(path, []byte("from-file\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	value, err = (LocalResolver{}).Resolve("file://" + path)
	if err != nil || string(value) != "from-file" {
		t.Fatalf("file secret = %q, %v", value, err)
	}
}
