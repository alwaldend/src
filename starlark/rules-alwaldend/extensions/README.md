---
title: Module extensions
---

## Api

<!-- STARDOC START -->
<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="shellcheck"></a>

## shellcheck

<pre>
shellcheck = use_extension("@com-alwaldend-src//starlark/bazel/extensions:defs.bzl", "shellcheck")
shellcheck.install(<a href="#shellcheck.install-name">name</a>, <a href="#shellcheck.install-integrity">integrity</a>, <a href="#shellcheck.install-urls">urls</a>)
</pre>

Rule to download shellcheck


**TAG CLASSES**

<a id="shellcheck.install"></a>

### install

**Attributes**

| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="shellcheck.install-name"></a>name |  name   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="shellcheck.install-integrity"></a>integrity |  integrity for the urls   | <a href="https://bazel.build/rules/lib/dict">Dictionary: String -> String</a> | required |  |
| <a id="shellcheck.install-urls"></a>urls |  urls   | <a href="https://bazel.build/rules/lib/dict">Dictionary: String -> String</a> | required |  |



<!-- STARDOC END -->
