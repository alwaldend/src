load(
    "private/utils.bzl",
    _genrule_src = "genrule_src",
    _genrule_with_wheels = "genrule_with_wheels",
    _pkg_tar_combined = "pkg_tar_combined",
    _py_binary_shell = "py_binary_shell",
    _apply_patches = "apply_patches",
)

genrule_with_wheels = _genrule_with_wheels
py_binary_shell = _py_binary_shell
genrule_src = _genrule_src
pkg_tar_combined = _pkg_tar_combined
apply_patches = _apply_patches
