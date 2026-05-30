package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
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

type reqState struct {
	w          http.ResponseWriter
	r          *http.Request
	ctx        context.Context
	vault      *api.Client
	state      *api.KVSecret
	statePath  string
	stateMount string
	lock       *api.KVSecret
	lockPath   string
	lockMount  string
	username   string
	password   string
	reqType    requestType
	reqBody    map[string]any
	respBody   []byte
}

func (self *reqState) left(err error) fp.Either[*reqState] {
	return fp.Left[*reqState](err)
}

func (self *reqState) right() fp.Either[*reqState] {
	return fp.Right(self)
}

type requestMonad = fp.Monad[*reqState]

func setType(rs *reqState) requestMonad {
	path := rs.r.URL.Path
	method := rs.r.Method
	if path == "/lock" && method == http.MethodPost {
		rs.reqType = reqTypeLock
	} else if path == "/unlock" && method == http.MethodPost {
		rs.reqType = reqTypeUnlock
	} else if path == "/state" && method == http.MethodPost {
		rs.reqType = reqTypeUpdate
	} else if path == "/state" && method == http.MethodGet {
		rs.reqType = reqTypeGet
	} else {
		return rs.left(fmt.Errorf("invalid request path or method"))
	}
	return rs.right()
}

func checkAuth(rs *reqState) requestMonad {
	username, password, ok := rs.r.BasicAuth()
	if !ok {
		return rs.left(fmt.Errorf("missng basic auth"))
	}
	if username != rs.username || password != rs.password {
		return rs.left(fmt.Errorf("invalid auth"))
	}
	return rs.right()
}

func readBody(rs *reqState) requestMonad {
	bodyRaw, err := io.ReadAll(rs.r.Body)
	if err != nil {
		return rs.left(fmt.Errorf("could not read request body: %w", err))
	}
	body := make(map[string]any)
	rs.reqBody = body
	if len(bodyRaw) != 0 {
		err = json.Unmarshal(bodyRaw, &body)
		if err != nil {
			return rs.left(fmt.Errorf("could not unmarshal request body: %w", err))
		}
	}
	return rs.right()
}

func fetchSecret(rs *reqState, st secretType) requestMonad {
	var mount, path string
	switch st {
	case secretTypeLock:
		mount, path = rs.lockMount, rs.lockPath
	case secretTypeState:
		mount, path = rs.stateMount, rs.statePath
	default:
		return rs.left(fmt.Errorf("invalid secret type: %s", st))
	}
	kv := rs.vault.KVv2(mount)
	list, err := rs.vault.Logical().List(fmt.Sprintf("%s/metadata/%s", mount, filepath.Dir(path)))
	if err != nil {
		return rs.left(fmt.Errorf("could not list secrets: %w", err))
	}
	create := true
	if list != nil {
		keys, ok := list.Data["keys"]
		if ok {
			keysList, ok := keys.([]any)
			if !ok {
				return rs.left(fmt.Errorf("invalid keys field for some reason: %s", keys))
			}
			for _, keyAny := range keysList {
				key, ok := keyAny.(string)
				if !ok {
					return rs.left(fmt.Errorf("invalid key type for some reason: %s", key))
				}
				if key == filepath.Base(path) {
					create = false
					break
				}
			}
		}
	}
	if create {
		_, err := kv.Put(rs.ctx, path, map[string]any{})
		if err != nil {
			return rs.left(fmt.Errorf("could not create secret: %w", err))
		}
	}
	secret, err := kv.Get(rs.ctx, path)
	if err != nil {
		return rs.left(fmt.Errorf("could not fetch secret: %w", err))
	}
	switch st {
	case secretTypeLock:
		rs.lock = secret
	case secretTypeState:
		rs.state = secret
	default:
		return rs.left(fmt.Errorf("invalid secret type: %s", st))
	}
	return rs.right()
}

func checkLock(rs *reqState) requestMonad {
	lock, ok := rs.lock.Data[lockKey]
	curId := rs.r.URL.Query().Get("ID")
	if rs.reqType == reqTypeUpdate {
		if !ok {
			return rs.left(fmt.Errorf("missing lock field %s", lockKey))
		}
		lockMap, ok := lock.(map[string]any)
		if !ok {
			return rs.left(fmt.Errorf("invalid lock format: %s", lock))
		}
		targetId, ok := lockMap["ID"]
		if !ok {
			return rs.left(fmt.Errorf("lock is missing ID"))
		}
		if curId != targetId {
			return rs.left(fmt.Errorf("have lock id %s, want lock id %s", curId, targetId))
		}
	} else if ok && rs.reqType == reqTypeLock {
		return rs.left(fmt.Errorf("trying to lock a locked state: %s", lock))
	} else if !ok && rs.reqType == reqTypeUnlock {
		return rs.left(fmt.Errorf("trying to unlock an unlocked state"))
	}
	return rs.right()
}

func updateState(rs *reqState) requestMonad {
	logger.Printf("Updating state %s", rs.statePath)
	_, err := rs.vault.KVv2(rs.stateMount).Put(
		rs.ctx,
		rs.statePath,
		map[string]any{stateKey: rs.reqBody},
		api.WithCheckAndSet(rs.state.VersionMetadata.Version),
	)
	if err != nil {
		return rs.left(fmt.Errorf("could not update state secret: %w", err))
	}
	return rs.right()
}

