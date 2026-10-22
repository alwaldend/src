## 1. Establish component inputs and acceptance

- [ ] 1.1 Inspect current provider inventory, VM images, DNS ownership, and role packaging; record selected non-secret environment inputs and verify they resolve unambiguously.
- [x] 1.2 Write the end-to-end scenario/failure matrix before implementation, covering upload, listings, site activation, reselection, disk exhaustion, reboot, and split-horizon TLS; verify every new requirement has an observable check.
- [x] 1.3 Add the component README and build/OpenSpec packaging through the owning workflows; verify discovery and strict owner validation include the new workspace.
- [x] 1.4 Complete the linked Nginx-role and Vault-identity plans; verify their owner-scoped checks and record their exact candidate revisions.

## 2. Provisioning, identity, and DNS

- [x] 2.1 Declare one protected local VM and one Yandex VM with 100 GB content disks and independent state; retain the Yandex disk across replacement and defer XCP-ng replacement; verify offline Terraform checks and that selected-environment inputs cannot destroy the other VM.
- [x] 2.2 Add the component AL wiring and provider-side folder/resource-set assignments through their owners; verify packaged plugin labels, identity references, and source permissions without live writes.
- [x] 2.3 Declare the download hostname in both DNS views and prepare apex updates while preserving its existing `www` alias in the DNS owner; verify the DNS linter reports one owner per name and the expected global/dc1 destinations.
- [x] 2.4 Give download DNS a single state owner and preserve separate environment VM states; verify declared backend keys and absence of duplicate DNS resource ownership.

## 3. Host configuration and static serving

- [ ] 3.1 Configure mounted storage and publisher/read-only service permissions; verify configuration reruns preserve fixture release files and selected links.
- [x] 3.2 Configure the standard Traefik service for direct HTTP/HTTPS, loopback Nginx routing, and the `www` permanent redirect; verify intended sockets, anonymous access, and alias path/query preservation without shared ingress.
- [x] 3.3 Configure HTTP-01 with Let's Encrypt externally, Vault ACME/EAB internally, and persistent per-host certificate state covering both content names and `www`; verify rendered configuration and document staged issuance/renewal checks for both DNS views.
- [x] 3.4 Configure JSON directory routes, file downloads, CORS, and the selected apex website; verify index-file listing behavior, byte ranges, correct MIME types, and non-public staging/service paths in the isolated fixture.

## 4. Publication and website integration

- [x] 4.1 Complete the linked SSH release-tool change; verify an ordinary file and website archive arrive at the exact public paths without a `files/` layer.
- [ ] 4.2 Complete the linked website change; verify `/downloads/`, header navigation, and embedded release listings against the isolated Nginx endpoint.
- [ ] 4.3 Run the integrated SSH-to-HTTP/browser fixture, including failure preservation and existing-release reselection; retain a repeatable command, candidate identity, checksums, and browser artifacts.
- [x] 4.4 Run package-scoped formatting/build checks, affected semantic lint, and repository quality gates; verify results against the exact implementation candidate.

## 5. Authorized rollout and evidence

- [ ] 5.1 Obtain explicit live-operation scope for Vault/provider bootstrap, provisioning, Ansible, certificates, publication, and DNS cutover; record the approved environments before running those operations.
- [ ] 5.2 Apply the authorized bootstrap and host stages, then deploy a concrete site release to each selected host; verify HTTP content using address overrides, then HTTPS and `www` redirects once each issuer can validate the routed DNS names.
- [ ] 5.3 Apply reviewed split-horizon DNS through its owning stages; verify public answers reach Yandex and local answers reach XCP-ng for both hostnames, including obsolete A/AAAA handling.
- [ ] 5.4 Verify authorized reboot persistence and certificate renewal, document the no-backup/manual-retention policy, and retain sanitized operational evidence before marking rollout complete.

## Recorded implementation scope

This implementation pass delivers IaC and the architecture diagram. The user
confirmed protecting the XCP-ng VM and handling replacement in a later change.
See [evidence.md](evidence.md) for completed source and fixture checks.

The initial source pass left live inventory/disk/host verification in tasks 1.1
and 3.1. Later progress is reconciled below; neither these historical records
nor this task-state review authorizes a new infrastructure operation.

### Progress reconciliation (2026-10-22)

