package versioning

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// ReleaseRefPlan is a deterministic, reviewed plan for publishing immutable
// release and nightly refs. It never represents authorization: a publisher
// must consume the exact plan, require distinct release-refs authority, and
// emit a ReleaseRefReceipt after verified publication.
type ReleaseRefPlan struct {
	Schema string `json:"schema"`
	Kind   string `json:"kind"`

	// Version is the exact calculated version resolved at plan time.
	Version string `json:"version"`
	// Channel is one of development, nightly, or release.
	Channel string `json:"channel"`
	// Commit is the exact commit that all plan refs will point to.
	Commit string `json:"commit"`
	// TreeState is clean or modified at plan resolution.
	TreeState string `json:"treeState"`

	// Refs lists the exact refs this plan will create or verify, sorted.
	// A nightly plan contains only tags; a release plan contains the branch
	// and its patch-zero tag.
	Refs []ReleasePlanRef `json:"refs"`

	// ExpectedRemote is the publisher's precondition snapshot at plan time.
	// It is advisory context, not an authority: the publisher re-fetches and
	// verifies remote state immediately before publication.
	ExpectedRemote []RemoteRefExpectation `json:"expectedRemote,omitempty"`

	// Atomic declares whether the target publisher supports atomic
	// multi-ref publication. Unsupported atomicity is an explicit refusal,
	// never generic success.
	Atomic bool `json:"atomic"`

	// Lease is the optional explicit release-refs lease used by the guarded
	// publisher. When absent, the publisher requires an empty lease from the
	// separate release-refs authority.
	Lease string `json:"lease,omitempty"`
}

// ReleasePlanRef is one exact ref expectation in a verified publication
// plan. Create refs must be absent; Verify refs must already point at the
// plan commit.
type ReleasePlanRef struct {
	Ref       string `json:"ref"`
	Expected  string `json:"expected"`
	Operation string `json:"operation"`
}

// RemoteRefExpectation binds an exact remote ref to an expected commit. It
// records only stable identities, never credentials or URLs.
type RemoteRefExpectation struct {
	Ref    string `json:"ref"`
	Commit string `json:"commit"`
}

// ReleaseRefReceipt proves that one guarded publication consumed an exact
// release-ref plan and left the remote refs in the expected state. It is
// local audit evidence, not an attestation.
type ReleaseRefReceipt struct {
	Schema      string                 `json:"schema"`
	PlanDigest  string                 `json:"plan_digest"`
	Version     string                 `json:"version"`
	Channel     string                 `json:"channel"`
	Commit      string                 `json:"commit"`
	TreeState   string                 `json:"treeState"`
	PublishedAt string                 `json:"published_at"`
	Atomic      bool                   `json:"atomic"`
	Lease       string                 `json:"lease,omitempty"`
	Refs        []RemoteRefExpectation `json:"refs"`
	Verified    bool                   `json:"verified"`
}

const (
	ReleaseRefPlanSchema    = "versioning/release_ref_plan/v1"
	ReleaseRefReceiptSchema = "versioning/release_ref_receipt/v1"
	releaseRefsAuthorityID  = "release-refs"
)

// Validate checks plan invariants. The plan is deterministic and must be
// reviewed before consumption.
func (p ReleaseRefPlan) Validate() error {
	if p.Schema != ReleaseRefPlanSchema || p.Kind != "ReleaseRefPlan" {
		return fmt.Errorf("release ref plan has an unsupported schema or kind")
	}
	if p.Version == "" || p.Channel == "development" {
		return fmt.Errorf("release ref plan must name a concrete version channel")
	}
	if !oneOf(p.Channel, "nightly", "release") {
		return fmt.Errorf("release ref plan has unknown channel %q", p.Channel)
	}
	if !validObjectID(p.Commit) {
		return fmt.Errorf("release ref plan commit is not a full object ID")
	}
	if p.TreeState != "clean" {
		return fmt.Errorf("release ref plan requires a clean tree")
	}
	if len(p.Refs) == 0 {
		return fmt.Errorf("release ref plan lists no refs")
	}
	previous := ""
	for _, ref := range p.Refs {
		if ref.Ref == "" || ref.Ref == previous {
			return fmt.Errorf("release ref plan refs must be sorted and unique")
		}
		previous = ref.Ref
		if !isOwnedRef(ref.Ref) {
			return fmt.Errorf("release ref plan ref %q is not owned", ref.Ref)
		}
		if !validObjectID(ref.Expected) {
			return fmt.Errorf("release ref plan ref %q has no verified expected commit", ref.Ref)
		}
		if ref.Operation != "create" && ref.Operation != "verify" {
			return fmt.Errorf("release ref plan ref %q has unknown operation", ref.Ref)
		}
		if ref.Operation == "create" && ref.Expected != p.Commit {
			return fmt.Errorf("release ref plan create ref %q must match the plan commit", ref.Ref)
		}
		if ref.Operation == "verify" && ref.Expected != p.Commit {
			return fmt.Errorf("release ref plan verify ref %q must match the plan commit", ref.Ref)
		}
	}
	for _, expectation := range p.ExpectedRemote {
		if !validObjectID(expectation.Commit) {
			return fmt.Errorf("release ref plan expected remote has an invalid ref or commit")
		}
	}
	return nil
}

