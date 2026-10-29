import * as THREE from "three";
import { OrbitControls } from "three/controls";
import { createDiagram, elementKey, cardinalities } from "diagram/scene";

const payload = JSON.parse(
    document.getElementById("diagram-data").textContent,
);
const { model, appearances, fontFamily, fontName } = payload;
await document.fonts.load(`16px ${JSON.stringify(fontName)}`);
await document.fonts.ready;
const viewport = document.getElementById("viewport");
const surface = document.getElementById("scene");
const scene = new THREE.Scene();
const renderer = new THREE.WebGLRenderer({
    antialias: true,
    alpha: false,
    preserveDrawingBuffer: true,
});
renderer.setPixelRatio(Math.min(devicePixelRatio, 2));
renderer.outputColorSpace = THREE.SRGBColorSpace;
surface.append(renderer.domElement);
const camera = new THREE.OrthographicCamera(-500, 500, 400, -400);
const controls = new OrbitControls(camera, renderer.domElement);
controls.enableDamping = false;
controls.minZoom = 0.1;
controls.maxZoom = 12;
scene.add(new THREE.AmbientLight(0xffffff, 2.4));
const light = new THREE.DirectionalLight(0xffffff, 1.7);
light.position.set(-400, 900, 1200);
scene.add(light);
let diagram;
let appearance = matchMedia("(prefers-color-scheme: dark)").matches
    ? "dark"
    : "light";
let selected = "";
let palette;
let baseCamera;
let baseTarget;
let viewMode = "readable";
let snapshotSize;
const vector = (point) => new THREE.Vector3(...point);
const selection = document.getElementById("selection");
const all = new Map();
for (const [kind, label, items] of [
    ["node", "Nodes", model.nodes],
    ["group", "Groups", model.groups],
    ["edge", "Connections", model.edges],
    ["note", "Notes", model.notes],
]) {
    if (!items.length) continue;
    const group = document.createElement("optgroup");
    group.label = label;
    for (const item of items) {
        const option = document.createElement("option");
        option.value = elementKey(kind, item.id);
        all.set(option.value, { kind, item });
        option.textContent = item.label || `${item.start} → ${item.end}`;
        group.append(option);
    }
    selection.append(group);
}
document.getElementById("source").textContent = payload.source;
document.getElementById("summary").textContent =
    `${model.nodes.length} nodes / ${model.edges.length} connections / ${model.groups.length} groups`;

