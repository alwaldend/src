---
title: SSH
description: Setup ssh
tags:
  - ansible_role
---

Normal setup installs fail2ban, installs the trusted user CA, discovers and
signs host keys, and invokes `devsec.hardening.ssh_hardening` with those signed
certificate paths. The upstream role retains its native tasks and restart
handler.

The `ssh_hardening_match_user` default supplies `Match User ansible` with
`AllowTcpForwarding remote`, `PermitListen none`, and
`AllowStreamLocalForwarding remote` through the upstream template. Set
`ssh_hardening_match_user: []` to omit this match, or replace the list with
custom entries. The wrapper explicitly maps this variable to upstream
`ssh_server_match_user`; use the wrapper variable for inventory overrides.

Select normal SSH setup through the owning deployment's `ssh` tag and an
explicit host limit, for example:

```sh
al bazel run //infra/ingress/ansible:ansible -- \
  --limit host1.ingress.alwaldend.com --tags ssh
```

This uses the owner's authentication and runs the full SSH setup described
above, including the native restart handler. The `ssh_hardening` tag alone
skips the required certificate discovery. The upstream template validates its
candidate during apply; `--check --diff` previews do not validate the proposed
configuration. The scanner does not invoke this role or change server policy.

Replacing a signed host certificate notifies a role-owned handler. It runs
`sshd -t` before reloading the daemon so new connections receive the renewed
certificate even when the hardening configuration is unchanged. The service
name follows `sshd_service_name` when provided, otherwise the OS-family default.
Existing sessions survive the reload. This is deployment-time renewal; the role
does not install a certificate renewal timer.
