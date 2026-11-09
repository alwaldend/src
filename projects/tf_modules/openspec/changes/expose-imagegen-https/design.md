## Context

Image generation needs SSH-only and password-authenticated HTTPS access with a 32 GiB VRAM floor.

## Goals / Non-Goals

Retain ssh_only by default; exposed publishes only 7860/tcp through a dedicated TLS and Basic-auth Caddy configuration. The application stays on loopback. Reject unsafe template options and weak credentials; disable outbound tunnels and synchronization while retaining loopback-only portal, Jupyter, and TensorBoard services.

## Decisions

Reuse the existing provider, pinned Terraform providers, and AL/Vault configuration flow.
The personal consumer declares a direct published ComfyUI/CUDA 12.9 tag in Terraform. The user explicitly requested removal of the image-discovery executable and selection by tag rather than digest. The reusable module takes the image as an
explicit input and supports declarative public model downloads. No stock template or unrelated provisioning configuration is inherited
by default. Keep the portal, Jupyter, and TensorBoard local; disable outbound
tunnels/synchronization and use a dedicated application-only public proxy.
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

## Runtime separation

Proceed with three modules: vastai_instance_imagegen owns the rental, certified
host constraints, sensitive login, native download manifest and dedicated HTTPS
proxy. vastai_instance_imagegen_forge owns Forge paths and native selection;
vastai_instance_imagegen_comfyui owns ComfyUI paths, loopback arguments and
a local-only API-wrapper service. Runtime-specific settings must not remain in the
base. Both wrappers forward the existing rental options and connection outputs.
The image remains an explicit caller-selected tag, without registry discovery.

This reuses Vast's startup/provisioner rather than installing ComfyUI into a
Forge image. The ComfyUI container is a new dependency requiring specific user
approval before declaration. Disable custom nodes, partner APIs and browser
internet access using the native offline flag; do not enable Manager. Preserve
the existing common authentication boundary. Offline acceptance cannot establish
GPU inference or actual image boot; those require a separately authorized rental.

## Chroma component selection

Add the native models/text_encoder destination and an optional structured model
selection. Encode native Forge settings as base64 JSON and pass the existing
--ui-settings-file flag. Set global and preset-specific checkpoint/module keys
because Forge restores preset-specific selections on page load. Selected paths
are derived from the validated manifest. Service exposure and credentials are
unchanged. The existing Terraform CLI fixture checks the staged settings.

The user explicitly approved vastai/comfy:v0.44.0-cuda-12.9-py312 on 2026-11-10.
Use its native ComfyUI service and retain the additional api-wrapper service on loopback.

The approved image declares diagnostic port 3000 in Docker metadata. Keep its
optional instance test runner explicitly disabled; request no publish-all flag.
