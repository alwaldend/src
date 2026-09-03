package main

import (
	"context"
	"errors"
	"fmt"
)

type rebaseOptions struct{}

type rebaseReport struct {
	Status          string `json:"status"`
	Branch          string `json:"branch"`
	Base            string `json:"base"`
	BaseOID         string `json:"base_oid"`
	OldHeadOID      string `json:"old_head_oid"`
	HeadOID         string `json:"head_oid"`
	TreeOID         string `json:"tree_oid"`
	RemoteHeadLease string `json:"remote_head_lease"`
	RemoteHeadOID   string `json:"remote_head_oid"`
	Pushed          bool   `json:"pushed"`
}

func (d *delivery) rebase(
	ctx context.Context,
	options rebaseOptions,
) (result *rebaseReport, returnErr error) {
	report, err := d.inspect(ctx)
	if err != nil {
		return nil, err
	}
	if err := ensureNoRefusals(report); err != nil {
		return nil, err
	}
	if report.RemoteHeadDiverged {
		return nil, fmt.Errorf(
			"remote feature tip differs from local HEAD; use prepare to authorize a remote replacement",
		)
	}
	if report.UniqueCommitCount != 1 || len(report.MergeCommits) != 0 {
		return nil, fmt.Errorf("rebase requires exactly one non-merge feature commit")
	}
	snapshot, err := d.repository.snapshot(
		ctx,
		d.fetchEndpoint,
		report.Base,
		report.Branch,
	)
	if err != nil {
		return nil, err
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
	}
	baseIsAncestor, err := d.repository.isAncestor(
		ctx,
		snapshot.BaseOID,
		report.LocalHeadOID,
	)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	if err := d.repository.rangeDiffCheck(
		ctx,
		snapshot.BaseOID,
		headOID,
	); err != nil {
		return nil, err
	}
	if err := d.repository.verifyRequiredSignature(ctx, headOID); err != nil {
		return nil, err
	}
	currentBranch, currentHead, err := d.repository.branchHead(ctx)
	if err != nil {
		return nil, err
	}
	if currentBranch != report.Branch || currentHead != headOID {
		return nil, fmt.Errorf("current branch or HEAD changed after rebase")
	}
	status, err := d.repository.status(ctx)
	if err != nil {
		return nil, err
	}
	if !status.clean() {
		return nil, fmt.Errorf("worktree changed after rebase")
	}
	tree, err := d.repository.tree(ctx, headOID)
	if err != nil {
		return nil, err
	}
	currentPullRequest, err := d.currentOpenPullRequest(ctx, report.Branch)
	if err != nil {
		return nil, err
	}
	if (currentPullRequest == nil) != (report.PullRequest == nil) {
		return nil, fmt.Errorf("pull request presence changed after inspection")
	}
	if currentPullRequest != nil &&
		!samePullRequestMetadata(*currentPullRequest, *report.PullRequest) {
		return nil, fmt.Errorf(
			"pull request metadata changed after inspection; preserve possible human edits",
		)
	}
	if err := d.repository.push(
		ctx,
		d.pushEndpoint,
		headOID,
		snapshot.RemoteHeadRef,
		snapshot.RemoteHeadOID,
		snapshot.RemoteBaseRef,
		snapshot.BaseOID,
	); err != nil {
		return nil, err
	}
	check, err := d.repository.snapshot(
		ctx,
		d.fetchEndpoint,
		report.Base,
		report.Branch,
	)
	if err != nil {
		return nil, err
	}
	checkClosed := false
	defer func() {
		if checkClosed {
			return
		}
		if cleanupErr := check.close(context.Background()); cleanupErr != nil {
			returnErr = errors.Join(
				returnErr,
				fmt.Errorf("clean rebase verification snapshot: %w", cleanupErr),
			)
		}
	}()
	if check.BaseOID != snapshot.BaseOID {
		return nil, fmt.Errorf("base changed during rebase push; review delivery state")
	}
	if check.RemoteHeadOID != headOID {
		return nil, fmt.Errorf("post-push feature ref does not match the pushed head")
	}
	if cleanupErr := check.close(context.Background()); cleanupErr != nil {
		return nil, fmt.Errorf("clean rebase verification snapshot: %w", cleanupErr)
	}
	checkClosed = true
	result.Status = "rebased"
	result.HeadOID = headOID
	result.TreeOID = tree
	result.RemoteHeadOID = headOID
	result.Pushed = true
	return result, nil
}
