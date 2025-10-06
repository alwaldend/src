load("@bazel_skylib//lib:structs.bzl", "structs")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

def _dep_from_dict(dep):
    """
    Convert dep dict to a struct

    Args:
        dep: dep dict

    Returns:
        dep struct if valid, None otherwise
    """
    repository = dep.get("repository")
    name = dep.get("name")
    version = dep.get("version")
    if not repository or not name or not version:
        return None
    file_name = "{}-{}.tgz".format(name, version)
    bazel_name = repository.split("/")[2].split(".")
    bazel_name = reversed(bazel_name)
    bazel_name.append(name)
    bazel_name_str = "_".join(bazel_name)
    return struct(
        name = name,
        repository = repository,
        version = version,
        bazel_name = bazel_name_str,
        http_url = "{}/{}".format(repository, file_name),
        file_name = file_name,
        integrity = dep.get("integrity", ""),
    )

def _dep_path(dep):
    """
    Construct a dep path

    Args:
        dep: dep struct

    Returns:
        path (list[str])
    """
    return ["repos", dep.repository, "packages", dep.name, "versions", dep.version]

def _helm_lock_find_dep(lock, dep):
    """
    Find a dep in helm_lock.json

    Args:
        lock: lock dict
        dep: dep struct

    Returns:
        dep dict if found, None otherwise
    """
    for elem in _dep_path(dep):
        if not val:
            return None
        val = val.get(elem)
    return val

def _helm_lock_update_dep(lock, dep):
    """
    Update the dep in helm_lock.json

    Args:
        lock: lock dict
        dep: dep struct
    """
    for elem in _dep_path(dep):
        cur = lock.get(elem)
        if not cur:
            cur = {}
            lock[elem] = cur
        lock = cur
    lock.update(structs.to_dict(dep))

def _handle_tag(mctx, tag):
    # lock = json.decode(mctx.read("helm_lock.json"))
    update_lock = False
    mctx.watch(tag.chart_lock)
    contents = mctx.read(tag.chart_lock)
    deps = [{}]
    for line in contents.splitlines():
        if line.startswith("-"):
            line = line[1:]
            deps.append({})
        name, _, value = line.partition(":")
        deps[-1][name.strip()] = value.strip()

    for dep_dict in deps:
        dep = _dep_from_dict(dep_dict)
        if not dep:
            if tag.debug:
                print("invalid dep: {}".format(dep_dict))
            continue

        if not dep.integrity:
            res = mctx.download(
                url = dep.http_url,
                output = dep.file_name,
                integrity = dep.integrity,
            )
            if not res.success:
                fail("could not download dep '{}': {}".format(dep, res))
            dep_dict["integrity"] = res.integrity
            dep = _dep_from_dict(dep_dict)
            # _helm_lock_update_dep(lock, dep)
            # update_lock = True

        http_file(
            name = dep.bazel_name,
            integrity = dep.integrity,
            downloaded_file_path = dep.file_name,
            url = dep.http_url,
        )

    # if update_lock:
    #     mctx.file("helm_lock.json", content = json.encode_indent(lock, indent = "  "))

def _impl(mctx):
    metadata = mctx.extension_metadata(reproducible = True)

    for mod in mctx.modules:
        for tag in mod.tags.from_lock:
            _handle_tag(mctx, tag)

    return metadata

al_helm_deps = module_extension(
    implementation = _impl,
    doc = "Extension to download helm dependencies",
    tag_classes = {
        "from_lock": tag_class({
            "helm_lock": attr.label(
                mandatory = True,
                doc = "Helm lock label",
            ),
            "chart_lock": attr.label(
                mandatory = True,
                doc = "Lock label",
            ),
            "debug": attr.bool(
                default = False,
                doc = "Enable debug",
            ),
        }),
    },
)
