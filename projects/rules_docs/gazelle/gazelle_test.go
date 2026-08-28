package gazelle

import (
	"os"
	"path/filepath"
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

	build := rule.EmptyFile("BUILD.bazel", "parent/child")
	result := NewLanguage().GenerateRules(language.GenerateArgs{
		Rel:          "parent/child",
		File:         build,
		RegularFiles: []string{"README.md", "guide.md"},
	})

	if len(result.Gen) != 1 || len(result.Imports) != 1 {
		t.Fatalf("GenerateRules() generated %d rules and %d imports, want 1 each",
			len(result.Gen), len(result.Imports))
	}
	docs := result.Gen[0]
	if docs.Kind() != docsKind || docs.Name() != docsName {
		t.Fatalf("generated rule = %s(%q), want %s(%q)",
			docs.Kind(), docs.Name(), docsKind, docsName)
	}
	glob, ok := rule.ParseGlobExpr(docs.Attr("srcs"))
	if !ok || !reflect.DeepEqual(glob.Patterns, []string{"*.md"}) {
		t.Fatalf("srcs = %#v, want glob([\"*.md\"])", glob)
	}
	if got := docs.AttrStrings("visibility"); got != nil {
		t.Errorf("visibility = %q, want omitted without repository context", got)
	}
	if len(result.Empty) != 0 {
		t.Errorf("Empty contains %d rules, want none", len(result.Empty))
	}
}

func TestParentVisibilityUsesNearestPackage(t *testing.T) {
	repoRoot := t.TempDir()
	parent := filepath.Join(repoRoot, "parent")
	child := filepath.Join(parent, "not_a_package", "child")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(parent, "BUILD.bazel"),
		[]byte(""),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	got := parentVisibility(language.GenerateArgs{
		Config: &config.Config{
			RepoRoot:            repoRoot,
			ValidBuildFileNames: []string{"BUILD.bazel"},
		},
		Dir: child,
	})
	want := []string{"//parent:__pkg__"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parentVisibility() = %q, want %q", got, want)
	}
}

func TestParentVisibilityDoesNotEscapeNestedWorkspace(t *testing.T) {
	outerRoot := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(outerRoot, "BUILD.bazel"),
		[]byte(""),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Join(outerRoot, "nested")
	if err := os.Mkdir(repoRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	got := parentVisibility(language.GenerateArgs{
		Config: &config.Config{
			RepoRoot:            repoRoot,
			ValidBuildFileNames: []string{"BUILD.bazel"},
		},
		Dir: repoRoot,
	})
	if got != nil {
		t.Fatalf("parentVisibility() = %q, want omitted at workspace root", got)
	}
}

func TestGenerateRulesRequiresBuildAndReadme(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args language.GenerateArgs
	}{
		{
			name: "no BUILD file",
			args: language.GenerateArgs{RegularFiles: []string{"README.md"}},
		},
		{
			name: "no README",
			args: language.GenerateArgs{File: rule.EmptyFile("BUILD.bazel", "pkg")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewLanguage().GenerateRules(tt.args)
			if len(result.Gen) != 0 || len(result.Empty) != 0 {
				t.Fatalf("GenerateRules() = %d generated, %d empty; want none",
					len(result.Gen), len(result.Empty))
			}
		})
	}
}

func TestMergePreservesManualAttributes(t *testing.T) {
	t.Parallel()

	const existing = `load("//docs:defs.bzl", "docs_filegroup")

docs_filegroup(
    name = "docs",
    srcs = [],
    deps = [":extra"],
    prefix = "custom/prefix",
    visibility = ["//custom:__pkg__"],
)
`
	build, err := rule.LoadData("BUILD.bazel", "parent/child", []byte(existing))
	if err != nil {
		t.Fatal(err)
	}
	lang := NewLanguage()
	result := lang.GenerateRules(language.GenerateArgs{
		Rel:          "parent/child",
		File:         build,
		RegularFiles: []string{"README.md"},
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
		"srcs = []",
		`deps = [":extra"]`,
		`prefix = "custom/prefix"`,
		`visibility = ["//custom:__pkg__"]`,
	} {
		if !strings.Contains(got, preserved) {
			t.Errorf("merged BUILD does not preserve %q:\n%s", preserved, got)
		}
	}
}

func TestNoReadmeDoesNotDeleteManualRule(t *testing.T) {
	t.Parallel()

	build, err := rule.LoadData(
		"BUILD.bazel",
		"pkg",
		[]byte("docs_filegroup(name = \"docs\", srcs = [])\n"),
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

	if got := string(build.Format()); !strings.Contains(got, "docs_filegroup(") {
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
		{name: "module mapping", apparent: "company_docs", want: "@company_docs//docs:defs.bzl"},
		{name: "legacy fallback", want: "@rules_docs//docs:defs.bzl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loads := lang.ApparentLoads(func(module string) string {
				if module == "rules_docs" {
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
