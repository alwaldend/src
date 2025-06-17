_common_deps = [
    "@bazel_tools//tools/build_defs/repo:cache.bzl",
    "@bazel_tools//tools/build_defs/repo:http.bzl",
    "@bazel_tools//tools/build_defs/repo:utils.bzl",
    "@bazel_skylib//:bzl_library",
    "@bazel_skylib//rules:run_binary",
    "@bazel_skylib//rules:native_binary",
    "@rules_pkg//pkg:bzl_srcs",
    "@rules_python//python:pip_bzl",
    "@rules_python//python:py_binary_bzl",
    "@rules_python//python:py_library_bzl",
    "@rules_shell//shell:rules_bzl",
    "@stardoc//stardoc:stardoc_lib",
]

BZL_LIBS = {
    "aspects": {
        "libs": {
            "al_print_deps": {},
            "al_toml_validate": {"deps": ["//bzl/providers:al_toml_info"]},
        },
        "common_deps": _common_deps,
    },
    "extensions": {
        "libs": {
            "al_shellcheck": {},
        },
        "common_deps": _common_deps,
    },
    "macros": {
        "common_deps": _common_deps + [
            "//bzl/rules:al_genrule",
            "//bzl/rules:al_template_files",
            "//bzl/rules:al_write_script",
        ],
        "libs": {
            "al_bzl_library": {},
            "al_bzl_library_map": {"deps": [
                ":al_bzl_library",
                ":al_md_data",
            ]},
            "al_lua_library": {"deps": [":al_run_tool"]},
            "al_run_tool": {},
            "al_template_files": {},
            "al_install_file": {},
            "al_genrule_with_wheels": {},
            "al_py_binary_shell": {},
            "al_genrule_src": {},
            "al_pkg_tar_combined": {},
            "al_apply_patches": {},
            "al_combine_files": {},
            "al_compile_pip_requirements_combined": {"deps": [":al_combine_files"]},
            "al_sh_library": {},
            "al_py_checker": {"deps": [":al_run_tool"]},
            "al_py_checkers": {"deps": [":al_py_checker"]},
            "al_genrule": {},
            "al_go_checkers": {},
            "al_alias_map": {},
            "al_md_data": {},
            "al_readme": {"deps": [":al_md_data"]},
        },
    },
    "providers": {
        "common_deps": _common_deps,
        "libs": {
            "al_transitive_sources": {},
            "al_toml_info": {},
        },
    },
    "toolchains": {
        "common_deps": _common_deps + ["//bzl/vars:toolchain_types"],
        "libs": {
            "al_current_qt_toolchain": {},
            "al_preinstalled_qt_toolchain": {},
        },
    },
    "rules": {
        "common_deps": _common_deps,
        "libs": {
            "al_genrule": {},
            "al_toml_data": {"deps": ["//bzl/providers:al_toml_info"]},
            "al_template_files": {},
            "al_write_script": {},
            "al_unpack_archives": {},
        },
    },
    "vars": {
        "common_deps": _common_deps,
        "libs": {
            "bzl_libs": {},
            "toolchain_types": {},
        },
    },
}
"""
Config for bzl_library targets in //bzl/...
"""
