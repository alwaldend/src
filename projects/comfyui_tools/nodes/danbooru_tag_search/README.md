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
underscores and the query is case-insensitive. Exact tags appear first, then
prefix matches and other substring matches, with more frequently used tags
first in each group. Use **Limit** to set how many results to return and
**Include deprecated** to search retired tags. Copy the comma-separated names
from the result or download the text file. Searches run locally with no model
or GPU and save their text results through ComfyUI's native Save Text node.

The bundled CSV is the unchanged Danbooru vocabulary from
[DeepGHS site_tags](https://huggingface.co/datasets/deepghs/site_tags) under
[CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). The
[dataset download target](https://forgejo.alwaldend.com/alwaldend/src/src/branch/master/third_party/co_huggingface_deepghs_site_tags/include.MODULE.bazel)
records its pinned revision and checksum.
