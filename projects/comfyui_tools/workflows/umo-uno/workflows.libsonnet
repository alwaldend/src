// UMO-UNO retains ordered subject references without USO's style branch.
local shared = import '../internal/flux1_reference.libsonnet';

function(models)
  local workflows = {
    'single-reference': shared.reference(models, 1, 0, 25, 4.0, 'UMO-UNO', referenceSize=512),
    'multi-reference': shared.reference(models, 3, 0, 25, 4.0, 'UMO-UNO', referenceSize=512),
  };
  workflows {
    'single-reference-prompt-sequence': shared.sequence(workflows['single-reference'], 'UMO-UNO-subject'),
    'multi-reference-prompt-sequence': shared.sequence(workflows['multi-reference'], 'UMO-UNO-multiple'),
  }
