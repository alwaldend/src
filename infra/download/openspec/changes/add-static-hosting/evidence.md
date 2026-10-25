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

## Local SSH inspection (2026-10-25)

At 18:58:19 UTC, the new `//infra/download:ssh.local` target connected to
`host1.dc1.download.alwaldend.com` as `ansible`, using the component AppRole,
an injected Vault-signed client key, and strict host verification against the
repository SSH server CA. Both environment targets built; only local ran.
The user explicitly authorized injected SSH access and this inspection.

Read-only `systemctl show`, `findmnt`, `df`, `btrfs filesystem usage`,
`readlink`, `stat`, `ss`, `getenforce`, and `uptime` observations established:

- Nginx and Traefik were active/running, with successful main-process status.
- `/dev/xvdb` was mounted at `/srv/download` as 100 GiB Btrfs with
  `nodev,nosuid,noexec`; approximately 1.3 GiB was used and 97 GiB available.
- The deduplication timer was active/waiting. Its last run at 00:20 UTC on
  October 25 exited successfully; the next run was scheduled for October 26.
- `sites/alwaldend.com/current` selected `releases/0.0.0-dev`. Publication
  roots belonged to `download`; staging remained mode `0700`.
- Nginx listened on `127.0.0.1:8008`; ports 80/443 were listening. SELinux
  was enforcing. The host had been up since September 26.
- No injected SSH key directories or temporary known-hosts files remained in
  the invocation's task-private temporary directory after successful exit.

The repeatable remote command input is `out/download-ssh/inspect-local.sh`,
executed with `bazel_agent bazel run //infra/download:ssh.local -- bash -s`.
Its SHA-256 is
`135e385c41ecf751409ae027f1085c0d1f7fc9edfe8f3adf44de22944e73edba`.
The sanitized `out/download-ssh/local-observation.log` has SHA-256
`1078f85db50a57e99cb2390862e74e45a98806dae008dad6cba2a32832630c21`.
`TMPDIR`, `TMP`, and `TEMP` pointed to ignored task-private scratch. These
observations establish local health at that time, not Yandex acceptance,
reboot persistence, an Ansible rerun, or a comparison against a newly built
website bundle. No service, storage, or persistent host configuration changed.

## Authorized local deployment 2026-10-26

The user authorized a full local deployment to test the hardening changes.
Only `host1.dc1.download.alwaldend.com` was deployed. Local Terraform's saved
plan and apply completed with zero additions, changes, or deletions. No DNS,
Yandex, publication, release-selection, or resource-replacement operation ran.
Source baseline was `1d2100e8ea8bb69678fa5ec96088cfcda202d72e`; the three
runtime corrections below belong to this delivery. Private receipts and
`out/ld/deployed-inputs.json` bind their input digests to the delivered candidate.

The live run exposed three defects missed by the earlier bounded checks:

- `download-ansible-host-trust`: Ansible lacked the host-CA entry and failed
  before gathering facts. Its wrapper now reuses the operator target's
  temporary, namespace-scoped CA trust with strict verification enabled.
- `traefik-first-enrollment-cache`: the enrollment block reevaluated its
  condition after replacing the cached-value fact, skipping the cache write.
  The role now records the decision before changing that fact. The missing
  cache was created with root ownership and mode `0600`; subsequent runs
  skipped the Vault enrollment request.
- `ufw-boot-activation`: `ENABLED=yes` did not enable `ufw.service`. The first
  reboot left filtering inactive, and the next logging task failed against
  missing chains. The role now enables/starts the native service before rule
  updates, after stopping conflicting firewalld. A second reboot demonstrated
  that filtering returned without another firewall task.

The first full attempt installed package updates and rebooted on DNF's hint,
then failed at the missing enrollment cache. After the corrections, the full
playbook completed with `ok=237 changed=15 failed=0 unreachable=0`. It included
a deliberate second reboot by overriding the existing reboot-result variable:

```sh
bazel_agent bazel run //infra/download/ansible:ansible.local -- \
  --extra-vars '{"update_all_packages_reboot_required":true}'
bazel_agent bazel run //infra/download/ansible:ansible.local -- \
  --tags firewall,nginx,traefik
```

