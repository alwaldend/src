// Persist complete, closed workspaces in one atomic IndexedDB record. Keeping
// the previous revision in the transaction prevents stale tabs overwriting it.
const databaseName = "ankiascode-browser";
const storeName = "workspaces";

async function openDatabase() {
    if (typeof indexedDB === "undefined") {
        throw new Error(
            "Browser storage is unavailable; download your workspace instead",
        );
    }
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(databaseName, 1);
        let blocked = false;
        request.onupgradeneeded = () => {
            request.result.createObjectStore(storeName);
        };
        request.onerror = () => reject(request.error);
        request.onblocked = () => {
            blocked = true;
            reject(new Error("Browser storage is blocked by another tab"));
        };
        request.onsuccess = () => {
            if (blocked) {
                request.result.close();
                return;
            }
            request.result.onversionchange = () => request.result.close();
            resolve(request.result);
        };
    });
}

function validateSnapshot(value) {
    if (
        !value ||
        value.version !== 1 ||
        !Number.isSafeInteger(value.revision) ||
        value.revision < 1 ||
        typeof value.savedAt !== "string" ||
        !Number.isFinite(Date.parse(value.savedAt)) ||
        !Array.isArray(value.entries)
    ) {
        throw new Error(
            "Saved workspace has an unsupported or invalid format",
        );
    }
}

export async function loadWorkspace(key) {
    const database = await openDatabase();
    try {
        return await new Promise((resolve, reject) => {
            const transaction = database.transaction(storeName, "readonly");
            const request = transaction.objectStore(storeName).get(key);
            transaction.onabort = () =>
                reject(
                    transaction.error ??
                        new Error("Read saved workspace was aborted"),
                );
            transaction.onerror = () => {};
            transaction.oncomplete = () => {
                if (request.result === undefined) {
                    resolve(null);
                    return;
                }
                try {
                    validateSnapshot(request.result);
                    resolve(request.result);
                } catch (error) {
                    reject(error);
                }
            };
        });
    } finally {
        database.close();
    }
}

export async function saveWorkspace(key, expectedRevision, entries) {
    const database = await openDatabase();
    try {
        return await new Promise((resolve, reject) => {
            const transaction = database.transaction(storeName, "readwrite");
            const store = transaction.objectStore(storeName);
            const request = store.get(key);
            let failure;
            let snapshot;
            transaction.onabort = () =>
                reject(
                    failure ??
                        transaction.error ??
                        new Error("Save workspace was aborted"),
                );
            transaction.onerror = () => {};
            transaction.oncomplete = () =>
                resolve({
                    revision: snapshot.revision,
                    savedAt: snapshot.savedAt,
                });
            request.onsuccess = () => {
                try {
                    if (request.result !== undefined)
                        validateSnapshot(request.result);
                    const current = request.result?.revision ?? 0;
                    if (current !== expectedRevision) {
                        throw new Error(
                            "Saved workspace changed in another tab; restore it before saving",
                        );
                    }
                    if (current === Number.MAX_SAFE_INTEGER) {
                        throw new Error(
                            "Saved workspace revision limit reached",
                        );
                    }
                    snapshot = {
                        version: 1,
                        revision: current + 1,
                        savedAt: new Date().toISOString(),
                        entries,
                    };
                    store.put(snapshot, key);
                } catch (error) {
                    failure = error;
                    transaction.abort();
                }
            };
        });
    } finally {
        database.close();
    }
}
