def al_txt_data(name, srcs, **kwargs):
    """
    Text data

    Args:
       name (str): target name
       srcs (list[str]): sources
       **kwargs: filegroup kwargs
    """
    native.filegroup(
        name = name,
        srcs = srcs,
        **kwargs
    )
