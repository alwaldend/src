package al

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	Name               string
	Args               []string
	DisableRunfilesEnv bool
	SecretEnv          []string
	VaultEnv        []string
	Vault              *Vault
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
    if args.Vault == nil && (len(args.VaultEnv) + len(args.SecretEnv)) != 0 {
        return nil, fmt.Errorf("cannot apply vault options: missing vault")
    }
    for _, vaultEnv := range args.VaultEnv {
        vaultEnvSplit := strings.SplitN(vaultEnv, ":", 2)
        if len(vaultEnvSplit) != 2 {
            return nil, fmt.Errorf("invalid vault env: %s", vaultEnv)
        }
        vaultEnvVars, err := args.Vault.Env(args.Ctx, vaultEnvSplit[0], vaultEnvSplit[1])
        if err != nil {
            return nil, fmt.Errorf("could not create vault env: %w", err)
        }
        cmd.Env = append(cmd.Env, vaultEnvVars...)
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
