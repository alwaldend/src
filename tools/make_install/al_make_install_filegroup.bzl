load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo")
load("//tools/make_install:al_make_install_filegroup_info.bzl", "AlMakeInstallFilegroupInfo")
load("//tools/make_install:al_make_install_info.bzl", "AlMakeInstallInfo")

def _impl(ctx):
    pkg_files = []
    pkg_dirs = []
    pkg_symlinks = []
    for src in ctx.attr.srcs:
        info = src[PackageFilegroupInfo]
        pkg_files.extend(info.pkg_files)
        pkg_dirs.extend(info.pkg_dirs)
        pkg_symlinks.extend(info.pkg_symlinks)
    filegroup = AlMakeInstallFilegroupInfo(
        install_dir = ctx.attr.install_dir,
        origin = ctx.label,
        install_args = ctx.attr.install_args,
        diff_args = ctx.attr.diff_args,
        pkg_files = depset(pkg_files),
        pkg_dirs = depset(pkg_dirs),
        pkg_symlinks = depset(pkg_symlinks),
    )
    return [
        DefaultInfo(
            files = depset(
                ctx.files.srcs,
                transitive = [dep[DefaultInfo].files for dep in ctx.attr.deps],
            ),
        ),
        PackageFilegroupInfo(
            pkg_files = pkg_files,
            pkg_dirs = pkg_dirs,
            pkg_symlinks = pkg_symlinks,
        ),
        AlMakeInstallInfo(
            filegroups = depset(
                [filegroup],
                transitive = [dep[AlMakeInstallInfo].filegroups for dep in ctx.attr.deps],
            ),
        ),
    ]

al_make_install_filegroup = rule(
    doc = "Create a make install filegroup",
    implementation = _impl,
    provides = [AlMakeInstallInfo, PackageFilegroupInfo],
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
        "deps": attr.label_list(
            providers = [AlMakeInstallInfo],
            doc = "Deps",
        ),
        "install_dir": attr.string(
            default = "${HOME}",
            doc = "Install directory",
        ),
    },
)
