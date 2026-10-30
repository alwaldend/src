"use strict";

const element = (id) => document.getElementById(id);
const source = element("policy");
const dialog = element("confirmation");
const fragment = location.hash.slice(1);
if (/^[a-f0-9]{64}$/.test(fragment)) {
    sessionStorage.setItem("cleaner-access", fragment);
    history.replaceState(null, "", location.pathname);
}
const capability = sessionStorage.getItem("cleaner-access") || "";
let state = { phase: "idle", job_id: "", actions: 0, completed: 0 };
let connected = false;
let pending = false;
let currentPlan = null;
let currentPlanID = "";
let plannedSource = null;
let selected = new Set();
let cursor = 0;
let lines = [];
let requestError = "";
let connectionError = "";

async function api(route, body) {
    const response = await fetch(route, {
        method: body === undefined ? "GET" : "POST",
        headers: {
            "X-Cleaner-Token": capability,
            ...(body === undefined
                ? {}
                : { "Content-Type": "application/json" }),
        },
        body: body === undefined ? undefined : JSON.stringify(body),
    });
    const result = await response.json();
    if (!response.ok)
        throw new Error(result.error || `Request failed (${response.status})`);
    return result;
}

function busy() {
    return pending || state.phase === "planning" || state.phase === "applying";
}

function canApply() {
    return (
        connected &&
        !busy() &&
        state.phase === "ready" &&
        currentPlanID === state.plan_id &&
        selected.size > 0 &&
        source.value === plannedSource
    );
}

function canResume() {
    return (
        connected &&
        !busy() &&
        state.remaining > 0 &&
        ["failed", "cancelled"].includes(state.phase) &&
        currentPlanID === state.plan_id
    );
}

function acceptState(next) {
    const order = {
        idle: 0,
        planning: 1,
        applying: 1,
        ready: 2,
        completed: 2,
        failed: 2,
        cancelled: 2,
    };
    if (Number(next.job_id) < Number(state.job_id)) return;
    if (next.job_id === state.job_id && order[next.phase] < order[state.phase])
        return;
    state = next;
}

function render() {
    const labels = {
        idle: "Ready to plan",
        planning: "Preparing your plan",
        ready: "",
        applying: `Completed ${state.completed} of ${state.actions} actions`,
        completed: `Completed ${state.completed} of ${state.actions} actions`,
        failed: state.remaining
            ? `Stopped · ${state.remaining} actions remaining`
            : "Operation failed",
        cancelled: state.remaining
            ? `Cancelled · ${state.remaining} actions remaining`
            : "Operation cancelled",
    };
    element("status").textContent = connected
        ? labels[state.phase]
        : "Not connected";
    element("status").dataset.phase = state.phase;
    element("operation-status").hidden =
        connected && ["idle", "ready"].includes(state.phase);
    element("plan").disabled = !connected || busy();
    element("plan").textContent =
        state.phase === "planning" ? "Planning…" : "Run plan";
    source.disabled = !connected || busy();
    element("policy-file").disabled = !connected || busy();
    element("load-policy").disabled = !connected || busy();
    element("cancel-operation").hidden =
        state.phase !== "planning" && state.phase !== "applying";
    element("cancel-operation").disabled = pending;
    element("apply").disabled = !canApply() && !canResume();
    element("apply").textContent = canResume() ? "Resume" : "Apply selected";
    element("confirm-apply").disabled = !canApply();
    const executing = state.plan_id && state.phase !== "ready";
    element("selected-count").textContent = executing
        ? `${state.completed} of ${state.actions} completed`
        : currentPlan
          ? `${selected.size} of ${currentPlan.actions.length} selected`
          : "";
    const dirty = currentPlan && source.value !== plannedSource;
    element("review-note").textContent = dirty
        ? executing
            ? "Policy changed. Resume uses the original selection; Run plan replaces it."
            : "Policy changed. Run plan again to update this list."
        : "";
    element("review-note").hidden = !dirty;
    const message = connectionError || requestError || state.error || "";
    element("error").textContent = message;
    element("error").hidden = !message;
    const rows = [...document.querySelectorAll("#items tbody tr")];
    const editable =
        connected &&
        !busy() &&
        !dirty &&
        state.phase === "ready" &&
        currentPlanID === state.plan_id;
    const matching = rows.filter((row) => !row.hidden).length;
    let hiddenSelected = 0;
    const completed = new Set(
        executing && currentPlanID === state.plan_id
            ? state.completed_indexes || []
            : [],
    );
    rows.forEach((row, index) => {
        if (completed.has(index)) selected.delete(index);
        const included = selected.has(index);
        const input = row.querySelector("input");
        input.disabled = !editable;
        input.checked = included;
        row.classList.toggle("excluded", !included && !completed.has(index));
        const action = currentPlan.actions[index];
        row.querySelector(".action-label").textContent = completed.has(index)
            ? "Completed"
            : included
              ? action.kind === "unrepost"
                  ? "Undo repost"
                  : "Delete post"
              : "Excluded";
        const forced = row.querySelector(".forced-label");
        if (forced) forced.hidden = !included || completed.has(index);
        if (row.hidden && included) hiddenSelected++;
    });
    element("select-all").disabled = !editable || matching === 0;
    element("exclude-all").disabled = !editable || matching === 0;
    const filtered = element("filter").value !== "";
    element("select-all").textContent = filtered
        ? "Select matching"
        : "Select all";
    element("exclude-all").textContent = filtered
        ? "Exclude matching"
        : "Exclude all";
    element("filter-summary").hidden = !filtered || rows.length === 0;
    element("filter-summary").textContent =
        `${matching} of ${rows.length} items match the filter.`;
    element("selection-note").textContent = hiddenSelected
        ? `${quantity(hiddenSelected, "selected item")} hidden by the filter. Apply includes them.`
        : "";
    element("selection-note").hidden = hiddenSelected === 0;
    if (!canApply() && dialog.open) dialog.close();
}

