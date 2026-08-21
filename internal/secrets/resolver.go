package secrets

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrNotFound = errors.New("secret not found")

type Resolver interface {
	Resolve(reference string) ([]byte, error)
}

type LocalResolver struct{}

func (LocalResolver) Resolve(reference string) ([]byte, error) {
	switch {
	case strings.HasPrefix(reference, "env://"):
		name := strings.TrimPrefix(reference, "env://")
		value, ok := os.LookupEnv(name)
		if !ok {
			return nil, fmt.Errorf("%w: environment variable %s", ErrNotFound, name)
		}
		return []byte(value), nil
	case strings.HasPrefix(reference, "file://"):
		path := strings.TrimPrefix(reference, "file://")
		if !filepath.IsAbs(path) {
			return nil, errors.New("secret file path must be absolute")
		}
		value, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%w: read secret file: %v", ErrNotFound, err)
		}
		return []byte(strings.TrimSuffix(string(value), "\n")), nil
	default:
		return nil, errors.New("unsupported secret reference")
	}
}
