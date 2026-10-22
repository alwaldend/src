import { createAnkiApp } from "./index.js";

const root = document.querySelector("#anki-app");
const status = root.querySelector("#app-status");
const openConfig = root.querySelector("#open-config");
const importConfig = root.querySelector("#import-config");
const openCollection = root.querySelector("#open-collection");
const importCollection = root.querySelector("#import-collection");
const planCollection = root.querySelector("#plan-collection");
const applyCollection = root.querySelector("#apply-collection");
const downloadCollection = root.querySelector("#download-collection");
const saveWorkspace = root.querySelector("#save-workspace");
const restoreWorkspace = root.querySelector("#restore-workspace");
const collectionStatus = root.querySelector("#collection-status");
const collectionReport = root.querySelector("#collection-report");
const files = root.querySelector("#workspace-files");
const editor = root.querySelector("#file-content");
const checkConfig = root.querySelector("#check-config");
const downloadConfig = root.querySelector("#download-config");
const controls = root.querySelectorAll("button, input, select, textarea");
const encoder = new TextEncoder();
const decoder = new TextDecoder("utf-8", { fatal: true });
const workspace = "/workspace";
const collection = "/collection.colpkg";
const storageKey = "anki-as-code-workspace";
const downloadURLs = new Set();

const example = {
    config: { version: "1" },
    decks: [
        {
            path: "decks/Lesson/deck.anki.toml",
            deck: { normal: {} },
            notes: [
                {
                    path: "decks/Lesson/notes/greeting.note.anki.toml",
                    config: {
                        note_type: "note_types/Basic.notetype.anki.toml",
                        tags: ["browser"],
                        fields: {
                            Front: "<b>你好</b>\nSay hello.",
                            Back: "Hello",
                        },
                        cards: [{}],
                    },
                },
            ],
        },
    ],
    note_types: [
        {
            path: "note_types/Basic.notetype.anki.toml",
            config: {
                name: "Basic",
                fields: [{ name: "Front" }, { name: "Back" }],
                templates: [
                    {
                        name: "Card",
                        config: {
                            q_format: "{{Front}}",
                            a_format: "{{FrontSide}}<hr>{{Back}}",
                        },
                    },
                ],
            },
        },
    ],
};

let application;
let candidate;
let busy = false;
let closed = false;
let selectedPath = "";
let hasCollection = false;
let savedRevision = 0;

function setStatus(text, error = false) {
    status.textContent = text;
    status.classList.toggle("text-danger", error);
    status.classList.toggle("text-body-secondary", !error);
}

function updateControls() {
    const unavailable = busy || closed || !application;
    for (const control of controls) control.disabled = unavailable;
    for (const control of [
        planCollection,
        applyCollection,
        downloadCollection,
    ])
        control.disabled = unavailable || !hasCollection;
}

function clearReport() {
    collectionReport.replaceChildren();
    collectionReport.hidden = true;
}

async function perform(description, operation) {
    if (busy || closed) return;
    busy = true;
    setStatus(description);
    updateControls();
    try {
        await operation();
    } catch (error) {
        setStatus(
            error instanceof Error ? error.message : String(error),
            true,
        );
    } finally {
        busy = false;
        updateControls();
    }
}

async function readEditor(instance, path) {
    return decoder.decode(await instance.workspace.readFile(path));
}

function showEditor(path, contents) {
    selectedPath = path;
    files.value = path;
    editor.value = contents;
}

async function storeEditor() {
    if (!selectedPath) return;
    await application.workspace.writeFile(
        selectedPath,
        encoder.encode(editor.value),
    );
}

async function checkConfiguration() {
    await storeEditor();
    return application.app.configDump({
        text: workspace,
        format: "json",
    });
}

