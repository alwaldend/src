LABELS = struct(
    black = "//py:black",
    isort = "//py:isort",
    mypy = "//py:mypy",
    pyproject = "//:pyproject.toml",
    run_args = "//sh/scripts:run-args-lib",
    go = "@rules_go//go",
    flake8 = "//py:flake8",
    tomlv = "//go:tomlv",
    editorconfig = "//:.editorconfig",
    stylua = "@com-alwaldend-src-cargo//:stylua__stylua",
    shfmt = "@cc_mvdan_sh_v3//cmd/shfmt:shfmt",
    stylua_config = "//lua:stylua.toml",
    shellcheck = "@com-github-koalaman-shellcheck//:shellcheck",
    install_file = "//py/install-file:lib",
    replace_section = "//py/replace-section",
    bazel_python_shell = "//py/bazel-python-shell:library",
)

TOOLCHAIN_TYPES = struct(
    qt = "//bzl/qt:qt-toolchain",
)
