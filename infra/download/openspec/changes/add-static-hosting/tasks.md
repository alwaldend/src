## 1. Establish component inputs and acceptance

- [ ] 1.1 Inspect current provider inventory, VM images, DNS ownership, and role packaging; record selected non-secret environment inputs and verify they resolve unambiguously.
- [ ] 1.2 Write the end-to-end scenario/failure matrix before implementation, covering upload, listings, site activation, reselection, disk exhaustion, reboot, and split-horizon TLS; verify every new requirement has an observable check.
- [ ] 1.3 Add the component README and build/OpenSpec packaging through the owning workflows; verify discovery and strict owner validation include the new workspace.
- [ ] 1.4 Complete the linked Nginx-role and Vault-identity plans; verify their owner-scoped checks and record their exact candidate revisions.

## 2. Provisioning, identity, and DNS

- [ ] 2.1 Declare one local VM and one Yandex VM with retained 100 GB content disks and independent state; verify offline Terraform checks and that selected-environment inputs cannot destroy the other VM.
- [ ] 2.2 Add the component AL wiring and provider-side folder/resource-set assignments through their owners; verify packaged plugin labels, identity references, and source permissions without live writes.
- [ ] 2.3 Declare the download hostname in both DNS views and prepare apex updates while preserving its existing `www` alias in the DNS owner; verify the DNS linter reports one owner per name and the expected global/dc1 destinations.
- [ ] 2.4 Give download DNS a single state owner and preserve separate environment VM states; verify declared backend keys and absence of duplicate DNS resource ownership.

## 3. Host configuration and static serving

- [ ] 3.1 Configure mounted storage and publisher/read-only service permissions; verify configuration reruns preserve fixture release files and selected links.
- [ ] 3.2 Configure the standard Traefik service for direct HTTP/HTTPS, loopback Nginx routing, and the `www` permanent redirect; verify intended sockets, anonymous access, and alias path/query preservation without shared ingress.
- [ ] 3.3 Configure public ACME DNS-01 credentials and persistent per-host certificate state covering both content names and `www`; verify rendered configuration and document staged issuance/renewal checks for both DNS views.
- [ ] 3.4 Configure JSON directory routes, file downloads, CORS, and the selected apex website; verify index-file listing behavior, byte ranges, correct MIME types, and non-public staging/service paths in the isolated fixture.

## 4. Publication and website integration

- [ ] 4.1 Complete the linked SSH release-tool change; verify an ordinary file and website archive arrive at the exact public paths without a `files/` layer.
- [ ] 4.2 Complete the linked website change; verify `/downloads/`, header navigation, and embedded release listings against the isolated Nginx endpoint.
- [ ] 4.3 Run the integrated SSH-to-HTTP/browser fixture, including failure preservation and existing-release reselection; retain a repeatable command, candidate identity, checksums, and browser artifacts.
- [ ] 4.4 Run package-scoped formatting/build checks, affected semantic lint, and repository quality gates; verify results against the exact implementation candidate.

## 5. Authorized rollout and evidence

- [ ] 5.1 Obtain explicit live-operation scope for Vault/provider bootstrap, provisioning, Ansible, certificates, publication, and DNS cutover; record the approved environments before running those operations.
- [ ] 5.2 Apply the authorized bootstrap and host stages, then deploy a concrete site release to each selected host; verify HTTPS, artifacts, and `www` TLS/redirect behavior using address overrides before DNS cutover.
- [ ] 5.3 Apply reviewed split-horizon DNS through its owning stages; verify public answers reach Yandex and local answers reach XCP-ng for both hostnames, including obsolete A/AAAA handling.
- [ ] 5.4 Verify authorized reboot persistence and certificate renewal, document the no-backup/manual-retention policy, and retain sanitized operational evidence before marking rollout complete.
