// Approved SDXL checkpoint using the shared generation/reference graphs.
local g = import '../internal/graph.libsonnet';
local n = import '../internal/nodes.libsonnet';

function(models)
  local settings = { name: 'NoobAI XL 1.1', prefix: 'NoobAI', steps: 30, cfg: 5.5, sampler: 'euler_ancestral', scheduler: 'normal', prompt: 'masterpiece, best quality, newest, highres, watercolor illustration, robot, blue body, red chest plate, reading, open book, park bench, snowy park, clear silhouettes', negative: 'worst quality, low quality, lowres, blurry, signature, username, logo, text, watermark', sequence_context: 'masterpiece, best quality, newest, highres, watercolor illustration, consistent small blue robot with a red chest plate, snowy park, one scene', reference_context: 'masterpiece, best quality, newest, highres, illustration, consistent referenced characters and visual style, one scene' };
  local workflows = (import '../internal/sdxl.libsonnet')(models, settings);
  workflows {
    'image-to-image-prompt-sequence': g.sequence(workflows['image-to-image'],
                                                 n.node(10, 'ImagegenPromptSequence', ['The referenced character stands beside a bench.\n\nThe referenced character sits on the bench.\n\nThe referenced character waves beside the bench.', settings.reference_context, 'line', 'img2img'], [380, 850], [680, 380], title='Img2img prompts; same reference for every image', order=9),
                                                 2,
                                                 3,
                                                 7,
                                                 [13, 14, 15]),
  }
