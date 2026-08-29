package gazelle

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/merger"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func TestGenerateRules(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "skills", "answer-question")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":                              "",
		"SKILL.md":                                 "---\nname: answer-question\n---\n",
		"evals/custom/BUILD.custom":                "",
		"evals/custom/context.txt":                 "hidden",
		"evals/README.md":                          "eval documentation",
		"evals/cases.yaml":                         "[]",
		"evals/fixtures/context.txt":               "context",
		"evals/legacy/BUILD":                       "",
		"evals/legacy/context.txt":                 "hidden",
		"evals/promptfooconfig.no_skill.yaml":      "tests: []",
		"evals/promptfooconfig.routing.yml":        "tests: []",
		"evals/promptfooconfig.yaml":               "tests: []",
		"evals/unrelated.json":                     "{}",
		"evals/vendor/BUILD.bazel":                 "",
		"evals/vendor/promptfooconfig.hidden.yaml": "tests: []",
	})

	args := generateArgs(repoRoot, packageDir)
	args.Config.ValidBuildFileNames = []string{"BUILD.custom"}
	args.RegularFiles = []string{"SKILL.md"}
	result := NewLanguage().GenerateRules(args)
	if len(result.Gen) != 3 || len(result.Imports) != 3 {
		t.Fatalf(
			"GenerateRules() generated %d rules and %d imports, want 3 each",
			len(result.Gen),
			len(result.Imports),
		)
	}

	wantData := []string{
		"evals/README.md",
		"evals/cases.yaml",
		"evals/fixtures/context.txt",
		"evals/unrelated.json",
	}
	want := []struct {
		name   string
		config string
		skills []string
	}{
		{
			name:   "eval_no_skill_config_test",
			config: "evals/promptfooconfig.no_skill.yaml",
		},
		{
			name:   "eval_routing_config_test",
			config: "evals/promptfooconfig.routing.yml",
			skills: []string{":skill"},
		},
		{
			name:   "eval_config_test",
			config: "evals/promptfooconfig.yaml",
			skills: []string{":skill"},
		},
	}

	for index, expectation := range want {
		generated := result.Gen[index]
		if generated.Kind() != promptfooValidateKind ||
			generated.Name() != expectation.name {
			t.Errorf(
				"generated rule %d = %s(%q), want %s(%q)",
				index,
				generated.Kind(),
				generated.Name(),
				promptfooValidateKind,
				expectation.name,
			)
		}
		if got := generated.AttrString("config"); got != expectation.config {
			t.Errorf("%s config = %q, want %q", expectation.name, got,
				expectation.config)
		}
		if got := generated.AttrStrings("data"); !reflect.DeepEqual(got, wantData) {
			t.Errorf("%s data = %q, want %q", expectation.name, got,
				wantData)
		}
		if got := generated.AttrStrings("skills"); !reflect.DeepEqual(
			got,
			expectation.skills,
		) {
			t.Errorf("%s skills = %q, want %q", expectation.name, got,
				expectation.skills)
		}
	}
	if len(result.Empty) != 0 {
		t.Errorf("Empty contains %d rules, want none", len(result.Empty))
	}
}

func TestGenerateRulesUsesStableCollisionSuffixes(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "pkg")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":                                  "",
		"evals/promptfooconfig.foo-bar.yaml":           "tests: []",
		"evals/promptfooconfig.foo_bar.yaml":           "tests: []",
		"evals/promptfooconfig.foo_bar__1ac77dfe.yaml": "tests: []",
	})

	result := NewLanguage().GenerateRules(generateArgs(repoRoot, packageDir))
	got := make([]string, 0, len(result.Gen))
	for _, generated := range result.Gen {
		got = append(got, generated.Name())
	}
	want := []string{
		"eval_foo_bar__1ac77dfe_config_test",
		"eval_foo_bar__ae45dea1_config_test",
		"eval_foo_bar_1ac77dfe_config_test",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("generated names = %q, want %q", got, want)
	}
}

