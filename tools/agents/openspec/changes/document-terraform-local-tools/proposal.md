## Why

The infrastructure skill needs a supported pattern for small API bootstraps
when a pinned provider lacks a resource, without automatically adding a provider
or routing the operation through Ansible.

## What Changes

Document packaged local tools invoked through `terraform_data` and `local-exec`,
including execution timing, secret handling, retry behavior, lifecycle limits,
and offline acceptance checks. Add representative offline evaluation cases.

## Impact

The canonical `repo-infra` Terraform reference and its evaluation configuration.
This documentation correction changes no deployed infrastructure.
