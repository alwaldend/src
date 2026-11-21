local g = import '../internal/graph.libsonnet';
local n = import '../internal/nodes.libsonnet';

function(models) {
  'text-to-video': {
    config: {},
    extra: { ds: { offset: [80, 60], scale: 0.55 } },
    groups: [],
    last_link_id: 17,
    last_node_id: 14,
    version: 0.4,
    nodes: [
      n.node(1, 'UNETLoader', [models.checkpoint, 'default'], [0, 0], [400, 120], title='H3 Base FL2VA'),
      n.node(2, 'CLIPLoader', [models.text_encoder, 'minimax', 'default'], [0, 160], [400, 140]),
      n.node(3, 'VAELoader', [models.vae], [0, 340], [400, 100], title='Video VAE'),
      n.node(4, 'VAELoader', [models.audio_vae], [0, 480], [400, 100], title='Stereo audio VAE'),
      n.node(5, 'MiniMaxH3ImageToVideo', ['A small sailboat glides across a calm lake at sunrise. The camera follows from shore. Gentle water and distant birds are audible.', 864, 480, 107], [460, 0], [460, 360], title='Video prompt and dimensions'),
      n.node(6, 'BasicGuider', [], [980, 0], [280, 100]),
      n.node(7, 'BasicScheduler', ['simple', 20, 1], [980, 140], [280, 140]),
      n.node(8, 'KSamplerSelect', ['res_multistep'], [980, 320], [280, 90]),
      n.node(9, 'RandomNoise', [42, 'randomize'], [980, 450], [280, 100]),
      n.node(10, 'SamplerCustomAdvanced', [], [1320, 0], [300, 160]),
      n.node(11, 'VAEDecode', [], [1680, 0], [260, 100], title='Decode video'),
      n.node(12, 'VAEDecodeAudio', [], [1680, 150], [260, 100], title='Decode stereo audio'),
      n.node(13, 'CreateVideo', [24, 8, 'sRGB', 'auto', 'none'], [2000, 0], [280, 200]),
      n.node(14, 'SaveVideo', ['video/minimax-h3-text-to-video', 'mp4', 'h264'], [2340, 0], [340, 220]),
    ],
    links: [
      [1, 2, 0, 5, 0],
      [2, 3, 0, 5, 1],
      [3, 1, 0, 6, 0],
      [4, 5, 0, 6, 1],
      [5, 1, 0, 7, 0],
      [6, 9, 0, 10, 0],
      [7, 6, 0, 10, 1],
      [8, 8, 0, 10, 2],
      [9, 7, 0, 10, 3],
      [10, 5, 1, 10, 4],
      [11, 10, 0, 11, 0],
      [12, 3, 0, 11, 1],
      [13, 10, 0, 12, 0],
      [14, 4, 0, 12, 1],
      [15, 11, 0, 13, 0],
      [16, 12, 0, 13, 1],
      [17, 13, 0, 14, 0],
    ],
  },
  'first-last-frame': g.variant(self['text-to-video'], {
    nodes: [
      n.node(15, 'LoadImage', ['first-frame.png', 'image'], [0, 630], [400, 300], title='Upload the first frame'),
      n.node(16, 'LoadImage', ['last-frame.png', 'image'], [460, 630], [400, 300], title='Upload the last frame'),
    ],
    updates: {
      '5': { widgets_values: ['Animate one continuous transition from the opening image to the closing image. Preserve the subject and scene. Use subtle ambient sound that matches the scene.', 864, 480, 107] },
      '14': { widgets_values: ['video/minimax-h3-first-last-frame', 'mp4', 'h264'] },
    },
    links: [
      [18, 15, 0, 5, 2],
      [19, 16, 0, 5, 3],
    ],
    metadata: { last_node_id: 16, last_link_id: 19 },
  }),
}
