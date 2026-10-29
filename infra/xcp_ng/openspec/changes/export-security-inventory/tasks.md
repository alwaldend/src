## Inventory reuse

- [x] Expose the existing inventory with narrow security-package visibility.
- [x] Verify the original inventory export and patch whitespace.
- [x] Validate the shared consumer through the parent implementation.

The packaged `//infra/security/ansible:ansible` consumer built successfully
and its composed inventories were used for the authorized maintenance run.
Inventory membership remains owned by the source inventory.

## XCP-ng interpreter selection

- [x] Select the side-by-side interpreter at host scope in the owner inventory.
- [x] Document the user-authorized manual bootstrap exception and existing pin.
- [x] Validate the owner source diff, inventory scope, and documentation links.
- [x] Validate integrated packages, merged interpreter precedence, and strict
      OpenSpec structure offline.
- [x] Parent: review the candidate before deployment.
- [x] Parent: complete the authorized manual bootstrap without reboot and
      verify the installation and host preservation checks.
- [x] Parent: validate native Ansible modules using the selected interpreter.
- [x] Parent: validate the shared-role audit and record live evidence if successful.

Offline source validation parsed the inventory YAML and compared it with
`0ce1dbb4`: only the existing XCP-ng host's interpreter selection changed.
The version matches `tools/py/include.MODULE.bazel`; local documentation
links resolve and `git diff --check` passes.

The parent reports successful installation at the declared path. The controller
and remote SHA verification, imports, and 125-certificate CA store passed;
the vendor Python digest and boot identity were unchanged, and all four guests
remained running. Evidence is private at
`out/host-security-updates-20261029/python-bootstrap/install.log`.
The parent's native `ping`, minimal `setup`, SHA-256 `stat`, and `command`
checks passed on CPython 3.14.4, with vendor Python digest, boot identity,
and all four running guests unchanged. The private receipt is
`out/host-security-updates-20261029/python-bootstrap/verified.json`.
The initial temporary verification play's play-scoped interpreter also reached
its delegated controller receipt task; explicitly selecting
`ansible_playbook_python` for that task corrected the second run. The checked-in
inventory uses host scope; a packaged local-only canary confirmed delegated
`setup` and evidence `copy` use `ansible_playbook_python` for both appliance
hosts. Its receipts are under
`out/host-security-updates-20261029/appliance-simplification/`.
The final shared-role audit completed on 2026-10-29: XCP-ng retained a finished
245-test report with warning status 78, all artifacts mode 0600, and unchanged
boot identity. XOA completed through the same role with 257 tests. Both Ansible
recaps had zero failed or unreachable tasks. Private evidence is under
`out/host-security-updates-20261029/lynis-review/live-v2/`; the wrapper log is
`lynis-review/live-audit-v2.log`. The source-only delegate performed no live
operations; the parent performed the authorized validation.

- [x] Use the shared fleet inventory consumer package group without widening owner access; verify consumer packaging.
