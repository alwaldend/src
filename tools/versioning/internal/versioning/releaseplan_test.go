package versioning

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestBuildReleaseRefPlanNightly(t *testing.T) {
	plan, err := BuildReleaseRefPlan(State{
		Version:   "2026.35.0-nightly.20260830",
		Channel:   "nightly",
		Branch:    "master",
		Commit:    abc123(),
		TreeState: "clean",
	}, "master")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Refs) != 1 || plan.Refs[0].Ref != "refs/tags/v2026.35.0-nightly.20260830" ||
		plan.Refs[0].Operation != "create" {
		t.Fatalf("unexpected nightly plan: %+v", plan.Refs)
	}
	if !plan.Atomic || plan.Lease != "" {
		t.Fatalf("unexpected atomic/lease: %+v", plan)
	}
	if err := plan.Validate(); err != nil {
		t.Fatalf("plan.Validate() error = %v", err)
	}
}

func TestBuildReleaseRefPlanRelease(t *testing.T) {
	plan, err := BuildReleaseRefPlan(State{
		Version:   "2026.35.3",
		Channel:   "release",
		Branch:    "releases/2026.35",
		Commit:    abc123(),
		TreeState: "clean",
	}, "master")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Refs) != 2 {
		t.Fatalf("release plan refs = %d, want 2", len(plan.Refs))
	}
	var tag, branch bool
	for _, ref := range plan.Refs {
		if ref.Ref == "refs/tags/v2026.35.3" && ref.Operation == "create" {
			tag = true
		}
		if ref.Ref == "refs/heads/releases/2026.35" && ref.Operation == "create" {
			branch = true
		}
	}
	if !tag || !branch {
		t.Fatalf("release plan missing tag/branch refs: %+v", plan.Refs)
	}
	if err := plan.Validate(); err != nil {
		t.Fatalf("plan.Validate() error = %v", err)
	}
}

func TestBuildReleaseRefPlanRejectsDevelopment(t *testing.T) {
	if _, err := BuildReleaseRefPlan(State{
		Version:   DevelopmentVersion,
		Channel:   "development",
		Branch:    "master",
		Commit:    abc123(),
		TreeState: "clean",
	}, "master"); err == nil {
		t.Fatal("development plan unexpectedly accepted")
	}
}

func TestReleaseRefPlanDigestIsDeterministic(t *testing.T) {
	state := State{
		Version: "2026.35.0", Channel: "release",
		Branch: "releases/2026.35", Commit: abc123(), TreeState: "clean",
	}
	first, err := BuildReleaseRefPlan(state, "master")
	if err != nil {
		t.Fatal(err)
	}
	firstDigest, err := first.Digest()
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildReleaseRefPlan(state, "master")
	if err != nil {
		t.Fatal(err)
	}
	secondDigest, err := second.Digest()
	if err != nil {
		t.Fatal(err)
	}
	if firstDigest != secondDigest {
		t.Fatalf("digest changed: %s != %s", firstDigest, secondDigest)
	}
	if !strings.HasPrefix(firstDigest, "sha256:") || len(firstDigest) != len("sha256:")+64 {
		t.Fatalf("digest = %q", firstDigest)
	}
}

func TestPublishReleaseRefsRequiresLeaseForAtomic(t *testing.T) {
	state := State{
		Version: "2026.35.0", Channel: "release",
		Branch: "releases/2026.35", Commit: abc123(), TreeState: "clean",
	}
	plan, err := BuildReleaseRefPlan(state, "master")
	if err != nil {
		t.Fatal(err)
	}
	plan.Lease = ""
	publisher := &staticPublisher{plan: plan, atomicSupport: true}
	_, err = PublishReleaseRefs(context.Background(), publisher, plan)
	if err == nil || !strings.Contains(err.Error(), "explicit lease") {
		t.Fatalf("PublishReleaseRefs() error = %v, want lease refusal", err)
	}
	plan.Lease = "release-refs:lease-1"
	receipt, err := PublishReleaseRefs(context.Background(), publisher, plan)
	if err != nil {
		t.Fatalf("PublishReleaseRefs() with lease error = %v", err)
	}
	if !receipt.Verified {
		t.Fatalf("receipt not verified: %+v", receipt)
	}
}

