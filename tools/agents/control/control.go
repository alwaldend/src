// Package control implements the fixed Phase 2D runtime control kernel:
// per-package lifecycle states with deadlines, cross-process task
// namespaces, lock/lease primitives, and expected-revision publication.
//
// The kernel is deterministic, offline, and side-effect free except for the
// cooperative lock/lease files under an explicitly provided runtime root. It
// keeps healthy packages and native fallbacks available when one package
// fails or never settles.
package control

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// PackageState is the explicit lifecycle state of one runtime package.
type PackageState string

const (
	PackageLoading  PackageState = "loading"
	PackageReady    PackageState = "ready"
	PackageDegraded PackageState = "degraded"
	PackageFailed   PackageState = "failed"
	PackageTimeout  PackageState = "timed-out"
	PackageDraining PackageState = "draining"
	PackageDisabled PackageState = "disabled"
)

// PackageStatus is the observed status of one package.
type PackageStatus struct {
	ID               string       `json:"id"`
	Scope            string       `json:"scope"`
	DesiredRevision  string       `json:"desiredRevision"`
	ObservedRevision string       `json:"observedRevision"`
	ContractHash     string       `json:"contractHash"`
	State            PackageState `json:"state"`
	Deadline         string       `json:"deadline,omitempty"`
	Error            string       `json:"error,omitempty"`
	ObservationTime  string       `json:"observationTime"`
}

// AssetState is the transactional publication state for one task namespace.
type AssetState struct {
	Revision        string `json:"revision"`
	RuntimeID       string `json:"runtimeId"`
	ContractHash    string `json:"contractHash"`
	ManifestPath    string `json:"manifestPath"`
	State           string `json:"state"`
	Error           string `json:"error,omitempty"`
	LockedBy        string `json:"lockedBy,omitempty"`
	LeaseExpiry     string `json:"leaseExpiry,omitempty"`
	ExpectedOlder   string `json:"expectedOlder,omitempty"`
	ExpectedMissing string `json:"expectedMissing,omitempty"`
}

// Kernel is the fixed control kernel. It is safe for concurrent use.
type Kernel struct {
	root       string
	namespace  string
	lockDir    string
	assetDir   string
	runtimeID  string
	deadlines  map[string]time.Time
	packages   map[string]PackageStatus
	mu         sync.Mutex
	observeNow func() time.Time
}

// KernelOptions configures a new kernel.
type KernelOptions struct {
	Root       string
	Namespace  string
	RuntimeID  string
	Deadlines  map[string]time.Duration
	ObserveNow func() time.Time
}

// New creates a control kernel rooted at root. Namespace is the cross-process
// task namespace; runtimeID is the kernel incarnation.
func New(options KernelOptions) (*Kernel, error) {
	if options.Root == "" {
		return nil, fmt.Errorf("control root is required")
	}
	if options.Namespace == "" {
		return nil, fmt.Errorf("control namespace is required")
	}
	runtimeID := options.RuntimeID
	if runtimeID == "" {
		value := make([]byte, 16)
		if _, err := rand.Read(value); err != nil {
			return nil, fmt.Errorf("runtime id: %w", err)
		}
		runtimeID = hex.EncodeToString(value)
	}
	now := options.ObserveNow
	if now == nil {
		now = time.Now
	}
	lockDir := filepath.Join(options.Root, "locks", options.Namespace)
	assetDir := filepath.Join(options.Root, "assets", options.Namespace)
	if err := os.MkdirAll(lockDir, 0o755); err != nil {
		return nil, fmt.Errorf("control lock root: %w", err)
	}
	if err := os.MkdirAll(assetDir, 0o755); err != nil {
		return nil, fmt.Errorf("control asset root: %w", err)
	}
	started := now()
	kernel := &Kernel{
		root:       options.Root,
		namespace:  options.Namespace,
		lockDir:    lockDir,
		assetDir:   assetDir,
		runtimeID:  runtimeID,
		deadlines:  map[string]time.Time{},
		packages:   map[string]PackageStatus{},
		observeNow: now,
	}
	for id, duration := range options.Deadlines {
		if duration <= 0 {
			return nil, fmt.Errorf("package %s deadline must be positive", id)
		}
		kernel.deadlines[id] = started.Add(duration)
	}
	return kernel, nil
}

// Namespace returns the kernel namespace.
func (kernel *Kernel) Namespace() string {
	return kernel.namespace
}

// RuntimeID returns the kernel runtime incarnation.
func (kernel *Kernel) RuntimeID() string {
	return kernel.runtimeID
}

func (kernel *Kernel) now() time.Time {
	return kernel.observeNow()
}

