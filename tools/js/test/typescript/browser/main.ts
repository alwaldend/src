import type { Vector3 } from "three";
import type { Coordinates } from "../model/coordinates.js";

export function drawCoordinates(position: Coordinates): HTMLCanvasElement {
    const vector: Pick<Vector3, "x" | "y" | "z"> = position;
    const canvas = document.createElement("canvas");
    canvas.width = 320;
    canvas.height = 120;
    const context = canvas.getContext("2d");
    if (!context) throw new Error("Canvas context unavailable");
    context.fillStyle = "#fff";
    context.fillRect(0, 0, canvas.width, canvas.height);
    context.fillStyle = "#000";
    context.font = "24px monospace";
    context.fillText(`${vector.x}, ${vector.y}, ${vector.z}`, 16, 64);
    return canvas;
}
