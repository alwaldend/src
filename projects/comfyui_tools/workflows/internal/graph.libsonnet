// Compose native graphs with each connection defined in one edge table.
local byId(nodes) = { [std.toString(node.id)]: node for node in nodes };
local get(nodes, id) = byId(nodes)[std.toString(id)];
{
  // An edge is [link ID, source node ID, source slot, target node ID, target slot].
  variant(base, patch)::
    local retained = [
      std.mergePatch(node, std.get(std.get(patch, 'updates', {}), std.toString(node.id), {}))
      for node in base.nodes
      if !std.member(std.get(patch, 'removeNodes', []), node.id)
    ];
    local nodes = retained + std.get(patch, 'nodes', []);
    local links = [edge for edge in base.links if !std.member(std.get(patch, 'removeLinks', []), edge[0])]
                  + std.get(patch, 'links', []);
    base + std.get(patch, 'metadata', {}) + {
      nodes: if std.objectHas(patch, 'order') then [get(nodes, id) for id in patch.order] else nodes,
      links: if std.objectHas(patch, 'linkOrder') then
        [std.filter(function(edge) edge[0] == id, links)[0] for id in patch.linkOrder]
      else links,
    },

  // Preserve model-specific encoders while centralizing editable sequence text.
  sequence(base, sequenceNode, positiveId, negativeId, saveId, linkIds, textInput='text')::
    local positive = get(base.nodes, positiveId);
    local negative = if negativeId == null then null else get(base.nodes, negativeId);
    local save = get(base.nodes, saveId);
    local sampler = [node for node in base.nodes if std.member(['KSampler', 'ClownsharKSampler_Beta', 'BasicGuider', 'CFGGuider', 'ImagegenChainedImageToImageGIF'], node.type)][0];
    local modelLink = [edge for edge in base.links if edge[3] == sampler.id && sampler.inputs[edge[4]].type == 'MODEL'][0];
    local progressLinkId = std.foldl(function(maximum, id) std.max(maximum, id), [edge[0] for edge in base.links] + linkIds, 0) + 1;
    local widget(name) = { name: name, type: 'STRING', widget: { name: name } };
    self.variant(base, {
      nodes: [sequenceNode {
        widgets_values: sequenceNode.widgets_values + [if negative == null then '' else negative.widgets_values[0]],
        size: [sequenceNode.size[0], sequenceNode.size[1] + 140],
      }],
      updates: {
        [std.toString(positiveId)]: { inputs: positive.inputs + [widget(textInput)] },
        [std.toString(saveId)]: { inputs: save.inputs + [widget('filename_prefix')] },
      } + (if negative == null then {} else {
             [std.toString(negativeId)]: { inputs: negative.inputs + [widget(textInput)], title: 'Encode shared negative prompt' },
           }),
      links: [
        [modelLink[0], modelLink[1], modelLink[2], sequenceNode.id, 0],
        [progressLinkId, sequenceNode.id, 3, sampler.id, modelLink[4]],
        [linkIds[0], sequenceNode.id, 0, positiveId, std.length(positive.inputs)],
        [linkIds[1], sequenceNode.id, 1, saveId, std.length(save.inputs)],
      ] + (if negative == null then [] else [[linkIds[2], sequenceNode.id, 2, negativeId, std.length(negative.inputs)]]),
      removeLinks: [modelLink[0]],
      metadata: { last_node_id: sequenceNode.id, last_link_id: progressLinkId },
    }),

  render(spec)::
    local nodes = byId(spec.nodes);
    local ids = [node.id for node in spec.nodes];
    local linkIds = [edge[0] for edge in spec.links];
    local expanded = [
      local source = nodes[std.toString(edge[1])].outputs[edge[2]];
      local target = nodes[std.toString(edge[3])].inputs[edge[4]];
      assert source.type == target.type : 'Port type mismatch in link ' + std.toString(edge[0]);
      edge + [source.type]
      for edge in spec.links
    ];
    assert std.length(std.set(ids)) == std.length(ids) : 'Duplicate node IDs';
    assert std.length(std.set(linkIds)) == std.length(linkIds) : 'Duplicate link IDs';
    assert std.all([
      std.length([edge for edge in spec.links if edge[3] == node.id && edge[4] == slot]) <= 1
      for node in spec.nodes
      for slot in std.range(0, std.length(node.inputs) - 1)
    ])
           : 'Multiple connections to the same input';
    spec {
      links: expanded,
      nodes: [node {
        inputs: [
          port { link: local matches = [edge[0] for edge in spec.links if edge[3] == node.id && edge[4] == slot];
                      if std.length(matches) == 0 then null else matches[0] }
          for slot in std.range(0, std.length(node.inputs) - 1)
          for port in [node.inputs[slot]]
        ],
        outputs: [
          port { links: [edge[0] for edge in spec.links if edge[1] == node.id && edge[2] == slot] }
          for slot in std.range(0, std.length(node.outputs) - 1)
          for port in [node.outputs[slot]]
        ],
      } for node in spec.nodes],
    },
}
