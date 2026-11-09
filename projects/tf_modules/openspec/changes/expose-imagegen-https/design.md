## Context

Image generation needs SSH-only and password-authenticated HTTPS access with a 32 GiB VRAM floor.

## Goals / Non-Goals

Retain ssh_only by default; exposed publishes only 7860/tcp through a dedicated TLS and Basic-auth Caddy configuration. Forge stays on loopback. Reject unsafe template options and weak credentials; disable outbound tunnels and synchronization while retaining loopback-only portal, Jupyter, and TensorBoard services.

## Decisions

Reuse the existing provider, pinned Terraform providers, and AL/Vault configuration flow.
The personal consumer resolves the newest matching Forge Neo/CUDA 12.9 tag
from Docker Hub with a bounded read-only external-data adapter. The user explicitly
requested unpinned image selection. The reusable module takes the image as an
explicit input and supports declarative public model downloads. No stock template or unrelated provisioning configuration is inherited
by default. Keep the portal, Jupyter, and TensorBoard local; disable outbound
tunnels/synchronization and use a dedicated Forge-only public proxy.
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
