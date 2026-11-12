// Limit native GIF sampling to an editable, feathered image-space rectangle.
local g = import 'graph.libsonnet';
local n = import 'nodes.libsonnet';

{
  masked(base, imageId)::
    local sampler = [node for node in base.nodes if std.member(['ImagegenChainedImageToImageGIF', 'ImagegenQwenChainedImageToImageGIF'], node.type)][0];
    local latent = [edge for edge in base.links if edge[3] == sampler.id && edge[4] == 3][0];
    local nodeId = std.foldl(std.max, [node.id for node in base.nodes], 0) + 1;
    local linkId = std.foldl(std.max, [edge[0] for edge in base.links], 0) + 1;
    g.variant(base, {
      nodes: [
        n.node(nodeId, 'SolidMask', [0, 1024, 1024], [0, 1850], [300, 180], title='Unchanged image area'),
        n.node(nodeId + 1, 'SolidMask', [1, 512, 256], [0, 2100], [300, 180], title='Edit region size'),
        n.node(nodeId + 2, 'FeatherMask', [32, 32, 32, 32], [350, 2100], [300, 180], title='Soft edit boundary'),
        n.node(nodeId + 3, 'MaskComposite', [256, 256, 'add'], [700, 1950], [300, 180], title='Edit region position'),
        n.node(nodeId + 4, 'SetLatentNoiseMask', [], [1050, 1950], [300, 140]),
        n.node(nodeId + 5, 'GetImageSize', [], [0, 1600], [300, 100]),
      ],
      removeLinks: [latent[0]],
      links: [
        [latent[0], latent[1], latent[2], nodeId + 4, 0],
        [linkId, nodeId + 4, 0, sampler.id, 3],
        [linkId + 1, nodeId, 0, nodeId + 3, 0],
        [linkId + 2, nodeId + 1, 0, nodeId + 2, 0],
        [linkId + 3, nodeId + 2, 0, nodeId + 3, 1],
        [linkId + 4, nodeId + 3, 0, nodeId + 4, 1],
        [linkId + 5, imageId, 0, nodeId + 5, 0],
        [linkId + 6, nodeId + 5, 0, nodeId, 0],
        [linkId + 7, nodeId + 5, 1, nodeId, 1],
      ],
      metadata: { last_node_id: nodeId + 5, last_link_id: linkId + 7 },
    }),
}
