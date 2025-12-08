AlMakeInstallFilegroupInfo = provider(
    doc = "Describe install info for a single filegrop",
    fields = {
        "srcs": "depset of PackageFilegroupInfo",
        "deps": "depset of AlMakeInstallFilegroupInfo",
        "install_dir": "Install directory",
        "origin": "Rule label",
        "install_args": "Args for the install command",
        "diff_args": "Args for the diff command",
        "pkg_prefix": "Ignore that prefix for srcs",
    },
)