function quantity(count, singular, plural = `${singular}s`) {
    return `${count} ${count === 1 ? singular : plural}`;
}

function node(tag, className, text) {
    const result = document.createElement(tag);
    if (className) result.className = className;
    if (text !== undefined) result.textContent = text;
    return result;
}

function showPlan(plan, id, policy) {
    currentPlan = plan;
    currentPlanID = id;
    plannedSource = policy;
    selected = new Set(plan.actions.map((_, index) => index));
    const decisions = new Map(
        (plan.decisions || []).map((decision) => [
            `${decision.kind === "repost" ? "unrepost" : "delete_post"}:${decision.source_id || decision.item_id}`,
            decision,
        ]),
    );
    const rows = document.createDocumentFragment();
    plan.actions.forEach((action, index) => {
        const decision = decisions.get(`${action.kind}:${action.target_id}`);
        const row = node("tr");
        row.dataset.target = action.target_id;
        row.dataset.search =
            `${decision?.text || ""} ${action.target_id}`.toLowerCase();
        const checkCell = node("td");
        const checkbox = node("input", "form-check-input");
        checkbox.type = "checkbox";
        checkbox.checked = true;
        checkbox.setAttribute(
            "aria-label",
            `Include ${action.kind === "unrepost" ? "repost" : "post"} ${action.target_id}`,
        );
        checkbox.addEventListener("change", () => {
            if (checkbox.checked) selected.add(index);
            else selected.delete(index);
            render();
        });
        const selectionControl = node("label", "selection-control");
        selectionControl.append(checkbox);
        checkCell.append(selectionControl);
        const activity = node("td");
        const text = decision?.text || "Content unavailable.";
        const characters = Array.from(text.replace(/\s+/g, " ").trim());
        const preview =
            characters.slice(0, 240).join("") +
            (characters.length > 240 ? "…" : "");
        const paragraph = node("div", "item-preview");
        const content = node("span", "item-text", preview);
        paragraph.append(content);
        activity.append(paragraph);
        if (preview !== text) {
            const expand = node(
                "button",
                "btn btn-link expand-post",
                "Show more",
            );
            expand.type = "button";
            expand.setAttribute("aria-expanded", "false");
            expand.addEventListener("click", () => {
                const expanded =
                    expand.getAttribute("aria-expanded") === "true";
                content.textContent = expanded ? preview : text;
                expand.textContent = expanded ? "Show more" : "Show less";
                expand.setAttribute("aria-expanded", String(!expanded));
            });
            paragraph.append(expand);
        }
        const metadata = node("div", "item-meta");
        if (decision?.created_at)
            metadata.append(
                node(
                    "span",
                    "",
                    new Date(decision.created_at).toLocaleString(),
                ),
            );
        const link = node("a", "", action.target_id);
        link.href = `https://x.com/i/status/${action.target_id}`;
        link.target = "_blank";
        link.rel = "noreferrer noopener";
        metadata.append(link);
        activity.append(metadata);
        const reasons = decision?.reasons || [];
        if (reasons.length)
            metadata.append(
                node(
                    "span",
                    "item-rule",
                    reasons
                        .map((reason) => reason.name || reason.operation)
                        .join(" · "),
                ),
            );
        const actionCell = node("td", "item-action");
        actionCell.append(
            node(
                "span",
                "action-label",
                action.kind === "unrepost" ? "Undo repost" : "Delete post",
            ),
        );
        if (action.forced)
            actionCell.append(
                node("span", "forced-label", "Forced by policy"),
            );
        row.append(checkCell, activity, actionCell);
        rows.append(row);
    });
    element("items").querySelector("tbody").replaceChildren(rows);
    element("items").hidden = plan.actions.length === 0;
    element("empty-state").hidden = plan.actions.length !== 0;
    if (!plan.actions.length)
        element("empty-state").replaceChildren(
            node("p", "", "No items selected for removal."),
            node(
                "span",
                "",
                "Your policy leaves the collected items unchanged.",
            ),
        );
    element("review-tools").hidden = plan.actions.length === 0;
    filterRows();
}

