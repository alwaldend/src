package gazelle

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// docsKind identifies the rule that packages Markdown files.
	docsKind = "docs_filegroup"
	// docsName is the conventional documentation target in each package.
	docsName = "docs"
)

// docsLanguage implements the Gazelle extension and its module-aware load hook.
type docsLanguage struct {
	// BaseLang supplies the default implementations of the remaining Gazelle hooks.
	language.BaseLang
	// loadLabel overrides the published rule label for an in-tree checkout.
	loadLabel string
}

var (
	_ language.Language            = (*docsLanguage)(nil)
	_ language.ModuleAwareLanguage = (*docsLanguage)(nil)
)

// NewLanguage returns the documentation Gazelle extension.
func NewLanguage() language.Language {
	return &docsLanguage{}
}

// Name identifies this extension to Gazelle's language filter.
func (*docsLanguage) Name() string {
	return "docs"
}

// Kinds preserves attributes on existing documentation rules.
func (*docsLanguage) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{
		docsKind: {
			// Attributes are intentionally non-mergeable. Gazelle may populate
			// missing attributes, but manual values always win.
			MergeableAttrs: map[string]bool{},
		},
	}
}

// Loads rejects legacy load resolution, which cannot account for Bzlmod repository mappings.
func (*docsLanguage) Loads() []rule.LoadInfo {
	panic("ApparentLoads should be called instead")
}

// RegisterFlags allows a consumer workspace to select an in-tree rule label.
func (lang *docsLanguage) RegisterFlags(flags *flag.FlagSet, _ string, _ *config.Config) {
	flags.StringVar(&lang.loadLabel, "docs_load_label", "", "Bazel label for the docs rule definitions")
}

// ApparentLoads resolves the rule definitions in the consumer workspace.
func (lang *docsLanguage) ApparentLoads(
	moduleToApparentName func(string) string,
) []rule.LoadInfo {
	rulesDocs := moduleToApparentName("rules_docs")
	if rulesDocs == "" {
		rulesDocs = "rules_docs"
	}
	loadLabel := lang.loadLabel
	if loadLabel == "" {
		loadLabel = fmt.Sprintf("@%s//docs:defs.bzl", rulesDocs)
	}
	return []rule.LoadInfo{{
		Name:    loadLabel,
		Symbols: []string{docsKind},
	}}
}

// GenerateRules adds a Markdown file group to an existing Bazel package with a README.
func (*docsLanguage) GenerateRules(
	args language.GenerateArgs,
) language.GenerateResult {
	if args.File == nil || !slices.Contains(args.RegularFiles, "README.md") {
		return language.GenerateResult{}
	}

	docs := rule.NewRule(docsKind, docsName)
	docs.SetAttr("srcs", rule.GlobValue{Patterns: []string{"*.md"}})
	visibility, err := parentVisibility(args)
	if err != nil {
		log.Printf("skip documentation generation in %q: %v", args.Dir, err)
		return language.GenerateResult{}
	}
	if visibility != nil {
		docs.SetAttr("visibility", visibility)
	}

	return language.GenerateResult{
		Gen:     []*rule.Rule{docs},
		Imports: []interface{}{nil},
	}
}

// parentVisibility locates the nearest enclosing package within the workspace for documentation visibility.
func parentVisibility(args language.GenerateArgs) (visibility []string, resultErr error) {
	if args.Config == nil || args.Config.RepoRoot == "" || args.Dir == "" {
		return nil, nil
	}
	relative, err := filepath.Rel(args.Config.RepoRoot, args.Dir)
	if err != nil {
		return nil, fmt.Errorf("locate package %q in workspace %q: %w", args.Dir, args.Config.RepoRoot, err)
	}
	if !filepath.IsLocal(relative) {
		return nil, fmt.Errorf("package %q is outside workspace %q", args.Dir, args.Config.RepoRoot)
	}
	if relative == "." {
		return nil, nil
	}

	root, err := os.OpenRoot(args.Config.RepoRoot)
	if err != nil {
		return nil, fmt.Errorf("open documentation workspace %q: %w", args.Config.RepoRoot, err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close documentation workspace: %w", err))
		}
	}()
	buildFileNames := append([]string{"BUILD.bazel", "BUILD"}, args.Config.ValidBuildFileNames...)
	for parent := filepath.Dir(relative); ; parent = filepath.Dir(parent) {
		for _, name := range buildFileNames {
			if name == "" {
				continue
			}
			boundary := filepath.Join(parent, name)
			info, err := root.Stat(boundary)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("inspect package marker %q: %w", boundary, err)
			}
			if info.Mode().IsRegular() {
				if parent == "." {
					return []string{"//:__pkg__"}, nil
				}
				return []string{"//" + filepath.ToSlash(parent) + ":__pkg__"}, nil
			}
		}
		if parent == "." {
			return nil, nil
		}
	}
}