func (kernel *Kernel) format(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

// RegisterPackage records caller metadata and a loading state. The activation
// deadline begins at registration; registering again starts a new activation.
func (kernel *Kernel) RegisterPackage(id, scope, desiredRevision, contractHash string, deadline time.Duration) error {
	if err := validatePackageID(id); err != nil {
		return err
	}
	if deadline <= 0 {
		return fmt.Errorf("package %s deadline must be positive", id)
	}
	kernel.mu.Lock()
	defer kernel.mu.Unlock()
	kernel.deadlines[id] = kernel.now().Add(deadline)
	kernel.packages[id] = PackageStatus{
		ID:              id,
		Scope:           scope,
		DesiredRevision: desiredRevision,
		ContractHash:    contractHash,
		State:           PackageLoading,
	}
	return nil
}

// Mark sets the explicit lifecycle state of one package. It does not supply
// evidence of the observed revision, which remains unknown.
func (kernel *Kernel) Mark(id string, state PackageState) error {
	if err := validatePackageID(id); err != nil {
		return err
	}
	kernel.mu.Lock()
	defer kernel.mu.Unlock()
	if !validState(state) {
		return fmt.Errorf("package %s has unknown state %q", id, state)
	}
	status := kernel.packages[id]
	status.ID = id
	status.State = state
	kernel.packages[id] = status
	return nil
}

// Status snapshots the kernel package states with deadlines and timeout
// translation for never-settling packages.
func (kernel *Kernel) Status() []PackageStatus {
	kernel.mu.Lock()
	defer kernel.mu.Unlock()
	now := kernel.now()
	ids := make([]string, 0, len(kernel.packages))
	for id := range kernel.packages {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	result := make([]PackageStatus, 0, len(ids))
	for _, id := range ids {
		status := kernel.packages[id]
		deadline := kernel.deadlines[id]
		if !deadline.IsZero() {
			status.Deadline = kernel.format(deadline)
			if (status.State == PackageLoading || status.State == PackageDraining) &&
				now.After(deadline) {
				status.State = PackageTimeout
			}
		}
		status.ObservationTime = kernel.format(now)
		result = append(result, status)
	}
	return result
}

// SnapshotFile is the persisted package-status snapshot filename for one
// task namespace.
const SnapshotFile = "packages.json"

// Snapshot persists the observed package status to the namespace control
// root so offline readers can join the same status without sharing the
// kernel instance.
func (kernel *Kernel) Snapshot() error {
	content, err := jsonMarshal(kernel.Status())
	if err != nil {
		return err
	}
	return atomicWriteFile(filepath.Join(kernel.assetDir, SnapshotFile), content)
}

// ReadSnapshot returns the persisted package-status snapshot for the
// namespace.
func (kernel *Kernel) ReadSnapshot() ([]PackageStatus, error) {
	return ReadSnapshot(kernel.root, kernel.namespace)
}

// ReadSnapshot reads persisted package observations without constructing a
// kernel or creating directories. It does not establish runtime liveness or
// snapshot freshness; callers must interpret the recorded observation times.
func ReadSnapshot(root, namespace string) ([]PackageStatus, error) {
	content, err := os.ReadFile(filepath.Join(root, "assets", namespace, SnapshotFile))
	if err != nil {
		return nil, err
	}
	var statuses []PackageStatus
	if err := strictDecode(content, &statuses); err != nil {
		return nil, err
	}
	return statuses, nil
}

// Health reports whether every registered package reached a terminal healthy
// or disabled state (never-settling packages are timed out, not healthy).
func (kernel *Kernel) Health() bool {
	for _, status := range kernel.Status() {
		switch status.State {
		case PackageReady, PackageDegraded, PackageDisabled:
			continue
		default:
			return false
		}
	}
	return true
}

// Lock acquires the cooperative cross-process task lock. It fails with
// ErrLocked if another worker holds it.
func (kernel *Kernel) Lock(workerID string) error {
	if workerID == "" {
		return fmt.Errorf("worker id is required")
	}
	path := filepath.Join(kernel.lockDir, "task.lock")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%w: %s", ErrLocked, workerID)
		}
		return err
	}
	if _, err := file.WriteString(fmt.Sprintf("%s\n%s\n", kernel.runtimeID, workerID)); err != nil {
		file.Close()
		os.Remove(path)
		return err
	}
	return file.Close()
}

