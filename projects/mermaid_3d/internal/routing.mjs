const distance = (a, b) =>
    Math.hypot(...a.map((value, axis) => b[axis] - value));
const interpolate = (a, b, t) =>
    a.map((value, axis) => value + (b[axis] - value) * t);

// Sample the curve once at build time. The viewer, bounds, marker tangents,
// and collision verification all consume this same path without spline overshoot.
function bezier(controls) {
    const length = controls
        .slice(1)
        .reduce(
            (sum, point, index) => sum + distance(controls[index], point),
            0,
        );
    const steps = Math.max(24, Math.ceil(length / 3));
    return Array.from({ length: steps + 1 }, (_, index) => {
        let level = controls;
        while (level.length > 1)
            level = level
                .slice(1)
                .map((point, i) =>
                    interpolate(level[i], point, index / steps),
                );
        return level[0];
    });
}

export function roundBends(points, radius) {
    const route = points.filter(
        (point, i) => !i || distance(point, points[i - 1]) > 0.001,
    );
    const result = [route[0]];
    for (let i = 1; i + 1 < route.length; i += 1) {
        const previous = route[i - 1],
            corner = route[i],
            next = route[i + 1];
        const incoming = distance(previous, corner),
            outgoing = distance(corner, next);
        const cut = Math.min(radius, incoming * 0.45, outgoing * 0.45);
        const entry = interpolate(corner, previous, cut / incoming);
        const exit = interpolate(corner, next, cut / outgoing);
        result.push(...bezier([entry, corner, exit]));
    }
    result.push(route.at(-1));
    return result.filter(
        (point, i) => !i || distance(point, result[i - 1]) > 0.001,
    );
}

