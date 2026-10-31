package cli

import (
	"context"
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

// newRunCmd coordinates configured plugins around a child process.
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
			runCmd := exec.Command(args[0], args[1:]...)
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
			if err := lc.Add(pluginManager.Lifecycle()); err != nil {
				return fmt.Errorf("register plugin lifecycle: %w", err)
			}
			if err := pluginManager.AddLabels(pluginLabels); err != nil {
				return fmt.Errorf("could not add plugin labels: %w", err)
			}
			if err := lc.Start(ctx.Ctx); err != nil {
				return fmt.Errorf("could not start the lifecycle: %w", err)
			}
			pluginEnv, err := pluginManager.Env()
			if err != nil {
				return fmt.Errorf("collect plugin environment: %w", err)
			}
			runCmd.Env = append(runCmd.Env, pluginEnv...)
			if err := runCommand(ctx.Ctx, runCmd); err != nil {
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

// runCommand lets the child finish its own shutdown while its injected files and plugin
// services are still available. The enclosing lifecycle stops plugins afterward.
func runCommand(ctx context.Context, cmd *exec.Cmd) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start cancelled command: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}
	wait := make(chan error, 1)
	go func() { wait <- cmd.Wait() }()
	select {
	case err := <-wait:
		if err != nil {
			return fmt.Errorf("wait for command: %w", err)
		}
		return nil
	case <-ctx.Done():
		stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return fmt.Errorf("stop cancelled command: %w", errors.Join(ctx.Err(), lifecycle.StopProcess(stopCtx, cmd, wait)))
	}
}

// overrideStd creates or truncates output files, opens an existing input file,
// and registers owned descriptors while borrowing ordinary standard streams.
func overrideStd(lc *lifecycle.Manager, cmd *exec.Cmd, stdout string, stderr string, stdin string) error {
	var err error
	var stdoutFile, stderrFile, stdinFile *os.File
	if stdout == "" {
		stdoutFile = os.Stdout
	} else {
		stdoutFile, err = os.OpenFile(stdout, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			return fmt.Errorf("could not open stdout %s: %w", stdout, err)
		}
		if err := lc.Add(lifecycle.StoppableFunc0(stdoutFile.Close)); err != nil {
			return fmt.Errorf("register stdout file: %w", errors.Join(err, stdoutFile.Close()))
		}
	}
	if stderr == "" {
		stderrFile = os.Stderr
	} else {
		stderrFile, err = os.OpenFile(stderr, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			return fmt.Errorf("could not open stderr %s: %w", stderr, err)
		}
		if err := lc.Add(lifecycle.StoppableFunc0(stderrFile.Close)); err != nil {
			return fmt.Errorf("register stderr file: %w", errors.Join(err, stderrFile.Close()))
		}
	}
	if stdin == "" {
		stdinFile = os.Stdin
	} else {
		stdinFile, err = os.OpenFile(stdin, os.O_RDONLY, 0o600)
		if err != nil {
			return fmt.Errorf("could not open stdin %s: %w", stdin, err)
		}
		if err := lc.Add(lifecycle.StoppableFunc0(stdinFile.Close)); err != nil {
			return fmt.Errorf("register stdin file: %w", errors.Join(err, stdinFile.Close()))
		}
	}
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	cmd.Stdin = stdinFile
	return nil
}
