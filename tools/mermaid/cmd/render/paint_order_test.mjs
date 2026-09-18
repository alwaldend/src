import assert from "node:assert/strict";
import path from "node:path";
import puppeteer from "puppeteer";
import { liftClusterLabels } from "./paint_order.mjs";

// Dagre emits empty layers and independently transformed roots for nested
// subgraphs. Verify visible coordinates and stacking in the pinned browser.
const source = `<svg xmlns="http://www.w3.org/2000/svg" width="400" height="300">
  <g transform="translate(8,12)"><g class="root">
    <g class="clusters"/><g class="edgePaths"/><g class="edgeLabels"/>
    <g class="nodes"><g class="root" transform="translate(30,50) scale(1.2)">
      <g class="clusters"><g class="cluster" transform="translate(2,3)">
        <g class="cluster-label" transform="translate(4,6)">
          <rect id="outer-title" width="40" height="18" fill="white"/>
        </g>
      </g></g>
      <g class="nodes"><g class="root" transform="translate(60,50)">
        <g class="clusters"><g class="cluster">
          <g class="cluster-label ">
            <rect id="inner-title" width="40" height="18" fill="white"/>
          </g>
        </g></g><g class="nodes"/>
      </g></g>
    </g></g>
    <rect id="cover" width="390" height="280" fill="red"/>
  </g></g>
</svg>`;

const browser = await puppeteer.launch({
    executablePath: path.resolve(process.argv[2]),
    headless: "shell",
});
try {
    const page = await browser.newPage();
    const measure = async (svg) => {
        await page.setContent(svg);
        return page.evaluate(() =>
            ["outer-title", "inner-title"].map((id) => {
                const rect = document
                    .getElementById(id)
                    .getBoundingClientRect();
                return {
                    id,
                    x: rect.x,
                    y: rect.y,
                    width: rect.width,
                    height: rect.height,
                    top: document.elementFromPoint(
                        rect.x + rect.width / 2,
                        rect.y + rect.height / 2,
                    ).id,
                };
            }),
        );
    };
    const before = await measure(source);
    const lifted = liftClusterLabels(source);
    const after = await measure(lifted);
    for (let index = 0; index < before.length; index += 1) {
        assert.equal(before[index].top, "cover");
        assert.equal(after[index].top, after[index].id);
        for (const dimension of ["x", "y", "width", "height"]) {
            assert.ok(
                Math.abs(before[index][dimension] - after[index][dimension]) <
                    0.001,
                `${after[index].id}: ${dimension} changed during the lift`,
            );
        }
    }
    assert.equal(liftClusterLabels(lifted), lifted);
    assert.equal(liftClusterLabels("<svg/>"), "<svg/>");
    console.log("Title stacking and nested coordinates verified");
} finally {
    await browser.close();
}
