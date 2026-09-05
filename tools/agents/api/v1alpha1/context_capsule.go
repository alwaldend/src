// Package v1alpha1 capsule types.
package v1alpha1

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CapsuleIdentity is the repository, workspace, worktree, revision, and
// dirty-input identity for one zero-context read.
type CapsuleIdentity struct {
	Repository        string                 `json:"repository"`
	WorkspaceRoot     string                 `json:"workspaceRoot"`
	WorktreePath      string                 `json:"worktreePath"`
	Revision          string                 `json:"revision"`
	DirtyInputs       bool                   `json:"dirtyInputs"`
	RevisionSource    string                 `json:"revisionSource,omitempty"`
	DirtyInputsSource string                 `json:"dirtyInputsSource,omitempty"`
	Git               *CapsuleGitObservation `json:"git,omitempty"`
	SourceDigest      string                 `json:"sourceDigest,omitempty"`
	SourceTime        string                 `json:"sourceTime,omitempty"`
	InputDigest       string                 `json:"inputDigest"`
	ByteSize          int64                  `json:"byteSize"`
	InputLanguages    []string               `json:"inputLanguages,omitempty"`
}

// CapsuleGitObservation separates local Git evidence from caller declarations.
// Dirty covers tracked and untracked paths under WorkspaceRoot, excluding
// ignored files. Nil means unknown, and a missing revision can mean an unborn
// branch. Neither field establishes the freshness of generated catalogs.
type CapsuleGitObservation struct {
	Revision    string   `json:"revision,omitempty"`
	Dirty       *bool    `json:"dirty"`
	ObservedAt  string   `json:"observedAt"`
	Unavailable []string `json:"unavailable,omitempty"`
}

// CapsuleTask is the task/session, coordinator, worker, and run identity.
type CapsuleTask struct {
	TaskID      string `json:"taskId,omitempty"`
	SessionID   string `json:"sessionId,omitempty"`
	Coordinator string `json:"coordinator,omitempty"`
	WorkerID    string `json:"workerId,omitempty"`
	RunID       string `json:"runId,omitempty"`
}

// CapsuleOutcome is the requested outcome, authority, and budget binding.
type CapsuleOutcome struct {
	RequestedOutcome string `json:"requestedOutcome,omitempty"`
	Authority        string `json:"authority,omitempty"`
	Budget           Budget `json:"budget"`
	GoalBinding      string `json:"goalBinding,omitempty"`
}

// CapsuleDocumentSource is one applicable instruction/owner document.
type CapsuleDocumentSource struct {
	Path   string `json:"path"`
	Digest string `json:"digest"`
}

// CapsuleComponent is the component, workspace, and review-owner slice.
type CapsuleComponent struct {
	Path           string `json:"path"`
	ComponentID    string `json:"componentId,omitempty"`
	OwnerReadme    string `json:"ownerReadme,omitempty"`
	Workspace      string `json:"workspace,omitempty"`
	Lifecycle      string `json:"lifecycle,omitempty"`
	ReviewOwners   string `json:"reviewOwners,omitempty"`
	Classification string `json:"classification,omitempty"`
}

// CapsuleCapability is one candidate capability.
type CapsuleCapability struct {
	ID        string   `json:"id"`
	Kind      string   `json:"kind"`
	Owner     string   `json:"owner"`
	Role      string   `json:"role"`
	Effects   []string `json:"effects,omitempty"`
	Cost      string   `json:"cost,omitempty"`
	Access    string   `json:"access,omitempty"`
	Excluded  string   `json:"excluded,omitempty"`
	Providers []string `json:"providers,omitempty"`
}

// CapsuleCheck is one relevant check in the applicable workspace slice.
type CapsuleCheck struct {
	WorkspaceID     string   `json:"workspaceId,omitempty"`
	PhaseID         string   `json:"phaseId"`
	ProviderRef     string   `json:"providerRef,omitempty"`
	CommandTemplate string   `json:"commandTemplate"`
	Selectors       []string `json:"selectors,omitempty"`
}

