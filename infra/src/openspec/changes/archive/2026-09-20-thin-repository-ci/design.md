## Decisions

The workflow owns scheduling, checkout, environment, and one command invocation.
The new `repo-ci` skill owns the authoring rule and links mutable settings to
their current owners. Discovery is generated through `.agents:write_skills`.

Preserve the existing Python authentication implementation in a committed
script, packaged by `py_binary`. Rewriting the working authentication protocol
in another language would add risk without helping this extraction. Runtime
credentials stay out of Bazel build actions, flags, logs, and artifacts.
The checked-out CI configuration is a declared runfile, replacing an HTTP read.

Use official `actions/checkout` v4.4.0 at its immutable commit and the event SHA,
with credential persistence disabled. Its Node.js runtime prerequisite belongs
in the runner role. Node and pinned Python use the host's managed CA bundle.

## Verification and continuation

Offline tests exercise resource rejection, exact ref selection, both login
outcomes, policy/TTL failures, token cleanup, redirect rejection and safe errors.
Skill configuration validation does not claim live model behavior coverage.
Validate workflow syntax with the pinned runner, then deploy the role package
change and run the exact candidate through the existing protected-branch
validation deployment. This continues the user's runner deployment authority.

Review verdict: proceed. Independent source comparison found the authentication
flow preserved. Runtime prerequisites are declared, and the CI-mode build,
offline workflow validator, smoke tests, skill discovery/configuration checks,
OpenSpec checks and Buildifier pass. Cold runner execution remains a live check.

On 2026-09-16 the scoped runner deployment stopped before executing any task:
Vault returned HTTP 403, and token lookup confirmed an invalid controller
session. A login refresh was requested without requesting credential content.
Candidate `7f1b96eeefff476aeae17d99b12e6aa2fa7e4d07` passed the full selected
validation plan (94 tests and affected-package semantic lint) and was published
to [PR 88](https://github.com/alwaldend/src/pull/88). At that point, prerequisite deployment and the protected-branch check still
required a refreshed login. The runtime acceptance below supersedes that
blocker; the earlier inline workflow run did not validate this extraction.

### Cold runner prerequisites (2026-09-20)

Run 5 exposed `runner-android-toolchain-missing` and
`runner-module-mirror-stall`: registered NDK toolchains require a local NDK even
for the Python smoke target, while the optional Cloudflare mirror transferred
module archives too slowly to reach Bazel's fallback URLs. Direct BCR and Google
endpoints worked in the recorded runner probes. The refreshed controller Vault
session was verified by a read-only TTL lookup.

Decision review verdict: proceed with declaring Android prerequisites in Ansible
and omitting the optional mirror through the preset generator. Removing NDK
registration would change Android build behavior; command-line overrides had
not avoided resolution. A later empty mirror flag cannot remove accumulated
values. The installer archive reuses the pinned command-line tools and a shared
package list; Ansible performs runtime installation and verifies package
receipts. This avoids requiring a working Android toolchain to bootstrap Bazel
on the runner. SDK installation adds download and disk cost, but preserves the
existing development build contract. Validate the packaged installer, deploy
only the Android/Bazel tasks to secure, and run the exact candidate in CI.

Run 6 (`223f56a3ddc8c63f159c7f0b5a6dda12442a92ef`, source tree
`796a16a275bf52c896a4a01fad54ffe8fb25ac2b`) built the smoke command and passed
resource checks after fetching the pinned LLVM toolchain. It exposed
`runner-ca-bundle-missing`: the workflow's legacy Fedora bundle path did not
exist, so pinned Python loaded zero CA certificates. Read-only probes with the
runner's same Python 3.13.9 runtime loaded 124 CAs from Fedora's canonical PEM
bundle and received HTTP 200 from both Forgejo and Vault. Select that bundle
for Node and Python. This repairs the trust input without weakening certificate
verification. The prerequisite deployment is idempotent (secure: 12 tasks OK,
zero changed, zero failed); explicit archive permissions account for the
runner's restrictive umask.

### Acceptance (2026-09-20)

[Live run 7](https://git.alwaldend.com/alwaldend/src/actions/runs/7) and job 7
succeeded on validation commit `35d5e911148522ff428c436da867959d2b5791d5`.
Its tree `90dbe23076dd3bdf6c5a2056cafeedbc00af3634` exactly matches source
candidate `9acf5679c99a5da8ec03a36adde92e93a6d80975`. The controller verified
`releases/*` protection and the expected tip before publishing with an explicit
lease; the receipt finished at stage `complete` with both run and job successful.

The live output confirms resource capacity, rejection of an unauthorized
audience, successful Vault OIDC login, the exact restricted policies and TTL,
and token revocation. This validates the extracted Bazel command rather than
the earlier inline workflow. All 96 selected tests, affected-package semantic
lint, and installer/Ansible/smoke builds pass. The installer archive launches
independently, and the scoped Ansible deployment reports zero changes and no
failures on repeat. The final archival/documentation tree receives the same
publication gates and an additional exact-tree live check; delivery receipts
remain under ignored `out/repo_ci`.
