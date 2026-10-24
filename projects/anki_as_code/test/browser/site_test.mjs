import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import fs from "node:fs/promises";
import http from "node:http";
import { createRequire } from "node:module";
import path from "node:path";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const [root, executablePath, fixtureTool] = process.argv
    .slice(2)
    .map((filename) => path.resolve(filename));
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(output, "Bazel must provide an artifact directory");
const scratch = await fs.mkdtemp(path.join(process.env.TEST_TMPDIR, "site-"));
const downloadDirectory = path.join(scratch, "downloads");
await fs.mkdir(downloadDirectory);
// Failure cases: invalid replacements retain edits and collection bytes;
// saving includes the live editor without applying; stale tabs cannot overwrite
// newer snapshots; failed Apply leaves the package intact; restored packages
// retain native identities, scheduling and binary media. Escaped JSON paths
// must neither create a collection nor overwrite a retained package, including
// after save/reload; failed replacement must preserve the current editor.
// Large restored workspaces must expose every file without leaving the editor
// attached to the previous worker or leaking that worker after replacement.
const fixtureArchive = path.join(output, "site-native-input.colpkg");
const fixtureBaseline = path.join(output, "site-native-baseline.json");
const fixture = spawnSync(
    fixtureTool,
    ["create", fixtureArchive, fixtureBaseline],
    { encoding: "utf8", timeout: 120000 },
);
assert.equal(fixture.status, 0, fixture.stderr || fixture.error?.message);
const landingPath = "/projects/anki_as_code/";
const appPath = `${landingPath}app/`;
const report = { cases: [], errors: [], blockedRequests: [], downloads: [] };
report.workerEvents = [];
report.failedRequests = [];
const requests = new Set();
let storageUnavailable = false;

