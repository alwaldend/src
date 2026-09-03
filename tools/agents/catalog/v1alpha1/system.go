package catalogv1alpha1

import (
	"encoding/json"
	"fmt"
)

const (
	// KindCapabilityCatalog identifies the capability catalog kind.
	KindCapabilityCatalog = "capability-catalog"
	// KindWorkspaceCheckCatalog identifies the workspace-check catalog kind.
	KindWorkspaceCheckCatalog = "workspace-check-catalog"
	// KindGoalCatalog identifies the goal catalog kind.
	KindGoalCatalog = "goal-catalog"
	// KindAgentSystemIndex identifies the agent system index kind.
	KindAgentSystemIndex = "agent-system-index"
)

// CapabilitySkill is one registered skill with its routing and dependency
// metadata.
type CapabilitySkill struct {
	ID                   string   `json:"id"`
	Owner                string   `json:"owner"`
	CanonicalPath        string   `json:"canonicalPath"`
	DiscoveryPath        string   `json:"discoveryPath"`
	Layer                string   `json:"layer"`
	Activation           string   `json:"activation"`
	Exclusions           []string `json:"exclusions"`
	CapabilityRefs       []string `json:"capabilityRefs"`
	Dependencies         []string `json:"dependencies"`
	Conflicts            []string `json:"conflicts"`
	ProviderRequirements []string `json:"providerRequirements"`
	ContextCost          string   `json:"contextCost"`
	EvaluationMaturity   string   `json:"evaluationMaturity"`
}

// CapabilityProvider is one executable capability provider.
type CapabilityProvider struct {
	ID             string   `json:"id"`
	Owner          string   `json:"owner"`
	Kind           string   `json:"kind"`
	SourcePath     string   `json:"sourcePath"`
	Classification string   `json:"classification"`
	ActionRefs     []string `json:"actionRefs"`
}

// CapabilityCatalog is the bounded capability projection over registry
// skills, runtime tools, direct binaries, operation providers, and discovery
// links.
type CapabilityCatalog struct {
	CatalogEnvelope
	Skills    []CapabilitySkill    `json:"skills"`
	Providers []CapabilityProvider `json:"providers"`
}

// Validate checks capability-specific invariants and the shared envelope.
func (catalog CapabilityCatalog) Validate() error {
	if catalog.Skills == nil || catalog.Providers == nil {
		return fmt.Errorf("capability catalog arrays must be non-null")
	}
	itemIDs := make([]string, 0, len(catalog.Skills)+len(catalog.Providers))
	seenSkill := map[string]bool{}
	for _, skill := range catalog.Skills {
		if skill.ID == "" || skill.Owner == "" || skill.CanonicalPath == "" ||
			skill.DiscoveryPath == "" || skill.Layer == "" ||
			skill.Activation == "" || skill.ContextCost == "" ||
			skill.EvaluationMaturity == "" {
			return fmt.Errorf("capability skill %q is incomplete", skill.ID)
		}
		if seenSkill[skill.ID] {
			return fmt.Errorf("duplicate capability skill %q", skill.ID)
		}
		seenSkill[skill.ID] = true
		if absPath(skill.CanonicalPath) || absPath(skill.DiscoveryPath) ||
			escapes(skill.CanonicalPath) || escapes(skill.DiscoveryPath) {
			return fmt.Errorf("capability skill %s has malformed path", skill.ID)
		}
		for _, ref := range skill.Dependencies {
			if !seenSkill[ref] && ref != skill.ID && !knownExternalRef(ref) {
				return fmt.Errorf("capability skill %s refers to unknown dependency %q",
					skill.ID, ref)
			}
		}
		itemIDs = append(itemIDs, "skill."+skill.ID)
	}
	seenProvider := map[string]bool{}
	for _, provider := range catalog.Providers {
		if provider.ID == "" || provider.Owner == "" || provider.Kind == "" ||
			provider.SourcePath == "" || provider.Classification == "" {
			return fmt.Errorf("capability provider %q is incomplete", provider.ID)
		}
		if seenProvider[provider.ID] {
			return fmt.Errorf("duplicate capability provider %q", provider.ID)
		}
		seenProvider[provider.ID] = true
		if absPath(provider.SourcePath) || escapes(provider.SourcePath) {
			return fmt.Errorf("capability provider %s has malformed source path %q",
				provider.ID, provider.SourcePath)
		}
		itemIDs = append(itemIDs, "provider."+provider.ID)
	}
	return catalog.CatalogEnvelope.Validate(itemIDs)
}

