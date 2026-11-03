## Context

See proposal.md for motivation. The host deployment uses Vault-backed SSH and Bazel-packaged Ansible;
VM provisioning is delivered separately. The requested V100 has
compute capability 7.0 and requires a compatible CUDA 12 runtime.

## Goals / Non-Goals

**Goals:** Configure an existing GPU host and an unprivileged local coding account.
**Non-Goals:** Deploy, authenticate cloud services, install on the current
host, add ComfyUI, or promise model quality or live GPU availability.

## Decisions

Verdict: use root-owned Codex requirements to pin the local provider and
disable hosted capabilities. The user's follow-up explicitly chooses normal
runtime networking and Codex's own command permissions. The T3/mTLS extension
owns browser ingress and Auto/Full defaults; see `../add-t3code-mtls/`.

The separately delivered Terraform stage owns the reviewed immutable GPU image,
V100 shape, storage encryption, destruction guards and SSH network inputs.
Use checksum-pinned official Codex and
llama.cpp archives, with CUDA 12 runtime libraries separately verified.
The user explicitly approved the bundled CUDA 12.8 runtime during this task.
The model is an operator-supplied GGUF with a required checksum, so model
selection and redistribution are not hidden provisioning dependencies.

The host runtime uses the provisioned encrypted disk. Inference stays on
loopback, and the T3 extension owns HTTP redirection and HTTPS with client
authentication. Cloud resources and their credentials are outside this change.

Keep component-specific Ansible tasks under this project. No current host
configuration changes are needed; host_bot remains the owner of its host.
The existing component AppRole supplies signed SSH access; provisioning
authentication is owned by the separate Terraform stage.

## Risks / Trade-offs

- GPU quota, capacity, image compatibility, and tool-call quality remain
  unverified until a separately authorized deployment and live acceptance.
- Root can change restrictions. The shared T3/Codex coding account has no sudo privileges.
- Full access removes Codex's command sandbox when explicitly selected.
- Local inference is on a Yandex VM: the cloud provider still operates the
  hardware. KMS encryption protects storage at rest, not running memory.

## Deployment prerequisites

Provision the GPU host through the separately delivered Terraform stage before
configuring it with Ansible. These remain separately authorized live operations.
Live host acceptance must cover GPU loading, a fresh Codex response and tool
round trip, local Auto review, Full access selection, mTLS, local API access
and reboot behavior. Model files and their reviewed checksum are operator inputs.
