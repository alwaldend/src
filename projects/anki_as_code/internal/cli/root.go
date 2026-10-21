package cli

import (
	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"github.com/spf13/cobra"
)

// New assembles commands around their workflow dependency.
func New(application ankiascodeapi.App) *cobra.Command {
	command := &cobra.Command{
		Use:           "anki_as_code",
		Short:         "Export and reconcile Anki collections as TOML",
		SilenceErrors: true,
		SilenceUsage:  true,
		// Parse host flags before entering opaque Rust server arguments.
		TraverseChildren: true,
	}
	command.PersistentFlags().String("server-engine", "", "Engine used by server commands: cgo or wasm (default: cgo, then wasm)")
	command.CompletionOptions.DisableDefaultCmd = true
	command.AddCommand(newConfigCommand(application), newCollectionCommand(application), newServerCommand(application), newAPICommand(application))
	return command
}
