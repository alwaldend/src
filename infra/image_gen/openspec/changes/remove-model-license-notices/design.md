## Context

Model license statements in this owner are prose in the operator READMEs plus
supporting clauses inside completed change records. The publisher terms live
upstream; the local text is a summary, not a duplicate of maintained source.
See [proposal.md](proposal.md) for the motivation.

## Goals / Non-Goals

**Goals:**

- Remove the model license notices the operator named from the owner's
  documentation and change records.
- Keep the change-record rationale and history coherent.

**Non-Goals:**

- Changing runtime behavior, model downloads, profiles, or Terraform/Ansible
  configuration.
- Removing unrelated license references: the rendered-diagram font license in
  `infra/image_gen/docs/BUILD.bazel`, and the third-party license requirements
  owned by `third_party/`.
- Removing non-license operational text: the gated Flux Hugging Face token
  requirement, immutable publisher revisions, and download sizes.

## Decisions

- Edit prose in place rather than deleting sections, so each affected paragraph
  keeps its operational content (token requirement, revision pinning, sampling
  settings) and loses only the license clause.
- Preserve completed checkboxes, approvals, publisher revisions and verification
  results in the active change records; strike only license wording. History
  stays intact and the records continue to describe what was done.
- Reconcile the `projects/comfyui_tools/README.md` sentence that both states a
  NoobAI restriction and links to the removed license record, so no dangling
  reference remains.

## Risks / Trade-offs

- Removing the notes removes the local warning of publisher restrictions
  (MiniMax territorial exclusions, Flux/NoobAI non-commercial terms, Llama
  community terms). The underlying obligations are unchanged; operators must
  consult the publishers directly.
- A broad text removal could delete adjacent operational facts. Task-level
  verification re-reads each edited paragraph to confirm the retained content.
