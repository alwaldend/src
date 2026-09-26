## 1. Prepare role fixtures and the runner dependency

- [ ] 1.1 Map the collection requirements to role inputs and assertions before
      role corrections. Verify all three roles cover unchanged runs, configuration
      updates, persistence, and explicit external-integration exclusions.
- [ ] 1.2 Integrate the linked tools/molecule runner after its smoke acceptance.
      Verify each scenario declares packaged role inputs, the supported guest,
      and only its required services without duplicating lifecycle implementation.

## 2. Establish the host scenario

- [ ] 2.1 Add representative users, firewall rules, test trust material,
      locally signed SSH certificates, enabled hardening, and storage inputs.
      Verify only external SSH signing is excluded, the real `host` entry point
      runs, and the configured administrator reconnects with privilege access.
- [ ] 2.2 Verify allowed and denied network traffic, users, trust, timezone,
      hardening, and the virtual-disk mount. Reboot and verify reconnect and
      marker persistence; retain the assertions and guest diagnostics.
- [ ] 2.3 Run the unchanged host convergence check and fix only demonstrated
      in-scope blockers. Verify zero changes within the declared fixture
      coverage, unchanged deployment defaults, and successful final cleanup.

## 3. Establish the Traefik scenario

- [ ] 3.1 Add a local backend, generated test TLS material, and static/dynamic
      configuration with EAB disabled. Verify real service enablement, trusted
      HTTPS access, and rejection of an unmatched route on a fresh VM.
- [ ] 3.2 Exercise unchanged and changed configurations before correcting
      lifecycle tasks. Make only required handler/reload fixes; verify zero
      changes and no unnecessary reload on an unchanged run, new routing after
      an update, removal of the old route, and subsequent idempotence.

## 4. Establish the Forgejo scenario

- [ ] 4.1 Add SQLite configuration and disposable account/repository setup
      separate from role convergence. Verify real service enablement and HTTP
      access, then push and freshly clone a known commit and compare its content.
- [ ] 4.2 Exercise unchanged and changed application configuration before
      correcting lifecycle tasks. Make only required fixes; verify zero changes
      and no unnecessary reload on an unchanged run, an observable configuration
      update, and subsequent idempotence.
- [ ] 4.3 Restart Forgejo and verify the existing account, repository, and
      commit remain accessible without recreating them. Retain sanitized
      application results and verify VM cleanup.

## 5. Integrate, document, and deliver

- [ ] 5.1 Expose the three role test aliases and aggregate suite. Query the
      selected labels, run each role independently and the complete suite, and
      inspect runtime package mappings to verify fixtures are excluded from
      production payloads.
- [ ] 5.2 Document supported tools/guest, exact commands, network and cache
      limitations, fixture exclusions, artifacts, and cleanup recovery. Verify a
      fresh invocation can reproduce the documented workflow without production
      credentials or persistent host changes.
- [ ] 5.3 Run the complete behavioral acceptance against the exact candidate,
      inspect representative retained artifacts, and run affected semantic lint,
      collection packaging, OpenSpec validation, and repository quality gates.
      Record actual results and unresolved environmental failures separately.
- [ ] 5.4 Review the final scope, commit and publish through repo-delivery,
      and archive this change only after implementation acceptance is complete.
      Verify publication receipts, retained evidence links, and cross-change links
      after either change is archived.
