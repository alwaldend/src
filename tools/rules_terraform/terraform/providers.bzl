"""Pinned Terraform provider archives and their filesystem-mirror identities."""

TerraformProviderInfo = provider(
    doc = "A provider ZIP and its canonical packed filesystem-mirror path.",
    fields = {
        "archive": "The provider distribution ZIP File.",
        "mirror_path": "Relative HOST/NAMESPACE/TYPE/terraform-provider-TYPE_VERSION_OS_ARCH.zip path.",
        "platform": "Terraform platform in OS_ARCH notation, such as linux_amd64.",
        "source": "Canonical lowercase HOST/NAMESPACE/TYPE provider address.",
        "version": "Exact provider semantic version, without a leading v.",
    },
)

_LOWER = "abcdefghijklmnopqrstuvwxyz"
_DIGITS = "0123456789"
_ALNUM = _LOWER + _DIGITS

def _characters(value, allowed):
    return all([character in allowed for character in value.elems()])

def _dns_label(value):
    return bool(value) and len(value) <= 63 and value[0] in _ALNUM and value[-1] in _ALNUM and _characters(value, _ALNUM + "-")

def _numeric_identifier(value):
    return bool(value) and _characters(value, _DIGITS) and (len(value) == 1 or value[0] != "0")

def _semantic_version(value):
    metadata = value.split("+")
    if len(metadata) > 2:
        return False
    if len(metadata) == 2:
        for identifier in metadata[1].split("."):
            if not identifier or not _characters(identifier, _ALNUM + _LOWER.upper() + "-"):
                return False
    release = metadata[0].split("-", 1)
    core = release[0].split(".")
    if len(core) != 3 or not all([_numeric_identifier(part) for part in core]):
        return False
    if len(release) == 2:
        for identifier in release[1].split("."):
            if not identifier or not _characters(identifier, _ALNUM + _LOWER.upper() + "-"):
                return False
            if _characters(identifier, _DIGITS) and not _numeric_identifier(identifier):
                return False
    return True

def terraform_provider_archive_error(source, version, platform):
    """Returns an identity validation error, or None for a safe mirror path.

    Args:
      source: Canonical lowercase ASCII HOST/NAMESPACE/TYPE provider address.
      version: Exact semantic version, including optional prerelease/build suffixes.
      platform: Terraform OS_ARCH platform; architecture uses Go naming.

    Returns:
      A diagnostic string, or None when the identity is valid.
    """
    parts = source.split("/")
    if len(parts) != 3 or not all([_dns_label(part) for part in parts[1:]]):
        return "source must be a canonical lowercase HOST/NAMESPACE/TYPE provider address"
    host = parts[0]
    if len(host) > 253 or not all([_dns_label(part) for part in host.split(".")]):
        return "source hostname must contain canonical lowercase ASCII DNS labels"
    if source == "terraform.io/builtin/terraform":
        return "Terraform's built-in provider has no downloadable archive"
    if not _semantic_version(version):
        return "version must be an exact semantic version without a leading v or a constraint"
    parts = platform.split("_")
    if len(parts) != 2 or not parts[0] or parts[0][0] not in _LOWER or not all([part and _characters(part, _ALNUM) for part in parts]):
        return "platform must use Terraform OS_ARCH notation, such as linux_amd64"
    return None

def terraform_provider_mirror_path(source, version, platform):
    """Returns the packed filesystem-mirror path for a validated provider identity.

    Args:
      source: Canonical provider address.
      version: Exact semantic version.
      platform: Terraform OS_ARCH platform.

    Returns:
      A relative ZIP path following Terraform's packed mirror layout.
    """
    error = terraform_provider_archive_error(source, version, platform)
    if error:
        fail(error)
    return "{}/terraform-provider-{}_{}_{}.zip".format(source, source.split("/")[-1], version, platform)

def _terraform_provider_impl(ctx):
    mirror_path = terraform_provider_mirror_path(ctx.attr.source, ctx.attr.version, ctx.attr.platform)
    archive = ctx.file.src
    return [
        TerraformProviderInfo(
            archive = archive,
            mirror_path = mirror_path,
            platform = ctx.attr.platform,
            source = ctx.attr.source,
            version = ctx.attr.version,
        ),
        DefaultInfo(
            files = depset([archive]),
            runfiles = ctx.runfiles(files = [archive]),
        ),
    ]

terraform_provider = rule(
    implementation = _terraform_provider_impl,
    attrs = {
        "platform": attr.string(mandatory = True, doc = "Terraform OS_ARCH platform, such as linux_amd64."),
        "source": attr.string(mandatory = True, doc = "Canonical lowercase HOST/NAMESPACE/TYPE provider address."),
        "src": attr.label(mandatory = True, allow_single_file = [".zip"], doc = "A checksum-pinned provider distribution ZIP."),
        "version": attr.string(mandatory = True, doc = "Exact semantic version, without a leading v."),
    },
    doc = "Describes a provider ZIP for Terraform rules to install in a packed filesystem mirror. Use terraform_providers.archive to download publisher releases with a mandatory integrity pin.",
)
