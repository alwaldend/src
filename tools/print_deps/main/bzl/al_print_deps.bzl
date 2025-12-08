def _describe(name, obj, exclude = []):
    """
    Print the properties of the given struct obj

    Args:
      name: the name of the struct we are introspecting.
      obj: the struct to introspect
      exclude: a list of names *not* to print (function names)
    """
    for k in dir(obj):
        if hasattr(obj, k) and k not in exclude:
            v = getattr(obj, k)
            t = type(v)
            print("%s.%s<%r> = %s" % (name, k, t, v))

def _print_deps_impl(target, ctx):
    if ctx.label.repo_name or not hasattr(ctx.rule.attr, "srcs"):
        return []

    print("LABEL: {}, {}".format(ctx.label, ctx.label.repo_name))

    # _describe("target", target)
    # _describe("ctx", ctx, ["build_setting_value", "outputs", "split_attr"])
    for src in ctx.rule.files.srcs:
        print("SRC: {}".format(src))

    return []

al_print_deps = aspect(
    implementation = _print_deps_impl,
    attr_aspects = ["deps"],
)
