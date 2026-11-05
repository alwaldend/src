import { readFile, writeFile } from "node:fs/promises";
import { magnitude, type Coordinates } from "../model/coordinates.js";

const destination = process.argv[2];
if (!destination) throw new Error("Expected a result destination");
const point: Coordinates = { x: 3, y: 4, z: 12 };
const metadata = JSON.parse(
    await readFile(new URL("../package.json", import.meta.url), "utf8"),
);
await writeFile(
    destination,
    JSON.stringify({
        node: process.versions.node,
        magnitude: magnitude(point),
        moduleType: metadata.type,
    }) + "\n",
);
