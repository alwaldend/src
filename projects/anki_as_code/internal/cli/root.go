package cli

import "github.com/spf13/cobra"

// New assembles commands around their workflow dependency.
func New(application application) *cobra.Command {
	command := &cobra.Command{
		Use:           "anki_as_code",
		Short:         "Export and reconcile Anki collections as TOML",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	command.CompletionOptions.DisableDefaultCmd = true
	command.AddCommand(newConfigCommand(application), newCollectionCommand(application), newServerCommand(application))
	return command
}
