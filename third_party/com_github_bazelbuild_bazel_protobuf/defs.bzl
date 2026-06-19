load("@bazel_skylib//lib:paths.bzl", "paths")

VERSIONS = {
    "7c096ba03fed206d9e6fe3afccbfe5c91a71ec46": {
        "base_url": "https://raw.githubusercontent.com/bazelbuild/bazel/7c096ba03fed206d9e6fe3afccbfe5c91a71ec46/src/main/protobuf",
        "integrity": {
            "action_cache.proto": "sha256-nu/0jzV+wtwHYXoKMm/die9hCNa5O8nupVkuy3/4SbM=",
            "analysis_v2.proto": "sha256-LcPeuXjeCKTNa5W4JLTC0YAdv5RaszKHfEGUlbWekIw=",
            "bazel_flags.proto": "sha256-ruSXO1Gq2iBkRaoowPzBhxWVzWWU1j6PO+RlvLWHiQo=",
            "bazel_output_service.proto": "sha256-kaa38LSovd/NrP7AmMoqTtfHKtpgJGtFjDyaZpcxmdQ=",
            "bazel_output_service_rev2.proto": "sha256-R+LAaLpZRqkb03viy6tw/33whTPYY2XuXnX4e2g2MAI=",
            "build.proto": "sha256-tvn8YIdTXb8zoC/JiOD6mI/QYdggI20AkJ5T2joZ/DE=",
            "builtin.proto": "sha256-8Ilfl5OYypA16/qQG1WDqps3EdtVDpfgLsHz7XKDA/k=",
            "cache_salt.proto": "sha256-j4b+OsTjcR7iU31wHGbQ7AOhI4fckJJuq76tp5M5tX0=",
            "command_line.proto": "sha256-utTR/EO+ks16cP80FNTkx8PxhD7Vfsp124S0rjDH+jM=",
            "command_server.proto": "sha256-DL0UgNFqU44trIkaJ++yZIsnoxfXS0VCVWqxr05pm0o=",
            "crash_debugging.proto": "sha256-fTIyCps0ngY4zZk6ptoXAXXr3WURj4NeX1zEmmDeDmo=",
            "crosstool_config.proto": "sha256-4zEFfSXObox4FGkM9TkyK73BGGtXILCS51d1qwpBFYs=",
            "deps.proto": "sha256-uGHbzgQXffnktyBIdrLyfhj0DrbSCz3/7+zdK688/pI=",
            "desugar_deps.proto": "sha256-tzNRLkAsY2jBDs82H+UlsoXrtvvWAtH5VEnxV6jNRR4=",
            "execution_graph.proto": "sha256-0Zml8D0V3bYAApinbnIL4SoKPR+Xm83rlsBvvd04Osc=",
            "execution_statistics.proto": "sha256-oe3fCk1eeFejDZ2K1ciTd87qVFgLBaHEjh8VZ2vfPD0=",
            "extra_actions_base.proto": "sha256-V0TM9DL7AwBsVoPMlA53Sy/ikkM1LKLn/TYNM+ttpjM=",
            "failure_details.proto": "sha256-1zSVqz252TskO4g/vH73tPc8Kqzzl2fs4ZMXLff7oC4=",
            "file_invalidation_data.proto": "sha256-G/uNiOtmv95V9RckBNkZ2D8d/4SRTUmhvksf3cRXfc4=",
            "invocation_policy.proto": "sha256-vtAAjmn7xTe8kSdTu/I4ldnP/Uocl62/QfKQgidPUjs=",
            "java_compilation.proto": "sha256-mEgFbT7CSqDsstmW1xRiM6792GA6Seo5K9YV9BEMraA=",
            "memory_pressure.proto": "sha256-lMvzcmPALapLruADGp2qysjRjRHuyXKq0T0NZoFYCbY=",
            "option_filters.proto": "sha256-yOkqilvkVIT3BeSr5PKj1LXHfJlHe1ziKGrdh+1k8sw=",
            "project/buildable_unit.proto": "sha256-tXZ/KdhrHTVeqqO8Yw3vHASUGs1GERY6APk5oaohuGk=",
            "project/enforcement_policy.proto": "sha256-0ZS+XAvHaA9mKhZyRiN9G/FeeH9ybKDayy/he4nl+kw=",
            "project/project.proto": "sha256-cqc18ktp7p9c2NV7/0Wf4bvY5GtYoqb6MX2Gcx6XM9Y=",
            "remote_execution_log.proto": "sha256-vrr09lacqE9jBCvUsdjrcOS6Tczq6r20jMyDOQuCtLQ=",
            "remote_scrubbing.proto": "sha256-MV8HnERxPP7hjxPnsNzY4ieblrycup+ijVAVEl8Z2Ws=",
            "spawn.proto": "sha256-p1JIyFhqq8yTjJF3jP+20QcQGTJqjmvlcixO1O+d+ZA=",
            "stardoc_output.proto": "sha256-S9VmuO+ug66EzJZGfWgqd2DSobA+Y1ZC6qudXbdt2Aw=",
            "strategy_policy.proto": "sha256-AGjiMiO+aPVd3DbWORbvrpuPZPHLrlsCusYFm1sDLU0=",
            "test_status.proto": "sha256-7jxtqLf6SuBdabnl/4BUCnAbQsZ/NU6Dj/1VOcQx2rY=",
            "worker_protocol.proto": "sha256-QvnyJRiNHs3P2rGMyyN3+Z3U/1p54WhYP5LZ6o0M0VI=",
            "xcode_config.proto": "sha256-hGwrGMAePLq5mom6MqjIBhLarNGPjjXAJV4TXfhjCMs=",
        },
    },
}

