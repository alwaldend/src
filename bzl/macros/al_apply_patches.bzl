def al_apply_patches(
        name,
        src,
        patches,
        visibility = ["//visibility:public"],
        **kwargs):
    """
    Create a genrule applying patches

    Args:
        name: genrule name
        src: source archive label
        patches: patches label
        visibility: visibility
        **kwargs: other genrule kwargs
    """
    native.genrule(
        name = name,
        srcs = [src, patches],
        outs = ["{}.tar".format(name)],
        visibility = visibility,
        cmd = """
            set -eux
            mkdir _src _patches
            cp $(rootpaths {patches}) _patches/
            tar -xf "$(rootpath {src})" -C _src --strip-components 2
            cd _src
            git apply -v ../_patches/*
            cd ../
            tar -cf "$(@)" -C "_src"
        """.format(src = src, patches = patches),
    )
