local g = import '../internal/graph.libsonnet';
local n = import '../internal/nodes.libsonnet';

{
  'tag-search': {
    config: {},
    extra: { ds: { offset: [60, 60], scale: 0.8 } },
    groups: [],
    last_link_id: 1,
    last_node_id: 2,
    version: 0.4,
    nodes: [
      n.node(1, 'ImagegenDanbooruTagSearch', ['blue hair', 20, false, true, 'relevance'], [0, 0], [330, 212]),
      n.node(2, 'SaveText', ['Danbooru-tag-search', 'txt'], [400, 0], [520, 240], title='Matching tags'),
    ],
    links: [[1, 1, 0, 2, 0]],
  },
  'tag-search.app': g.app(self['tag-search'], [
    [1, 'query', { description: 'Enter a Danbooru tag or part of a tag. Spaces match underscores.' }],
    [1, 'limit'],
    [1, 'include_deprecated'],
    [1, 'prompt_spaces', { description: 'Use spaces in multiword tags for NoobAI prompts.' }],
    [1, 'sort_by', { description: 'Choose relevance or descending Danbooru post count.' }],
  ], 2),
}
