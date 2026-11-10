/** Display panel prompts as an automatically growing list of text boxes. */
import { app } from "../../scripts/app.js";

const editors = new WeakMap();

function installEditor(node) {
    const original = node.widgets.find((widget) => widget.name === "prompts");
    const separator = node.widgets.find(
        (widget) => widget.name === "separator",
    );
    if (!original || !separator || editors.has(node)) return;
    let raw = original.value;
    const container = document.createElement("div");
    Object.assign(container.style, {
        display: "flex",
        flexDirection: "column",
        gap: "6px",
        overflowY: "auto",
        padding: "4px",
        boxSizing: "border-box",
    });

    function appendField(value = "") {
        const field = document.createElement("textarea");
        const index = container.children.length + 1;
        field.value = value;
        field.setAttribute("aria-label", `Panel ${index}`);
        field.placeholder = `Panel ${index}`;
        field.rows = 3;
        Object.assign(field.style, {
            width: "100%",
            minHeight: "64px",
            flexShrink: "0",
            boxSizing: "border-box",
            resize: "vertical",
            color: "var(--input-text, #ddd)",
            background: "var(--comfy-input-bg, #222)",
            border: "1px solid var(--border-color, #555)",
            borderRadius: "4px",
            padding: "8px",
            font: "inherit",
        });
        field.addEventListener("input", () => {
            if (field === container.lastElementChild && field.value.trim()) {
                appendField();
            }
            raw = JSON.stringify(
                Array.from(container.children)
                    .map((input) => input.value)
                    .filter((value) => value.trim()),
            );
            separator.value = "JSON";
            widget.callback?.(raw);
            node.graph?.change();
            node.setDirtyCanvas(true, true);
        });
        container.append(field);
    }

    function render() {
        let panels;
        if (separator.value === "JSON") {
            try {
                panels = JSON.parse(raw);
                if (
                    !Array.isArray(panels) ||
                    !panels.every((item) => typeof item === "string")
                ) {
                    throw new Error("Invalid panel list");
                }
            } catch {
                // Keep malformed input intact so backend validation can report it.
                panels = [raw];
            }
        } else {
            const normalized = raw.replace(/\r\n?/g, "\n");
            panels = normalized.split(
                separator.value === "blank line" ? /\n\s*\n/ : /\n/,
            );
        }
        container.replaceChildren();
        panels
            .filter((panel) => panel.trim())
            .forEach((panel) => appendField(panel));
        appendField();
    }

    const index = node.widgets.indexOf(original);
    const widget = node.addDOMWidget("prompts", "custom", container, {
        getValue: () => raw,
        setValue: (value) => {
            raw = value;
            render();
        },
        getMinHeight: () => 220,
        getMaxHeight: () => 420,
    });
    node.widgets.splice(node.widgets.indexOf(widget), 1);
    node.widgets.splice(index, 1, widget);
    original.onRemove?.();
    original.element?.remove();
    separator.hidden = true;
    editors.set(node, render);
    render();
    node.setSize(node.computeSize());
}

app.registerExtension({
    name: "alwaldend.prompt_sequence",
    nodeCreated(node) {
        if (node.comfyClass === "ImagegenPromptSequence") installEditor(node);
    },
    loadedGraphNode(node) {
        editors.get(node)?.();
    },
});
