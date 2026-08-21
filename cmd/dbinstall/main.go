package main

import (
	"os"

	"github.com/vatebur/dbinstall/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(cli.ExitCode(err))
	}
}
