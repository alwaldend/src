## Context

See proposal.md for motivation. The user explicitly requested: "Halve the limit for now. Dont run ansible to set it, just modify the host in-place".

## Decisions

Apply the exact requested imperative exception only to the two scheduling lines in `/var/home/simeonwarrenbot/.bazelrc`, and reconcile checked-in Ansible desired state in this change. Do not run a playbook or restart existing servers. New invocations read the updated configuration.

## Risks / Trade-offs

Lower concurrency can increase build duration. Independent servers still have independent budgets; this change does not guarantee safety for arbitrary build concurrency.

## Migration Plan

The live configuration was changed directly on 2026-10-11 from eight jobs and `HOST_CPUS*0.5` to four jobs and `HOST_CPUS*0.25`. Verify the live file against checked-in desired state and inspect Bazel's parsed options. A future rollback would restore the two original settings in both locations.
