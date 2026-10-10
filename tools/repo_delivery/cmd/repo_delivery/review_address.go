package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/cobra"
)

const reviewAddressSchema = "repo_delivery/review_address/v1"

type reviewAddressOptions struct {
	ReceiptFile      string
	ThreadID         string
	BodyFile         string
	ReplyReceiptFile string
	DefectID         string
}

// This checkpoint prevents replaying a public reply after an interrupted or
// uncertain mutation. Only the separate, expiring reply receipt grants the
// existing one-use resolution authority.
type reviewAddressState struct {
	Schema        string `json:"schema"`
	HeadOID       string `json:"head_oid"`
	PullRequestID string `json:"pull_request_id"`
	ThreadID      string `json:"thread_id"`
	BodyDigest    string `json:"body_digest"`
	DefectID      string `json:"defect_id,omitempty"`
}

type reviewAddressWorkflow struct {
	delivery *delivery
	options  reviewAddressOptions
}

func newReviewAddressCommand(ctx context.Context, config *deliveryConfig,
	getenv func(string) string, stdout io.Writer, runner commandRunner,
) *cobra.Command {
	options := reviewAddressOptions{}
	command := &cobra.Command{
		Use: "address", Short: "Reply and resolve feedback on an exact validated published candidate",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return fmt.Errorf("configure review delivery: %w", err)
			}
			workflow := reviewAddressWorkflow{delivery: delivery, options: options}
			report, err := workflow.run(cmd.Context())
			if err != nil {
				return fmt.Errorf("address review thread %q: %w", options.ThreadID, err)
			}
			if err := writeJSON(stdout, report); err != nil {
				return fmt.Errorf("write addressed review state: %w", err)
			}
			return nil
		},
	}
	command.SetContext(ctx)
	flags := command.Flags()
	flags.StringVar(&options.ReceiptFile, "receipt-file", "", "published preparation receipt with recorded passing validation")
	flags.StringVar(&options.ThreadID, "thread-id", "", "review thread node ID to address")
	flags.StringVar(&options.BodyFile, "body-file", "", "reasoned reply in the receipt's ignored out/<task>/ directory")
	flags.StringVar(&options.ReplyReceiptFile, "reply-receipt-file", "", "optional durable reply receipt path; default derives from preparation receipt and thread")
	flags.StringVar(&options.DefectID, "defect-id", "", "optional stable defect identity")
	for _, flag := range []string{"receipt-file", "thread-id", "body-file"} {
		_ = command.MarkFlagRequired(flag)
	}
	return command
}

