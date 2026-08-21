package host

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadOSReleaseAndFamily(t *testing.T) {
	path := filepath.Join(t.TempDir(), "os-release")
	content := "ID=custom\nID_LIKE=\"debian ubuntu\"\nVERSION_ID='24.04'\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	values, err := readOSRelease(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := osFamily(values); got != "debian" {
		t.Fatalf("family = %q", got)
	}
}

func TestCertification(t *testing.T) {
	certified, _ := certification(Info{OSID: "ubuntu", OSVersion: "24.04", InitSystem: "systemd"})
	if !certified {
		t.Fatal("Ubuntu 24.04 with systemd should match the matrix")
	}
	certified, _ = certification(Info{OSID: "centos", OSVersion: "7", InitSystem: "systemd"})
	if certified {
		t.Fatal("CentOS 7 should not match the matrix")
	}
}
