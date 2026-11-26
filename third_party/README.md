---
title: Third party
description: Third party tree
weight: 7
cascade:
  - categories:
      - third_party
---

This tree contains vendored and externally sourced code. All tracked content
follows the repository's public-source policy and retains its upstream license
and provenance requirements.

## Requirements

- Bazel targets MUST use repository-internal visibility unless an owning
  toolchain contract requires otherwise.
- Third-party content MUST NOT be represented or published as first-party
  source or artifacts.
- Third-party content MAY be used in builds.
- An explicit repackaging or mirroring target MAY publish an upstream artifact
  under its owner and license policy.

## Blender assets

- [Pipkin Pippa mesh and diffuse texture](ai_darlin_pippa/README.md)
- [Human base meshes and demo scenes](org_blender_demo_assets/README.md)
- [Snow character rig](org_blender_studio_snow/README.md)
- [CMU breakdance clips](com_github_shriinivas_cmubvh/README.md)
- [Mixamo dance clips from the Pippa album](co_huggingface_tanish434_mixamo/README.md)
- [Kenney furniture](nl_kenney_furniture_kit/README.md)

Build an individual target to fetch its asset. The example below locates
the downloaded file under `external/` in the Bazel output base.

```sh
bazel build //third_party/com_github_shriinivas_cmubvh:helicopter
bazel info output_base
bazel cquery --output=files //third_party/com_github_shriinivas_cmubvh:helicopter
```

## Publish vendored helm

This is a state-changing operator example. An agent needs explicit publication
authority and must invoke the target through `bazel`.

```sh
bazel run //third_party:publish_helm.io_goharbor_helm_harbor
```
