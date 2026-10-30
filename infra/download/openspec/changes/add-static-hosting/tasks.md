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

## 8. Deployment and operations documentation

- [x] 8.1 Review the checked-in serving, provisioning, publication, and maintenance paths; distinguish declared behavior from outstanding live acceptance.
- [x] 8.2 Extend the maintained diagrams with deployment order and operations, including issuer reachability, release reselection, deduplication, health checks, and storage recovery limits.
- [x] 8.3 Render and visually inspect all views, verify documentation/build integration, and prepare the documentation for guarded delivery without executing live operations.

This documentation pass follows the 2026-10-25 request to review `infra/download`
and show its deployment and operations. It preserves the incomplete rollout
tasks above and adds no infrastructure behavior or deployment authority.

Review covered the three Terraform roots, stage-scoped AL wiring, Ansible
inventory/templates, the shared `download_host` role, and release-tool
publication contracts at baseline `77f046644`. No new configuration defect
was established by this source review. Recovery limits are intentional;
tasks 1.1, 3.1, and 5.2–5.4 still require their own live evidence. The diagrams
and Ansible package built successfully, and the rendered views were visually
inspected. Exact-candidate validation and publication receipts are retained
under ignored `out/download-diagram/`; generated SVGs remain build outputs.

## 9. Injected SSH access and local inspection

- [x] 9.1 Add explicit local and Yandex SSH targets through `tools/ssh/defs.bzl`, with temporary Vault-signed client credentials and strict repository-CA host verification.
- [x] 9.2 Build both targets and exercise the local target with read-only service and storage checks; retain a sanitized observation and confirm temporary credential cleanup.
- [x] 9.3 Document the targets and reconcile the observed local deployment with the remaining acceptance work; validate and prepare the updated candidate for delivery.

Before implementation, cover these failure cases: the wrong environment is
selected; the injector label or runtime dependency is missing; SSH receives
literal environment placeholders or omits the signed client certificate; an
untrusted host is accepted; inherited agent identities are used; temporary
credential files remain after exit. Use packaged target inspection, the
existing AL lifecycle check, and the authorized local SSH read-only run.
This request authorizes credential injection and local inspection, not a
host deployment, service restart, storage change, or Yandex connection.

