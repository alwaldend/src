package main

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	preparationReceiptSchema   = "repo_delivery/preparation/v1"
	rewriteAuthorizationSchema = "repo_delivery/rewrite_authorization/v1"
	receiptFileLimit           = 256 * 1024
	receiptRevisionBytes       = 32
	scopeModePaths             = "paths"
	scopeModeUseIndex          = "use_index"
	scopeModeMessageOnly       = "message_only"
)

type refExpectation struct {
	Ref     string `json:"ref"`
	Present bool   `json:"present"`
	OID     string `json:"oid"`
}

// pullRequestExpectation binds either the explicit absence of a pull request
// or the stable identity and pre-publication metadata of the selected pull
// request. Bodies are represented only by a domain-separated digest.
type pullRequestExpectation struct {
	Present               bool   `json:"present"`
	ID                    string `json:"id"`
	Number                int    `json:"number"`
	IdentityDigest        string `json:"identity_sha256"`
	AuthorLogin           string `json:"author_login"`
	IsDraft               bool   `json:"is_draft"`
	PriorHeadOID          string `json:"prior_head_oid"`
	PriorProjectionDigest string `json:"prior_projection_sha256"`
}

type aggregateScope struct {
	Mode            string   `json:"mode"`
	AuthorizedPaths []string `json:"authorized_paths"`
	AggregatePaths  []string `json:"aggregate_paths"`
}

// rewriteAuthorization is a typed, non-authorizing synchronization record
// for an explicitly authorized remote feature-ref replacement. It proves the
// tool observed an exact old remote OID and an exact replacement head OID
// within one provider-owned repository, alongside the task scope and any
// source receipts that established ownership. It never contains credentials
// or remote endpoints; the far ends are identified by digest.
type rewriteAuthorization struct {
	Schema                string           `json:"schema"`
	IssuedAt              string           `json:"issued_at"`
	RepositoryFingerprint string           `json:"repository_fingerprint"`
	RemoteRepository      remoteRepository `json:"remote_repository"`
	Provider              string           `json:"provider"`
	OldRemoteOID          string           `json:"old_remote_oid"`
	NewHeadOID            string           `json:"new_head_oid"`
	OwnerRoot             string           `json:"owner_root"`
	TaskPaths             []string         `json:"task_paths"`
	SourceReceiptDigest   string           `json:"source_receipt_digest,omitempty"`
}

func (a rewriteAuthorization) validate() error {
	if a.Schema != rewriteAuthorizationSchema {
		return fmt.Errorf(
			"unsupported rewrite authorization schema %q",
			a.Schema,
		)
	}
	if _, err := parseCanonicalReviewReceiptTime(
		"issued_at",
		a.IssuedAt,
	); err != nil {
		return err
	}
	if !validSHA256Digest(a.RepositoryFingerprint) {
		return fmt.Errorf(
			"rewrite authorization has an invalid repository fingerprint",
		)
	}
	if a.RemoteRepository.Host == "" || a.RemoteRepository.Owner == "" ||
		a.RemoteRepository.Name == "" {
		return fmt.Errorf(
			"rewrite authorization has incomplete provider identity",
		)
	}
	if a.Provider == "" {
		return fmt.Errorf("rewrite authorization has an empty provider")
	}
	if !isObjectID(a.OldRemoteOID) || !isObjectID(a.NewHeadOID) {
		return fmt.Errorf(
			"rewrite authorization must identify exact old and new OIDs",
		)
	}
	if strings.TrimSpace(a.OwnerRoot) == "" ||
		len(a.OwnerRoot) > 253 {
		return fmt.Errorf(
			"rewrite authorization has an invalid owner root",
		)
	}
	if len(a.TaskPaths) == 0 || !sortedUniqueStrings(a.TaskPaths) {
		return fmt.Errorf(
			"rewrite authorization task paths must be sorted and unique",
		)
	}
	for _, path := range a.TaskPaths {
		validated, err := validateTaskPaths([]string{path})
		if err != nil || len(validated) != 1 || validated[0] != path {
			return fmt.Errorf(
				"rewrite authorization contains an invalid task path %q",
				path,
			)
		}
	}
	if a.SourceReceiptDigest != "" &&
		!validSHA256Digest(a.SourceReceiptDigest) {
		return fmt.Errorf(
			"rewrite authorization has an invalid source receipt digest",
		)
	}
	return nil
}

type preparationReceipt struct {
	Schema                string                 `json:"schema"`
	RevisionNonce         string                 `json:"revision_nonce"`
	RepositoryFingerprint string                 `json:"repository_fingerprint"`
	RemoteName            string                 `json:"remote_name"`
	FetchEndpointDigest   string                 `json:"fetch_endpoint_sha256"`
	PushEndpointDigest    string                 `json:"push_endpoint_sha256"`
	RemoteRepository      remoteRepository       `json:"remote_repository"`
	Forge                 string                 `json:"forge"`
	BaseRef               string                 `json:"base_ref"`
	BaseOID               string                 `json:"base_oid"`
	HeadRef               string                 `json:"head_ref"`
	PreparedHeadOID       string                 `json:"prepared_head_oid"`
	PreparedTreeOID       string                 `json:"prepared_tree_oid"`
	ExpectedRemoteHead    refExpectation         `json:"expected_remote_head"`
	ExpectedPullRequest   pullRequestExpectation `json:"expected_pull_request"`
	Scope                 aggregateScope         `json:"scope"`
	RewriteAuthorization  *rewriteAuthorization  `json:"rewrite_authorization,omitempty"`
}

type receiptPullRequestState string

const (
	receiptPullRequestAbsent                        receiptPullRequestState = "absent"
	receiptPullRequestHeadPushed                    receiptPullRequestState = "head_pushed_without_pull_request"
	receiptPullRequestPrior                         receiptPullRequestState = "prior"
	receiptPullRequestHeadPushedWithPriorProjection receiptPullRequestState = "head_pushed_with_prior_projection"
	receiptPullRequestDesired                       receiptPullRequestState = "desired"
	receiptPullRequestCreatedDesired                receiptPullRequestState = "created_desired"
)

