package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type stringFlag struct {
	name  string
	usage string
}

var (
	inputFlag = stringFlag{
		name:  "input",
		usage: "Collection archive or offline SQLite database",
	}
	textFlag = stringFlag{
		name:  "text",
		usage: "Collection root config or its containing directory",
	}
	outputFlag = stringFlag{
		name:  "output",
		usage: "Destination text directory (export) or collection archive (build)",
	}
	baseFlag = stringFlag{
		name:  "base",
		usage: "Original collection archive",
	}
	collectionFlag = stringFlag{
		name:  "collection",
		usage: "Root collection config for identity reservation, including external note paths",
	}
)

func (flag stringFlag) bind(command *cobra.Command, value *string) {
	command.Flags().StringVar(value, flag.name, "", flag.usage)
}

func (flag stringFlag) bindRequired(command *cobra.Command, value *string) {
	flag.bind(command, value)
	// The definition registers this fixed name immediately before requiring it.
	if err := command.MarkFlagRequired(flag.name); err != nil {
		panic(fmt.Errorf("register required flag %q: %w", flag.name, err))
	}
}
