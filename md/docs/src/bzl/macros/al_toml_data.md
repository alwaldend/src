<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_toml_data"></a>

## al_toml_data

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_toml_data.bzl", "al_toml_data")

al_toml_data(<a href="#al_toml_data-name">name</a>, <a href="#al_toml_data-srcs">srcs</a>, <a href="#al_toml_data-visibility">visibility</a>, <a href="#al_toml_data-tomlv_label">tomlv_label</a>, <a href="#al_toml_data-size">size</a>, <a href="#al_toml_data-kwargs">**kwargs</a>)
</pre>

Create toml data targets

Targets:
- ${name}: filegroup
- ${name}-test-load: al_run_tool test whether the file can be loaded


**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_toml_data-name"></a>name |  Target name   |  none |
| <a id="al_toml_data-srcs"></a>srcs |  Toml files   |  none |
| <a id="al_toml_data-visibility"></a>visibility |  visibility   |  `["//visibility:public"]` |
| <a id="al_toml_data-tomlv_label"></a>tomlv_label |  <p align="center"> - </p>   |  `"//tools:tomlv"` |
| <a id="al_toml_data-size"></a>size |  <p align="center"> - </p>   |  `"small"` |
| <a id="al_toml_data-kwargs"></a>kwargs |  kwargs   |  none |


