// Native FLUX.1 USO/UMO references share loaders, conditioning and sampling.
local g = import 'graph.libsonnet';
local n = import 'nodes.libsonnet';

{
  reference(models, subjects, styles, steps, guidance, prefix, referenceSize=1024)::
    assert subjects >= 0 && styles >= 0 && subjects + styles > 0;
    local subjectIds = [20 + index * 4 for index in std.range(0, subjects - 1)];
    local styleIds = [52 + index * 3 for index in std.range(0, styles - 1)];
    local lastConditioning = if subjects == 0 then 4 else subjectIds[subjects - 1] + 3;
    local lastModel = if styles == 0 then 2 else styleIds[styles - 1] + 2;
    local scene = if subjects > 1 then 'All referenced characters read books together beside a tree' else if subjects == 1 then 'The referenced character reads a book beside a tree' else 'A small blue robot reads a book beside a tree';
    {
      config: {},
      extra: { ds: { offset: [60, 60], scale: 0.55 } },
      groups: [],
      last_node_id: std.max(10, std.max(lastConditioning, lastModel)),
      last_link_id: if styles > 0 then 199 + styles * 5 else if subjects > 0 then 99 + subjects * 5 else 12,
      version: 0.4,
      nodes: [
        n.node(1, 'CheckpointLoaderSimple', [models.checkpoint], [0, 0], [340, 140], title='FLUX.1 dev FP8 with local encoders and VAE'),
        n.node(2, 'LoraLoaderModelOnly', [models.lora, 1.0], [390, 0], [340, 140], title=prefix + ' adapter'),
        n.node(3, 'CLIPTextEncode', [scene + ', a wide comic panel with clear silhouettes.'], [790, 0], [420, 260], title='Scene and visual style'),
        n.node(4, 'FluxGuidance', [guidance], [1260, 0], [300, 120]),
        n.node(5, 'ConditioningZeroOut', [], [1260, 190], [300, 90], title='CFG 1: unused negative branch'),
        n.node(6, 'FluxKontextMultiReferenceLatentMethod', ['uxo/uno'], [1620, 0], [320, 120], title='Native USO/UNO reference positioning'),
        n.node(7, 'EmptySD3LatentImage', [1024, 1024, 1], [1260, 370], [300, 150], title='Output size; fresh noise'),
        n.node(8, 'KSampler', [123456789, 'randomize', steps, 1.0, 'euler', 'simple', 1.0], [2010, 0], [320, 300]),
        n.node(9, 'VAEDecode', [], [2400, 0], [300, 120]),
        n.node(10, 'SaveImage', [prefix], [2760, 0], [340, 300]),
      ] + std.flattenArrays([
        [
          n.node(id, 'LoadImage', ['reference-' + (index + 1) + '.png', 'image'], [0, 600 + index * 410], [340, 300], title='Subject reference ' + (index + 1)),
          n.node(id + 1, 'ImageScaleToMaxDimension', ['area', referenceSize], [390, 600 + index * 410], [300, 120], title='Reference size; preserve aspect'),
          n.node(id + 2, 'VAEEncode', [], [790, 600 + index * 410], [300, 120]),
          n.node(id + 3, 'ReferenceLatent', [], [1260, 600 + index * 410], [300, 120]),
        ]
        for index in std.range(0, subjects - 1)
        for id in [subjectIds[index]]
      ]) + (if styles == 0 then [] else [
              n.node(50, 'ModelPatchLoader', [models.projector], [0, 600 + subjects * 410], [340, 120], title='USO style projector'),
              n.node(51, 'CLIPVisionLoader', [models.clip_vision], [390, 600 + subjects * 410], [340, 120], title='SigLIP style encoder'),
            ]) + std.flattenArrays([
        [
          n.node(id, 'LoadImage', ['style-' + (index + 1) + '.png', 'image'], [0, 820 + (subjects + index) * 410], [340, 300], title='Style reference ' + (index + 1)),
          n.node(id + 1, 'CLIPVisionEncode', ['center'], [790, 820 + (subjects + index) * 410], [300, 120]),
          n.node(id + 2, 'USOStyleReference', [], [1620, 820 + (subjects + index) * 410], [320, 150]),
        ]
        for index in std.range(0, styles - 1)
        for id in [styleIds[index]]
      ]),
      links: [
        [1, 1, 0, 2, 0],
        [2, 1, 1, 3, 0],
        [3, 3, 0, 4, 0],
        [4, 3, 0, 5, 0],
        [5, lastConditioning, 0, 6, 0],
        [6, 6, 0, 8, 1],
        [7, 5, 0, 8, 2],
        [8, 7, 0, 8, 3],
        [9, lastModel, 0, 8, 0],
        [10, 8, 0, 9, 0],
        [11, 1, 2, 9, 1],
        [12, 9, 0, 10, 0],
      ] + std.flattenArrays([
        [
          [link, id, 0, id + 1, 0],
          [link + 1, id + 1, 0, id + 2, 0],
          [link + 2, 1, 2, id + 2, 1],
          [link + 3, if index == 0 then 4 else subjectIds[index - 1] + 3, 0, id + 3, 0],
          [link + 4, id + 2, 0, id + 3, 1],
        ]
        for index in std.range(0, subjects - 1)
        for id in [subjectIds[index]]
        for link in [100 + index * 5]
      ]) + std.flattenArrays([
        [
          [link, 51, 0, id + 1, 0],
          [link + 1, id, 0, id + 1, 1],
          [link + 2, if index == 0 then 2 else styleIds[index - 1] + 2, 0, id + 2, 0],
          [link + 3, 50, 0, id + 2, 1],
          [link + 4, id + 1, 0, id + 2, 2],
        ]
        for index in std.range(0, styles - 1)
        for id in [styleIds[index]]
        for link in [200 + index * 5]
      ]),
    },

  sequence(base, prefix)::
    local subjects = std.length([node for node in base.nodes if node.type == 'ReferenceLatent']);
    local prompts = if subjects > 1 then 'All referenced characters enter a quiet park.\n\nAll referenced characters read books together on a bench.\n\nAll referenced characters wave beside a tree.' else if subjects == 1 then 'The referenced character enters a quiet park.\n\nThe referenced character reads a book on a bench.\n\nThe referenced character waves beside a tree.' else 'A small blue robot enters a quiet park.\n\nThe small blue robot reads a book on a bench.\n\nThe small blue robot waves beside a tree.';
    g.sequence(base,
               n.node(80, 'ImagegenPromptSequence', [prompts, 'An illustrated comic panel with consistent subjects and visual style.', 'line', prefix], [0, 2400], [680, 380], title='Panel prompts and shared context'),
               3,
               null,
               10,
               [300, 301, 302]),
}
