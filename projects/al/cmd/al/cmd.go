package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"github.com/spf13/cobra"
)

func Execute(ctx *al.CmdCtx) int {
	root, err := newRootCommand(ctx)
	if err != nil {
		fmt.Fprintf(ctx.Stderr, "could not create commands: %s\n", err)
		return 2
	}
	if err := root.Execute(); err != nil {
		fmt.Fprintf(ctx.Stderr, "could not execute: %s\n", err)
		return 1
	}
	return 0
}

func newRootCommand(ctx *al.CmdCtx) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "al",
		Short:         "Al repository tool",
		Long:          "Al repository tool",
		SilenceErrors: true,
	}
	cmd.SetOut(ctx.Stdout)
	cmd.SetErr(ctx.Stderr)
	cmd.SetIn(ctx.Stdin)
	cmd.SetArgs(ctx.Args[1:])
	cmd.AddCommand(newConfigCmd(ctx))
	cmd.AddCommand(newRunCmd(ctx))
	return cmd, nil
}

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
			err := al.DumpConfigs(ctx.Ctx, out, configs...)
			return err
		},
	}
	flags := cmdDump.PersistentFlags()
	flags.StringArrayVarP(&configs, "config", "c", nil, "Config path")
	flags.StringVarP(&out, "out", "o", "", "Output path")
	cmd.AddCommand(cmdDump)
	return cmd
}

func newRunCmd(ctx *al.CmdCtx) *cobra.Command {
	var configs []string
	var disableRunfilesEnv bool
	var cmdStdout string
	var cmdStdin string
	var cmdStderr string
	var pluginLabels []string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a command",
		Long:  "Run a command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Need at least one argument")
			}
			cmdCtx, cancel := context.WithCancel(ctx.Ctx)
			defer cancel()
			cfg, err := al.LoadConfigs(cmdCtx, configs...)
			if err != nil {
				return fmt.Errorf("could not load configs: %w", err)
			}
			cmdArgs := al.CommandArgs{
				Ctx:                cmdCtx,
				Name:               args[0],
				Args:               args[1:],
				DisableRunfilesEnv: disableRunfilesEnv,
			}
			if cmdStdout != "" {
				stdoutFile, err := os.OpenFile(cmdStdout, os.O_WRONLY, 0600)
				if err != nil {
					return fmt.Errorf("could not open stdout %s: %w", cmdStdout, err)
				}
				defer stdoutFile.Close()
				cmdArgs.Stdout = stdoutFile
			}
			if cmdStderr != "" {
				stderrFile, err := os.OpenFile(cmdStderr, os.O_WRONLY, 0600)
				if err != nil {
					return fmt.Errorf("could not open stderr %s: %w", cmdStderr, err)
				}
				defer stderrFile.Close()
				cmdArgs.Stderr = stderrFile
			}
			if cmdStdin != "" {
				stdinFile, err := os.OpenFile(cmdStdin, os.O_RDONLY, 0600)
				if err != nil {
					return fmt.Errorf("could not open stdin %s: %w", cmdStdin, err)
				}
				defer stdinFile.Close()
				cmdArgs.Stdin = stdinFile
			}
			runCmd, err := al.Command(cmdArgs)
			if err != nil {
				return fmt.Errorf("could not create command: %w", err)
			}
			pluginManager, err := al_plugin.NewManager(cfg, cmdArgs.Stderr)
			if err != nil {
				return fmt.Errorf("could not create plugin manager: %w", err)
			}
			pluginStartResponses, err := pluginManager.StartPlugins(cmdCtx, pluginLabels)
			if err != nil {
				return fmt.Errorf("could not prepare plugins: %w", err)
			}
			envs := []string{}
			for _, response := range pluginStartResponses {
				for key, value := range response.Env {
					envs = append(envs, key)
					runCmd.Env = append(runCmd.Env, fmt.Sprintf("%s=%s", key, value))
				}
			}
			if len(envs) > 0 {
				fmt.Fprintf(ctx.Stderr, "Setting envs: %s\n", strings.Join(envs, ", "))
			}
			var res error
			if err := runCmd.Run(); err != nil {
				res = fmt.Errorf("could not run the command: %w", err)
			}
			cancel()
			if err := pluginManager.Shutdown(); err != nil {
				res = errors.Join(res, fmt.Errorf("could not shutdown plugins: %w", err))
			}
			return res
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&configs, "config", nil, "Config path")
	flags.StringVar(&cmdStdout, "stdout", "", "Override stdout")
	flags.StringVar(&cmdStderr, "stderr", "", "Override stderr")
	flags.StringVar(&cmdStdin, "stdin", "", "Override stdin")
	flags.StringArrayVar(&pluginLabels, "plugin_label", nil, "Plugin labels to run")
	flags.BoolVar(&disableRunfilesEnv, "disable_runfiles_env", false, "If set, do not set bazel runfiles variables")
	return cmd
}
