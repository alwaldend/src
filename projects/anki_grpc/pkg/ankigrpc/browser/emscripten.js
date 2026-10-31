import createEmscriptenMemfs from "./emscripten_memfs.js";

// Classic Emscripten uses musl's generic/Linux open flags and POSIX mode bits.
// Callers use these constants rather than flags from their host platform.
const constants = Object.freeze({
    O_RDONLY: 0,
    O_WRONLY: 1,
    O_RDWR: 2,
    O_CREAT: 0o100,
    O_EXCL: 0o200,
    O_TRUNC: 0o1000,
    O_APPEND: 0o2000,
    O_NONBLOCK: 0o4000,
    O_DSYNC: 0o10000,
    O_SYNC: 0o4010000,
    O_DIRECTORY: 0o200000,
    O_NOFOLLOW: 0o400000,
    O_CLOEXEC: 0o2000000,
    S_IFMT: 0o170000,
    S_IFREG: 0o100000,
    S_IFDIR: 0o40000,
    S_IFLNK: 0o120000,
    S_IFCHR: 0o20000,
    S_IFBLK: 0o60000,
    S_IFIFO: 0o10000,
    S_IFSOCK: 0o140000,
});

const stringFlags = Object.freeze({
    r: constants.O_RDONLY,
    "r+": constants.O_RDWR,
    rs: constants.O_RDONLY | constants.O_SYNC,
    "rs+": constants.O_RDWR | constants.O_SYNC,
    w: constants.O_WRONLY | constants.O_CREAT | constants.O_TRUNC,
    wx:
        constants.O_WRONLY |
        constants.O_CREAT |
        constants.O_TRUNC |
        constants.O_EXCL,
    "w+": constants.O_RDWR | constants.O_CREAT | constants.O_TRUNC,
    "wx+":
        constants.O_RDWR |
        constants.O_CREAT |
        constants.O_TRUNC |
        constants.O_EXCL,
    a: constants.O_WRONLY | constants.O_CREAT | constants.O_APPEND,
    ax:
        constants.O_WRONLY |
        constants.O_CREAT |
        constants.O_APPEND |
        constants.O_EXCL,
    "a+": constants.O_RDWR | constants.O_CREAT | constants.O_APPEND,
    "ax+":
        constants.O_RDWR |
        constants.O_CREAT |
        constants.O_APPEND |
        constants.O_EXCL,
    as:
        constants.O_WRONLY |
        constants.O_CREAT |
        constants.O_APPEND |
        constants.O_SYNC,
    "as+":
        constants.O_RDWR |
        constants.O_CREAT |
        constants.O_APPEND |
        constants.O_SYNC,
});

function stats(value, options) {
    const result = { ...value };
    const convert = options?.bigint ? BigInt : Number;
    for (const key of [
        "dev",
        "ino",
        "mode",
        "nlink",
        "uid",
        "gid",
        "rdev",
        "size",
        "blksize",
        "blocks",
    ]) {
        result[key] = convert(value[key]);
    }
    // MEMFS has no birthtime. Node permits the ctime fallback on such stores.
    for (const key of ["atime", "mtime", "ctime", "birthtime"]) {
        const milliseconds = (value[key] ?? value.ctime).getTime();
        result[key] = new Date(milliseconds);
        result[`${key}Ms`] = convert(milliseconds);
        if (options?.bigint)
            result[`${key}Ns`] = BigInt(milliseconds) * 1000000n;
    }
    for (const [name, type] of Object.entries({
        isFile: constants.S_IFREG,
        isDirectory: constants.S_IFDIR,
        isSymbolicLink: constants.S_IFLNK,
        isCharacterDevice: constants.S_IFCHR,
        isBlockDevice: constants.S_IFBLK,
        isFIFO: constants.S_IFIFO,
        isSocket: constants.S_IFSOCK,
    })) {
        result[name] = () => (value.mode & constants.S_IFMT) === type;
    }
    return result;
}

