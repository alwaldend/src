# Printer firewall access

This Terraform root defines one IPv4 forwarding exception on router1. Clients
in `192.168.2.0/24` arriving on wireless `bridge2` may reach the wired printer
through `bridge1`. The destination is the single IPv4 address from `printer1_a` in the
owner's canonical [DNS declaration](../dnsconfig.json). All IPv4 services on
that device are included; no other wired destination matches this rule.
The bare address matches RouterOS's canonical representation of a host destination.

The rule is inserted before the enabled forward drop with comment `drop forward`.
Planning fails if that insertion point is absent or ambiguous. The root does
not manage existing firewall rules, router exports, DNS, or a print web app.
Existing connection tracking handles replies. The printer must retain its
declared address and have a working return route through `192.168.1.1`.

## Credentials and state

The `tf=firewall` wrappers select the existing `src_infra_mikrotik` AppRole and
inject `routeros_hosturl`, `routeros_username`, and `routeros_password`. The
endpoint comes from the shared infrastructure AL configuration. The injector
reuses the existing non-DNS RouterOS API credential at
`secrets/alwaldend.com/vault1/approles/src_infra_ingress/mikrotik`, with fields
`mikrotik_username` and `mikrotik_password`. The
[paired Vault policy](../../vault/tf/approles/src_infra_mikrotik/main.tf) grants
the MikroTik AppRole read access to that exact secret. It does not assume the
ingress identity, copy the credential, or reuse the DNS account. Ingress retains
ownership of the account and its rotation; the account's router permissions
are shared between workflows.

State uses the owner's `tf_backend/firewall` Vault HTTP backend, independently
of the DNS root's `tf_backend/tf`. The backend plugin initializes its own
entries. These definitions do not create a RouterOS account or store credentials.
Deploy the named read policy through the owning Vault Terraform workflow before
planning the firewall. Policy and firewall applies each require explicit live authority.

## Validation

```sh
bazel test //infra/mikrotik/firewall:tf_tests.fmt_test \
  //infra/mikrotik/firewall:tf_tests.validate_test \
  //infra/dns:config_test
```

These checks validate formatting and the configuration against the pinned
provider schema without authentication or router access. The lookup is inline;
no module or runtime input exists solely to support mocks. The precondition
rejects a live lookup that does not return exactly one enabled terminal drop.
Offline schema validation does not execute that lookup or its failure cases;
review the live plan before applying. Test packaging omits only `backend.tf`
so initialization does not need the operational Vault backend.

## Deployment and acceptance

With explicit authority for the live operation, review normal text plans:

```sh
bazel run //infra/vault/tf:tf.plan -- -target=module.src_infra_mikrotik
bazel run //infra/mikrotik/tf:tf.plan
bazel run //infra/mikrotik/firewall:tf.plan
```

The scoped Vault plan must add the named read policy and append it to the
existing MikroTik AppRole, preserving its current policies. It does not audit
unrelated Vault resources. Its existing DNS policy also gains the named
Cloudflare-secret read needed for global DNS. Deploy these grants through the same owning
Vault apply target and scope, with explicit live approval.

The printer A record is declared in both dc1 and global DNS. The existing local
record should remain unchanged; the global record is unproxied and points to
the same private address. Public resolution does not provide network access.
Expect the global DNS addition and the single firewall exception in their
respective roots. Investigate any unrelated change, replacement, or deletion. Apply
each root with its owning `tf.apply` target and review the new interactive
plan before confirming. Never import an entire router export to deploy this rule.

After applying, verify the rule precedes the forward drop. From an actual
Wi-Fi client using the site DNS resolver, run:

```sh
nslookup printer1.dc1.alwaldend.com 192.168.2.1
curl --connect-timeout 5 --max-time 10 -I http://printer1.dc1.alwaldend.com/
```

Confirm the name resolves to the declared printer IP and its web interface
opens, then submit a test job with that printer's supported protocol. These
checks require a powered-on printer and its actual services. Retain a sanitized
deployment and acceptance receipt with the change; offline success alone does
not establish connectivity. DNS and routed access do not relay mDNS/AirPrint
discovery between the subnets.

Rollback by reverting the printer DNS declaration and removing this exception
through their owning Terraform workflows. Existing tracked connections may
remain until they expire. Manual rule reordering can change behavior; verify
the insertion order after any manual firewall changes.
