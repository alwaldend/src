---
title: Qt wrapper for bazel
---

## Usage

- Install qt: `aqt install-qt -O /opt/qt linux desktop 6.9.0`
- Register toolchain: `register_toolchains("//bazel/qt:preinstalled-qt-toolchain")`

## Api

<!-- STARDOC START -->
<!-- Generated with Stardoc: http://skydoc.bazel.build -->

<a id="current_qt_toolchain"></a>

## current_qt_toolchain

<pre>
load("@com-alwaldend-src//starlark/bazel/qt:defs.bzl", "current_qt_toolchain")

current_qt_toolchain(<a href="#current_qt_toolchain-name">name</a>)
</pre>

Get current selected qt toolchain

**ATTRIBUTES**

| Name                                       | Description                    | Type                                                                | Mandatory | Default |
| :----------------------------------------- | :----------------------------- | :------------------------------------------------------------------ | :-------- | :------ |
| <a id="current_qt_toolchain-name"></a>name | A unique name for this target. | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required  |         |

<a id="preinstalled_qt_toolchain"></a>

## preinstalled_qt_toolchain

<pre>
load("@com-alwaldend-src//starlark/bazel/qt:defs.bzl", "preinstalled_qt_toolchain")

preinstalled_qt_toolchain(<a href="#preinstalled_qt_toolchain-name">name</a>, <a href="#preinstalled_qt_toolchain-dir">dir</a>, <a href="#preinstalled_qt_toolchain-platform">platform</a>, <a href="#preinstalled_qt_toolchain-version">version</a>)
</pre>

Preinstalled qt toolchain

**ATTRIBUTES**

| Name                                                    | Description                    | Type                                                                | Mandatory | Default     |
| :------------------------------------------------------ | :----------------------------- | :------------------------------------------------------------------ | :-------- | :---------- |
| <a id="preinstalled_qt_toolchain-name"></a>name         | A unique name for this target. | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required  |             |
| <a id="preinstalled_qt_toolchain-dir"></a>dir           | Root qt directory              | String                                                              | optional  | `"/opt/qt"` |
| <a id="preinstalled_qt_toolchain-platform"></a>platform | Qt platform                    | String                                                              | optional  | `"gcc_64"`  |
| <a id="preinstalled_qt_toolchain-version"></a>version   | Qt version                     | String                                                              | optional  | `"6.9.0"`   |

<a id="register_qt_toolchains"></a>

## register_qt_toolchains

<pre>
load("@com-alwaldend-src//starlark/bazel/qt:defs.bzl", "register_qt_toolchains")

register_qt_toolchains()
</pre>

Register qt toolchains

<!-- STARDOC END -->
