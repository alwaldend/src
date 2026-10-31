# Rsync publication verification

The source changes were exercised on 2026-09-26 in the dedicated download
feature worktree, based on `43b5cb1923cb0887574120ccf1bcac4314983b84`.
Final candidate and validation receipts are task-local under
`out/infra-download-split/ssh-publication/`. No production VM, DNS record,
Vault value, or website was changed.

## Isolated SSH-to-HTTP acceptance

The task-local `e2e.py` invokes the built release executable and both packaged
infra publishing targets against OpenSSH and Nginx in a disposable Fedora
container. The test administrator authenticates with an ephemeral key and
uses real sudo to run rsync and selection as `download`. That account has no
authorized keys. The content filesystem is a 48 MiB tmpfs with `noexec`;
the fixture does not pretend to test Btrfs or SELinux enforcement.

The observed fixture image is
`50e8c70641a2f172d28798130df13da2f3735bff9c4aa29ad0fd2276ec9efe9a`, with
Nginx 1.30.5, OpenSSH server 10.2p1, rsync 3.5.0, and sudo 1.9.17.
Its container and ephemeral keys are removed after each run.

All 33 cases passed: host-key refusal; administrator/content-account boundary;
ordinary archives, Unicode and punctuation names, JSON listings, HTTP bytes,
permissions, repeat/update preservation; tar.gz and zip sites and assets;
retained-site selection without local payload; invalid paths, links, devices,
duplicate/malformed entries, and missing index; failed, interrupted, full-disk,
and concurrent transfers; both component wrappers; generated manifest/page
URLs; and the existing OCI executable argument contract.

Repeat in the prepared worktree:

```sh
bazel_agent bazel build //tools/release/main/go
bazel_agent bazel run --script_path=out/infra-download-split/ssh-publication/publish.local //infra/download:publish.local
bazel_agent bazel run --script_path=out/infra-download-split/ssh-publication/publish.yandex //infra/download:publish.yandex
python3 out/infra-download-split/ssh-publication/e2e.py
```

`e2e-result.json` retains the cases, package/image identity, and source hashes.
This is task-local acceptance evidence, not a newly maintained CI or Molecule
suite. It preserves the earlier decision to develop maintained host fixtures
through Molecule separately.

## Other checks

The release Go tests, strict release/download OpenSpec validation, affected
semantic lint, and all 25 repository quality checks passed before preparation.
The final delivery validation repeats required gates on its committed candidate.
The existing OCI `head_deploy` wrapper also builds, while the SSH wrappers
contain neither a Soras argument nor its runfiles dependency.

Both packaged Ansible inventories pass syntax/task-list checks and render
Nginx, Traefik, and deduplication configuration without a publisher-key
variable or environment input. The role retains the content account and
filesystem ownership while removing its authorized-key assertion/task.

## Corrected findings and limits

The draft custom receiver was removed when the user selected rsync. Ordinary
files now follow rsync update semantics; retained website versions are selected
without a digest database. A Bazel label default initially defeated `oras=None`;
the public macro now owns that default. An actual release-page fixture also
covers generation without Git metadata, which previously dereferenced nil.

Fixture image acquisition initially failed because its loopback proxy was
unreachable; only the disposable image build uses host networking. A direct
wrapper invocation before materializing runfiles failed, so fixture invocation
uses Bazel-generated launchers. Lint alone did not refresh the executable for
one fixture attempt; rebuilding the executable and rerunning all cases passed.
These failure logs remain next to the passing result.

Live certificate authentication, target-host sudo configuration, SELinux,
Btrfs persistence, and public routing remain authorized rollout checks.

## Independent publication PR

Publication is extracted from the aggregate candidate
`e2d4f59aee637445f996aac6e05c87704b8fa636` into its own branch from
`43b5cb1923cb0887574120ccf1bcac4314983b84`. The release implementation and
component publishing targets are unchanged. Host configuration, including
the content account without authorized keys and the server rsync package,
remains in [PR #102](https://github.com/alwaldend/src/pull/102). Neither branch
is stacked on the other; host setup is a prerequisite for live publication.

Split validation receipts, source-preservation evidence, and the isolated
acceptance result are retained under `out/publication-extraction/` in the
publication worktree. Its fixture copies the Nginx template from the original
aggregate commit into ignored scratch, recording its source and hash; the
publication build has no dependency on unmerged host source. The same 33-case
fixture runs against the standalone release executable and publishing targets.
The original fixture commands above describe their original worktree; the
split uses launchers and `e2e.py` under `out/publication-extraction/`.

## Review follow-up after host merge

The host branch advances to merged trunk
`63bc5d924c31a808fee12173872e253675e5f0bb`; its tree is unchanged and its remote
feature branch was already deleted. Publication was rebased independently to
`b8099584fc583c27867c8d448cfa0189889e8b2c` before addressing the three review
threads. The user confirmed independently selectable upload, extract, and link
steps, and requested the local-dependency-first rule in AGENTS.md.

The expanded isolated SSH/HTTP fixture passes 44 cases, including upload alone,
extraction from a published archive without local payloads, link-only operation
without an archive or rsync, retained extraction, invalid step sequences, error
context, and literal quoted version/key paths. Previous archive-safety, disk
exhaustion, interrupted/concurrent transfers, both wrapper targets, URL
projection, and OCI checks still pass. The fixture image remains the immutable
image recorded above; no live host is used.

The release Go tests, owner OpenSpec validation, generated skill-discovery test,
and repo-go offline Promptfoo configuration test pass. The latter checks skill
staging and configuration, not model behavior. The repo-go skill is packaged
canonically under tools/agents/skills and discovered through the generated link.
The release package's propagated errors retain causes through contextual %w
wrapping, including filesystem, parsing, subprocess, and archive-reader errors.

The initial review follow-up added chhongzh/shlex for a quoting API, even
though the repository already pinned anmitsu/go-shlex. The user rejected that
unapproved addition. The correction removes it and restores the existing
POSIX quoting approach; dependency manifests, checksums, and module imports
return to their trunk contents. AGENTS.md now strongly discourages new
external dependencies and requires approval of the specific dependency before
installation, vendoring, or source/build declaration. The existing 44-case
fixture covers literal apostrophes, Unicode, punctuation, and quoted key paths;
its new result and final-candidate hashes accompany the correction receipt.

Final candidate validation, fixture source hashes, commands, and guarded review
receipts are retained under `out/rebase-review-113/` in the publication worktree.
The final delivery gates rerun repository quality, affected semantic lint,
Go/OpenSpec/skill checks, and both SSH plus existing OCI wrapper builds.
