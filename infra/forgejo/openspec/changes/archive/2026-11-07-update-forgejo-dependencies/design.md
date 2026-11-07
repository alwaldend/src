## Context

See [proposal](proposal.md). Existing Bazel archive declarations own the pins;
Ansible packages their binaries. The Forgejo Terraform consumer uses an exact
provider constraint, which must match the shared archive declaration.

## Goals / Non-Goals

Preserve the existing platforms, packaging, and provider interfaces while
updating stable releases. Live deployment and database migration are excluded.

## Decisions

- Select the latest stable releases reported by the publishers on 2026-11-07:
  [Forgejo 17.0.1](https://codeberg.org/forgejo/forgejo/releases/tag/v17.0.1),
  [runner 13.2.1](https://code.forgejo.org/forgejo/runner/releases/tag/v13.2.1),
  [CLI 0.7.0](https://codeberg.org/forgejo-contrib/forgejo-cli/releases/tag/v0.7.0),
  and [provider 1.7.0](https://github.com/svalabs/terraform-provider-forgejo/releases/tag/v1.7.0).
- Use publisher SHA256 files for the server, source, runner, and provider.
  The CLI release has no checksum asset; calculate SHA256 from the release
  archive obtained directly from its publisher.
- Keep source and binary Forgejo versions identical. Reconcile the Bzlmod lock
  using Bazel; inspect it for unrelated changes.
- Run binary version commands directly through dependency aliases, avoiding
  the authenticated service wrappers. Build deployment packages without running
  playbooks. Check Terraform formatting and provider packaging locally.

## Risks / Trade-offs

- Major server upgrades can change configuration and database behavior. Check
  the 16 and 17 release notes and configured settings; local packaging and version
  checks do not establish live upgrade readiness.
- CLI 0.7 removes positional `help`; existing wrappers forward arguments. Its
  local version command is `version`, while server and runner use `--version`.
- New provider write-only secret arguments are optional; existing configurations
  retain their current arguments and secret-handling behavior.

## Migration Plan

Only update repository inputs. A future deployment requires separate authority
and an upgrade procedure appropriate to the then-deployed version and data.
