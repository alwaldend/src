## Context

See proposal.md. `rules_js` 2.8.3 supplies Node execution and npm linking;
its README points to the companion `rules_ts` for TypeScript. The latter's
3.10.1 release accepts that existing rules_js pin and mirrors TypeScript 7.0.2
with integrity-checked native platform packages.

## Goals / Non-Goals

Provide strict type checking and emitted ES modules/declarations through the
maintained `ts_project` API. Keep compiler configurations with their consumers.
This prerequisite does not migrate product code or add a bundler or custom
TypeScript rule.

## Decisions

- Pin approved npm packages in tools/package.json and regenerate its pnpm lock.
  Configure the TypeScript extension using `version_from` so the compiler
  version has one owner. Preserve the existing Node and JavaScript rule pins.
- Use the upstream compiler, validator, and declaration providers directly.
  Select `transpiler = "tsc"` and disable the JavaScript persistent worker for
  the native compiler. Avoid a repository-specific compilation wrapper.
- Use ordinary `.ts` inputs and emitted `.js`/`.d.ts` files with a project-owned
  `package.json` declaring `"type": "module"`. Supply the manifest as both a
  compiler input and runtime data.
  Separate Node and DOM compiler contexts; a shared fixture verifies actual
  declaration and runtime dependency propagation across packages.
- Write E2E acceptance before fixtures or build implementation. Execute the
  compiled Node program and browser module, then invoke the pinned compiler
  on an invalid assignment and require its semantic type-error diagnostic.
  Retain JSON results, compiler diagnostics, hashes, and a browser canvas PNG.
- The browser fixture consumes Three.js declarations through a type-only
  import. Adding the runtime scene engine remains with the diagram PR.

## Risks / Trade-offs

- Compiler and declaration packages add build dependencies and upstream
  transitive packages. Pin them through the existing owners, disable package
  lifecycle scripts, and review generated lock changes.
- Native compiler packaging differs from prior JavaScript TypeScript releases.
  Exercise both successful compilation and a real error through the packaged
  compiler, including its declared runtime files.
- A successful compile alone does not establish runnable output. Exercise
  emitted modules in pinned Node and Chrome, with external requests blocked.

## Migration Plan

Deliver this trunk-based prerequisite first. The renderer and snapshot
migrations consume it after it merges, retaining their separate PR scopes.
Existing JavaScript tools and website consumers remain part of validation.
