## Why

The operator explicitly approved Google go-jsonnet v0.22.0, including its formatter,
to remove repeated workflow JSON. Shared execution tooling belongs in tools/jsonnet;
workflow definitions and generation contracts belong to projects/comfyui_tools.

## What Changes

- Pin the upstream Linux x86_64 compiler/formatter archive and publisher checksum
  through the existing binary toolchain extension. Add no Bazel rules dependency.
- Provide a sandboxed, declared-input Jsonnet-to-JSON rule and runnable tools.
- Document tool use and verify compiler, formatter and generated build outputs.

## Capabilities

### New Capabilities

- `jsonnet-tooling`: reproducible Jsonnet execution during builds.

## Impact

No installation on the shared host, rental runtime dependency, live state or
infrastructure mutation. Only the specifically approved upstream release is added.
