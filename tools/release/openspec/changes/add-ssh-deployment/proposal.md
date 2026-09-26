## Why

Release artifacts need direct SSH publication to the download hosts. Static
websites use public archive uploads plus publication of their extracted
content and selection, so this belongs in the existing release tool's deployment model.

## What Changes

- Extend deployment metadata, Bazel rule inputs, and the existing `deploy`
  command with an SSH deployment type alongside OCI deployment. Use rsync
  over the existing administrator SSH access, with sudo to the content account.
- Upload ordinary files without extraction into
  `projects/<project>/releases/<version>/` on the selected host.
- Make upload, extraction of a published website archive, and atomic selection
  independently selectable steps of deploy.
- Make redeployment of an existing site release select that release by
  changing the link; do not add a rollback command.
- Preserve local release-package layout and OCI behavior. SSH-only deployment
  must not require an OCI executable or Ansible.

## Capabilities

### New Capabilities

- `ssh-release-deployment`: Direct authenticated SSH upload, site archive
  activation, and existing-release selection.

### Modified Capabilities

None; this owner currently has no baseline OpenSpec specifications.

## Impact

The existing deployment protobuf, release rules, CLI, deployer, and generated
release metadata need coordinated changes. The [hosting plan](../../../../../infra/download/openspec/changes/add-static-hosting/proposal.md)
owns remote storage and access. The [website plan](../../../../../projects/alwaldend.com/openspec/changes/add-download-browser-and-ssh-hosting/proposal.md)
consumes the new deployment type and displays release directories.

Acceptance must use an isolated SSH server and real HTTP downloads, cover
publication failures and reselection, and retain a repeatable result artifact.
This plan does not upload any release or change a running site.
