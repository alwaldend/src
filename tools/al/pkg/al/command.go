package al

import (
	"context"
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

type CommandArgs struct {
	Ctx                context.Context
	Config             *al_proto.Config
	Name               string
	Args               []string
	SecretEnv          []string
	Vault              *Vault
	SetVaultEnv        bool
	DisableRunfilesEnv bool
}

func Command(args CommandArgs) (*exec.Cmd, error) {
	cmd := exec.Command(args.Name, args.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	if !args.DisableRunfilesEnv {
		err := SetRunfilesEnv(cmd)
		if err != nil {
			return nil, fmt.Errorf("could not set runfiles env: %w", err)
		}
	}
	if args.SetVaultEnv {
		vaultEnv, err := args.Vault.DefaultEnv(args.Ctx)
		if err != nil {
			return nil, fmt.Errorf("could not create vault env: %w", err)
		}
		cmd.Env = append(cmd.Env, vaultEnv...)
	}
	for _, secretName := range args.SecretEnv {
		env, err := args.Vault.SecretEnv(args.Ctx, secretName)
		if err != nil {
			return nil, fmt.Errorf("could not template env for %s: %w", secretName, err)
		}
		cmd.Env = append(cmd.Env, env...)
	}
	return cmd, nil
}
