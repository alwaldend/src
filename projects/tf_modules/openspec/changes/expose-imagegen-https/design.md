## Gated model follow-up

Live preparation exposed VAST_OFFER_ID_QUERY: an id equality filter returned no
offers even though a GPU-model search returned the selected offer. The native
API accepted ask_contract_id and returned the exact same offer ID. Use that
field while retaining the returned-ID, hosting, memory and price postconditions.
The failed interactive apply made no remote-resource changes.

The user approved FLUX.2 klein 9B and requested a rental. Its publisher gates the
checkpoint. Add an optional sensitive Hugging Face token to the base and ComfyUI
interfaces. The caller reads it from its Vault namespace. Encode it into the
existing sensitive startup command and export it only for the native boot
process, without putting it into the download manifest or persistent global
environment file. Native Hugging Face downloading already consumes HF_TOKEN.
No new downloader or Python package is added to the rental. The selected host
administrator and rental processes can access this token, so use a read-only
token restricted to the model repository. Existing host trust remains accepted.

Failure cases: null token preserves public-download behavior; arbitrary token
characters must not become shell syntax; native provisioning must receive the
value; neither Terraform public outputs nor fixture diagnostics may expose it.

## Context

The user also requested a first-party prompt-list node. Extend native write_files
with a map of approved source modules written only as
`ComfyUI/custom_nodes/<validated-directory>/__init__.py`. Merge their directory names
into the existing custom-node whitelist; do not install another pack or package.
Verify default-empty behavior and prevent directory traversal and name collisions.

Image generation needs SSH-only and password-authenticated HTTPS access with a 32 GiB VRAM floor.

## Goals / Non-Goals

Retain ssh_only by default; exposed publishes only 7860/tcp through a dedicated TLS and Basic-auth Caddy configuration. The application stays on loopback. Reject unsafe template options and weak credentials; disable outbound tunnels and synchronization while retaining loopback-only portal, Jupyter, and TensorBoard services.

## Decisions

Reuse the existing provider, pinned Terraform providers, and AL/Vault configuration flow.
The personal consumer declares a direct published ComfyUI/CUDA 12.9 tag in Terraform. The user explicitly requested removal of the image-discovery executable and selection by tag rather than digest. The reusable module takes the image as an
explicit input and supports declarative public model downloads. No stock template or unrelated provisioning configuration is inherited
by default. Keep the portal, Jupyter, and TensorBoard local; disable outbound
tunnels/synchronization and use a dedicated application-only public proxy.
Secrets remain sensitive in Terraform and reside in the Vault-backed backend; they are never public outputs.
The caller owns template trust and deployment authority. The 32 GB requirement is interpreted as 32 GiB VRAM.

## Risks / Trade-offs

Offline fixtures verify request and lifecycle behavior, not actual image startup or marketplace availability.
HTTPS uses the Vast-provided instance certificate; clients must verify its trust before submitting credentials.

## Migration Plan

Default SSH behavior is retained; changed startup settings intentionally replace a rental.
Create a dedicated alwaldend API token through the existing alwaldend root and store it in the regular user namespace before renting.

## Open Questions

None for implementation. A live acceptance run must inspect the then-selected image and actual service bindings.

## Runtime separation

Proceed with three modules: vastai_instance_imagegen owns the rental, certified
host constraints, sensitive login, native download manifest and dedicated HTTPS
proxy. vastai_instance_imagegen_forge owns Forge paths and native selection;
vastai_instance_imagegen_comfyui owns ComfyUI paths, loopback arguments and
a local-only API-wrapper service. Runtime-specific settings must not remain in the
base. Both wrappers forward the existing rental options and connection outputs.
The image remains an explicit caller-selected tag, without registry discovery.

This reuses Vast's startup/provisioner rather than installing ComfyUI into a
Forge image. The ComfyUI container is a new dependency requiring specific user
approval before declaration. Disable custom nodes, partner APIs and browser
internet access using the native offline flag; do not enable Manager. Preserve
the existing common authentication boundary. Offline acceptance cannot establish
GPU inference or actual image boot; those require a separately authorized rental.

## Chroma component selection

Add the native models/text_encoder destination and an optional structured model
selection. Encode native Forge settings as base64 JSON and pass the existing
--ui-settings-file flag. Set global and preset-specific checkpoint/module keys
because Forge restores preset-specific selections on page load. Selected paths
are derived from the validated manifest. Service exposure and credentials are
unchanged. The existing Terraform CLI fixture checks the staged settings.

The user explicitly approved vastai/comfy:v0.44.0-cuda-12.9-py312 on 2026-11-10.
Use its native ComfyUI service and retain the additional api-wrapper service on loopback.

The approved image declares diagnostic port 3000 in Docker metadata. Keep its
optional instance test runner explicitly disabled; request no publish-all flag.

## Workflow provisioning follow-up

The user requested deployment and then rebase, and identified the missing workflow
from the model README. Keep workflow JSON as ordinary user documents in the personal imagegen owner.
The user subsequently requested including them in startup again. The root reads
those documents and the wrapper compacts and stages them through runtime_files
after native workspace initialization. Changes to startup replace the rental. Reference-style extensions require separate specific dependency
approval and evidence of compatibility with Chroma Flash before declaration.

The user explicitly approved RES4LYF and its required dependencies on 2026-11-10,
then selected Vast's native fetcher. Use typed git_repos, write_files and
pip_packages entries rather than another installer. Pin the exact Git revision
and all Python dependencies with hashes; allowlist only the approved directory.
OpenCV's headless distribution supplies the same required cv2 module without GUI
system libraries. The CLI fixture verifies the actual rendered manifest and allowlist; it does
not install fake fixture packages. Larger startup payloads use gzip/base64 with
a root-private decoded file to reduce Vast request argument size.

Vast rejected the initial expanded lock because startup arguments exceeded its
16,384-byte limit. Use exact Python 3.12 Linux x86_64 wheel URLs and hashes rather
than hashes for unused platforms. Reject oversized rendered startup before rental.

## Provider-owned startup transport

The expanded Juggernaut configuration exceeded the API limit before creating a
rental. The module now supplies plain JSON through quoted, content-derived
heredocs rather than Base64-encoding each file before compression. It opts into
the provider's gzip startup transport; the provider owns encoding and the API
byte-limit guard. Native boot ordering, model provisioner inputs, authentication
and service exposure remain unchanged. The complete selected configuration must
pass its owning interactive plan before the separately authorized rental.

## Native image-adapter model paths

Extend the validated model kinds with ipadapter and clip_vision and map them
only in ComfyUI to models/ipadapter and models/clip_vision. Other runtime maps
still reject unsupported kinds before creation. Native Git/download phases
retain approved immutable references; no custom downloader or public service is
added. The real Terraform CLI fixture covers both access modes and both
checkpoint layouts with the new destinations.
