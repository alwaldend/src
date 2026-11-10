---
title: Vast.ai image generation instance module
description: Shared Vast rentals with SSH-only or password-protected HTTPS access
---

This module rents one GPU with at least **32 GiB VRAM**, using a caller-selected Vast-native image.
The default mode is `ssh_only`; `exposed` publishes the application through HTTPS and a
caller-supplied username/password. The caller supplies provider authentication
and secure state. The [image-generation infrastructure project](../../../infra/image_gen/README.md)
generates and stores its login in Vault before creating the paid instance.
Building and testing these definitions does not deploy or rent capacity.

## Runtime modules

Use the [Forge wrapper](../vastai_instance_imagegen_forge/README.md) or
[ComfyUI wrapper](../vastai_instance_imagegen_comfyui/README.md) for application
configuration. This base owns offer selection, rental lifecycle, authentication,
and the native download manifest. Runtime wrappers own model directories,
loopback arguments, optional JSON configuration and local portal entries.
No image discovery or registry lookup runs during Terraform planning. The caller
selects an explicit published image tag; changing it replaces an existing rental.
The base supplies a plain startup script and opts into the provider's gzip
transport. The provider owns encoding and the 16,384-byte API command limit;
the selected Vast images provide the required shell utilities.

By default no `template_hash_id` is sent, avoiding the stock template's public
ports and unrelated provisioning configuration. Configure `model_downloads` for
automatic first-boot downloads rather than editing a public template.

An optional reviewed Vast-native `template_hash_id` can provide environment
settings, including model provisioning. Copy and review it first, removing
all Docker flags except literal `-e NAME=value` or `--env NAME=value` assignments.
Missing/malformed metadata, published ports, host networking, and other Docker
flags fail closed. The explicit image overrides the template image. Credentials
must never be stored in a saved/public template. Mutable scripts/URLs and
arbitrary extensions remain trusted code, outside the network guarantees below.

## Access and local services

| Mode                 | Public application access                  | Authentication                        |
| -------------------- | ------------------------------------------ | ------------------------------------- |
| `ssh_only` (default) | No web port requested; Caddy disabled      | Vast account SSH key and forwarding   |
| `exposed`            | Only application proxy, container 7860/tcp | TLS plus HTTP Basic username/password |

Each wrapper fixes its application listener to loopback in both modes. Local tools are retained:
Instance Portal at `127.0.0.1:11111`, Jupyter at `127.0.0.1:18080`, and TensorBoard
at `127.0.0.1:16006`. Access them through SSH forwarding to the corresponding
internal port. They are not published by the proxy. These local management
interfaces can provide unauthenticated shell/file access, so only trusted local
processes and SSH users belong on this instance.

Cloudflare tunnel management, Syncthing (including outbound relays), and the
serverless worker are disabled. Other bundled services retain their native defaults. The inspected ComfyUI image declares diagnostic `3000/tcp` in its Docker metadata;
the instance test runner is explicitly disabled and no publish-all flag is used. The module requests only the SSH service, plus
`-p 7860:7860` in exposed mode. Verify actual listeners during live acceptance;
reviewed provisioning code can still start additional listeners or outbound traffic.

Exposed mode requires `web_password` with 32–64 alphanumeric characters and a
`web_username` starting with a letter and containing letters, digits, or underscores
(maximum 64 characters). Generate a random password in the caller. The restricted
alphabet prevents shell substitution. The dedicated proxy hashes the password
with bcrypt, authenticates every path, and has no portal, URL-token/cookie bypass,
plaintext redirect, access log, or Caddy admin API. HTTP/3 is disabled. Missing
or invalid TLS files fail proxy startup instead of falling back to plaintext.
It uses the selected image's bundled Caddy, not an additional downloaded proxy.

Portal/CLI auto-updates are disabled so boot uses the selected image's bundled
services. Startup controls service exposure; it is not an outbound firewall or a
privacy audit of arbitrary application code and extensions. Credentials occur
in sensitive `onstart`, are sent to Vast, and remain in Terraform state/history.
There are no password outputs; restrict the Vault-backed state accordingly.

## Compose

Compose a runtime wrapper rather than configuring the low-level runtime contract
here. Package its filegroup in the consuming Terraform target's `data`, and select
`//projects/vastai_terraform_provider:provider`. Both wrappers forward the base's
rental inputs and connection outputs. No Vault resources belong in these shared
modules; callers generate and persist their own login before rental creation.

