import type {
    BoxElement,
    DiagramEdge,
    DiagramModel,
    DiagramNode,
    ElementKind,
    Palette,
    Point3,
    SceneInspection,
} from "../types.js";
import * as THREE from "three";
import { mergeGeometries } from "three/addons/utils/BufferGeometryUtils.js";
import { createLabels } from "diagram/labels";

type SceneMaterial = THREE.Material & { map?: THREE.Texture | null };
interface RenderObject extends THREE.Object3D {
    geometry?: THREE.BufferGeometry;
    material?: SceneMaterial | SceneMaterial[];
    isMesh?: boolean;
    isSprite?: boolean;
}
interface DiagramOptions {
    model: DiagramModel;
    palette: Palette;
    fontFamily: string;
    camera: THREE.OrthographicCamera;
    surface: HTMLElement;
}

class PathSection extends THREE.Curve<THREE.Vector3> {
    constructor(
        private readonly evaluate: (
            t: number,
            target?: THREE.Vector3,
        ) => THREE.Vector3,
    ) {
        super();
    }
    override getPoint(t: number, target?: THREE.Vector3): THREE.Vector3 {
        return this.evaluate(t, target);
    }
}

export const elementKey = (kind: ElementKind, id: string) => `${kind}:${id}`;
export const cardinalities: Record<string, string> = {
    only_one: "1",
    zero_or_one: "0..1",
    zero_or_more: "0..*",
    one_or_more: "1..*",
};

const vector = (point: Point3) =>
    new THREE.Vector3(point[0], point[1], point[2]);
