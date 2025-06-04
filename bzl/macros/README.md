---
title: Bazel macros
---

## Api

<!-- STARDOC START -->
<!-- Generated with Stardoc: http://skydoc.bazel.build -->

<a id="al_apply_patches"></a>

## al_apply_patches

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_apply_patches")

al_apply_patches(<a href="#al_apply_patches-name">name</a>, <a href="#al_apply_patches-src">src</a>, <a href="#al_apply_patches-patches">patches</a>, <a href="#al_apply_patches-visibility">visibility</a>)
</pre>

Apply patches

**PARAMETERS**

| Name                                               | Description          | Default Value             |
| :------------------------------------------------- | :------------------- | :------------------------ |
| <a id="al_apply_patches-name"></a>name             | target name          | `"patched"`               |
| <a id="al_apply_patches-src"></a>src               | source archive label | `""`                      |
| <a id="al_apply_patches-patches"></a>patches       | patches label        | `""`                      |
| <a id="al_apply_patches-visibility"></a>visibility | visibility           | `["//visibility:public"]` |

<a id="al_bzl_library"></a>

## al_bzl_library

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_bzl_library")

al_bzl_library(<a href="#al_bzl_library-name">name</a>, <a href="#al_bzl_library-deps">deps</a>, <a href="#al_bzl_library-srcs">srcs</a>, <a href="#al_bzl_library-out">out</a>, <a href="#al_bzl_library-input">input</a>, <a href="#al_bzl_library-replace_section_src">replace_section_src</a>, <a href="#al_bzl_library-run_args_src">run_args_src</a>, <a href="#al_bzl_library-visibility">visibility</a>)
</pre>

Generate stardoc and bzl_library targets

**PARAMETERS**

| Name                                                               | Description               | Default Value                    |
| :----------------------------------------------------------------- | :------------------------ | :------------------------------- |
| <a id="al_bzl_library-name"></a>name                               | library name              | none                             |
| <a id="al_bzl_library-deps"></a>deps                               | library dependencies      | `[]`                             |
| <a id="al_bzl_library-srcs"></a>srcs                               | library sources           | `["defs.bzl"]`                   |
| <a id="al_bzl_library-out"></a>out                                 | generated md file         | `"api.md"`                       |
| <a id="al_bzl_library-input"></a>input                             | input file                | `"defs.bzl"`                     |
| <a id="al_bzl_library-replace_section_src"></a>replace_section_src | <p align="center"> - </p> | `"//python/replace-section"`     |
| <a id="al_bzl_library-run_args_src"></a>run_args_src               | <p align="center"> - </p> | `"//shell/scripts:run-args-lib"` |
| <a id="al_bzl_library-visibility"></a>visibility                   | <p align="center"> - </p> | `["//visibility:public"]`        |

<a id="al_combine_files"></a>

## al_combine_files

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_combine_files")

al_combine_files(<a href="#al_combine_files-name">name</a>, <a href="#al_combine_files-srcs">srcs</a>, <a href="#al_combine_files-genrule_kwargs">**genrule_kwargs</a>)
</pre>

Combine several files into one

**PARAMETERS**

| Name                                                       | Description               | Default Value |
| :--------------------------------------------------------- | :------------------------ | :------------ |
| <a id="al_combine_files-name"></a>name                     | target name               | `""`          |
| <a id="al_combine_files-srcs"></a>srcs                     | sources to combine        | `[]`          |
| <a id="al_combine_files-genrule_kwargs"></a>genrule_kwargs | <p align="center"> - </p> | none          |

<a id="al_compile_pip_requirements_combined"></a>

## al_compile_pip_requirements_combined

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_compile_pip_requirements_combined")

al_compile_pip_requirements_combined(<a href="#al_compile_pip_requirements_combined-name">name</a>, <a href="#al_compile_pip_requirements_combined-srcs">srcs</a>, <a href="#al_compile_pip_requirements_combined-compile_kwargs">**compile_kwargs</a>)
</pre>

Compile seveal requirement files

**PARAMETERS**

| Name                                                                           | Description                           | Default Value |
| :----------------------------------------------------------------------------- | :------------------------------------ | :------------ |
| <a id="al_compile_pip_requirements_combined-name"></a>name                     | name                                  | `""`          |
| <a id="al_compile_pip_requirements_combined-srcs"></a>srcs                     | requirement files to combine          | `[]`          |
| <a id="al_compile_pip_requirements_combined-compile_kwargs"></a>compile_kwargs | kwargs for `compile_pip_requirements` | none          |

<a id="al_genrule_src"></a>

## al_genrule_src

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_genrule_src")

al_genrule_src(<a href="#al_genrule_src-name">name</a>, <a href="#al_genrule_src-patterns">patterns</a>, <a href="#al_genrule_src-visibility">visibility</a>)
</pre>

