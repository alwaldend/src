package al

import (
	"fmt"
	"os"
	"os/exec"

	"git.alwaldend.com/alwaldend/src/tools/al/api/al_proto"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

// Set runfiles info for a command
func SetRunfilesEnv(cmd *exec.Cmd) error {
	runfilesEnv, err := runfiles.Env()
	if err != nil {
		return fmt.Errorf("could not create runfiles env: %w", err)
	}
	cmd.Env = append(cmd.Env, runfilesEnv...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	return err
}

type RunCommandArgs struct {
	config *al_proto.Config
	name string
	args []string
	env []string
}

// func RunCommand(args RunCommandArgs) (*exec.Cmd, error) {
// 	cmd, err := Command(args.name, args.args...)
// 	for
// }

func Command(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	err := SetRunfilesEnv(cmd)
	if err != nil {
		return nil, fmt.Errorf("could not set runfiles env: %w", err)
	}
	return cmd, nil
}