CONTRACTS = [
    "action_cache.proto",
    # "analysis_v2.proto",
    "bazel_flags.proto",
    # "bazel_output_service.proto",
    # "bazel_output_service_rev2.proto",
    # "build.proto",
    "builtin.proto",
    "cache_salt.proto",
    # "command_line.proto",
    # "command_server.proto",
    "crash_debugging.proto",
    "crosstool_config.proto",
    "deps.proto",
    "desugar_deps.proto",
    "execution_graph.proto",
    "execution_statistics.proto",
    "extra_actions_base.proto",
    # "failure_details.proto",
    "file_invalidation_data.proto",
    # "invocation_policy.proto",
    "java_compilation.proto",
    "memory_pressure.proto",
    "option_filters.proto",
    # "remote_execution_log.proto",
    "remote_scrubbing.proto",
    # "spawn.proto",
    "stardoc_output.proto",
    "strategy_policy.proto",
    "test_status.proto",
    "worker_protocol.proto",
    "xcode_config.proto",
    # "project/buildable_unit.proto",
    "project/enforcement_policy.proto",
    # "project/project.proto",
]

_BUILD = """
load("@protobuf//bazel:java_proto_library.bzl", "java_proto_library")
load("@protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@protobuf//bazel:proto_library.bzl", "proto_library")
load("@rules_go//go:def.bzl", "go_library")
load("@rules_go//proto:def.bzl", "go_proto_library")
load("@rules_java//java:defs.bzl", "java_library")

exports_files(["{contract_file}"])

proto_library(
    name = "{contract_name}_proto",
    srcs = ["{contract_file}"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "{contract_name}_go_proto",
    importpath = "git.alwaldend.com/alwaldend/src/third_party/com_github_bazelbuild_bazel_protobuf/{contract_path}",
    proto = ":{contract_name}_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "{contract_name}",
    embed = [":{contract_name}_go_proto"],
    importpath = "git.alwaldend.com/alwaldend/src/third_party/com_github_bazelbuild_bazel_protobuf/{contract_path}",
    visibility = ["//visibility:public"],
)

java_proto_library(
    name = "{contract_name}_java_proto",
    deps = [":{contract_name}_proto"],
)

java_library(
    name = "{contract_name}_java_library",
    visibility = ["//visibility:public"],
    exports = [":{contract_name}_java_proto"],
)

py_proto_library(
    name = "{contract_name}_py_pb2",
    visibility = ["//visibility:public"],
    deps = [":{contract_name}_proto"],
)
"""

