## 1. Prepare the recovery workflow

- [x] 1.1 Add scoped, fail-closed certificate tasks and verify the packaged
      Ansible task listing contains only issuance, certificate/key writes, and reload.
- [x] 1.2 Document login prerequisites, the temporary override, and restored TLS
      verification; review the instructions against the existing wrappers.
- [x] 1.3 Validate the OpenSpec change, package builds, repository quality, and
      affected semantic lint without contacting inventory hosts.

## 2. Execute and verify the authorized recovery

- [x] 2.1 Establish a valid bootstrap login without printing credentials; confirm
      the existing injector can authenticate and start Ansible.
- [x] 2.2 Run the documented Ansible renewal for host1 and confirm successful
      issuance, both certificate writes, and service reload in restricted evidence.
- [x] 2.3 Verify the served certificate's new validity and the unsealed Vault
      service with normal TLS verification restored.
