package cli

import (
	"fmt"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/pkg/al"
	"github.com/spf13/cobra"
)

// newConfigCmd exposes configuration inspection.
func newConfigCmd(ctx *al.CmdCtx) *cobra.Command {
	var configs []string
	var out string
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config commands",
		Long:  "Config commands",
	}
	cmdDump := &cobra.Command{
		Use:   "dump",
		Short: "Dump configs",
		Long:  "Dump merged configs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := al.DumpConfigs(ctx.Ctx, out, configs...); err != nil {
				return fmt.Errorf("dump configuration: %w", err)
			}
			return nil
		},
	}
	flags := cmdDump.PersistentFlags()
	flags.StringArrayVarP(&configs, "config", "c", nil, "Config path")
	flags.StringVarP(&out, "out", "o", "", "Output path")
	cmd.AddCommand(cmdDump)
	return cmd
}
