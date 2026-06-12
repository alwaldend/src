package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/tf_backend/tf_backend_proto"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
)

type requestType string
type secretType string

const (
	reqTypeUpdate   = requestType("update")
	reqTypeGet      = requestType("get")
	reqTypeLock     = requestType("lock")
	reqTypeUnlock   = requestType("unlock")
	secretTypeState = secretType("state")
	secretTypeLock  = secretType("lock")
	stateKey        = "state"
	lockKey         = "lock"
)

type TfBackend struct {
	username   string
	password   string
	server     *Server
	config     *tf_backend_proto.Config
	statePath  string
	stateMount string
	lockPath   string
	lockMount  string
	vault      *api.Client
	ctx        *al.CmdCtx
}

func NewTfBackend(ctx *al.CmdCtx, config *al_proto.Config, backend *tf_backend_proto.Config) (*TfBackend, error) {
	if backend.VaultSecret == "" || backend.VaultSecretMount == "" {
		return nil, fmt.Errorf("missing secret config")
	}
	res := &TfBackend{
		username: uuid.NewString(),
		password: uuid.NewString(),
		config:   backend,
		ctx:      ctx,
	}
	res.statePath = fmt.Sprintf("%s/state", backend.VaultSecret)
	res.stateMount = backend.VaultSecretMount
	res.lockPath = fmt.Sprintf("%s/lock", backend.VaultSecret)
	res.lockMount = backend.VaultSecretMount
	vault := al.NewVault(config)
	client, err := vault.Client(ctx.Ctx, backend.VaultConn, backend.VaultAuth).Get()
	if err != nil {
		return nil, fmt.Errorf("could not create vault client: %w", err)
	}
	res.vault = client.Client
	res.server = NewServer(ctx, res.reqHandler)
	return res, nil
}

func (self *TfBackend) Start(ctx context.Context) error {
	return self.server.Start(ctx)
}

func (self *TfBackend) Stop(ctx context.Context) error {
	return self.server.Stop(ctx)
}

func (self *TfBackend) Env() (map[string]string, error) {
	addr, err := self.server.Address()
	if err != nil {
		return nil, fmt.Errorf("could not get server address: %w", err)
	}
	return map[string]string{
		"TF_HTTP_USERNAME":       self.username,
		"TF_HTTP_PASSWORD":       self.password,
		"TF_HTTP_ADDRESS":        fmt.Sprintf("http://%s/state", addr),
		"TF_HTTP_UPDATE_METHOD":  http.MethodPost,
		"TF_HTTP_LOCK_ADDRESS":   fmt.Sprintf("http://%s/lock", addr),
		"TF_HTTP_LOCK_METHOD":    http.MethodPost,
		"TF_HTTP_UNLOCK_ADDRESS": fmt.Sprintf("http://%s/unlock", addr),
		"TF_HTTP_UNLOCK_METHOD":  http.MethodPost,
	}, nil
}

func (self *TfBackend) reqHandler(w http.ResponseWriter, r *http.Request) {
	if err := self.handleRequest(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		self.ctx.Logger.Printf("error: %s", err)
	}
}

func (self *TfBackend) handleRequest(w http.ResponseWriter, r *http.Request) error {
	username, password, ok := r.BasicAuth()
	if !ok {
		return fmt.Errorf("missng basic auth")
	}
	if username != self.username || password != self.password {
		return fmt.Errorf("invalid auth")
	}
	reqType, err := self.requestType(r)
	if err != nil {
		return fmt.Errorf("could not determine request type: %w", err)
	}
	bodyRaw, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("could not read request body: %w", err)
	}
	body := make(map[string]any)
	if len(bodyRaw) != 0 {
		if err := json.Unmarshal(bodyRaw, &body); err != nil {
			return fmt.Errorf("could not unmarshal request body: %w", err)
		}
	}
	lock, err := self.fetchSecret(r, secretTypeLock)
	if err != nil {
		return fmt.Errorf("could not fetch lock: %w", err)
	}
	if err := self.checkLock(r, reqType, lock); err != nil {
		return fmt.Errorf("invalid lock: %w", err)
	}
	switch reqType {
	case reqTypeGet:
		state, err := self.fetchSecret(r, secretTypeState)
		if err != nil {
			return fmt.Errorf("could not fetch state secret: %w", err)
		}
		if err := self.writeState(r, w, state); err != nil {
			return fmt.Errorf("could not write state: %w", err)
		}
	case reqTypeUpdate:
		state, err := self.fetchSecret(r, secretTypeState)
		if err != nil {
			return fmt.Errorf("could not fetch state secret: %w", err)
		}
		if err := self.updateState(r, state, body); err != nil {
			return fmt.Errorf("could not update state: %w", err)
		}
	case reqTypeLock:
		if err := self.lockState(r, lock, body); err != nil {
			return fmt.Errorf("could not lock state: %w", err)
		}
	case reqTypeUnlock:
		if err := self.unlockState(r, lock); err != nil {
			return fmt.Errorf("could not unlock state: %w", err)
		}
	default:
		return fmt.Errorf("invalid request type")
	}
	return nil
}

