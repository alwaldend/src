load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

def _impl(ctx):
    for mod in ctx.modules:
        for tag in mod.tags.download:
            for name, repo in tag.repos.items():
                repo_name = "{}-{}".format(tag.name, name)
                kwargs = {}
                for line in repo:
                    key, value = line.split("=", 1)
                    kwargs[key] = value

                if tag.download_type == "http_archive":
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
                mandatory = True,
                values = ["http_archive", "http_file"],
                doc = "Download type",
            ),
            "executable": attr.bool(
                default = False,
                doc = "Field executable for http_file",
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
