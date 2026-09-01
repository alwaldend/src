// Package plan implements the advisory Phase 3A impact planner.
//
// The planner is deterministic and offline: identical canonical inputs
// produce identical plans and digests. It never grants authority; selected
// effects, targets, and checks remain advisory and subject to goal criteria,
// admission, and agent judgment.
package plan

import (
	"fmt"
	"path"
	"sort"
	"strings"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

const (
	// KindImpactPlan identifies an advisory impact plan.
	KindImpactPlan = "ImpactPlan"

	// PlanNamespace is the identity namespace for generated plans.
	PlanNamespace = "plan"

	// CapabilityValidation is the generic validation capability bound when
	// the capability catalog is unavailable. Capability selection remains
	// advisory until admission binds a provider.
	CapabilityValidation = "validation"
)

// ImpactProfile selects the conservative validation profile.
type ImpactProfile string

const (
	// ProfileChangedFast validates only the changed packages.
	ProfileChangedFast ImpactProfile = "changed/fast"
	// ProfileWorkspace validates the applicable tracked workspace.
	ProfileWorkspace ImpactProfile = "workspace"
	// ProfileFreshEvidence validates the changed scope and binds fresh
	// evidence references.
	ProfileFreshEvidence ImpactProfile = "fresh/evidence"
	// ProfileFullAudit validates the whole repository with the most checks.
	ProfileFullAudit ImpactProfile = "full/audit"
	// ProfileDiagnose performs minimal read-only diagnosis.
	ProfileDiagnose ImpactProfile = "diagnose"
)

// PlanOptions are the canonical planner inputs.
type PlanOptions struct {
	Profile       ImpactProfile
	Intent        string
	AffectedPaths []string
	Labels        []string
	// RequestRemote explicitly requests remote.write and remote.destroy
	// effects. The plan still never grants authority: when requested, the
	// plan notes the authority requirement in Escalation.
	RequestRemote bool
}

// Index is the planner's bounded read model over the checked catalogs. It
// groups topology, capability, and workspace-check projections so the
// planner can stay deterministic while the underlying catalog documents
// evolve.
type Index struct {
	Topology       *catalogv1alpha1.TopologyCatalog
	Capability     *catalogv1alpha1.CapabilityCatalog
	WorkspaceCheck *catalogv1alpha1.WorkspaceCheckCatalog
}

// Plan computes the smallest-sufficient advisory plan for one intent.
//
// Identity, digest, selections, and checks are deterministic functions of
// the canonical inputs and the index. Unknown or unscoped inputs
// conservatively expand or refuse scope and always carry a reason in
// CoverageGaps.
func Plan(
	intent v1alpha1.Reference,
	opts PlanOptions,
	idx *Index,
) (v1alpha1.ImpactPlan, error) {
	profile, err := resolveProfile(opts.Profile)
	if err != nil {
		return v1alpha1.ImpactPlan{}, err
	}
	if strings.TrimSpace(opts.Intent) == "" {
		return v1alpha1.ImpactPlan{}, fmt.Errorf("intent identity is required")
	}
	scope := sortScope(opts.AffectedPaths, opts.Labels)
	effectiveScope := scopeForProfile(profile, scope, idx)
	targets, gaps := planTargets(profile, effectiveScope, idx)
	if len(targets) == 0 {
		return v1alpha1.ImpactPlan{}, fmt.Errorf(
			"no targets resolvable for %q: %s",
			opts.Intent,
			strings.Join(gaps, "; "),
		)
	}
	effects, forbidden := planEffects(profile, opts.RequestRemote)
	checks := planChecks(profile, targets, scope)
	cost := profileCost(profile)
	capability, reason := selectCapability(idx)
	plan := v1alpha1.ImpactPlan{
		APIVersion: v1alpha1.APIVersion,
		Kind:       KindImpactPlan,
		ID:         PlanNamespace + "/" + opts.Intent,
		IntentRef:  intent,
		Profile:    v1alpha1.ImpactProfile(profile),
		Capabilities: []v1alpha1.Reference{
			capability,
		},
		CapabilityReasons: []string{reason},
		Effects:           effects,
		ForbiddenEffects:  forbidden,
		Targets:           targets,
		Checks:            checks,
		CoverageGaps:      gaps,
		CostClass:         cost.class,
		ExpectedMaxCost:   cost.budget,
		Cacheability:      profileCacheability(profile),
		Concurrency:       int(cost.budget.Concurrency),
	}
	if profile == ProfileFreshEvidence {
		plan.EvidenceRefs = []v1alpha1.InputReference{{
			Reference: intent,
			Role:      "fresh-evidence-candidate",
		}}
	}
	if profile == ProfileDiagnose {
		plan.CoverageGaps = append(plan.CoverageGaps,
			"diagnose is read-only: execution, write, and remote coverage "+
				"are intentionally absent")
	}
	if opts.RequestRemote {
		plan.Escalation = "remote.write/remote.destroy requested: authority " +
			"requirement must be satisfied by admission or operator approval"
	}
	plan.Ordering = profileOrdering(profile)
	if err := plan.Validate(); err != nil {
		return v1alpha1.ImpactPlan{}, fmt.Errorf("invalid plan: %w", err)
	}
	canonical, err := v1alpha1.CanonicalImpactPlanJSON(plan)
	if err != nil {
		return v1alpha1.ImpactPlan{}, err
	}
	plan.Digest = v1alpha1.DigestOfCanonicalJSON(canonical)
	return plan, nil
}

func resolveProfile(value ImpactProfile) (ImpactProfile, error) {
	switch value {
	case ProfileChangedFast, ProfileWorkspace, ProfileFreshEvidence,
		ProfileFullAudit, ProfileDiagnose:
		return value, nil
	}
	return "", fmt.Errorf("unknown impact profile %q", value)
}

// sortScope canonicalizes affected paths and labels: trimmed,
// deduplicated, repository-relative, and lexically sorted.
func sortScope(paths, labels []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(paths)+len(labels))
	add := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		value = path.Clean(strings.TrimPrefix(value, "./"))
		if value == "." {
			return
		}
		if seen[value] {
			return
		}
		seen[value] = true
		result = append(result, value)
	}
	for _, value := range paths {
		add(value)
	}
	for _, value := range labels {
		add(strings.TrimPrefix(value, "//"))
	}
	sort.Strings(result)
	return result
}

