---
title: SSH
description: Setup ssh
tags:
  - ansible_role
---

Replacing a signed host certificate notifies a role-owned handler. It runs
`sshd -t` before reloading the daemon so new connections receive the renewed
certificate even when the hardening configuration is unchanged. The service
name follows `sshd_service_name` when provided, otherwise the OS-family default.
Existing sessions survive the reload. This is deployment-time renewal; the role
does not install a certificate renewal timer.