async function adoptCandidate(metadata) {
    const paths = await candidate.workspace.listFiles(workspace);
    const path =
        paths.find((path) => path.endsWith(".note.anki.toml")) || paths[0];
    if (!path)
        throw new Error(
            "The saved workspace contains no configuration files.",
        );
    const contents = await readEditor(candidate, path);
    const collectionPresent = (
        await candidate.workspace.listFiles("/")
    ).includes(collection);
    if (closed) return;
    const previous = application;
    application = candidate;
    candidate = undefined;
    if (metadata !== undefined) savedRevision = metadata?.revision || 0;
    hasCollection = collectionPresent;
    files.replaceChildren(
        ...paths.map(
            (path) => new Option(path.slice(workspace.length + 1), path),
        ),
    );
    showEditor(path, contents);
    collectionStatus.textContent = hasCollection
        ? "Collection open. Plan your changes before applying them."
        : "Open an Anki collection package to plan and apply changes.";
    clearReport();
    await previous?.dispose();
}

// Imports retain the page's saved revision; restores return new save metadata.
async function replaceWorkspace(prepare) {
    candidate = await createAnkiApp();
    try {
        if (closed) return;
        const metadata = await prepare(candidate);
        await adoptCandidate(metadata);
    } finally {
        await candidate?.dispose();
        candidate = undefined;
    }
}

async function openConfiguration(configs) {
    await replaceWorkspace(async (instance) => {
        await instance.app.configExport({ output: workspace, configs });
        const paths = await instance.workspace.listFiles("/");
        if (paths.some((path) => !path.startsWith(`${workspace}/`)))
            throw new Error(
                "Imported files must not be outside the configuration directory.",
            );
        await instance.app.configDump({ text: workspace, format: "json" });
        if (hasCollection)
            await instance.workspace.writeFile(
                collection,
                await application.workspace.readFile(collection),
            );
    });
    setStatus(
        hasCollection
            ? "Configuration opened. Your collection is kept."
            : "Configuration opened and ready to edit.",
    );
}

async function restoreSaved() {
    await replaceWorkspace(async (instance) => {
        const metadata = await instance.session.restore(storageKey);
        if (!metadata)
            throw new Error(
                "There is no saved workspace in this browser yet.",
            );
        return metadata;
    });
    setStatus("Saved workspace restored.");
}

async function start() {
    try {
        await replaceWorkspace(async (instance) => {
            const metadata = await instance.session.restore(storageKey);
            if (!metadata)
                await instance.app.configExport({
                    output: workspace,
                    configs: example,
                });
            return metadata;
        });
        setStatus(
            savedRevision
                ? "Saved workspace restored."
                : "Example ready to edit.",
        );
    } catch (error) {
        await replaceWorkspace((instance) =>
            instance.app.configExport({ output: workspace, configs: example }),
        );
        setStatus(
            `Could not restore saved work: ${error.message || error}. The example is available; your saved work has not been replaced.`,
            true,
        );
    }
}

function showPlan(plan) {
    const changes = plan.changes || [];
    collectionReport.replaceChildren();
    const summary = document.createElement("p");
    summary.className = "mb-2";
    summary.textContent = changes.length
        ? `${changes.length} collection changes planned.`
        : plan.state_changed
          ? "Collection bindings will be updated."
          : "No collection changes needed.";
    collectionReport.append(summary);
    if (changes.length) {
        const list = document.createElement("ul");
        list.className = "mb-2";
        for (const change of changes) {
            const item = document.createElement("li");
            item.className = "text-break";
            item.textContent = `${change.action} ${change.resource}: ${change.path}`;
            if (change.from) item.textContent += ` (from ${change.from})`;
            list.append(item);
        }
        collectionReport.append(list);
    }
    if (plan.card_generation_deferred) {
        const note = document.createElement("p");
        note.className = "mb-2";
        note.textContent =
            "Anki determines card generation when changes are applied.";
        collectionReport.append(note);
    }
    collectionReport.hidden = false;
}

function download(name, contents, type) {
    const url = URL.createObjectURL(new Blob([contents], { type }));
    downloadURLs.add(url);
    const link = document.createElement("a");
    link.href = url;
    link.download = name;
    document.body.append(link);
    link.click();
    link.remove();
    setTimeout(() => {
        URL.revokeObjectURL(url);
        downloadURLs.delete(url);
    }, 1000);
}

