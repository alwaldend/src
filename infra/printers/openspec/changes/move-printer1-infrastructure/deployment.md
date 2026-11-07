# Authorized deployment

Completed on 2026-11-07, with final scoped convergence checks at 18:11 UTC.
The user authorized apply after checking the plans, then supplied the
MikroTik-owned administrator credential and requested removal of the ingress
dependency. Operations used owning repository targets and interactive apply
approval. This receipt contains no credentials, raw state or saved plans.

The initial published source was daf2e4dd716f8e9d99daa0cbdb2e07db7123d41f.
Deployment also used the administrator-secret reference, REST API permission
correction and import-placement handling recorded in this change's source.

## Provisioning

| Completed (UTC) | Owning workflow                                                                      | Result                                                                                                                    |
| --------------- | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------- |
| 18:01:48        | `//infra/vault/tf:tf.apply`, targeted to the two owner modules and DNS AppRole group | 10 added, 2 changed, 1 destroyed; provisioned printers access and removed only the obsolete MikroTik ingress-read policy. |
| 18:02:29        | `//infra/mikrotik/accounts:tf.apply`                                                 | 4 added, 0 changed, 0 destroyed; group `*7`, user `*6`, stable generated password and Vault credential.                   |
| 18:05:15        | `//infra/mikrotik/accounts:tf.apply`                                                 | 0 added, 1 changed, 0 destroyed; added `api` to the group alongside `rest-api`, without rotating the password.            |

The scoped Vault targets were `module.src_infra_printers`,
`module.src_infra_mikrotik` and `vault_identity_group.dns_approles`. This does
not establish convergence of unrelated Vault configuration.

RouterOS rejected REST requests with only `rest-api`. Its documented
[REST API permission requirement](https://manual.mikrotik.com/docs/developer-guides/rest-api/)
requires both `api` and `rest-api`. The group now has these plus `read` and
`write`; it has no account-policy administration or sensitive-value access.
The separate administrator credential remains in the MikroTik owner's Vault
prefix. Printers receives only its generated credential.

## Existing object adoption

Each transfer removed only the verified old state binding before importing
the exact existing object through the new owner's target. No DNS record or
firewall rule was changed, deleted or recreated.

| Resource address                                                                | Preserved provider ID              | Import completed (UTC) |
| ------------------------------------------------------------------------------- | ---------------------------------- | ---------------------- |
| `module.dns.module.dc1[0].routeros_ip_dns_record.records["printer1_a/A/dc1"]`   | `*8E`                              | 18:04:32               |
| `module.dns.module.global.cloudflare_dns_record.records["printer1_a/A/global"]` | `ac55acb4eae5f5ddf980ca787196b6d5` | 18:05:42               |
| `routeros_ip_firewall_filter.wifi_printer`                                      | `*29`                              | 18:06:31               |

DNS moved from the MikroTik `tf_backend/tf` backend to printers
`tf_backend/dns`. The firewall moved from MikroTik `tf_backend/firewall` to
printers `tf_backend/tf`. Both DNS views retain `192.168.1.221`, TTL 300 and
the existing name; the global record remains unproxied.

The former firewall backend's resource and data-source bindings were removed.
Its empty address listing was verified at 18:07:11 UTC. The temporary
state-only migration target and AL backend selector were then removed.
Historical protected backend entries were not deleted.

Import does not restore the provider's creation-only `place_before` field.
An initial plan proposed replacing the rule solely for that metadata; it was
not applied. Source now ignores updates to this field and retains the guarded
insertion point for new rules. Read-only API inspection at 18:08:36 UTC
confirmed rule `*29` precedes the enabled terminal drop `*16`, with the
expected source, destination, interfaces and accept action.

## Acceptance

Normal text plans reported **no changes**:

| Target                                                          | Observed (UTC) |
| --------------------------------------------------------------- | -------------- |
| `//infra/printers/dns:tf.plan`                                  | 18:06:22       |
| `//infra/printers/tf:tf.plan`                                   | 18:08:40       |
| `//infra/mikrotik/tf:tf.plan`                                   | 18:08:50       |
| `//infra/mikrotik/accounts:tf.plan`                             | 18:11:46       |
| `//infra/vault/tf:tf.plan`, with the provisioning targets above | 18:11:59       |

Vault capability checks at 18:09 UTC returned `deny` for the MikroTik
AppRole reading the ingress credential and for printers reading
`routeros_admin`. The dedicated printers account authenticated successfully
for its firewall plan and the rule-order inspection. The retained MikroTik
DNS root has no printer bindings and preserves its other records.

Read-only host checks resolved the printer name to `192.168.1.221` and fetched
its status page with HTTP 200. The printer rejects HEAD requests with HTTP
411, so a GET is the appropriate web reachability check. No printer settings
were changed. A new Wi-Fi client test and print job were not performed in
this deployment; cross-subnet service discovery remains outside this change.
