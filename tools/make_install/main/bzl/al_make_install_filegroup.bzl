load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo")
load("//tools/make_install/main/bzl:al_make_install_filegroup_info.bzl", "AlMakeInstallFilegroupInfo")

def _impl(ctx):
    return [
        DefaultInfo(
            files = depset(
                ctx.files.srcs,
                transitive = [dep[DefaultInfo].files for dep in ctx.attr.deps],
            ),
        ),
        AlMakeInstallFilegroupInfo(
            install_dir = ctx.attr.install_dir,
            origin = ctx.label,
            install_args = ctx.attr.install_args,
            diff_args = ctx.attr.diff_args,
            pkg_prefix = ctx.attr.pkg_prefix,
            srcs = depset([src[PackageFilegroupInfo] for src in ctx.attr.srcs]),
            deps = depset(transitive = [dep[AlMakeInstallFilegroupInfo] for dep in ctx.attr.deps]),
        ),
    ]

al_make_install_filegroup = rule(
    doc = "Create a make install filegroup",
    implementation = _impl,
    provides = [AlMakeInstallFilegroupInfo],
    attrs = {
        "srcs": attr.label_list(
            providers = [PackageFilegroupInfo],
            doc = "Sources to install",
        ),
        "install_args": attr.string_list(
            default = ["--compare", "-D", "--backup=numbered"],
            doc = "Install args",
        ),
        "diff_args": attr.string_list(
            default = [],
            doc = "Diff args",
        ),
        "pkg_prefix": attr.string(
            mandatory = True,
            doc = "Ignore that prefix for srcs",
        ),
        "deps": attr.label_list(
            providers = [AlMakeInstallFilegroupInfo],
            doc = "Deps",
        ),
        "install_dir": attr.string(
            default = "${HOME}",
            doc = "Install directory",
        ),
    },
)