func TestRemovedDuplicateConfigsPreserveExistingTargets(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "pkg")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":                "",
		"evals/promptfooconfig.yaml": "tests: []",
		"evals/promptfooconfig.yml":  "tests: []",
	})

	args := generateArgs(repoRoot, packageDir)
	initial := NewLanguage().GenerateRules(args)
	if len(initial.Gen) != 2 {
		t.Fatalf("GenerateRules() generated %d rules, want 2", len(initial.Gen))
	}
	var buildSource strings.Builder
	for _, generated := range initial.Gen {
		if !strings.HasPrefix(generated.Name(), "eval__") {
			t.Fatalf("duplicate default name %q lacks collision suffix",
				generated.Name())
		}
		buildSource.WriteString("promptfoo_validate_test(\n    name = \"")
		buildSource.WriteString(generated.Name())
		buildSource.WriteString("\",\n    config = \"")
		buildSource.WriteString(generated.AttrString("config"))
		buildSource.WriteString("\",\n)\n\n")
	}
	build, err := rule.LoadData(
		"BUILD.bazel",
		"pkg",
		[]byte(buildSource.String()),
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"evals/promptfooconfig.yaml",
		"evals/promptfooconfig.yml",
	} {
		if err := os.Remove(filepath.Join(packageDir, name)); err != nil {
			t.Fatal(err)
		}
	}

	args.File = build
	removed := NewLanguage().GenerateRules(args)
	if len(removed.Gen) != 0 || len(removed.Empty) != 0 {
		t.Fatalf(
			"after removal generated %d rules and %d empty rules, want 0 and 0",
			len(removed.Gen),
			len(removed.Empty),
		)
	}
	merger.MergeFile(
		build,
		removed.Empty,
		removed.Gen,
		merger.PreResolve,
		NewLanguage().Kinds(),
		nil,
	)
	for _, generated := range initial.Gen {
		findRule(t, build, promptfooValidateKind, generated.Name())
	}
}

func TestGenericPackageDoesNotInventSkillTarget(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "pkg")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":                "",
		"evals/promptfooconfig.yaml": "tests: []",
	})

	result := NewLanguage().GenerateRules(generateArgs(repoRoot, packageDir))
	if len(result.Gen) != 1 {
		t.Fatalf("GenerateRules() generated %d rules, want 1", len(result.Gen))
	}
	if got := result.Gen[0].AttrStrings("skills"); got != nil {
		t.Fatalf("generic validation skills = %q, want omitted", got)
	}
}

func TestRootPackageDoesNotInferUnsupportedSkillTarget(t *testing.T) {
	repoRoot := t.TempDir()
	writeFiles(t, repoRoot, map[string]string{
		"BUILD.bazel":                "",
		"SKILL.md":                   "---\nname: fixture\n---\n",
		"evals/promptfooconfig.yaml": "tests: []",
	})

	args := generateArgs(repoRoot, repoRoot)
	args.Rel = ""
	args.RegularFiles = []string{"SKILL.md"}
	result := NewLanguage().GenerateRules(args)
	if len(result.Gen) != 1 {
		t.Fatalf("GenerateRules() generated %d rules, want 1", len(result.Gen))
	}
	if got := result.Gen[0].AttrStrings("skills"); got != nil {
		t.Fatalf("root validation skills = %q, want omitted", got)
	}
}

func TestOnlyExactNoSkillVariantOmitsSkill(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "skill")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":                         "",
		"SKILL.md":                            "---\nname: skill\n---\n",
		"evals/promptfooconfig.no-skill.yaml": "tests: []",
	})

	args := generateArgs(repoRoot, packageDir)
	args.RegularFiles = []string{"SKILL.md"}
	result := NewLanguage().GenerateRules(args)
	if len(result.Gen) != 1 {
		t.Fatalf("GenerateRules() generated %d rules, want 1", len(result.Gen))
	}
	if got := result.Gen[0].AttrStrings("skills"); !reflect.DeepEqual(
		got,
		[]string{":skill"},
	) {
		t.Fatalf("non-control validation skills = %q, want [:skill]", got)
	}
}

