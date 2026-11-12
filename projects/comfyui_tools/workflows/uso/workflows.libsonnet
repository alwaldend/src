// USO separates subject latents and SigLIP/projector style references.
local shared = import '../internal/flux1_reference.libsonnet';

function(models)
  local workflows = {
    'single-reference': shared.reference(models, 1, 0, 20, 3.5, 'USO'),
    'style-reference': shared.reference(models, 0, 1, 20, 3.5, 'USO'),
    'character-style': shared.reference(models, 1, 2, 20, 3.5, 'USO'),
  };
  workflows {
    'single-reference-prompt-sequence': shared.sequence(workflows['single-reference'], 'USO-subject'),
    'style-prompt-sequence': shared.sequence(workflows['style-reference'], 'USO-style'),
    'prompt-sequence': shared.sequence(workflows['character-style'], 'USO-combined'),
  }
