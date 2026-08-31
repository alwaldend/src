import { spawn } from "node:child_process";
import { readdir, readFile } from "node:fs/promises";
import process from "node:process";

const PROCESS_POLL_INTERVAL_MS = 10;
const PROCESS_INSPECTION_FAILURE_LIMIT = 3;

function codedError(code, message, ErrorType = Error) {
    const error = new ErrorType(message);
    error.code = code;
    return error;
}

function deferred() {
    let resolve;
    let reject;
    const promise = new Promise((resolvePromise, rejectPromise) => {
        resolve = resolvePromise;
        reject = rejectPromise;
    });
    return { promise, reject, resolve };
}

function validUtf8PrefixLength(buffer) {
    let index = 0;
    while (index < buffer.length) {
        const first = buffer[index];
        if (first <= 0x7f) {
            index += 1;
            continue;
        }

        let width;
        let secondMinimum = 0x80;
        let secondMaximum = 0xbf;
        if (first >= 0xc2 && first <= 0xdf) {
            width = 2;
        } else if (first === 0xe0) {
            width = 3;
            secondMinimum = 0xa0;
        } else if (first >= 0xe1 && first <= 0xec) {
            width = 3;
        } else if (first === 0xed) {
            width = 3;
            secondMaximum = 0x9f;
        } else if (first >= 0xee && first <= 0xef) {
            width = 3;
        } else if (first === 0xf0) {
            width = 4;
            secondMinimum = 0x90;
        } else if (first >= 0xf1 && first <= 0xf3) {
            width = 4;
        } else if (first === 0xf4) {
            width = 4;
            secondMaximum = 0x8f;
        } else {
            return index;
        }

        if (index + width > buffer.length) return index;
        const second = buffer[index + 1];
        if (second < secondMinimum || second > secondMaximum) return index;
        for (let offset = 2; offset < width; offset += 1) {
            const continuation = buffer[index + offset];
            if (continuation < 0x80 || continuation > 0xbf) return index;
        }
        index += width;
    }
    return index;
}

function decodeUtf8Prefix(output) {
    const buffer = Buffer.concat(output.parts, output.bytes);
    const length = validUtf8PrefixLength(buffer);
    return {
        dropped: length < buffer.length,
        text: buffer.subarray(0, length).toString("utf8"),
    };
}

function parseProcStat(stat) {
    const closingParenthesis = stat.lastIndexOf(")");
    if (closingParenthesis < 0) return undefined;
    const fields = stat
        .slice(closingParenthesis + 2)
        .trim()
        .split(/\s+/u);
    if (fields.length < 3) return undefined;
    const group = Number(fields[2]);
    if (!Number.isSafeInteger(group) || group <= 0) return undefined;
    return { group, state: fields[0] };
}

async function processGroupHasLiveMember(group) {
    const entries = await readdir("/proc", { withFileTypes: true });
    for (const entry of entries) {
        if (!entry.isDirectory() || !/^\d+$/u.test(entry.name)) continue;
        const stat = await readFile(`/proc/${entry.name}/stat`, "utf8").catch(
            (error) => {
                if (["EACCES", "ENOENT", "EPERM"].includes(error?.code)) {
                    return undefined;
                }
                throw error;
            },
        );
        if (stat === undefined) continue;
        const parsed = parseProcStat(stat);
        if (parsed?.group !== group) continue;
        if (parsed.state !== "Z" && parsed.state !== "X") return true;
    }
    return false;
}

function signalProcessGroup(child) {
    if (!Number.isSafeInteger(child.pid) || child.pid <= 0) return undefined;
    try {
        process.kill(-child.pid, "SIGKILL");
        return undefined;
    } catch (error) {
        if (error?.code === "ESRCH") return undefined;
        if (child.exitCode === null && child.signalCode === null) {
            try {
                child.kill("SIGKILL");
            } catch {
                // Preserve the process-group failure as the useful error.
            }
        }
        return error;
    }
}

export async function waitForProcessGroup(
    group,
    child,
    inspect = processGroupHasLiveMember,
) {
    if (!Number.isSafeInteger(group) || group <= 0) return;
    let inspectionFailures = 0;
    while (true) {
        signalProcessGroup(child);
        try {
            if (!(await inspect(group))) return;
            inspectionFailures = 0;
        } catch (error) {
            inspectionFailures += 1;
            if (inspectionFailures >= PROCESS_INSPECTION_FAILURE_LIMIT) {
                return codedError(
                    "EXEC_CLEANUP",
                    `could not verify process-group cleanup: ${error.message}`,
                );
            }
        }
        await new Promise((resolve) => {
            setTimeout(resolve, PROCESS_POLL_INTERVAL_MS);
        });
    }
}

function destroyOutput(child) {
    child.stdout?.destroy();
    child.stderr?.destroy();
}

