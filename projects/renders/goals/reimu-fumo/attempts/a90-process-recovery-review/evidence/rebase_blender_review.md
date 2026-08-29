# Repo-Blender Rebase Conflict Review

## Scope and state

Reviewed only the four requested unmerged files while the rebase was paused:

- `projects/agents/skills/repo-blender/SKILL.md`
- `projects/agents/skills/repo-blender/evals/README.md`
- `projects/agents/skills/repo-blender/evals/cases.yaml`
- `projects/agents/skills/repo-blender/references/execution-modes.md`

In this rebase, stage 2 is the upstream/base (`HEAD`) side and stage 3 is the
replayed feature commit (`9fabf532`). The working files contain ordinary
conflict markers exactly around the differences summarized below.

## Finding

Stage 2 is a strict semantic superset of stage 3 in all four files. A direct
stage-3-to-stage-2 diff contains only additions or a qualification; stage 3
contains no unique guidance. Therefore the intended semantic union is exactly
the stage-2 content, not a hand-interleaved version.

The upstream additions consistently introduce one narrowly bounded exception:
an already-installed Flatpak Blender may be used only when the user explicitly
requests it as a disposable, persistent live MCP mutation host. The
repository-pinned Blender remains mandatory for batch work, reproducible
verification, rendering of record, and deliverables.

There is no contradiction in stage 2: the toolchain prohibition is explicitly
qualified to batch and deliverable work, while the exception is limited to
interactive mutation and requires pinned clean-reopen and pixel gates. The
stage-3 sentence banning Flatpak without qualification is obsolete because it
predates that explicit exception. Keeping both variants would create an actual
ambiguity.

## File-by-file resolution

| File | Stage-2 addition | Recommendation |
| --- | --- | --- |
| `SKILL.md` | Qualifies the Flatpak ban and links the explicit live-host exception. | Take stage 2 verbatim. |
| `evals/README.md` | Declares that the suite covers the narrow exception. | Take stage 2 verbatim so the suite description matches policy. |
| `evals/cases.yaml` | Adds the regression case for the exception and its safety gates. | Take stage 2 verbatim; dropping it would leave new policy unevaluated. |
| `references/execution-modes.md` | Defines the complete exception, version check, source-compatibility boundary, immutable snapshots, pinned verification, and cleanup. | Take stage 2 verbatim. |

## Exact recommended resolution

While the same rebase remains paused, select stage 2 (`ours`) for all four
paths, stage them explicitly, and continue only after checking that no conflict
markers remain:

```sh
git checkout --ours -- \
  projects/agents/skills/repo-blender/SKILL.md \
  projects/agents/skills/repo-blender/evals/README.md \
  projects/agents/skills/repo-blender/evals/cases.yaml \
  projects/agents/skills/repo-blender/references/execution-modes.md
git add -- \
  projects/agents/skills/repo-blender/SKILL.md \
  projects/agents/skills/repo-blender/evals/README.md \
  projects/agents/skills/repo-blender/evals/cases.yaml \
  projects/agents/skills/repo-blender/references/execution-modes.md
```

Then validate the resolved skill package using its owning Bazel targets and
the repository Buildifier test if BUILD/Starlark files in the overall rebase
changed. Do not use stage 3 or delete only the marker lines: either would
silently discard the newer upstream exception and its regression coverage.
