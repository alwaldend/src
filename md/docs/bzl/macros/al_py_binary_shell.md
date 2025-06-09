<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_py_binary_shell"></a>

## al_py_binary_shell

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_py_binary_shell.bzl", "al_py_binary_shell")

al_py_binary_shell(<a href="#al_py_binary_shell-name">name</a>, <a href="#al_py_binary_shell-deps">deps</a>, <a href="#al_py_binary_shell-srcs">srcs</a>, <a href="#al_py_binary_shell-shell_type">shell_type</a>, <a href="#al_py_binary_shell-shell_label">shell_label</a>, <a href="#al_py_binary_shell-kwargs">**kwargs</a>)
</pre>

Create a py_binary target that allows you to run commands in proper python environment

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_py_binary_shell-name"></a>name |  target name   |  none |
| <a id="al_py_binary_shell-deps"></a>deps |  py_binary deps   |  `[]` |
| <a id="al_py_binary_shell-srcs"></a>srcs |  py_binary srcs   |  `[]` |
| <a id="al_py_binary_shell-shell_type"></a>shell_type |  ${BAZEL_PYTHON_SHELL_TYPE}   |  `"python"` |
| <a id="al_py_binary_shell-shell_label"></a>shell_label |  <p align="center"> - </p>   |  `"//py/bazel-python-shell:library"` |
| <a id="al_py_binary_shell-kwargs"></a>kwargs |  other py_binary kwargs   |  none |


