load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    files = ctx.files.site
    transitive_files = []
    for tool in ctx.attr.tools:
        transitive_files.append(tool[DefaultInfo].default_runfiles.files)
        transitive_files.append(tool[DefaultInfo].files)
    transitive_files.append(ctx.attr.postcss[DefaultInfo].default_runfiles.files)
    runfiles = ctx.runfiles(
        files = files,
        transitive_files = depset(transitive = transitive_files),
    )

    path_entries = [
        "$${{PWD}}/$$(dirname '{}')".format(ctx.executable.postcss.short_path),
        "$${{PWD}}/$$(dirname '{}')".format(ctx.executable.postcss.path),
        "$${{PWD}}/$${{0}}.runfiles/{}/$$(dirname '{}')".format(ctx.workspace_name, ctx.executable.postcss.short_path),
    ]
    for tool in ctx.attr.tools:
        executable = tool[DefaultInfo].files_to_run.executable
        if executable != None:
            path_entries.append("$${{PWD}}/$$(dirname '{}')".format(executable.short_path))
            path_entries.append("$${{PWD}}/$${{0}}.runfiles/{}/$$(dirname '{}')".format(ctx.workspace_name, executable.short_path))
        for f in tool[DefaultInfo].files.to_list():
            path_entries.append("$${{PWD}}/$$(dirname '{}')".format(f.short_path))
            path_entries.append("$${{PWD}}/$${{0}}.runfiles/{}/$$(dirname '{}')".format(ctx.workspace_name, f.short_path))
            path_entries.append("$${{PWD}}/$$(dirname '{}')".format(f.path))
            path_entries.append("$$(dirname '{}')".format(f.path))
        for f in tool[DefaultInfo].default_runfiles.files.to_list():
            path_entries.append("$${{PWD}}/$$(dirname '{}')".format(f.short_path))
            path_entries.append("$${{PWD}}/$${{0}}.runfiles/{}/$$(dirname '{}')".format(ctx.workspace_name, f.short_path))
            path_entries.append("$${{PWD}}/$$(dirname '{}')".format(f.path))
            path_entries.append("$$(dirname '{}')".format(f.path))
    path = ":".join(path_entries + [
        ctx.attr.env.get("PATH", ""),
        "$${PATH}",
    ])
    has_sass = False
    for tool in ctx.attr.tools:
        for f in tool[DefaultInfo].files.to_list() + tool[DefaultInfo].default_runfiles.files.to_list():
            if f.basename == "sass":
                has_sass = True
                break
        if has_sass:
            break
    env = {}
    env.update(ctx.attr.env)
    env["PATH"] = path
    env["BAZEL_BINDIR"] = "$${{PWD}}/{}".format(
        ctx.executable.postcss.root.path,
    )
    if has_sass:
        env["DART_SASS_BINARY"] = "sass"
    env_script = " && ".join([
        cmd
        for name, value in env.items()
        for cmd in ['{}="{}"'.format(name, value), "export {}".format(name)]
    ])

    return [
        DefaultInfo(
            files = depset(files),
            runfiles = runfiles,
        ),
        AlHugoSiteInfo(
            site_archive = ctx.file.site,
            postcss = ctx.executable.postcss,
            postcss_files_to_run = ctx.attr.postcss[DefaultInfo].files_to_run,
            env = ctx.attr.env,
            env_script = env_script,
        ),
    ]

al_hugo_site = rule(
    implementation = _impl,
    doc = "Define a hugo site",
    provides = [AlHugoSiteInfo],
    attrs = {
        "site": attr.label(
            mandatory = True,
            allow_single_file = [".tar"],
            doc = "Hugo site archive",
        ),
        "postcss": attr.label(
            mandatory = True,
            doc = "Postcss target",
            executable = True,
            cfg = "exec",
        ),
        "tools": attr.label_list(
            doc = "Tools that should be available for the build",
            default = [],
        ),
        "env": attr.string_dict(
            default = {},
            doc = """
                Hugo environment variables
                (support location statements, support make variables, support
                shell commands)
            """,
        ),
    },
)