function render() {
    diagram.update();
    fitCameraDepth();
    renderer.render(scene, camera);
}
function fitCameraDepth() {
    const depths = diagram.framePoints(true).map((point) => -point.z);
    let near = Math.min(...depths);
    let far = Math.max(...depths);
    const margin = Math.max(1, (far - near) * 0.05);
    if (near <= margin) {
        const retreat = margin - near + 1;
        camera.position.addScaledVector(
            camera.getWorldDirection(new THREE.Vector3()),
            -retreat,
        );
        near += retreat;
        far += retreat;
    }
    camera.near = Math.max(0.01, near - margin);
    camera.far = far + margin;
    camera.updateProjectionMatrix();
}
function rebuild() {
    palette = appearances[appearance];
    const root = document.documentElement;
    for (const [key, value] of Object.entries({
        canvas: palette.background,
        ink: palette.textColor,
        border: palette.clusterBorder,
        surface: palette.primaryColor,
        secondary: palette.secondaryColor,
        muted: palette.actorLineColor,
        "diagram-font": fontFamily,
    }))
        root.style.setProperty(`--${key}`, value);
    root.style.colorScheme = appearance;
    scene.background = new THREE.Color(palette.background);
    if (diagram) {
        diagram.dispose();
        scene.remove(diagram.root);
    }
    diagram = createDiagram({ model, palette, fontFamily, camera, surface });
    scene.add(diagram.root);
    select(selected);
}
function select(id) {
    selected = id;
    selection.value = id;
    // Resolve containment first. Adding a connection's endpoints must not
    // turn those endpoints into new seeds for the next connection.
    const members = new Set([id]);
    const connections = new Set([id]);
    const endpointKey = (id) =>
        elementKey(all.has(elementKey("node", id)) ? "node" : "group", id);
    for (const member of members) {
        const entry = all.get(member);
        if (entry?.kind !== "group") continue;
        for (const child of entry.item.members) {
            const kind = all.has(elementKey("group", child))
                ? "group"
                : model.type === "sequence"
                  ? "note"
                  : "node";
            members.add(elementKey(kind, child));
        }
        for (const edge of entry.item.edgeIds ?? [])
            connections.add(elementKey("edge", edge));
    }
    const related = new Set(members);
    for (const edge of model.edges)
        if (
            connections.has(elementKey("edge", edge.id)) ||
            members.has(endpointKey(edge.start)) ||
            members.has(endpointKey(edge.end))
        ) {
            related.add(elementKey("edge", edge.id));
            related.add(endpointKey(edge.start));
            related.add(endpointKey(edge.end));
        }
    diagram.highlight(id ? related : null);
    const item = all.get(id)?.item;
    document.getElementById("details").textContent = item
        ? [
              item.label || "Connection",
              ...(item.lines?.slice(1) ?? []),
              item.start ? `${item.start} → ${item.end}` : "",
              item.startMarker && cardinalities[item.startMarker]
                  ? `Source: ${cardinalities[item.startMarker]}`
                  : "",
              item.endMarker && cardinalities[item.endMarker]
                  ? `Target: ${cardinalities[item.endMarker]}`
                  : "",
          ]
              .filter(Boolean)
              .join("\n")
        : "Select a node, group, or connection to trace its relationships.";
}
function resize() {
    if (!baseCamera) return;
    camera.updateMatrixWorld(true);
    const points = diagram.framePoints();
    const low = new THREE.Vector2(
        Math.min(...points.map((p) => p.x)),
        Math.min(...points.map((p) => p.y)),
    );
    const high = new THREE.Vector2(
        Math.max(...points.map((p) => p.x)),
        Math.max(...points.map((p) => p.y)),
    );
    const span = high.clone().sub(low);
    const middle = high.clone().add(low).multiplyScalar(0.5);
    const shift = new THREE.Vector3(middle.x, middle.y, 0).applyQuaternion(
        camera.quaternion,
    );
    camera.position.add(shift);
    controls.target.add(shift);
    const availableWidth = snapshotSize?.width ?? viewport.clientWidth;
    const availableHeight =
        snapshotSize?.height ??
        parseFloat(
            getComputedStyle(viewport).getPropertyValue("--view-height"),
        );
    const fitScale = Math.min(
        (availableWidth - 96) / span.x,
        (availableHeight - 96) / span.y,
    );
    const scale =
        viewMode === "fit"
            ? fitScale
            : Math.max(0.9, Math.min(1.15, fitScale));
    const width =
        viewMode === "fit"
            ? availableWidth
            : Math.max(availableWidth, Math.ceil(span.x * scale + 96));
    const height =
        snapshotSize?.height ??
        (viewMode === "fit"
            ? Math.min(availableHeight, Math.ceil(span.y * scale + 96))
            : Math.max(availableHeight, Math.ceil(span.y * scale + 96)));
    if (viewMode === "fit")
        viewport.style.setProperty("--overview-height", height + "px");
    else viewport.style.removeProperty("--overview-height");
    surface.style.width = width + "px";
    surface.style.height = height + "px";
    viewport.dataset.mode = viewMode;
    const frameHeight = height / scale;
    camera.left = -width / scale / 2;
    camera.right = width / scale / 2;
    camera.top = frameHeight / 2;
    camera.bottom = -frameHeight / 2;
    camera.updateProjectionMatrix();
    renderer.setSize(width, height);
    const touchScroll =
        viewMode === "readable" && matchMedia("(pointer: coarse)").matches;
    controls.enabled = !snapshotSize && !touchScroll;
    renderer.domElement.style.touchAction = touchScroll
        ? "pan-x pan-y"
        : "none";
    document
        .getElementById("fit")
        .setAttribute("aria-pressed", String(viewMode === "fit"));
    document.getElementById("view-status").textContent =
        viewMode === "fit"
            ? "Overview · Reset view restores readable text"
            : "Readable view · Scroll to explore the full diagram";
    render();
}
function reset(front = false, mode = "readable") {
    viewMode = mode;
    camera.position.copy(baseCamera);
    controls.target.copy(baseTarget);
    camera.zoom = 1;
    if (front) camera.position.set(baseTarget.x, baseTarget.y, baseCamera.z);
    camera.lookAt(controls.target);
    controls.update();
    resize();
    viewport.scrollTo(0, 0);
}
rebuild();
const box = new THREE.Box3().setFromPoints(diagram.points);
const center = box.getCenter(new THREE.Vector3());
const size = box.getSize(new THREE.Vector3());
const extent = Math.max(size.x, size.y, size.z, 200);
baseTarget = center;
baseCamera = center
    .clone()
    .add(vector(model.viewDirection).multiplyScalar(extent));
