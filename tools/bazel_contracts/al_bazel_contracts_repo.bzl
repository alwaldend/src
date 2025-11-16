load("@bazel_skylib//lib:paths.bzl", "paths")

_BUILD = """
load("@protobuf//bazel:java_proto_library.bzl", "java_proto_library")
load("@protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_go//go:def.bzl", "go_library")
load("@rules_go//proto:def.bzl", "go_proto_library")
load("@rules_java//java:defs.bzl", "java_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

exports_files(["{contract_file}"])

proto_library(
    name = "{contract_name}_proto",
    srcs = ["{contract_file}"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "{contract_name}_go_proto",
    importpath = "git.alwaldend.com/src/tools/bazel_contracts/{contract_path}",
    proto = ":{contract_name}_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "{contract_name}",
    embed = [":{contract_name}_go_proto"],
    importpath = "git.alwaldend.com/src/tools/bazel_contracts/{contract_path}",
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

def _impl(ctx):
    downloads = []
    for contract in ctx.attr.contracts:
        url = "{}/{}/{}".format(ctx.attr.base_url, ctx.attr.prefix, contract)
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
        ctx.file(
            "{}/BUILD.bazel".format(contract_path),
            content = _BUILD.format(
                contract_file = contract_file,
                contract_path = contract_path,
                contract_name = contract_name,
            ),
        )
        downloads.append((download, contract, integrity))

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

    ctx.file("BUILD.bazel", "")

    return ctx.repo_metadata(
        reproducible = True,
    )

al_bazel_contracts_repo = repository_rule(
    doc = "Bazel contracts repo",
    implementation = _impl,
    attrs = {
        "contracts": attr.string_list(
            mandatory = True,
            doc = "Contract paths",
        ),
        "prefix": attr.string(
            mandatory = True,
            doc = "Prefix for the contracts inside the repo",
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
