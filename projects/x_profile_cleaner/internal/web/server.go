package web

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/app"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
)

// Config supplies command-owned resources and the initial editor contents.
type Config struct {
	// Address is an explicit loopback listen address; port zero is permitted.
	Address string
	// Policy initializes the editor without evaluating it.
	Policy string
	// Output receives the private access link once, separately from progress logs.
	Output io.Writer
	// Logger shares application progress between stderr and Logs.
	Logger *slog.Logger
	// Logs retains bounded records for authenticated browser polling.
	Logs *Log
}

// status describes the current operation without repeatedly transferring a plan.
type status struct {
	// Phase is idle, planning, ready, applying, completed, failed, or cancelled.
	Phase string `json:"phase"`
	// JobID identifies the operation for cancellation or explicit resumption.
	JobID string `json:"job_id"`
	// PlanID identifies the current review or execution plan, if one exists.
	PlanID string `json:"plan_id"`
	// Actions counts planned or submitted actions in the current operation.
	Actions int `json:"actions"`
	// Completed counts successful mutations and targets confirmed absent.
	Completed int `json:"completed"`
	// Remaining counts original selected actions without confirmed completion.
	Remaining int `json:"remaining"`
	// CompletedIndexes identifies completed actions in the retained execution plan.
	CompletedIndexes []int `json:"completed_indexes,omitempty"`
	// Error describes a failed or cancelled operation without granting authority.
	Error string `json:"error,omitempty"`
}

// execution retains the fixed submitted subset and its progress for this server session.
type execution struct {
	// plan contains only the actions selected and confirmed before initial execution.
	plan model.Plan
	// policy is the original source shown with the selected plan, never reevaluated here.
	policy string
	// id identifies this execution plan across resumed jobs.
	id string
	// completed records successful or already-absent targets by original plan index.
	completed map[int]bool
}

// server owns a single ephemeral account session and its immutable current plan.
type server struct {
	// ctx ends all jobs when serve stops.
	ctx context.Context
	// application supplies existing planning and exact execution behavior.
	application *app.App
	// config holds the command-owned logger and editor setup.
	config Config
	// host is the exact listener authority accepted by the HTTP boundary.
	host string
	// token is a fresh per-run capability, absent from request URLs and logs.
	token string
	// mu protects session state and job admission.
	mu sync.Mutex
	// state is copied for browser status requests.
	state status
	// policy is the latest submitted editor source.
	policy string
	// plan is retained only until replacement or a validated apply starts.
	plan *model.Plan
	// execution survives a failed or cancelled apply until a new plan replaces it.
	execution *execution
	// sequence identifies each admitted operation.
	sequence uint64
	// cancel interrupts the active job; nil means no job is running.
	cancel context.CancelFunc
	// jobs prevents server shutdown from abandoning mutation goroutines.
	jobs sync.WaitGroup
}

// Serve binds a private loopback interface and cancels its jobs on command shutdown.
func Serve(parent context.Context, application *app.App, config Config) error {
	host, port, err := net.SplitHostPort(config.Address)
	if err != nil {
		return fmt.Errorf("parse listen address: %w", err)
	}
	if host == "localhost" {
		host = "127.0.0.1"
	}
	if ip := net.ParseIP(host); ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("listen address must use a loopback IP or localhost")
	}
	if len(config.Policy) > 1<<20 {
		return fmt.Errorf("browser policy exceeds 1 MiB")
	}
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return fmt.Errorf("create browser session: %w", err)
	}
	listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("listen for browser requests: %w", err)
	}
	defer listener.Close()
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	s := &server{ctx: ctx, application: application, config: config, host: listener.Addr().String(), token: hex.EncodeToString(secret[:]), policy: config.Policy, state: status{Phase: "idle"}}
	if _, err := fmt.Fprintf(config.Output, "Open http://%s/#%s\n", s.host, s.token); err != nil {
		return fmt.Errorf("write private browser access link: %w", err)
	}
	config.Logger.InfoContext(ctx, "Serving profile cleaner", "address", s.host)
	httpServer := &http.Server{Handler: http.HandlerFunc(s.handle), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 8192}
	stopped := make(chan error, 1)
	go func() {
		<-ctx.Done()
		shutdown, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		stopped <- httpServer.Shutdown(shutdown)
	}()
	serveErr := httpServer.Serve(listener)
	cancel()
	shutdownErr := <-stopped
	s.jobs.Wait()
	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		return fmt.Errorf("serve browser interface: %w", serveErr)
	}
	if shutdownErr != nil {
		return fmt.Errorf("stop browser interface: %w", shutdownErr)
	}
	config.Logger.Info("Browser interface stopped")
	return nil
}

