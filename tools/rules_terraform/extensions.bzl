"""Hermetic provider downloads through terraform_providers.archive tags.

Each tag declares a complete provider selection and SHA256 SRI integrity pin.
The generated repository exposes :provider as a TerraformProviderInfo target.
Downloads occur during repository fetching; no registry resolution, host tools,
or provider installation commands run in the extension.
"""

load("//terraform:providers.bzl", "terraform_provider_archive_error", "terraform_provider_mirror_path")

_LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
_BASE64 = _LETTERS + "0123456789+/"

def terraform_provider_pin_error(name, source, version, platform, urls, integrity):
    """Returns a pin validation error, or None for an explicit provider pin.

    Args:
      name: Generated Bazel repository name.
      source: Canonical provider address.
      version: Exact semantic version.
      platform: Terraform OS_ARCH platform.
      urls: HTTPS URLs for the same immutable release archive.
      integrity: SHA256 subresource integrity string.

    Returns:
      A diagnostic string, or None when the pin is valid.
    """
    if not name or name[0] not in _LETTERS or not all([character in _LETTERS + "0123456789._-" for character in name.elems()]):
        return "name must be a Bazel repository name beginning with a letter"
    error = terraform_provider_archive_error(source, version, platform)
    if error:
        return error
    if not integrity.startswith("sha256-"):
        return "integrity must contain one SHA256 SRI digest (sha256-BASE64)"
    digest = integrity.removeprefix("sha256-")
    if len(digest) != 44 or digest[-1] != "=" or not all([character in _BASE64 for character in digest[:-1].elems()]) or digest[-2] not in "AEIMQUYcgkosw048":
        return "integrity must contain a canonical base64-encoded 32-byte SHA256 digest"
    if not urls:
        return "urls must contain at least one immutable HTTPS archive URL"
    for url in urls:
        if not url.startswith("https://"):
            return "provider archive URLs must use HTTPS"
        remainder = url.removeprefix("https://")
        parts = remainder.split("/", 1)
        if len(parts) != 2 or not parts[0] or not parts[1] or any([character in url for character in [" ", "\t", "\r", "\n", "\\", "#", "?"]]) or "@" in parts[0]:
            return "provider archive URLs must have a host and immutable path without credentials, queries, or fragments"
        if any([part in ["latest", "HEAD"] for part in parts[1].split("/")]):
            return "provider archive URLs must identify immutable releases, not latest or HEAD"
    return None

def terraform_provider_archive_conflict_error(name, source, version, platform, names, selections, versions):
    """Checks repository, archive, and provider version selections before fetching.

    Args:
      name: Repository name being registered.
      source: Canonical provider address.
      version: Exact semantic version.
      platform: Terraform OS_ARCH platform.
      names: Previously registered repository names, mapped to declaring modules.
      selections: Previously registered (source, version, platform) tuples, mapped to names.
      versions: Previously selected provider sources, mapped to versions.

    Returns:
      A diagnostic string, or None when this declaration does not collide.
    """
    if name in names:
        return "duplicate provider repository name '{}' (already declared by {})".format(name, names[name])
    previous_version = versions.get(source)
    if previous_version and previous_version != version:
        return "provider '{}' selects conflicting versions '{}' and '{}'; select one version per provider across the extension graph".format(source, previous_version, version)
    previous = selections.get((source, version, platform))
    if previous:
        return "provider archive {} {} {} is already declared as '{}'".format(source, version, platform, previous)
    return None

def _provider_repository_impl(ctx):
    mirror_path = terraform_provider_mirror_path(ctx.attr.source, ctx.attr.version, ctx.attr.platform)
    ctx.download(
        url = ctx.attr.urls,
        output = mirror_path,
        integrity = ctx.attr.integrity,
        canonical_id = "{} {} {} {}".format(ctx.attr.source, ctx.attr.version, ctx.attr.platform, ctx.attr.integrity),
    )
    ctx.file("BUILD.bazel", """load({rule}, "terraform_provider")

terraform_provider(
    name = "provider",
    src = {archive},
    source = {source},
    version = {version},
    platform = {platform},
    visibility = ["//visibility:public"],
)
""".format(
        rule = json.encode(str(ctx.attr._provider_rule)),
        archive = json.encode(mirror_path),
        source = json.encode(ctx.attr.source),
        version = json.encode(ctx.attr.version),
        platform = json.encode(ctx.attr.platform),
    ))

_provider_repository = repository_rule(
    implementation = _provider_repository_impl,
    attrs = {
        "integrity": attr.string(mandatory = True),
        "platform": attr.string(mandatory = True),
        "source": attr.string(mandatory = True),
        "urls": attr.string_list(mandatory = True),
        "version": attr.string(mandatory = True),
        "_provider_rule": attr.label(default = Label("//terraform:providers.bzl")),
    },
)

def _terraform_providers_impl(ctx):
    names = {}
    selections = {}
    versions = {}
    root_deps = []
    root_dev_deps = []
    archives = []
    for mod in ctx.modules:
        for tag in mod.tags.archive:
            error = terraform_provider_pin_error(tag.name, tag.source, tag.version, tag.platform, tag.urls, tag.integrity)
            if error:
                fail("provider archive '{}' in module '{}': {}".format(tag.name, mod.name, error))
            error = terraform_provider_archive_conflict_error(tag.name, tag.source, tag.version, tag.platform, names, selections, versions)
            if error:
                fail("module '{}': {}".format(mod.name, error))
            names[tag.name] = mod.name
            selections[(tag.source, tag.version, tag.platform)] = tag.name
            versions[tag.source] = tag.version
            archives.append(tag)
            if mod.is_root:
                if ctx.is_dev_dependency(tag):
                    root_dev_deps.append(tag.name)
                else:
                    root_deps.append(tag.name)

    # Validate every declaration before registering any repository downloads.
    for tag in archives:
        _provider_repository(
            name = tag.name,
            source = tag.source,
            version = tag.version,
            platform = tag.platform,
            urls = tag.urls,
            integrity = tag.integrity,
        )
    return ctx.extension_metadata(
        root_module_direct_deps = root_deps,
        root_module_direct_dev_deps = root_dev_deps,
        reproducible = True,
    )

terraform_providers = module_extension(
    implementation = _terraform_providers_impl,
    tag_classes = {
        "archive": tag_class(attrs = {
            "integrity": attr.string(mandatory = True, doc = "Required SHA256 SRI digest verified against the publisher's release checksum."),
            "name": attr.string(mandatory = True, doc = "Unique Bazel repository name; exposes :provider."),
            "platform": attr.string(mandatory = True, doc = "Terraform OS_ARCH platform, such as linux_amd64."),
            "source": attr.string(mandatory = True, doc = "Canonical lowercase HOST/NAMESPACE/TYPE provider address."),
            "urls": attr.string_list(mandatory = True, doc = "HTTPS URLs for the same immutable release archive."),
            "version": attr.string(mandatory = True, doc = "Exact semantic version, without a leading v."),
        }),
    },
    doc = "Downloads explicitly pinned Terraform provider ZIPs with Bazel's verified repository downloader. Select one version per provider across the extension graph; multiple platforms are supported. Declare one archive per source/version/platform and import each generated :provider repository with use_repo.",
)
