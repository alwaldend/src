## Why

The operator asked to remove the model license notices from this owner's
documentation. The notices describe upstream publisher terms rather than
runtime behavior, and the operator no longer wants them maintained here.

## What Changes

- Remove model license statements from the operator-facing
  [rental README](../../README.md) and
  [model configuration README](../tf_setup/README.md): the Flux non-commercial
  descriptor and license-acceptance clause, the MiniMax H3 Community License
  territorial exclusion, the MOSS Apache 2.0 declaration, the FLUX.1 dev
  non-commercial and Apache 2.0 clauses, the NoobAI commercialization
  restriction, the Illustrious CreativeML OpenRAIL-M name, and the HiDream
  MIT/Apache 2.0/Llama component-license list.
- Strike the matching license clauses from the active change-record rationale
  while preserving approvals, publisher revisions, verification results and
  completed-task history.
- Reconcile the `projects/comfyui_tools/README.md` cross-reference so it no
  longer points to a removed license record.

Operational requirements that are not license statements remain: the gated Flux
Hugging Face token requirement, immutable publisher revisions, and download
sizes.

## Capabilities

None. This is a documentation change with no spec-level behavior delta.

## Impact

Documentation only: `infra/image_gen/README.md`,
`infra/image_gen/tf_setup/README.md`, the `add-imagegen-rental` and
`add-minimax-h3-video` change records, and
`projects/comfyui_tools/README.md`. No source, configuration, model download,
rental or deployment state changes.
