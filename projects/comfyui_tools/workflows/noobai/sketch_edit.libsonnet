// Separate starting pixels and desired structural guidance; no runtime extension.
local g = import '../internal/graph.libsonnet';
local n = import '../internal/nodes.libsonnet';

function(models, settings, quality)
  local canvasInputs = [
    { name: 'width', type: 'INT', widget: { name: 'width' } },
    { name: 'height', type: 'INT', widget: { name: 'height' } },
  ];
  local base = {
    version: 0.4,
    config: {},
    last_node_id: 15,
    last_link_id: 22,
    groups: [],
    nodes: [
      n.node(1, 'CheckpointLoaderSimple', [models.checkpoint], [0, 0], [320, 150]),
      n.node(2, 'LoadImage', ['source.png', 'image'], [0, 220], [320, 310], title='Source image: starting appearance'),
      n.node(3, 'ImageScale', ['lanczos', 1024, 1024, 'disabled'], [370, 220], [320, 180], title='Output canvas: source dimensions'),
      n.node(4, 'GetImageSize', [], [370, 450], [320, 100]),
      n.node(5, 'LoadImage', ['sketch.png', 'image'], [0, 620], [320, 310], title='Separate sketch: desired final outlines'),
      n.node(6, 'ImageScale', ['lanczos', 1024, 1024, 'disabled'], [370, 650], [320, 180], title='Align sketch to output canvas') { inputs+: canvasInputs },
      n.node(7, 'Canny', [0.1, 0.3], [740, 650], [320, 150], title='Extract sketch edges'),
      n.node(8, 'ControlNetLoader', [models.controlnet], [740, 850], [320, 100]),
      n.node(9, 'CLIPTextEncode', [quality + ', illustration, robot, blue body, red chest plate, standing, raised right hand, waving, snowy park'], [740, 0], [400, 230], title='Describe the complete desired result'),
      n.node(10, 'CLIPTextEncode', [settings.negative + ', extra limbs, duplicate, text, watermark'], [740, 270], [400, 180], title='Negative prompt'),
      n.node(11, 'ControlNetApplyAdvanced', [0.7, 0.0, 0.8], [1190, 300], [320, 240], title='Sketch influence: zero disables guidance'),
      n.node(12, 'VAEEncode', [], [740, 1040], [320, 120], title='Encode source pixels'),
      n.node(13, 'KSampler', [42, 'fixed', settings.steps, settings.cfg, settings.sampler, settings.scheduler, 0.65], [1580, 0], [320, 320], title='Source change strength: denoise'),
      n.node(14, 'VAEDecode', [], [1960, 0], [320, 140]),
      n.node(15, 'SaveImage', ['NoobAI-sketch-edit'], [2340, 0], [360, 350]),
    ],
    links: [
      [1, 2, 0, 3, 0],
      [2, 3, 0, 4, 0],
      [3, 5, 0, 6, 0],
      [4, 4, 0, 6, 1],
      [5, 4, 1, 6, 2],
      [6, 6, 0, 7, 0],
      [7, 1, 1, 9, 0],
      [8, 1, 1, 10, 0],
      [9, 9, 0, 11, 0],
      [10, 10, 0, 11, 1],
      [11, 8, 0, 11, 2],
      [12, 7, 0, 11, 3],
      [13, 1, 2, 11, 4],
      [14, 3, 0, 12, 0],
      [15, 1, 2, 12, 1],
      [16, 1, 0, 13, 0],
      [17, 11, 0, 13, 1],
      [18, 11, 1, 13, 2],
      [19, 12, 0, 13, 3],
      [20, 13, 0, 14, 0],
      [21, 1, 2, 14, 1],
      [22, 14, 0, 15, 0],
    ],
    extra: {
      ds: { offset: [40, 60], scale: 0.6 },
      linearMode: false,
      linearData: {
        inputs: [
          [2, 'image', { description: 'Original image to edit. The source and sketch must use matching framing.' }],
          [5, 'image', { description: 'Separate sketch showing the desired final outlines. Remove obsolete outlines; omit arrows and annotations.' }],
          [9, 'text', { height: 180, description: 'Describe the finished image with NoobAI tags, including the desired pose and appearance.' }],
          [13, 'denoise', { description: 'Higher values permit larger changes; lower values preserve more of the source.' }],
          [11, 'strength', { description: 'Sketch influence. Zero disables it for a fixed-seed comparison.' }],
          [11, 'end_percent'],
          [7, 'low_threshold'],
          [7, 'high_threshold'],
          [3, 'width', { description: 'Output width; use multiples of eight. Source, sketch and mask share this canvas.' }],
          [3, 'height', { description: 'Output height; use multiples of eight and preserve the reference aspect ratio.' }],
          [13, 'seed'],
          [13, 'steps'],
          [13, 'cfg'],
          [10, 'text', { height: 100 }],
          [15, 'filename_prefix'],
        ],
        outputs: [15],
      },
    },
  };
  local edit = base;
  local masked = g.variant(edit, {
    nodes: [
      n.node(16, 'MaskToImage', [], [0, 1100], [320, 120], title='Source painted mask: white permits editing'),
      n.node(17, 'ImageScale', ['nearest-exact', 1024, 1024, 'disabled'], [370, 1100], [320, 180], title='Align source mask to canvas') { inputs+: canvasInputs },
      n.node(18, 'ImageToMask', ['red'], [370, 1340], [320, 120]),
      n.node(19, 'SetLatentNoiseMask', [], [1190, 1040], [320, 140]),
      n.node(20, 'ImageCompositeMasked', [0, 0, false], [2340, 450], [320, 180], title='Preserve resized source outside the mask'),
    ],
    updates: {
      '2': { title: 'Source image: paint the allowed edit area in Mask Editor' },
      '15': { widgets_values: ['NoobAI-sketch-edit-masked'], pos: [2740, 0] },
    },
    removeLinks: [19, 22],
    links: [
      [23, 2, 1, 16, 0],
      [24, 16, 0, 17, 0],
      [25, 4, 0, 17, 1],
      [26, 4, 1, 17, 2],
      [27, 17, 0, 18, 0],
      [28, 12, 0, 19, 0],
      [29, 18, 0, 19, 1],
      [30, 19, 0, 13, 3],
      [31, 3, 0, 20, 0],
      [32, 14, 0, 20, 1],
      [33, 18, 0, 20, 4],
      [34, 20, 0, 15, 0],
    ],
    metadata: {
      last_node_id: 20,
      last_link_id: 34,
      extra: edit.extra {
        linearData+: {
          inputs: [[2, 'image', { description: 'In Graph Mode, open this source in Mask Editor and save the painted edit area. White permits edits; an empty mask preserves the source.' }]] + edit.extra.linearData.inputs[1:],
        },
      },
    },
  });
  {
    'sketch-edit': edit,
    'sketch-edit.app': edit { extra+: { linearMode: true } },
    'sketch-edit-masked': masked,
    'sketch-edit-masked.app': masked { extra+: { linearMode: true } },
  }
