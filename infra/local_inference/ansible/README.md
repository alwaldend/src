---
title: Local inference configuration
description: Local llama.cpp, Codex, and T3 Code behind Traefik mTLS
---

This stage configures the dedicated Ubuntu 24.04 GPU VM. The image must
already supply a proprietary V100-compatible NVIDIA driver at least 570.26.
The approved CUDA 12.8 runtime is installed from llama.cpp's checksum-pinned
archive. Runtime versions and checksums live in `group_vars/local_inference.yaml`.

Set `LOCAL_INFERENCE_HOST` from Terraform's `ansible_host` output. Supply
`LOCAL_INFERENCE_MODEL_SRC` as an absolute controller-side GGUF path and
`LOCAL_INFERENCE_MODEL_SHA256` as its SHA-256 digest. No model is downloaded
implicitly. The default context is one 32k conversation; reduce
`local_inference_context_size` when the model's memory requirements demand it.

The prepared Vault server PKI role, component AppRole, DNS permissions, and
public delegation must be deployed separately before certificate enrollment.
No production token or certificate is created by source validation.

T3 Code runs as its ordinary unprivileged `t3code` account on
`127.0.0.1:3773`. The complete checksum-pinned official standalone bundle
is installed before startup. Traefik exposes
`https://t3code.local-inference.alwaldend.com`, requires a certificate from
the repository's client CA, and obtains its server certificate through
Vault ACME. Port 80 serves ACME validation and HTTPS redirection. T3's native
pairing/session authentication remains enabled. The standalone CLI is
`/opt/t3code/releases/<t3code_version>/t3`; use its full path as the service
account to create pairing grants after deployment. Treat
pairing codes and startup URLs as secrets. This release does not consume the
legacy `T3CODE_AUTH_TOKEN` variable.

Initial T3 settings select the existing Codex executable, its dedicated
service-account home, the `local-coder` model, and Auto mode. They seed only
an absent settings file, preserving later UI changes. Auto uses Codex's
selected local provider for approval review. Full access remains selectable
and removes the command sandbox when chosen. The actual local model must
support tool calls and approval review; fitting in VRAM does not establish
agent reliability.

Root-owned `/etc/codex/requirements.toml` pins Codex to
the Responses bridge at `127.0.0.1:8080/v1` and disables OpenAI-hosted search, apps, plugins, image
generation, and remote-control capabilities. The MCP allowlist permits only
`t3-code` at `http://127.0.0.1:3773/mcp`; its port follows `t3code_port`.
T3 supplies per-session authentication. A changed URL or another MCP name
remains blocked. Local shell/edit tools, hooks, and subagents remain available.
Accounts retain normal networking, including network access from commands
allowed by Codex. This configuration does not add an outbound
firewall or custom restrictions to T3's other normal functions.

The packaged first-party bridge runs as `llama` in `responses-bridge.service`.
It translates namespaced function schemas, call history, and JSON/SSE replies
for the pinned llama.cpp server on `127.0.0.1:8081`. Unsupported tool schemas
fail explicitly instead of disappearing from the model's tools. The bridge
has a fixed loopback upstream and no authentication state. Readiness checks
reach llama.cpp through the bridge. Archive comparison repairs incomplete
runtime extraction on subsequent Ansible runs.

T3 and direct Codex sessions share the same unprivileged coding account,
home, and workspace. An administrative SSH session can run Codex directly:

```sh
sudo -iu t3code codex -C /opt/t3code/work
```

Build the package without connecting to a host:

```sh
bazel_agent bazel build //infra/local_inference/ansible:ansible_bin
```

The normal `//infra/local_inference/ansible` target deploys and requires
separate authorization. Future live acceptance must verify GPU model loading,
a fresh Codex response and tool follow-up, local Auto review, Full selection,
T3 MCP discovery, mTLS rejection/acceptance, application authentication, and
reboot persistence.
Keep non-secret verification artifacts under ignored `out/` and record
model/runtime revisions with the results.
