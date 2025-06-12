<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_write_script"></a>

## al_write_script

<pre>
load("@com-alwaldend-git-src//bzl/rules:al_write_script.bzl", "al_write_script")

al_write_script(<a href="#al_write_script-name">name</a>, <a href="#al_write_script-out">out</a>, <a href="#al_write_script-content">content</a>, <a href="#al_write_script-make_vars">make_vars</a>, <a href="#al_write_script-set_flags">set_flags</a>, <a href="#al_write_script-shebang">shebang</a>)
</pre>

Write a script and make it executable

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="al_write_script-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="al_write_script-out"></a>out |  Output file   | <a href="https://bazel.build/concepts/labels">Label</a>; <a href="https://bazel.build/reference/be/common-definitions#configurable-attributes">nonconfigurable</a> | required |  |
| <a id="al_write_script-content"></a>content |  Script content   | String | required |  |
| <a id="al_write_script-make_vars"></a>make_vars |  Additional make vars   | <a href="https://bazel.build/rules/lib/dict">Dictionary: String -> String</a> | optional |  `{}`  |
| <a id="al_write_script-set_flags"></a>set_flags |  Flags to pass to set   | List of strings | optional |  `["-eu"]`  |
| <a id="al_write_script-shebang"></a>shebang |  Sheband to use   | String | optional |  `"#!/usr/bin/env sh"`  |


