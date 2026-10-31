## Context

The collection packages reusable roles through `al_ansible_role`; consuming
playbooks own host selection and data paths. CISOfy publishes Lynis as shell
source in a tar archive with a `lynis/` directory.

## Goals / Non-Goals

Provide a bounded on-demand root audit and private evidence. Implementation
validation does not authorize remote execution or authentication. Host playbook
wiring and unrelated package-maintenance changes belong to the parent task.

## Decisions

- Use `http_file` and a public-to-repository `archive` alias to preserve the
  official tarball and its license instead of compiling or repackaging scripts.
- Use GNU `timeout`, a fresh root-owned temporary directory, and Lynis's own
  command options. Preserve stock upstream source and its private `/tmp/lynis.*`
  files; the user explicitly accepted these on remote hosts.
- Require an existing controller-owned 0700 destination. Create a unique result
  directory per invocation to avoid concurrent-host and repeated-run collisions.
- Reuse Ansible's private temporary-directory creation and direct `unarchive`
  transfer. Define receipt fields in one role-local template and render it at
  the existing recovery and collection checkpoints.
- Write an incomplete private controller receipt before remote work and save
  the allocated remote recovery path before preparing the scanner. Keep command
  transport errors in that receipt, including when no exit code is returned.
- Collect available evidence in an `always` block. Ignore individual stat/fetch
  errors only to attempt the other files; their aggregate failure prohibits
  cleanup and fails the role. Clean only role-owned remote scratch after all
  existing artifacts and the collection receipt have been retained. A terminated
  scanner can leave its own private temporary files; never use broad `/tmp`
  cleanup patterns.
- Require an accepted exit (0 or 78), successful collection, `finish=true`, and
  a positive `lynis_tests_done` in the retained report for completion. Preserve
  partial files on all failure paths; exit status alone is insufficient.
- Keep `audit system` with all groups and `--cronjob`; omit redundant `--quick`,
  which only suppresses input pauses. Share the DNF and YUM cache-mutation
  exclusions in `audit.prf`; the role README lists their exact coverage loss.
- Require remote Python supported by the packaged Ansible runtime. Appliance
  Python provisioning and removal of legacy transport belong to the parent.

## Risks / Trade-offs

- Audit evidence contains host details: restrict access and do not log contents.
- A timeout can produce incomplete reports: preserve them with the exit code.
- Upstream checks observe live system state: a package build verifies packaging,
  not live audit success. Report that validation boundary explicitly.

## Native simplification acceptance

Run Lynis with native async polling, remove its stdout exit protocol, and retain scratch when terminal status is unknown.
Failure cases: warnings, timeout, missing or unfinished native status, incomplete
reports and failed evidence transfers. Preserve partial evidence and recovery
paths; remove only terminated scanner scratch after successful retention.
Production scanning is stopped; validate only bounded fixtures.
