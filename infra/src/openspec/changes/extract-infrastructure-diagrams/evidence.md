# Conversion evidence

Observed 2026-10-11 against the task candidate. The source revision is
`1729c90d42167131cccf2fd3fd93e43b77da35a8`; 3x-ui is based on
`1e42fa210913917c86f89d8ccda132a4d6392043` from PR #141.

| Original page     | Authoritative replacement                    |
| ----------------- | -------------------------------------------- |
| DC1               | `infra/mikrotik/docs/dc1.mmd`                |
| Vault             | `infra/vault/docs/architecture.mmd`          |
| Flux              | `infra/flux/docs/architecture.mmd`           |
| Forgejo           | `infra/forgejo/docs/architecture.mmd`        |
| www               | `projects/alwaldend.com/docs/legacy-www.mmd` |
| Threexui          | `infra/threexui/docs/architecture.mmd`       |
| DNS               | `infra/dns/docs/architecture.mmd`            |
| Ingress           | `infra/ingress/docs/architecture.mmd`        |
| Truenas           | `infra/nas/docs/architecture.mmd`            |
| Archive/Harvester | `infra/archive/docs/harvester.mmd`           |
| Archive/Opencode  | `infra/archive/docs/opencode.mmd`            |
| Archive/Hermes    | `infra/archive/docs/hermes.mmd`              |
| Archive/Proxmox   | `infra/pve/docs/archive-proxmox.mmd`         |
| Archive/Harbor    | `infra/harbor/docs/archive-harbor.mmd`       |

The initial 14 site Mermaid targets rendered light/dark SVG pairs. The initial
documentation site build packaged all 28 SVGs at the image paths referenced by the READMEs.
Original XML node labels, protocols, containment, and arrowheads were inspected;
coordinate-only www, DNS, and Hermes edges were associated with their visible
endpoints. Harvester hardware attributes are consolidated into their host nodes;
Ingress's hostname moved from the regional VM caption to its Traefik service.
Both preserve the original topology and information. The www architecture
remains explicitly historical, and all five archived pages retain that status.

The rendered gallery was inspected through the collaborative browser. Across
14 SVGs and 31 containers, measured title rectangles showed no title overlap or
horizontal escape. Sampling each edge at 201 points found no edge passing
through a title rectangle. This bounded geometric check does not prove every
possible visual property. A representative Flux raster was also rendered with
the pinned native renderer and inspected. Sources, gallery, measurements, and
logs are task scratch under `out/convert-infra-diagrams`; delivery receipts own
final candidate and publication evidence.

## Session ergonomics

- `DIAGRAM_DISCOVERY_OVERREAD`: broad combined policy/owner reads and a full
  worktree inventory caused truncated output. Subsequent discovery used bounded
  owning sections and exact paths. Future continuation can use this page map.
- `DIAGRAM_PREVIEW_SNAPSHOT_FAILURE`: collaborative browser DOM inspection
  worked, but snapshots failed. The repository's pinned raster renderer supplied
  representative visual evidence; no alternate browser automation was used.
- `DIAGRAM_SCRATCH_PACKAGE_IGNORED`: Bazel intentionally ignores `out/`, so a
  scratch render package could not load. The existing renderer executable was
  reused with workspace-absolute inputs and task-owned output. Its launcher and
  declared external artifact roots resolved runfiles and font dependencies.

No shared policies or runtime configuration changed because of these findings.

## Validation

Candidate `6463bd2ea03f5b2770685240aa7c76708ce39fcb` passed all 28 selected
repository quality, repository OpenSpec, and diagnostics skill offline tests;
`//projects/alwaldend.com:site` built successfully. Semantic lint passed across
all 30 affected package selections after refreshing the configured Python
3.13.9 repositories through Bazel. No dependency pins or host configuration
changed. Final delivery repeats these gates against the prepared candidate.

- `DIAGRAM_REPOSITORY_CONTENTS_SYMLINK`: a cached interpreter symlink referred
  to another output base that no longer existed. A bounded, forced fetch of the
  two existing pinned Python repositories without repository-contents reuse
  rebuilt correct symlinks and allowed semantic lint to pass.
- `DIAGRAM_GIT_INDEX_CONTENTION`: polling Git status during validation briefly
  contended with the delivery tool's index-tree verification. The lock cleared;
  subsequent monitoring read receipt/log state without concurrent Git commands.
- `DIAGRAM_DELETED_DIRECTORY_STAGE`: explicit preparation of an already staged,
  deleted directory did not match a Git pathspec. Preparation used the inspected
  fully task-owned index, preserving the complete deletion scope.

## 3x-ui architecture views

The original 3x-ui page is represented by one diagram with three separate sections: overview,
detailed control plane with one simple data-plane node, and detailed data plane
with one simple control-plane node. The overview retains all three node names,
their provider assignments, and Vault, Let's Encrypt, and Mullvad connections.
The overview attaches shared protocol links to one data-plane group boundary,
with unconnected njalla1, work1, and yc1 inventory nodes inside provider groups
captioned "njal.la" and "Yandex Cloud". Each node uses a two-line label, with
"Node" above its full hostname.
The intermediate three sources rendered as light/dark SVG pairs and native WebPs.
They are now combined into one authoritative source. The current collection
contains 14 Mermaid sources representing 14 original pages, with 28 light/dark
SVGs selected by the documentation build.

The completed remote review identified that Forgejo's extracted Proxmox topology
predates its documented XCP-ng deployment. Its README and diagram caption now
explicitly classify that preserved topology as historical, with the deployment
section retaining the current contract.

The overview hierarchy revision was inspected as a native WebP: User appears
above Admin, Master belongs to Proxmox inside Control Plane, and the Proxmox
caption has two lines. All three data-node hostnames match the checked-in
Ansible inventory and DNS declarations. No edge crosses the nested captions.

The overview render was inspected after moving control-plane management and ACME
links to the Control Plane boundary. Their paths start at that boundary and do
not traverse its nested Master or Proxmox elements. User remains above Admin.

A render differing only in the Admin HTTPS endpoint confirmed that connecting
to the Control Plane boundary instead of Master removes excess vertical space
in its nested groups. All external control-plane links now attach to the outer
boundary, while User remains above Admin and all captions remain readable.

A native WebP render of the detailed control-plane view confirms two-line
Proxmox and 3x-ui VM captions with full hostnames. Both title plates fit their
containers, remain separate, and avoid the service edges.

The complete native WebP was inspected after combining the three sections on
one canvas. The initial combined render placed all three sections in vertical order,
with left-to-right flow inside each section and no visible edges between views.
User remains above Admin in the overview. Full node hostnames, provider captions,
boundary connections, Master, and two-line Proxmox/3x-ui VM captions are retained.
The former three source files and render targets are replaced by architecture.mmd
and one light/dark asset pair to avoid maintaining duplicate sources.

The wide-canvas revision places the overview above the control-plane and data-plane
sections side by side. The maintained architecture_webp target uses the shared
5:2 default and produces 2850 × 1140 pixels. The complete WebP was inspected with
all prior captions, group-boundary edges, and node hostnames retained.

The full WebP now selects the site's dark palette. Its maintained target built
successfully and the complete raster was inspected: light labels and edges remain
visible on the dark canvas, with the same 5:2 proportions and section layout.

The preview-only WebP target and README reference were subsequently removed at
the user's request. Documentation retains only its build-time light/dark SVG
asset pair. The reviewed dark 5:2 WebP remains solely in ignored task scratch;
no diagram WebP asset was committed or selected by the documentation consumer.
