package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/tf_backend/tf_backend_proto"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
)

type (
	// requestType identifies a Terraform HTTP backend operation.
	requestType string
	// secretType selects the versioned state or lock record.
	secretType string
)

const (
	// maxRequestBytes bounds state and lock payloads buffered for JSON decoding.
	maxRequestBytes = 32 << 20
	// reqTypeUpdate replaces state under the current lock.
	reqTypeUpdate = requestType("update")
	// reqTypeGet reads the current state.
	reqTypeGet = requestType("get")
	// reqTypeLock acquires an unlocked state record.
	reqTypeLock = requestType("lock")
	// reqTypeUnlock releases an existing state lock.
	reqTypeUnlock = requestType("unlock")
	// secretTypeState selects the Vault record containing Terraform state.
	secretTypeState = secretType("state")
	// secretTypeLock selects the Vault record containing Terraform lock information.
	secretTypeLock = secretType("lock")
	// stateKey names the state payload within a versioned KV document.
	stateKey = "state"
	// lockKey names the lock payload within a versioned KV document.
	lockKey = "lock"
)

// backend exposes one Vault-backed state and lock pair to Terraform.
type backend struct {
	// username identifies the generated HTTP Basic authentication principal.
	username string
	// password authenticates the Terraform process using this invocation.
	password string
	// server owns the loopback endpoint used by Terraform.
	server *httpBackendServer
	// statePath locates the state record within the KV mount.
	statePath string
	// stateMount identifies the KV engine holding state.
	stateMount string
	// lockPath locates the lock record within the KV mount.
	lockPath string
	// lockMount identifies the KV engine holding the lock.
	lockMount string
	// vault borrows the authenticated client from vaultStore.
	vault *api.Client
	// vaultStore owns and releases invocation-scoped Vault credentials.
	vaultStore *al.VaultStore
	// ctx provides the invocation logger.
	ctx *al.CmdCtx
}

// newBackend acquires a Vault client and prepares an invocation-specific loopback endpoint.
func newBackend(requestCtx context.Context, ctx *al.CmdCtx, config *al_proto.Config, settings *tf_backend_proto.Config) (*backend, error) {
	if settings.VaultSecret == "" || settings.VaultSecretMount == "" {
		return nil, fmt.Errorf("missing secret config")
	}
	res := &backend{
		username: uuid.NewString(),
		password: uuid.NewString(),
		ctx:      ctx,
	}
	res.statePath = fmt.Sprintf("%s/state", settings.VaultSecret)
	res.stateMount = settings.VaultSecretMount
	res.lockPath = fmt.Sprintf("%s/lock", settings.VaultSecret)
	res.lockMount = settings.VaultSecretMount
	vault := al.NewVault(config)
	client, err := vault.Client(requestCtx, settings.VaultConn, settings.VaultAuth).Get()
	if err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return nil, errors.Join(fmt.Errorf("could not create vault client: %w", al.SanitizeVaultError(err)), vault.Stop(cleanupCtx))
	}
	res.vaultStore = vault
	res.vault = client.Client
	res.server = newServer(res.reqHandler)
	return res, nil
}

// start opens the backend listener after authentication succeeds.
func (b *backend) start(ctx context.Context) error {
	if err := b.server.start(ctx); err != nil {
		return fmt.Errorf("start Terraform HTTP listener: %w", err)
	}
	return nil
}

// stop drains the listener before releasing invocation-owned credentials.
func (b *backend) stop(ctx context.Context) error {
	serverErr := b.server.stop(ctx)
	if err := b.vaultStore.Stop(ctx); err != nil {
		serverErr = errors.Join(serverErr, fmt.Errorf("release backend Vault credentials: %w", err))
	}
	if serverErr != nil {
		return fmt.Errorf("stop Terraform backend: %w", serverErr)
	}
	return nil
}

