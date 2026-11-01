## Inventory reuse

- [x] Expose the existing inventory with narrow security-package visibility.
- [x] Verify the export, unchanged inventory bytes, and patch whitespace.
- [x] Validate the shared consumer through the parent implementation.

The packaged `//infra/security/ansible:ansible` consumer built successfully
and its composed inventories were used for the authorized maintenance run.
Inventory membership remains owned by the source inventory.

- [x] Use the shared fleet inventory consumer package group without widening owner access; verify consumer packaging.
