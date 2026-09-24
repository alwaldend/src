# IaC validation

The implementation delivers infrastructure source, the reusable Nginx role,
Vault AppRole wiring, and the architecture diagram. It does not deploy hosts
or implement the linked release-tool and website-browser changes. The delivery
receipt under ignored `out/infra-download-implementation/` binds the final
candidate and repository gates; fixture evidence below records the tested
service inputs. The acceptance matrix preceded implementation.

## Offline provisioning

All three roots (`local`, `yandex`, `dns`) passed provider-pinned
`init -backend=false`, `validate`, and `test` with mock providers. The mocked
plans exercise the actual resource graphs and canonical DNS inputs without
provider credentials or backend access. Sources were copied to task-owned
scratch with their declared local modules and generated cloud-init input;
`offline --chdir <scratch-root>` selected those copies. Logs and summarized
exit statuses are in `tf-validation.json` and `mock-plans.json` beside the
receipt. These checks establish configuration validity, not available live
inventory or a successful real provider plan.

The roots have distinct Vault state keys. Local owns only its VM and disks;
Yandex owns its VM, network, retained disk/address and delegated DNS zone;
DNS owns the shared aliases/delegation. The apex owner retains its records.
XCP-ng `prevent_destroy` is the user-confirmed protection until a replacement
workflow exists; the cloud disk and address have their own destruction guards.

## Native service fixture

