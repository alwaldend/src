AlMakeInstallFilegroupInfo = provider(
    doc = "Describe install info for a single filegrop",
    fields = {
        "filegroup": "PackageFilegroupInfo",
        "install_dir": "Install directory",
        "origin": "Rule label",
        "install_args": "Args for the install command",
        "diff_args": "Args for the diff command",
    },
)