The override exercised the owning Ansible reboot task; it is not a persistent
setting or a claim that DNF requested the second reboot. The subsequent service
rerun reported `ok=39 changed=0 failed=0 unreachable=0`. Traefik and Nginx
retained their process IDs and start timestamps across that rerun.

Final read-only observation: **2026-10-26T09:53:54Z**.

- UFW was active and its systemd service enabled. Default incoming policy was
  deny; the retained IPv4/IPv6 allows were TCP 22, 80, and 443. Retired SSH/mDNS
  profiles were absent. Controller TCP probes reached 22/80/443 and timed out
  for 5355 and 8008; Nginx was bound only to `127.0.0.1:8008`.
- Nginx, Traefik, SSH, and the deduplication timer were active. Nginx and the
  maintenance service required `/srv/download`; Traefik did not. Btrfs was
  mounted by UUID with `nodev,nosuid,noexec,nofail`; SELinux was enforcing.
- Traefik's parent/configuration directories were root-owned `0750`, its
  configuration files `0640`, binary `0755`, and enrollment cache `0600`.
  The service identity could neither write configuration/binary nor read the
  private cache. Runtime data remained service-owned `0700`, ACME state `0600`.
  Nginx could read the selected index but could not write the project root.
- Sixteen HTTP checks passed before and after deployment with identical
  response bytes and selected headers: repository-CA-verified HTTPS for all
  three serving names, HTTP and `www` redirects preserving path/query, JSON
  listings and CORS, the complete 93,464,617-byte release archive, HEAD and
  matching byte ranges, and private-path rejection. Existing Vault certificates
  remained trusted; fresh TLS issuance and timed renewal were not exercised.
- `current` still selected `releases/0.0.0-dev`. File counts and total sizes
  stayed unchanged for the archive and 6,374 site files; twelve sampled site
  checksums and the complete archive's HTTP checksum matched their baseline.

Reproduction inputs and sanitized observations are in ignored `out/ld/`:
`inspect.py`, `http_check.py`, `summary.json`, `network-probes.json`, and the
before/after observations. Inspection uses `ssh.local -- sudo -n python3 -`
with the script on stdin. HTTP checks pin `192.168.10.66` while verifying each
requested hostname against the repository server CA. Raw Terraform plan files
were removed after applying and retaining their non-secret summary.

Runtime correction SHA-256 values:

| Source                                                          | SHA-256                                                            |
| --------------------------------------------------------------- | ------------------------------------------------------------------ |
| `infra/download/al.lua`                                         | `319e5fa6125cec9da3064f65b3d2adf7d162a9b48a429fa4c67d872ee7e9a6bb` |
| `roles/traefik/tasks/enrollment.yaml` in the Ansible collection | `a20dd4e84918991d2e9eb9093c293f84c33bab2b6b282ec537558515169eac52` |
| `roles/firewall/tasks/main.yml` in the Ansible collection       | `657bfa704f5b2aa6569c66b5917d50dc7802d3d804dea1de58bc3115bc25b2d3` |

This verifies the existing local release and host lifecycle, not publication
of a new bundle, physical disk-loss recovery, all packet types/interfaces, or
the deferred Yandex deployment. Those broader acceptance tasks remain open.
The initial fixture had covered cache reuse and UFW rule files; the new live
observations specifically close first-cache and boot-activation gaps without
relabeling the earlier fixture as broader coverage.

## Yandex provisioning 2026-10-26

The user reported the HDD quota fixed and authorized Yandex deployment.
Terraform source at `f24934cd48ad028ec67ac97c8d83fc93994105c2` successfully
applied the reviewed saved plan: **2 added, 0 changed, 0 destroyed**. The
100 GiB content disk and VM with its 20 GiB boot disk were created. The
existing network, subnet, security group, reserved address, and host DNS
remained unchanged. A fresh full plan subsequently reported no changes.

Observed provisioning facts:

- Instance `fv42ifm6l98jtq2g2lh3` is `RUNNING`, using Fedora 43 image
  `fd8iku26nnkveh8s4di1` in `ru-central1-d`.
- Content disk: `fv4d6phbfap3jp554cfr`; private address: `10.9.0.4`;
  reserved public address: `158.160.255.130`.
