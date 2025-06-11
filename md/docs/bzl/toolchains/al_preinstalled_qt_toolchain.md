<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_preinstalled_qt_toolchain"></a>

## al_preinstalled_qt_toolchain

<pre>
load("@com-alwaldend-git-src//bzl/toolchains:al_preinstalled_qt_toolchain.bzl", "al_preinstalled_qt_toolchain")

al_preinstalled_qt_toolchain(<a href="#al_preinstalled_qt_toolchain-name">name</a>, <a href="#al_preinstalled_qt_toolchain-dir">dir</a>, <a href="#al_preinstalled_qt_toolchain-platform">platform</a>, <a href="#al_preinstalled_qt_toolchain-version">version</a>)
</pre>

Preinstalled qt toolchain

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="al_preinstalled_qt_toolchain-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="al_preinstalled_qt_toolchain-dir"></a>dir |  Root qt directory   | String | optional |  `"/opt/qt"`  |
| <a id="al_preinstalled_qt_toolchain-platform"></a>platform |  Qt platform   | String | optional |  `"gcc_64"`  |
| <a id="al_preinstalled_qt_toolchain-version"></a>version |  Qt version   | String | optional |  `"6.9.0"`  |


