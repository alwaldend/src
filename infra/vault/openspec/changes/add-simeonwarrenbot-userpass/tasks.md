## 1. Selected userpass entries

- [x] 1.1 Generate 72-character passwords and userpass entries for an explicit selected-user list containing only simeonwarrenbot; validate pinned-provider formatting and schema.
- [x] 1.2 Remove user/group storage grants without adding deny rules, retain exact infrastructure management and the human bot-userpass read grant, and verify owner denials, selection, length, and lifecycle against disposable Vault.
- [x] 1.3 Reconcile specifications and documentation and review the full candidate; retain exact-candidate quality, lint, and publication receipts.

## 2. Authorized deployment

- [x] 2.1 Obtain explicit approval of the revised infrastructure apply scope; the user authorized it on 2026-11-16.
- [x] 2.2 Review and apply the two policies and two assignments, then the generator, stored entry, and existing login through the owning interactive Terraform target.
- [x] 2.3 Verify the deployed password, original identity, human read exception, absent bot access, shared helper integrity, and a full no-change plan.

## Continuation

The user specified userpass as the entry name, length 72, and an explicit
selected-user list containing simeonwarrenbot. They then clarified that users
must receive no grant, rather than an explicit deny. Remove the proposed
users-group namespace policy. Keep only infrastructure management and one
human read exception. The authorized configuration is now deployed; prior
designs were never deployed.

The provider metadata finding is addressed by metadata read/update on the
selected-entry writer. An outdated review-thread reply was rejected by
automatic approval and remains untouched; do not retry it here.

The local test driver uses -dev-no-store-token and checks helper integrity;
keep those safeguards. A previous fixture omitted that flag and the user
restored the host login. Clean fixture credentials, logs, and state. Reuse
this worktree and branch, with receipts in out/simeonwarrenbot-auth.
