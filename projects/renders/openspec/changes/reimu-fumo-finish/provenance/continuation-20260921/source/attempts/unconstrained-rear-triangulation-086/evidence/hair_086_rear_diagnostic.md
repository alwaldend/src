# Four inherited rear-head marks: causal diagnosis

The four rear marks are inherited from the front brow/eyelid-constrained triangulation copied onto the curved back. There are both real triangle-interior chord depressions and biased smooth normals. They are not rear faces assigned eyebrow material or visible separate eyebrow objects.

One read-only pinned Blender 5.2.1 LTS (`9e2066aef7ef`) process inspected exact saved 084 and exited 0. No object, normal, material, model file or render was changed. No 085 asset was inspected.

## Source chain and spatial match

`macro_066_build.py:98` adds the four narrow graphic domains at image rows 415 and 441. Lines 111–128 constrain the common CDT to these graphic boundaries. Lines 150–165 reuse that triangulation for both bulged skins, reverse the rear winding, and assign every rear face material 0. The saved head confirms all 6,474 back triangles use `Fabric041 brown fabric` and have corresponding front triangles. There are no head modifiers or custom normals; the inspected faces use smooth shading.

The two X-negative strips project to rear pixels X277–313; the reflected strips to X197–233. Upper-strip Y is approximately 245.26–247.64; lower-strip Y259.69–262.07. This matches the four inherited lines in the reported fixed rear view, including the dark lower pair. Named brow/eyelid/expression objects are all `hide_render=true`; the actual Macro038 brow objects are also located on the front, at negative world Y.

## Geometry versus shading

The original analytic back-bulge formula matches every sampled region vertex within 1.87 nanometres. However, curved analytic positions at vertices do not make the planar triangle interiors curved. The inherited graphic boundaries contain edges 16–21 mm long across strips only about 0.4 mm tall. Native ray hits reveal inward chord error at strip interiors:

| Region, world X | Centerline samples | Maximum inward chord error | Mean inward chord error | Maximum smooth-normal error |
|---|---:|---:|---:|---:|
| Negative X, upper | 9 | 0.246 mm | 0.132 mm | 14.01° |
| Negative X, lower | 9 | 0.182 mm | 0.099 mm | 20.75° |
| Positive X, upper | 9 | 0.483 mm | 0.328 mm | 10.31° |
| Positive X, lower | 9 | 0.270 mm | 0.158 mm | 22.48° |

Errors compare the actual back ray hit and interpolated smooth vertex normal against the builder's analytic back position/differential normal at the same X/Z. The source has no custom normals, so this measures its existing smooth-normal field, not a hypothetical replacement.

An additional 18 control samples per strip at ±1 mm in Z give mean chord errors of 0.043/0.070/0.085/0.085 mm, respectively. Mean smooth-normal errors on the strips are 7.82/11.14/7.55/11.23°, versus 5.41/2.54/10.22/8.69° on those controls. The affected normal field therefore extends into neighboring triangles; a centerline-only shading edit would not isolate it completely. Maximum neighboring facet dihedrals at strip boundaries reach 161–172° despite both incident faces using the same brown material.

Strong physical witness: world `(23.268342, 37.453815, 135.447085)` mm, rear pixel approximately `(215.20,246.45)`. The analytic back at that X/Z has Y37.936893 mm, leaving a 0.483078-mm inward chord depression. Smooth-normal deviation there is 10.31°, while its facet normal differs by 49.84°.

## Actionable repair distinction

Analytic custom split normals on the affected rear region could correct much of the biased shading field without changing front coordinates or topology. That is a plausible low-impact visual diagnostic, not an established complete repair: normals do not remove the measured 0.18–0.48-mm physical depressions, geometric shadowing, or already constructed pile placement/orientation. A normal-only render was not authorized or performed.

For a geometric repair, the smallest targeted scope is the four rear strip neighborhoods and their adjacent long graphic edges: remove the unnecessary rear graphic constraints, retessellate with adequate interior samples, and evaluate new points from the unchanged analytic back formula. Simply subdividing the existing triangles by linear interpolation would preserve their depressions. Preserve the patch boundaries, all front geometry/material domains, the outer head outline, and unrelated surfaces. No full-head shape redesign is indicated by this evidence.

The visible 130,000-strand head pile is generated from head normals; subsequent 076 reseating uses receiver facet normals. It can amplify the existing folds or shading discontinuities. Its independent shadow contribution was not isolated, and this report does not claim every dark pixel would disappear after either proposed repair.

## Scope and receipts

Coverage: four copied strip domains, 36 centerline and 144 neighboring surface samples in total, their incident boundary edges, named eyebrow-object visibility, and head-pile presence. No render variant, pile collision audit, full-head self-intersection audit, or visual acceptance. 076 was hash-verified but not reopened; unchanged head identity between 076 and 084 is coordinator evidence.

- Saved 084: `cb6bff1fa0e654f1720236a5a90dab56346f41e36f999dfc08fe746ea26294cb`.
- Baseline 076: `c7aeaf157f7d451d658c302c6a9300125ab145a7551b5fa8713288052c747050`.
- Source builder: `d7616fedb77b121de36ca9a1cffd080e4267b07f5d4dab6b2d3e9c61e0eb8bc4`.
- Diagnostic code: `990fe584b1ddf20feee1ef415b6621b9d7c2af8ea87954d6c018be739d84b3bf`.
- Diagnostic JSON: `56e93f855daa91ee8cf0a28e7954d3851291e0278e281ca6889d605d12dc74ac`.

Both model hashes were unchanged after inspection. Full spatial/normal samples and edge witnesses are preserved in `hair_086_rear_diagnostic.json`.
