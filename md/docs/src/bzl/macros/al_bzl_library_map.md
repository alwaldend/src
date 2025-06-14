<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_bzl_library_map"></a>

## al_bzl_library_map

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_bzl_library_map.bzl", "al_bzl_library_map")

al_bzl_library_map(<a href="#al_bzl_library_map-name">name</a>, <a href="#al_bzl_library_map-libs">libs</a>, <a href="#al_bzl_library_map-common_deps">common_deps</a>, <a href="#al_bzl_library_map-visibility">visibility</a>)
</pre>

Create al_bzl_library targets from a map

Targets:
- ${name}-stardoc: all stardoc markdown
- ${name}: all bzl_library targets


**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_bzl_library_map-name"></a>name |  al_md_data name   |  none |
| <a id="al_bzl_library_map-libs"></a>libs |  map of al_bzl_library, keys are names, values are kwargs for al_bzl_library   |  none |
| <a id="al_bzl_library_map-common_deps"></a>common_deps |  list of common blz_library deps   |  `[]` |
| <a id="al_bzl_library_map-visibility"></a>visibility |  visibility   |  `["//visibility:public"]` |


