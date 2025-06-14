def al_md_data(name, srcs, visibility = ["//visibility:public"], **kwargs):
    """
    Markdown data backed by a filegroup

    Targets:
    - ${name}: filegroup

    Args:
        name: filegroup name
        srcs: markdown files
        visibility: visibility
        **kwargs: filegroup kwargs
    """
    native.filegroup(
        name = name,
        srcs = srcs,
        visibility = visibility,
        **kwargs
    )
