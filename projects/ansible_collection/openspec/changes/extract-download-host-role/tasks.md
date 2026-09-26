## 1. Role extraction

- [x] 1.1 Move tasks, handlers, defaults, and maintenance templates into the role;
      compare moved content with source revision `adce867e` and build the collection.
- [x] 1.2 Verify both consumer inventories with Ansible syntax/task expansion,
      reject a missing disk before host mutation, and render both packaged units.
- [x] 1.3 Prepare an eight-file consumer conversion patch and verify it applies to
      PR #102 with `git apply --check`.

- [x] 1.4 Remove the reviewed capacity assertion and discovery; inspect the diff
      to verify attachment validation remains before host setup.

- [x] 1.5 Create publication roots only and protect privileged directory operations;
      verify the sentinel-symlink fixture and unchanged second execution pass.
- [x] 1.6 Remove Traefik's mount override and update the prepared consumer to use
      system-disk state; inspect role tasks, shared service template, and patch.
- [x] 1.7 Document the SELinux labeling rationale using the Nginx deployment guide.

- [x] 1.8 Remove the single-host restriction and obsolete service overrides;
      verify two hosts reach the attachment precondition before host setup.
- [x] 1.9 Move filesystem tasks into `filesystem.yaml` and use the pinned
      `btrfs_info` module; verify direct-device and symlink UUID selection,
      missing-match rejection, and both packaged consumer task expansions.
- [x] 1.10 Compare the moved directory tasks with the prior sentinel fixture and
      verify the updated consumer patch still applies to PR #102.

## 2. Consumer follow-up

- [ ] 2.1 After this role merges, rebase PR #102 onto trunk, apply its prepared
      conversion, and validate/publish the role-based consumer against trunk.
