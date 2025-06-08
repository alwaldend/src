<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_md_data"></a>

## al_md_data

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_md_data.bzl", "al_md_data")

al_md_data(<a href="#al_md_data-name">name</a>, <a href="#al_md_data-srcs">srcs</a>, <a href="#al_md_data-visibility">visibility</a>, <a href="#al_md_data-kwargs">**kwargs</a>)
</pre>

Markdown data backed by a filegroup

Targets:
    ${name}: filegroup


**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_md_data-name"></a>name |  filegroup name   |  none |
| <a id="al_md_data-srcs"></a>srcs |  markdown files   |  none |
| <a id="al_md_data-visibility"></a>visibility |  visibility   |  `["//visibility:public"]` |
| <a id="al_md_data-kwargs"></a>kwargs |  kwargs   |  none |


