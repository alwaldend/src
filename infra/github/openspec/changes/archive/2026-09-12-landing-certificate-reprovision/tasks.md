## 1. Terraform

- [x] 1.1 Add `certificate_reprovision_projects` with a description that distinguishes it from first-creation bootstrap
- [x] 1.2 Omit the Pages block when a project is in either `bootstrap_projects` or `certificate_reprovision_projects`
- [x] 1.3 Confirm with a read-only `tf.plan` filtered to the recovery variable that only the named project's Pages block changes

## 2. Documentation

- [x] 2.1 Document the two filtered applies and HTTPS verification in `infra/github/tf/README.md`
- [x] 2.2 Keep the `add-project-site` skill's rollout guidance consistent with the separately named variable

## 3. Verification

- [x] 3.1 `bazel run //infra/github/tf:tf.fmt` and `bazel test //infra/github/tf:tf_tests.fmt_test`
- [x] 3.2 `bazel run //infra/github/tf:tf.plan` with both sets empty reports no unrelated drift
- [x] 3.3 `OPENSPEC_PROJECT=infra/github bazel run //tools/openspec -- validate landing-certificate-reprovision --type change --strict --no-interactive`
