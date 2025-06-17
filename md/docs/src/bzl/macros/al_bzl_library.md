<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_bzl_library"></a>

## al_bzl_library

<pre>
load("@com-alwaldend-git-src//bzl/macros:al_bzl_library.bzl", "al_bzl_library")

al_bzl_library(<a href="#al_bzl_library-name">name</a>, <a href="#al_bzl_library-visibility">visibility</a>, <a href="#al_bzl_library-src">src</a>, <a href="#al_bzl_library-deps">deps</a>)
</pre>

Generate targets for a bzl library

Targets:
- ${name}: bzl_library
- ${name}-stardoc: stardoc target (if src is passed)


**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="al_bzl_library-name"></a>name |  library name   |  none |
| <a id="al_bzl_library-visibility"></a>visibility |  visibility   |  none |
| <a id="al_bzl_library-src"></a>src |  (Optinal[str]): bzl source file   |  `None` |
| <a id="al_bzl_library-deps"></a>deps |  bzl_library deps   |  `[]` |


