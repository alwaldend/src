<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_sh_library"></a>

## al_sh_library

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_sh_library.bzl", "al_sh_library")

al_sh_library(<a href="#al_sh_library-name">name</a>, <a href="#al_sh_library-shfmt_src">shfmt_src</a>, <a href="#al_sh_library-editorconfig_src">editorconfig_src</a>, <a href="#al_sh_library-shellcheck_src">shellcheck_src</a>, <a href="#al_sh_library-run_args_src">run_args_src</a>, <a href="#al_sh_library-visibility">visibility</a>,
              <a href="#al_sh_library-common_kwargs">**common_kwargs</a>)
</pre>

Create targets for a shell library

Targets:
- ${name}-shfmt-fix: executable to run shfmt
- ${name}-shfmt-test: test whether the script is formatted
- ${name}-shellcheck-test: shellcheck test


**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_sh_library-name"></a>name |  target name   |  none |
| <a id="al_sh_library-shfmt_src"></a>shfmt_src |  <p align="center"> - </p>   |  `"@cc_mvdan_sh_v3//cmd/shfmt:shfmt"` |
| <a id="al_sh_library-editorconfig_src"></a>editorconfig_src |  <p align="center"> - </p>   |  `"//:.editorconfig"` |
| <a id="al_sh_library-shellcheck_src"></a>shellcheck_src |  <p align="center"> - </p>   |  `"@com-github-koalaman-shellcheck//:shellcheck"` |
| <a id="al_sh_library-run_args_src"></a>run_args_src |  <p align="center"> - </p>   |  `"//sh/scripts:run-args-lib"` |
| <a id="al_sh_library-visibility"></a>visibility |  <p align="center"> - </p>   |  `["//visibility:public"]` |
| <a id="al_sh_library-common_kwargs"></a>common_kwargs |  kwargs for both targets   |  none |


