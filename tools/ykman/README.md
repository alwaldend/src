---
title: Ykman
description: Ykman (Yubikey Manager)
tags:
  - yubikey
---

## Links

Python dependencies use the shared [uv locking workflow](../py/README.md).
Dependencies are declared in `pyproject.toml` and locked in `uv.lock`.
Run `//tools/ykman:uv_lock.update` to update the native lock and review its
diff. Build `//tools/ykman:uv_lock` to generate it without source writes.
The sandboxed resolution action requires network access to PyPI.

- Install: https://docs.yubico.com/software/yubikey/tools/ykman/Install_ykman.html
- Import certificate: https://docs.yubico.com/software/yubikey/tools/authenticator/auth-guide/smartc-import-certificates.html
