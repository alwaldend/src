---
name: repo-openspec
description: >-
  Integrate the pinned OpenSpec skills with this repository: select owner-local
  workspaces, invoke the Bazel CLI wrapper, preserve maintained work, and route
  operations to upstream skills. Use alongside the applicable upstream skill.
---

# Repository OpenSpec integration

`AGENTS.md` owns when a change record is required and the user's request owns
scope and authority. This skill owns repository integration; the upstream
`openspec-*` skills own operation and artifact procedures.

## Select the workspace and maintained work

Use `<owner>/openspec/` for the component that owns the affected behavior.
Read its README, `openspec/config.yaml`, and relevant existing specifications.
For example, parser work belongs in `tools/rules_docs_gazelle/openspec/` and
Vault module work in `infra/vault/openspec/`. `infra/src/openspec/` owns the
repository's shared evolution and contracts. Keep requirements with their
owners when work spans components.

Reuse a matching active change before creating another. Resume from its tasks,
design, and linked evidence, then recheck mutable inputs. Keep task state and
acceptance evidence there; raw logs and caches belong under ignored `out/<task>/`.
For a record without a behavioral delta, set `skip_specs: true` in its
`.openspec.yaml` rather than inventing requirements. Legacy migrated history
is context, not a second writable task record.

## Invoke the pinned CLI

Run from the root Bazel workspace with `bazel-agent` and `repo-bazel`.
Translate every upstream `openspec <arguments>` example into:

```sh
OPENSPEC_PROJECT=<owner> bazel_agent bazel run //tools/openspec -- <arguments>
```

`OPENSPEC_PROJECT` selects the owner's working directory without entering a
nested Bazel workspace. Keep it set for all operations in that workflow;
without it the wrapper defaults to `infra/src`. Use this wrapper instead of
installing or invoking a host-global CLI.

## Select the upstream operation skill

Load the skill for the actual phase and follow its procedure through the wrapper:

| Operation                                         | Upstream skill                                            |
| ------------------------------------------------- | --------------------------------------------------------- |
| Investigate requirements or design                | `openspec-explore`                                        |
| Create a proposal and planning artifacts together | `openspec-propose`                                        |
| Start a change step by step                       | `openspec-new-change`                                     |
| Create the next artifact                          | `openspec-continue-change`                                |
| Generate all planning artifacts                   | `openspec-ff-change`                                      |
| Revise existing planning artifacts                | `openspec-update-change`                                  |
| Implement maintained tasks                        | `openspec-apply-change`                                   |
| Verify implementation against artifacts           | `openspec-verify-change`                                  |
| Sync delta specifications without archive         | `openspec-sync-specs`                                     |
| Archive one or several completed changes          | `openspec-archive-change`, `openspec-bulk-archive-change` |
| Guided workflow onboarding                        | `openspec-onboard`                                        |

User authority and repository policy continue to apply while following an
upstream procedure. Artifact generation grants no deployment authority or new
approval requirement. Passing CLI structural validation alone does not prove
implementation acceptance; use the owning checks against the actual candidate.
Use `repo-delivery` for publication.

## Preserve managed discovery

`//.agents:write_skills` owns discovery. The repository-specific skill lives
under `tools/openspec/skills/repo-openspec`; upstream skills are generated unchanged
from the pinned archive in `third_party/org_fissionai_openspec/`. Do not run
OpenSpec tool-integration generation over `.agents/skills` or edit those
upstream projections. Update their owning pin and regenerate discovery instead.
