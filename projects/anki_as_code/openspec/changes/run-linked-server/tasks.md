## 1. Linked worker

- [x] 1.1 Write installed CLI and server E2E cases before production changes, covering help, invalid arguments, stdio, unavailable temporary storage, interrupt handling and repeatable native/JSON artifacts.
- [x] 1.2 Add server argument forwarding and explicit self-worker configuration, preserving collection behavior; verify installed execution and existing application workflows.

- [x] 1.3 Remove production runfiles/PATH discovery, require explicit application/collection worker configuration, update native fixtures, and verify missing configuration fails without opening a collection.

- [x] 1.4 Consume the ankiserver and ankiclient Go packages and retain native CLI delegation.

- [x] 1.5 Expose generic server run forwarding, add server serve as a delegate to server run serve, self-launch with server serve, and verify root/serve help, native validation and interrupts through installed E2E.

## 2. Delivery

- [x] 2.1 Update current docs, validate OpenSpec, inspect the optimized installed executable and select quality/lint/native gates for independent master-based delivery; receipts own publication.
