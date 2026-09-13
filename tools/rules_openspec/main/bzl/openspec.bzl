"""Validate owner-local OpenSpec workspaces with the root-pinned CLI."""

load("@aspect_rules_js//js:defs.bzl", "js_library")
load("@rules_openspec_npm//:@fission-ai/openspec/package_json.bzl", openspec = "bin")

OPENSPEC_ENV = {
    "DO_NOT_TRACK": "1",
    "OPENSPEC_NO_UPDATE_CHECK": "1",
    "OPENSPEC_TELEMETRY": "0",
}

def openspec_validation(name, source, archived = False):
    """Declare sandboxed validation for one owner workspace.

    Args:
        name: Test target name.
        source: Filegroup in the owner's openspec package, including config,
            baseline specs, and complete native active and archived artifacts.
        archived: Check archived task completion instead of current artifacts.
    """
    source = native.package_relative_label(source)
    if source.package != "openspec" and not source.package.endswith("/openspec"):
        fail("OpenSpec source must belong to an openspec package: {}".format(source))

    # Source belongs to another package, including standalone nested modules.
    # Preserve those declared source runfiles instead of copying their paths
    # into this package or adding a JavaScript dependency to each nested module.
    js_library(
        name = name + "_source",
        testonly = True,
        data = [source],
        no_copy_to_bin = [source],
    )

    owner = source.package[:-len("openspec")].rstrip("/")
    repository = source.workspace_root[len("external/"):]
    chdir = "../{}/{}".format(repository, owner) if repository else (owner or ".")
    openspec.openspec_test(
        name = name,
        size = "small",
        args = [
            "validate",
            "--archived" if archived else "--all",
            "--strict",
            "--no-interactive",
        ],
        chdir = chdir,
        data = [":" + name + "_source"],
        env = OPENSPEC_ENV | {
            "XDG_CONFIG_HOME": "$${TEST_TMPDIR}/config",
            "XDG_DATA_HOME": "$${TEST_TMPDIR}/data",
        },
    )