openConfig.addEventListener("click", () => importConfig.click());
openCollection.addEventListener("click", () => importCollection.click());

importCollection.addEventListener("change", () => {
    const file = importCollection.files[0];
    if (!file) return;
    void perform("Opening collection…", async () => {
        try {
            if (!file.name.toLowerCase().endsWith(".colpkg"))
                throw new Error(
                    "Choose an Anki collection package (.colpkg).",
                );
            const bytes = new Uint8Array(await file.arrayBuffer());
            await replaceWorkspace(async (instance) => {
                const state = "/import-bindings.toml";
                await instance.workspace.writeFile(collection, bytes);
                await instance.app.collectionExport({
                    collection,
                    output: workspace,
                    stateOutput: state,
                });
                await instance.app.collectionStateUpdate({
                    collection,
                    state,
                });
                await instance.workspace.remove(state);
            });
            setStatus(
                "Collection opened. Its configuration is ready to edit.",
            );
        } finally {
            importCollection.value = "";
        }
    });
});

importConfig.addEventListener("change", () => {
    const file = importConfig.files[0];
    if (!file) return;
    void perform("Opening configuration…", async () => {
        try {
            await openConfiguration(JSON.parse(await file.text()));
        } finally {
            importConfig.value = "";
        }
    });
});

files.addEventListener("change", () => {
    const path = files.value;
    void perform("Loading file…", async () => {
        try {
            await storeEditor();
            showEditor(path, await readEditor(application, path));
            setStatus("File loaded.");
        } catch (error) {
            files.value = selectedPath;
            throw error;
        }
    });
});

editor.addEventListener("input", () => {
    clearReport();
    setStatus("Unsaved changes. Save locally to keep them after reloading.");
});

checkConfig.addEventListener("click", () => {
    void perform("Checking configuration…", async () => {
        await checkConfiguration();
        setStatus("Configuration format is valid.");
    });
});

downloadConfig.addEventListener("click", () => {
    void perform("Preparing configuration download…", async () => {
        const json = await checkConfiguration();
        download("configuration.json", json, "application/json;charset=utf-8");
        setStatus("Configuration JSON download prepared.");
    });
});

planCollection.addEventListener("click", () => {
    void perform("Planning collection changes…", async () => {
        await storeEditor();
        const plan = await application.app.collectionPlan({
            collection,
            text: workspace,
        });
        showPlan(plan);
        setStatus("Plan ready. Review the changes before applying.");
    });
});

applyCollection.addEventListener("click", () => {
    void perform("Applying collection changes…", async () => {
        await storeEditor();
        const output = "/applied.colpkg";
        try {
            await application.app.collectionBuild({
                base: collection,
                text: workspace,
                output,
            });
            await application.workspace.rename(output, collection);
        } finally {
            await application.workspace.remove(output);
        }
        clearReport();
        setStatus(
            "Changes applied. Save locally or download the collection to keep them.",
        );
    });
});

downloadCollection.addEventListener("click", () => {
    void perform("Preparing collection download…", async () => {
        const bytes = await application.workspace.readFile(collection);
        download("collection.colpkg", bytes, "application/octet-stream");
        setStatus(
            "Collection download prepared. Unapplied configuration edits are not included.",
        );
    });
});

saveWorkspace.addEventListener("click", () => {
    void perform("Saving workspace locally…", async () => {
        await storeEditor();
        const metadata = await application.session.save(
            storageKey,
            savedRevision,
        );
        savedRevision = metadata.revision;
        setStatus("Workspace saved locally. Apply is separate from saving.");
    });
});

restoreWorkspace.addEventListener("click", () => {
    void perform("Restoring saved workspace…", restoreSaved);
});

window.addEventListener("pagehide", (event) => {
    if (event.persisted) return;
    closed = true;
    updateControls();
    if (application) void application.dispose();
    if (candidate) void candidate.dispose();
    for (const url of downloadURLs) URL.revokeObjectURL(url);
    downloadURLs.clear();
});

await perform("Loading saved workspace…", start);
