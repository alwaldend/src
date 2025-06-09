<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_genrule_executable"></a>

## al_genrule_executable

<pre>
load("@com-alwaldend-git-src//bzl/rules:al_genrule.bzl", "al_genrule_executable")

al_genrule_executable(<a href="#al_genrule_executable-name">name</a>, <a href="#al_genrule_executable-srcs">srcs</a>, <a href="#al_genrule_executable-data">data</a>, <a href="#al_genrule_executable-outs">outs</a>, <a href="#al_genrule_executable-cmd">cmd</a>, <a href="#al_genrule_executable-set_flags">set_flags</a>, <a href="#al_genrule_executable-shell">shell</a>, <a href="#al_genrule_executable-tools">tools</a>, <a href="#al_genrule_executable-worker">worker</a>)
</pre>

Build executable using shell worker

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="al_genrule_executable-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="al_genrule_executable-srcs"></a>srcs |  Sources, will not be added to runfiles   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_genrule_executable-data"></a>data |  Data, will be added to runfiles   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_genrule_executable-outs"></a>outs |  Outputs   | List of labels; <a href="https://bazel.build/reference/be/common-definitions#configurable-attributes">nonconfigurable</a> | required |  |
| <a id="al_genrule_executable-cmd"></a>cmd |  Script to execute   | String | required |  |
| <a id="al_genrule_executable-set_flags"></a>set_flags |  set flags   | List of strings | optional |  `["-eux"]`  |
| <a id="al_genrule_executable-shell"></a>shell |  shell to use   | String | optional |  `"/bin/sh"`  |
| <a id="al_genrule_executable-tools"></a>tools |  Tools, will be added to runfiles   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_genrule_executable-worker"></a>worker |  Worker binary   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `"@com-alwaldend-git-src//go/bazel-shell-worker"`  |


<a id="al_genrule_regular"></a>

## al_genrule_regular

<pre>
load("@com-alwaldend-git-src//bzl/rules:al_genrule.bzl", "al_genrule_regular")

al_genrule_regular(<a href="#al_genrule_regular-name">name</a>, <a href="#al_genrule_regular-srcs">srcs</a>, <a href="#al_genrule_regular-data">data</a>, <a href="#al_genrule_regular-outs">outs</a>, <a href="#al_genrule_regular-cmd">cmd</a>, <a href="#al_genrule_regular-set_flags">set_flags</a>, <a href="#al_genrule_regular-shell">shell</a>, <a href="#al_genrule_regular-tools">tools</a>, <a href="#al_genrule_regular-worker">worker</a>)
</pre>

Build shell worker rule

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="al_genrule_regular-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="al_genrule_regular-srcs"></a>srcs |  Sources, will not be added to runfiles   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_genrule_regular-data"></a>data |  Data, will be added to runfiles   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_genrule_regular-outs"></a>outs |  Outputs   | List of labels; <a href="https://bazel.build/reference/be/common-definitions#configurable-attributes">nonconfigurable</a> | required |  |
| <a id="al_genrule_regular-cmd"></a>cmd |  Script to execute   | String | required |  |
| <a id="al_genrule_regular-set_flags"></a>set_flags |  set flags   | List of strings | optional |  `["-eux"]`  |
| <a id="al_genrule_regular-shell"></a>shell |  shell to use   | String | optional |  `"/bin/sh"`  |
| <a id="al_genrule_regular-tools"></a>tools |  Tools, will be added to runfiles   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_genrule_regular-worker"></a>worker |  Worker binary   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `"@com-alwaldend-git-src//go/bazel-shell-worker"`  |


<a id="al_genrule_test"></a>

## al_genrule_test

<pre>
load("@com-alwaldend-git-src//bzl/rules:al_genrule.bzl", "al_genrule_test")

al_genrule_test(<a href="#al_genrule_test-name">name</a>, <a href="#al_genrule_test-srcs">srcs</a>, <a href="#al_genrule_test-data">data</a>, <a href="#al_genrule_test-outs">outs</a>, <a href="#al_genrule_test-cmd">cmd</a>, <a href="#al_genrule_test-set_flags">set_flags</a>, <a href="#al_genrule_test-shell">shell</a>, <a href="#al_genrule_test-tools">tools</a>, <a href="#al_genrule_test-worker">worker</a>)
</pre>

Test using shell worker

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="al_genrule_test-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="al_genrule_test-srcs"></a>srcs |  Sources, will not be added to runfiles   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_genrule_test-data"></a>data |  Data, will be added to runfiles   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_genrule_test-outs"></a>outs |  Outputs   | List of labels; <a href="https://bazel.build/reference/be/common-definitions#configurable-attributes">nonconfigurable</a> | required |  |
| <a id="al_genrule_test-cmd"></a>cmd |  Script to execute   | String | required |  |
| <a id="al_genrule_test-set_flags"></a>set_flags |  set flags   | List of strings | optional |  `["-eux"]`  |
| <a id="al_genrule_test-shell"></a>shell |  shell to use   | String | optional |  `"/bin/sh"`  |
| <a id="al_genrule_test-tools"></a>tools |  Tools, will be added to runfiles   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_genrule_test-worker"></a>worker |  Worker binary   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `"@com-alwaldend-git-src//go/bazel-shell-worker"`  |


