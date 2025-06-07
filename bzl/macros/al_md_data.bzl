def al_md_data(name, srcs, visibility = ["//visibility:public"], **kwargs):
    """
    Markdowd data backed by a filegroup

    Args:
        name: filegroup name
        srcs: markdown files
        visibility: visibility
        **kwargs: kwargs

    Targets:
        ${name}: filegroup
    """
    native.filegroup(
        name = name,
        srcs = srcs,
        visibility = visibility,
    )
