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
registers the runner, installs the repository Android prerequisites and Bazel
tools for that account, and checks the
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

`ansible_bin` is the packaged entry point for offline `--syntax-check` and
fixture rendering. The normal `ansible` target is a live deployment.

The `android` and `bazel` deployment tags install the prerequisites used even by
non-Android builds: root module toolchain registration resolves the SDK and
NDK during analysis. Package versions remain owned by `tools/android/packages.txt`;
the runner rc selects the installed SDK and the managed `ndk/current` link.
