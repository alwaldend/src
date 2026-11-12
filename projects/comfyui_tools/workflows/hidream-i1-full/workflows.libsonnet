// Native HiDream Full conditioning with shared latent, sequence and GIF graphs.
local g = import '../internal/graph.libsonnet';
local n = import '../internal/nodes.libsonnet';

function(models)
  local source = (import '../chroma/workflows.libsonnet')(
    { checkpoint: models.checkpoint, text_encoder: models.t5, vae: models.vae },
    { steps: 50, cfg: 5, sampler: 'uni_pc', scheduler: 'simple', min_padding: 0, filename_prefix: 'HiDreamI1Full' },
    sequences=true
  );
  local workflows = {
    [purpose]: g.variant(source[purpose], {
      removeNodes: [2, 5],
      removeLinks: [2, 3, 4],
      nodes: [n.node(2, 'QuadrupleCLIPLoader', [models.clip_l, models.clip_g, models.t5, models.llama], [0, 240], [400, 240], title='HiDream: four text encoders', order=2)],
      updates: {
        '4': { type: 'ModelSamplingSD3', properties: { 'Node name for S&R': 'ModelSamplingSD3' }, widgets_values: [3] },
        '16': { type: 'ReHiDreamPatcher', properties: { 'Node name for S&R': 'ReHiDreamPatcher' }, widgets_values: ['float64', true] },
      },
      links: [[2, 2, 0, 6, 0], [3, 2, 0, 7, 0]],
    })
    for purpose in std.objectFields(source)
  };
  workflows {
    'image-to-image-prompt-sequence': g.sequence(workflows['image-to-image'],
                                                 n.node(17, 'ImagegenPromptSequence', ['The referenced character stands beside a bench.\n\nThe referenced character sits on the bench.\n\nThe referenced character waves beside the bench.', 'An illustrated scene, consistent referenced subject and visual style.', 'line', 'img2img'], [0, 1400], [680, 380], title='Img2img prompts; same reference for every image', order=12),
                                                 6,
                                                 7,
                                                 11,
                                                 [18, 19, 20]),
  }
