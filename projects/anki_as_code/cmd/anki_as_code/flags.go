package main

import (
	"fmt"

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
	// inputFlag describes the source collection option shared by collection commands.
	inputFlag = stringFlag{
		name:  "input",
		usage: "Collection archive or offline SQLite database",
	}
	// textFlag describes the desired declaration input used by reconciliation commands.
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
	// collectionFlag describes the root config option used for identity reservation.
	collectionFlag = stringFlag{
		name:  "collection",
		usage: "Root collection config for identity reservation, including external note paths",
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
