## Ownership and preservation

The baseline is `infra/arch/arch.drawio` at 1729c90d4 (14 pages), observed 2026-10-11. The 3x-ui replacement is the source from PR #141 at 1e42fa210. This branch incorporates that diagram without modifying its original worktree or published branch.

Each source belongs under its owner's `docs/` directory. Project READMEs consume generated light/dark SVG siblings through docs_filegroup, with no checked-in render projections. Archive labels describe the source document's historical classification, not deployed state. The legacy www diagram remains explicitly historical because the website has a newer deployment contract.

The archived pages without component projects share `infra/archive`, as explicitly requested. Proxmox and Harbor retain existing owners. The removed central project's obsolete requirements are removed with it; the repository change record preserves the migration and source revision.

## Acceptance

Account for all 14 pages, preserve named topology and protocol labels, render every source, inspect representative outputs and all nested containers, package documentation successfully, remove infra/arch and active stale references, and pass scoped delivery checks.

## Rendered layout decisions

Container titles use short captions; the 3x-ui overview and control-plane view split hostnames onto a second line. Service nodes may also use line breaks. Vault's nested containers reserve 12 pixels below captions. Ingress uses a short regional VM caption, with the ingress hostname on Traefik. Harvester keeps each host's storage and RAM as node attributes rather than separate nested inventory boxes, preserving the bidirectional host link without routing through a title plate.

## 3x-ui architecture views

The 3x-ui documentation publishes one source and render with three separate sections.
The overview
shows the simple control plane in Proxmox, njalla1 and work1 in njal.la, yc1 in Yandex
Cloud, and the external Vault, Let's Encrypt, and Mullvad services. A shared data-plane
boundary contains both provider groups; traffic, management, ACME, and WireGuard
links connect to that boundary, while its node inventory remains unconnected.
Provider captions use their proper names without the word "datacenter". The control-plane view expands its
Traefik and panel and uses one simple data-plane node. The data-plane view expands
one node's shared Traefik, 3x-ui, and Xray implementation and uses one simple
control-plane node. The sections share one canvas but no visible architecture
edges. Invisible links put the overview above the two detail sections, while
each section lays out its connections from left to right. Repository documentation
consumes build-time light/dark SVGs only. WebPs are user previews under ignored
task scratch, using the shared 5:2 raster default and the site's dark palette.

The overview orders user traffic above administration. Its Control Plane group
contains Proxmox (with a two-line hypervisor/hostname caption) and the Master
node. Data-node labels show "Node" above full hostnames from the checked-in
Ansible inventory and DNS declarations. Source ordering of groups, service
nodes, and edges keeps the user/data route above the administration/master route.

Outgoing overview links from Control Plane attach to its group boundary: management
to Data plane and ACME to Vault. Admin HTTPS also attaches to Control Plane. No external edge reaches Master,
so Mermaid can lay out the nested control-plane contents independently.

The detailed control-plane view labels its nested VM "3x-ui VM". Both VM and
Proxmox container captions put their full hostnames on a second line.
