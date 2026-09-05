// Command agent_system renders one bounded offline context capsule.
//
// It joins checked owner-local catalogs with applicable source documents.
// It performs no network or stateful operation, never mutates
// source, and does not depend on Cordis or MCP. A missing optional provider
// or catalog becomes a structured unavailable field; it never fails the
// whole capsule.
package agent_system

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
	dirtyDeclared bool
	producerRef   string
	maxBytes      int
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
		"caller-declared revision override (default: observed Git HEAD or input digest)")
	flags.BoolVar(&opts.dirty, "dirty-inputs", true,
		"caller-declared dirty source state (not verified)")
	flags.StringVar(&opts.producerRef, "producer-ref", "repository.agent-system",
		"producer reference for the capsule")
	flags.IntVar(&opts.maxBytes, "max-bytes", 65536,
		"maximum serialized output bytes (1024 to 1048576)")
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	flags.Visit(func(value *flag.Flag) {
		if value.Name == "dirty-inputs" {
			opts.dirtyDeclared = true
		}
	})
	if opts.workspaceRoot == "" {
		opts.workspaceRoot = os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	}
	if opts.workspaceRoot == "" {
		return options{}, fmt.Errorf("--workspace-root is required")
	}
	if flags.NArg() != 0 {
		return options{}, fmt.Errorf("unexpected positional arguments")
	}
	if opts.maxBytes < 1024 || opts.maxBytes > 1048576 {
		return options{}, fmt.Errorf("--max-bytes must be between 1024 and 1048576")
	}
	opts.markdown = opts.markdown || !opts.json
	if strings.HasPrefix(opts.path, "//") {
		opts.path = strings.SplitN(strings.TrimPrefix(opts.path, "//"), ":", 2)[0]
	} else if strings.HasPrefix(opts.path, "@") {
		return options{}, fmt.Errorf("external Bazel labels are not workspace paths")
	}
	opts.path = filepath.ToSlash(filepath.Clean(opts.path))
	if filepath.IsAbs(opts.path) || opts.path == ".." || strings.HasPrefix(opts.path, "../") {
		return options{}, fmt.Errorf("--path must stay inside the workspace")
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
		Role: role + ":" + path,
	})
}

func (c *capsuleBuilder) resolvePath() string {
	return c.opts.path
}

