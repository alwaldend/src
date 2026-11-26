## Context

Terraform provisions the rental. Ansible installs the selected model profile.

## Decisions

Use a dedicated `moss-image-heal` profile with KBlueLeaf's pinned INT8 ConvRot
checkpoint, the existing Qwen3.5 text encoder and the small FLUX.2 VAE. Pass
the profile name through the existing Ansible selector. Reuse the first-party
Prompt Sequence node and the workflow package's generated files.

## Validation

The pinned Hugging Face model metadata confirms the checkpoint digest. The
`//infra/image_gen/test/ansible:assets_test` fixture copied all five Heal files
and verified their contents against the package. GPU execution remains unverified.
