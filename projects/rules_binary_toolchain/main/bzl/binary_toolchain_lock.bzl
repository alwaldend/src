BinaryToolchainLockArchive = provider(
    doc = "Binary toolchain archive",
    fields = {
        "download": "Kwargs for repository_ctx.download",
        "download_and_extract": "Kwargs for repository_ctx.download_and_extract",
        "extract":  "Kwargs for repository_ctx.extract",
        "toolchain": "Kwargs for toolchain rules",
        "bin_path": "Path to the binary",
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
    res_archives = {}
    res_version_archives = None
    for cur_version, cur_archives in archives.items():
        res_archives[cur_version] = {}
        if cur_version == version:
            res_version_archives = res_archives[cur_version]
        for archive_name, archive in cur_archives.items():
            download = archive.get("download", {})
            download_and_extract = archive.get("download_and_extract", {})
            if not download and not download_and_extract:
                fail("archive {}/{} is missing download options".format(cur_version, archive_name))
            bin_path = archive.get("bin_path")
            if not bin_path:
                fail("archive {}/{} is missing bin_path".format(cur_version, archive_name))
            extract = archive.get("extract", {})
            toolchain = archive.get("toolchain", {})
            res_archives[cur_version][archive_name] = BinaryToolchainLockArchive(
                bin_path = bin_path,
                download = download,
                download_and_extract = download_and_extract,
                extract = extract,
                toolchain = toolchain,
            )

    res = BinaryToolchainLock(
        version = version,
        archives = res_archives,
        version_archives = res_version_archives,
    )
    return res

def binary_toolchain_parse_lock_attr(ctx, attr_name = "lock"):
    """
    Args:
        ctx: repository ctx
        attr_name: attributes containing the lock label

    Returns:
        BinaryToolchainLock
    """
    return binary_toolchain_parse_lock(json.decode(ctx.read(getattr(ctx.attr, attr_name))))
