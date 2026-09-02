// Command agent_system renders one bounded offline context capsule.
//
// It is a join over the checked owner-local catalogs and optional runtime
// observations. It performs no network or stateful operation, never mutates
// source, and does not depend on Cordis or MCP. A missing optional provider
// or catalog becomes a structured unavailable field; it never fails the
// whole capsule.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	"git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

type options struct {
	workspaceRoot string
	path          string
	taskID        string
	repository    string
	json          bool
	markdown      bool
	revision      string
	dirty         bool
	producerRef   string
}

func parseFlags(args []string) (options, error) {
	var opts options
	flags := flag.NewFlagSet("agent_system", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required; defaults to BUILD_WORKSPACE_DIRECTORY)")
	flags.StringVar(&opts.path, "path", ".",
		"path or label relative to the workspace root to orient on")
	flags.StringVar(&opts.taskID, "task", "",
		"optional task/session identity to bind")
	flags.StringVar(&opts.repository, "repository", "alwaldend/src",
		"repository identity")
	flags.BoolVar(&opts.json, "json", true, "emit the portable JSON capsule")
	flags.BoolVar(&opts.markdown, "markdown", false,
		"emit the human Markdown projection instead of JSON")
	flags.StringVar(&opts.revision, "revision", "",
		"exact Git revision to bind (default: content-derived source digest)")
	flags.BoolVar(&opts.dirty, "dirty-inputs", true,
		"whether the source inputs are dirty (uncommitted)")
	flags.StringVar(&opts.producerRef, "producer-ref", "repository.agent-system",
		"producer reference for the capsule")
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if opts.workspaceRoot == "" {
		opts.workspaceRoot = os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	}
	if opts.workspaceRoot == "" {
		return options{}, fmt.Errorf("--workspace-root is required")
	}
	if opts.json && opts.markdown {
		return options{}, fmt.Errorf("--json and --markdown are mutually exclusive")
	}
	return opts, nil
}

// catalogSnapshot is the bounded read of all checked catalog projections.
// Missing optional catalogs are structured unavailability.
type catalogSnapshot struct {
	topology       *catalogv1alpha1.TopologyCatalog
	policy         *catalogv1alpha1.PolicyCatalog
	action         *catalogv1alpha1.ActionCatalog
	capability     *catalogv1alpha1.CapabilityCatalog
	workspaceCheck *catalogv1alpha1.WorkspaceCheckCatalog
	goal           *catalogv1alpha1.GoalCatalog
	index          *catalogv1alpha1.AgentSystemIndex
	limitations    []string
}

func (c *capsuleBuilder) loadCatalogs() *catalogSnapshot {
	snapshot := &catalogSnapshot{}
	if value, err := c.readTopology(); err == nil {
		snapshot.topology = value
	} else {
		snapshot.limitations = append(snapshot.limitations,
			"topology catalog unavailable: "+bounded(err))
	}
	if value, err := c.readPolicy(); err == nil {
		snapshot.policy = value
	} else {
		snapshot.limitations = append(snapshot.limitations,
			"policy catalog unavailable: "+bounded(err))
	}
	if value, err := c.readAction(); err == nil {
		snapshot.action = value
	} else {
		snapshot.limitations = append(snapshot.limitations,
			"action catalog unavailable: "+bounded(err))
	}
	if value, err := c.readCapability(); err == nil {
		snapshot.capability = value
	} else {
		snapshot.limitations = append(snapshot.limitations,
			"capability catalog unavailable: "+bounded(err))
	}
	if value, err := c.readWorkspaceCheck(); err == nil {
		snapshot.workspaceCheck = value
	} else {
		snapshot.limitations = append(snapshot.limitations,
			"workspace-check catalog unavailable: "+bounded(err))
	}
	if value, err := c.readGoal(); err == nil {
		snapshot.goal = value
	} else {
		snapshot.limitations = append(snapshot.limitations,
			"goal catalog unavailable: "+bounded(err))
	}
	if value, err := c.readIndex(); err == nil {
		snapshot.index = value
	} else {
		snapshot.limitations = append(snapshot.limitations,
			"index catalog unavailable: "+bounded(err))
	}
	return snapshot
}

// capsuleBuilder collects the bounded inputs and builds the capsule.
type capsuleBuilder struct {
	root   string
	opts   options
	inputs []v1alpha1.InputReference
}

func (c *capsuleBuilder) input(path, role, digestValue string) {
	c.inputs = append(c.inputs, v1alpha1.InputReference{
		Reference: v1alpha1.Reference{
			Kind:   v1alpha1.ReferenceRepository,
			ID:     c.opts.repository,
			Digest: digestValue,
		},
		Role: role,
	})
	_ = path
}

