---
title: Js
description: Js
---

Repository JavaScript tooling dependencies belong to the pnpm workspace in
`tools/`. Its `package.json` and `pnpm-lock.yaml` own their requested and resolved
versions. `//tools:node_modules` owns the Bazel package store and links; loading
the repository root package does not load the npm extension's generated rules.

Run pnpm with `//tools/pnpm`, passing `--dir` with the absolute path to `tools/`
for dependency updates. Keep lifecycle scripts disabled (`--ignore-scripts`)
when updating the lockfile.

A security override selects patched KaTeX where Mermaid's upstream range
retains a vulnerable release. PostCSS is pinned explicitly to supply a secure
version to its CLI and plugins' peers. Mermaid and Three.js are exact pins;
the Dagre spacing patch is maintained with each Mermaid upgrade.

## TypeScript

Use `ts_project` from `@aspect_rules_ts//ts:defs.bzl` to type-check and compile
TypeScript. It integrates with the existing `js_binary`, `js_test`, and npm
providers. `tools/package.json` owns the compiler pin; the TypeScript module
extension reads it with `version_from` and uses upstream integrity-checked
compiler archives. The Node runtime and `rules_js` retain their existing pins.

Keep compiler options with the consuming project. Set `strict: true` and
`noEmitOnError: true`, choose `transpiler = "tsc"`, and set
`supports_workers = 0` for the native compiler. Node programs use NodeNext
module resolution and explicit Node declarations; browser programs declare
their DOM libraries. Add the required npm type targets to `deps`.

Use `.ts` source and a project-owned `package.json` with `"type": "module"`
for ES modules. Include that manifest in `srcs` for compiler resolution and
in `data` for runtime packaging. When multiple Bazel packages share it, expose
one `copy_to_bin` target from the manifest's owning package and consume that
target in both attributes. Compilation emits `.js`, and
`declaration = True` additionally emits `.d.ts` for dependent projects.
Consumers depend on the `ts_project` target and execute its generated JavaScript
through `js_binary` or package it for the browser. Generated JavaScript remains
a build output. The native compiler provides `tsc`; it does not provide the
former JavaScript `tsserver` entry point.

The [Node](https://github.com/alwaldend/src/blob/master/tools/js/test/typescript/node/BUILD.bazel),
[browser](https://github.com/alwaldend/src/blob/master/tools/js/test/typescript/browser/BUILD.bazel), and
[shared model](https://github.com/alwaldend/src/blob/master/tools/js/test/typescript/model/BUILD.bazel) fixtures demonstrate the direct
rule configuration and declaration dependencies. Run their full acceptance with:

```sh
bazel test //tools/js/test/typescript:typescript_test
```

The test runs emitted modules in pinned Node and Chrome, exercises a Three.js
vector calculation, checks a real compiler type error, and retains
`node-result.json`, `browser.png`, `type-error.log`, and
`verification.json` with output hashes in its test output directory.
