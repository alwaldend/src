---
title: Bazel rules
---

## Api

<!-- STARDOC START -->
<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="al_template_files"></a>

## al_template_files

<pre>
load("@com-alwaldend-src//bzl/rules:al-template-files.bzl", "al_template_files")

al_template_files(<a href="#al_template_files-name">name</a>, <a href="#al_template_files-srcs">srcs</a>, <a href="#al_template_files-data">data</a>, <a href="#al_template_files-outs">outs</a>, <a href="#al_template_files-templater">templater</a>)
</pre>

Load data files, then template the template and write the output

**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="al_template_files-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |
| <a id="al_template_files-srcs"></a>srcs |  Template files   | <a href="https://bazel.build/concepts/labels">List of labels</a> | required |  |
| <a id="al_template_files-data"></a>data |  Data files   | <a href="https://bazel.build/concepts/labels">List of labels</a> | required |  |
| <a id="al_template_files-outs"></a>outs |  Output files   | List of labels; <a href="https://bazel.build/reference/be/common-definitions#configurable-attributes">nonconfigurable</a> | required |  |
| <a id="al_template_files-templater"></a>templater |  Templater to use   | <a href="https://bazel.build/concepts/labels">Label</a> | optional |  `"@com-alwaldend-src//go/template-files"`  |



<!-- STARDOC END -->
