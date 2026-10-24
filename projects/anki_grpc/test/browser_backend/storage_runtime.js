/** Run the same generated-client storage workflow in Node and Chromium. */
export async function runStorage({
    fs,
    root,
    module,
    application,
    createBrowserServer,
    installGoFilesystem,
}) {
    const evidence = {
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
        hostCalls: {},
        timezoneOffsets: [],
        instances: 0,
        remainingDescriptors: 0,
    };
    let running;
    try {
        installGoFilesystem(fs);
        WebAssembly.Suspending = undefined;
        WebAssembly.promising = undefined;
        globalThis.ankiBrowserBackendStorageMode = true;
        globalThis.ankiBrowserBackendRoot = root;
        globalThis.ankiBrowserBackendFactory = async (timezone) => {
            const offsets = [
                "2026-01-15T12:00:00Z",
                "2026-07-15T12:00:00Z",
            ].map((date) =>
                timezone.offsetForUTC(String(Date.parse(date) / 1000)),
            );
            evidence.timezoneOffsets.push(offsets);
            const observedTimezone = {
                offsetForUTC(seconds) {
                    evidence.hostCalls.timezoneUTC =
                        (evidence.hostCalls.timezoneUTC || 0) + 1;
                    return timezone.offsetForUTC(seconds);
                },
                offsetsForLocal(seconds) {
                    evidence.hostCalls.timezoneLocal =
                        (evidence.hostCalls.timezoneLocal || 0) + 1;
                    return timezone.offsetsForLocal(seconds);
                },
            };
            const observed = new Proxy(fs, {
                get(target, name) {
                    const value = target[name];
                    if (typeof value !== "function") return value;
                    return (...args) => {
                        evidence.hostCalls[name] =
                            (evidence.hostCalls[name] || 0) + 1;
                        const result = Reflect.apply(value, target, args);
                        if (name === "openSync")
                            evidence.remainingDescriptors++;
                        if (name === "closeSync")
                            evidence.remainingDescriptors--;
                        return result;
                    };
                },
            });
            const server = await createBrowserServer(module, observed, {
                root,
                timezone: observedTimezone,
            });
            evidence.instances++;
            return server;
        };
        let verified;
        const verification = new Promise((resolve) => {
            verified = resolve;
        });
        globalThis.ankiBrowserBackendVerified = () => {
            evidence.workflow = JSON.parse(
                globalThis.ankiBrowserBackendReport,
            );
            verified();
        };
        const go = new Go();
        go.exit = (code) => {
            evidence.exitCode = code;
        };
        const { instance } = await WebAssembly.instantiate(
            application,
            go.importObject,
        );
        running = go.run(instance);
        await Promise.race([
            verification,
            running.then(() => {
                if (!evidence.workflow)
                    throw new Error("Go exited before reporting verification");
            }),
        ]);
        globalThis.ankiBrowserBackendFinish();
        await running;
    } catch (error) {
        evidence.error = error.stack || error.message;
    }
    return evidence;
}
