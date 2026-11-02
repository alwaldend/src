# Goal: describe the AL command-line product

Read README.md, glossary.md, foundations/invocation.md, foundations/state.md, and
releases/bundle.md. Match the pilot's depth and eight-section skeleton.

## Source of truth

Use the single source snapshot named in README.md. Read the owning command,
state handling, defaults, and behavior tests before writing. The Reference list
is the reading map. Drafting describes executable behavior; it does not authorize
source changes by itself. When implementation fixes are separately authorized,
refresh the shared source pin and affected descriptions together, retaining the
original provenance of earlier observations. Do not treat prose as runtime evidence.
The project path is `projects/al`, while Bazel consumers retain `@rules_al`.
Keep current source links and retained evidence locations accurate after moves;
original runtime paths inside historical records remain evidence of those runs.

## Writing rules

Describe the operator's experience, with technical mechanisms only in
`> Technical note:` blocks where they explain an observable consequence. Use
sentence-case headings and glossary terms. Explain surprising behavior directly;
record suspected defects instead of smoothing them over. Keep the fixed modifier,
interrupt, and cross-cutting rows from README.md in every feature document.
Use relative links to foundations. Include one `stateDiagram-v2` for each
interaction. End with Open questions and verification and the shared source
footer. Mark scripted observations separately from manual verification.

## Foundation references

- [Entry, streams, and interruption](foundations/invocation.md).
- [State boundaries and receipts](foundations/state.md).
- [Bundle generation pilot](releases/bundle.md).
- [Publication state handoffs](releases/publication.md#the-interaction-event-by-event).

## State ownership and order

The bundle document owns local generation and minimal source identity. The
publication document owns artifact upload/deletion, nested `release website`
extraction/activation and their dry runs, remote staging, and operation receipts.
Versioning owns source refs and stamping. Delivery owns candidate, push, PR,
and review state. The operator's ordinary Bazel commands own build and test
outcomes. Do not mix their meanings of completion.

Write the pilot, foundations, and publication yourself before delegating remaining
documents. Then check every result for consistent vocabulary, complete tables,
working relative links, source footer, and explicit unknowns. Add checklist rows
for observable claims and each nontrivial modifier/interrupt cell. Keep stable IDs.

## Working rules

Use the existing owner-local OpenSpec change and guarded delivery workflow.
The repository's single aggregate commit convention supersedes per-document
commits in the upstream template. Keep scratch and private evidence under ignored
`out/unify-release-versioning/rebase-description`. Do not initialize a nested
repository or change product code to fit a description. Authorized source fixes
belong to the implementation workflow; update affected claims, checklist results,
and triage history afterward. Update structure and coverage when files change.
No coverage row may remain not started at handoff; drafted is honest when manual
checks remain. Never infer rollback from a failure exit.
