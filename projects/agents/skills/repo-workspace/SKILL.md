---
name: repo-workspace
description: >-
  Establish an isolated task workspace before the first mutation or
  task-scratch write. Use to verify checkout isolation, select or reuse a
  dedicated feature branch in its own linked worktree, keep task downloads,
  reports, logs, caches, and temporary files under ignored workspace scratch,
  and preserve unrelated work. It does not own secret handling, delivery,
  layout, or versioning.
---

# Prepare a task workspace

Confirm the workspace is isolated before the first mutation or task-scratch
write. Inspection alone does not authorize a change; the task owns that
authority.

## Verify checkout isolation first

1. Identify the current checkout before writing: record the worktree root, the
   current branch, and both Git directories.

   ```sh
   git rev-parse --show-toplevel --abbrev-ref HEAD --git-dir --git-common-dir
   ```

2. Require a dedicated feature branch in its own linked worktree for any task
   that can modify repository files or task-owned scratch. Do not create
   another worktree or switch branches when the task already supplies one.
3. Query the current checkout before requesting a full inventory. Only when the
   current checkout is ambiguous, or the task requires placement elsewhere,
   inspect bounded state with `git worktree list` and `git rev-parse`.
4. Reuse the existing task worktree and branch instead of creating duplicates.
   Reconfirm isolation after an interruption, when the branch, candidate, or
   authority may have changed.
5. Treat the default branch and default checkout as read-only unless the user
   explicitly authorized that exact task there. When isolation is missing and
   no authorized location exists, ask instead of mutating the default checkout.
6. Preserve unrelated work. Never discard, auto-stash, commit, or rewrite
   shared, human-owned, unrelated, or ambiguous history while preparing a
   workspace.

## Keep scratch under ignored out

- Keep every task-owned download, report, log, cache, extracted archive,
  temporary file, and tool input or output under ignored `out/<task>/` in the
  applicable workspace. The repository ignores `out/`; see `.gitignore`.
- Point configurable temp and cache variables at that directory too, including
  `TMPDIR`, `TMP`, and the tool-specific cache variables a task's tools honor.
- Never write task scratch to `/tmp`, `/var/tmp`, or another shared system
  location unless the tool cannot use workspace scratch at all. State that
  reason and remove only your own residue before handoff.
- Never delete unrelated temporary files, including those created by other
  tasks, users, or agents.
- Keep secret-bearing temporary material task-private, access-restricted,
  short-lived, and explicitly cleaned up; never track, stage, commit, or
  promote it as evidence. "Outside the repository" means outside tracked or
  committable source, which ignored task scratch satisfies. Follow
  `repo-secrets` for what counts as secret-bearing and how to handle it.

## Clean up on handoff

- Remove only task-owned scratch residue under the task's `out/<task>/`
  directory, and only once its evidence is no longer needed.
- Leave unrelated files, worktrees, and branches alone.
- Delivery, commit, push, and pull-request state belong to `repo-delivery`;
  repository layout and path selection belong to `project-layout`.