func endpointDigest(kind string, endpoint string) string {
	sum := sha256.Sum256([]byte(
		"repo_delivery endpoint v1\x00" + kind + "\x00" + endpoint,
	))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func digestStrings(domain string, values ...string) string {
	hash := sha256.New()
	_, _ = io.WriteString(hash, domain)
	for _, value := range values {
		_, _ = io.WriteString(hash, "\x00")
		_, _ = io.WriteString(hash, strconv.Itoa(len(value)))
		_, _ = io.WriteString(hash, ":")
		_, _ = io.WriteString(hash, value)
	}
	return "sha256:" + hex.EncodeToString(hash.Sum(nil))
}

func newReceiptRevisionNonce() (string, error) {
	value := make([]byte, receiptRevisionBytes)
	if _, err := io.ReadFull(cryptorand.Reader, value); err != nil {
		return "", fmt.Errorf("generate preparation receipt revision nonce: %w", err)
	}
	return hex.EncodeToString(value), nil
}

func validReceiptRevisionNonce(value string) bool {
	if len(value) != receiptRevisionBytes*2 || strings.ToLower(value) != value {
		return false
	}
	decoded, err := hex.DecodeString(value)
	return err == nil && len(decoded) == receiptRevisionBytes
}

func pullRequestIdentityDigest(value pullRequest) string {
	return digestStrings(
		"repo_delivery pull request identity v1",
		value.ID,
		strconv.Itoa(value.Number),
		value.URL,
		strings.ToUpper(value.State),
		value.BaseRefName,
		value.HeadRefName,
		strings.ToLower(value.HeadRepositoryOwner),
		strings.ToLower(value.HeadRepositoryName),
		strconv.FormatBool(value.IsCrossRepository),
	)
}

func pullRequestProjectionDigest(title string, body string) string {
	return digestStrings(
		"repo_delivery pull request projection v1",
		normalizeText(title),
		normalizeText(body),
	)
}

func newPullRequestExpectation(
	value *pullRequest,
) (pullRequestExpectation, error) {
	if value == nil {
		return pullRequestExpectation{}, nil
	}
	if err := validatePullRequest(*value); err != nil {
		return pullRequestExpectation{}, fmt.Errorf(
			"build pull-request expectation: %w",
			err,
		)
	}
	if !strings.EqualFold(value.State, "OPEN") {
		return pullRequestExpectation{}, fmt.Errorf(
			"build pull-request expectation: pull request is not open",
		)
	}
	expectation := pullRequestExpectation{
		Present:        true,
		ID:             value.ID,
		Number:         value.Number,
		IdentityDigest: pullRequestIdentityDigest(*value),
		AuthorLogin:    value.AuthorLogin,
		IsDraft:        value.IsDraft,
		PriorHeadOID:   value.HeadRefOID,
		PriorProjectionDigest: pullRequestProjectionDigest(
			value.Title,
			value.Body,
		),
	}
	if err := expectation.validate(); err != nil {
		return pullRequestExpectation{}, err
	}
	return expectation, nil
}

func (e pullRequestExpectation) validate() error {
	if !e.Present {
		if e.ID != "" || e.Number != 0 || e.IdentityDigest != "" ||
			e.AuthorLogin != "" || e.IsDraft ||
			e.PriorProjectionDigest != "" {
			return fmt.Errorf(
				"absent expected pull request contains identity or metadata",
			)
		}
		if e.PriorHeadOID != "" && !isObjectID(e.PriorHeadOID) {
			return fmt.Errorf(
				"absent expected pull request has an invalid original remote head",
			)
		}
		return nil
	}
	if e.ID == "" || e.Number <= 0 {
		return fmt.Errorf("expected pull request has an invalid stable identity")
	}
	if !validSHA256Digest(e.IdentityDigest) {
		return fmt.Errorf("expected pull request has an invalid identity digest")
	}
	if !isObjectID(e.PriorHeadOID) {
		return fmt.Errorf("expected pull request has an invalid prior head OID")
	}
	if !validSHA256Digest(e.PriorProjectionDigest) {
		return fmt.Errorf("expected pull request has an invalid projection digest")
	}
	return nil
}

func (e pullRequestExpectation) matchesInvariant(value pullRequest) bool {
	return e.Present &&
		e.ID == value.ID &&
		e.Number == value.Number &&
		e.IdentityDigest == pullRequestIdentityDigest(value) &&
		e.AuthorLogin == value.AuthorLogin &&
		e.IsDraft == value.IsDraft
}

func (g *gitRepository) repositoryFingerprint(
	ctx context.Context,
) (string, error) {
	root, err := canonicalExistingPath(g.directory)
	if err != nil {
		return "", fmt.Errorf("canonicalize repository root: %w", err)
	}
	gitDirectory, err := g.text(ctx, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return "", fmt.Errorf("resolve absolute Git directory: %w", err)
	}
	gitDirectory, err = canonicalExistingPath(gitDirectory)
	if err != nil {
		return "", fmt.Errorf("canonicalize Git directory: %w", err)
	}
	commonDirectory, err := g.text(ctx, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("resolve Git common directory: %w", err)
	}
	if !filepath.IsAbs(commonDirectory) {
		commonDirectory = filepath.Join(g.directory, commonDirectory)
	}
	commonDirectory, err = canonicalExistingPath(commonDirectory)
	if err != nil {
		return "", fmt.Errorf("canonicalize Git common directory: %w", err)
	}
	objectFormat, err := g.text(ctx, "rev-parse", "--show-object-format")
	if err != nil {
		return "", fmt.Errorf("resolve repository object format: %w", err)
	}
	sum := sha256.Sum256([]byte(strings.Join([]string{
		"repo_delivery repository v1",
		root,
		gitDirectory,
		commonDirectory,
		objectFormat,
	}, "\x00")))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func canonicalExistingPath(value string) (string, error) {
	absolute, err := filepath.Abs(value)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func (d *delivery) newPreparationReceipt(
	ctx context.Context,
	baseBranch string,
	baseOID string,
	branch string,
	headOID string,
	treeOID string,
	remoteHeadOID string,
	pullRequest *pullRequest,
	scope aggregateScope,
) (preparationReceipt, error) {
	expectation, err := newPullRequestExpectation(pullRequest)
	if err != nil {
		return preparationReceipt{}, err
	}
	if pullRequest == nil {
		expectation.PriorHeadOID = remoteHeadOID
	}
	if pullRequest != nil {
		if remoteHeadOID == "" || pullRequest.HeadRefOID != remoteHeadOID {
			return preparationReceipt{}, fmt.Errorf(
				"selected pull request does not identify the expected remote head",
			)
		}
		if pullRequest.BaseRefName != baseBranch ||
			pullRequest.HeadRefName != branch ||
			pullRequest.IsCrossRepository ||
			!strings.EqualFold(
				pullRequest.HeadRepositoryOwner,
				d.remoteRepository.Owner,
			) ||
			!strings.EqualFold(
				pullRequest.HeadRepositoryName,
				d.remoteRepository.Name,
			) {
			return preparationReceipt{}, fmt.Errorf(
				"selected pull request does not match the prepared repository and refs",
			)
		}
	}
	return d.newPreparationReceiptWithExpectation(
		ctx,
		baseBranch,
		baseOID,
		branch,
		headOID,
		treeOID,
		remoteHeadOID,
		expectation,
		scope,
	)
}

// newPreparationReceiptWithExpectation is used when a base rebase derives a
// new receipt. The original pull-request expectation must be carried forward;
// recapturing the forge after a partial publication would bless that partial
// or replacement state.
func (d *delivery) newPreparationReceiptWithExpectation(
	ctx context.Context,
	baseBranch string,
	baseOID string,
	branch string,
	headOID string,
	treeOID string,
	remoteHeadOID string,
	expectedPullRequest pullRequestExpectation,
	scope aggregateScope,
) (preparationReceipt, error) {
	revisionNonce, err := newReceiptRevisionNonce()
	if err != nil {
		return preparationReceipt{}, err
	}
	fingerprint, err := d.repository.repositoryFingerprint(ctx)
	if err != nil {
		return preparationReceipt{}, err
	}
	receipt := preparationReceipt{
		Schema:                preparationReceiptSchema,
		RevisionNonce:         revisionNonce,
		RepositoryFingerprint: fingerprint,
		RemoteName:            d.remote,
		FetchEndpointDigest: endpointDigest(
			"fetch",
			d.fetchEndpoint,
		),
		PushEndpointDigest: endpointDigest(
			"push",
			d.pushEndpoint,
		),
		RemoteRepository: d.remoteRepository,
		Forge:            d.forge.Name(),
		BaseRef:          "refs/heads/" + baseBranch,
		BaseOID:          baseOID,
		HeadRef:          "refs/heads/" + branch,
		PreparedHeadOID:  headOID,
		PreparedTreeOID:  treeOID,
		ExpectedRemoteHead: refExpectation{
			Ref:     "refs/heads/" + branch,
			Present: remoteHeadOID != "",
			OID:     remoteHeadOID,
		},
		ExpectedPullRequest: expectedPullRequest,
		Scope:               scope,
	}
	if err := receipt.validate(); err != nil {
		return preparationReceipt{}, err
	}
	return receipt, nil
}

func (r preparationReceipt) validate() error {
	if r.Schema != preparationReceiptSchema {
		return fmt.Errorf("unsupported preparation receipt schema %q", r.Schema)
	}
	if !validReceiptRevisionNonce(r.RevisionNonce) {
		return fmt.Errorf("preparation receipt has an invalid revision nonce")
	}
	for name, value := range map[string]string{
		"repository fingerprint": r.RepositoryFingerprint,
		"fetch endpoint digest":  r.FetchEndpointDigest,
		"push endpoint digest":   r.PushEndpointDigest,
	} {
		if !validSHA256Digest(value) {
			return fmt.Errorf("preparation receipt has invalid %s", name)
		}
	}
	if r.RemoteName == "" || strings.HasPrefix(r.RemoteName, "-") ||
		strings.ContainsRune(r.RemoteName, '\x00') {
		return fmt.Errorf("preparation receipt has an invalid remote name")
	}
	if r.RemoteRepository.Host == "" || r.RemoteRepository.Owner == "" ||
		r.RemoteRepository.Name == "" {
		return fmt.Errorf("preparation receipt has an incomplete repository identity")
	}
	if r.Forge == "" {
		return fmt.Errorf("preparation receipt has an empty forge adapter")
	}
	if err := validateFullBranchRef(r.BaseRef); err != nil {
		return fmt.Errorf("invalid receipt base ref: %w", err)
	}
	if err := validateFullBranchRef(r.HeadRef); err != nil {
		return fmt.Errorf("invalid receipt head ref: %w", err)
	}
	if r.BaseRef == r.HeadRef {
		return fmt.Errorf("preparation receipt base and head refs are equal")
	}
	for name, oid := range map[string]string{
		"base":          r.BaseOID,
		"prepared head": r.PreparedHeadOID,
		"prepared tree": r.PreparedTreeOID,
	} {
		if !isObjectID(oid) {
			return fmt.Errorf("preparation receipt has invalid %s OID", name)
		}
	}
	if r.ExpectedRemoteHead.Ref != r.HeadRef {
		return fmt.Errorf("expected remote head ref differs from the prepared head ref")
	}
	if r.ExpectedRemoteHead.Present {
		if !isObjectID(r.ExpectedRemoteHead.OID) {
			return fmt.Errorf("expected remote head is present without a valid OID")
		}
	} else if r.ExpectedRemoteHead.OID != "" {
		return fmt.Errorf("absent expected remote head has an OID")
	}
	if err := r.ExpectedPullRequest.validate(); err != nil {
		return fmt.Errorf("invalid pull-request expectation: %w", err)
	}
	if r.RewriteAuthorization != nil {
		if err := r.RewriteAuthorization.validate(); err != nil {
			return fmt.Errorf("invalid rewrite authorization: %w", err)
		}
		if !validSHA256Digest(r.RepositoryFingerprint) ||
			r.RewriteAuthorization.RepositoryFingerprint !=
				r.RepositoryFingerprint {
			return fmt.Errorf(
				"rewrite authorization belongs to a different repository",
			)
		}
		if !sameRemoteRepository(
			r.RemoteRepository,
			r.RewriteAuthorization.RemoteRepository,
		) {
			return fmt.Errorf(
				"rewrite authorization belongs to a different provider repository",
			)
		}
	}
	return r.Scope.validate()
}

func validSHA256Digest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != 71 {
		return false
	}
	_, err := hex.DecodeString(strings.TrimPrefix(value, "sha256:"))
	return err == nil
}

func validateFullBranchRef(value string) error {
	if !strings.HasPrefix(value, "refs/heads/") {
		return fmt.Errorf("%q is not a full branch ref", value)
	}
	branch := strings.TrimPrefix(value, "refs/heads/")
	if branch == "" || strings.HasPrefix(branch, "-") ||
		strings.ContainsRune(branch, '\x00') {
		return fmt.Errorf("%q is not a safe branch ref", value)
	}
	return nil
}

func (s aggregateScope) validate() error {
	switch s.Mode {
	case scopeModePaths, scopeModeUseIndex, scopeModeMessageOnly:
	default:
		return fmt.Errorf("invalid aggregate scope mode %q", s.Mode)
	}
	if len(s.AuthorizedPaths) == 0 || len(s.AggregatePaths) == 0 {
		return fmt.Errorf("aggregate scope paths must not be empty")
	}
	for name, values := range map[string][]string{
		"authorized": s.AuthorizedPaths,
		"aggregate":  s.AggregatePaths,
	} {
		if !sortedUniqueStrings(values) {
			return fmt.Errorf("aggregate scope %s paths are not sorted and unique", name)
		}
		for _, value := range values {
			validated, err := validateTaskPaths([]string{value})
			if err != nil || len(validated) != 1 || validated[0] != value {
				return fmt.Errorf("aggregate scope contains invalid %s path %q", name, value)
			}
		}
	}
	if err := refusePathsOutside(
		s.AggregatePaths,
		s.AuthorizedPaths,
		"aggregate",
	); err != nil {
		return err
	}
	if s.Mode != scopeModePaths &&
		!reflect.DeepEqual(s.AuthorizedPaths, s.AggregatePaths) {
		return fmt.Errorf("frozen aggregate scope differs from its exact paths")
	}
	return nil
}

func sortedUniqueStrings(values []string) bool {
	for index, value := range values {
		if index != 0 && value <= values[index-1] {
			return false
		}
	}
	return true
}

func newAggregateScope(
	mode string,
	authorized []string,
	aggregate []string,
) (aggregateScope, error) {
	authorized = append([]string(nil), authorized...)
	aggregate = append([]string(nil), aggregate...)
	sort.Strings(authorized)
	sort.Strings(aggregate)
	if mode != scopeModePaths {
		authorized = append([]string(nil), aggregate...)
	}
	scope := aggregateScope{
		Mode:            mode,
		AuthorizedPaths: authorized,
		AggregatePaths:  aggregate,
	}
	if err := scope.validate(); err != nil {
		return aggregateScope{}, err
	}
	return scope, nil
}

func (d *delivery) buildAggregateScope(
	ctx context.Context,
	parentOID string,
	treeOID string,
	mode string,
	authorized []string,
) (aggregateScope, error) {
	if err := d.repository.rangeDiffCheck(ctx, parentOID, treeOID); err != nil {
		return aggregateScope{}, err
	}
	paths, err := d.repository.changedPaths(ctx, parentOID, treeOID)
	if err != nil {
		return aggregateScope{}, err
	}
	if err := refuseScratchPaths(paths); err != nil {
		return aggregateScope{}, err
	}
	return newAggregateScope(mode, authorized, paths)
}

func (d *delivery) verifyAggregateScope(
	ctx context.Context,
	baseOID string,
	headOID string,
	scope aggregateScope,
) error {
	return verifyAggregateScopeInRepository(
		ctx,
		d.repository,
		baseOID,
		headOID,
		scope,
	)
}

func verifyAggregateScopeInRepository(
	ctx context.Context,
	repository *gitRepository,
	baseOID string,
	headOID string,
	scope aggregateScope,
) error {
	if err := scope.validate(); err != nil {
		return err
	}
	if err := repository.rangeDiffCheck(ctx, baseOID, headOID); err != nil {
		return err
	}
	paths, err := repository.changedPaths(ctx, baseOID, headOID)
	if err != nil {
		return err
	}
	sort.Strings(paths)
	if !reflect.DeepEqual(paths, scope.AggregatePaths) {
		return fmt.Errorf(
			"aggregate paths changed: got %q, expected %q",
			paths,
			scope.AggregatePaths,
		)
	}
	return refusePathsOutside(paths, scope.AuthorizedPaths, "aggregate")
}

func reconcileRebasedAggregateScope(
	ctx context.Context,
	repository *gitRepository,
	baseOID string,
	headOID string,
	priorHeadOID string,
	scope aggregateScope,
) (aggregateScope, error) {
	if err := scope.validate(); err != nil {
		return aggregateScope{}, err
	}
	if err := repository.rangeDiffCheck(ctx, baseOID, headOID); err != nil {
		return aggregateScope{}, err
	}
	paths, err := repository.changedPaths(ctx, baseOID, headOID)
	if err != nil {
		return aggregateScope{}, err
	}
	expected := make(map[string]bool, len(scope.AggregatePaths))
	for _, path := range scope.AggregatePaths {
		expected[path] = true
	}
	observed := make(map[string]bool, len(paths))
	for _, path := range paths {
		if !expected[path] {
			return aggregateScope{}, fmt.Errorf(
				"rebased aggregate introduced path %q",
				path,
			)
		}
		observed[path] = true
	}
	for _, path := range scope.AggregatePaths {
		if observed[path] {
			continue
		}
		priorEntry, err := repository.pathEntry(ctx, priorHeadOID, path)
		if err != nil {
			return aggregateScope{}, err
		}
		baseEntry, err := repository.pathEntry(ctx, baseOID, path)
		if err != nil {
			return aggregateScope{}, err
		}
		if priorEntry != baseEntry {
			return aggregateScope{}, fmt.Errorf(
				"rebased aggregate lost non-identical path %q",
				path,
			)
		}
	}
	if len(paths) == 0 {
		return aggregateScope{}, fmt.Errorf("rebased aggregate path scope is empty")
	}
	if err := refusePathsOutside(paths, scope.AuthorizedPaths, "aggregate"); err != nil {
		return aggregateScope{}, err
	}
	scope.AggregatePaths = paths
	if err := scope.validate(); err != nil {
		return aggregateScope{}, err
	}
	return scope, nil
}

func (d *delivery) validateReceiptContext(
	ctx context.Context,
	receipt preparationReceipt,
	report *inspection,
) error {
	if err := receipt.validate(); err != nil {
		return err
	}
	fingerprint, err := d.repository.repositoryFingerprint(ctx)
	if err != nil {
		return err
	}
	wantBaseRef := "refs/heads/" + report.Base
	wantHeadRef := report.HeadRef
	comparisons := []struct {
		name string
		got  string
		want string
	}{
		{"repository fingerprint", receipt.RepositoryFingerprint, fingerprint},
		{"remote name", receipt.RemoteName, d.remote},
		{"fetch endpoint", receipt.FetchEndpointDigest, endpointDigest("fetch", d.fetchEndpoint)},
		{"push endpoint", receipt.PushEndpointDigest, endpointDigest("push", d.pushEndpoint)},
		{"forge adapter", receipt.Forge, d.forge.Name()},
		{"base ref", receipt.BaseRef, wantBaseRef},
		{"head ref", receipt.HeadRef, wantHeadRef},
		{"prepared head", receipt.PreparedHeadOID, report.LocalHeadOID},
		{"prepared tree", receipt.PreparedTreeOID, report.LocalTreeOID},
	}
	for _, comparison := range comparisons {
		if comparison.got != comparison.want {
			if comparison.name == "base ref" {
				return fmt.Errorf(
					"preparation receipt base ref %q differs from current delivery state %q; "+
						"refusal reason: missing_base (no pull request exists and no explicit --base was supplied to publish)",
					comparison.got,
					comparison.want,
				)
			}
			return fmt.Errorf(
				"preparation receipt %s differs from current delivery state",
				comparison.name,
			)
		}
	}
	if !sameRemoteRepository(receipt.RemoteRepository, d.remoteRepository) {
		return fmt.Errorf("preparation receipt identifies a different remote repository")
	}
	parents, err := d.repository.commitParents(ctx, receipt.PreparedHeadOID)
	if err != nil {
		return err
	}
	if len(parents) != 1 || parents[0] != receipt.BaseOID {
		return fmt.Errorf("prepared commit parent differs from the receipt base OID")
	}
	if err := d.verifyAggregateScope(
		ctx,
		receipt.BaseOID,
		receipt.PreparedHeadOID,
		receipt.Scope,
	); err != nil {
		return fmt.Errorf("verify receipt aggregate scope: %w", err)
	}
	return nil
}

func receiptAllowsRemoteHead(
	receipt preparationReceipt,
	remoteHeadOID string,
) bool {
	if remoteHeadOID == receipt.PreparedHeadOID {
		return true
	}
	if receipt.ExpectedRemoteHead.Present {
		return remoteHeadOID == receipt.ExpectedRemoteHead.OID
	}
	return remoteHeadOID == ""
}

func (r preparationReceipt) pullRequestState(
	current *pullRequest,
	remoteHeadOID string,
	desired commitProjection,
) (receiptPullRequestState, error) {
	if err := r.validate(); err != nil {
		return "", err
	}
	if !receiptAllowsRemoteHead(r, remoteHeadOID) {
		return "", fmt.Errorf(
			"remote head differs from the preparation receipt lease",
		)
	}
	expectation := r.ExpectedPullRequest
	if current == nil {
		if expectation.Present {
			return "", fmt.Errorf(
				"the expected pull request is absent; it may have been closed or replaced",
			)
		}
		if remoteHeadOID == r.PreparedHeadOID {
			return receiptPullRequestHeadPushed, nil
		}
		return receiptPullRequestAbsent, nil
	}
	if err := validatePullRequest(*current); err != nil {
		return "", fmt.Errorf("validate current pull request: %w", err)
	}
	if !strings.EqualFold(current.State, "OPEN") ||
		!r.pullRequestMatchesTopology(*current) ||
		current.HeadRefOID != remoteHeadOID {
		return "", fmt.Errorf(
			"current pull request does not match the receipt repository, refs, and remote head",
		)
	}
	desiredDigest := pullRequestProjectionDigest(desired.Title, desired.Body)
	currentDigest := pullRequestProjectionDigest(current.Title, current.Body)
	if !expectation.Present {
		partialTransition := remoteHeadOID == r.PreparedHeadOID ||
			(r.ExpectedRemoteHead.Present &&
				remoteHeadOID == r.ExpectedRemoteHead.OID &&
				remoteHeadOID != expectation.PriorHeadOID)
		if !partialTransition || current.IsDraft ||
			currentDigest != desiredDigest {
			return "", fmt.Errorf(
				"a pull request appeared after explicit absence and is not the exact desired recovery state",
			)
		}
		return receiptPullRequestCreatedDesired, nil
	}
	if !expectation.matchesInvariant(*current) {
		return "", fmt.Errorf(
			"pull-request identity or invariant metadata changed after preparation",
		)
	}
	if remoteHeadOID == r.PreparedHeadOID {
		switch currentDigest {
		case desiredDigest:
			return receiptPullRequestDesired, nil
		case expectation.PriorProjectionDigest:
			return receiptPullRequestHeadPushedWithPriorProjection, nil
		default:
			return "", fmt.Errorf(
				"pull-request projection changed after the prepared head was pushed",
			)
		}
	}
	if !r.ExpectedRemoteHead.Present ||
		remoteHeadOID != r.ExpectedRemoteHead.OID {
		return "", fmt.Errorf(
			"pull request differs from its exact pre-publication expectation",
		)
	}
	if currentDigest == expectation.PriorProjectionDigest {
		return receiptPullRequestPrior, nil
	}
	if r.ExpectedRemoteHead.OID != expectation.PriorHeadOID &&
		currentDigest == desiredDigest {
		return receiptPullRequestDesired, nil
	}
	return "", fmt.Errorf(
		"pull request differs from its exact pre-publication expectation",
	)
}

func (r preparationReceipt) pullRequestMatchesTopology(value pullRequest) bool {
	return value.BaseRefName == strings.TrimPrefix(r.BaseRef, "refs/heads/") &&
		value.HeadRefName == strings.TrimPrefix(r.HeadRef, "refs/heads/") &&
		strings.EqualFold(
			value.HeadRepositoryOwner,
			r.RemoteRepository.Owner,
		) &&
		strings.EqualFold(
			value.HeadRepositoryName,
			r.RemoteRepository.Name,
		) &&
		!value.IsCrossRepository
}

// bindCreatedPullRequest transitions an explicit-absence expectation to the
// first exact created/recovered pull request. The updated receipt must be
// durably installed before any later fallible step so subsequent retries
// reject a closed-and-recreated replacement with the same projection.
func (r preparationReceipt) bindCreatedPullRequest(
	current *pullRequest,
	remoteHeadOID string,
	desired commitProjection,
) (preparationReceipt, error) {
	state, err := r.pullRequestState(current, remoteHeadOID, desired)
	if err != nil {
		return preparationReceipt{}, err
	}
	if r.ExpectedPullRequest.Present {
		return r, nil
	}
	if state != receiptPullRequestCreatedDesired || current == nil {
		return preparationReceipt{}, fmt.Errorf(
			"cannot bind a pull request before exact creation is verified",
		)
	}
	expectation, err := newPullRequestExpectation(current)
	if err != nil {
		return preparationReceipt{}, err
	}
	revisionNonce, err := newReceiptRevisionNonce()
	if err != nil {
		return preparationReceipt{}, err
	}
	bound := r
	bound.RevisionNonce = revisionNonce
	bound.ExpectedPullRequest = expectation
	if err := bound.validate(); err != nil {
		return preparationReceipt{}, err
	}
	return bound, nil
}

// validateTaskOutputPath checks a repository-relative lexical path. Callers
// still need to prove that the parent contains no symlinks and that Git marks
// the path ignored and untracked before trusting its contents.
func validateTaskOutputPath(relative string) error {
	if relative == "" || filepath.IsAbs(relative) ||
		strings.ContainsRune(relative, '\x00') ||
		strings.ContainsRune(relative, '\\') {
		return fmt.Errorf("path must be a safe repository-relative path")
	}
	canonical := filepath.ToSlash(filepath.Clean(filepath.FromSlash(relative)))
	if canonical != relative {
		return fmt.Errorf("path must use a canonical lexical form")
	}
	parts := strings.Split(relative, "/")
	if len(parts) < 3 || parts[0] != "out" || parts[1] == "" ||
		parts[1] == "." || parts[1] == ".." {
		return fmt.Errorf("path must be under an out/<task>/ directory")
	}
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return fmt.Errorf("path contains an unsafe component")
		}
	}
	return nil
}

