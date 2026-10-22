## 1. Linked entry point

- [x] 1.1 Specify linked/native CLI failure and success scenarios before production changes, with repeatable output artifacts.
- [x] 1.2 Extract shared Rust server logic and link one C ABI entry point into Go using existing toolchains; verify standalone transport/backend tests and linked CLI cases.

- [x] 1.3 Name the Go packages ankiserver and ankiclient, updating import paths, targets, consumers and documentation while retaining existing lifecycle behavior.

- [x] 1.4 Add native serve dispatch and separate root/serve help; specify failure cases before code and verify standalone/linked native transport and lifecycle E2E.

- [x] 1.5 Remove the implicit standalone-worker runtime input from the client library; native tests declare their own executable fixtures.

## 2. Delivery

- [x] 2.1 Document the callable server and build boundary, validate OpenSpec, inspect optimized installed output and select exact-candidate quality/lint/native gates for a master-based PR; receipts own publication.