func TestGenerateRulesRequiresDirectoryAndConfigs(t *testing.T) {
	tests := []struct {
		name       string
		build      bool
		files      map[string]string
		evalsBuild bool
		want       int
		withDir    bool
	}{
		{
			name:    "no BUILD file",
			files:   map[string]string{"evals/promptfooconfig.yaml": ""},
			want:    1,
			withDir: true,
		},
		{
			name:    "no evals directory",
			build:   true,
			withDir: true,
		},
		{
			name:    "no Promptfoo config",
			build:   true,
			files:   map[string]string{"evals/cases.yaml": "[]"},
			withDir: true,
		},
		{
			name:       "evals is a subpackage",
			build:      true,
			files:      map[string]string{"evals/promptfooconfig.yaml": ""},
			evalsBuild: true,
			withDir:    true,
		},
		{
			name:  "missing physical directory",
			build: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoRoot := t.TempDir()
			packageDir := filepath.Join(repoRoot, "pkg")
			files := make(map[string]string, len(tt.files)+2)
			for path, content := range tt.files {
				files[path] = content
			}
			if tt.build {
				files["BUILD.bazel"] = ""
			}
			if tt.evalsBuild {
				files["evals/BUILD.bazel"] = ""
			}
			if tt.withDir {
				writeFiles(t, packageDir, files)
			}

			args := generateArgs(repoRoot, packageDir)
			if !tt.build {
				args.File = nil
			}
			if !tt.withDir {
				args.Dir = ""
			}
			result := NewLanguage().GenerateRules(args)
			if len(result.Gen) != tt.want || len(result.Empty) != 0 {
				t.Fatalf(
					"GenerateRules() = %d generated, %d empty; want %d generated",
					len(result.Gen),
					len(result.Empty),
					tt.want,
				)
			}
		})
	}
}

func TestEvalScanErrorDoesNotMarkValidationTargetsObsolete(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "pkg")
	writeFiles(t, packageDir, map[string]string{"BUILD.bazel": ""})
	if err := os.Symlink("evals", filepath.Join(packageDir, "evals")); err != nil {
		t.Skipf("cannot create symlink loop: %v", err)
	}

	const existing = `promptfoo_validate_test(
    name = "eval_config_test",
    config = "evals/promptfooconfig.yaml",
)
`
	build, err := rule.LoadData("BUILD.bazel", "pkg", []byte(existing))
	if err != nil {
		t.Fatal(err)
	}
	lang := NewLanguage()
	args := generateArgs(repoRoot, packageDir)
	args.File = build
	result := lang.GenerateRules(args)
	if len(result.Gen) != 0 || len(result.Empty) != 0 {
		t.Fatalf(
			"scan failure generated %d rules and %d empty rules, want no-op",
			len(result.Gen),
			len(result.Empty),
		)
	}
	merger.MergeFile(
		build,
		result.Empty,
		result.Gen,
		merger.PreResolve,
		lang.Kinds(),
		nil,
	)
	findRule(t, build, promptfooValidateKind, "eval_config_test")
}

func TestEvalSymlinkDoesNotMarkValidationTargetsObsolete(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "pkg")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":                       "",
		"linked-evals/promptfooconfig.yaml": "tests: []",
	})
	evalsDir := filepath.Join(packageDir, "evals")
	if err := os.Symlink("linked-evals", evalsDir); err != nil {
		t.Skipf("cannot create evals symlink: %v", err)
	}

	const existing = `promptfoo_validate_test(
    name = "eval_config_test",
    config = "evals/promptfooconfig.yaml",
)
`
	build, err := rule.LoadData("BUILD.bazel", "pkg", []byte(existing))
	if err != nil {
		t.Fatal(err)
	}
	lang := NewLanguage()
	args := generateArgs(repoRoot, packageDir)
	args.File = build
	result := lang.GenerateRules(args)
	if len(result.Gen) != 0 || len(result.Empty) != 0 {
		t.Fatalf(
			"symlinked evals generated %d rules and %d empty rules, want no-op",
			len(result.Gen),
			len(result.Empty),
		)
	}
	merger.MergeFile(
		build,
		result.Empty,
		result.Gen,
		merger.PreResolve,
		lang.Kinds(),
		nil,
	)
	findRule(t, build, promptfooValidateKind, "eval_config_test")
}

