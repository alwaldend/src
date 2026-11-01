---
title: Local inference
description: Yandex Cloud V100 host for local Codex inference and tools
---

This project owns the Terraform provisioning and Ansible configuration for
a Yandex Cloud V100 host running llama.cpp, Codex CLI, and T3 Code. These definitions
do not establish deployed state. Deployment is a separate operator action.

## Layout

- `tf_setup/` declares one `gpu-standard-v1` VM: 1 V100 with 32 GB VRAM,
  8 vCPUs, and 96 GiB RAM. Its 200 GiB SSD is encrypted by a rotating
  AES-256 KMS key. SSH requires explicit administrator CIDRs.
- `ansible/` installs checksum-pinned Codex, llama.cpp, and the approved
  CUDA 12.8 runtime; stages a checksum-verified model; and configures
  local Codex inference and tools. It also runs T3 Code behind Traefik mTLS.
  Runtime pins live in its group variables.
- `al.lua` wires the component's Vault AppRole, Yandex credentials,
  Terraform HTTP state backend, DNS credentials, and the
  signed Ansible SSH key.
- `openspec/` records the component contract and implementation evidence.

There is no service-side `tf/` stage: llama.cpp and Codex configuration lives
in host files managed by Ansible, rather than a separately provisioned API.

## Prerequisites for a later deployment

Apply the prepared component AppRole declaration through `infra/vault/tf`,
then synchronize `infra/yandex_cloud/org1/tf`. The existing folder automation
discovers `src_infra_local_inference` through the AppRole inventory. These
are separate live operations requiring explicit authorization. Backend state
and locks are created through the existing Vault HTTP backend flow.

Supply the required immutable Ubuntu GPU image ID and administrator CIDRs
described in [provisioning](tf_setup/README.md). GPU quota and capacity must
be available in the selected zone. Ansible requires the VM address, a
controller-side GGUF, and its SHA-256 digest, as described in
[host configuration](ansible/README.md). No model is selected or downloaded
implicitly, and no cloud inference credential is needed.

## Inference and access

Codex connects to the local Responses bridge at `127.0.0.1:8080/v1`.
The bridge translates function namespaces for llama.cpp on port 8081 and
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
T3 and llama.cpp listen only on loopback. The canonical DNS declarations
are in `dnsconfig.json`.

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
`//infra/local_inference/test/protocol:protocol_test`. It starts the actual bridge
against a deterministic function-only backend and checks JSON/streaming routing,
replayed calls, and explicit rejection of unsupported tools. See
[protocol acceptance](test/protocol/README.md) for the pinned Codex round trip.
