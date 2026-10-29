export type Point3 = number[];
export type Point2 = number[];
export type ColorScheme = "light" | "dark";
export type ElementKind = "node" | "group" | "edge";

export interface TextOptions {
    size?: number;
    heading?: boolean;
    compartments?: number[];
}
export interface TextLayout {
    width: number;
    height: number;
    size: number;
    heading: boolean;
    rows: { text: string; bold: boolean; y: number }[];
    rules: number[];
}
export type Measurements = Record<"normal" | "bold", Record<string, number>>;

export interface BoxElement {
    id: string;
    label: string;
    parentId?: string | null;
    position: Point3;
    size: Point3;
    lines?: string[];
    textLayout?: TextLayout;
    members?: string[];
    direction?: string;
}
export interface DiagramNode extends BoxElement {
    shape: string;
    lines: string[];
    textLayout: TextLayout;
}
export interface DiagramGroup extends BoxElement {
    members: string[];
}
export interface DiagramEdge {
    id: string;
    label: string;
    start: string;
    end: string;
    startLabel?: string;
    endLabel?: string;
    startMarker: string;
    endMarker: string;
    dashed: boolean;
    points: Point3[];
    labelPosition?: Point3;
    startLabelPosition?: Point3;
    endLabelPosition?: Point3;
}
export interface DiagramModel {
    type: string;
    direction: string;
    viewDirection: Point3;
    nodes: DiagramNode[];
    groups: DiagramGroup[];
    edges: DiagramEdge[];
}

// The fields needed from Mermaid's public layout data. Mermaid owns grammar
// and source semantics; the browser normalizes display text before measurement.
export interface GraphNode {
    id: string | number;
    label: string;
    text?: string;
    shape?: string;
    parentId?: string | number | null;
    isGroup?: boolean;
    dir?: string;
    annotations?: string[];
    members?: { text: string }[];
    methods?: { text: string }[];
    attributes?: { type: string; name: string; keys?: string[] }[];
}
export interface GraphEdge {
    id?: string | number;
    start: string | number;
    end: string | number;
    label?: string;
    startLabelRight?: string;
    startLabelLeft?: string;
    endLabelRight?: string;
    endLabelLeft?: string;
    arrowTypeStart?: string;
    arrowTypeEnd?: string;
    pattern?: string;
    thickness?: string;
    isLayoutOnly?: boolean;
}
export interface GraphData {
    nodes: GraphNode[];
    edges: GraphEdge[];
    direction?: string;
}
export interface ParsedDiagram extends GraphData {
    type: string;
}

export type Palette = Record<
    | "background"
    | "textColor"
    | "nodeBorder"
    | "nodeTextColor"
    | "secondaryColor"
    | "primaryColor"
    | "lineColor"
    | "actorLineColor"
    | "clusterBorder",
    string
> &
    Record<string, string>;
export interface RenderPayload {
    model: DiagramModel;
    appearances: Record<ColorScheme, Palette>;
    fontFamily: string;
    fontName: string;
    fontLicense: string;
    source: string;
    title: string;
    mermaidVersion: string;
}
export interface Bounds {
    left: number;
    right: number;
    top: number;
    bottom: number;
}
export interface LabelInspection {
    text: string;
    kind: string;
    ownerId?: string;
    bounds: Bounds;
    fontPixels: number;
}
export interface SceneInspection {
    sceneCorners: Point2[];
    nodeBounds: { id: string; bounds: Bounds }[];
    labels: LabelInspection[];
    meshCount: number;
    elements: {
        key: string;
        id: string;
        kind: ElementKind;
        opacities: number[];
    }[];
}
export interface ViewInspection extends SceneInspection {
    model: DiagramModel;
    appearance: ColorScheme;
    canvas: { width: number; height: number };
    container: { width: number; height: number };
    palette: Palette;
    selected: string;
    drawCalls: number;
    camera: Point3;
    target: Point3;
    zoom: number;
    background: string;
    fontFamily: string;
    fontLoaded: boolean;
}
export interface DiagramView {
    ready: boolean;
    inspect(): ViewInspection;
}
declare global {
    interface Window {
        diagramView: DiagramView;
    }
}
