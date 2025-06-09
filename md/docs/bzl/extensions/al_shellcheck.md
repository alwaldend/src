<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_shellcheck"></a>

## al_shellcheck

<pre>
al_shellcheck = use_extension("@com-alwaldend-git-src//bzl/extensions:al_shellcheck.bzl", "al_shellcheck")
al_shellcheck.install(<a href="#al_shellcheck.install-name">name</a>, <a href="#al_shellcheck.install-integrity">integrity</a>, <a href="#al_shellcheck.install-urls">urls</a>)
</pre>

Extension to download shellcheck


**TAG CLASSES**

<a id="al_shellcheck.install"></a>

### install

**Attributes**

| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="al_shellcheck.install-name"></a>name |  name   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="al_shellcheck.install-integrity"></a>integrity |  integrity for the urls   | <a href="https://bazel.build/rules/lib/dict">Dictionary: String -> String</a> | required |  |
| <a id="al_shellcheck.install-urls"></a>urls |  urls   | <a href="https://bazel.build/rules/lib/dict">Dictionary: String -> String</a> | required |  |