// scopeForProfile expands the sorted input scope for profiles that validate
// more than the changed paths alone. The workspace profile validates one
// tracked workspace; full/audit validates the repository root.
func scopeForProfile(
	profile ImpactProfile,
	scope []string,
	idx *Index,
) []string {
	switch profile {
	case ProfileWorkspace:
		return workspaceScope(scope, idx)
	case ProfileFullAudit:
		return []string{"."}
	}
	return scope
}

// workspaceScope maps each affected path to the longest tracked workspace
// containing it, falling back to the repository root.
func workspaceScope(scope []string, idx *Index) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(scope))
	add := func(value string) {
		if value == "" {
			value = "."
		}
		if seen[value] {
			return
		}
		seen[value] = true
		result = append(result, value)
	}
	if idx != nil && idx.WorkspaceCheck != nil {
		for _, item := range scope {
			best := ""
			for _, workspace := range idx.WorkspaceCheck.Workspaces {
				prefix := workspace.Path
				if prefix == "." {
					prefix = ""
				}
				if item == prefix || (prefix != "" &&
					strings.HasPrefix(item, prefix+"/")) {
					if len(prefix) > len(best) {
						best = prefix
					}
				}
			}
			add(best)
		}
	}
	if len(result) == 0 {
		result = append(result, ".")
	}
	sort.Strings(result)
	return result
}

// planTargets maps the effective scope to Bazel labels at the narrowest
// practical scope. Ambiguous root-module, rc, toolchain, and generated
// changes conservatively expand or refuse scope with an explicit gap.
func planTargets(
	profile ImpactProfile,
	effectiveScope []string,
	idx *Index,
) ([]v1alpha1.PlanTarget, []string) {
	var targets []v1alpha1.PlanTarget
	var gaps []string
	seen := map[string]bool{}
	forceRoot := false
	add := func(target v1alpha1.PlanTarget) {
		if seen[target.Label] {
			return
		}
		seen[target.Label] = true
		targets = append(targets, target)
	}
	for _, item := range effectiveScope {
		switch {
		case item == ".":
			add(v1alpha1.PlanTarget{
				Label:      "//...",
				Path:       ".",
				AffectedBy: profileAffectedBy(profile),
			})
		case scopeExpansionReason(item) != "":
			gaps = append(gaps, scopeExpansionReason(item))
			forceRoot = true
		case isUnknownGeneratedPath(item):
			gaps = append(gaps,
				"unknown-generated scope "+item+": conservatively refused")
		case strings.Contains(item, ":"):
			add(v1alpha1.PlanTarget{
				Label:      "//" + item,
				AffectedBy: profileAffectedBy(profile),
			})
		default:
			label := "//" + item + ":all"
			component := componentForItem(item, idx)
			if component != "" {
				label = "//" + component + ":all"
			} else {
				gaps = append(gaps,
					"path outside registered topology: "+item)
			}
			add(v1alpha1.PlanTarget{
				Label:      label,
				Path:       item,
				AffectedBy: profileAffectedBy(profile),
			})
		}
	}
	if forceRoot && profile != ProfileFullAudit {
		root := v1alpha1.PlanTarget{
			Label:      "//...",
			Path:       ".",
			AffectedBy: profileAffectedBy(profile),
		}
		if !seen[root.Label] {
			targets = append([]v1alpha1.PlanTarget{root}, targets...)
		}
	}
	if profile != ProfileChangedFast && profile != ProfileDiagnose {
		gaps = append(gaps,
			"reverse dependency graph unavailable: no reverse-affected "+
				"targets claimed")
	}
	sort.Slice(targets, func(i, j int) bool {
		if targets[i].Label != targets[j].Label {
			return targets[i].Label < targets[j].Label
		}
		return targets[i].Path < targets[j].Path
	})
	return targets, dedupe(gaps)
}