// CapsuleProviderStatus is one runtime/boot/incarnation observation.
type CapsuleProviderStatus struct {
	ProviderID       string `json:"providerId"`
	BootID           string `json:"bootId,omitempty"`
	CatalogETag      string `json:"catalogEtag,omitempty"`
	DesiredRevision  string `json:"desiredRevision,omitempty"`
	ObservedRevision string `json:"observedRevision,omitempty"`
	ContractHash     string `json:"contractHash,omitempty"`
	State            string `json:"state"`
	ObservationTime  string `json:"observationTime,omitempty"`
	ExpiresAt        string `json:"expiresAt,omitempty"`
	Unavailable      string `json:"unavailable,omitempty"`
}

// CapsuleProvenance is the capsule's observation provenance.
type CapsuleProvenance struct {
	Inputs       []InputReference `json:"inputs,omitempty"`
	ObservedAt   string           `json:"observedAt,omitempty"`
	Freshness    string           `json:"freshness"`
	Completeness Completeness     `json:"completeness"`
	Truncated    bool             `json:"truncated"`
	Limitations  []string         `json:"limitations,omitempty"`
	NextActions  []string         `json:"nextActions,omitempty"`
	AuthorizedBy string           `json:"authorizedBy,omitempty"`
}

// ContextCapsule is the bounded read-only projection over owner-local
// catalogs and optional runtime observations.
type ContextCapsule struct {
	APIVersion   string                  `json:"apiVersion"`
	Kind         string                  `json:"kind"`
	ID           string                  `json:"id"`
	Identity     CapsuleIdentity         `json:"identity"`
	Task         CapsuleTask             `json:"task,omitempty"`
	Outcome      CapsuleOutcome          `json:"outcome"`
	Documents    []CapsuleDocumentSource `json:"documents,omitempty"`
	Component    CapsuleComponent        `json:"component,omitempty"`
	Capabilities []CapsuleCapability     `json:"capabilities,omitempty"`
	Checks       []CapsuleCheck          `json:"checks,omitempty"`
	Providers    []CapsuleProviderStatus `json:"providers,omitempty"`
	Provenance   CapsuleProvenance       `json:"provenance"`
}

const (
	// KindContextCapsule identifies the bounded context capsule.
	KindContextCapsule = "context-capsule"
)

// Validate enforces the bounded capsule invariants.
func (capsule ContextCapsule) Validate() error {
	if capsule.APIVersion != APIVersion || capsule.Kind != KindContextCapsule {
		return fmt.Errorf("malformed capsule identity")
	}
	if capsule.Identity.Repository == "" ||
		capsule.Identity.WorkspaceRoot == "" ||
		capsule.Identity.WorktreePath == "" ||
		capsule.Identity.InputDigest == "" {
		return fmt.Errorf("capsule identity is incomplete")
	}
	if capsule.Identity.ByteSize <= 0 {
		return fmt.Errorf("capsule byte size must be positive")
	}
	if git := capsule.Identity.Git; git != nil {
		if _, err := time.Parse(time.RFC3339Nano, git.ObservedAt); err != nil {
			return fmt.Errorf("capsule Git observation time is invalid")
		}
		if git.Revision != "" {
			if _, err := hex.DecodeString(git.Revision); err != nil || len(git.Revision) != 40 && len(git.Revision) != 64 {
				return fmt.Errorf("capsule Git revision is invalid")
			}
		}
		for _, reason := range git.Unavailable {
			if strings.TrimSpace(reason) == "" {
				return fmt.Errorf("capsule Git unavailable reasons must not be blank")
			}
		}
		if (git.Dirty == nil || git.Revision == "") && len(git.Unavailable) == 0 {
			return fmt.Errorf("unavailable capsule Git fields require a reason")
		}
	}
	if err := capsule.Identity.validateSources(); err != nil {
		return err
	}
	if capsule.Provenance.Completeness == "" {
		return fmt.Errorf("capsule provenance completeness is required")
	}
	if len(capsule.Documents) > 0 {
		seen := map[string]bool{}
		for _, document := range capsule.Documents {
			if document.Path == "" || document.Digest == "" {
				return fmt.Errorf("capsule document is incomplete")
			}
			if seen[document.Path] {
				return fmt.Errorf("duplicate capsule document path %q", document.Path)
			}
			seen[document.Path] = true
		}
	}
	return nil
}

