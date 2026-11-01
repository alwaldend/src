import { roundBends, routeGraph } from "./routing.mjs";
import { sequenceModel } from "./sequence.mjs";
import { measuredText } from "./text.mjs";

// Mermaid owns grammar and semantic data. This adapter consumes the pinned
// database interfaces with display text normalized by the parsing browser.
function graphText(item) {
    const label = item.label ?? item.id;
    const lines = [
        label,
        ...(item.annotations ?? []).map((value) => `«${value}»`),
        ...(item.attributes ?? []).map(
            (value) =>
                `${value.type} ${value.name}${value.keys?.length ? `  ${value.keys.join(", ")}` : ""}`,
        ),
        ...(item.members ?? []).map((value) => value.text),
        ...(item.methods ?? []).map((value) => value.text),
    ].flatMap((line) => line.split("\n"));
    const compartments =
        lines.length > 1
            ? [
                  label.split("\n").length + (item.annotations?.length ?? 0),
                  ...(item.members?.length && item.methods?.length
                      ? [
                            label.split("\n").length +
                                (item.annotations?.length ?? 0) +
                                item.members.length,
                        ]
                      : []),
              ]
            : [];
    return { label, lines, compartments };
}

export function diagramText(raw) {
    if (raw.type === "sequence")
        return [
            ...Object.values(raw.values.getActors).map(
                (actor) => actor.description ?? actor.name,
            ),
            ...raw.values.getMessages.map((event) => event.message),
        ].flatMap((value) => (value ?? "").split("\n"));
    return (raw.values.getData?.nodes ?? []).flatMap(
        (item) => graphText(item).lines,
    );
}

export function adaptDiagram(raw, measurements) {
    const viewDirection = [0.28, 0.2, 1.8];
    if (raw.type === "sequence") {
        const model = {
            ...sequenceModel(raw.values, measurements),
            viewDirection,
        };
        for (const edge of model.edges)
            edge.points = roundBends(edge.points, 36);
        return model;
    }
    if (
        ![
            "flowchart-v2",
            "stateDiagram",
            "classDiagram",
            "er",
            "mindmap",
        ].includes(raw.type)
    )
        throw new Error(`Unsupported Mermaid diagram type: ${raw.type}`);
    const data = raw.values.getData;
    if (!data?.nodes || !data?.edges)
        throw new Error(`Unsupported Mermaid database: ${raw.type}`);
    const nodes = [];
    const groups = [];
    for (const item of data.nodes) {
        const { label, lines, compartments } = graphText(item);
        const terminal = ["stateStart", "stateEnd"].includes(item.shape);
        // A rectangle fits inside a diamond with twice its width and height.
        const textScale = item.shape === "diamond" ? 2 : 1;
        const text = measuredText(lines, measurements, {
            size: item.isGroup ? 14 : 16,
            heading: !item.isGroup && lines.length > 1,
            compartments: item.isGroup ? [] : compartments,
        });
        const node = {
            id: String(item.id),
            label: terminal
                ? item.shape === "stateStart"
                    ? "Start"
                    : "End"
                : label,
            lines: terminal ? [] : lines,
            textLayout: text,
            shape: item.shape ?? "rect",
            parentId: item.parentId ? String(item.parentId) : null,
            position: [0, 0, 0],
            size: terminal
                ? [28, 28, 28]
                : [
                      Math.max(144, (text.width + 24) * textScale),
                      Math.max(58, (text.height + 18) * textScale),
                      40,
                  ],
        };
        if (item.isGroup)
            groups.push({ ...node, members: [], direction: item.dir });
        else nodes.push(node);
    }
    const edges = data.edges
        .filter(
            (edge) =>
                !edge.isLayoutOnly &&
                edge.pattern !== "invisible" &&
                edge.thickness !== "invisible",
        )
        .map((edge, index) => ({
            id: String(edge.id ?? `edge-${index}`),
            start: String(edge.start),
            end: String(edge.end),
            label: edge.label ?? "",
            startLabel: edge.startLabelRight ?? edge.startLabelLeft ?? "",
            endLabel: edge.endLabelLeft ?? edge.endLabelRight ?? "",
            startMarker: edge.arrowTypeStart ?? "none",
            endMarker: edge.arrowTypeEnd ?? "none",
            dashed: ["dotted", "dashed"].includes(edge.pattern),
        }));
    const model = {
        type: raw.type,
        viewDirection,
        direction: raw.values.getDirection ?? data.direction ?? "LR",
        nodes,
        groups,
        edges,
        notes: [],
        guides: [],
    };
    placeGraph(model, data.edges);
    routeGraph(model);
    return model;
}

function bounds(items, padding = [0, 0, 0]) {
    const lower = [Infinity, Infinity, Infinity];
    const upper = [-Infinity, -Infinity, -Infinity];
    for (const item of items)
        for (let axis = 0; axis < 3; axis += 1) {
            lower[axis] = Math.min(
                lower[axis],
                item.position[axis] - item.size[axis] / 2 - padding[axis],
            );
            upper[axis] = Math.max(
                upper[axis],
                item.position[axis] + item.size[axis] / 2 + padding[axis],
            );
        }
    if (!items.length) return { position: [0, 0, 0], size: [120, 80, 80] };
    return {
        position: lower.map((value, axis) => (value + upper[axis]) / 2),
        size: lower.map((value, axis) => upper[axis] - value),
    };
}