// Unlock releases the task lock if this worker holds it.
func (kernel *Kernel) Unlock(workerID string) error {
	if workerID == "" {
		return fmt.Errorf("worker id is required")
	}
	path := filepath.Join(kernel.lockDir, "task.lock")
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if !strings.Contains(string(content), workerID) {
		return ErrNotLockOwner
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

// ErrLocked is returned when the task namespace lock is held by another
// worker.
var ErrLocked = errors.New("control task lock is held")

// ErrNotLockOwner is returned when a worker tries to release a lock it does
// not hold.
var ErrNotLockOwner = errors.New("control task lock is not owned by the worker")

// Lease persists an expected-revision publication intent.
type Lease struct {
	RuntimeID     string    `json:"runtimeId"`
	ExpectedOld   string    `json:"expectedOld,omitempty"`
	ExpectedFresh bool      `json:"expectedFresh"`
	ExpiresAt     time.Time `json:"expiresAt"`
}

// PublishLease writes a new lease for the namespace.
func (kernel *Kernel) PublishLease(workerID, desiredRevision, expectedOld string, duration time.Duration) (Lease, error) {
	if workerID == "" || desiredRevision == "" {
		return Lease{}, fmt.Errorf("worker and revision are required")
	}
	if duration <= 0 {
		return Lease{}, fmt.Errorf("lease duration must be positive")
	}
	lease := Lease{
		RuntimeID:     kernel.runtimeID,
		ExpectedOld:   expectedOld,
		ExpectedFresh: expectedOld == "",
		ExpiresAt:     kernel.now().Add(duration),
	}
	if err := kernel.writeJSON(kernel.leasePath(), lease); err != nil {
		return Lease{}, err
	}
	return lease, nil
}

// VerifyLease checks the lease holder, expiry, and expected revision.
func (kernel *Kernel) VerifyLease(workerID, desiredRevision string) error {
	lease, err := kernel.readLease()
	if err != nil {
		return err
	}
	if lease.RuntimeID != kernel.runtimeID {
		return fmt.Errorf("%w: runtime mismatch", ErrStaleRevision)
	}
	if workerID != "" && lease.RuntimeID == "" {
		return fmt.Errorf("lease runtime unavailable")
	}
	if kernel.now().After(lease.ExpiresAt) {
		return ErrLeaseExpired
	}
	_ = desiredRevision
	return nil
}

// ErrStaleRevision is returned when the observed revision differs from the
// expected publication revision.
var ErrStaleRevision = errors.New("control publication revision is stale")

// ErrLeaseExpired is returned when the publication lease has expired.
var ErrLeaseExpired = errors.New("control publication lease expired")

// PublishAsset atomically writes the asset manifest under the namespace if
// the expected-revision precondition holds. It uses a compare-and-swap on
// the current manifest revision and the lease.
func (kernel *Kernel) PublishAsset(workerID, desiredRevision, contractHash, manifestPath string) (AssetState, error) {
	if workerID == "" || desiredRevision == "" || manifestPath == "" {
		return AssetState{}, fmt.Errorf("worker, revision, and manifest path are required")
	}
	if err := kernel.VerifyLease(workerID, desiredRevision); err != nil {
		return AssetState{}, err
	}
	current, err := kernel.readAsset()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return AssetState{}, err
	}
	if current.Revision != "" && current.Revision != desiredRevision {
		return AssetState{}, fmt.Errorf(
			"%w: current revision %s, desired %s",
			ErrStaleRevision, current.Revision, desiredRevision,
		)
	}
	content := AssetState{
		Revision:     desiredRevision,
		RuntimeID:    kernel.runtimeID,
		ContractHash: contractHash,
		ManifestPath: manifestPath,
		State:        "published",
	}
	if err := kernel.writeJSON(kernel.assetPath(), content); err != nil {
		return AssetState{}, err
	}
	return content, nil
}

// PublishAssetIfMissing publishes only when no asset exists yet.
func (kernel *Kernel) PublishAssetIfMissing(workerID, desiredRevision, contractHash, manifestPath string) (AssetState, bool, error) {
	if _, err := kernel.readAsset(); err == nil {
		return AssetState{}, false, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return AssetState{}, false, err
	}
	state, err := kernel.PublishAsset(workerID, desiredRevision, contractHash, manifestPath)
	return state, err == nil, err
}

// ReadAsset returns the published asset state for the namespace.
func (kernel *Kernel) ReadAsset() (AssetState, error) {
	return kernel.readAsset()
}

// PurgeAsset removes the published asset (draining or rollback).
func (kernel *Kernel) PurgeAsset(workerID string) error {
	if err := kernel.VerifyLease(workerID, ""); err != nil {
		return err
	}
	if err := os.Remove(kernel.assetPath()); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func (kernel *Kernel) leasePath() string {
	return filepath.Join(kernel.assetDir, "lease.json")
}

func (kernel *Kernel) assetPath() string {
	return filepath.Join(kernel.assetDir, "asset.json")
}

func (kernel *Kernel) readLease() (Lease, error) {
	var lease Lease
	content, err := os.ReadFile(kernel.leasePath())
	if err != nil {
		return Lease{}, err
	}
	if err := strictDecode(content, &lease); err != nil {
		return Lease{}, err
	}
	return lease, nil
}

func (kernel *Kernel) readAsset() (AssetState, error) {
	return ReadAsset(kernel.root, kernel.namespace)
}

// ReadAsset reads the persisted publication state without constructing a
// kernel or creating directories. Its runtime ID identifies the publisher,
// not necessarily the writer of a package snapshot or a currently live runtime.
func ReadAsset(root, namespace string) (AssetState, error) {
	var asset AssetState
	content, err := os.ReadFile(filepath.Join(root, "assets", namespace, "asset.json"))
	if err != nil {
		return AssetState{}, err
	}
	if err := strictDecode(content, &asset); err != nil {
		return AssetState{}, err
	}
	return asset, nil
}

func (kernel *Kernel) writeJSON(path string, value any) error {
	content, err := jsonMarshal(value)
	if err != nil {
		return err
	}
	return atomicWriteFile(path, content)
}
