package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
)

var logger = log.New(os.Stderr, "com.alwaldend.src.tools.vault.injector ", log.Flags())

type Plugin struct {
}

func (self Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func main() {
	fp.PipeE(
		al_plugin.ServeDefault,
		func(err error) error { logger.Printf("could not serve the plugin: %s", err); os.Exit(1); return nil },
	)(Plugin{})
}

type ResourceHandler struct {
	ctx      context.Context
	secrets  map[string]map[string]any
	vaultOps map[string]map[string]any
	files    map[string]*ResourceContextFile
	vault    *al.Vault
	config   *al_proto.Config
	logger   *log.Logger
}

type ResourceContextFile struct {
	Path string
}

type ResourceContext struct {
	Files    map[string]*ResourceContextFile
	File     *ResourceContextFile
	VaultOps map[string]map[string]any
	VaultOp  map[string]any
	Secrets  map[string]map[string]any
	Secret   map[string]any
}

func NewResourceHandler(ctx context.Context, config *al_proto.Config, vault *al.Vault, stderr io.Writer) *ResourceHandler {
	if stderr == nil {
		stderr = os.Stderr
	}
	return &ResourceHandler{
		ctx:      ctx,
		config:   config,
		secrets:  make(map[string]map[string]any),
		vaultOps: make(map[string]map[string]any),
		files:    make(map[string]*ResourceContextFile),
		vault:    vault,
		logger:   log.New(stderr, "com.alwaldend.src.projects.al.pkg.al.ResourceHandler ", log.Flags()),
	}
}

func (self *ResourceHandler) Clean() error {
	var res error
	for _, file := range self.files {
		self.logger.Printf("removing file: %s", file.Path)
		res = errors.Join(res, os.RemoveAll(file.Path))
	}
	return res
}

type PrepareCommandArgs struct {
	Env         []string
	EnvVault    []string
	EnvLabels   []string
	EnvDisabled bool
	Files       []string
}

func (self *ResourceHandler) PrepareCommand(ctx context.Context, cmd *exec.Cmd, args *PrepareCommandArgs) error {
	for _, file := range args.Files {
		_, err := self.getFile(ctx, file)
		if err != nil {
			return fmt.Errorf("could not prepare file %s: %w", file, err)
		}
	}
	if !args.EnvDisabled {
		err := self.prepareEnv(ctx, cmd, args)
		if err != nil {
			return fmt.Errorf("could not prepare env: %w", err)
		}
	}
	return nil
}

func (self *ResourceHandler) prepareEnv(ctx context.Context, cmd *exec.Cmd, args *PrepareCommandArgs) error {
	for _, envVault := range args.EnvVault {
		split := strings.SplitN(envVault, ":", 2)
		vaultName, authName := "", ""
		if len(split) == 1 {
			vaultName = split[0]
		} else if len(split) == 2 {
			vaultName, authName = split[0], split[1]
		}
		env, err := fp.Get(self.vault.Env(ctx, vaultName, authName))
		if err != nil {
			return fmt.Errorf("could not prepare vault env %s: %w", envVault, err)
		}
		cmd.Env = append(cmd.Env, env...)
	}
	for _, envLabel := range args.EnvLabels {
		labelName, labelVal := "", ""
		split := strings.SplitN(envLabel, "=", 2)
		if len(split) > 0 {
			labelName = split[0]
		}
		if len(split) > 1 {
			labelVal = split[1]
		}
		envs, err := al.EnvByLabel(self.config, labelName, labelVal)
		if err != nil {
			return fmt.Errorf("could not get env for label %s: %w", envLabel, err)
		}
		for _, envConfig := range envs {
			env, err := self.getEnv(ctx, envConfig.Name)
			if err != nil {
				return fmt.Errorf("could not prepare env %s: %w", envConfig.Name, err)
			}
			cmd.Env = append(cmd.Env, env)
		}
	}
	for _, envName := range args.Env {
		env, err := self.getEnv(ctx, envName)
		if err != nil {
			return fmt.Errorf("could not prepare env %s: %w", envName, err)
		}
		cmd.Env = append(cmd.Env, env)
	}
	return nil
}

func (self *ResourceHandler) getVaultOp(ctx context.Context, name string) (map[string]any, error) {
	self.logger.Printf("executing a vault operation: %s", name)
	res, ok := self.vaultOps[name]
	if ok {
		return res, nil
	}
	op, err := al.VaultOpByName(self.config, name)
	if err != nil {
		return nil, fmt.Errorf("could not find by name: %w", err)
	}
	res, err = fp.Get(self.vault.VaultOp(ctx, op))
	if err != nil {
		return nil, fmt.Errorf("could not execute vault op: %w", err)
	}
	self.secrets[name] = res
	return res, nil
}

func (self *ResourceHandler) getSecret(ctx context.Context, name string) (map[string]any, error) {
	self.logger.Printf("fetching a vault secret: %s", name)
	res, ok := self.secrets[name]
	if ok {
		return res, nil
	}
	res, err := fp.Get(self.vault.FetchSecret(ctx, name))
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret %s: %w", name, err)
	}
	self.secrets[name] = res
	return res, nil
}

