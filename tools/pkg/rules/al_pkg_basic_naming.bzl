load("@rules_pkg//pkg:providers.bzl", "PackageVariablesInfo")

def _basic_naming_impl(ctx):
    values = {}
    for dep in ctx.attr.deps:
        values.update(dep[PackageVariablesInfo].values)
    for key, value in ctx.var.items():
        values["VAR_{}".format(key)] = value
    return PackageVariablesInfo(values = values)

al_pkg_basic_naming = rule(
    implementation = _basic_naming_impl,
    doc = "Variables for @rules_pkg",
    attrs = {
        "deps": attr.label_list(
            doc = "Deps to merge",
            providers = [PackageVariablesInfo],
        ),
    },
)