// handle enforces local authority, capability authentication, and explicit routes.
func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; img-src 'self' data:; base-uri 'none'; frame-ancestors 'none'; form-action 'none'")
	if r.Host != s.host || (r.Header.Get("Origin") != "" && r.Header.Get("Origin") != "http://"+s.host) {
		s.fail(w, http.StatusForbidden, fmt.Errorf("foreign host or origin"))
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/") {
		if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Cleaner-Token")), []byte(s.token)) != 1 {
			s.fail(w, http.StatusForbidden, fmt.Errorf("open the private access link printed by serve"))
			return
		}
		s.handleAPI(w, r)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/")
	if name == "" {
		name = "index.html"
	}
	types := map[string]string{"index.html": "text/html; charset=utf-8", "app.js": "text/javascript; charset=utf-8", "app.css": "text/css; charset=utf-8", "theme.css": "text/css; charset=utf-8"}
	contentType, ok := types[name]
	if !ok || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	content, err := assets.ReadFile(name)
	if err != nil {
		s.fail(w, http.StatusInternalServerError, fmt.Errorf("read interface asset: %w", err))
		return
	}
	w.Header().Set("Content-Type", contentType)
	if _, err := w.Write(content); err != nil {
		s.config.Logger.Debug("Browser asset delivery ended", "error", err)
	}
}

