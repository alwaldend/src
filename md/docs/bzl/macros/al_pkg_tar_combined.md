<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_pkg_tar_combined"></a>

## al_pkg_tar_combined

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_pkg_tar_combined.bzl", "al_pkg_tar_combined")

al_pkg_tar_combined(<a href="#al_pkg_tar_combined-name">name</a>, <a href="#al_pkg_tar_combined-srcs">srcs</a>, <a href="#al_pkg_tar_combined-strip_components">strip_components</a>, <a href="#al_pkg_tar_combined-kwargs">**kwargs</a>)
</pre>

Create a genrule combining several tars into one

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_pkg_tar_combined-name"></a>name |  genrule name   |  none |
| <a id="al_pkg_tar_combined-srcs"></a>srcs |  dicts tar archives ({"label": "tar label", "dir": "target dir"})   |  `[]` |
| <a id="al_pkg_tar_combined-strip_components"></a>strip_components |  value of --stip-components   |  `2` |
| <a id="al_pkg_tar_combined-kwargs"></a>kwargs |  other genrule kwargs   |  none |