_BUILD_ROOT = """
load("@protobuf//bazel:proto_library.bzl", "proto_library")
load("@rules_go//go:def.bzl", "go_library")

DEPS = {deps}

proto_library(
    name = "{name}_proto",
    deps = ["{{}}_proto".format(dep) for dep in DEPS],
    visibility = ["//visibility:public"],
)

go_library(
    name = "{name}",
    deps = DEPS,
    visibility = ["//visibility:public"],
)
"""

def _repo_impl(ctx):
    downloads = []
    root_deps = []
    for contract in ctx.attr.contracts:
        url = "{}/{}".format(ctx.attr.base_url, contract)
        integrity = ctx.attr.integrity.get(contract, "")
        contract_file = paths.basename(contract)
        contract_path = contract.removesuffix(".proto")
        contract_name = paths.basename(contract_path)
        download = ctx.download(
            url = url,
            output = "{}/{}".format(contract_path, contract_file),
            block = False,
            integrity = integrity,
        )
        root_deps.append("//{}:{}".format(contract_path, contract_name))
        ctx.file(
            "{}/BUILD.bazel".format(contract_path),
            content = _BUILD.format(
                contract_file = contract_file,
                contract_path = contract_path,
                contract_name = contract_name,
            ),
        )
        downloads.append((download, contract, integrity))

    ctx.file("BUILD.bazel", content = _BUILD_ROOT.format(name = ctx.original_name, deps = json.encode(root_deps)))

    cur_integrity = {}
    target_integrity = {}
    for download, contract, integrity in downloads:
        res = download.wait()
        cur_integrity[contract] = integrity
        target_integrity[contract] = res.integrity
    if cur_integrity != target_integrity:
        fail(
            "invalid integrity, got:\n{}\nwant:\n{}".format(
                json.encode_indent(cur_integrity, indent = "    "),
                json.encode_indent(target_integrity, indent = "    "),
            ),
        )

    return ctx.repo_metadata(
        reproducible = True,
    )

contracts_repo = repository_rule(
    doc = "Bazel contracts repo",
    implementation = _repo_impl,
    attrs = {
        "contracts": attr.string_list(
            mandatory = True,
            doc = "Contract paths",
        ),
        "base_url": attr.string(
            mandatory = True,
            doc = "Base url",
        ),
        "integrity": attr.string_dict(
            mandatory = True,
            doc = "Integrity, keys are contract paths, values are integrity",
        ),
    },
)

def _extension_impl(ctx):
    root_module_direct_deps = []
    for mod in ctx.modules:
        for tag in mod.tags.download:
            root_module_direct_deps.append(tag.name)
            version = VERSIONS.get(tag.version, {})
            contracts_repo(
                name = tag.name,
                contracts = tag.contracts or version["contracts"],
                base_url = tag.base_url or version["base_url"],
                integrity = tag.integrity or version.get("integrity", {}),
            )
    return ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

contracts_extension = module_extension(
    implementation = _extension_impl,
    doc = "Create a repo of bazel contracts",
    tag_classes = {
        "download": tag_class({
            "name": attr.string(
                mandatory = True,
                doc = "Name",
            ),
            "version": attr.string(
                doc = "Version",
            ),
            "integrity": attr.string_dict(
                doc = "Integrity",
            ),
            "base_url": attr.string(
                doc = "Base url",
            ),
            "prefix": attr.string(
                doc = "Prefix for the contracts inside the repo",
            ),
            "contracts": attr.string_list(
                default = CONTRACTS,
                doc = "List of contract paths",
            ),
        }),
    },
)
