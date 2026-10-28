## Candidate verification

The plugin is a standalone Gazelle language using the existing Gazelle 0.53.0
and rules_go 0.60.0 dependencies. The repository remains pinned to Bazel 9.3.1;
all first-party module compatibility declarations now require that baseline.

- Standalone `bazel_agent bazel test //...`: passed. The actual Gazelle binaries
  exercise inheritance, overrides, relative labels, malformed directives,
  attribute/comment preservation, new Go packages, workspace boundaries, and
  repeat generation. Generated fixtures are retained in test undeclared outputs.
- Root AL command test, module generator test, full-check runner test, and its
  offline skill configuration test: all four passed with enforcement enabled.
- Repeated policy generation: zero BUILD-file changes. Three existing
  `gazelle:ignore` packages retain their directives and carry manually maintained
  package declarations. All other root-workspace BUILD packages are generated.
- Direct Traefik certificate consumer: rejected by native Bazel analysis.
- Full-language Gazelle: blocked before writing by existing missing Go embed
  outputs/Anki generated imports and unresolved Python runfiles imports.
  Policy-only generation through the registered extension succeeds.

## Bazel 9.3.1 indirect-policy defect

Stable defect: `BAZEL-TRANSITIVE-POLICY-COLLISION` (fixture-tested).

An isolated workspace uses a public synthetic text file and native filegroups;
no real certificate or live service is involved. The protected package selects
this group:

```starlark
package_group(
    name = "consumers",
    packages = ["//infra/...", "//tools/...", "//website/..."],
)
```

`//infra/cert:certificate` has public ordinary visibility and
`transitive_visibility = "//infra:consumers"`. `//tools/bridge:certificate`
depends on it and also has public ordinary visibility. Root policy is a group
allowing `//...`, inherited by the bridge and project fixture packages through
Gazelle. Analysis observations:

| Consumer                                                         | Expected | Observed  |
| ---------------------------------------------------------------- | -------- | --------- |
| `//tools/bridge:certificate`                                     | Allow    | Allow     |
| `//website:documentation` through bridge                         | Allow    | Allow     |
| `//projects/direct:artifact` through protected target            | Reject   | Reject    |
| `//projects/source:artifact` through protected source file       | Reject   | Reject    |
| `//projects/indirect:artifact` through bridge with local policy  | Reject   | **Allow** |
| Same indirect consumer after removing only bridge package policy | Reject   | Reject    |

Commands use `bazel_agent bazel build --nobuild <consumer>`. Raw fixture BUILD
files, logs, and `enforcement-results.json` are under the ignored task directory
`out/transitive-visibility/`. The failed indirect assertion remains a validation
gap; it is not converted into a passing test.

