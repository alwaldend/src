## Acceptance

The local inference consumer's packaged Ansible build and localhost-only
syntax check passed. Its rendered standalone unit invokes the installed
executable, and the official complete 0.1.6 bundle served client assets and
native pairing/session endpoints in a disposable network-isolated fixture.
See the consuming change's `validation.md` under
`infra/local_inference/openspec/changes/add-t3code-mtls/` for the shared evidence.

Source inspection confirms `force: false` preserves an existing settings file,
changed startup inputs notify the restart handler, and steady state uses
`state: started`. Environment writes use `no_log: true`; standalone mode omits
the unsupported legacy token variable. Actual user creation, extraction,
systemd activation, and repeat deployment on a VM were not performed.

The follow-up removed the executable-existence extraction shortcut. A local
fixture using the checked-in extraction task's module arguments repaired a
partial synthetic executable and missing support library, then reported no
change when reapplied. This verifies archive repair/idempotency without
creating service users or activating systemd on a host.
