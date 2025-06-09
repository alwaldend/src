<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_run_tool"></a>

## al_run_tool

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_run_tool.bzl", "al_run_tool")

al_run_tool(<a href="#al_run_tool-name">name</a>, <a href="#al_run_tool-tool">tool</a>, <a href="#al_run_tool-executable">executable</a>, <a href="#al_run_tool-test">test</a>, <a href="#al_run_tool-run_args_label">run_args_label</a>, <a href="#al_run_tool-kwargs">**kwargs</a>)
</pre>

Generate either sh_test, sh_binary, or run_binary target

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_run_tool-name"></a>name |  Target name (required)   |  none |
| <a id="al_run_tool-tool"></a>tool |  Tool label to run (required)   |  none |
| <a id="al_run_tool-executable"></a>executable |  If True, generate sh_binary   |  `False` |
| <a id="al_run_tool-test"></a>test |  If True, generate sh_test   |  `False` |
| <a id="al_run_tool-run_args_label"></a>run_args_label |  <p align="center"> - </p>   |  `"//sh/scripts:run-args-lib"` |
| <a id="al_run_tool-kwargs"></a>kwargs |  kwargs for rules   |  none |


