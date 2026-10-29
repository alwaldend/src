## Decisions

XOA transport cleanup must surround all work after connection allocation, not
only the updater. A failed fact query, Debian assertion, or boot comparison
must still close and remove the invocation's connection. Use nested native
`always` blocks and skip cleanup when no private directory was allocated;
preserve the original failure and do not reset an unrelated caller connection.

Both appliances use `alwaldend.main.lynis_audit` and the owner's installed
interpreter. Ansible owns runtime compatibility errors; no raw version matrix,
additional Python bootstrap, optional collector, or duplicate Lynis transport
is maintained. The shared role owns private result allocation, scanner evidence,
failure reporting, and remote scratch cleanup. Keep appliance boot verification.

Updates remain separately selected operations. XOA uses the unified updater
with automatic updates enabled and reboot disabled; native APT dependency
installation remains disabled. XO application automation is outside OS updates.
The extra appliance_updates_approved flag is removed; selecting a mutation
target does not authorize an agent to deploy it without the user's instruction.

Retain XCP-ng's single-host pool/coordinator, HA, task, XOSTOR, attached-media,
and yum-cron checks. Preserve placement disable/restore, the vendor toolstack
restart and readiness checks, and unchanged boot identity. Generic DNF modules,
RPU, unattended yum scheduling, and guest power operations are excluded.

Vendor yum retains the two stable repository selectors and explicit global and
repository-specific signature/TLS checks. The optional origin override keeps
the demonstrated baseurl-only guard through vendor Python's native yum loader;
a modern Ansible interpreter does not supply vendor YUM bindings.

## Transport evidence

Keep the raw command status sentinels and exit-100 guards. Pinned ansible-core
2.20.7 `command` failures call `fail_json`, which exits 1; the password SSH
plugin treats sshpass exits 1–6 as connection/authentication failures. Native
module JSON alone therefore does not establish equivalent failure handling.
Source evidence is retained in
`out/host-security-updates-20261029/simplification/updater/transport-source-evidence.md`.
No live module canary was executed.

## Validation

Failure cases were recorded before edits in the same private directory.
Validate the packaged candidate, task selection, private result paths, boot
comparison, and update failure/recovery boundaries offline. Previous successful
live runs remain historical evidence; the parent owns any newly authorized
canaries. Scanner completion never establishes absence of compromise.

Vendor references: [updates](https://docs.xcp-ng.org/management/updates/),
[CLI](https://docs.xcp-ng.org/appendix/cli_reference/), and
[XOA application updater](https://docs.xen-orchestra.com/updater/).
