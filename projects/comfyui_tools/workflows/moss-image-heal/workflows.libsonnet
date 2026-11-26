local g = import '../internal/graph.libsonnet';
local moss = import '../moss-image/workflows.libsonnet';

function(models)
  local generation = moss(models, {
    title: 'MOSS-Image-heal INT8',
    image_prefix: 'MOSS-Image-heal',
    sequence_prefix: 'moss-image-heal',
  });
  {
    'text-to-image': generation['text-to-image'],
    'text-to-image.app': generation['text-to-image.app'],
    'layout-guided': generation['layout-guided'],
    'prompt-sequence': generation['prompt-sequence'],
    'prompt-sequence.app': g.app(self['prompt-sequence'], [
      [11, 'prompts', { height: 180, description: 'Enter one panel prompt per nonempty line.' }],
      [11, 'context', { height: 100, description: 'Add details shared by every panel.' }],
      [11, 'negative_prompt', { height: 100, description: 'Describe what to avoid in every panel.' }],
      [6, 'width'],
      [6, 'height'],
      [8, 'seed'],
      [8, 'steps'],
      [8, 'cfg'],
    ], 10),
  }
