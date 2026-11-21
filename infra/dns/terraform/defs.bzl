"""Bazel entry points for owner-local Terraform DNS roots."""

load("//infra:al_config.bzl", "infra_dns")
load("//projects/al:defs.bzl", "al_config")
load("//tools/al:al_terraform.bzl", "al_terraform_package")
load("//tools/rules_terraform/terraform:defs.bzl", "terraform_binary")
load("//tools/vault/tf_backend:defs.bzl", "vault_tf_backend")

def _terraform_file_impl(ctx):
    ctx.actions.symlink(output = ctx.outputs.out, target_file = ctx.file.src)
    return [DefaultInfo(files = depset([ctx.outputs.out]))]

_terraform_file = rule(
    implementation = _terraform_file_impl,
    attrs = {
        "out": attr.output(mandatory = True),
        "src": attr.label(allow_single_file = [".tf"], mandatory = True),
    },
)

def _root_file(name, source):
    _terraform_file(
        name = name,
        out = name + ".tf",
        src = "//infra/dns/terraform:" + source,
        visibility = ["//visibility:private"],
    )
    return ":" + name

def infra_terraform_dns_config(name, auth, backend_path, dc1 = False, global_records = True, vault_auth = None, visibility = None):
    """Compose the DNS backend and provider credentials under an owner's AppRole.

    Args:
        name: AL config target name.
        auth: Owner AppRole config label.
        backend_path: Unique Vault KV backend path for the DNS root.
        dc1: Inject RouterOS credentials.
        global_records: Inject Cloudflare credentials.
        vault_auth: Optional named Vault authentication.
        visibility: AL config visibility.
    """
    backend = name + "_backend"
    providers = name + "_providers"
    labels = {"tf": "dns"}
    vault_tf_backend(
        name = backend,
        labels = labels,
        path = backend_path,
        vault_auth = vault_auth,
    )
    infra_dns(
        name = providers,
        dc1 = dc1,
        global_records = global_records,
        labels = labels,
        vault_auth = vault_auth,
    )
    al_config(
        name = name,
        configs = [auth, ":" + backend, ":" + providers, "//infra:al"],
        visibility = visibility,
    )

def infra_terraform_dns(
        name,
        configs,
        dc1 = False,
        data = [],
        plugin_labels = {"tf": "dns"},
        offline = False,
        validate = False,
        tests_use_al = False,
        import_addresses = False):
    """Declare operations and offline checks for an owner-local DNS root.

    Args:
        name: Operation prefix, such as tf.
        configs: Owner AL configuration selecting the DNS credentials and backend.
        dc1: Include the RouterOS provider and local DNS records.
        data: Additional owner data such as Terraform test fixtures.
        plugin_labels: Operation selectors for the supplied AL configuration.
        offline: Preserve a direct offline Terraform launcher.
        validate: Add a provider schema check without the HTTP backend.
        tests_use_al: Retain AL wrapping for the package checks.
        import_addresses: Expose the DNS module's import addresses as a root output.
    """
    root_data = data + [
        ":dnsconfig",
        _root_file(name + "_dns_common", "common.tf"),
        _root_file(name + "_dns_enabled", "enabled.tf"),
        _root_file(name + "_dns_records", "global_dc1.tf" if dc1 else "global.tf"),
        "//projects/tf_modules/dns_records" if dc1 else "//projects/tf_modules/dns_records/global",
    ]
    if import_addresses:
        root_data.append(_root_file(name + "_dns_outputs", "outputs.tf"))

    providers = ["@terraform_provider_cloudflare_cloudflare//:provider"]
    if dc1:
        providers.append("@terraform_provider_terraform_routeros_routeros//:provider")
    runtime_data = root_data + [_root_file(name + "_dns_backend", "backend.tf")]

    test_kwargs = {}
    if validate:
        test_kwargs = {
            "data": root_data,
            "validate": True,
        }

    al_terraform_package(
        name = name,
        configs = configs,
        data = runtime_data,
        plugin_labels = plugin_labels,
        providers = providers,
        test_kwargs = test_kwargs,
        tests_use_al = tests_use_al,
    )
    if offline:
        terraform_binary(
            name = "offline",
            testonly = True,
            arguments = ["--direct"],
            data = runtime_data,
            providers = providers,
        )