// resolveExistingPrefix follows symlinks in the existing part of a path.
// Missing descendants remain useful for new-file tasks, but an existing
// dangling or cyclic symlink is an error rather than guessed ownership.
func resolveExistingPrefix(path string) (string, error) {
	var suffix []string
	for depth := 0; depth < 256; depth++ {
		_, err := os.Lstat(path)
		if err == nil {
			resolved, err := filepath.EvalSymlinks(path)
			if err != nil {
				return "", err
			}
			for index := len(suffix) - 1; index >= 0; index-- {
				resolved = filepath.Join(resolved, suffix[index])
			}
			return resolved, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(path)
		if parent == path {
			return "", err
		}
		suffix = append(suffix, filepath.Base(path))
		path = parent
	}
	return "", fmt.Errorf("path exceeds the 256-component resolution limit")
}

func resolveContextPath(root, path string) (string, string, error) {
	resolvedRoot, err := resolveExistingPrefix(root)
	if err != nil {
		return "", "", fmt.Errorf("resolve workspace root: %w", err)
	}
	resolvedPath, err := resolveExistingPrefix(filepath.Join(resolvedRoot, filepath.FromSlash(path)))
	if err != nil {
		return "", "", fmt.Errorf("resolve requested path: %w", err)
	}
	relative, err := filepath.Rel(resolvedRoot, resolvedPath)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("resolved --path must stay inside the workspace")
	}
	return resolvedRoot, filepath.ToSlash(relative), nil
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
) ([]v1alpha1.CapsuleDocumentSource, string) {
	var documents []v1alpha1.CapsuleDocumentSource
	observed := map[string]bool{}
	reportedUnavailable := map[string]bool{}
	reportUnavailable := func(path string) {
		if !reportedUnavailable[path] {
			reportedUnavailable[path] = true
			snapshot.limitations = append(snapshot.limitations,
				"document unavailable: "+path)
		}
	}
	add := func(path string, required bool) bool {
		if path == "" {
			return false
		}
		if available, exists := observed[path]; exists {
			if required && !available {
				reportUnavailable(path)
			}
			return false
		}
		digest := c.digestFile(path)
		observed[path] = digest != ""
		if digest == "" {
			if required {
				reportUnavailable(path)
			}
			return false
		}
		documents = append(documents, v1alpha1.CapsuleDocumentSource{
			Path: path, Digest: digest,
		})
		c.input(path, "document", digest)
		return true
	}
	// Discovery may advance to an ancestor only after confirmed absence.
	// A present or inaccessible path still owns its place in the hierarchy,
	// even when its bytes cannot be safely included in this capsule.
	addDiscovered := func(path string) bool {
		_, err := os.Lstat(filepath.Join(c.root, filepath.FromSlash(path)))
		if os.IsNotExist(err) {
			return false
		}
		add(path, true)
		return true
	}
	path := c.resolvePath()
	if info, err := os.Stat(filepath.Join(c.root, filepath.FromSlash(path))); err == nil && !info.IsDir() {
		path = filepath.ToSlash(filepath.Dir(path))
	}
	var ancestors []string
	for {
		ancestors = append(ancestors, path)
		if path == "." {
			break
		}
		path = filepath.ToSlash(filepath.Dir(path))
	}
	for index := len(ancestors) - 1; index >= 0; index-- {
		candidate := filepath.ToSlash(filepath.Join(ancestors[index], "AGENTS.md"))
		if index == len(ancestors)-1 {
			add(candidate, true)
		} else {
			addDiscovered(candidate)
		}
	}
	ownerReadme := ""
	for _, ancestor := range ancestors {
		candidate := filepath.ToSlash(filepath.Join(ancestor, "README.md"))
		if addDiscovered(candidate) {
			ownerReadme = candidate
			break
		}
	}
	// These are source navigation links, not parsed build instructions.
	for _, names := range [][]string{{"BUILD.bazel", "BUILD"}, {"include.MODULE.bazel"}, {"MODULE.bazel"}} {
		found := false
		for _, ancestor := range ancestors {
			for _, name := range names {
				if addDiscovered(filepath.ToSlash(filepath.Join(ancestor, name))) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
	for _, policy := range c.policyRecords(snapshot, c.resolvePath()) {
		add(policy.AgentPolicySource, true)
		add(policy.OwnerBoundaryRef, true)
		add(policy.ReviewSource, true)
		for _, axis := range policy.Axes {
			add(axis.Source, true)
		}
	}
	addDiscovered("CODEOWNERS")
	return documents, ownerReadme
}

func (c *capsuleBuilder) digestFile(relative string) string {
	path, err := filepath.EvalSymlinks(filepath.Join(c.root, filepath.FromSlash(relative)))
	if err != nil {
		return ""
	}
	relativePath, err := filepath.Rel(c.root, path)
	if err != nil || relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return ""
	}
	content, err := os.ReadFile(path)
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
	if snapshot.action != nil {
		for _, provider := range snapshot.action.Providers {
			if !c.ownerInWorkspace(snapshot, workspace, provider.Owner) {
				continue
			}
			result = append(result, v1alpha1.CapsuleCapability{
				ID:    provider.ID,
				Kind:  "provider",
				Owner: provider.Owner,
				Role:  "provider",
			})
		}
	}
	if snapshot.capability != nil {
		for _, skill := range snapshot.capability.Skills {
			if !c.ownerInWorkspace(snapshot, workspace, skill.Owner) {
				continue
			}
			capability := v1alpha1.CapsuleCapability{
				ID:        skill.ID,
				Kind:      "skill",
				Owner:     skill.Owner,
				Role:      skill.Layer,
				Cost:      skill.ContextCost,
				Access:    skill.Activation,
				Providers: append([]string(nil), skill.ProviderRequirements...),
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

// Root-workspace capabilities remain candidates in nested workspaces. A
// capability owned by another standalone workspace is outside this inventory.
// This is an ownership filter, not task-intent routing or authority.
func (c *capsuleBuilder) ownerInWorkspace(
	snapshot *catalogSnapshot,
	workspace *catalogv1alpha1.WorkspaceRecord,
	owner string,
) bool {
	if snapshot.workspaceCheck == nil {
		return true
	}
	var best *catalogv1alpha1.WorkspaceRecord
	for index := range snapshot.workspaceCheck.Workspaces {
		candidate := &snapshot.workspaceCheck.Workspaces[index]
		if candidate.Path == "." || !pathContains(candidate.Path, owner) {
			continue
		}
		if best == nil || len(candidate.Path) > len(best.Path) {
			best = candidate
		}
	}
	return best == nil || workspace != nil && best.Path == workspace.Path
}

func pathContains(parent, child string) bool {
	return parent == "." || child == parent || strings.HasPrefix(child, parent+"/")
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
	if snapshot.action != nil {
		for _, provider := range snapshot.action.Providers {
			result = append(result, v1alpha1.CapsuleProviderStatus{
				ProviderID:  provider.ID,
				State:       "unavailable",
				Unavailable: "live provider observation unavailable in offline capsule",
			})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ProviderID < result[j].ProviderID
	})
	if len(result) == 0 {
		result = append(result, v1alpha1.CapsuleProviderStatus{
			ProviderID:  "runtime",
			State:       "unavailable",
			Unavailable: "no declared providers or live observations available",
		})
	}
	return result
}

// nextActions are the safe, bounded discovery actions for the caller.
func (c *capsuleBuilder) nextActions(
	snapshot *catalogSnapshot,
	documents []v1alpha1.CapsuleDocumentSource,
) []string {
	actions := []string{}
	actions = append(actions,
		"inspect the applicable AGENTS.md and owner README for exact authority")
	buildPath := "the owner BUILD file"
	for _, document := range documents {
		if name := filepath.Base(document.Path); name == "BUILD.bazel" || name == "BUILD" {
			buildPath = document.Path
			break
		}
	}
	actions = append(actions,
		"inspect "+buildPath+" and select the narrowest checks for the requested change")
	if snapshot.goal != nil && len(snapshot.goal.Goals) > 0 {
		actions = append(actions, "select a goal explicitly if this task needs durable continuation")
	}
	for _, limitation := range snapshot.limitations {
		actions = append(actions, "inspect unavailable input: "+limitation)
	}
	return actions
}

func (c *capsuleBuilder) sourceDigest(snapshot *catalogSnapshot) string {
	hasher := sha256.New()
	for _, input := range c.inputs {
		hasher.Write([]byte(input.Reference.ID))
		hasher.Write([]byte{0})
		hasher.Write([]byte(input.Role))
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
	component := v1alpha1.CapsuleComponent{
		Path:      path,
		Workspace: workspace.ID,
	}
	if snapshot.topology != nil {
		bestLen := -1
		for _, record := range snapshot.topology.Components {
			if pathContains(record.Path, path) && len(record.Path) > bestLen {
				bestLen = len(record.Path)
				component.ComponentID = record.ID
				component.OwnerReadme = record.OwnerReadme
				component.Lifecycle = record.Lifecycle
			}
		}
		for _, tree := range snapshot.topology.Trees {
			if pathContains(tree.Path, path) {
				component.Classification = string(tree.Boundary)
			}
		}
	}
	documents, ownerReadme := c.documentSources(snapshot)
	component.OwnerReadme = ownerReadme
	capabilities := c.capabilities(snapshot, workspace)
	closeChecks := c.checks(snapshot, workspace)
	providers := c.providerStatus(snapshot)
	inputDigest := c.sourceDigest(snapshot)
	revision := c.opts.revision
	if revision == "" {
		revision = inputDigest
	}
	limitations := append([]string(nil), snapshot.limitations...)
	limitations = append(limitations,
		"catalog freshness against owning sources is unknown; catalog self-digests verify only stored bytes",
		"runtime health, task authority, goal binding, and effective CODEOWNERS are not observed",
		"capabilities are workspace candidates; task-intent routing and cross-workspace dependencies are not resolved",
	)

	capsule := v1alpha1.ContextCapsule{
		APIVersion: v1alpha1.APIVersion,
		Kind:       v1alpha1.KindContextCapsule,
		Identity: v1alpha1.CapsuleIdentity{
			Repository:        c.opts.repository,
			WorkspaceRoot:     c.root,
			WorktreePath:      path,
			Revision:          revision,
			DirtyInputs:       c.opts.dirty,
			RevisionSource:    "input-digest",
			DirtyInputsSource: "conservative-default",
			SourceDigest:      inputDigest,
			InputDigest:       inputDigest,
			ByteSize:          1,
		},
		Task: v1alpha1.CapsuleTask{
			TaskID: c.opts.taskID,
		},
		Documents:    documents,
		Component:    component,
		Capabilities: capabilities,
		Checks:       closeChecks,
		Providers:    providers,
		Provenance: v1alpha1.CapsuleProvenance{
			Inputs:       append([]v1alpha1.InputReference(nil), c.inputs...),
			ObservedAt:   observedAt.UTC().Format(time.RFC3339Nano),
			Freshness:    "unknown",
			Completeness: v1alpha1.CompletenessPartial,
			Limitations:  limitations,
			NextActions:  c.nextActions(snapshot, documents),
		},
	}
	if c.opts.revision != "" {
		capsule.Identity.RevisionSource = "caller-declared"
	}
	if c.opts.dirtyDeclared {
		capsule.Identity.DirtyInputsSource = "caller-declared"
	}
	return capsule, nil
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

func Run(args []string, stdout io.Writer) error {
	return runWithClock(args, stdout, time.Now)
}

func runWithClock(args []string, stdout io.Writer, now func() time.Time) error {
	if len(args) > 0 && args[0] == "plan" {
		return runPlan(args[1:], stdout)
	}
	if len(args) > 0 && args[0] == "normalize" {
		return runNormalize(args[1:], stdout)
	}
	if len(args) > 0 && args[0] == "coverage" {
		return runCoverage(args[1:], stdout)
	}
	if len(args) > 0 && args[0] == "aggregate" {
		return runAggregate(args[1:], stdout)
	}
	if len(args) > 0 && args[0] == "friction" {
		return runFriction(args[1:], stdout)
	}
	opts, err := parseFlags(args)
	if err != nil {
		return err
	}
	root, err := commandWorkspaceRoot(opts.workspaceRoot)
	if err != nil {
		return err
	}
	root, opts.path, err = resolveContextPath(root, opts.path)
	if err != nil {
		return err
	}
	builder := &capsuleBuilder{root: root, opts: opts}
	snapshot := builder.loadCatalogs()
	capsule, err := builder.build(snapshot, now())
	if err != nil {
		return err
	}
	applyGitObservation(&capsule, observeGit(root, now, runGitCommand))
	content, err := encodeMeasuredCapsule(&capsule)
	if err != nil {
		return err
	}
	if opts.markdown {
		content = []byte(v1alpha1.RenderContextMarkdown(capsule))
	}
	if len(content) > opts.maxBytes {
		return fmt.Errorf("capsule output is %d bytes, exceeding --max-bytes %d; narrow --path or explicitly raise the output limit", len(content), opts.maxBytes)
	}
	_, err = stdout.Write(content)
	return err
}

func main() {
	if err := Run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "agent_system:", err)
		os.Exit(1)
	}
}

// The size includes the canonical JSON, ID, byteSize field, and newline.
// Re-encode until the decimal byteSize width stops changing; ID width is fixed.
func encodeMeasuredCapsule(capsule *v1alpha1.ContextCapsule) ([]byte, error) {
	for iteration := 0; iteration < 8; iteration++ {
		capsule.ID = ""
		id, err := v1alpha1.CapsuleID(*capsule)
		if err != nil {
			return nil, err
		}
		capsule.ID = id
		content, err := v1alpha1.CanonicalContextJSON(*capsule)
		if err != nil {
			return nil, err
		}
		if capsule.Identity.ByteSize == int64(len(content)) {
			return content, nil
		}
		capsule.Identity.ByteSize = int64(len(content))
	}
	return nil, fmt.Errorf("capsule byte-size calculation did not converge")
}
