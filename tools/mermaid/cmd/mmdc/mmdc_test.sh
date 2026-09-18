#!/usr/bin/env bash
set -euo pipefail

mmdc="$1"
input="$2"
mono_anchor="$3"
output="${TEST_TMPDIR}/smoke.svg"

"${mmdc}" -i "${input}" -o "${output}"
test -s "${output}"
grep -q '<svg' "${output}"

# Interactive runs apply the checked-in repository theme by default.
grep -q 'fill:#ffffff' "${output}"
grep -q 'stroke:#1e1e1e' "${output}"
grep -q 'rough-node' "${output}"

if grep -q '#ECECFF\|#9370DB' "${output}"; then
    echo "rendered SVG still contains the Mermaid default theme colors" >&2
    exit 1
fi

width_of() {
    grep -o 'max-width: [0-9.]*px' "$1" | head -1
}

# A CLI run must resolve labels through the repository's pinned fonts. Point
# the ambient Fontconfig at a monospace-only pool: labels measured with those
# metrics render wider, so an unchanged canvas shows the run ignored them.
poison="${TEST_TMPDIR}/poison"
mkdir -p "${poison}/lib"
cp "${mono_anchor}" "${poison}/lib/"
cat >"${poison}/fonts.conf" <<FONTCONF
<?xml version="1.0"?>
<fontconfig>
  <dir>${poison}/lib</dir>
  <cachedir>${poison}/cache</cachedir>
  <alias><family>Helvetica</family><prefer><family>Liberation Mono</family></prefer></alias>
  <alias><family>Arial</family><prefer><family>Liberation Mono</family></prefer></alias>
  <alias><family>sans-serif</family><prefer><family>Liberation Mono</family></prefer></alias>
</fontconfig>
FONTCONF

# Guard against a vacuous check: the pool must exclude the proportional font,
# so a run that consulted it could not produce the same metrics.
if compgen -G "${poison}/lib/LiberationSans*" >/dev/null; then
    echo "poisoned font pool unexpectedly contains a proportional font" >&2
    exit 1
fi

poisoned="${TEST_TMPDIR}/poisoned.svg"
FONTCONFIG_FILE="${poison}/fonts.conf" FONTCONFIG_PATH="${poison}" \
    "${mmdc}" -i "${input}" -o "${poisoned}"

if ! cmp -s "${output}" "${poisoned}"; then
    echo "host Fontconfig changed the CLI canvas; rendering is not hermetic" >&2
    exit 1
fi
