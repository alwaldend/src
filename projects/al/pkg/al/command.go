package al

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

type CmdCtx struct {
	Ctx             context.Context
	RequestShutdown context.CancelFunc
	Args            []string
	Getenv          func(string) string
	Stdin           io.Reader
	Stdout          io.Writer
	Stderr          io.Writer
	Logger          *log.Logger
}

func NewCmdCtx(ctx context.Context, prefix string) *CmdCtx {
	ctx, cancel := context.WithCancel(ctx)
	return &CmdCtx{
		Ctx:             ctx,
		RequestShutdown: cancel,
		Args:            os.Args,
		Getenv:          os.Getenv,
		Stdin:           os.Stdin,
		Stdout:          os.Stdout,
		Stderr:          os.Stderr,
		Logger:          log.New(os.Stderr, prefix, log.Flags()),
	}
}

// FindRunfiles discovers declared or executable-adjacent Bazel runtime files.
// A binary-only installation returns nil without error. Explicitly configured
// runtime files and discovered metadata must remain valid; errors are not treated
// as an absent optional context.
func FindRunfiles() (*runfiles.Runfiles, error) {
	manifest := os.Getenv("RUNFILES_MANIFEST_FILE")
	directory := os.Getenv("RUNFILES_DIR")
	if manifest == "" && directory != "" {
		info, err := os.Stat(directory)
		if err != nil {
			return nil, fmt.Errorf("inspect declared runfiles directory %q: %w", directory, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("declared runfiles directory %q is not a directory", directory)
		}
	}
	if manifest == "" && directory == "" {
		program := os.Args[0]
		if filepath.Base(program) == program {
			var err error
			program, err = os.Executable()
			if err != nil {
				return nil, fmt.Errorf("locate executable for runfiles discovery: %w", err)
			}
		}
		found := false
		for _, suffix := range []string{".runfiles_manifest", ".runfiles"} {
			path := program + suffix
			info, err := os.Stat(path)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("inspect adjacent runfiles %q: %w", path, err)
			}
			if (suffix == ".runfiles_manifest" && info.Mode().IsRegular()) ||
				(suffix == ".runfiles" && info.IsDir()) {
				found = true
				break
			}
		}
		if !found {
			return nil, nil
		}
	}
	runtimeFiles, err := runfiles.New()
	if err != nil {
		return nil, fmt.Errorf("load Bazel runfiles: %w", err)
	}
	return runtimeFiles, nil
}

// SetRunfilesEnv propagates available Bazel runtime files without requiring them
// for an installed command that has no runfiles-backed inputs.
func SetRunfilesEnv(cmd *exec.Cmd) error {
	runtimeFiles, err := FindRunfiles()
	if err != nil {
		return fmt.Errorf("could not create runfiles env: %w", err)
	}
	if runtimeFiles != nil {
		cmd.Env = append(cmd.Env, runtimeFiles.Env()...)
	}
	return nil
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
	ctx := args.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, args.Name, args.Args...)
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
