---
title: T3 Code
description: Run the T3 Code server as a systemd service
tags:
  - ansible_role
  - ai
---

The default `npm` installation mode retains the npm launcher. Callers can
select `t3code_install_method: standalone` and supply a release URL and
`sha256:` checksum. The role installs the complete official bundle under an
administrator-owned release directory and invokes its executable directly.
State, work, and home remain owned by the ordinary T3 service account.
Archive comparison repairs missing or incomplete bundle files on reapplication;
the existence of the executable alone does not skip extraction.
Standalone mode uses native T3 pairing/session authentication and omits the
legacy npm mode's `T3CODE_AUTH_TOKEN` environment entry.

Changed runtime, environment, or service inputs restart T3 Code. An unchanged
run ensures the service is started. Token-bearing environment writes are
excluded from Ansible logs. Optional `t3code_initial_settings` seeds first
startup without overwriting later UI settings. Deployment and provider setup
remain the caller's responsibility.
