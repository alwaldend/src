## Context

See proposal.md for motivation. The repository uses AL-injected Vault state,
Yandex provider 0.203.0, and Bazel-packaged Ansible. The requested V100 has
compute capability 7.0 and requires a compatible CUDA 12 runtime.

## Goals / Non-Goals

**Goals:** Prepare an encrypted GPU host and an unprivileged local coding account.
**Non-Goals:** Deploy, authenticate cloud services, install on the current
host, add ComfyUI, or promise model quality or live GPU availability.

## Decisions

Verdict: use root-owned Codex requirements to pin the local provider and
disable hosted capabilities. The user's follow-up explicitly chooses normal
runtime networking and Codex's own command permissions. The T3/mTLS extension
owns browser ingress and Auto/Full defaults; see `../add-t3code-mtls/`.

Use a required immutable Ubuntu 24.04 GPU image ID with a V100-compatible
proprietary driver already installed. Avoid a moving image family and avoid
changing kernel drivers implicitly. Use checksum-pinned official Codex and
llama.cpp archives, with CUDA 12 runtime libraries separately verified.
The user explicitly approved the bundled CUDA 12.8 runtime during this task.
The model is an operator-supplied GGUF with a required checksum, so model
selection and redistribution are not hidden provisioning dependencies.

Use a single 200 GiB encrypted boot SSD for OS, model, and workspaces.
Protect the key and disk from accidental destruction. Restrict SSH to explicit
administrator CIDRs and keep inference on loopback. The T3 extension adds
HTTP validation/redirection and HTTPS with client authentication. Do not
attach a cloud service account to the workload VM.

Keep component-specific Ansible tasks under this project. No current host
configuration changes are needed; host_bot remains the owner of its host.
Add the component AppRole in Vault's owning workspace so the existing
Yandex folder automation can discover it; preserve that owner's change record.

## Risks / Trade-offs

- GPU quota, capacity, image compatibility, and tool-call quality remain
  unverified until a separately authorized deployment and live acceptance.
- Root can change restrictions. The shared T3/Codex coding account has no sudo privileges.
- Full access removes Codex's command sandbox when explicitly selected.
- Local inference is on a Yandex VM: the cloud provider still operates the
  hardware. KMS encryption protects storage at rest, not running memory.

## Migration Plan

There is no existing deployment to migrate. Later, bootstrap the component
AppRole/folder, select an immutable compatible image, review a Terraform plan,
provision only with authorization, and configure via Ansible. Live acceptance
must cover GPU loading, a fresh Codex response and tool round trip, local
Auto review, Full access selection, mTLS, local API access, and a reboot. Keep disk and KMS key
on rollback; destroying encrypted data requires an explicit separate decision.
