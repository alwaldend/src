## 1. Configure native HTTPS mirrors

- [x] 1.1 Declare the GitHub and GitLab destinations in the repository catalog.
- [x] 1.2 Package the Forgejo API bootstrap in a dedicated Terraform root.
- [x] 1.3 Store destination credentials in Vault and verify both live HTTPS mirrors.
- [x] 1.4 Remove the unused App and SSH mirror configuration paths.
- [x] 1.5 Apply the two pending Terraform state-address moves.
- [x] 1.6 Document the push-mirroring flow in the Forgejo architecture docs.

## 2. Restrict mirrored refs

- [x] 2.1 Configure the trunk and release branch filter on new and existing mirrors.
- [x] 2.2 Verify the bootstrap and filtered push behavior offline.
- [x] 2.3 Apply and verify the filters and pruning on both destinations.