function placeGraph(model, layoutEdges) {
    const all = new Map(
        [...model.nodes, ...model.groups].map((item) => [item.id, item]),
    );
    for (const item of all.values())
        if (item.parentId && !all.has(item.parentId))
            throw new Error(`Missing parent ${item.parentId}`);
    const translate = (item, delta) => {
        item.position = item.position.map(
            (value, axis) => value + delta[axis],
        );
        for (const id of item.members ?? []) translate(all.get(id), delta);
    };
    const active = new Set();
    function place(parentId, direction) {
        if (active.has(parentId)) throw new Error("Cyclic group containment");
        active.add(parentId);
        const children = [...all.values()].filter(
            (item) => item.parentId === parentId,
        );
        for (const child of children)
            if (child.members) {
                place(child.id, child.direction ?? direction);
                const members = [...all.values()].filter(
                    (item) => item.parentId === child.id,
                );
                child.members = members.map((item) => item.id);
                const box = bounds(members, [
                    42,
                    Math.max(55, child.textLayout.height / 2 + 28),
                    28,
                ]);
                box.size[0] = Math.max(
                    box.size[0],
                    child.textLayout.width + 44,
                );
                // Reserve a front routing corridor inside each volume.
                const corridor = 48 + model.edges.length * 4;
                child.position = box.position.map(
                    (value, axis) => value + (axis === 2 ? corridor / 2 : 0),
                );
                child.size = box.size.map(
                    (value, axis) => value + (axis === 2 ? corridor : 0),
                );
            }
        const direct = (id) => {
            let current = all.get(id);
            while (current && current.parentId !== parentId)
                current = all.get(current.parentId);
            return current?.id;
        };
        const adjacency = new Map(children.map((item) => [item.id, []]));
        for (const edge of layoutEdges) {
            const a = direct(String(edge.start));
            const b = direct(String(edge.end));
            if (a && b && a !== b) adjacency.get(a).push(b);
        }
        // Condense strongly connected components before ranking. Feedback
        // cycles remain intact and cannot make longest-path ranking diverge.
        let next = 0;
        const indices = new Map();
        const low = new Map();
        const stack = [];
        const onStack = new Set();
        const components = [];
        function visit(id) {
            indices.set(id, next);
            low.set(id, next++);
            stack.push(id);
            onStack.add(id);
            for (const target of adjacency.get(id)) {
                if (!indices.has(target)) {
                    visit(target);
                    low.set(id, Math.min(low.get(id), low.get(target)));
                } else if (onStack.has(target))
                    low.set(id, Math.min(low.get(id), indices.get(target)));
            }
            if (indices.get(id) === low.get(id)) {
                const component = [];
                let member;
                do {
                    member = stack.pop();
                    onStack.delete(member);
                    component.push(member);
                } while (member !== id);
                components.push(component.reverse());
            }
        }
        for (const child of children)
            if (!indices.has(child.id)) visit(child.id);
        const componentOf = new Map(
            components.flatMap((component, i) =>
                component.map((id) => [id, i]),
            ),
        );
        const ranks = components.map(() => 0);
        // Tarjan emits sinks first, so reverse order is topological.
        for (let i = components.length - 1; i >= 0; i -= 1)
            for (const id of components[i])
                for (const target of adjacency.get(id)) {
                    const j = componentOf.get(target);
                    if (i !== j) ranks[j] = Math.max(ranks[j], ranks[i] + 1);
                }
        const columns = new Map();
        for (const child of children) {
            const rank = ranks[componentOf.get(child.id)];
            if (!columns.has(rank)) columns.set(rank, []);
            columns.get(rank).push(child);
        }
        const vertical = ["TB", "TD", "BT"].includes(direction);
        const main = vertical ? 1 : 0;
        const cross = vertical ? 0 : 1;
        let cursor = 0;
        for (const [rank, column] of [...columns].sort(
            (a, b) => a[0] - b[0],
        )) {
            const extent = Math.max(...column.map((item) => item.size[main]));
            const total = column.reduce(
                (sum, item) => sum + item.size[cross] + 76,
                -76,
            );
            let across = -total / 2;
            for (let i = 0; i < column.length; i += 1) {
                const item = column[i];
                const position = [0, 0, (((rank + i) % 3) - 1) * 44];
                position[main] = cursor + extent / 2;
                position[cross] = across + item.size[cross] / 2;
                if (vertical) position[1] *= -1;
                if (direction === "RL" || direction === "BT")
                    position[main] *= -1;
                translate(
                    item,
                    position.map((value, axis) => value - item.position[axis]),
                );
                across += item.size[cross] + 76;
            }
            cursor += extent + 110;
        }
        active.delete(parentId);
    }
    place(null, model.direction);
    const center = bounds(
        [...all.values()].filter((item) => !item.parentId),
    ).position;
    for (const item of all.values())
        if (!item.parentId)
            translate(
                item,
                center.map((value) => -value),
            );
}
