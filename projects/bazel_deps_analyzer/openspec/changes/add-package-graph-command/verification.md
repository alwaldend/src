# Package graph command verification

Observed 2026-10-28 in the task worktree, based on
`c4f8470da5c4ea07e3799e332bfb8fb2db4d634a`. Capture revision and implementation
identity are recorded separately below. Publication and final candidate checks
are recorded by the repository delivery receipts.

## Candidate identity

- `internal/graph/packages.go` SHA-256:
  `5ebb0c9911a8619d2fa96699a4f933ca592346365b724e13dc0217746bfa7cb4`.
- `cmd/bazel_deps_analyzer/main.go` SHA-256:
  `accf3b9c5c456ff52c32e71a8f11b9d964c94c858eeb1e02799386781515b72d`.
- Built analyzer SHA-256:
  `a9fbde4cc633eb7b18962ba5dd4376754ec9a8706bc1a9d10e8ce6c734057653`.

The binary is at
`bazel-bin/projects/bazel_deps_analyzer/cmd/bazel_deps_analyzer/bazel_deps_analyzer_/bazel_deps_analyzer`.

## Checks

Failure cases and E2E coverage were written before implementation. The command
hierarchy is `analyze target` and `analyze package`; bare root and bare
`analyze` show help. Package depth, namespace separation, membership, counts,
filters, reciprocal-direction preservation, JSON/Mermaid equivalence, invalid
input/flags, output preservation, and existing target modes are covered.

```sh
bazel_agent bazel test \
  //projects/bazel_deps_analyzer/test/e2e:e2e_test \
  //projects/bazel_deps_analyzer/openspec:validate_test
```

Both test targets passed in one invocation, which completed at 07:03:28 UTC
with exit status zero. The package fixture passes through the analyzer's
JSON-input rule and the native dark Mermaid WebP rule. Its retained artifacts
are `packages.mmd`, `packages.json`, `packages.dark.webp`, and
`packages-report.json` under
`bazel-testlogs/projects/bazel_deps_analyzer/test/e2e/e2e_test/test.outputs/`.

The pinned gofumpt and buildifier runs succeeded. `git diff --check` passed.
Detailed logs remain under ignored `out/bazel-dependency-diagram/`:
`package-validation.log`, `package-gofumpt.log`, and `package-buildifier.log`.

## Recorded root capture

The rebuilt analyzer processed `root-query.ndjson` using
`analyze package --package-depth 1`, once with default Mermaid output and once
with `--format json`. It generated `analyzer-packages.mmd` and
`analyzer-packages.json` under the same ignored task directory.

The independent comparison with the previously verified root overview passed:

- 287,492 selected targets and 390,589 selected target edges.
- Nine package groups; every group count and internal-edge count matched.
- All 42 directed package relationships and their target-edge counts matched.
- Mermaid draws 30 connections, including 12 reciprocal connections; expanding
  them reproduces exactly the JSON directed relationships.
- Ten unresolved dependency references remain reported. This is a partial
  capture, not a claim of complete repository reachability.

`analyzer-packages-verification.json` records the comparison. Generated Mermaid
SHA-256 is `1dc71d2c332413b46379d447e6b98dc52c142f85a2201744c2b5af27f0dcb168`.
The generated source was rendered through the repository's pinned Mermaid
renderer as `analyzer-packages.dark.webp` (3322 by 704 pixels). Visual inspection
confirmed readable captions for `root`, seven top-level packages, and one
`External` node, with no edge captions. The renderer command and exit status
are recorded in `analyzer-packages.render-command.json`; it completed
successfully. No Mermaid nodes, edges, or placement were edited manually.

The generated viewer consumes analyzer JSON memberships for depths one through
six across all fourteen workspace contexts. Independent comparison verified
membership, internal-edge counts, and every directed relationship at all 84
workspace/depth combinations. Evidence is retained in
`viewer-package-depth-verification.json`. At depth two the root capture has
237 groups; its target and edge counts are unchanged.

Independent source review reported no findings against the implementation
hashes above. Repository formatting and `git diff --check` passed; the reviewed
implementation hashes remained unchanged afterward.

## Explicit deps/srcs filtering verification

Observed 2026-10-28 on base `e710c4324c36d75328cc6894d730d6590b64291e`.
Failure cases and command/native-query E2E coverage preceded implementation.
The parent owns independent review, root overview rendering, repository-wide
gates, and delivery. This implementation made no commits or PR changes.

Pinned gofumpt, buildifier, and Prettier completed successfully. The final
behavioral validation command completed at 08:39:51 UTC with exit status zero:

```sh
bazel_agent bazel test \
  //projects/bazel_deps_analyzer/test/e2e:e2e_test \
  //projects/bazel_deps_analyzer/openspec:validate_test \
  //projects/bazel_deps_analyzer:bazel_deps_analyzer
```

Both tests passed and the standalone binary was built. New E2E evidence covers
both CLI modes and JSON-input rules, real native genquery capture, explicit and
non-explicit attributes, tools/data/visibility, overlapping labels, selector
conditions versus branches, all supported label encodings, alias omission,
generated producers, metadata groups, filtered traversal/depth, unchanged
defaults, deterministic aggregate counts, conflicting metadata, malformed
encodings, and output preservation. The hand fixture selects 18 package-view
targets and six edges; its target closure has six nodes and five edges.

