"""Hugo build rule backed by the Hugo persistent worker."""

load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _al_hugo_worker_impl(ctx):
    hugo = ctx.toolchains["//pkg/bzl:toolchain_type"]
    site_info = ctx.attr.site[AlHugoSiteInfo]
    out_dir = ctx.actions.declare_directory("{}.dest".format(ctx.label.name))
    flagfile = ctx.actions.declare_file("{}.flagfile.json".format(ctx.label.name))
    ctx.actions.write(
        output = flagfile,
        content = json.encode_indent({
            "site_archive": site_info.site_archive.path,
            "arguments": [ctx.expand_location(arg) for arg in ctx.attr.arguments],
            "hugo": hugo.hugo.path,
            "postcss": site_info.postcss.path,
            "postcss_bindir": site_info.postcss.root.path,
            "tool_dirs": site_info.tool_dirs,
            "env": {
                key: value
                for key, value in site_info.env.items()
                if key != "PATH"
            },
            "output_dir": out_dir.path,
        }),
    )

    ctx.actions.run(
        mnemonic = "HugoSite",
        executable = ctx.executable.worker,
        arguments = ["run", "--flagfile={}".format(flagfile.path)],
        inputs = depset(
            direct = [
                flagfile,
                hugo.hugo,
                site_info.postcss,
                site_info.site_archive,
            ],
            transitive = [
                hugo.hugo_target[DefaultInfo].files,
                hugo.hugo_target[DefaultInfo].default_runfiles.files,
                ctx.attr.site[DefaultInfo].default_runfiles.files,
            ],
        ),
        outputs = [out_dir],
        tools = [
            hugo.hugo_target[DefaultInfo].files_to_run,
            hugo.hugo_target[DefaultInfo].files,
            site_info.postcss_files_to_run,
        ],
        execution_requirements = {
            "supports-workers": "1",
            "requires-worker-protocol": "proto",
        },
        progress_message = "Hugo site (worker): %{label}",
    )

    return [
        DefaultInfo(
            files = depset([out_dir]),
        ),
    ]

al_hugo_worker = rule(
    implementation = _al_hugo_worker_impl,
    doc = "Build a Hugo site with the Hugo persistent worker",
    toolchains = [
        "//pkg/bzl:toolchain_type",
    ],
    attrs = {
        "out_dir": attr.string(
            mandatory = True,
            doc = "Stable destination name retained for API compatibility",
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
        "worker": attr.label(
            default = "//cmd/hugo_worker",
            executable = True,
            cfg = "exec",
            doc = "Hugo worker binary",
        ),
    },
)
