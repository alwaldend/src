## Context

The user explicitly requested certificate renewal through Ansible with Vault
certificate verification disabled, and documentation of that procedure. The
active inventory selects `host1.vault.alwaldend.com`. Its certificate expired
at 2026-10-02T01:16:21Z; this blocked the pending X Article upload before any X
request.

## Goals / Non-Goals

Renew the existing server certificate through its owning role and restore
verified TLS. Retain the existing authentication, inventory identity, and service
configuration. Do not deploy Terraform, change quorum, restart Vault, or
persist a verification bypass.

## Decisions

- Select the existing issuance, certificate/key copy, and reload tasks with
  `vault_tls`; the broad `vault` tag would also deploy unrelated host state.
- Fail closed when `vault_tls` is explicitly requested, while preserving the
  existing full-bootstrap tolerance. Hide the certificate block's result,
  which contains a private key.
- Use `VAULT_SKIP_VERIFY=true` only in the recovery process environment. The
  pinned Vault client and Ansible subprocess inherit it; no AL configuration
  change is required. Certificate login uses curl's separate `--insecure` flag
  only if an expired bootstrap token also needs replacement.
- Reload the existing service with its configured SIGHUP operation after both
  certificate files are installed. Verify normal TLS independently afterward.
- Set `ansible_host` to the existing DNS-owned connection hostname without
  changing `inventory_hostname`, which supplies the Raft node identity. Include
  the connection hostname in certificate alternate names so the existing
  unseal endpoint remains covered.
- Derive temporary SSH host trust from the checked-in public server CA and use
  the existing certificate's host alias. Keep strict SSH verification enabled
  and scope these options to the Ansible invocation.

## Risks / Trade-offs

- Temporarily disabling server verification weakens endpoint authentication;
  the user authorized this recovery exception for the configured Vault endpoint.
  Keep the override scoped and restore validation for follow-up operations.
- Certificate and key installation are separate tasks. A failed write stops
  before reload; a subsequent authorized recovery must complete both files.
- An invalid host token or unavailable YubiKey can block renewal even when the
  verification override works. Do not infer successful renewal from source
  validation or a command that never reached Ansible.

## Execution state

The initial lookup reached Vault but rejected the saved token. The user
completed YubiKey login in their terminal after making the device available;
the agent did not collect the PIN. Two Ansible attempts then stopped during
fact gathering, first on the obsolete connection hostname and then on missing
SSH host trust. Neither reached certificate issuance. The inventory mapping
and invocation-scoped repository CA trust resolved these distinct failures.

On 2026-10-02 at 20:20:47Z, the authorized host-limited `vault_tls` run started
with temporary Vault verification bypass and strict SSH verification. It
completed at 20:21:04Z with `ok=5 changed=4 failed=0 unreachable=0 ignored=0`:
fact gathering, issuance, both file writes, and reload. No other host setup
tasks ran. The packaged task listing and syntax check also passed.

At 20:21:24Z, `//:vault.status` succeeded with `VAULT_SKIP_VERIFY` removed and
reported an initialized, unsealed service. At 20:21:51Z, ordinary OpenSSL
chain, time, and hostname checks verified both `vault.alwaldend.com` and
`host1.vault.dc1.alwaldend.com` against the repository root CA. The new leaf is
valid from 2026-10-02T20:20:27Z through 2026-11-02T06:49:58Z and includes both
host identities and both shared service names. These are point-in-time
observations, not a continuing health guarantee.

Sanitized receipts and validation logs are under
`out/vault-certificate-recovery/`; credential-bearing raw renewal output was
removed. Source publication is recorded by the repository delivery receipts
and pull request, separately from this operational evidence.