func (w *reviewAddressWorkflow) run(ctx context.Context) (report *reviewInspection, returnErr error) {
	d := w.delivery
	if err := validateOpaqueID("review thread ID", w.options.ThreadID); err != nil {
		return nil, fmt.Errorf("validate review address: %w", err)
	}
	if err := validateReviewReplyJoinRefs("", "", w.options.DefectID); err != nil {
		return nil, fmt.Errorf("validate defect identity: %w", err)
	}
	replyFile := w.options.ReplyReceiptFile
	if replyFile == "" {
		key := digestStrings("repo_delivery review address path v1", w.options.ThreadID)
		replyFile = w.options.ReceiptFile + ".review-" + key[:16] + ".json"
	}
	stateFile := replyFile + ".address.json"
	if err := d.requireSameTaskOutputs(w.options.ReceiptFile, w.options.BodyFile, replyFile, stateFile); err != nil {
		return nil, fmt.Errorf("validate review output scope: %w", err)
	}
	for _, path := range []string{replyFile, stateFile} {
		if filepath.Clean(path) == filepath.Clean(w.options.ReceiptFile) || filepath.Clean(path) == filepath.Clean(w.options.BodyFile) {
			return nil, fmt.Errorf("review workflow outputs must be distinct from preparation receipt and body")
		}
	}
	preparation, err := d.beginReceiptTransaction(ctx, w.options.ReceiptFile, true)
	if err != nil {
		return nil, fmt.Errorf("lock published candidate receipt: %w", err)
	}
	defer func() { returnErr = errors.Join(returnErr, preparation.close()) }()
	receipt, err := preparation.read()
	if err != nil {
		return nil, fmt.Errorf("read published candidate receipt: %w", err)
	}
	publication, err := w.checkPublishedCandidate(ctx, receipt)
	if err != nil {
		return nil, fmt.Errorf("check review candidate: %w", err)
	}

	body, err := d.readReviewBody(ctx, w.options.BodyFile)
	if err != nil {
		return nil, fmt.Errorf("read reasoned review reply: %w", err)
	}
	body, err = withCommentDisclaimer(body)
	if err != nil {
		return nil, fmt.Errorf("add review reply disclaimer: %w", err)
	}
	checkpoint, err := d.beginReceiptTransaction(ctx, stateFile, false)
	if err != nil {
		return nil, fmt.Errorf("lock review workflow checkpoint: %w", err)
	}
	defer func() { returnErr = errors.Join(returnErr, checkpoint.close()) }()
	inspection, err := d.inspectReviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("inspect review address target: %w", err)
	}
	if inspection.PullRequest.ID != publication.Inspection.PullRequest.ID || inspection.PullRequest.HeadRefOID != receipt.PreparedHeadOID {
		return nil, fmt.Errorf("review inventory belongs to a different published candidate")
	}
	expected := reviewAddressState{
		Schema:        reviewAddressSchema,
		HeadOID:       receipt.PreparedHeadOID,
		PullRequestID: inspection.PullRequest.ID,
		ThreadID:      w.options.ThreadID,
		BodyDigest:    reviewReplyBodyDigest(body),
		DefectID:      w.options.DefectID,
	}
	target := reviewTargetOptions{
		PullRequestID:             inspection.PullRequest.ID,
		ExpectedHeadOID:           receipt.PreparedHeadOID,
		ExpectedPullRequestDigest: inspection.PullRequestExpectationDigest,
	}
	if checkpoint.version.present {
		if err := expected.requireCheckpoint(checkpoint.version.contents); err != nil {
			return nil, fmt.Errorf("verify saved review workflow: %w", err)
		}
	} else {
		if err := w.postReply(ctx, preparation, checkpoint, inspection, target, expected, body, replyFile, stateFile); err != nil {
			return nil, fmt.Errorf("post review reply once: %w", err)
		}
	}
	reply, err := d.readReviewReplyReceipt(ctx, replyFile)
	if err != nil {
		return nil, fmt.Errorf("workflow already attempted its reply; no usable receipt remains; inspect remote state before recovery: %w", err)
	}
	if reply.ReplyBodyDigest != expected.BodyDigest || reply.DefectID != expected.DefectID {
		return nil, fmt.Errorf("durable reply receipt does not match this workflow's reply")
	}
	if err := preparation.requireCurrent(ctx); err != nil {
		return nil, fmt.Errorf("recheck review preparation before resolution: %w", err)
	}
	expectation := reviewThreadExpectation{
		ThreadID:              w.options.ThreadID,
		ExpectedLastCommentID: reply.ReplyCommentID,
		ExpectedDigest:        reply.ResultThreadDigest,
	}
	if err := d.requireValidationCandidate(ctx, receipt); err != nil {
		return nil, fmt.Errorf("recheck exact candidate before resolution: %w", err)
	}
	resolved, err := d.resolveReviewThread(ctx, target, expectation, replyFile)
	if err != nil {
		return nil, fmt.Errorf("resolve review using saved reply receipt %q: %w", replyFile, err)
	}
	return resolved, nil
}

func (w *reviewAddressWorkflow) checkPublishedCandidate(ctx context.Context, receipt preparationReceipt) (*verifyReport, error) {
	d := w.delivery

	state, err := d.readValidationState(ctx, w.options.ReceiptFile)
	if err != nil {
		return nil, fmt.Errorf("read candidate validation: %w", err)
	}
	if state.Status != "verified" && state.Status != "publication_attempted" {
		return nil, fmt.Errorf("review address requires attempted publication and passing validation; state is %q", state.Status)
	}
	if state.HeadOID != receipt.PreparedHeadOID || state.TreeOID != receipt.PreparedTreeOID || state.RepositoryFingerprint != receipt.RepositoryFingerprint {
		return nil, fmt.Errorf("validation belongs to a different prepared candidate")
	}
	if err := d.requireValidationCandidate(ctx, receipt); err != nil {
		return nil, fmt.Errorf("check exact review candidate: %w", err)
	}
	if err := state.requirePassingValidation(ctx, d, w.options.ReceiptFile, receipt); err != nil {
		return nil, fmt.Errorf("check recorded review validation: %w", err)
	}
	publication, err := d.verifyPublication(ctx, &receipt)
	if err != nil {
		return nil, fmt.Errorf("check exact published review candidate: %w", err)
	}
	if publication.Inspection.PullRequest == nil {
		return nil, fmt.Errorf("review address requires an open published pull request")
	}
	return publication, nil
}

