import { readFile, writeFile } from "node:fs/promises";
import { REVISION, Vector3 } from "three";
import { magnitude, type Coordinates } from "../model/coordinates.js";

const destination = process.argv[2];
if (!destination) throw new Error("Expected a result destination");
const vector = new Vector3(3, 4, 12);
const point: Coordinates = vector;
const metadata = JSON.parse(
    await readFile(new URL("../package.json", import.meta.url), "utf8"),
);
await writeFile(
    destination,
    JSON.stringify({
        node: process.versions.node,
        magnitude: magnitude(point),
        vectorMagnitude: vector.length(),
        engineRevision: REVISION,
        moduleType: metadata.type,
    }) + "\n",
);
