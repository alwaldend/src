load(":binary_toolchain_parse_lock.bzl", "binary_toolchain_parse_lock_attr")
load(":binary_toolchain_repo.bzl", "binary_toolchain_repo")

_BUILD = """
load("@binary_toolchain//main/bzl:binary_toolchain_toolchain.bzl", "binary_toolchain")

CHILDREN = {children}

toolchain_type(
    name = "toolchain_type",
    visibility = ["//visibility:public"],
)

[
    [
        toolchain(
            name = "toolchain_{{}}".format(child),
            toolchain = ":toolchain_{{}}_impl".format(child),
            toolchain_type = "@com_alwaldend_src//tools/trufflehog/main/bzl:toolchain_type",
            exec_compatible_with = platforms,
            visibility = ["//visibility:public"],
        ),
        al_trufflehog_toolchain(
            name = "toolchain_{{}}_impl".format(child),
            exec_compatible_with = platforms,
            trufflehog = "@{{}}//:trufflehog_bin".format(child),
        )
    ]
    for child, platforms in CHILDREN
]
"""

def _impl(ctx):
    lock = binary_toolchain_parse_lock_attr(ctx)
    repos = []
    for archive_name, _ in lock.version_archives.items():
        name = "{}_{}".format(ctx.original_name, archive_name)
        repos.append(name)
        binary_toolchain_repo(
            name = name,
            lock = ctx.attr.lock,
            root_repo = ctx.name,
            archive_name = archive_name,
        )
    ctx.file(
        "BUILD.bazel",
        _BUILD.format(
            repos = repos,
        ),
    )
    return ctx.repo_metadata(
        reproducible = True
    )

binary_toolchain_root_repo = repository_rule(
    implementation = _impl,
    doc = "Root binary toolchain repository",
    attrs = {
        "lock": attr.label(
            mandatory = True,
            doc = "Lock file label",
        ),
    },
)