The removed [fixture instructions](https://github.com/alwaldend/src/blob/8a6898c6b3e794124fafa33da11d08aaf759b0a8/infra/download/test/README.md) describe the earlier workflow. On 2026-09-24,
a fresh rootless Fedora 44 container with networking disabled passed both
fixtures against the checked-in Nginx role and consumer templates:

- Nginx `1.30.5-1.fc44`, Ansible Core `2.20.7-1.fc44`, systemd `259.9-1.fc44`.
- Repository-pinned Traefik `3.7.1`.
- Fedora base digest `sha256:f4daff1c8e5c94219817f8b0c39f9c4510231e9b1cf806a3bfe575e92737a815`.
- Fixture image ID `204b578345e4d7985dd5e7e35fe6d3e8b708bb2edaad7a81e714e414c30605a7`.
- Validated Nginx configuration SHA-256 `ab211c207dbbab93331fde3035dac9c3e6c124b28f6e5dc21f7e368d38cae401`.

`result.json`, `traefik-result.json`, and four Ansible lifecycle logs retain
HTTP outcomes and versions. Verified behaviors include JSON with `index.html`
and archives, encoded filenames, exact bytes, HEAD, ranges, CORS, relative
trailing-slash redirects, private-path rejection, symlink-escape rejection,
unknown hosts, rejected writes, site HTML/CSS, atomic selected-link changes,
loopback-only Nginx, and service enablement. An unchanged role rerun performs
no reload; invalid configuration preserves the prior file and running site;
a valid changed configuration reloads without altering content.

Traefik served verified fixture HTTPS for the apex and download host and
redirected `www` permanently while preserving path/query. HTTP redirects
retain path/query and reach HTTPS. Only the certificate issuer was replaced
with an ephemeral fixture certificate. Both issuer configurations were rendered
with HTTP-01, and EAB was present only for the local Vault configuration. Neither
ACME issuer was contacted.

The packaged Ansible playbook passed syntax validation without AL credential
injection. Host disk formatting/mounting, real SSH accounts and SELinux
execution are reviewed declarations, not covered by the native service fixture.
A controller-only Ansible inspection resolved the packaged inventory and both
environment groups. Host FQDNs match the permitted SSH principals and the
canonical DNS-derived names.

## Btrfs deduplication fixture

The user selected Btrfs and daily maintenance after the first PR review.
A fresh rootless container mounted only a task-owned Btrfs fixture directory,
with networking disabled. The same native serving fixtures and the real
`download-deduplicate` systemd service/timer passed there. Artifacts are under
`out/infra-download-implementation/dedup-fixture/`.

Two separately written 4 MiB files in project/site release directories initially
had separate extents. Linux FIEMAP reported 4 MiB of shared extents per file
after maintenance. Byte checksums, inode numbers, and modes were unchanged.
An independent write to one file preserved the other's original checksum.
Three service runs succeeded; HTTP stayed available during maintenance,
private/staging paths and their symlinks were absent from the hash database,
and the daily timer was enabled with a scheduled next run.

`deduplication-result.json` records versions, before/after extents, checksums,
HTTP observations, copy-on-write, and timer state; the service journal and
immutable fixture image ID are retained alongside it. This verifies the
maintenance unit on Btrfs, not live VM formatting or SELinux enforcement.

## Identity review

Source composition uses the existing modules and membership conventions:

| Allowed path or operation                  | Consumer                                                  |
| ------------------------------------------ | --------------------------------------------------------- |
| Own `src_infra_download` KV subtree        | Three state backends and publisher keys                   |
| Own Yandex folder credential subtree       | Cloud VM and delegated DNS provisioning                   |
| Existing Cloudflare DNS-token entry        | Global record provisioning                                |
| Existing RouterOS DNS entry                | dc1 records                                               |
| `ssh/servers/sign/src_infra_download_ssh`  | Host certificates under `download.alwaldend.com`          |
| Existing Ansible admin-client signing role | Configuration connections                                 |
| Existing infrastructure identity group     | XO/Yandex identity discovery and established token policy |

The new policy composition has no read grant for another AppRole's state,
another Yandex folder key, or general provider administrator secrets. The
shared group retains its pre-existing token-creation and identity-scoped PKI
policy. The new `src_infra_download_pki_server` role adds EAB access for internal
HTTP-01 issuance, limited to the three serving names without descendants or
client certificates. Its maximum lifetime is seven days; existing consumers
retain their defaults. This is an offline policy review, not a token-capability
probe. The existing DNS credential's
upstream scope still needs verification during authorized bootstrap.

## Repository checks

Component packaging, Ansible syntax, all three owner OpenSpec validations,
DNS lint/generated-page checks, Terraform formatting in the affected owners,
affected-package semantic lint, and `//:repo_quality_test` passed. The SVG
built through the maintained Mermaid rule and its rendered preview was visually
reviewed. Candidate-bound repetition of mandatory gates is recorded by the
delivery receipt. Generated DNS pages were refreshed with their owning dump
command after general formatting; their currentness check passed.

## Remaining acceptance

Confirm the local address and named XO inventory, provision each selected
host, exercise storage/SELinux enforcement and reboot persistence, test public
ACME staging/production and local Vault issuance/renewal after each issuer can
resolve the relevant DNS view and reach HTTP port 80. The release tool still needs safe
SSH upload/extraction/selection and failure cases; the website still needs
its downloads UI and reusable browser component. No live deployment or cutover
was performed by this implementation pass.

## Bounded workflow review

Observed fixture friction: a container build inherited an unusable loopback
proxy; package acquisition used host networking while the runtime fixture
remained disconnected. Reusing an active Nginx instance observed old workers
during asynchronous reload; the documented procedure now requires a fresh
container. Pinned launchers change directory, so formatter inputs use absolute
paths. These are reproducibility notes, not changes to shared host settings.

## Review refactors and independent PRs

The three Terraform roots now live directly under `infra/download`: `local`,
`yandex`, and `dns`. Offline initialization, validation, and their existing
mock plans passed after relocation. No live deployment or state migration
has occurred. The reusable `infra.xo_login` helper is consumed by download
and OpenHands; both built AL configurations contain exactly one XO plugin
with the expected endpoints and execution labels.

The shared Ansible roles and Vault identity are delivered in independent PRs
against `master`. The download PR consumes them as merge prerequisites.
Service fixture inputs are unchanged by the source relocation; the native
Nginx/Traefik and Btrfs results above retain their original provenance.

The content-device preflight follows symlinks, requires a block device, and
checks the provisioned 100 GiB capacity before configuring the host or creating
the filesystem. Read-only controller probes against a missing path and
`/dev/null` both failed at the attachment assertion before host configuration.
Their inventory, exact command inputs, logs, and results are retained under
`out/infra-download-split/content-device-*`. Live attachment identity and the
successful disk/mount path remain part of authorized rollout verification.

## Terraform extraction

[PR #108](https://github.com/alwaldend/src/pull/108) carries the three Terraform
roots, their existing mock plans, provider assignments, and shared XO helper.
Its candidate `5f37f03b79b6067f5fdd6873746fed98375f2c77` passed all three
backend-disabled initialization/validation/mock-plan sequences, provisioning
AL checks, DNS/OpenSpec checks, formatting, semantic lint, and repository
quality on trunk `636f0baef81f1051f6a796b391b6b85fec80f4d1`.

The host PR removes those roots and provider changes. Its AL configuration
selects only host credentials; the shared canonical DNS inputs remain packaged
for the Ansible inventory. All twelve removed root files are byte-identical
to their extracted counterparts. Service templates, playbooks, and fixtures
are unchanged by this extraction. The earlier service fixture results retain
their provenance. Generated DNS pages now remain at the trunk revision; the
merged manual refresh workflow supersedes the earlier freshness checks.

The subsequent hostname review removes the Ansible address/HostKeyAlias
overrides, DNS-file lookups and packaging, and their derived-name assertion.
All DNS declarations now belong exclusively to #108. That PR gives the local
host its canonical A record and directs the download alias to it, so the named
SSH endpoint can resolve before apex cutover. The storage and serving tasks
and service templates remain unchanged. The earlier split receipt above
precedes this review correction; final receipts bind the revised candidates.

The endpoint-ownership review in #108 subsequently centralizes the local
address in the apex owner and uses a unique host CNAME. Its runbook requires
local DNS preparation before Ansible, preserving public Pages records until
the later public cutover. This replaces the earlier host-A description above;
the hostname-only Ansible implementation and service fixtures are unchanged.

## Handler failure recovery

The play now sets `force_handlers: true`, so notified reload/restart handlers
are not skipped solely because a later task fails. Nginx still validates its
candidate before installation. This uses the existing Ansible play keyword;
packaged syntax and candidate-bound repository checks cover the configuration
change. It does not establish recovery when a host is unreachable.

## Download host role integration

[PR #109](https://github.com/alwaldend/src/pull/109) merged as `43b5cb1`, owning
`alwaldend.main.download_host`. The main host change rebases onto that trunk and
uses the prepared role-based consumer. Its play retains privilege escalation
and `force_handlers`; inventory, Vault injection, and routing templates remain
with the component. Duplicate tasks, defaults, and maintenance templates are
removed. Traefik uses `/opt/traefik` on the system disk.

The reviewed role supersedes the earlier capacity and service-drop-in evidence:
it validates attachment, discovers the Btrfs UUID through the pinned collection,
creates publication roots without privileged traversal into site directories,
and leaves service units to the shared roles. Terraform owns disk sizing.
Role-owned UUID/directory fixtures remain recorded with that change. This
integration checks both packaged inventories and rendered component/maintenance
configuration; no live host or disk operation is performed. Candidate-bound
validation and output records are under `out/infra-download-split/rebase-current/`.

## Final consumer review

The ACME contact uses the component identity, `src_infra_download@alwaldend.com`,
matching the existing deployment convention. The architecture diagram is the
first content below the README metadata. The custom container fixture scripts
and their Bazel/documentation package are removed at the user's request;
proper Molecule coverage is deferred to separate work. Earlier fixture results
above remain historical evidence, not currently maintained test coverage.

## Administrator rsync publication

The current source retains the download account and content ownership but no
longer configures publisher authorized keys or reads the old publisher Vault
entry. Historical credential tables above describe their original candidate.
The release tool and explicit component targets are extracted into an independent
[PR #113](https://github.com/alwaldend/src/pull/113) against `master`. This host
PR retains only the account, access,
storage, service configuration, and server rsync prerequisite. The publication
PR carries the release verification record and isolated SSH-to-HTTP evidence.

The split preserves host implementation from aggregate candidate
`e2d4f59aee637445f996aac6e05c87704b8fa636`. Split validation receipts and
source-preservation evidence are retained under
`out/infra-download-split/publication-extraction/`. No live operation is run.