// Digest returns the canonical content-address of the plan with the digest
// omitted. It is deterministic across encoders.
func (p ReleaseRefPlan) Digest() (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}
	canonical, err := CanonicalJSON(p)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// CanonicalJSON encodes value deterministically with sorted keys for use as
// a content-addressable body.
func CanonicalJSON(value any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}

func validObjectID(value string) bool {
	if len(value) != 40 && len(value) != 64 {
		return false
	}
	for _, character := range value {
		if !(character >= '0' && character <= '9') &&
			!(character >= 'a' && character <= 'f') {
			return false
		}
	}
	return true
}

func isOwnedRef(ref string) bool {
	return strings.HasPrefix(ref, "refs/tags/v") ||
		strings.HasPrefix(ref, "refs/heads/releases/")
}

// Validate checks receipt invariants.
func (r ReleaseRefReceipt) Validate() error {
	if r.Schema != ReleaseRefReceiptSchema {
		return fmt.Errorf("release ref receipt has an unsupported schema")
	}
	if !validSHA256Hex(r.PlanDigest) || !validObjectID(r.Commit) {
		return fmt.Errorf("release ref receipt has an invalid plan digest or commit")
	}
	if !oneOf(r.Channel, "nightly", "release") || r.TreeState != "clean" {
		return fmt.Errorf("release ref receipt has an invalid channel or tree state")
	}
	if len(r.Refs) == 0 {
		return fmt.Errorf("release ref receipt lists no refs")
	}
	if r.PublishedAt == "" || !r.Verified {
		return fmt.Errorf("release ref receipt must record verified publication")
	}
	for _, expectation := range r.Refs {
		if !isOwnedRef(expectation.Ref) || !validObjectID(expectation.Commit) {
			return fmt.Errorf("release ref receipt has an invalid bound ref")
		}
	}
	return nil
}

// ReleaseRefPublisher is the provider-neutral guarded publisher boundary. It
// fetches expected remote state, acquires an explicit lease, atomically
// publishes the refs when support exists, and verifies the remote result.
// Unsupported atomicity or observation is an explicit refusal or unknown.
type ReleaseRefPublisher interface {
	// FetchExpectedRemoteState returns the current remote refs that the plan
	// expects. It must not mutate anything.
	FetchExpectedRemoteState(ctx context.Context, plan ReleaseRefPlan) ([]RemoteRefExpectation, error)
	// AcquireLease acquires a distinct release-refs lease. When refused, the
	// caller must stop before any remote mutation.
	AcquireLease(ctx context.Context, plan ReleaseRefPlan) (string, error)
	// ReleaseLease releases a previously acquired lease.
	ReleaseLease(ctx context.Context, lease string) error
	// Publish applies the plan refs to the remote. When the remote does not
	// support atomic multi-ref publication and plan.Atomic is true, it must
	// return ErrAtomicityUnsupported before mutating anything.
	Publish(ctx context.Context, plan ReleaseRefPlan, lease string) error
	// Verify confirms the remote refs match the plan exactly.
	Verify(ctx context.Context, plan ReleaseRefPlan) ([]RemoteRefExpectation, error)
}

// ErrAtomicityUnsupported is the explicit refusal for providers that cannot
// guarantee atomic multi-ref publication.
var ErrAtomicityUnsupported = errors.New("remote does not support atomic multi-ref publication")

