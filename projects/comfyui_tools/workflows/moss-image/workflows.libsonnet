// Native MOSS examples adapted from Comfy-Org/workflow_templates at
// 36ecd0663961c5e1e83f7c0ea521658b33b36bbc; see the project README.
local g = import '../internal/graph.libsonnet';
local n = import '../internal/nodes.libsonnet';
local caption = 'A small blue robot reading a book on a wooden bench in a snowy park, watercolor illustration, clear silhouettes, soft winter light.';
local negative = 'worst quality, bad anatomy, watermark, jpeg artifact';
local boxes = [
  { x: 200, y: 160, width: 620, height: 720, metadata: { type: 'obj', text: '', desc: 'A small blue robot reading a book on a wooden bench', palette: [] } },
];
local panels = 'A small blue robot reads a book on a snowy park bench at dawn.\n\nA small blue robot reads a book on a snowy park bench at noon.\n\nA small blue robot reads a book on a snowy park bench at sunset.';
local sequenceNode(id, prefix, prompts=panels, context='Watercolor illustration, clear silhouettes, soft light.') =
  n.node(id, 'ImagegenPromptSequence', [prompts, context, 'line', prefix], [0, 1640], [680, 380], title='Panel prompts and shared negative text');

function(models) {
  'text-to-image': {
    config: {},
    extra: { ds: { offset: [60, 60], scale: 0.6 } },
    groups: [],
    last_link_id: 10,
    last_node_id: 10,
    version: 0.4,
    nodes: [
      n.node(1, 'UNETLoader', [models.checkpoint, 'default'], [0, 0], [340, 120], title='MOSS-Image SFT INT8'),
      n.node(2, 'CLIPLoader', [models.text_encoder, 'moss_image', 'default'], [0, 200], [340, 140]),
      n.node(3, 'VAELoader', [models.vae], [0, 420], [340, 100]),
      n.node(4, 'CLIPTextEncode', [caption], [420, 0], [420, 260], title='Example scene'),
      n.node(5, 'CLIPTextEncode', [negative], [420, 340], [420, 180]),
      n.node(6, 'EmptyFlux2LatentImage', [1024, 1024, 1], [420, 620], [320, 140]),
      n.node(7, 'SkipLayerGuidanceDiTSimple', ['', '22', 0, 0.5], [900, 0], [340, 160]),
      n.node(8, 'KSampler', [42, 'fixed', 25, 2.5, 'er_sde', 'simple', 1], [1300, 0], [340, 300]),
      n.node(9, 'VAEDecode', [], [1700, 0], [320, 100]),
      n.node(10, 'SaveImage', ['MOSS-Image-example'], [2080, 0], [360, 360]),
    ],
    links: [
      [1, 1, 0, 7, 0],
      [2, 2, 0, 4, 0],
      [3, 2, 0, 5, 0],
      [4, 7, 0, 8, 0],
      [5, 4, 0, 8, 1],
      [6, 5, 0, 8, 2],
      [7, 6, 0, 8, 3],
      [8, 8, 0, 9, 0],
      [9, 3, 0, 9, 1],
      [10, 9, 0, 10, 0],
    ],
  },
  'base-text-to-image': g.variant(self['text-to-image'], {
    updates: {
      '1': { widgets_values: [models.base_checkpoint, 'default'], title: 'MOSS-Image Base INT8' },
      '10': { widgets_values: ['MOSS-Image-Base-example'] },
    },
  }),
  'layout-guided': g.variant(self['text-to-image'], {
    nodes: [
      n.node(11, 'CreateBoundingBoxes', [1024, 1024, boxes, boxes], [0, 850], [600, 650], title='Example layout: edit the object box'),
      n.node(12, 'MOSSImageLayoutPrompt', [caption, 1024, 1024, '1'], [750, 850], [440, 280]),
      n.node(13, 'GetImageSize', [], [750, 1240], [320, 100]),
    ],
    updates: {
      '4': { inputs: [{ name: 'clip', type: 'CLIP' }, { name: 'text', type: 'STRING', widget: { name: 'text' } }] },
      '10': { widgets_values: ['MOSS-Image-layout-example'] },
    },
    links: [
      [11, 11, 2, 12, 0],
      [12, 12, 0, 4, 1],
      [13, 11, 0, 13, 0],
      [14, 13, 0, 6, 0],
      [15, 13, 1, 6, 1],
      [16, 13, 0, 12, 1],
      [17, 13, 1, 12, 2],
    ],
    metadata: { last_node_id: 13, last_link_id: 17 },
  }),
  'base-layout-guided': g.variant(self['layout-guided'], {
    updates: {
      '1': { widgets_values: [models.base_checkpoint, 'default'], title: 'MOSS-Image Base INT8' },
      '10': { widgets_values: ['MOSS-Image-Base-layout-example'] },
    },
  }),
  'image-edit': g.variant(self['text-to-image'], {
    nodes: [
      n.node(11, 'LoadImage', ['reference.png', 'image'], [0, 850], [340, 320], title='Select the image to edit'),
      n.node(12, 'ImageScaleToTotalPixels', ['lanczos', 1, 16], [420, 850], [320, 120]),
      n.node(13, 'GetImageSize', [], [800, 850], [320, 100]),
      n.node(14, 'LoraLoaderModelOnly', [models.edit_lora, 1], [0, -200], [400, 120]),
    ],
    updates: {
      '4': n.node(4, 'TextEncodeMOSSImageEdit', ['Change the robot to red. Keep the scene, pose and illustration style unchanged.', 'Moss-Image-Edit', 1], [420, 0], [420, 260]),
      '5': n.node(5, 'TextEncodeMOSSImageEdit', [negative, 'Moss-Image-Edit', 1], [420, 340], [420, 180]),
      '7': { widgets_values: ['', '22', 0, 1] },
      '8': { widgets_values: [42, 'fixed', 30, 4, 'euler', 'normal', 1] },
      '10': { widgets_values: ['MOSS-Image-edit-example'] },
    },
    removeLinks: [1],
    links: [
      [11, 1, 0, 14, 0],
      [12, 14, 0, 7, 0],
      [13, 11, 0, 12, 0],
      [14, 12, 0, 13, 0],
      [15, 13, 0, 6, 0],
      [16, 13, 1, 6, 1],
      [17, 12, 0, 4, 1],
      [18, 12, 0, 5, 1],
      [19, 3, 0, 4, 3],
      [20, 3, 0, 5, 3],
    ],
    metadata: { last_node_id: 14, last_link_id: 20 },
  }),
  'color-depth-control': g.variant(self['text-to-image'], {
    nodes: [
      n.node(11, 'LoadImage', ['color-map.png', 'image'], [0, 850], [340, 320], title='Prepared flat region colors, 1024 square'),
      n.node(12, 'LoadImage', ['depth-map.png', 'image'], [420, 850], [340, 320], title='Prepared depth: near is white, 1024 square'),
      n.node(13, 'LoraLoaderModelOnly', [models.control_lora, 1], [0, -200], [420, 120]),
      n.node(14, 'MOSSImageControl', [1, 0, 0.6], [900, 850], [340, 180]),
      n.node(15, 'TemporalScoreRescaling', [0.95, 1], [900, -400], [340, 120]),
      n.node(16, 'ModelSamplingAuraFlow', [5, 'flow'], [900, -200], [340, 120]),
      n.node(17, 'SamplerER_SDE', ['ER-SDE', 3, 1, 1.05], [900, 300], [340, 160]),
      n.node(18, 'BasicScheduler', ['simple', 25, 1], [900, 550], [340, 160]),
    ],
    updates: {
      '7': { widgets_values: ['', '24', 0, 1] },
      '8': n.node(8, 'SamplerCustom', [true, 42, 'fixed', 3], [1300, 0], [340, 300]),
      '10': { widgets_values: ['MOSS-Image-control-example'] },
    },
    removeLinks: [1, 4, 5, 6, 7],
    links: [
      [11, 1, 0, 7, 0],
      [12, 7, 0, 15, 0],
      [13, 15, 0, 16, 0],
      [14, 16, 0, 13, 0],
      [15, 13, 0, 8, 0],
      [16, 13, 0, 18, 0],
      [17, 4, 0, 14, 0],
      [18, 5, 0, 14, 1],
      [19, 3, 0, 14, 2],
      [20, 11, 0, 14, 3],
      [21, 12, 0, 14, 4],
      [22, 14, 0, 8, 1],
      [23, 14, 1, 8, 2],
      [24, 17, 0, 8, 3],
      [25, 18, 0, 8, 4],
      [26, 6, 0, 8, 5],
    ],
    metadata: { last_node_id: 18, last_link_id: 26 },
  }),
  'prompt-sequence': g.sequence(self['text-to-image'], sequenceNode(11, 'moss-image'), 4, 5, 10, [11, 12, 13]),
  'base-prompt-sequence': g.sequence(self['base-text-to-image'], sequenceNode(11, 'moss-image-base'), 4, 5, 10, [11, 12, 13]),
  'layout-prompt-sequence': g.sequence(self['layout-guided'], sequenceNode(14, 'moss-image-layout'), 12, 5, 10, [18, 19, 20], textInput='prompt', negativeTextInput='text'),
  'base-layout-prompt-sequence': g.sequence(self['base-layout-guided'], sequenceNode(14, 'moss-image-base-layout'), 12, 5, 10, [18, 19, 20], textInput='prompt', negativeTextInput='text'),
  'image-edit-prompt-sequence': g.sequence(self['image-edit'], sequenceNode(15, 'moss-image-edit', prompts='Change the robot to red. Keep the scene, pose and illustration style unchanged.\n\nChange the robot to green. Keep the scene, pose and illustration style unchanged.\n\nChange the robot to yellow. Keep the scene, pose and illustration style unchanged.', context=''), 4, 5, 10, [21, 22, 23], textInput='prompt'),
  'color-depth-prompt-sequence': g.sequence(self['color-depth-control'], sequenceNode(19, 'moss-image-control'), 4, 5, 10, [27, 28, 29]),
}