- Public DNS for both `host1.yc.download.alwaldend.com` and
  `download.alwaldend.com` resolves to that reserved address.
- The attached active security group allows incoming TCP 22/80/443 from
  `0.0.0.0/0` and all outgoing traffic. The subnet has no custom route table.
- Authenticated provider serial output records a default route through
  `10.9.0.1`, a started `sshd.service`, and cloud-init completion at
  **2026-10-26T10:21:26Z**. Initial SSH host keys were extracted from its
  host-key section into a private temporary known-hosts file. The Ed25519
  fingerprint is `SHA256:klrtd53HuGjypuLRjsSdtzydibPmAhw5YDaQ0d7KWYw`.

Host configuration is **not complete**. The initial `ansible.yandex` invocation
used the verified host-key file with strict checking enabled, but timed out
before gathering facts: `ok=0 changed=0 unreachable=1 failed=0`. Controller
TCP probes to 22, 80, and 443 timed out; a TCP probe from the local download VM
also timed out. The controller's existing HTTP proxy accepted CONNECT but
delivered no SSH banner within the bounded timeout. A later authenticated
provider read still showed the VM running and the same successful boot log.
The packet-loss boundary remains undetermined; a running guest and configured
allow rules do not establish working public reachability.

The source now gates apex/www routes with `download_site_enabled`: true locally,
false on Yandex while those public names remain on GitHub Pages. This prevents
premature Let's Encrypt requests for names whose challenges cannot reach the
new VM. The packaged Ansible syntax check and controller-side rendering of
both environments passed. Render checks verified the download router in both,
apex/www only locally, loopback Nginx, and valid dynamic TOML. This source has
not reached the Yandex guest. No content was published, filesystem formatted,
service configured, certificate issued, or root-site DNS changed by Ansible.

