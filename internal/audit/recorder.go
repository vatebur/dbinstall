package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vatebur/dbinstall/internal/plan"
)

func WritePlan(stateDirectory string, document plan.Document) (string, error) {
	directory := filepath.Join(stateDirectory, "operations")
	if err := os.MkdirAll(directory, 0o750); err != nil {
		return "", fmt.Errorf("create operation directory: %w", err)
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return "", fmt.Errorf("serialize plan: %w", err)
	}
	data = append(data, '\n')
	path := filepath.Join(directory, document.OperationID+"-plan.json")
	temporary, err := os.CreateTemp(directory, ".plan-*.tmp")
	if err != nil {
		return "", fmt.Errorf("create temporary plan: %w", err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if err := temporary.Chmod(0o640); err != nil {
		temporary.Close()
		return "", fmt.Errorf("secure temporary plan: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return "", fmt.Errorf("write temporary plan: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return "", fmt.Errorf("sync temporary plan: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return "", fmt.Errorf("close temporary plan: %w", err)
	}
	if err := os.Rename(temporaryName, path); err != nil {
		return "", fmt.Errorf("publish plan: %w", err)
	}
	return path, nil
}
