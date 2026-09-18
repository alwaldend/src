#!/usr/bin/env bash
set -euo pipefail

svg="$1"

test -s "${svg}"
grep -q '<svg' "${svg}"

# The renderer must apply the repository theme rather than Mermaid's built-in
# default, whose lavender nodes render poorly on the documentation canvas.
grep -q 'fill:#ffffff' "${svg}"
grep -q 'stroke:#1e1e1e' "${svg}"

# The theme requests the hand-drawn look, so nodes carry Mermaid's rough
# renderer class instead of the plain vector node class.
grep -q 'rough-node' "${svg}"

if grep -q '#ECECFF\|#9370DB' "${svg}"; then
    echo "rendered SVG still contains the Mermaid default theme colors" >&2
    exit 1
fi

# Label geometry follows the pinned fonts, not the host Fontconfig. The action
# resolves labels through the repository's Liberation fonts, so a canvas
# produced by other host fonts indicates the font isolation stopped working.
expected_width="688.087px"
actual_width="$(grep -o 'max-width: [0-9.]*px' "${svg}" | head -1)"
if [[ "${actual_width}" != "max-width: ${expected_width}" ]]; then
    echo "expected canvas width ${expected_width}, got: ${actual_width}" >&2
    exit 1
fi
