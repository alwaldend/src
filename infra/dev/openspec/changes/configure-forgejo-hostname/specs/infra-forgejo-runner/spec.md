## MODIFIED Requirements

### Requirement: Packaged setup and deployment configuration

The owner SHALL expose Terraform setup and packaged Ansible entry points
through its `al` configuration and Vault injection dependencies. The runner
configuration SHALL reference `https://forgejo.alwaldend.com`, `/dev/sdb`, and
an environment-supplied `FORGEJO_RUNNER_TOKEN` rather than a checked-in
registration token.

Sources: [setup BUILD](../../../../../tf_setup/BUILD.bazel),
[Ansible BUILD](../../../../../ansible/BUILD.bazel), and
[runner variables](../../../../../ansible/group_vars/all.yaml).

#### Scenario: Resolve the registration credential source

- **WHEN** a maintainer inspects the runner configuration scaffold
- **THEN** its token value is an environment lookup and its instance URL is
  the canonical Forgejo service
- **AND** the token reference alone does not enable the inactive runner role
