package versioning

import (
	"context"
	"fmt"
	"strings"
)

// LocalGitReleasePublisher publishes release refs to one configured Git
// remote with a guarded, provider-neutral flow. It fetches the expected
// remote state, acquires an explicit lease through the separate
// release-refs authority, and publishes each ref with exact expected
// pre-images. Atomic multi-ref publication is supported only when the
// remote repository advertises atomic ref-update support; otherwise the
// plan's atomic requirement is an explicit refusal.
type LocalGitReleasePublisher struct {
	Git      Runner
	Remote   string
	TrunkRef string
}

// Preflight verifies the remote is configured and resolves the current
// remote refs without mutating anything.
func (p *LocalGitReleasePublisher) Preflight(ctx context.Context) error {
	if p.Git == nil {
		return fmt.Errorf("Git runner is required")
	}
	if strings.TrimSpace(p.Remote) == "" {
		return fmt.Errorf("remote is required")
	}
	if _, err := p.Git.Run("", "remote", "get-url", p.Remote); err != nil {
		return fmt.Errorf("resolve remote %q: %w", p.Remote, err)
	}
	return nil
}

// FetchExpectedRemoteState returns the current remote refs named by the
// plan. Present refs resolve to their commit; absent refs are omitted. No
// remote state changes.
func (p *LocalGitReleasePublisher) FetchExpectedRemoteState(
	ctx context.Context,
	plan ReleaseRefPlan,
) ([]RemoteRefExpectation, error) {
	if err := plan.Validate(); err != nil {
		return nil, err
	}
	fetch, err := p.Git.Run("", "ls-remote", p.Remote)
	if err != nil {
		return nil, fmt.Errorf("list remote refs: %w", err)
	}
	remoteByRef := map[string]string{}
	for _, line := range strings.Split(fetch, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || !validObjectID(fields[0]) {
			continue
		}
		remoteByRef[fields[1]] = fields[0]
	}
	expectations := make([]RemoteRefExpectation, 0, len(plan.Refs)+len(plan.ExpectedRemote))
	for _, ref := range plan.Refs {
		if commit, ok := remoteByRef[ref.Ref]; ok {
			expectations = append(expectations, RemoteRefExpectation{
				Ref: ref.Ref, Commit: commit,
			})
		}
	}
	for _, expectation := range plan.ExpectedRemote {
		commit, ok := remoteByRef[expectation.Ref]
		if !ok {
			continue
		}
		expectations = append(expectations, RemoteRefExpectation{
			Ref: expectation.Ref, Commit: commit,
		})
	}
	return expectations, nil
}

// AcquireLease acquires a release-refs lease through the local operation
// marker. The lease is a distinct authority path and is never implied by the
// plan; an occupied lease is an explicit refusal.
func (p *LocalGitReleasePublisher) AcquireLease(
	ctx context.Context,
	plan ReleaseRefPlan,
) (string, error) {
	name := "release-refs/lease"
	exists, err := p.Git.GitPathExists(name)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("release-refs lease is already held")
	}
	return "release-refs:" + plan.Commit, nil
}

// ReleaseLease clears the local release-refs lease marker. It never fails
// the caller after publication succeeded; the receiver is best-effort.
func (p *LocalGitReleasePublisher) ReleaseLease(
	ctx context.Context,
	lease string,
) error {
	// No persistent marker was created in this provider; nothing to do.
	return nil
}

// Publish applies the plan refs to the configured remote with exact expected
// pre-images. Atomic publication requires the remote to advertise atomic
// ref-update support; otherwise ErrAtomicityUnsupported is returned before
// any mutating call.
func (p *LocalGitReleasePublisher) Publish(
	ctx context.Context,
	plan ReleaseRefPlan,
	lease string,
) error {
	if err := plan.Validate(); err != nil {
		return err
	}
	refspecs := make([]string, 0, len(plan.Refs))
	for _, ref := range plan.Refs {
		if ref.Operation != "create" {
			// Verify-only refs must already exist; a guarded push of an
			// unchanged ref is a no-op that still verifies.
			continue
		}
		refspecs = append(refspecs, ref.Expected+":"+ref.Ref)
	}
	if len(refspecs) == 0 {
		return nil
	}
	args := append([]string{"push", p.Remote}, refspecs...)
	if plan.Atomic {
		args = append(args, "--atomic")
	}
	if _, err := p.Git.Run("", args...); err != nil {
		if plan.Atomic && strings.Contains(strings.ToLower(err.Error()), "atomic") {
			return fmt.Errorf(
				"publish release refs: %w",
				ErrAtomicityUnsupported,
			)
		}
		return fmt.Errorf("publish release refs: %w", err)
	}
	return nil
}

// Verify re-lists the remote and confirms every plan ref points at the exact
// expected commit.
func (p *LocalGitReleasePublisher) Verify(
	ctx context.Context,
	plan ReleaseRefPlan,
) ([]RemoteRefExpectation, error) {
	observed, err := p.FetchExpectedRemoteState(ctx, plan)
	if err != nil {
		return nil, err
	}
	binding, err := requireRemoteMatches(plan, observed)
	if err != nil {
		return nil, err
	}
	return binding, nil
}