func knownExternalRef(ref string) bool {
	return ref == "" || len(ref) > 3 && ref[:3] == "ext:"
}

// ItemIDs returns the sorted set-like identity list for a capability catalog.
func (catalog CapabilityCatalog) ItemIDs() []string {
	ids := make([]string, 0, len(catalog.Skills)+len(catalog.Providers))
	for _, skill := range catalog.Skills {
		ids = append(ids, "skill."+skill.ID)
	}
	for _, provider := range catalog.Providers {
		ids = append(ids, "provider."+provider.ID)
	}
	return ids
}

// CanonicalJSONCapability encodes a complete capability catalog with a digest
// over the canonical bytes with the digest field omitted.
func CanonicalJSONCapability(catalog CapabilityCatalog) ([]byte, error) {
	if err := catalog.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := catalog
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode capability catalog: %w", err)
	}
	withoutDigest.Digest = digest(content)
	contentWithDigest, err := json.MarshalIndent(withoutDigest, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("encode capability catalog with digest: %w", err)
	}
	return append(contentWithDigest, '\n'), nil
}

// VerifySelfDigest recomputes the content-addressed digest over the canonical
// bytes of the complete capability document and compares it to the stored
// digest, rejecting an empty or mismatched digest.
func (catalog CapabilityCatalog) VerifySelfDigest() error {
	withoutDigest := catalog
	withoutDigest.Digest = ""
	return verifySelfDigest(catalog.Digest, withoutDigest)
}

// DecodeCapabilityStrict decodes and validates a complete capability catalog
// document with unknown-field and trailing-JSON rejection.
func DecodeCapabilityStrict(content []byte) (CapabilityCatalog, error) {
	var catalog CapabilityCatalog
	if err := DecodeStrict(content, &catalog); err != nil {
		return CapabilityCatalog{}, err
	}
	if err := catalog.Validate(); err != nil {
		return CapabilityCatalog{}, err
	}
	if err := catalog.VerifySelfDigest(); err != nil {
		return CapabilityCatalog{}, err
	}
	return catalog, nil
}

// WorkspaceProjection records which workspace-check feature sets apply.
type WorkspaceProjection struct {
	BazelIgnore     bool `json:"bazelIgnore"`
	RootOverride    bool `json:"rootOverride"`
	DocsAggregation bool `json:"docsAggregation"`
	FullCheck       bool `json:"fullCheck"`
}

// WorkspacePhase is one bound check phase template.
type WorkspacePhase struct {
	ID              string `json:"id"`
	ProviderRef     string `json:"providerRef"`
	CommandTemplate string `json:"commandTemplate"`
}

// WorkspaceRecord is one tracked Bzlmod workspace root with its projections.
type WorkspaceRecord struct {
	ID          string              `json:"id"`
	Path        string              `json:"path"`
	ModulePath  string              `json:"modulePath"`
	ModuleName  string              `json:"moduleName"`
	Projections WorkspaceProjection `json:"projections"`
	Phases      []WorkspacePhase    `json:"phases"`
}

// WorkspaceCheckCatalog is the bounded projection over tracked MODULE.bazel
// roots and their configured check projections.
type WorkspaceCheckCatalog struct {
	CatalogEnvelope
	Workspaces []WorkspaceRecord `json:"workspaces"`
}

// Validate checks workspace-check-specific invariants and the envelope.
func (catalog WorkspaceCheckCatalog) Validate() error {
	if catalog.Workspaces == nil {
		return fmt.Errorf("workspaces must be a non-null array")
	}
	itemIDs := make([]string, 0, len(catalog.Workspaces))
	seenPath := map[string]bool{}
	for _, workspace := range catalog.Workspaces {
		if workspace.ID == "" || workspace.Path == "" ||
			workspace.ModulePath == "" || workspace.ModuleName == "" {
			return fmt.Errorf("workspace record %q is incomplete", workspace.ID)
		}
		if seenPath[workspace.Path] {
			return fmt.Errorf("duplicate workspace path %q", workspace.Path)
		}
		seenPath[workspace.Path] = true
		itemIDs = append(itemIDs, "workspace."+workspace.ID)
	}
	return catalog.CatalogEnvelope.Validate(itemIDs)
}

// ItemIDs returns the sorted set-like identity list for a workspace catalog.
func (catalog WorkspaceCheckCatalog) ItemIDs() []string {
	ids := make([]string, 0, len(catalog.Workspaces))
	for _, workspace := range catalog.Workspaces {
		ids = append(ids, "workspace."+workspace.ID)
	}
	return ids
}

