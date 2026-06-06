package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"github.com/spf13/cobra"
)

func Execute(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	root, err := newRootCommand(ctx, args, stdin, stdout, stderr)
	if err != nil {
		return fmt.Errorf("could not create commands: %w", err)
	}
	err = root.Execute()
	if err != nil {
		return fmt.Errorf("could not execute commands: %w", err)
	}
	return nil
}

func newRootCommand(
	ctx context.Context,
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "al",
		Short:         "Al repository tool",
		Long:          "Al repository tool",
		SilenceErrors: true,
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)
	cmd.AddCommand(newConfigCmd(ctx))
	cmd.AddCommand(newRunCmd(ctx))
	return cmd, nil
}

func newConfigCmd(ctx context.Context) *cobra.Command {
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
			err := al.DumpConfigs(ctx, out, configs...)
			return err
		},
	}
	flags := cmdDump.PersistentFlags()
	flags.StringArrayVarP(&configs, "config", "c", nil, "Config path")
	flags.StringVarP(&out, "out", "o", "", "Output path")
	cmd.AddCommand(cmdDump)
	return cmd
}

func newRunCmd(ctx context.Context) *cobra.Command {
	var configs []string
	var disableRunfilesEnv bool
	var stdout string
	var stdin string
	var stderr string
	var pluginLabels []string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a command",
		Long:  "Run a command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Need at least one argument")
			}
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			cfg, err := al.LoadConfigs(ctx, configs...)
			if err != nil {
				return fmt.Errorf("could not load configs: %w", err)
			}
			cmdArgs := al.CommandArgs{
				Ctx:                ctx,
				Name:               args[0],
				Args:               args[1:],
				DisableRunfilesEnv: disableRunfilesEnv,
			}
			if stdout != "" {
				stdoutFile, err := os.OpenFile(stdout, os.O_WRONLY, 0600)
				if err != nil {
					return fmt.Errorf("could not open stdout %s: %w", stdout, err)
				}
				defer stdoutFile.Close()
				cmdArgs.Stdout = stdoutFile
			}
			if stderr != "" {
				stderrFile, err := os.OpenFile(stderr, os.O_WRONLY, 0600)
				if err != nil {
					return fmt.Errorf("could not open stderr %s: %w", stderr, err)
				}
				defer stderrFile.Close()
				cmdArgs.Stderr = stderrFile
			}
			if stdin != "" {
				stdinFile, err := os.OpenFile(stdin, os.O_RDONLY, 0600)
				if err != nil {
					return fmt.Errorf("could not open stdin %s: %w", stdin, err)
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
			pluginStartResponses, err := pluginManager.StartPlugins(ctx, pluginLabels).Get()
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
			fmt.Fprintf(os.Stderr, "Setting envs %s\n", strings.Join(envs, ", "))
			err = runCmd.Run()
			return err
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&configs, "config", nil, "Config path")
	flags.StringVar(&stdout, "stdout", "", "Override stdout")
	flags.StringVar(&stderr, "stderr", "", "Override stderr")
	flags.StringVar(&stdin, "stdin", "", "Override stdin")
	flags.StringArrayVar(&pluginLabels, "plugin_label", nil, "Plugin labels to run")
	flags.BoolVar(&disableRunfilesEnv, "disable_runfiles_env", false, "If set, do not set bazel runfiles variables")
	return cmd
}
