import { createServer } from "./server.js";

/** Initialize a Rust reactor over a borrowed Node-compatible synchronous fs.
 * The absolute provider root is mounted at guest /. The caller owns storage;
 * this server owns only its descriptors and never destroys the provider.
 */
export async function createBrowserServer(module, fs, { root, timezone }) {
    const host = createWASI(fs, root, timezone);
    let server;
    let closed = false;
    function close() {
        if (closed) return;
        closed = true;
        const failures = [];
        try {
            server?.close();
        } catch (error) {
            failures.push(error);
        }
        try {
            host.close();
        } catch (error) {
            failures.push(error);
        }
        if (failures.length)
            throw new AggregateError(failures, "Close Anki browser backend");
    }
    try {
        const instance = await WebAssembly.instantiate(module, host.imports);
        host.bind(instance.exports.memory);
        instance.exports._initialize();
        server = createServer(instance);
        return {
            command(method, request) {
                if (closed) throw new Error("Anki backend is closed");
                try {
                    return server.command(method, request);
                } catch (error) {
                    try {
                        close();
                    } catch (cleanup) {
                        throw new AggregateError(
                            [error, cleanup],
                            "Anki command and cleanup failed",
                        );
                    }
                    throw error;
                }
            },
            close,
        };
    } catch (error) {
        try {
            close();
        } catch (cleanup) {
            throw new AggregateError(
                [error, cleanup],
                "Anki initialization and cleanup failed",
            );
        }
        throw error;
    }
}

/** Report a WASI process exit as a terminal backend failure. */
class CommandExit extends Error {
    constructor(code) {
        super(`Anki WASI command exited with status ${code}`);
        this.code = code;
    }
}

/** Carry a checked ABI failure back as errno instead of trapping the guest. */
class WasiError extends Error {
    constructor(errno) {
        super(`WASI errno ${errno}`);
        this.errno = errno;
    }
}

