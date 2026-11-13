# Backup coverage correction verification

| Dimension    | Outcome                                                                                                                                                                            |
| ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Completeness | All three PR 273 findings are corrected within the existing complete-recovery contract.                                                                                            |
| Correctness  | The real-binary rehearsal independently rejects seven omitted roots before downtime, rejects covered storage exchanges and restores removed external persistence from the archive. |
| Coherence    | Conservative INI normalization, effective legacy storage and explicit startup coverage preserve the existing activation boundaries. No deployment or new dependency is introduced. |

## Review dispositions

- Quoted INI paths: simple single, double and backtick quotes are normalized;
  ambiguous syntax is rejected. The actual external Git repository uses a
  quoted absolute path and is read through Forgejo's API after restoration.
- Legacy storage: effective LFS, avatar and repository-avatar paths are included
  with modern replacement precedence. External legacy storage types and advanced
  named profiles are rejected. Effective subsystem identities reject exchanges
  between already covered roots.
- Startup coverage: executable, runtime configuration and unit must all be
  covered by overridden roots before the old service stops.

The strengthened test was written before the production correction. Against
the prior implementation it failed because an upgrade succeeded with omitted
external repository coverage. A fixture token parsing issue caused by legacy
configuration warnings was corrected before confirming that reproduction.

The accepted rehearsal separately omits four external persistence roots and
three startup roots, preserving the running old instance in each case. It then
changes unused legacy values while keeping modern effective locations, upgrades
15.0.3 to 17.0.1, checks unchanged behavior and downgrade rejection, removes the
actual data, external repositories, object trees, configuration and executable,
and service unit, extracts the archive and starts 15.0.3 with its matching state. The committed
repository README and retained fixture objects are recovered.

The preflight scenario also rejects malformed and multiline quoting, inline
comments, relative explicit storage paths and named storage profiles without
service actions. Existing drain, backup, version and doctor failure scenarios
continue to pass.
Explicit local subsystems retain application-data defaults. Global
`[storage].PATH`, including an empty setting, is rejected before downtime;
the pinned release's migration guard makes its interpretation dependent on
instance history. The design records the failing initial fixture and source.

## Validation and limits

The combined selection of the integration target, root repository quality,
buildifier and OpenSpec validation/archive tests passes all 124 Bazel tests.
Affected-target semantic lint, explicit Ruff checks of the filter and fixture
sources, and the Forgejo deployment package build pass. The owner-local change
passes strict OpenSpec validation with no specification delta because it repairs
implementation of the existing guarantees.

Commands and coverage boundaries are retained in the original
[verification](../2026-11-07-add-controlled-forgejo-upgrades/verification.md).
Sanitized scenario JSON and raw validation logs remain in ignored
`out/forgejo-upgrade/`; Git and delivery receipts own exact final candidate,
publication and remote review state. The service-manager fixture is unchanged;
real systemd, live inventory, production data and actual LFS/avatar HTTP upload
workflows remain outside coverage. Fixture objects verify archive preservation,
not those upload workflows.
