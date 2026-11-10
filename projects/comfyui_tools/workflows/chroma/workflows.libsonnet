// Model-specific native workflows; graph serialization is shared.
local g = import '../internal/graph.libsonnet';
local n = import '../internal/nodes.libsonnet';

function(models) {
  'text-to-image': {
    config: {},
    extra: { ds: { offset: [60, 60], scale: 0.65 } },
    groups: [],
    last_link_id: 16,
    last_node_id: 15,
    version: 0.4,
    nodes: [
      n.node(1, 'UNETLoader', ['Chroma1-HD-Flash.safetensors', 'default'], [0, 0], [315, 220], order=1),
      n.node(2, 'CLIPLoader', ['t5xxl_fp16.safetensors', 'chroma', 'default'], [0, 240], [315, 220], order=2),
      n.node(3, 'VAELoader', ['ae.safetensors'], [0, 480], [315, 220], order=3),
      n.node(4, 'ModelSamplingAuraFlow', [1.0], [350, 0], [315, 220], order=4),
      n.node(5, 'T5TokenizerOptions', [1, 0], [350, 240], [315, 220], order=5),
      n.node(6, 'CLIPTextEncode', ['A comic panel in a hand-painted watercolor style: a traveler in a red coat walking through a quiet mountain village, wide shot, clear silhouettes.'], [700, 200], [315, 220], order=6),
      n.node(7, 'CLIPTextEncode', [''], [700, 460], [315, 220], order=7),
      n.node(10, 'VAEDecode', [], [1400, 0], [315, 220], order=10),
      n.node(11, 'SaveImage', ['ChromaFlash'], [1750, 0], [315, 220], order=11),
      n.node(8, 'EmptySD3LatentImage', [1024, 1024, 1], [700, 720], [315, 220], order=8),
      n.node(9, 'KSampler', [0, 'randomize', 10, 1.0, 'heun', 'beta', 1.0], [1050, 0], [315, 220], order=9),
    ],
    links: [
      [1, 1, 0, 4, 0],
      [2, 2, 0, 5, 0],
      [3, 5, 0, 6, 0],
      [4, 5, 0, 7, 0],
      [5, 4, 0, 9, 0],
      [6, 6, 0, 9, 1],
      [7, 7, 0, 9, 2],
      [8, 8, 0, 9, 3],
      [9, 9, 0, 10, 0],
      [10, 3, 0, 10, 1],
      [11, 10, 0, 11, 0],
    ],
  },
  'image-to-image': g.variant(self['text-to-image'], {
    removeNodes: [8, 9],
    nodes: [
      n.node(8, 'VAEEncode', [], [700, 760], [315, 220], order=8),
      n.node(12, 'LoadImage', ['reference.png', 'image'], [0, 760], [315, 220], order=12),
      n.node(13, 'ImageScale', ['lanczos', 1024, 1024, 'center'], [350, 760], [315, 220], order=13),
      n.node(9, 'KSampler', [0, 'randomize', 10, 1.0, 'heun', 'beta', 0.7], [1050, 0], [315, 220], order=9),
    ],
    order: [1, 2, 3, 4, 5, 6, 7, 10, 11, 8, 9, 12, 13],
    links: [
      [12, 12, 0, 13, 0],
      [13, 13, 0, 8, 0],
      [14, 3, 0, 8, 1],
    ],
  }),
  'style-reference': g.variant(self['text-to-image'], {
    removeNodes: [9],
    nodes: [
      n.node(9, 'ClownsharKSampler_Beta', [0.0, 'multistep/res_2m', 'beta57', 10, -1, 1.0, 1.0, 0, 'randomize', 'standard', true], [1050, 0], [315, 220], order=9),
      n.node(12, 'LoadImage', ['reference.png', 'image'], [0, 1040], [315, 220], order=12),
      n.node(13, 'ImageScale', ['lanczos', 1024, 1024, 'center'], [350, 1040], [315, 220], order=13),
      n.node(14, 'VAEEncode', [], [700, 1040], [315, 220], order=14),
      n.node(15, 'ClownGuide_Style_Beta', ['positive', 'WCT', 1.0, 1.0, 'constant', 0, -1, false], [1050, 1040], [315, 220], order=15),
      n.node(16, 'ReChromaPatcher', ['float32', true], [700, 0], [315, 140], order=16),
    ],
    removeLinks: [5],
    links: [
      [5, 16, 0, 9, 0],
      [12, 12, 0, 13, 0],
      [13, 13, 0, 14, 0],
      [14, 3, 0, 14, 1],
      [15, 14, 0, 15, 0],
      [16, 15, 0, 9, 4],
      [17, 4, 0, 16, 0],
    ],
    metadata: { last_link_id: 17, last_node_id: 16 },
    linkOrder: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17],
  }),
}
