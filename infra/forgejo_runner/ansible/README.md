---
title: Ansible
description: Register and configure Forgejo runner hosts
---

After VM and DNS provisioning, run the controller-only registration playbook
and then the packaged deployment:

```sh
bazel_agent bazel run //infra/forgejo_runner/ansible:registration
bazel_agent bazel run //infra/forgejo_runner/ansible:ansible
```

Registration authenticates to Forgejo with a temporary controller token from
Vault OIDC, reads the repository registration credential, and stores it under
the component AppRole's config path. A compare-and-set write preserves other
fields and rejects concurrent changes. Credential tasks suppress output.

Deployment uses the component AppRole for signed SSH access and credential
injection. It waits for cloud-init, applies `host` and `dev_vm`, installs and
registers the runner, installs Bazel tools for that account, and checks the
service. Only the inventory's `secure` host is deployed. Existing registration
is preserved on subsequent runs and service changes restart through a handler.

The initial image has no signed SSH host certificate. For the first deployment
of a new VM, create a private ignored `out/forgejo_runner` directory and use
task-local trust on first use:

```sh
bazel_agent bazel run //infra/forgejo_runner/ansible:ansible -- \
  --extra-vars "ansible_ssh_common_args='-o StrictHostKeyChecking=accept-new -o UserKnownHostsFile=$PWD/out/forgejo_runner/bootstrap_known_hosts'"
```

The `host` role then configures Vault-signed host certificates. Repeat the
deployment with the normal command to verify those certificates. Deployment
packages the public [server CA](../../../data/ssh/README.md#ca_servers) and
prepares a private operational cache at
`$BUILD_WORKSPACE_DIRECTORY/out/forgejo_runner`. Its `known_hosts` file trusts
that CA only for inventory hostnames. Normal connections require strict host
verification using that file; the initial explicit SSH override above remains
available for bootstrap. Set `runner_controller_state_dir` to override the
cache directory. No persistent controller SSH configuration is changed.

After committing the candidate, the controller-only CI validation deployment
creates the branch named in that commit's [CI configuration](../ci.json):

```sh
bazel_agent bazel run //infra/forgejo_runner/ansible:ci_validation -- \
  --extra-vars "runner_ci_workspace=$PWD runner_ci_commit=<full-commit-sha>"
```

It checks the live `releases/*` protection rule, creates the branch with an
empty expected-value lease, and verifies the published SHA and effective
protection. An existing branch at the same commit resumes read-only validation;
a different commit requires the explicit retry option below. It then polls
that commit's smoke workflow and jobs for up to `runner_ci_timeout` seconds
(default 240, maximum 300), reporting only their identifiers, statuses, and run
URL. A created branch
remains available for inspection if CI fails; this target never rewrites history
or deletes the branch. CI Vault authentication accepts the protected default
branch and the exact validation branch according to the
[Vault role](../../vault/tf/forgejo_ci.tf).

To test a committed fix on the existing validation branch, also pass
`runner_ci_previous_commit=<full-observed-branch-sha>`. The helper verifies
that exact protected tip, requires it to be an ancestor of the new candidate,
and pushes with an explicit lease for the previous SHA. A stale tip or a
candidate that would rewrite history is rejected before pushing.

The Python standard-library helper runs inside the packaged Ansible controller,
where it can invoke Git and poll the Forgejo API without writing an auth file
or putting credentials in module arguments. It inherits the temporary Forgejo
token supplied by the same named service authentication as registration. The
token enters Git only through its child environment; credential tasks suppress
output, and the helper emits a bounded receipt without response bodies or logs.

`ansible_bin` is the packaged entry point for offline `--syntax-check` and
fixture rendering. The normal `ansible` target is a live deployment.
