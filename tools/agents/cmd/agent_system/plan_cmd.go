package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	"git.alwaldend.com/alwaldend/src/tools/agents/plan"
)

type planOptions struct {
	workspaceRoot string
	profile       string
	intent        string
	paths         multiFlag
	labels        multiFlag
	remote        bool
}

// multiFlag collects repeated string flags.
type multiFlag []string

func (m *multiFlag) String() string {
	return fmt.Sprint([]string(*m))
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func parsePlanFlags(args []string) (planOptions, error) {
	var opts planOptions
	flags := flag.NewFlagSet("agent_system plan", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required; defaults to BUILD_WORKSPACE_DIRECTORY)")
	flags.StringVar(&opts.profile, "profile", "changed/fast",
		"impact profile: changed/fast, workspace, fresh/evidence, "+
			"full/audit, diagnose")
	flags.StringVar(&opts.intent, "intent", "",
		"intent identity (required)")
	flags.Var(&opts.paths, "path",
		"affected path relative to the workspace root (repeatable)")
	flags.Var(&opts.labels, "label",
		"affected Bazel label (repeatable)")
	flags.BoolVar(&opts.remote, "request-remote", false,
		"explicitly request remote.write/remote.destroy effects "+
			"(still advisory; escalation notes the authority requirement)")
	if err := flags.Parse(args); err != nil {
		return planOptions{}, err
	}
	if opts.workspaceRoot == "" {
		opts.workspaceRoot = os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	}
	if opts.workspaceRoot == "" {
		return planOptions{}, fmt.Errorf("--workspace-root is required")
	}
	if opts.intent == "" {
		return planOptions{}, fmt.Errorf("--intent is required")
	}
	if len(opts.paths) == 0 && len(opts.labels) == 0 {
		return planOptions{}, fmt.Errorf("at least one --path or --label is required")
	}
	return opts, nil
}

// runPlan loads the planner index from the workspace catalogs and emits the
// advisory impact plan as canonical JSON. It never mutates source.
func runPlan(args []string, stdout io.Writer) error {
	opts, err := parsePlanFlags(args)
	if err != nil {
		return err
	}
	root, err := filepath.Abs(opts.workspaceRoot)
	if err != nil {
		return err
	}
	builder := &capsuleBuilder{root: root}
	idx, err := loadPlanIndex(builder)
	if err != nil {
		return err
	}
	profile, err := resolveProfileName(opts.profile)
	if err != nil {
		return err
	}
	plan, err := plan.Plan(v1alpha1.Reference{
		Kind: v1alpha1.ReferenceTask,
		ID:   opts.intent,
	}, plan.PlanOptions{
		Profile:       profile,
		Intent:        opts.intent,
		AffectedPaths: []string(opts.paths),
		Labels:        []string(opts.labels),
		RequestRemote: opts.remote,
	}, idx)
	if err != nil {
		return err
	}
	content, err := v1alpha1.CanonicalImpactPlanJSON(plan)
	if err != nil {
		return err
	}
	_, err = stdout.Write(content)
	return err
}

// loadPlanIndex reads the bounded planner projections. A missing catalog is
// a structured gap, not a failure: the planner conservatively falls back to
// generic validation capability and records coverage gaps.
func loadPlanIndex(builder *capsuleBuilder) (*plan.Index, error) {
	idx := &plan.Index{}
	if catalog, err := builder.readTopology(); err == nil {
		idx.Topology = catalog
	}
	if catalog, err := builder.readCapability(); err == nil {
		idx.Capability = catalog
	}
	if catalog, err := builder.readWorkspaceCheck(); err == nil {
		idx.WorkspaceCheck = catalog
	}
	return idx, nil
}

func resolveProfileName(value string) (plan.ImpactProfile, error) {
	switch value {
	case string(plan.ProfileChangedFast):
		return plan.ProfileChangedFast, nil
	case string(plan.ProfileWorkspace):
		return plan.ProfileWorkspace, nil
	case string(plan.ProfileFreshEvidence):
		return plan.ProfileFreshEvidence, nil
	case string(plan.ProfileFullAudit):
		return plan.ProfileFullAudit, nil
	case string(plan.ProfileDiagnose):
		return plan.ProfileDiagnose, nil
	}
	return "", fmt.Errorf("unknown profile %q", value)
}
