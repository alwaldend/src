---
title: Danbooru tag search
description: Local prompt-tag lookup in ComfyUI
---

Build the optional installation bundle and the workflow catalog.

```sh
bazel build //projects/comfyui_tools/nodes/danbooru_tag_search:installation \
  //projects/comfyui_tools/workflows:workflows
```

Extract `bazel-bin/projects/comfyui_tools/nodes/danbooru_tag_search/danbooru_tag_search.tar.gz`
into your ComfyUI directory and restart ComfyUI. The archive installs the search
node with its full tag vocabulary under `custom_nodes/danbooru_tag_search/`.
Import `bazel-bin/projects/comfyui_tools/workflows/danbooru/danbooru-tag-search.app.json`
for App mode or the matching `danbooru-tag-search.json` for Graph Mode.

Enter a tag or a fragment such as `blue hair` and run the App. Spaces match
underscores and the query is case-insensitive. **Sort by** defaults to relevance,
which places exact matches first, then prefixes and other substring matches.
Popularity ranks all matches by descending post count. Turn on **Prompt spaces**
to copy tags such as `blue hair` into NoobAI, or turn it off to retain CSV names
such as `blue_hair`. The bundled App and graph enable prompt spaces by default.
Use **Limit** to set the number of results and **Include deprecated** to search
retired tags. Copy the results or download the text file. Searches run locally with no model
or GPU and save their text results through ComfyUI's native Save Text node.

The bundled CSV is the unchanged Danbooru vocabulary from
[DeepGHS site_tags](https://huggingface.co/datasets/deepghs/site_tags) under
[CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). The
[dataset download target](https://forgejo.alwaldend.com/alwaldend/src/src/branch/master/third_party/co_huggingface_deepghs_site_tags/include.MODULE.bazel)
records its pinned revision and checksum.
