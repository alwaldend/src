package al

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

type CmdCtx struct {
	Ctx      context.Context
	Args     []string
	Getenv   func(string) string
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	LogFlags int
}

func NewCmdCtx(ctx context.Context) *CmdCtx {
	return &CmdCtx{
		Ctx:      ctx,
		Args:     os.Args,
		Getenv:   os.Getenv,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		LogFlags: log.Flags(),
	}
}

// Set runfiles info for a command
func SetRunfilesEnv(cmd *exec.Cmd) error {
	runfilesEnv, err := runfiles.Env()
	if err != nil {
		return fmt.Errorf("could not create runfiles env: %w", err)
	}
	cmd.Env = append(cmd.Env, runfilesEnv...)
	return err
}

type CommandArgs struct {
	Name               string
	Args               []string
	Ctx                context.Context
	Stdout             io.Writer
	Stderr             io.Writer
	Stdin              io.Reader
	DisableRunfilesEnv bool
}

func RunCommand(args CommandArgs) error {
	cmd, err := Command(args)
	if err != nil {
		return fmt.Errorf("could not create command: %w", err)
	}
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not run command: %w", err)
	}
	return nil
}

func Command(args CommandArgs) (*exec.Cmd, error) {
	cmd := exec.Command(args.Name, args.Args...)
	if args.Stdout == nil {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = args.Stdout
	}
	if args.Stderr == nil {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = args.Stderr
	}
	if args.Stdin == nil {
		cmd.Stdin = os.Stdin
	} else {
		cmd.Stdin = args.Stdin
	}
	cmd.Env = append(cmd.Env, os.Environ()...)
	if !args.DisableRunfilesEnv {
		err := SetRunfilesEnv(cmd)
		if err != nil {
			return nil, fmt.Errorf("could not set runfiles env: %w", err)
		}
	}
	return cmd, nil
}