// PublishReleaseRefs consumes an exact reviewed plan through a
// provider-neutral guarded publisher and emits a verified ReleaseRefReceipt.
// It never falls back to generic success when observation is unavailable.
func PublishReleaseRefs(
	ctx context.Context,
	publisher ReleaseRefPublisher,
	plan ReleaseRefPlan,
) (ReleaseRefReceipt, error) {
	if publisher == nil {
		return ReleaseRefReceipt{}, fmt.Errorf("release ref publisher is required")
	}
	digest, err := plan.Digest()
	if err != nil {
		return ReleaseRefReceipt{}, err
	}
	remoteState, err := publisher.FetchExpectedRemoteState(ctx, plan)
	if err != nil {
		return ReleaseRefReceipt{}, fmt.Errorf("fetch expected remote state: %w", err)
	}
	if err := requireRemotePreconditions(plan, remoteState); err != nil {
		return ReleaseRefReceipt{}, err
	}
	if plan.Atomic && strings.TrimSpace(plan.Lease) == "" {
		return ReleaseRefReceipt{}, fmt.Errorf(
			"release ref plan requires an explicit lease for atomic publication",
		)
	}
	lease, err := publisher.AcquireLease(ctx, plan)
	if err != nil {
		return ReleaseRefReceipt{}, fmt.Errorf("acquire release-refs lease: %w", err)
	}
	released := false
	defer func() {
		if !released {
			_ = publisher.ReleaseLease(ctx, lease)
		}
	}()
	if plan.Atomic {
		if err := publisher.Publish(ctx, plan, lease); err != nil {
			if errors.Is(err, ErrAtomicityUnsupported) {
				return ReleaseRefReceipt{}, fmt.Errorf(
					"refusing publication: %w", err,
				)
			}
			return ReleaseRefReceipt{}, fmt.Errorf("publish release refs: %w", err)
		}
	} else {
		// A non-atomic publisher may still apply each ref only after the
		// remote lease is held; partial failure is explicit, never success.
		if err := publisher.Publish(ctx, plan, lease); err != nil {
			return ReleaseRefReceipt{}, fmt.Errorf("publish release refs: %w", err)
		}
	}
	verified, err := publisher.Verify(ctx, plan)
	if err != nil {
		return ReleaseRefReceipt{}, fmt.Errorf("verify release refs: %w", err)
	}
	binding, err := requireRemoteMatches(plan, verified)
	if err != nil {
		return ReleaseRefReceipt{}, err
	}
	released = true
	if err := publisher.ReleaseLease(ctx, lease); err != nil {
		return ReleaseRefReceipt{}, fmt.Errorf(
			"release refs published and verified, but lease release failed: %w", err,
		)
	}
	return ReleaseRefReceipt{
		Schema:      ReleaseRefReceiptSchema,
		PlanDigest:  digest,
		Version:     plan.Version,
		Channel:     plan.Channel,
		Commit:      plan.Commit,
		TreeState:   plan.TreeState,
		PublishedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Atomic:      plan.Atomic,
		Lease:       lease,
		Refs:        binding,
		Verified:    true,
	}, nil
}

// requireRemotePreconditions checks that the observed remote refs are
// consistent with the plan before any mutation. A create ref must be absent
// or already at its exact expected commit; a verify ref must already point at
// the expected commit.
func requireRemotePreconditions(
	plan ReleaseRefPlan,
	observed []RemoteRefExpectation,
) error {
	observedByRef := make(map[string]string, len(observed))
	for _, expectation := range observed {
		observedByRef[expectation.Ref] = expectation.Commit
	}
	for _, expectation := range plan.ExpectedRemote {
		commit, present := observedByRef[expectation.Ref]
		if !present || commit != expectation.Commit {
			return fmt.Errorf(
				"remote ref %s = %q, want expected %s",
				expectation.Ref,
				commit,
				expectation.Commit,
			)
		}
	}
	for _, ref := range plan.Refs {
		commit, present := observedByRef[ref.Ref]
		if ref.Operation == "verify" && !present {
			return fmt.Errorf(
				"remote ref %s is absent, want %s", ref.Ref, ref.Expected,
			)
		}
		if present && commit != ref.Expected {
			return fmt.Errorf(
				"remote ref %s = %q, want %s", ref.Ref, commit, ref.Expected,
			)
		}
	}
	return nil
}

