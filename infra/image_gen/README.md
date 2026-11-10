---
title: Image generation
description: Vast.ai ComfyUI rental with Vault-managed HTTPS credentials
---

![Image-generation architecture](assets/architecture.dark.webp)

This project owns one image-generation rental in [tf_setup](tf_setup/README.md).
It uses the reusable [ComfyUI module](../../projects/tf_modules/vastai_instance_imagegen_comfyui/README.md)
in exposed mode: a single HTTPS proxy with a generated username/password login,
plus account-key SSH access. The GPU must have at least 32 GiB VRAM.
The rental's Terraform defaults to the native Qwen-Image-Edit-2511 model bundle,
downloaded on first boot; see [model configuration](tf_setup/README.md#download-models).

## Identity and state

The component authenticates as `src_infra_image_gen`, declared by
[Vault Terraform](../vault/tf/approle_src_infra_image_gen.tf). Its own namespace
at mount `secrets` is `alwaldend.com/vault1/approles/src_infra_image_gen`:

| Entry              | Owner                     | Purpose                                  |
| ------------------ | ------------------------- | ---------------------------------------- |
| `vastai`           | `infra/vastai/alwaldend`  | Dedicated generated team key (`api_key`) |
| `web`              | This project's `tf_setup` | Generated HTTPS username/password        |
| `tf_backend/setup` | This project's backend    | Rental Terraform state and lock          |

[The existing team root](../vastai/alwaldend/README.md) creates and delivers the
key using its own issuer and Vault identity. Its delivery policy covers only
this component's `vastai` entry. The rental reads that entry through its own
AppRole; it never reads the team's operator seed or passes an API key to the
container. The operator can bootstrap the role and browse its KV entries through
existing Vault identity membership. No general infrastructure group is required.

## Deployment status

The former personal stack was destroyed on 2026-11-10. The replacement AppRole,
scoped policies, dedicated team key, generated web login and component state
are deployed. No GPU instance has been rented and there is no active URL.
The [deployment receipt](openspec/changes/add-imagegen-rental/design.md#credential-only-deployment)
records the reviewed scope. No old state or credential was reused. Historical
Vault versions can remain protected; the shared team and operator seed remain.

Deployment is staged: authorize the component role and delivery policy through
Vault Terraform, create the dedicated key through the team root, then rent
through this project's `tf_setup`. Each live apply requires its own reviewed
interactive plan. Building and testing the source performs none of these steps.

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
`bazel build //infra/image_gen/docs:architecture`. The documentation
consumes the rendered target directly; generated images are not checked in.
The diagram describes the configured
architecture, not an active rental. The simplified certificate label refers to
Vast's non-public instance trust chain; see the base module's certificate guidance.
