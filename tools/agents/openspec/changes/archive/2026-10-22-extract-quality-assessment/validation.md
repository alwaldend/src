# Extraction validation

Initial extraction observed on 2026-10-20 in the task worktree.

- `//.agents:write_skills` generated the new direct relative discovery link.
  `//.agents:write_skills_test` passed.
- Both skills passed `skill-creator/scripts/quick_validate.py`.
- Both skills' `eval_config_test` and `eval_routing_config_test` passed: four
  offline configurations, with no model or judge calls.
- Inspected the generated decision `eval_config_test.sh` and its runfiles
  manifest: both SKILL.md files and UI metadata are staged as siblings.
  The decision skill's relative link resolves in that layout.
- Inspected the source against the acceptance scenarios: quality assessment
  owns general scrutiny and calibrated judgment; decision review retains
  alternatives, trade-offs, verdicts, authority distinctions, and reopening.
- `openspec validate extract-quality-assessment --strict` passed through the
  repository wrapper. The configured repository formatter completed; its only
  additional diff was task-list indentation.

Live answer quality and routing remain unverified. The checked-in suites
contain seven quality cases, three standalone routing cases, five retained
decision cases, and one composition case. Offline results establish harness
validity, not correct judgment or successful live research.

Final mandatory repository quality, semantic lint, skill builds, and repeated
configuration checks are bound to the prepared candidate by repo-delivery.
Its ignored receipts and logs own exact commit identity and publication state.

## Question skill removal

The user subsequently requested removal of answer-question. Its canonical
package, evaluations, discovery declaration, generated link, and routing row
are removed. The maintained rules_skills and rules_promptfoo examples and the
agent ergonomics assessment fixture now refer to quality-assessment.

Regeneration reported 44 discovery entries, and the exact-state discovery test
passed. Strict OpenSpec validation passed again. A bounded reference search of
`AGENTS.md`, `.agents`, and `tools` found no remaining live package consumers;
remaining mentions describe removal, preserve archived history, or name a
self-contained Gazelle fixture that does not load the deleted package.
The final delivery validation also builds that updated ergonomics fixture.

## Session ergonomics

`assessment-extraction-discovery-overread`: one unnecessary full worktree
inventory and two oversized batched read responses caused truncation. A bounded
checkout inspection was sufficient; subsequent reads used live Cordis handlers
and explicit ranges. One generic tool-search call was unavailable; live tool
metadata exposed the supported handlers. This is an observed session cost,
not a defect in the assessment skills or evidence of their runtime behavior.
No shared workflow changes are proposed.
