package gazelle

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/merger"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// languageName identifies the extension in Gazelle's language registry.
	languageName = "enforce_transitive_visibility"
	// attributeName is the package attribute maintained by this extension.
	attributeName = "transitive_visibility"
	// configFlag selects the workspace policy configuration.
	configFlag = "transitive_visibility_config"
	// defaultConfig is relative to the selected workspace, never its parent.
	defaultConfig = "tools/transitive_visibility/transitive_visibility.toml"
)

// policyConfig decodes the public TOML configuration schema.
type policyConfig struct {
	// Paths maps canonical directory paths to absolute package-group labels.
	Paths map[string]string `toml:"paths"`
}

// visibilityLanguage generates package attributes from validated path mappings.
type visibilityLanguage struct {
	// BaseLang supplies unused language interfaces.
	language.BaseLang
	// filename is the selected configuration path.
	filename string
	// policy contains immutable mappings shared by all packages.
	policy policyConfig
}

// NewLanguage constructs the standalone Gazelle extension.
func NewLanguage() language.Language { return &visibilityLanguage{} }

// Name identifies this extension to Gazelle.
func (*visibilityLanguage) Name() string { return languageName }

// RegisterFlags exposes an optional workspace-relative configuration filename.
func (l *visibilityLanguage) RegisterFlags(fs *flag.FlagSet, _ string, _ *config.Config) {
	fs.StringVar(&l.filename, configFlag, defaultConfig, "TOML path-to-package-group configuration, relative to the repository root")
}

// CheckFlags validates all mappings before Gazelle walks or writes BUILD files.
func (l *visibilityLanguage) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	filename := l.filename
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(c.RepoRoot, filename)
	}
	explicit := false
	fs.Visit(func(f *flag.Flag) { explicit = explicit || f.Name == configFlag })
	content, err := os.ReadFile(filename)
	if os.IsNotExist(err) && !explicit {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read transitive_visibility configuration %q: %w", filename, err)
	}
	metadata, err := toml.Decode(string(content), &l.policy)
	if err != nil {
		return fmt.Errorf("decode transitive_visibility configuration %q: %w", filename, err)
	}
	if unknown := metadata.Undecoded(); len(unknown) != 0 {
		return fmt.Errorf("transitive_visibility configuration %q: unknown fields %v", filename, unknown)
	}
	if len(l.policy.Paths) == 0 {
		return fmt.Errorf("transitive_visibility configuration %q: paths must not be empty", filename)
	}
	paths := make([]string, 0, len(l.policy.Paths))
	for directory := range l.policy.Paths {
		paths = append(paths, directory)
	}
	sort.Strings(paths)
	for _, directory := range paths {
		if directory == "" || path.IsAbs(directory) || path.Clean(directory) != directory || directory == ".." || strings.HasPrefix(directory, "../") || strings.ContainsAny(directory, "\\:*") {
			return fmt.Errorf("transitive_visibility configuration %q: invalid workspace-relative path %q", filename, directory)
		}
		value := l.policy.Paths[directory]
		parsed, err := label.Parse(value)
		if err != nil {
			return fmt.Errorf("transitive_visibility configuration %q, path %q: parse label %q: %w", filename, directory, value, err)
		}
		if parsed.Relative || value == "" || strings.ContainsAny(value, " \t\r\n") {
			return fmt.Errorf("transitive_visibility configuration %q, path %q: expected one absolute package-group label, got %q", filename, directory, value)
		}
	}
	return nil
}

// KnownDirectives recognizes the removed directive only to report a migration error.
func (*visibilityLanguage) KnownDirectives() []string { return []string{attributeName} }

// Configure rejects obsolete overrides instead of silently ignoring policy.
func (*visibilityLanguage) Configure(_ *config.Config, _ string, f *rule.File) {
	if f == nil {
		return
	}
	for _, directive := range f.Directives {
		if directive.Key == attributeName {
			log.Fatalf("%s: gazelle:transitive_visibility is obsolete; move the mapping to the transitive_visibility TOML configuration", f.Path)
		}
	}
}

// Kinds owns only the transitive visibility attribute of the package declaration.
func (*visibilityLanguage) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{"package": {MatchAny: true, MergeableAttrs: map[string]bool{attributeName: true}}}
}

// policyFor selects the longest directory-prefix match without crossing siblings.
func (l *visibilityLanguage) policyFor(rel string) string {
	if rel == "" {
		rel = "."
	}
	for {
		if policy := l.policy.Paths[rel]; policy != "" {
			return policy
		}
		if rel == "." {
			return ""
		}
		rel = path.Dir(rel)
	}
}

// Fix moves managed package declarations after loads, preserving their syntax trees.
func (l *visibilityLanguage) Fix(_ *config.Config, f *rule.File) {
	if l.policyFor(f.Pkg) == "" {
		return
	}
	f.Sync()
	position := packagePosition(f)
	for _, declaration := range f.Rules {
		if declaration.Kind() != "package" || declaration.Index() <= position {
			continue
		}
		index := declaration.Index()
		statement := f.File.Stmt[index]
		copy(f.File.Stmt[position+1:index+1], f.File.Stmt[position:index])
		f.File.Stmt[position] = statement
		content := f.Content
		*f = *rule.ScanAST(f.Pkg, f.File)
		f.Content = content
		return
	}
}

// GenerateRules adds policy only to existing packages or packages another language creates.
func (l *visibilityLanguage) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	policy := l.policyFor(args.Rel)
	if policy == "" || (args.File == nil && len(args.OtherGen) == 0) {
		return language.GenerateResult{}
	}
	declaration := rule.NewRule("package", "")
	declaration.SetAttr(attributeName, policy)
	if args.File != nil {
		declaration.SetPrivateAttr(merger.UnstableInsertIndexKey, packagePosition(args.File))
	}
	return language.GenerateResult{Gen: []*rule.Rule{declaration}, Imports: []interface{}{nil}}
}

// packagePosition returns the first statement position after all load statements.
func packagePosition(f *rule.File) int {
	position := 0
	for _, load := range f.Loads {
		if load.Index() >= position {
			position = load.Index() + 1
		}
	}
	return position
}

// Ensure the extension implements Gazelle's language contract.
var _ language.Language = (*visibilityLanguage)(nil)
