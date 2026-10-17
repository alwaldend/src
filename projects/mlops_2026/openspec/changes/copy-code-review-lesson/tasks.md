## 1. Import

- [x] 1.1 Copy the seven non-CSV lesson files unchanged and retain the upstream MIT notice; verify the file inventory and SHA-256 hashes against the pinned archive.
- [x] 1.2 Copy the first 100 lines of upstream `data/uber.csv`, including its header; verify exact prefix bytes and line count, with lesson BUILD files and dependency changes absent.
- [x] 1.3 Protect imported files from first-party formatting and import checks; verify repository quality checks preserve their bytes.

## 2. Delivery

- [x] 2.1 Verify strict change validation and the import inventory; select repository quality and semantic lint gates for publication through `repo-delivery`.