Private task artifacts are under `out/yd/`: plan/apply summaries, sanitized
provider observations, the initial host-key projection, failed Ansible log,
and rendered configurations. The provider inspection uses the component's
existing AL-injected service account and the
[serial-output API](https://yandex.cloud/en/docs/compute/api-ref/Instance/getSerialPortOutput);
credential values and raw serial output were not copied into evidence.
Raw saved-plan/state JSON is removed after retaining its non-secret summary.

Resume by checking public SSH reachability, reconfirming the instance/address
and pinned host keys, then running the complete owning Ansible target. After
host certificates activate, return to its normal CA-trusted invocation and
verify public HTTPS, Btrfs, UFW startup, permissions, and an unchanged service
rerun. The public root remains on Pages until a separately approved cutover.

### Authorized VM replacement

The user subsequently requested recreating the unreachable Yandex VM. Using
source at `398a04ceb6b90545058294aef34af0b2ea09e819`, the owning `tf.plan`
target saved a plan with `-replace=yandex_compute_instance.download`. Its
JSON review showed exactly one delete/create pair: the VM and its automatically
deleted boot disk. All seven other Terraform resources were no-ops, including
the independently retained content disk and reserved address. The owning
`tf.apply` target applied that exact saved plan: **1 added, 0 changed,
1 destroyed**. A fresh full plan reported no changes.

Provider observations at **2026-10-26T13:27:42Z** confirmed:

- New running VM: `fv48uhoida702mi6vc8g`, replacing
  `fv42ifm6l98jtq2g2lh3`; the Fedora image and zone are unchanged.
- New boot disk: `fv4ordjvms4248khl67b`; retained content disk:
  `fv4d6phbfap3jp554cfr`, attached with `autoDelete: false`.
- Reserved public IP remains `158.160.255.130`, with provider `used: true`;
  the new private address is `10.9.0.27`.
- Cloud-init finished at **2026-10-26T13:27:03Z**, the default route uses
  `10.9.0.1`, and the serial log records `sshd.service` started.
- The new Ed25519 host fingerprint is
  `SHA256:AzmZobd4S1XPxngMqvjoyi+KmfYKw7ZgPOOZGn1SvqY`. Initial host public
  keys were obtained from authenticated serial output for this exact new VM;
  the old VM's bootstrap trust file was not reused.

Recreation **did not resolve SSH connectivity**. A direct TCP 22 probe after
cloud-init completion timed out after ten seconds. The complete owning
`ansible.yandex` target was retried with the new verified known-hosts file and
strict checking. It timed out before gathering facts:
`ok=0 changed=0 unreachable=1 failed=0`. No host configuration, content
publication, or certificate issuance occurred. The guest firewall remains
uninspected, and the packet-loss boundary is still undetermined. The attached
security-group rules, reserved-address mapping, content disk, and DNS were
preserved; no local resource or root-site DNS was changed.

Private artifacts use the `out/yd/replacement-*` prefix: reviewed plan summary,
apply and no-change plan logs, sanitized provider observations, and the failed
Ansible log. `replacement_known_hosts` contains only the new public host keys.
Raw saved-plan material is removed after its non-secret summary is retained.
At this stage, the deployment diagram still marked configuration as pending.

## Yandex zone change 2026-10-26

The user explicitly requested another availability zone and then approved
destroying the unused content disk and reserved IP. The default zone changed
from `ru-central1-d` to `ru-central1-a`. Both destruction guards were removed
only for the approved replacement, then restored immediately after apply.

Before the move, the live security group still matched Terraform: incoming
TCP 22/80/443, unrestricted IPv4 egress, correct attachment, and no custom
subnet route table. Controller DNS, public DNS, and both authoritative Yandex
nameservers agreed on the old IP; there was no AAAA record or hosts-file
override. A TCP MTR trace received replies through `10.255.255.57`, then no
further replies; this does not identify the packet-loss boundary. A probe from
the existing `threexui` VM, whose metadata confirmed `ru-central1-d`, connected
to the old download IP in 32 ms and received `SSH-2.0-OpenSSH_10.0`. SSH worked
from inside Yandex while controller connections timed out.

Using baseline `140281a04f75e8555d373f3969ff705ed2fe7a8d` plus the zone and
temporary lifecycle edits, the owning Terraform targets saved and reviewed
the complete plan before applying it: **4 added, 1 changed, 4 destroyed**.
The VM, subnet, content disk, and reserved address were replaced; the host A
record changed. The network, security group, and DNS zone remained unchanged.
After restoring both guards, a full plan reported no changes.

Provider observations at **2026-10-26T14:14:57Z** confirmed:

- Running VM `fhm0hcpj1uvtal4d21c4` in `ru-central1-a`.
- Boot disk `fhme4es3qrukh0rvlpv0` and 100 GiB content disk
  `fhmlkb2eal2efsd06hdm`, with content `autoDelete: false`.
- Reserved address `e9bbdr8mgqjs527gt51f`: **`51.250.85.209`**, in use.
- Subnet `e9bhke8mmu6nmqod6hps`, still `10.9.0.0/24`.
- Cloud-init completed at **2026-10-26T14:14:40Z**. The new Ed25519 host
  fingerprint, obtained from authenticated serial output, is
  `SHA256:K7LTl5Fe6kIGut8v3dA6gu0IfDji8C19MczYFKmypEU`.

The controller connected directly to TCP 22 in 8.664 ms. Both controller and
public DNS resolved the inventory hostname to the new address. The public
download name also resolves there; apex A/AAAA records remain the eight GitHub
Pages addresses, and `www` remains a CNAME to the apex. Both zone and public
address changed, so this outcome does not establish which caused the previous
network-path failure.

The complete `ansible.yandex` invocation used the verified initial host keys
with strict checking and completed all host tasks:
`ok=244 changed=97 unreachable=0 failed=0`. Package updates required and
performed a reboot. The wrapper subsequently exited 1 during cleanup because
its injector could not revoke an owned Vault token before its shutdown timeout.
The configured AppRole token lifetime is 600 seconds and the deployment took
longer; expiry is a plausible explanation, but the sanitized error does not
establish the underlying response. No temporary credential files remained.
Subsequent SSH and shorter Ansible invocations exited successfully.

Post-deployment inspection exposed inactive UFW after that reboot. Firewalld
was absent at the initial firewall check and installed by a later dependency;
both conflicting units were then enabled. The shared firewall role now masks
firewalld even when initially absent, after stopping/disabling it when present.
The focused firewall/Nginx/Traefik run passed with `ok=33 changed=3 failed=0`.
A verification reboot through `ansible.yandex --tags update_all_packages
--extra-vars '{"update_all_packages_reboot_required":true}'` passed with
`ok=5 changed=1 failed=0`. Before any subsequent reconciliation, read-only
inspection confirmed a new boot ID, active enabled UFW, Nginx, Traefik, SSH,
and the deduplication timer. Firewalld remained masked and inactive. The final
firewall/Nginx/Traefik repeat passed with `ok=33 changed=0 failed=0`; Nginx and
Traefik PIDs remained unchanged.

Public acceptance passed before and after reboot and after the unchanged run:

- Ten HTTP checks verified system-trusted Let's Encrypt TLS for
  `download.alwaldend.com`, HTTP-to-HTTPS redirects, the projects redirect,
  JSON listings, HEAD, CORS, private-path rejection, and absent apex/www routes
  on the Yandex proxy. The certificate issuer is Let's Encrypt `YR1`, expiring
  **2027-01-24T13:30:20Z**.
- The 100 GiB content disk is mounted as Btrfs at `/srv/download` with
  `nosuid,nodev,noexec`; Nginx requires that mount and listens only on
  `127.0.0.1:8008`. SELinux is enforcing.
- UFW allows only TCP 22/80/443 for both address families and denies other
  incoming traffic. External TCP probes to 8008 and 5355 timed out.
- Nginx can read but cannot write the projects directory. Traefik cannot
  modify its root-owned configuration or binary; ACME data is private (0600).
- No release was uploaded or selected. The public projects listing is empty,
  so file-byte and range checks on Yandex await content publication. The apex
  stays on GitHub Pages. Certificate renewal and disk-loss recovery remain
  unverified.

Private `out/yd/zone-a-*` artifacts retain the reviewed plan summary and digest,
apply/no-change logs, provider observations, deployment output, host reports,
and HTTP results. `http_check.py` repeats the ten address-pinned HTTP checks;
`inspect_host.py` runs read-only through `ssh.yandex`. Raw saved-plan material
was removed after summarization. The deployment diagram no longer marks the
Yandex host as awaiting configuration.

## Local-to-public synchronization command

`sync.yandex` builds and uses the existing publication records and SSH injector.
The generated launcher was inspected for local-to-Yandex direction, required
mountpoints, temporary client certificate and host trust references, and the
existing administrator/content-account boundary. The release owner's
[isolated acceptance](../../../../../tools/release/openspec/changes/sync-published-content/evidence.md)
records source-authoritative replacement and deletion tests. This task adds
the command without running it against either deployment. The empty public
store and public apex routing have not been changed by this implementation.

## Public Git redirect and Traefik update

The user requested a public git.alwaldend.com redirect and a Traefik update.
The Yandex-only route uses Traefik's permanent RedirectRegex middleware and
existing Let's Encrypt resolver. The public Forgejo-owned CNAME was updated
in place; local Forgejo DNS remains 192.168.10.40. The root site remains on
GitHub Pages.

The scoped public Traefik deployment completed 20 tasks with four changes and
zero failures, and its wrapper exited successfully. The running binary is
3.7.16; the service is active and enabled. Twelve isolated routing checks and
fourteen live HTTP/HTTPS checks passed, including exact path/query preservation
and trusted Git/download certificates. Content remains unchanged and no
production synchronization was run.

See the [DNS acceptance](../../../../forgejo/openspec/changes/redirect-public-git/evidence.md)
and [dependency update PR #203](https://github.com/alwaldend/src/pull/203)
for ownership, checksums, certificate dates, and private evidence locations.

## Ansible synchronization migration (2026-10-26)

`//infra/download/test/sync:sync_test` now executes the packaged Ansible
playbook and release CLI through a loopback OpenSSH server with an ephemeral
client certificate. All 17 receipt cases pass. The fixture covers mounted-root
requirements, host trust, unsafe layouts, missing indexes, authoritative
replacement and deletion, unchanged repeats, failed staging without public
changes, local collisions and concurrency, retained archive hashes, extraction
cleanup and rollback, empty source, and missing source. Its declared output
`sync-receipt.json` records content hashes, modes, links, and completed cases.

A new payload-filter case first failed because synchronize's implicit `-F`
interpreted a published `.rsync-filter` as transfer policy and omitted a valid
payload. Explicit filter reset preserves these files as ordinary content. The
complete fixture then passed. Paths containing spaces use the supported, quoted
`--link-dest` rsync option because the pinned module's `link_dest` parameter
does not quote them. Both directions use actual SSH even for loopback hosts.

The Go sync command, implementation, and old build/test target are removed.
Release publication remains in tools/release; infra/download owns the mirror.
Read-only live preflight confirms both Btrfs mounts and sufficient space,
local selection `releases/0.0.0-dev`, an empty public store, and GNU coreutils
9.7 on Yandex. No production synchronization has run at this checkpoint.

## Versioned site publication and mirror (2026-10-26)

The approved feature-branch release exception produced local tag v2026.44.0
and releases/2026.44 at published PR commit
`b8ecc2f8690098e0441c42f05bf29b1072f62f61`. The exact release archive is
34,575,710 bytes compressed and 271,246,267 bytes extracted: 6,524 files,
3,503 HTML files, and all 127 tag pages. Archive inspection found no project
binary payloads. Its SHA-256 is
`b6dddb3e672eba687f10c0f5fe55470b56929b339af598ac02cb6723d0069758`.

Remote release-ref publication was refused. A read-only Git dry run accepts
the proposed atomic ref update, while GitHub reports active creation/update
restrictions and required signatures for releases/2026.44. The versioning
wrapper masks all push errors containing its own --atomic argument as
unsupported atomicity. Neither remote release ref was created; no protection
was changed. The authorized site rollout uses the fixed locally tagged version
whose source commit is already published in the PR. Raw diagnostics remain in
ignored task evidence. The standard remote release procedure remains incomplete.

Local upload and extraction succeeded. Selection succeeded but initial cleanup
reported an inaccessible inherited login directory. The release owner's real
SSH fixture reproduced that boundary; the corrected tool then completed a
link-only retry, leaving only sites/alwaldend.com/releases/2026.44.0. Five
local HTTPS pages and four assets matched the archive exactly. The old archive
remains unchanged with SHA-256
`1d7b12cd719c0449cbde042dcf8df6b035a69b2747dd94e705b0e38d7debce63`.

The Ansible check-mode preview proposed only additions to the empty public
store. Live synchronization completed with ok=47, changed=10, failed=0. Yandex
has the selected version, 6,524 site files, and both matching archive hashes.
A direct, proxy-free HTTPS range request to the public host returned 206 with
the expected 34,575,710-byte total and matching first 4,096 bytes. Site ingress
and root DNS acceptance follow in the DNS owner's activation record.

The restrictive-umask fixture failed before the archive-directory correction
and passed afterward. A final link-only retry restored the selected local
archive directory to 0755; trusted direct local HTTPS returned the same
archive range as Yandex. The public repeat sync preview was unchanged for
both content trees. Final root routing, DNS, and certificates are verified in
[the DNS activation record](../../../../dns/openspec/changes/activate-public-download-site/evidence.md).

The two live-only fixture gaps were private login-directory traversal and
remote umask inheritance for rsync --files-from. Their stable evidence names
are private-home-red/green and archive-mode-red/green under task scratch.
The fixture now exercises both conditions before the longer transfer cases.
The release-ref-error-classification defect remains outside this deployment
change: the existing publisher hides a failed push's cause when its formatted
command contains --atomic. The read-only dry run and rules response narrow
the failure; unavailable rule-suite audit access limits the server-side cause.

## Traefik dependency extraction (2026-10-26)

PR #203 owns the 3.7.16 dependency pin, checksum, and toolchain-schema correction.
PR #202 retains all role, routing, hardening, and deployment changes. The live
service remains on 3.7.16; extracting the source change performs no deployment.
Merge the dependency update before redeploying the download branch.

## Approved deployment diagram publication (2026-10-26)

The approved preview is now the maintained deployment.mmd: public Yandex above
local XCP-ng, a separate Publisher in each deployment, distinct external-service
groups, one HTTPS connection listing each environment's names, the GitHub
redirect, and separate site/file paths into one storage node. GitHub has its
service label above github.com. The shared-Publisher layout was rejected and
is not published. No operations, DNS, or GitHub Pages were added to this view.

The existing dark WebP render is packaged with the documentation's SVG variants
and linked from the README. The maintained build exactly matches the approved
preview: SHA-256
`0613c11503a0cccf9ed59919163ae4f4c7030a496e68a97f1d6f471b09f15851`.
The README explains the abbreviated site path, controller-side Publisher role,
and browser redirect. Generated images remain build outputs rather than
independently maintained binaries. This source publication does not redeploy
the website or change live infrastructure.

Resume correction: the preview-only restriction was lifted explicitly before
this publication. Earlier generic continuation requests had been treated as
requests to redisplay the image; the remaining work is diagram delivery and
the site owner's first-interaction investigation, not another host rollout.

## Go discovery before browser redirection (2026-10-30)

Stable defect: `go-discovery-redirect`. After rebasing PR #202 onto master
74bbbde92e0cc79366347dc5b40f679e0d5536de, candidate
615ba1903454e8f8e91cbc88063ba1528d1ff71f returned HTTP 301 from
/alwaldend/src?go-get=1 to GitHub in the real isolated proxy fixture. GitHub's
metadata names its own hostname and cannot establish this vanity prefix.

The new higher-priority Traefik Host/Query router sends go-get=1 to Nginx.
Its fixed HTML identifies git.alwaldend.com/alwaldend/src, git, and
https://github.com/alwaldend/src.git. Root, trailing-slash, nested-module, and
extra-query requests return the same metadata without redirection. Unknown
repositories and /alwaldend/src-other return 404 without metadata. Ordinary
and go-get=0 requests retain permanent path/query-preserving GitHub redirects.
Both public Git routers and the metadata virtual host are absent locally;
discovery also works when apex-site routing is disabled.

The fixture rendered the actual templates through the packaged Ansible
launcher and ran the repository-pinned Traefik with real Nginx 1.30.5 from
existing local acceptance image
64718a7a597ce17a1a52c11f388021b9b2f7c40bd8981c46263e16a108149665.
All 42 request checks passed. Download bytes, site serving, unknown-host
rejection, redirects, and discovery were checked through HTTPS. A real pinned
Go client resolved git.alwaldend.com/alwaldend/src/tools/example@v0.0.1 and
reported Git origin https://github.com/alwaldend/src.git, subdirectory
tools/example, and tag refs/tags/tools/example/v0.0.1. Only Git's fetch URL
was rewritten to an isolated local repository; the discovery hostname and
metadata remained the production values.

The retained repeatable harness is out/git-redirect/verify_discovery.py;
DISCOVERY_RUN selects its isolated output directory. The failing baseline and
passing source run are in discovery-before/receipt.json and
discovery-after/receipt.json. It uses already built discovery-ansible,
discovery-traefik, and discovery-go launchers under out/hover-analysis/input-delay.
It uses a loopback-only CONNECT proxy, temporary Git configuration, generated
Traefik TLS, and fixture-scoped Go certificate relaxation. It neither changes
host DNS/trust nor reaches GitHub. Container and proxy processes are removed
after each run. This is local routing/discovery evidence; public network
reachability, the observed timeout, certificate issuance, and live deployment
are outside this fix's validation scope. No live configuration was applied.

## Rebase review, 2026-11-01

PR #202's reconciliation preserves trunk's third-party namespace and explicit
SSH controls. Synchronization reuses its common SSH injector and continues to
mirror projects and sites only. Release fixture paths follow the migrated
`//tools/release/cmd/release` target. The private SSH/mount acceptance verifies
that mounted extracted releases are refused before selection or deletion.
The operations diagram now describes complete project upload promotion and
rollback by extracting a retained archive; the DNS runbook describes the
completed public CNAME cutover. Rendered labels were inspected. The release
owner's [review evidence](../../../../../tools/release/openspec/changes/sync-published-content/evidence.md#rebase-and-cleanup-review-2026-11-01)
records the focused acceptance. Delivery receipts bind final checks to the
rebased candidate; no new deployment, content publication or DNS write occurred.
