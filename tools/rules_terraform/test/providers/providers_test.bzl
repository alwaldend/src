"""Provider identity, archive pin, collision, and runfiles contract tests."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts", "unittest")
load("//:extensions.bzl", "terraform_provider_archive_conflict_error", "terraform_provider_pin_error")
load("//terraform:providers.bzl", "TerraformProviderInfo", "terraform_provider_archive_error", "terraform_provider_mirror_path")

_SOURCE = "registry.terraform.io/hashicorp/local"
_VERSION = "2.9.0"
_PLATFORM = "linux_amd64"
_URL = "https://releases.hashicorp.com/terraform-provider-local/2.9.0/terraform-provider-local_2.9.0_linux_amd64.zip"
_INTEGRITY = "sha256-" + "A" * 43 + "="

def _identity_test_impl(ctx):
    env = unittest.begin(ctx)
    for version in [_VERSION, "3.0.2-rc07", "1.2.3-alpha.1+build.09"]:
        asserts.equals(env, None, terraform_provider_archive_error(_SOURCE, version, _PLATFORM))
    for source in ["hashicorp/local", "registry.terraform.io/HashiCorp/local", "registry.terraform.io/../local", "registry.terraform.io/hashicorp/", "registry.terraform.io/hashicorp/local/extra", "https://registry.terraform.io/hashicorp/local", "terraform.io/builtin/terraform"]:
        asserts.true(env, terraform_provider_archive_error(source, _VERSION, _PLATFORM) != None, source)
    for version in ["", "v1.2.3", ">= 1", "1.2", "01.2.3", "1.2.3-", "1.2.3-01", "1.2.3+", "1.2.3/../../escape"]:
        asserts.true(env, terraform_provider_archive_error(_SOURCE, version, _PLATFORM) != None, version)
    for platform in ["", "linux", "linux_x86_64", "Linux_amd64", "linux_../amd64"]:
        asserts.true(env, terraform_provider_archive_error(_SOURCE, _VERSION, platform) != None, platform)
    asserts.equals(env, "registry.terraform.io/hashicorp/local/terraform-provider-local_2.9.0_linux_amd64.zip", terraform_provider_mirror_path(_SOURCE, _VERSION, _PLATFORM))
    return unittest.end(env)

identity_test = unittest.make(_identity_test_impl)

def _pin_test_impl(ctx):
    env = unittest.begin(ctx)
    asserts.equals(env, None, terraform_provider_pin_error("local_linux_amd64", _SOURCE, _VERSION, _PLATFORM, [_URL], _INTEGRITY))
    for integrity in ["", "sha512-" + "A" * 44, "sha256-" + "A" * 44, "sha256-" + "A" * 42 + "B=", "sha256-" + "A" * 43 + "= extra"]:
        asserts.true(env, terraform_provider_pin_error("local", _SOURCE, _VERSION, _PLATFORM, [_URL], integrity) != None, integrity)
    for urls in [[], ["http://example.com/archive.zip"], ["file:///archive.zip"], ["https://example.com"], ["https://example.com/latest/archive.zip"], ["https://user@example.com/archive.zip"], ["https://example.com/archive.zip?token=example"]]:
        asserts.true(env, terraform_provider_pin_error("local", _SOURCE, _VERSION, _PLATFORM, urls, _INTEGRITY) != None, str(urls))
    for name in ["", "../local", "@local", "local/provider"]:
        asserts.true(env, terraform_provider_pin_error(name, _SOURCE, _VERSION, _PLATFORM, [_URL], _INTEGRITY) != None, name)
    return unittest.end(env)

pin_test = unittest.make(_pin_test_impl)

def _collision_test_impl(ctx):
    env = unittest.begin(ctx)
    names = {"local": "root_module"}
    selections = {(_SOURCE, _VERSION, _PLATFORM): "local"}
    versions = {_SOURCE: _VERSION}
    asserts.true(env, "duplicate provider repository name" in terraform_provider_archive_conflict_error("local", _SOURCE, "2.8.0", _PLATFORM, names, selections, versions))
    asserts.true(env, "already declared as 'local'" in terraform_provider_archive_conflict_error("other", _SOURCE, _VERSION, _PLATFORM, names, selections, versions))
    asserts.true(env, "conflicting versions '2.9.0' and '2.8.0'" in terraform_provider_archive_conflict_error("older", _SOURCE, "2.8.0", _PLATFORM, names, selections, versions))
    asserts.true(env, "conflicting versions" in terraform_provider_archive_conflict_error("arm_older", _SOURCE, "2.8.0", "linux_arm64", names, selections, versions))
    asserts.equals(env, None, terraform_provider_archive_conflict_error("arm", _SOURCE, _VERSION, "linux_arm64", names, selections, versions))
    asserts.equals(env, None, terraform_provider_archive_conflict_error("another", "registry.terraform.io/hashicorp/random", "3.9.0", _PLATFORM, names, selections, versions))
    return unittest.end(env)

collision_test = unittest.make(_collision_test_impl)

def _archive_fixture_impl(ctx):
    archive = ctx.actions.declare_file(ctx.label.name + ".zip")
    ctx.actions.write(archive, "")
    return [DefaultInfo(files = depset([archive]))]

archive_fixture = rule(implementation = _archive_fixture_impl)

def _provider_contract_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    info = target[TerraformProviderInfo]
    asserts.equals(env, _SOURCE, info.source)
    asserts.equals(env, _VERSION, info.version)
    asserts.equals(env, _PLATFORM, info.platform)
    asserts.equals(env, terraform_provider_mirror_path(_SOURCE, _VERSION, _PLATFORM), info.mirror_path)
    asserts.equals(env, [info.archive], target[DefaultInfo].files.to_list())
    asserts.true(env, info.archive in target[DefaultInfo].default_runfiles.files.to_list())
    return analysistest.end(env)

provider_contract_test = analysistest.make(_provider_contract_test_impl)