Create a filegroup and a genrule generating a tar archive

**PARAMETERS**

| Name                                             | Description               | Default Value             |
| :----------------------------------------------- | :------------------------ | :------------------------ |
| <a id="al_genrule_src-name"></a>name             | <p align="center"> - </p> | `"src"`                   |
| <a id="al_genrule_src-patterns"></a>patterns     | <p align="center"> - </p> | `["**"]`                  |
| <a id="al_genrule_src-visibility"></a>visibility | <p align="center"> - </p> | `["//visibility:public"]` |

<a id="al_genrule_with_wheels"></a>

## al_genrule_with_wheels

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_genrule_with_wheels")

al_genrule_with_wheels(<a href="#al_genrule_with_wheels-wheels">wheels</a>, <a href="#al_genrule_with_wheels-srcs">srcs</a>, <a href="#al_genrule_with_wheels-cmd">cmd</a>, <a href="#al_genrule_with_wheels-kwargs">**kwargs</a>)
</pre>

Regular genrule with wheels added to ${PYTHONPATH}

**PARAMETERS**

| Name                                             | Description               | Default Value |
| :----------------------------------------------- | :------------------------ | :------------ |
| <a id="al_genrule_with_wheels-wheels"></a>wheels | <p align="center"> - </p> | `None`        |
| <a id="al_genrule_with_wheels-srcs"></a>srcs     | <p align="center"> - </p> | `[]`          |
| <a id="al_genrule_with_wheels-cmd"></a>cmd       | <p align="center"> - </p> | `[]`          |
| <a id="al_genrule_with_wheels-kwargs"></a>kwargs | <p align="center"> - </p> | none          |

<a id="al_install_file"></a>

## al_install_file

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_install_file")

al_install_file(<a href="#al_install_file-name">name</a>, <a href="#al_install_file-args">args</a>, <a href="#al_install_file-install_file_src">install_file_src</a>, <a href="#al_install_file-visibility">visibility</a>, <a href="#al_install_file-py_binary_kwargs">**py_binary_kwargs</a>)
</pre>

**PARAMETERS**

| Name                                                          | Description               | Default Value                 |
| :------------------------------------------------------------ | :------------------------ | :---------------------------- |
| <a id="al_install_file-name"></a>name                         | <p align="center"> - </p> | `None`                        |
| <a id="al_install_file-args"></a>args                         | <p align="center"> - </p> | `[]`                          |
| <a id="al_install_file-install_file_src"></a>install_file_src | <p align="center"> - </p> | `"//python/install-file:lib"` |
| <a id="al_install_file-visibility"></a>visibility             | <p align="center"> - </p> | `["//visibility:public"]`     |
| <a id="al_install_file-py_binary_kwargs"></a>py_binary_kwargs | <p align="center"> - </p> | none                          |

<a id="al_lua_library"></a>

## al_lua_library

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_lua_library")

al_lua_library(<a href="#al_lua_library-name">name</a>, <a href="#al_lua_library-srcs">srcs</a>, <a href="#al_lua_library-check">check</a>, <a href="#al_lua_library-run_args_src">run_args_src</a>, <a href="#al_lua_library-stylua_config_src">stylua_config_src</a>, <a href="#al_lua_library-stylua_src">stylua_src</a>, <a href="#al_lua_library-pkg_tar_kwargs">pkg_tar_kwargs</a>,
               <a href="#al_lua_library-visibility">visibility</a>)
</pre>

Generate targets for a lua library

**PARAMETERS**

| Name                                                           | Description                              | Default Value                    |
| :------------------------------------------------------------- | :--------------------------------------- | :------------------------------- |
| <a id="al_lua_library-name"></a>name                           | library name                             | none                             |
| <a id="al_lua_library-srcs"></a>srcs                           | library sources                          | `[]`                             |
| <a id="al_lua_library-check"></a>check                         | if set, only these files will be checked | `[]`                             |
| <a id="al_lua_library-run_args_src"></a>run_args_src           | <p align="center"> - </p>                | `"//shell/scripts:run-args-lib"` |
| <a id="al_lua_library-stylua_config_src"></a>stylua_config_src | <p align="center"> - </p>                | `"//lua:stylua.toml"`            |
| <a id="al_lua_library-stylua_src"></a>stylua_src               | <p align="center"> - </p>                | `"@cargo//:stylua__stylua"`      |
| <a id="al_lua_library-pkg_tar_kwargs"></a>pkg_tar_kwargs       | kwargs for pkg_tar                       | `{}`                             |
| <a id="al_lua_library-visibility"></a>visibility               | visibility                               | `["//visibility:public"]`        |

