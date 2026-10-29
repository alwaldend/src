## 1. Reusable workflow

- [x] 1.1 Add scoped operator authentication and strict SSH trust.
- [x] 1.2 Compose existing owner inventories without copying host lists.
- [x] 1.3 Run pinned Lynis and enforce enabled automatic updates with reboots disabled.

## 2. Validation and execution

- [x] 2.1 Validate the packaged playbook and an audit/update canary with unchanged boot identity.
- [x] 2.2 Complete active-host audits and authorized updates, including appliance transport or explicit limitations.
- [x] 2.3 Verify service health, automatic-update configuration, installed/running kernels and unchanged boots.

Exact-candidate source gates and separate master-based PR publication are
tracked by the delivery receipts, which own publication state.

## 3. PR simplification

- [x] Put audit and updater roles in a block with unconditional boot checks in
      `always`, preserving role tags and the original failure.
- [x] Reuse the injected maintenance SSH configuration in the standalone SSH
      target, preserving its injector dependency and strict trust policy.
- [x] Validate source/package identity, pinned Ansible syntax and offline
      failure/tag behavior, and packaged SSH arguments without authentication.

Failure cases recorded before implementation: audit or updater failure skips
`post_tasks`, concealing a changed boot; successful cleanup can accidentally
rescue the original failure; inherited block tags can run an unselected role;
transport loss can be mistaken for a verified unchanged boot. Test success,
audit failure, update failure, changed boot, and unreachable outcomes, plus
audit-only and update-only selection. Duplicate SSH arguments can drift from
the generated configuration or lose literal environment expansion; verify
the generated wrapper passes `-F` with `HOST_MAINTENANCE_SSH_CONFIG` and retains
the maintenance plugin selector. Final review and delivery remain with the
parent; this task performs no live operations, commits, or push.

The parent authorized changing this target's rule/load to the existing
`al_binary_run` after the literal `-F` expansion failure. Before implementation,
extend offline validation to forwarded `-oHostKeyAlias`, arguments containing
spaces or shell metacharacters, missing/empty configuration variables, and
unchanged host-PATH SSH selection. Preserve AL configuration, injector data,
and the maintenance selector; introduce no new rule or dependency.

## Complete Ansible review

- [x] Centralize the repeated fleet inventory consumer permission in one named package group, preserving exact existing consumer access.
- [x] Validate all changed Ansible sources, role callers, inventory exports, and packaged workflows.
