---
title: Fumo render review
description: Fast measurable gates for Blender review packets
---

# Fumo render review

This package has three pinned commands and one snapshot-only live helper.
`render_packet` uses the repository's pinned Blender to produce fixed-camera
PNGs and optional component-ID masks. `fumo_review` checks those files before a
person spends time reviewing them. `live_checkpoint` proves exact ancestry and
hands a saved candidate to `checkpoint_render`, which produces the controlling
low-resolution beauty and silhouette views in pinned background Blender.

Visual review remains mandatory. The measurements catch known regressions; they
cannot decide whether a plush resembles the references, reads as sewn fabric,
or looks appealing.

## Freeze and check a live scratch scene quickly

Use this two-phase path only for early visual rejection. The first phase runs
inside Blender's embedded Python because the current scene and `bpy` data do
not exist in an external process. It proves that the open file is the capsule's
exact hashed parent, captures state before saving, suspends save handlers, and
saves a copy under repository-root `out/`. It does not render or call an
external modeling service. The second phase reopens that copy in the repository-
pinned background Blender with automatic script execution disabled.

Create a source-owned task capsule under `out/`:

```json
{
  "schema_version": 1,
  "checkpoint": "hair_stage3c_02",
  "scene": "Hair scratch",
  "resolution": { "width": 320, "height": 320 },
  "parent": {
    "path": "out/reimu_head/parent.blend",
    "sha256": "<64 lowercase hex characters>"
  },
  "source_packet": {
    "path": "out/reimu_head/source_trace/README.md",
    "sha256": "<64 lowercase hex characters>"
  },
  "hypothesis": "The traced panel stays outside the support.",
  "views": {
    "front": { "camera": "Review Front" },
    "rear": { "camera": "Review Rear" }
  },
  "vetoes": ["exposed beige support", "symmetric V-shaped fringe"],
  "scene_boolean_requirements": {
    "stage3d_clearance_gate_pass": true
  }
}
```

At most two distinct cameras and dimensions from 16 through 640 pixels are
allowed. Import and call the helper from the live Blender process, including a
Blender MCP code execution:

```python
import sys

workspace = "/absolute/path/to/repository"
if workspace not in sys.path:
    sys.path.insert(0, workspace)

from projects.renders.cmd.fumo_review.live_checkpoint import (
    open_checkpoint_parent,
    snapshot_checkpoint,
)

open_checkpoint_parent(
    "out/reimu_head/checkpoint_02.json",
    workspace,
)

# Apply the checkpoint's builder edits here.

request = snapshot_checkpoint(
    "out/reimu_head/checkpoint_02.json",
    "out/reimu_head/checkpoint_02",
    workspace,
)
print(request)
```

`open_checkpoint_parent` is mandatory at the start of the branch. It reloads
the frozen parent with auto-execution and Python load handlers disabled, hashes
it, and issues an in-process attestation. A later file load clears that
attestation. It requires the live Blender version to be no newer than the
hashed repository pin and rejects a parent authored by a newer Blender
major/minor than the live process. The latter prevents an older live process
from silently discarding unsupported newer-file data before the snapshot.
This distinguishes real parent ancestry from merely observing the same
filepath after its bytes have changed. Same-minor Blender file subversions are
not patch versions, so the pinned reopen and frozen scene checks remain the
final compatibility proof.

The output path must be a new, nonexistent directory strictly below repository-
root `out/`. A successful live handoff contains only `candidate.blend` and
`snapshot_request.json`; the request is canonical JSON capped at 16 MiB and is
written last. Render it from the
workspace root:

```sh
bazel_agent run //projects/renders/cmd/fumo_review:checkpoint_render -- \
  --snapshot-request out/reimu_head/checkpoint_02/snapshot_request.json
```