// componentForItem returns the Bazel package directory of the topology
// component whose path is the nearest registered ancestor of the item,
// preferring the longest registered prefix.
func componentForItem(item string, idx *Index) string {
	if idx == nil || idx.Topology == nil {
		return ""
	}
	best := ""
	bestLen := -1
	for _, component := range idx.Topology.Components {
		prefix := component.Path
		if item == prefix || strings.HasPrefix(item, prefix+"/") {
			if len(prefix) > bestLen {
				best = componentPackage(component)
				bestLen = len(prefix)
			}
		}
	}
	return best
}

// componentPackage is the Bazel package directory for one topology
// component, derived from its build file when present.
func componentPackage(component catalogv1alpha1.TopologyComponent) string {
	if component.BuildPath != "" {
		return path.Dir(component.BuildPath)
	}
	return component.Path
}

// scopeExpansionReason names configuration whose change invalidates narrow
// scope. The planner expands such changes to the repository root.
func scopeExpansionReason(item string) string {
	base := path.Base(item)
	lower := strings.ToLower(base)
	switch {
	case base == "MODULE.bazel", base == "MODULE.bazel.lock",
		base == "WORKSPACE", base == "WORKSPACE.bazel",
		base == ".bazelrc", base == ".bazelversion":
		return "root-module/rc change " + item +
			": scope conservatively expanded"
	case strings.HasSuffix(lower, ".bzl"):
		return "build rule change " + item +
			": scope conservatively expanded"
	case strings.HasPrefix(item, "tools/toolchains/"),
		strings.HasPrefix(item, "third_party/"):
		return "toolchain/dependency change " + item +
			": scope conservatively expanded"
	}
	return ""
}

