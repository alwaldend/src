# References and landmark contract

[Back to goal](README.md)

## Controlling image

- Original Sisyphus PNG:
  `/var/home/simeonwarrenbot/.t3/userdata/attachments/136a19c9-b38c-4c43-959e-141b9e2224e3-d73fe574-eda6-4832-82e1-e3d0d4252fa4.png`
- Resolution: `1090 × 1127`.
- SHA-256:
  `3d40e2726ae5ff84983f642e20809bb6689c77ecffe5060c6aa760bdee121519`.
- The quote is source context, not part of the quote-free scaffold camera.

## Normalized image-space landmarks

Coordinates use `(x / width, y / height)` from the upper-left. Values are
visual measurements with about `±0.025` uncertainty, not photogrammetry.

- Boulder bounding box: approximately `(0.14, 0.11)` to `(0.61, 0.57)`;
  center `(0.37, 0.34)`.
- Boulder bottom/contact neighborhood: approximately `(0.37, 0.55)`.
- Primary open-sky field: left of the diagonal rock face, strongest below the
  quote and around the figure.
- Main climb direction: lower-left to upper-right, approximately `45–55°`.
- Lower slope enters near `(0.00, 0.96)` and reaches the boulder contact near
  `(0.55, 0.55)`; the large right rock face owns most of the right half.
- Human slot in the source: approximately `(0.14, 0.45)` to `(0.31, 0.87)`.
  This locates only the future character envelope; it does not authorize a
  humanoid or Fumo replacement in this attempt.

## Authority

The controlling image governs composition, boulder-to-frame scale, climb
angle, light direction, warm value structure, and negative sky space. Existing
scene renders are failure evidence only. The original image must be packed
into the candidate and placed in a `REFERENCE_ONLY` collection at the side of
the working set; it must not overlay the beauty render.
