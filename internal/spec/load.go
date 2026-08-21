package spec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadFile(path string) (DatabaseInstance, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DatabaseInstance{}, fmt.Errorf("read specification: %w", err)
	}
	return Load(data)
}

func Load(data []byte) (DatabaseInstance, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var document DatabaseInstance
	if err := decoder.Decode(&document); err != nil {
		return DatabaseInstance{}, fmt.Errorf("decode specification: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return DatabaseInstance{}, errors.New("decode specification: multiple YAML documents are not allowed")
		}
		return DatabaseInstance{}, fmt.Errorf("decode specification trailer: %w", err)
	}

	Default(&document)
	if err := Validate(document); err != nil {
		return DatabaseInstance{}, err
	}
	return document, nil
}

func Default(document *DatabaseInstance) {
	name := document.Metadata.Name
	if document.Spec.Installation.Method == "" {
		document.Spec.Installation.Method = "archive"
	}
	if document.Spec.Installation.Source == "" {
		document.Spec.Installation.Source = "online"
	}
	if document.Spec.Instance.BindAddress == "" {
		document.Spec.Instance.BindAddress = "0.0.0.0"
	}
	if document.Spec.Instance.Port == 0 {
		document.Spec.Instance.Port = defaultPort(document.Spec.Provider)
	}
	root := filepath.Join("/var/lib/dbinstall/instances", name)
	if document.Spec.Instance.Paths.Data == "" {
		document.Spec.Instance.Paths.Data = filepath.Join(root, "data")
	}
	if document.Spec.Instance.Paths.Logs == "" {
		document.Spec.Instance.Paths.Logs = filepath.Join("/var/log/dbinstall", name)
	}
	if document.Spec.Instance.Paths.Temp == "" {
		document.Spec.Instance.Paths.Temp = filepath.Join(root, "tmp")
	}
	if document.Spec.Workload.Profile == "" {
		document.Spec.Workload.Profile = "development"
	}
	if document.Spec.Workload.Durability == "" {
		document.Spec.Workload.Durability = "strict"
	}
}

func defaultPort(provider string) int {
	if provider == "postgresql" {
		return 5432
	}
	return 3306
}