function materialsOf(object: RenderObject): SceneMaterial[] {
    return Array.isArray(object.material)
        ? object.material
        : object.material
          ? [object.material]
          : [];
}
function corners(item: Pick<BoxElement, "position" | "size">) {
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
function dispose(object: THREE.Object3D) {
    object.traverse((child: RenderObject) => {
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
}: DiagramOptions) {
    const root = new THREE.Group();
    const objects = new Map<string, THREE.Group>();
    const pickable: THREE.Object3D[] = [];
    const frameStrokes: {
        material: THREE.MeshBasicMaterial;
        offset: THREE.Vector3;
    }[] = [];
    const points = [
        ...[...model.nodes, ...model.groups].flatMap(corners),
        ...model.edges.flatMap((item) => item.points.map(vector)),
    ];
    const labelLayer = createLabels({
        root,
        camera,
        surface,
        fontFamily,
        palette,
        occluders: model.nodes.map((item) => ({
            id: item.id,
            points: corners(item),
        })),
    });
    const label = labelLayer.add;
    function createElement(kind: ElementKind, id: string) {
        const object = new THREE.Group();
        object.userData = { id, kind, key: elementKey(kind, id) };
        return object;
    }
    function addElement(
        object: THREE.Group,
        pickingSurfaces?: THREE.Object3D[],
    ) {
        if (pickingSurfaces) pickable.push(...pickingSurfaces);
        else
            object.traverse((child: RenderObject) => {
                child.userData = object.userData;
                if (child.isMesh) pickable.push(child);
            });
        root.add(object);
        objects.set(object.userData.key, object);
    }
    function meshMaterial(color: THREE.ColorRepresentation) {
        return new THREE.MeshStandardMaterial({
            color,
            roughness: 0.8,
            metalness: 0,
            transparent: true,
        });
    }
    function lineMaterial(color: THREE.ColorRepresentation) {
        return new THREE.MeshBasicMaterial({ color, transparent: true });
    }
    function pathTube(
        parent: THREE.Object3D,
        points: Point3[],
        color: THREE.ColorRepresentation,
        radius = 1.15,
        dashed = false,
    ) {
        const path = new THREE.CurvePath<THREE.Vector3>();
        for (let i = 1; i < points.length; i += 1)
            path.add(
                new THREE.LineCurve3(vector(points[i - 1]), vector(points[i])),
            );
        const length = path.getLength();
        if (length < 0.001) return null;
        const parts: THREE.BufferGeometry[] = [];
        // Arc-length intervals retain dash spacing through bends. Merge the
        // disjoint sections so each path requires just one mesh and draw call.
        for (let start = 0; start < length; start += dashed ? 14 : length) {
            const span = dashed ? Math.min(8, length - start) : length;
            const section = new PathSection((t, target) =>
                path.getPoint((start + span * t) / length, target),
            );
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
        const mesh = new THREE.Mesh(geometry!, lineMaterial(color));
        parent.add(mesh);
        return mesh;
    }
    function boxOutline(
        parent: THREE.Object3D,
        size: Point3,
        position: Point3,
        color: THREE.ColorRepresentation,
        radius = 0.8,
    ) {
        const corners: Point3[] = [];
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
                        material: mesh!.material,
                        offset: vector(corners[i])
                            .add(vector(corners[j]))
                            .multiplyScalar(0.5)
                            .sub(vector(position)),
                    });
                }
            }
    }
    function extrude(shape: THREE.Shape, depth: number) {
        const geometry = new THREE.ExtrudeGeometry(shape, {
            depth,
            bevelEnabled: false,
            curveSegments: 12,
        });
        geometry.translate(0, 0, -depth / 2);
        return geometry;
    }
    function roundedBox(
        width: number,
        height: number,
        depth: number,
        radius: number,
    ) {
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
    function drawNode(node: DiagramNode) {
        const object = createElement("node", node.id);
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
        addElement(object, [mesh]);
    }
    function marker(
        parent: THREE.Object3D,
        type: string,
        tip: Point3,
        previous: Point3,
    ) {
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
    function drawEdge(edge: DiagramEdge) {
        const object = createElement("edge", edge.id);
        pathTube(object, edge.points, palette.lineColor, 1.15, edge.dashed);
        marker(object, edge.startMarker, edge.points[0], edge.points[1]);
        marker(
            object,
            edge.endMarker,
            edge.points.at(-1)!,
            edge.points.at(-2)!,
        );
        if (edge.label)
            label(object, edge.label.split("\n"), edge.labelPosition!, {
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
        ] as [string | undefined, Point3 | undefined][])
            if (text)
                label(object, [text], point!, { size: 14, background: true });
        addElement(object);
    }
    for (const group of model.groups) {
        const object = createElement("group", group.id);
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
                layout: group.textLayout,
            },
        );
        addElement(object);
    }
    for (const node of model.nodes) drawNode(node);
    for (const edge of model.edges) drawEdge(edge);

    function geometryCorners() {
        const bounds = new THREE.Box3();
        root.updateWorldMatrix(true, true);
        root.traverse((item: RenderObject) => {
            if (!item.geometry || item.isSprite || !item.visible) return;
            item.geometry.computeBoundingBox();
            bounds.union(
                item.geometry
                    .boundingBox!.clone()
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
        highlight(related: Set<string> | null) {
            for (const [key, object] of objects)
                object.traverse((child: RenderObject) => {
                    for (const material of materialsOf(child))
                        material.opacity =
                            !related || related.has(key) ? 1 : 0.14;
                });
        },
        pick(raycaster: THREE.Raycaster): string {
            return (
                raycaster.intersectObjects(pickable, false)[0]?.object.userData
                    .key ?? ""
            );
        },
        inspect(): SceneInspection {
            let meshCount = 0;
            root.traverse((item: RenderObject) => {
                if (item.isMesh) meshCount += 1;
            });
            return {
                ...labelLayer.inspect(points),
                meshCount,
                elements: [...objects].map(([key, object]) => {
                    const opacities = new Set<number>();
                    object.traverse((child: RenderObject) => {
                        for (const material of materialsOf(child))
                            opacities.add(material.opacity);
                    });
                    return {
                        key,
                        id: object.userData.id,
                        kind: object.userData.kind as ElementKind,
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
