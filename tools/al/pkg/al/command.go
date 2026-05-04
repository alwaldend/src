package al

import (
	"context"
	"errors"
	"fmt"
	"git.alwaldend.com/alwaldend/src/tools/al/api/al_proto"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

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
	Name               string
	Args               []string
	Ctx                context.Context
	Stdout             io.Writer
	Stderr             io.Writer
	Stdin              io.Reader
	DisableRunfilesEnv bool
	SecretEnv          []string
	VaultEnv           []string
	Vault              *Vault
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
	if args.Vault == nil && (len(args.VaultEnv)+len(args.SecretEnv)) != 0 {
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
	// for _, secretName := range args.SecretEnv {
	// 	env, err := args.Vault.SecretEnv(args.Ctx, secretName)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("could not template env for %s: %w", secretName, err)
	// 	}
	// 	cmd.Env = append(cmd.Env, env...)
	// }
	return cmd, nil
}

type resourceHandler struct {
	secrets map[string]map[string]any
	files   map[string]*ResourceContextFile
	vault   *Vault
	config  *al_proto.Config
}

type ResourceContextFile struct {
	Path string
}

type ResourceContext struct {
	Files   map[string]*ResourceContextFile
	File    *ResourceContextFile
	Secrets map[string]map[string]any
	Secret  map[string]any
}

func newResourceHandler(config *al_proto.Config, vault *Vault) *resourceHandler {
	return &resourceHandler{
		config:  config,
		secrets: make(map[string]map[string]any),
		files:   make(map[string]*ResourceContextFile),
		vault:   vault,
	}
}

func (self *resourceHandler) get_secret(ctx context.Context, name string) (map[string]any, error) {
	res, ok := self.secrets[name]
	if ok {
		return res, nil
	}
	res, err := self.vault.FetchSecret(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret %s: %w", name, err)
	}
	self.secrets[name] = res
	return res, nil
}

func (self *resourceHandler) get_secrets(ctx context.Context, templateCtx *ResourceContext, names []string) error {
	templateCtx.Secrets = map[string]map[string]any{}
	for _, secretName := range names {
		secret, err := self.get_secret(ctx, secretName)
		if err != nil {
			return fmt.Errorf("could not prepare secret %s: %w", secretName, err)
		}
		templateCtx.Secrets[secretName] = secret
		templateCtx.Secret = secret
	}
	return nil
}

func (self *resourceHandler) get_files(ctx context.Context, templateCtx *ResourceContext, names []string) error {
	templateCtx.Files = map[string]*ResourceContextFile{}
	for _, fileName := range names {
		file, err := self.get_file(ctx, fileName)
		if err != nil {
			return fmt.Errorf("could not prepare file %s: %w", fileName, err)
		}
		templateCtx.Files[fileName] = file
		templateCtx.File = file
	}
	return nil
}

func (self *resourceHandler) get_file(ctx context.Context, name string) (*ResourceContextFile, error) {
	res, ok := self.files[name]
	if ok {
		return res, nil
	}
	fileConfig, err := FileByName(self.config, name)
	if err != nil {
		return nil, fmt.Errorf("could not find file config: %w", err)
	}
	templateCtx := &ResourceContext{}
	err = self.get_secrets(ctx, templateCtx, fileConfig.Secrets)
	if err != nil {
		return nil, fmt.Errorf("could not prepare secrets: %w", err)
	}
	err = self.get_files(ctx, templateCtx, fileConfig.Files)
	if err != nil {
		return nil, fmt.Errorf("could not prepare files: %w", err)
	}
	tmpl, err := template.New(fmt.Sprintf("al_file_%s", name)).Parse(fileConfig.Value)
	if err != nil {
		return nil, fmt.Errorf("could not template secret file %s: %w", name, err)
	}
	tmp, err := os.CreateTemp("", fmt.Sprintf("al_file_%s_*.txt", name))
	if err != nil {
		return nil, fmt.Errorf("could not create temporary file: %w", err)
	}
	defer tmp.Close()
	err = tmpl.Execute(tmp, templateCtx)
	if err != nil {
		return nil, fmt.Errorf("could not template file %s: %w", name, err)
	}
	self.files[name] = res
	return &ResourceContextFile{Path: tmp.Name()}, nil
}

func (self *resourceHandler) Clean() error {
	var res error
	for _, file := range self.files {
		res = errors.Join(res, os.Remove(file.Path))
	}
	return res
}
