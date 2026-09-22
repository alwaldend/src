# Phase 6B skill-validation specificity evidence

Learning-proposal defect signature
`skill-validation-error-underspecified` (30 s baseline): the skill validation
aspect ran per target but failures did not name the offending skill, so
agents had to identify which skill failed among many aspect actions.

Change: the validator now prefixes every error line with the skill name,
e.g. `ERROR: skill spellcheck: interface has unexpected key(s):
short_descriptionX`.

Live regression (fixture-tested + live): deliberately corrupted
`projects/agents/skills/spellcheck/agents/openai.yaml` (renamed
`short_description` key), built the skill target, and confirmed the streamed
error names `skill spellcheck` and the unexpected key. Restored the file and
confirmed the target builds clean. Nested-workspace test suite passes 17/17.

Scope note: per-error file-path context (e.g. `SKILL.md` vs
`agents/openai.yaml`) was already added in an earlier candidate on this
branch; the skill-name prefix removes cross-skill ambiguity, which was the
measured failure mode.
