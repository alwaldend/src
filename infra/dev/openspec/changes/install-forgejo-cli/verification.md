# Installation and authentication verification

On 2026-11-16 the user authorized fj installation on the current dev1 host and
launching login. The owning `//infra/dev/ansible:fj_install` target runs the
packaged playbook locally, selects only the `fj` tasks, and disables privilege
escalation. Task listing and syntax checks passed before installation.

The first copy selected a Bazel executable launcher, and standalone execution
failed to locate runfiles. The source was corrected to the existing pinned
Linux x86-64 archive's `fj_filegroup`, exposed through its owning
`include.MODULE.bazel`. Reconciliation replaced only that task-owned executable.
The binary pin and archive checksum were unchanged; `bazel mod deps` passed.

Corrected installation reported `ok=3 changed=1 failed=0`. Reapplication reported
`ok=3 changed=0 failed=0`. The installed command reported `fj v0.7.0`, and a
byte comparison with the packaged standalone payload passed.

Native `auth login` reported that no instance client was configured. The browser
was opened at Forgejo's applications page and redirected to sign-in. The user
then supplied an application token through a private credential handoff. Native
`auth add-token` consumed it through stdin, and `whoami` verified the expected
`simeonwarrenbot` account. The handoff file was deleted; the CLI credential file
is mode `0600`. No credential contents or OAuth state are retained here.

The existing bot userpass secret remains owned by Vault at its existing exact
path. Its link was supplied to the user without reading the password. No server
OAuth application, Vault configuration, runner activation or unrelated service
deployment was performed. Browser OAuth remains available after configuring an
appropriate instance client through its owning workflow; token login is the
observed authenticated state for this installation.

Local logs and delivery receipts remain in ignored
`out/simplify-forgejo-runner/`. Source validation and final local review are
recorded with the follow-up delivery evidence; their results do not establish
unexercised remote PR mutation behavior.
