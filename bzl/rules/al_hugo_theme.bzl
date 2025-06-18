load("//bzl/providers:al_hugo_theme_info.bzl", "AlHugoThemeInfo")

def _impl(ctx):
    srcs = [
        [src, file, is_node]
        for attr, is_node in [[ctx.attr.srcs, False], [ctx.attr.node_modules, True]]
        for src in attr
        for file in src[DefaultInfo].files.to_list()
    ]
    theme_name = ctx.label.name
    files = []
    output = []
    for src, file, is_node in srcs:
        root = src.label.workspace_root
        relative_path = file.path
        if root:
            relative_path = relative_path.replace("{}/".format(root), "")
        if is_node and "node_modules/" in relative_path:
            relative_path = "node_modules/{}".format(relative_path.split("node_modules/", 1)[-1])
        relative_path = "{}/{}".format(theme_name, relative_path)
        output_file = ctx.actions.declare_symlink(relative_path)
        ctx.actions.symlink(output = output_file, target_path = file.path)
        files.append(output_file)

    return [
        DefaultInfo(
            files = depset(files, transitive = [src[DefaultInfo].files for src in ctx.attr.srcs]),
        ),
        AlHugoThemeInfo(
            files = depset(files),
            name = theme_name,
        ),
    ]

al_hugo_theme = rule(
    implementation = _impl,
    doc = "Define a hugo theme",
    attrs = {
        "srcs": attr.label_list(
            default = [],
            doc = "Theme sources",
            allow_files = True,
        ),
        "node_modules": attr.label_list(
            default = [],
            doc = "Node modules for the theme",
        ),
    },
)