func (d *delivery) receiptPath(
	ctx context.Context,
	value string,
	requireExisting bool,
) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("--receipt-file is required")
	}
	if filepath.IsAbs(value) {
		if filepath.Clean(value) != value {
			return "", fmt.Errorf("receipt file path must use a canonical lexical form")
		}
	} else if err := validateTaskOutputPath(filepath.ToSlash(value)); err != nil {
		return "", fmt.Errorf(
			"receipt file must be under an ignored out/<task>/ directory: %w",
			err,
		)
	}
	absolute := value
	if !filepath.IsAbs(absolute) {
		absolute = filepath.Join(d.repository.directory, absolute)
	}
	absolute = filepath.Clean(absolute)
	relative, err := filepath.Rel(d.repository.directory, absolute)
	if err != nil {
		return "", fmt.Errorf("locate receipt file in repository: %w", err)
	}
	relative = filepath.ToSlash(relative)
	if err := validateTaskOutputPath(relative); err != nil {
		return "", fmt.Errorf(
			"receipt file must be under an ignored out/<task>/ directory: %w",
			err,
		)
	}
	parent := filepath.Dir(absolute)
	resolvedParent, err := canonicalExistingPath(parent)
	if err != nil {
		return "", fmt.Errorf("resolve receipt directory: %w", err)
	}
	if resolvedParent != filepath.Clean(parent) {
		return "", fmt.Errorf("receipt file path must not contain symlinks")
	}
	tracked, err := d.repository.pathTracked(ctx, relative)
	if err != nil {
		return "", err
	}
	if tracked {
		return "", fmt.Errorf("receipt file must be untracked")
	}
	ignored, err := d.repository.pathIgnored(ctx, relative)
	if err != nil {
		return "", err
	}
	if !ignored {
		return "", fmt.Errorf("receipt file must be ignored by Git")
	}
	info, err := os.Lstat(absolute)
	if err == nil {
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("receipt file must be a regular non-symlink file")
		}
		if info.Mode().Perm() != 0o600 {
			return "", fmt.Errorf(
				"receipt file %s has mode %04o; require 0600",
				relative,
				info.Mode().Perm(),
			)
		}
	} else if requireExisting || !os.IsNotExist(err) {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("receipt file does not exist")
		}
		return "", fmt.Errorf("inspect receipt file: %w", err)
	}
	return absolute, nil
}

