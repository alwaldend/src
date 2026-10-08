# Licensed Reimu Fumo fallback research

Observed: 2026-09-05

Scope: one bounded search, followed only by inspection of the three resulting
creator listings and Sketchfab's public metadata for their exact model IDs.
No model archive was downloaded; no account, purchase, upload, or replacement
was attempted.

## Decision

Suitable, explicitly licensed alternatives do exist. The strongest physical
geometry benchmark is the red V2 scan by `kryik1023`; the only plausible clean
reusable modeling base in this set is Gorgonych's authored and rigged model,
but its CC BY-NC 4.0 restriction and unverified archive contents make it a
conditional option rather than an automatic replacement.

The two scans are useful chiefly for macro silhouette, depth, contact, cloth
thickness, and asymmetry. Neither listing demonstrates clean sewn-panel
topology, complete hidden surfaces, or a rig, so neither should be treated as
a drop-in Blender construction base merely because it is downloadable.

## Candidates

### 1. V2 Reimu Hakurei Fumo 3D scan — `kryik1023`

- Creator listing: [Sketchfab model page](https://sketchfab.com/3d-models/project-v2-reimu-hakurei-fumo-3d-scan-ff926e77d0564c028bde86fd32a487c1)
- Primary metadata: [Sketchfab API record](https://api.sketchfab.com/v3/models/ff926e77d0564c028bde86fd32a487c1)
- License: [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).
  Sketchfab states that attribution is required and commercial use is
  allowed. This is explicit permission, separate from download availability.
- Access: public metadata reports `isDownloadable: true`.
- Format: the accessible listing/API does not disclose archive extensions.
  No archive was fetched, so OBJ, FBX, glTF, or `.blend` availability is not
  verified.
- Recorded mesh size: 5,053 vertices and 9,999 faces.
- Visible resemblance/construction: the thumbnail is a direct red Reimu Fumo
  scan in a seated three-quarter pose. It preserves the squat physical plush
  proportions, fabric surface, stuffed limb/foot volumes, layered sleeves and
  skirt, and real-world asymmetry better than the authored candidate below.
- Limitations: scan-like texture and surface noise, occluded or incomplete
  underside/side regions, a baked seated pose, and no demonstrated rig or
  editable panel structure. Best use: a reference-fidelity geometry benchmark
  and measurement source, not an assumed production mesh.

### 2. Fumo Reimu — Gorgonych

- Creator listing: [Sketchfab model page](https://sketchfab.com/3d-models/fumo-reimu-4cb1dec5f8a447079c2fed94bcfdbee4)
- Primary metadata: [Sketchfab API record](https://api.sketchfab.com/v3/models/4cb1dec5f8a447079c2fed94bcfdbee4)
- License: [CC BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/).
  Attribution is required and commercial use is prohibited.
- Access: the listing shows "Download 3D Model" and public metadata reports
  `isDownloadable: true`.
- Format: the creator explicitly says it was modeled and rigged in Blender
  and textured in Substance Painter, with 2048 x 2048 PBR textures. That
  establishes the authoring tool, not the archive format: the accessible
  listing/API does not prove that the downloadable archive contains a
  `.blend` file.
- Recorded mesh size: 18,573 vertices and 35,288 faces.
- Visible resemblance/construction: a clean, recognizable seated Reimu with
  the principal Fumo silhouette, face, bow, sleeves, skirt, and small feet.
  It is the strongest candidate here for an editable/rigged base if the
  archive and noncommercial constraint are acceptable.
- Limitations: the visible model reads as a smooth stylized 3D character more
  than a scanned stuffed toy; the hair, clothing, and body are comparatively
  regular and hard-surfaced. The page gives no evidence of seam-faithful sewn
  panels. Its silhouette could still be useful, but it is not automatically a
  fidelity upgrade over the current model.

### 3. Blue Reimu Hakurei Fumo 3D scan — DAsh9986

- Creator listing: [Sketchfab model page](https://sketchfab.com/3d-models/project-blue-reimu-hakurei-fumo-3d-scan-fdc1dab1cdf04d63b3493f4247279e6e)
- Primary metadata: [Sketchfab API record](https://api.sketchfab.com/v3/models/fdc1dab1cdf04d63b3493f4247279e6e)
- License: [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).
  Sketchfab states that attribution is required and commercial use is
  allowed; the creator also says it may be used with credit.
- Access: public metadata reports `isDownloadable: true`.
- Format: the accessible listing/API does not disclose archive extensions.
  The creator identifies it as a Polycam scan; no Blender source or rig is
  claimed.
- Recorded mesh size: 20,445 vertices and 40,748 faces.
- Visible resemblance/construction: a direct physical Fumo scan with useful
  stuffed head, face, body, sleeve, skirt, and foot volumes. It is credible as
  a second physical-shape cross-check.
- Limitations: this is a blue costume variant rather than the target red
  Reimu, and the preview shows stronger scan incompleteness/distortion around
  the bow, hair, and lower layered parts. Its higher face count does not imply
  cleaner topology. Best use: corroborating macro volume and depth, not as
  the primary base.

## Bounded next action if separately authorized

Use the V2 red scan first as a read-only measurement benchmark against the
current canonical views. Only if replacement reuse is still desired, fetch
Gorgonych's archive in task scratch and verify its actual files, mesh/object
structure, rig, textures, and attribution requirements before choosing it as
a base. A download button alone was not used as license evidence here.

These creator-applied licenses describe reuse of the uploaded models. They do
not by themselves resolve any separate rights in Touhou, Reimu, or the
commercial Fumo design.
