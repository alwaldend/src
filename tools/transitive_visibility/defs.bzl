"""Checks of ordinary artifact dependencies, complementary to package visibility."""

load("@bazel_skylib//rules:diff_test.bzl", "diff_test")

def artifact_boundary_test(name, target):
    """Reject tools, infrastructure and user artifacts reachable from a runtime root.

    Execution-tool and implicit dependency edges are excluded. The query is
    unconfigured, so select branches remain part of its conservative coverage.
    """
    target = str(native.package_relative_label(target))
    native.genquery(
        name = name + "_query",
        expression = "filter('^//(tools|infra|users)(/|:)', deps({}))".format(target),
        opts = ["--noimplicit_deps", "--notool_deps", "--nonodep_deps"],
        scope = [target],
        testonly = True,
    )
    native.genquery(
        name = name + "_empty",
        expression = "set()",
        scope = [],
        testonly = True,
    )
    diff_test(
        name = name,
        file1 = ":" + name + "_empty",
        file2 = ":" + name + "_query",
    )
