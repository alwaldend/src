## 1. Operator documentation

- [x] 1.1 Remove the model license clauses from `infra/image_gen/README.md` and
      `infra/image_gen/tf_setup/README.md`, retaining the gated-token
      requirement, revision pins and sampling settings; verify no model license
      term remains in either file.
- [x] 1.2 Re-read every edited paragraph and confirm the non-license
      operational text is intact.

## 2. Change records and cross-reference

- [x] 2.1 Strike the license clauses in the `add-imagegen-rental` and
      `add-minimax-h3-video` change records while preserving approvals,
      revisions and completed tasks; verify the records still validate.
- [x] 2.2 Update the `projects/comfyui_tools/README.md` cross-reference so it no
      longer points to a removed license record; verify no dangling link.

## 3. Validation

- [x] 3.1 Run the owner OpenSpec suite
      (`bazel test //infra/image_gen/openspec:validate_test`) and verify it
      passes.
- [x] 3.2 Verify a repository search finds no model license statement left in
      the affected files.
