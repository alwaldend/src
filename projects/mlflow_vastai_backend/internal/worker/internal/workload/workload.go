package workload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"github.com/spf13/afero"
)

// Outcome records the facts observed while running the workload.
type Outcome struct {
	// ExitCode is set when the command reports an exit code.
	ExitCode *int32
	// Started reports that the workload command launched.
	Started bool
	// Observed reports that the command exit was confirmed.
	Observed bool
	// Stopped reports that the cancellation signal stopped the command.
	Stopped bool
}

// prepared holds a staged project and its configured child command.
type prepared struct {
	// filesystem opens the workload logs.
	filesystem afero.Fs
	// directory contains the workload outputs.
	directory string
	// command is configured but has not started.
	command *exec.Cmd
}

// prepare stages the project and configures its command.
func prepare(ctx context.Context, filesystem afero.Fs, directory, id string, specification *jobpb.Workload) (*prepared, error) {
	if err := filesystem.MkdirAll(directory, 0o700); err != nil {
		return nil, fmt.Errorf("create job directory: %w", err)
	}
	project, err := prepareProject(ctx, filesystem, directory, id, specification)
	if err != nil {
		return nil, fmt.Errorf("prepare project: %w", err)
	}
	command := exec.Command("/bin/sh", "-c", specification.GetCommand())
	command.Dir = project
	command.Env = workloadEnvironment(command, id, specification)
	return &prepared{filesystem: filesystem, directory: directory, command: command}, nil
}

// run starts the prepared child command and records its observed outcome.
func (prepared *prepared) run(ctx context.Context) (outcome Outcome, err error) {
	stdout, err := prepared.filesystem.OpenFile(filepath.Join(prepared.directory, "stdout.log"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return Outcome{}, fmt.Errorf("open workload stdout: %w", err)
	}
	defer func() { err = errors.Join(err, closeResource(stdout, "workload stdout")) }()
	stderr, err := prepared.filesystem.OpenFile(filepath.Join(prepared.directory, "stderr.log"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return Outcome{}, fmt.Errorf("open workload stderr: %w", err)
	}
	defer func() { err = errors.Join(err, closeResource(stderr, "workload stderr")) }()
	prepared.command.Stdout = stdout
	prepared.command.Stderr = stderr
	outcome, err = runCommand(ctx, prepared.command)
	if err != nil {
		return outcome, fmt.Errorf("run project entry point: %w", err)
	}
	return outcome, nil
}

// closeResource adds context to a close error.
func closeResource(resource io.Closer, name string) error {
	if err := resource.Close(); err != nil {
		return fmt.Errorf("close %s: %w", name, err)
	}
	return nil
}
