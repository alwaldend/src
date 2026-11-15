## Decision

Proceed with stable ansible-core 2.22.0 and the matching Ansible 15.0.0rc1 bundle. Core 2.22 supports target Python 3.9–3.15 and controller Python 3.13–3.15. The matching community bundle is still a release candidate as of 2026-11-15; Python 3.15 compatibility is the task-specific reason for using it. Keep unrelated dependency versions unchanged where resolution permits.

The alternative of installing a separate host Python would add a second runtime and lose access to system Python's native SELinux bindings. Patching Ansible's private dataclasses integration would maintain an unsupported fork. Updating the existing tooling preserves normal host module execution and keeps pins with their current owner.

Sources: [core releases](https://pypi.org/project/ansible-core/), [community releases](https://pypi.org/project/ansible/), [support matrix](https://docs.ansible.com/projects/ansible-core/2.22/reference_appendices/release_and_maintenance.html).

## Acceptance

Regenerate the lock through uv_lock.update, inspect package changes, validate representative packaged playbooks and lint, and repeat live fact gathering on secure1 and dev1 with strict host verification. Preserve the initial Python failure as sanitized evidence. Do not proceed with host configuration if the replacement tooling fails these checks.

The XO appliance preflight exposed a launcher gap: native SSH askpass executes
`ansible`, which was missing from PATH. Package the pinned adhoc CLI with every
Ansible wrapper beside its generated entry points. Verify
through the owning read-only appliance inspection before live certificate work.

The first real connection additionally showed the upstream askpass dispatcher
requires basename ansible at import time. A small typed entry point restores
that name before importing the CLI; both normal adhoc operation and native
askpass continue to use the pinned upstream implementation.

The sibling executable layout resolves native askpass without changing PATH.
An isolated packaged SSH transport verified shared-memory password delivery and
module execution with PATH restricted to /usr/bin:/bin. A fresh live XO attempt
was refused on port 22 before authentication; it does not establish live acceptance.
