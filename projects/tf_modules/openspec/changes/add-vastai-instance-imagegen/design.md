## Context

Vast.ai's image-generation guide recommends the AI-Dock Forge template.
Template environment options are merged by the rental API, including published
ports; overriding launch mode does not remove them.

## Goals / Non-Goals

Provide reusable source and offline evidence. Credentials, state backends,
offer discovery, SSH keys, template creation, and deployment belong to callers.

## Decisions

Proceed with a required template hash and bounded automatic datacenter offer
selection. Optional offer ID narrows the same verified/datacenter search. The module reads
the template and rejects Docker options other than literal environment settings.
The provider retains ordinary image-based rental compatibility and adds a
creation-only template hash. The module supplies the resolved image explicitly.

Override startup with AI-Dock's documented `init.sh` argument mechanism:
disable Cloudflare and auxiliary network services, disable updates, and bind
Forge to loopback on port 17860. Preserve the selected template's provisioning
environment. The forwarded local port defaults to 7860 and binds to 127.0.0.1.

A new pinned application image was considered but requires a separate dependency
decision and duplicates the template selection the caller already makes. Using
an unmodified public template was rejected because its published ports survive
request merging. No arbitrary startup or Docker option escape hatch is exposed.

## Risks / Trade-offs

- AI-Dock startup contract changes → document the required template family and
  leave GPU boot and image generation explicitly unverified until an approved rental.
- Untrusted templates can execute arbitrary provisioning → callers must review
  template content; port validation is not application-code isolation.
- Offer availability and price can change → enforce an hourly ceiling and
  ignore changes to offer ID on existing contracts; a replacement selects a
  currently matching offer.
- Destroy deletes local outputs; stopping retains storage charges → document
  download and cleanup requirements without executing them.

## Validation failure cases

Before implementation, cover merged public-port options, malformed or ambiguous
template responses, quoted environment values, unsupported Docker flags,
template replacement, a pending SSH address, and loopback forwarding after
refresh. The fixture must reject creation unless template hash, SSH mode, and
controlled startup are present. Retain a non-secret lifecycle artifact.

## Shared evaluation inputs

The ephemeral evaluator uses the existing hardening template as a pure input, compiled without providers, plans or state. Extract fixed ComfyUI runtime settings and common offer filters into JSON declarations owned by these modules and consumed by both paths. Keep dynamic hardware/port/whitelist inputs in their callers and preserve personal provisioning. Validate module and provider consumers offline.

The single shared offer declaration caps both inet_down_cost and inet_up_cost
at USD 0.04/GB. Terraform and the evaluator inherit the same filters; retain
verified/datacenter and resource constraints. This is a rate ceiling, not a
total account spending guarantee, and does not create persistent storage.

The operator widens the initial USD 0.015/GB ceiling to USD 0.025/GB for both directions to allow more market capacity while retaining a bounded transfer price.

On 2026-11-12 the operator approves USD 0.04/GB for the corpse evaluation and
then explicitly requests that ceiling permanently. A read-only market query
found no eligible offers at USD 0.025/GB and three verified datacenter offers
below USD 1/hour charging USD 0.0390625/GB in each direction. Change the
existing shared declaration rather than maintaining an evaluator override.
This source change affects future searches; it performs no personal Terraform
apply or replacement and does not change host configuration.

## Native USO projector destination

The discussed USO model uses a native ComfyUI model patch, distinct from a
LoRA. Extend the existing allowed model kinds with model_patch and map it to
models/model_patches in the shared ComfyUI runtime declaration. Preserve public
HTTPS/plain safetensors validation and require every selected runtime to own
its destination. Extend the configured-profile fixture's destination contract
before implementation; do not change Forge's supported directories or rent
capacity for validation.
