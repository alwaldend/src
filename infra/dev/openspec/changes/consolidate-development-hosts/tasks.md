## 1. Implementation

- [x] 1.1 Consolidate owner and maintained OpenSpec work; verify old source paths have no active consumers.
- [x] 1.2 Align DNS, routes, certificate scope and state partitions; verify offline Terraform and inventory outputs.
- [x] 1.3 Update consumers and documentation; run syntax, packaging, quality and semantic lint checks.
- [x] 1.4 Replace nested deployment orchestration with one standard Ansible invocation and a scoped deployment identity; verify packaged syntax, selected hosts and credential scope.

- [x] 1.5 Remove the legacy runner VM, DNS and inventory declarations; verify only secure1 and dev1 remain.
- [x] 1.6 Merge DNS into one root and shared identity, remove unused runner credentials and old inventory names, and document the pending state handoff; verify packaged configuration and host separation offline.
