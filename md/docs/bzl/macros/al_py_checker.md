<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_py_checker"></a>

## al_py_checker

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_py_checker.bzl", "al_py_checker")

al_py_checker(<a href="#al_py_checker-name">name</a>, <a href="#al_py_checker-tool">tool</a>, <a href="#al_py_checker-args_bin">args_bin</a>, <a href="#al_py_checker-args_test">args_test</a>, <a href="#al_py_checker-test_size">test_size</a>, <a href="#al_py_checker-disable_fix">disable_fix</a>, <a href="#al_py_checker-kwargs">**kwargs</a>)
</pre>

Create -fix and -test targets for a python checker

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_py_checker-name"></a>name |  Name prefix   |  `None` |
| <a id="al_py_checker-tool"></a>tool |  Tool label   |  `None` |
| <a id="al_py_checker-args_bin"></a>args_bin |  Args for the binary target   |  `None` |
| <a id="al_py_checker-args_test"></a>args_test |  Args for the test   |  `None` |
| <a id="al_py_checker-test_size"></a>test_size |  <p align="center"> - </p>   |  `"small"` |
| <a id="al_py_checker-disable_fix"></a>disable_fix |  If set, do not create fix target   |  `False` |
| <a id="al_py_checker-kwargs"></a>kwargs |  Kwargs for both targets   |  none |


