## Decisions and evidence

Use the user-approved upstream requirements and immutable source archive, with
SHA-256 integrity. Expose only requested datasets and training code through
project-scoped targets. Generated locks and Gazelle import mappings remain owned
by their standard update targets. No host package installation is needed.

This companion record reconciles the shared dependency changes already made
under the project plan; it does not duplicate project requirements. Requirements
and manifest tests passed, and the HTTP and full model E2E consumers passed.
Next: validate the exact delivery candidate and publish using repo-delivery.

The user also requested namespaced imports in repo-python. Update its canonical
skill and offline eval case. The upstream archive receives an optional loader
path compatibility patch for Bazel runfiles; training behavior remains intact.

Ruff ICN003 checks only explicitly listed module names, not a universal ban.
Use the already-pinned ast-grep import_from_statement syntax rule through the
quality suite, covering tracked first-party Python files and avoiding downloaded
external source. Validate rejection of ordinary, aliased, multiline, relative
and star imports; accept namespaced imports, comments, and string literals.
Migrate existing first-party imports, preserving package export APIs.

Acceptance evidence: 36 checks passed, including real HTTP/model consumers,
shared lock/manifest consistency, offline skill validation, and quality. The
existing Blender audit passed with the migrated imports, and the Ansible bundle
built. Semantic lint passed for the application, Python skill/checker, and
Autoscroll packages. Publication validates the archived exact candidate again.