The likely mechanism is deduplication of distinct package specifications:
[RuleConfiguredTargetBuilder](https://github.com/bazelbuild/bazel/blob/9.3.1/src/main/java/com/google/devtools/build/lib/analysis/RuleConfiguredTargetBuilder.java)
adds the local policy first to an immutable set;
[PackageSpecificationProvider](https://github.com/bazelbuild/bazel/blob/9.3.1/src/main/java/com/google/devtools/build/lib/analysis/PackageSpecificationProvider.java)
inherits [NativeInfo equality](https://github.com/bazelbuild/bazel/blob/9.3.1/src/main/java/com/google/devtools/build/lib/packages/NativeInfo.java)
without exposing the package specifications as Starlark fields. This explanation
is inferred from source; the behavioral discrepancy is observed with the pinned
executable. No Bazel patch, host change, or alternate enforcement system was
introduced.

## Session ergonomics

- `BAZEL-RUNTIME-VERIFICATION`: versioned source and upstream tests were
  insufficient to establish combined-policy behavior. The isolated executable
  experiment caught a missing upstream case. Retain this fixture when upgrading
  Bazel and retire the limitation only after both indirect cases reject.
- `GAZELLE-BOOTSTRAP-RUNFILES`: directly invoking the root Gazelle binary failed
  Java runfile lookup. The standalone plugin executable, run through Bazel,
  regenerated corrected policies while toolchain registration was blocked.
- `BOUNDED-LOG-LINES`: one formatter log line contained a very large filename
  list. Subsequent inspection bounded both line count and line length. Keep
  that output discipline; no shared instruction or runtime changes are needed.

## Full repository matrix

All 28 nested-module commands passed, including the new plugin workspace.
Root build failed after building 6,495 targets. Root test failed before execution
because the Bazel JVM crashed with SIGSEGV. The Kotlin toolchain exception was
corrected after the root build analysis; focused validation must check it again.

| Workspace | Phase | Affected targets                           | Root cause                                                              | How to fix                                                                                              | Evidence                                            |
| --------- | ----- | ------------------------------------------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | --------------------------------------------------- |
| root      | build | Traefik and SSH Ansible role file packages | Existing production dependencies on infra CA outputs are now prohibited | Move deployment-specific certificate inputs to infrastructure-owned packaging; do not weaken the policy | `root.build.log`, transitive visibility diagnostics |
| root      | build | `//:requirements`, `//tools/bzlenv:*`      | Existing references to removed tool targets/packages                    | Update aliases and tool environment entries to their current owners                                     | `root.build.log`, loading errors                    |
| root      | build | `//tools/gradle:gradle-wrapper`            | Existing indexing of an empty Bash runfiles file list                   | Correct runfiles lookup in `tools/gradle/defs.bzl`                                                      | `root.build.log`, index-out-of-range at line 42     |
| root      | build | `//tools/kt:kotlin_toolchain_impl`         | External Kotlin rule needs the documented public build API exception    | Corrected the package directive; focused recheck pending                                                | `root.build.log`, toolchain alias diagnostic        |
| root      | test  | Normal `//...` test expansion              | JVM SIGSEGV before tests ran                                            | Retry after server restart; no host configuration changes                                               | `root.test.log`, server `jvm.out`                   |

The matrix below is copied from the checker's report. Logs remain local under
`out/full-repo-check/run.1566079259/logs/`. Normal `//...` excludes manual and
incompatible targets and does not enable optional configurations such as Harbor.

| Workspace                                            | Phase | Command                         | Result         | Duration | Log                                                                                                                                          |
| ---------------------------------------------------- | ----- | ------------------------------- | -------------- | -------: | -------------------------------------------------------------------------------------------------------------------------------------------- |
| root                                                 | build | `bazel_agent bazel build //...` | FAIL (exit 1)  |  1172.8s | [logs/root.build.log](logs/root.build.log)                                                                                                   |
| root                                                 | test  | `bazel_agent bazel test //...`  | FAIL (exit 37) |    26.9s | [logs/root.test.log](logs/root.test.log)                                                                                                     |
| projects/bazel_gazelle_enforce_transitive_visibility | build | `bazel_agent bazel build //...` | PASS (exit 0)  |     0.6s | [logs/projects__bazel_gazelle_enforce_transitive_visibility.build.log](logs/projects__bazel_gazelle_enforce_transitive_visibility.build.log) |
| projects/bazel_gazelle_enforce_transitive_visibility | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.3s | [logs/projects__bazel_gazelle_enforce_transitive_visibility.test.log](logs/projects__bazel_gazelle_enforce_transitive_visibility.test.log)   |
| tools/rules_binary_toolchain                         | build | `bazel_agent bazel build //...` | PASS (exit 0)  |     1.9s | [logs/tools__rules_binary_toolchain.build.log](logs/tools__rules_binary_toolchain.build.log)                                                 |
| tools/rules_binary_toolchain                         | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.3s | [logs/tools__rules_binary_toolchain.test.log](logs/tools__rules_binary_toolchain.test.log)                                                   |
| tools/rules_docs                                     | build | `bazel_agent bazel build //...` | PASS (exit 0)  |     1.1s | [logs/tools__rules_docs.build.log](logs/tools__rules_docs.build.log)                                                                         |
| tools/rules_docs                                     | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.3s | [logs/tools__rules_docs.test.log](logs/tools__rules_docs.test.log)                                                                           |
| tools/rules_docs_gazelle                             | build | `bazel_agent bazel build //...` | PASS (exit 0)  |    13.4s | [logs/tools__rules_docs_gazelle.build.log](logs/tools__rules_docs_gazelle.build.log)                                                         |
| tools/rules_docs_gazelle                             | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.4s | [logs/tools__rules_docs_gazelle.test.log](logs/tools__rules_docs_gazelle.test.log)                                                           |
| tools/rules_dnscontrol                               | build | `bazel_agent bazel build //...` | PASS (exit 0)  |     1.2s | [logs/tools__rules_dnscontrol.build.log](logs/tools__rules_dnscontrol.build.log)                                                             |
| tools/rules_dnscontrol                               | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.3s | [logs/tools__rules_dnscontrol.test.log](logs/tools__rules_dnscontrol.test.log)                                                               |
| tools/rules_hugo                                     | build | `bazel_agent bazel build //...` | PASS (exit 0)  |    20.5s | [logs/tools__rules_hugo.build.log](logs/tools__rules_hugo.build.log)                                                                         |
| tools/rules_hugo                                     | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.5s | [logs/tools__rules_hugo.test.log](logs/tools__rules_hugo.test.log)                                                                           |
| tools/rules_iso                                      | build | `bazel_agent bazel build //...` | PASS (exit 0)  |     1.3s | [logs/tools__rules_iso.build.log](logs/tools__rules_iso.build.log)                                                                           |
| tools/rules_iso                                      | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.3s | [logs/tools__rules_iso.test.log](logs/tools__rules_iso.test.log)                                                                             |
| tools/rules_openspec                                 | build | `bazel_agent bazel build //...` | PASS (exit 0)  |     4.6s | [logs/tools__rules_openspec.build.log](logs/tools__rules_openspec.build.log)                                                                 |
| tools/rules_openspec                                 | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.3s | [logs/tools__rules_openspec.test.log](logs/tools__rules_openspec.test.log)                                                                   |
| tools/rules_promptfoo                                | build | `bazel_agent bazel build //...` | PASS (exit 0)  |    16.9s | [logs/tools__rules_promptfoo.build.log](logs/tools__rules_promptfoo.build.log)                                                               |
| tools/rules_promptfoo                                | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.6s | [logs/tools__rules_promptfoo.test.log](logs/tools__rules_promptfoo.test.log)                                                                 |
| tools/rules_promptfoo_gazelle                        | build | `bazel_agent bazel build //...` | PASS (exit 0)  |    13.2s | [logs/tools__rules_promptfoo_gazelle.build.log](logs/tools__rules_promptfoo_gazelle.build.log)                                               |
| tools/rules_promptfoo_gazelle                        | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     1.9s | [logs/tools__rules_promptfoo_gazelle.test.log](logs/tools__rules_promptfoo_gazelle.test.log)                                                 |
| tools/rules_skills                                   | build | `bazel_agent bazel build //...` | PASS (exit 0)  |    14.7s | [logs/tools__rules_skills.build.log](logs/tools__rules_skills.build.log)                                                                     |
| tools/rules_skills                                   | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     1.0s | [logs/tools__rules_skills.test.log](logs/tools__rules_skills.test.log)                                                                       |
| tools/rules_skill_gazelle                            | build | `bazel_agent bazel build //...` | PASS (exit 0)  |    15.4s | [logs/tools__rules_skill_gazelle.build.log](logs/tools__rules_skill_gazelle.build.log)                                                       |
| tools/rules_skill_gazelle                            | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.5s | [logs/tools__rules_skill_gazelle.test.log](logs/tools__rules_skill_gazelle.test.log)                                                         |
| tools/rules_template                                 | build | `bazel_agent bazel build //...` | PASS (exit 0)  |    13.5s | [logs/tools__rules_template.build.log](logs/tools__rules_template.build.log)                                                                 |
| tools/rules_template                                 | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.5s | [logs/tools__rules_template.test.log](logs/tools__rules_template.test.log)                                                                   |
| tools/rules_terraform                                | build | `bazel_agent bazel build //...` | PASS (exit 0)  |     2.6s | [logs/tools__rules_terraform.build.log](logs/tools__rules_terraform.build.log)                                                               |
| tools/rules_terraform                                | test  | `bazel_agent bazel test //...`  | PASS (exit 0)  |     0.5s | [logs/tools__rules_terraform.test.log](logs/tools__rules_terraform.test.log)                                                                 |

## Public certificate ownership follow-up

The user selected independent Terraform roots under `data/ssh/tf` and
`data/ssl/alwaldend.com/tf`; infrastructure Terraform must not refresh their
files. The four checked-in public snapshots have been moved with identical
SHA-256 hashes, and every current label and file-path consumer was updated.
The former Vault `local_file` resources now have non-destructive `removed`
blocks. No live Terraform or Vault operation ran.

On 2026-10-28, before the new refresh roots were implemented:

- `bazel_agent bazel query 'filter("^//infra/", deps(//projects/ansible_collection:collections))'` passed with no targets.
- `bazel_agent bazel build //projects/ansible_collection:collections //infra/forgejo/ansible:ansible_bin //infra/ingress/ansible:ansible_bin //infra/vault/ansible:ansible_bin //infra/download:al` passed.
- All four Kustomize certificate references resolve to the data-owned PKI snapshot.

The repeatable commands, build log, query output, and original file hashes are
under `out/transitive-visibility/certificate-migration/`. This is partial
working-tree validation, not final candidate or publication evidence. The
refresh roots now use the existing Vault and Local providers. The proposed
HTTP dependency is unnecessary and was not installed or declared. The user
requested standard Terraform checks rather than a custom integration harness.

The follow-up implementation passed the two data roots' native
`tf_tests.fmt_test` and `validate_test`, plus the existing Vault Terraform
format test (five tests). Both roots' `tf.plan` and `tf.apply` wrappers build,
as do the collection and representative Forgejo, ingress, Vault, and download
consumers. The combined dependency query for the collection and both refresh
commands returns no infra targets. The Forgejo Ansible runfiles manifest was
checked against the unchanged source bytes for the Traefik, SSH, and server
CA files. These checks did not run a live refresh. Native Vault data sources
read only the `certificate` and `public_key` fields; the HTTP-provider request
is withdrawn.

The user also requested a rebase onto current master. Delivery receipts and
post-rebase gate logs under `out/transitive-visibility/certificate-migration/`
record the resulting candidate and publication status. The unrelated missing
tool labels and Gradle runfiles defect remain outside this certificate fix.

Session review: speculative provider selection and a proposed custom test
harness added avoidable work. Pinned provider data-source documentation and
native Terraform validation were sufficient for this configuration change.
Use discovered file paths for batched reads; a nonexistent path rejects the
whole Cordis read batch. No shared skill or host configuration was changed.

The rebase onto `7dea6a8b5f0f94047e0a9e58a784e705345825e8` completed, producing
candidate `b5ee00c6987a25b852bdbbd411c7bb78806d465f`. Its first validation check
passed all five Terraform tests but found one formatting issue where the
rebased Mermaid BUILD file placed `package()` between rules. The configured
formatter moved that declaration above the rules without semantic changes.
The amended candidate and subsequent gate results are in the delivery receipts.

The collection's unfiltered graph contains only the `tools/nogo` build-analyzer
targets under tools. With `--notool_deps`, there are no tools targets. These
execution-tool edges and `.bzl` loads are not artifact-boundary violations.
Source inspection of Bazel 9.3.1 nevertheless shows no execution-attribute
exemption in the native transitive validator; the plugin README now records
this additional scope limitation rather than promising exact edge filtering.

## Bzlenv repair

The user authorized fixing the stale tool references uncovered by lint. The
eight deleted tools packages now resolve to their existing public binary
targets in third_party. No dependencies or compatibility packages were added.

A generated-output check found that adding the executable directories to PATH
exposed wrapper filenames such as hugo_native_binary.exe and ignored the
configured command names. None of the eight expected names resolved in those
directories. The rule now adds a dedicated runfiles directory with symlinks
named from the tools map, while retaining the executables' runtime files.

Focused lint and the bzlenv build passed. A build with --build_runfile_links
materialized the runfiles tree for inspection, as required by the repository's
default lazy-runfiles setting. All 43 configured command aliases resolve to
executable files. Non-mutating version commands passed through the generated
aliases for Blender, Buf, DNSControl, Hugo, Minisign, ORAS, ShellCheck, and
Taplo. The check did not activate a host environment or load a .env file.
Evidence is under out/transitive-visibility/certificate-migration in
bzlenv-command-names-before.json, bzlenv-command-names-after.json, and the
corresponding build and version logs.

The preceding candidate passed root quality, native Terraform checks, the
plugin test, representative builds, and nested-module documentation builds.
Delivery lint also reported the separate existing Gradle wrapper's empty
Bash-runfiles indexing error. The bzlenv repair does not change Gradle. Final
candidate validation and publication state are recorded in delivery receipts.

## Gradle wrapper repair

The user authorized fixing the remaining Gradle blocker. The wrapper now uses
the existing rules_shell runfiles provider and merges its declared files and
default runfiles without indexing an empty default file list. Focused semantic
lint passes.

Execution then exposed the hardcoded obsolete +_repo_rules+net_gradle_gradle
repository name. The launcher now derives Gradle and Java runfile locations
from the declared Gradle file and Java toolchain provider. An offline version
command passed and reported Gradle 8.13 and the pinned Azul Java 21.0.9 runtime.
No dependency or pin was added or changed; caches were confined to the ignored
task directory. The generated launcher and runfiles manifest were inspected.

Evidence is under out/transitive-visibility/certificate-migration:
gradle-lint.log records the analysis fix, gradle-version.log preserves the
subsequent path failure, and gradle-version-fixed.log plus
gradle-fixed-lint.log record successful startup and lint after both repairs.
Final candidate checks and publication state are owned by delivery receipts.

## Publication and review correction

The candidate passed all 29 delivery checks and was published as PR 210.
Publication and validation must inherit the same process environment; the
first publication guard rejected environment drift before pushing. Repeating
validation and publication within one process passed all guards. The trusted
receipts and logs remain under out/transitive-visibility/certificate-migration.

PR review identified that both certificate refresh examples omitted the
backend key in AL_TF_BACKEND_CONFIG_path. The Terraform launcher passes each
matching environment value verbatim as a backend-config argument; the suffix
does not supply the key. Both examples now set the value to
path=$PWD/out/<owner>/terraform.tfstate. This correction changes documentation
only and performs no live Terraform operation.

## Full dependency capture cleanup

Removed the unused downloader OCI image and its image-only build test, then
removed the sole Alpine declaration and regenerated MODULE.bazel through
tools/bazel_module. The surviving downloader build test and module freshness
test pass. The obsolete requirements alias had no consumers; the root README
now names the existing uv_lock.update command. That command target and the
Python configuration target both build without changing dependency locks.

Complete root and external dependency captures contain 25,602 targets with
43,763 edges and 256,520 targets with 324,220 edges, respectively. Independent
comparison of the captured input edges with the diagram edges found no missing
in-scope references. Both scopes include files, preserve disconnected targets,
omit execution-tool and visibility edges, and have no depth limit or grouping.
Six dark 8192-by-8192 WebPs cover force-directed, circular and matrix layouts.
Each decoded image was inspected; render counters match the full captured
target and edge counts. Raster annotations are sparse at this scale; companion
Mermaid files preserve every full label. All captures, scripts and diagram
artifacts remain ignored under out/transitive-visibility and are not committed.
