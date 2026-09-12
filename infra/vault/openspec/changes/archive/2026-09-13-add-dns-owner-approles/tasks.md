## 1. Owner identities

- [x] 1.1 Add the missing per-owner AppRole modules and isolated DNS group;
      verify complete source-inventory coverage and absence of extra group policies.
- [x] 1.2 Add provider read policies and attach them to existing component
      identities; verify global-only roles have no dc1 access and existing addresses
      remain unchanged.

## 2. Offline acceptance

- [x] 2.1 Package all local module sources and run the owning Terraform format
      test plus BUILD formatting checks without contacting Vault or DNS providers.
- [x] 2.2 Verify linked OpenSpec artifacts with the pinned strict validator and
      record implementation evidence before archiving the completed source change.
