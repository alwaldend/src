# ClamAV dependency validation

Observed on 2026-10-29 in the existing `ops/host-security-updates-20261029`
worktree. Evidence is under ignored
`out/host-security-updates-20261029/clamav/dependency/`. The source scope is this
dependency package and the generated root `MODULE.bazel` include. Publication
and live operations are withheld by the task instruction.

## Artifact evidence

- The downloaded DEB is 105819396 bytes, SHA-256
  `28d6efc5b4423e7830c3559339552eb53870a9eac51ac4efb37d60530d329886`.
  Both values match the publisher's release API response in
  `publisher-release.json`.
- The Bazel `archive` output compares byte-for-byte with that download.
- The Bazel `payload_tar` output compares byte-for-byte with the DEB's
  `data.tar.gz` member, extracted independently with `ar p`. Its SHA-256 is
  `0f4f8db22849120b57ea9f16fb1011a68ed5ea44f88268ea6384341c7f09873a`.
- `artifact-sha256.txt`, `artifact-paths.txt`, and
  `publisher-payload-files.txt` preserve the actual artifact observations.
  The payload retains upstream paths and relative shared-library symlinks.
- The detached signature was not verified; the integrity evidence is the
  independently calculated hash and publisher metadata match.

## Build and runtime checks

The following repository commands completed successfully through `bazel_agent`:

```sh
al bazel run //tools/bazel_module:update
al bazel mod deps --lockfile_mode=update
al bazel build \
  //third_party/com_github_cisco_talos_clamav:archive \
  //third_party/com_github_cisco_talos_clamav:payload_tar \
  //third_party/com_github_cisco_talos_clamav:docs
al bazel build --config=lint \
  //third_party/com_github_cisco_talos_clamav:archive \
  //third_party/com_github_cisco_talos_clamav:payload_tar \
  //third_party/com_github_cisco_talos_clamav:docs
```

Module generation added exactly one include; module reconciliation produced no
lockfile change. Module freshness passed. The combined freshness/buildifier
invocation initially failed only on parent-owned
`infra/security/ansible/BUILD.bazel` formatting. The dependency's definitions
were formatted with repository buildifier; the parent was notified of the
unrelated workspace result. Full invocation logs are `module-update.log`,
`module-deps.log`, `build.log`, `validation-tests.log`, and `lint.log`.

The targeted buildifier check passed for all three dependency Bazel files.
Repository Prettier formatted the documentation and YAML. Strict OpenSpec
validation passed (`openspec-validate.log`). No test or source changes were
made outside the assigned dependency and generated module include.

The Bazel payload tar was unpacked under the evidence directory's `payload/`.
On Linux x86_64 with glibc 2.43, all three commands returned `ClamAV 1.5.4`:

```sh
clamav_root="$PWD/out/host-security-updates-20261029/clamav/dependency/payload"
clamav_config="$PWD/out/host-security-updates-20261029/clamav/dependency/version-only.conf"
# version-only.conf is an empty task-local file, suitable only for this check.
LD_LIBRARY_PATH="$clamav_root/usr/local/lib" \
  "$clamav_root/usr/local/bin/clamscan" --version
LD_LIBRARY_PATH="$clamav_root/usr/local/lib" \
  "$clamav_root/usr/local/bin/freshclam" --config-file="$clamav_config" --version
LD_LIBRARY_PATH="$clamav_root/usr/local/lib" \
  "$clamav_root/usr/local/sbin/clamd" --config-file="$clamav_config" --version
```

`version-verification.txt` records the successful execution at 19:57:52 UTC.
The initial freshclam invocation failed because the default configuration did
not exist; `version-default-config-failure.txt` preserves that evidence. Passing
an explicit empty configuration resolved this version-only check. No signature
database was downloaded, package installed, or service started.

`clamscan-linkage.txt` records the upstream `/usr/local/lib` RUNPATH.
`glibc-symbol-versions.txt` records a maximum required GLIBC symbol version of
2.30 across the packaged executables and shared libraries. This is not a test
on an older endpoint OS or proof of scan behavior.

## Session ergonomics

`CLAMAV-DISCOVERY-OUTPUT`: Initial catalog and file reads were too broad and
produced truncated outputs. Focused tool-name selection and owner-local reads
resolved this routed discovery defect; no shared policy change is proposed.

`CLAMAV-SOURCE-ARTIFACT-PATH`: File-only aliases did not create execroot symlinks.
Source artifacts were resolved from Bazel's external repository directory after
`cquery`, avoiding another download or alternate extraction implementation.

The owner-local tasks and this record identify the source, interface, exact
artifact hashes, successful execution, and remaining publication boundary.

## Separate dependency pull request

On 2026-10-30, the user authorized extracting the dependencies from security
PR #214 into a separate master-based pull request. The isolated
`t3/extract-third-party-deps` worktree preserves this package from source commit
`5e4bf7ccb7e7948e12c12ae12b92d048bbba06e3`; only this validation addendum changes.
The earlier publication boundary above describes the original observation.

Fresh evidence is under ignored `out/dependency-extraction-20261030/`.
The package's `:all` build and strict OpenSpec validation pass. Both archive
hashes match the artifact evidence above, and `ar p <deb> data.tar.gz` compares
byte-for-byte with the Bazel payload. Isolated `clamscan`, `freshclam`, and
`clamd` version commands each report `ClamAV 1.5.4`, using the extracted library
path and explicit empty version-only configuration. Module generation preserves
the two requested scanner includes and lock reconciliation makes no change.
The final delivery receipt binds the committed candidate and publication gates.
