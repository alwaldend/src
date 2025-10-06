load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

_NATIVE_BINARY = """
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")
native_binary(
    name = "{name}",
    src = glob(["**/{src}"])[0],
    visibility = ["//visibility:public"],
)
"""
_NATIVE_BINARY_DEFAULTS = {"name": "bin"}

def _impl(ctx):
    for mod in ctx.modules:
        for tag in mod.tags.download:
            for name, repo in tag.repos.items():
                repo_name = "{}_{}".format(tag.name, name)
                kwargs = {}
                for line in repo:
                    key, value = line.split("=", 1)
                    kwargs[key] = value
                if tag.download_type == "http_archive":
                    if tag.build_file_native_binary:
                        print(_NATIVE_BINARY_DEFAULTS | tag.build_file_native_binary)
                        kwargs.setdefault(
                            "build_file_content",
                            _NATIVE_BINARY.format(**(_NATIVE_BINARY_DEFAULTS | tag.build_file_native_binary)),
                        )
                    elif tag.build_file_content:
                        kwargs.setdefault("build_file_content", tag.build_file_content)
                    kwargs.setdefault("strip_prefix", tag.strip_prefix)
                    http_archive(
                        name = repo_name,
                        **kwargs
                    )
                elif tag.download_type == "http_file":
                    http_file(
                        name = repo_name,
                        executable = tag.executable,
                        **kwargs
                    )
                else:
                    fail("invalid download type: {}".format(tag.download_type))

al_repo_map = module_extension(
    implementation = _impl,
    doc = "Extension to create several repos from a map",
    tag_classes = {
        "download": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Name",
            ),
            "download_type": attr.string(
                values = ["http_archive", "http_file"],
                default = "http_archive",
                doc = "Download type",
            ),
            "executable": attr.bool(
                default = False,
                doc = "Field executable for http_file",
            ),
            "build_file_native_binary": attr.string_dict(
                doc = "Args for a native binary build file",
            ),
            "build_file_content": attr.string(
                doc = "Build file content",
            ),
            "strip_prefix": attr.string(
                doc = "Strip prefix",
            ),
            "repos": attr.string_list_dict(
                mandatory = True,
                doc = "Map of repos",
            ),
        }),
    },
)
