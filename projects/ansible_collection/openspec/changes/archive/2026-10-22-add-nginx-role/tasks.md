## 1. Define and package the role

- [x] 1.1 Write failure scenarios for invalid configuration, unwanted listeners, unchanged reruns, and content preservation before implementation; verify coverage against the role specification.
- [x] 1.2 Add the native Nginx role and its documented consumer inputs using collection conventions; verify the selected guest distribution can install the declared package.
- [x] 1.3 Add Bazel role packaging and aggregate collection membership; verify the built collection includes the role once at the expected path.

## 2. Validate configuration and lifecycle

- [x] 2.1 Implement complete-candidate validation before activation; verify a rejected candidate leaves the running service and persistent configuration unchanged.
- [x] 2.2 Implement initial enable/start and change-driven reload handlers; verify an unchanged second run performs no restart or reload.
- [x] 2.3 Exercise the role with the download consumer's loopback configuration; verify JSON listing, normal file serving, and selected-site rendering through HTTP.
- [x] 2.4 Run role/collection packaging and semantic checks; retain the isolated fixture's configuration, commands, and HTTP results as repeatable evidence.

Evidence: the download component's
[validation record](https://github.com/alwaldend/src/blob/39b2dd9cea5fea84705a36cc6003fb5f46ba0547/infra/download/openspec/changes/add-static-hosting/evidence.md)
and [repeatable native fixture](https://github.com/alwaldend/src/blob/39b2dd9cea5fea84705a36cc6003fb5f46ba0547/infra/download/test/README.md).
Package installation was exercised on Fedora 44; the selected Fedora 43 cloud
image still needs its authorized rollout checks.

The role implementation is split from the consumer PR. The immutable fixture
links above preserve the tested consumer and evidence before the split. The
shared Traefik tasks also suppress EAB credentials and rendered private
configuration in Ansible logs.
