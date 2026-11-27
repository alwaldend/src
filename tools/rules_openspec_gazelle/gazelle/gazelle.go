package gazelle

import (
	"flag"
	"fmt"
	"path"
	"slices"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// workspaceKind declares the runfiles of an OpenSpec workspace.
	workspaceKind = "openspec_workspace"
	// workspaceName is the default source target name.
	workspaceName = "source"
	// validationKind declares strict checks for the current and archived artifacts.
	validationKind = "openspec_validate_test"
	// validationName is the default workspace validation target name.
	validationName = "validate_test"
)

// openspecLanguage generates declarations using package-local Gazelle inputs.
type openspecLanguage struct {
	// BaseLang supplies the default implementations of the remaining Gazelle hooks.
	language.BaseLang
	// loadLabel overrides the rule module load for an in-tree build.
	loadLabel string
}

var (
	_ language.Language            = (*openspecLanguage)(nil)
	_ language.ModuleAwareLanguage = (*openspecLanguage)(nil)
)

// NewLanguage returns the OpenSpec Gazelle extension.
func NewLanguage() language.Language {
	return &openspecLanguage{}
}

// Name identifies the extension for Gazelle's language filter.
func (*openspecLanguage) Name() string {
	return "openspec"
}

// Kinds registers the generated declarations while preserving manual attributes.
func (*openspecLanguage) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{
		workspaceKind:  {},
		validationKind: {},
	}
}

// RegisterFlags exposes the rule load location used in the monorepo.
func (lang *openspecLanguage) RegisterFlags(flags *flag.FlagSet, _ string, _ *config.Config) {
	flags.StringVar(&lang.loadLabel, "openspec_load_label", "", "Bazel label for the OpenSpec rule definitions")
}

// Loads rejects initialization without module-aware load resolution.
func (*openspecLanguage) Loads() []rule.LoadInfo {
	panic("ApparentLoads should be called instead")
}

// ApparentLoads resolves the rules module's name as seen by the consumer.
func (lang *openspecLanguage) ApparentLoads(moduleToApparentName func(string) string) []rule.LoadInfo {
	loadLabel := lang.loadLabel
	if loadLabel == "" {
		repository := moduleToApparentName("rules_openspec")
		if repository == "" {
			repository = "rules_openspec"
		}
		loadLabel = fmt.Sprintf("@%s//rules:defs.bzl", repository)
	}
	return []rule.LoadInfo{{
		Name:    loadLabel,
		Symbols: []string{workspaceKind, validationKind},
	}}
}

// GenerateRules creates the two declarations for a configured OpenSpec workspace.
func (*openspecLanguage) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	if path.Base(args.Rel) != "openspec" || !slices.Contains(args.RegularFiles, "config.yaml") {
		return language.GenerateResult{}
	}
	source := rule.NewRule(workspaceKind, workspaceName)
	source.SetAttr("config", "config.yaml")
	if args.File != nil {
		for _, existing := range args.File.Rules {
			switch existing.Name() {
			case workspaceName:
				if existing.Kind() != workspaceKind {
					return language.GenerateResult{}
				}
				// Existing source declarations may already include the config.
				source.DelAttr("config")
			case validationName:
				if existing.Kind() != validationKind {
					return language.GenerateResult{}
				}
			}
		}
	}

	source.SetAttr("srcs", rule.GlobValue{
		Patterns:   []string{"changes/**", "specs/**/*.md"},
		AllowEmpty: true,
	})
	validation := rule.NewRule(validationKind, validationName)
	validation.SetAttr("source", ":"+workspaceName)
	return language.GenerateResult{
		Gen:     []*rule.Rule{source, validation},
		Imports: []interface{}{nil, nil},
	}
}
