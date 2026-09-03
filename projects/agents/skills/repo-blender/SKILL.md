---
name: repo-blender
description: >-
  Operate Blender assets in this repository using the repository-pinned
  Blender and Bazel entry points, task-owned outputs, safe background
  execution, and explicit GUI or MCP boundaries. Use when an agent opens,
  inspects, scripts, renders, saves, converts, or exports a .blend file here;
  use alongside blender-reference-fidelity when supplied-reference likeness
  controls acceptance. Do not use as a modeling, sculpting, style, or
  visual-quality guide.
---

# Operate Blender in this repository

Use the repository-pinned Blender, contain every task-owned side effect, and
verify the saved or rendered artifact that the user will actually receive.
Blender command success proves execution, not a correct scene or image.

## Use the repository toolchain

Read the nearest `AGENTS.md`, owning project README and BUILD file, and
`tools/blender/README.md`. Read the current version and archive contract from
`tools/blender/binary_toolchain.json`; do not hard-code a remembered version.

Follow `bazel-agent` for every Bazel invocation. Prefer an owning package's
purpose-built `blender_binary` target when it implements the operation;
otherwise use `//tools/blender:blender`. Do not substitute `blender` from
`PATH`, Flatpak Blender, a separately downloaded binary, or an extracted ELF
for batch work, reproducible verification, or deliverables. The explicitly
requested, disposable Flatpak MCP-host exception in
[execution-modes.md](references/execution-modes.md) does not replace the pinned
toolchain. Changing the toolchain or Bazel targets additionally requires the
owning Bazel and dependency skills.

Pass Blender arguments after Bazel's separator:

```sh
bazel_agent bazel run //tools/blender:blender -- --version
```

## Select the least-stateful execution mode

Use background Blender for ordinary inspection, scripted changes, simulation,
conversion, validation, saving, and rendering. Use a foreground GUI or MCP
session only when the operation genuinely requires a live window, viewport,
modal operator, or interactive add-on. Configured MCP client tools do not mean
that Blender is running or that a compatible listener exists.

Read [execution-modes.md](references/execution-modes.md) before choosing or
starting a Blender process. For live UI-context operators, also follow its
settled-state and effect-verification protocol.

Do not enable add-ons, change preferences or startup files, start a listener,
or launch a second Blender merely because the task mentions Blender. Probe an
existing interactive connection read-only, verify its Blender identity and
open-file identity, and reuse it only when it matches the task. Starting a GUI
still follows the environment's GUI and host-execution approval boundary.

## Protect inputs and contain outputs

Put one-off scripts, copied candidates, Blender configuration, temporary
storage, caches, bakes, logs, renders, exports, manifests, and backup files in
a task-specific repository-root `out/<task>/` directory. Point configurable
Blender paths there; do not allow backups or caches to appear beside tracked
sources.

Treat tracked, user-supplied, baseline, and declared protected `.blend` files
as immutable inputs unless the user explicitly authorizes their promotion or
replacement. Hash important inputs, edit an explicit candidate copy, and
rehash the inputs after Blender exits. Read-only inspection must not save.
Never let a foreground editor and background renderer mutate or read the same
file while it is being saved; publish a checkpoint copy and render that frozen
snapshot.

Before mutating an existing asset, rendering review evidence, or exporting a
deliverable, read
[asset safety](references/asset-safety-and-verification.md).

## Make automation explicit and inspectable

Set the scene, view layer, active object, mode, camera, render engine,
resolution, output path, save path, and export path explicitly whenever they
affect the result. Do not rely on user preferences, an unsaved UI state, or an
implicitly active object unless that live state is itself the named input.

Prefer dedicated Blender MCP summary, inspection, navigation, screenshot, and
render tools over arbitrary code. Use code execution only when the focused
tools cannot perform the operation. Inspect progressively instead of dumping a
large scene, set selection and mode deliberately, and return structured data.

After a mutating operation, clean-reopen the exact saved candidate with the
pinned background Blender and verify properties material to the request. For
images, inspect the final pixels. Record the Blender/tool identity, input and
output paths, relevant hashes, and exit state. Artistic likeness, construction,
and presentation quality remain the responsibility of the applicable
specialist workflow; technical reopen or render success cannot pass them.
