# Acceptance matrix

Write these cases before implementation. Source and isolated fixture checks
do not establish live deployment. This pass implements IaC and its diagram;
the linked release-tool and browser changes remain separate work.

| Boundary                 | Success                                                               | Failure or preservation case                                                                     | Evidence                                                     |
| ------------------------ | --------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------ |
| Environment provisioning | One explicitly selected VM, 100 GB content disk                       | Other environment absent from the selected state; destructive changes rejected                   | Offline Terraform validation and mocked plans                |
| Storage                  | Mount content before starting services                                | No filesystem overwrite; rerun preserves files and selected link                                 | Ansible fixture and declared disk lifecycle                  |
| Nginx role               | Full configuration validated, service starts                          | Invalid candidate leaves active configuration intact; unchanged rerun does not reload            | Isolated native service fixture                              |
| HTTP downloads           | JSON directory entries, file bytes, HEAD, ranges, CORS                | `index.html` cannot hide listing; staging, state, symlink escapes, and unknown hosts unavailable | HTTP result artifact                                         |
| Website                  | Root and assets follow `current`                                      | Repointing selects another release; absent site returns an error rather than listing             | HTTP fixture                                                 |
| Traefik                  | Direct public HTTP/HTTPS and loopback upstream                        | No ingress hop, mTLS requirement, dashboard, or public Nginx listener                            | Rendered config, sockets, HTTP/TLS fixture                   |
| Website alias            | `www` redirects to apex                                               | Path and query preserved over HTTP and trusted HTTPS                                             | Fixture plus authorized live TLS check                       |
| Credentials              | Dedicated AppRole, own state, selected DNS reads, scoped host signing | No unrelated state or credential paths; no secrets in source or evidence                         | Policy composition and packaging checks                      |
| DNS                      | Public selects Yandex; dc1 selects XCP-ng                             | One owner per name; old apex A/AAAA removed; mail and `www` preserved                            | Offline normalization, source linter, generated declarations |
| Certificate lifecycle    | HTTP-01: Let's Encrypt externally, Vault internally                   | Issuer-specific DNS reachability, EAB scope, CA trust, and renewal                               | Separately authorized live check                             |
| SSH publication          | Dedicated writer; services read content                               | Writer cannot alter service config or ACME secrets                                               | Account/filesystem declarations and fixture                  |
| Release activation       | Verified upload, safe extraction, atomic selection                    | Interrupted/conflicting upload, unsafe archive, disk exhaustion preserve selection               | Linked release-tool E2E (separate implementation)            |
| Browser                  | Downloads/header and reusable release listing                         | Missing service, unsafe names, multiple instances                                                | Linked website E2E (separate implementation)                 |
| Host lifecycle           | Reboot and Ansible rerun preserve releases                            | No backups, replication, automatic deletion, or rollback command                                 | Authorized rollout evidence                                  |

Required identity paths: own `src_infra_download` KV subtree for state and
publisher inputs; its Yandex folder credentials; the existing Cloudflare DNS
token and RouterOS DNS credentials; `src_infra_download_ssh` for host signing.
Internal ACME uses a component PKI role limited to the three serving names. Unrelated AppRole state and provider
administrative credentials are forbidden. Existing infrastructure group
membership supplies the established XO/Yandex identity discovery contract.

## Deduplication extension

The user selected Btrfs with daily background deduplication. Before implementing
it, add these observable cases to the native fixture:

- Two independently uploaded files with identical blocks retain their paths,
  inode identities, permissions, and byte checksums while sharing extents.
- Modifying one deduplicated file leaves the other file's bytes intact.
- A repeated maintenance run uses the persistent hash database and succeeds.
- Scan only published `projects/` and extracted `sites/` content; keep staging,
  ACME state, and the deduplication database outside the scan roots.
- A filesystem mismatch must fail without forced formatting. Storage must be
  mounted before the maintenance service starts. Schedule the service daily
  with low I/O priority and one job at a time.
- Keep HTTP serving available while maintenance runs. Failure of maintenance
  must not delete or activate releases.

Retain the actual maintenance unit, command/version, allocation comparison,
checksums, and copy-on-write result. Live filesystem creation and scheduling
remain deployment checks; an isolated Btrfs fixture is not live deployment.
