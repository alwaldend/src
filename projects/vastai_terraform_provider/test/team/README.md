---
title: Vast.ai provider test root
description: Non-compute provider validation in a Terraform-created test team
---

This root manages three test objects in the team created by
the infra/vastai bootstrap root: a generated `test_` API key, a
`TEST_VASTAI_TERRAFORM_PROVIDER` environment marker, and an unassigned
`test_vastai_terraform_provider` role requesting current-account reads.
Vast.ai adds baseline role permissions; `effective_permissions` records the
actual grants. SSH keys require personal context and cannot be created in a team.
It has no compute, storage rental, credit transfer, team creation, or member
invitation resources. It adds no team-ID guard.

The marker value is non-secret. Generated keys remain sensitive in Terraform
state and are not outputs. The root uses a separate Vault HTTP
backend under the provider project's own `tf_backend/test` entry.

The project authenticates as `src_projects_vastai_terraform_provider`.
Its `vastai=1` injector reads the Terraform-generated token from its own
`vastai` credential entry into `VAST_API_KEY` for provider authentication;
`tf=test` selects the isolated backend.
The provider is built from source and installed
from its declared Bazel ZIP; no Terraform Registry publication is required.

```sh
bazel test //projects/vastai_terraform_provider/test/team:tf_tests.fmt_test
bazel run //projects/vastai_terraform_provider/test/team:tf.plan -- -out=/absolute/private/path/create.tfplan
bazel run //projects/vastai_terraform_provider/test/team:tf.show -- /absolute/private/path/create.tfplan
bazel run //projects/vastai_terraform_provider/test/team:tf.apply -- /absolute/private/path/create.tfplan
```

For an update, plan with `-var=test_marker=updated` and apply that saved plan.
For cleanup, use `tf.destroy_plan -- -out=/absolute/private/path/delete.tfplan`,
inspect the deletion plan, and pass it to `tf.apply`. Apply always requires
a saved plan. These commands contact the live API and mutate only when an
authorized saved plan is applied. Keep plans private and remove them after use.