func TestFollowedPackageSymlinkCannotEscapeRepository(t *testing.T) {
	repoRoot := t.TempDir()
	outsideRoot := t.TempDir()
	outsidePackage := filepath.Join(outsideRoot, "skill")
	writeFiles(t, outsidePackage, map[string]string{
		"evals/promptfooconfig.yaml": "tests: []",
	})
	packageDir := filepath.Join(repoRoot, "linked-skill")
	if err := os.Symlink(outsidePackage, packageDir); err != nil {
		t.Skipf("cannot create package symlink: %v", err)
	}

	const existing = `promptfoo_validate_test(
    name = "eval_config_test",
    config = "evals/promptfooconfig.yaml",
)
`
	build, err := rule.LoadData(
		"BUILD.bazel",
		"linked-skill",
		[]byte(existing),
	)
	if err != nil {
		t.Fatal(err)
	}
	args := generateArgs(repoRoot, packageDir)
	args.File = build
	args.Rel = "linked-skill"
	result := NewLanguage().GenerateRules(args)
	if len(result.Gen) != 0 || len(result.Empty) != 0 {
		t.Fatalf(
			"escaped package generated %d rules and %d empty rules, want no-op",
			len(result.Gen),
			len(result.Empty),
		)
	}
}

func TestBuildBoundaryScanErrorIsNonDestructive(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "pkg")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":                "",
		"evals/promptfooconfig.yaml": "tests: []",
		"evals/vendor/context.txt":   "hidden",
	})
	boundary := filepath.Join(packageDir, "evals", "vendor", "BUILD")
	if err := os.Symlink("BUILD", boundary); err != nil {
		t.Skipf("cannot create build-file symlink loop: %v", err)
	}

	const existing = `promptfoo_validate_test(
    name = "eval_config_test",
    config = "evals/promptfooconfig.yaml",
)
`
	build, err := rule.LoadData("BUILD.bazel", "pkg", []byte(existing))
	if err != nil {
		t.Fatal(err)
	}
	args := generateArgs(repoRoot, packageDir)
	args.File = build
	result := NewLanguage().GenerateRules(args)
	if len(result.Gen) != 0 || len(result.Empty) != 0 {
		t.Fatalf(
			"boundary scan failure generated %d rules and %d empty rules, want no-op",
			len(result.Gen),
			len(result.Empty),
		)
	}
}

