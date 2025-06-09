<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_genrule_src"></a>

## al_genrule_src

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_genrule_src.bzl", "al_genrule_src")

al_genrule_src(<a href="#al_genrule_src-name">name</a>, <a href="#al_genrule_src-srcs">srcs</a>, <a href="#al_genrule_src-visibility">visibility</a>)
</pre>

Create a filegroup and a genrule generating a tar archive

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_genrule_src-name"></a>name |  genrule name   |  none |
| <a id="al_genrule_src-srcs"></a>srcs |  source labels   |  `[]` |
| <a id="al_genrule_src-visibility"></a>visibility |  <p align="center"> - </p>   |  `["//visibility:public"]` |