func (self *TfBackend) requestType(r *http.Request) (requestType, error) {
	path := r.URL.Path
	method := r.Method
	var reqType requestType
	switch {
	case path == "/lock" && method == http.MethodPost:
		reqType = reqTypeLock
	case path == "/unlock" && method == http.MethodPost:
		reqType = reqTypeUnlock
	case path == "/state" && method == http.MethodPost:
		reqType = reqTypeUpdate
	case path == "/state" && method == http.MethodGet:
		reqType = reqTypeGet
	default:
		return "", fmt.Errorf("invalid request path or method")
	}
	return reqType, nil
}

func (self *TfBackend) fetchSecret(r *http.Request, st secretType) (*api.KVSecret, error) {
	var mount, path string
	switch st {
	case secretTypeLock:
		mount, path = self.lockMount, self.lockPath
	case secretTypeState:
		mount, path = self.stateMount, self.statePath
	default:
		return nil, fmt.Errorf("invalid secret type: %s", st)
	}
	kv := self.vault.KVv2(mount)
	list, err := self.vault.Logical().List(fmt.Sprintf("%s/metadata/%s", mount, filepath.Dir(path)))
	if err != nil {
		return nil, fmt.Errorf("could not list secrets: %w", err)
	}
	create := true
	if list != nil {
		keys, ok := list.Data["keys"]
		if ok {
			keysList, ok := keys.([]any)
			if !ok {
				return nil, fmt.Errorf("invalid keys field for some reason: %s", keys)
			}
			for _, keyAny := range keysList {
				key, ok := keyAny.(string)
				if !ok {
					return nil, fmt.Errorf("invalid key type for some reason: %s", key)
				}
				if key == filepath.Base(path) {
					create = false
					break
				}
			}
		}
	}
	if create {
		if _, err := kv.Put(r.Context(), path, map[string]any{}); err != nil {
			return nil, fmt.Errorf("could not create secret: %w", err)
		}
	}
	state, err := kv.Get(r.Context(), path)
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret: %w", err)
	}
	return state, nil
}

func (self *TfBackend) checkLock(r *http.Request, reqType requestType, lockObj *api.KVSecret) error {
	lock, ok := lockObj.Data[lockKey]
	curId := r.URL.Query().Get("ID")
	switch {
	case reqType == reqTypeUpdate:
		if !ok {
			return fmt.Errorf("missing lock field %s", lockKey)
		}
		lockMap, ok := lock.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid lock format: %s", lock)
		}
		targetId, ok := lockMap["ID"]
		if !ok {
			return fmt.Errorf("lock is missing ID")
		}
		if curId != targetId {
			return fmt.Errorf("have lock id %s, want lock id %s", curId, targetId)
		}
	case ok && reqType == reqTypeLock:
		return fmt.Errorf("trying to lock a locked state: %s", lock)
	case !ok && reqType == reqTypeUnlock:
		return fmt.Errorf("trying to unlock an unlocked state")
	default:
		return nil
	}
	return nil
}

func (self *TfBackend) updateState(r *http.Request, state *api.KVSecret, body map[string]any) error {
	self.ctx.Logger.Printf("updating state %s", self.statePath)
	if _, err := self.vault.KVv2(self.stateMount).Put(
		r.Context(),
		self.statePath,
		map[string]any{stateKey: body},
		api.WithCheckAndSet(state.VersionMetadata.Version),
	); err != nil {
		return fmt.Errorf("could not update state secret: %w", err)
	}
	return nil
}

func (self *TfBackend) writeState(_ *http.Request, w http.ResponseWriter, stateObj *api.KVSecret) error {
	self.ctx.Logger.Printf("returning state %s", self.statePath)
	state, ok := stateObj.Data[stateKey]
	if !ok {
		return nil
	}
	stateJson, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("could not marshal state: %w", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(stateJson)
	return nil
}

func (self *TfBackend) unlockState(r *http.Request, lock *api.KVSecret) error {
	self.ctx.Logger.Printf("unlocking state %s using lock %s", self.statePath, self.lockPath)
	_, ok := lock.Data[lockKey]
	if !ok {
		return fmt.Errorf("missing lock info")
	}
	if _, err := self.vault.KVv2(self.lockMount).Put(
		r.Context(),
		self.lockPath,
		map[string]any{},
		api.WithCheckAndSet(lock.VersionMetadata.Version),
	); err != nil {
		return fmt.Errorf("could not update lock: %w", err)
	}
	return nil
}

func (self *TfBackend) lockState(r *http.Request, lock *api.KVSecret, body map[string]any) error {
	self.ctx.Logger.Printf("locking state %s using lock %s", self.statePath, self.lockPath)
	if _, err := self.vault.KVv2(self.lockMount).Put(
		r.Context(),
		self.lockPath,
		map[string]any{lockKey: body},
		api.WithCheckAndSet(lock.VersionMetadata.Version),
	); err != nil {
		return fmt.Errorf("could not update lock secret: %w", err)
	}
	return nil
}
