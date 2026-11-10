// Native port contracts for the nodes used by these workflow families.
local ports = {
  CFGGuider: { inputs: [{ name: 'model', type: 'MODEL' }, { name: 'positive', type: 'CONDITIONING' }, { name: 'negative', type: 'CONDITIONING' }], outputs: [{ name: 'GUIDER', type: 'GUIDER' }] },
  CFGNorm: { inputs: [{ name: 'model', type: 'MODEL' }], outputs: [{ name: 'MODEL', type: 'MODEL' }] },
  CLIPLoader: { inputs: [], outputs: [{ name: 'CLIP', type: 'CLIP' }] },
  CLIPTextEncode: { inputs: [{ name: 'clip', type: 'CLIP' }], outputs: [{ name: 'CONDITIONING', type: 'CONDITIONING' }] },
  CheckpointLoaderSimple: { inputs: [], outputs: [{ name: 'MODEL', type: 'MODEL' }, { name: 'CLIP', type: 'CLIP' }, { name: 'VAE', type: 'VAE' }] },
  ClownGuide_Style_Beta: { inputs: [{ name: 'guide', type: 'LATENT' }], outputs: [{ name: 'GUIDES', type: 'GUIDES' }] },
  ClownsharKSampler_Beta: { inputs: [{ name: 'model', type: 'MODEL' }, { name: 'positive', type: 'CONDITIONING' }, { name: 'negative', type: 'CONDITIONING' }, { name: 'latent_image', type: 'LATENT' }, { name: 'guides', type: 'GUIDES' }], outputs: [{ name: 'LATENT', type: 'LATENT' }, { name: 'LATENT', type: 'LATENT' }, { name: 'OPTIONS', type: 'OPTIONS' }] },
  EmptyFlux2LatentImage: { inputs: [{ name: 'width', type: 'INT', widget: { name: 'width' } }, { name: 'height', type: 'INT', widget: { name: 'height' } }], outputs: [{ name: 'LATENT', type: 'LATENT' }] },
  EmptyLatentImage: { inputs: [], outputs: [{ name: 'LATENT', type: 'LATENT' }] },
  EmptySD3LatentImage: { inputs: [], outputs: [{ name: 'LATENT', type: 'LATENT' }] },
  Flux2Scheduler: { inputs: [{ name: 'width', type: 'INT', widget: { name: 'width' } }, { name: 'height', type: 'INT', widget: { name: 'height' } }], outputs: [{ name: 'SIGMAS', type: 'SIGMAS' }] },
  FluxKontextImageScale: { inputs: [{ name: 'image', type: 'IMAGE' }], outputs: [{ name: 'IMAGE', type: 'IMAGE' }] },
  FluxKontextMultiReferenceLatentMethod: { inputs: [{ name: 'conditioning', type: 'CONDITIONING' }], outputs: [{ name: 'CONDITIONING', type: 'CONDITIONING' }] },
  GetImageSize: { inputs: [{ name: 'image', type: 'IMAGE' }], outputs: [{ name: 'INT', type: 'INT' }, { name: 'INT', type: 'INT' }, { name: 'INT', type: 'INT' }] },
  ImageScale: { inputs: [{ name: 'image', type: 'IMAGE' }], outputs: [{ name: 'IMAGE', type: 'IMAGE' }] },
  ImageScaleToTotalPixels: { inputs: [{ name: 'image', type: 'IMAGE' }], outputs: [{ name: 'IMAGE', type: 'IMAGE' }] },
  ImagegenPromptSequence: { inputs: [], outputs: [{ name: 'prompt', type: 'STRING' }, { name: 'filename_prefix', type: 'STRING' }] },
  KSampler: { inputs: [{ name: 'model', type: 'MODEL' }, { name: 'positive', type: 'CONDITIONING' }, { name: 'negative', type: 'CONDITIONING' }, { name: 'latent_image', type: 'LATENT' }], outputs: [{ name: 'LATENT', type: 'LATENT' }] },
  KSamplerSelect: { inputs: [], outputs: [{ name: 'SAMPLER', type: 'SAMPLER' }] },
  LoadImage: { inputs: [], outputs: [{ name: 'IMAGE', type: 'IMAGE' }, { name: 'MASK', type: 'MASK' }] },
  ModelSamplingAuraFlow: { inputs: [{ name: 'model', type: 'MODEL' }], outputs: [{ name: 'MODEL', type: 'MODEL' }] },
  RandomNoise: { inputs: [], outputs: [{ name: 'NOISE', type: 'NOISE' }] },
  ReChromaPatcher: { inputs: [{ name: 'model', type: 'MODEL' }], outputs: [{ name: 'model', type: 'MODEL' }] },
  ReSDPatcher: { inputs: [{ name: 'model', type: 'MODEL' }], outputs: [{ name: 'model', type: 'MODEL' }] },
  ReferenceLatent: { inputs: [{ name: 'conditioning', type: 'CONDITIONING' }, { name: 'latent', type: 'LATENT' }], outputs: [{ name: 'CONDITIONING', type: 'CONDITIONING' }] },
  SamplerCustomAdvanced: { inputs: [{ name: 'noise', type: 'NOISE' }, { name: 'guider', type: 'GUIDER' }, { name: 'sampler', type: 'SAMPLER' }, { name: 'sigmas', type: 'SIGMAS' }, { name: 'latent_image', type: 'LATENT' }], outputs: [{ name: 'LATENT', type: 'LATENT' }, { name: 'LATENT', type: 'LATENT' }] },
  SaveImage: { inputs: [{ name: 'images', type: 'IMAGE' }], outputs: [] },
  T5TokenizerOptions: { inputs: [{ name: 'clip', type: 'CLIP' }], outputs: [{ name: 'CLIP', type: 'CLIP' }] },
  TextEncodeQwenImageEditPlus: { inputs: [{ name: 'clip', type: 'CLIP' }, { name: 'vae', type: 'VAE' }, { name: 'image1', type: 'IMAGE' }, { name: 'image2', type: 'IMAGE' }, { name: 'image3', type: 'IMAGE' }], outputs: [{ name: 'CONDITIONING', type: 'CONDITIONING' }] },
  UNETLoader: { inputs: [], outputs: [{ name: 'MODEL', type: 'MODEL' }] },
  VAEDecode: { inputs: [{ name: 'samples', type: 'LATENT' }, { name: 'vae', type: 'VAE' }], outputs: [{ name: 'IMAGE', type: 'IMAGE' }] },
  VAEEncode: { inputs: [{ name: 'pixels', type: 'IMAGE' }, { name: 'vae', type: 'VAE' }], outputs: [{ name: 'LATENT', type: 'LATENT' }] },
  VAELoader: { inputs: [], outputs: [{ name: 'VAE', type: 'VAE' }] },
};

{
  // Link IDs are derived by graph.render, not repeated in node definitions.
  node(id, kind, values, pos, size, indexed=true, title=null, order=id - 1)::
    assert std.objectHas(ports, kind) : 'Unknown workflow node type: ' + kind;
    {
      id: id,
      type: kind,
      pos: pos,
      size: size,
      flags: {},
      order: order,
      mode: 0,
      properties: { 'Node name for S&R': kind },
      widgets_values: values,
      inputs: ports[kind].inputs,
      outputs: [
        port + (if indexed then { slot_index: index } else {})
        for index in std.range(0, std.length(ports[kind].outputs) - 1)
        for port in [ports[kind].outputs[index]]
      ],
    } + (if title == null then {} else { title: title }),
}
