load("@aspect_rules_js//js:providers.bzl", "JsInfo")
load("//bzl/providers:al_hugo_module_info.bzl", "AlHugoModuleInfo")
load("//bzl/providers:al_hugo_theme_info.bzl", "AlHugoThemeInfo")

def _impl(ctx):
    theme_name = ctx.attr.theme_name or ctx.label.name
    return [
        DefaultInfo(
            files = depset([dir]),
        ),
        AlHugoThemeInfo(
            name = theme_name,
            archive = ctx.file.src,
            modules = depset([module[AlHugoModuleInfo] for module in ctx.attr.modules]),
        ),
    ]

al_hugo_theme = rule(
    implementation = _impl,
    doc = "Define a hugo theme",
    attrs = {
        "theme_name": attr.string(
            doc = "Theme name (default: label name)",
        ),
        "modules": attr.label_list(
            doc = "Module list for this theme",
            default = [],
            providers = [AlHugoModuleInfo],
        ),
        "src": attr.label(
            doc = "Theme archive",
            allow_single_file = [".tar"],
        ),
    },
)
