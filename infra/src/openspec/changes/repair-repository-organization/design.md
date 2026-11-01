## Decisions

Decision review: proceed with the existing directory taxonomy and dependency managers. Wholesale relocation or new dependency management would not repair the demonstrated faults.

Third-party packages retain executable pins. Legacy rules consume those executable targets rather than maintaining independent versions. Nested modules use native filegroups for specification source exports; only CLI and validation consumers depend on rules_openspec. Documentation dependencies remain regular when public BUILD files load their rules.

The representative Drawio desktop action failed after selecting the canonical executable. Decision review: migrate its only two consumers to the already maintained sandboxed browser renderer, preserving target names and packaged SVG filenames, and remove the unused desktop toolchain. Retaining another desktop pin would restore the reviewed duplication rather than repair ownership. The original failure evidence is in out/repository-organization/drawio-consumer.log.

Owner declarations remain authoritative. A repository-quality inventory updater discovers source-export packages and module identities, and renders central validation membership and the project README index. It detects stale output without contacting external services. CI and the full checker share discovery from .bazelignore. Existing project landing membership remains an explicit publication decision, not inferred from directory existence.

Artifact auditing uses Bazel genquery with execution-tool and implicit edges excluded, starting from explicit representative runtime targets. This complements native enforcement without claiming a complete configured runtime or action graph. Existing documentation and developer package exceptions remain unchanged.

## Failure cases and validation

Before implementation, cover newly added owners and workspaces, stale inventories, absent source exports, ignored non-workspace directories, unsafe boundaries, continuation after a failed build, forbidden indirect artifact dependencies, and execution tools excluded from the audit. Preserve repeatable fixture output under Bazel test outputs. Check canonical executable selection and the absence of workflow dependencies in nested module declarations. Binary schema tests and migrated constraints belong to the linked owner change.

## Candidate and continuation

Worktree: t3/review-repo-organization. Initial source revision: 4e10a19e32f51c93f9c56a3f872951725a8a3f57. Validation receipts and logs belong under out/repository-organization. Next action: write behavioral fixtures, then implement and validate the exact candidate.

## Validation evidence

Binary schema and runtime fixtures pass (10 tests). The initial integration
checks demonstrated indirect artifact inclusion and execution-tool exclusion.
All affected standalone module graphs were reconciled through `bazel mod deps`.
The Helm command fixture executes the canonical binary successfully. Both
archived Drawio outputs now have finite positive viewboxes (2281 by 2349 and
1376 by 1686), and the complete Hugo site builds with their packaged output.
The renderer fixture reproduces the pinned exporter's missing shadow-padding
return value and verifies the corrected adapter. Inventory freshness now agrees
with the configured Markdown formatter. Exact candidate checks and publication
are recorded separately by al-delivery.

Session review: an overbroad staging proposal was rejected by automatic approval
review. Only explicitly inspected new task files were staged afterward. The
Drawio migration required source inspection after successful actions produced
invalid SVG dimensions; acceptance checks therefore inspect representative
output as well as command results. Raw failure and validation logs remain under
ignored out/repository-organization.
