load("@bazel_skylib//lib:shell.bzl", "shell")
load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    hugo = ctx.toolchains["//tools/hugo/main/bzl:toolchain_type"]
    site_info = ctx.attr.site[AlHugoSiteInfo]
    out_dir = ctx.actions.declare_directory(ctx.attr.out_dir)
    args = ctx.actions.args()
    for arg in ctx.attr.arguments:
        args.add(ctx.expand_location(arg))
    args.add("--destination")
    args.add(out_dir.path)
    env_script = ctx.expand_make_variables(
        "env_script",
        ctx.expand_location(site_info.env_script),
        {},
    )
    ctx.actions.run_shell(
        command = """\
            set -eu
            {env_script}
            tar -xf {site_archive}
            mkdir -p static
            cp {env_file} static/hugo_env.txt
            exec {hugo} "$@"
        """.format(
            env_file = shell.quote(hugo.env_file.path),
            env_script = env_script,
            hugo = shell.quote(hugo.hugo.path),
            site_archive = shell.quote(site_info.site_archive.path),
        ),
        inputs = depset(
            direct = [
                hugo.env_file,
                site_info.site_archive,
            ],
            transitive = [
                ctx.attr.site[DefaultInfo].default_runfiles.files,
            ],
        ),
        outputs = [out_dir] + ctx.outputs.outs,
        arguments = [args],
        progress_message = "Running hugo action: %{label}",
        tools = [hugo.hugo_target[DefaultInfo].files_to_run],
    )
    return [
        DefaultInfo(
            files = depset([out_dir] + ctx.outputs.outs),
        ),
    ]

al_hugo_run_binary = rule(
    implementation = _impl,
    doc = "Build a Hugo site with the registered Hugo toolchain",
    toolchains = [
        "//tools/hugo/main/bzl:toolchain_type",
    ],
    attrs = {
        "outs": attr.output_list(
            doc = "Output files",
        ),
        "out_dir": attr.string(
            mandatory = True,
            doc = "Hugo destination output directory",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Hugo arguments",
        ),
        "site": attr.label(
            mandatory = True,
            providers = [AlHugoSiteInfo],
            doc = "Hugo site",
        ),
    },
)
