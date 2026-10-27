// Compile Docsy's document corpus and package the matching browser runtime.
import fs from "node:fs/promises";
import path from "node:path";
import { createRequire } from "node:module";
import { parseArgs } from "node:util";

const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const lunr = require("lunr");
const { values } = parseArgs({
    options: { documents: { type: "string" }, output: { type: "string" } },
});
if (!values.documents || !values.output)
    throw new Error("--documents and --output are required");
const documents = JSON.parse(
    await fs.readFile(
        path.join(values.documents, "search-documents.json"),
        "utf8",
    ),
);
if (!Array.isArray(documents))
    throw new Error("Search documents must be an array");
const references = new Set();
for (const document of documents) {
    for (const field of ["ref", "title", "body", "description", "excerpt"]) {
        if (typeof document[field] !== "string")
            throw new Error(`Invalid document ${field}`);
    }
    if (references.has(document.ref))
        throw new Error(`Duplicate search reference: ${document.ref}`);
    references.add(document.ref);
}
// Keep Docsy's fields and index-time boosts. Query-time boosts belong to the worker.
const index = lunr(function () {
    this.ref("ref");
    this.field("title", { boost: 5 });
    this.field("categories", { boost: 3 });
    this.field("tags", { boost: 3 });
    this.field("description", { boost: 2 });
    this.field("body");
    documents.forEach((document) => this.add(document));
});
await fs.mkdir(path.join(values.output, "assets/json"), { recursive: true });
await fs.mkdir(path.join(values.output, "assets/js"), { recursive: true });
await fs.writeFile(
    path.join(values.output, "assets/json/prebuilt-search-index.json"),
    JSON.stringify({
        index: index.toJSON(),
        documents: documents.map(({ ref, title, excerpt }) => ({
            ref,
            title,
            excerpt,
        })),
    }),
);
await fs.copyFile(
    require.resolve("lunr/lunr.min.js"),
    path.join(values.output, "assets/js/lunr.js"),
);
await fs.mkdir(path.join(values.output, "static/licenses"), {
    recursive: true,
});
await fs.copyFile(
    require.resolve("lunr/LICENSE"),
    path.join(values.output, "static/licenses/lunr.txt"),
);
console.log(`Indexed ${documents.length} documents with Lunr ${lunr.version}`);
