import assert from "node:assert/strict";
import fs from "node:fs/promises";
import path from "node:path";
import { createHash } from "node:crypto";
import { spawnSync } from "node:child_process";
import puppeteer, { type Browser } from "puppeteer";

const [
    chrome,
    nodeProgram,
    browserModule,
    invalidSource,
    compiler,
    declarations,
] = process.argv.slice(2);
if (
    !chrome ||
    !nodeProgram ||
    !browserModule ||
    !invalidSource ||
    !compiler ||
    !declarations
)
    throw new Error(
        "Expected Chrome, Node/browser programs, invalid source, compiler, and declarations",
    );
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR ?? process.cwd();
await fs.mkdir(output, { recursive: true });
const report: Record<string, unknown> = { status: "failed" };
const hash = (value: string | Uint8Array): string =>
    createHash("sha256").update(value).digest("hex");
let browser: Browser | undefined;
try {
    const destination = path.join(output, "node-result.json");
    const node = spawnSync(nodeProgram, [destination], {
        encoding: "utf8",
        timeout: 30000,
    });
    assert.equal(node.status, 0, node.stdout + node.stderr);
    const nodeResult = JSON.parse(await fs.readFile(destination, "utf8"));
    assert.equal(nodeResult.magnitude, 13);
    assert.equal(nodeResult.moduleType, "module");
    report.node = nodeResult;

    const invalid = spawnSync(
        compiler,
        [
            "--ignoreConfig",
            "--pretty",
            "false",
            "--strict",
            "--noEmit",
            "--target",
            "ES2024",
            "--module",
            "NodeNext",
            invalidSource,
        ],
        { encoding: "utf8", timeout: 30000 },
    );
    const diagnostics = invalid.stdout + invalid.stderr;
    await fs.writeFile(path.join(output, "type-error.log"), diagnostics);
    assert.notEqual(
        invalid.status,
        0,
        "The compiler must reject a type error",
    );
    assert.deepEqual(
        [...diagnostics.matchAll(/error TS(\d+):/g)].map((match) => match[1]),
        ["2322"],
        diagnostics,
    );
    report.compiler = { status: invalid.status, diagnostics };

    const requests: string[] = [];
    browser = await puppeteer.launch({
        executablePath: chrome,
        headless: true,
        args: ["--no-sandbox", "--disable-dev-shm-usage"],
    });
    const page = await browser.newPage();
    await page.setRequestInterception(true);
    page.on("request", (request) => {
        if (/^(data:|about:)/.test(request.url())) void request.continue();
        else {
            requests.push(request.url());
            void request.abort();
        }
    });
    const program = await fs.readFile(browserModule, "utf8");
    const result = await page.evaluate(
        async (url) => {
            const module = await import(url);
            const canvas: HTMLCanvasElement = module.drawCoordinates({
                x: 3,
                y: 4,
                z: 12,
            });
            document.body.append(canvas);
            const context = canvas.getContext("2d");
            if (!context) throw new Error("Canvas context unavailable");
            const pixels = context.getImageData(
                0,
                0,
                canvas.width,
                canvas.height,
            ).data;
            let ink = 0;
            for (let i = 0; i < pixels.length; i += 4)
                if (pixels[i] < 128 && pixels[i + 3] === 255) ink += 1;
            return {
                width: canvas.width,
                height: canvas.height,
                ink,
                png: canvas.toDataURL("image/png"),
            };
        },
        `data:text/javascript;base64,${Buffer.from(program).toString("base64")}`,
    );
    assert.ok(result.width > 0 && result.height > 0);
    assert.ok(
        result.ink > 20,
        "The emitted browser module must draw visible content",
    );
    assert.deepEqual(requests, [], "The emitted module must work offline");
    const png = Buffer.from(result.png.split(",")[1], "base64");
    await fs.writeFile(path.join(output, "browser.png"), png);
    report.browser = {
        width: result.width,
        height: result.height,
        ink: result.ink,
        sha256: hash(png),
    };
    report.outputs = {
        browserModule: hash(program),
        declarations: hash(await fs.readFile(declarations)),
    };
    report.status = "passed";
} catch (error) {
    report.error = error instanceof Error ? error.message : String(error);
    throw error;
} finally {
    await browser?.close();
    await fs.writeFile(
        path.join(output, "verification.json"),
        JSON.stringify(report, null, 2) + "\n",
    );
}