// environment constructs Terraform HTTP backend settings for the running listener.
func (b *backend) environment() (map[string]string, error) {
	addr, err := b.server.address()
	if err != nil {
		return nil, fmt.Errorf("could not get server address: %w", err)
	}
	return map[string]string{
		"TF_HTTP_USERNAME":       b.username,
		"TF_HTTP_PASSWORD":       b.password,
		"TF_HTTP_ADDRESS":        fmt.Sprintf("http://%s/state", addr),
		"TF_HTTP_UPDATE_METHOD":  http.MethodPost,
		"TF_HTTP_LOCK_ADDRESS":   fmt.Sprintf("http://%s/lock", addr),
		"TF_HTTP_LOCK_METHOD":    http.MethodPost,
		"TF_HTTP_UNLOCK_ADDRESS": fmt.Sprintf("http://%s/unlock", addr),
		"TF_HTTP_UNLOCK_METHOD":  http.MethodPost,
	}, nil
}

// reqHandler reports request failures and writes state only after successful retrieval.
func (b *backend) reqHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	state, err := b.handleRequest(r)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "Terraform request exceeds the body limit", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Terraform backend request failed", http.StatusInternalServerError)
		b.ctx.Logger.Printf("Terraform request failed: %s", al.SanitizeVaultError(err))
		return
	}
	if state != nil {
		if _, err := w.Write(state); err != nil {
			b.ctx.Logger.Printf("write Terraform state response: %v", err)
		}
	}
}

// handleRequest authenticates, checks lock ownership, and runs the requested operation.
func (b *backend) handleRequest(r *http.Request) ([]byte, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, fmt.Errorf("missing basic auth")
	}
	if username != b.username || password != b.password {
		return nil, fmt.Errorf("invalid auth")
	}
	reqType, err := parseRequestType(r)
	if err != nil {
		return nil, fmt.Errorf("could not determine request type: %w", err)
	}
	body := make(map[string]any)
	payload := bufio.NewReader(r.Body)
	_, err = payload.Peek(1)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("read Terraform request: %w", err)
	}
	if err == nil {
		decoder := json.NewDecoder(payload)
		decoder.UseNumber()
		if err := decoder.Decode(&body); err != nil {
			return nil, fmt.Errorf("decode Terraform request: %w", err)
		}
		if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
			if err != nil {
				return nil, fmt.Errorf("read end of Terraform request: %w", err)
			}
			return nil, fmt.Errorf("unexpected data after Terraform request")
		}
	}
	lock, err := b.fetchSecret(r, secretTypeLock)
	if err != nil {
		return nil, fmt.Errorf("could not fetch lock: %w", err)
	}
	if err := checkLock(r, reqType, lock); err != nil {
		return nil, fmt.Errorf("invalid lock: %w", err)
	}
	switch reqType {
	case reqTypeGet:
		state, err := b.fetchSecret(r, secretTypeState)
		if err != nil {
			return nil, fmt.Errorf("could not fetch state secret: %w", err)
		}
		data, err := b.encodeState(state)
		if err != nil {
			return nil, fmt.Errorf("encode Terraform state: %w", err)
		}
		return data, nil
	case reqTypeUpdate:
		state, err := b.fetchSecret(r, secretTypeState)
		if err != nil {
			return nil, fmt.Errorf("could not fetch state secret: %w", err)
		}
		if err := b.updateState(r, state, body); err != nil {
			return nil, fmt.Errorf("could not update state: %w", err)
		}
	case reqTypeLock:
		if err := b.lockState(r, lock, body); err != nil {
			return nil, fmt.Errorf("could not lock state: %w", err)
		}
	case reqTypeUnlock:
		if err := b.unlockState(r, lock); err != nil {
			return nil, fmt.Errorf("could not unlock state: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid request type")
	}
	return nil, nil
}

