load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("//tools/bzl/main/bzl:al_bzl_generate_repository.bzl", "al_bzl_generate_repository")
load(":al_binary_toolchain_repo.bzl", "al_binary_toolchain_repo")

_BUILD_ROOT = """
load("@com_alwaldend_src//tools/trufflehog/main/bzl:al_trufflehog_toolchain.bzl", "al_trufflehog_toolchain")
CHILDREN = {children}

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
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.toolchain:
            lock = json.decode(ctx.file(tag.lock))
            version = lock["version"]
            archives = lock["archives"].get(version, [])
            if not archives:
                fail("missing archives for version '{}'".format(version))
            for archive_name, archive in archives.items():
                name = "{}_{}".format(tag.name, archive_name)
                root_module_direct_deps.append(name)
                al_binary_toolchain_repo(
                    name = name,
                    lock = tag.lock,
                    version = version,
                    archive_name = archive_name,
                )
            root_module_direct_deps.append(tag.name)
            al_bzl_generate_repository(
                name = tag.name,
                files = {
                    "BUILD.bazel": _BUILD_ROOT.format(children = []),
                },
            )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_binary_toolchain_extension = module_extension(
    implementation = _impl,
    doc = "A simple binary toolchain extension",
    tag_classes = {
        "toolchain": tag_class(
            {
                "name": attr.string(
                    mandatory = True,
                    doc = "Repo name",
                ),
                "lock": attr.string(
                    mandatory = True,
                    doc = "Toolchain lock",
                ),
            },
        ),
    },
)
