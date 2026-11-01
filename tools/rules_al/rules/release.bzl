load("@rules_pkg//pkg:providers.bzl", "PackageFilesInfo")
load("//rules:git_current_state.bzl", "GIT_CURRENT_STATE")
load(":release_files_info.bzl", "ReleaseFilesInfo")
load(":release_info.bzl", "ReleaseInfo")

def _impl(ctx):
    doc = ctx.actions.declare_file("{}.index.md".format(ctx.label.name))
    release_page = ctx.actions.declare_file("{}.release_page.json".format(ctx.label.name))
    dest_src_map = {
        "index.md": doc,
        "release_page.json": release_page,
    }
    release_name = ctx.attr.release_name
    transitive = []
    manifest = None
    files = {}

    if ctx.attr.manifest:
        manifest = ctx.file.manifest
    elif ctx.attr.srcs:
        manifest = ctx.actions.declare_file("{}.release.json".format(ctx.label.name))
        extra_manifest = ctx.actions.declare_file("{}.release_extra.json".format(ctx.label.name))
        inputs = ctx.files.git_state + [extra_manifest]
        if ctx.file.git_bundle:
            inputs.append(ctx.file.git_bundle)
        transitive = [dep[DefaultInfo].files for dep in ctx.attr.srcs]
        for src in ctx.files.srcs:
            dest_src_map[src.basename] = src
        args = ctx.actions.args()
        args.add("release")
        args.add("generate")
        args.add_all(["--output_manifest", manifest.path])
        args.add_all(["--output_release_page", release_page.path])
        for src in ctx.attr.srcs:
            src_manifest = src[ReleaseFilesInfo].manifest
            files.update(src[ReleaseFilesInfo].files)
            inputs.append(src_manifest)
            args.add_all(["--merge_manifest", src_manifest.path])
        args.add_all(["--merge_manifest", extra_manifest.path])
        if ctx.file.git_bundle:
            args.add_all(["--git_bundle", ctx.file.git_bundle.path])
        metadata = {"name": ctx.attr.release_name, "project": {"subdir": ctx.attr.project}}
        if ctx.file.git_bundle:
            metadata["git"] = {"revision": {"hash": "HEAD"}}
        ctx.actions.write(
            output = extra_manifest,
            content = json.encode(metadata),
        )

        ctx.actions.run(
            executable = ctx.executable.release_tool,
            arguments = [args],
            inputs = inputs,
            outputs = [manifest, release_page],
            # The bundle excludes blobs but still contains private Git metadata.
            # Keep it in the persistent local action cache only.
            execution_requirements = {"no-remote": "1"},
        )
    else:
        fail("missing both manifest and srcs")

    if not release_name:
        fail("missing release_name")

    ctx.actions.write(
        output = doc,
        content = ctx.attr.page_template.format(release_name = release_name),
    )
    dest_src_map["release.json"] = manifest
    dest_src_map = {
        "{}/{}/releases/{}/{}".format(ctx.attr.root_prefix, ctx.attr.project, release_name, name): file
        for name, file in dest_src_map.items()
    }

    return [
        DefaultInfo(
            files = depset(
                direct = dest_src_map.values(),
                transitive = transitive,
            ),
        ),
        ReleaseInfo(
            files = files,
            manifest = manifest,
            project = ctx.attr.project,
            release_name = ctx.attr.release_name,
        ),
        PackageFilesInfo(
            attributes = {},
            dest_src_map = dest_src_map,
        ),
    ]

release = rule(
    implementation = _impl,
    doc = "Rule describing a release",
    provides = [PackageFilesInfo, DefaultInfo, ReleaseInfo],
    attrs = {
        "srcs": attr.label_list(
            providers = [ReleaseFilesInfo],
            doc = "Sources",
        ),
        "release_name": attr.string(
            mandatory = True,
            doc = "Release name",
        ),
        "project": attr.string(
            mandatory = True,
            doc = "Project package_name",
        ),
        "page_template": attr.string(default = "# Release {release_name}\n", doc = "Release page template with a release_name placeholder."),
        "root_prefix": attr.string(
            default = "content/docs",
            doc = "Root prefix",
        ),
        "manifest": attr.label(
            doc = "Load manifest from a file instead of generating it from srcs",
        ),
        "release_tool": attr.label(
            executable = True,
            cfg = "exec",
            default = Label("//cmd/al:al"),
            doc = "Release tool",
        ),
        "git_state": attr.label_list(
            default = GIT_CURRENT_STATE,
            doc = "Additional declared Git state used to generate release metadata",
        ),
        "git_bundle": attr.label(
            allow_single_file = True,
            doc = "Deterministic Git bundle used to generate release metadata",
        ),
    },
)
