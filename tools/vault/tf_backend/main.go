package main

import (
	"context"
	"encoding/json"
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
	"git.alwaldend.com/alwaldend/src/tools/vault/tf_backend/tf_backend_proto"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
)

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

type pluginState struct {
	ctx        context.Context
	cfg        *tf_backend_proto.Config
	listener   net.Listener
	req        *al_proto.PluginStartRequest
	statePath  string
	stateMount string
	lockPath   string
	lockMount  string
	vault      *api.Client
	username   string
	password   string
}

func (self *pluginState) left(err error) fp.Either[*pluginState] {
	return fp.Left[*pluginState](err)
}

func (self *pluginState) right() fp.Either[*pluginState] {
	return fp.Right(self)
}

type pluginMonad = fp.Monad[*pluginState]

type reqState struct {
	w       http.ResponseWriter
	r       *http.Request
	ctx     context.Context
	ps      *pluginState
	state   *api.KVSecret
	lock    *api.KVSecret
	reqType requestType
	reqBody map[string]any
}

func (self *reqState) left(err error) fp.Either[*reqState] {
	return fp.Left[*reqState](err)
}

func (self *reqState) right() fp.Either[*reqState] {
	return fp.Right(self)
}

type requestMonad = fp.Monad[*reqState]
type requestFunc = func(*reqState) fp.Monad[*reqState]

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
	if username != rs.ps.username || password != rs.ps.password {
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

func fetchSecret(st secretType) requestFunc {
	return func(rs *reqState) requestMonad {
		var mount, path string
		switch st {
		case secretTypeLock:
			mount, path = rs.ps.lockMount, rs.ps.lockPath
		case secretTypeState:
			mount, path = rs.ps.stateMount, rs.ps.statePath
		default:
			return rs.left(fmt.Errorf("invalid secret type: %s", st))
		}
		kv := rs.ps.vault.KVv2(mount)
		list, err := rs.ps.vault.Logical().List(fmt.Sprintf("%s/metadata/%s", mount, filepath.Dir(path)))
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
	logger.Printf("Updating state %s", rs.ps.statePath)
	_, err := rs.ps.vault.KVv2(rs.ps.stateMount).Put(
		rs.ctx,
		rs.ps.statePath,
		map[string]any{stateKey: rs.reqBody},
		api.WithCheckAndSet(rs.state.VersionMetadata.Version),
	)
	if err != nil {
		return rs.left(fmt.Errorf("could not update state secret: %w", err))
	}
	return rs.right()
}

func writeState(rs *reqState) requestMonad {
	logger.Printf("Returning state %s", rs.ps.statePath)
	state, ok := rs.state.Data[stateKey]
	if !ok {
		return rs.right()
	}
	stateJson, err := json.Marshal(state)
	if err != nil {
		return rs.left(fmt.Errorf("could not marshal the state: %w", err))
	}
	rs.w.WriteHeader(http.StatusOK)
	rs.w.Write(stateJson)
	return rs.right()
}

func unlockState(rs *reqState) requestMonad {
	logger.Printf("Unlocking state %s using lock %s", rs.ps.statePath, rs.ps.lockPath)
	_, ok := rs.lock.Data[lockKey]
	if !ok {
		return rs.left(fmt.Errorf("missing lock info"))
	}
	_, err := rs.ps.vault.KVv2(rs.ps.lockMount).Put(
		rs.ctx,
		rs.ps.lockPath,
		map[string]any{},
		api.WithCheckAndSet(rs.lock.VersionMetadata.Version),
	)
	if err != nil {
		return rs.left(fmt.Errorf("could not update lock: %w", err))
	}
	return rs.right()
}

func lockState(rs *reqState) requestMonad {
	logger.Printf("Locking state %s using lock %s", rs.ps.statePath, rs.ps.lockPath)
	_, err := rs.ps.vault.KVv2(rs.ps.lockMount).Put(
		rs.ctx,
		rs.ps.lockPath,
		map[string]any{lockKey: rs.reqBody},
		api.WithCheckAndSet(rs.lock.VersionMetadata.Version),
	)
	if err != nil {
		return rs.left(fmt.Errorf("could not update lock secret: %w", err))
	}
	return rs.right()
}

func rootHandler(ps *pluginState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "panic: %s", r)
				logger.Printf("panic: %s\n%s", r, string(debug.Stack()))
			}
		}()
		_, err := fp.Get(handleRequest(ps, w, r))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %s", err)
			logger.Printf("error: %s", err)
		}
	}
}

