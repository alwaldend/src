package main

import (
	"fmt"
	"os"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/secrets/main/proto/contracts"
	"google.golang.org/protobuf/encoding/protojson"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

type BuilderSecretOpts struct {
	Env     []string
	Backend []string
	Secret  []string
	Output  []string
}

func (self *Builder) BuildSecret(opts *BuilderSecretOpts) error {
	secret := &contracts.Secret{
		Env: make(map[string]string),
	}
	for _, env := range opts.Env {
		envSplit := strings.SplitN(env, "=", 2)
		secret.Env[envSplit[0]] = envSplit[1]
	}
	for _, backendPath := range opts.Backend {
		backend := &contracts.Backend{}
		backendContent, err := os.ReadFile(backendPath)
		if err != nil {
			return fmt.Errorf("could not read backend file %s: %w", backendPath, err)
		}
		err = protojson.Unmarshal(backendContent, backend)
		if err != nil {
			return fmt.Errorf("could not unmarshal backend file %s: %w", backendPath, err)
		}
		secret.Backends = append(secret.Backends, backend)
	}
	for _, secretPath := range opts.Secret {
		curSecret := &contracts.Secret{}
		curSecretContent, err := os.ReadFile(secretPath)
		if err != nil {
			return fmt.Errorf("could not read secret file %s: %w", secretPath, err)
		}
		err = protojson.Unmarshal(curSecretContent, curSecret)
		if err != nil {
			return fmt.Errorf("could not unmarshal secret file %s: %w", secretPath, err)
		}
		secret.Secrets = append(secret.Secrets, curSecret)
	}
	res, err := protojson.MarshalOptions{Indent: "    "}.Marshal(secret)
	if err != nil {
		return fmt.Errorf("could not marshal the secret: %w", err)
	}
	for _, filePath := range opts.Output {
		err := os.WriteFile(filePath, res, os.FileMode(0o555))
		if err != nil {
			return fmt.Errorf("could not write the secret to %s: %w", filePath, err)
		}
	}
	return nil
}
