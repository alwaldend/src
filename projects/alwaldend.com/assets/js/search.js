// Site-local search uses Docsy's Bootstrap presentation and a prebuilt worker index.
(() => {
    function initialize() {
        const inputs = [
            ...document.querySelectorAll(".td-search--offline input"),
        ];
        if (!inputs.length) return;
        let worker;
        let ready;
        let sequence = 0;
        let generation = 0;
        const pending = new Map();
        const submitted = new WeakMap();
        const failureMessage = "Search is unavailable. Press Enter to retry.";

        function reset(error) {
            worker?.terminate();
            worker = undefined;
            ready = undefined;
            for (const request of pending.values()) request.reject(error);
            pending.clear();
        }
        function request(type, data = {}) {
            return new Promise((resolve, reject) => {
                const id = ++sequence;
                pending.set(id, { resolve, reject });
                worker.postMessage({ id, type, ...data });
            });
        }
        function load() {
            if (ready) return ready;
            try {
                worker = new Worker(inputs[0].dataset.searchWorkerSrc);
                worker.onmessage = ({ data }) => {
                    const request = pending.get(data.id);
                    if (!request) return;
                    pending.delete(data.id);
                    if (data.error) request.reject(new Error(data.error));
                    else request.resolve(data.results);
                };
                worker.onerror = () => reset(new Error(failureMessage));
                ready = request("initialize", {
                    url: new URL(
                        inputs[0].dataset.searchIndexSrc,
                        location.href,
                    ).href,
                }).catch((error) => {
                    reset(error);
                    throw error;
                });
                return ready;
            } catch (error) {
                reset(error);
                return Promise.reject(error);
            }
        }
        function dismiss() {
            for (const input of inputs) {
                bootstrap.Popover.getInstance(input)?.dispose();
                input.removeAttribute("aria-busy");
            }
        }
        function show(input, query, message, results = []) {
            dismiss();
            const content = document.createElement("div");
            content.setAttribute("role", "status");
            const header = document.createElement("div");
            header.className = "d-flex justify-content-between mb-3";
            const title = document.createElement("strong");
            title.textContent = "Search results";
            const close = document.createElement("button");
            close.type = "button";
            close.className = "btn-close";
            close.setAttribute("aria-label", "Close search results");
            close.addEventListener("click", () => {
                ++generation;
                input.value = "";
                dismiss();
            });
            header.append(title, close);
            content.append(header);
            const body = document.createElement("div");
            body.style.maxHeight = `calc(100vh - ${input.getBoundingClientRect().top + 180}px)`;
            body.style.overflowY = "auto";
            if (message) {
                const status = document.createElement("p");
                status.textContent = message;
                body.append(status);
            } else if (!results.length) {
                const status = document.createElement("p");
                status.textContent = `No results found for query "${query}"`;
                body.append(status);
            } else {
                for (const result of results) {
                    const entry = document.createElement("div");
                    entry.className = "mt-4";
                    const reference = document.createElement("small");
                    reference.className = "d-block text-body-secondary";
                    reference.textContent = result.ref;
                    const link = document.createElement("a");
                    link.className = "d-block";
                    link.style.fontSize = "1.2rem";
                    link.href = result.ref;
                    link.textContent = result.title;
                    const excerpt = document.createElement("p");
                    excerpt.textContent = result.excerpt;
                    entry.append(reference, link, excerpt);
                    body.append(entry);
                }
            }
            content.append(body);
            new bootstrap.Popover(input, {
                content,
                html: true,
                customClass: "td-offline-search-results",
                placement: "bottom",
                trigger: "manual",
            }).show();
        }
        async function search(input) {
            const current = ++generation;
            const query = input.value.trim();
            submitted.set(input, query);
            if (!query) {
                dismiss();
                return;
            }
            show(input, query, "Loading search…");
            input.setAttribute("aria-busy", "true");
            try {
                await load();
                if (current !== generation) return;
                const limit =
                    parseInt(input.dataset.offlineSearchMaxResults, 10) || 10;
                const results = await request("query", { query, limit });
                if (current === generation) show(input, query, null, results);
            } catch (error) {
                console.warn("Local search failed:", error);
                if (current === generation) show(input, query, failureMessage);
            }
        }
        for (const input of inputs) {
            input.addEventListener("focus", () => {
                load().catch(() => {});
            });
            input.addEventListener("change", () => {
                // Enter submits without committing the native change event.
                // Keep result links intact when focus leaves that same query.
                if (submitted.get(input) !== input.value.trim()) search(input);
            });
            input.addEventListener("input", () => {
                submitted.delete(input);
                ++generation;
                dismiss();
            });
            input.addEventListener("keydown", (event) => {
                if (event.key === "Enter") {
                    event.preventDefault();
                    search(input);
                }
                if (event.key === "Escape") {
                    ++generation;
                    dismiss();
                }
            });
            input.closest("form")?.addEventListener("submit", (event) => {
                event.preventDefault();
                event.stopPropagation();
                search(input);
            });
        }
    }
    if (document.readyState === "loading")
        document.addEventListener("DOMContentLoaded", initialize, {
            once: true,
        });
    else initialize();
})();
