---
name: repo-delivery
description: >-
  Finalize task-owned repository changes by verifying, committing, pushing the
  feature branch, and creating or updating its pull request. Use when an
  implementation is ready for delivery; do not use for read-only reviews.
---

# Deliver repository changes

Treat the root `AGENTS.md` as the policy source; this skill supplies the
delivery mechanics.

1. Inspect the branch, status, and task-scoped diff. Preserve unrelated and
   pre-existing changes. Stage explicit paths or safe task-owned hunks, never a
   blanket worktree.
2. Run the relevant validations and `git diff --check`, then review the exact
   staged diff.
3. Confirm the branch is a feature branch. Commit all task-owned changes with
   the required LLM footer, without amending or rewriting existing history.
4. Push normally to at least one configured remote. Stop if the push fails.
5. Detect the forge from that remote's URL and use its authenticated
   integration. Use GitHub tooling for a GitHub remote and the repository's
   `//tools/fj` wrapper for Forgejo. Never create a duplicate pull request on a
   different forge.
6. Find the pull request for the exact head branch. Update only task-owned
   summary and verification information while preserving unrelated human
   content, or create the pull request when none exists. Add the required LLM
   footer.
7. Verify that the remote branch contains the commit and no requested changes
   remain uncommitted.
