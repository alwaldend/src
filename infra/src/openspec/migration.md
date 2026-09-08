# Goal migration to OpenSpec

On 2026-09-08, all nine tracked maintained project goal records at source commit `550d7e79b1f5fdbc2b6017b75178471d6914082f` were migrated into native OpenSpec changes. Each change now lives in its owning project’s `openspec/` workspace. Eight completed records are imported as archived history. The Reimu Fumo change remains open with blocked execution and no accepted candidate.

The [machine-readable inventory](migration.json) maps all 222 tracked source files (834,320 bytes) to byte-identical provenance snapshots. Test fixtures and ignored task scratch are outside the maintained-record inventory.

## Record mapping

| Former maintained goal                                    | OpenSpec change                                                                                                                                 | Preserved outcome | Preserved execution   |
| --------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- | --------------------- |
| `projects/agents/goals/agent-system-phase-1`              | [agent-system-phase-1](../../../projects/agents/openspec/changes/archive/2026-09-08-agent-system-phase-1/proposal.md)                           | achieved          | paused                |
| `projects/agents/goals/agent-system-phase-2`              | [agent-system-phase-2](../../../projects/agents/openspec/changes/archive/2026-09-08-agent-system-phase-2/proposal.md)                           | achieved          | paused                |
| `projects/agents/goals/agent-system-phase-3`              | [agent-system-phase-3](../../../projects/agents/openspec/changes/archive/2026-09-08-agent-system-phase-3/proposal.md)                           | achieved          | paused                |
| `projects/agents/goals/agent-system-phase-4`              | [agent-system-phase-4](../../../projects/agents/openspec/changes/archive/2026-09-08-agent-system-phase-4/proposal.md)                           | achieved          | paused                |
| `projects/agents/goals/agent-system-phase-five-followup`  | [agent-system-phase-five-followup](../../../projects/agents/openspec/changes/archive/2026-09-08-agent-system-phase-five-followup/proposal.md)   | achieved          | paused                |
| `projects/agents/goals/agent-system-phase-six-ergonomics` | [agent-system-phase-six-ergonomics](../../../projects/agents/openspec/changes/archive/2026-09-08-agent-system-phase-six-ergonomics/proposal.md) | achieved          | paused                |
| `projects/agents/goals/repo-agent-system`                 | [repo-agent-system](../../../projects/agents/openspec/changes/archive/2026-09-08-repo-agent-system/proposal.md)                                 | achieved          | paused                |
| `projects/renders/goals/reimu-fumo-finish`                | [reimu-fumo-finish](../../../projects/renders/openspec/changes/reimu-fumo-finish/proposal.md)                                                   | open              | blocked               |
| `projects/mcp_cordis/goals/runtime_extensions`            | [runtime_extensions](../../../projects/mcp_cordis/openspec/changes/archive/2026-09-08-runtime-extensions/proposal.md)                           | achieved          | unavailable in source |

## Mapping and acceptance

- Objective and scope map to `proposal.md`; accepted plans, decisions, rejected approaches, and constraints map to `design.md`.
- The latest acceptance criteria map to `tasks.md` and meaningful capability delta requirements. Each YAML criterion keeps its stable ID and revision in `migration.json`; original revisions, attempts, review verdicts, and evidence remain in provenance.
- Completed task checkboxes reproduce explicit historical acceptance reviews. They do not revalidate today's source. The Markdown Cordis record declares Complete and no failing criteria; its absent execution field stays unavailable rather than being guessed.
- Reimu acceptance and continuation tasks remain unchecked. A supported authorized unblock can change execution while the outcome remains open until every required criterion passes.
- Completed records were placed directly in their owner’s `openspec/changes/archive/2026-09-08-<id>/`. No archive or spec-sync command applied their historical deltas. Current baseline specifications were derived separately from maintained owner source. Do not reapply historical deltas.
- Historical Phase 5 acceptance includes an explicit fixture-gap inventory. Its accepted review and still-pending prose are both retained; migration does not claim the missing infrastructure fixtures were implemented.

## Provenance and navigation

Each change has `migration.json`, a checksum inventory at `provenance/manifest.json`, original record bytes at `provenance/source/`, and a provenance navigation guide. The old agent goal landing and BUILD file are also retained under the repository-agent-system archive. The historical BUILD snapshot uses `.original` to preserve bytes without creating a new Bazel package.

Among preserved Markdown links, 47 links depend on the former location; companion navigation tables resolve them where the referenced artifact exists. 46 historical targets are unavailable in the checkout. Original bytes and historical Git identities are not rewritten to disguise unavailable scratch or superseded paths.

The former maintained goal directories are removed; native OpenSpec work is the continuation surface. Provenance files are excluded from generic formatting so their source digests remain stable. The migration inventory is import evidence, not an alternate work scheduler or mutable status catalog.
