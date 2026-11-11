// Native FLUX.2 dev uses distilled guidance and optional reference latents.
local g = import '../internal/graph.libsonnet';
local n = import '../internal/nodes.libsonnet';

function(models) {
  'text-to-image': {
    config: {},
    extra: { ds: { offset: [60, 60], scale: 0.55 } },
    groups: [],
    last_link_id: 12,
    last_node_id: 13,
    version: 0.4,
    nodes: [
      n.node(1, 'UNETLoader', [models.checkpoint, 'default'], [0, 0], [320, 150], title='FLUX.2 dev 32B FP8'),
      n.node(2, 'CLIPLoader', [models.text_encoder, 'flux2', 'default'], [0, 220], [320, 150], title='Local Mistral text encoder'),
      n.node(3, 'VAELoader', [models.vae], [0, 440], [320, 120]),
      n.node(4, 'CLIPTextEncode', ['A watercolor comic panel of a small blue robot reading a book on a bench in a snowy park, clear silhouettes, wide shot.'], [390, 0], [420, 280], title='Scene and visual style'),
      n.node(5, 'FluxGuidance', [4.0], [880, 250], [310, 100]),
      n.node(6, 'BasicGuider', [], [1250, 0], [310, 110], title='Distilled guidance: no negative branch'),
      n.node(7, 'EmptyFlux2LatentImage', [1024, 1024, 1], [390, 520], [310, 150]),
      n.node(8, 'RandomNoise', [123456789, 'randomize'], [880, 0], [310, 120]),
      n.node(9, 'Flux2Scheduler', [28, 1024, 1024], [880, 470], [310, 160]),
      n.node(10, 'KSamplerSelect', ['euler'], [1250, 240], [310, 100]),
      n.node(11, 'SamplerCustomAdvanced', [], [1620, 0], [310, 220]),
      n.node(12, 'VAEDecode', [], [1990, 0], [310, 120]),
      n.node(13, 'SaveImage', ['FLUX2-dev'], [2360, 0], [360, 360]),
    ],
    links: [
      [1, 1, 0, 6, 0],
      [2, 2, 0, 4, 0],
      [3, 4, 0, 5, 0],
      [4, 5, 0, 6, 1],
      [5, 8, 0, 11, 0],
      [6, 6, 0, 11, 1],
      [7, 10, 0, 11, 2],
      [8, 9, 0, 11, 3],
      [9, 7, 0, 11, 4],
      [10, 11, 0, 12, 0],
      [11, 3, 0, 12, 1],
      [12, 12, 0, 13, 0],
    ],
  },
  'multi-reference': g.variant(self['text-to-image'], {
    updates: { '4': { widgets_values: ['Use the character designs and illustration style from reference images 1, 2 and 3. Create a new comic panel of the characters reading books together on a park bench in winter. Change their poses and composition to show the requested scene.'] } },
    nodes: std.flattenArrays([
      local id = 14 + 4 * index;
      local x = 390 * index;
      [
        n.node(id, 'LoadImage', ['reference-' + std.toString(index + 1) + '.png', 'image'], [x, 880], [320, 300], title='Reference ' + std.toString(index + 1) + ': character and style'),
        n.node(id + 1, 'ImageScaleToTotalPixels', ['lanczos', 1, 1], [x, 1250], [320, 120]),
        n.node(id + 2, 'VAEEncode', [], [x, 1440], [320, 120]),
        n.node(id + 3, 'ReferenceLatent', [], [1250 + x, 880], [320, 120]),
      ]
      for index in std.range(0, 2)
    ]),
    removeLinks: [3],
    links: std.flattenArrays([
      local id = 14 + 4 * index;
      local edge = 13 + 4 * index;
      [
        [edge, id, 0, id + 1, 0],
        [edge + 1, id + 1, 0, id + 2, 0],
        [edge + 2, 3, 0, id + 2, 1],
        [edge + 3, id + 2, 0, id + 3, 1],
      ]
      for index in std.range(0, 2)
    ]) + [
      [25, 4, 0, 17, 0],
      [26, 17, 0, 21, 0],
      [27, 21, 0, 25, 0],
      [28, 25, 0, 5, 0],
    ],
    metadata: { last_node_id: 25, last_link_id: 28 },
  }),
  'single-reference': g.variant(self['multi-reference'], {
    removeNodes: std.range(18, 25),
    removeLinks: std.range(17, 28),
    updates: {
      '4': { widgets_values: ['Use the character design and illustration style from reference image 1. Create a new comic panel of the character reading a book on a park bench in winter. Change the pose and composition to show the requested scene.'] },
    },
    links: [
      [17, 4, 0, 17, 0],
      [18, 17, 0, 5, 0],
    ],
    metadata: { last_node_id: 17, last_link_id: 18 },
  }),
  'single-reference-prompt-sequence': g.sequence(self['single-reference'],
                                                 n.node(18, 'ImagegenPromptSequence', ['The referenced robot character enters a snowy park.\n\nThe referenced robot character reads a book on a bench.\n\nThe referenced robot character waves beside a tree.', 'Use the character design and illustration style from reference image 1. Create a new comic panel with clear silhouettes and the requested action.', 'line', 'comic'], [0, 1660], [680, 380], title='Shared reference and panel prompts'),
                                                 4,
                                                 null,
                                                 13,
                                                 [19, 20]),
  'multi-reference-prompt-sequence': g.sequence(self['multi-reference'],
                                                n.node(26, 'ImagegenPromptSequence', ['The referenced robot characters enter a snowy park.\n\nThe referenced robot characters read books on a bench.\n\nThe referenced robot characters wave beside a tree.', 'Use character designs and the illustration style from all three reference images. Create a new comic panel with clear silhouettes and the requested action.', 'line', 'comic'], [0, 1660], [680, 380], title='Shared references and panel prompts'),
                                                4,
                                                null,
                                                13,
                                                [29, 30]),
}
