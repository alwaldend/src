---
title: Vast.ai ComfyUI instance module
description: Loopback ComfyUI with native model provisioning and authenticated HTTPS
---

This module composes [the base rental](../vastai_instance_imagegen/README.md)
with the native Vast ComfyUI service at `127.0.0.1:18188`. Supply a reviewed
`vastai/comfy` image through `image`. It forwards the base's access, credential,
offer, storage and connection interfaces, retaining SSH-only access by default.

ComfyUI uses `--disable-all-custom-nodes --offline`; Manager is not enabled and unapproved
custom nodes are disabled, partner API nodes are disabled, and the frontend
receives ComfyUI's restrictive offline Content-Security-Policy. No wildcard CORS
flag is supplied. The bundled API wrapper remains on loopback `127.0.0.1:18288` for scripted
jobs. It is included in the local portal inventory but never routed by the
public proxy. Portal, Jupyter and TensorBoard also remain local; tunnels, sync
and worker services remain disabled. These controls limit exposure; they are not an outbound firewall.

Model kinds map below `/workspace/ComfyUI`:

| Kind           | Directory                 |
| -------------- | ------------------------- |
| `checkpoint`   | `models/diffusion_models` |
| `lora`         | `models/loras`            |
| `vae`          | `models/vae`              |
| `embedding`    | `models/embeddings`       |
| `text_encoder` | `models/text_encoders`    |
| `ipadapter`    | `models/ipadapter`        |
| `clip_vision`  | `models/clip_vision`      |

Here `checkpoint` means the standalone diffusion model, loaded with ComfyUI's
Load Diffusion Model node. For an all-in-one SDXL checkpoint loaded by CheckpointLoaderSimple, set
`checkpoint_directory = "models/checkpoints"`. The only allowed alternatives
are that directory and the default `models/diffusion_models`; other model kinds
retain the paths above.
`workflows` accepts named native JSON documents and stages them in
`user/default/workflows` after Vast initializes the workspace. The wrapper
compacts the JSON; the base supplies plain JSON in its startup script and enables
the provider's gzip transport. The provider rejects encoded commands exceeding
Vast's 16 KB limit before creating a rental.
`custom_node_sources` accepts approved first-party Python modules keyed by
validated directory names. The native file provisioner stages each module as
`custom_nodes/<name>/__init__.py`; only those directories and approved Git packs
are whitelisted. No packages are installed for first-party source modules.
`custom_node_web_scripts` optionally supplies one `web/index.js` per declared
first-party module. ComfyUI serves it when the module exports `WEB_DIRECTORY`;
script keys cannot introduce additional node directories or change the whitelist.
`custom_nodes` accepts explicitly approved public GitHub repositories at full
commit IDs. Vast's built-in `git_repos` fetcher installs them, and ComfyUI's
whitelist enables only those directories. `python_requirements` supplies a
version- and hash-locked requirements file installed with `--require-hashes
--no-deps` in the existing application environment. No custom installer or image
rebuild is needed. Downloads are served
by Vast's provisioner inside the rented instance, with no model cache retained
after destruction.

`hugging_face_token` optionally passes a sensitive, Vault-sourced read-only
token to native first-boot downloading. Public URLs and the manifest contain
no credentials. The token's account must already have gated repository access;
this module does not accept licenses. The token is available to the rented
host and its processes. It is encoded into sensitive startup configuration,
without adding it to `/etc/environment` or a public Terraform output.

```hcl
module "imagegen" {
  source = "../../../projects/tf_modules/vastai_instance_imagegen_comfyui"
  image = "vastai/comfy:v0.44.0-cuda-12.9-py312"
}
```

The user approved this Vast image specifically. Its public source and bundled
services still require maintenance review when upgrading; no registry checker
runs during planning. [Vast's source](https://github.com/vast-ai/base-image/tree/main/derivatives/pytorch/derivatives/comfyui)
owns the boot contract. [ComfyUI's native arguments](https://github.com/Comfy-Org/ComfyUI/blob/v0.44.0/comfy/cli_args.py)
own the offline and custom-node behavior. Offline Terraform acceptance covers
both access modes, native model paths and proxy/forwarding ports; actual image
boot and GPU generation require separately authorized live acceptance.

Fixed runtime settings are declared in runtime.json and merged with the
selected checkpoint directory and approved node whitelist. The ephemeral
evaluator consumes the same declaration instead of maintaining another set of
ports, command-line flags and native model directories.
