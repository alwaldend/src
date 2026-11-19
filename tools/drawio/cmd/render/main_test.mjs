import assert from "node:assert/strict";
import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";

const [webapp, browser, fontAnchor] = process.argv
    .slice(2)
    .map((p) => path.resolve(p));
const scratch = await fs.mkdtemp(
    path.join(process.env.TEST_TMPDIR, "drawio-"),
);
const input = path.join(scratch, "diagram.drawio");
const script = fileURLToPath(new URL("main.mjs", import.meta.url));
const xml = `<mxfile><diagram name="Test"><mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/><mxCell id="box" value="Service" vertex="1" parent="1"><mxGeometry x="20" y="20" width="160" height="80" as="geometry"/></mxCell></root></mxGraphModel></diagram></mxfile>`;
await fs.writeFile(input, xml);
async function run(name, pageName = "Test") {
    const output = path.join(scratch, name + ".svg");
    const manifest = path.join(scratch, name + ".json");
    await fs.writeFile(manifest, JSON.stringify({ [pageName]: output }));
    const result = spawnSync(
        process.execPath,
        [script, webapp, browser, input, manifest, fontAnchor],
        { encoding: "utf8", timeout: 45000 },
    );
    return { ...result, output };
}
try {
    const first = await run("first");
    assert.equal(first.status, 0, first.stderr);
    const second = await run("second");
    assert.equal(second.status, 0, second.stderr);
    const svg = await fs.readFile(first.output, "utf8");
    assert.equal(svg, await fs.readFile(second.output, "utf8"));
    assert.match(svg, /Service/);
    assert.match(
        svg,
        /<svg\b[^>]*><rect width="100%" height="100%" fill="#ffffff"\/>/,
    );
    assert.match(svg, /Source SHA-256: [0-9a-f]{64}/);
    assert.match(svg, /viewBox="0 0 [1-9][0-9.]* [1-9][0-9.]*"/);
    const missing = await run("missing", "Absent");
    assert.notEqual(missing.status, 0);
    assert.match(missing.stderr, /Expected exactly one page named Absent/);
    await assert.rejects(fs.access(missing.output));
    await fs.writeFile(input, "<invalid");
    const malformed = await run("malformed");
    assert.notEqual(malformed.status, 0);
    assert.match(malformed.stderr, /Invalid Drawio XML/);
    await assert.rejects(fs.access(malformed.output));
} finally {
    await fs.rm(scratch, { recursive: true, force: true });
}
