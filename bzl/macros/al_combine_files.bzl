def al_combine_files(name, srcs, **kwargs):
    """
    Create a genrule combining several files into one

    Args:
        name: genrule target
        srcs: list of labels to combine
        **kwargs: other genrule kwargs
    """
    native.genrule(
        name = name,
        srcs = srcs,
        outs = ["{}-output".format(name)],
        cmd = "cat $(SRCS) >$(@)",
        **kwargs
    )
