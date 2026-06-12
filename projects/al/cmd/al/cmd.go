package main

import (
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
		ctx.Logger.Printf("could not create commands: %s", err)
		return 2
	}
	if err := root.Execute(); err != nil {
		ctx.Logger.Printf("could not execute: %s", err)
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
	var cmdStdout string
	var cmdStdin string
	var cmdStderr string
	var pluginLabels []string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a command",
		Long:  "Run a command",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 {
				return fmt.Errorf("Need at least one argument")
			}
			var lc lifecycle.Manager
			defer func() {
				if curErr := lc.StopTimeout(time.Second * 10); curErr != nil {
					err = errors.Join(err, fmt.Errorf("shut down with an error: %w", curErr))
				}
			}()
			cfg, err := al.LoadConfigs(ctx.Ctx, configs...)
			if err != nil {
				return fmt.Errorf("could not load configs: %w", err)
			}
			runCmd := exec.CommandContext(ctx.Ctx, args[0], args[1:]...)
			runCmd.Env = os.Environ()
			if err := al.SetRunfilesEnv(runCmd); err != nil {
				return fmt.Errorf("could not add runfiles env: %w", err)
			}
			if err := overrideStd(&lc, runCmd, cmdStdout, cmdStderr, cmdStdin); err != nil {
				return fmt.Errorf("could not override std: %w", err)
			}
			pluginManager, err := al_plugin.NewManager(ctx, cfg)
			if err != nil {
				return fmt.Errorf("could not create plugin manager: %w", err)
			}
			if err := pluginManager.AddLabels(pluginLabels); err != nil {
				return fmt.Errorf("could not add plugin labels: %w", err)
			}
			lc.Add(pluginManager.Lifecycle())
			if err := lc.Start(ctx.Ctx); err != nil {
				return fmt.Errorf("could not start the lifecycle: %w", err)
			}
			runCmd.Env = append(runCmd.Env, pluginManager.Env()...)
			if err := runCmd.Run(); err != nil {
				return fmt.Errorf("could not run the command: %w", err)
			}
			return nil
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&configs, "config", nil, "Config path")
	flags.StringVar(&cmdStdout, "stdout", "", "Override stdout")
	flags.StringVar(&cmdStderr, "stderr", "", "Override stderr")
	flags.StringVar(&cmdStdin, "stdin", "", "Override stdin")
	flags.StringArrayVar(&pluginLabels, "plugin_label", nil, "Plugin labels to run")
	return cmd
}

func overrideStd(lc *lifecycle.Manager, cmd *exec.Cmd, stdout string, stderr string, stdin string) error {
	var err error
	var stdoutFile, stderrFile, stdinFile *os.File
	if stdout == "" {
		stdoutFile = os.Stdout
	} else {
		stdoutFile, err = os.OpenFile(stdout, os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("could not open stdout %s: %w", stdout, err)
		}
		lc.Add(lifecycle.StoppableFunc0(stdoutFile.Close))
	}
	if stderr == "" {
		stderrFile = os.Stderr
	} else {
		stderrFile, err = os.OpenFile(stderr, os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("could not open stderr %s: %w", stderr, err)
		}
		lc.Add(lifecycle.StoppableFunc0(stderrFile.Close))
	}
	if stdin == "" {
		stdinFile = os.Stdin
	} else {
		stdinFile, err = os.OpenFile(stdin, os.O_RDONLY, 0600)
		if err != nil {
			return fmt.Errorf("could not open stdin %s: %w", stdin, err)
		}
		lc.Add(lifecycle.StoppableFunc0(stdinFile.Close))
	}
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	cmd.Stdin = stdinFile
	return nil
}