// writeAtomicIgnoredJSON validates the ignored task output path and installs
// the JSON value atomically with restrictive 0600 permissions. It is used for
// local typed receipt files that are never treated as credentials.
func (d *delivery) writeAtomicIgnoredJSON(
	ctx context.Context,
	path string,
	description string,
	value any,
) error {
	absolute, err := d.receiptPath(ctx, path, false)
	if err != nil {
		return err
	}
	contents, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", description, err)
	}
	contents = append(contents, '\n')
	if len(contents) > receiptFileLimit {
		return fmt.Errorf("%s exceeds 256 KiB", description)
	}
	temporary, err := os.CreateTemp(
		filepath.Dir(absolute),
		"."+filepath.Base(absolute)+".tmp-",
	)
	if err != nil {
		return fmt.Errorf("create atomic %s: %w", description, err)
	}
	temporaryPath := temporary.Name()
	cleanup := true
	defer func() {
		_ = temporary.Close()
		if cleanup {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(0o600); err != nil {
		return fmt.Errorf("set %s permissions: %w", description, err)
	}
	if _, err := temporary.Write(contents); err != nil {
		return fmt.Errorf("write %s: %w", description, err)
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync %s: %w", description, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close %s: %w", description, err)
	}
	if err := os.Rename(temporaryPath, absolute); err != nil {
		return fmt.Errorf("install %s atomically: %w", description, err)
	}
	cleanup = false
	return nil
}

