_RFC_INDEX = """
---
title: Rfc{rfc}
description: Rfc{rfc}
tags:
  - generated
  - rfc
---

## Html

- Link: {html_url}

## Txt

- Original: {txt_url}
- Local: [./rfc.txt](./rfc.txt)

{{{{< readfile file="rfc.txt" code="true" lang="txt" >}}}}
"""

_ROOT_BUILD = """
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files", "strip_prefix")

pkg_filegroup(
    name = "filegroup",
    srcs = [":files"],
    visibility = ["//visibility:public"],
)

pkg_files(
    name = "files",
    srcs = glob(["**/*.txt", "**/*.md"]),
    strip_prefix = strip_prefix.from_pkg(),
    visibility = ["//visibility:public"],
)
"""

def _impl(ctx):
    downloads = []
    for rfc in ctx.attr.rfcs:
        format = "txt"
        url = ctx.attr.url.format(postfix = "rfc{}.{}".format(rfc, format))
        html_url = ctx.attr.url.format(postfix = "rfc{}.{}".format(rfc, "html"))
        integrity = ctx.attr.integrity.get("{}/{}".format(rfc, format), "")
        download = ctx.download(
            url = url,
            integrity = integrity,
            block = False,
            output = "rfc{}/rfc.{}".format(rfc, format),
        )
        downloads.append(["{}/{}".format(rfc, format), integrity, download])
        ctx.file(
            "rfc{}/_index.md".format(rfc),
            content = _RFC_INDEX.format(
                rfc = rfc,
                url = url,
                html_url = html_url,
                txt_url = url,
                file = "rfc.{}".format(format),
            ),
        )
    failures = []
    for rfc, integrity, download in downloads:
        integrity_target = download.wait().integrity
        if integrity != integrity_target:
            failures.append("missing integrity for rfc '{}': want '{}', got '{}'".format(rfc, integrity_target, integrity))
    if failures:
        fail("\n".join(failures))
    ctx.file(
        "BUILD.bazel",
        content = _ROOT_BUILD,
    )

al_rfc_repo = repository_rule(
    doc = "Rfc repository",
    implementation = _impl,
    attrs = {
        "rfcs": attr.string_list(
            doc = "Rfcs",
        ),
        "integrity": attr.string_dict(
            doc = "Rfc integrity",
        ),
        "url": attr.string(
            default = "https://www.rfc-editor.org/rfc/{postfix}",
            doc = "Url format",
        ),
    },
)
