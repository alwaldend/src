import * as THREE from "three";
import { textLayout } from "diagram/text";

export function createLabels({
    root,
    camera,
    surface,
    fontFamily,
    palette,
    occluders,
}) {
    const labels = [];
    const vector = (point) => new THREE.Vector3(...point);
    function label(
        parent,
        lines,
        position,
        {
            size = 16,
            background = false,
            color = palette.textColor,
            border = false,
            kind = "annotation",
            heading = false,
            compartments = [],
            layout: preparedLayout,
        } = {},
    ) {
        if (!lines.length) return;
        const canvas = document.createElement("canvas");
        const context = canvas.getContext("2d");
        const scale = 3;
        const layout =
            preparedLayout ??
            textLayout(
                lines,
                (text, bold) => {
                    context.font = `${bold ? "bold " : ""}${size}px ${fontFamily}`;
                    return context.measureText(text).width;
                },
                { size, heading, compartments },
            );
        const { width, height, rows, rules } = layout;
        canvas.width = Math.ceil(width * scale);
        canvas.height = Math.ceil(height * scale);
        context.scale(scale, scale);
        context.font = `${size}px ${fontFamily}`;
        if (background) {
            context.fillStyle = palette.background;
            context.beginPath();
            context.roundRect(0.5, 0.5, width - 1, height - 1, 4);
            context.fill();
            if (border) {
                context.strokeStyle = palette.clusterBorder;
                context.lineWidth = 1;
                context.stroke();
            }
        }
        context.fillStyle = color;
        context.textAlign = heading ? "left" : "center";
        context.textBaseline = "middle";
        for (const y of rules) {
            context.strokeStyle = palette.clusterBorder;
            context.lineWidth = 1;
            context.beginPath();
            context.moveTo(4, y);
            context.lineTo(width - 4, y);
            context.stroke();
        }
        for (const row of rows) {
            context.font = `${row.bold ? "bold " : ""}${layout.size}px ${fontFamily}`;
            context.fillText(row.text, layout.heading ? 10 : width / 2, row.y);
        }
        const texture = new THREE.CanvasTexture(canvas);
        texture.colorSpace = THREE.SRGBColorSpace;
        // Billboards must remain readable when their anchor face turns away from
        // the camera; otherwise the solid node slices through the label texture.
        const sprite = new THREE.Sprite(
            new THREE.SpriteMaterial({
                map: texture,
                transparent: true,
                depthWrite: false,
                depthTest: false,
            }),
        );
        sprite.renderOrder = 10;
        sprite.position.copy(vector(position));
        sprite.scale.set(width, height, 1);
        parent.add(sprite);
        const leader = new THREE.Line(
            new THREE.BufferGeometry().setFromPoints([
                new THREE.Vector3(),
                new THREE.Vector3(),
            ]),
            new THREE.LineBasicMaterial({
                color: palette.actorLineColor,
                transparent: true,
                depthTest: false,
            }),
        );
        leader.visible = false;
        parent.add(leader);
        labels.push({
            sprite,
            leader,
            anchor: vector(position),
            width,
            height,
            size,
            kind,
            text: lines.join("\n"),
            ownerId: parent.userData.id,
        });
    }

    function project(point) {
        const value = point.clone().project(camera);
        return [
            ((value.x + 1) * surface.clientWidth) / 2,
            ((1 - value.y) * surface.clientHeight) / 2,
        ];
    }
    function screenBounds(points) {
        const projected = points.map(project);
        return {
            left: Math.min(...projected.map((p) => p[0])),
            right: Math.max(...projected.map((p) => p[0])),
            top: Math.min(...projected.map((p) => p[1])),
            bottom: Math.max(...projected.map((p) => p[1])),
        };
    }
    function labelBounds(item) {
        const center = project(
            item.sprite.getWorldPosition(new THREE.Vector3()),
        );
        const factor =
            (surface.clientHeight / (camera.top - camera.bottom)) *
            camera.zoom;
        return {
            left: center[0] - (item.width * factor) / 2,
            right: center[0] + (item.width * factor) / 2,
            top: center[1] - (item.height * factor) / 2,
            bottom: center[1] + (item.height * factor) / 2,
        };
    }
    function placeLabels() {
        for (const item of labels) item.sprite.position.copy(item.anchor);
        root.updateWorldMatrix(true, true);
        camera.updateMatrixWorld(true);
        const occupied = occluders.map((item) => screenBounds(item.points));
        occupied.push(
            ...labels.filter((item) => item.kind === "node").map(labelBounds),
        );
        const overlaps = (a, b) =>
            a.left < b.right + 5 &&
            a.right > b.left - 5 &&
            a.top < b.bottom + 5 &&
            a.bottom > b.top - 5;
        for (const item of labels) {
            item.leader.visible = false;
            if (item.kind === "node") continue;
            const original = labelBounds(item);
            let placed = original;
            let offset = [0, 0];
            // Prefer the original anchor, then nearby positions in stable rings.
            const candidates = [[0, 0]];
            for (let distance = 12; distance <= 240; distance += 12)
                for (const [x, y] of [
                    [0, -1],
                    [0, 1],
                    [-1, 0],
                    [1, 0],
                    [-1, -1],
                    [1, -1],
                    [-1, 1],
                    [1, 1],
                ])
                    candidates.push([x * distance, y * distance]);
            for (const [dx, dy] of candidates) {
                const candidate = {
                    left: original.left + dx,
                    right: original.right + dx,
                    top: original.top + dy,
                    bottom: original.bottom + dy,
                };
                if (
                    candidate.left < 6 ||
                    candidate.top < 6 ||
                    candidate.right > surface.clientWidth - 6 ||
                    candidate.bottom > surface.clientHeight - 6
                )
                    continue;
                if (occupied.some((other) => overlaps(candidate, other)))
                    continue;
                placed = candidate;
                offset = [dx, dy];
                break;
            }
            occupied.push(placed);
            if (!offset[0] && !offset[1]) continue;
            const world = item.sprite.getWorldPosition(new THREE.Vector3());
            const ndc = world.clone().project(camera);
            ndc.x += (2 * offset[0]) / surface.clientWidth;
            ndc.y -= (2 * offset[1]) / surface.clientHeight;
            const moved = item.sprite.parent.worldToLocal(
                ndc.unproject(camera),
            );
            item.sprite.position.copy(moved);
            const attribute = item.leader.geometry.attributes.position;
            attribute.setXYZ(0, ...item.anchor.toArray());
            attribute.setXYZ(1, ...moved.toArray());
            attribute.needsUpdate = true;
            item.leader.geometry.computeBoundingSphere();
            item.leader.visible = Math.hypot(...offset) > 14;
        }
        root.updateWorldMatrix(true, true);
    }

    return {
        add: label,
        place: placeLabels,
        cameraPoints() {
            const points = [];
            // Labels face the camera, so include their anchored rectangles in camera
            // space. Geometry alone cannot bound a long or multiline group title.
            for (const item of labels) {
                const center = item.sprite.parent
                    .localToWorld(item.anchor.clone())
                    .applyMatrix4(camera.matrixWorldInverse);
                for (const x of [-1, 1])
                    for (const y of [-1, 1])
                        points.push(
                            center
                                .clone()
                                .add(
                                    new THREE.Vector3(
                                        (x * item.width) / 2,
                                        (y * item.height) / 2,
                                        0,
                                    ),
                                ),
                        );
            }

            return points;
        },
        inspect(geometryPoints) {
            return {
                sceneCorners: geometryPoints.map(project),
                nodeBounds: occluders.map((item) => ({
                    id: item.id,
                    bounds: screenBounds(item.points),
                })),
                labels: labels.map((item) => ({
                    text: item.text,
                    kind: item.kind,
                    ownerId: item.ownerId,
                    bounds: labelBounds(item),
                    fontPixels:
                        ((item.size * surface.clientHeight) /
                            (camera.top - camera.bottom)) *
                        camera.zoom,
                })),
            };
        },
    };
}
