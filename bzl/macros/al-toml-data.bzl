def al_toml_data(name, srcs, visibility = ["//visibility:public"], **kwargs):
    """
    Create toml data targets

    Args:
        name: Target name
        srcs: Toml files
    """
    native.filegroup(name = name, srcs = srcs, visibility = visibility)
