package versioning

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
	// leasePath is the exclusively created operation marker owned by this publisher.
	leasePath string
	// leaseID identifies the current local operation.
	leaseID string
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

// AcquireLease atomically creates an operation marker shared by linked worktrees.
// The marker coordinates local publishers; Git serializes remote ref updates.
func (p *LocalGitReleasePublisher) AcquireLease(
	ctx context.Context,
	plan ReleaseRefPlan,
) (string, error) {
	if p.leasePath != "" {
		return "", fmt.Errorf("publisher already holds a lease")
	}
	common, err := p.Git.Run("", "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("resolve operation lock directory: %w", err)
	}
	if !filepath.IsAbs(common) {
		return "", fmt.Errorf("Git common directory must be absolute")
	}
	marker := filepath.Join(common, "release-refs.lock")
	if err := os.Mkdir(marker, 0o700); err != nil {
		return "", fmt.Errorf("acquire release-refs operation lock %q: %w", marker, err)
	}
	p.leasePath, p.leaseID = marker, "release-refs:"+plan.Commit
	return p.leaseID, nil
}

// ReleaseLease removes only this publisher's exclusively acquired operation marker.
func (p *LocalGitReleasePublisher) ReleaseLease(ctx context.Context, lease string) error {
	if p.leasePath == "" || lease != p.leaseID {
		return fmt.Errorf("release-refs lease is not held by this publisher")
	}
	if err := os.Remove(p.leasePath); err != nil {
		return fmt.Errorf("remove operation lock: %w", err)
	}
	p.leasePath, p.leaseID = "", ""
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
	if lease != p.leaseID || p.leasePath == "" {
		return fmt.Errorf("publication requires an acquired operation lock")
	}
	observed, err := p.FetchExpectedRemoteState(ctx, plan)
	if err != nil {
		return fmt.Errorf("refresh remote under lock: %w", err)
	}
	if err := requireRemotePreconditions(plan, observed); err != nil {
		return fmt.Errorf("verify locked preconditions: %w", err)
	}
	present := map[string]bool{}
	for _, ref := range observed {
		present[ref.Ref] = true
	}
	args := []string{"push", p.Remote}
	leases := []string{}
	for _, ref := range plan.Refs {
		if ref.Operation != "create" || present[ref.Ref] {
			continue
		}
		args = append(args, ref.Expected+":"+ref.Ref)
		leases = append(leases, "--force-with-lease="+ref.Ref+":")
	}
	if len(leases) == 0 {
		return nil
	}
	if plan.Atomic {
		args = append(args, "--atomic")
	}
	args = append(args, leases...)
	if _, err := p.Git.Run("", args...); err != nil {
		message := strings.ToLower(err.Error())
		if plan.Atomic && (strings.Contains(message, "does not support --atomic push") || strings.Contains(message, "atomic push not supported") || strings.Contains(message, "atomic push is not supported")) {
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