func (expected reviewAddressState) requireCheckpoint(contents []byte) error {
	var saved reviewAddressState
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&saved); err != nil {
		return fmt.Errorf("decode review workflow checkpoint: %w", err)
	}
	if saved != expected {
		return fmt.Errorf("review workflow candidate, thread, body, or defect identity changed")
	}
	var trailing any
	if err := decoder.Decode(&trailing); err == nil {
		return fmt.Errorf("review workflow checkpoint contains trailing data")
	} else if !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode trailing review workflow data: %w", err)
	}
	return nil
}

func (w *reviewAddressWorkflow) postReply(
	ctx context.Context,
	preparation, checkpoint *receiptTransaction,
	inspection *reviewInspection,
	target reviewTargetOptions,
	expected reviewAddressState,
	body, replyFile, stateFile string,
) error {
	d := w.delivery

	// Never overwrite a receipt from another workflow or replay an unknown reply.
	absolute, err := d.receiptPath(ctx, replyFile, false)
	if err != nil {
		return fmt.Errorf("validate durable reply receipt: %w", err)
	}
	version, err := captureReceiptFileVersion(absolute)
	if err != nil {
		return fmt.Errorf("inspect durable reply receipt: %w", err)
	}
	if version.present {
		return fmt.Errorf("reply receipt already exists without this workflow's checkpoint")
	}
	thread, err := findReviewThread(inspection, w.options.ThreadID)
	if err != nil {
		return fmt.Errorf("find review address thread: %w", err)
	}
	if thread.IsResolved || thread.IsOutdated || len(thread.Comments) == 0 {
		return fmt.Errorf("review address requires a current unresolved thread with feedback")
	}
	if err := preparation.requireCurrent(ctx); err != nil {
		return fmt.Errorf("recheck review candidate receipt: %w", err)
	}
	if err := checkpoint.requireCurrent(ctx); err != nil {
		return fmt.Errorf("recheck review checkpoint: %w", err)
	}
	receipt, err := preparation.read()
	if err != nil {
		return fmt.Errorf("reread candidate receipt before reply: %w", err)
	}
	if err := d.requireValidationCandidate(ctx, receipt); err != nil {
		return fmt.Errorf("recheck exact candidate before reply: %w", err)
	}
	if err := d.writeAtomicIgnoredJSON(ctx, stateFile, "review workflow checkpoint", expected); err != nil {
		return fmt.Errorf("record reply attempt before mutation: %w", err)
	}
	priorDigest := thread.ExpectationDigest
	expectation := reviewThreadExpectation{
		ThreadID:              thread.ID,
		ExpectedLastCommentID: thread.Comments[len(thread.Comments)-1].ID,
		ExpectedDigest:        priorDigest,
	}
	replied, err := d.replyToReviewThread(ctx, target, expectation, body)
	if err != nil {
		return fmt.Errorf("reply outcome is unconfirmed; checkpoint prevents automatic reposting; inspect remote state: %w", err)
	}
	reply, err := d.newReviewReplyReceipt(ctx, target, thread.ID, priorDigest, body, replied, "", "candidate-"+expected.HeadOID, w.options.DefectID)
	if err != nil {
		return fmt.Errorf("reply succeeded but authority was not recorded; do not repost: %w", err)
	}
	if err := d.writeReviewReplyReceipt(ctx, replyFile, reply); err != nil {
		return fmt.Errorf("reply succeeded but receipt installation failed; do not repost: %w", err)
	}
	return nil
}
