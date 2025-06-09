<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_apply_patches"></a>

## al_apply_patches

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_apply_patches.bzl", "al_apply_patches")

al_apply_patches(<a href="#al_apply_patches-name">name</a>, <a href="#al_apply_patches-src">src</a>, <a href="#al_apply_patches-patches">patches</a>, <a href="#al_apply_patches-visibility">visibility</a>, <a href="#al_apply_patches-kwargs">**kwargs</a>)
</pre>

Create a genrule applying patches

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_apply_patches-name"></a>name |  genrule name   |  none |
| <a id="al_apply_patches-src"></a>src |  source archive label   |  none |
| <a id="al_apply_patches-patches"></a>patches |  patches label   |  none |
| <a id="al_apply_patches-visibility"></a>visibility |  visibility   |  `["//visibility:public"]` |
| <a id="al_apply_patches-kwargs"></a>kwargs |  other genrule kwargs   |  none |


