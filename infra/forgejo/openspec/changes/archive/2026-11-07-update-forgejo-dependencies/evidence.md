# Dependency update verification

Observed on 2026-11-07. This verifies repository inputs and local outputs;
it does not establish deployed versions or live upgrade readiness.

## Release integrity

The publisher release APIs reported stable, non-draft releases for Forgejo
17.0.1, runner 13.2.1, CLI 0.7.0, and Terraform provider 1.7.0. Links are in
[design](design.md). Bazel successfully fetched and verified the pinned server,
source archive, runner, CLI archive, and provider archive. The source archive
extracted with `strip_prefix = "forgejo-src-17.0.1"`.

Server, source, and runner checksums matched their publisher `.sha256` assets.
The provider checksum matched its publisher `SHA256SUMS` and GitHub release
asset digest. The CLI publisher supplies no separate checksum asset; SHA256
of its downloaded release archive was
`6803c12eed17bfadb40ba95da2a6ab58f9b4a7dc10eabfe3512fce7e2fe6216f`.

`bazel mod deps --lockfile_mode=update` succeeded without a lockfile diff.
It reported an unrelated existing OpenCode toolchain import warning.

## Representative output

- `bazel run //third_party/org_codeberg_forgejo_bin -- --version` reported
  `17.0.1+gitea-1.22.0`, including SQLite support.
- `bazel run //third_party/org_codeberg_forgejo_runner_bin -- --version`
  reported `v13.2.1`.
- `bazel run //third_party/org_codeberg_forgejo_contrib_forgejo_cli_bin -- version`
  reported `fj v0.7.0`. This CLI uses `version`, not `--version`.
- Runner `validate --workflow --path` accepted a task-local workflow with a
  push trigger, one job, and one run step. Validation did not execute jobs.
- Both `//infra/forgejo/ansible:ansible_bin` and
  `//infra/forgejo_runner/ansible:ansible_bin` built successfully.
- `//infra/forgejo/tf:tf.direct.terraform` initialized declared providers with
  `init -backend=false -input=false` against a task-local copy of the exact
  configuration and repository catalog, then returned
  `Success! The configuration is valid.` from `validate`.

Terraform validation required a shorter task-local `TMPDIR` under
`out/forgejo-deps/tmp`: the original 91-byte absolute temporary directory
exceeded the Unix socket path limit after the provider socket filename was
appended. Neither validation command accessed the configured backend.

## Compatibility and limits

The Forgejo 16 notes preserve the configured OAuth2 JWT secret behavior and add
asymmetric LFS JWT options; the pinned 17 source still documents the configured
secret and OAuth2 registration settings. The Docker trusted-proxy default
change does not control this binary deployment. Forgejo 17 introduces strict
duration parsing; the checked-in service configuration has no duration
overrides. Upstream changes to Git hooks, mirroring redirects, rendering, and
database migrations still require consideration before a separately authorized
live upgrade. No live plan, playbook, server startup, or deployment ran.

## Repository checks

Passed `git diff --check`, `bazel run //:format`, and semantic lint of the four
Forgejo dependency aliases, Terraform dependency documentation, service
Terraform launcher, and both Ansible package targets.

The following combined test invocation passed all 125 selected tests:

```sh
bazel test //:repo_quality_test //:buildifier_test \
  //infra/forgejo/tf:tf_tests.fmt_test \
  //third_party/terraform/test:provider_installation_test \
  //infra/src/openspec/validation:validate_test \
  //infra/src/openspec/validation:archive_test
```

`OPENSPEC_PROJECT=infra/forgejo bazel run //tools/openspec -- validate
update-forgejo-dependencies --strict` also passed. Publication state and the
exact reviewed commit belong to Git and local delivery receipts.
