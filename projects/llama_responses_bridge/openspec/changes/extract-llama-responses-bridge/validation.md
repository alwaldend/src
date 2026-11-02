## Source acceptance

The project was extracted from local inference candidate
`559d33c39235b8f6f8a91c9fb5444532bc4903a3` onto master
`d15e849913f42839df69d7af8c09a3174cb58012`. Comparing the adapter, command,
and protocol test bodies after owner-path substitution found no behavior changes.
The command doc comment now describes the project rather than a component.

The existing command-level protocol suite passed from the new owner. It
starts the relocated executable and verifies JSON/SSE namespace routing,
forced selection, replayed calls, concurrent mappings, error paths, health,
and cancellation. The optional real Codex case was not rerun; it is opt-in
and remains available in the relocated suite. No GPU or model ran.

Thirty selected repository quality, formatting, protocol, and new-owner
OpenSpec tests passed. Semantic lint passed for the new project, project
inventory owner, and generated OpenSpec validation packages. Configured
formatters and the inventory generator updated only the reviewed project
registration and documentation projections. Guarded delivery reruns mandatory
gates against the exact committed candidate; receipts own publication state.

No infrastructure was deployed. The existing host PR must consume the root
`//projects/llama_responses_bridge` command after this prerequisite merges.
A task-local consumer patch removes the old implementation and rewires the
Ansible package and acceptance documentation; `git apply --check` passed.

## Session ergonomics

`BRIDGE-EXTRACTION-OWNER-PATHS` (fixture-tested): a separate master-based
worktree verifies the project without hidden dependencies on the infrastructure
branch. The command-level suite checks the relocated executable runfile.
The initial scratch scaffold path was corrected before project generation.
The existing managed Bazel cache required sandbox access; no duplicate cache,
new dependency, shared configuration change, or agent delegation was introduced.

## Landing-page follow-up

Added the project-owned site/content landing and registered it in
projects/projects.bzl. The generated project index now links the landing.
The shared Hugo site and site build test passed. Inspection of candidate HTML
confirmed the landing identity, purpose, limitations, documentation and source
links, project index entry, and AI, Go, and active-status taxonomy entries.
The task-local artifact is out/extraction/landing-html-verification.json.
No site deployment ran; the page becomes public through the ordinary site
publication workflow after merge.

The project-layout skill previously exempted repository-only work from landing
pages. The tools/agents change require-project-landing-pages replaces that
exemption for new first-party projects, with a dedicated offline eval case.

## Compression review correction

`BRIDGE-UPSTREAM-COMPRESSION` (fixture-tested): a review identified that
DisableCompression does not remove a client's forwarded Accept-Encoding.
The command-level backend now honors gzip negotiation. Before the fix the
fixture failed with HTTP 502, retained in out/extraction/compression-before.log.
The proxy Director now removes Accept-Encoding before forwarding; JSON and
SSE namespace round trips pass, retained in compression-after.log. The fixture
and correction use existing standard-library capabilities and no dependencies.
