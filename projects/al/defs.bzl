"""Public AL configuration, command wrappers, release rules, and providers."""

load("//projects/al/rules:al_binary.bzl", _al_binary_run = "al_binary_run", _al_binary_run_map = "al_binary_run_map", _al_binary_run_test = "al_binary_run_test")
load("//projects/al/rules:al_config.bzl", _AlConfigInfo = "AlConfigInfo", _al_config = "al_config", _al_plugin = "al_plugin", _al_plugin_config = "al_plugin_config", _al_vault_auth = "al_vault_auth", _al_vault_conn = "al_vault_conn")
load("//projects/al/rules:al_git_binary.bzl", _al_git_binary = "al_git_binary")
load("//projects/al/rules:al_inject.bzl", _al_inject = "al_inject")
load("//projects/al/rules:al_release.bzl", _al_release = "al_release")
load("//projects/al/rules:al_release_binary.bzl", _al_release_binary = "al_release_binary")
load("//projects/al/rules:al_release_bundle.bzl", _al_release_bundle = "al_release_bundle")
load("//projects/al/rules:al_release_deployment.bzl", _al_release_deployment = "al_release_deployment")
load("//projects/al/rules:al_release_deployment_info.bzl", _AlReleaseDeploymentInfo = "AlReleaseDeploymentInfo")
load("//projects/al/rules:al_release_deps.bzl", _al_release_deps = "al_release_deps")
load("//projects/al/rules:al_release_files.bzl", _al_release_files = "al_release_files")
load("//projects/al/rules:al_release_files_info.bzl", _AlReleaseFilesInfo = "AlReleaseFilesInfo")
load("//projects/al/rules:al_release_info.bzl", _AlReleaseInfo = "AlReleaseInfo")

al_release = _al_release
al_release_binary = _al_release_binary
al_release_bundle = _al_release_bundle
al_release_deployment = _al_release_deployment
al_release_deps = _al_release_deps
al_release_files = _al_release_files
AlReleaseDeploymentInfo = _AlReleaseDeploymentInfo
AlReleaseFilesInfo = _AlReleaseFilesInfo
AlReleaseInfo = _AlReleaseInfo
al_git_binary = _al_git_binary

al_config = _al_config
AlConfigInfo = _AlConfigInfo
al_binary_run = _al_binary_run
al_binary_run_map = _al_binary_run_map
al_binary_run_test = _al_binary_run_test

al_plugin_config = _al_plugin_config
al_plugin = _al_plugin
al_vault_auth = _al_vault_auth
al_vault_conn = _al_vault_conn

al_inject = _al_inject
