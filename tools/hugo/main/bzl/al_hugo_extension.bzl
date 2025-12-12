load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load(":al_hugo_archives.bzl", "AL_HUGO_ARCHIVES")

_BUILD = """
load("@rules_pkg//pkg:mappings.bzl", "pkg_files", "strip_prefix")

alias(
    name = "pkg_files",
    actual = "{name}",
    visibility = ["//visibility:public"],
)

pkg_files(
    name = "{name}",
    srcs = glob(["**"]),
    prefix = "{prefix}",
    strip_prefix = strip_prefix.from_pkg(),
    visibility = ["//visibility:public"],
)
"""

_BUILD_REPO = """
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup")

alias(
    name = "themes",
    actual = "{name}",
    visibility = ["//visibility:public"],
)

pkg_filegroup(
    name = "{name}",
    srcs = {srcs},
    prefix = "themes",
    visibility = ["//visibility:public"],
)
"""

def _repo_impl(ctx):
    ctx.file(
        "BUILD.bazel",
        _BUILD_REPO.format(
            name = ctx.original_name,
            srcs = json.encode([str(val) for val in ctx.attr.themes]),
        ),
    )
    return ctx.repo_metadata(
        reproducible = True,
    )

_repo = repository_rule(
    implementation = _repo_impl,
    doc = "Hugo repository",
    attrs = {
        "themes": attr.label_list(
            default = [],
            doc = "List of hugo themes",
        ),
    },
)

_BUILD_TOOLCHAIN = """
load("@com_alwaldend_src//tools/hugo/main/bzl:al_hugo_toolchain.bzl", "al_hugo_toolchain")
load("@bazel_skylib//rules:native_binary.bzl", "native_binary")

native_binary(
    name = "bin",
    src = "hugo",
    exec_compatible_with = {platforms},
    visibility = ["//visibility:public"],
)

alias(
    name = "{name}",
    actual = ":toolchain",
    visibility = ["//visibility:public"],
)

toolchain(
    name = "toolchain",
    toolchain = "toolchain_impl",
    toolchain_type = "@com_alwaldend_src//tools/hugo/main/bzl:toolchain_type",
    exec_compatible_with = {platforms},
    visibility = ["//visibility:public"],
)

al_hugo_toolchain(
    name = "toolchain_impl",
    exec_compatible_with = {platforms},
    hugo = ":bin",
)
"""

def _impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.toolchains:
            for archive in AL_HUGO_ARCHIVES.get(tag.version, []):
                name = "_".join([tag.name] + archive["platforms"]).replace(":", "_")
                root_module_direct_deps.append(name)
                http_archive(
                    name = name,
                    url = archive["url"],
                    strip_prefix = archive.get("strip_prefix"),
                    build_file_content = _BUILD_TOOLCHAIN.format(
                        name = name,
                        platforms = ["@platforms//{}".format(platform) for platform in archive["platforms"]],
                    ),
                    integrity = archive.get("integrity", ""),
                )

        for tag in mod.tags.lock:
            lock = json.decode(ctx.read(tag.lock))
            names = []
            for data in lock.get("remote_themes", []):
                name = "{}_{}".format(tag.name, data["path"].replace("/", "_").replace("-", "_").replace(".", "_").lower())
                root_module_direct_deps.append(name)
                names.append("@{}".format(name))
                http_archive(
                    name = name,
                    url = data["url"],
                    strip_prefix = data["strip_prefix"],
                    build_file_content = _BUILD.format(prefix = data["path"], name = name),
                    integrity = data["integrity"],
                )
            _repo(
                name = tag.name,
                themes = names,
            )
            root_module_direct_deps.append(tag.name)
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

al_hugo_extension = module_extension(
    implementation = _impl,
    doc = "Hugo extension",
    tag_classes = {
        "lock": tag_class(
            {
                "name": attr.string(
                    mandatory = True,
                    doc = "Name prefix",
                ),
                "lock": attr.label(
                    doc = "Hugo lock path",
                    mandatory = True,
                ),
            },
        ),
        "toolchains": tag_class(
            {
                "name": attr.string(
                    mandatory = True,
                    doc = "Name",
                ),
                "version": attr.string(
                    mandatory = True,
                    doc = "Release version",
                ),
            },
        ),
    },
)
