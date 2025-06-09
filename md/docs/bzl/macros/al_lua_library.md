<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_lua_library"></a>

## al_lua_library

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_lua_library.bzl", "al_lua_library")

al_lua_library(<a href="#al_lua_library-name">name</a>, <a href="#al_lua_library-srcs">srcs</a>, <a href="#al_lua_library-check">check</a>, <a href="#al_lua_library-stylua_config_label">stylua_config_label</a>, <a href="#al_lua_library-stylua_label">stylua_label</a>, <a href="#al_lua_library-pkg_tar_kwargs">pkg_tar_kwargs</a>, <a href="#al_lua_library-visibility">visibility</a>)
</pre>

Generate targets for a lua library

Targets:
    ${name}: pkg_tar
    ${name}-stylua-fix: al_run_tool executable
    ${name}-stylua-test: al_run_tool test


**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_lua_library-name"></a>name |  library name   |  none |
| <a id="al_lua_library-srcs"></a>srcs |  library sources   |  none |
| <a id="al_lua_library-check"></a>check |  if set, only these files will be checked   |  `[]` |
| <a id="al_lua_library-stylua_config_label"></a>stylua_config_label |  <p align="center"> - </p>   |  `"//lua:stylua.toml"` |
| <a id="al_lua_library-stylua_label"></a>stylua_label |  <p align="center"> - </p>   |  `"@com-alwaldend-src-cargo//:stylua__stylua"` |
| <a id="al_lua_library-pkg_tar_kwargs"></a>pkg_tar_kwargs |  kwargs for pkg_tar   |  `{}` |
| <a id="al_lua_library-visibility"></a>visibility |  visibility   |  `["//visibility:public"]` |


