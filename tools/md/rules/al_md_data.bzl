def al_md_data(name, srcs, deps = [], **kwargs):
    """
    Markdown data backed by a filegroup

    Targets:
    - ${name}: filegroup

    Args:
        name (string): filegroup name
        srcs (list[str]): markdown files
        deps (list[str]): deps
        **kwargs: filegroup kwargs
    """
    native.filegroup(
        name = name,
        srcs = srcs + deps,
        **kwargs
    )
