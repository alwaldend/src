---
title: Image generation
description: Personal Vast.ai ComfyUI rental with Vault-managed HTTPS credentials
---

![Image-generation architecture](assets/architecture.dark.webp)

This project owns one image-generation rental in [tf_setup](tf_setup/README.md).
It uses the reusable [ComfyUI module](../../../projects/tf_modules/vastai_instance_imagegen_comfyui/README.md)
in exposed mode: a single HTTPS proxy with a generated username/password login,
plus account-key SSH access. The GPU must have at least 32 GiB VRAM.
The rental's Terraform defaults to the publisher's full Chroma1-HD-Flash checkpoint with its T5 encoder and VAE,
downloaded on first boot; see [model configuration](tf_setup/README.md#download-models).

The project reuses the existing `user_simeonwarren` identity and regular secret
namespace. [Infrastructure credential setup](../../../infra/vastai/alwaldend/README.md) creates a dedicated API token in the
alwaldend team and writes it to `imagegen/vastai` in that user namespace. Token generation belongs to infra/vastai; this user project only reads the
result. The existing alwaldend root uses its team issuer and infrastructure Vault
identity to generate and deliver the key. Its delivery permission covers only
this KV entry. The user project authenticates as itself to read that key and
provision the instance. No extra AppRole or shared rental-token copy is created.

The [rental root](tf_setup/README.md) reads that user-owned token. Its generated
HTTPS login is stored at `imagegen/web`, with independent rental state at
`tf_backend/imagegen`. The alwaldend root owns token generation in its existing team state and never
rents a GPU. The existing key was preserved during its completed state handoff. No shared host configuration is changed.

## Privacy boundary

The configured Vast ComfyUI image runs ComfyUI on the rented host. It is not a separate hosted
inference API, and this configuration adds no prompt/image upload service or
access logging. ComfyUI can retain generated images and prompt metadata locally;
provisioning scripts and extensions can access those files and the internet.
The image is selected directly by tag in the rental Terraform. Review optional templates, scripts, models, and extensions before use.

[Vast's privacy policy](https://vast.ai/privacy) covers its account and website
processing; it does not establish a blanket no-retention or no-training promise
for this stack. [Vast's security FAQ](https://github.com/vast-ai/docs/blob/main/guides/reference/faq/security.mdx)
describes tenant isolation and recommends Secure Cloud providers for sensitive
workloads. The module requests datacenter/verified offers and checks returned
`hosting_type=1` and `verification=verified`, rejecting absent tier metadata.
This selects Vast's certified datacenter tier; it does not establish
confidentiality against the host administrator. SSH and HTTPS
protect transport, not a workload from its machine's administrator.

Keep wanted images locally and export them before releasing the rental. Stopping
retains billable storage. Destruction removes the contract's local storage but
this project cannot verify physical erasure or exclude host-made copies.
Generated login credentials also remain in the Vault-backed Terraform state
and may remain in historical Vault secret/state versions. There are no public
password outputs. Destroying the rental does not revoke copies already disclosed.

## Architecture diagram

The diagram source is [docs/architecture.mmd](docs/architecture.mmd). It uses the
shared Alwaldend dark palette and typography, composed with the existing Mermaid
actor artwork for the requested browser actor. The local template derives both
from their owners. Build the WebP with
`bazel build //users/simeonwarren/imagegen/docs:architecture`. The documentation
consumes the rendered target directly; generated images are not checked in.
The diagram describes the configured
architecture, not an active rental. The simplified certificate label refers to
Vast's non-public instance trust chain; see the base module's certificate guidance.
