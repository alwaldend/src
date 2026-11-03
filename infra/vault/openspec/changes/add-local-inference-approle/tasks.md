## 1. Component identity

- [x] 1.1 Declare the component AppRole and inventory memberships; verify
      Terraform formatting and the packaged root without accessing live Vault.
- [x] 1.2 Add global DNS access and scoped server PKI enrollment; verify
      formatting and packaged configuration without a Vault write.

## 2. Independent delivery and authorized apply

- [x] 2.1 Extract the Vault changes and canonical DNS prerequisite into a
      master-based branch; verify formatting, package builds, and owner checks.
- [x] 2.2 Save and inspect the owning Terraform plan, then apply its reviewed
      in-scope changes; retain only sanitized results and remove private artifacts.

- [x] 2.3 Correct the shared injector and backend canonical runfile paths after
      the standalone AL runner could not resolve the consuming root name.
- [x] 2.4 Renew the expired Vault server certificate through host-limited
      `vault_tls` Ansible tasks, then verify normal TLS before both applies.
- [x] 2.5 Review and apply the local-inference DNS saved plan through its
      existing owner root; retain sanitized evidence.

## Execution evidence

On 2026-11-03, the initial Vault plan stopped before plugin startup because
`com_alwaldend_src` was not mapped by the external AL binary. The `_main`
plugin paths resolved the injector and backend; the next plan reached client
initialization. A read-only status check then identified the expired server
certificate (2026-11-02T06:49:58Z).

The user authorized temporary insecure Vault certificate recovery. The existing
`vault_tls` workflow was invoked only for `host1.vault.alwaldend.com`, with
strict SSH host verification against the checked-in server CA. It stopped
before Ansible started: AppRole SecretID creation returned HTTP 403. A root
read-only TTL lookup also returned HTTP 403 and invalid token. The user was
asked to refresh host authentication through the existing YubiKey login.
No certificate, Vault Terraform, or DNS changes were applied by these attempts.

Offline checks passed before the runtime prerequisite edit: 34 tests, including
repository quality, Terraform formatting, OpenSpec structure, and canonical DNS
ownership; affected Vault and DNS-export semantic lint also passed. Re-run the
mandatory checks on the final prepared candidate before publication.

After the user refreshed login, the host-limited certificate recovery completed
on 2026-11-03 with `ok=5 changed=4 failed=0 unreachable=0 ignored=0`.
A root status request succeeded with `VAULT_SKIP_VERIFY` removed.

The exact saved Vault plan was reviewed and applied: 15 resources added,
2 identity groups updated, and none destroyed. Four obsolete `local_file`
certificate snapshot entries were forgotten without destroying their files.
A fresh plan exited 0 with no changes. The PKI roles remain restricted to the
canonical T3 Code hostname, without subdomains or client certificates.

The first DNS-only attempt authenticated through the new AppRole but stopped
on the unrelated Yandex credential injection. Its source correction and DNS
execution evidence belong to the local-inference change in PR #233.

The published prerequisite candidate passed 35 offline tests (including the
existing AL injector E2E), affected semantic lint, and Terraform launcher builds.
Automated code review completed for that candidate without findings.

The corrected DNS-only owner command in PR #233 then applied exactly three
Cloudflare records, with no changes or destroys. Complete API inventory and
public authoritative DNS checks verified the declarations, and the post-apply
scoped plan reported no changes. Its detailed sanitized evidence remains with
the local-inference owner. No GPU resources or host deployment were applied.