reset();
controls.addEventListener("change", render);
new ResizeObserver(() => {
    if (!snapshotSize) resize();
}).observe(viewport);
document.getElementById("appearance").value = appearance;
document.getElementById("appearance").addEventListener("change", (event) => {
    appearance = event.target.value;
    rebuild();
    render();
});
selection.addEventListener("change", (event) => {
    select(event.target.value);
    render();
});
document.getElementById("reset").addEventListener("click", () => reset());
document
    .getElementById("front")
    .addEventListener("click", () => reset(true, viewMode));
document
    .getElementById("fit")
    .addEventListener("click", () => reset(false, "fit"));
let pointerStart;
renderer.domElement.addEventListener("pointerdown", (event) => {
    pointerStart = [event.clientX, event.clientY];
});
renderer.domElement.addEventListener("pointerup", (event) => {
    if (
        !pointerStart ||
        Math.hypot(
            event.clientX - pointerStart[0],
            event.clientY - pointerStart[1],
        ) > 5
    )
        return;
    const rect = renderer.domElement.getBoundingClientRect();
    const raycaster = new THREE.Raycaster();
    raycaster.setFromCamera(
        new THREE.Vector2(
            ((event.clientX - rect.left) / rect.width) * 2 - 1,
            (-(event.clientY - rect.top) / rect.height) * 2 + 1,
        ),
        camera,
    );
    select(diagram.pick(raycaster));
    render();
});
window.diagramView = {
    ready: true,
    prepareSnapshot({ width, height, colorScheme, camera: view }) {
        if (
            ![width, height].every(
                (value) =>
                    Number.isInteger(value) && value >= 128 && value <= 8192,
            ) ||
            width * height > 16777216
        )
            throw new Error(
                "Snapshot dimensions must be integers from 128 to 8192, at most 16777216 pixels total",
            );
        if (
            !["light", "dark"].includes(colorScheme) ||
            !["isometric", "front"].includes(view)
        )
            throw new Error("Invalid snapshot appearance or camera");
        snapshotSize = { width, height };
        appearance = colorScheme;
        selected = "";
        renderer.setPixelRatio(1);
        rebuild();
        reset(view === "front", "fit");
        // Label placement depends on the projection. Fit the rendered extents,
        // including displaced labels, rather than geometry alone.
        const padding = 24;
        for (let attempt = 0; attempt < 16; attempt += 1) {
            const state = diagram.inspect(true);
            const points = [
                ...state.sceneCorners,
                ...state.labels.flatMap(({ bounds }) => [
                    [bounds.left, bounds.top],
                    [bounds.right, bounds.bottom],
                ]),
            ];
            const left = Math.min(...points.map((point) => point[0]));
            const right = Math.max(...points.map((point) => point[0]));
            const top = Math.min(...points.map((point) => point[1]));
            const bottom = Math.max(...points.map((point) => point[1]));
            if (
                left >= padding - 1 &&
                top >= padding - 1 &&
                right <= width - padding + 1 &&
                bottom <= height - padding + 1
            )
                return this.inspect();
            const shift = new THREE.Vector3(
                (((left + right - width) / 2) * (camera.right - camera.left)) /
                    width,
                (((height - top - bottom) / 2) *
                    (camera.top - camera.bottom)) /
                    height,
                0,
            ).applyQuaternion(camera.quaternion);
            camera.position.add(shift);
            controls.target.add(shift);
            const scale = Math.max(
                1.02,
                (right - left) / (width - padding * 2),
                (bottom - top) / (height - padding * 2),
            );
            camera.left *= scale;
            camera.right *= scale;
            camera.top *= scale;
            camera.bottom *= scale;
            camera.updateProjectionMatrix();
            render();
        }
        throw new Error(
            "Could not fit the complete diagram within snapshot bounds",
        );
    },
    inspect: () => ({
        ...diagram.inspect(Boolean(snapshotSize)),
        model,
        appearance,
        canvas: { width: surface.clientWidth, height: surface.clientHeight },
        container: {
            width: viewport.clientWidth,
            height: viewport.clientHeight,
        },
        palette,
        selected,
        drawCalls: renderer.info.render.calls,
        camera: camera.position.toArray(),
        target: controls.target.toArray(),
        zoom: camera.zoom,
        background: `#${scene.background.getHexString()}`,
        fontFamily,
        fontLoaded: document.fonts.check(`16px ${JSON.stringify(fontName)}`),
    }),
};
