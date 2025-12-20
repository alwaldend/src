package main

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func Execute(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
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
		Use:           "secrets",
		Short:         "secrets",
		Long:          "Secrets tool",
		SilenceErrors: true,
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)
	cmdExport, err := newExportCommand(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create exprot command: %w", err)
	}
	cmdBuild, err := newBuildCommand(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create build command: %w", err)
	}
	cmdEdit, err := newEditCommand(ctx, stdout, stderr)
	if err != nil {
		return nil, fmt.Errorf("could not create edit command: %w", err)
	}
	cmd.AddCommand(cmdExport, cmdBuild, cmdEdit)
	return cmd, nil
}

func newBuildCommand(ctx context.Context) (*cobra.Command, error) {
	builder := NewBuilder()
	build := &cobra.Command{
		Use:   "build",
		Short: "build",
		Long:  "Build",
	}

	secretOpts := &BuilderSecretOpts{}
	secret := &cobra.Command{
		Use:   "secret",
		Short: "secret",
		Long:  "Secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			return builder.BuildSecret(secretOpts)
		},
	}
	secretFlags := secret.PersistentFlags()
	secretFlags.StringArrayVar(&secretOpts.Secret, "secret", nil, "Embed an additional secret")
	secretFlags.StringArrayVar(&secretOpts.Backend, "backend", nil, "Backend information")
	secretFlags.StringArrayVar(&secretOpts.Env, "env", nil, "Env variables in the form of key=value")
	secretFlags.StringArrayVar(&secretOpts.Output, "output", nil, "Output secret to this file")
	build.AddCommand(secret)

	return build, nil
}

func newExportCommand(ctx context.Context) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "export",
		Long:  "Export a secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented")
		},
	}
	return cmd, nil
}

func newEditCommand(
	ctx context.Context,
	stdout, stderr io.Writer,
) (*cobra.Command, error) {
	editor := NewEditor(stdout, stderr)
	opts := &EditorOpts{}
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "edit",
		Long:  "Edit",
		RunE: func(cmd *cobra.Command, args []string) error {
			return editor.Edit(opts)
		},
	}
	secretFlags := cmd.PersistentFlags()
	secretFlags.StringArrayVar(&opts.Backend, "backend", nil, "Edit a backend")
	secretFlags.StringArrayVar(&opts.Output, "output", nil, "Output the result")
	return cmd, nil
}
