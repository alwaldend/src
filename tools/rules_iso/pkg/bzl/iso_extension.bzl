"""Module extension that downloads ISO images and exposes flash targets."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")
load(":iso_flash_repo.bzl", "iso_flash_repo")

_BUILD_TEMPLATE = """load("@rules_iso//pkg/bzl:iso_flash_binary.bzl", "iso_flash_binary")

alias(
    name = "file",
    actual = "{alias_actual}",
    visibility = ["//visibility:public"],
)

iso_flash_binary(
    name = "iso_flash",
    image = "{alias_actual}",
    visibility = ["//visibility:public"],
)
"""

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.image:
            if tag.downloaded_file_path and not tag.downloaded_file_path.endswith(".iso"):
                fail("downloaded_file_path must end with .iso for {} image".format(tag.name))
            if not tag.urls:
                fail("image {} must declare at least one URL".format(tag.name))
            file_name = tag.downloaded_file_path or "{}.iso".format(tag.name)
            http_file(
                name = tag.name,
                downloaded_file_path = file_name,
                integrity = tag.integrity,
                urls = tag.urls,
            )
            root_module_direct_deps.append(tag.name)
            flash_repo_name = "{}_flash".format(tag.name)
            iso_flash_repo(
                name = flash_repo_name,
                image_repo = tag.name,
            )
            root_module_direct_deps.append(flash_repo_name)
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

_image = tag_class(
    {
        "downloaded_file_path": attr.string(
            doc = "File name inside the repository, must end with .iso",
        ),
        "integrity": attr.string(
            mandatory = True,
            doc = "Expected integrity hash of the downloaded ISO",
        ),
        "name": attr.string(
            mandatory = True,
            doc = "Repository name for the downloaded image",
        ),
        "urls": attr.string_list(
            mandatory = True,
            doc = "Mirror URLs for the ISO",
        ),
    },
)

iso_extension = module_extension(
    implementation = _impl,
    doc = "Downloads ISO images and adds an iso_flash run target per image",
    tag_classes = {
        "image": _image,
    },
)