// isUnknownGeneratedPath reports whether a path looks like generated output
// rather than ordinary source.
func isUnknownGeneratedPath(item string) bool {
	lower := strings.ToLower(item)
	for _, marker := range []string{"bazel-", "/out/", "/generated/"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// planEffects computes the required and forbidden effects for one scope.
// Plans are advisory: remote write and destroy are forbidden unless
// explicitly requested, in which case the escalation note records the
// authority requirement.
func planEffects(
	profile ImpactProfile,
	requestRemote bool,
) ([]v1alpha1.Effect, []v1alpha1.Effect) {
	effects := []v1alpha1.Effect{v1alpha1.EffectSourceRead}
	if profile != ProfileDiagnose {
		effects = append(effects, v1alpha1.EffectCodeExecute)
	}
	if profile == ProfileFullAudit {
		effects = append(effects, v1alpha1.EffectHostWrite)
	}
	if requestRemote {
		effects = append(effects,
			v1alpha1.EffectRemoteWrite, v1alpha1.EffectRemoteDestroy)
	}
	sort.Slice(effects, func(i, j int) bool {
		return effects[i] < effects[j]
	})
	var forbidden []v1alpha1.Effect
	if !requestRemote {
		forbidden = []v1alpha1.Effect{
			v1alpha1.EffectRemoteWrite,
			v1alpha1.EffectRemoteDestroy,
		}
	}
	sort.Slice(forbidden, func(i, j int) bool {
		return forbidden[i] < forbidden[j]
	})
	return effects, forbidden
}

// planChecks emits the conservative minimum checks in repository order:
// git diff --check, affected test, affected build, then profile-specific
// checks (buildifier for BUILD changes, full-repo-check for full/audit).
func planChecks(
	profile ImpactProfile,
	targets []v1alpha1.PlanTarget,
	scope []string,
) []v1alpha1.PlanCheck {
	var checks []v1alpha1.PlanCheck
	add := func(identifier, checkScope, covers string) {
		checks = append(checks, v1alpha1.PlanCheck{
			Identifier: identifier,
			Scope:      checkScope,
			Covers:     covers,
		})
	}
	add("git/diff/check", "", "whitespace and patch hygiene")
	if profile == ProfileDiagnose {
		checks = append(checks[:0], v1alpha1.PlanCheck{
			Identifier: "diagnose/read-only",
			Scope:      "affected",
			Covers:     "minimal read-only diagnosis",
		})
		return checks
	}
	if profile == ProfileFullAudit {
		add("bazel_agent/test/all", "//...",
			"repository-wide test coverage")
		add("bazel_agent/build/all", "//...",
			"repository-wide build coverage")
		add("full-repo-check", "", "structured exhaustive repository check")
	} else {
		for _, target := range targets {
			checkScope := identifierLabel(target.Label)
			add("bazel_agent/test/"+checkScope, target.Label,
				"affected package tests at "+target.Label)
			add("bazel_agent/build/"+checkScope, target.Label,
				"affected package build at "+target.Label)
		}
	}
	if hasBuildChange(scope) {
		add("buildifier/affected", strings.Join(scope, ","),
			"BUILD file formatting")
	}
	return checks
}

func hasBuildChange(scope []string) bool {
	for _, item := range scope {
		base := path.Base(item)
		if base == "BUILD" || base == "BUILD.bazel" ||
			strings.HasSuffix(base, ".bzl") {
			return true
		}
	}
	return false
}

// identifierLabel maps a Bazel label to the lowercase-safe check identifier
// suffix.
func identifierLabel(label string) string {
	switch label {
	case "", "//...":
		return "all"
	}
	value := strings.TrimPrefix(label, "//")
	return strings.ReplaceAll(value, ":", "/")
}

// selectCapability picks the bounded validation capability. The planner
// prefers the registered bazel-agent provider and falls back to a generic
// artifact capability when the capability catalog is unavailable.
func selectCapability(idx *Index) (v1alpha1.Reference, string) {
	if idx != nil && idx.Capability != nil {
		for _, provider := range idx.Capability.Providers {
			if provider.ID == "bazel-agent" {
				return v1alpha1.Reference{
					Kind: v1alpha1.ReferenceProvider,
					ID:   "bazel-agent",
				}, "bazel-agent: hermetic Bazel validation provider"
			}
		}
	}
	return v1alpha1.Reference{
			Kind: v1alpha1.ReferenceArtifact,
			ID:   CapabilityValidation,
		}, "validation: generic capability (capability catalog " +
			"unavailable or no bazel-agent provider)"
}

type costProfile struct {
	class  v1alpha1.CostClass
	budget v1alpha1.Budget
}

func profileCost(profile ImpactProfile) costProfile {
	switch profile {
	case ProfileChangedFast, ProfileDiagnose:
		return costProfile{
			class: v1alpha1.CostClassLow,
			budget: v1alpha1.Budget{
				Calls:       8,
				Bytes:       1 << 20,
				DurationMS:  60000,
				Concurrency: 1,
			},
		}
	case ProfileWorkspace, ProfileFreshEvidence:
		return costProfile{
			class: v1alpha1.CostClassMedium,
			budget: v1alpha1.Budget{
				Calls:       24,
				Bytes:       8 << 20,
				DurationMS:  300000,
				Concurrency: 2,
			},
		}
	default:
		return costProfile{
			class: v1alpha1.CostClassHigh,
			budget: v1alpha1.Budget{
				Calls:       64,
				Bytes:       64 << 20,
				DurationMS:  1800000,
				Concurrency: 4,
			},
		}
	}
}

func profileCacheability(profile ImpactProfile) v1alpha1.Cacheability {
	switch profile {
	case ProfileChangedFast, ProfileWorkspace, ProfileFreshEvidence:
		return v1alpha1.Cacheable
	case ProfileFullAudit:
		return v1alpha1.NotCacheable
	}
	return v1alpha1.CacheabilityUnknown
}

func profileOrdering(profile ImpactProfile) []string {
	if profile == ProfileDiagnose {
		return []string{"diagnose"}
	}
	return []string{"git-diff-check", "test", "build"}
}

func profileAffectedBy(profile ImpactProfile) string {
	if profile == ProfileChangedFast {
		return "source change"
	}
	return "scope match"
}

func dedupe(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}
