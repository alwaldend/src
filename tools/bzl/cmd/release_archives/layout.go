package main

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"
)

// releaseModuleName reads the literal module identity from a release manifest.
var releaseModuleName = regexp.MustCompile(`(?m)^module\(\s*name\s*=\s*"([a-z][a-z0-9._-]*)"`)

// nativeReleaseArchive identifies a package's explicit archive build target.
var nativeReleaseArchive = regexp.MustCompile(`(?m)^\s*name\s*=\s*"release_archive",`)

// archiveLayout describes first-party module identities and the source trees
// whose archives are built by this command.
type archiveLayout struct {
	// modules maps each source directory to its standalone module name.
	modules map[string]string
	// sources lists the owners that use the generic archive builder.
	sources []string
}

// readArchiveLayout discovers the project-owned release manifests in a commit.
func readArchiveLayout(ctx context.Context, root, commit string) (archiveLayout, error) {
	tree, err := gitOutput(ctx, root, "ls-tree", "-r", "--name-only", "-z", commit)
	if err != nil {
		return archiveLayout{}, fmt.Errorf("list release sources at %s: %w", commit, err)
	}
	layout := archiveLayout{modules: make(map[string]string)}
	seen := make(map[string]string)
	for _, file := range strings.Split(string(tree), "\x00") {
		if !strings.HasSuffix(file, "/release.MODULE.bazel") {
			continue
		}
		data, err := gitOutput(ctx, root, "show", commit+":"+file)
		if err != nil {
			return archiveLayout{}, fmt.Errorf("read release manifest %q: %w", file, err)
		}
		match := releaseModuleName.FindSubmatch(data)
		if match == nil {
			return archiveLayout{}, fmt.Errorf("missing module name in %q", file)
		}
		owner, name := path.Dir(file), string(match[1])
		if previous, exists := seen[name]; exists {
			return archiveLayout{}, fmt.Errorf("module %q is declared by %q and %q", name, previous, owner)
		}
		seen[name] = owner
		layout.modules[owner] = name
		build, err := gitOutput(ctx, root, "show", commit+":"+owner+"/BUILD.bazel")
		if err != nil {
			return archiveLayout{}, fmt.Errorf("read release build definitions for %q: %w", owner, err)
		}
		if !nativeReleaseArchive.Match(build) {
			layout.sources = append(layout.sources, owner)
		}
	}
	return layout, nil
}

// archiveLabels maps repository-relative Bazel labels to the release layout.
func archiveLabels(owner string, modules map[string]string) *strings.Replacer {
	var pairs []string
	for project, name := range modules {
		prefix := "@" + name + "//"
		if project == owner {
			prefix = "//"
			pairs = append(pairs, `strip_import_prefix = "/`+project+`"`, `strip_import_prefix = "/"`)
		}
		pairs = append(pairs, "//"+project+":", prefix+":", "//"+project+"/", prefix)
	}
	return strings.NewReplacer(pairs...)
}

// containsBazelLabels identifies archive source files with Bazel label strings.
func containsBazelLabels(name string) bool {
	return path.Base(name) == "BUILD" || strings.HasSuffix(name, ".bazel") || strings.HasSuffix(name, ".bzl")
}

// monorepoOpenSpecVisibility identifies the generated package directive for a
// monorepo tree whose package group is outside the archived module.
var monorepoOpenSpecVisibility = regexp.MustCompile(`(?m)^package\(transitive_visibility = "//(?:projects|tools):transitive_consumers"\)\n(?:\n|$)`)

// relocateOpenSpecPackage removes the reference to the enclosing monorepo's
// package group while preserving the module's own OpenSpec targets and visibility.
func relocateOpenSpecPackage(source string) string {
	return monorepoOpenSpecVisibility.ReplaceAllString(source, "")
}
