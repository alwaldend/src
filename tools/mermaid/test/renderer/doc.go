// Package renderer asserts the observable contract of the Mermaid theme:
// which colors and shapes the renderer applies, how it isolates fonts, and the
// geometry it gives boxed labels and nested containers.
//
// The checks read the rendered SVG rather than the theme JSON wherever the
// renderer projects the setting into the document, so a theme edit that stops
// taking effect fails here even when the config file still looks correct.
package renderer
