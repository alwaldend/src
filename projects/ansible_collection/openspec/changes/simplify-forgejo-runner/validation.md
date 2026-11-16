# Validation

The local acceptance fixture executes the packaged registration tasks and renders
the actual configuration and unit. Its sanitized `registration.json` records
initial registration, repeat deployment without bootstrap inputs, `0600` mode,
managed paths, advanced settings, missing-input rejection, suppressed failure
output and the registration deadline. The rendered YAML and unit are retained
alongside it under the test's Bazel undeclared outputs.

| Check | Result |
| --- | --- |
| `git diff --check` | Passed |
| `bazel run //:format` and pinned Ruff formatting | Passed; no unrelated source changes |
| `bazel test //:repo_quality_test //projects/ansible_collection/test/forgejo_runner:registration_test //infra/src/openspec/validation:projects_ansible_collection_validate_test` | 29 selected tests passed |
| `bazel build --config=lint //projects/ansible_collection/roles/forgejo_runner:role //projects/ansible_collection/test/forgejo_runner:registration_test //projects/ansible_collection/openspec:source //tools:transitive_consumers` | Passed |
| `bazel build //projects/ansible_collection/roles/forgejo_runner:role //infra/dev/ansible:ansible_bin` | Passed |
| `bazel run //infra/dev/ansible:syntax` | Passed; dedicated offline syntax entry point |

Local review covers the complete prepared diff against base
`de17c7905c559aef7957d5d4ea6cfeebd3fd5bcd`: no actionable findings. The exact
candidate, final review and publication state are retained in task-local delivery
receipts under `out/simplify-forgejo-runner/`.

The fixture uses a disposable registration executable and a recording handler.
It does not prove live Forgejo protocol behavior, real systemd activation,
account creation, package installation or cross-account ACL transfer. The
`infra/dev` syntax check does not activate its disabled runner role. No live
infrastructure or credentials were changed.

Authentication was traced through the current role, `infra/dev/BUILD.bazel`,
`infra/vault/tf/approle_src_infra_dev.tf` and the existing consolidation deployment
record. The old runner identity and token lookup were already retired. The
worker requires its registration file after bootstrap, not Vault/OIDC login.
Shared provisioning and service-management authentication remains separately owned.

The fixture's initial setup failures were resolved locally: add the exact test
package to the tool consumer allowlist, replace the handler only inside a
disposable role copy, and use a private Python executable lookup to avoid Linux's
shebang length limit. Reproduction uses the checked-in Bazel target; no host
configuration or external dependency was added.
