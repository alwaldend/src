---
name: repo-gazelle-plugin
description: >-
  Create or update Gazelle language plugins and integrate them with
  this Bazel monorepo. Use when changing BUILD-rule generation behavior; use
  repo-bazel for ordinary generated BUILD updates that do not change a plugin.
---

# Add a repository Gazelle plugin

Use `tools/rules_docs_gazelle` as the reference implementation.
Follow `repo-bazel` and `al-bazel` for builds and release verification.

## Place the plugin

- Put the plugin under `tools/rules_*_gazelle` with its own Bazel packages
  and README. Use workspace-relative labels when building the plugin and its
  rules in the monorepo.
- Keep production rules separate from their Gazelle plugin. The plugin may
  depend on the rules it generates. Use a project-owned
  `release.MODULE.bazel` when the plugin is published independently.
- Keep discovery within `GenerateArgs.Config.RepoRoot`.
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
  returns an empty string. The root Gazelle target passes the in-tree label
  through the plugin's load-label flag. Make legacy `Loads`
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
`tools/rules_docs_gazelle/gazelle/gazelle_test.go`. Cover:

- the positive generation case and exact rule kind, name, and owned values;
- missing or irrelevant inputs, including the absence of an existing package
  when packages are a prerequisite;
- workspace-root and nearest-package boundary behavior without escaping a
  consumer workspace;
- apparent repository-name mapping and the canonical fallback;
- a real Gazelle merge that preserves every supported manual attribute;
- input removal without accidental deletion; and
- repeat generation or an integration fixture that proves idempotence.

Test observable BUILD output and merge behavior rather than only calling
helpers. Export a small plugin-only `gazelle_binary`, like
`//gazelle:gazelle_docs`, so the standalone module can exercise its extension
without initializing unrelated languages from the parent repository.

## Integrate with the monorepo

1. Add the public language library to the root `gazelle_binary`.
2. Set its local load label on the root Gazelle target. External consumers
   use the plugin's module-aware load.
3. Run the root Gazelle target over an affected scope. Review all generated
   changes and repeat the run to verify that the output is stable.

## Verify with Bazel

From the repository root, run the plugin tests and check its integration:

```sh
bazel test //tools/rules_docs_gazelle/...
bazel run //:gazelle
bazel test //:buildifier_test
```

Update the root lock through `bazel mod deps --lockfile_mode=update` when
module dependencies change. For a published plugin, build its executable and
use it in a separate consumer workspace through its registry archive.