func (c *capsuleBuilder) resolvePath() string {
	target := filepath.Clean(filepath.Join(c.root, filepath.FromSlash(c.opts.path)))
	relative, err := filepath.Rel(c.root, target)
	if err != nil {
		return "."
	}
	if strings.HasPrefix(relative, "..") || relative == ".." {
		return "."
	}
	return filepath.ToSlash(relative)
}

// applicableWorkspace is the longest tracked workspace root containing the
// path, or the root workspace.
func (c *capsuleBuilder) applicableWorkspace(snapshot *catalogSnapshot) *catalogv1alpha1.WorkspaceRecord {
	path := c.resolvePath()
	var best *catalogv1alpha1.WorkspaceRecord
	bestLen := -1
	if snapshot.workspaceCheck != nil {
		for index := range snapshot.workspaceCheck.Workspaces {
			workspace := &snapshot.workspaceCheck.Workspaces[index]
			candidate := workspace.Path
			if candidate == "." {
				candidate = ""
			}
			if candidate == "" || path == candidate || strings.HasPrefix(path, candidate+"/") {
				if len(candidate) > bestLen {
					best = workspace
					bestLen = len(candidate)
				}
			}
		}
	}
	if best == nil {
		best = &catalogv1alpha1.WorkspaceRecord{
			ID:   "root",
			Path: ".",
			Phases: []catalogv1alpha1.WorkspacePhase{
				{
					ID:              "root.check",
					ProviderRef:     "repository.bazel-operations",
					CommandTemplate: "bazel_agent test //...",
				},
			},
		}
	}
	return best
}

// policyRecords returns the policy records whose prefix covers the path.
func (c *capsuleBuilder) policyRecords(snapshot *catalogSnapshot, path string) []catalogv1alpha1.PolicyRecord {
	if snapshot.policy == nil {
		return nil
	}
	var records []catalogv1alpha1.PolicyRecord
	for _, record := range snapshot.policy.Policies {
		prefix := record.PathPrefix
		if prefix == "/" || prefix == "" {
			records = append(records, record)
			continue
		}
		trimmed := strings.TrimPrefix(prefix, "/")
		if path == trimmed || strings.HasPrefix(path, trimmed+"/") {
			records = append(records, record)
		}
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].Precedence != records[j].Precedence {
			return records[i].Precedence < records[j].Precedence
		}
		return records[i].ID < records[j].ID
	})
	return records
}

// documentSources aggregates the applicable instruction and owner documents.
func (c *capsuleBuilder) documentSources(
	snapshot *catalogSnapshot,
	workspace *catalogv1alpha1.WorkspaceRecord,
) []v1alpha1.CapsuleDocumentSource {
	var documents []v1alpha1.CapsuleDocumentSource
	if digest := c.digestFile("AGENTS.md"); digest != "" {
		documents = append(documents, v1alpha1.CapsuleDocumentSource{
			Path:   "AGENTS.md",
			Digest: digest,
		})
	}
	if snapshot.topology != nil {
		for _, component := range snapshot.topology.Components {
			if workspace == nil || component.Path == workspace.Path {
				documents = append(documents, v1alpha1.CapsuleDocumentSource{
					Path:   component.OwnerReadme,
					Digest: c.digestFile(component.OwnerReadme),
				})
			}
		}
	}
	return documents
}

func (c *capsuleBuilder) digestFile(relative string) string {
	content, err := os.ReadFile(filepath.Join(c.root, filepath.FromSlash(relative)))
	if err != nil {
		return ""
	}
	value := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(value[:])
}

// capabilities joins provider/action and skill records relevant to the path.
func (c *capsuleBuilder) capabilities(
	snapshot *catalogSnapshot,
	workspace *catalogv1alpha1.WorkspaceRecord,
) []v1alpha1.CapsuleCapability {
	result := []v1alpha1.CapsuleCapability{}
	if snapshot.action == nil {
		return result
	}
	for _, provider := range snapshot.action.Providers {
		owner := provider.Owner
		if workspace != nil && workspace.Path != "." &&
			!strings.HasPrefix(owner, workspace.Path) {
			continue
		}
		result = append(result, v1alpha1.CapsuleCapability{
			ID:    provider.ID,
			Kind:  "provider",
			Owner: provider.Owner,
			Role:  "provider",
		})
	}
	if snapshot.capability != nil {
		for _, skill := range snapshot.capability.Skills {
			owner := strings.ReplaceAll(skill.Owner, "/", ".")
			if workspace != nil && workspace.Path != "." &&
				!strings.HasPrefix(owner, workspace.Path) {
				continue
			}
			capability := v1alpha1.CapsuleCapability{
				ID:      skill.ID,
				Kind:    "skill",
				Owner:   skill.Owner,
				Role:    skill.Layer,
				Effects: append([]string(nil), skill.CapabilityRefs...),
				Cost:    skill.ContextCost,
				Access:  skill.Activation,
			}
			if len(skill.Exclusions) > 0 {
				capability.Excluded = strings.Join(skill.Exclusions, "; ")
			}
			result = append(result, capability)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Kind != result[j].Kind {
			return result[i].Kind < result[j].Kind
		}
		return result[i].ID < result[j].ID
	})
	return result
}

