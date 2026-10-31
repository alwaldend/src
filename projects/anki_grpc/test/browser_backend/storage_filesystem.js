import { createEmscriptenFilesystem } from "/emscripten.js";

// The test application owns its provider and root; the backend borrows both.
export const fs = await createEmscriptenFilesystem();
export const root = "/browser-backend";
fs.mkdirSync(root);
