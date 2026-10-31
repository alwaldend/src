import assert from "node:assert/strict";
import { spawn } from "node:child_process";
import { once } from "node:events";
import fs from "node:fs/promises";
import http from "node:http";
import { createRequire } from "node:module";
import path from "node:path";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const [cleaner, chrome] = process.argv
    .slice(2)
    .map((file) => path.resolve(file));
const artifacts = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(artifacts, "Bazel must provide an artifact directory");
const scratch = await fs.mkdtemp(path.join(process.env.TEST_TMPDIR, "serve-"));
const initialPolicy =
    "delete_post_when(function(post) return post.age_days > 30 end)\n";
const pickedPolicy =
    initialPolicy +
    "delete_repost_when(function(post) return post.age_days > 30 end)\n";
const initialPath = path.join(scratch, "initial.lua");
const pickedPath = path.join(scratch, "picked.lua");
await fs.writeFile(initialPath, initialPolicy);
await fs.writeFile(pickedPath, pickedPolicy);

const records = [
    {
        id: "101",
        text:
            "A draft thought to remove. ".repeat(45) +
            "Complete text ends here.",
    },
    { id: "102", text: "Keep this discussion for later" },
    {
        id: "103",
        text: '<img src=x onerror="window.profileScriptExecuted=true">',
    },
    {
        id: "104",
        text: "A repost to undo",
        referenced_posts: [{ type: "retweeted", id: "901" }],
    },
    { id: "105", text: "Foreign activity", author_id: "200" },
].map((record) => ({
    author_id: "100",
    created_at: new Date(Date.now() - 60 * 86400000).toISOString(),
    ...record,
}));
const mutations = [];
const mutationAttempts = [];
const missingTargets = new Set();
let rejectMutation = "";
let loseResponse = "";
let holdMutation = "";
let mutationHeld = () => {};
let reads = 0;
let holdTimeline = false;
let rejectTimeline = false;
let exhaustedDelete = false;
const fixture = http.createServer((request, response) => {
    response.setHeader("Content-Type", "application/json");
    const url = new URL(request.url, "http://fixture");
    if (request.method === "GET" && url.pathname === "/2/users/me") {
        response.end(JSON.stringify({ data: { id: "100" } }));
    } else if (
        request.method === "GET" &&
        url.pathname === "/2/users/100/tweets"
    ) {
        reads++;
        if (holdTimeline) return;
        if (rejectTimeline) {
            response.writeHead(429);
            response.end("{}");
            return;
        }
        response.end(JSON.stringify({ data: records }));
    } else if (request.method === "DELETE") {
        mutationAttempts.push(url.pathname);
        if (missingTargets.has(url.pathname)) {
            response.writeHead(404);
            response.end("{}");
            return;
        }
        if (url.pathname === rejectMutation) {
            rejectMutation = "";
            response.writeHead(503);
            response.end("{}");
            return;
        }
        if (url.pathname === holdMutation) {
            mutationHeld();
            return;
        }
        mutations.push(url.pathname);
        if (url.pathname === loseResponse) {
            loseResponse = "";
            missingTargets.add(url.pathname);
            response.destroy();
            return;
        }
        if (url.pathname.startsWith("/2/tweets/") && !exhaustedDelete) {
            exhaustedDelete = true;
            response.setHeader("X-Rate-Limit-Remaining", "0");
            response.setHeader(
                "X-Rate-Limit-Reset",
                String(Math.floor(Date.now() / 1000) + 1),
            );
        }
        response.end(
            JSON.stringify({
                data: url.pathname.includes("/retweets/")
                    ? { retweeted: false }
                    : { deleted: true },
            }),
        );
    } else {
        response.writeHead(404);
        response.end("{}");
    }
});
fixture.listen(0, "127.0.0.1");
await once(fixture, "listening");
const environment = Object.fromEntries(
    Object.entries(process.env).filter(([name]) => !name.startsWith("X_")),
);
Object.assign(environment, {
    X_API_KEY: "browser-fixture-key",
    X_API_SECRET: "browser-fixture-secret",
    X_ACCESS_TOKEN: "browser-fixture-token",
    X_ACCESS_TOKEN_SECRET: "browser-fixture-token-secret",
});
const child = spawn(
    cleaner,
    [
        "serve",
        "--listen",
        "127.0.0.1:0",
        "--config",
        initialPath,
        "--api-base-url",
        `http://127.0.0.1:${fixture.address().port}/2`,
        "--policy-timeout",
        "5s",
    ],
    { env: environment, stdio: ["ignore", "pipe", "pipe"] },
);
let diagnostics = "";
child.stderr.on("data", (chunk) => {
    diagnostics += chunk;
});
const exited = once(child, "exit");
let browser;
try {
    const accessLink = await new Promise((resolve, reject) => {
        let output = "";
        const timer = setTimeout(
            () => reject(new Error("serve did not publish its access link")),
            5000,
        );
        child.stdout.on("data", (chunk) => {
            output += chunk;
            const match = output.match(
                /http:\/\/127\.0\.0\.1:\d+\/#[a-f0-9]{64}/,
            );
            if (match) {
                clearTimeout(timer);
                resolve(match[0]);
            }
        });
        child.once("exit", () => {
            clearTimeout(timer);
            reject(new Error("serve exited before publishing an access link"));
        });
    });
    const address = new URL(accessLink);
    const capability = address.hash.slice(1);
    const api = async (route, body, overrides = {}) => {
        const response = await fetch(address.origin + route, {
            method: body === undefined ? "GET" : "POST",
            headers: {
                "X-Cleaner-Token": capability,
                Origin: address.origin,
                ...(body === undefined
                    ? {}
                    : { "Content-Type": "application/json" }),
                ...overrides,
            },
            body: body === undefined ? undefined : JSON.stringify(body),
        });
        return { status: response.status, body: await response.json() };
    };
    assert.equal((await fetch(address.origin + "/api/state")).status, 403);
    assert.equal(
        (
            await api("/api/state", undefined, {
                Origin: "https://foreign.example",
            })
        ).status,
        403,
    );
    const foreignHostStatus = await new Promise((resolve, reject) => {
        const request = http.get(
            address.origin + "/api/state",
            {
                headers: {
                    Host: "foreign.example",
                    "X-Cleaner-Token": capability,
                },
            },
            (response) => {
                response.resume();
                resolve(response.statusCode);
            },
        );
        request.on("error", reject);
    });
    assert.equal(foreignHostStatus, 403);
    assert.equal(
        (await api("/api/plan", { policy: pickedPolicy, path: "/etc/passwd" }))
            .status,
        400,
    );
    assert.equal((await api("/api/plan", {})).status, 400);
    assert.equal((await api("/api/plan", null)).status, 400);
    assert.equal(mutations.length, 0);

    browser = await puppeteer.launch({
        executablePath: chrome,
        headless: "shell",
        userDataDir: path.join(scratch, "browser"),
    });
    const page = await browser.newPage();
    const failures = [];
    page.on("pageerror", (error) => failures.push(error.message));
    await page.setViewport({ width: 1360, height: 1000 });
    await page.goto(accessLink);
    await page.waitForFunction(
        (source) => document.querySelector("#policy")?.value === source,
        {},
        initialPolicy,
    );
    assert.equal(
        new URL(page.url()).hash,
        "",
        "the access link must leave the address bar",
    );
    const picker = await page.$("#policy-file");
    await picker.uploadFile(pickedPath);
    await page.waitForFunction(
        (source) => document.querySelector("#policy").value === source,
        {},
        pickedPolicy,
    );
    await page.click("#plan");
    await page.waitForSelector('#status[data-phase="ready"]');
    await page.waitForFunction(
        () => document.querySelectorAll("#items tbody tr").length === 4,
    );
    assert.equal(mutations.length, 0, "planning must never mutate X");
    assert.equal(
        await page.evaluate(() => window.profileScriptExecuted),
        undefined,
    );
    const stored = await api("/api/plan");
    assert.equal(stored.status, 200);
    const planID = stored.body.plan_id;
    assert.equal(
        (await api("/api/apply", { plan_id: planID, selected: [999] })).status,
        400,
    );
    assert.equal(
        (await api("/api/apply", { plan_id: planID, selected: [0, 0] }))
            .status,
        400,
    );
    assert.equal(
        (await api("/api/apply", { plan_id: planID, selected: [null] }))
            .status,
        400,
    );
    assert.equal(
        (
            await api("/api/apply", {
                plan_id: planID,
                selected: [],
                actions: [{ kind: "delete_post", target_id: "105" }],
            })
        ).status,
        400,
    );
    assert.equal(
        (await api("/api/apply", { plan_id: "stale", selected: [0] })).status,
        409,
    );
    assert.equal(
        (await api("/api/apply", { plan_id: planID, selected: [] })).status,
        400,
    );
    assert.equal(mutations.length, 0);
    await page.$eval("#policy", (input) => {
        input.value += "-- changed\n";
        input.dispatchEvent(new Event("input"));
    });
    assert.ok(
        await page.$eval("#apply", (button) => button.disabled),
        "editing the policy requires a new plan",
    );
    await page.$eval(
        "#policy",
        (input, value) => {
            input.value = value;
            input.dispatchEvent(new Event("input"));
        },
        pickedPolicy,
    );
    await page.click("#exclude-all");
    assert.equal(
        await page.$$eval("#items input:checked", (inputs) => inputs.length),
        0,
    );
    await page.click("#select-all");
    assert.equal(
        await page.$$eval("#items input:checked", (inputs) => inputs.length),
        4,
    );
    await page.type("#filter", "102");
    await page.click("#exclude-all");
    assert.equal(
        await page.$$eval("#items input:checked", (inputs) => inputs.length),
        3,
        "filtered bulk exclusion must preserve hidden selections",
    );
    assert.match(
        await page.$eval("#selection-note", (el) => el.textContent),
        /3 selected items hidden/,
    );
    assert.equal(
        await page.$eval(
            'tr[data-target="102"] .item-action',
            (el) => el.textContent,
        ),
        "Excluded",
    );
    await page.click("#select-all");
    assert.equal(
        await page.$$eval("#items input:checked", (inputs) => inputs.length),
        4,
    );
    await page.$eval("#filter", (input) => {
        input.value = "no matching item";
        input.dispatchEvent(new Event("input"));
    });
    assert.ok(await page.$eval("#select-all", (button) => button.disabled));
    assert.match(
        await page.$eval("#filter-summary", (el) => el.textContent),
        /0 of 4/,
    );
    await page.$eval("#filter", (input) => {
        input.value = "";
        input.dispatchEvent(new Event("input"));
    });
    await page.click('tr[data-target="102"] input[type="checkbox"]');
    assert.match(
        await page.$eval("#selected-count", (element) => element.textContent),
        /3/,
    );
    assert.ok(await page.$eval("#operation-status", (el) => el.hidden));
    const previewText = await page.$eval(
        'tr[data-target="101"] .item-text',
        (el) => el.textContent,
    );
    assert.ok(
        previewText.length < records[0].text.length &&
            previewText.endsWith("…"),
    );
    await page.click('tr[data-target="101"] .expand-post');
    assert.equal(
        await page.$eval(
            'tr[data-target="101"] .item-text',
            (el) => el.textContent,
        ),
        records[0].text,
    );
    await page.click('tr[data-target="101"] .expand-post');
    assert.equal(
        await page.$eval(
            'tr[data-target="101"] .item-text',
            (el) => el.textContent,
        ),
        previewText,
    );
    await page.click('tr[data-target="101"] .expand-post');
    await page.setViewport({ width: 1360, height: 720 });
    const scrolling = await page.evaluate(() => {
        const list = document.querySelector(".items-container");
        const apply = document.querySelector("#apply");
        const before = apply.getBoundingClientRect().top;
        list.scrollTop = list.scrollHeight;
        const last = document.querySelector("#items tbody tr:last-child");
        return {
            moved: list.scrollTop > 0,
            lastVisible:
                last.getBoundingClientRect().bottom <=
                list.getBoundingClientRect().bottom + 1,
            controlsStayVisible:
                apply.getBoundingClientRect().top === before &&
                before >= 0 &&
                apply.getBoundingClientRect().bottom <= innerHeight,
            pageFits: document.documentElement.scrollHeight <= innerHeight,
        };
    });
    assert.ok(scrolling.moved, "the removal list must scroll independently");
    assert.ok(scrolling.lastVisible, "scrolling must reveal the last item");
    assert.ok(
        scrolling.controlsStayVisible,
        "review controls must stay visible",
    );
    assert.ok(scrolling.pageFits, "desktop review must fit the viewport");
    await page.screenshot({
        path: path.join(artifacts, "serve-review-scrolled.png"),
        fullPage: true,
    });
    await page.$eval(".items-container", (list) => {
        list.scrollTop = 0;
    });
    await page.setViewport({ width: 1360, height: 1000 });
    await page.click('tr[data-target="101"] .expand-post');
    const themes = {};
    for (const mode of ["light", "dark"]) {
        await page.select("#theme", mode);
        await page.evaluate(async () => {
            await Promise.allSettled(
                document
                    .getAnimations()
                    .map((animation) => animation.finished),
            );
        });
        themes[mode] = await page.evaluate(
            () => getComputedStyle(document.body).backgroundColor,
        );
        const checkboxColors = await page.$eval(
            "#items input:checked",
            (el) => ({
                fill: getComputedStyle(el).backgroundColor,
                canvas: getComputedStyle(document.body).backgroundColor,
            }),
        );
        assert.notEqual(
            checkboxColors.fill,
            checkboxColors.canvas,
            "checked boxes must remain visible",
        );
        await page.screenshot({
            path: path.join(artifacts, `serve-review-${mode}.png`),
            fullPage: true,
        });
    }
    assert.notEqual(themes.light, themes.dark);
    await page.setViewport({ width: 390, height: 844 });
    await page.screenshot({
        path: path.join(artifacts, "serve-review-mobile.png"),
        fullPage: true,
    });
    assert.ok(
        await page.evaluate(
            () => document.documentElement.scrollWidth <= window.innerWidth,
        ),
        "page must fit a narrow viewport",
    );
    await page.setViewport({ width: 1360, height: 1000 });
    await page.click("#apply");
    await page.waitForSelector("#confirmation[open]");
    await page.click("#confirm-apply");
    await page.waitForFunction(() =>
        document
            .querySelector("#logs")
            .textContent.includes("Waiting for API quota"),
    );
    await page.waitForSelector('#status[data-phase="completed"]');
    assert.deepEqual(
        [...mutations].sort(),
        ["/2/tweets/101", "/2/tweets/103", "/2/users/100/retweets/901"].sort(),
    );
    const completed = await api("/api/state");
    assert.equal(completed.body.completed, 3);
    assert.equal(
        (await api("/api/apply", { plan_id: planID, selected: [0] })).status,
        409,
    );
    assert.match(
        await page.$eval("#logs", (element) => element.textContent),
        /Cleanup completed/,
    );
    await page.reload();
    await page.waitForSelector('#status[data-phase="completed"]');
    assert.equal(mutations.length, 3, "reload must not replay apply");

    holdTimeline = true;
    await page.click("#plan");
    await page.waitForSelector('#status[data-phase="planning"]');
    assert.equal(
        (await api("/api/plan", { policy: pickedPolicy })).status,
        409,
    );
    assert.equal(
        (await api("/api/apply", { plan_id: planID, selected: [0] })).status,
        409,
    );
    await page.click("#cancel-operation");
    await page.waitForSelector('#status[data-phase="cancelled"]');
    holdTimeline = false;
    assert.equal(mutations.length, 3);
    assert.equal((await api("/api/plan")).status, 409);

    rejectTimeline = true;
    await page.click("#plan");
    await page.waitForSelector('#status[data-phase="failed"]');
    assert.match(
        await page.$eval("#error", (element) => element.textContent),
        /429/,
    );
    assert.equal((await api("/api/plan")).status, 409);
    assert.equal(mutations.length, 3);
    rejectTimeline = false;
    const readsBeforeEmptyPolicy = reads;
    assert.equal(
        (
            await api("/api/plan", {
                policy: "-- a different tab's empty policy\n",
            })
        ).status,
        202,
    );
    await page.waitForSelector('#status[data-phase="ready"]');
    await page.waitForFunction(
        () =>
            document.querySelector("#selected-count").textContent ===
            "0 of 0 selected",
    );
    assert.equal(
        reads,
        readsBeforeEmptyPolicy,
        "empty policy must not fetch timeline data",
    );
    assert.match(
        await page.$eval("#review-note", (element) => element.textContent),
        /Policy changed/,
    );
    assert.ok(await page.$eval("#apply", (button) => button.disabled));
    await page.reload();
    await page.waitForSelector('#status[data-phase="ready"]');
    await page.waitForFunction(() =>
        document.querySelector("#policy").value.includes("different tab"),
    );
    // Resume the exact selection after errors, a lost success response, and cancellation.
    for (const interrupted of ["failure", "cancellation"]) {
        missingTargets.clear();
        await page.$eval(
            "#policy",
            (input, value) => {
                input.value = value;
                input.dispatchEvent(new Event("input"));
            },
            pickedPolicy,
        );
        await page.click("#plan");
        await page.waitForSelector('#status[data-phase="ready"]');
        await page.waitForFunction(
            () => document.querySelectorAll("#items tbody tr").length === 4,
        );
        await page.click('tr[data-target="102"] input[type="checkbox"]');
        const attemptsBefore = mutationAttempts.length;
        const mutationsBefore = mutations.length;
        const readsBefore = reads;
        const held = new Promise((resolve) => {
            mutationHeld = resolve;
        });
        exhaustedDelete = false;
        if (interrupted === "failure") rejectMutation = "/2/tweets/103";
        else holdMutation = "/2/tweets/103";
        await page.click("#apply");
        await page.waitForSelector("#confirmation[open]");
        await page.click("#confirm-apply");
        if (interrupted === "cancellation") {
            await held;
            await page.waitForFunction(
                () =>
                    document.querySelector("#selected-count").textContent ===
                    "2 of 3 completed",
            );
            await page.click("#cancel-operation");
        }
        const phase = interrupted === "failure" ? "failed" : "cancelled";
        await page.waitForSelector(`#status[data-phase="${phase}"]`);
        await page.reload();
        await page.waitForFunction(
            () =>
                document.querySelector("#apply").textContent === "Resume" &&
                !document.querySelector("#apply").disabled,
        );
        const stopped = await api("/api/state");
        assert.equal(stopped.body.completed, 2);
        assert.equal(stopped.body.remaining, 1);
        assert.equal(
            (await api("/api/resume", { job_id: "stale" })).status,
            409,
        );
        assert.equal(
            (
                await api("/api/resume", {
                    job_id: stopped.body.job_id,
                    selected: [0],
                })
            ).status,
            400,
        );
        assert.equal(
            (
                await api("/api/apply", {
                    plan_id: stopped.body.plan_id,
                    selected: [0],
                })
            ).status,
            409,
        );
        assert.equal(
            reads,
            readsBefore,
            "reload must not replan an interrupted run",
        );
        await page.screenshot({
            path: path.join(artifacts, `serve-resume-${phase}.png`),
            fullPage: true,
        });
        if (interrupted === "failure") {
            loseResponse = "/2/tweets/103";
            await page.click("#apply");
            await page.waitForFunction(
                () =>
                    document.querySelector("#status").dataset.phase ===
                        "failed" &&
                    document
                        .querySelector("#error")
                        .textContent.includes("103"),
            );
            await page.waitForFunction(
                () => !document.querySelector("#apply").disabled,
            );
        } else holdMutation = "";
        await page.click("#apply");
        await page.waitForSelector('#status[data-phase="completed"]');
        const resumed = await api("/api/state");
        assert.equal(resumed.body.completed, 3);
        assert.equal(resumed.body.remaining, 0);
        assert.equal(
            reads,
            readsBefore,
            "resume must not fetch or evaluate a new policy",
        );
        assert.equal(mutations.length - mutationsBefore, 3);
        const attempts = mutationAttempts.slice(attemptsBefore);
        assert.equal(attempts.filter((p) => p === "/2/tweets/101").length, 1);
        assert.equal(
            attempts.filter((p) => p === "/2/users/100/retweets/901").length,
            1,
        );
        assert.equal(
            attempts.filter((p) => p === "/2/tweets/103").length,
            interrupted === "failure" ? 3 : 2,
        );
        assert.equal(attempts.filter((p) => p === "/2/tweets/102").length, 0);
        assert.equal(
            (await api("/api/resume", { job_id: resumed.body.job_id })).status,
            409,
        );
    }
    assert.deepEqual(failures, []);
    for (const value of [
        capability,
        ...Object.entries(environment)
            .filter(([key]) => key.startsWith("X_"))
            .map(([, value]) => value),
    ]) {
        assert.ok(
            !diagnostics.includes(value),
            "diagnostics must not disclose credentials",
        );
    }
    await fs.writeFile(
        path.join(artifacts, "serve-acceptance.json"),
        JSON.stringify(
            {
                mutations,
                reads,
                themes,
                checks: [
                    "preloaded policy",
                    "file picker",
                    "non-mutating plan",
                    "exact subset",
                    "stale and fabricated selections",
                    "access control",
                    "live quota logs",
                    "reload",
                    "resume exact unfinished selection",
                    "resume after cancellation",
                    "missing targets and lost responses",
                    "cancellation",
                    "API failure",
                    "escaped profile text",
                    "responsive theme",
                    "filtered bulk selection and hidden counts",
                    "readable exclusions and expandable truncated previews",
                ],
            },
            null,
            2,
        ),
    );
} finally {
    if (browser) await browser.close();
    child.kill("SIGINT");
    let timer;
    await Promise.race([
        exited,
        new Promise((resolve) => {
            timer = setTimeout(() => {
                child.kill("SIGKILL");
                resolve();
            }, 3000);
        }),
    ]);
    clearTimeout(timer);
    fixture.closeAllConnections();
    await new Promise((resolve) => fixture.close(resolve));
    await fs.rm(scratch, { recursive: true, force: true });
}