type receiptFileVersion struct {
	present  bool
	contents []byte
	digest   [sha256.Size]byte
}

func newReceiptFileVersion(contents []byte) receiptFileVersion {
	copyOfContents := append([]byte(nil), contents...)
	return receiptFileVersion{
		present:  true,
		contents: copyOfContents,
		digest:   sha256.Sum256(copyOfContents),
	}
}

func (v receiptFileVersion) equal(other receiptFileVersion) bool {
	if v.present != other.present {
		return false
	}
	if !v.present {
		return true
	}
	return v.digest == other.digest && bytes.Equal(v.contents, other.contents)
}

func captureReceiptFileVersion(path string) (receiptFileVersion, error) {
	contents, err := readStableReceiptFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return receiptFileVersion{}, nil
	}
	if err != nil {
		return receiptFileVersion{}, err
	}
	return newReceiptFileVersion(contents), nil
}

type receiptTransaction struct {
	delivery *delivery
	path     string
	absolute string
	lockPath string
	lockFile *os.File
	version  receiptFileVersion
	closed   bool
}

func (d *delivery) beginReceiptTransaction(
	ctx context.Context,
	path string,
	requireExisting bool,
) (transaction *receiptTransaction, returnErr error) {
	// Validate the destination before opening the lock, but defer the
	// existence requirement until after acquiring it. A publisher that starts
	// while the first preparation is installing the receipt must serialize
	// behind that preparation rather than race its absent destination.
	absolute, err := d.receiptPath(ctx, path, false)
	if err != nil {
		return nil, err
	}
	lockInput := path + ".lock"
	lockPath, err := d.receiptPath(ctx, lockInput, false)
	if err != nil {
		return nil, fmt.Errorf("validate preparation receipt lock path: %w", err)
	}
	lockFile, err := openReceiptLockFile(lockPath)
	if err != nil {
		return nil, err
	}
	locked := false
	defer func() {
		if transaction != nil {
			return
		}
		if locked {
			if unlockErr := syscall.Flock(
				int(lockFile.Fd()),
				syscall.LOCK_UN,
			); unlockErr != nil {
				returnErr = errors.Join(
					returnErr,
					fmt.Errorf("unlock preparation receipt: %w", unlockErr),
				)
			}
		}
		if closeErr := lockFile.Close(); closeErr != nil {
			returnErr = errors.Join(
				returnErr,
				fmt.Errorf("close preparation receipt lock: %w", closeErr),
			)
		}
	}()
	if err := lockReceiptFile(ctx, lockFile); err != nil {
		return nil, err
	}
	locked = true
	confirmedLockPath, err := d.receiptPath(ctx, lockInput, true)
	if err != nil {
		return nil, fmt.Errorf("revalidate preparation receipt lock: %w", err)
	}
	if confirmedLockPath != lockPath {
		return nil, fmt.Errorf("preparation receipt lock path changed")
	}
	if err := requireOpenFileIdentity(lockFile, lockPath, "receipt lock"); err != nil {
		return nil, err
	}
	confirmed, err := d.receiptPath(ctx, path, requireExisting)
	if err != nil {
		return nil, fmt.Errorf("revalidate preparation receipt under lock: %w", err)
	}
	if confirmed != absolute {
		return nil, fmt.Errorf("preparation receipt path changed before locking")
	}
	version, err := captureReceiptFileVersion(absolute)
	if err != nil {
		return nil, err
	}
	if requireExisting && !version.present {
		return nil, fmt.Errorf("preparation receipt disappeared while locking")
	}
	transaction = &receiptTransaction{
		delivery: d,
		path:     path,
		absolute: absolute,
		lockPath: lockPath,
		lockFile: lockFile,
		version:  version,
	}
	return transaction, nil
}