func TestMergeReconcilesGeneratedAttrsAndPreservesManualAttrs(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "pkg")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":                              "",
		"SKILL.md":                                 "---\nname: pkg\n---\n",
		"evals/cases.yaml":                         "[]",
		"evals/context.txt":                        "context",
		"evals/promptfooconfig.kept.yaml":          "tests: []",
		"evals/promptfooconfig.no_skill.yaml":      "tests: []",
		"evals/promptfooconfig.yaml":               "tests: []",
		"evals/vendor/BUILD.bazel":                 "",
		"evals/vendor/promptfooconfig.hidden.yaml": "tests: []",
	})

	const existing = `load(
    "@rules_promptfoo//promptfoo:defs.bzl",
    "promptfoo_test",
    "promptfoo_validate_test",
)

promptfoo_validate_test(
    name = "eval_config_test",
    args = ["--manual-validation-option"],
    config = "evals/manual.yaml",
    data = [":manual_data"],
    env = {"MODE": "manual"},
    skills = [":manual_skill"],
    tags = ["manual-validation-tag"],
)

promptfoo_validate_test(
    name = "eval_no_skill_config_test",
    config = "evals/kept.yaml",  # keep
    data = [
        ":kept_data",  # keep
        "evals/obsolete.txt",
    ],
    skills = [":kept_skill"],  # keep
    visibility = ["//custom:__pkg__"],
)

# keep
promptfoo_validate_test(
    name = "eval_kept_config_test",
    config = "evals/rule-kept.yaml",
    data = [":rule_kept_data"],
    skills = [":rule_kept_skill"],
)

promptfoo_test(
    name = "eval",
    args = ["--repeat", "3"],
    config = "evals/manual-live.yaml",
    env_inherit = ["CODEX_HOME"],
    isolate_codex_home = False,
    skills = [":skill"],
    tags = ["manual-live-tag"],
)
`
	build, err := rule.LoadData("BUILD.bazel", "pkg", []byte(existing))
	if err != nil {
		t.Fatal(err)
	}
	lang := NewLanguage()
	args := generateArgs(repoRoot, packageDir)
	args.File = build
	args.RegularFiles = []string{"SKILL.md"}
	result := lang.GenerateRules(args)
	merger.MergeFile(
		build,
		result.Empty,
		result.Gen,
		merger.PreResolve,
		lang.Kinds(),
		nil,
	)

	got := string(build.Format())
	for _, preserved := range []string{
		`args = ["--manual-validation-option"]`,
		`env = {"MODE": "manual"}`,
		`tags = ["manual-validation-tag"]`,
		`visibility = ["//custom:__pkg__"]`,
		`name = "eval"`,
		`env_inherit = ["CODEX_HOME"]`,
		`isolate_codex_home = False`,
		`tags = ["manual-live-tag"]`,
	} {
		if !strings.Contains(got, preserved) {
			t.Errorf("merged BUILD does not preserve %q:\n%s", preserved, got)
		}
	}
	if count := strings.Count(got, "promptfoo_test("); count != 1 {
		t.Errorf("merged BUILD has %d live tests, want 1:\n%s", count, got)
	}

	validation := findRule(
		t,
		build,
		promptfooValidateKind,
		"eval_config_test",
	)
	if config := validation.AttrString("config"); config !=
		"evals/promptfooconfig.yaml" {
		t.Errorf("validation config = %q, want generated config", config)
	}
	wantData := []string{"evals/cases.yaml", "evals/context.txt"}
	if data := validation.AttrStrings("data"); !reflect.DeepEqual(
		data,
		wantData,
	) {
		t.Errorf("validation data = %q, want %q", data, wantData)
	}
	if skills := validation.AttrStrings("skills"); !reflect.DeepEqual(
		skills,
		[]string{":skill"},
	) {
		t.Errorf("validation skills = %q, want [:skill]", skills)
	}

	noSkill := findRule(
		t,
		build,
		promptfooValidateKind,
		"eval_no_skill_config_test",
	)
	if config := noSkill.AttrString("config"); config != "evals/kept.yaml" {
		t.Errorf("kept config = %q, want evals/kept.yaml", config)
	}
	wantKeptData := []string{
		"evals/cases.yaml",
		"evals/context.txt",
		":kept_data",
	}
	if data := noSkill.AttrStrings("data"); !reflect.DeepEqual(
		data,
		wantKeptData,
	) {
		t.Errorf("no-skill data = %q, want %q", data, wantKeptData)
	}
	if skills := noSkill.AttrStrings("skills"); !reflect.DeepEqual(
		skills,
		[]string{":kept_skill"},
	) {
		t.Errorf("kept no-skill skills = %q, want [:kept_skill]", skills)
	}

	kept := findRule(
		t,
		build,
		promptfooValidateKind,
		"eval_kept_config_test",
	)
	if config := kept.AttrString("config"); config != "evals/rule-kept.yaml" {
		t.Errorf("rule-kept config = %q, want evals/rule-kept.yaml", config)
	}
	if data := kept.AttrStrings("data"); !reflect.DeepEqual(
		data,
		[]string{":rule_kept_data"},
	) {
		t.Errorf("rule-kept data = %q, want [:rule_kept_data]", data)
	}
	if skills := kept.AttrStrings("skills"); !reflect.DeepEqual(
		skills,
		[]string{":rule_kept_skill"},
	) {
		t.Errorf("rule-kept skills = %q, want [:rule_kept_skill]", skills)
	}
	if count := strings.Count(got, "# keep"); count < 4 {
		t.Errorf("merged BUILD lost keep comments:\n%s", got)
	}

	for _, existingRule := range build.Rules {
		if existingRule.Kind() == promptfooTestKind &&
			existingRule.Name() == "eval" {
			wantArgs := []string{"--repeat", "3"}
			if liveArgs := existingRule.AttrStrings("args"); !reflect.DeepEqual(
				liveArgs,
				wantArgs,
			) {
				t.Errorf("live args = %q, want %q", liveArgs, wantArgs)
			}
			return
		}
	}
	t.Error("merged BUILD has no eval promptfoo_test")
}

