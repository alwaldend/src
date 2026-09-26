# Role acceptance matrix

This matrix precedes role changes. Every scenario applies its real packaged
role independently, using only disposable guest prerequisites.

| Role | Fresh and unchanged applications | Update and persistence | Exclusions |
| --- | --- | --- | --- |
| host | Administrator sudo/reconnect, regular user, allowed/denied traffic, UTC, test CA, enabled SSH/OS hardening, nonempty LVM/filesystem/mount; second application reports zero changes | Add a regular user through updated role inputs; reconnect after hardening and reboot; managed mount and marker survive | External SSH certificate issuance only (`ssh_sign_key`); locally signed host certificates remain active |
| traefik | Real enabled systemd service, trusted TLS to local backend, unmatched route rejected; unchanged application reports zero changes and unchanged lifecycle journal/timestamps | Change route, observe new route and reject old route; repeat idempotence; restart retains routing | Vault EAB/public ACME (`traefik_disable_eab: true`) |
| forgejo | Real enabled systemd service, HTTP, disposable local account/repository, known Git commit pushed and freshly cloned; unchanged application reports zero changes and unchanged lifecycle journal/timestamps | Supported application configuration becomes visible, repeat idempotence; restart preserves account/repository/commit without recreation | External database, OIDC |

Retain named assertions and non-secret service diagnostics through the shared
runner. Never retain passwords, tokens, private keys, or full temporary trees.
All scenarios exclude production inventory, full deployment playbooks, and
deployment-specific Xen disk assertions. Failure, interruption, and successful
completion must destroy the guest. No blanket idempotence exclusions or false
`changed_when` results may conceal role behavior.