func openReceiptLockFile(path string) (*os.File, error) {
	for attempts := 0; attempts != 4; attempts++ {
		flags := syscall.O_RDWR | syscall.O_CLOEXEC | syscall.O_NOFOLLOW
		fileDescriptor, err := syscall.Open(
			path,
			flags|syscall.O_CREAT|syscall.O_EXCL,
			0o600,
		)
		created := err == nil
		if errors.Is(err, syscall.EEXIST) {
			fileDescriptor, err = syscall.Open(path, flags, 0)
		}
		if errors.Is(err, syscall.ENOENT) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("open preparation receipt lock: %w", err)
		}
		file := os.NewFile(uintptr(fileDescriptor), path)
		if file == nil {
			_ = syscall.Close(fileDescriptor)
			return nil, fmt.Errorf("open preparation receipt lock: invalid file")
		}
		if created {
			if err := file.Chmod(0o600); err != nil {
				_ = file.Close()
				return nil, fmt.Errorf(
					"set preparation receipt lock permissions: %w",
					err,
				)
			}
			if err := syncReceiptDirectory(filepath.Dir(path)); err != nil {
				_ = file.Close()
				return nil, fmt.Errorf(
					"persist preparation receipt lock identity: %w",
					err,
				)
			}
		}
		if err := requireOpenFileIdentity(file, path, "receipt lock"); err != nil {
			_ = file.Close()
			return nil, err
		}
		return file, nil
	}
	return nil, fmt.Errorf("preparation receipt lock path changed repeatedly")
}

func requireOpenFileIdentity(file *os.File, path string, description string) error {
	opened, err := file.Stat()
	if err != nil {
		return fmt.Errorf("inspect opened %s: %w", description, err)
	}
	named, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect named %s: %w", description, err)
	}
	if !opened.Mode().IsRegular() || opened.Mode().Perm() != 0o600 ||
		named.Mode()&os.ModeSymlink != 0 || !os.SameFile(opened, named) {
		return fmt.Errorf(
			"%s must be the same regular non-symlink 0600 file that was opened",
			description,
		)
	}
	return nil
}