/**
 * Create one initialized, caller-owned classic MEMFS filesystem. Each call
 * copies options and lifecycle arrays before passing them to the generated
 * factory; initialization errors reject unchanged.
 * The returned Node API subset supports string paths, byte-buffer read/write,
 * numeric/string open flags, mkdir options and bigint stat options. It uses
 * browser-native Uint8Array, Date and queueMicrotask; no globals are installed.
 *
 * FS owns all bytes, nodes and descriptors. Its module stays private and alive
 * with this object. Standard descriptors 0–2 remain reserved. Hard links fail
 * with upstream EMLINK; ownership changes fail ENOTSUP because MEMFS cannot
 * represent them. Metadata (including nlink/blocks) is upstream MEMFS metadata.
 * fsync only completes in memory; durable snapshots remain caller-owned.
 */
export async function createEmscriptenFilesystem(options) {
    const moduleOptions = { ...options };
    for (const name of ["preInit", "preRun", "postRun"]) {
        if (Array.isArray(moduleOptions[name]))
            moduleOptions[name] = moduleOptions[name].slice();
    }
    const module = await createEmscriptenMemfs(moduleOptions);
    const FS = module.FS;
    const codes = new Map(
        Object.entries(module.ERRNO_CODES).map(([code, number]) => [
            number,
            code,
        ]),
    );
    const fs = { constants };

    function fail(code) {
        throw Object.assign(new Error(code), { code });
    }

    if (!FS.initialized || [0, 1, 2].some((fd) => !FS.getStream(fd))) {
        fail("EINVAL");
    }

    function filename(value) {
        if (typeof value !== "string" || value.includes("\0")) fail("EINVAL");
        return value;
    }

    function length(value) {
        if (!Number.isSafeInteger(value) || value < 0) fail("EINVAL");
        return value;
    }

    function stream(fd) {
        if (!Number.isInteger(fd) || fd < 0) fail("EBADF");
        return FS.getStreamChecked(fd);
    }

    function milliseconds(value) {
        const number =
            value instanceof Date ? value.getTime() : Number(value) * 1000;
        if (!Number.isFinite(number)) fail("EINVAL");
        return number;
    }

    function io(method, fd, buffer, offset, count, position) {
        const opened = stream(fd);
        if (!ArrayBuffer.isView(buffer)) fail("EINVAL");
        length(offset);
        length(count);
        if (offset + count > buffer.byteLength) fail("EINVAL");
        if (position != null) length(position);
        const bytes = new Uint8Array(
            buffer.buffer,
            buffer.byteOffset,
            buffer.byteLength,
        );
        const writing = method === "write";
        const append = writing && opened.isAppend;
        const previous = opened.position;
        const effective = append ? FS.fstat(fd).size : (position ?? previous);
        const extent = effective + count;
        if (writing && count && !Number.isSafeInteger(extent)) fail("EFBIG");
        try {
            // Node pwrite on O_APPEND ignores its position but retains the fd
            // cursor. FS.write seeks first and otherwise honors the position.
            return FS[method](
                opened,
                bytes,
                offset,
                count,
                position == null ? undefined : effective,
            );
        } catch (error) {
            if (writing) opened.position = previous;
            // Only identify native allocation failures from a growing regular
            // file. Other RangeErrors and programming exceptions propagate.
            if (
                writing &&
                FS.isFile(opened.node.mode) &&
                error instanceof RangeError &&
                extent > FS.fstat(fd).size
            ) {
                if (
                    /^Invalid (?:typed array|array buffer) length(?::|$)/.test(
                        error.message,
                    )
                )
                    fail("EFBIG");
                if (error.message === "Array buffer allocation failed")
                    fail("ENOMEM");
            }
            throw error;
        } finally {
            if (append && position != null) opened.position = previous;
        }
    }

    function mkdir(path, options = 0o777) {
        filename(path);
        const mode =
            typeof options === "object" ? (options.mode ?? 0o777) : options;
        if (!options?.recursive) {
            FS.mkdir(path, mode);
            return;
        }
        if (!path) fail("ENOENT");
        let current = path.startsWith("/") ? "/" : "";
        let first;
        for (const component of path.split("/")) {
            if (!component || component === ".") continue;
            current += `${current && !current.endsWith("/") ? "/" : ""}${component}`;
            if (component === "..") {
                FS.stat(current);
                continue;
            }
            try {
                FS.mkdir(current, mode);
                first ??= current;
            } catch (error) {
                if (error.errno !== module.ERRNO_CODES.EEXIST) throw error;
                if (!FS.isDir(FS.stat(current).mode)) fail("ENOTDIR");
            }
        }
        return first;
    }

    const operations = {
        open(path, flags, mode = 0o666) {
            if (typeof flags === "string") flags = stringFlags[flags];
            const supported =
                3 |
                Object.entries(constants)
                    .filter(([key]) => key.startsWith("O_"))
                    .reduce((bits, [, value]) => bits | value, 0);
            if (!Number.isInteger(flags) || flags < 0 || (flags & 3) === 3)
                fail("EINVAL");
            if (flags & ~supported) fail("ENOTSUP");
            // Exclusive creation must see the link itself, including dangling
            // links, before upstream FS.open can create or truncate a target.
            if (flags & constants.O_CREAT && flags & constants.O_EXCL)
                flags |= constants.O_NOFOLLOW;
            return FS.open(filename(path), flags, mode).fd;
        },
        close(fd) {
            // Closing standard streams would let FS reuse Go's reserved fds.
            if (fd >= 0 && fd <= 2) fail("EBADF");
            FS.close(stream(fd));
        },
        read: (...args) => io("read", ...args),
        write: (...args) => io("write", ...args),
        stat: (path, options) => stats(FS.stat(filename(path)), options),
        lstat: (path, options) => stats(FS.lstat(filename(path)), options),
        fstat(fd, options) {
            stream(fd);
            return stats(FS.fstat(fd), options);
        },
        readdir(path) {
            return FS.readdir(filename(path)).filter(
                (name) => name !== "." && name !== "..",
            );
        },
        mkdir,
        rmdir: (path) => FS.rmdir(filename(path)),
        unlink: (path) => FS.unlink(filename(path)),
        rename: (from, to) => FS.rename(filename(from), filename(to)),
        truncate: (path, size = 0) =>
            FS.truncate(filename(path), length(size)),
        ftruncate(fd, size = 0) {
            stream(fd);
            FS.ftruncate(fd, length(size));
        },
        fsync(fd) {
            const opened = stream(fd);
            // This is Emscripten's fsync convention: synchronous backends have
            // already completed writes when no stream fsync operation exists.
            const error = opened.stream_ops.fsync?.(opened);
            if (error) throw new FS.ErrnoError(error < 0 ? -error : error);
        },
        utimes: (path, atime, mtime) =>
            FS.utime(filename(path), milliseconds(atime), milliseconds(mtime)),
        futimes(fd, atime, mtime) {
            const opened = stream(fd);
            FS.doSetAttr(opened, opened.node, {
                atime: milliseconds(atime),
                mtime: milliseconds(mtime),
            });
        },
        chmod: (path, mode) => FS.chmod(filename(path), mode),
        fchmod(fd, mode) {
            stream(fd);
            FS.fchmod(fd, mode);
        },
        chown(path) {
            FS.stat(filename(path));
            fail("ENOTSUP");
        },
        fchown(fd) {
            stream(fd);
            fail("ENOTSUP");
        },
        lchown(path) {
            FS.lstat(filename(path));
            fail("ENOTSUP");
        },
        readlink: (path) => FS.readlink(filename(path)),
        link: (from, to) => FS.link(filename(from), filename(to)),
        symlink(target, path, type) {
            if (type != null && type !== "file" && type !== "dir")
                fail("ENOTSUP");
            FS.symlink(filename(target), filename(path));
        },
    };

    for (const [name, operation] of Object.entries(operations)) {
        const sync = (...args) => {
            try {
                return operation(...args);
            } catch (error) {
                // Do not modify upstream's shared/cached ErrnoError instances.
                const code = error.code ?? codes.get(error.errno);
                if (!code) throw error;
                throw Object.assign(
                    new Error(`${code}: ${name}`, { cause: error }),
                    { code, syscall: name },
                );
            }
        };
        fs[`${name}Sync`] = sync;
        fs[name] = (...args) => {
            const callback = args.pop();
            if (typeof callback !== "function")
                throw new TypeError(`${name} requires a callback`);
            queueMicrotask(() => {
                let result;
                try {
                    result = sync(...args);
                } catch (error) {
                    if (name === "read" || name === "write")
                        callback(error, 0, args[1]);
                    else callback(error);
                    return;
                }
                // Invoke outside the try: a caller's exception is not an I/O
                // failure and must never cause a second callback invocation.
                if (name === "read" || name === "write")
                    callback(null, result, args[1]);
                else if (result === undefined) callback(null);
                else callback(null, result);
            });
        };
    }
    return fs;
}
