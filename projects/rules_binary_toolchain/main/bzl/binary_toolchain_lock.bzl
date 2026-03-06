BinaryToolchainLockArchive = provider(
    doc = "Binary toolchain archive",
    fields = {
        "download": "Kwargs for repository_ctx.download",
        "extract":  "Kwargs for repository_ctx.extract",
        "toolchain": "Kwargs for toolchain rules",
    }
)

BinaryToolchainLock = provider(
    doc = "Lock for a binary toolchain",
    fields = {
        "version": "Version to use",
        "version_archives": "Map of version archives",
        "archives": "Map of archive versions where values are maps of per-platorm archives",
    }
)

def binary_toolchain_parse_lock(lock):
    """
    Args:
        lock: lock dict

    Returns:
        BinaryToolchainLock
    """
    if not lock:
        fail("lock is empty")
    if type(lock) != "dict":
        fail("lock is not a dict: {}".format(type(lock)))
    version = lock.get("version")
    if not version:
        fail("lock is missing version")
    archives = lock.get("archives", {})
    if version not in archives:
        fail("missing archives for version {}, have {}".format(version, archives.keys()))
    res = BinaryToolchainLock(
        version = version,
        archives = {},
    )
    for cur_version, cur_archives in archives.items():
        res.archives[cur_version] = {}
        if cur_version == version:
            res.version_archives = res.archives[cur_version]
        for archive_name, archive in cur_archives.items():
            download = archive.get("download")
            if not download:
                fail("archive {}/{} is missing download".format(cur_version, archive_name))
            extract = archive.get("extract")
            toolchain = archive.get("toolchain")
            res.archives[cur_version][archive_name] = BinaryToolchainArchive(
                download = download,
                extract = extract,
                toolchain = toolchain,
            )
    return res

def binary_toolchain_parse_lock_attr(ctx):
    """
    Args:
        ctx: repository ctx

    Returns:
        BinaryToolchainLock
    """
    return binary_toolchain_parse_lock(json.decode(ctx.file(ctx.attr.lock)))
