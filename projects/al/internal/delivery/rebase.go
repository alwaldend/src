package delivery

import (
	"context"
	"errors"
	"fmt"
)

// rebaseOptions selects local synchronization without authorizing publication.
type rebaseOptions struct{}

// rebaseReport binds the local result to the fetched source and remote lease.
type rebaseReport struct {
	// Status distinguishes a replayed candidate from an already based branch.
	Status string `json:"status"`
	// Branch identifies the task-owned feature branch.
	Branch string `json:"branch"`
	// Base names the pull request's actual target branch.
	Base string `json:"base"`
	// BaseOID identifies the fetched base used for synchronization.
	BaseOID string `json:"base_oid"`
	// OldHeadOID preserves the pre-rebase local candidate identity.
	OldHeadOID string `json:"old_head_oid"`
	// HeadOID identifies the resulting local candidate.
	HeadOID string `json:"head_oid"`
	// TreeOID identifies the resulting local content.
	TreeOID string `json:"tree_oid"`
	// RemoteHeadLease records the fetched feature ref, empty when absent.
	RemoteHeadLease string `json:"remote_head_lease"`
	// RemoteHeadOID is the observed remote feature head, unchanged by this operation.
	RemoteHeadOID string `json:"remote_head_oid"`
	// Pushed remains false because publication requires its separate guarded workflow.
	Pushed bool `json:"pushed"`
	// NextAction describes the required continuation after a local replay.
	NextAction string `json:"next_action,omitempty"`
}

// rebase replays one owned commit locally and preserves the remote publication state.
func (d *delivery) rebase(
	ctx context.Context,
	options rebaseOptions,
) (result *rebaseReport, returnErr error) {
	report, err := d.inspect(ctx)
	if err != nil {
		return nil, fmt.Errorf("inspect rebase candidate: %w", err)
	}
	if err := ensureNoRefusals(report); err != nil {
		return nil, fmt.Errorf("check rebase eligibility: %w", err)
	}
	if report.RemoteHeadDiverged {
		return nil, fmt.Errorf(
			"remote feature tip differs from local HEAD; use prepare to authorize a remote replacement",
		)
	}
	if report.UniqueCommitCount != 1 || len(report.MergeCommits) != 0 {
		return nil, fmt.Errorf(
			"rebase requires exactly one non-merge feature commit; observed %d commits and %d merge commits; use prepare with --consolidate when explicitly authorized",
			report.UniqueCommitCount,
			len(report.MergeCommits),
		)
	}
	snapshot, err := d.repository.snapshot(
		ctx,
		d.fetchEndpoint,
		report.Base,
		report.Branch,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch rebase refs: %w", err)
	}
	defer func() {
		if cleanupErr := snapshot.close(context.Background()); cleanupErr != nil {
			returnErr = errors.Join(
				returnErr,
				fmt.Errorf("clean rebase snapshot: %w", cleanupErr),
			)
		}
	}()
	if snapshot.RemoteHeadOID != report.RemoteHeadOID {
		return nil, fmt.Errorf(
			"remote feature ref changed after inspection; retry from inspect",
		)
	}
	result = &rebaseReport{
		Branch:          report.Branch,
		Base:            report.Base,
		BaseOID:         snapshot.BaseOID,
		OldHeadOID:      report.LocalHeadOID,
		RemoteHeadLease: snapshot.RemoteHeadOID,
		RemoteHeadOID:   snapshot.RemoteHeadOID,
	}
	baseIsAncestor, err := d.repository.isAncestor(
		ctx,
		snapshot.BaseOID,
		report.LocalHeadOID,
	)
	if err != nil {
		return nil, fmt.Errorf("compare rebase ancestry: %w", err)
	}
	if baseIsAncestor {
		result.Status = "already_based"
		result.HeadOID = report.LocalHeadOID
		result.TreeOID = report.LocalTreeOID
		return result, nil
	}
	headOID, err := d.repository.rebase(
		ctx,
		snapshot.BaseOID,
		report.Branch,
		report.LocalHeadOID,
	)
	if err != nil {
		return nil, fmt.Errorf("replay feature commit: %w", err)
	}
	if err := d.repository.rangeDiffCheck(
		ctx,
		snapshot.BaseOID,
		headOID,
	); err != nil {
		return nil, fmt.Errorf("check rebased diff: %w", err)
	}
	if err := d.repository.verifyRequiredSignature(ctx, headOID); err != nil {
		return nil, fmt.Errorf("verify rebased signature: %w", err)
	}
	currentBranch, currentHead, err := d.repository.branchHead(ctx)
	if err != nil {
		return nil, fmt.Errorf("read rebased branch: %w", err)
	}
	if currentBranch != report.Branch || currentHead != headOID {
		return nil, fmt.Errorf("current branch or HEAD changed after rebase")
	}
	status, err := d.repository.status(ctx)
	if err != nil {
		return nil, fmt.Errorf("inspect rebased worktree: %w", err)
	}
	if !status.clean() {
		return nil, fmt.Errorf("worktree changed after rebase")
	}
	tree, err := d.repository.tree(ctx, headOID)
	if err != nil {
		return nil, fmt.Errorf("read rebased tree: %w", err)
	}
	result.Status = "rebased"
	result.HeadOID = headOID
	result.TreeOID = tree
	result.NextAction = "prepare the rebased candidate, validate it, then publish explicitly"
	return result, nil
}