// Permit only CDN scripts declared with integrity in the actual built pages.
const scripts = new Set();
for (const pagePath of [landingPath, appPath]) {
    const html = await fs.readFile(
        path.join(root, pagePath, "index.html"),
        "utf8",
    );
    for (const [tag] of html.matchAll(/<script\b[^>]*>/g)) {
        const url =
            tag.includes('integrity="') &&
            tag.match(/src="(https:\/\/[^"]+)"/)?.[1];
        if (url) scripts.add(url);
    }
}
report.externalScripts = [...scripts].sort();
// Enforce the network allowlist in Chromium without intercepting worker
// requests: Puppeteer's Fetch interception can leave them paused on reload.
const contentPolicy = [
    "default-src 'self' data: blob:",
    `script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval' ${[...scripts].join(" ")}`,
    "style-src 'self' 'unsafe-inline'",
    "connect-src 'self'",
    "worker-src 'self'",
].join("; ");
const server = http.createServer(async (request, response) => {
    const pathname = decodeURIComponent(
        new URL(request.url, "http://localhost").pathname,
    );
    requests.add(pathname);
    const filename = path.resolve(
        root,
        `.${pathname.endsWith("/") ? `${pathname}index.html` : pathname}`,
    );
    if (!filename.startsWith(`${root}${path.sep}`)) {
        response.writeHead(404).end();
        return;
    }
    try {
        response.setHeader("Content-Security-Policy", contentPolicy);
        const types = {
            ".html": "text/html",
            ".js": "text/javascript",
            ".css": "text/css",
            ".wasm": "application/wasm",
            ".svg": "image/svg+xml",
            ".webp": "image/webp",
            ".png": "image/png",
            ".json": "application/json",
        };
        response.setHeader(
            "Content-Type",
            types[path.extname(filename)] || "application/octet-stream",
        );
        const contents = await fs.readFile(filename);
        if (storageUnavailable && pathname === `${appPath}worker.js`)
            response.end(
                'Object.defineProperty(globalThis, "indexedDB", { value: undefined });\n' +
                    contents.toString("utf8"),
            );
        else response.end(contents);
    } catch {
        response.writeHead(404).end();
    }
});
await new Promise((resolve) => server.listen(0, "127.0.0.1", resolve));
const base = `http://127.0.0.1:${server.address().port}`;
const proxy = process.env.HTTPS_PROXY || process.env.https_proxy;
let browser;
let page;
try {
    browser = await puppeteer.launch({
        executablePath,
        headless: "shell",
        args: proxy ? [`--proxy-server=${proxy}`] : [],
    });
    report.browser = await browser.version();
    page = await browser.newPage();
    page.setDefaultTimeout(120000);
    page.on("pageerror", (error) => report.errors.push(error.message));
    page.on("workercreated", (worker) =>
        report.workerEvents.push({ event: "created", url: worker.url() }),
    );
    page.on("workerdestroyed", (worker) =>
        report.workerEvents.push({ event: "destroyed", url: worker.url() }),
    );
    page.on("requestfailed", (request) =>
        report.failedRequests.push({
            url: request.url(),
            error: request.failure()?.errorText,
        }),
    );
    page.on("request", (request) => {
        const url = request.url();
        if (
            !url.startsWith(`${base}/`) &&
            !url.startsWith("data:") &&
            !url.startsWith("blob:") &&
            !scripts.has(url)
        )
            report.blockedRequests.push(url);
    });
    await page.setViewport({ width: 1280, height: 1000 });
    await page.emulateMediaFeatures([
        { name: "prefers-color-scheme", value: "light" },
    ]);
    const ready = async () => {
        await page.waitForSelector("#anki-app");
        await page.waitForFunction(
            () =>
                !document.querySelector("#check-config").disabled ||
                document
                    .querySelector("#app-status")
                    .classList.contains("text-danger"),
        );
        assert.equal(
            await page.$eval("#check-config", (element) => element.disabled),
            false,
            await page.$eval("#app-status", (element) => element.textContent),
        );
    };
    const action = async (selector) => {
        await page.click(selector);
        await ready();
    };
    const value = (selector) =>
        page.$eval(selector, (element) => element.value);
    const status = () =>
        page.$eval("#app-status", (element) => element.textContent);
    const capture = async (filename) => {
        await page.evaluate(async () => {
            window.scrollTo({ top: 0, left: 0, behavior: "instant" });
            await new Promise((resolve) =>
                requestAnimationFrame(() => requestAnimationFrame(resolve)),
            );
        });
        await page.screenshot({
            path: path.join(output, filename),
            fullPage: true,
        });
    };
    const shell = () =>
        page.evaluate(() => {
            const body = getComputedStyle(document.body);
            const navbar = document.querySelector("header .td-navbar");
            if (!navbar) throw new Error("Shared navigation is missing");
            return {
                background: body.backgroundColor,
                foreground: body.color,
                font: body.fontFamily,
                fontSize: body.fontSize,
                navbarFont: getComputedStyle(navbar).fontFamily,
                navigation: [...navbar.querySelectorAll("a[href]")].map(
                    (link) => ({
                        text: link.textContent.trim(),
                        href: link.getAttribute("href"),
                    }),
                ),
                stylesheets: [
                    ...document.querySelectorAll('link[rel="stylesheet"]'),
                ]
                    .map((link) => link.href)
                    .sort(),
                footer: Boolean(document.querySelector("footer")),
            };
        });
    const controlPalette = () =>
        page.$$eval("#file-content, #workspace-files", (elements) =>
            elements.map((element) => {
                const style = getComputedStyle(element);
                return {
                    id: element.id,
                    background: style.backgroundColor,
                    foreground: style.color,
                };
            }),
        );
    const selectTheme = async (mode) => {
        const option = await page.$(`[data-bs-theme-value="${mode}"]`);
        assert.ok(option, `Shared theme menu offers ${mode}`);
        await option.evaluate((element) =>
            element
                .closest(".dropdown")
                .querySelector('[data-bs-toggle="dropdown"]')
                .click(),
        );
        await option.click();
        await page.waitForFunction(
            (expected) =>
                document.documentElement.dataset.bsTheme === expected,
            {},
            mode,
        );
    };
    await page.goto(`${base}${landingPath}`, { waitUntil: "networkidle0" });
    await page.waitForFunction(
        () => document.documentElement.dataset.bsTheme === "light",
    );
    report.landingLight = await shell();
    assert.ok(report.landingLight.stylesheets.length > 0);
    assert.ok(report.landingLight.footer);
    const links = await page.$$("a");
    const tryLinks = [];
    for (const link of links)
        if (
            await link.evaluate(
                (element) => element.textContent.trim() === "Try it now",
            )
        )
            tryLinks.push(link);
    assert.equal(
        tryLinks.length,
        1,
        "Landing must have one Try it now action",
    );
    await Promise.all([
        page.waitForNavigation({ waitUntil: "networkidle0" }),
        tryLinks[0].click(),
    ]);
    assert.equal(new URL(page.url()).pathname, appPath);
    await ready();
    report.appLight = await shell();
    assert.deepEqual(
        report.appLight,
        report.landingLight,
        "App must reuse the site shell",
    );
    assert.equal(
        await page.$eval("#file-content", (element) =>
            Boolean(element.labels.length),
        ),
        true,
    );
    const lifetime = await page.$eval(
        ".landing-content",
        (element) => element.textContent,
    );
    assert.match(
        lifetime,
        /session|tab|reload/i,
        "Storage lifetime must be visible",
    );
    report.cases.push("landing navigation and shared presentation");

    const files = await page.$$eval("#workspace-files option", (options) =>
        options.map((option) => option.value),
    );
    const notePath = files.find((filename) =>
        filename.endsWith(".note.anki.toml"),
    );
    assert.ok(notePath, "Example must generate an editable note");
    assert.equal(
        await value("#workspace-files"),
        notePath,
        "The example note must open without setup actions",
    );
    const originalNote = await value("#file-content");
    assert.match(originalNote, /你好/, "Generated note must preserve Unicode");
    const editedNote = originalNote.replace("你好", "edited Ж café");
    await page.locator("#file-content").fill(editedNote);
    await page.select(
        "#workspace-files",
        files.find((file) => file !== notePath),
    );
    await ready();
    await page.select("#workspace-files", notePath);
    await ready();
    assert.equal(await value("#file-content"), editedNote);
    await action("#check-config");
    assert.match(await status(), /valid/i);
    report.files = files;
    report.editedNote = notePath;
    await fs.writeFile(path.join(output, "edited-note.anki.toml"), editedNote);
    report.cases.push("automatic example, editing and retention across files");

    const client = await browser.target().createCDPSession();
    await client.send("Browser.setDownloadBehavior", {
        behavior: "allowAndName",
        downloadPath: downloadDirectory,
        eventsEnabled: true,
    });
    const download = async (selector) => {
        let began;
        let timer;
        let rejectDownload;
        let onBegin;
        let onProgress;
        const completed = new Promise((resolve, reject) => {
            rejectDownload = reject;
            onBegin = (event) => {
                began = event;
            };
            onProgress = (event) => {
                if (event.guid !== began?.guid) return;
                if (event.state === "completed") resolve(began);
                if (event.state === "canceled")
                    reject(new Error("Download canceled"));
            };
            client.on("Browser.downloadWillBegin", onBegin);
            client.on("Browser.downloadProgress", onProgress);
            timer = setTimeout(
                () => reject(new Error("Download timed out")),
                15000,
            );
        });
        try {
            await page.click(selector).catch(rejectDownload);
            const details = await completed;
            const bytes = await fs.readFile(
                path.join(downloadDirectory, details.guid),
            );
            const contents = bytes.toString("utf8");
            report.downloads.push({
                filename: details.suggestedFilename,
                bytes: bytes.length,
            });
            return { ...details, contents, bytes };
        } finally {
            clearTimeout(timer);
            client.off("Browser.downloadWillBegin", onBegin);
            client.off("Browser.downloadProgress", onProgress);
        }
    };
    const jsonDownload = await download("#download-config");
    await ready();
    assert.match(jsonDownload.suggestedFilename, /\.json$/);
    const parsedOutput = JSON.parse(jsonDownload.contents);
    assert.equal(parsedOutput.config.version, "1");
    const notes = [
        ...(parsedOutput.notes || []),
        ...(parsedOutput.decks || []).flatMap((deck) => deck.notes || []),
    ];
    assert.ok(
        notes.some((note) =>
            Object.values(note.config.fields).some((field) =>
                field.includes("edited Ж café"),
            ),
        ),
    );
    await fs.writeFile(
        path.join(output, "downloaded-configuration.json"),
        jsonDownload.contents,
    );
    // Download must also include an edit made after the last Check.
    const latestNote = editedNote.replace("edited Ж café", "latest Ж café");
    await page.locator("#file-content").fill(latestNote);
    const latestDownload = await download("#download-config");
    await ready();
    assert.match(latestDownload.contents, /latest Ж café/);
    report.cases.push(
        "complete JSON download includes the current editor text",
    );

    const imported = path.join(scratch, "imported-configuration.json");
    const openJSON = async (contents) => {
        await fs.writeFile(imported, contents);
        const [chooser] = await Promise.all([
            page.waitForFileChooser(),
            page.click("#open-config"),
        ]);
        await chooser.accept([imported]);
        await ready();
    };
    report.invalidImports = [];
    const escapedConfiguration = {
        config: { version: "1" },
        notes: [
            {
                path: "../collection.colpkg",
                config: { fields: { Front: "overwritten" } },
            },
        ],
    };
    const temporaryConfiguration = {
        config: { version: "1" },
        notes: [
            {
                path: "../tmp/nested/hidden.note.anki.toml",
                config: { fields: { Front: "must not be silently omitted" } },
            },
        ],
    };
    for (const invalid of [
        "{",
        JSON.stringify({ config: { version: "1" }, unknown_field: true }),
        JSON.stringify({ config: { version: "999" } }),
        JSON.stringify(escapedConfiguration),
        JSON.stringify(temporaryConfiguration),
    ]) {
        await openJSON(invalid);
        assert.equal(
            await page.$eval("#app-status", (element) =>
                element.classList.contains("text-danger"),
            ),
            true,
        );
        report.invalidImports.push(await status());
        assert.equal(
            await value("#file-content"),
            latestNote,
            "Failed imports must preserve current edits",
        );
        assert.deepEqual(
            await page.$$eval("#workspace-files option", (options) =>
                options.map((option) => option.value),
            ),
            files,
        );
        assert.equal(
            await page.$eval(
                "#download-collection",
                (element) => element.disabled,
            ),
            true,
            "Configuration imports must not create a collection package",
        );
    }
    await openJSON(jsonDownload.contents);
    assert.equal(await value("#workspace-files"), notePath);
    // Re-export makes inherited deck references explicit; compare the complete
    // configuration below rather than requiring identical TOML formatting.
    assert.match(
        await value("#file-content"),
        /edited Ж café/,
        "Opening JSON restores the edited note directly",
    );
    const reopenedDownload = await download("#download-config");
    await ready();
    assert.deepEqual(JSON.parse(reopenedDownload.contents), parsedOutput);
    report.cases.push(
        "failed imports preserve edits and valid JSON reopens directly",
    );

    const recoveryNote = await value("#file-content");
    await page.locator("#file-content").fill("[");
    await action("#check-config");
    report.invalidTOML = await status();
    assert.match(report.invalidTOML, /TOML|toml|parse|decode|syntax|line/i);
    assert.equal(
        await value("#file-content"),
        "[",
        "Malformed TOML must remain editable",
    );
    await page.select(
        "#workspace-files",
        files.find((file) => file !== notePath),
    );
    await ready();
    await action("#check-config");
    assert.equal(
        await page.$eval("#app-status", (element) =>
            element.classList.contains("text-danger"),
        ),
        true,
        "Check must validate other files too",
    );
    await page.select("#workspace-files", notePath);
    await ready();
    assert.equal(await value("#file-content"), "[");
    await action("#download-config");
    assert.equal(
        await page.$eval("#app-status", (element) =>
            element.classList.contains("text-danger"),
        ),
        true,
        "Download must reject malformed configuration",
    );
    await page.locator("#file-content").fill(recoveryNote);
    await action("#check-config");
    assert.match(await status(), /valid/i);
    report.cases.push("malformed TOML recovery");

    const savedNote = recoveryNote.replace("edited Ж café", "saved Ж café");
    await page.locator("#file-content").fill(savedNote);
    await action("#save-workspace");
    assert.match(await status(), /saved/i);
    await page.reload({ waitUntil: "networkidle0" });
    await ready();
    assert.equal(await value("#file-content"), savedNote);
    report.cases.push(
        "explicit save includes editor text and restores on reload",
    );

    const competing = await browser.newPage();
    competing.setDefaultTimeout(120000);
    await competing.goto(`${base}${appPath}`, { waitUntil: "networkidle0" });
    await competing.waitForFunction(
        () => !document.querySelector("#save-workspace").disabled,
    );
    const newerNote = savedNote.replace("saved Ж café", "newer Ж café");
    await competing.locator("#file-content").fill(newerNote);
    await competing.click("#save-workspace");
    await competing.waitForFunction(
        () => !document.querySelector("#save-workspace").disabled,
    );
    assert.match(
        await competing.$eval("#app-status", (element) => element.textContent),
        /saved/i,
    );
    await competing.close();
    const staleNote = savedNote.replace(
        "saved Ж café",
        "unsaved stale Ж café",
    );
    await page.locator("#file-content").fill(staleNote);
    await action("#save-workspace");
    assert.match(await status(), /changed|conflict|another|newer/i);
    assert.equal(await value("#file-content"), staleNote);
    const staleDownload = await download("#download-config");
    await ready();
    assert.match(staleDownload.contents, /unsaved stale Ж café/);
    // Replacing a worker must retain the page's observed revision, rather than
    // adopting a competing tab's newer revision or depending on a storage read.
    await openJSON(staleDownload.contents);
    assert.match(await status(), /opened/i);
    await action("#save-workspace");
    assert.match(await status(), /changed|conflict|another|newer/i);
    assert.match(await value("#file-content"), /unsaved stale Ж café/);
    await action("#restore-workspace");
    assert.equal(await value("#file-content"), newerNote);
    report.cases.push(
        "conflicting saves preserve edits and explicit restore recovers",
    );

    // Seed through the persistence module, then exercise the real UI, worker
    // and Go filesystem. The file count exceeds Chromium's argument limit.
    const beforeLargeRestore = await download("#download-config");
    await ready();
    const largeWorkspace = await page.evaluate(async (appPath) => {
        const { loadWorkspace, saveWorkspace } = await import(
            `${appPath}storage.js`
        );
        const key = "anki-as-code-workspace";
        const snapshot = await loadWorkspace(key);
        const data = new TextEncoder().encode("# Large workspace note\n");
        for (let index = 0; index < 200000; index++)
            snapshot.entries.push({
                path: `/workspace/large/${String(index).padStart(6, "0")}.note.anki.toml`,
                directory: false,
                data,
            });
        const metadata = await saveWorkspace(
            key,
            snapshot.revision,
            snapshot.entries,
        );
        return {
            fileCount: snapshot.entries.filter((entry) => !entry.directory)
                .length,
            lastPath: snapshot.entries.at(-1).path,
            revision: metadata.revision,
        };
    }, appPath);
    const previousWorkers = page.workers();
    assert.equal(previousWorkers.length, 1);
    const previousDisposed = new Promise((resolve, reject) => {
        const destroyed = (worker) => {
            if (worker !== previousWorkers[0]) return;
            clearTimeout(timer);
            page.off("workerdestroyed", destroyed);
            resolve();
        };
        const timer = setTimeout(() => {
            page.off("workerdestroyed", destroyed);
            reject(
                new Error("The previous workspace worker was not disposed"),
            );
        }, page.getDefaultTimeout());
        page.on("workerdestroyed", destroyed);
    });
    await Promise.all([action("#restore-workspace"), previousDisposed]);
    assert.match(await status(), /saved workspace restored/i);
    assert.equal(
        await page.$eval("#workspace-files", (element) => element.length),
        largeWorkspace.fileCount,
    );
    assert.equal(page.workers().length, 1);
    assert.ok(!page.workers().includes(previousWorkers[0]));
    await page.select("#workspace-files", largeWorkspace.lastPath);
    await ready();
    assert.equal(await value("#file-content"), "# Large workspace note\n");
    const largeWorkspaceEdit = "# Edited restored workspace note\n";
    await page.locator("#file-content").fill(largeWorkspaceEdit);
    await action("#save-workspace");
    assert.match(await status(), /workspace saved locally/i);
    const largeSaved = await page.evaluate(
        async (appPath, lastPath) => {
            const { loadWorkspace } = await import(`${appPath}storage.js`);
            const snapshot = await loadWorkspace("anki-as-code-workspace");
            return {
                revision: snapshot.revision,
                contents: new TextDecoder().decode(
                    snapshot.entries.find((entry) => entry.path === lastPath)
                        .data,
                ),
            };
        },
        appPath,
        largeWorkspace.lastPath,
    );
    assert.equal(largeSaved.contents, largeWorkspaceEdit);
    assert.equal(largeSaved.revision, largeWorkspace.revision + 1);
    report.largeWorkspace = { ...largeWorkspace, ...largeSaved };
    report.cases.push(
        "large workspace restore, file selection, editing and worker disposal",
    );
    await openJSON(beforeLargeRestore.contents);
    assert.match(await status(), /opened/i);

    const openCollection = async (filename) => {
        const [chooser] = await Promise.all([
            page.waitForFileChooser(),
            page.click("#open-collection"),
        ]);
        await chooser.accept([filename]);
        await ready();
    };
    const invalidCollection = path.join(scratch, "invalid.colpkg");
    await fs.writeFile(invalidCollection, "invalid collection");
    await openCollection(invalidCollection);
    assert.equal(await value("#file-content"), newerNote);
    assert.equal(
        await page.$eval("#app-status", (element) =>
            element.classList.contains("text-danger"),
        ),
        true,
    );
    await openCollection(fixtureArchive);
    assert.match(await value("#file-content"), /browser fixture original/);
    await action("#plan-collection");
    assert.match(
        await page.$eval(
            "#collection-report",
            (element) => element.textContent,
        ),
        /no .*changes/i,
        "Imported bindings must make the unchanged collection converge",
    );
    const originalCollection = await download("#download-collection");
    await ready();
    const collectionNote = await value("#file-content");
    await page.locator("#file-content").fill("[");
    await action("#apply-collection");
    assert.equal(await value("#file-content"), "[");
    const unchangedCollection = await download("#download-collection");
    await ready();
    assert.deepEqual(unchangedCollection.bytes, originalCollection.bytes);
    const appliedNote = collectionNote.replace(
        "browser fixture original",
        "browser edited 你好",
    );
    await page.locator("#file-content").fill(appliedNote);
    await action("#plan-collection");
    assert.match(
        await page.$eval(
            "#collection-report",
            (element) => element.textContent,
        ),
        /update.*note|note.*update/i,
    );
    const beforeSaveCollection = await download("#download-collection");
    await ready();
    await action("#save-workspace");
    const unappliedCollection = await download("#download-collection");
    await ready();
    assert.deepEqual(
        unappliedCollection.bytes,
        beforeSaveCollection.bytes,
        "Save must not apply configuration changes",
    );
    await action("#apply-collection");
    assert.match(await status(), /applied/i);
    await action("#plan-collection");
    assert.match(
        await page.$eval(
            "#collection-report",
            (element) => element.textContent,
        ),
        /no .*changes/i,
    );
    const appliedCollection = await download("#download-collection");
    await ready();
    assert.match(appliedCollection.suggestedFilename, /\.colpkg$/);
    const collectionOutput = path.join(output, "site-browser-output.colpkg");
    await fs.writeFile(collectionOutput, appliedCollection.bytes);
    const verification = spawnSync(
        fixtureTool,
        [
            "verify",
            collectionOutput,
            fixtureBaseline,
            path.join(output, "site-native-verification.json"),
        ],
        { encoding: "utf8", timeout: 120000 },
    );
    assert.equal(
        verification.status,
        0,
        verification.stderr || verification.error?.message,
    );
    const collectionJSON = await download("#download-config");
    await ready();
    await openJSON(collectionJSON.contents);
    const retainedCollection = await download("#download-collection");
    await ready();
    assert.deepEqual(retainedCollection.bytes, appliedCollection.bytes);
    const retainedEditor = await value("#file-content");
    report.escapedImports = [];
    for (const escaped of [
        escapedConfiguration,
        temporaryConfiguration,
        {
            config: { version: "1" },
            decks: [
                {
                    path: "../collection.colpkg/deck.anki.toml",
                    deck: { normal: {} },
                },
            ],
        },
        {
            config: { version: "1" },
            noteTypes: [
                {
                    path: "../../unexpected.notetype.anki.toml",
                    config: { name: "Outside" },
                },
            ],
        },
    ]) {
        await openJSON(JSON.stringify(escaped));
        assert.match(await status(), /outside.*configuration/i);
        report.escapedImports.push(await status());
        assert.equal(await value("#file-content"), retainedEditor);
        const preserved = await download("#download-collection");
        await ready();
        assert.deepEqual(preserved.bytes, appliedCollection.bytes);
    }
    await openCollection(invalidCollection);
    assert.match(await value("#file-content"), /browser edited 你好/);
    const retainedAfterFailure = await download("#download-collection");
    await ready();
    assert.deepEqual(retainedAfterFailure.bytes, appliedCollection.bytes);
    await action("#save-workspace");
    await page.reload({ waitUntil: "networkidle0" });
    await ready();
    assert.match(await value("#file-content"), /browser edited 你好/);
    const restoredCollection = await download("#download-collection");
    await ready();
    assert.deepEqual(restoredCollection.bytes, appliedCollection.bytes);
    await action("#plan-collection");
    assert.match(
        await page.$eval(
            "#collection-report",
            (element) => element.textContent,
        ),
        /no .*changes/i,
    );
    report.cases.push(
        "collection import, plan, atomic apply, download and native verification",
        "JSON replacement and failed collection imports preserve the package",
        "escaped JSON paths preserve editor and package through save/reload",
        "local save retains unapplied edits and restores completed collection bytes",
    );
    report.controlsLight = await controlPalette();
    await capture("site-desktop-light.png");

    await selectTheme("dark");
    report.appDark = await shell();
    assert.notEqual(report.appDark.background, report.appLight.background);
    report.controlsDark = await controlPalette();
    for (const control of report.controlsDark) {
        const light = report.controlsLight.find(
            (entry) => entry.id === control.id,
        );
        assert.notEqual(
            control.background,
            light.background,
            `${control.id} must follow dark mode`,
        );
        assert.notEqual(
            control.foreground,
            control.background,
            `${control.id} text must contrast with its background`,
        );
    }
    await capture("site-desktop-dark.png");
    await page.goto(`${base}${landingPath}`, { waitUntil: "networkidle0" });
    await page.waitForFunction(
        () => document.documentElement.dataset.bsTheme === "dark",
    );
    report.landingDark = await shell();
    assert.deepEqual(
        report.appDark,
        report.landingDark,
        "Dark mode must use the same site shell",
    );
    await page.goto(`${base}${appPath}`, { waitUntil: "networkidle0" });
    await ready();
    await page.waitForFunction(
        () => document.documentElement.dataset.bsTheme === "dark",
    );
    await page.reload({ waitUntil: "networkidle0" });
    await ready();
    await page.waitForFunction(
        () => document.documentElement.dataset.bsTheme === "dark",
    );
    report.cases.push("theme menu and persisted dark preference");
    await selectTheme("light");
    await page.setViewport({ width: 390, height: 844 });
    report.mobile = await page.evaluate(() => ({
        viewport: document.documentElement.clientWidth,
        content: document.documentElement.scrollWidth,
        controls: [
            ...document.querySelectorAll(
                "#anki-app textarea, #anki-app select, #anki-app input",
            ),
        ].map((element) => {
            const box = element.getBoundingClientRect();
            return { id: element.id, left: box.left, right: box.right };
        }),
    }));
    assert.ok(
        report.mobile.content <= report.mobile.viewport + 1,
        "Mobile page must not overflow horizontally",
    );
    for (const control of report.mobile.controls) {
        assert.ok(
            control.left >= -1,
            `${control.id} begins outside the viewport`,
        );
        assert.ok(
            control.right <= report.mobile.viewport + 1,
            `${control.id} exceeds the viewport`,
        );
    }
    await capture("site-mobile-light.png");
    report.cases.push("mobile controls and page fit");

    // Storage is optional for local file work. Exercise both imports using the
    // real worker with IndexedDB unavailable, then retain portable artifacts.
    storageUnavailable = true;
    await page.reload({ waitUntil: "networkidle0" });
    await ready();
    report.unavailableStorage = { startup: await status() };
    assert.match(report.unavailableStorage.startup, /storage.*unavailable/i);
    await openJSON(jsonDownload.contents);
    assert.match(await status(), /opened/i);
    assert.match(await value("#file-content"), /edited Ж café/);
    await action("#save-workspace");
    report.unavailableStorage.save = await status();
    assert.match(report.unavailableStorage.save, /storage.*unavailable/i);
    const portableJSON = await download("#download-config");
    await ready();
    assert.deepEqual(JSON.parse(portableJSON.contents), parsedOutput);
    await openCollection(fixtureArchive);
    assert.match(await value("#file-content"), /browser fixture original/);
    await action("#plan-collection");
    assert.match(
        await page.$eval(
            "#collection-report",
            (element) => element.textContent,
        ),
        /no .*changes/i,
    );
    const portableCollection = await download("#download-collection");
    await ready();
    assert.equal(portableCollection.bytes.subarray(0, 2).toString(), "PK");
    await fs.writeFile(
        path.join(output, "unavailable-storage-collection.colpkg"),
        portableCollection.bytes,
    );
    report.cases.push("JSON and collection imports work without IndexedDB");
    for (const resource of [
        "app.js",
        "index.js",
        "worker.js",
        "wasm_exec.js",
        "app.wasm",
    ])
        assert.ok(
            requests.has(`${appPath}${resource}`),
            `Site must serve its own ${resource}`,
        );
    assert.deepEqual(report.errors, [], "Site must not emit uncaught errors");
    assert.deepEqual(
        report.blockedRequests,
        [],
        "Site must not require undeclared external resources",
    );
    console.log(
        `Verified ${report.cases.length} rendered-site workflow groups`,
    );
} catch (error) {
    report.failure = error.stack;
    if (page) {
        report.finalStatus = await page
            .$eval("#app-status", (element) => element.textContent)
            .catch(() => null);
        report.finalControls = await page
            .$$eval(
                "#anki-app button, #anki-app select, #anki-app textarea",
                (elements) =>
                    elements.map((element) => ({
                        id: element.id,
                        disabled: element.disabled,
                    })),
            )
            .catch(() => []);
        await page
            .screenshot({
                path: path.join(output, "site-failure.png"),
                fullPage: true,
            })
            .catch(() => {});
    }
    throw error;
} finally {
    report.requests = [...requests].sort();
    await fs.writeFile(
        path.join(output, "site-report.json"),
        JSON.stringify(report, null, 2),
    );
    if (browser) await browser.close();
    await new Promise((resolve, reject) =>
        server.close((error) => (error ? reject(error) : resolve())),
    );
    await fs.rm(scratch, { recursive: true, force: true });
}
