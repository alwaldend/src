export interface Coordinates {
    x: number;
    y: number;
    z: number;
}

export function magnitude(point: Coordinates): number {
    return Math.hypot(point.x, point.y, point.z);
}