/** Adapt Preview 1 imports directly to the caller's synchronous Node fs API. */
function createWASI(fs, root, timezone) {
    if (!fs || !fs.constants || typeof fs.openSync !== "function") {
        throw new TypeError(
            "A Node-compatible synchronous filesystem is required",
        );
    }
    if (
        typeof root !== "string" ||
        !root.startsWith("/") ||
        root.includes("\0")
    ) {
        throw new TypeError(
            "An absolute provider filesystem root is required",
        );
    }
    if (
        !timezone ||
        typeof timezone.offsetForUTC !== "function" ||
        typeof timezone.offsetsForLocal !== "function"
    ) {
        throw new TypeError("A Go timezone host is required");
    }
    root = root.replace(/\/+$/, "") || "/";
    const constants = fs.constants;

    const encoder = new TextEncoder();
    const decoder = new TextDecoder("utf-8", { fatal: true });
    const argumentsBytes = ["anki_browser_backend"].map((text) =>
        encoder.encode(`${text}\0`),
    );
    const environmentBytes = [encoder.encode("TMPDIR=/tmp\0")];
    const allFileRights = (1n << 28n) - 1n;
    const descriptors = new Map([
        [0, { base: 2n, inheriting: 0n }],
        [1, { base: 64n, inheriting: 0n }],
        [2, { base: 64n, inheriting: 0n }],
        [
            3,
            {
                base: allFileRights,
                inheriting: allFileRights,
                name: "/",
                directory: true,
                readable: true,
                writable: false,
                flags: 0,
            },
        ],
    ]);
    let memory;
    let nextDescriptor = 4;
    if (!fs.statSync(root).isDirectory())
        throw Object.assign(new Error("Anki root is not a directory"), {
            code: "ENOTDIR",
        });
    fs.mkdirSync(providerPath("/tmp"), { recursive: true, mode: 0o700 });
    descriptors.get(3).provider = fs.openSync(root, constants.O_RDONLY);

    function providerPath(name) {
        return root === "/" ? name : root + name;
    }

    function resolve(fd, name) {
        const directory = requireDescriptor(fd);
        if (!directory.directory) throw new WasiError(54);
        if (!name) throw new WasiError(44);
        if (name.startsWith("/") || name.includes("\0"))
            throw new WasiError(28);
        const parts = directory.name.split("/").filter(Boolean);
        for (const part of name.split("/")) {
            if (part === "..") parts.pop();
            else if (part && part !== ".") parts.push(part);
        }
        return `/${parts.join("/")}`;
    }

    function number(value) {
        const result = Number(value);
        if (!Number.isSafeInteger(result)) throw new WasiError(61);
        if (result < 0) throw new WasiError(28);
        return result;
    }

    function fileType(info) {
        if (info.isDirectory()) return 3;
        if (info.isFile()) return 4;
        if (info.isSymbolicLink()) return 7;
        if (info.isCharacterDevice()) return 2;
        if (info.isBlockDevice()) return 1;
        if (info.isSocket()) return 6;
        return 0;
    }

    function stat(fd) {
        return fs.fstatSync(requireDescriptor(fd).provider, { bigint: true });
    }

    function closeDescriptor(fd) {
        const descriptor = requireDescriptor(fd);
        // A failed close can still release the provider fd. Invalidate the guest
        // mapping before one attempt; never retry an uncertain provider close.
        descriptors.delete(fd);
        // WASI returns one errno: a close failure takes precedence over a prior
        // sync failure. Keep unknown close exceptions terminal via checked().
        if (descriptor.provider !== undefined)
            fs.closeSync(descriptor.provider);
        throwSyncError(descriptor);
    }

    // A write count reports transferred bytes. Preserve a later sync failure
    // until the next write, explicit sync or close can report it without replay.
    function throwSyncError(descriptor) {
        const error = descriptor.syncError;
        descriptor.syncError = undefined;
        if (error) throw error;
    }

    function seek(fd, offset, whence) {
        const descriptor = requireDescriptor(fd);
        if (descriptor.directory) throw new WasiError(31);
        let position;
        if (whence === 0) position = 0n;
        else if (whence === 1) position = BigInt(descriptor.position);
        else if (whence === 2) position = BigInt(stat(fd).size);
        else throw new WasiError(28);
        descriptor.position = number(position + offset);
        return BigInt(descriptor.position);
    }

    function open(fd, name, lookup, flags, base, inheriting, fdFlags) {
        name = resolve(fd, name);
        const readable = Boolean(base & 2n);
        const writable = Boolean(base & 64n);
        if (flags & 8 && !writable) throw new WasiError(2);
        let info;
        try {
            info =
                lookup & 1
                    ? fs.statSync(providerPath(name))
                    : fs.lstatSync(providerPath(name));
        } catch (error) {
            if (error.code !== "ENOENT") throw error;
        }
        if (info && flags & 1 && flags & 4) throw new WasiError(20);
        if (!(lookup & 1) && info?.isSymbolicLink()) throw new WasiError(32);
        if (flags & 2) {
            if (!info) throw new WasiError(44);
            if (!info.isDirectory()) throw new WasiError(54);
        }
        if (info?.isDirectory() && (writable || flags & (1 | 8)))
            throw new WasiError(31);
        let mode = writable
            ? readable
                ? constants.O_RDWR
                : constants.O_WRONLY
            : constants.O_RDONLY;
        if (flags & 1) mode |= constants.O_CREAT;
        if (flags & 1 && flags & 4) mode |= constants.O_EXCL;
        if (flags & 8) mode |= constants.O_TRUNC;
        if (flags & 2) mode |= constants.O_DIRECTORY || 0;
        if (!(lookup & 1)) mode |= constants.O_NOFOLLOW || 0;
        // Guest positions and append/sync flags are owned here, not by a second
        // provider cursor. In particular, Linux O_APPEND would ignore pwrite's position.
        const provider = fs.openSync(providerPath(name), mode, 0o600);
        try {
            info = fs.fstatSync(provider);
            if (flags & 2 && !info.isDirectory()) throw new WasiError(54);
            const opened = nextDescriptor++;
            descriptors.set(opened, {
                base,
                inheriting,
                provider,
                name,
                directory: info.isDirectory(),
                readable,
                writable,
                flags: fdFlags,
                position: 0,
            });
            return opened;
        } catch (error) {
            try {
                fs.closeSync(provider);
            } catch (cleanup) {
                throw new AggregateError(
                    [error, cleanup],
                    "Open descriptor cleanup failed",
                );
            }
            throw error;
        }
    }

    function read(fd, size) {
        const descriptor = requireDescriptor(fd);
        if (!descriptor.readable) throw new WasiError(8);
        if (descriptor.directory) throw new WasiError(31);
        const buffer = new Uint8Array(size);
        const count = fs.readSync(
            descriptor.provider,
            buffer,
            0,
            size,
            descriptor.position,
        );
        descriptor.position = number(
            BigInt(descriptor.position) + BigInt(count),
        );
        return buffer.subarray(0, count);
    }

    function write(fd, buffer) {
        const descriptor = requireDescriptor(fd);
        if (!descriptor.writable) throw new WasiError(8);
        if (descriptor.directory) throw new WasiError(31);
        throwSyncError(descriptor);
        if (descriptor.flags & 1) descriptor.position = number(stat(fd).size);
        const count = fs.writeSync(
            descriptor.provider,
            buffer,
            0,
            buffer.length,
            descriptor.position,
        );
        descriptor.position = number(
            BigInt(descriptor.position) + BigInt(count),
        );
        if (descriptor.flags & (2 | 8 | 16)) {
            try {
                fs.fsyncSync(descriptor.provider);
            } catch (error) {
                if (nodeErrno(error) === undefined) throw error;
                if (count === 0) throw error;
                descriptor.syncError = error;
            }
        }
        return count;
    }

    function times(info, access, modification, flags) {
        return [
            flags & 1
                ? Number(access) / 1e9
                : flags & 2
                  ? Date.now() / 1000
                  : Number(info.atimeMs) / 1000,
            flags & 4
                ? Number(modification) / 1e9
                : flags & 8
                  ? Date.now() / 1000
                  : Number(info.mtimeMs) / 1000,
        ];
    }

    function readdir(fd, restart) {
        const descriptor = requireDescriptor(fd);
        if (!descriptor.directory) throw new WasiError(54);
        if (!descriptor.entries || restart) {
            descriptor.entries = fs
                .readdirSync(providerPath(descriptor.name))
                .sort()
                .map((name) => {
                    const info = fs.lstatSync(
                        providerPath(resolve(fd, name)),
                        { bigint: true },
                    );
                    return { name, inode: info.ino, type: fileType(info) };
                });
        }
        return descriptor.entries;
    }

    function bytes(pointer, length) {
        pointer >>>= 0;
        if (
            !memory ||
            !Number.isSafeInteger(length) ||
            length < 0 ||
            pointer > memory.buffer.byteLength ||
            length > memory.buffer.byteLength - pointer
        ) {
            throw new WasiError(21);
        }
        return new Uint8Array(memory.buffer, pointer, length);
    }

    function data(pointer, length) {
        const value = bytes(pointer, length);
        return new DataView(value.buffer, value.byteOffset, value.byteLength);
    }

    function path(pointer, length) {
        const value = bytes(pointer, length >>> 0);
        try {
            return decoder.decode(value);
        } catch {
            throw new WasiError(28);
        }
    }

    // Node error codes retain POSIX distinctions instead of collapsing to EIO.
    const errnos = {
        E2BIG: 1,
        EACCES: 2,
        EAGAIN: 6,
        EBADF: 8,
        EBUSY: 10,
        ECANCELED: 11,
        EDQUOT: 19,
        EEXIST: 20,
        EFAULT: 21,
        EFBIG: 22,
        EINTR: 27,
        EINVAL: 28,
        EIO: 29,
        EISDIR: 31,
        ELOOP: 32,
        EMFILE: 33,
        EMLINK: 34,
        ENAMETOOLONG: 37,
        ENFILE: 41,
        ENODEV: 43,
        ENOENT: 44,
        ENOMEM: 48,
        ENOSPC: 51,
        ENOSYS: 52,
        ENOTDIR: 54,
        ENOTEMPTY: 55,
        ENOTSUP: 58,
        EOPNOTSUPP: 58,
        ENOTTY: 59,
        ENXIO: 60,
        EOVERFLOW: 61,
        EPERM: 63,
        EPIPE: 64,
        // Preview 1 has no remote-I/O errno; retain the concrete I/O failure.
        EREMOTEIO: 29,
        EROFS: 69,
        ESPIPE: 70,
        ESTALE: 72,
        ETIMEDOUT: 73,
        ETXTBSY: 74,
        EXDEV: 75,
    };

    function nodeErrno(error) {
        return Object.hasOwn(errnos, error?.code)
            ? errnos[error.code]
            : undefined;
    }

    function checked(operation) {
        return (...args) => {
            try {
                operation(...args);
                return 0;
            } catch (error) {
                if (error instanceof WasiError) return error.errno;
                const errno = nodeErrno(error);
                if (errno !== undefined) return errno;
                throw error;
            }
        };
    }

    function requireDescriptor(fd, right = 0n) {
        const descriptor = descriptors.get(fd);
        if (!descriptor) throw new WasiError(8);
        if ((descriptor.base & right) !== right) throw new WasiError(76);
        return descriptor;
    }

    function stringsSizes(values, count, size) {
        data(count, 4).setUint32(0, values.length, true);
        data(size, 4).setUint32(
            0,
            values.reduce((total, value) => total + value.length, 0),
            true,
        );
    }

    function stringsGet(values, pointers, buffer) {
        const table = data(pointers, values.length * 4);
        let offset = buffer >>> 0;
        bytes(
            offset,
            values.reduce((total, value) => total + value.length, 0),
        );
        for (const [index, value] of values.entries()) {
            table.setUint32(index * 4, offset, true);
            bytes(offset, value.length).set(value);
            offset += value.length;
        }
    }

    function writeStat(pointer, value) {
        const output = data(pointer, 64);
        bytes(pointer, 64).fill(0);
        output.setBigUint64(0, BigInt(value.dev), true);
        output.setBigUint64(8, BigInt(value.ino), true);
        output.setUint8(16, fileType(value));
        output.setBigUint64(24, BigInt(value.nlink), true);
        output.setBigUint64(32, BigInt(value.size), true);
        for (const [offset, field] of [
            [40, "atime"],
            [48, "mtime"],
            [56, "ctime"],
        ]) {
            output.setBigUint64(
                offset,
                value[`${field}Ns`] ??
                    BigInt(Math.trunc(Number(value[`${field}Ms`]) * 1e6)),
                true,
            );
        }
    }

    function checkTimeFlags(flags) {
        if (flags & ~15 || (flags & 3) === 3 || (flags & 12) === 12) {
            throw new WasiError(28);
        }
    }

    function vectors(pointer, count) {
        count >>>= 0;
        const table = data(pointer, count * 8);
        const result = [];
        for (let index = 0; index < count; index++) {
            const address = table.getUint32(index * 8, true);
            const length = table.getUint32(index * 8 + 4, true);
            bytes(address, length);
            result.push({ address, length });
        }
        return result;
    }

    const wasi = {
        args_sizes_get: checked((count, size) =>
            stringsSizes(argumentsBytes, count, size),
        ),
        args_get: checked((pointers, buffer) =>
            stringsGet(argumentsBytes, pointers, buffer),
        ),
        environ_sizes_get: checked((count, size) =>
            stringsSizes(environmentBytes, count, size),
        ),
        environ_get: checked((pointers, buffer) =>
            stringsGet(environmentBytes, pointers, buffer),
        ),
        clock_res_get: checked((clock, pointer) => {
            if (clock !== 0 && clock !== 1) throw new WasiError(28);
            data(pointer, 8).setBigUint64(0, 1_000_000n, true);
        }),
        clock_time_get: checked((clock, _precision, pointer) => {
            let milliseconds;
            if (clock === 0) milliseconds = Date.now();
            else if (clock === 1) milliseconds = performance.now();
            else throw new WasiError(28);
            data(pointer, 8).setBigUint64(
                0,
                BigInt(Math.floor(milliseconds * 1_000_000)),
                true,
            );
        }),
        random_get: checked((pointer, length) => {
            const output = bytes(pointer, length >>> 0);
            for (let offset = 0; offset < output.length; offset += 65_536) {
                crypto.getRandomValues(
                    output.subarray(offset, offset + 65_536),
                );
            }
        }),
        fd_prestat_get: checked((fd, pointer) => {
            requireDescriptor(fd);
            if (fd !== 3) throw new WasiError(8);
            const output = data(pointer, 8);
            output.setUint32(0, 0, true);
            output.setUint32(4, 1, true);
        }),
        fd_prestat_dir_name: checked((fd, pointer, length) => {
            requireDescriptor(fd);
            if (fd !== 3) throw new WasiError(8);
            if (length >>> 0 < 1) throw new WasiError(37);
            bytes(pointer, length >>> 0)[0] = 47;
        }),
        fd_close: checked((fd) => {
            closeDescriptor(fd);
        }),
        fd_fdstat_get: checked((fd, pointer) => {
            const rights = requireDescriptor(fd);
            const output = data(pointer, 24);
            const info =
                fd <= 2
                    ? { type: 2, flags: 0 }
                    : { type: fileType(stat(fd)), flags: rights.flags };
            bytes(pointer, 24).fill(0);
            output.setUint8(0, info.type);
            output.setUint16(2, info.flags, true);
            output.setBigUint64(8, rights.base, true);
            output.setBigUint64(16, rights.inheriting, true);
        }),
        fd_fdstat_set_flags: checked((fd, flags) => {
            requireDescriptor(fd, 8n);
            if (flags & ~31) throw new WasiError(28);
            requireDescriptor(fd).flags = flags;
        }),
        fd_filestat_get: checked((fd, pointer) => {
            requireDescriptor(fd, 1n << 21n);
            bytes(pointer, 64);
            writeStat(pointer, stat(fd));
        }),
        fd_filestat_set_size: checked((fd, size) => {
            requireDescriptor(fd, 1n << 22n);
            const descriptor = requireDescriptor(fd);
            if (!descriptor.writable) throw new WasiError(8);
            if (descriptor.directory) throw new WasiError(31);
            fs.ftruncateSync(
                descriptor.provider,
                number(BigInt.asUintN(64, size)),
            );
        }),
        fd_filestat_set_times: checked((fd, access, modification, flags) => {
            requireDescriptor(fd, 1n << 23n);
            checkTimeFlags(flags);
            fs.futimesSync(
                requireDescriptor(fd).provider,
                ...times(
                    stat(fd),
                    BigInt.asUintN(64, access),
                    BigInt.asUintN(64, modification),
                    flags,
                ),
            );
        }),
        fd_seek: checked((fd, offset, whence, pointer) => {
            requireDescriptor(fd, 4n);
            const output = data(pointer, 8);
            output.setBigUint64(0, seek(fd, offset, whence), true);
        }),
        fd_tell: checked((fd, pointer) => {
            requireDescriptor(fd, 32n);
            const output = data(pointer, 8);
            output.setBigUint64(0, seek(fd, 0n, 1), true);
        }),
        fd_sync: checked((fd) => {
            const descriptor = requireDescriptor(fd, 16n);
            throwSyncError(descriptor);
            fs.fsyncSync(descriptor.provider);
        }),
        fd_datasync: checked((fd) => {
            const descriptor = requireDescriptor(fd, 1n);
            throwSyncError(descriptor);
            fs.fsyncSync(descriptor.provider);
        }),
        fd_read: checked((fd, pointer, count, result) => {
            requireDescriptor(fd, 2n);
            const output = data(result, 4);
            const buffers = vectors(pointer, count);
            let total = 0;
            for (const { address, length } of buffers) {
                if (fd === 0) break;
                for (let offset = 0; offset < length; ) {
                    const size = Math.min(length - offset, 65_536);
                    let received;
                    try {
                        received = read(fd, size);
                    } catch (error) {
                        if (total === 0 || nodeErrno(error) === undefined)
                            throw error;
                        output.setUint32(0, total, true);
                        return;
                    }
                    bytes(address + offset, received.length).set(received);
                    total += received.length;
                    offset += received.length;
                    if (received.length < size) {
                        output.setUint32(0, total, true);
                        return;
                    }
                }
            }
            output.setUint32(0, total, true);
        }),
        fd_write: checked((fd, pointer, count, result) => {
            requireDescriptor(fd, 64n);
            const output = data(result, 4);
            const buffers = vectors(pointer, count);
            let total = 0;
            for (const { address, length } of buffers) {
                for (let offset = 0; offset < length; ) {
                    const source = bytes(
                        address + offset,
                        Math.min(length - offset, 65_536),
                    );
                    let written;
                    if (fd === 1 || fd === 2) {
                        const text = new TextDecoder().decode(source);
                        if (fd === 2) console.error(text);
                        else console.log(text);
                        written = source.length;
                    } else {
                        try {
                            written = write(fd, source);
                        } catch (error) {
                            if (total === 0 || nodeErrno(error) === undefined)
                                throw error;
                            output.setUint32(0, total, true);
                            return;
                        }
                    }
                    total += written;
                    offset += written;
                    if (
                        written < source.length ||
                        requireDescriptor(fd).syncError
                    ) {
                        output.setUint32(0, total, true);
                        return;
                    }
                }
            }
            output.setUint32(0, total, true);
        }),
        fd_readdir: checked((fd, pointer, length, cookie, result) => {
            requireDescriptor(fd, 1n << 14n);
            const output = bytes(pointer, length >>> 0);
            const count = data(result, 4);
            cookie = BigInt.asUintN(64, cookie);
            const entries = readdir(fd, cookie === 0n);
            let used = 0;
            if (cookie <= BigInt(entries.length)) {
                for (
                    let index = Number(cookie);
                    index < entries.length;
                    index++
                ) {
                    const entry = entries[index];
                    const name = encoder.encode(entry.name);
                    const record = new Uint8Array(24 + name.length);
                    const fields = new DataView(record.buffer);
                    fields.setBigUint64(0, BigInt(index + 1), true);
                    fields.setBigUint64(8, BigInt(entry.inode), true);
                    fields.setUint32(16, name.length, true);
                    fields.setUint8(20, entry.type);
                    record.set(name, 24);
                    const size = Math.min(record.length, output.length - used);
                    output.set(record.subarray(0, size), used);
                    used += size;
                    if (size < record.length || used === output.length) break;
                }
            }
            count.setUint32(0, used, true);
        }),
        path_open: checked(
            (
                fd,
                lookup,
                pointer,
                length,
                flags,
                base,
                inheriting,
                fdFlags,
                result,
            ) => {
                const directory = requireDescriptor(fd, 1n << 13n);
                if (lookup & ~1 || flags & ~15 || fdFlags & ~31) {
                    throw new WasiError(28);
                }
                base = BigInt.asUintN(64, base);
                inheriting = BigInt.asUintN(64, inheriting);
                if ((base & directory.inheriting) !== base) {
                    throw new WasiError(76);
                }
                if ((inheriting & directory.inheriting) !== inheriting) {
                    throw new WasiError(76);
                }
                if (flags & 1) requireDescriptor(fd, 1n << 10n);
                const output = data(result, 4);
                const opened = open(
                    fd,
                    path(pointer, length),
                    lookup,
                    flags,
                    base,
                    inheriting,
                    fdFlags,
                );
                output.setUint32(0, opened, true);
            },
        ),
        path_filestat_get: checked((fd, lookup, pointer, length, result) => {
            requireDescriptor(fd, 1n << 18n);
            if (lookup & ~1) throw new WasiError(28);
            bytes(result, 64);
            const name = providerPath(resolve(fd, path(pointer, length)));
            writeStat(
                result,
                lookup & 1
                    ? fs.statSync(name, { bigint: true })
                    : fs.lstatSync(name, { bigint: true }),
            );
        }),
        path_filestat_set_times: checked(
            (fd, lookup, pointer, length, access, modification, flags) => {
                requireDescriptor(fd, 1n << 20n);
                if (lookup & ~1) throw new WasiError(28);
                checkTimeFlags(flags);
                const name = providerPath(resolve(fd, path(pointer, length)));
                const info = fs.lstatSync(name, { bigint: true });
                if (!(lookup & 1) && info.isSymbolicLink())
                    throw new WasiError(58);
                fs.utimesSync(
                    name,
                    ...times(
                        lookup & 1
                            ? fs.statSync(name, { bigint: true })
                            : info,
                        BigInt.asUintN(64, access),
                        BigInt.asUintN(64, modification),
                        flags,
                    ),
                );
            },
        ),
        path_create_directory: checked((fd, pointer, length) => {
            requireDescriptor(fd, 1n << 9n);
            fs.mkdirSync(providerPath(resolve(fd, path(pointer, length))), {
                mode: 0o700,
            });
        }),
        path_remove_directory: checked((fd, pointer, length) => {
            requireDescriptor(fd, 1n << 25n);
            fs.rmdirSync(providerPath(resolve(fd, path(pointer, length))));
        }),
        path_unlink_file: checked((fd, pointer, length) => {
            requireDescriptor(fd, 1n << 26n);
            fs.unlinkSync(providerPath(resolve(fd, path(pointer, length))));
        }),
        path_rename: checked(
            (oldFD, oldPath, oldLength, newFD, newPath, newLength) => {
                requireDescriptor(oldFD, 1n << 16n);
                requireDescriptor(newFD, 1n << 17n);
                const oldName = resolve(oldFD, path(oldPath, oldLength));
                const newName = resolve(newFD, path(newPath, newLength));
                fs.renameSync(providerPath(oldName), providerPath(newName));
                for (const descriptor of descriptors.values()) {
                    if (
                        descriptor.name === oldName ||
                        descriptor.name?.startsWith(`${oldName}/`)
                    ) {
                        descriptor.name =
                            newName + descriptor.name.slice(oldName.length);
                    }
                }
            },
        ),
        path_readlink: checked((fd, pointer, length, buffer, size, result) => {
            requireDescriptor(fd, 1n << 15n);
            bytes(buffer, size >>> 0);
            bytes(result, 4);
            const info = fs.lstatSync(
                providerPath(resolve(fd, path(pointer, length))),
            );
            throw new WasiError(info.isSymbolicLink() ? 58 : 28);
        }),
        path_link: () => 58,
        path_symlink: () => 58,
        poll_oneoff: () => 58,
        sched_yield: () => 0,
        proc_exit(code) {
            throw new CommandExit(code >>> 0);
        },
    };

    return {
        imports: {
            wasi_snapshot_preview1: wasi,
            anki_timezone: {
                offset_for_utc(seconds) {
                    const offset = timezone.offsetForUTC(seconds.toString());
                    if (!Number.isInteger(offset))
                        throw new TypeError("Invalid Go UTC timezone offset");
                    return offset;
                },
                offsets_for_local(seconds) {
                    const offsets = timezone.offsetsForLocal(
                        seconds.toString(),
                    );
                    if (typeof offsets !== "string")
                        throw new TypeError(
                            "Invalid Go local timezone offsets",
                        );
                    return BigInt(offsets);
                },
            },
        },
        close() {
            const failures = [];
            for (const fd of descriptors.keys()) {
                try {
                    closeDescriptor(fd);
                } catch (error) {
                    failures.push(error);
                }
            }
            if (failures.length)
                throw new AggregateError(
                    failures,
                    "Close WASI provider descriptors",
                );
        },
        bind(value) {
            if (!(value instanceof WebAssembly.Memory)) {
                throw new TypeError("The Anki command must export its memory");
            }
            memory = value;
        },
    };
}
