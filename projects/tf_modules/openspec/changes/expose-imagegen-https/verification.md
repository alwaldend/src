## Candidate and scope

Task branch: `t3/image-generation-exposed-mode`. The module takes an explicit
Vast-native Forge image reference. The personal root resolves the newest published
Forge Neo/CUDA 12.9 tag through a bounded read-only Docker Hub adapter, reusing the
repository's existing Terraform external provider. The user explicitly requested
unpinned selection after initially approving a digest pin. The temporary pin
package was removed. Public model declarations produce a local native-provisioner
manifest. No shared host configuration or paid rental is part of this candidate.
The separately authorized credential apply is recorded below.

## Source and runtime evidence

- Vast's live recommended catalog selected `vastai/sd-forge`, Forge Neo revision
  `d880e8b`, rather than AI-Dock. The then-selected digest was verified against Docker Hub tag
  metadata on 2026-11-09. Count-created values on copied templates do not establish
  a reliable popularity ranking.
- Downloaded selected linux/amd64 layers from that manifest into ignored task
  scratch and verified every downloaded compressed-layer SHA-256. Inspected the
  actual boot, Supervisor, portal inventory, and service scripts. Its image config
  declares no EXPOSE ports. No image/container was started or published.
- Portal and Jupyter bind loopback; TensorBoard's command uses its localhost
  default. The configured startup keeps those tools local and disables outbound tunnel management, Syncthing, and the serverless worker.
  Other local services retain native defaults; SSH-only mode disables Caddy. Tunnels
  and Syncthing outbound relays are disabled. Arbitrary optional provisioning code
  remains trusted and can change these behaviors.
- The bundled Caddy `v2.11.7` accepted the actual proxy configuration using fixture
  TLS files. Loopback HTTPS returned 401 without credentials, 401 for a wrong
  password, 200 for the correct login, and 401 for a URL-token attempt. Plain HTTP
  did not serve the backend. This check found that stdin hashing requires a newline;
  startup now supplies one. Missing TLS files remain a fail-closed startup condition.
- The inspected Caddy configuration manager CLI prints generated credentials.
  Startup imports only its `load_config` inventory function, avoiding that CLI.
  The dedicated proxy provides no portal/token/cookie auth route or admin API.
- Current Vast documentation identifies Secure Cloud as the certified datacenter
  tier. The query uses datacenter/verified constraints and returned offers must
  confirm `hosting_type=1` and `verification=verified`; absent metadata fails closed.

## Reproducible offline coverage

The provider's Terraform CLI suite exercises both modes, immutable runtime
selection, optional/no template, unsafe template rejection, memory and tier
rejection, sensitive credential output, template/mode replacement, legacy-state
upgrade, and observed port URLs. It emits `imagegen-lifecycle.json` and
`imagegen-https.json` in Bazel test outputs. The personal root test uses the real
random provider and mocked Vault with an overridden child module to verify
credential generation/storage without live systems. Owner OpenSpec validation
and Terraform schema/format checks run separately.

Ignored `out/imagegen/` holds bounded logs, inspected public source, local proxy
fixture script and sanitized `proxy-evidence.json`; those fixtures use no real
credentials. Fixture certificate/key/config were removed after the check.

## Limits

No full GPU image boot, actual marketplace selection, management UI behavior,
model loading, or image generation was tested. Live acceptance must verify those
and the externally observed service bindings. Models are supplied through `model_downloads` or SCP/SFTP; no stock provisioning script is inherited. There is no spending timer.
The existing alwaldend root declares a dedicated API token and stores it in the regular user namespace before a rental.
HTTPS/SSH do not protect a running workload against its host administrator.

## Session ergonomics

Observed friction: initial reads exceeded useful output bounds; later reads/logs
were scoped. Two initial catalog calls used unsupported query fields/operators;
reading the official endpoint contract resolved that. Terraform's nested-list
mocking was unsuitable for this root test, so real rental behavior remains in the
CLI fixture and the credential test overrides only the child module. No shared
policy, dependency, or provider configuration was changed to work around failures.

## Latest image and model follow-up

