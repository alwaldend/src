---
title: Agent ergonomics comparison fixtures
---

# Agent ergonomics comparison fixtures

These four public cases measure bounded question answering, a local prose
edit, resumption after a changed postcondition, and a delivery validation
decision. They complement the existing Promptfoo skill suites. A successful
configuration check does not run these cases; a successful delivery decision
does not establish that Git publication or review handling works.

## Run a paired comparison

Use the same case definitions, subject model, tool permissions, and prompt
template for each revision. Freeze both instruction bundles before starting.
Stage only `AGENTS.md`, a catalog of skill names and descriptions, and complete
skill bundles, preserving relative reference paths. Record the originating
commit and a digest of the exact staged bytes. The candidate may initially be
identified by a source digest; bind its result to the final commit only after
confirming that its relevant inputs did not change.

For each case and revision, create a fresh task-local run directory containing
`instructions/` and `fixture/`. Copy only the case's `files` map into the
fixture. Keep `cases.json`, acceptance criteria, baseline answers, and other
runs outside the subject's input scope. Substitute the absolute run directory
and the case's `request` into `subject-prompt.txt`; retain that exact prompt.
Start a fresh subject without conversation history. Prefer alternating
revision order or running each pair concurrently. Record the actual order;
an unbalanced smoke run cannot support a duration comparison.

The coordinator can use an already available isolated agent provider; this
fixture package does not start another service, install tools, or manage
credentials. An in-session fresh agent is suitable for an exploratory smoke
comparison. Its behavioral path restriction is not an enforced filesystem
sandbox: record ambient instructions, tools, and any observable contamination.
Do not call this a fully isolated model evaluation or compare it directly with
the Promptfoo provider's enforced read-only runs.

The existing `rules_promptfoo` login-reuse mode requires exclusive use of the
persistent login by all participating processes. Do not start that mode while
unrelated agents use the same login. A separate authorized provider can run
the existing suites; these fixtures do not authorize changing host defaults.

## Judge outcomes independently

The coordinator, rather than the subject, checks:

- Every key in `expectedAnswer` matches the final JSON answer; additional
  prose belongs in its `reason` field.
- Every `expectedFiles` entry matches exactly, and all other fixture and
  instruction bytes remain unchanged. Inspect for newly created paths too.
- The explanation satisfies `humanCheck`. Record the reviewer and any
  disagreement; the same-model coordinator is not an independent model judge.
- Tool traces, when available, support `expectedSkillReads` and show no
  unauthorized operation. A correct answer does not prove skill routing, and
  unchanged bytes do not prove that no write was attempted.

Record separate answer, artifact, routing, and authority verdicts. Use `null`
with a reason for unobservable verdicts; do not turn an unknown into a pass.

## Report measurements with their coverage

One result per case and revision should include the prompt and fixture
digests, instruction-bundle digest, model identity when available, start/end
observation times, the final answer, artifact comparison, and reviewer verdict.
Preserve the smallest sufficient public evidence under the task's `out/` area.

| Metric                   | Acceptable source and limit                                                            |
| ------------------------ | -------------------------------------------------------------------------------------- |
| Task outcome             | Coordinator's answer and file checks plus explanation review                           |
| Skill routing            | Observed successful skill-file reads; otherwise `null`                                 |
| Tool calls and commands  | Complete provider trace; distinguish calls from shell commands, otherwise `null`       |
| Context bytes            | Exact bytes actually returned by observed reads plus prompt bytes; otherwise `null`    |
| Staged instruction bytes | Filesystem byte count; available-input size only, not context consumed                 |
| Tokens                   | Provider usage fields; otherwise `null`                                                |
| Duration                 | Coordinator or provider clock; identify whether queue and scheduling time are included |

Never use the subject's claimed success, counts, timing, or skill use as an
observation. If the environment exposes only the final answer, retain outcome
and filesystem evidence and leave trace metrics unavailable. One paired run
per case can find regressions; it cannot establish a reliable speedup or a
general success-rate improvement. Repeat and examine disagreements before
making those claims.
