import assert from "node:assert/strict";
import fs from "node:fs/promises";
import path from "node:path";
import { promisify } from "node:util";
import { execFile } from "node:child_process";
const run = promisify(execFile);
const cli = path.resolve(process.argv[2]);
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
const input = path.join(output, "literal.html");
await fs.writeFile(
    input,
    '<script>throw new Error("must not execute")</script>\n' +
        "long line ".repeat(160) +
        "\n" +
        "    indented & <literal>\n".repeat(60),
);
for (const theme of ["light", "dark"]) {
    const image = path.join(output, `${theme}.png`);
    await run(cli, [
        "--input",
        input,
        "--output",
        image,
        "--theme",
        theme,
        "--scale",
        "1",
    ]);
    const png = await fs.readFile(image);
    assert.equal(png.subarray(1, 4).toString(), "PNG");
    assert.ok(
        png.readUInt32BE(16) === 960,
        "long lines must wrap within the canvas",
    );
    assert.ok(png.readUInt32BE(20) > 1000, "all lines must be captured");
}
for (const args of [
    ["--scale", "0"],
    ["--theme", "invalid"],
    ["--unknown", "value"],
]) {
    await assert.rejects(
        run(cli, [
            "--input",
            input,
            "--output",
            path.join(output, "invalid.png"),
            ...args,
        ]),
    );
}
await assert.rejects(
    run(cli, [
        "--input",
        input + ".missing",
        "--output",
        path.join(output, "missing.png"),
    ]),
);
await fs.writeFile(
    path.join(output, "report.json"),
    JSON.stringify(
        {
            themes: ["light", "dark"],
            literalSource: true,
            invalidOptionsRejected: true,
        },
        null,
        2,
    ),
);

// Exercise multiline tokens, HTML-like strings, unknown extensions, and overrides through the command.
const javascript = path.join(output, "highlight.js");
await fs.writeFile(
    javascript,
    `/* This comment spans
   two numbered lines. */
export function describe(name) {
    const markup = '<img src="invalid" onerror="throw new Error(1)">';
    return "Hello, " + name + ": " + markup;
}
`,
);
for (const theme of ["light", "dark"]) {
    await run(cli, [
        "--input",
        javascript,
        "--output",
        path.join(output, `highlight-${theme}.png`),
        "--theme",
        theme,
        "--width",
        "640",
    ]);
}
const unknown = path.join(output, "unknown.extension");
await fs.copyFile(javascript, unknown);
await run(cli, [
    "--input",
    unknown,
    "--output",
    path.join(output, "plain.png"),
]);
await run(cli, [
    "--input",
    unknown,
    "--output",
    path.join(output, "override.png"),
    "--language",
    "javascript",
]);
await assert.rejects(
    run(cli, [
        "--input",
        unknown,
        "--output",
        path.join(output, "invalid-language.png"),
        "--language",
        "unsupported-language",
    ]),
);

assert.notDeepEqual(
    await fs.readFile(path.join(output, "plain.png")),
    await fs.readFile(path.join(output, "override.png")),
    "language override must visibly highlight identical source and filename",
);
