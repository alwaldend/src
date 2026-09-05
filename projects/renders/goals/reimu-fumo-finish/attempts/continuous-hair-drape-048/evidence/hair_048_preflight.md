# Ribbon receiving-surface preflight repair

First048 save was withheld: one right-band ray missed its front fabric
footprint and hit a rear panel, requesting65.381mm movement. The12mm maximum
gate rejected it. Restrict ray hits to10mm of the existing front band depth.
For a narrow edge overhang only, use a nearest fabric point within2mm; do not
use a far back-surface hit. Keep the movement limit and record fallback count
and distance. If the2mm nearest support gate fails, stop and repair the band
pattern rather than relaxing support distance. Source047 remains unchanged.
