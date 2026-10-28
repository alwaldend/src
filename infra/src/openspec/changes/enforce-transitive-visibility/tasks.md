## 1. Directive generator

- [x] 1.1 Create the standalone project module and write behavioral tests before implementation; verify the fixtures cover the design's failure cases.
- [x] 1.2 Implement inherited directive generation and conservative package updates; verify observable output and idempotence through the real Gazelle executable.

## 2. Repository integration

- [x] 2.1 Register the plugin in Gazelle, module discovery, documentation, and the workspace validation matrix; build the standalone module and repository Gazelle binary.
- [x] 2.2 Enable repository-local enforcement and update first-party Bazel compatibility declarations through their owning workflows; verify module freshness and locks.
- [x] 2.3 Finalize tree consumer groups and directives under the selected exception policy, generate declarations, and require an unchanged second Gazelle run.

## 3. Validation and delivery

- [x] 3.1 Verify direct rejection and indirect behavior with repeatable fixtures; document the Bazel 9.3.1 intermediary-policy limitation accepted by the user.
- [x] 3.2 Run repository formatting and required build/test validation; classify existing policy violations separately from implementation defects and record results.
- [x] 3.3 Complete repository delivery with task-owned commits, remote publication, and a linked pull request. PR 210 is published; delivery receipts own subsequent revisions and verification status.

## 4. Data-owned public certificate exports

- [x] 4.1 Move the four public CA outputs into their data owners, preserve their contents, and update every build and file-path consumer.
- [x] 4.2 Add separate Terraform refresh configurations under data/ssh/tf and data/ssl/alwaldend.com/tf; read public values through the existing Vault provider without consuming infrastructure targets or state, and retire the infra-owned local_file resources without destroying files.
- [x] 4.3 Run native Terraform formatting and schema validation, confirm collection dependencies exclude infra, build affected deployment packages, and record candidate validation. Rebase onto current master as requested; publication remains subject to task 3.3.

## 5. Bzlenv tool references

- [x] 5.1 Replace references to deleted tool packages with their existing third-party binary targets, preserving the environment's tool selection and exposing each executable under its configured command name.
- [x] 5.2 Build and lint bzlenv, inspect the generated executable paths and runfiles, and record validation. Final candidate delivery remains subject to task 3.3.

## 6. Gradle wrapper runfiles

- [x] 6.1 Replace the empty-file-list assumption with the existing shell runtime's declared files and runfiles; derive Gradle and Java launch paths from the declared targets.
- [x] 6.2 Build and lint the wrapper and run its non-mutating version command. Final candidate delivery remains subject to task 3.3.

## 7. Full dependency capture cleanup

- [x] 7.1 Remove the unused LeetCode downloader image and its image-only test, retire the Alpine declaration, and regenerate the root module through its owner.
- [x] 7.2 Remove the obsolete requirements alias and update the root documentation to the existing uv lock command; validate the surviving downloader and Python build interfaces.
- [x] 7.3 Capture root and external target graphs without depth limits, verify dark WebP diagram variants, and keep all diagram artifacts under ignored out/ as requested.

## 8. Central policy configuration and declaration ordering

- [x] 8.1 Replace the current directive contract with validated TOML path mappings; preserve the existing consumer groups and exceptions.
- [x] 8.2 Update real-Gazelle failure fixtures before implementation, covering configuration errors, path selection, workspace boundaries, declaration ordering, preservation, and idempotence.
- [x] 8.3 Read tools/transitive_visibility/transitive_visibility.toml, remove BUILD comment overrides, and place package declarations after loads.
- [x] 8.4 Regenerate affected workspaces, validate the exact outputs and delivery gates, and update the existing pull request without committing diagram artifacts.