The renderer atomically claims the request, rechecks every source, tool, parent,
dependency, candidate, and request hash, then validates the scene, camera
projection, and boolean preconditions again. It forces a bounded Eevee preview,
renders beauty and black-on-white silhouette PNGs, checks their dimensions and
foreground, rehashes all inputs, and writes `manifest.json`, `VERDICT.md`, and
finally `READY`. A directory without `READY` is incomplete and must not be
reviewed or resumed in place.

The final manifest distinguishes live and render Blender versions and records
the exact capsule, parent, source packet, dependencies, tool sources, candidate,
camera projection, view layers, engine, output hashes, timings, and T0–T5
timeline. T1 records the completed candidate save and T2 the completed snapshot
hash/reverification. Request publication is instead proven by its last atomic
rename and the helper's successful return; a JSON payload cannot truthfully
timestamp its own later publication. T4 and T5 remain empty until a person
supplies the absolute verdict and updates the goal log.

Both phases fail closed on stale hashes, unsafe or symbolic-link escape paths,
missing scene or camera names, duplicate cameras, false or missing declared
scene booleans, an empty scene, external dependencies outside the workspace,
or output ownership conflicts. Builders should publish cheap topology,
clearance, or convergence gates as scene boolean properties and put their
expected values in the capsule. The helper does not infer builder success from
object names or veto prose. A quick checkpoint may reject a bad representation;
it can never accept the model or replace the pinned full packet, component
masks, regressions, or implementation-blind reviews.

## Render a packet

Keep the source `.blend`, JSON spec, and output packet under the repository-root
`out/` directory. The command resolves its three paths from the workspace root
and never saves the source `.blend`:

```sh
bazel_agent run //projects/renders/cmd/fumo_review:render_packet -- \
  --blend-file out/reimu_attempt/body_coupon.blend \
  --spec out/reimu_attempt/render_packet.json \
  --output-dir out/reimu_attempt/packet
```

Schema version 1 names the exact Blender scene, camera objects, and mesh
objects to use:

```json
{
  "schema_version": 1,
  "scene": "Review Scene",
  "resolution": { "width": 512, "height": 512 },
  "beauty_views": {
    "front": { "camera": "Review Front", "output": "front.png" },
    "side": { "camera": "Review Side", "output": "side.png" }
  },
  "component_id_passes": {
    "front_ids": {
      "camera": "Review Front",
      "output": "front_ids.png",
      "objects": {
        "Front Skirt Panel": [255, 0, 0],
        "Left Sleeve": [0, 255, 0]
      },
      "fallback": { "mode": "color", "rgb": [0, 0, 0] }
    },
    "side_ids": {
      "camera": "Review Side",
      "output": "side_ids.png",
      "objects": { "Front Skirt Panel": [255, 0, 0] },
      "fallback": { "mode": "hide" }
    }
  }
}
```

Every output must be a relative `.png` path beneath the packet directory, and
no two views or passes may use the same path. A component pass must explicitly
choose whether unlisted renderable meshes are hidden or assigned a fallback
color. Non-mesh geometry is hidden in component passes. The driver validates
all scene, camera, and mapped mesh names before rendering anything.

Component colors are 8-bit sRGB requests. The ID renderer uses unlit emission,
inverse-sRGB scene-linear values, the Standard view transform, no look, and no
dithering. Consequently flat interiors decode to the requested bytes; normal
anti-aliasing can blend colors along silhouette edges.

The packet's `manifest.json` records Blender's version, source `.blend` and
render-spec paths and SHA-256 values, dimensions, camera names, component
mapping, output names, and output SHA-256 values. Beauty views render before
any temporary in-memory material changes. The driver does not call a Blender
save operation.

## Audit a packet

Keep attempt inputs and generated reports under the repository-root `out/`
directory. Paths in the configuration are resolved relative to the JSON file.

```sh
bazel_agent run //projects/renders/cmd/fumo_review:fumo_review -- \
  --config out/reimu_attempt/review.json \
  --output-dir out/reimu_attempt/review
```

