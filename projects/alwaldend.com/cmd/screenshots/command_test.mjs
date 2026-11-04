import assert from "node:assert/strict";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { spawn } from "node:child_process";

const command = path.resolve(process.argv[2]);
const chrome = path.resolve(process.argv[3]);
const root = await fs.mkdtemp(
    path.join(process.env.TEST_TMPDIR || os.tmpdir(), "site-screenshots-"),
);
const site = path.join(root, "site");
await fs.mkdir(site);
await fs.writeFile(
    path.join(site, "style.css"),
    "body{min-height:1700px;background:#fff;color:#111}h1{font-size:40px}",
);
for (const [page, title] of [
    ["", "Home"],
    ["projects", "Project landings"],
    ["docs/projects", "Project documentation"],
]) {
    const directory = path.join(site, page);
    await fs.mkdir(directory, { recursive: true });
    await fs.writeFile(
        path.join(directory, "index.html"),
        `<!doctype html><title>${title}</title><link rel="stylesheet" href="/style.css"><h1>${title}</h1><img src="https://example.invalid/blocked.png"><script>document.title += " " + innerWidth + "x" + innerHeight + " " + (matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light")</script>`,
    );
}
const invoke = async (args) => {
    const child = spawn(process.execPath, [command, site, chrome, ...args], {
        env: process.env,
    });
    let stdout = "";
    let stderr = "";
    child.stdout.on("data", (chunk) => {
        stdout += chunk;
    });
    child.stderr.on("data", (chunk) => {
        stderr += chunk;
    });
    const code = await new Promise((resolve, reject) => {
        child.on("error", reject);
        child.on("close", resolve);
    });
    return { code, stdout, stderr };
};
const verify = async (directory, pages, width, height, scheme) => {
    const manifest = JSON.parse(
        await fs.readFile(path.join(directory, "screenshots.json"), "utf8"),
    );
    assert.deepEqual(
        manifest.pages.map((page) => page.path),
        pages,
    );
    assert.equal(manifest.colorScheme, scheme);
    assert.deepEqual(manifest.viewport, { width, height });
    for (const page of manifest.pages) {
        assert.ok(
            page.title.endsWith(`${width}x${height} ${scheme}`),
            "the browser must observe the selected viewport and color scheme",
        );
        const png = await fs.readFile(path.join(directory, page.file));
        assert.equal(png.subarray(1, 4).toString(), "PNG");
        assert.equal(png.readUInt32BE(16), width);
        assert.ok(
            png.readUInt32BE(20) >= 1700,
            "local CSS and full-page capture must both work",
        );
        assert.equal(page.image.width, width);
        assert.equal(page.image.height, png.readUInt32BE(20));
        assert.ok(
            page.blockedExternalRequests > 0,
            "external image requests must be intercepted",
        );
    }
    assert.equal(
        (await fs.readdir(directory)).some((name) =>
            name.startsWith(".browser-"),
        ),
        false,
    );
    return manifest;
};
const report = [];
try {
    const standard = path.join(root, "standard");
    let result = await invoke(["--output-dir", standard]);
    assert.equal(result.code, 0, result.stderr);
    report.push({
        name: "default-indexes",
        manifest: await verify(
            standard,
            ["/", "/projects/", "/docs/projects/"],
            1440,
            1000,
            "light",
        ),
    });
    const custom = path.join(root, "custom");
    result = await invoke([
        "--output-dir",
        custom,
        "--path",
        "/projects/",
        "--width",
        "800",
        "--height",
        "600",
        "--color-scheme",
        "dark",
    ]);
    assert.equal(result.code, 0, result.stderr);
    report.push({
        name: "custom-configuration",
        manifest: await verify(custom, ["/projects/"], 800, 600, "dark"),
    });
    for (const [name, args] of [
        ["missing-page", ["--path", "/missing/"]],
        ["external-page", ["--path", "https://example.invalid/"]],
        ["invalid-width", ["--width", "0"]],
        ["invalid-scheme", ["--color-scheme", "sepia"]],
        ["duplicate-page", ["--path", "/projects/", "--path", "/projects/"]],
    ]) {
        const directory = path.join(root, name);
        result = await invoke(["--output-dir", directory, ...args]);
        assert.notEqual(result.code, 0, `${name} must fail`);
        const entries = await fs.readdir(directory).catch(() => []);
        assert.equal(
            entries.some((entry) => entry.startsWith(".browser-")),
            false,
            `${name} left a browser profile`,
        );
        report.push({ name, result: "rejected" });
    }
    const artifacts = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
    if (artifacts) {
        await fs.mkdir(artifacts, { recursive: true });
        await fs.cp(standard, path.join(artifacts, "standard"), {
            recursive: true,
        });
        await fs.cp(custom, path.join(artifacts, "custom"), {
            recursive: true,
        });
        await fs.writeFile(
            path.join(artifacts, "screenshot-results.json"),
            JSON.stringify(report, null, 2) + "\n",
        );
    }
} finally {
    await fs.rm(root, { recursive: true, force: true });
}
