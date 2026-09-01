package control

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testNowValue() time.Time {
	return time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
}

func newTestKernel(t *testing.T, namespace string) *Kernel {
	t.Helper()
	kernel, err := New(KernelOptions{
		Root:       t.TempDir(),
		Namespace:  namespace,
		RuntimeID:  "runtime-test",
		ObserveNow: func() time.Time { return testNowValue() },
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = kernel.Unlock("test-worker")
	})
	return kernel
}

func TestPackageStatesAndTimeout(t *testing.T) {
	kernel := newTestKernel(t, "states")
	if err := kernel.RegisterPackage("pkg.a", "project", "rev-1", "hash-a", 10*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if err := kernel.Mark("pkg.a", PackageReady); err != nil {
		t.Fatal(err)
	}
	status := kernel.Status()
	if len(status) != 1 || status[0].State != PackageReady {
		t.Fatalf("unexpected status: %+v", status)
	}
	if !kernel.Health() {
		t.Fatalf("healthy kernel reported unhealthy")
	}
}

func TestNeverSettlingPackageTimesOut(t *testing.T) {
	clock := testNowValue()
	kernel, err := New(KernelOptions{
		Root:       t.TempDir(),
		Namespace:  "timeout",
		RuntimeID:  "runtime-test",
		ObserveNow: func() time.Time { return clock },
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := kernel.RegisterPackage("pkg.b", "scratch", "rev-1", "hash-b", 2*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	clock = clock.Add(5 * time.Millisecond)
	status := kernel.Status()
	if len(status) != 1 || status[0].State != PackageTimeout {
		t.Fatalf("expected timed-out state, got %+v", status)
	}
	if kernel.Health() {
		t.Fatalf("timed-out package must not report kernel health")
	}
}

func TestLockIsExclusiveAcrossKernels(t *testing.T) {
	root := t.TempDir()
	first, err := New(KernelOptions{Root: root, Namespace: "cross", RuntimeID: "runtime-1", ObserveNow: testNowValue})
	if err != nil {
		t.Fatal(err)
	}
	second, err := New(KernelOptions{Root: root, Namespace: "cross", RuntimeID: "runtime-2", ObserveNow: testNowValue})
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Lock("worker-a"); err != nil {
		t.Fatal(err)
	}
	if err := second.Lock("worker-b"); !errors.Is(err, ErrLocked) {
		t.Fatalf("second kernel lock error = %v, want ErrLocked", err)
	}
	if err := first.Unlock("worker-a"); err != nil {
		t.Fatal(err)
	}
	if err := second.Lock("worker-b"); err != nil {
		t.Fatalf("lock after release: %v", err)
	}
}

func TestLeaseAndExpectedRevisionPublication(t *testing.T) {
	kernel := newTestKernel(t, "lease")
	lease, err := kernel.PublishLease("worker-a", "rev-2", "rev-1", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if lease.RuntimeID != "runtime-test" {
		t.Fatalf("lease runtime mismatch: %+v", lease)
	}
	if err := kernel.VerifyLease("worker-a", "rev-2"); err != nil {
		t.Fatalf("verify lease: %v", err)
	}
	asset, err := kernel.PublishAsset("worker-a", "rev-2", "hash-2", "manifests/a.json")
	if err != nil {
		t.Fatal(err)
	}
	if asset.State != "published" || asset.Revision != "rev-2" {
		t.Fatalf("unexpected asset: %+v", asset)
	}
	read, err := kernel.ReadAsset()
	if err != nil {
		t.Fatal(err)
	}
	if read.Revision != "rev-2" {
		t.Fatalf("read revision = %s, want rev-2", read.Revision)
	}
}

func TestStaleRevisionPublicationRejected(t *testing.T) {
	kernel := newTestKernel(t, "stale")
	if _, err := kernel.PublishLease("worker-a", "rev-9", "rev-1", time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := kernel.VerifyLease("worker-a", "rev-9"); err != nil {
		t.Fatalf("verify lease: %v", err)
	}
	expiredClock := testNowValue()
	expiredKernel, err := New(KernelOptions{
		Root:       t.TempDir(),
		Namespace:  "expired",
		RuntimeID:  "runtime-test",
		ObserveNow: func() time.Time { return expiredClock },
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := expiredKernel.PublishLease("worker-a", "rev-9", "rev-1", time.Minute); err != nil {
		t.Fatal(err)
	}
	expiredClock = expiredClock.Add(2 * time.Minute)
	if _, err := expiredKernel.PublishAsset("worker-a", "rev-9", "hash-9", "manifests/a.json"); !errors.Is(err, ErrLeaseExpired) {
		t.Fatalf("expired lease error = %v, want ErrLeaseExpired", err)
	}

	clock := testNowValue()
	kernel2, err := New(KernelOptions{
		Root:       t.TempDir(),
		Namespace:  "stale-2",
		RuntimeID:  "runtime-test",
		ObserveNow: func() time.Time { return clock },
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := kernel2.PublishLease("worker-a", "rev-2", "rev-1", time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, err := kernel2.PublishAsset("worker-a", "rev-1", "hash-1", "manifests/a.json"); err != nil {
		t.Fatal(err)
	}
	if _, err := kernel2.PublishAsset("worker-a", "rev-2", "hash-2", "manifests/a.json"); !errors.Is(err, ErrStaleRevision) {
		t.Fatalf("stale current publication error = %v, want ErrStaleRevision", err)
	}

	expired := testNowValue()
	kernel3, err := New(KernelOptions{
		Root:       t.TempDir(),
		Namespace:  "stale-3",
		RuntimeID:  "runtime-test",
		ObserveNow: func() time.Time { return expired },
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := kernel3.PublishLease("worker-a", "rev-5", "rev-1", time.Minute); err != nil {
		t.Fatal(err)
	}
	expired = expired.Add(2 * time.Minute)
	if _, err := kernel3.PublishAsset("worker-a", "rev-5", "hash-5", "manifests/a.json"); !errors.Is(err, ErrLeaseExpired) {
		t.Fatalf("expired lease error = %v, want ErrLeaseExpired", err)
	}
}

func TestPublishAssetIfMissingIsIdempotent(t *testing.T) {
	kernel := newTestKernel(t, "missing")
	if _, err := kernel.PublishLease("worker-a", "rev-1", "", time.Minute); err != nil {
		t.Fatal(err)
	}
	first, created, err := kernel.PublishAssetIfMissing("worker-a", "rev-1", "hash-1", "manifests/a.json")
	if err != nil || !created {
		t.Fatalf("first publish: created=%t err=%v", created, err)
	}
	if first.Revision != "rev-1" {
		t.Fatalf("first revision = %s", first.Revision)
	}
	_, secondCreated, err := kernel.PublishAssetIfMissing("worker-a", "rev-1", "hash-1", "manifests/a.json")
	if err != nil {
		t.Fatal(err)
	}
	if secondCreated {
		t.Fatalf("second publish must not create")
	}
}

func TestKernelRejectsMalformedPackage(t *testing.T) {
	kernel := newTestKernel(t, "malformed")
	if err := kernel.RegisterPackage("../escape", "project", "rev", "hash", time.Second); err == nil ||
		!strings.Contains(err.Error(), "malformed") {
		t.Fatalf("RegisterPackage() error = %v, want malformed rejection", err)
	}
}

func TestNamespacesAreIsolated(t *testing.T) {
	root := t.TempDir()
	first, err := New(KernelOptions{Root: root, Namespace: "task-a", RuntimeID: "runtime-a", ObserveNow: testNowValue})
	if err != nil {
		t.Fatal(err)
	}
	second, err := New(KernelOptions{Root: root, Namespace: "task-b", RuntimeID: "runtime-b", ObserveNow: testNowValue})
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Lock("worker-a"); err != nil {
		t.Fatal(err)
	}
	if err := second.Lock("worker-b"); err != nil {
		t.Fatalf("different namespace lock must not conflict: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "assets", "task-b", "lease.json")); !os.IsNotExist(err) {
		t.Fatalf("task-b must have no lease")
	}
	if _, err := first.PublishLease("worker-a", "rev-1", "", time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "assets", "task-a", "lease.json")); err != nil {
		t.Fatalf("task-a lease must be isolated: %v", err)
	}
}
