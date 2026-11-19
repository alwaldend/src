// Package gazelle generates documentation targets for Bazel packages.
package gazelle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	docsKind = "docs_filegroup"
	docsName = "docs"
)

type docsLanguage struct {
	language.BaseLang
}

var (
	_ language.Language            = (*docsLanguage)(nil)
	_ language.ModuleAwareLanguage = (*docsLanguage)(nil)
)

// NewLanguage returns the documentation Gazelle extension.
func NewLanguage() language.Language {
	return &docsLanguage{}
}

func (*docsLanguage) Name() string {
	return "docs"
}

func (*docsLanguage) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{
		docsKind: {
			// Attributes are intentionally non-mergeable. Gazelle may populate
			// missing attributes, but manual values always win.
			MergeableAttrs: map[string]bool{},
		},
	}
}

func (*docsLanguage) Loads() []rule.LoadInfo {
	panic("ApparentLoads should be called instead")
}

func (*docsLanguage) ApparentLoads(
	moduleToApparentName func(string) string,
) []rule.LoadInfo {
	rulesDocs := moduleToApparentName("rules_docs")
	if rulesDocs == "" {
		rulesDocs = "rules_docs"
	}
	return []rule.LoadInfo{{
		Name:    fmt.Sprintf("@%s//docs:defs.bzl", rulesDocs),
		Symbols: []string{docsKind},
	}}
}

func (*docsLanguage) GenerateRules(
	args language.GenerateArgs,
) language.GenerateResult {
	if args.File == nil || !contains(args.RegularFiles, "README.md") {
		return language.GenerateResult{}
	}

	docs := rule.NewRule(docsKind, docsName)
	docs.SetAttr("srcs", rule.GlobValue{Patterns: []string{"*.md"}})
	if visibility := parentVisibility(args); visibility != nil {
		docs.SetAttr("visibility", visibility)
	}

	return language.GenerateResult{
		Gen:     []*rule.Rule{docs},
		Imports: []interface{}{nil},
	}
}

func contains(files []string, name string) bool {
	for _, file := range files {
		if file == name {
			return true
		}
	}
	return false
}

func parentVisibility(args language.GenerateArgs) []string {
	if args.Config == nil || args.Config.RepoRoot == "" || args.Dir == "" {
		return nil
	}

	repoRoot := filepath.Clean(args.Config.RepoRoot)
	dir := filepath.Clean(args.Dir)
	if dir == repoRoot {
		return nil
	}
	parent := filepath.Dir(dir)
	buildFileNames := args.Config.ValidBuildFileNames
	if len(buildFileNames) == 0 {
		buildFileNames = []string{"BUILD.bazel", "BUILD"}
	}

	for {
		rel, err := filepath.Rel(repoRoot, parent)
		if err != nil || rel == ".." || strings.HasPrefix(
			rel,
			".."+string(filepath.Separator),
		) {
			return nil
		}
		if hasBuildFile(parent, buildFileNames) {
			if rel == "." {
				return []string{"//:__pkg__"}
			}
			return []string{
				"//" + filepath.ToSlash(rel) + ":__pkg__",
			}
		}
		if parent == repoRoot {
			return nil
		}
		next := filepath.Dir(parent)
		if next == parent {
			return nil
		}
		parent = next
	}
}

func hasBuildFile(dir string, names []string) bool {
	for _, name := range names {
		info, err := os.Stat(filepath.Join(dir, name))
		if err == nil && info.Mode().IsRegular() {
			return true
		}
	}
	return false
}
