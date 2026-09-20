## 1. Define compatible role inputs

- [ ] 1.1 Inventory existing consumers and document the supported distribution,
      packages, service precedence, and opt-in management inputs; verify fixtures
      cover documented service selection, legacy overrides, and conflicts.
- [ ] 1.2 Verify required PostgreSQL modules and Python driver against the pinned
      dependencies; update dependencies through the owning lock workflow if needed
      and verify packaged module availability without a live host.

## 2. Implement reusable database management

- [ ] 2.1 Correct service-variable handling and add guarded initialization;
      verify empty, initialized, missing-storage, unrecognized-data, and
      incompatible-version cases preserve existing data and reject unsafe startup.
- [ ] 2.2 Add opt-in listeners and ordered authentication rules with bounded
      handlers; verify rendered local-access configuration, unchanged reruns, and
      redaction of secret-bearing inputs and diffs using non-secret fixtures.
- [ ] 2.3 Reconcile declared application roles, owned databases, and extensions;
      verify idempotence, unavailable-extension failure, absence of unintended
      superuser grants, and preservation of undeclared objects.

## 3. Validate compatibility and hand off integration

- [ ] 3.1 Run collection packaging, role syntax/rendering, semantic lint,
      OpenSpec, and repository quality checks; record the exact candidate and
      verify an initialized-host fixture with only existing inputs remains
      unchanged apart from requested package/service convergence.
- [ ] 3.2 Document input migration and a disposable-host acceptance procedure;
      verify it covers fresh initialization, reboot, unchanged rerun, application
      login, invalid-password and remote-access rejection, and data preservation.
      Distinguish unexecuted procedures from separately authorized live evidence.
- [ ] 3.3 Document the tested role inputs and acceptance evidence for the linked
      Nexus consumer; verify its integration references this owner contract and
      keeps Nexus-specific defaults outside the reusable role.
