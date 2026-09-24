## 1. Define and package the role

- [ ] 1.1 Write failure scenarios for invalid configuration, unwanted listeners, unchanged reruns, and content preservation before implementation; verify coverage against the role specification.
- [ ] 1.2 Add the native Nginx role and its documented consumer inputs using collection conventions; verify the selected guest distribution can install the declared package.
- [ ] 1.3 Add Bazel role packaging and aggregate collection membership; verify the built collection includes the role once at the expected path.

## 2. Validate configuration and lifecycle

- [ ] 2.1 Implement complete-candidate validation before activation; verify a rejected candidate leaves the running service and persistent configuration unchanged.
- [ ] 2.2 Implement initial enable/start and change-driven reload handlers; verify an unchanged second run performs no restart or reload.
- [ ] 2.3 Exercise the role with the download consumer's loopback configuration; verify JSON listing, normal file serving, and selected-site rendering through HTTP.
- [ ] 2.4 Run role/collection packaging and semantic checks; retain the isolated fixture's configuration, commands, and HTTP results as repeatable evidence.
