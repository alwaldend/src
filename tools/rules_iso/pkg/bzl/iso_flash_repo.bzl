"""Repository rule for the ISO flash target."""

_FLASH_BINARY = Label(":iso_flash_binary.bzl")

_BUILD = """\
load("{flash_binary}", "iso_flash_binary")

iso_flash_binary(
    name = "iso_flash",
    image = "@{image_repo}//file:file",
    visibility = ["//visibility:public"],
)
"""

def _impl(ctx):
    ctx.file("BUILD.bazel", _BUILD.format(
        flash_binary = _FLASH_BINARY,
        image_repo = ctx.attr.image_repo,
    ))
    return ctx.repo_metadata(reproducible = True)

iso_flash_repo = repository_rule(
    implementation = _impl,
    doc = "Flash target repository for an ISO image",
    attrs = {
        "image_repo": attr.string(
            mandatory = True,
            doc = "Name of the repository holding the downloaded ISO",
        ),
    },
)
