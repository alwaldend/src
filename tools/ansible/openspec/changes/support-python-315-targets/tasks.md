## 1. Compatibility

- [x] 1.1 Update existing Ansible pins and regenerate the lock; verify no new external package is introduced and generation is stable.
- [x] 1.2 Build and syntax-check representative deployments, run semantic lint, and verify live fact gathering on both development hosts.
- [x] 1.3 Package the pinned ansible CLI for native SSH askpass and verify the real password-authenticated appliance connection without sshpass.

- [x] 1.4 Remove the wrapper PATH override and verify native askpass resolves its packaged sibling helper through the isolated transport integration. Record the refused live XO connection separately.
