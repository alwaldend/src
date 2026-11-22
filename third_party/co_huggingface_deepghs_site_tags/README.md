---
title: Danbooru tag vocabulary
description: Checksum-pinned CSV of Danbooru tag names and metadata
---

The [DeepGHS site_tags dataset](https://huggingface.co/datasets/deepghs/site_tags)
provides the Danbooru tag vocabulary, including tag categories, post counts, and
deprecation flags. The dataset card lists the license as CC BY 4.0.

Download the full 204 MB tag vocabulary with Bazel:

```sh
bazel build //third_party/co_huggingface_deepghs_site_tags:danbooru_tags
```