func lockReceiptFile(ctx context.Context, file *os.File) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		err := syscall.Flock(
			int(file.Fd()),
			syscall.LOCK_EX|syscall.LOCK_NB,
		)
		if err == nil {
			return nil
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) &&
			!errors.Is(err, syscall.EAGAIN) {
			return fmt.Errorf("lock preparation receipt: %w", err)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for preparation receipt lock: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func (t *receiptTransaction) close() error {
	if t == nil || t.closed {
		return nil
	}
	t.closed = true
	unlockErr := syscall.Flock(int(t.lockFile.Fd()), syscall.LOCK_UN)
	closeErr := t.lockFile.Close()
	var err error
	if unlockErr != nil {
		err = errors.Join(
			err,
			fmt.Errorf("unlock preparation receipt: %w", unlockErr),
		)
	}
	if closeErr != nil {
		err = errors.Join(
			err,
			fmt.Errorf("close preparation receipt lock: %w", closeErr),
		)
	}
	return err
}

func (t *receiptTransaction) read() (preparationReceipt, error) {
	if t == nil || t.closed {
		return preparationReceipt{}, fmt.Errorf(
			"preparation receipt transaction is not active",
		)
	}
	if !t.version.present {
		return preparationReceipt{}, fmt.Errorf("preparation receipt does not exist")
	}
	return decodePreparationReceipt(t.version.contents)
}

func (t *receiptTransaction) requireCurrent(ctx context.Context) error {
	if t == nil || t.closed {
		return fmt.Errorf("preparation receipt transaction is not active")
	}
	if err := requireOpenFileIdentity(
		t.lockFile,
		t.lockPath,
		"receipt lock",
	); err != nil {
		return err
	}
	confirmed, err := t.delivery.receiptPath(ctx, t.path, t.version.present)
	if err != nil {
		return fmt.Errorf("revalidate preparation receipt version: %w", err)
	}
	if confirmed != t.absolute {
		return fmt.Errorf("preparation receipt path changed")
	}
	current, err := captureReceiptFileVersion(t.absolute)
	if err != nil {
		return err
	}
	return requireReceiptFileVersion(t.version, current)
}

func requireReceiptFileVersion(
	expected receiptFileVersion,
	current receiptFileVersion,
) error {
	if expected.equal(current) {
		return nil
	}
	return fmt.Errorf(
		"preparation receipt changed since its exact byte version was read; refusing a stale operation",
	)
}

func encodePreparationReceipt(receipt preparationReceipt) ([]byte, error) {
	if err := receipt.validate(); err != nil {
		return nil, err
	}
	contents, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode preparation receipt: %w", err)
	}
	contents = append(contents, '\n')
	if len(contents) > receiptFileLimit {
		return nil, fmt.Errorf("preparation receipt exceeds 256 KiB")
	}
	return contents, nil
}

func (t *receiptTransaction) write(
	ctx context.Context,
	receipt preparationReceipt,
) error {
	return t.writeWithDirectorySync(ctx, receipt, syncReceiptDirectory)
}

func (t *receiptTransaction) writeWithDirectorySync(
	ctx context.Context,
	receipt preparationReceipt,
	syncDirectory func(string) error,
) error {
	contents, err := encodePreparationReceipt(receipt)
	if err != nil {
		return err
	}
	if t.version.present {
		prior, decodeErr := decodePreparationReceipt(t.version.contents)
		if decodeErr == nil && prior.RevisionNonce == receipt.RevisionNonce {
			return fmt.Errorf(
				"preparation receipt update must use a fresh revision nonce; " +
					"delete the existing receipt file and re-run prepare",
			)
		}
	}
	directory := filepath.Dir(t.absolute)
	temporary, err := os.CreateTemp(
		directory,
		"."+filepath.Base(t.absolute)+".tmp-",
	)
	if err != nil {
		return fmt.Errorf("create atomic receipt file: %w", err)
	}
	temporaryPath := temporary.Name()
	cleanup := true
	defer func() {
		_ = temporary.Close()
		if cleanup {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(0o600); err != nil {
		return fmt.Errorf("set receipt file permissions: %w", err)
	}
	if _, err := temporary.Write(contents); err != nil {
		return fmt.Errorf("write preparation receipt: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync preparation receipt: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close preparation receipt: %w", err)
	}
	if err := t.requireCurrent(ctx); err != nil {
		return err
	}
	installed, err := installReceiptAfterStableComparison(
		temporaryPath,
		t.absolute,
		t.version,
		syncDirectory,
	)
	if installed {
		// The temporary path no longer identifies the file after the rename.
		// In particular, do not try to roll the rename back if syncing the
		// directory reports an uncertain durability outcome.
		cleanup = false
		t.version = newReceiptFileVersion(contents)
	}
	return err
}

// installReceiptAfterStableComparison must be called while the stable
// per-receipt lock is held. That lock supplies cross-process exclusion for
// cooperating repo_delivery processes. The adjacent comparison additionally
// detects a lock-ignoring write that completed before it, but it is not a
// portable atomic pathname-content CAS: such a writer can still race the
// following rename.
func installReceiptAfterStableComparison(
	temporaryPath string,
	absolute string,
	expected receiptFileVersion,
	syncDirectory func(string) error,
) (bool, error) {
	current, err := captureReceiptFileVersion(absolute)
	if err != nil {
		return false, err
	}
	if err := requireReceiptFileVersion(expected, current); err != nil {
		return false, err
	}
	return installReceiptAtomically(temporaryPath, absolute, syncDirectory)
}

func installReceiptAtomically(
	temporaryPath string,
	absolute string,
	syncDirectory func(string) error,
) (bool, error) {
	if err := os.Rename(temporaryPath, absolute); err != nil {
		return false, fmt.Errorf(
			"install preparation receipt atomically: %w",
			err,
		)
	}
	if err := syncDirectory(filepath.Dir(absolute)); err != nil {
		return true, fmt.Errorf(
			"preparation receipt install outcome unknown: atomic rename succeeded, but containing-directory sync failed; re-inspect the receipt and repository state before retrying: %w",
			err,
		)
	}
	return true, nil
}

func syncReceiptDirectory(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open containing directory: %w", err)
	}
	if err := directory.Sync(); err != nil {
		_ = directory.Close()
		return fmt.Errorf("sync containing directory: %w", err)
	}
	if err := directory.Close(); err != nil {
		return fmt.Errorf("close containing directory: %w", err)
	}
	return nil
}

type canonicalJSONShape struct {
	fields     map[string]*canonicalJSONShape
	element    *canonicalJSONShape
	scalarKind string
	optional   bool
}

func preparationReceiptJSONShape() *canonicalJSONShape {
	stringValue := &canonicalJSONShape{scalarKind: "string"}
	boolValue := &canonicalJSONShape{scalarKind: "boolean"}
	numberValue := &canonicalJSONShape{scalarKind: "number"}
	stringsArray := &canonicalJSONShape{element: stringValue}
	remoteRepository := &canonicalJSONShape{fields: map[string]*canonicalJSONShape{
		"host":  stringValue,
		"owner": stringValue,
		"name":  stringValue,
	}}
	return &canonicalJSONShape{fields: map[string]*canonicalJSONShape{
		"schema":                 stringValue,
		"revision_nonce":         stringValue,
		"repository_fingerprint": stringValue,
		"remote_name":            stringValue,
		"fetch_endpoint_sha256":  stringValue,
		"push_endpoint_sha256":   stringValue,
		"remote_repository":      remoteRepository,
		"forge":                  stringValue,
		"base_ref":               stringValue,
		"base_oid":               stringValue,
		"head_ref":               stringValue,
		"prepared_head_oid":      stringValue,
		"prepared_tree_oid":      stringValue,
		"expected_remote_head": {fields: map[string]*canonicalJSONShape{
			"ref":     stringValue,
			"present": boolValue,
			"oid":     stringValue,
		}},
		"expected_pull_request": {fields: map[string]*canonicalJSONShape{
			"present":                 boolValue,
			"id":                      stringValue,
			"number":                  numberValue,
			"identity_sha256":         stringValue,
			"author_login":            stringValue,
			"is_draft":                boolValue,
			"prior_head_oid":          stringValue,
			"prior_projection_sha256": stringValue,
		}},
		"scope": {fields: map[string]*canonicalJSONShape{
			"mode":             stringValue,
			"authorized_paths": stringsArray,
			"aggregate_paths":  stringsArray,
		}},
		"rewrite_authorization": {optional: true, fields: map[string]*canonicalJSONShape{
			"schema":                 stringValue,
			"issued_at":              stringValue,
			"repository_fingerprint": stringValue,
			"remote_repository":      remoteRepository,
			"provider":               stringValue,
			"old_remote_oid":         stringValue,
			"new_head_oid":           stringValue,
			"owner_root":             stringValue,
			"task_paths":             stringsArray,
			"source_receipt_digest":  {optional: true, scalarKind: "string"},
		}},
	}}
}

func validateCanonicalReceiptJSON(contents []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.UseNumber()
	if err := validateCanonicalJSONValue(
		decoder,
		preparationReceiptJSONShape(),
		"receipt",
	); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			return fmt.Errorf("preparation receipt contains multiple JSON values")
		}
		return fmt.Errorf("finish reading preparation receipt JSON: %w", err)
	}
	return nil
}