Docker Hub returned 404 for the literal `latest` tag. The personal adapter selects
the newest matching dated Neo/CUDA 12.9 tag from an update-ordered bounded page,
failing closed for registry errors, malformed metadata, or a missing matching
family. Fixture E2E evidence covers its JSON external-data protocol. Actual model
retrieval and future-image startup remain unverified. Generated manifest contents
and Forge destination paths are covered by a real Terraform CLI apply fixture. Previous binary/proxy
checks apply to the observed image only; they are not a guarantee about future tags.

## Dedicated team token follow-up

The user requested creation of a dedicated alwaldend token stored in the regular
user namespace. `infra/vastai/alwaldend` declares the dedicated API-key
resource and its user Vault KV entry alongside its existing team resources. Issuer authentication uses the existing
`src_infra_vastai` role; the infrastructure writer delivers the regular user secret through an exact-path
Vault policy.
The completed authorized handoff moved its resource ownership into the existing
team backend while preserving the key and user Vault entry. The previous
shared-key read policy and new-identity change record have been removed.

Offline credential-root tests passed. After the user renewed the shared host
Vault session, a read-only issuer lookup confirmed the `alwaldend` team identity.
On 2026-11-09, the then-owning interactive
`//users/simeonwarren/imagegen/tf:tf.apply` displayed exactly two additions:
`vastai_api_key.imagegen` and `vault_kv_secret_v2.vastai`, with no changes or
destructions. After review, that apply created the dedicated
`users_simeonwarren_imagegen` API key and stored its value under
`secrets/alwaldend.com/vault1/approles/user_simeonwarren/imagegen/vastai`.
The subsequent owning `:tf.plan` refreshed both resources and reported no changes.
No GPU was rented, no AppRole was changed during that handoff, and no credential
values were printed or committed. The temporary raw issuer response was removed;
only allowlisted operational metadata remains in ignored task evidence.

## Ownership and exposure follow-up

The public image resolver and protocol fixture moved to tools/vastai. The
existing infra/vastai/alwaldend root now declares the dedicated token and regular
user Vault write. The personal rental contains no team-issuer authentication.
The original key and Vault entry were transferred into the existing team state
through the explicitly authorized handoff, without regeneration.

Startup disables only the known outbound access services and omits unrelated
telemetry controls. Its dedicated proxy avoids stock credential logging and
plaintext fallback. The updated stdin launch passed the bundled Caddy loopback
TLS fixture: missing/wrong credentials and URL-token attempts returned 401, the
correct password returned 200, and plaintext did not serve the backend. Full GPU
boot and future-image service bindings remain unverified.

The corrected existing-team implementation passes 40 scoped checks, affected
semantic lint, and the actual catalog fixture for the account/team roots. The
prepared handoff was exercised with public loopback fixtures for success, wrong
key identity, conflicting destination ownership, an existing lock, and a failed
source write after the destination commit. That last case retains both locks
for reconciliation. The later live handoff processed state only in memory under the user's explicit
exception approval; no raw state or credential value was printed or saved.

Session correction: preserving the old backend by introducing a second
credential root added unnecessary ownership structure. The existing alwaldend
root is the authority; a separately authorized state handoff reconciles the
already-created resource without a new root or key rotation.

The authorized handoff completed on 2026-11-09. The normal
`//infra/vastai/alwaldend:tf.plan` refreshed key 36816419 and the same user Vault
entry and reported no changes. Both locks were released, the source state's
current managed resources are empty, and both state lineages were preserved.

The later wiring correction removes all recipient-token lookup and user Vault
provider configuration from the team root. Its existing infrastructure role
gets lifecycle access only to the generated-token KV entry; the user project
continues reading that entry using its own identity for provisioning. The writer
permission is deployed; the Vault owner records fresh live acceptance.

Both scoped permissions are deployed. The current human login discovers the
existing KV mount and can browse its own imagegen folder. The team writer can
manage only the exact generated-token entry in the user namespace; sibling
password and user-backend access remain denied. The existing key and version-1
KV entry are unchanged, and the team root reports no changes using its own
credentials. The reviewed refresh-only state update changes no remote resources.
