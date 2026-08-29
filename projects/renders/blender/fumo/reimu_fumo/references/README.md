# Reimu Fumo reference sources

These user-supplied references control the reusable Reimu Fumo asset. Binary
files are tracked with Git LFS so the reference packet survives beyond one
host's attachment and temporary-output directories.

## Canonical references

- `canonical_front_25cm.png`
  - Source attachment:
    `136a19c9-b38c-4c43-959e-141b9e2224e3-82f07f2f-0207-4723-8a9b-d76b2e7e0d49.png`
  - SHA-256:
    `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`
  - Role: primary authority for the exact variant, front silhouette, overall
    proportions, hair, bow, expression, clothing, seated pose, and real-world
    scale. The complete plush is 25 cm tall.
- `canonical_turn_180.gif`
  - Source URL: `https://c.tenor.com/wuTstMILarIAAAAC/tenor.gif`
  - SHA-256:
    `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30`
  - Role: primary authority for three-quarter, side, and rear silhouettes;
    depth and layer order; seated volume; hair and bow drape; skirt pooling;
    and foot occlusion. The preserved file contains 30 coalesced 498 by 498
    frames at 10 centiseconds per frame.

The canonical front wins direct conflicts about exact variant identity,
frontal proportions, graphics, and scale. The canonical turn wins direct
conflicts about hidden-side silhouette, depth, and layer order. Do not average
their different cloth poses into a fictional shape.

## Supporting references

- `clean_front.png`
  - Source attachment:
    `136a19c9-b38c-4c43-959e-141b9e2224e3-04ac3273-048f-4ee0-a605-12a7edd4c7bf.png`
  - SHA-256:
    `37813e03e04e4966f1dbe914e03a25a5f5ae561dcbf58b72677195c513ea48ca`
  - Role: supporting graphic landmarks and neutral front silhouette. Its
    isolated presentation suppresses contact and fabric evidence.
- `physical_front.png`
  - Source attachment:
    `136a19c9-b38c-4c43-959e-141b9e2224e3-c690c3e5-d072-4c0a-bda5-7d452a501519.png`
  - SHA-256:
    `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51`
  - Role: supporting real-fabric pile, stuffing, front construction, and
    seated-contact evidence. Its proportions do not override the canonical
    front.
- `physical_side.png`
  - Source attachment:
    `136a19c9-b38c-4c43-959e-141b9e2224e3-5a8e0eba-24f7-4d02-97db-5608d16966f9.png`
  - SHA-256:
    `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d`
  - Role: supporting fabric-panel thickness, overlap, hair and bow layering,
    skirt pooling, and foot construction. The camera is oblique, so this image
    is not an orthographic dimensional controller.
- `turn.gif`
  - Source attachment:
    `136a19c9-b38c-4c43-959e-141b9e2224e3-f8df027f-9508-4b59-8779-21310993ccfa.gif`
  - SHA-256:
    `b42368e921bd055d73fbbb7bf65c2509a9aaf190cab02f89824b92b4cb75ece4`
  - Role: qualitative cross-check for front-to-side-to-rear continuity where
    the canonical turn is occluded. It is low-resolution and shows a different
    physical pose, so it does not control exact dimensions.
- `sofa.gif`
  - Source attachment:
    `136a19c9-b38c-4c43-959e-141b9e2224e3-631a9796-a3a8-487c-baa6-89da2eb26598.gif`
  - SHA-256:
    `7c9173f91e6b6c801a1c77e50f9135e86fc89319f3c0262c10312320b1af8589`
  - Role: qualitative evidence for seated compression, pile, panel edges,
    applique seating, and soft contact. Motion blur, occlusion, pan, zoom, and
    sofa perspective exclude exact silhouette measurements.

Supporting sources may veto an impossible fabric transition or construction,
but they do not override the canonical pair or supply averaged proportions.
The flattened preview of the canonical turn and unrelated Sisyphus imagery
are intentionally excluded from this packet.