func TestNoConfigsPreserveConventionallyNamedManualRules(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "pkg")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":      "",
		"evals/cases.yaml": "[]",
	})

	const existing = `promptfoo_validate_test(
    name = "eval_config_test",
    config = "evals/promptfooconfig.yaml",
    data = ["evals/cases.yaml"],
)

promptfoo_validate_test(
    name = "manual_validation",
    config = "evals/manual.yaml",
)

promptfoo_test(
    name = "eval",
    config = "evals/manual-live.yaml",
)
`
	build, err := rule.LoadData("BUILD.bazel", "pkg", []byte(existing))
	if err != nil {
		t.Fatal(err)
	}
	lang := NewLanguage()
	args := generateArgs(repoRoot, packageDir)
	args.File = build
	result := lang.GenerateRules(args)
	if len(result.Gen) != 0 || len(result.Empty) != 0 {
		t.Fatalf(
			"GenerateRules() = %d generated, %d empty; want a no-op",
			len(result.Gen),
			len(result.Empty),
		)
	}
	merger.MergeFile(
		build,
		result.Empty,
		result.Gen,
		merger.PreResolve,
		lang.Kinds(),
		nil,
	)

	findRule(t, build, promptfooValidateKind, "eval_config_test")
	findRule(t, build, promptfooValidateKind, "manual_validation")
	findRule(t, build, promptfooTestKind, "eval")
}

func TestRenamedConfigPreservesExistingTargetsAndHonorsKeep(t *testing.T) {
	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "pkg")
	writeFiles(t, packageDir, map[string]string{
		"BUILD.bazel":                        "",
		"evals/promptfooconfig.quality.yaml": "tests: []",
	})

	const existing = `promptfoo_validate_test(
    name = "eval_config_test",
    config = "evals/promptfooconfig.yaml",
)

# keep
promptfoo_validate_test(
    name = "eval_kept_config_test",
    config = "evals/removed-kept.yaml",
)

promptfoo_validate_test(
    name = "eval_pinned_config_test",
    config = "evals/removed-pinned.yaml",  # keep
)
`
	build, err := rule.LoadData("BUILD.bazel", "pkg", []byte(existing))
	if err != nil {
		t.Fatal(err)
	}
	lang := NewLanguage()
	args := generateArgs(repoRoot, packageDir)
	args.File = build
	result := lang.GenerateRules(args)
	merger.MergeFile(
		build,
		result.Empty,
		result.Gen,
		merger.PreResolve,
		lang.Kinds(),
		nil,
	)

	findRule(t, build, promptfooValidateKind, "eval_config_test")
	quality := findRule(
		t,
		build,
		promptfooValidateKind,
		"eval_quality_config_test",
	)
	if config := quality.AttrString("config"); config !=
		"evals/promptfooconfig.quality.yaml" {
		t.Errorf("replacement config = %q, want quality config", config)
	}
	findRule(t, build, promptfooValidateKind, "eval_kept_config_test")
	pinned := findRule(
		t,
		build,
		promptfooValidateKind,
		"eval_pinned_config_test",
	)
	if config := pinned.AttrString("config"); config !=
		"evals/removed-pinned.yaml" {
		t.Errorf("pinned config = %q, want preserved config", config)
	}
}