// checks joins the workspace-check phases for the applicable workspace.
func (c *capsuleBuilder) checks(
	snapshot *catalogSnapshot,
	workspace *catalogv1alpha1.WorkspaceRecord,
) []v1alpha1.CapsuleCheck {
	if workspace == nil || snapshot.workspaceCheck == nil {
		return nil
	}
	result := make([]v1alpha1.CapsuleCheck, 0, len(workspace.Phases))
	for _, phase := range workspace.Phases {
		result = append(result, v1alpha1.CapsuleCheck{
			WorkspaceID:     workspace.ID,
			PhaseID:         phase.ID,
			ProviderRef:     phase.ProviderRef,
			CommandTemplate: phase.CommandTemplate,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].PhaseID < result[j].PhaseID
	})
	return result
}

// providerStatus produces structured runtime/provider observations. It never
// contacts a provider; observations are bounded and unavailable when not
// passed in.
func (c *capsuleBuilder) providerStatus(
	snapshot *catalogSnapshot,
) []v1alpha1.CapsuleProviderStatus {
	var result []v1alpha1.CapsuleProviderStatus
	if snapshot.index != nil {
		for _, descriptor := range snapshot.index.Catalogs {
			result = append(result, v1alpha1.CapsuleProviderStatus{
				ProviderID:       descriptor.ID,
				CatalogETag:      descriptor.Digest,
				DesiredRevision:  descriptor.DerivationVersion,
				ObservedRevision: descriptor.DerivationVersion,
				State:            string(descriptor.Completeness),
			})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ProviderID < result[j].ProviderID
	})
	if len(result) == 0 {
		result = append(result, v1alpha1.CapsuleProviderStatus{
			ProviderID:  "catalogs",
			State:       "unavailable",
			Unavailable: "no index descriptors observed",
		})
	}
	return result
}

// goalBinding returns the active owner-local goal, if any.
func (c *capsuleBuilder) goalBinding(snapshot *catalogSnapshot) string {
	if snapshot.goal == nil {
		return ""
	}
	for _, goal := range snapshot.goal.Goals {
		if goal.Availability == "available" && goal.CoarseStatus != nil &&
			goal.CoarseStatus.Execution == "active" {
			return goal.CandidatePath
		}
	}
	return ""
}

// nextActions are the safe, bounded discovery actions for the caller.
func (c *capsuleBuilder) nextActions(
	snapshot *catalogSnapshot,
	workspace *catalogv1alpha1.WorkspaceRecord,
) []string {
	actions := []string{}
	actions = append(actions,
		"inspect the applicable AGENTS.md and owner README for exact authority")
	if workspace != nil {
		actions = append(actions,
			"run the applicable workspace check phase: "+workspace.Phases[0].CommandTemplate)
	}
	if snapshot.goal != nil && len(snapshot.goal.Goals) > 0 {
		actions = append(actions, "review the active goal record for continuation state")
	}
	for _, limitation := range snapshot.limitations {
		actions = append(actions, "resolve unavailable catalog: "+limitation)
	}
	if len(actions) == 0 {
		actions = append(actions, "read the architecture roadmap for the next phase")
	}
	return actions
}

func (c *capsuleBuilder) sourceDigest(snapshot *catalogSnapshot) string {
	hasher := sha256.New()
	for _, input := range c.inputs {
		hasher.Write([]byte(input.Reference.ID))
		hasher.Write([]byte{0})
		hasher.Write([]byte(input.Reference.Digest))
		hasher.Write([]byte{0})
	}
	if snapshot.index != nil {
		hasher.Write([]byte(snapshot.index.Digest))
	}
	return "sha256:" + hex.EncodeToString(hasher.Sum(nil))
}

