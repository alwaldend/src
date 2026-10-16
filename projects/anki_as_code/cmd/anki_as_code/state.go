package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/app"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"github.com/spf13/cobra"
)

// newStateCommand groups state operations using application-managed collection access.
func newStateCommand(application *app.App, options *app.StateOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "state",
		Short: "Inspect and manage collection-backed path bindings",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return fmt.Errorf("state requires a subcommand")
		},
	}
	command.PersistentFlags().StringVar(&application.Options.Input, "input", "", "Existing .anki2 collection database")
	if err := command.MarkPersistentFlagRequired("input"); err != nil {
		panic(fmt.Errorf("register required state input: %w", err))
	}
	command.AddCommand(newStateDumpCommand(application, options), newStateImportCommand(application, options),
		newStateWriteCommand(application, options, "replace", "Replace the whole stored state object", (*state.Store).Replace),
		newStateWriteCommand(application, options, "update", "Merge an external state object into stored state", (*state.Store).Update))
	return command
}

// newStateDumpCommand prints stored bindings through application-managed access.
func newStateDumpCommand(application *app.App, options *app.StateOptions) *cobra.Command {
	var format string
	command := &cobra.Command{
		Use: "dump", Short: "Dump the stored state object", Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := application.Start(command.Context()); err != nil {
				return fmt.Errorf("initialize state application: %w", err)
			}
			err := application.WithState(command.Context(), *options, func(store *state.Store) error {
				value, err := store.Dump(command.Context())
				if err != nil {
					return fmt.Errorf("dump state: %w", err)
				}
				data, err := codec.Encode(value, codec.Format(format))
				if err != nil {
					return fmt.Errorf("encode dumped state: %w", err)
				}
				if _, err := command.OutOrStdout().Write(data); err != nil {
					return fmt.Errorf("write dumped state: %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("state dump: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&format, "format", string(codec.TOML), "State format: toml or json")
	return command
}

// newStateImportCommand passes a complete PATH=ID batch to one managed state operation.
func newStateImportCommand(application *app.App, options *app.StateOptions) *cobra.Command {
	var text string
	command := &cobra.Command{
		Use: "import PATH=ID [PATH=ID...]", Short: "Bind external objects to configuration paths", Args: cobra.MinimumNArgs(1),
		RunE: func(command *cobra.Command, arguments []string) error {
			bindings := map[string]string{}
			for _, argument := range arguments {
				separator := strings.LastIndexByte(argument, '=')
				if separator <= 0 {
					return fmt.Errorf("binding %q must be PATH=ID", argument)
				}
				address := argument[:separator]
				if _, exists := bindings[address]; exists {
					return fmt.Errorf("binding path %q was supplied more than once", address)
				}
				bindings[address] = argument[separator+1:]
			}
			configs, err := config.Load(text)
			if err != nil {
				return fmt.Errorf("load import configs: %w", err)
			}
			if err := application.Start(command.Context()); err != nil {
				return fmt.Errorf("initialize state application: %w", err)
			}
			err = application.WithState(command.Context(), *options, func(store *state.Store) error {
				if _, err := store.Import(command.Context(), configs, bindings); err != nil {
					return fmt.Errorf("import collection bindings: %w", err)
				}
				if _, err := fmt.Fprintln(command.ErrOrStderr(), "Imported collection bindings:", len(bindings)); err != nil {
					return fmt.Errorf("write import status: %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("state import: %w", err)
			}
			return nil
		},
	}
	textFlag.bindRequired(command, &text)
	return command
}

// newStateWriteCommand decodes an external object and invokes its selected state operation.
func newStateWriteCommand(application *app.App, options *app.StateOptions, name, summary string, write func(*state.Store, context.Context, *collectionpb.State) (bool, error)) *cobra.Command {
	var format string
	command := &cobra.Command{
		Use: name + " [FILE]", Short: summary, Args: cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, arguments []string) error {
			var data []byte
			var err error
			if len(arguments) == 0 || arguments[0] == "-" {
				data, err = io.ReadAll(command.InOrStdin())
			} else {
				data, err = os.ReadFile(arguments[0])
			}
			if err != nil {
				return fmt.Errorf("read state for %s: %w", name, err)
			}
			value := &collectionpb.State{}
			if err := codec.Decode(data, codec.Format(format), value); err != nil {
				return fmt.Errorf("decode state for %s: %w", name, err)
			}
			if err := application.Start(command.Context()); err != nil {
				return fmt.Errorf("initialize state application: %w", err)
			}
			err = application.WithState(command.Context(), *options, func(store *state.Store) error {
				if _, err := write(store, command.Context(), value); err != nil {
					return fmt.Errorf("%s state: %w", name, err)
				}
				if _, err := fmt.Fprintf(command.ErrOrStderr(), "Collection state %s completed\n", name); err != nil {
					return fmt.Errorf("write %s status: %w", name, err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("state %s: %w", name, err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&format, "format", string(codec.TOML), "State format: toml or json")
	return command
}
