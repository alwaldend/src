load("//tools/release/main/bzl:al_release_files.bzl", "al_release_files")

def al_release_deps(name, srcs, visibility = None, **kwargs):
    """
    Generate a dependency diagram using genquery

    Args:
        name: name
        srcs: list of labels to generate deps for (should be full labels)
        visibility: visibility
        **kwargs: kwargs for al_release_files
    """
    native.genquery(
        name = "{}.dot".format(name),
        expression = " union ".join(["deps({})".format(src) for src in srcs]),
        opts = [
            "--output",
            "graph",
        ],
        scope = srcs,
    )

    al_release_files(
        name = name,
        srcs = [":{}.dot".format(name)],
        visibility = visibility,
        **kwargs
    )