func handleRequest(ps *pluginState, w http.ResponseWriter, r *http.Request) requestMonad {
	return fp.Pipe7(
		fp.Just(&reqState{w: w, r: r, ctx: r.Context(), ps: ps}),
		checkAuth,
		setType,
		readBody,
		fetchSecret(secretTypeLock),
		checkLock,
		doAction,
	)
}

func doAction(rs *reqState) requestMonad {
	switch rs.reqType {
	case reqTypeGet:
		return fp.Pipe(
			rs.right(),
			fetchSecret(secretTypeState),
			writeState,
		)
	case reqTypeUpdate:
		return fp.Pipe(
			rs.right(),
			fetchSecret(secretTypeState),
			updateState,
		)
	case reqTypeUnlock:
		return fp.Pipe(
			rs.right(),
			unlockState,
		)
	case reqTypeLock:
		return fp.Pipe(
			rs.right(),
			lockState,
		)
	default:
		return rs.left(fmt.Errorf("invalid request type"))
	}
}

func parseConfig(ps *pluginState) pluginMonad {
	_, err := fp.Get(al_plugin.ParseConfig(ps.req.Plugin, ps.cfg))
	if err != nil {
		return ps.left(fmt.Errorf("could not parse plugin config: %w", err))
	}
	ps.statePath = fmt.Sprintf("%s/state", ps.cfg.VaultSecret)
	ps.stateMount = ps.cfg.VaultSecretMount
	ps.lockPath = fmt.Sprintf("%s/lock", ps.cfg.VaultSecret)
	ps.lockMount = ps.cfg.VaultSecretMount
	return ps.right()
}

func createVault(ps *pluginState) pluginMonad {
	return fp.Pipe4(
		fp.Right(ps.req.Config),
		al.NewVault,
		func(v *al.Vault) fp.Monad[*api.Client] {
			return v.Client(ps.ctx, ps.cfg.VaultConn, ps.cfg.VaultAuth)
		},
		func(v *api.Client) pluginMonad {
			ps.vault = v
			return fp.Right(ps)
		},
	)
}

func startServer(ps *pluginState) pluginMonad {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return ps.left(fmt.Errorf("could not listen on a port: %w", err))
	}
	ps.listener = listener
	handler := http.NewServeMux()
	handler.HandleFunc("/", rootHandler(ps))
	server := &http.Server{Handler: handler}
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
	return ps.right()
}

func createResponse(ps *pluginState) fp.Monad[*al_proto.PluginStartResponse] {
	return fp.Right(&al_proto.PluginStartResponse{
		Env: map[string]string{
			"TF_HTTP_USERNAME":       ps.username,
			"TF_HTTP_PASSWORD":       ps.password,
			"TF_HTTP_ADDRESS":        fmt.Sprintf("http://%s/state", ps.listener.Addr().String()),
			"TF_HTTP_UPDATE_METHOD":  http.MethodPost,
			"TF_HTTP_LOCK_ADDRESS":   fmt.Sprintf("http://%s/lock", ps.listener.Addr().String()),
			"TF_HTTP_LOCK_METHOD":    http.MethodPost,
			"TF_HTTP_UNLOCK_ADDRESS": fmt.Sprintf("http://%s/unlock", ps.listener.Addr().String()),
			"TF_HTTP_UNLOCK_METHOD":  http.MethodPost,
		},
	})
}

func (self Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	return fp.Get(fp.Pipe5(
		fp.Just(&pluginState{
			ctx:      ctx,
			cfg:      &tf_backend_proto.Config{},
			req:      req,
			username: uuid.New().String(),
			password: uuid.New().String(),
		}),
		parseConfig,
		createVault,
		startServer,
		createResponse,
	))
}

func main() {
	fp.GetOrExit(al_plugin.Serve(context.Background(), os.Stdin, os.Stdout, Plugin{}), 1)
}
