---
title: Bazel providers
---

<!-- STARDOC START -->
<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="AlTransitiveSources"></a>

## AlTransitiveSources

<pre>
load("@com-alwaldend-src//starlark/bazel/providers:defs.bzl", "AlTransitiveSources")

AlTransitiveSources(<a href="#AlTransitiveSources-transitive_sources">transitive_sources</a>)
</pre>

Provide transitive sources

- https://github.com/bazelbuild/examples/blob/main/rules/depsets/foo.bzl
- https://stackoverflow.com/a/57699683

**FIELDS**

| Name  | Description |
| :------------- | :------------- |
| <a id="AlTransitiveSources-transitive_sources"></a>transitive_sources |  -    |


<a id="al_transitive_sources"></a>

## al_transitive_sources

<pre>
load("@com-alwaldend-src//starlark/bazel/providers:defs.bzl", "al_transitive_sources")

al_transitive_sources(<a href="#al_transitive_sources-srcs">srcs</a>, <a href="#al_transitive_sources-deps">deps</a>)
</pre>

Obtain the source files for a target and its transitive dependencies.

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_transitive_sources-srcs"></a>srcs |  a list of source files   |  none |
| <a id="al_transitive_sources-deps"></a>deps |  a list of targets that are direct dependencies   |  none |

**RETURNS**

a collection of the transitive sources



<!-- STARDOC END -->
