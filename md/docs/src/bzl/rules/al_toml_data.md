<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_toml_data"></a>

## al_toml_data

<pre>
load("@com-alwaldend-git-src//bzl/rules:al_toml_data.bzl", "al_toml_data")

al_toml_data(<a href="#al_toml_data-name">name</a>, <a href="#al_toml_data-deps">deps</a>, <a href="#al_toml_data-srcs">srcs</a>, <a href="#al_toml_data-tomlv">tomlv</a>)
</pre>



**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="al_toml_data-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="al_toml_data-deps"></a>deps |  Toml data targets   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_toml_data-srcs"></a>srcs |  Toml files   | <a href="https://bazel.build/concepts/labels">List of labels</a> | optional |  `[]`  |
| <a id="al_toml_data-tomlv"></a>tomlv |  Tomlv label for the validation aspect   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `"@com-alwaldend-git-src//tools:tomlv"`  |


