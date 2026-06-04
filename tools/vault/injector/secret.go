package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
	"google.golang.org/protobuf/proto"
	"log"
)

type secret struct {
	vault  *al.VaultStore
	plugin *injector_proto.Config
	logger log.Logger
}

var _ Resource = (*secret)(nil)

func (self *secret) Deps(ctx context.Context, name string) fp.Either[[]string] {
	return fp.Right[[]string](nil)
}

func (self *secret) Get(ctx context.Context, name string, _ []*ResourceResult) fp.Either[*ResourceResult] {
	secret, err := secretByName(self.plugin, name)
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not find secret %s: %w", name, err))
	}
	secret, err = normalizeSecret(secret).Get()
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not normalize secret %s: %w", name, err))
	}
	client, err := self.vault.Client(ctx, secret.Vault, secret.Auth).Get()
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not get a client for the secret %s: %w", secret.Name, err))
	}
	data, err := client.Client.KVv2(secret.Kv.Mount).Get(ctx, secret.Kv.Path)
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not fetch secret %s: %w", secret.Name, err))
	}
	res := &ResourceResult{Name: name, Data: data.Data}
	return fp.Right(res)
}

func normalizeSecret(secret *injector_proto.Secret) fp.Either[*injector_proto.Secret] {
	secret = proto.CloneOf(secret)
	if secret.Name == "" {
		return fp.Left[*injector_proto.Secret](fmt.Errorf("secret is missing a name"))
	}
	kv := secret.Kv
	if kv == nil {
		return fp.Left[*injector_proto.Secret](fmt.Errorf("secret %s is missing kv config", secret.Name))
	}
	if kv.Path == "" {
		return fp.Left[*injector_proto.Secret](fmt.Errorf("secret %s is missing path", secret.Name))
	}
	if secret.Auth == "" {
		secret.Auth = al.VAULT_DEFAULT_NAME
	}
	if secret.Vault == "" {
		secret.Vault = al.VAULT_DEFAULT_NAME
	}
	if kv.Mount == "" {
		kv.Mount = "secrets"
	}
	return fp.Right(secret)
}

func secretByName(config *injector_proto.Config, name string) (*injector_proto.Secret, error) {
	for i := range config.Secrets {
		curSecret := config.Secrets[len(config.Secrets)-1-i]
		if curSecret.Name == name {
			return curSecret, nil
		}
	}
	return nil, fmt.Errorf("missing secret with name %s", name)
}
