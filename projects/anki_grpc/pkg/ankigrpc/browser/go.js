/**
 * Install a borrowed Node-compatible filesystem for Go js/wasm before Go.run.
 * Import wasm_exec.js first so its runtime output sink is available. Keep these
 * globals and the provider stable for the Go runtime's lifetime. Paths use cwd
 * '/'; callers select their workspace with BasePathFs, not with this helper.
 * This facade supplies no storage and reserves stdout/stderr for runtime output.
 */
export function installGoFilesystem(fs) {
    const output = globalThis.fs;
    const process = globalThis.process;
    if (!output || typeof output.writeSync !== "function" || !process) {
        throw new Error(
            "Import wasm_exec.js before installing the Go filesystem",
        );
    }
    if (!fs || !fs.constants || typeof fs.open !== "function") {
        throw new TypeError(
            "A Node-compatible callback filesystem is required",
        );
    }
    const pathArguments = new Map([
        ...[
            "open",
            "stat",
            "lstat",
            "readdir",
            "mkdir",
            "rmdir",
            "unlink",
            "chmod",
            "chown",
            "lchown",
            "utimes",
            "truncate",
            "readlink",
        ].map((name) => [name, [0]]),
        ["rename", [0, 1]],
        ["link", [0, 1]],
        ["symlink", [1]],
    ]);
    globalThis.fs = new Proxy(fs, {
        get(target, name) {
            if (name === "writeSync") {
                return (fd, buffer, ...args) =>
                    fd === 1 || fd === 2
                        ? output.writeSync(fd, buffer, ...args)
                        : target.writeSync(fd, buffer, ...args);
            }
            if (name === "write") {
                return (fd, buffer, offset, length, position, callback) => {
                    if (fd !== 1 && fd !== 2)
                        return target.write(
                            fd,
                            buffer,
                            offset,
                            length,
                            position,
                            callback,
                        );
                    let count;
                    try {
                        count = output.writeSync(
                            fd,
                            buffer.subarray(offset, offset + length),
                        );
                    } catch (error) {
                        callback(error);
                        return;
                    }
                    callback(null, count, buffer);
                };
            }
            const value = Reflect.get(target, name, target);
            if (typeof value !== "function") return value;
            return (...args) => {
                for (const index of pathArguments.get(name) || [])
                    args[index] = resolve(args[index]);
                return Reflect.apply(value, target, args);
            };
        },
    });
    globalThis.path = { resolve };
    globalThis.process = new Proxy(process, {
        get(target, name) {
            if (name === "cwd") return () => "/";
            if (name === "chdir")
                return () => {
                    throw Object.assign(
                        new Error("Go filesystem cwd is fixed at /"),
                        { code: "ENOSYS" },
                    );
                };
            const value = Reflect.get(target, name, target);
            return typeof value === "function" ? value.bind(target) : value;
        },
    });
}

/** Resolve POSIX paths against a fixed root, matching Node path.posix.resolve. */
function resolve(...paths) {
    const parts = [];
    for (const path of paths) {
        if (typeof path !== "string")
            throw new TypeError("Expected a string path");
        if (path.startsWith("/")) parts.length = 0;
        for (const part of path.split("/")) {
            if (part === "..") parts.pop();
            else if (part && part !== ".") parts.push(part);
        }
    }
    return `/${parts.join("/")}`;
}