<a id="al_pkg_tar_combined"></a>

## al_pkg_tar_combined

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_pkg_tar_combined")

al_pkg_tar_combined(<a href="#al_pkg_tar_combined-name">name</a>, <a href="#al_pkg_tar_combined-tars">tars</a>, <a href="#al_pkg_tar_combined-strip_components">strip_components</a>, <a href="#al_pkg_tar_combined-genrule_kwargs">**genrule_kwargs</a>)
</pre>

Combine several tars into one

**PARAMETERS**

| Name                                                              | Description               | Default Value |
| :---------------------------------------------------------------- | :------------------------ | :------------ |
| <a id="al_pkg_tar_combined-name"></a>name                         | <p align="center"> - </p> | `None`        |
| <a id="al_pkg_tar_combined-tars"></a>tars                         | <p align="center"> - </p> | `[]`          |
| <a id="al_pkg_tar_combined-strip_components"></a>strip_components | <p align="center"> - </p> | `2`           |
| <a id="al_pkg_tar_combined-genrule_kwargs"></a>genrule_kwargs     | <p align="center"> - </p> | none          |

<a id="al_py_binary_shell"></a>

## al_py_binary_shell

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_py_binary_shell")

al_py_binary_shell(<a href="#al_py_binary_shell-name">name</a>, <a href="#al_py_binary_shell-cmd">cmd</a>, <a href="#al_py_binary_shell-deps">deps</a>, <a href="#al_py_binary_shell-srcs">srcs</a>, <a href="#al_py_binary_shell-shell_type">shell_type</a>, <a href="#al_py_binary_shell-kwargs">**kwargs</a>)
</pre>

Create a shell with linked python deps

**PARAMETERS**

| Name                                                 | Description               | Default Value |
| :--------------------------------------------------- | :------------------------ | :------------ |
| <a id="al_py_binary_shell-name"></a>name             | <p align="center"> - </p> | `""`          |
| <a id="al_py_binary_shell-cmd"></a>cmd               | <p align="center"> - </p> | `""`          |
| <a id="al_py_binary_shell-deps"></a>deps             | <p align="center"> - </p> | `[]`          |
| <a id="al_py_binary_shell-srcs"></a>srcs             | <p align="center"> - </p> | `[]`          |
| <a id="al_py_binary_shell-shell_type"></a>shell_type | <p align="center"> - </p> | `"python"`    |
| <a id="al_py_binary_shell-kwargs"></a>kwargs         | <p align="center"> - </p> | none          |

<a id="al_py_black"></a>

## al_py_black

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_py_black")

al_py_black(<a href="#al_py_black-name">name</a>, <a href="#al_py_black-srcs">srcs</a>, <a href="#al_py_black-black_src">black_src</a>, <a href="#al_py_black-pyproject_src">pyproject_src</a>, <a href="#al_py_black-run_args_src">run_args_src</a>)
</pre>

Generate -fix and -test targets for black

**PARAMETERS**

| Name                                                | Description               | Default Value |
| :-------------------------------------------------- | :------------------------ | :------------ |
| <a id="al_py_black-name"></a>name                   | <p align="center"> - </p> | `None`        |
| <a id="al_py_black-srcs"></a>srcs                   | <p align="center"> - </p> | `None`        |
| <a id="al_py_black-black_src"></a>black_src         | <p align="center"> - </p> | `None`        |
| <a id="al_py_black-pyproject_src"></a>pyproject_src | <p align="center"> - </p> | `None`        |
| <a id="al_py_black-run_args_src"></a>run_args_src   | <p align="center"> - </p> | `None`        |

<a id="al_py_checker"></a>

## al_py_checker

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_py_checker")

al_py_checker(<a href="#al_py_checker-name">name</a>, <a href="#al_py_checker-args_bin">args_bin</a>, <a href="#al_py_checker-args_test">args_test</a>, <a href="#al_py_checker-disable_fix">disable_fix</a>, <a href="#al_py_checker-kwargs">**kwargs</a>)
</pre>

Create -fix and -test targets for a python checker

**PARAMETERS**

| Name                                              | Description               | Default Value |
| :------------------------------------------------ | :------------------------ | :------------ |
| <a id="al_py_checker-name"></a>name               | <p align="center"> - </p> | `None`        |
| <a id="al_py_checker-args_bin"></a>args_bin       | <p align="center"> - </p> | `None`        |
| <a id="al_py_checker-args_test"></a>args_test     | <p align="center"> - </p> | `None`        |
| <a id="al_py_checker-disable_fix"></a>disable_fix | <p align="center"> - </p> | `False`       |
| <a id="al_py_checker-kwargs"></a>kwargs           | <p align="center"> - </p> | none          |

<a id="al_py_checkers"></a>

