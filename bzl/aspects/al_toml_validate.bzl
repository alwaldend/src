def _print_deps_impl(target, ctx):
    # Make sure the rule has a srcs attribute.
    if hasattr(ctx.rule.attr, "srcs"):
        # Iterate through the files that make up the sources and
        # print their paths.
        for src in ctx.rule.attr.srcs:
            for f in src.files.to_list():
                print(f.path)
    return []

al_print_deps = aspect(
    implementation = _print_deps_impl,
    attr_aspects = ["srcs"],
)