func (self *ResourceHandler) addSecrets(ctx context.Context, templateCtx *ResourceContext, names []string) error {
	templateCtx.Secrets = map[string]map[string]any{}
	for _, secretName := range names {
		secret, err := self.getSecret(ctx, secretName)
		if err != nil {
			return fmt.Errorf("could not prepare secret %s: %w", secretName, err)
		}
		templateCtx.Secrets[secretName] = secret
		templateCtx.Secret = secret
	}
	return nil
}

func (self *ResourceHandler) addFiles(ctx context.Context, templateCtx *ResourceContext, names []string) error {
	templateCtx.Files = map[string]*ResourceContextFile{}
	for _, fileName := range names {
		file, err := self.getFile(ctx, fileName)
		if err != nil {
			return fmt.Errorf("could not prepare file %s: %w", fileName, err)
		}
		templateCtx.Files[fileName] = file
		templateCtx.File = file
	}
	return nil
}

func (self *ResourceHandler) addVaultOps(ctx context.Context, templateCtx *ResourceContext, names []string) error {
	templateCtx.VaultOps = make(map[string]map[string]any)
	for _, opName := range names {
		op, err := self.getVaultOp(ctx, opName)
		if err != nil {
			return fmt.Errorf("could not prepare vault op %s: %w", opName, err)
		}
		templateCtx.VaultOps[opName] = op
		templateCtx.VaultOp = op
	}
	return nil
}

func (self *ResourceHandler) getEnv(ctx context.Context, name string) (string, error) {
	self.logger.Printf("setting an env variable: %s", name)
	env, err := al.EnvByName(self.config, name)
	if err != nil {
		return "", fmt.Errorf("could not find env %s: %w", name, err)
	}
	content, err := self.template(ctx, env.Value, env.Files, env.Secrets, env.VaultOps)
	if err != nil {
		return "", fmt.Errorf("could not template env %s: %w", name, err)
	}
	return fmt.Sprintf("%s=%s", name, content), nil
}

func (self *ResourceHandler) getFile(ctx context.Context, name string) (*ResourceContextFile, error) {
	self.logger.Printf("creating a file: %s", name)
	res, ok := self.files[name]
	if ok {
		return res, nil
	}
	fileConfig, err := al.FileByName(self.config, name)
	if err != nil {
		return nil, fmt.Errorf("could not find file config: %w", err)
	}
	value := ""
	if fileConfig.Value != "" {
		value = fileConfig.Value
	} else if fileConfig.FromFile != "" {
		valueBytes, err := os.ReadFile(fileConfig.FromFile)
		if err != nil {
			return nil, fmt.Errorf("could not read from_file %s: %w", fileConfig.FromFile, err)
		}
		value = string(valueBytes)
	}
	if value == "" {
		return nil, fmt.Errorf("missing file contents")
	}
	content, err := self.template(ctx, value, fileConfig.Files, fileConfig.Secrets, fileConfig.VaultOps)
	if err != nil {
		return nil, fmt.Errorf("could not template file %s: %w", name, err)
	}
	tmp, err := os.CreateTemp("", fmt.Sprintf("al_file_%s_*.txt", name))
	if err != nil {
		return nil, fmt.Errorf("could not create temporary file: %w", err)
	}
	defer tmp.Close()
	res = &ResourceContextFile{Path: tmp.Name()}
	self.files[name] = res
	_, err = tmp.WriteString(content)
	if err != nil {
		return nil, fmt.Errorf("could not write to the temp file: %w", err)
	}
	return res, nil
}

func (self *ResourceHandler) template(ctx context.Context, value string, files []string, secrets []string, vaultOps []string) (string, error) {
	templateCtx := &ResourceContext{}
	err := self.addSecrets(ctx, templateCtx, secrets)
	if err != nil {
		return "", fmt.Errorf("could not prepare secrets: %w", err)
	}
	err = self.addFiles(ctx, templateCtx, files)
	if err != nil {
		return "", fmt.Errorf("could not prepare files: %w", err)
	}
	err = self.addVaultOps(ctx, templateCtx, vaultOps)
	if err != nil {
		return "", fmt.Errorf("could not prepare vault ops: %w", err)
	}
	tmpl, err := template.New("resource_handler.template").Funcs(
		template.FuncMap{
			"join": func(elems []any, sep string) string {
				elemsString := make([]string, len(elems))
				for _, elem := range elems {
					elemsString = append(elemsString, elem.(string))
				}
				return strings.Join(elemsString, sep)
			},
		},
	).Option("missingkey=error").Parse(value)
	if err != nil {
		return "", fmt.Errorf("could not parse template '%s': %w", value, err)
	}
	var buff bytes.Buffer
	err = tmpl.Execute(&buff, templateCtx)
	if err != nil {
		return "", fmt.Errorf("could not execute template '%s': %w", value, err)
	}
	return buff.String(), nil
}
