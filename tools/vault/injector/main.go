package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

var logger = log.New(os.Stderr, "com.alwaldend.src.tools.vault.injector ", log.Flags())

type mapa = map[string]any

type plugin struct {
	cleanConsume chan string
	cleanFiles   []string
	cleanStart   chan struct{}
	ctx          context.Context
	v            *al.VaultStore
}

func (self plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func clean(p *plugin) {
	<-p.cleanStart
	for _, c := range p.cleanFiles {
		logger.Printf("removing path %s", c)
		err := os.RemoveAll(c)
		if err != nil {
			logger.Printf("could not clean: could not remove path %s: %s", c, err)
		}
	}
}

func consume(p *plugin) {
	select {
	case c := <-p.cleanConsume:
		p.cleanFiles = append(p.cleanFiles, c)
	case <-p.ctx.Done():
		close(p.cleanConsume)
		close(p.cleanStart)
	}
}

func setup(p *plugin) {
	p.cleanConsume = make(chan string)
	p.cleanStart = make(chan struct{}, 1)
	p.ctx, _ = signal.NotifyContext(p.ctx, syscall.SIGINT, syscall.SIGTERM)
}

func main() {
	fp.Pipe4E(
		fp.Compute0(setup),
		fp.Compute0G(clean),
		fp.Compute0G(consume),
		al_plugin.ServeDefault,
		func(err error) error { logger.Printf("could not serve the plugin: %s", err); os.Exit(1); return nil },
	)(&plugin{})
}

type resourceHandler struct {
	ctx      context.Context
	secrets  map[string]map[string]any
	vaultOps map[string]map[string]any
	files    map[string]*resourceContextFile
	vault    *al.VaultStore
	config   *al_proto.Config
	plugin   *injector_proto.Config
	logger   *log.Logger
}

type resourceContextFile struct {
	Path string
}

type resourceContext struct {
	Files    map[string]*resourceContextFile
	File     *resourceContextFile
	VaultOps map[string]map[string]any
	VaultOp  map[string]any
	Secrets  map[string]map[string]any
	Secret   map[string]any
}

func newResourceHandler(ctx context.Context, config *al_proto.Config, vault *al.VaultStore, stderr io.Writer) *resourceHandler {
	if stderr == nil {
		stderr = os.Stderr
	}
	return &resourceHandler{
		ctx:      ctx,
		config:   config,
		secrets:  make(map[string]map[string]any),
		vaultOps: make(map[string]map[string]any),
		files:    make(map[string]*resourceContextFile),
		vault:    vault,
		logger:   log.New(stderr, "com.alwaldend.src.projects.al.pkg.al.ResourceHandler ", log.Flags()),
	}
}

type prepareCommandArgs struct {
	Env         []string
	EnvVault    []string
	EnvLabels   []string
	EnvDisabled bool
	Files       []string
}

func (self *resourceHandler) PrepareCommand(ctx context.Context, cmd *exec.Cmd, args *prepareCommandArgs) error {
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

func (self *resourceHandler) prepareEnv(ctx context.Context, cmd *exec.Cmd, args *prepareCommandArgs) error {
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
		envs, err := EnvByLabel(self.plugin, labelName, labelVal)
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

func (self *resourceHandler) getSecret(ctx context.Context, name string) (map[string]any, error) {
	self.logger.Printf("fetching a vault secret: %s", name)
	res, ok := self.secrets[name]
	if ok {
		return res, nil
	}
	res, err := self.fetchSecret(ctx, name).Get()
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret %s: %w", name, err)
	}
	self.secrets[name] = res
	return res, nil
}

func (self *resourceHandler) addSecrets(ctx context.Context, templateCtx *resourceContext, names []string) error {
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

func (self *resourceHandler) addFiles(ctx context.Context, templateCtx *resourceContext, names []string) error {
	templateCtx.Files = map[string]*resourceContextFile{}
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

func (self *resourceHandler) addVaultOps(ctx context.Context, templateCtx *resourceContext, names []string) error {
	templateCtx.VaultOps = make(map[string]map[string]any)
	for _, opName := range names {
		op, err := self.getVaultOp(ctx, opName).Get()
		if err != nil {
			return fmt.Errorf("could not prepare vault op %s: %w", opName, err)
		}
		templateCtx.VaultOps[opName] = op
		templateCtx.VaultOp = op
	}
	return nil
}

func (self *resourceHandler) getEnv(ctx context.Context, name string) (string, error) {
	self.logger.Printf("setting an env variable: %s", name)
	env, err := EnvByName(self.plugin, name)
	if err != nil {
		return "", fmt.Errorf("could not find env %s: %w", name, err)
	}
	content, err := self.template(ctx, env.Value, env.Files, env.Secrets, env.VaultOps)
	if err != nil {
		return "", fmt.Errorf("could not template env %s: %w", name, err)
	}
	return fmt.Sprintf("%s=%s", name, content), nil
}

func EnvByLabel(config *injector_proto.Config, name string, value string) ([]*injector_proto.Env, error) {
	if name == "" {
		return nil, fmt.Errorf("missing label name")
	}
	res := []*injector_proto.Env{}
	for i := range config.Env {
		curEnv := config.Env[len(config.Env)-1-i]
		actualVal, ok := curEnv.Labels[name]
		if ok && (value == "" || value == actualVal) {
			res = append(res, curEnv)
		}
	}
	return res, nil

}

func EnvByName(config *injector_proto.Config, name string) (*injector_proto.Env, error) {
	for i := range config.Env {
		curEnv := config.Env[len(config.Env)-1-i]
		if curEnv.Name == name {
			return curEnv, nil
		}
	}
	return nil, fmt.Errorf("missing env with name %s", name)
}
