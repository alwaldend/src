package cli

import (
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/internal/bazelcmd"
	"git.alwaldend.com/alwaldend/src/tools/rules_al/internal/delivery"
	"git.alwaldend.com/alwaldend/src/tools/rules_al/internal/release"
	"git.alwaldend.com/alwaldend/src/tools/rules_al/internal/versioncmd"
	"git.alwaldend.com/alwaldend/src/tools/rules_al/pkg/al"
	"github.com/spf13/cobra"
)

// Execute constructs AL commands and returns the selected operation's exit status.
func Execute(ctx *al.CmdCtx) int {
	root, err := newRootCommand(ctx)
	if err != nil {
		ctx.Logger.Printf("could not create commands: %s", err)
		return 2
	}
	if err := root.Execute(); err != nil {
		ctx.Logger.Printf("could not execute: %s", err)
		var coded interface{ ExitCode() int }
		if errors.As(err, &coded) && coded.ExitCode() > 0 {
			return coded.ExitCode()
		}
		return 1
	}
	return 0
}

// newRootCommand assembles direct lifecycle commands without starting plugins.
func newRootCommand(ctx *al.CmdCtx) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "al",
		Short:         "Repository lifecycle, commands, and Bazel rules",
		Long:          "Repository lifecycle, commands, and Bazel rules",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetOut(ctx.Stdout)
	cmd.SetErr(ctx.Stderr)
	cmd.SetIn(ctx.Stdin)
	cmd.SetArgs(ctx.Args[1:])
	cmd.AddCommand(newConfigCmd(ctx))
	cmd.AddCommand(newRunCmd(ctx))
	releaseCommand, err := release.NewCommand(ctx.Ctx, ctx.Stdin, ctx.Stdout, ctx.Stderr)
	if err != nil {
		return nil, fmt.Errorf("create release commands: %w", err)
	}
	cmd.AddCommand(releaseCommand, delivery.NewCommand(ctx.Ctx, ctx.Getenv, ctx.Stdout), versioncmd.NewCommand(), bazelcmd.NewCommand(), bazelcmd.NewDoctorCommand())
	cmd.SetContext(ctx.Ctx)
	return cmd, nil
}
