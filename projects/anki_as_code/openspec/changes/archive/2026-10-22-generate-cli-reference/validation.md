## Implementation evidence

Observed on 2026-10-18 before delivery preparation:

- The generator executable E2E and existing CLI command tests pass. The E2E
  retains a `cli/` tree in its undeclared outputs and confirms a short index,
  all five root-group links/pages, isolated complete help, identical successive
  generations, inherited namespace flags, aliases, request fields, native help
  guidance, and invalid invocation/output failures. The updated E2E was written
  before the implementation and failed against the previous single-file tool.
- The sandboxed `cli_reference` action and complete Hugo site build pass.
  Rendered `/docs/projects/anki_as_code/docs/cli/index.html` has no help blocks
  and links once to each child in its content. Child help sections are `api`
  (542), `col` (10), `config` (3), `help` (1), and `server` (3). Each page contains
  only its group's headings and retains parent navigation. Collection flags and
  nested API anchors are present. The README link remains valid. Generated
  Markdown is absent from Git.
- Exact-candidate quality, semantic lint, E2E, and site checks are selected in
  `out/anki-cli-docs/tree/checks.json`; delivery receipts own final results.
  The authorized redeployment uses the stamped release archive and existing
  Pages target, followed by build-status and public HTTPS checks of all six pages.

## Session observations

- `CLI-DOCS-NIL-APPLICATION` (fixture-tested): the initial generator passed a nil
  interface, which failed while state commands captured workflow method values.
  A documentation-only application embeds the interface so construction works
  while workflow execution remains unavailable; the executable E2E passes.
- `CLI-DOCS-DISCOVERY-SCOPE` (live): an initial broad filename search matched
  collection fixtures and exceeded the output budget. Subsequent searches use
  the owning project and bounded paths.
- `CORDIS-ARGUMENT-SHAPE` (live): three rejected discovery calls established
  that invocation uses `scope`, `package`, `tool`, and `arguments`. Subsequent
  reads used that shape. These observations require no shared policy changes.
- `CLI-TREE-DUPLICATE-LIST` (fixture-tested): Docsy's automatic section list
  duplicated the generated child links. The index uses the supported `no_list`
  parameter; rendered content now has exactly one link per child.
- `DELIVERY-TOOL-CALL-ENVIRONMENT` (live): the harness changes
  `CODEX_TOOL_CALL_ID` between tool calls. Run validation and publication in one
  invocation so the existing environment guard checks a stable environment.

Resume from this change and the task-local delivery receipt, rather than
reconstructing validation or publication state from conversation history.
