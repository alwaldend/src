AlTransitiveSources = provider(
    fields = ["transitive_sources"],
    doc = """
        Provide transitive sources

        - https://github.com/bazelbuild/examples/blob/main/rules/depsets/foo.bzl
        - https://stackoverflow.com/a/57699683
    """,
)

def al_transitive_sources(srcs, deps):
    """
    Obtain the source files for a target and its transitive dependencies.

    Args:
      srcs: a list of source files
      deps: a list of targets that are direct dependencies

    Returns:
      a collection of the transitive sources
    """
    return depset(
        srcs,
        transitive = [dep[AlTransitiveSources].transitive_sources for dep in deps if AlTransitiveSources in dep],
    )
