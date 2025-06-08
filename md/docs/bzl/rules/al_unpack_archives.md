<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_unpack_archives"></a>

## al_unpack_archives

<pre>
load("@com-alwaldend-git-src//bzl/rules:al_unpack_archives.bzl", "al_unpack_archives")

al_unpack_archives(<a href="#al_unpack_archives-name">name</a>, <a href="#al_unpack_archives-srcs">srcs</a>, <a href="#al_unpack_archives-out">out</a>)
</pre>

Unpack several archives using tar into a directory

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="al_unpack_archives-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="al_unpack_archives-srcs"></a>srcs |  -   | <a href="https://bazel.build/concepts/labels">List of labels</a> | required |  |
| <a id="al_unpack_archives-out"></a>out |  -   | String | optional |  `""`  |


