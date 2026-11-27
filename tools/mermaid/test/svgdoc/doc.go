// Package svgdoc inspects rendered Mermaid artifacts: the paint layers of an
// SVG, the canvas an SVG draws, and the dimensions of the raster projection a
// service-upload consumer renders.
//
// Mermaid writes one group per layer under the root group, and SVG paints in
// document order, so an element's offset in the document is its paint order.
// The checks that need these scans live in more than one workspace, so the
// scans live here rather than in each.
package svgdoc