func TestCombinedGazelleCreatesFreshSkillPackage(t *testing.T) {
	gazellePath, ok := bazel.FindBinary(
		"gazelle",
		"gazelle_promptfoo_with_skills",
	)
	if !ok {
		t.Fatal("could not find combined Gazelle binary")
	}

	repoRoot := t.TempDir()
	packageDir := filepath.Join(repoRoot, "skills", "example")
	writeFiles(t, repoRoot, map[string]string{
		"MODULE.bazel": "module(name = \"fixture\")\n",
	})
	writeFiles(t, packageDir, map[string]string{
		"SKILL.md":                   "---\nname: example\n---\n",
		"evals/cases.yaml":           "[]",
		"evals/promptfooconfig.yaml": "tests: []",
	})

	runGazelle(t, gazellePath, repoRoot)
	buildPath := filepath.Join(packageDir, "BUILD.bazel")
	firstBytes, err := os.ReadFile(buildPath)
	if err != nil {
		t.Fatalf("read generated BUILD: %v", err)
	}
	first := string(firstBytes)
	generated, err := rule.LoadFile(buildPath, "skills/example")
	if err != nil {
		t.Fatalf("parse generated BUILD: %v", err)
	}

	skill := findRule(t, generated, "skill_library", "skill")
	glob, ok := rule.ParseGlobExpr(skill.Attr("srcs"))
	if !ok || !reflect.DeepEqual(glob.Patterns, []string{"**"}) {
		t.Fatalf("generated skill srcs = %#v, want glob([\"**\"])", glob)
	}
	wantExcludes := []string{"BUILD", "BUILD.bazel", "evals/**"}
	if !reflect.DeepEqual(glob.Excludes, wantExcludes) {
		t.Errorf("generated skill excludes = %q, want %q", glob.Excludes,
			wantExcludes)
	}

	validation := findRule(
		t,
		generated,
		promptfooValidateKind,
		"eval_config_test",
	)
	if config := validation.AttrString("config"); config !=
		"evals/promptfooconfig.yaml" {
		t.Errorf("generated config = %q, want Promptfoo config", config)
	}
	if data := validation.AttrStrings("data"); !reflect.DeepEqual(
		data,
		[]string{"evals/cases.yaml"},
	) {
		t.Errorf("generated data = %q, want [evals/cases.yaml]", data)
	}
	if skills := validation.AttrStrings("skills"); !reflect.DeepEqual(
		skills,
		[]string{":skill"},
	) {
		t.Errorf("generated skills = %q, want [:skill]", skills)
	}
	for _, load := range []string{
		`load("@rules_promptfoo//promptfoo:defs.bzl", "promptfoo_validate_test")`,
		`load("@rules_skill//skill:defs.bzl", "skill_library")`,
	} {
		if !strings.Contains(first, load) {
			t.Errorf("generated BUILD does not contain %q:\n%s", load, first)
		}
	}

	runGazelle(t, gazellePath, repoRoot)
	secondBytes, err := os.ReadFile(buildPath)
	if err != nil {
		t.Fatalf("read regenerated BUILD: %v", err)
	}
	if second := string(secondBytes); second != first {
		t.Fatalf("second Gazelle run changed BUILD:\nfirst:\n%s\nsecond:\n%s",
			first, second)
	}
}

func TestApparentLoads(t *testing.T) {
	lang := NewLanguage().(language.ModuleAwareLanguage)
	tests := []struct {
		name     string
		apparent string
		want     string
	}{
		{
			name:     "module mapping",
			apparent: "company_promptfoo",
			want:     "@company_promptfoo//promptfoo:defs.bzl",
		},
		{
			name: "legacy fallback",
			want: "@rules_promptfoo//promptfoo:defs.bzl",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loads := lang.ApparentLoads(func(module string) string {
				if module == "rules_promptfoo" {
					return tt.apparent
				}
				return ""
			})
			if len(loads) != 1 || loads[0].Name != tt.want {
				t.Fatalf(
					"ApparentLoads() = %#v, want load %q",
					loads,
					tt.want,
				)
			}
			wantSymbols := []string{
				promptfooTestKind,
				promptfooValidateKind,
			}
			if !reflect.DeepEqual(loads[0].Symbols, wantSymbols) {
				t.Errorf(
					"load symbols = %q, want %q",
					loads[0].Symbols,
					wantSymbols,
				)
			}
		})
	}
}

func generateArgs(repoRoot, packageDir string) language.GenerateArgs {
	return language.GenerateArgs{
		Config: &config.Config{
			RepoRoot:            repoRoot,
			ValidBuildFileNames: []string{"BUILD.bazel", "BUILD"},
		},
		Dir:  packageDir,
		File: rule.EmptyFile("BUILD.bazel", "pkg"),
		Rel:  "pkg",
	}
}

func findRule(
	t *testing.T,
	build *rule.File,
	kind string,
	name string,
) *rule.Rule {
	t.Helper()
	for _, candidate := range build.Rules {
		if candidate.Kind() == kind && candidate.Name() == name {
			return candidate
		}
	}
	t.Fatalf("BUILD has no %s(%q):\n%s", kind, name, build.Format())
	return nil
}

func hasRule(build *rule.File, kind, name string) bool {
	for _, candidate := range build.Rules {
		if candidate.Kind() == kind && candidate.Name() == name {
			return true
		}
	}
	return false
}

func runGazelle(t *testing.T, gazellePath, repoRoot string) {
	t.Helper()
	command := exec.Command(gazellePath)
	command.Dir = repoRoot
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("run combined Gazelle: %v\n%s", err, output)
	}
}

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	for path, content := range files {
		fullPath := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
