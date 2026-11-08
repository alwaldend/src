---
title: Local inference
description: Independent GPU provisioning, DNS, and local coding host configuration
---

This component owns independent GPU provisioning, DNS, and host configuration
for llama.cpp, Codex CLI, and T3 Code. These definitions do not establish
deployed state or service health. Deployment requires separate authorization.

## Layout

- `tf_setup/` owns the GPU VM, encrypted storage, rotating KMS key, network
  and Yandex child DNS zone. Its stage-local AL configuration selects only
  Yandex credentials and the provisioning state backend. Its
  [change record](openspec/changes/extract-local-inference-provisioning/)
  records provisioning source acceptance and allocation outcomes.
- `dns/` owns the public Cloudflare delegation and CNAME with independent
  state and DNS-only credentials. Its [change record](openspec/changes/separate-dns-terraform-root/)
  records verified DNS ownership.
- `ansible/` installs checksum-pinned Codex, llama.cpp, and the approved
  CUDA 12.8 runtime; stages a checksum-verified model; and configures
  local Codex inference and tools. It also runs T3 Code behind Traefik mTLS.
  Bazel owns runtime pins in `third_party/`; the deployment packages the
  verified archives for Ansible. The shared `llama` role owns inference.
- `BUILD.bazel` wires the component's Vault AppRole and signed Ansible SSH key.
- `openspec/` records the component contract and implementation evidence.

There is no service-side `tf/` stage: llama.cpp and Codex configuration lives
in host files managed by Ansible, rather than a separately provisioned API.

## Prerequisites for a later deployment

The [provisioning stage](tf_setup/README.md) supplies the VM address.
The host deployment consumes its VM address; it does not create cloud resources
or load Yandex credentials. Ansible requires the VM address, a
controller-side GGUF, and its SHA-256 digest, as described in
[host configuration](ansible/README.md). No model is selected or downloaded
implicitly, and no cloud inference credential is needed. The documented
provisioning uses Fedora, while the current Ansible stage requires Ubuntu
24.04; host configuration needs adaptation before deployment.

## Inference and access

Codex connects to the local Responses bridge at `127.0.0.1:8080/v1`.
The [project-owned bridge](../../projects/llama_responses_bridge/README.md)
translates function namespaces for llama.cpp on port 8081 and
restores their identities in replies. It performs no inference or tool execution.
Root-owned requirements constrain
its provider and disable OpenAI-hosted search, apps, plugins, image
generation, and remote control. T3's authenticated loopback MCP is allowed;
unapproved MCP identities remain blocked. Local shell/edit tools remain
enabled. Accounts retain normal networking; Codex's command permissions
govern their tool use.

T3 Code and direct Codex sessions share an ordinary coding account. T3
initially selects `local-coder` with Auto approval review through local inference. Full access
remains selectable. Seed settings preserve subsequent UI changes. Traefik
exposes `https://t3code.local-inference.alwaldend.com` using Vault ACME and
requires a trusted client certificate plus native T3 pairing/session authentication.
T3 and llama.cpp listen only on loopback. The canonical `dnsconfig.json` is
consumed by Vault's exact-hostname PKI role and the independent
[DNS Terraform root](dns/README.md). Provisioning retains only the Yandex
child zone and VM address.

Local means inference runs on the VM. Browser interactions reach that VM, and Yandex
Cloud operates its hardware. Encryption protects persistent storage at rest.

## Source validation

Formatting, packaged Ansible syntax, rendered configuration, and
backend-disabled Terraform validation can run without deploying. Never use
the normal Ansible target or a Terraform plan as an offline check.
The maintained change records the actual source-validation evidence and
the GPU, fresh-session tool-call, approval, mTLS, and reboot checks required
after separately authorized deployment.

The checked-in protocol acceptance target is
`//projects/llama_responses_bridge/test/protocol:protocol_test`. It starts the actual bridge
against a deterministic function-only backend and checks JSON/streaming routing,
replayed calls, and explicit rejection of unsupported tools. See
[protocol acceptance](../../projects/llama_responses_bridge/test/protocol/README.md) for the pinned Codex round trip.
