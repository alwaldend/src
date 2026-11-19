---
title: Android
description: Android
---

`packages.txt` owns the SDK and NDK package selection used by `:install` and
`:installer_tar`. The latter packages the installer and the pinned Linux x86-64
command-line tools for Ansible deployment without bootstrapping Bazel on the
remote host. It requires a host Java runtime; downloading SDK packages happens
only when the installed command runs, never in a Bazel build action.
