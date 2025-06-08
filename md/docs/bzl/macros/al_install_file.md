<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_install_file"></a>

## al_install_file

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_install_file.bzl", "al_install_file")

al_install_file(<a href="#al_install_file-name">name</a>, <a href="#al_install_file-args">args</a>, <a href="#al_install_file-install_file_label">install_file_label</a>, <a href="#al_install_file-visibility">visibility</a>, <a href="#al_install_file-py_binary_kwargs">**py_binary_kwargs</a>)
</pre>

Create py_binary target to install file

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_install_file-name"></a>name |  target name   |  none |
| <a id="al_install_file-args"></a>args |  install-file args   |  `[]` |
| <a id="al_install_file-install_file_label"></a>install_file_label |  <p align="center"> - </p>   |  `"//py/install-file:lib"` |
| <a id="al_install_file-visibility"></a>visibility |  <p align="center"> - </p>   |  `["//visibility:public"]` |
| <a id="al_install_file-py_binary_kwargs"></a>py_binary_kwargs |  <p align="center"> - </p>   |  none |


