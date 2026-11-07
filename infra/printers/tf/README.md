# Printer network access

This root defines the printer1 IPv4 forwarding exception on router1. Clients
in `192.168.2.0/24` arriving on Wi-Fi `bridge2` may reach the printer through
wired `bridge1`. The destination comes from `printer1_a` in the canonical
[declaration](../dnsconfig.json). All IPv4 services on that host are included;
other wired destinations do not match. The printer needs stable addressing
and a return gateway through `192.168.1.1`.

The rule precedes the enabled forward drop with comment `drop forward`.
Planning fails unless exactly one such drop matches. The existing rule's
`tf[infra/mikrotik/firewall]` comment is retained as its deployed identifier,
avoiding an incidental router update during state adoption. No router export
or unrelated firewall rule is applied. Existing connection tracking handles
replies; manual rule reordering requires another order check.

`place_before` is creation-only in the pinned provider and is not returned by
RouterOS on import. Its updates are ignored after creation so adopting an
existing rule does not replace it. New rules still use the guarded insertion
point; verify actual ordering after adoption or manual firewall changes.

## Credentials and state

The `tf=main` wrapper uses the `src_infra_printers` AppRole and its
`tf_backend/tf` backend. It injects the RouterOS endpoint from shared infra AL
configuration and the dedicated `mikrotik_username` and `mikrotik_password`
fields from `secrets/alwaldend.com/vault1/approles/src_infra_mikrotik/printers`.
The [MikroTik accounts root](../../mikrotik/accounts/README.md) owns account
creation and rotation; the printers AppRole has only read access to that secret.
Its [Vault definition](../../vault/tf/approles/src_infra_printers/main.tf) also
grants the DNS credentials used independently by the DNS root.

The [state transfer](../docs/state-migration.md) is complete; the existing rule
retains its provider identity.
Operational targets are `//infra/printers/tf:tf.<operation>` and require explicit
live authority.

The shared injection helper supplies `MIKROTIK_HOST`, `MIKROTIK_USER`, and
`MIKROTIK_PASSWORD` directly to the default RouterOS provider. DNS retains
its separate aliased provider inputs.

## Verification

```sh
bazel test //infra/printers/tf:tf_tests.fmt_test \
  //infra/printers/tf:tf_tests.validate_test
```

These checks validate formatting and provider schema offline, excluding only
the backend. They do not execute the live drop lookup or verify connectivity.
After an authorized migration, expect a no-change plan and verify the existing
rule still precedes the terminal drop. From a real Wi-Fi client, run:

```sh
nslookup printer1.dc1.alwaldend.com 192.168.2.1
curl --connect-timeout 5 --max-time 10 --location --output /dev/null \
  --write-out '%{http_code}\n' http://printer1.dc1.alwaldend.com/
```

Confirm the declared address resolves and the printer page opens, then submit
a test job through its supported print protocol. Retain sanitized acceptance
evidence. Routed access does not relay mDNS or AirPrint discovery across subnets.
