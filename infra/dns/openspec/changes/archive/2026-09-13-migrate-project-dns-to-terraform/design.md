## Context and outcome

The baseline at `cb4b5bd3a37d29dc778ed140ab72d05efef7816a` declares 168
records: 83 global and 85 dc1. Canonical declarations stay with their owners;
shared apex/mail records move from JavaScript into `infra/dns/dnsconfig.json`.

The user rejected both the exporter and the central ownership configuration.
The final implementation discovers files at runtime and prints their domain
ownership as a table. It retains neither an inventory manifest nor BIND
generation. The Terraform normalizer remains solely responsible for converting
canonical declarations into provider resource values. The linter compares
canonical DNS names to enforce single-file ownership; it does not translate
values, TTLs, or provider settings.

## Terraform and credentials

Each owner uses `tf_setup` when present, otherwise `tf`. Root-workspace wrappers
package nested owners through source runfile aliases, preserving their own
root and state without copying module implementation. Global-only owners use
the shared module's global entrypoint, so they do not require RouterOS.
The combined entrypoint groups RouterOS resources in its dc1 child module.

Pins and import semantics are documented in [provider compatibility](provider-compatibility.md).
Related resource values preserve the original effective TTLs and multiplicity.
New AppRole modules give each owner its state subtree and only DNS credential
reads for its views; new roles do not join the broad infrastructure group.

Before adoption, `dns_enabled=false` prevents resource creation. It does not promise that
provider initialization is offline: the pinned RouterOS provider probes its
version even with no DNS resources. Operational root commands therefore
require the reviewed Vault policy and provider credential prerequisites. The
ordinary AL injection path supplies real credentials; tests use provider mocks
without contacting infrastructure.

Cloudflare's shared module resolves the zone by name when the optional zone ID
is absent or empty. It performs no zone lookup while disabled or when the
owner has no global records. The existing token must permit zone listing and
reading. DNS and ingress share their configured RouterOS endpoint through
`infra/al_lib.lua` and retain separate Vault credentials. Existing credential
entries need no new metadata fields. The configured router hostname does not
replace the independent recovery requirement for endpoint-owning scopes.

## Transfer and rollback

Current source removes central write entrypoints. The recorded writer audit
and serialized owner operations control competing writes within their stated
coverage; unidentified external writers and historical checkouts are not
proven disabled. Existing DNS records remain live while exact-ID imports
transfer their state ownership. Verified missing declarations use a reviewed
additions-only plan that preserves every existing provider record.

See [cutover](cutover.md) for concrete prerequisites, activation and rollback.
Source validation is recorded in [validation](validation.md). The user has
authorized deployment one owner at a time while preserving existing records
and services. [Operational evidence](operations.md) records shared prerequisite
verification and [completed rollout](operations.md#verified-completion); each
owner's adoption change records its actual transfer. All nonempty owners have
verified deployment receipts and retain enabled defaults. The coordination's
operational work is complete; exact-candidate quality and publication remain
with repository delivery.
