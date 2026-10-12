"""Standard rules_rust targets for members of Anki's imported Cargo workspace."""

load("@com_github_ankitects_anki_crates//:crates.bzl", "aliases", "all_crate_deps", "crate_edition")
load("@rules_rust//cargo:defs.bzl", "cargo_build_script", "cargo_toml_env_vars")
load("@rules_rust//rust:defs.bzl", "rust_library")

def anki_crate(name, path, deps = [], build_deps = None, build_data = [], build_env = {}, build_tools = [], features = [], compile_data = []):
    """Declare a workspace library with Cargo-owned external dependencies.

    Args:
        name: Upstream crate name and library target name.
        path: Upstream workspace-relative directory containing Cargo.toml.
        deps: Workspace-member library edges, omitted by all_crate_deps.
        build_deps: Workspace-member build-script edges, or None for no script.
        build_data: Source and generated files read by the upstream build script.
        build_env: Explicit environment for upstream build-script tools/outputs.
        build_tools: Executable inputs used by the upstream build script.
        features: Features enabled for this library and its build script.
        compile_data: Additional files embedded by the library's source.
    """
    env = name + "_env"
    cargo_toml_env_vars(
        name = env,
        src = path + "/Cargo.toml",
        workspace = "Cargo.toml",
    )
    script = []
    if build_deps != None:
        script_name = name + "_build"
        cargo_build_script(
            name = script_name,
            srcs = native.glob([path + "/*.rs"]),
            crate_root = path + "/build.rs",
            crate_features = features,
            build_script_env = build_env,
            build_script_env_files = [":" + env],
            crate_name = "build_script_build",
            pkg_name = name,
            edition = crate_edition(package_name = path),
            aliases = aliases(build = True, build_proc_macro = True, package_name = path),
            deps = all_crate_deps(build = True, package_name = path) + build_deps,
            proc_macro_deps = all_crate_deps(build_proc_macro = True, package_name = path),
            data = [path + "/Cargo.toml"] + build_data,
            tools = build_tools,
            rundir = "external/" + native.repo_name() + "/" + path,
            use_cc_toolchain = False,
            tags = ["noclippy", "norustfmt"],
        )
        script = [":" + script_name]
    rust_library(
        name = name,
        crate_name = name,
        crate_root = path + "/src/lib.rs",
        srcs = native.glob([path + "/src/**/*.rs"]),
        compile_data = native.glob([path + "/src/**"], exclude = ["**/*.rs"], allow_empty = True) + compile_data,
        crate_features = features,
        edition = crate_edition(package_name = path),
        rustc_env_files = [":" + env],
        aliases = aliases(package_name = path),
        deps = all_crate_deps(normal = True, package_name = path) + deps + script,
        proc_macro_deps = all_crate_deps(proc_macro = True, package_name = path),
        tags = ["noclippy", "norustfmt"],
    )
