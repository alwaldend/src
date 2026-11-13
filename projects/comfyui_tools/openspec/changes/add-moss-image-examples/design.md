## Context

The existing graph and node libraries generate native JSON from model-specific
Libsonnet. Upstream examples are pinned for provenance to workflow_templates
revision `36ecd0663961c5e1e83f7c0ea521658b33b36bbc`.

## Goals / Non-Goals

**Goals:** Preserve model-specific sampling and conditioning in compact examples.

**Non-Goals:** Copy large upstream subgraphs, add prompt enhancement, train adapters
or add a geometry-estimation dependency.

## Decisions

- Generate six flat graphs with the existing graph library. SFT/Base use 25
  ER-SDE/simple steps, CFG 2.5 and layer-22 skip guidance through 50% of sampling.
- Layout examples add native bounding boxes and MOSS layout prompt encoding.
- Edit uses the SFT model and edit LoRA, the `Moss-Image-Edit` reference format,
  fresh noise sized from the reference, and 30 Euler/normal steps at CFG 4.
- Control uses prepared region/depth maps, the SFT/control LoRA, layer-24 skip
  guidance, shift 5, temporal score rescaling and upstream ER-SDE sampling.
  Automatic MoGe depth extraction is excluded because it is an additional model.
- Use the INT8 encoder consistently across examples, avoiding a second GPTQ
  download; weights and runtime remain those requested or already approved.

## Risks / Trade-offs

- Examples are not quality guarantees → label them examples and record which
  graphs actually execute on the GPU.
- Input maps or image files are user-supplied → document required filenames,
  dimensions and the depth convention.
- Native dynamic edit inputs can differ across releases → compare generated
  input names against the deployed server's object information.

## Validation

All six graphs generated successfully with graph port and reciprocal-link
validation. Jsonnet formatting and affected semantic lint passed. The existing
offline provider acceptance suite verified that all six example graphs were
staged without content changes, within the startup payload limit. The native
MOSS, bounding-box and skip-guidance contracts were compared with ComfyUI
v0.44.0 source. Root quality, workflow-format and buildifier checks passed
29 tests. Live SFT text-to-image and Base layout-guided examples completed on
the requested RTX PRO 5000 rental. All six examples and both adapters were
discovered by the server. Edit/control execution remains untested; see the
matching infrastructure change for the deployment and output evidence.