// CanonicalJSONWorkspaceCheck encodes a complete workspace-check catalog with
// a digest over the canonical bytes with the digest field omitted.
func CanonicalJSONWorkspaceCheck(catalog WorkspaceCheckCatalog) ([]byte, error) {
	if err := catalog.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := catalog
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode workspace-check catalog: %w", err)
	}
	withoutDigest.Digest = digest(content)
	contentWithDigest, err := json.MarshalIndent(withoutDigest, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("encode workspace-check catalog with digest: %w", err)
	}
	return append(contentWithDigest, '\n'), nil
}

// VerifySelfDigest recomputes the content-addressed digest over the canonical
// bytes of the complete workspace-check document and compares it to the
// stored digest, rejecting an empty or mismatched digest.
func (catalog WorkspaceCheckCatalog) VerifySelfDigest() error {
	withoutDigest := catalog
	withoutDigest.Digest = ""
	return verifySelfDigest(catalog.Digest, withoutDigest)
}

// DecodeWorkspaceCheckStrict decodes and validates a complete workspace-check
// catalog document with unknown-field and trailing-JSON rejection.
func DecodeWorkspaceCheckStrict(content []byte) (WorkspaceCheckCatalog, error) {
	var catalog WorkspaceCheckCatalog
	if err := DecodeStrict(content, &catalog); err != nil {
		return WorkspaceCheckCatalog{}, err
	}
	if err := catalog.Validate(); err != nil {
		return WorkspaceCheckCatalog{}, err
	}
	if err := catalog.VerifySelfDigest(); err != nil {
		return WorkspaceCheckCatalog{}, err
	}
	return catalog, nil
}

// GoalCoarseIdentity is the validated stable goal identity.
type GoalCoarseIdentity struct {
	Name      string `json:"name"`
	OwnerRoot string `json:"ownerRoot"`
	Scope     string `json:"scope"`
}

// GoalCoarseStatus is the validated coarse goal state.
type GoalCoarseStatus struct {
	Outcome   string `json:"outcome"`
	Execution string `json:"execution"`
}

// GoalContinuation is the bounded resume projection for an open goal. It
// carries the exact structured continuation state a fresh agent needs to
// resume without free-form archaeology. It never embeds plan, result, or
// evidence bytes; those live only in the goal record.
type GoalContinuation struct {
	ActiveAttempt    string   `json:"activeAttempt,omitempty"`
	StableDefect     string   `json:"stableDefect,omitempty"`
	Hypothesis       string   `json:"hypothesis,omitempty"`
	Subject          string   `json:"subject,omitempty"`
	AffectedCriteria []string `json:"affectedCriteria,omitempty"`
	RegressionRefs   []string `json:"regressionRefs,omitempty"`
	PriorAttemptID   string   `json:"priorAttemptID,omitempty"`
	DominantFailure  string   `json:"dominantFailure,omitempty"`
	MeasurableDelta  string   `json:"measurableDelta,omitempty"`
	NextAction       string   `json:"nextAction,omitempty"`
	Blocker          string   `json:"blocker,omitempty"`
	ResumeCondition  string   `json:"resumeCondition,omitempty"`
}

// GoalRecord is one eligible goal directory projection.
type GoalRecord struct {
	CandidatePath string              `json:"candidatePath"`
	Availability  string              `json:"availability"`
	Reason        string              `json:"reason,omitempty"`
	Identity      *GoalCoarseIdentity `json:"identity,omitempty"`
	CoarseStatus  *GoalCoarseStatus   `json:"coarseStatus,omitempty"`
	Continuation  *GoalContinuation   `json:"continuation,omitempty"`
}

// GoalCatalog is the bounded goal projection over the registered owner-local
// goals roots. Only completely validated records receive identity and coarse
// status; invalid or interrupted records are available only as eligible
// unavailable candidates.
type GoalCatalog struct {
	CatalogEnvelope
	Goals []GoalRecord `json:"goals"`
}

