# Development DNS

This root manages every record in the canonical [declarations](../dnsconfig.json):
secure1, dev1 and dev1's services. It uses the shared `src_infra_dev` AppRole,
DNS-only provider credentials and one Vault HTTP backend:
`alwaldend.com/vault1/approles/src_infra_dev/tf_backend/dns`.

Validate offline with `//infra/dev/dns:tf_tests.fmt_test`. Operational entry
points are `:tf.plan` and interactive `:tf.apply`. Source validation does not
authenticate, transfer state or apply DNS.

## Pending state handoff

The merged root has not been adopted live. Do not apply it against an empty
backend without reconciling existing records. The prior backends remain intact:

| Legacy AppRole                    | Backend path below `alwaldend.com/vault1/approles/` |
| --------------------------------- | --------------------------------------------------- |
| `src_users_simeonwarren_host_bot` | `src_users_simeonwarren_host_bot/tf_backend/tf`     |
| `src_infra_forgejo_runner`        | `src_infra_forgejo_runner/tf_backend/tf_setup`      |

Under separate authorization, follow the [shared handoff procedure](../../dns/README.md#migration-and-recovery):

1. Provision the shared Vault identity and DNS read policy. Freeze writers using
   both old roots and record the source revisions and backend ownership.
2. Inspect both legacy states and fresh provider inventories through their
   owning workflows. A previous runner-backend initialization referenced Xen
   Orchestra resources; its complete contents remain unverified. Preserve and
   reconcile unrelated resources separately. Do not assume that state is DNS-only
   or that the absent Proxmox runner requires a live deletion.
3. Match existing records by exact provider identity and reviewed name mapping.
   Adopt DNS resources into the new state and release old ownership without
   destroying records. The temporary `dns_routeros_import_ids` and
   `dns_cloudflare_import_ids` maps accept normalized keys such as
   `dev1_dc1/A/dc1`; Cloudflare values use the provider's import-ID syntax.
   The former logical keys `host_bot_dc1` and `host_bot_dc1_target` are now
   `dev1_dc1` and `dev1_dc1_target`. Other current keys are unchanged.
   Never infer a secure1 record from the obsolete runner1 declaration.
4. Separate ownership adoption from the approved host and service DNS renames.
   Review the normal text plan for imports, additions, changes and deletions;
   reject duplicate ownership and destructive recreation as migration effects.
5. After adoption, remove the temporary import maps and blocks, verify the
   resulting no-change plan, and retire legacy identities only after confirming
   that no remaining state or credential consumer needs them. Keep DNS enabled.

No step above has run as part of this source change. Rollback after a handoff
requires freezing writers and reconciling ownership; reverting source alone
does not restore the previous states.