// Omitted source labels preserve legacy capsules. Explicit labels must name
// a supported basis and agree with the corresponding value when derived.
// Caller declarations remain independent of any accompanying observation.
func (identity CapsuleIdentity) validateSources() error {
	switch identity.RevisionSource {
	case "":
	case "caller-declared":
		if strings.TrimSpace(identity.Revision) == "" {
			return fmt.Errorf("caller-declared revision source requires a nonblank revision")
		}
	case "input-digest":
		if identity.Revision != identity.InputDigest {
			return fmt.Errorf("input-digest revision source requires revision to match inputDigest")
		}
	case "observed-git-head":
		if identity.Git == nil || identity.Git.Revision == "" || identity.Revision != identity.Git.Revision {
			return fmt.Errorf("observed-git-head revision source requires a matching Git revision observation")
		}
	default:
		return fmt.Errorf("unknown capsule revision source %q", identity.RevisionSource)
	}
	switch identity.DirtyInputsSource {
	case "", "caller-declared":
	case "conservative-default":
		if !identity.DirtyInputs {
			return fmt.Errorf("conservative-default dirty-input source requires dirtyInputs to be true")
		}
	case "observed-git-status":
		if identity.Git == nil || identity.Git.Dirty == nil || identity.DirtyInputs != *identity.Git.Dirty {
			return fmt.Errorf("observed-git-status dirty-input source requires a matching Git dirty observation")
		}
	default:
		return fmt.Errorf("unknown capsule dirty-input source %q", identity.DirtyInputsSource)
	}
	return nil
}

// CapsuleID computes the deterministic identity of a complete capsule over
// its canonical bytes with the ID field omitted.
func CapsuleID(capsule ContextCapsule) (string, error) {
	if err := capsule.Validate(); err != nil {
		return "", err
	}
	withoutID := capsule
	withoutID.ID = ""
	content, err := json.Marshal(withoutID)
	if err != nil {
		return "", fmt.Errorf("encode context capsule: %w", err)
	}
	value := sha256.Sum256(content)
	return "capsule." + hex.EncodeToString(value[:]), nil
}

// CanonicalContextJSON assigns the deterministic capsule ID if absent and
// encodes a complete context capsule deterministically.
func CanonicalContextJSON(capsule ContextCapsule) ([]byte, error) {
	if capsule.ID == "" {
		id, err := CapsuleID(capsule)
		if err != nil {
			return nil, err
		}
		capsule.ID = id
	}
	if err := capsule.Validate(); err != nil {
		return nil, err
	}
	return append(append([]byte(nil), mustMarshal(capsule)...), '\n'), nil
}

func mustMarshal(capsule ContextCapsule) []byte {
	content, err := json.Marshal(capsule)
	if err != nil {
		panic(fmt.Sprintf("encode context capsule: %v", err))
	}
	return content
}

