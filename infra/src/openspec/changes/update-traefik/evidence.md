# Traefik update acceptance

On 2026-10-26 the original upgrade verified upstream
[release 3.7.16](https://github.com/traefik/traefik/releases/tag/v3.7.16).
The published checksum manifest and release asset metadata both reported Linux
AMD64 archive SHA-256
`759a7bc0d66b74c588255ed9e41c7f9d4aba60bdb7bce9dde551b1268a55eb07`.
The extracted lock preserves its SRI equivalent and moves platform constraints
from `rules` to `toolchain` without changing their values.

The original upgrade checked the [migration notes](https://doc.traefik.io/traefik/migrate/v3/),
ran twelve isolated routing checks, and verified the public service running
3.7.16. Historical live acceptance belongs to the download change in
[PR #202](https://github.com/alwaldend/src/pull/202). This independent extraction
performs no new deployment.

Independent extraction validation passed against master on 2026-10-26:

- The built binary reports `Version: 3.7.16` and `OS/Arch: linux/amd64`.
- `//infra/download/ansible:ansible_bin` packages successfully.
- Repository quality, Buildifier, repository OpenSpec validation, and semantic
  lint of both affected packages pass.
- The implementation diff contains only `binary_toolchain.json`; the other
  files are this change record.

Logs and guarded delivery receipts are under `out/traefik-version/` in the
extraction worktree. Publication state is owned by those receipts and GitHub.