Task 1.4 is complete. [PR #104](https://github.com/alwaldend/src/pull/104)
merged the Nginx role at `14e51b750cd342a2fcd3178ebd3b9dd4c6fd05b7`, with
recorded lifecycle/HTTP fixtures. [PR #105](https://github.com/alwaldend/src/pull/105)
merged the Vault identity at `5ee276f1bdaa6435fb6b905664f545ee98faa4ee`;
[PR #107](https://github.com/alwaldend/src/pull/107), candidate
`1758ada2b212bd45273a440550fd8cc6d5746cd2`, records its authorized bootstrap,
capability checks and no-change follow-up plans.

The remaining checkboxes retain these specific gaps:

- 1.1 and 3.1: both-environment inventory resolution and mounted-host/storage
  acceptance remain unproven by the retained source and fixture evidence.
- 4.2 and 4.3: the Downloads UI and browser acceptance are absent. The SSH/HTTP,
  reselection and failure portions already passed 44 isolated scenarios in
  [PR #113](https://github.com/alwaldend/src/pull/113), candidate
  `b216b45686ac2aef98feef77b118a39c5f895767`; they need not be rediscovered.
- 5.1 and 5.2: retain the recorded Vault authority and the later explicitly
  authorized XCP-ng publication scope in the website design. Vault bootstrap
  and local-deployment history exist, but complete two-environment concrete-site
  HTTP/HTTPS/redirect acceptance and public/Yandex authority are not established.
- 5.3: [PR #114](https://github.com/alwaldend/src/pull/114), candidate
  `a60c613ca17f5f72ebd8651d3b9893d9770312b4`, records the observed local AAAA
  fallback and a dc1-only correction plan during an XCP-ng deployment. It does
  not establish successful public/Yandex cutover or full two-view acceptance.
- 5.4: the no-backup/manual-retention policy is documented; reboot persistence
  and certificate-renewal acceptance remain unevidenced.

These are historical source and publication observations, not a fresh claim
about live service health. The initial evidence file retains its original scope.

## 6. Content deduplication

- [x] 6.1 Configure Btrfs and daily incremental deduplication over published content with private hash state, mount ordering, and bounded resource usage.
- [x] 6.2 Exercise real Btrfs extent sharing, unchanged bytes/inodes/permissions, copy-on-write isolation, repeat execution, and serving during maintenance; retain repeatable fixture artifacts.

## 7. PR review and separation

- [x] 7.1 Move the independent Terraform roots to `infra/download/local`, `yandex`, and `dns`, keeping independent state and updating source/build/documentation references.
- [x] 7.2 Extract `infra.xo_login` and use it for download and the existing OpenHands consumer.
- [x] 7.3 Move Vault identity composition to the standalone `infra/vault/approles/src_infra_download` root with its own backend and name-based lookups.
- [x] 7.4 Separate the shared Ansible roles and Vault identity into independent PRs against `master`, with the component PR retaining only its own scope.
- [x] 7.5 Guard filesystem creation with a block-device check; verify missing and non-block paths fail before host configuration. The reviewed role in #109 removes the capacity assertion; Terraform owns sizing.

The shared Ansible roles are delivered in [PR #104](https://github.com/alwaldend/src/pull/104)
and the Vault identity in [PR #105](https://github.com/alwaldend/src/pull/105).
Both target `master`; merge prerequisites do not form a PR stack.

- [x] 7.6 Extract the Terraform roots, provider assignments, and XO helper into [PR #108](https://github.com/alwaldend/src/pull/108), targeting `master`; retain only host configuration and its canonical inputs in #102.

The merged Terraform PR owns all DNS declarations. The host PR connects by
inventory FQDN without parsing DNS files, address overrides, or HostKeyAlias.
Its rebase onto #108 combines the stage-scoped provisioning and host AL calls.
The [Terraform evidence](terraform-evidence.md) retains the prerequisite
validation and review history. Generated DNS snapshots remain at the trunk
revision under the manual refresh workflow.

- [x] 7.7 Address hostname review comments by using inventory FQDNs directly and preparing local DNS before configuration while keeping public cutover separate.

- [x] 7.8 Force notified service handlers after later task failures so a retry cannot silently leave successfully installed configuration unapplied.

- [x] 7.9 Rebase the host configuration onto merged Terraform PR #108, preserving its provider parameter, folder permissions, DNS input packaging, and deployment prerequisites.

- [x] 7.10 Consume the merged `download_host` role from #109, removing duplicated
      tasks/defaults/templates and keeping Traefik state on the system disk;
      validate both packaged inventories and rendered configurations.

- [x] 7.11 Apply the final consumer review: use the AppRole-based ACME contact,
      remove custom fixture scripts for future Molecule coverage, and put the
      architecture diagram first in the README.

- [x] 7.12 Retain the download account without configuring authorized keys,
      implement rsync deployment in the release tool, and expose explicit
      component publishing targets; verify SSH-to-HTTP behavior in isolation.

- [x] 7.13 Extract the release tool and component publishing targets into an
      independent [PR #113](https://github.com/alwaldend/src/pull/113) against
      `master`, preserving the host account
      and rsync prerequisites here; validate and retain evidence for both branches.