func validateCanonicalJSONValue(
	decoder *json.Decoder,
	shape *canonicalJSONShape,
	location string,
) error {
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("read %s: %w", location, err)
	}
	if shape.scalarKind != "" {
		if _, structured := token.(json.Delim); structured {
			return fmt.Errorf("%s must be a scalar JSON value", location)
		}
		validType := false
		switch shape.scalarKind {
		case "string":
			_, validType = token.(string)
		case "boolean":
			_, validType = token.(bool)
		case "number":
			_, validType = token.(json.Number)
		}
		if !validType {
			return fmt.Errorf(
				"%s must be a JSON %s",
				location,
				shape.scalarKind,
			)
		}
		return nil
	}
	opening, ok := token.(json.Delim)
	if !ok {
		return fmt.Errorf("%s has the wrong JSON shape", location)
	}
	if shape.fields != nil {
		if opening != '{' {
			return fmt.Errorf("%s must be a JSON object", location)
		}
		seen := make(map[string]struct{}, len(shape.fields))
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return fmt.Errorf("read %s object key: %w", location, err)
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("%s contains a non-string object key", location)
			}
			child, allowed := shape.fields[key]
			if !allowed {
				return fmt.Errorf(
					"%s contains non-canonical key %q",
					location,
					key,
				)
			}
			if _, duplicate := seen[key]; duplicate {
				return fmt.Errorf("%s contains duplicate key %q", location, key)
			}
			seen[key] = struct{}{}
			if err := validateCanonicalJSONValue(
				decoder,
				child,
				location+"."+key,
			); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("close %s object: %w", location, err)
		}
		if closing != json.Delim('}') {
			return fmt.Errorf("%s has an invalid object terminator", location)
		}
		requiredCount := 0
		for key, child := range shape.fields {
			if !child.optional {
				requiredCount++
			}
			_ = key
		}
		if seenCount := len(seen); seenCount < requiredCount {
			var missing []string
			for key := range shape.fields {
				if shape.fields[key].optional {
					continue
				}
				if _, present := seen[key]; !present {
					missing = append(missing, key)
				}
			}
			sort.Strings(missing)
			return fmt.Errorf(
				"%s has %d keys, want %d (missing: %q)",
				location,
				seenCount,
				requiredCount,
				missing,
			)
		}
		return nil
	}
	if shape.element != nil {
		if opening != '[' {
			return fmt.Errorf("%s must be a JSON array", location)
		}
		index := 0
		for decoder.More() {
			if err := validateCanonicalJSONValue(
				decoder,
				shape.element,
				fmt.Sprintf("%s[%d]", location, index),
			); err != nil {
				return err
			}
			index++
		}
		closing, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("close %s array: %w", location, err)
		}
		if closing != json.Delim(']') {
			return fmt.Errorf("%s has an invalid array terminator", location)
		}
		return nil
	}
	return fmt.Errorf("%s has no canonical JSON shape", location)
}

func stableFileInfo(left os.FileInfo, right os.FileInfo) bool {
	return os.SameFile(left, right) &&
		left.Mode() == right.Mode() &&
		left.Size() == right.Size() &&
		left.ModTime().Equal(right.ModTime())
}

func readStableReceiptFile(path string) ([]byte, error) {
	pathBefore, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect preparation receipt: %w", err)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open preparation receipt: %w", err)
	}
	defer file.Close()
	openedBefore, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("inspect opened preparation receipt: %w", err)
	}
	if !openedBefore.Mode().IsRegular() ||
		openedBefore.Mode().Perm() != 0o600 ||
		!stableFileInfo(pathBefore, openedBefore) {
		return nil, fmt.Errorf(
			"preparation receipt changed while it was being opened",
		)
	}
	first, err := io.ReadAll(io.LimitReader(file, receiptFileLimit+1))
	if err != nil {
		return nil, fmt.Errorf("read preparation receipt: %w", err)
	}
	if len(first) > receiptFileLimit {
		return nil, fmt.Errorf("preparation receipt exceeds 256 KiB")
	}
	afterFirst, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("reinspect opened preparation receipt: %w", err)
	}
	if !stableFileInfo(openedBefore, afterFirst) {
		return nil, fmt.Errorf("preparation receipt changed while it was read")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("rewind preparation receipt: %w", err)
	}
	second, err := io.ReadAll(io.LimitReader(file, receiptFileLimit+1))
	if err != nil {
		return nil, fmt.Errorf("reread preparation receipt: %w", err)
	}
	afterSecond, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("inspect reread preparation receipt: %w", err)
	}
	pathAfter, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("reinspect preparation receipt path: %w", err)
	}
	if len(second) > receiptFileLimit ||
		!bytes.Equal(first, second) ||
		!stableFileInfo(afterFirst, afterSecond) ||
		!stableFileInfo(afterSecond, pathAfter) ||
		pathAfter.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("preparation receipt changed while it was read")
	}
	return first, nil
}

func decodePreparationReceipt(
	contents []byte,
) (preparationReceipt, error) {
	if err := validateCanonicalReceiptJSON(contents); err != nil {
		return preparationReceipt{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	var receipt preparationReceipt
	if err := decoder.Decode(&receipt); err != nil {
		return preparationReceipt{}, fmt.Errorf(
			"decode preparation receipt: %w",
			err,
		)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return preparationReceipt{}, fmt.Errorf(
				"preparation receipt contains multiple JSON values",
			)
		}
		return preparationReceipt{}, fmt.Errorf(
			"finish decoding preparation receipt: %w",
			err,
		)
	}
	if err := receipt.validate(); err != nil {
		return preparationReceipt{}, err
	}
	return receipt, nil
}

func (d *delivery) readReceipt(
	ctx context.Context,
	path string,
) (result preparationReceipt, returnErr error) {
	transaction, err := d.beginReceiptTransaction(ctx, path, true)
	if err != nil {
		return preparationReceipt{}, err
	}
	defer func() {
		if closeErr := transaction.close(); closeErr != nil {
			returnErr = errors.Join(returnErr, closeErr)
		}
	}()
	result, err = transaction.read()
	if err != nil {
		return preparationReceipt{}, err
	}
	if err := transaction.requireCurrent(ctx); err != nil {
		return preparationReceipt{}, err
	}
	return result, nil
}
