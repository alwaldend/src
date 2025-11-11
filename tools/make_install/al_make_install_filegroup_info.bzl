AlMakeInstallFilegroupInfo = provider(
    doc = "Describe install info for a single filegrop",
    fields = {
        "pkg_files": "pkg_files from PackageFilegroupInfo",
        "pkg_dirs": "pkg_dirs from PackageFilegroupInfo",
        "pkg_symlinks": "pkg_symlinks from PackageFilegroupInfo",
        "install_dir": "Install directory",
        "origin": "Rule label",
        "install_args": "Args for the install command",
        "diff_args": "Args for the diff command",
    },
)
