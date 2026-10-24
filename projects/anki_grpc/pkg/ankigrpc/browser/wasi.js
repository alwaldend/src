import { createServer } from "./server.js";

/** Initialize a callable Rust backend over the caller's shared Afero store. */
export async function createBrowserServer(module, filesystem) {
    const host = createWASI(filesystem);
    const instance = await WebAssembly.instantiate(module, host.imports);
    host.bind(instance.exports.memory);
    instance.exports._initialize();
    return createServer(instance);
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

/** Marshal Preview 1 imports into synchronous, memory-only Afero operations. */
function createWASI(filesystem) {
    if (!filesystem || typeof filesystem.call !== "function") {
        throw new TypeError("A Go Afero filesystem endpoint is required");
    }
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
        [3, { base: allFileRights, inheriting: allFileRights }],
    ]);
    let memory;

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

    function call(operation, ...args) {
        const result = filesystem.call(operation, ...args);
        if (result.errno) throw new WasiError(result.errno);
        return result.value;
    }

    function checked(operation) {
        return (...args) => {
            try {
                operation(...args);
                return 0;
            } catch (error) {
                if (error instanceof WasiError) return error.errno;
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
        output.setBigUint64(0, 1n, true);
        output.setBigUint64(8, BigInt(value.inode), true);
        output.setUint8(16, value.type);
        output.setBigUint64(24, 1n, true);
        output.setBigUint64(32, BigInt(value.size), true);
        for (const offset of [40, 48, 56]) {
            output.setBigUint64(offset, BigInt(value.time), true);
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
            requireDescriptor(fd);
            if (fd > 2) call("close", fd);
            descriptors.delete(fd);
        }),
        fd_fdstat_get: checked((fd, pointer) => {
            const rights = requireDescriptor(fd);
            const output = data(pointer, 24);
            const info = fd <= 2 ? { type: 2, flags: 0 } : call("fdstat", fd);
            bytes(pointer, 24).fill(0);
            output.setUint8(0, info.type);
            output.setUint16(2, info.flags, true);
            output.setBigUint64(8, rights.base, true);
            output.setBigUint64(16, rights.inheriting, true);
        }),
        fd_fdstat_set_flags: checked((fd, flags) => {
            requireDescriptor(fd, 8n);
            if (flags & ~31) throw new WasiError(28);
            call("flags", fd, flags);
        }),
        fd_filestat_get: checked((fd, pointer) => {
            requireDescriptor(fd, 1n << 21n);
            bytes(pointer, 64);
            writeStat(pointer, call("stat", fd));
        }),
        fd_filestat_set_size: checked((fd, size) => {
            requireDescriptor(fd, 1n << 22n);
            call("truncate", fd, BigInt.asUintN(64, size).toString());
        }),
        fd_filestat_set_times: checked((fd, access, modification, flags) => {
            requireDescriptor(fd, 1n << 23n);
            checkTimeFlags(flags);
            call(
                "times",
                fd,
                BigInt.asUintN(64, access).toString(),
                BigInt.asUintN(64, modification).toString(),
                flags,
            );
        }),
        fd_seek: checked((fd, offset, whence, pointer) => {
            requireDescriptor(fd, 4n);
            const output = data(pointer, 8);
            output.setBigUint64(
                0,
                BigInt(call("seek", fd, offset.toString(), whence)),
                true,
            );
        }),
        fd_tell: checked((fd, pointer) => {
            requireDescriptor(fd, 32n);
            const output = data(pointer, 8);
            output.setBigUint64(0, BigInt(call("seek", fd, "0", 1)), true);
        }),
        fd_sync: checked((fd) => {
            requireDescriptor(fd, 16n);
            call("sync", fd);
        }),
        fd_datasync: checked((fd) => {
            requireDescriptor(fd, 1n);
            call("sync", fd);
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
                    const received = call("read", fd, size);
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
                        written = call("write", fd, source);
                    }
                    total += written;
                    offset += written;
                    if (written < source.length) {
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
            const entries = call("readdir", fd, cookie === 0n);
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
                const opened = call(
                    "open",
                    fd,
                    path(pointer, length),
                    flags,
                    Boolean(base & 2n),
                    Boolean(base & 64n),
                    fdFlags,
                );
                descriptors.set(opened, { base, inheriting });
                output.setUint32(0, opened, true);
            },
        ),
        path_filestat_get: checked((fd, lookup, pointer, length, result) => {
            requireDescriptor(fd, 1n << 18n);
            if (lookup & ~1) throw new WasiError(28);
            bytes(result, 64);
            writeStat(result, call("pathStat", fd, path(pointer, length)));
        }),
        path_filestat_set_times: checked(
            (fd, lookup, pointer, length, access, modification, flags) => {
                requireDescriptor(fd, 1n << 20n);
                if (lookup & ~1) throw new WasiError(28);
                checkTimeFlags(flags);
                call(
                    "pathTimes",
                    fd,
                    path(pointer, length),
                    BigInt.asUintN(64, access).toString(),
                    BigInt.asUintN(64, modification).toString(),
                    flags,
                );
            },
        ),
        path_create_directory: checked((fd, pointer, length) => {
            requireDescriptor(fd, 1n << 9n);
            call("mkdir", fd, path(pointer, length));
        }),
        path_remove_directory: checked((fd, pointer, length) => {
            requireDescriptor(fd, 1n << 25n);
            call("removeDirectory", fd, path(pointer, length));
        }),
        path_unlink_file: checked((fd, pointer, length) => {
            requireDescriptor(fd, 1n << 26n);
            call("unlink", fd, path(pointer, length));
        }),
        path_rename: checked(
            (oldFD, oldPath, oldLength, newFD, newPath, newLength) => {
                requireDescriptor(oldFD, 1n << 16n);
                requireDescriptor(newFD, 1n << 17n);
                call(
                    "rename",
                    oldFD,
                    path(oldPath, oldLength),
                    newFD,
                    path(newPath, newLength),
                );
            },
        ),
        path_readlink: checked((fd, pointer, length, buffer, size, result) => {
            requireDescriptor(fd, 1n << 15n);
            bytes(buffer, size >>> 0);
            bytes(result, 4);
            call("pathStat", fd, path(pointer, length));
            // The shared Afero memory filesystem contains no symbolic links.
            throw new WasiError(28);
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
                    return call("timezoneUTC", seconds.toString());
                },
                offsets_for_local(seconds) {
                    return BigInt(call("timezoneLocal", seconds.toString()));
                },
            },
        },
        bind(value) {
            if (!(value instanceof WebAssembly.Memory)) {
                throw new TypeError("The Anki command must export its memory");
            }
            memory = value;
        },
    };
}