// parseRequestType maps the Terraform HTTP paths and methods to backend operations.
func parseRequestType(r *http.Request) (requestType, error) {
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

// fetchSecret reads a KV record and initializes missing data using a create-only write.
func (b *backend) fetchSecret(r *http.Request, st secretType) (*api.KVSecret, error) {
	var mount, path string
	switch st {
	case secretTypeLock:
		mount, path = b.lockMount, b.lockPath
	case secretTypeState:
		mount, path = b.stateMount, b.statePath
	default:
		return nil, fmt.Errorf("invalid secret type: %s", st)
	}
	kv := b.vault.KVv2(mount)
	state, err := kv.Get(r.Context(), path)
	if err == nil {
		return state, nil
	}
	if !errors.Is(err, api.ErrSecretNotFound) {
		return nil, fmt.Errorf("read %s secret: %w", st, al.SanitizeVaultError(err))
	}
	_, createErr := kv.Put(r.Context(), path, map[string]any{}, api.WithCheckAndSet(0))
	if createErr != nil {
		var response *api.ResponseError
		if !errors.As(createErr, &response) || response.StatusCode != http.StatusBadRequest {
			return nil, fmt.Errorf("initialize %s secret: %w", st, al.SanitizeVaultError(createErr))
		}
		// KV returns HTTP 400 if a concurrent caller won the creation race.
		// Read once to identify its record without repeating the write.
	}
	state, err = kv.Get(r.Context(), path)
	if err != nil {
		return nil, fmt.Errorf("read initialized %s secret: %w", st, errors.Join(al.SanitizeVaultError(createErr), al.SanitizeVaultError(err)))
	}
	return state, nil
}

// checkLock validates the current lock against the requested state operation.
func checkLock(r *http.Request, reqType requestType, lockObj *api.KVSecret) error {
	lock, ok := lockObj.Data[lockKey]
	currentID := r.URL.Query().Get("ID")
	switch {
	case reqType == reqTypeUpdate:
		if !ok {
			return fmt.Errorf("missing lock field %s", lockKey)
		}
		lockMap, ok := lock.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid lock format")
		}
		targetID, ok := lockMap["ID"]
		if !ok {
			return fmt.Errorf("lock is missing ID")
		}
		if currentID != targetID {
			return fmt.Errorf("lock ID does not match")
		}
	case ok && reqType == reqTypeLock:
		return fmt.Errorf("trying to lock a locked state")
	case !ok && reqType == reqTypeUnlock:
		return fmt.Errorf("trying to unlock an unlocked state")
	default:
		return nil
	}
	return nil
}

// updateState uses the observed KV version to prevent a concurrent state overwrite.
func (b *backend) updateState(r *http.Request, state *api.KVSecret, body map[string]any) error {
	b.ctx.Logger.Printf("updating state %s", b.statePath)
	if _, err := b.vault.KVv2(b.stateMount).Put(
		r.Context(),
		b.statePath,
		map[string]any{stateKey: body},
		api.WithCheckAndSet(state.VersionMetadata.Version),
	); err != nil {
		return fmt.Errorf("could not update state secret: %w", al.SanitizeVaultError(err))
	}
	return nil
}

// encodeState serializes the state field, returning an empty response for a fresh record.
func (b *backend) encodeState(stateObj *api.KVSecret) ([]byte, error) {
	b.ctx.Logger.Printf("returning state %s", b.statePath)
	state, ok := stateObj.Data[stateKey]
	if !ok {
		return nil, nil
	}
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("could not marshal state: %w", al.SanitizeVaultError(err))
	}
	return stateJSON, nil
}

// unlockState clears the lock only at the observed KV version.
func (b *backend) unlockState(r *http.Request, lock *api.KVSecret) error {
	b.ctx.Logger.Printf("unlocking state %s using lock %s", b.statePath, b.lockPath)
	_, ok := lock.Data[lockKey]
	if !ok {
		return fmt.Errorf("missing lock info")
	}
	if _, err := b.vault.KVv2(b.lockMount).Put(
		r.Context(),
		b.lockPath,
		map[string]any{},
		api.WithCheckAndSet(lock.VersionMetadata.Version),
	); err != nil {
		return fmt.Errorf("could not update lock: %w", al.SanitizeVaultError(err))
	}
	return nil
}

// lockState writes Terraform lock metadata at the observed KV version.
func (b *backend) lockState(r *http.Request, lock *api.KVSecret, body map[string]any) error {
	b.ctx.Logger.Printf("locking state %s using lock %s", b.statePath, b.lockPath)
	if _, err := b.vault.KVv2(b.lockMount).Put(
		r.Context(),
		b.lockPath,
		map[string]any{lockKey: body},
		api.WithCheckAndSet(lock.VersionMetadata.Version),
	); err != nil {
		return fmt.Errorf("could not update lock secret: %w", al.SanitizeVaultError(err))
	}
	return nil
}
