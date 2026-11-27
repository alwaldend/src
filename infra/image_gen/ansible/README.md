---
title: Image-generation application installation
---

The playbook loads the profile registry and runtime declaration, translates them
into ordinary role variables, then calls the shared
[comfyui_tools role](../../../projects/ansible_collection/roles/comfyui_tools/README.md).
It installs every selected workflow/App, model file (including LoRAs and
ControlNet), approved Git node pack, first-party node, hash-locked Python
requirement and optional pinned Video2X/RIFE engine. The existing image supplies
ComfyUI and its Python environment. Both fresh installation and updates use
this same role; Terraform only provisions the rental and secured access.

The [model catalog](../../../projects/comfyui_tools/workflows/models.json) defines the profiles.
[Deployment settings](../deployment.json) select the default profile and `workspace_directory`.
The playbook combines that directory with the application name from the
shared runtime declaration. The same workspace setting is passed to Terraform's
native bootstrap. Changing it is not an in-place migration of existing data.

The role preserves inputs, outputs and unselected files. It waits for aggregate
idle counts before installation, loaded-code changes and restart; it never
reads job prompts, history or queue entries. Changed application inputs notify
one restart. Coordinate an idle maintenance window because another client can
still submit work during installation. Workflows can be copied without restart.

After separately authorized account SSH registration and rental creation, use
the observed Terraform SSH host/port outputs and the corresponding local host
identity. The account registration reads its public key from Vault metadata;
it neither transfers nor creates a private key. Pin the server host key in
operator-owned storage before strict SSH checking.

```sh
export IMAGEGEN_SSH_HOST=PROVIDER_REPORTED_HOST
export IMAGEGEN_SSH_PORT=PROVIDER_REPORTED_PORT
export IMAGEGEN_KNOWN_HOSTS="$PWD/out/TASK/known_hosts"
bazel test //infra/image_gen/ansible:ansible.syntax
bazel run //infra/image_gen/ansible:ansible
```

Pass `-- --extra-vars imagegen_model_profile=moss-image-heal` to select a
registered model for one installation. Without this argument the playbook uses
the model selection in `deployment.json`.

Pass `-- --skip-tags workflows` to preserve later App/workflow edits. The syntax
entry point does not authenticate or connect to hosts. Public profiles need no Vault model-token entry. For a gated profile, use
`bazel run //infra/image_gen/ansible:ansible.huggingface`: the existing AL injector
reads the optional model token through the component identity and exports only
`COMFYUI_HUGGING_FACE_TOKEN` to Ansible. The role requires it when the selected
profile declares `hugging_face_auth` and has model downloads. Workflow-only runs
need no token. No Vault or Vast API credential is exported to Ansible or installed
on the rental.

The local role acceptance fixture verifies installation and update safety
without a rental. No deployment is implied by source validation. Changing a
retained instance's Terraform startup can replace it; use this Ansible target
for application updates instead.