func TestPublishReleaseRefsRefusesUnsupportedAtomicity(t *testing.T) {
	state := State{
		Version: "2026.35.0", Channel: "release",
		Branch: "releases/2026.35", Commit: abc123(), TreeState: "clean",
	}
	plan, err := BuildReleaseRefPlan(state, "master")
	if err != nil {
		t.Fatal(err)
	}
	plan.Lease = "lease"
	publisher := &staticPublisher{plan: plan, atomicSupport: false}
	_, err = PublishReleaseRefs(context.Background(), publisher, plan)
	if err == nil || !strings.Contains(err.Error(), "atomic multi-ref") {
		t.Fatalf("PublishReleaseRefs() error = %v, want atomicity refusal", err)
	}
	if publisher.publishCalled {
		t.Fatal("publisher mutated before refusing unsupported atomicity")
	}
}

func TestPublishReleaseRefsEmitsVerifiedReceipt(t *testing.T) {
	state := State{
		Version: "2026.35.0", Channel: "release",
		Branch: "releases/2026.35", Commit: abc123(), TreeState: "clean",
	}
	plan, err := BuildReleaseRefPlan(state, "master")
	if err != nil {
		t.Fatal(err)
	}
	plan.Lease = "lease"
	publisher := &staticPublisher{plan: plan, atomicSupport: true}
	receipt, err := PublishReleaseRefs(context.Background(), publisher, plan)
	if err != nil {
		t.Fatal(err)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("receipt.Validate() error = %v", err)
	}
	if !receipt.Verified || receipt.PlanDigest == "" || len(receipt.Refs) < 2 {
		t.Fatalf("unexpected receipt: %+v", receipt)
	}
}

func TestPublishReleaseRefsRejectsStaleRemoteState(t *testing.T) {
	state := State{
		Version: "2026.35.0", Channel: "release",
		Branch: "releases/2026.35", Commit: abc123(), TreeState: "clean",
	}
	plan, err := BuildReleaseRefPlan(state, "master")
	if err != nil {
		t.Fatal(err)
	}
	plan.Lease = "lease"
	publisher := &staticPublisher{plan: plan, atomicSupport: true, staleRemote: true}
	_, err = PublishReleaseRefs(context.Background(), publisher, plan)
	if err == nil || !strings.Contains(err.Error(), "remote ref") {
		t.Fatalf("PublishReleaseRefs() error = %v, want stale-remote refusal", err)
	}
	if publisher.publishCalled {
		t.Fatal("publisher mutated on stale remote state")
	}
}

func abc123() string {
	return "0123456789abcdef0123456789abcdef01234567"
}

type staticPublisher struct {
	plan          ReleaseRefPlan
	atomicSupport bool
	staleRemote   bool
	publishCalled bool
}

func (p *staticPublisher) FetchExpectedRemoteState(
	ctx context.Context,
	plan ReleaseRefPlan,
) ([]RemoteRefExpectation, error) {
	expectations := make([]RemoteRefExpectation, 0, len(plan.Refs))
	for _, ref := range plan.Refs {
		commit := ref.Expected
		if p.staleRemote && !p.publishCalled {
			commit = "ffffffffffffffffffffffffffffffffffffffff"
		}
		expectations = append(expectations, RemoteRefExpectation{
			Ref: ref.Ref, Commit: commit,
		})
	}
	// Always include the trunk ref at its expected commit for plan
	// preconditions to pass or fail deterministically.
	for _, expectation := range plan.ExpectedRemote {
		if !isPresent(expectations, expectation.Ref) {
			commit := expectation.Commit
			if p.staleRemote && !p.publishCalled {
				commit = "ffffffffffffffffffffffffffffffffffffffff"
			}
			expectations = append(expectations, RemoteRefExpectation{
				Ref: expectation.Ref, Commit: commit,
			})
		}
	}
	return expectations, nil
}

func isPresent(expectations []RemoteRefExpectation, ref string) bool {
	for _, expectation := range expectations {
		if expectation.Ref == ref {
			return true
		}
	}
	return false
}

func (p *staticPublisher) AcquireLease(
	ctx context.Context,
	plan ReleaseRefPlan,
) (string, error) {
	if plan.Lease == "" {
		return "", errors.New("plan lease required")
	}
	return "lease", nil
}

func (p *staticPublisher) ReleaseLease(ctx context.Context, lease string) error {
	return nil
}

func (p *staticPublisher) Publish(
	ctx context.Context,
	plan ReleaseRefPlan,
	lease string,
) error {
	if p.staleRemote {
		return nil
	}
	if plan.Atomic && !p.atomicSupport {
		p.publishCalled = false
		return ErrAtomicityUnsupported
	}
	p.publishCalled = true
	return nil
}

func (p *staticPublisher) Verify(
	ctx context.Context,
	plan ReleaseRefPlan,
) ([]RemoteRefExpectation, error) {
	return p.FetchExpectedRemoteState(ctx, plan)
}

var _ = time.Now
