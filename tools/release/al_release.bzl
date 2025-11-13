load("@rules_pkg//pkg:providers.bzl", "PackageFilesInfo")
load("//tools/release:al_release_files_info.bzl", "AlReleaseFilesInfo")

_DOC_TEMPLATE = """\
---
title: {release_name}
description: Release {release_name}
tags:
  - generated
  - release
---

## Files

{{{{< alwaldend/release_files >}}}}
"""

def _impl(ctx):
    doc = ctx.actions.declare_file("{}.index.md".format(ctx.label.name))
    dest_src_map = {"index.md": doc}
    release_name = ""
    transient = []

    if ctx.attr.manifest:
        release_name = ctx.file.manifest.basename
        dest_src_map["release.json"] = ctx.file.manifest
    elif ctx.attr.srcs:
        release_name = ctx.attr.release_name
        manifest = ctx.actions.declare_file("{}.release.json".format(ctx.label.name))
        transient = [dep[DefaultInfo].files for dep in ctx.attr.srcs]
        dest_src_map["release.json"] = manifest
        for src in ctx.files.srcs:
            dest_src_map[src.basename] = src
        args = ctx.actions.args()
        args.add("generate")
        inputs = []
        args.add_all(["--output", manifest.path])
        for src in ctx.attr.srcs:
            merge_manifest = src[AlReleaseFilesInfo].manifest
            inputs.append(merge_manifest)
            args.add_all(["--item", "{}:{}".format("manifest", merge_manifest.path)])
        ctx.actions.run(
            executable = ctx.executable.release_tool,
            arguments = [args],
            inputs = inputs,
            outputs = [manifest],
        )
    else:
        fail("missing both manifest and srcs")

    if not release_name:
        fail("missing release_name")

    ctx.actions.write(
        output = doc,
        content = _DOC_TEMPLATE.format(release_name = release_name),
    )

    return [
        DefaultInfo(
            files = depset(
                direct = dest_src_map.values(),
                transitive = transient,
            ),
        ),
        PackageFilesInfo(
            attributes = {},
            dest_src_map = {
                "{}/{}/{}".format(ctx.attr.prefix, release_name, name): file
                for name, file in dest_src_map.items()
            },
        ),
    ]

al_release = rule(
    implementation = _impl,
    doc = "Rule describing a release",
    provides = [PackageFilesInfo],
    attrs = {
        "srcs": attr.label_list(
            providers = [AlReleaseFilesInfo],
            doc = "Sources",
        ),
        "release_name": attr.string(
            mandatory = True,
            doc = "Release name",
        ),
        "prefix": attr.string(
            mandatory = True,
            doc = "Prefix to apply to all generated files",
        ),
        "manifest": attr.label(
            doc = "Load manifest from a file instead of generating it from srcs",
        ),
        "release_tool": attr.label(
            executable = True,
            cfg = "exec",
            default = "//tools/release",
            doc = "Release tool",
        ),
    },
)
