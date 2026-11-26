## Context

The Heal checkpoint uses the same native ComfyUI loader, encoder and latent
format as MOSS-Image Base. The model registry supplies its filenames.

## Decisions

Reuse the MOSS generation, layout and sequence builders with Heal-specific
names and output prefixes. Keep the native MOSS sampling preset. Use the
existing native App metadata helpers for text-to-image and prompt sequences.

## Validation

The workflow build and formatting tests passed. Inspection confirmed that all
five Heal files use the registered model assets and that the native Apps expose
the expected controls and save node. All 16 existing MOSS files match the
baseline build byte for byte.
