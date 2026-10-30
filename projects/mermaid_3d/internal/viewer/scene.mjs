import * as THREE from "three";
import { mergeGeometries } from "three/geometry";
import { createLabels } from "diagram/labels";

export const elementKey = (kind, id) => `${kind}:${id}`;
export const cardinalities = {
    only_one: "1",
    zero_or_one: "0..1",
    zero_or_more: "0..*",
    one_or_more: "1..*",
};

const vector = (point) => new THREE.Vector3(...point);
function materialsOf(object) {
    return Array.isArray(object.material)
        ? object.material
        : object.material
          ? [object.material]
          : [];
}
function corners(item) {
    return Array.from({ length: 8 }, (_, i) =>
        vector(
            item.size.map(
                (value, axis) =>
                    item.position[axis] +
                    (((i >> axis) & 1 ? 1 : -1) * value) / 2,
            ),
        ),
    );
}
function dispose(object) {
    object.traverse((child) => {
        child.geometry?.dispose();
        const materials = materialsOf(child);
        for (const material of materials) {
            material.map?.dispose();
            material.dispose();
        }
    });
}

export function createDiagram({
    model,
    palette,
    fontFamily,
    camera,
    surface,
}) {
    const root = new THREE.Group();
    const objects = new Map();
    const pickable = [];
    const frameStrokes = [];
    const points = [
        ...[...model.nodes, ...model.groups, ...model.notes].flatMap(corners),
        ...[...model.edges, ...model.guides].flatMap((item) =>
            item.points.map(vector),
        ),
    ];
    const labelLayer = createLabels({
        root,
        camera,
        surface,
        fontFamily,
        palette,
        occluders: [...model.nodes, ...model.notes].map((item) => ({
            id: item.id,
            points: corners(item),
        })),
    });
    const label = labelLayer.add;
    function meshMaterial(color) {
        return new THREE.MeshStandardMaterial({
            color,
            roughness: 0.8,
            metalness: 0,
            transparent: true,
        });
    }
    function lineMaterial(color) {
        return new THREE.MeshBasicMaterial({ color, transparent: true });
    }
    function pathTube(parent, points, color, radius = 1.15, dashed = false) {
        const path = new THREE.CurvePath();
        for (let i = 1; i < points.length; i += 1)
            path.add(
                new THREE.LineCurve3(vector(points[i - 1]), vector(points[i])),
            );
        const length = path.getLength();
        if (length < 0.001) return null;
        const parts = [];
        // Arc-length intervals retain dash spacing through bends. Merge the
        // disjoint sections so each path requires just one mesh and draw call.
        for (let start = 0; start < length; start += dashed ? 14 : length) {
            const span = dashed ? Math.min(8, length - start) : length;
            const section = new THREE.Curve();
            section.getPoint = (t, target) =>
                path.getPoint((start + span * t) / length, target);
            parts.push(
                new THREE.TubeGeometry(
                    section,
                    Math.max(2, Math.ceil(span / 2)),
                    radius,
                    8,
                    false,
                ),
            );
        }
        const geometry =
            parts.length === 1 ? parts[0] : mergeGeometries(parts);
        if (parts.length > 1) for (const part of parts) part.dispose();
        const mesh = new THREE.Mesh(geometry, lineMaterial(color));
        parent.add(mesh);
        return mesh;
    }
    function boxOutline(parent, size, position, color, radius = 0.8) {
        const corners = [];
        for (let i = 0; i < 8; i += 1)
            corners.push(
                size.map(
                    (value, axis) =>
                        position[axis] +
                        (((i >> axis) & 1 ? 1 : -1) * value) / 2,
                ),
            );
        for (let i = 0; i < 8; i += 1)
            for (let axis = 0; axis < 3; axis += 1) {
                const j = i ^ (1 << axis);
                if (i < j) {
                    const mesh = pathTube(
                        parent,
                        [corners[i], corners[j]],
                        color,
                        radius,
                    );
                    frameStrokes.push({
                        material: mesh.material,
                        offset: vector(corners[i])
                            .add(vector(corners[j]))
                            .multiplyScalar(0.5)
                            .sub(vector(position)),
                    });
                }
            }
    }
    function extrude(shape, depth) {
        const geometry = new THREE.ExtrudeGeometry(shape, {
            depth,
            bevelEnabled: false,
            curveSegments: 12,
        });
        geometry.translate(0, 0, -depth / 2);
        return geometry;
    }
    function roundedBox(width, height, depth, radius) {
        const x = -width / 2;
        const y = -height / 2;
        const shape = new THREE.Shape();
        shape.moveTo(x + radius, y);
        shape.lineTo(x + width - radius, y);
        shape.quadraticCurveTo(x + width, y, x + width, y + radius);
        shape.lineTo(x + width, y + height - radius);
        shape.quadraticCurveTo(
            x + width,
            y + height,
            x + width - radius,
            y + height,
        );
        shape.lineTo(x + radius, y + height);
        shape.quadraticCurveTo(x, y + height, x, y + height - radius);
        shape.lineTo(x, y + radius);
        shape.quadraticCurveTo(x, y, x + radius, y);
        return extrude(shape, depth);
    }
    function drawNode(node, kind = "node") {
        const object = new THREE.Group();
        object.userData = {
            id: node.id,
            kind,
            key: elementKey(kind, node.id),
        };
        const [w, h, d] = node.size;
        let geometry;
        if (node.shape === "diamond") {
            const shape = new THREE.Shape();
            shape.moveTo(0, h / 2);
            shape.lineTo(w / 2, 0);
            shape.lineTo(0, -h / 2);
            shape.lineTo(-w / 2, 0);
            shape.closePath();
            geometry = extrude(shape, d);
        } else if (/circle|stateStart|stateEnd/i.test(node.shape)) {
            geometry = new THREE.SphereGeometry(1, 32, 20);
            geometry.scale(w / 2, h / 2, d / 2);
        } else if (node.shape === "cylinder") {
            geometry = new THREE.CylinderGeometry(w / 2, w / 2, h, 40);
            geometry.scale(1, 1, d / w);
        } else
            geometry = roundedBox(
                w,
                h,
                d,
                node.shape === "stadium" ? h / 2 : 10,
            );
        const fill = /stateStart|stateEnd/.test(node.shape)
            ? palette.nodeBorder
            : palette.secondaryColor;
        const mesh = new THREE.Mesh(geometry, meshMaterial(fill));
        mesh.userData = object.userData;
        object.add(mesh);
        pickable.push(mesh);
        const outline = new THREE.LineSegments(
            new THREE.EdgesGeometry(geometry, 25),
            new THREE.LineBasicMaterial({
                color: palette.nodeBorder,
                transparent: true,
            }),
        );
        object.add(outline);
        if (node.shape === "stateEnd") {
            const ring = new THREE.Mesh(
                new THREE.TorusGeometry(w * 0.38, 1.3, 8, 36),
                lineMaterial(palette.background),
            );
            ring.position.z = d / 2 + 1;
            object.add(ring);
        }
        label(object, node.lines, [0, 0, d / 2 + 5], {
            color: palette.nodeTextColor,
            kind: "node",
            heading: node.lines.length > 1,
            layout: node.textLayout,
        });
        object.position.copy(vector(node.position));
        root.add(object);
        objects.set(object.userData.key, object);
    }
    function marker(parent, type, tip, previous) {
        if (!type || type === "none" || cardinalities[type]) return;
        const direction = vector(tip).sub(vector(previous)).normalize();
        const diamond = ["composition", "aggregation"].includes(type);
        const geometry = diamond
            ? new THREE.OctahedronGeometry(6)
            : new THREE.ConeGeometry(5, 12, 3);
        const material = lineMaterial(
            ["extension", "aggregation"].includes(type)
                ? palette.background
                : palette.lineColor,
        );
        const mesh = new THREE.Mesh(geometry, material);
        mesh.position.copy(vector(tip).addScaledVector(direction, -6));
        mesh.quaternion.setFromUnitVectors(
            new THREE.Vector3(0, 1, 0),
            direction,
        );
        parent.add(mesh);
        const outline = new THREE.LineSegments(
            new THREE.EdgesGeometry(geometry),
            new THREE.LineBasicMaterial({
                color: palette.lineColor,
                transparent: true,
            }),
        );
        outline.position.copy(mesh.position);
        outline.quaternion.copy(mesh.quaternion);
        parent.add(outline);
    }
    function drawEdge(edge) {
        const object = new THREE.Group();
        object.userData = {
            id: edge.id,
            kind: "edge",
            key: elementKey("edge", edge.id),
        };
        pathTube(object, edge.points, palette.lineColor, 1.15, edge.dashed);
        marker(object, edge.startMarker, edge.points[0], edge.points[1]);
        marker(object, edge.endMarker, edge.points.at(-1), edge.points.at(-2));
        if (edge.label)
            label(object, edge.label.split("\n"), edge.labelPosition, {
                size: 14,
                background: true,
            });
        for (const [text, point] of [
            [
                edge.startLabel || cardinalities[edge.startMarker],
                edge.startLabelPosition,
            ],
            [
                edge.endLabel || cardinalities[edge.endMarker],
                edge.endLabelPosition,
            ],
        ])
            if (text)
                label(object, [text], point, { size: 14, background: true });
        object.traverse((child) => {
            child.userData = object.userData;
            if (child.isMesh) pickable.push(child);
        });
        root.add(object);
        objects.set(object.userData.key, object);
    }
    for (const group of model.groups) {
        const object = new THREE.Group();
        object.userData = {
            id: group.id,
            kind: "group",
            key: elementKey("group", group.id),
        };
        boxOutline(
            object,
            group.size,
            group.position,
            palette.actorLineColor,
            1,
        );
        label(
            object,
            group.label.split("\n"),
            [
                group.position[0],
                group.position[1] + group.size[1] / 2,
                group.position[2] + group.size[2] / 2 + 5,
            ],
            {
                size: 14,
                background: true,
                border: true,
                layout: group.textLayout,
            },
        );
        for (const section of group.sections ?? []) {
            const z = group.position[2] + group.size[2] / 2;
            pathTube(
                object,
                [
                    [group.position[0] - group.size[0] / 2, section.y, z],
                    [group.position[0] + group.size[0] / 2, section.y, z],
                ],
                palette.clusterBorder,
                0.8,
                true,
            );
            label(
                object,
                section.label.split("\n"),
                [group.position[0], section.y, z + 5],
                { size: 14, background: true, border: true },
            );
        }
        object.traverse((child) => {
            child.userData = object.userData;
            if (child.isMesh) pickable.push(child);
        });
        root.add(object);
        objects.set(object.userData.key, object);
    }
    for (const guide of model.guides)
        pathTube(
            root,
            guide.points,
            palette.actorLineColor,
            guide.activation ? 4 : 0.7,
            !guide.activation,
        );
    for (const node of model.nodes) drawNode(node);
    for (const edge of model.edges) drawEdge(edge);
    for (const note of model.notes) drawNode(note, "note");

    function geometryCorners() {
        const bounds = new THREE.Box3();
        root.updateWorldMatrix(true, true);
        root.traverse((item) => {
            if (!item.geometry || item.isSprite || !item.visible) return;
            item.geometry.computeBoundingBox();
            bounds.union(
                item.geometry.boundingBox
                    .clone()
                    .applyMatrix4(item.matrixWorld),
            );
        });
        return bounds.isEmpty()
            ? points
            : corners({
                  position: bounds.getCenter(new THREE.Vector3()).toArray(),
                  size: bounds.getSize(new THREE.Vector3()).toArray(),
              });
    }

    return {
        root,
        points,
        update() {
            labelLayer.place();
            const towardCamera = camera
                .getWorldDirection(new THREE.Vector3())
                .negate();
            const frameColor = new THREE.Color(palette.actorLineColor);
            const background = new THREE.Color(palette.background);
            for (const stroke of frameStrokes)
                stroke.material.color
                    .copy(frameColor)
                    .lerp(
                        background,
                        stroke.offset.dot(towardCamera) < 0 ? 0.6 : 0,
                    );
        },
        framePoints(includeGeometry = false) {
            return [
                ...(includeGeometry ? geometryCorners() : points).map(
                    (point) =>
                        point.clone().applyMatrix4(camera.matrixWorldInverse),
                ),
                ...labelLayer.cameraPoints(),
            ];
        },
        highlight(related) {
            for (const [key, object] of objects)
                object.traverse((child) => {
                    for (const material of materialsOf(child))
                        material.opacity =
                            !related || related.has(key) ? 1 : 0.14;
                });
        },
        pick(raycaster) {
            return (
                raycaster.intersectObjects(pickable, false)[0]?.object.userData
                    .key ?? ""
            );
        },
        inspect(includeGeometry = false) {
            let meshCount = 0;
            root.traverse((item) => {
                if (item.isMesh) meshCount += 1;
            });
            return {
                ...labelLayer.inspect(
                    includeGeometry ? geometryCorners() : points,
                ),
                meshCount,
                elements: [...objects].map(([key, object]) => {
                    const opacities = new Set();
                    object.traverse((child) => {
                        for (const material of materialsOf(child))
                            opacities.add(material.opacity);
                    });
                    return {
                        key,
                        id: object.userData.id,
                        kind: object.userData.kind,
                        opacities: [...opacities],
                    };
                }),
            };
        },
        dispose() {
            dispose(root);
        },
    };
}
