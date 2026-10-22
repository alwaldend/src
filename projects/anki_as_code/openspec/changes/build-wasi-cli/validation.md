## WASI acceptance

Observed on 2026-10-22 with the repository-pinned Go 1.26.5 and wazero 1.12.0.

```sh
bazel_agent bazel build -c opt --//projects/anki_as_code:wasm=true //projects/anki_as_code
bazel_agent bazel cquery -c opt --//projects/anki_as_code:wasm=true //projects/anki_as_code --output=files
bazel_agent bazel run --script_path=out/anki-as-code-wasm/revision/wazero @com_github_tetratelabs_wazero//cmd/wazero
```

The output has WebAssembly version-1 magic and imports WASI Preview 1. Neither
`ankigrpc/native` nor `ankigrpc/wasm` appears in its configured dependency graph.
The optimized module is 27,902,544 bytes, SHA-256
`814294d3e1753dde68b4c8a73002db97c2f834155074d33aaf04552f97ef094c`.
The task copy is `out/anki-as-code-wasm/revision/anki_as_code.wasm`.

The actual module passed root help, stdin JSON config export, config dump with
semantic JSON equality, malformed-input rejection before directory creation,
and `server serve` rejection with `no local engines are compiled`. The fixture
includes a deck, a note type, a card template, Unicode, quotes, and a newline.
The exported tree, input JSON, stdout/stderr, and compiled runtime cache remain
under ignored `out/anki-as-code-wasm/revision/`.

Repeat the configuration round-trip using the retained public fixture:

```sh
out/anki-as-code-wasm/revision/wazero run -cachedir out/anki-as-code-wasm/revision/runtime-cache \
  -mount "$PWD/out/anki-as-code-wasm/revision/workflow:/" \
  out/anki-as-code-wasm/revision/anki_as_code.wasm config export --output /configs \
  < out/anki-as-code-wasm/revision/input.json
out/anki-as-code-wasm/revision/wazero run -cachedir out/anki-as-code-wasm/revision/runtime-cache \
  -mount "$PWD/out/anki-as-code-wasm/revision/workflow:/" \
  out/anki-as-code-wasm/revision/anki_as_code.wasm config dump --text /configs
```

## Build scope

```sh
bazel_agent bazel build -c opt --//projects/anki_as_code:wasm=true //projects/anki_as_code //tools/repo_delivery
```

The same invocation produced a WebAssembly CLI and an ELF host tool. With both
`//projects/anki_grpc:cgo=true` and `//projects/anki_grpc:wasm=true`, configured
dependency queries found neither adapter under the selected public WASI target
and both adapters under the sibling `cmd/anki_as_code:configured` target. This
checks platform scope and server selection independently; no invocation-wide
platform or engine overrides are needed.

## Resolved build defects

- `wasi-schema-generator-platform`: the previous `run_binary.srcs` analyzed
  Anki's Rust descriptor generator for WASI, without a matching Rust toolchain.
  The schema rule now selects execution configuration for that input.
- `grpc-wasi-build-metadata`: Gazelle omits `wasip1` from its platform catalog,
  dropping the existing `grpclog` import from the generated dependency selection.
  The maintained post-Gazelle patch repairs that selection without changing pins.
- `go-platform-select-restriction`: the pinned rules_go macro rejects explicit
  `goos`/`goarch` attributes together with existing `select()` link options.
  The supported `go_cross_binary` wrapper retains those link options and scopes
  its platform transition to the wrapped server-free command.

`bazel_agent bazel mod deps --lockfile_mode=update` succeeded without lockfile drift.

## Native compatibility

```sh
bazel_agent bazel test --//projects/anki_as_code:wasm=true //projects/anki_as_code/cmd/anki_as_code:cli_workflow_test //tools/bazelrc:agent_environment_test
```

Both native test targets passed with the project WASM flag enabled, including real native worker, external-only executable,
configuration, collection, and API workflows. The native suite retains its
own undeclared output artifacts through Bazel.

The first attempt encountered `cached-crate-absolute-symlinks`: the cached
`percent-encoding-iri` repository pointed into a removed output base. A bounded
`bazel_agent bazel fetch --force --repo=@@rules_rust++crate+com_github_ankitects_anki_crates__percent-encoding-iri-2.2.0`
refreshed that exact pinned repository; the subsequent suite passed. No cache
files were edited manually and no persistent host settings changed.

## Repository checks

`bazel_agent bazel run //:format` completed, `//:repo_quality_test` passed all
26 checks, and `bazel_agent bazel build --config=lint` passed for the affected
project, CLI, schema-generation, Bazelrc, Go configuration, and gRPC patch packages.
Delivery reruns the mandatory checks against the committed candidate and records
their results in the ignored preparation receipt and validation logs.

## Session reproducibility

Build investigation used the repository's owner docs, OpenSpec wrapper,
Bazel runner, and live Cordis read handlers. One metadata-regeneration attempt
was ineffective because Gazelle's own platform catalog lacks WASI; the final
post-generation patch addresses the observed missing dependency directly.
Early guesses at Bazel external paths added avoidable reads. Follow the
workspace's `bazel-out` symlink to its output base before inspecting downloaded
sources, and check a batched read's paths first because one absent file rejects
the whole Cordis read. These are task-local observations, without policy changes.
