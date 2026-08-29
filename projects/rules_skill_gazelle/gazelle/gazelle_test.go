package gazelle

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/merger"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func TestGenerateRules(t *testing.T) {
	t.Parallel()

	result := NewLanguage().GenerateRules(language.GenerateArgs{
		Rel:          "projects/agents/skills/example",
		RegularFiles: []string{"SKILL.md", "reference.md"},
	})

	if len(result.Gen) != 1 || len(result.Imports) != 1 {
		t.Fatalf("GenerateRules() generated %d rules and %d imports, want 1 each",
			len(result.Gen), len(result.Imports))
	}
	skill := result.Gen[0]
	if skill.Kind() != skillKind || skill.Name() != skillName {
		t.Fatalf("generated rule = %s(%q), want %s(%q)",
			skill.Kind(), skill.Name(), skillKind, skillName)
	}
	glob, ok := rule.ParseGlobExpr(skill.Attr("srcs"))
	if !ok || !reflect.DeepEqual(glob.Patterns, []string{"**"}) {
		t.Fatalf("srcs patterns = %#v, want glob([\"**\"])", glob.Patterns)
	}
	wantExcludes := []string{"BUILD.bazel", "BUILD", "evals/**"}
	if !reflect.DeepEqual(glob.Excludes, wantExcludes) {
		t.Fatalf("srcs excludes = %#v, want %#v", glob.Excludes, wantExcludes)
	}
	if len(result.Empty) != 0 {
		t.Errorf("Empty contains %d rules, want none", len(result.Empty))
	}
}

func TestGenerateRulesUsesConfiguredBuildFileNames(t *testing.T) {
	t.Parallel()

	result := NewLanguage().GenerateRules(language.GenerateArgs{
		Config: &config.Config{
			ValidBuildFileNames: []string{
				"BUILD.custom",
				"BUILD.bazel",
				"BUILD.custom",
			},
		},
		Rel:          "skills/example",
		RegularFiles: []string{"SKILL.md"},
	})

	if len(result.Gen) != 1 {
		t.Fatalf("GenerateRules() generated %d rules, want 1", len(result.Gen))
	}
	glob, ok := rule.ParseGlobExpr(result.Gen[0].Attr("srcs"))
	if !ok {
		t.Fatalf("srcs = %#v, want a glob", result.Gen[0].Attr("srcs"))
	}
	wantExcludes := []string{
		"BUILD.bazel",
		"BUILD",
		"BUILD.custom",
		"evals/**",
	}
	if !reflect.DeepEqual(glob.Excludes, wantExcludes) {
		t.Fatalf("srcs excludes = %#v, want %#v", glob.Excludes, wantExcludes)
	}
}

func TestSkillSourceExcludesFallsBackForEmptyConfig(t *testing.T) {
	t.Parallel()

	got := skillSourceExcludes(language.GenerateArgs{Config: &config.Config{}})
	want := []string{"BUILD.bazel", "BUILD", "evals/**"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("skillSourceExcludes() = %#v, want %#v", got, want)
	}
}

func TestGenerateRulesRequiresSkillInstructions(t *testing.T) {
	t.Parallel()

	result := NewLanguage().GenerateRules(language.GenerateArgs{
		File:         rule.EmptyFile("BUILD.bazel", "pkg"),
		RegularFiles: []string{"README.md"},
	})
	if len(result.Gen) != 0 || len(result.Empty) != 0 {
		t.Fatalf("GenerateRules() = %d generated, %d empty; want none",
			len(result.Gen), len(result.Empty))
	}
}

func TestGenerateRulesSkipsRepositoryRoot(t *testing.T) {
	t.Parallel()

	result := NewLanguage().GenerateRules(language.GenerateArgs{
		File:         rule.EmptyFile("BUILD.bazel", ""),
		RegularFiles: []string{"SKILL.md"},
	})
	if len(result.Gen) != 0 || len(result.Empty) != 0 {
		t.Fatalf("GenerateRules() = %d generated, %d empty; want none",
			len(result.Gen), len(result.Empty))
	}
}

func TestMergePreservesManualAttributes(t *testing.T) {
	t.Parallel()

	const existing = `load("//skill:defs.bzl", "skill_library")

skill_library(
    name = "skill",
    srcs = ["SKILL.md"],
    tags = ["manual"],
    visibility = ["//custom:__pkg__"],
)
`
	build, err := rule.LoadData("BUILD.bazel", "skills/example", []byte(existing))
	if err != nil {
		t.Fatal(err)
	}
	lang := NewLanguage()
	result := lang.GenerateRules(language.GenerateArgs{
		Rel:          "skills/example",
		File:         build,
		RegularFiles: []string{"SKILL.md"},
	})
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
		`srcs = ["SKILL.md"]`,
		`tags = ["manual"]`,
		`visibility = ["//custom:__pkg__"]`,
	} {
		if !strings.Contains(got, preserved) {
			t.Errorf("merged BUILD does not preserve %q:\n%s", preserved, got)
		}
	}
}

func TestNoSkillInstructionsDoesNotDeleteManualRule(t *testing.T) {
	t.Parallel()

	build, err := rule.LoadData(
		"BUILD.bazel",
		"skills/example",
		[]byte("skill_library(name = \"skill\", srcs = [])\n"),
	)
	if err != nil {
		t.Fatal(err)
	}
	lang := NewLanguage()
	result := lang.GenerateRules(language.GenerateArgs{File: build})
	merger.MergeFile(
		build,
		result.Empty,
		result.Gen,
		merger.PreResolve,
		lang.Kinds(),
		nil,
	)

	if got := string(build.Format()); !strings.Contains(got, "skill_library(") {
		t.Fatalf("manual rule was deleted:\n%s", got)
	}
}

func TestApparentLoads(t *testing.T) {
	t.Parallel()

	lang := NewLanguage().(language.ModuleAwareLanguage)
	tests := []struct {
		name     string
		apparent string
		want     string
	}{
		{
			name:     "module mapping",
			apparent: "company_skills",
			want:     "@company_skills//skill:defs.bzl",
		},
		{
			name: "legacy fallback",
			want: "@rules_skill//skill:defs.bzl",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loads := lang.ApparentLoads(func(module string) string {
				if module == "rules_skill" {
					return tt.apparent
				}
				return ""
			})
			if len(loads) != 1 || loads[0].Name != tt.want {
				t.Fatalf("ApparentLoads() = %#v, want load %q", loads, tt.want)
			}
		})
	}
}
