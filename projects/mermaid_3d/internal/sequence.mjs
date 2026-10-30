// The pinned Mermaid database owns event codes and their order. Time remains
// vertical; depth separates participant planes and message routing corridors.
import { measuredText } from "./text.mjs";

export function sequenceModel(values, measurements) {
    for (const field of ["getBoxes", "getCreatedActors", "getDestroyedActors"])
        if (Object.keys(values[field] ?? {}).length)
            throw new Error(
                `Sequence ${field} is not supported by the spatial adapter`,
            );
    const types = new Map(
        Object.entries(values.lineTypes).map(([name, code]) => [code, name]),
    );
    const nodes = values.getActorKeys.map((id, index) => {
        const actor = values.getActors[id];
        const label = actor.description ?? actor.name;
        const lines = label.split("\n");
        const text = measuredText(lines, measurements, {
            heading: lines.length > 1,
        });
        return {
            id,
            label,
            lines,
            textLayout: text,
            shape: "rect",
            parentId: null,
            position: [0, 0, index * 24],
            size: [
                Math.max(164, text.width + 24),
                Math.max(62, text.height + 22),
                48,
            ],
        };
    });
    for (let i = 1; i < nodes.length; i += 1)
        nodes[i].position[0] =
            nodes[i - 1].position[0] +
            (nodes[i - 1].size[0] + nodes[i].size[0]) / 2 +
            126;
    const actors = new Map(nodes.map((node) => [node.id, node]));
    const model = {
        type: "sequence",
        direction: "TB",
        nodes,
        groups: [],
        edges: [],
        notes: [],
        guides: [],
    };
    const blocks = [];
    const activations = new Map();
    let y = -Math.max(100, ...nodes.map((node) => node.size[1] / 2 + 50));
    let eventY = y;
    const front = Math.max(...nodes.map((node) => node.position[2])) + 24;
    const endpoint = (id, time, depth = 0) => {
        const actor = actors.get(id);
        if (!actor) throw new Error(`Missing sequence participant ${id}`);
        return [
            actor.position[0] + depth * 10,
            time,
            actor.position[2] + actor.size[2] / 2,
        ];
    };
    const activeDepth = (id) =>
        Math.max(0, (activations.get(id)?.length ?? 0) - 1);
    for (const [index, event] of values.getMessages.entries()) {
        const type = types.get(event.type);
        const label = event.message;
        if (type === "ACTIVE_START") {
            const stack = activations.get(event.from) ?? [];
            stack.push({ y: eventY, depth: stack.length });
            activations.set(event.from, stack);
        } else if (type === "ACTIVE_END") {
            const start = activations.get(event.from)?.pop();
            if (start === undefined)
                throw new Error("Unmatched sequence activation");
            model.guides.push({
                activation: true,
                actor: event.from,
                points: [
                    endpoint(event.from, start.y, start.depth),
                    endpoint(event.from, eventY, start.depth),
                ],
            });
        } else if (
            /^(LOOP|ALT|OPT|PAR|CRITICAL|BREAK|RECT)_START$/.test(type)
        ) {
            const group = {
                id: `block-${event.id}`,
                label: `${type.split("_")[0].toLowerCase()}: ${label}`,
                parentId: blocks.at(-1)?.id ?? null,
                members: [],
                edgeIds: [],
                sections: [],
                top: y,
                bottom: y,
                position: [0, 0, 0],
                size: [1, 1, 1],
            };
            model.groups.push(group);
            blocks.push(group);
            eventY = y;
            y -= 70;
        } else if (["ALT_ELSE", "PAR_AND", "CRITICAL_OPTION"].includes(type)) {
            const group = blocks.at(-1);
            if (!group)
                throw new Error("Sequence branch has no enclosing block");
            group.sections.push({ label, y });
            eventY = y;
            y -= 70;
        } else if (/^(LOOP|ALT|OPT|PAR|CRITICAL|BREAK|RECT)_END$/.test(type)) {
            const group = blocks.pop();
            if (!group) throw new Error("Unmatched sequence block end");
            group.bottom = y + 15;
            eventY = group.bottom;
            y -= 40;
        } else if (type === "NOTE") {
            const lines = label.split("\n");
            const text = measuredText(lines, measurements, {
                heading: lines.length > 1,
            });
            const width = Math.max(208, text.width + 24);
            const height = Math.max(64, text.height + 24);
            const centerY = y - (height - 64) / 2;
            const a = endpoint(event.from, centerY);
            const b = endpoint(event.to ?? event.from, centerY);
            let x = (a[0] + b[0]) / 2;
            if (event.placement === 0) x = a[0] - width / 2 - 25;
            if (event.placement === 1) x = a[0] + width / 2 + 25;
            model.notes.push({
                id: `note-${event.id}`,
                parentId: blocks.at(-1)?.id ?? null,
                label,
                lines,
                textLayout: text,
                shape: "rect",
                position: [x, centerY, front + 65],
                size: [width, height, 24],
            });
            eventY = centerY - height / 2;
            y -= height + 31;
        } else if (
            [
                "SOLID",
                "DOTTED",
                "SOLID_OPEN",
                "DOTTED_OPEN",
                "BIDIRECTIONAL_SOLID",
                "BIDIRECTIONAL_DOTTED",
            ].includes(type)
        ) {
            const start = endpoint(event.from, y, activeDepth(event.from));
            const self = event.from === event.to;
            const next = values.getMessages[index + 1];
            const activatesTarget =
                next?.from === event.to &&
                types.get(next.type) === "ACTIVE_START";
            const endDepth = activatesTarget
                ? (activations.get(event.to)?.length ?? 0)
                : activeDepth(event.to);
            const end = endpoint(event.to, y - (self ? 32 : 0), endDepth);
            const z = start[2] + 28;
            const points = self
                ? [
                      start,
                      [start[0], y, z],
                      [start[0] + 110, y, z],
                      [start[0] + 110, end[1], z],
                      [end[0], end[1], z],
                      end,
                  ]
                : [start, end];
            const edge = {
                id: `message-${event.id}`,
                label,
                start: event.from,
                end: event.to,
                points,
                startMarker: type.startsWith("BIDIRECTIONAL")
                    ? "point"
                    : "none",
                endMarker: type.endsWith("OPEN") ? "none" : "point",
                dashed: type.includes("DOTTED"),
                labelPosition: [
                    self ? start[0] + 110 : (start[0] + end[0]) / 2,
                    y + 22,
                    self ? z + 5 : (start[2] + end[2]) / 2,
                ],
            };
            model.edges.push(edge);
            for (const group of blocks) group.edgeIds.push(edge.id);
            eventY = end[1];
            y -= self ? 125 : 95;
        } else {
            throw new Error(
                `Unsupported sequence event: ${type ?? event.type}`,
            );
        }
    }
    if (
        blocks.length ||
        [...activations.values()].some((stack) => stack.length)
    )
        throw new Error("Unclosed sequence interval");
    for (const actor of nodes)
        model.guides.unshift({
            actor: actor.id,
            activation: false,
            points: [
                endpoint(actor.id, -actor.size[1] / 2),
                endpoint(actor.id, y),
            ],
        });
    // Outer blocks include the complete volumes of nested blocks, so their
    // boundaries remain distinct when the camera is rotated.
    for (const group of [...model.groups].reverse()) {
        const points = model.edges
            .filter((edge) => group.edgeIds.includes(edge.id))
            .flatMap((edge) => edge.points);
        const nested = model.groups.filter(
            (child) => child.parentId === group.id,
        );
        const notes = model.notes.filter((note) => note.parentId === group.id);
        const members = [...nested, ...notes];
        group.members = members.map((child) => child.id);
        for (const child of members)
            for (const sign of [-1, 1])
                points.push(
                    child.position.map(
                        (value, axis) => value + (sign * child.size[axis]) / 2,
                    ),
                );
        if (!points.length)
            throw new Error("Empty sequence blocks are not supported");
        const low = [
            Math.min(...points.map((point) => point[0])) - 35,
            group.bottom,
            Math.min(...points.map((point) => point[2])) - 25,
        ];
        const high = [
            Math.max(...points.map((point) => point[0])) + 35,
            group.top,
            Math.max(...points.map((point) => point[2])) + 25,
        ];
        group.position = low.map((value, axis) => (value + high[axis]) / 2);
        group.size = low.map((value, axis) => high[axis] - value);
        delete group.top;
        delete group.bottom;
    }
    return model;
}
