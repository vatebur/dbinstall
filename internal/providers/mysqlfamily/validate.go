package mysqlfamily

import (
	"fmt"
	"strings"

	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/spec"
)

func Validate(document spec.DatabaseInstance, machine host.Info, supportedSeries string) error {
	if !strings.HasPrefix(document.Spec.Version, supportedSeries+".") {
		return fmt.Errorf("version %s is outside supported series %s", document.Spec.Version, supportedSeries)
	}
	if document.Spec.Installation.Method == "package" && document.Metadata.Name != document.Spec.Provider {
		return fmt.Errorf("package installation requires the canonical instance name %q", document.Spec.Provider)
	}
	if machine.OSFamily != "rhel" && machine.OSFamily != "debian" {
		return fmt.Errorf("unsupported OS family %q", machine.OSFamily)
	}
	return nil
}
