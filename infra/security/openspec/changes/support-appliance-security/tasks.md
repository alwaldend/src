## Source simplification

- [x] Record failure cases before source edits.
- [x] Remove the optional collector, historical runtime survey, and redundant interpreter probes.
- [x] Remove the extra update-approval flag; preserve explicit mutation targets and no-reboot policy.
- [x] Inspect pinned password-transport failure semantics and retain raw vendor status handling.
- [x] Use the shared Lynis report lifecycle while preserving boot evidence and existing connection preparation.
- [x] Validate the exact candidate offline and return the remaining live acceptance boundary to the parent.

## Evidence

## Complete Ansible review

- [x] Ensure XOA connection cleanup runs after evidence preparation, fact gathering, updater failure, or a failed boot comparison, while touching only allocated transport.
- [x] Exercise failure cleanup and retained boot evidence offline; validate all changed appliance playbooks and owning records.

## Retained evidence

Current source evidence belongs under
`out/host-security-updates-20261029/simplification/updater/`.
No live operation, authentication, Git command, new dependency, or captured
runtime modification is authorized to this child. XCP-ng vendor preflight,
repository verification, mirror selection, placement recovery, and native
no-reboot behavior remain intact.

On 2026-10-30, packaging, all four appliance playbook syntax checks, source YAML
and vendor shell parsing, and strict owner OpenSpec validation passed. Shared
Lynis result allocation replaces the nested appliance audit directory. Boot
records and the existing XOA connection recovery path are added after the role
returns; connection cleanup remains in an always block. Interrupted execution
does not establish complete boot evidence. Authentication policy is unchanged.
The parent owns combined formatting/package verification and any authorized
runtime acceptance. No synthetic boot-task execution or live audit was run by
this child; only parsing and task-selection checks were executed.

Earlier shared-role audits completed on both appliances: XCP-ng retained 245
tests and XOA 257, both with warning status 78, private artifacts, and unchanged
boot identities. Receipts remain under
`out/host-security-updates-20261029/lynis-review/live-v2/` and
`lynis-review/live-audit-v2.log`. Earlier package/syntax, interpreter, controller
canary, and scoped formatting receipts remain under `appliance-simplification/`.
Those runs predate this candidate.

Earlier authorized XOA OS and XCP-ng vendor updates retained boot identity and
verified host placement restoration. Native repository-option verification is
recorded under `clamav/pr-simplification/`; operational receipts remain under
`appliances/`, all within the same original private task directory. The owner
inventory change records the existing Python bootstrap exception. Hypervisor
unattended updates and XO application automation remain unsupported.

The follow-up fixture executed the candidate XOA playbook with its actual
evidence, boot comparison and cleanup tasks in a private user namespace.
Preparation, facts, distribution, updater, boot-mismatch and success cases
passed: allocated connection directories were removed and failures stayed
failures; missing allocations skipped connection reset. Facts and the updater
were controlled fixture boundaries, so no OS packages or SSH policy changed.
The initial boot fault injector did not reach the built-in action; the retained
retry supplied an explicit mismatched original identity and exercised the
comparison failure. Evidence: `ansible-review/validation-summary.json` under
the original private task root.
