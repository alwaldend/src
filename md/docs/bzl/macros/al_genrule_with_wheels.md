<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_genrule_with_wheels"></a>

## al_genrule_with_wheels

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_genrule_with_wheels.bzl", "al_genrule_with_wheels")

al_genrule_with_wheels(<a href="#al_genrule_with_wheels-name">name</a>, <a href="#al_genrule_with_wheels-wheels">wheels</a>, <a href="#al_genrule_with_wheels-srcs">srcs</a>, <a href="#al_genrule_with_wheels-cmd">cmd</a>, <a href="#al_genrule_with_wheels-kwargs">**kwargs</a>)
</pre>

Regular genrule with wheels added to ${PYTHONPATH}

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_genrule_with_wheels-name"></a>name |  genrule name   |  none |
| <a id="al_genrule_with_wheels-wheels"></a>wheels |  list of wheel labels   |  none |
| <a id="al_genrule_with_wheels-srcs"></a>srcs |  srcs for the genrule   |  `[]` |
| <a id="al_genrule_with_wheels-cmd"></a>cmd |  genrule cmd   |  `[]` |
| <a id="al_genrule_with_wheels-kwargs"></a>kwargs |  other genrule kwargs   |  none |


