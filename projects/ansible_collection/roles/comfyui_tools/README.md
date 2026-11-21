---
title: ComfyUI tools role
description: Install selected ComfyUI models, node packs, workflows and dependencies
---

This role configures an existing ComfyUI installation from ordinary role
variables. Supply the application directory, runtime arguments and model-directory
map, plus the selected model entries, Git node packs, first-party source paths
and workflow files. Required inputs default to null and are validated before
installation; optional inputs have static defaults.
Application and state directory paths may contain hidden components such as
`.cache`. Supply normalized absolute paths with no `.` or `..` components.

The role knows no profile registry or JSON schema. The deployment playbook
loads its JSON and translates it into those variables once, without mutable
profile facts. Other callers can supply the same variables directly. Large
artifacts remain caller dependencies rather than dependencies of every
collection consumer.

`defaults/main.yaml` documents the complete input contract. Public model URLs
require a publisher commit or SHA-256. Git packs require full commit IDs.
Loopback HTTP model fixtures and local Git repositories are supported for local
acceptance. Gated Hugging Face downloads receive the optional token only on
the publisher request; tasks hide it and redirects do not forward it. The role
accepts `comfyui_tools_hugging_face_token` as a caller input and never contacts
Vault itself. Set `comfyui_tools_require_hugging_face_token` when selected model
downloads require authentication. This deployment obtains the token through AL
injection and maps its environment value in the playbook.
Requirements use the existing Python environment with `--require-hashes --no-deps`.

The role waits for the aggregate `/prompt` idle count before installation,
before changing loaded code and before restart. It reads no prompts, queue
entries, job history, inputs or outputs. This is an operator-coordinated idle
window, not a lock against another client submitting new work. Changed model,
node, dependency or argument inputs notify one restart. Workflows do not require
a restart. Pass `--skip-tags workflows` to preserve operator workflow edits.
Unselected files and generated images are retained.

The role manages `COMFYUI_ARGS` in the supplied environment file. Its baseline
arguments come from the caller; only installed node directories are appended
to the whitelist. The Vast bootstrap preserves this line across reboot.
The optional pinned Video2X installation currently uses `/opt/imagegen`, the
sequence node's established tool contract. Other application paths are configurable.

The [local acceptance fixture](../../test/comfyui_tools/README.md) runs the
packaged role without a GPU, SSH connection, account credential or live service.