export class ProcessSupervisor {
    #accepting = true;
    #active = new Set();

    execute({ file, args, options }) {
        if (!this.#accepting) {
            return Promise.reject(
                codedError(
                    "EXEC_DISPOSED",
                    "ctx.exec is unavailable because the package is disposing",
                ),
            );
        }

        let child;
        try {
            child = spawn(file, args, {
                cwd: options.cwd,
                detached: true,
                env: options.env,
                shell: false,
                stdio: ["ignore", "pipe", "pipe"],
            });
        } catch (error) {
            return Promise.reject(error);
        }

        const outcome = deferred();
        const done = deferred();
        const record = {
            child,
            done: done.promise,
            exitCode: null,
            exitSignal: null,
            finalizeStarted: false,
            forcedError: undefined,
            groupStopped: undefined,
            outputLimitExceeded: false,
            promise: outcome.promise,
            spawnError: undefined,
            stderr: { bytes: 0, parts: [] },
            stdout: { bytes: 0, parts: [] },
            outputBytes: 0,
        };
        this.#active.add(record);

        const waitForGroupStopped = () => {
            record.groupStopped ??= waitForProcessGroup(child.pid, child);
            return record.groupStopped;
        };

        const stop = (error = undefined) => {
            if (error !== undefined) record.forcedError ??= error;
            signalProcessGroup(child);
            destroyOutput(child);
        };
        record.stop = stop;

        const collect = (destination, chunk) => {
            if (record.forcedError || record.outputLimitExceeded) return;
            const retainedBytes = Math.min(
                chunk.length,
                options.maxBytes - record.outputBytes,
            );
            if (retainedBytes > 0) {
                destination.parts.push(
                    Buffer.from(chunk.subarray(0, retainedBytes)),
                );
                destination.bytes += retainedBytes;
                record.outputBytes += retainedBytes;
            }
            if (retainedBytes === chunk.length) return;

            record.outputLimitExceeded = true;
            if (!options.allowTruncatedOutput) {
                record.forcedError = codedError(
                    "EXEC_OUTPUT_LIMIT",
                    `process output exceeded the ` +
                        `${options.maxBytes}-byte limit`,
                    RangeError,
                );
            }
            stop();
        };

        const finalize = async () => {
            if (record.finalizeStarted) return;
            record.finalizeStarted = true;
            clearTimeout(timeoutTimer);
            const cleanupError = await waitForGroupStopped();
            this.#active.delete(record);

            const decodedStdout = decodeUtf8Prefix(record.stdout);
            const decodedStderr = decodeUtf8Prefix(record.stderr);
            const invalidUtf8 = decodedStdout.dropped || decodedStderr.dropped;
            const error =
                record.spawnError ??
                record.forcedError ??
                cleanupError ??
                (invalidUtf8 && !options.allowTruncatedOutput
                    ? codedError(
                          "EXEC_INVALID_UTF8",
                          "process output was not valid UTF-8",
                      )
                    : undefined);

            if (error !== undefined) {
                outcome.reject(error);
            } else {
                outcome.resolve({
                    code: record.exitCode,
                    signal: record.exitSignal,
                    stdout: decodedStdout.text,
                    stderr: decodedStderr.text,
                    truncated: record.outputLimitExceeded || invalidUtf8,
                    outputLimitExceeded: record.outputLimitExceeded,
                });
            }
            done.resolve();
        };

        const timeoutTimer = setTimeout(() => {
            stop(
                codedError(
                    "EXEC_TIMEOUT",
                    `process exceeded the ${options.timeoutMs} ms timeout`,
                ),
            );
        }, options.timeoutMs);
        timeoutTimer.unref();

        child.stdout?.on("data", (chunk) => collect(record.stdout, chunk));
        child.stderr?.on("data", (chunk) => collect(record.stderr, chunk));
        const outputError = (error) => {
            stop(
                codedError(
                    "EXEC_OUTPUT_ERROR",
                    `failed to read process output: ${error.message}`,
                ),
            );
        };
        child.stdout?.on("error", outputError);
        child.stderr?.on("error", outputError);
        child.once("error", (error) => {
            record.spawnError = error;
            stop();
        });
        child.once("exit", (code, signal) => {
            record.exitCode = code;
            record.exitSignal = signal;
            signalProcessGroup(child);
            void waitForGroupStopped().then(async () => {
                await new Promise((resolve) => setImmediate(resolve));
                destroyOutput(child);
            });
        });
        child.once("close", (code, signal) => {
            record.exitCode = code;
            record.exitSignal = signal;
            void finalize();
        });

        return record.promise;
    }

    async close(error) {
        this.#accepting = false;
        const records = [...this.#active];
        for (const record of records) record.stop(error);
        await Promise.allSettled(records.map((record) => record.done));
    }
}
