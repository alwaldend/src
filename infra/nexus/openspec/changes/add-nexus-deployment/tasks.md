## 1. Resolve deployment inputs and package boundaries

- [ ] 1.1 Select and record compatible Nexus CE, bundled Java, PostgreSQL, and
      `datadrivers/nexus` versions; verify official checksums and required provider
      resources against that exact combination.
- [ ] 1.2 Define hostname, address, DNS-view, VM-sizing, and storage inputs without
      guessing live allocations; verify that required inputs and their validation
      are documented in the component README.
- [ ] 1.3 Add executable component packaging alongside the existing planning
      scaffold; verify strict owner validation continues to pass and runtime
      targets keep repository-internal visibility.
- [ ] 1.4 Add Nexus distribution and provider pins through their owning workflows;
      verify checksum enforcement, generated metadata freshness, and provider
      installation from the declared Bazel mirror.
- [ ] 1.5 Before task 3.1, implement and validate the linked
      [PostgreSQL companion change](../../../../../projects/ansible_collection/openspec/changes/extend-native-postgresql-role/tasks.md)
      under its collection owner; verify its reusable contract and compatibility
      acceptance evidence are available for Nexus integration.

## 2. Wire identity and provisioning

- [ ] 2.1 Define `src_infra_nexus`, required group membership, scoped DNS access,
      and SSH/ingress certificate prerequisites in their existing Vault owners;
      verify policy fixtures grant only the intended secret paths and capabilities.
- [ ] 2.2 Add Nexus XO resource-set inventory using the existing identity flow;
      verify named inventory and pending-first-login behavior without adding a
      lookup requiring the new VM to exist before provisioning.
- [ ] 2.3 Add `al.lua` stage selectors, separate HTTP state paths, and credential
      injection; verify non-secret fixtures and packaged plugin/run-label alignment.
- [ ] 2.4 Implement `tf_setup` for one VM, guarded data disks, network, and canonical
      DNS records; verify Terraform formatting/schema, DNS lint, source ownership,
      and consistent rendered guest hostname/address inputs without a live plan.

## 3. Deploy native services

- [ ] 3.1 Configure the PostgreSQL role delivered by the companion change with
      Nexus-specific initialization, loopback authentication, database ownership,
      and `pg_trgm` inputs; verify fresh-host and initialized-host integration
      fixtures against the collection-owned contract.
- [ ] 3.2 Add the reusable Nexus role with pinned archive verification, bundled
      Java, dedicated account, systemd unit, configuration, and restart handlers;
      verify package contents and rendered service/configuration fixtures.
- [ ] 3.3 Configure and guard Nexus/PostgreSQL data mounts and service ordering;
      verify missing/wrong-storage fixtures fail before formatting or service writes.
- [ ] 3.4 Implement bounded readiness and initial credential bootstrap; verify
      first-start, already-complete, interrupted-handoff, and invalid-credential
      cases without password disclosure or repeated credential resets.
- [ ] 3.5 Compose the deployment playbook with host, PostgreSQL, Nexus, and the
      existing Traefik role; verify merged variables, backend routes, listener
      restrictions, and that selecting only Nexus changes no ingress configuration.

## 4. Configure Nexus through Terraform

- [ ] 4.1 Add the service root and pinned provider with explicit URL, CA trust,
      authentication injection, and independent backend; verify offline provider
      schema validation and absence of literal credentials.
- [ ] 4.2 Declare filesystem blob stores and the three specified proxies; verify
      names, formats, upstream URLs, Docker Hub index selection, cache settings,
      and disabled destructive cleanup through configuration fixtures.
- [ ] 4.3 Declare anonymous-access policy, read-only client permissions/account,
      and Docker Bearer Token Realm; verify privilege scope, secret handling, and
      password/state behavior for the pinned provider.
- [ ] 4.4 Configure the private Docker connector and its separate HTTPS ingress
      input; verify client URL examples and rendered routing agree with the selected
      Nexus routing strategy without exposing an unauthenticated backend port.
- [ ] 4.5 Document and encode required adoption of built-in singleton resources;
      verify a fixture containing an unrelated repository is preserved and a repeat
      reconciliation does not manufacture resource drift.

## 5. Validate source and document operations

- [ ] 5.1 Run scoped Terraform, Ansible packaging/rendering, DNS, OpenSpec,
      semantic lint, and repository quality checks; record exact candidate/results
      and inspect representative outputs without invoking live infrastructure.
- [ ] 5.2 Document the ordered Vault/XO/setup/Ansible/service deployment sequence,
      authenticated pip/npm/Docker usage, capacity settings, and failure recovery;
      verify each command maps to an actual packaged target and stage prerequisite.
- [ ] 5.3 Document consistent backup/restore and version-upgrade procedures for
      PostgreSQL, blobs, configuration, and encryption material; verify the procedure
      accounts for VM replacement and irreversible schema migration.
- [ ] 5.4 Deliver a bounded acceptance procedure for an explicitly authorized
      disposable environment covering reboot, unchanged rerun, credential failure,
      unauthorized clients, pinned downloads, and cache reuse with empty client
      caches; verify all specification scenarios map to checks and distinguish
      executed evidence from procedures not yet run.
- [ ] 5.5 After source implementation is explicitly requested and complete,
      validate and deliver the task-owned implementation through repository delivery;
      record any outstanding operational acceptance without claiming deployment.