// requireRemoteMatches checks that every plan ref is present with the exact
// expected commit after verified publication and returns the plan-bound
// expectations. Unrelated observed refs are ignored.
func requireRemoteMatches(
	plan ReleaseRefPlan,
	observed []RemoteRefExpectation,
) ([]RemoteRefExpectation, error) {
	observedByRef := make(map[string]string, len(observed))
	for _, expectation := range observed {
		observedByRef[expectation.Ref] = expectation.Commit
	}
	binding := make([]RemoteRefExpectation, 0, len(plan.Refs))
	for _, ref := range plan.Refs {
		commit, ok := observedByRef[ref.Ref]
		if !ok || commit != ref.Expected {
			return nil, fmt.Errorf(
				"remote ref %s = %q, want %s", ref.Ref, commit, ref.Expected,
			)
		}
		binding = append(binding, RemoteRefExpectation{
			Ref: ref.Ref, Commit: commit,
		})
	}
	sort.Slice(binding, func(i, j int) bool { return binding[i].Ref < binding[j].Ref })
	return binding, nil
}

// BuildReleaseRefPlan resolves the exact release-ref plan for a reviewed
// version state. The state must be concrete (nightly or release), clean, and
// consistent with the owned refs. It never mutates anything.
func BuildReleaseRefPlan(state State, trunkBranch string) (ReleaseRefPlan, error) {
	if state.Channel != "nightly" && state.Channel != "release" {
		return ReleaseRefPlan{}, fmt.Errorf(
			"cannot build a release ref plan for channel %q", state.Channel,
		)
	}
	if state.TreeState != "clean" {
		return ReleaseRefPlan{}, fmt.Errorf(
			"release ref plan requires a clean tree",
		)
	}
	if !validObjectID(state.Commit) {
		return ReleaseRefPlan{}, fmt.Errorf(
			"release ref plan requires a verified commit",
		)
	}
	tag := "refs/tags/v" + state.Version
	refs := []ReleasePlanRef{{Ref: tag, Expected: state.Commit, Operation: "create"}}
	expectedRemote := []RemoteRefExpectation{}
	if state.Channel == "release" {
		const releasePrefix = "releases/"
		if !strings.HasPrefix(state.Branch, releasePrefix) {
			return ReleaseRefPlan{}, fmt.Errorf(
				"release state %+v is not on a release branch", state,
			)
		}
		branchRef := "refs/heads/" + state.Branch
		refs = append(refs, ReleasePlanRef{
			Ref: branchRef, Expected: state.Commit, Operation: "create",
		})
		// The branch and patch-zero tag are both created by this plan. Prefer
		// an explicit expected-remote that verifies the caller's reviewed
		// source ref (the trunk branch) is unchanged at the plan commit.
		expectedRemote = append(expectedRemote, RemoteRefExpectation{
			Ref: "refs/heads/" + trunkBranch, Commit: state.Commit,
		})
	} else if state.Channel == "nightly" && state.Branch == trunkBranch {
		// A nightly tag is created on the verified trunk commit. The trunk
		// branch is not moved, so the plan binds its current remote state.
		expectedRemote = append(expectedRemote, RemoteRefExpectation{
			Ref: "refs/heads/" + trunkBranch, Commit: state.Commit,
		})
	} else {
		return ReleaseRefPlan{}, fmt.Errorf(
			"nightly state %+v is not on trunk %q", state, trunkBranch,
		)
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Ref < refs[j].Ref })
	sort.Slice(expectedRemote, func(i, j int) bool {
		return expectedRemote[i].Ref < expectedRemote[j].Ref
	})
	return ReleaseRefPlan{
		Schema:         ReleaseRefPlanSchema,
		Kind:           "ReleaseRefPlan",
		Version:        state.Version,
		Channel:        state.Channel,
		Commit:         state.Commit,
		TreeState:      state.TreeState,
		Refs:           refs,
		ExpectedRemote: expectedRemote,
		Atomic:         true,
	}, nil
}

func validSHA256Hex(value string) bool {
	value = strings.TrimPrefix(value, "sha256:")
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func oneOf(value string, choices ...string) bool {
	for _, choice := range choices {
		if value == choice {
			return true
		}
	}
	return false
}
