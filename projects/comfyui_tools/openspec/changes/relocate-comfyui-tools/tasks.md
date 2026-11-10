## 1. Relocate the assets

- [x] 1.1 Move only custom nodes and workflows, preserve runtime contents and restore the CLI to tools/comfyui.
- [x] 1.2 Update Terraform file paths, Bazel data inputs, target visibility and documentation ownership.
- [x] 1.3 Regenerate project and owner inventories and verify a second run is stable.

## 2. Project integration and verification

- [x] 2.1 Register the project and its landing content without adding a release target.
- [x] 2.2 Verify node behavior, Terraform packaging and representative rendered documentation without accessing a rental.
- [x] 2.3 Pass repository quality and affected semantic lint checks on the prepared candidate.

## 3. Delivery

- [x] 3.1 Review the complete prepared candidate and retain the comparison and finding dispositions before guarded PR publication.

## 4. Node directories

- [x] 4.1 Place Prompt Sequence in nodes/prompt_sequence with its own source, test, web script and README; verify imports, Bazel consumers and Terraform packaging after the move.

- [x] 4.2 Group workflow documents by model, retaining native filenames and JSON contents; verify both configured profiles and documentation links.