// Routes leave a node through its side, avoiding front-face text. Direct
// spatial segments are preferred; a depth corridor is needed only when an
// unrelated solid obstructs the connection.
export function routeGraph(model) {
    const all = new Map(
        [...model.nodes, ...model.groups].map((item) => [item.id, item]),
    );
    const vertical = ["TB", "TD", "BT"].includes(model.direction);
    const primary = vertical ? 1 : 0;
    const attachments = new Map();
    const endpoints = new Map();
    for (const edge of model.edges) {
        const a = all.get(edge.start);
        const b = all.get(edge.end);
        if (!a || !b) throw new Error("Missing connection endpoint");
        const axis =
            Math.abs(b.position[primary] - a.position[primary]) >
            (a.size[primary] + b.size[primary]) / 2 + 20
                ? primary
                : 1 - primary;
        const sign = Math.sign(b.position[axis] - a.position[axis]) || 1;
        const pair = [];
        for (const [node, other, side] of [
            [a, b, sign],
            [b, a, a === b ? sign : -sign],
        ]) {
            const key = node.id + ":" + axis + ":" + side;
            if (!attachments.has(key)) attachments.set(key, []);
            const attachment = { node, other, axis, side, point: null };
            attachments.get(key).push(attachment);
            pair.push(attachment);
        }
        endpoints.set(edge.id, pair);
    }
    for (const list of attachments.values()) {
        const cross = 1 - list[0].axis;
        list.sort((a, b) => a.other.position[cross] - b.other.position[cross]);
        list.forEach((item, index) => {
            const { node, axis, side } = item;
            const point = [...node.position];
            point[axis] += side * (node.size[axis] / 2 + 2);
            point[cross] +=
                ((index + 1) / (list.length + 1) - 0.5) *
                node.size[cross] *
                0.65;
            item.point = point;
        });
    }
    function intersects(a, b, node) {
        let enter = 0,
            leave = 1;
        for (let axis = 0; axis < 3; axis += 1) {
            const low = node.position[axis] - node.size[axis] / 2 - 6;
            const high = node.position[axis] + node.size[axis] / 2 + 6;
            const delta = b[axis] - a[axis];
            if (Math.abs(delta) < 1e-8) {
                if (a[axis] < low || a[axis] > high) return false;
            } else {
                const times = [
                    (low - a[axis]) / delta,
                    (high - a[axis]) / delta,
                ].sort((x, y) => x - y);
                enter = Math.max(enter, times[0]);
                leave = Math.min(leave, times[1]);
                if (enter > leave) return false;
            }
        }
        return enter <= leave;
    }
    const [vx, vy, vz] = model.viewDirection;
    const horizontal = Math.hypot(vx, vz);
    const magnitude = Math.hypot(vx, vy, vz);
    const project = (point) => [
        (vz * point[0] - vx * point[2]) / horizontal,
        (-vx * vy * point[0] +
            horizontal * horizontal * point[1] -
            vy * vz * point[2]) /
            (horizontal * magnitude),
        0,
    ];
    const projectedNodes = new Map(
        model.nodes.map((node) => {
            const points = Array.from({ length: 8 }, (_, i) =>
                project(
                    node.position.map(
                        (value, axis) =>
                            value +
                            (((i >> axis) & 1 ? 1 : -1) * node.size[axis]) / 2,
                    ),
                ),
            );
            const low = [0, 1].map((axis) =>
                Math.min(...points.map((point) => point[axis])),
            );
            const high = [0, 1].map((axis) =>
                Math.max(...points.map((point) => point[axis])),
            );
            return [
                node.id,
                {
                    position: [
                        (low[0] + high[0]) / 2,
                        (low[1] + high[1]) / 2,
                        0,
                    ],
                    size: [high[0] - low[0], high[1] - low[1], 1],
                },
            ];
        }),
    );
    const clear = (points, edge) =>
        model.nodes.every(
            (node) =>
                node.id === edge.start ||
                node.id === edge.end ||
                points
                    .slice(1)
                    .every(
                        (point, index) =>
                            !intersects(points[index], point, node) &&
                            !intersects(
                                project(points[index]),
                                project(point),
                                projectedNodes.get(node.id),
                            ),
                    ),
        );
    let lane = 0;
    const front = Math.max(
        ...model.nodes.map((node) => node.position[2] + node.size[2] / 2),
    );
    const offset = (point, axis, amount) =>
        point.map((value, index) => value + (index === axis ? amount : 0));
    const along = (a, b, distance) => {
        const length = Math.hypot(...a.map((value, i) => b[i] - value));
        return a.map(
            (value, i) =>
                value + (b[i] - value) * Math.min(0.35, distance / length),
        );
    };
    for (const edge of model.edges) {
        const [a, b] = endpoints.get(edge.id);
        const start = a.point,
            end = b.point;
        let points =
            edge.start === edge.end
                ? [
                      start,
                      offset(start, a.axis, a.side * 80),
                      offset(end, b.axis, b.side * 80),
                      end,
                  ]
                : [start, end];
        if (!clear(points, edge)) {
            const exit = offset(start, a.axis, a.side * 28);
            const entry = offset(end, b.axis, b.side * 28);
            const z = front + 32 + lane++ * 12;
            points = [
                start,
                exit,
                [exit[0], exit[1], z],
                [entry[0], entry[1], z],
                entry,
                end,
            ];
            if (!clear(points, edge)) {
                const cross = 1 - a.axis;
                const candidates = model.nodes
                    .flatMap((node) =>
                        [-1, 1].map((sign) => {
                            const coordinate =
                                node.position[cross] +
                                sign * (node.size[cross] / 2 + 42);
                            const turnA = [...exit],
                                turnB = [...entry];
                            turnA[cross] = coordinate;
                            turnB[cross] = coordinate;
                            turnA[2] = z;
                            turnB[2] = z;
                            return [start, exit, turnA, turnB, entry, end];
                        }),
                    )
                    .filter((candidate) => clear(candidate, edge));
                const length = (candidate) =>
                    candidate
                        .slice(1)
                        .reduce(
                            (sum, point, i) =>
                                sum +
                                Math.hypot(
                                    ...point.map(
                                        (value, axis) =>
                                            value - candidate[i][axis],
                                    ),
                                ),
                            0,
                        );
                candidates.sort((left, right) => length(left) - length(right));
                if (!candidates.length)
                    throw new Error("No clear spatial route for " + edge.id);
                points = candidates[0];
            }
        }
        edge.points = points;
        let longest = 0,
            middle = 1;
        for (let i = 1; i < points.length; i += 1) {
            const distance = Math.hypot(
                points[i][0] - points[i - 1][0],
                points[i][1] - points[i - 1][1],
            );
            if (distance > longest) {
                longest = distance;
                middle = i;
            }
        }
        const left = points[middle - 1],
            right = points[middle];
        edge.labelPosition = left.map(
            (value, i) => (value + right[i]) / 2 + (i === 1 ? 18 : 0),
        );
        if (edge.start === edge.end)
            edge.labelPosition = [
                (Math.min(...points.map((point) => point[0])) +
                    Math.max(...points.map((point) => point[0]))) /
                    2,
                Math.max(...points.map((point) => point[1])) + 28,
                start[2],
            ];
        edge.startLabelPosition = offset(along(start, points[1], 42), 1, -26);
        edge.endLabelPosition = offset(along(end, points.at(-2), 42), 1, -26);
        // Keep the established route and ports, reducing the bend radius near
        // obstacles instead of letting an unconstrained spline cut through them.
        for (const scale of [1, 0.5, 0.25, 0.125, 0.0625]) {
            const handle = Math.min(distance(start, end) * 0.4, 180) * scale;
            const curved =
                points.length === 2
                    ? bezier([
                          start,
                          offset(start, a.axis, a.side * handle),
                          offset(end, b.axis, b.side * handle),
                          end,
                      ])
                    : roundBends(points, 64 * scale);
            if (clear(curved, edge)) {
                edge.points = curved;
                break;
            }
        }
    }
}