// Validate checks goal-specific invariants and the shared envelope.
func (catalog GoalCatalog) Validate() error {
	if catalog.Goals == nil {
		return fmt.Errorf("goals must be a non-null array")
	}
	itemIDs := make([]string, 0, len(catalog.Goals))
	for _, goal := range catalog.Goals {
		if goal.CandidatePath == "" || absPath(goal.CandidatePath) ||
			escapes(goal.CandidatePath) {
			return fmt.Errorf("goal record has malformed candidate path %q",
				goal.CandidatePath)
		}
		if !oneOf(goal.Availability, "available", "unavailable") {
			return fmt.Errorf("goal %s has unknown availability %q",
				goal.CandidatePath, goal.Availability)
		}
		needsIdentity := goal.Availability == "available"
		if needsIdentity && (goal.Identity == nil || goal.CoarseStatus == nil) {
			return fmt.Errorf("available goal %s lacks identity or coarse status",
				goal.CandidatePath)
		}
		if !needsIdentity && (goal.Identity != nil || goal.CoarseStatus != nil) {
			return fmt.Errorf("unavailable goal %s carries unvalidated identity",
				goal.CandidatePath)
		}
		if !needsIdentity && goal.Reason == "" {
			return fmt.Errorf("unavailable goal %s lacks a reason",
				goal.CandidatePath)
		}
		if goal.Continuation != nil && !needsIdentity {
			return fmt.Errorf("unavailable goal %s carries a continuation",
				goal.CandidatePath)
		}
		if goal.Continuation != nil {
			if goal.CoarseStatus == nil ||
				goal.CoarseStatus.Outcome != "open" {
				return fmt.Errorf("goal %s carries a continuation without open outcome",
					goal.CandidatePath)
			}
			if err := validateGoalContinuation(*goal.Continuation); err != nil {
				return fmt.Errorf("goal %s: %w", goal.CandidatePath, err)
			}
		}
		itemIDs = append(itemIDs, "goal."+goal.CandidatePath)
	}
	return catalog.CatalogEnvelope.Validate(itemIDs)
}

func validateGoalContinuation(continuation GoalContinuation) error {
	if continuation.ActiveAttempt == "" {
		return fmt.Errorf("continuation lacks an active attempt")
	}
	for _, field := range []struct {
		name  string
		value string
	}{
		{"stableDefect", continuation.StableDefect},
		{"hypothesis", continuation.Hypothesis},
		{"subject", continuation.Subject},
		{"dominantFailure", continuation.DominantFailure},
		{"measurableDelta", continuation.MeasurableDelta},
		{"nextAction", continuation.NextAction},
		{"blocker", continuation.Blocker},
		{"resumeCondition", continuation.ResumeCondition},
	} {
		value := field.value
		if value == "" || len(value) <= 4096 {
			continue
		}
		return fmt.Errorf(
			"continuation %s exceeds 4096 bytes",
			field.name,
		)
	}
	for _, field := range []struct {
		name   string
		values []string
	}{
		{"affectedCriteria", continuation.AffectedCriteria},
		{"regressionRefs", continuation.RegressionRefs},
	} {
		previous := ""
		for _, value := range field.values {
			if value <= previous || len(value) > 64 {
				return fmt.Errorf(
					"continuation %s must be unique and sorted",
					field.name,
				)
			}
			previous = value
		}
	}
	return nil
}

// ItemIDs returns the sorted set-like identity list for a goal catalog.
func (catalog GoalCatalog) ItemIDs() []string {
	ids := make([]string, 0, len(catalog.Goals))
	for _, goal := range catalog.Goals {
		ids = append(ids, "goal."+goal.CandidatePath)
	}
	return ids
}

// CanonicalJSONGoal encodes a complete goal catalog with a digest over the
// canonical bytes with the digest field omitted.
func CanonicalJSONGoal(catalog GoalCatalog) ([]byte, error) {
	if err := catalog.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := catalog
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode goal catalog: %w", err)
	}
	withoutDigest.Digest = digest(content)
	contentWithDigest, err := json.MarshalIndent(withoutDigest, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("encode goal catalog with digest: %w", err)
	}
	return append(contentWithDigest, '\n'), nil
}

// VerifySelfDigest recomputes the content-addressed digest over the canonical
// bytes of the complete goal document and compares it to the stored digest,
// rejecting an empty or mismatched digest.
func (catalog GoalCatalog) VerifySelfDigest() error {
	withoutDigest := catalog
	withoutDigest.Digest = ""
	return verifySelfDigest(catalog.Digest, withoutDigest)
}

