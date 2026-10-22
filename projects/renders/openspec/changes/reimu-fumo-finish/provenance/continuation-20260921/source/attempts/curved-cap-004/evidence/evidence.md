# Curved-cap reset evidence

Blender MCP 5.1.1 hid only the A157 continuous hair cap, retained the
existing soft head cushion and separate fringe and lock objects, created one
partial curved padded cap, and saved candidate
`58b77d791b64ad0cb2034f26dd5caa603f9bd9f3e493d646e011813b9789d6aa`
without changing protected A157.

Pinned Blender 5.2.1 clean-reopened that exact candidate and rendered the
frozen front and side cameras. The immutable packet manifest binds the
candidate before and after rendering and records:

- front:
  `8e9604f7fcf620aba53f046d49f95a65087871a9360cf5b422266a1df580b304`;
- side:
  `0dc7b196580ac3d71a9f02ed8ce3f9cd8dbe8e22c3a04d2278fbbf31db3905c7`.

Pixel inspection rejected the candidate before measurements or full-view
rendering. In front view, the new shell remains a continuous raised band
around the face. In side view, it is a single deep, nearly rectangular brown
slab spanning crown to nape. These are the explicit helmet and card failures
for `head-hair-helmet-read` in `PROCESS.md`.

This is the second consecutive reviewed failure of a continuous head-covering
representation. The stop rule therefore retires that representation family;
changing its dimensions would not be an authorized correction.
