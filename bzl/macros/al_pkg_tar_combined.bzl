def al_pkg_tar_combined(name, srcs = [], strip_components = 2, **kwargs):
    """
    Create a genrule combining several tars into one

    Args:
        name: genrule name
        srcs: dicts tar archives ({"label": "tar label", "dir": "target dir"})
        strip_components: value of --stip-components
        **kwargs: other genrule kwargs
    """
    cmd = "set -eux"
    out = "{}.tar".format(name)
    cmd += "\n".join(["""
        mkdir -p '{dir}'
        tar -xf $(execpath {label}) --strip-components '{strip_components}' -C '{dir}'
    """.format(strip_components = strip_components, **tar) for tar in srcs])
    cmd += """
        tar -cf $(execpath {output}) {dirs}
    """.format(output = out, dirs = " ".join(["'{}'".format(tar["dir"]) for tar in srcs]))
    native.genrule(
        name = name,
        srcs = [tar["label"] for tar in srcs],
        outs = [out],
        cmd = cmd,
        **kwargs
    )
