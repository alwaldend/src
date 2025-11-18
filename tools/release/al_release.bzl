load("@rules_pkg//pkg:providers.bzl", "PackageFilesInfo")
load("//tools/release:al_release_files_info.bzl", "AlReleaseFilesInfo")

_DOC_TEMPLATE = """\
---
title: {release_name}
description: Release {release_name}
toc_hide: true
hide_summary: true
tags:
    - generated
    - release
---

{{{{< alwaldend/release standalone=true >}}}}
"""

def _impl(ctx):
    doc = ctx.actions.declare_file("{}.index.md".format(ctx.label.name))
    git = ctx.toolchains["//tools/git/main/bzl:toolchain_type"]
    dest_src_map = {
        "index.md": doc,
    }
    release_name = ctx.attr.release_name
    transitive = []
    manifest = None

    if ctx.attr.manifest:
        manifest = ctx.file.manifest
    elif ctx.attr.srcs:
        manifest = ctx.actions.declare_file("{}.release.json".format(ctx.label.name))
        extra_manifest = ctx.actions.declare_file("{}.release_extra.json".format(ctx.label.name))
        inputs = ctx.files.git_state + [extra_manifest]
        transitive = [dep[DefaultInfo].files for dep in ctx.attr.srcs]
        for src in ctx.files.srcs:
            dest_src_map[src.basename] = src
        args = ctx.actions.args()
        args.add("generate")
        args.add_all(["--output", manifest.path])
        for src in ctx.attr.srcs:
            src_manifest = src[AlReleaseFilesInfo].manifest
            inputs.append(src_manifest)
            args.add_all(["--manifest", src_manifest.path])
        args.add_all(["--manifest", extra_manifest.path])
        args.add_all(["--git_root", git.git_root])
        ctx.actions.write(
            output = extra_manifest,
            content = json.encode(
                {
                    "name": ctx.attr.release_name,
                    "project": {
                        "subdir": ctx.attr.project,
                    },
                    "git": {
                        "revision": "HEAD",
                    },
                },
            ),
        )
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
    dest_src_map["release.json"] = manifest
    dest_src_map = {
        "{}/{}/releases/{}/{}".format(ctx.attr.root_prefix, ctx.attr.project, release_name, name): file
        for name, file in dest_src_map.items()
    }

    return [
        DefaultInfo(
            files = depset(
                direct = dest_src_map.values(),
                transitive = transitive,
            ),
        ),
        PackageFilesInfo(
            attributes = {},
            dest_src_map = dest_src_map,
        ),
    ]

al_release = rule(
    implementation = _impl,
    doc = "Rule describing a release",
    provides = [PackageFilesInfo],
    toolchains = ["//tools/git/main/bzl:toolchain_type"],
    attrs = {
        "srcs": attr.label_list(
            providers = [AlReleaseFilesInfo],
            doc = "Sources",
        ),
        "release_name": attr.string(
            mandatory = True,
            doc = "Release name",
        ),
        "project": attr.string(
            mandatory = True,
            doc = "Project package_name",
        ),
        "root_prefix": attr.string(
            default = "content/docs",
            doc = "Root prefix",
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
        "git_state": attr.label(
            default = "@com_alwaldend_src_tools_git//:git_state",
            doc = "Files that should invalidate the cache on new commit",
        ),
    },
)
