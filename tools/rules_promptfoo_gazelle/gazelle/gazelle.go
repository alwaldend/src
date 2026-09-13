// Package gazelle generates offline Promptfoo validation targets.
package gazelle

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	promptfooConfigPrefix = "promptfooconfig"
	promptfooTestKind     = "promptfoo_test"
	promptfooValidateKind = "promptfoo_validate_test"
)

var configExtensions = map[string]bool{
	".json": true,
	".yaml": true,
	".yml":  true,
}

type promptfooLanguage struct {
	language.BaseLang
}

type promptfooConfig struct {
	path             string
	sanitizedVariant string
	variant          string
}

var (
	_ language.Language            = (*promptfooLanguage)(nil)
	_ language.ModuleAwareLanguage = (*promptfooLanguage)(nil)
)

// NewLanguage returns the Promptfoo Gazelle extension.
func NewLanguage() language.Language {
	return &promptfooLanguage{}
}

func (*promptfooLanguage) Name() string {
	return "promptfoo"
}

func (*promptfooLanguage) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{
		promptfooTestKind: {
			MergeableAttrs: map[string]bool{},
		},
		promptfooValidateKind: {
			NonEmptyAttrs: map[string]bool{
				"config": true,
			},
			MergeableAttrs: map[string]bool{
				"config": true,
				"data":   true,
				"skills": true,
			},
		},
	}
}

func (*promptfooLanguage) Loads() []rule.LoadInfo {
	panic("ApparentLoads should be called instead")
}

func (*promptfooLanguage) ApparentLoads(
	moduleToApparentName func(string) string,
) []rule.LoadInfo {
	rulesPromptfoo := moduleToApparentName("rules_promptfoo")
	if rulesPromptfoo == "" {
		rulesPromptfoo = "rules_promptfoo"
	}
	return []rule.LoadInfo{{
		Name: fmt.Sprintf(
			"@%s//promptfoo:defs.bzl",
			rulesPromptfoo,
		),
		Symbols: []string{
			promptfooTestKind,
			promptfooValidateKind,
		},
	}}
}

