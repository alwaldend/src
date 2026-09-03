---
name: repo-gazelle-plugin
description: >-
  Create or update standalone Gazelle language plugins and integrate them with
  this Bazel monorepo. Use when changing BUILD-rule generation behavior; use
  repo-bazel for ordinary generated BUILD updates that do not change a plugin.
---

# Add a repository Gazelle plugin

Use `projects/rules_docs_gazelle` as the reference implementation. Also follow
the `bazel-nested-module`, `repo-bazel`, and `bazel-agent` skills: a repository
Gazelle plugin is both a Go language extension and a standalone Bzlmod module.

## Preserve workspace boundaries

- Put the plugin in an underscore-named `projects/rules_*_gazelle` module with
  its own `MODULE.bazel`, lockfile, root `BUILD.bazel`, README, and `go.mod`.
  Reuse the nested-module Bazel configuration symlinks and documentation
  conventions rather than treating it as a root package.
- Keep production rules and their Gazelle plugin as separate modules. The
  plugin may depend on the rules module whose symbols it generates, while the
  rules module must not gain a runtime dependency on its generator.
- Declare Gazelle and `rules_go` only where the plugin needs them. Declare a
  generator as `dev_dependency = True` in consumers when it is needed only to
  edit BUILD files. Use sibling `local_path_override` entries for unpublished
  modules in this checkout and update module locks through Bazel.
- Do not traverse out of `GenerateArgs.Config.RepoRoot`. A plugin running in a
  nested workspace must not infer packages, visibility, or files from its
  parent workspace. Root Gazelle also does not cross directories listed in
  `.bazelignore`.
- For any physical directory walk, treat both `BUILD` and `BUILD.bazel`, plus
  configured valid BUILD filenames, as package boundaries. Propagate stat and
  walk errors into a non-destructive no-op; an unreadable boundary is not
  evidence that traversal is safe.

## Implement the smallest language surface

- Expose a public Go library whose `NewLanguage` returns
  `language.Language`, and add compile-time interface assertions. Implement
  `language.Configurer` or `language.Resolver` only when the plugin genuinely
  needs directives, import resolution, or dependency generation.
- When a generated rule loads a symbol from a Bzlmod dependency, implement
  `language.ModuleAwareLanguage`. Resolve the canonical module name with the
  `moduleToApparentName` callback in `ApparentLoads`; a consumer may rename the
  repository. Retain the canonical name as the fallback when the callback
  returns an empty string. As in `rules_docs_gazelle`, make legacy `Loads`
  fail loudly so a new call path cannot silently emit the wrong repository
  name.
- Keep rule kinds, names, source matching, loads, and imports deterministic.
  Return one `Imports` entry for every generated rule when the language uses
  Gazelle's import-resolution pipeline.

## Make merging conservative

- In `Kinds`, mark an attribute mergeable only when the generator owns its
  value and can safely reconcile existing content. Prefer non-mergeable
  attributes when users may customize `srcs`, `deps`, `visibility`, or other
  policy-bearing values.
- Generate only from explicit, package-local evidence. Do not invent a BUILD
  package merely to make a rule possible unless that behavior is part of the
  plugin's contract.
- Do not return an `Empty` rule just because generation inputs disappeared if
  doing so could delete a hand-maintained rule. Preserve unknown attributes
  and manual values during merge. Generated names must not collide with
  unrelated targets.
- Make a second Gazelle run produce no diff. Idempotence and preservation of
  manual edits are part of correctness, not cosmetic properties.

## Test generation and merging

Add focused Go tests patterned after
`projects/rules_docs_gazelle/gazelle/gazelle_test.go`. Cover:

- the positive generation case and exact rule kind, name, and owned values;
- missing or irrelevant inputs, including the absence of an existing package
  when packages are a prerequisite;
- workspace-root and nearest-package boundary behavior without escaping a
  nested repository;
- apparent repository-name mapping and the canonical fallback;
- a real Gazelle merge that preserves every supported manual attribute;
- input removal without accidental deletion; and
- repeat generation or an integration fixture that proves idempotence.

Test observable BUILD output and merge behavior rather than only calling
helpers. Export a small plugin-only `gazelle_binary`, like
`//gazelle:gazelle_docs`, so the standalone module can exercise its extension
without initializing unrelated languages from the parent repository.

## Integrate with the monorepo

1. Follow `bazel-nested-module` for `.bazelignore`, parent dependency and local
   override entries, documentation aggregation, and project exclusions.
2. Add the public language library to the `languages` of the root
   `//:gazelle_bin`. Keep the plugin dependency development-only unless root
   BUILD files load runtime symbols from that module.
3. Register the new workspace in the `full-repo-check` runner. Update its
   tests, documentation, workspace count, and expected command count so both
   `build //...` and `test //...` run inside the module.
4. Run the plugin against its own nested workspace when appropriate, then run
   the root `//:gazelle` target. Inspect every generated diff, including files
   outside the intended package, and fix the generator rather than accepting
   unrelated churn. Run each generator again and require an empty second diff.

## Verify with Bazel

Invoke every Bazel command through `bazel_agent`. From the plugin workspace,
update its lock and validate its complete, intentionally small graph:

```sh
bazel_agent mod deps --lockfile_mode=update
bazel_agent bazel test //...
bazel_agent bazel build //...
```

From the repository root, update the parent lock when its module graph changed,
run Gazelle, inspect the resulting diff, and validate the integration:

```sh
bazel_agent mod deps --lockfile_mode=update
bazel_agent bazel run //:gazelle
bazel_agent bazel test //path/to/affected/package:all
bazel_agent bazel build //path/to/affected/package:all
bazel_agent bazel test //:buildifier_test
```

Run root Gazelle a second time and confirm that it changes nothing. Finish with
the `full-repo-check` skill when a new nested workspace was added. Never use a
direct `bazel` invocation or hand-edit a generated lockfile.
