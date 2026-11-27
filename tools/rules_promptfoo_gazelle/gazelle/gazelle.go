package gazelle

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// maxEvalEntries bounds directory entries examined across one package.
	maxEvalEntries = 65536
	// maxEvalDirectoryDepth bounds the number of nested directories below evals.
	maxEvalDirectoryDepth = 64
	// maxValidationInputs bounds the total config and data labels emitted per package.
	maxValidationInputs = 1 << 20

	// promptfooConfigPrefix starts every recognized evaluation configuration basename.
	promptfooConfigPrefix = "promptfooconfig"
	// promptfooTestKind identifies live tests whose attributes belong to the author.
	promptfooTestKind = "promptfoo_test"
	// promptfooValidateKind identifies the offline checks managed by this generator.
	promptfooValidateKind = "promptfoo_validate_test"
)

// configExtensions lists configuration formats recognized by the validation rule.
var configExtensions = map[string]bool{
	".json": true,
	".yaml": true,
	".yml":  true,
}

// promptfooLanguage implements the Gazelle extension and its module-aware load hook.
type promptfooLanguage struct {
	// BaseLang supplies the default implementations of the remaining Gazelle hooks.
	language.BaseLang
	// loadLabel overrides the published rule label for an in-tree checkout.
	loadLabel string
}

// promptfooConfig records the names used to generate a validation target for one file.
type promptfooConfig struct {
	// path is the package-relative config label.
	path string
	// sanitizedVariant is the normalized suffix for the target name.
	sanitizedVariant string
	// variant retains the original spelling for exact control-variant matching.
	variant string
}

var (
	_ language.Language            = (*promptfooLanguage)(nil)
	_ language.ModuleAwareLanguage = (*promptfooLanguage)(nil)
)

// NewLanguage returns the Promptfoo Gazelle extension.
func NewLanguage() language.Language {
	return &promptfooLanguage{}
}

// Name identifies this extension to Gazelle's language filter.
func (*promptfooLanguage) Name() string {
	return "promptfoo"
}

// Kinds allows Gazelle to reconcile only the validation inputs it owns.
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

// Loads rejects legacy load resolution, which cannot account for Bzlmod repository mappings.
func (*promptfooLanguage) Loads() []rule.LoadInfo {
	panic("ApparentLoads should be called instead")
}

// RegisterFlags allows a consumer workspace to select an in-tree rule label.
func (lang *promptfooLanguage) RegisterFlags(flags *flag.FlagSet, _ string, _ *config.Config) {
	flags.StringVar(&lang.loadLabel, "promptfoo_load_label", "", "Bazel label for the promptfoo rule definitions")
}

// ApparentLoads resolves the rule definitions in the consumer workspace.
func (lang *promptfooLanguage) ApparentLoads(
	moduleToApparentName func(string) string,
) []rule.LoadInfo {
	rulesPromptfoo := moduleToApparentName("rules_promptfoo")
	if rulesPromptfoo == "" {
		rulesPromptfoo = "rules_promptfoo"
	}
	loadLabel := lang.loadLabel
	if loadLabel == "" {
		loadLabel = fmt.Sprintf("@%s//promptfoo:defs.bzl", rulesPromptfoo)
	}
	return []rule.LoadInfo{{
		Name: loadLabel,
		Symbols: []string{
			promptfooTestKind,
			promptfooValidateKind,
		},
	}}
}

