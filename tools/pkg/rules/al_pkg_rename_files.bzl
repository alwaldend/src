load("@bazel_skylib//lib:paths.bzl", "paths")
load("@rules_pkg//pkg:providers.bzl", "PackageDirsInfo", "PackageFilegroupInfo", "PackageFilesInfo", "PackageSymlinkInfo")

def _rename(dest, renames):
    """
    Args:
        dest: path string
        renames: dict of renames
    Returns:
        renamed dest
    """
    rename = renames.get(paths.basename(dest))
    if rename:
        dirname = paths.dirname(dest)
        if dirname:
            dest = "{}/{}".format(dirname, rename)
        else:
            dest = rename
    return dest

def _impl(ctx):
    pkg_files = []
    pkg_dirs = []
    pkg_symlinks = []
    info = ctx.attr.src[PackageFilegroupInfo]
    for files, origin in info.pkg_files:
        dest_src_map = {}
        for dest, source in files.dest_src_map.items():
            dest_src_map[_rename(dest, ctx.attr.basenames)] = source
        pkg_files.append(
            (
                PackageFilesInfo(
                    dest_src_map = dest_src_map,
                ),
                origin,
            ),
        )
    for dir_info, origin in info.pkg_dirs:
        dirs = []
        for dir in dir_info.dirs:
            dirs.append(_rename(dir, ctx.attr.basenames))
        pkg_dirs.append(
            (
                PackageDirsInfo(
                    attributes = dir_info.attributes,
                    dirs = dirs,
                ),
                origin,
            ),
        )
    for sym_info, origin in info.pkg_symlinks:
        destination = _rename(sym_info.destination, ctx.attr.basenames)
        pkg_symlinks.append(
            (
                PackageSymlinkInfo(
                    attributes = sym_info.attributes,
                    destination = destination,
                    target = sym_info.target,
                ),
                origin,
            ),
        )

    return [
        ctx.attr.src[DefaultInfo],
        PackageFilegroupInfo(
            pkg_files = pkg_files,
            pkg_dirs = pkg_dirs,
            pkg_symlinks = pkg_symlinks,
        ),
    ]

al_pkg_rename_files = rule(
    implementation = _impl,
    provides = [PackageFilegroupInfo],
    doc = "Rename files in PackageFilegroupInfo",
    attrs = {
        "src": attr.label(
            providers = [PackageFilegroupInfo],
            doc = "PackageFilegroupInfo",
        ),
        "basenames": attr.string_dict(
            doc = "Map of file names to their new names",
        ),
    },
)