Offer selection requests one verified datacenter NVIDIA GPU, an amd64 host with
CUDA 12.9 compatibility, at least 32 GiB VRAM, 99% reported reliability, 100 GB disk
by default, and one direct port for SSH or two for SSH plus HTTPS. Vast calls its
certified tier [Secure Cloud (Datacenter)](https://docs.vast.ai/guides/instances/choosing/find-and-rent).
The returned offer must confirm `hosting_type=1` and `verification=verified`;
missing tier data, insufficient memory, and excessive hourly price are rejected.
There is no fallback to community hosts or smaller GPUs. `min_gpu_ram_gb` can
raise the floor. `offer_id` retains all constraints.

The default USD 1/hour ceiling includes requested disk; it does not cap transfer
charges, continuing storage, or total duration. Marketplace changes do not replace
an existing contract, but each normal plan rechecks availability; an empty eligible
market can block planning. Mode, credentials, image, template, startup, or disk
changes replace the instance. Export wanted images before replacement/destruction.

## Model downloads

`model_downloads` is a list of `{url, filename, kind}` objects. The URL must be
public HTTPS without embedded credentials; filenames end in `.safetensors` and
cannot traverse directories. `kind` selects a directory supplied by the runtime wrapper. The base validates
public URLs, filenames, uniqueness and the presence of a destination for each
kind; it contains no application-specific paths.

Startup writes a local JSON/YAML-compatible manifest and passes it to Vast's
[native provisioner](https://github.com/vast-ai/base-image/blob/main/ROOT/opt/instance-tools/lib/provisioner/README.md).
The native application launcher waits while `/.provisioning` exists. Provisioning
runs on the rented machine and downloads to its billable local storage; there
is no model cache preserved after destruction. An empty list downloads no
configured models. SCP/SFTP remains available for manual file transfer.
This interface currently supports public files; do not put gated-model tokens
in URLs or source. Gated downloads need a separately configured Vault-backed
`HF_TOKEN`/`CIVITAI_TOKEN` injection. Review model licenses and required space.

## Connect and certificates

Configure account SSH public keys before renting. Creation only acknowledges a
contract. After boot, refresh through the caller's owning workflow. The
`ssh_forward_command` remains null until SSH details are observed; run it locally
and open `local_url` (default `http://127.0.0.1:7860`). Both forwarding endpoints
bind loopback, using the caller's SSH agent/normal identity configuration.

In exposed mode, `https_url` remains null until refresh reports both public IP
and the actual mapped `7860/tcp` host port. It does not establish application
readiness.

## Certificates

The proxy uses `/etc/instance.crt` and `/etc/instance.key`. Follow
[Vast's certificate guidance](https://docs.vast.ai/guides/instances/connect/jupyter)
and verify issuer/fingerprint through a trusted channel before entering the login.
Do not disable TLS validation or blindly bypass warnings. No DNS or publicly
trusted domain certificate is configured by this module.

## Privacy

Applications run on the rented machine; output files and prompt/image metadata
may persist locally. Extensions and provisioning code may access them and make
network requests. This module adds no prompt-sharing service. Runtime-specific
controls are described by each wrapper; none establish confidentiality from the
host administrator or a no-training commitment for the complete stack.

[Vast's privacy policy](https://vast.ai/privacy) covers account/website processing.
Certified datacenter status does not prevent the host administrator accessing the
workload. See the [infrastructure project's privacy boundary](../../../infra/image_gen/README.md#privacy-boundary)
for host trust and retention limits. Stopping retains billable storage; export
wanted files before destruction. This project cannot verify physical erasure or
exclude host-made copies.

## Validation

```sh
bazel test //projects/tf_modules/vastai_instance_imagegen:all
bazel test //projects/vastai_terraform_provider/test/acceptance:acceptance_test
```

The base and both wrappers run through pinned Terraform/provider against a loopback API
fixture: both modes, unsafe templates/credentials, memory/tier rejection, sensitive
output, replacement, and mapped URLs. `imagegen-lifecycle.json` and
`imagegen-https.json` are repeatable undeclared test outputs. No GPU is rented.
Actual image boot, TLS/authentication, management-service bindings, model loading,
and image generation remain unverified without a scoped live acceptance run.
