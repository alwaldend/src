"""Root-workspace adapters for project-owned DNS in nested Bazel modules."""

load("//projects/al/rules/al:al_config.bzl", "al_config")
load("//tools/terraform:defs.bzl", "terraform_binary_map", "terraform_test_map")
load("//tools/vault:defs.bzl", "vault_binary_map")

def _project_files_impl(ctx):
    symlinks = {}
    runfiles = ctx.runfiles()
    for target, directory in ctx.attr.directories.items():
        if not directory or directory.startswith("/") or ".." in directory.split("/"):
            fail("runfiles directories must be nonempty relative paths: " + directory)
        for source in target[DefaultInfo].files.to_list():
            path = directory + "/" + source.basename
            if path in symlinks:
                fail("duplicate project runfiles path: " + path)
            symlinks[path] = source
        runfiles = runfiles.merge(target[DefaultInfo].default_runfiles)
    return [DefaultInfo(runfiles = runfiles.merge(ctx.runfiles(symlinks = symlinks)))]

project_files = rule(
    implementation = _project_files_impl,
    doc = "Expose flat canonical source bundles at their project checkout paths.",
    attrs = {
        "directories": attr.label_keyed_string_dict(
            allow_files = True,
            doc = "Source filegroups mapped to logical project directories.",
        ),
    },
)

def project_dns(name):
    """Package a nested project's own root without a monorepo module dependency.

    Args:
        name: The nested Bazel module and owning projects/ directory name.
    """
    project = "projects/" + name
    al_config(
        name = name + "_al",
        srcs = ["@{}//:dns_al_source".format(name)],
        deps = ["//infra:al"],
    )
    project_files(
        name = name + "_files",
        directories = {
            "@{}//:dnsconfig".format(name): project,
            "@{}//tf:sources".format(name): project + "/tf",
        },
    )
    sources = [
        ":" + name + "_files",
        "//projects/tf_modules/dns_records/global:global",
    ]
    terraform_binary_map(
        name = name + "_tf",
        chdir = project + "/tf",
        configs = [":" + name + "_al"],
        data = sources + [
            "//tools/vault/injector",
            "//tools/vault/tf_backend",
        ],
        run_args = ["--plugin_label", "tf=main"],
    )
    terraform_test_map(
        name = name + "_tf_tests",
        chdir = project + "/tf",
        data = sources,
    )
    vault_binary_map(
        name = name + "_vault",
        configs = [":" + name + "_al"],
    )
