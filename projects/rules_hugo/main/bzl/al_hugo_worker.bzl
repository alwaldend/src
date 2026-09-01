"""Hugo build rule backed by the Hugo persistent worker."""

load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _al_hugo_worker_impl(ctx):
    hugo = ctx.toolchains["//main/bzl:toolchain_type"]
    site_info = ctx.attr.site[AlHugoSiteInfo]
    out_dir = ctx.actions.declare_directory(ctx.attr.out_dir)
    env_script = ctx.expand_make_variables(
        "env_script",
        ctx.expand_location(site_info.env_script),
        {},
    )

    flagfile = ctx.actions.declare_file("{}.flagfile.json".format(ctx.label.name))
    ctx.actions.write(
        output = flagfile,
        content = json.encode_indent({
            "site_archive": site_info.site_archive.path,
            "arguments": [ctx.expand_location(arg) for arg in ctx.attr.arguments],
            "env_script": env_script,
            "hugo": hugo.hugo.path,
            "env_file": hugo.env_file.path,
            "postcss": site_info.postcss.path,
            "postcss_bindir": site_info.postcss.root.path,
            "output_dir": out_dir.path,
            "shell": "/bin/sh",
        }),
    )

    ctx.actions.run(
        mnemonic = "HugoSite",
        executable = ctx.executable.worker,
        arguments = ["run", "--flagfile={}".format(flagfile.path)],
        inputs = depset(
            direct = [
                flagfile,
                hugo.env_file,
                site_info.postcss,
                site_info.site_archive,
            ],
            transitive = [
                ctx.attr.site[DefaultInfo].default_runfiles.files,
            ],
        ),
        outputs = [out_dir],
        tools = [
            hugo.hugo_target[DefaultInfo].files_to_run,
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
        "//main/bzl:toolchain_type",
    ],
    attrs = {
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
        "worker": attr.label(
            mandatory = True,
            executable = True,
            cfg = "exec",
            doc = "Hugo worker binary",
        ),
    },
)
