// Package gazelle generates Bazel skill library targets.
package gazelle

import (
	"fmt"

	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	skillKind = "skill_library"
	skillName = "skill"
)

type skillLanguage struct {
	language.BaseLang
}

var (
	_ language.Language            = (*skillLanguage)(nil)
	_ language.ModuleAwareLanguage = (*skillLanguage)(nil)
)

// NewLanguage returns the skill Gazelle extension.
func NewLanguage() language.Language {
	return &skillLanguage{}
}

func (*skillLanguage) Name() string {
	return "skill"
}

func (*skillLanguage) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{
		skillKind: {
			// Attributes are intentionally non-mergeable. Gazelle may populate
			// missing attributes, but manual values always win.
			MergeableAttrs: map[string]bool{},
		},
	}
}

func (*skillLanguage) Loads() []rule.LoadInfo {
	panic("ApparentLoads should be called instead")
}

func (*skillLanguage) ApparentLoads(
	moduleToApparentName func(string) string,
) []rule.LoadInfo {
	rulesSkill := moduleToApparentName("rules_skill")
	if rulesSkill == "" {
		rulesSkill = "rules_skill"
	}
	return []rule.LoadInfo{{
		Name:    fmt.Sprintf("@%s//skill:defs.bzl", rulesSkill),
		Symbols: []string{skillKind},
	}}
}

func (*skillLanguage) GenerateRules(
	args language.GenerateArgs,
) language.GenerateResult {
	if args.Rel == "" || !contains(args.RegularFiles, "SKILL.md") {
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

func contains(files []string, name string) bool {
	for _, file := range files {
		if file == name {
			return true
		}
	}
	return false
}
