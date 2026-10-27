"""Build the complete search index before rendering the visitor-facing site."""

load("@rules_hugo//pkg/bzl:al_hugo_site.bzl", "al_hugo_site")
load("@rules_hugo//pkg/bzl:al_hugo_worker.bzl", "al_hugo_worker")
load("@rules_pkg//pkg:mappings.bzl", "REMOVE_BASE_DIRECTORY", "pkg_files")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

def _compile_index_impl(ctx):
    output = ctx.actions.declare_directory(ctx.label.name + ".source")
    args = ctx.actions.args()
    args.add("--documents", ctx.file.documents.path)
    args.add("--output", output.path)
    ctx.actions.run(
        executable = ctx.executable._compiler,
        arguments = [args],
        env = {
            "BAZEL_BINDIR": ctx.bin_dir.path,
            "JS_BINARY__NO_CD_BINDIR": "1",
        },
        inputs = [ctx.file.documents],
        outputs = [output],
        mnemonic = "CompileSiteSearchIndex",
    )
    return [DefaultInfo(files = depset([output]))]

_compile_index = rule(
    implementation = _compile_index_impl,
    attrs = {
        "documents": attr.label(allow_single_file = True, mandatory = True),
        "_compiler": attr.label(
            default = "//projects/alwaldend.com/cmd/search_index",
            cfg = "exec",
            executable = True,
        ),
    },
)

def indexed_hugo_site(name, source, arguments, visibility = None):
    """Render matching search and site content with the same build flags."""
    tool_args = {
        "postcss": "//tools/postcss",
        "tools": ["//third_party/com_github_sass_dart_sass:dart-sass"],
    }
    al_hugo_site(name = name + "_document_source", site = source, **tool_args)
    al_hugo_worker(
        name = name + "_documents",
        site = ":" + name + "_document_source",
        arguments = arguments + ["--renderSegments", "search"],
        out_dir = name + ".documents",
    )
    _compile_index(name = name + "_index", documents = ":" + name + "_documents")
    pkg_files(
        name = name + "_index_assets",
        srcs = [":" + name + "_index"],
        renames = {":" + name + "_index": REMOVE_BASE_DIRECTORY},
    )
    pkg_tar(
        name = name + "_indexed_archive",
        deps = [source],
        srcs = [":" + name + "_index_assets"],
    )
    al_hugo_site(name = name + "_site", site = ":" + name + "_indexed_archive", **tool_args)
    al_hugo_worker(
        name = name,
        site = ":" + name + "_site",
        arguments = arguments + ["--renderSegments", "site"],
        out_dir = name + ".destination",
        visibility = visibility,
    )
