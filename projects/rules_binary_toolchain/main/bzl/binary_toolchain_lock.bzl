BinaryToolchainLockArchiveBinary = provider(
    doc = "Binary config",
    fields = {
        "name": "Binary name",
        "path": "Binary path",
    },
)

BinaryToolchainLockArchiveAction = provider(
    doc = "Archive action",
    fields = {
        "download": "Kwargs for repository_ctx.download",
        "download_and_extract": "Kwargs for repository_ctx.download_and_extract",
        "execute": "List of kwargs for repository_ctx.execute",
        "extract": "Kwargs for repository_ctx.extract",
    },
)

BinaryToolchainLockArchive = provider(
    doc = "Binary toolchain archive",
    fields = {
        "toolchain_name": "Toolchain name",
        "toolchain_version": "Toolchain version",
        "archive_name": "Archive name",
        "actions": "List of BinaryToolchainLockArchiveAction",
        "toolchain": "Kwargs for toolchain rules",
        "binaries": "List of BinaryToolchainLockArchiveBinary",
    },
)

BinaryToolchainLockToolchain = provider(
    doc = "Toolchain config",
    fields = {
        "name": "Toolchain name",
        "version": "Toolchain version",
    },
)

BinaryToolchainLock = provider(
    doc = "Binary toolchain lock",
    fields = {
        "toolchains": "List of BinaryToolchainLockToolchain",
        "archives": "List of BinaryToolchainLockArchive",
    },
)

def binary_toolchain_lock_find_archives(lock, toolchain_name):
    """
    Args:
        lock: BinaryToolchainLock
        toolchain_name: toolchain name

    Returns:
        list[BinaryToolchainLockArchive]
    """
    toolchain = None
    for cur_toolchain in lock.toolchains:
        if cur_toolchain.name == toolchain_name:
            toolchain = cur_toolchain
            break
    if not toolchain:
        fail("could not find toolchain with name '{}'".format(toolchain_name))
    res = []
    for archive in lock.archives:
        if archive.toolchain_name == toolchain.name and archive.toolchain_version == toolchain.version:
            res.append(archive)
    return res

def binary_toolchain_lock_parse_toolchain(toolchain):
    """
    Args:
        toolchain: toolchain dict

    Returns:
         BinaryToolchainLockToolchain
    """
    if type(toolchain) != "dict":
        fail("toolchain is not a dict")
    name = toolchain.get("name")
    if not name:
        fail("toolchain is missing name")
    version = toolchain.get("version")
    if not version:
        fail("toolchain is missing version")
    res = BinaryToolchainLockToolchain(
        name = name,
        version = version,
    )
    return res

def binary_toolchain_lock_parse_archive(archive):
    """
    Args:
        archive: archive dict

    Returns:
        BinaryToolchainLockArchive
    """
    if type(archive) != "dict":
        fail("archive is not a dict")
    if not archive:
        fail("empty archive")
    toolchain_name = archive.get("toolchain_name")
    if not toolchain_name:
        fail("archive is missing toolchain_name")
    toolchain_version = archive.get("toolchain_version")
    if not toolchain_version:
        fail("archive is missing toolchain_version")
    archive_name = archive.get("archive_name")
    if not archive_name:
        fail("archive is missing archive_name")
    actions = []
    for action in archive.get("actions", []):
        if type(action) != "dict":
            fail("invalid action type: {}".format(action))
        download = action.get("download", {})
        download_and_extract = action.get("download_and_extract", {})
        execute = action.get("execute", {})
        extract = action.get("extract", {})
        actions.append(
            BinaryToolchainLockArchiveAction(
                download = download,
                download_and_extract = download_and_extract,
                execute = execute,
                extract = extract,
            ),
        )
    binaries = []
    for bin in archive.get("binaries", []):
        if type(bin) != "dict":
            fail("invalid binary type: {}".format(bin))
        name = bin.get("name")
        if not name:
            fail("binary is missing name: {}".format(bin))
        path = bin.get("path")
        if not path:
            fail("binary is missing path: {}".format(bin))
        binaries.append(
            BinaryToolchainLockArchiveBinary(
                name = name,
                path = path,
            ),
        )
    toolchain = archive.get("toolchain", {})
    res = BinaryToolchainLockArchive(
        toolchain_name = toolchain_name,
        toolchain_version = toolchain_version,
        archive_name = archive_name,
        toolchain = toolchain,
        actions = actions,
        binaries = binaries,
    )
    return res

def binary_toolchain_lock_parse(lock):
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
    toolchains = []
    archives = []
    for toolchain in lock.get("toolchains", []):
        toolchains.append(binary_toolchain_lock_parse_toolchain(toolchain))
    for archive in lock.get("archives", []):
        archives.append(binary_toolchain_lock_parse_archive(archive))
    res = BinaryToolchainLock(
        archives = archives,
        toolchains = toolchains,
    )
    return res
