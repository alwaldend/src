---
title: Local Responses protocol acceptance
---

Run the command-level fixture without a GPU, model, credentials, or deployment:

```sh
bazel_agent bazel test //infra/local_inference/test/protocol:protocol_test --test_output=all
```

The fixture starts the built bridge on an ephemeral loopback port and uses a
function-only inference backend. It checks duplicate local names in separate
namespaces, plain functions, JSON and SSE replies, replayed calls, forced tool
selection, and rejection of unsupported schemas. Bazel retains `test.log` as
the repeatable acceptance artifact.

The same target accepts `--test_arg=-codex-binary=/absolute/path/to/codex`.
Use the pinned Codex 0.166.0 executable in an isolated environment with no
system Codex requirements. The optional case starts a local MCP fixture,
asks Codex to use its tool, checks the model request contains the flattened
schema, returns a deterministic tool call, and verifies that Codex executes
the tool and submits its result before completing the follow-up response.
This tests the protocol and execution path, not model quality.

On a host with shared managed Codex requirements, build first, then run the
compiled test in an isolated network/mount namespace. For example, using the
host's existing `bwrap` installation:

```sh
mkdir -p out/local_inference/compatibility/tmp
bwrap --unshare-net --die-with-parent --ro-bind / / --dev-bind /dev /dev \
  --proc /proc --tmpfs /etc/codex \
  --bind "$PWD/out/local_inference/compatibility" "$PWD/out/local_inference/compatibility" \
  --setenv TMPDIR "$PWD/out/local_inference/compatibility/tmp" \
  --setenv RUNFILES_MANIFEST_FILE "$PWD/bazel-bin/infra/local_inference/test/protocol/protocol_test_/protocol_test.runfiles_manifest" \
  "$PWD/bazel-bin/infra/local_inference/test/protocol/protocol_test_/protocol_test" \
  -bridge=_main/infra/local_inference/cmd/responses_bridge/responses_bridge_/responses_bridge \
  -codex-binary=/absolute/path/to/codex -test.v -test.timeout=60s
```

No namespace tool conversion is currently available in the checked llama.cpp
b11960 or b12032 releases. The bridge can be removed when a pinned upstream
release passes these acceptance cases directly.