## al_py_checkers

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_py_checkers")

al_py_checkers(<a href="#al_py_checkers-srcs">srcs</a>, <a href="#al_py_checkers-black_src">black_src</a>, <a href="#al_py_checkers-isort_src">isort_src</a>, <a href="#al_py_checkers-mypy_src">mypy_src</a>, <a href="#al_py_checkers-pyproject_src">pyproject_src</a>, <a href="#al_py_checkers-run_args_src">run_args_src</a>)
</pre>

Generate -fix and -test targets for python checkers

**PARAMETERS**

| Name                                                   | Description               | Default Value                    |
| :----------------------------------------------------- | :------------------------ | :------------------------------- |
| <a id="al_py_checkers-srcs"></a>srcs                   | <p align="center"> - </p> | `[]`                             |
| <a id="al_py_checkers-black_src"></a>black_src         | <p align="center"> - </p> | `"//python:black"`               |
| <a id="al_py_checkers-isort_src"></a>isort_src         | <p align="center"> - </p> | `"//python:isort"`               |
| <a id="al_py_checkers-mypy_src"></a>mypy_src           | <p align="center"> - </p> | `"//python:mypy"`                |
| <a id="al_py_checkers-pyproject_src"></a>pyproject_src | <p align="center"> - </p> | `"//:pyproject"`                 |
| <a id="al_py_checkers-run_args_src"></a>run_args_src   | <p align="center"> - </p> | `"//shell/scripts:run-args-lib"` |

<a id="al_py_isort"></a>

## al_py_isort

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_py_isort")

al_py_isort(<a href="#al_py_isort-name">name</a>, <a href="#al_py_isort-srcs">srcs</a>, <a href="#al_py_isort-isort_src">isort_src</a>, <a href="#al_py_isort-pyproject_src">pyproject_src</a>, <a href="#al_py_isort-run_args_src">run_args_src</a>)
</pre>

Generate -fix and -test targets for isort

**PARAMETERS**

| Name                                                | Description               | Default Value |
| :-------------------------------------------------- | :------------------------ | :------------ |
| <a id="al_py_isort-name"></a>name                   | <p align="center"> - </p> | `None`        |
| <a id="al_py_isort-srcs"></a>srcs                   | <p align="center"> - </p> | `None`        |
| <a id="al_py_isort-isort_src"></a>isort_src         | <p align="center"> - </p> | `None`        |
| <a id="al_py_isort-pyproject_src"></a>pyproject_src | <p align="center"> - </p> | `None`        |
| <a id="al_py_isort-run_args_src"></a>run_args_src   | <p align="center"> - </p> | `None`        |

<a id="al_py_mypy"></a>

## al_py_mypy

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_py_mypy")

al_py_mypy(<a href="#al_py_mypy-name">name</a>, <a href="#al_py_mypy-srcs">srcs</a>, <a href="#al_py_mypy-mypy_src">mypy_src</a>, <a href="#al_py_mypy-pyproject_src">pyproject_src</a>, <a href="#al_py_mypy-run_args_src">run_args_src</a>)
</pre>

Generate -fix and -test targets for mypy

**PARAMETERS**

| Name                                               | Description               | Default Value |
| :------------------------------------------------- | :------------------------ | :------------ |
| <a id="al_py_mypy-name"></a>name                   | <p align="center"> - </p> | `None`        |
| <a id="al_py_mypy-srcs"></a>srcs                   | <p align="center"> - </p> | `None`        |
| <a id="al_py_mypy-mypy_src"></a>mypy_src           | <p align="center"> - </p> | `None`        |
| <a id="al_py_mypy-pyproject_src"></a>pyproject_src | <p align="center"> - </p> | `None`        |
| <a id="al_py_mypy-run_args_src"></a>run_args_src   | <p align="center"> - </p> | `None`        |

<a id="al_sh_script"></a>

## al_sh_script

<pre>
load("@com-alwaldend-src//starlark/bazel/macros:defs.bzl", "al_sh_script")

al_sh_script(<a href="#al_sh_script-name">name</a>, <a href="#al_sh_script-visibility">visibility</a>, <a href="#al_sh_script-common_kwargs">**common_kwargs</a>)
</pre>

Generate sh_library and sh_binary targets

**PARAMETERS**

| Name                                                 | Description               | Default Value             |
| :--------------------------------------------------- | :------------------------ | :------------------------ |
| <a id="al_sh_script-name"></a>name                   | name                      | `""`                      |
| <a id="al_sh_script-visibility"></a>visibility       | <p align="center"> - </p> | `["//visibility:public"]` |
| <a id="al_sh_script-common_kwargs"></a>common_kwargs | kwargs for both targets   | none                      |

<!-- STARDOC END -->
