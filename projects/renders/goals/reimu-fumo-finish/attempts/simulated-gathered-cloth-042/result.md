# Real cloth gives folds but fails garment acceptance

042b SHA819207b0d96958e6e86110c30399db45db183772bf66cb03296671aef4ce764b
saved, clean-reopened and rendered fixed front/side. Real solver folds replace
the former prescribed cone, but the white hem becomes a flat low band and
the fabric hangs too far around the toes. Whole-character appearance still
fails. Initialization history is in cloth_042_mechanics.md.

042c SHA0cfd8341b54bc0d78c43e0ebb55b1f7534585a19505181c86b5bc35be04e6686
changes requested rest width53mm to36mm and adds one post-solve subdivision.
Clean-reopened front/side show only smoothing, not a shorter garment.
Direct saved-mesh comparison:2016 vertices each, maximum position delta0.0m.
The rest-shape parameter did not affect simulation, so no length improvement
is claimed. Stop parameter tweaking; a read-only source/API diagnosis is
delegated to determine why rest_shape_key was ineffective.

Decision reset.042c remains an unaccepted experimental base for a bounded
visible-pile and sewn-detail test, not retained macro geometry. Fresh visual
review must keep silhouette and contact failures separate from material read.
No animation, reusable-export or final criterion passes are claimed.
