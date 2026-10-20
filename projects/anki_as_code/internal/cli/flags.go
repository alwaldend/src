package cli

import (
	"fmt"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"

	"github.com/spf13/cobra"
)

// stringFlag defines a reusable string flag and its help text.
type stringFlag struct {
	// name is the long option name registered with Cobra.
	name string
	// usage is the help description shown for the option.
	usage string
}

var (
	// addressFlag selects an owned native strategy or an external server.
	addressFlag = stringFlag{
		name:  "address",
		usage: "Anki backend (subprocess[+ENGINE]://NAME, inprocess[+ENGINE]://NAME, host:port, tcp:// or unix://); empty uses subprocess",
	}
	// stateKeyFlag selects the namespace used by collection operations.
	stateKeyFlag = stringFlag{
		name:  "state-key",
		usage: "Collection state namespace key",
	}
	// collectionFlag describes the source collection option shared by collection commands.
	collectionFlag = stringFlag{
		name:  "collection",
		usage: "Anki .colpkg archive or offline .anki2 database",
	}
	// textFlag describes the declaration input shared by configuration and collection commands.
	textFlag = stringFlag{
		name:  "text",
		usage: "Collection root config or its containing directory",
	}
	// outputFlag describes the destination option used by export and build.
	outputFlag = stringFlag{
		name:  "output",
		usage: "Destination text directory (export) or collection archive (build)",
	}
	// baseFlag describes the base archive option used by build.
	baseFlag = stringFlag{
		name:  "base",
		usage: "Original collection archive",
	}
)

// bind registers the string flag and binds it to the destination value.
func (flag stringFlag) bind(command *cobra.Command, value *string) {
	command.Flags().StringVar(value, flag.name, "", flag.usage)
}

// bindRequired registers the string flag and requires its presence for command execution.
func (flag stringFlag) bindRequired(command *cobra.Command, value *string) {
	flag.bind(command, value)
	// The definition registers this fixed name immediately before requiring it.
	if err := command.MarkFlagRequired(flag.name); err != nil {
		panic(fmt.Errorf("register required flag %q: %w", flag.name, err))
	}
}

// commandCollectionOptions combines the workflow path with parsed transport and namespace flags.
func commandCollectionOptions(command *cobra.Command, path string) ankiascodeapi.CollectionOptions {
	return ankiascodeapi.CollectionOptions{
		Path:     path,
		Address:  command.Flag(addressFlag.name).Value.String(),
		StateKey: command.Flag(stateKeyFlag.name).Value.String(),
	}
}