func (c *capsuleBuilder) build(snapshot *catalogSnapshot, observedAt time.Time) (v1alpha1.ContextCapsule, error) {
	path := c.resolvePath()
	workspace := c.applicableWorkspace(snapshot)
	revision := c.opts.revision
	if revision == "" {
		revision = c.sourceDigest(snapshot)
	}
	component := v1alpha1.CapsuleComponent{
		Path:         path,
		Workspace:    workspace.ID,
		ReviewOwners: "@simeonwarren",
	}
	if snapshot.topology != nil {
		for _, record := range snapshot.topology.Components {
			if record.Path == path || strings.HasPrefix(path, record.Path+"/") {
				component.ComponentID = record.ID
				component.OwnerReadme = record.OwnerReadme
				component.Lifecycle = record.Lifecycle
				component.Classification = record.DocsState
				break
			}
		}
	}
	documents := c.documentSources(snapshot, workspace)
	capabilities := c.capabilities(snapshot, workspace)
	closeChecks := c.checks(snapshot, workspace)
	providers := c.providerStatus(snapshot)
	goal := c.goalBinding(snapshot)

	byteSize := c.estimateByteSize(snapshot, documents, capabilities, closeChecks, providers)

	capsule := v1alpha1.ContextCapsule{
		APIVersion: v1alpha1.APIVersion,
		Kind:       v1alpha1.KindContextCapsule,
		Identity: v1alpha1.CapsuleIdentity{
			Repository:    c.opts.repository,
			WorkspaceRoot: c.opts.workspaceRoot,
			WorktreePath:  c.opts.path,
			Revision:      revision,
			DirtyInputs:   c.opts.dirty,
			SourceDigest:  revision,
			InputDigest:   c.sourceDigest(snapshot),
			ByteSize:      byteSize,
		},
		Task: v1alpha1.CapsuleTask{
			TaskID: c.opts.taskID,
		},
		Outcome: v1alpha1.CapsuleOutcome{
			GoalBinding: goal,
		},
		Documents:    documents,
		Component:    component,
		Capabilities: capabilities,
		Checks:       closeChecks,
		Providers:    providers,
		Provenance: v1alpha1.CapsuleProvenance{
			ObservedAt:   observedAt.UTC().Format(time.RFC3339Nano),
			Freshness:    "fresh",
			Completeness: v1alpha1.CompletenessComplete,
			Limitations:  append([]string(nil), snapshot.limitations...),
			NextActions:  c.nextActions(snapshot, workspace),
			AuthorizedBy: "repository.agent-system",
		},
	}
	if len(snapshot.limitations) > 0 {
		capsule.Provenance.Completeness = v1alpha1.CompletenessPartial
	}
	return capsule, nil
}

func (c *capsuleBuilder) estimateByteSize(
	snapshot *catalogSnapshot,
	documents []v1alpha1.CapsuleDocumentSource,
	capabilities []v1alpha1.CapsuleCapability,
	checks []v1alpha1.CapsuleCheck,
	providers []v1alpha1.CapsuleProviderStatus,
) int64 {
	var size int64
	if snapshot.topology != nil {
		size += int64(len(snapshot.topology.Components) * 32)
	}
	for _, document := range documents {
		size += int64(len(document.Path) + len(document.Digest))
	}
	for _, capability := range capabilities {
		size += int64(len(capability.ID) + len(capability.Kind) + len(capability.Owner))
	}
	for _, check := range checks {
		size += int64(len(check.PhaseID) + len(check.CommandTemplate))
	}
	for _, provider := range providers {
		size += int64(len(provider.ProviderID) + len(provider.State))
	}
	return size + 512
}

func bounded(err error) string {
	if err == nil {
		return "unavailable"
	}
	message := strings.TrimSpace(err.Error())
	if len(message) > 200 {
		message = message[:200]
	}
	return message
}

func run(args []string, stdout io.Writer) error {
	if len(args) > 0 && args[0] == "plan" {
		return runPlan(args[1:], stdout)
	}
	if len(args) > 0 && args[0] == "normalize" {
		return runNormalize(args[1:], stdout)
	}
	if len(args) > 0 && args[0] == "coverage" {
		return runCoverage(args[1:], stdout)
	}
	opts, err := parseFlags(args)
	if err != nil {
		return err
	}
	root, err := filepath.Abs(opts.workspaceRoot)
	if err != nil {
		return err
	}
	builder := &capsuleBuilder{root: root, opts: opts}
	snapshot := builder.loadCatalogs()
	observedAt := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	capsule, err := builder.build(snapshot, observedAt)
	if err != nil {
		return err
	}
	if err := builder.setCapsuleID(&capsule); err != nil {
		return err
	}
	if opts.markdown {
		_, err := io.WriteString(stdout, v1alpha1.RenderContextMarkdown(capsule))
		return err
	}
	content, err := v1alpha1.CanonicalContextJSON(capsule)
	if err != nil {
		return err
	}
	_, err = stdout.Write(content)
	return err
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "agent_system:", err)
		os.Exit(1)
	}
}

func (c *capsuleBuilder) setCapsuleID(capsule *v1alpha1.ContextCapsule) error {
	if capsule.ID != "" {
		return nil
	}
	id, err := v1alpha1.CapsuleID(*capsule)
	if err != nil {
		return err
	}
	capsule.ID = id
	return nil
}
