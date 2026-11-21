"""Validate owner-local OpenSpec workspaces with the root-pinned CLI."""

load("@aspect_rules_js//js:defs.bzl", "js_library", "js_test")

OPENSPEC_ENTRY_POINT = Label("//tools/rules_openspec:openspec_entry_point")
OPENSPEC_NPM_PACKAGE = Label("//tools/rules_openspec:node_modules/@fission-ai/openspec")

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

    # Preserve the declared runfile paths, including for cross-package callers.
    js_library(
        name = name + "_source",
        testonly = True,
        data = [source],
        no_copy_to_bin = [source],
    )

    owner = source.package[:-len("openspec")].rstrip("/")
    repository = source.workspace_root[len("external/"):] or "_main"
    chdir = "$${TEST_SRCDIR}/" + repository + "/" + owner
    js_test(
        entry_point = OPENSPEC_ENTRY_POINT,
        name = name,
        size = "small",
        args = [
            "validate",
            "--archived" if archived else "--all",
            "--strict",
            "--no-interactive",
        ],
        chdir = chdir,
        data = [":" + name + "_source", OPENSPEC_NPM_PACKAGE],
        env = OPENSPEC_ENV | {
            "XDG_CONFIG_HOME": "$${TEST_TMPDIR}/config",
            "XDG_DATA_HOME": "$${TEST_TMPDIR}/data",
        },
    )

def openspec_validate_test(name, source):
    """Declare strict current and archived checks for a single OpenSpec workspace.

    Args:
        name: Validation suite target name.
        source: The workspace source target.
    """
    current_test = name + "_current_test"
    archive_test = name + "_archive_test"
    openspec_validation(name = current_test, source = source)
    openspec_validation(name = archive_test, source = source, archived = True)
    native.test_suite(
        name = name,
        tests = [":" + current_test, ":" + archive_test],
    )
