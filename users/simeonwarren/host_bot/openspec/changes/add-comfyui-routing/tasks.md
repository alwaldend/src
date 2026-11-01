## 1. Host routing

- [x] 1.1 Add ComfyUI variables, local Traefik routing, DNS declarations, and documentation; verify rendered TOML and DNS views.

## 2. Validation

- [x] 2.1 Build the Ansible package, run DNS and repository checks, and verify the rendered ingress-to-loopback route preserves T3 Code.

## Evidence

Validated on 2026-11-01: both owner Ansible binaries built; DNS lint and all
four `//infra/dns:config_test` tests passed; strict OpenSpec validation,
`//:repo_quality_test`, and affected-package semantic lint passed. Rendering
with the packaged Jinja dependency and parsing with TOML verified the hostname,
HTTPS backend, backend SNI, Vault resolver, loopback port, DNS views, and
unchanged existing routes and DNS. Reproduction input and rendered outputs
are in ignored `out/comfyui-ingress-traefik/`; exact-candidate delivery checks
are recorded in its preparation and validation receipts. No live deployment
or service-health check was performed.

## 3. Authorized deployment

- [x] 3.1 Deploy the ComfyUI Traefik route through the owning Ansible playbook with `--tags traefik`; verify the active route.
- [x] 3.2 Review an additions-only DNS plan, apply the exact saved plan, and verify public and site-local records.

On 2026-11-01, the authorized deployment applied an exact saved DNS plan:
four ComfyUI additions, zero updates and zero deletions. Both DNS views
resolved to their declared targets. The owning Ansible playbook ran with
`--tags traefik --limit dc1.host-bot.simeonwarren.users.alwaldend.com
--connection local` after verifying that inventory resolves to this host.
Its recap was 12 successful tasks, 5 changes, and zero failures or
unreachable hosts. The local connection preserved SSH trust settings.
The public ingress deployment is recorded by its owning change. Reviewed
plan summaries and deployment logs remain under ignored
`out/comfyui-routing-deploy/`; secret-bearing saved plan artifacts were
removed after applying.

The moved main target subsequently returned 200 through the deployed
local HTTPS proxy, including its required frontend assets. The owned
verification instance was stopped after the check.

## 4. ComfyUI architecture diagram

- [x] 4.1 Add the maintained diagram, matching T3 Code's topology conventions, and integrate its rendered variants into Host Bot documentation.
- [x] 4.2 Build and inspect the rendered variants and documentation assets, verify the represented routing/process boundaries, and record delivery evidence.

The original-theme render and both native site variants built on 2026-11-01.
Visual inspection found readable labels, contained titles, and no edges crossing
title plates. The documentation asset target includes both ComfyUI variants
alongside the existing T3 Code variants. Source review confirmed public mTLS,
WireGuard transport, local Traefik/Vault forwarding, and the shared Python
process containing the API and generation thread. Reproduction outputs and
review evidence remain under ignored `out/comfyui-diagram/`; exact-candidate
quality, build, lint, and publication results belong to its delivery receipts.
No deployment, application restart, or inference run was performed.

## 5. Current Mermaid render defaults

- [x] 5.1 Make the normal ComfyUI and T3 Code targets use the current native site macro, retain compatibility aliases, and route the historical blog consumer explicitly.
- [x] 5.2 Build and inspect both native diagram pairs, verify alias output identity and historical asset freshness, and record verification evidence for delivery.

Both normal render targets now produce native light/dark pairs matching the
current shared site style. The former site-specific targets resolve to the
same configured render targets and output files. All four default SVGs are
byte-identical to the prior site variants. Visual comparison with Flux and
Download confirmed matching typography, shapes, and captions. All six blog
image freshness tests pass without publication-image changes. Source, alias,
image, and delivery evidence remain under ignored
`out/comfyui-diagram-style/`. No live operation occurred.
