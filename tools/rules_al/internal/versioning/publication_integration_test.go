package versioning

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestPatchPublication runs the source-ref lifecycle against an isolated Git remote.
func TestPatchPublication(t *testing.T) {
	f := newGitFixture(t, "")
	remote := filepath.Join(t.TempDir(), "remote.git")
	runFixtureGit(t, f.directory, "init", "--bare", remote)
	runFixtureGit(t, f.directory, "remote", "add", "origin", remote)
	runFixtureGit(t, f.directory, "push", "origin", "master")
	if _, err := f.repository.StartRelease(time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC), true, false); err != nil {
		t.Fatal(err)
	}
	p := &LocalGitReleasePublisher{Git: f.repository.Git, Remote: "origin"}
	publish := func() ReleaseRefReceipt {
		t.Helper()
		state, err := f.repository.Inspect()
		if err != nil {
			t.Fatal(err)
		}
		plan, err := BuildReleaseRefPlan(state, "master")
		if err != nil {
			t.Fatal(err)
		}
		plan.Lease = "fixture"
		receipt, err := PublishReleaseRefs(context.Background(), p, plan)
		if err != nil {
			t.Fatal(err)
		}
		return receipt
	}
	publish()
	runFixtureGit(t, f.directory, "commit", "--allow-empty", "-m", "patch")
	if _, err := f.repository.TagRelease(false); err != nil {
		t.Fatal(err)
	}
	runFixtureGit(t, f.directory, "push", "origin", "releases/2026.44")
	receipt := publish()
	publish() // A second project can reuse the shared source refs.
	state, err := f.repository.Inspect()
	if err != nil {
		t.Fatal(err)
	}
	plan, err := BuildReleaseRefPlan(state, "master")
	if err != nil {
		t.Fatal(err)
	}
	lease, err := p.AcquireLease(context.Background(), plan)
	if err != nil {
		t.Fatal(err)
	}
	other := &LocalGitReleasePublisher{Git: f.repository.Git, Remote: "origin"}
	if _, err := other.AcquireLease(context.Background(), plan); err == nil {
		t.Fatal("concurrent publisher acquired the operation lock")
	}
	if err := p.ReleaseLease(context.Background(), lease); err != nil {
		t.Fatal(err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(receipt, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "patch-receipt.json"), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