The [local SSH observation](evidence.md#local-ssh-inspection-2026-10-25)
establishes current local service, storage, and timer health. It does not
complete the two-environment, reboot, configuration-rerun, or concrete-bundle
comparison criteria retained above.

## 10. Shared-role security fixes

- [x] 10.1 List explicit deletions of inherited SSH and mDNS rules, including multicast destinations, after declaring management access; document rule retirement, storage, and maintenance behavior.
- [x] 10.2 Reference the shared [deployment-hardening change](../../../../../projects/ansible_collection/openspec/changes/harden-download-deployment/proposal.md) for implementation and isolated acceptance. Existing live rollout tasks remain open.

## 11. Deployment diagram layout

- [x] 11.1 Use one Mermaid source with local and Yandex subgraphs at the previous host-level detail: Traefik, loopback Nginx, Btrfs, and the environment's own certificate issuer; omit the overview, operations, DNS, and GitHub Pages and mark cloud deployment deferred.
- [x] 11.2 Render the maintained dark WebP, inspect labels and connections at display size, and validate the diagram and documentation packaging.

The 2026-10-26 revision uses `infra/threexui/docs/architecture.mmd` as its
layout reference. The local view uses Vault HTTP-01 with EAB; the Yandex view
uses Let's Encrypt HTTP-01 without EAB. This changes documentation only and
does not authorize infrastructure operations.

The maintained `//infra/download/docs:deployment_dark_webp` target renders
the same source as the site SVG with the shared dark palette at double pixel
density. `//infra/download/docs:all` built and the owner OpenSpec check passed.
Visual inspection confirmed contained titles, readable labels, and connections
that avoid nodes and title plates. The rendered WebP and source/output digests
are retained under ignored `out/dltf/`; generated images are not checked in.

## 12. Authorized local deployment

The user requested a full local deployment on 2026-10-26 to test the shared
role changes. Scope is the existing XCP-ng VM and the complete local Ansible
playbook, including declared package updates and a reboot if required. Yandex
deployment remains deferred. No resource deletion or replacement is approved.

- [x] 12.1 Review and apply the local Terraform plan, then run the complete local Ansible deployment with the owning packaged targets.
- [x] 12.2 Verify trusted HTTPS, redirects, listings and file responses, service permissions, mount dependencies, firewall rules, and content preservation across a repeated deployment.
- [x] 12.3 Record the exact deployed candidate, observations, failures or remaining limitations, and publish the sanitized evidence.

Before execution, check for destructive Terraform actions, missing storage,
loss of SSH access, failed package or service activation, unintended public
listeners, unreadable or writable service content, certificate trust failures,
and changed release files or selection after a repeated run. Reuse the
existing read-only SSH target and HTTP checks; do not alter stored releases
or service state through imperative diagnostic commands.

## 13. Authorized Yandex deployment

On 2026-10-26 the user reported the quota fixed and requested Yandex deployment.
This supersedes the prior deployment deferral. Scope is the existing Yandex
Terraform root and complete Yandex Ansible playbook; root-site DNS remains on
GitHub Pages. That initial request did not approve resource deletion or
replacement; the later VM-only replacement approval is recorded below.

- [x] 13.1 Apply the reviewed creation-only plan and obtain verified initial host public keys from provider-owned evidence for the first Ansible connection.
- [x] 13.2 Configure public download serving and Let's Encrypt for names routed to this host, then run the complete Yandex playbook.
- [x] 13.3 Verify public HTTPS, storage, service permissions, firewall persistence, and unchanged service reruns; record publication and root-site limitations.
- [x] 13.4 Update deployment status and diagram, retain sanitized evidence, validate, and publish the changes.

Check for exhausted quota, an untrusted initial SSH key, unavailable content
disk, failed host dependencies, unintended listeners, invalid public TLS,
and certificate requests for the apex/www names that still route to Pages.
Use the owning Terraform and Ansible workflows for persistent changes.

Terraform completed with two creates and no deletion; a follow-up plan was
unchanged. The VM is running and cloud-init completed. The first Ansible run
failed before gathering facts because SSH timed out. Direct TCP probes from
the controller and local VM, plus a tunnel through the controller's existing
proxy, did not establish an SSH session. Host configuration and public TLS
checks were blocked at this stage; the later zone change below resolved access.

## 14. Authorized Yandex VM replacement

After the firewall and Terraform investigation, the user explicitly requested
recreating the unreachable VM. Scope is replacement of that VM and its boot
disk through the owning Terraform root. Retain the content disk, reserved IP,
network, security group, and DNS; the local deployment is outside this action.

- [x] 14.1 Review and apply a saved VM-only replacement plan, preserving the content disk and reserved address.
- [x] 14.2 Verify the replacement's boot, disk attachment, initial host keys, and SSH connectivity; retry the owning Ansible deployment and record the result.

The replacement completed with one resource added and one destroyed. The
follow-up plan is unchanged. The fresh VM completed cloud-init and started
SSH, but the post-boot TCP probe and Ansible connection still timed out.
Ansible changed nothing; tasks 13.2 and 13.3 remained incomplete at this stage. See the
[replacement evidence](evidence.md#authorized-vm-replacement).

## 15. Authorized availability-zone change

The user requested moving the deployment to another availability zone and
explicitly approved destroying the unused content disk and reserved IP.
The destination is `ru-central1-a`. Replace the VM, subnet, content disk,
and address through Terraform; update only the component's allocated-address
DNS record. The local deployment and public apex remain unchanged.

- [x] 15.1 Temporarily remove the two destruction guards, review the complete replacement plan, and apply only the authorized zone change.
- [x] 15.2 Restore both destruction guards, verify a no-change plan, and obtain the new VM's initial host keys from authenticated provider output.
- [x] 15.3 Verify hostname resolution and SSH, retry the complete Ansible deployment, and record public-service acceptance or its exact blocker.
- [x] 15.4 Update deployment evidence and current documentation, validate the guarded final configuration, and publish the changes.

Do not hand off or publish source with the destruction guards removed.
Before apply, reject changes beyond the four zone-dependent resources and
the host A record. Before Ansible, verify that the new public address, DNS,
and trusted initial keys all identify the same replacement VM.

## 16. Explicit local-to-public content synchronization

The user requested a command to promote local content to the public store
instead of publishing the same release twice. Implement the command and verify
it in isolation; this request does not execute a production synchronization.
Local is authoritative: public-only content is removed and differing bytes are
replaced, as clarified by the user during implementation.

- [x] 16.1 Expose `sync.yandex` using the existing publication records and SSH injector, with both content mountpoints required; verify the built target's arguments and inputs.
- [x] 16.2 Document local publication followed by explicit sync and dry run, link the release owner's behavior and evidence, and validate the component before publication.

## 18. Authorized website rollout

The user authorized site publication to local, verification, synchronization to
Yandex, public verification, and root DNS cutover in that order. Local version
collisions must fail; public sync overwrites from local. Native automatic-update
packages are approved on both download hosts. After learning that 0.0.0-dev
already exists, the user selected a different release procedure: prepare a
versioned release from the validated feature commit, recording this exception
to the normal trunk-based branch point without overwriting any existing tag.

- [x] 18.1 Authenticate normal publication through the existing Vault SSH injection workflow and verify duplicate local publication refuses without changes.
- [x] 18.2 Deploy and verify the approved automatic-update role on both hosts, retaining package, configuration, timer, and convergence evidence.
- [x] 18.3 Publish the verified versioned website locally, verify HTTP content, sync to Yandex, and verify the public copy.
- [x] 18.4 Enable public site routing and apply the authorized root DNS cutover through its owner; verify trusted apex/www HTTPS and redirects.

## 19. Infrastructure-owned synchronization

The user clarified that release publication belongs in tools/release but
local-to-public synchronization is an infrastructure operation. Move sync to
this component's packaged Ansible playbook and remove the Go sync command and
unused helpers, rather than retaining parallel implementations.

- [x] 19.1 Move and adapt the real SSH acceptance fixture before implementation.
- [x] 19.2 Implement sync with native Ansible modules, controller-held credentials, private staging, checked atomic project exchange, authoritative deletion, and check-mode preview.
- [x] 19.3 Remove the old sync implementation and update targets, documentation, and owner-local contracts.
- [x] 19.4 Pass the adapted fixture and deployment gates, then use this workflow for the authorized rollout.

## 17. Public git redirect

The user requested that public git.alwaldend.com point to this public host and
redirect to github.com, preferably in Traefik. Preserve paths and query strings,
the local Forgejo route, and the public apex's current Pages routing.

- [x] 17.1 Add a Yandex-only Traefik redirect with the existing public certificate resolver and verify both inventory renders and actual redirect responses in isolation.
- [x] 17.2 Coordinate the Forgejo-owned public CNAME update, deploy the scoped Traefik change, and verify public TLS, redirects, and unchanged download behavior.

## 20. Publish the approved deployment diagram

The user approved publication of the final preview after restoring separate
Publisher actors. Reuse the maintained Mermaid source and its existing render
targets; no live infrastructure operation is part of this documentation change.

- [x] 20.1 Publish the approved public-first stacked diagram, separate external services, hostname-labelled HTTPS connections, GitHub redirect, and site/file storage paths.
- [x] 20.2 Package the dark WebP beside the existing SVG variants and inspect the generated outputs. Delivery receipts record final validation and PR publication.

## 21. Preserve public Go module discovery

- [x] 21.1 Reproduce discovery redirection before implementation and record root/nested discovery, prefix-boundary rejection, query matching, browser redirects, and local exclusion failure cases.
- [x] 21.2 Route go-get=1 requests to narrowly scoped static Nginx metadata before Traefik's browser redirect; update both owners' contracts.
- [x] 21.3 Exercise real Traefik/Nginx and the pinned Go client against isolated local Git fixtures, record the receipt, and prepare the fix for repository delivery without a live deployment.
