"""Behavioral analysis fixtures for binary lock metadata."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")
load("//main/bzl:binary_toolchain_lock.bzl", "binary_toolchain_lock_parse")

ParsedLockInfo = provider(doc = "Parsed lock metadata exposed by the analysis fixture.", fields = ["lock"])

def _fixture_impl(ctx):
    return [ParsedLockInfo(lock = binary_toolchain_lock_parse(json.decode(ctx.attr.lock)))]

lock_fixture = rule(implementation = _fixture_impl, attrs = {"lock": attr.string()})

def _failure_impl(ctx):
    env = analysistest.begin(ctx)
    asserts.expect_failure(env, "unknown field")
    return analysistest.end(env)

lock_failure_test = analysistest.make(_failure_impl, expect_failure = True)

def _constraints_impl(ctx):
    env = analysistest.begin(ctx)
    lock = analysistest.target_under_test(env)[ParsedLockInfo].lock
    asserts.equals(env, ["@platforms//os:linux", "@platforms//cpu:x86_64"], lock.archives[0].toolchain["exec_compatible_with"])
    return analysistest.end(env)

lock_constraints_test = analysistest.make(_constraints_impl)

def schema_tests():
    """Declare executable analysis checks for supported and malformed locks."""
    archive = {
        "archive_name": "linux_x86_64",
        "toolchain_name": "fixture",
        "toolchain_version": "1",
        "binaries": [{"name": "fixture", "path": "bin/fixture"}],
        "actions": [],
        "toolchain": {"exec_compatible_with": ["@platforms//os:linux", "@platforms//cpu:x86_64"]},
    }
    lock = {"toolchains": [{"name": "fixture", "version": "1"}], "archives": [archive]}
    lock_fixture(name = "valid", lock = json.encode(lock))
    lock_constraints_test(name = "constraints_test", target_under_test = ":valid")
    cases = {
        "top_level": lock | {"unknown": True},
        "toolchain": lock | {"toolchains": [{"name": "fixture", "version": "1", "unknown": True}]},
        "obsolete_rules": lock | {"archives": [archive | {"rules": {"exec_compatible_with": []}}]},
        "constraint_field": lock | {"archives": [archive | {"toolchain": {"exec_compatibile_with": []}}]},
        "binary": lock | {"archives": [archive | {"binaries": [{"name": "fixture", "path": "bin/fixture", "unknown": True}]}]},
        "action": lock | {"archives": [archive | {"actions": [{"unknown": {}}]}]},
    }
    for name, value in cases.items():
        lock_fixture(name = name, lock = json.encode(value), tags = ["manual"])
        lock_failure_test(name = name + "_test", target_under_test = ":" + name)
