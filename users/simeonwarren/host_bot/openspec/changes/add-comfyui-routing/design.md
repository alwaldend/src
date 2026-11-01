## Context

T3 Code already has owner-local DNS declarations and a Vault-backed HTTPS
router. The existing ComfyUI launcher binds to `127.0.0.1:8188` by default.

## Goals / Non-Goals

Provide the same routing pattern for ComfyUI. Service installation, startup,
and live deployment are outside this source configuration change.

## Decisions

Reuse the existing DNS view pattern and local Traefik template, with
`comfyui_domain` and `comfyui_port` inventory variables. Router-level Vault
certificate resolution follows T3 Code; no static TLS default-domain change
is needed. A separate proxy or direct external listener would duplicate the
existing owner or bypass ingress authentication.

## Risks / Trade-offs

- A stopped ComfyUI server returns a gateway error; document the existing launcher.
- Source checks cannot establish live DNS or certificates; deployment requires
  separately authorized DNS reconciliation and Ansible execution.

## Architecture documentation

Keep the maintained source in `docs/comfyui-architecture.mmd` under Host Bot,
which owns the network topology. Follow the existing T3 Code diagram's
left-to-right layout, public/local actors, ingress, WireGuard path, local
Traefik, and Vault certificate flow. Show the HTTP API and generation thread
inside one ComfyUI process, with local model weights and PNG outputs on the
host. The process uses the Z-Image target's transient systemd scope; do not
imply a separate execution worker or remote inference service.

Use `mermaid_site_svg` for the normal `rendered` and `rendered_comfyui`
targets, matching Download, Flux, 3x-ui, Forgejo, and X Profile Cleaner.
Publish their light/dark variants through Host Bot's `docs_assets`; retain
the former `*_site` labels as aliases rather than duplicate render actions.
Link the diagrams from the owner README; keep generated SVGs out of source
control. The existing blog publication consumes an explicitly named
`rendered_historical` target with the original theme and one SVG output.

Before implementation, check these failure cases: incorrect authentication or
transport labels, a false process boundary, omitted documentation assets,
clipped text, overlapping container titles, and edges crossing title plates.
Build the render and documentation targets, inspect the rendered variants,
and compare the represented topology with its existing configuration owners.
This is a documentation addition with no live operation or behavior change.

For the default-style correction, also verify both canonical targets produce
native light/dark pairs, the aliases resolve to those same outputs, and the
blog's existing generated-image freshness test passes without changing its
published asset. The shared theme and Mermaid sources remain unchanged.