// handleAPI maps authenticated requests to state reads or serialized operations.
func (s *server) handleAPI(w http.ResponseWriter, r *http.Request) {
	switch r.Method + " " + r.URL.Path {
	case "GET /api/session":
		s.mu.Lock()
		policy := s.policy
		s.mu.Unlock()
		s.respond(w, http.StatusOK, map[string]string{"policy": policy})
	case "GET /api/state":
		var after uint64
		if value := r.URL.Query().Get("after"); value != "" {
			var err error
			after, err = strconv.ParseUint(value, 10, 64)
			if err != nil {
				s.fail(w, http.StatusBadRequest, fmt.Errorf("parse log cursor: %w", err))
				return
			}
		}
		s.mu.Lock()
		state := s.state
		state.CompletedIndexes = append([]int(nil), s.state.CompletedIndexes...)
		s.mu.Unlock()
		entries, cursor := s.config.Logs.since(after)
		s.respond(w, http.StatusOK, struct {
			// status is the operation snapshot.
			status
			// Logs are new retained lines since the supplied cursor.
			Logs []logEntry `json:"logs"`
			// Cursor is the next polling position.
			Cursor uint64 `json:"cursor"`
		}{state, entries, cursor})
	case "GET /api/plan":
		s.mu.Lock()
		plan, id, source := s.plan, s.state.PlanID, s.policy
		if plan == nil && s.execution != nil {
			plan, id, source = &s.execution.plan, s.execution.id, s.execution.policy
		}
		s.mu.Unlock()
		if plan == nil {
			s.fail(w, http.StatusConflict, fmt.Errorf("no reviewable plan; run plan first"))
			return
		}
		s.respond(w, http.StatusOK, map[string]any{"plan_id": id, "plan": plan, "policy": source})
	case "POST /api/plan":
		var request struct {
			// Policy is editor text, never a server-side path.
			Policy *string `json:"policy"`
		}
		if !s.decode(w, r, &request) {
			return
		}
		if request.Policy == nil {
			s.fail(w, http.StatusBadRequest, fmt.Errorf("policy text is required"))
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.admit(w) {
			return
		}
		s.policy = *request.Policy
		s.execution = nil
		ctx := s.start("planning", 0)
		go s.prepare(ctx, *request.Policy)
		s.respond(w, http.StatusAccepted, s.state)
	case "POST /api/apply":
		var request struct {
			// PlanID identifies the immutable server-held plan being reviewed.
			PlanID string `json:"plan_id"`
			// Selected contains unique indexes into that plan's action list.
			Selected []*int `json:"selected"`
		}
		if !s.decode(w, r, &request) {
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.admit(w) {
			return
		}
		if s.plan == nil || request.PlanID != s.state.PlanID {
			s.fail(w, http.StatusConflict, fmt.Errorf("plan is absent, replaced, or already applied; run plan again"))
			return
		}
		selected := map[int]bool{}
		for _, value := range request.Selected {
			if value == nil {
				s.fail(w, http.StatusBadRequest, fmt.Errorf("selection indexes must be integers"))
				return
			}
			index := *value
			if index < 0 || index >= len(s.plan.Actions) || selected[index] {
				s.fail(w, http.StatusBadRequest, fmt.Errorf("selection contains an invalid or repeated action index"))
				return
			}
			selected[index] = true
		}
		if len(selected) == 0 {
			s.fail(w, http.StatusBadRequest, fmt.Errorf("select at least one planned action"))
			return
		}
		plan := *s.plan
		plan.Actions = []model.Action{}
		for index, action := range s.plan.Actions {
			if selected[index] {
				plan.Actions = append(plan.Actions, action)
			}
		}
		s.execution = &execution{plan: plan, policy: s.policy, completed: map[int]bool{}}
		s.startApply()
		s.respond(w, http.StatusAccepted, s.state)
	case "POST /api/resume":
		var request struct {
			// JobID rejects duplicate or stale resume requests from another tab.
			JobID string `json:"job_id"`
		}
		if !s.decode(w, r, &request) {
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.admit(w) {
			return
		}
		if request.JobID != s.state.JobID || s.execution == nil || s.state.Remaining == 0 || (s.state.Phase != "failed" && s.state.Phase != "cancelled") {
			s.fail(w, http.StatusConflict, fmt.Errorf("no interrupted execution for this job"))
			return
		}
		s.startApply()
		s.respond(w, http.StatusAccepted, s.state)
	case "POST /api/cancel":
		var request struct {
			// JobID prevents a delayed tab from cancelling a later operation.
			JobID string `json:"job_id"`
		}
		if !s.decode(w, r, &request) {
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.cancel == nil || request.JobID != s.state.JobID {
			s.fail(w, http.StatusConflict, fmt.Errorf("operation is no longer running"))
			return
		}
		s.cancel()
		s.respond(w, http.StatusAccepted, s.state)
	default:
		s.fail(w, http.StatusNotFound, fmt.Errorf("unknown API operation"))
	}
}

// admit checks job exclusion and shutdown while the caller holds the session lock.
func (s *server) admit(w http.ResponseWriter) bool {
	if s.ctx.Err() != nil || s.cancel != nil {
		s.fail(w, http.StatusConflict, fmt.Errorf("another operation is running or the server is stopping"))
		return false
	}
	return true
}

// start invalidates any old plan and registers one job under the session lock.
func (s *server) start(phase string, actions int) context.Context {
	ctx, cancel := context.WithCancel(s.ctx)
	s.cancel, s.plan = cancel, nil
	s.sequence++
	s.state = status{Phase: phase, JobID: strconv.FormatUint(s.sequence, 10), Actions: actions}
	s.jobs.Add(1)
	return ctx
}

// prepare retains a complete plan only when planning and cancellation both succeed.
func (s *server) prepare(ctx context.Context, source string) {
	defer s.jobs.Done()
	plan, err := s.application.PrepareSource(ctx, source, "browser-policy.lua")
	s.mu.Lock()
	defer s.mu.Unlock()
	if err == nil {
		err = ctx.Err()
	}
	if s.finish(ctx, err) {
		s.plan = &plan
		s.state.Phase, s.state.PlanID, s.state.Actions = "ready", s.state.JobID, len(plan.Actions)
	}
}

// startApply projects unfinished actions without changing the confirmed selection.
// The caller holds the session lock, and admission has excluded another active job.
func (s *server) startApply() {
	run := s.execution
	plan := run.plan
	plan.Actions = []model.Action{}
	indexes := []int{}
	completed := []int{}
	for index, action := range run.plan.Actions {
		if run.completed[index] {
			completed = append(completed, index)
		} else {
			indexes = append(indexes, index)
			plan.Actions = append(plan.Actions, action)
		}
	}
	ctx := s.start("applying", len(run.plan.Actions))
	if run.id == "" {
		run.id = s.state.JobID
	}
	s.state.PlanID = run.id
	s.state.Completed, s.state.Remaining = len(completed), len(indexes)
	s.state.CompletedIndexes = completed
	go s.apply(ctx, run, plan, indexes)
}

// apply records confirmed progress by original index and retains unfinished actions.
func (s *server) apply(ctx context.Context, run *execution, plan model.Plan, indexes []int) {
	defer s.jobs.Done()
	_, err := s.application.Apply(ctx, plan, func(index int) {
		s.mu.Lock()
		defer s.mu.Unlock()
		original := indexes[index]
		run.completed[original] = true
		s.state.Completed++
		s.state.Remaining--
		s.state.CompletedIndexes = append(s.state.CompletedIndexes, original)
	})
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.finish(ctx, err) {
		s.state.Phase = "completed"
	}
}

// finish records a terminal outcome and releases job admission under the lock.
func (s *server) finish(ctx context.Context, err error) bool {
	cancelled := ctx.Err() != nil
	s.cancel()
	s.cancel = nil
	if err == nil && cancelled {
		err = ctx.Err()
	}
	if err != nil {
		s.state.Phase, s.state.Error = "failed", err.Error()
		if cancelled {
			s.state.Phase = "cancelled"
		}
		s.config.Logger.Error("Browser operation stopped", "job", s.state.JobID, "phase", s.state.Phase, "error", err)
		return false
	}
	return true
}

// decode rejects oversized, unknown, and trailing input before changing session state.
func (s *server) decode(w http.ResponseWriter, r *http.Request, target any) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		s.fail(w, http.StatusUnsupportedMediaType, fmt.Errorf("expected application/json"))
		return false
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		s.fail(w, http.StatusBadRequest, fmt.Errorf("decode request: %w", err))
		return false
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		s.fail(w, http.StatusBadRequest, fmt.Errorf("request must contain one JSON object"))
		return false
	}
	return true
}

// respond writes one JSON response without logging request URLs or credentials.
func (s *server) respond(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		s.config.Logger.Debug("Browser response delivery ended", "error", err)
	}
}

// fail exposes a bounded operation error as data for text-only rendering.
func (s *server) fail(w http.ResponseWriter, code int, err error) {
	s.respond(w, code, map[string]string{"error": err.Error()})
}
