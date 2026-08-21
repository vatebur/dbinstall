package builtin

import (
	"fmt"

	"github.com/vatebur/dbinstall/internal/provider"
	"github.com/vatebur/dbinstall/internal/providers/greatsql"
	"github.com/vatebur/dbinstall/internal/providers/mysql"
	"github.com/vatebur/dbinstall/internal/providers/postgresql"
)

func Registry() (*provider.Registry, error) {
	registry := provider.NewRegistry()
	for _, value := range []provider.Provider{mysql.Provider{}, greatsql.Provider{}, postgresql.Provider{}} {
		if err := registry.Register(value); err != nil {
			return nil, fmt.Errorf("register built-in provider: %w", err)
		}
	}
	return registry, nil
}
