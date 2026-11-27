package gazelle

import (
	"flag"
	"fmt"
	"slices"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// skillKind identifies a packaged skill rule.
	skillKind = "skill_library"
	// skillName is the conventional skill target name.
	skillName = "skill"
)

// skillLanguage implements the Gazelle extension and its module-aware load hook.
type skillLanguage struct {
	// BaseLang supplies the default implementations of the remaining Gazelle hooks.
	language.BaseLang
	// loadLabel overrides the published rule label for an in-tree checkout.
	loadLabel string
}

var (
	_ language.Language            = (*skillLanguage)(nil)
	_ language.ModuleAwareLanguage = (*skillLanguage)(nil)
)

// NewLanguage returns the skill Gazelle extension.
func NewLanguage() language.Language {
	return &skillLanguage{}
}

// Name identifies this extension to Gazelle's language filter.
func (*skillLanguage) Name() string {
	return "skill"
}

// Kinds keeps user-specified attributes when merging skill rules.
func (*skillLanguage) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{
		skillKind: {
			// Attributes are intentionally non-mergeable. Gazelle may populate
			// missing attributes, but manual values always win.
			MergeableAttrs: map[string]bool{},
		},
	}
}

// Loads rejects legacy load resolution, which cannot account for Bzlmod repository mappings.
func (*skillLanguage) Loads() []rule.LoadInfo {
	panic("ApparentLoads should be called instead")
}

// RegisterFlags allows a consumer workspace to select an in-tree rule label.
func (lang *skillLanguage) RegisterFlags(flags *flag.FlagSet, _ string, _ *config.Config) {
	flags.StringVar(&lang.loadLabel, "skill_load_label", "", "Bazel label for the skill rule definitions")
}

// ApparentLoads resolves the rule definitions in the consumer workspace.
func (lang *skillLanguage) ApparentLoads(
	moduleToApparentName func(string) string,
) []rule.LoadInfo {
	rulesSkill := moduleToApparentName("rules_skills")
	if rulesSkill == "" {
		rulesSkill = "rules_skills"
	}
	loadLabel := lang.loadLabel
	if loadLabel == "" {
		loadLabel = fmt.Sprintf("@%s//skill:defs.bzl", rulesSkill)
	}
	return []rule.LoadInfo{{
		Name:    loadLabel,
		Symbols: []string{skillKind},
	}}
}

// GenerateRules bundles skill sources in each named package containing SKILL.md.
func (*skillLanguage) GenerateRules(
	args language.GenerateArgs,
) language.GenerateResult {
	if args.Rel == "" || !slices.Contains(args.RegularFiles, "SKILL.md") {
		return language.GenerateResult{}
	}

	skill := rule.NewRule(skillKind, skillName)
	skill.SetAttr("srcs", rule.GlobValue{
		Patterns: []string{"**"},
		Excludes: skillSourceExcludes(args),
	})

	return language.GenerateResult{
		Gen:     []*rule.Rule{skill},
		Imports: []interface{}{nil},
	}
}

// skillSourceExcludes omits BUILD files and separately packaged evaluations from the skill bundle.
func skillSourceExcludes(args language.GenerateArgs) []string {
	buildFileNames := []string{"BUILD.bazel", "BUILD"}
	if args.Config != nil {
		buildFileNames = append(buildFileNames, args.Config.ValidBuildFileNames...)
	}

	excludes := make([]string, 0, len(buildFileNames)+1)
	seen := make(map[string]bool, len(buildFileNames))
	for _, name := range buildFileNames {
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		excludes = append(excludes, name)
	}
	return append(excludes, "evals/**")
}
