## Context

A personal image generation rental needs generated credentials retained in Vault before paid instance creation.

## Goals / Non-Goals

Add imagegen/tf_setup with random_password, Vault KV storage, and the reusable exposed-mode module. Reuse the existing user_simeonwarren AppRole, regular user secret namespace, and independent Vault-backed state. Use the newest published Vast Forge Neo/CUDA 12.9 image tag without a template by default; provide optional declarative public model downloads; allow an optional reviewed environment-only template; retain the USD 1/hour cap. Do not execute a live rental as implementation validation.

## Decisions

Reuse the existing provider, pinned Terraform providers, and AL/Vault configuration flow.
Secrets remain sensitive in Terraform and reside in the Vault-backed backend; they are never public outputs.
The caller owns template trust and deployment authority. The 32 GB requirement is interpreted as 32 GiB VRAM.

## Risks / Trade-offs

Offline fixtures verify request and lifecycle behavior, not actual image startup or marketplace availability.
HTTPS uses the Vast-provided instance certificate; clients must verify its trust before submitting credentials.

## Migration Plan

Default SSH behavior is retained; changed startup settings intentionally replace a rental.
Create a dedicated alwaldend API token through the existing alwaldend root and store it in the regular user namespace before renting.

## Open Questions

None for implementation. A live acceptance run must inspect the then-selected image and actual service bindings.

## Default checkpoint

The user explicitly approved RunDiffusion/Juggernaut-XL-v9 as the rental default. Set model_downloads to the publisher's standalone safetensors checkpoint at revision cf419233522daa0b9ea36c3aff98fa2cab1fb0fb. This keeps the reusable template's model list generic and empty. The native provisioner downloads the checkpoint before Forge starts; no model-serving package or authenticated downloader is added. The user may replace the list or supply an empty list.

The publisher describes this as its most popular SDXL checkpoint; this is not an independently verified global ranking of uncensored models. The 7,105,348,188-byte file fits the rental's 100 GB disk and its SDXL architecture fits Forge and the 32 GiB VRAM floor. License terms still apply. Failed public downloads or upstream availability can prevent first boot; a fixed source revision prevents silently selecting another release. Verification covers public URL availability and the existing manifest/lifecycle fixtures, not a live GPU generation.
