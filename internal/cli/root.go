package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/vatebur/dbinstall/internal/build"
	"github.com/vatebur/dbinstall/internal/host"
	"github.com/vatebur/dbinstall/internal/providers/builtin"
	"github.com/vatebur/dbinstall/internal/spec"
)

const (
	ExitUsage      = 2
	ExitValidation = 3
	ExitPreflight  = 4
	ExitInternal   = 10
)

type codedError struct {
	code int
	err  error
}

func (e *codedError) Error() string { return e.err.Error() }
func (e *codedError) Unwrap() error { return e.err }

func ExitCode(err error) int {
	var coded *codedError
	if errors.As(err, &coded) {
		return coded.code
	}
	return ExitInternal
}

func Execute() error {
	root := NewRootCommand(os.Stdout, os.Stderr)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

func NewRootCommand(stdout, stderr io.Writer) *cobra.Command {
	var output string
	root := &cobra.Command{
		Use:           "dbinstall",
		Short:         "Plan and manage native database installations",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.PersistentFlags().StringVar(&output, "output", "human", "output format: human or json")

	root.AddCommand(newVersionCommand(stdout, &output))
	root.AddCommand(newProvidersCommand(stdout, &output))
	root.AddCommand(newInspectCommand(stdout, &output))
	root.AddCommand(newValidateCommand(stdout, &output))
	return root
}

func newVersionCommand(stdout io.Writer, output *string) *cobra.Command {
	return &cobra.Command{
		Use: "version", Short: "Print build version",
		RunE: func(_ *cobra.Command, _ []string) error {
			value := map[string]string{"version": build.Version, "commit": build.Commit, "date": build.Date}
			if *output == "json" {
				return writeJSON(stdout, value)
			}
			_, err := fmt.Fprintf(stdout, "dbinstall %s (%s, %s)\n", build.Version, build.Commit, build.Date)
			return err
		},
	}
}

func newProvidersCommand(stdout io.Writer, output *string) *cobra.Command {
	return &cobra.Command{
		Use: "providers", Short: "List compiled database providers",
		RunE: func(_ *cobra.Command, _ []string) error {
			registry, err := builtin.Registry()
			if err != nil {
				return &codedError{ExitInternal, err}
			}
			descriptors := registry.Descriptors()
			if *output == "json" {
				return writeJSON(stdout, descriptors)
			}
			for _, descriptor := range descriptors {
				fmt.Fprintf(stdout, "%s\t%s\t%s\n", descriptor.Name, descriptor.FoundationStatus, descriptor.DisplayName)
			}
			return nil
		},
	}
}

func newInspectCommand(stdout io.Writer, output *string) *cobra.Command {
	return &cobra.Command{
		Use: "inspect", Short: "Inspect the local host without changing it",
		RunE: func(_ *cobra.Command, _ []string) error {
			machine, err := (host.LocalDetector{}).Detect("/")
			if err != nil {
				return &codedError{ExitPreflight, err}
			}
			if *output == "json" {
				return writeJSON(stdout, machine)
			}
			fmt.Fprintf(stdout, "OS: %s %s (%s)\nArchitecture: %s\nCPU: %d\nMemory: %d bytes\nInit: %s\nCertified: %t - %s\n",
				machine.OSID, machine.OSVersion, machine.OSFamily, machine.Architecture,
				machine.CPUCount, machine.MemoryBytes, machine.InitSystem, machine.Certified, machine.Certification)
			return nil
		},
	}
}

func newValidateCommand(stdout io.Writer, output *string) *cobra.Command {
	var file string
	command := &cobra.Command{
		Use: "validate", Short: "Validate and normalize an installation specification",
		RunE: func(_ *cobra.Command, _ []string) error {
			document, err := spec.LoadFile(file)
			if err != nil {
				return &codedError{ExitValidation, err}
			}
			if *output == "json" {
				return writeJSON(stdout, document)
			}
			fmt.Fprintf(stdout, "valid: %s (%s %s)\n", document.Metadata.Name, document.Spec.Provider, document.Spec.Version)
			return nil
		},
	}
	command.Flags().StringVarP(&file, "file", "f", "", "path to the YAML specification")
	_ = command.MarkFlagRequired("file")
	return command
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
