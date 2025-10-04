def al_genrule_src(name, srcs = [], visibility = ["//visibility:public"]):
    """
    Create a filegroup and a genrule generating a tar archive

    Args:
        name: genrule name
        srcs: source labels
    """
    filegroup = "{}-filegroup".format(name)
    native.filegroup(
        name = filegroup,
        srcs = srcs,
    )
    native.genrule(
        name = name,
        srcs = srcs,
        outs = ["{}.tar".format(name)],
        cmd = "tar -cf ${@} ${<}",
        visibility = visibility,
    )
