package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

// fdroidserverVersion pins the tool installed in the workspace virtual environment.
const fdroidserverVersion = "2.4.5"

// main reports wrapper or F-Droid failures to the caller.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	command := &cobra.Command{
		Use: "fdroid-wrapper", DisableFlagParsing: true, SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(cmd.Context(), args); err != nil {
				return fmt.Errorf("run F-Droid: %w", err)
			}
			return nil
		},
	}
	err := command.ExecuteContext(ctx)
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run locates the workspace cache and launches F-Droid in the caller’s data directory.
func run(ctx context.Context, args []string) error {
	workspaceRoot := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if workspaceRoot == "" {
		binaryPath, err := filepath.Abs(os.Args[0])
		if err != nil {
			return fmt.Errorf("resolve binary path: %w", err)
		}
		workspaceRoot = filepath.Dir(filepath.Dir(filepath.Dir(binaryPath)))
	}
	cacheDir := filepath.Join(workspaceRoot, "out", "fdroid-wrapper")
	venvDir := filepath.Join(cacheDir, "venv")
	fdroidBin := filepath.Join(venvDir, "bin", "fdroid")

	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("create cache directory: %w", err)
	}
	if info, err := os.Stat(fdroidBin); err != nil || info.Mode()&0o111 == 0 {
		fmt.Fprintln(os.Stderr, "provisioning fdroidserver "+fdroidserverVersion+"...")
		if err := provision(ctx, venvDir); err != nil {
			return fmt.Errorf("provision F-Droid environment %q: %w", venvDir, err)
		}
	}

	env := os.Environ()
	if os.Getenv("ANDROID_HOME") == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			sdk := filepath.Join(home, "Android", "Sdk")
			if info, err := os.Stat(sdk); err == nil && info.IsDir() {
				env = append(env, "ANDROID_HOME="+sdk)
			}
		}
	}

	jdk := filepath.Join(workspaceRoot, "out", "gradle-wrapper", "jdk21")
	if os.Getenv("JAVA_HOME") == "" {
		if info, err := os.Stat(filepath.Join(jdk, "bin", "java")); err == nil && info.Mode()&0o111 != 0 {
			env = append(env, "JAVA_HOME="+jdk, "PATH="+filepath.Join(jdk, "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
		}
	}

	cmd := exec.CommandContext(ctx, fdroidBin, args...)
	cmd.Env = env
	cmd.Dir = os.Getenv("BUILD_WORKING_DIRECTORY")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execute %q: %w", cmd.Path, err)
	}
	return nil
}

// provision installs the pinned F-Droid server into a Python virtual environment.
func provision(ctx context.Context, venvDir string) error {
	if err := runCommand(ctx, "python3", "-m", "venv", venvDir); err != nil {
		return fmt.Errorf("create virtualenv: %w", err)
	}
	pip := filepath.Join(venvDir, "bin", "pip")
	if err := runCommand(ctx, pip, "install", "--disable-pip-version-check", "fdroidserver=="+fdroidserverVersion); err != nil {
		return fmt.Errorf("install fdroidserver: %w", err)
	}
	return nil
}

// runCommand runs a provisioning step with its diagnostics attached.
func runCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execute %q: %w", cmd.Path, err)
	}
	return nil
}