func writeState(rs *reqState) requestMonad {
	logger.Printf("Returning state %s", rs.statePath)
	state, ok := rs.state.Data[stateKey]
	if !ok {
		return rs.right()
	}
	stateJson, err := json.Marshal(state)
	if err != nil {
		return rs.left(fmt.Errorf("could not marshal the state: %w", err))
	}
	rs.respBody = stateJson
	return rs.right()
}

func unlockState(rs *reqState) requestMonad {
	logger.Printf("Unlocking state %s using lock %s", rs.statePath, rs.lockPath)
	_, ok := rs.lock.Data[lockKey]
	if !ok {
		return rs.left(fmt.Errorf("missing lock info"))
	}
	_, err := rs.vault.KVv2(rs.lockMount).Put(
		rs.ctx,
		rs.lockPath,
		map[string]any{},
		api.WithCheckAndSet(rs.lock.VersionMetadata.Version),
	)
	if err != nil {
		return rs.left(fmt.Errorf("could not update lock: %w", err))
	}
	return rs.right()
}

func lockState(rs *reqState) requestMonad {
	logger.Printf("Locking state %s using lock %s", rs.statePath, rs.lockPath)
	_, err := rs.vault.KVv2(rs.lockMount).Put(
		rs.ctx,
		rs.lockPath,
		map[string]any{lockKey: rs.reqBody},
		api.WithCheckAndSet(rs.lock.VersionMetadata.Version),
	)
	if err != nil {
		return rs.left(fmt.Errorf("could not update lock secret: %w", err))
	}
	return rs.right()
}

func rootHandler(vault *api.Client, username string, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "panic: %s", r)
				logger.Printf("panic: %s\n%s", r, string(debug.Stack()))
			}
		}()
		res, err := fp.Get(fp.Right(&reqState{
			w:          w,
			r:          r,
			ctx:        r.Context(),
			vault:      vault,
			username:   username,
			password:   password,
			statePath:  fmt.Sprintf("%s/state", *vaultSecret),
			stateMount: *vaultSecretMount,
			lockPath:   fmt.Sprintf("%s/lock", *vaultSecret),
			lockMount:  *vaultSecretMount,
		}).
			FlatMap(checkAuth).
			FlatMap(setType).
			FlatMap(readBody).
			FlatMap(func(rs *reqState) requestMonad { return fetchSecret(rs, secretTypeLock) }).
			FlatMap(checkLock).
			FlatMap(func(rs *reqState) requestMonad {
				switch rs.reqType {
				case reqTypeGet:
					return fetchSecret(rs, secretTypeState).FlatMap(writeState)
				case reqTypeUpdate:
					return fetchSecret(rs, secretTypeState).FlatMap(updateState)
				case reqTypeUnlock:
					return unlockState(rs)
				case reqTypeLock:
					return lockState(rs)
				default:
					return rs.left(fmt.Errorf("invalid request type"))
				}
			}))
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(res.respBody)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %s", err)
			logger.Printf("error: %s", err)
		}
	}
}

func (self Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (res *al_proto.PluginStartResponse, err error) {
	logger.Printf("Creating a vault token")
	vault, err := al.NewVault(req.Config)
	if err != nil {
		return nil, fmt.Errorf("could not create vault client: %w", err)
	}
	vault_client, err := vault.Client(ctx, *vaultConn, *vaultAuth)
	if err != nil {
		return nil, fmt.Errorf("could not authorize vault client: %w", err)
	}
	logger.Println("Created a vault token")
	logger.Printf("Creating a server")
	username, password := uuid.New().String(), uuid.New().String()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("could not listen on a port: %w", err)
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/", rootHandler(vault_client, username, password))
	server := &http.Server{Handler: handler}
	logger.Printf("Created a server")
	go func() {
		curCtx, cancel1 := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel1()
		<-curCtx.Done()
		curCtx, cancel2 := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
		defer cancel2()
		err := server.Shutdown(curCtx)
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
	res = &al_proto.PluginStartResponse{
		Env: map[string]string{
			"TF_HTTP_USERNAME":       username,
			"TF_HTTP_PASSWORD":       password,
			"TF_HTTP_ADDRESS":        fmt.Sprintf("http://%s/state", listener.Addr().String()),
			"TF_HTTP_UPDATE_METHOD":  http.MethodPost,
			"TF_HTTP_LOCK_ADDRESS":   fmt.Sprintf("http://%s/lock", listener.Addr().String()),
			"TF_HTTP_LOCK_METHOD":    http.MethodPost,
			"TF_HTTP_UNLOCK_ADDRESS": fmt.Sprintf("http://%s/unlock", listener.Addr().String()),
			"TF_HTTP_UNLOCK_METHOD":  http.MethodPost,
		},
	}
	return res, nil
}

func main() {
	logger.Printf("Starting the plugin")
	flag.Parse()
	err := al_plugin.Serve(context.Background(), os.Stdin, os.Stdout, Plugin{})
	if err != nil {
		logger.Printf("failed to run plugin: %s", err.Error())
		os.Exit(1)
	}
}
