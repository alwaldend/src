## Context

Image generation requires per-instance environment overrides and observed mapped ports.

## Goals / Non-Goals

Pass optional Docker options through the generated create protocol. Read typed port bindings without retaining API-returned credentials; preserve omitted-field compatibility.

## Decisions

Reuse the existing provider, pinned Terraform providers, and AL/Vault configuration flow.
Secrets remain sensitive in Terraform and reside in the Vault-backed backend; they are never public outputs.
The caller owns template trust and deployment authority. The 32 GB requirement is interpreted as 32 GiB VRAM.

## Risks / Trade-offs

Offline fixtures verify request and lifecycle behavior, not actual image startup or marketplace availability.
HTTPS uses the Vast-provided instance certificate; clients must verify its trust before submitting credentials.

## Migration Plan

Default SSH behavior is retained; changed startup settings intentionally replace a rental.
Deploy identity prerequisites separately before using the new rental root.

## Rental rejection diagnostics

The authorized ComfyUI deployment returned HTTP 400 twice. Decode bounded rental
rejection metadata only for known categories and the documented numeric startup
limit; never render free-text messages, request content or credentials. Preserve
generic errors for unknown response shapes. This extends the existing constant
team-name rejection diagnostic rather than enabling raw response logging.

The bounded diagnostic identified the server's 16,384-byte limit, and the
corrected request created the authorized rental successfully.
