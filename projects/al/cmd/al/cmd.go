package main

import (
	"context"
    "log"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
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
    logger := log.New(ctx.Stderr, "com.alwaldend.projects.al.cmd.al.run ", ctx.LogFlags)
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a command",
		Long:  "Run a command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Need at least one argument")
			}
			lc := lifecycle.Handler{}
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
                if err := lc.Stop(ctx); err != nil {
                    logger.Printf("shutdown failed: %w", err)
                }
			}()
			cmdCtx, cancel := context.WithCancel(ctx.Ctx)
			defer cancel()
			cfg, err := al.LoadConfigs(cmdCtx, configs...)
			if err != nil {
				return fmt.Errorf("could not load configs: %w", err)
			}
			runCmd := exec.CommandContext(cmdCtx, args[0], args[1:]...)
			if err := al.SetRunfilesEnv(runCmd); err != nil {
				return fmt.Errorf("could not add runfiles env: %w", err)
			}
			pluginManager, err := al_plugin.NewManager(ctx, cfg)
			if err != nil {
				return fmt.Errorf("could not create plugin manager: %w", err)
			}
			pluginStart, err := pluginManager.Start(cmdCtx, pluginLabels)
			if err != nil {
				return fmt.Errorf("could not prepare plugins: %w", err)
			}
			envs := []string{}
			for _, response := range pluginStart {
				for key, value := range response.Env {
					envs = append(envs, key)
					runCmd.Env = append(runCmd.Env, fmt.Sprintf("%s=%s", key, value))
				}
			}
			if errShutdown := pluginManager.Shutdown(shutdownCtx); errShutdown != nil {
				err = errors.Join(err, fmt.Errorf("could not shut down plugins: %w", errShutdown))
			}
			return err
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
