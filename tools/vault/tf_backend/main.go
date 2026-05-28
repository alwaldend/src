package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"git.alwaldend.com/alwaldend/src/tools/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/al/pkg/al_plugin"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
)

var vaultConn = flag.String("vault_conn", "", "Vault connection")
var vaultAuth = flag.String("vault_auth", "", "Vault auth")
var vaultSecret = flag.String("vault_secret", "", "Vault secret path")
var vaultSecretMount = flag.String("vault_secret_mount", "", "Vault secret mount")
var logger = log.New(os.Stderr, "com.alwaldend.src.tools.vault.tf_backend ", log.Flags())

type Plugin struct {
}

func validate(targetUsername string, targetPasssword string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Write([]byte("missing basic auth"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != targetUsername || password != targetPasssword {
			w.Write([]byte("invalid auth"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func fetchSecret(vault_client *api.Client, w http.ResponseWriter, r *http.Request) (*api.KVSecret, bool) {
	kv := vault_client.KVv2(*vaultSecretMount)
	secret, err := kv.Get(r.Context(), *vaultSecret)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("could not get secret %s on mount %s: %s", *vaultSecret, *vaultSecretMount, err)))
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}
	return secret, true
}

func handlerState(vault *api.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret, ok := fetchSecret(vault, w, r)
		if !ok {
			return
		}
		lock, ok := secret.Data["lock"]
		if ok {
			w.Write([]byte(fmt.Sprintf("secret %s on mount %s is locked: %s", *vaultSecret, *vaultSecretMount, lock)))
			w.WriteHeader(http.StatusConflict)
			return
		}
		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("could not read request body: %s", err)))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			bodyEncoded := base64.StdEncoding.EncodeToString(body)
			_, err = vault.KVv2(*vaultSecretMount).Put(
				r.Context(),
				*vaultSecret,
				map[string]any{"state": bodyEncoded},
				api.WithCheckAndSet(secret.VersionMetadata.Version),
			)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("could not update secret %s on mount %s: %s", *vaultSecret, *vaultSecretMount, err)))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			state, ok := secret.Data["state"]
			if !ok {
				state = ""
			}
			body, err := base64.StdEncoding.DecodeString(state.(string))
			if err != nil {
				w.Write([]byte(fmt.Sprintf("could not decode state of secret %s on mount %s: %s", *vaultSecret, *vaultSecretMount, err)))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		}
	}
}

func handlerUnlock(vault *api.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret, ok := fetchSecret(vault, w, r)
		if !ok {
			return
		}
		_, ok = secret.Data["lock"]
		if !ok {
			w.Write([]byte(fmt.Sprintf("secret %s on mount %s is not locked", *vaultSecret, *vaultSecretMount)))
			w.WriteHeader(http.StatusConflict)
			return
		}
		state, ok := secret.Data["state"]
		if !ok {
			w.Write([]byte(fmt.Sprintf("secret %s on mount %s is missing state", *vaultSecret, *vaultSecretMount)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err := vault.KVv2(*vaultSecretMount).Put(
			r.Context(),
			*vaultSecret,
			map[string]any{"state": state},
			api.WithCheckAndSet(secret.VersionMetadata.Version),
		)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("could not update secret %s on mount %s: %s", *vaultSecret, *vaultSecretMount, err)))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
func handlerLock(vault *api.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret, ok := fetchSecret(vault, w, r)
		if !ok {
			return
		}
		lock, ok := secret.Data["lock"]
		if ok {
			w.Write([]byte(fmt.Sprintf("secret %s on mount %s is locked: %s", *vaultSecret, *vaultSecretMount, lock)))
			w.WriteHeader(http.StatusConflict)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("could not read request body: %s", err)))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = vault.KVv2(*vaultSecretMount).Patch(
			r.Context(),
			*vaultSecret,
			map[string]any{"lock": body},
			api.WithCheckAndSet(secret.VersionMetadata.Version),
		)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("could not update secret %s on mount %s: %s", *vaultSecret, *vaultSecretMount, err)))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (self Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (res *al_proto.PluginStartResponse, err error) {
	logger.Println("Creating a vault token")
	vault, err := al.NewVault(req.Config)
	if err != nil {
		return nil, fmt.Errorf("could not create vault client: %w", err)
	}
	vault_client, err := vault.Client(ctx, *vaultConn, *vaultAuth)
	if err != nil {
		return nil, fmt.Errorf("could not authorize vault client: %w", err)
	}
	username, password := uuid.New().String(), uuid.New().String()
	logger.Println("Created a vault token")
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, fmt.Errorf("could not listen on a port: %w", err)
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/lock", validate(username, password, handlerLock(vault_client)))
	handler.HandleFunc("/unlock", validate(username, password, handlerUnlock(vault_client)))
	handler.HandleFunc("/state", validate(username, password, handlerState(vault_client)))
	server := &http.Server{Handler: handler}
	go func() {
		curCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		for range curCtx.Done() {
		}
		err := server.Shutdown(context.Background())
		if err != nil {
			panic(fmt.Sprintf("could not stop the server: %s", err))
		}
	}()
	go func() {
		err := server.Serve(listener)
		if err != nil {
			panic(fmt.Sprintf("could not serve: %s", err))
		}
	}()
	path := fmt.Sprintf("%s/metadata/%s", *vaultSecretMount, filepath.Dir(*vaultSecret))
	list, err := vault_client.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("could not list secrets on path %s: %w", path, err)
	}
	create := true
	for _, key := range list.Data["keys"].([]any) {
		if key.(string) == filepath.Base(*vaultSecret) {
			create = false
			break
		}
	}
	if create {
		_, err := vault_client.KVv2(*vaultSecretMount).Put(ctx, *vaultSecret, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create secret %s on mount %s: %w", *vaultSecret, *vaultSecretMount, err)
		}
	}
	_, err = vault_client.KVv2(*vaultSecretMount).Get(ctx, *vaultSecret)
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret %s on mount %s: %w", *vaultSecret, *vaultSecretMount, err)
	}

	res = &al_proto.PluginStartResponse{
		Env: map[string]string{
			"TF_HTTP_USERNAME":       username,
			"TF_HTTP_PASSWORD":       password,
			"TF_HTTP_ADDRESS":        fmt.Sprintf("http://%s/state", listener.Addr().String()),
			"TF_HTTP_LOCK_ADDRESS":   fmt.Sprintf("http://%s/lock", listener.Addr().String()),
			"TF_HTTP_UNLOCK_ADDRESS": fmt.Sprintf("http://%s/unlock", listener.Addr().String()),
		},
	}
	return res, nil
}

func main() {
	flag.Parse()
	err := al_plugin.Serve(context.Background(), os.Stdin, os.Stdout, Plugin{})
	if err != nil {
		logger.Printf("failed to run plugin: %s", err.Error())
		os.Exit(1)
	}
}