// GenerateRules emits offline validation rules only after discovering the complete eval file set.
func (*promptfooLanguage) GenerateRules(
	args language.GenerateArgs,
) language.GenerateResult {
	if args.Dir == "" {
		return language.GenerateResult{}
	}

	configs, data, err := discoverEvalFiles(args)
	if err != nil {
		log.Printf("skip Promptfoo generation in %q: %v", args.Dir, err)
		return language.GenerateResult{}
	}
	if len(configs) == 0 {
		return language.GenerateResult{}
	}

	if len(configs) > maxValidationInputs/(len(data)+1) {
		log.Printf("skip Promptfoo generation in %q: more than %d validation inputs", args.Dir, maxValidationInputs)
		return language.GenerateResult{}
	}
	names := targetNames(configs)
	withSkill := args.Rel != "" && slices.Contains(args.RegularFiles, "SKILL.md")
	generated := make([]*rule.Rule, 0, len(configs))
	imports := make([]interface{}, 0, len(configs))
	for index, config := range configs {
		validation := rule.NewRule(promptfooValidateKind, names[index])
		validation.SetAttr("config", config.path)
		if len(data) > 0 {
			validation.SetAttr("data", data)
		}
		if withSkill && config.variant != "no_skill" {
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

// discoverEvalFiles collects package-local configuration and data paths in stable order.
func discoverEvalFiles(
	args language.GenerateArgs,
) ([]promptfooConfig, []string, error) {
	if args.Config == nil || args.Config.RepoRoot == "" || args.Dir == "" {
		return nil, nil, fmt.Errorf("repository root and package directory are required")
	}
	repoRoot, err := resolvedPath(args.Config.RepoRoot)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve repository root %q: %w", args.Config.RepoRoot, err)
	}
	packageDir, err := resolvedPath(args.Dir)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve package %q: %w", args.Dir, err)
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
		return nil, nil, fmt.Errorf("inspect evals directory %q: %w", evalsDir, err)
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

	buildFileNames := safeBuildFileNames(args.Config.ValidBuildFileNames)
	var configs []promptfooConfig
	var data []string
	remaining := maxEvalEntries
	var scan func(string, int) error
	scan = func(relative string, depth int) error {
		dir := filepath.Join(packageDir, relative)
		isPackage, err := hasBuildFile(dir, buildFileNames)
		if err != nil {
			return fmt.Errorf("check package boundary %q: %w", dir, err)
		}
		if isPackage {
			return nil
		}
		if depth > maxEvalDirectoryDepth {
			return fmt.Errorf("evaluation directory %q exceeds depth limit %d", dir, maxEvalDirectoryDepth)
		}
		entries, err := readEvalDirectory(dir, remaining)
		if err != nil {
			return fmt.Errorf("discover evaluation inputs: %w", err)
		}
		remaining -= len(entries)
		for _, entry := range entries {
			name := filepath.Join(relative, entry.Name())
			if entry.IsDir() {
				if err := scan(name, depth+1); err != nil {
					return fmt.Errorf("scan subdirectory %q: %w", entry.Name(), err)
				}
			} else if entry.Type().IsRegular() {
				variant, ok := configVariant(entry.Name())
				if depth == 0 && ok {
					configs = append(configs, promptfooConfig{
						path:             filepath.ToSlash(name),
						sanitizedVariant: sanitizeVariant(variant),
						variant:          variant,
					})
				} else {
					data = append(data, filepath.ToSlash(name))
				}
			}
		}
		return nil
	}
	if err := scan("evals", 0); err != nil {
		return nil, nil, fmt.Errorf("scan evaluation tree %q: %w", evalsDir, err)
	}

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].path < configs[j].path
	})
	sort.Strings(data)
	return configs, data, nil
}

// readEvalDirectory reads at most one entry beyond the remaining package budget to detect overflow.
func readEvalDirectory(dir string, remaining int) ([]os.DirEntry, error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("open evaluation directory %q: %w", dir, err)
	}
	var entries []os.DirEntry
	var readErr error
	for len(entries) <= remaining {
		batch, err := file.ReadDir(min(256, remaining+1-len(entries)))
		entries = append(entries, batch...)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				readErr = err
			}
			break
		}
	}
	if err := errors.Join(readErr, file.Close()); err != nil {
		return nil, fmt.Errorf("read evaluation directory %q: %w", dir, err)
	}
	if len(entries) > remaining {
		return nil, fmt.Errorf("evaluation tree exceeds limit of %d entries at %q", maxEvalEntries, dir)
	}
	return entries, nil
}

// resolvedPath resolves a directory through symlinks to an absolute path.
func resolvedPath(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("make %q absolute: %w", path, err)
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", fmt.Errorf("resolve symlinks in %q: %w", absolute, err)
	}
	return resolved, nil
}

// pathWithin reports whether path names root or one of its descendants.
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

// configVariant extracts the optional variant of a conventional Promptfoo filename.
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

// sanitizeVariant maps a variant to lowercase label characters separated by single underscores.
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

// targetNames uses a path-derived suffix when multiple configurations share a sanitized name.
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

// safeBuildFileNames combines conventional and configured package markers without duplicates.
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

// hasBuildFile checks the candidate package markers and reports failed probes.
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
		return false, fmt.Errorf("inspect marker %q: %w", filepath.Join(dir, name), err)
	}
	return false, nil
}
