# Guidance verification

The installed and repository-pinned CLI both identify as fj 0.7.0. Help was
inspected without invoking the Vault wrapper for the root, PR, create, view,
review, review-list, auth, login and token-import commands. It confirms explicit
worktree/remote selection, body-file creation, review/inline-comment listing and
stdin token import. It exposes no structured thread-reply/resolution workflow.

The pinned `src/auth.rs` was read through the official Codeberg contents API.
It confirms the custom client-ID file format, user/system search order and
`http://127.0.0.1:26218/` PKCE callback. The related host change records the
actual missing-client login result and user-selected token authentication.

The skill package and its validation aspect passed. The existing offline
Promptfoo configuration, including four new remote-routing scenarios, passed
`@rules_al//skills/al-delivery:eval_config_test`. This verifies staging and
configuration loading, not live agent behavior or remote PR mutation support.
Owner-local OpenSpec validation passed.

The required discovery check initially found a stale generated `al-bazel`
Host Bot link. `//.agents:write_skills` refreshed the archived projections,
including this skill, and its exact-state test then passed. The updater was
run twice using unchanged maintained inputs. Canonical skill source remains
the owner; no discovery copy was edited manually.

Final quality, semantic lint, complete-candidate review and publication are
bound to the task-local preparation and delivery receipts under ignored
`out/simplify-forgejo-runner/`. The existing PR remains on GitHub because its
verified origin and existing destination are GitHub; no mirror was selected.

The aggregate verification selected 34 passing tests: repository quality,
runner acceptance, the affected owner-local OpenSpec checks, offline skill eval
configuration and discovery. Affected semantic lint and skill validation passed.
The complete source/configuration/documentation diff was reviewed locally with
no actionable findings. The final candidate is rechecked through normal Bazel
commands before publication; its exact head remains in the local receipt.
