## 1. Implementation

- [x] 1.1 Remove AL cached-tool execution and declarations; verify AL and runner tests.
- [x] 1.2 Add atomic standalone-binary installation and repo_delivery and bazel_agent install targets; verify isolated installation, replacement and failure behavior.
- [x] 1.3 Update invocation guidance; verify that active instructions no longer recommend cached execution.

## 2. Validation and delivery

- [x] 2.1 Execute the installed tool and preserve a repeatable JSON acceptance artifact.
- [x] 2.2 Pass affected tests, repository quality, semantic lint and OpenSpec validation.

Delivery follows repo-delivery; its receipts own commit, push and
pull-request state.