The native capture confirms the `name` marker on rules, default metadata,
unflattened selectors, and upstream tool/implicit/nodep exclusion. Bazel 9.3.1
serializes LABEL_LIST_DICT through `stringListDictValue`; E2E uses that actual
encoding. Omitted and explicitly empty rule inputs and attribute label lists
are canonicalized before duplicate checks. The parent's empty-input finding
was reproduced before correction and verified afterward through both CLI modes.

The saved older root capture is rejected at its first RULE without a name
marker, with all required recapture flags. Its command returns status one,
writes no stdout, and preserves an existing output sentinel.

Pre-review implementation identities:

- `internal/graph/attributes.go` SHA-256:
  `ebe205336673cc82ea9990e785dbba833d2f746cdd526563cdb691f371f82a3d`.
- `internal/graph/graph.go` SHA-256:
  `150582aad253b2651952c0dfe926db5b49eaa138e17c79fb4b172e4635aecc8e`.
- Standalone analyzer SHA-256:
  `300e8a430a1ff0c488e5a1e88a62250c2fb99aadc81b6af04f15f3d3ed85cfd3`.

Seven repeatable filtered query/diagram/JSON/report artifacts are retained in
`out/bazel-build-dependency-filter/implementation/artifacts/` and the E2E
undeclared outputs. The same implementation scratch directory contains
`candidate.sha256`, `validation-handoff.log`, formatter logs,
`old-capture-rejection.json`, and the empty-input before/after reproducer.

One enriched root query ran from 08:34:49 to 08:37:41 UTC with the documented
metadata flags, `--lockfile_mode=off`, explicit `--keep_going`, and a 900-second
bound. Bazel returned status three without timing out. Its 104,258,963-byte
capture is partial: missing package/target references and Harbor errors remain.
The outer command session reported status 143 after the subprocess's status
three and completed receipt were written; it does not establish query success.

The capture is `out/bazel-build-dependency-filter/root-query.ndjson`, SHA-256
`224d339b49dd445ab96be6831ccdcb3c7b8f1e9898fde1c0e2d804d92f9adc20`.
`root-query.receipt.json` records the exact argv, query-start worktree diff and
untracked-file hashes, timestamps, effective keep-going setting, exit status,
and partial state. `root-query.restricted.log` retains diagnostics. The parent
owns analysis, comparisons, and rendering of this partial capture; these
implementation checks do not claim a complete repository graph.

### Native non-label selector correction

Independent review identified a valid Bazel selector None branch containing
`isDefaultValue: true` and `stringValue: ""` for a string-valued `srcs` attribute.
A native Starlark fixture and both-mode E2E checks were added before the fix.
The native diagram build failed as reported; its capture and failure log are
retained as `implementation/nonlabel-selector-native.ndjson` and
`implementation/nonlabel-selector-before.log` under the same task output root.

None defaults now pass through the underlying type decoder. Empty label
defaults and non-label defaults contribute no edges. Invalid label encodings
or nonempty dependency-label defaults still fail. The E2E fixture checks the
actual emitted string default and both modes with a disconnected string rule.
The corrected `attributes.go` SHA-256 is
`f95d82bc9f972ba600cb61fde5bdfe3499d86209d9d491e72cdd55f6e95f314f`.
Round-two review remains parent-owned. The root query was not repeated.

The final batch passed E2E and OpenSpec validation and built the standalone
analyzer, `//tools/mermaid:mmdc`, and `//tools/mermaid/cmd/render:render` at
08:49:12 UTC. The analyzer was immediately copied to
`out/bazel-build-dependency-filter/analyzer`, SHA-256
`d611786e09c74c6578bd78d1aa45bfe86e2f3b7943f1e51dc8d4438703024ab0`.
`implementation/final-artifacts.json` records executable and manifest paths;
renderer artifacts were verified present at 08:50:16 UTC. The validation log is
`implementation/validation-round1-final.log`. Eight retained filtered artifacts
include the native non-label selector capture. No output-base or cache cleanup
was performed, and Bazel was released to the parent after this batch.

### Final source and package overview acceptance

Independent round-two review reported no findings. The parent confirmed the
reviewed source/build/test manifest remained unchanged: SHA-256
`7bec56692592d206cb93bc815ea5db2550cb48414b7e5f8f240c2e12cd544bda`
for the path-sorted `sha256sum` manifest of eleven changed non-Markdown files.

The preserved analyzer converted the enriched root capture with
`analyze package --build-deps-only --package-depth=1` to both JSON and Mermaid.
It selected 255,407 targets and 232,036 target edges, aggregated into nine
package nodes and nineteen directed package relationships. No selected edge
references an absent record; the query's nonzero exit still makes the capture
partial. These counts cover the root workspace context, not separate nested
workspace queries.

Separate-route ELK and curved Dagre SVG/WebP variants preserve all nineteen
directed relationships. Source and rendered topology checks, label containment,
embedded fonts, viewer tabs, zoom, panning, and download-byte checks passed.
The parent inspected both WebPs: labels are readable and unobstructed; the
curved view provides more space between routes. Both retain the four explicit
`//projects` to `//infra` target edges, including documentation and CA inputs.

Artifacts and repeatable render commands are under
`out/bazel-build-dependency-filter/diagram/`. The accepted run is
`runs/2026-10-28T085536308Z/`, containing authoritative package JSON/Mermaid,
both rendered variants, standalone and inline viewers, and validation reports.
`out/bazel-build-dependency-filter/ergonomics-review.json` records the bounded
session review and artifact-availability mitigation. Root quality, semantic
lint, and publication are recorded by the parent delivery receipts.