A gate failure still writes both reports and returns exit code 1. Invalid JSON,
non-finite numbers, unsafe or colliding output paths, replaced inputs, and
other invalid configuration or publication requests return exit code 2 without
publishing a report. The output directory must not already exist or be a
symbolic link; the tool builds an owned staging directory beside it and
publishes the completed report with an atomic no-replace rename.

Race-safe publication uses Linux `renameat2(RENAME_NOREPLACE)` and is therefore
Linux-only. Put the output beneath the repository's task-owned `out/` tree with
a parent directory trusted against concurrent mutation. Accidental output
collisions and symbolic-link traversal are rejected; defending against a
malicious same-user process replacing that owned parent is outside this tool's
threat model.

## Audit configuration

```json
{
  "schema_version": 1,
  "title": "Reimu body coupon",
  "required_views": ["front", "side", "rear", "three_quarter"],
  "views": {
    "front": {
      "image": "front.png",
      "mask": "front_ids.png"
    },
    "side": { "image": "side.png" },
    "rear": { "image": "rear.png" },
    "three_quarter": { "image": "three_quarter.png" }
  },
  "required_files": ["blend"],
  "files": { "blend": "body_coupon.blend" },
  "components": [
    {
      "name": "front skirt panel",
      "view": "front",
      "rgb": [255, 0, 0],
      "tolerance": 0,
      "thresholds": {
        "pixel_fraction": { "min": 0.08, "max": 0.3 },
        "bbox_width": { "min": 0.35, "max": 0.7 },
        "connected_components": { "min": 1, "max": 1 },
        "max_horizontal_run": { "max": 0.65 }
      }
    }
  ]
}
```

If a view has a `mask`, component colors are sampled there; otherwise they are
sampled from its display image. Matching uses per-channel RGB tolerance. Alpha
is ignored. Bounding boxes and `max_horizontal_run` are fractions of image
width or height. Coordinates use a top-left origin and exclusive maximums.

Supported threshold names are `pixel_count`, `pixel_fraction`, `bbox_x_min`,
`bbox_y_min`, `bbox_x_max`, `bbox_y_max`, `bbox_width`, `bbox_height`,
`connected_components`, `max_horizontal_run`, and
`max_horizontal_run_pixels`. Every threshold accepts `min`, `max`, or both.

The decoder accepts 8-bit, non-interlaced grayscale, grayscale-alpha, RGB, and
RGBA PNGs. Use flat, unlit, color-managed component-ID renders for masks so
anti-aliasing and shading do not create accidental regions.

PNG files are limited to 64 MiB, each dimension to 8192 pixels, the decoded
image to 16,777,216 pixels, decompressed scanlines to 64 MiB, and the chunk
count to 4096. A packet may decode at most 33,554,432 pixels across its views
and masks; each image's dimensions are charged to that budget before scanline
decompression. The JSON configuration is limited to 1 MiB; it may name at most
32 views, 64 files, and 256 components. Named files are hashed with bounded
streaming reads, and are limited to 1 GiB both individually and in aggregate.
Component gates are limited to 268,435,456 aggregate source-pixel evaluations,
checked before constructing the next component mask.

Every result records the SHA-256 and byte size of the exact configuration,
view images, masks, and named files it audited. Publication revalidates that
evidence with one bounded read per record and validates the semantic
audit-result digest. Displayed view images are copied during that same
validation read; HTML is then rendered only from the already-validated
records. The aggregate
`evidence_sha256` covers stable `role`, `name`, `sha256`, and `size` fields;
paths and read limits remain metadata and do not make the digest
machine-dependent. Replacing an input before its publication-validation read
therefore fails instead of producing a stale PASS page, and later source-image
changes cannot alter the snapshotted images in an already-published report.
The tool does not lock producer files as a transaction: finish the render/
model producer before publication, because a hostile concurrent writer can
change a different source path after that path's sequential validation read.

The render spec and audit config remain separate on purpose: the producer
controls Blender scene names and output bytes, while the auditor controls the
review gates applied to those outputs.