function filterRows() {
    const query = element("filter").value.toLowerCase();
    for (const row of document.querySelectorAll("#items tbody tr"))
        row.hidden = !row.dataset.search.includes(query);
    render();
}

function appendLogs(result) {
    const fresh = result.logs.filter((entry) => entry.id > cursor);
    cursor = Math.max(cursor, result.cursor);
    if (!fresh.length) return;
    lines = [...lines, ...fresh.map((entry) => entry.text)].slice(-1000);
    element("logs").textContent = lines.join("\n");
    if (element("follow-logs").checked)
        element("logs").scrollTop = element("logs").scrollHeight;
}

async function poll() {
    try {
        const result = await api(`/api/state?after=${cursor}`);
        connected = true;
        connectionError = "";
        acceptState(result);
        appendLogs(result);
        if (state.plan_id && currentPlanID !== state.plan_id) {
            const planned = await api("/api/plan");
            if (state.plan_id === planned.plan_id)
                showPlan(planned.plan, planned.plan_id, planned.policy);
        }
    } catch (error) {
        connected = false;
        connectionError = `${error.message}. Open the private link again if the server was restarted.`;
    }
    render();
    setTimeout(poll, 750);
}

async function operate(route, body) {
    pending = true;
    requestError = "";
    render();
    try {
        acceptState(await api(route, body));
    } catch (error) {
        requestError = error.message;
    } finally {
        pending = false;
        render();
    }
}

source.addEventListener("input", render);
element("load-policy").addEventListener("click", () =>
    element("policy-file").click(),
);
element("policy-file").addEventListener("change", async (event) => {
    const file = event.target.files[0];
    if (!file || busy()) return;
    if (file.size > 1024 * 1024) {
        requestError = "The policy file exceeds 1 MiB.";
    } else {
        try {
            source.value = await file.text();
            requestError = "";
        } catch (error) {
            requestError = `Cannot read policy: ${error.message}`;
        }
    }
    render();
});
element("plan").addEventListener("click", async () => {
    plannedSource = source.value;
    currentPlanID = "";
    await operate("/api/plan", { policy: source.value });
});
element("cancel-operation").addEventListener("click", () =>
    operate("/api/cancel", { job_id: state.job_id }),
);
element("apply").addEventListener("click", () => {
    if (canResume()) {
        operate("/api/resume", { job_id: state.job_id });
        return;
    }
    if (!canApply()) return;
    const actions = currentPlan.actions.filter((_, index) =>
        selected.has(index),
    );
    const posts = actions.filter(
        (action) => action.kind === "delete_post",
    ).length;
    element("confirmation-summary").textContent =
        `Delete ${quantity(posts, "post or reply", "posts or replies")} and undo ${quantity(actions.length - posts, "repost")}.`;
    element("confirm-apply").textContent =
        `Apply ${quantity(actions.length, "action")}`;
    dialog.showModal();
});
element("dismiss-confirmation").addEventListener("click", () =>
    dialog.close(),
);
element("confirm-apply").addEventListener("click", async () => {
    if (!canApply()) return;
    const request = {
        plan_id: currentPlanID,
        selected: [...selected].sort((a, b) => a - b),
    };
    dialog.close();
    await operate("/api/apply", request);
});
for (const [id, checked] of [
    ["select-all", true],
    ["exclude-all", false],
]) {
    element(id).addEventListener("click", () => {
        for (const [index, input] of [
            ...document.querySelectorAll("#items input"),
        ].entries()) {
            if (input.disabled || input.closest("tr").hidden) continue;
            if (checked) selected.add(index);
            else selected.delete(index);
        }
        render();
    });
}
element("filter").addEventListener("input", filterRows);
const preferred = matchMedia("(prefers-color-scheme: dark)");
function applyTheme() {
    const choice = element("theme").value;
    document.documentElement.dataset.bsTheme =
        choice === "system" ? (preferred.matches ? "dark" : "light") : choice;
    localStorage.setItem("cleaner-theme", choice);
}
element("theme").value = localStorage.getItem("cleaner-theme") || "system";
element("theme").addEventListener("change", applyTheme);
preferred.addEventListener("change", applyTheme);
applyTheme();
async function initialize() {
    try {
        const session = await api("/api/session");
        source.value = session.policy;
        connected = true;
    } catch (error) {
        requestError = error.message;
    }
    await poll();
}
initialize();