// DecodeGoalStrict decodes and validates a complete goal catalog document
// with unknown-field and trailing-JSON rejection.
func DecodeGoalStrict(content []byte) (GoalCatalog, error) {
	var catalog GoalCatalog
	if err := DecodeStrict(content, &catalog); err != nil {
		return GoalCatalog{}, err
	}
	if err := catalog.Validate(); err != nil {
		return GoalCatalog{}, err
	}
	if err := catalog.VerifySelfDigest(); err != nil {
		return GoalCatalog{}, err
	}
	return catalog, nil
}

// IndexCatalogDescriptor references one catalog without embedding its body.
type IndexCatalogDescriptor struct {
	ID                string   `json:"id"`
	Kind              string   `json:"kind"`
	Schema            string   `json:"schema"`
	DerivationVersion string   `json:"derivationVersion"`
	Digest            string   `json:"digest"`
	InputDigests      []string `json:"inputDigests"`
	Completeness      string   `json:"completeness"`
	Limitations       []string `json:"limitations"`
	QueryRoutes       []string `json:"queryRoutes"`
}

// AgentSystemIndex is the bounded inventory of catalog descriptors and
// conflicts. It embeds no catalog payloads or bodies.
type AgentSystemIndex struct {
	CatalogEnvelope
	Catalogs []IndexCatalogDescriptor `json:"catalogs"`
}

// Validate checks index-specific invariants and the shared envelope.
func (index AgentSystemIndex) Validate() error {
	if index.Catalogs == nil {
		return fmt.Errorf("index arrays must be non-null")
	}
	itemIDs := make([]string, 0, len(index.Catalogs)+len(index.CatalogEnvelope.Conflicts))
	seenID := map[string]bool{}
	for _, catalog := range index.Catalogs {
		if catalog.ID == "" || catalog.Kind == "" ||
			catalog.Schema == "" || catalog.DerivationVersion == "" ||
			catalog.Digest == "" || len(catalog.InputDigests) == 0 {
			return fmt.Errorf("index catalog descriptor %q is incomplete", catalog.ID)
		}
		if seenID[catalog.ID] {
			return fmt.Errorf("duplicate index catalog descriptor %q", catalog.ID)
		}
		seenID[catalog.ID] = true
		if !digestPattern.MatchString(catalog.Digest) {
			return fmt.Errorf("index descriptor %s has malformed digest %q",
				catalog.ID, catalog.Digest)
		}
		for _, inputDigest := range catalog.InputDigests {
			if !digestPattern.MatchString(inputDigest) {
				return fmt.Errorf("index descriptor %s has malformed input digest %q",
					catalog.ID, inputDigest)
			}
		}
		itemIDs = append(itemIDs, "catalog."+catalog.ID)
	}
	return index.CatalogEnvelope.Validate(itemIDs)
}

// ItemIDs returns the sorted set-like identity list for an agent system index.
func (index AgentSystemIndex) ItemIDs() []string {
	ids := make([]string, 0, len(index.Catalogs)+len(index.CatalogEnvelope.Conflicts))
	for _, catalog := range index.Catalogs {
		ids = append(ids, "catalog."+catalog.ID)
	}
	for _, conflict := range index.CatalogEnvelope.Conflicts {
		ids = append(ids, "conflict."+conflict.ID)
	}
	return ids
}

// CanonicalJSONIndex encodes a complete agent system index with a digest over
// the canonical bytes with the digest field omitted.
func CanonicalJSONIndex(index AgentSystemIndex) ([]byte, error) {
	if err := index.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := index
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode agent system index: %w", err)
	}
	withoutDigest.Digest = digest(content)
	contentWithDigest, err := json.MarshalIndent(withoutDigest, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("encode agent system index with digest: %w", err)
	}
	return append(contentWithDigest, '\n'), nil
}

// VerifySelfDigest recomputes the content-addressed digest over the canonical
// bytes of the complete agent system index and compares it to the stored
// digest, rejecting an empty or mismatched digest.
func (index AgentSystemIndex) VerifySelfDigest() error {
	withoutDigest := index
	withoutDigest.Digest = ""
	return verifySelfDigest(index.Digest, withoutDigest)
}

// DecodeIndexStrict decodes and validates a complete agent system index
// document with unknown-field and trailing-JSON rejection.
func DecodeIndexStrict(content []byte) (AgentSystemIndex, error) {
	var index AgentSystemIndex
	if err := DecodeStrict(content, &index); err != nil {
		return AgentSystemIndex{}, err
	}
	if err := index.Validate(); err != nil {
		return AgentSystemIndex{}, err
	}
	if err := index.VerifySelfDigest(); err != nil {
		return AgentSystemIndex{}, err
	}
	return index, nil
}
