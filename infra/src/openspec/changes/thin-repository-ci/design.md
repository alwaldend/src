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
to [PR 88](https://github.com/alwaldend/src/pull/88). Remaining work is the
authorized prerequisite deployment and protected-branch CI check after login
refresh, followed by acceptance and archival. Do not equate the prior inline
workflow's successful run with validation of this extracted command.

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
