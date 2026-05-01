package main

import (
	"context"
	"fmt"
	"io"

	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
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
	var secretEnv []string
	var vaultEnv []string
	var disableRunfilesEnv bool
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a command",
		Long:  "Run a command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Need at least one argument")
			}
			cfg, err := al.LoadConfigs(ctx, configs...)
			if err != nil {
				return fmt.Errorf("could not load configs: %w", err)
			}
			vault, err := al.NewVault(cfg)
			if err != nil {
				return fmt.Errorf("could not create Vault: %w", err)
			}
            defer vault.Clean()
			runCmd, err := al.Command(
				al.CommandArgs{
					Ctx:                ctx,
					Name:               args[0],
					Args:               args[1:],
					SecretEnv:          secretEnv,
					Vault:              vault,
					VaultEnv:        vaultEnv,
					DisableRunfilesEnv: disableRunfilesEnv,
				},
			)
			if err != nil {
				return fmt.Errorf("could not create command: %w", err)
			}
			err = runCmd.Run()
			return err
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&configs, "config", nil, "Config path")
	flags.StringArrayVar(&secretEnv, "secret_env", nil, "Env variables to add (supports templating)")
	flags.StringArrayVar(&vaultEnv, "vault_env", nil, "If set, set Vault environment variables from the config")
	flags.BoolVar(&disableRunfilesEnv, "disable_runfiles_env", false, "If set, do not set bazel runfiles variables")
	return cmd
}
