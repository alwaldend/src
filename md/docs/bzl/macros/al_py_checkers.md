<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_py_checkers"></a>

## al_py_checkers

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_py_checkers.bzl", "al_py_checkers")

al_py_checkers(<a href="#al_py_checkers-srcs">srcs</a>, <a href="#al_py_checkers-isort_label">isort_label</a>, <a href="#al_py_checkers-black_label">black_label</a>, <a href="#al_py_checkers-mypy_label">mypy_label</a>, <a href="#al_py_checkers-flake8_label">flake8_label</a>, <a href="#al_py_checkers-pyproject_label">pyproject_label</a>)
</pre>

Generate -fix and -test targets for python checkers

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_py_checkers-srcs"></a>srcs |  list of source file labels   |  none |
| <a id="al_py_checkers-isort_label"></a>isort_label |  <p align="center"> - </p>   |  `"//py:isort"` |
| <a id="al_py_checkers-black_label"></a>black_label |  <p align="center"> - </p>   |  `"//py:black"` |
| <a id="al_py_checkers-mypy_label"></a>mypy_label |  <p align="center"> - </p>   |  `"//py:mypy"` |
| <a id="al_py_checkers-flake8_label"></a>flake8_label |  <p align="center"> - </p>   |  `"//py:flake8"` |
| <a id="al_py_checkers-pyproject_label"></a>pyproject_label |  <p align="center"> - </p>   |  `"//:pyproject"` |


