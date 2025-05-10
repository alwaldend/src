load(
    "private/utils.bzl",
    _apply_patches = "apply_patches",
    _combine_files = "combine_files",
    _compile_pip_requirements_combined = "compile_pip_requirements_combined",
    _genrule_src = "genrule_src",
    _genrule_with_wheels = "genrule_with_wheels",
    _install_file = "install_file",
    _pkg_tar_combined = "pkg_tar_combined",
    _py_binary_shell = "py_binary_shell",
    _sh_script = "sh_script",
)

genrule_with_wheels = _genrule_with_wheels
py_binary_shell = _py_binary_shell
genrule_src = _genrule_src
pkg_tar_combined = _pkg_tar_combined
apply_patches = _apply_patches
combine_files = _combine_files
compile_pip_requirements_combined = _compile_pip_requirements_combined
sh_script = _sh_script
install_file = _install_file
