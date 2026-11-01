"""Public AL configuration, command wrappers, release rules, and providers."""

load("//rules:al_binary.bzl", _al_binary_run = "al_binary_run", _al_binary_run_map = "al_binary_run_map", _al_binary_run_test = "al_binary_run_test")
load("//rules:al_config.bzl", _AlConfigInfo = "AlConfigInfo", _al_config = "al_config")
load("//rules:git_binary.bzl", _git_binary = "git_binary")
load("//rules:git_info_file.bzl", _git_info_file = "git_info_file")
load("//rules:release.bzl", _release = "release")
load("//rules:release_binary.bzl", _release_binary = "release_binary")
load("//rules:release_bundle.bzl", _release_bundle = "release_bundle")
load("//rules:release_deployment.bzl", _release_deployment = "release_deployment")
load("//rules:release_deployment_info.bzl", _ReleaseDeploymentInfo = "ReleaseDeploymentInfo")
load("//rules:release_deps.bzl", _release_deps = "release_deps")
load("//rules:release_files.bzl", _release_files = "release_files")
load("//rules:release_files_info.bzl", _ReleaseFilesInfo = "ReleaseFilesInfo")
load("//rules:release_info.bzl", _ReleaseInfo = "ReleaseInfo")

release = _release
release_binary = _release_binary
release_bundle = _release_bundle
release_deployment = _release_deployment
release_deps = _release_deps
release_files = _release_files
ReleaseDeploymentInfo = _ReleaseDeploymentInfo
ReleaseFilesInfo = _ReleaseFilesInfo
ReleaseInfo = _ReleaseInfo
git_binary = _git_binary
git_info_file = _git_info_file

al_config = _al_config
AlConfigInfo = _AlConfigInfo
al_binary_run = _al_binary_run
al_binary_run_map = _al_binary_run_map
al_binary_run_test = _al_binary_run_test