// DecodeContextCapsule decodes and validates a complete context capsule
// document with unknown-field and trailing-JSON rejection.
func DecodeContextCapsule(content []byte) (ContextCapsule, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var capsule ContextCapsule
	if err := decoder.Decode(&capsule); err != nil {
		return ContextCapsule{}, fmt.Errorf("decode context capsule: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return ContextCapsule{}, err
	}
	if err := capsule.Validate(); err != nil {
		return ContextCapsule{}, err
	}
	return capsule, nil
}

// RenderContextMarkdown renders a human projection of the context capsule
// from the same data. It is not a second authority.
func RenderContextMarkdown(capsule ContextCapsule) string {
	var builder strings.Builder
	builder.WriteString("# Context capsule\n\n")
	fmt.Fprintf(&builder, "- ID: `%s`\n", capsule.ID)
	fmt.Fprintf(&builder, "- Repository: `%s`\n", capsule.Identity.Repository)
	fmt.Fprintf(&builder, "- Workspace root: `%s`\n", capsule.Identity.WorkspaceRoot)
	fmt.Fprintf(&builder, "- Worktree: `%s`\n", capsule.Identity.WorktreePath)
	fmt.Fprintf(&builder, "- Revision: `%s`\n", capsule.Identity.Revision)
	fmt.Fprintf(&builder, "- Dirty inputs: %t\n", capsule.Identity.DirtyInputs)
	if capsule.Identity.RevisionSource != "" {
		fmt.Fprintf(&builder, "- Revision source: %s\n", capsule.Identity.RevisionSource)
	}
	if capsule.Identity.DirtyInputsSource != "" {
		fmt.Fprintf(&builder, "- Dirty-input source: %s\n", capsule.Identity.DirtyInputsSource)
	}
	if git := capsule.Identity.Git; git != nil {
		fmt.Fprintf(&builder, "- Git observed at: %s\n", git.ObservedAt)
		if git.Revision != "" {
			fmt.Fprintf(&builder, "- Observed Git HEAD: `%s`\n", git.Revision)
		}
		if git.Dirty == nil {
			builder.WriteString("- Observed Git dirty state: unavailable\n")
		} else {
			fmt.Fprintf(&builder, "- Observed Git dirty state: %t\n", *git.Dirty)
		}
		for _, reason := range git.Unavailable {
			fmt.Fprintf(&builder, "- Git limitation: %s\n", reason)
		}
	}
	fmt.Fprintf(&builder, "- Input digest: `%s`\n", capsule.Identity.InputDigest)
	fmt.Fprintf(&builder, "- Byte size: %d\n", capsule.Identity.ByteSize)
	builder.WriteString("\n## Component\n\n")
	if capsule.Component.Path == "" {
		builder.WriteString("None.\n")
	} else {
		fmt.Fprintf(&builder, "- Path: `%s`\n", capsule.Component.Path)
		fmt.Fprintf(&builder, "- Component: `%s`\n", capsule.Component.ComponentID)
		fmt.Fprintf(&builder, "- Workspace: `%s`\n", capsule.Component.Workspace)
		fmt.Fprintf(&builder, "- Owner README: `%s`\n", capsule.Component.OwnerReadme)
		fmt.Fprintf(&builder, "- Lifecycle: `%s`\n", capsule.Component.Lifecycle)
		fmt.Fprintf(&builder, "- Review owners: %s\n", capsule.Component.ReviewOwners)
	}
	builder.WriteString("\n## Documents\n\n")
	if len(capsule.Documents) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, document := range capsule.Documents {
			fmt.Fprintf(&builder, "- `%s` (%s)\n", document.Path, document.Digest)
		}
	}
	builder.WriteString("\n## Capabilities\n\n")
	if len(capsule.Capabilities) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, capability := range capsule.Capabilities {
			fmt.Fprintf(&builder, "- `%s` (%s, %s, owned by `%s`) — %s\n",
				capability.ID, capability.Kind, capability.Role,
				capability.Owner, capability.Access)
		}
	}
	builder.WriteString("\n## Checks\n\n")
	if len(capsule.Checks) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, check := range capsule.Checks {
			fmt.Fprintf(&builder, "- `%s` via `%s`: %s\n",
				check.PhaseID, check.ProviderRef, check.CommandTemplate)
		}
	}
	builder.WriteString("\n## Providers\n\n")
	if len(capsule.Providers) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, provider := range capsule.Providers {
			fmt.Fprintf(&builder, "- `%s`: %s (etag `%s`, revision `%s`)\n",
				provider.ProviderID, provider.State,
				provider.CatalogETag, provider.ObservedRevision)
		}
	}
	builder.WriteString("\n## Provenance\n\n")
	fmt.Fprintf(&builder, "- Freshness: %s\n", capsule.Provenance.Freshness)
	fmt.Fprintf(&builder, "- Completeness: %s\n", capsule.Provenance.Completeness)
	fmt.Fprintf(&builder, "- Truncated: %t\n", capsule.Provenance.Truncated)
	if len(capsule.Provenance.Limitations) > 0 {
		for _, limitation := range capsule.Provenance.Limitations {
			fmt.Fprintf(&builder, "- Limitation: %s\n", limitation)
		}
	}
	builder.WriteString("\n## Next actions\n\n")
	if len(capsule.Provenance.NextActions) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, action := range capsule.Provenance.NextActions {
			fmt.Fprintf(&builder, "- %s\n", action)
		}
	}
	return builder.String()
}