func (*promptfooLanguage) GenerateRules(
	args language.GenerateArgs,
) language.GenerateResult {
	if args.Dir == "" {
		return language.GenerateResult{}
	}

	configs, data, err := discoverEvalFiles(args)
	if err != nil {
		return language.GenerateResult{}
	}
	names := targetNames(configs)
	if len(configs) == 0 {
		return language.GenerateResult{}
	}

	generated := make([]*rule.Rule, 0, len(configs))
	imports := make([]interface{}, 0, len(configs))
	for index, config := range configs {
		validation := rule.NewRule(promptfooValidateKind, names[index])
		validation.SetAttr("config", config.path)
		if len(data) > 0 {
			validation.SetAttr("data", data)
		}
		if args.Rel != "" && contains(args.RegularFiles, "SKILL.md") &&
			config.variant != "no_skill" {
			validation.SetAttr("skills", []string{":skill"})
		}
		generated = append(generated, validation)
		imports = append(imports, nil)
	}

	return language.GenerateResult{
		Gen:     generated,
		Imports: imports,
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

func discoverEvalFiles(
	args language.GenerateArgs,
) ([]promptfooConfig, []string, error) {
	if args.Config == nil || args.Config.RepoRoot == "" || args.Dir == "" {
		return nil, nil, fmt.Errorf("repository root and package directory are required")
	}
	repoRoot, err := resolvedPath(args.Config.RepoRoot)
	if err != nil {
		return nil, nil, err
	}
	packageDir, err := resolvedPath(args.Dir)
	if err != nil {
		return nil, nil, err
	}
	if !pathWithin(repoRoot, packageDir) {
		return nil, nil, fmt.Errorf(
			"package directory %q resolves outside repository root %q",
			args.Dir,
			args.Config.RepoRoot,
		)
	}

	evalsDir := filepath.Join(packageDir, "evals")
	info, err := os.Lstat(evalsDir)
	if os.IsNotExist(err) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, nil, fmt.Errorf(
			"evals directory %q is a symbolic link",
			evalsDir,
		)
	}
	if !info.IsDir() {
		return nil, nil, nil
	}

	var configuredBuildFileNames []string
	if args.Config != nil && len(args.Config.ValidBuildFileNames) > 0 {
		configuredBuildFileNames = args.Config.ValidBuildFileNames
	}
	buildFileNames := safeBuildFileNames(configuredBuildFileNames)
	evalsIsPackage, err := hasBuildFile(evalsDir, buildFileNames)
	if err != nil {
		return nil, nil, err
	}
	if evalsIsPackage {
		return nil, nil, nil
	}

	var configs []promptfooConfig
	var data []string
	err = filepath.WalkDir(evalsDir, func(
		path string,
		entry fs.DirEntry,
		err error,
	) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path != evalsDir {
				isPackage, err := hasBuildFile(path, buildFileNames)
				if err != nil {
					return err
				}
				if isPackage {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(packageDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		variant, ok := configVariant(entry.Name())
		if ok && filepath.Dir(rel) == "evals" {
			configs = append(configs, promptfooConfig{
				path:             rel,
				sanitizedVariant: sanitizeVariant(variant),
				variant:          variant,
			})
			return nil
		}
		data = append(data, rel)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].path < configs[j].path
	})
	sort.Strings(data)
	return configs, data, nil
}

func resolvedPath(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(absolute)
}

func pathWithin(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	if err != nil || filepath.IsAbs(relative) {
		return false
	}
	return relative != ".." && !strings.HasPrefix(
		relative,
		".."+string(filepath.Separator),
	)
}

func configVariant(name string) (string, bool) {
	extension := filepath.Ext(name)
	if !configExtensions[extension] {
		return "", false
	}
	stem := strings.TrimSuffix(name, extension)
	if stem == promptfooConfigPrefix {
		return "", true
	}
	prefix := promptfooConfigPrefix + "."
	if !strings.HasPrefix(stem, prefix) {
		return "", false
	}
	variant := strings.TrimPrefix(stem, prefix)
	if variant == "" {
		return "", false
	}
	return variant, true
}

func sanitizeVariant(variant string) string {
	var result strings.Builder
	lastUnderscore := false
	for _, character := range strings.ToLower(variant) {
		if character >= 'a' && character <= 'z' ||
			character >= '0' && character <= '9' {
			result.WriteRune(character)
			lastUnderscore = false
			continue
		}
		if result.Len() > 0 && !lastUnderscore {
			result.WriteByte('_')
			lastUnderscore = true
		}
	}
	return strings.Trim(result.String(), "_")
}

func targetNames(configs []promptfooConfig) []string {
	bases := make([]string, len(configs))
	counts := make(map[string]int, len(configs))
	for index, config := range configs {
		variant := config.sanitizedVariant
		if variant == "" && config.variant != "" {
			variant = "variant"
		}
		base := "eval"
		if variant != "" {
			base += "_" + variant
		}
		bases[index] = base
		counts[base]++
	}

	names := make([]string, len(configs))
	for index, base := range bases {
		if counts[base] > 1 {
			digest := sha256.Sum256([]byte(configs[index].path))
			base += fmt.Sprintf("__%x", digest[:4])
		}
		names[index] = base + "_config_test"
	}
	return names
}

func safeBuildFileNames(configured []string) []string {
	names := make([]string, 0, len(configured)+2)
	seen := make(map[string]bool, len(configured)+2)
	for _, name := range append(
		[]string{"BUILD.bazel", "BUILD"},
		configured...,
	) {
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	return names
}

func hasBuildFile(dir string, names []string) (bool, error) {
	for _, name := range names {
		info, err := os.Stat(filepath.Join(dir, name))
		if err == nil {
			if info.Mode().IsRegular() {
				return true, nil
			}
			continue
		}
		if os.IsNotExist(err) {
			continue
		}
		return false, err
	}
	return false, nil
}
