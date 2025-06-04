---
title: Bazel rules
---

## Api

<!-- STARDOC START -->
<!-- Generated with Stardoc: http://skydoc.bazel.build -->

<a id="unpack_archives"></a>

## unpack_archives

<pre>
load("@com-alwaldend-src//starlark/bazel/rules:defs.bzl", "unpack_archives")

unpack_archives(<a href="#unpack_archives-name">name</a>, <a href="#unpack_archives-srcs">srcs</a>, <a href="#unpack_archives-out">out</a>)
</pre>

**ATTRIBUTES**

| Name                                  | Description                    | Type                                                                | Mandatory | Default |
| :------------------------------------ | :----------------------------- | :------------------------------------------------------------------ | :-------- | :------ |
| <a id="unpack_archives-name"></a>name | A unique name for this target. | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required  |         |
| <a id="unpack_archives-srcs"></a>srcs | -                              | <a href="https://bazel.build/concepts/labels">List of labels</a>    | required  |         |
| <a id="unpack_archives-out"></a>out   | -                              | String                                                              | optional  | `""`    |

<!-- STARDOC END -->
