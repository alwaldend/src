# Reject the first sewn lower assembly

Candidate0ef451284c084a12d71d0c097110929fcb5501362bd480837d9296ea71699b03
is not retained. Keep023 as the modeling baseline. Preserve024 only as
failure evidence: light legs with dark toes are a useful reference direction,
but this implementation introduces overly straight exposed cylinders and a
sharp triangular-tooth hem. Neither is softly stuffed/gathered construction.

Independent reviewer lower_024_blind inspected references first, then all five
fixed views before baseline023. Absolute scores /10: likeness5, silhouette5,
construction3, identity6, contact5, plush read3, presentation6. After baseline
comparison it rejected whole-candidate retention. Root inspected all five
views and agrees. Existing helmet hair, rigid sleeves/bow and bodice persist.

Independent technical audit head_019_technical: all80 non-target controls and
inspected rig record exactly match023. Three new meshes are finite, manifold,
weighted, without zero edges or degenerate faces. Hem-root to evaluated red
midsurface max5.59nanometers; reverse polyline max0.03268mm. The3080 tested
hem/red edge crossings stay in the intended seam-thickness band, within
0.384mm of seam. No tested hem/leg or red/leg edge crossings in either
direction. Hem/leg sampled minimum clearance0.327mm left/.282mm right.
Floor clearance hem0.825mm, legs0.100mm. These technical results cannot
override the decisive visual failure. Candidate and source bytes unchanged.

Root-cause review through attempts022–024: torso proportions and collar
contact improved, but code-defined cloth repeatedly reads rigid where curves
flatten at every control station. In024 each ruffle profile segment uses
zero-slope smoothstep endpoints, creating a sequence of lobes that projects
as triangular teeth rather than rolling cloth. Foot profile has too much
constant-radius exposed length and a sharp rim under diagnostic lighting.
Retain the demonstrated evaluated-edge/support method, not these surfaces.

Next hypothesis: shorter continuously tapered stuffed legs with gently domed
toe panels, plus a continuously curling excess-length ribbon. Use a smooth
closed fold trajectory in normal/tangent/height directions instead of flat
stationary profile sections. Preserve023 upper model and global height.
Do not add stitches, fibers or final materials to conceal construction flaws.
The unsaved palette test also showed that lighter skin removes a misleading
brown-face read but does not fix the helmet/cheek-card geometry.

Feedback route remains effective: save/reopen/five renders took about7.2s;
technical and independent pixel reviews settled the result quickly. Most
time was helper design, not render latency or screenshots. Root owns writes,
workers inspect immutable candidates. No GUI investigation is needed.

Artifacts under out/reimu_fumo_finish/desktop_astra/: lower_024_candidate.blend,
lower_024_review/, lower_024_writer_receipt.json, lower_024_technical_audit.json.
Front PNG87d1629f366a458297be08f615487a1a5eaf0456b1f504272d4d9830c53d84e6.
Three-quarter PNGc57aaf411079ffdf4608463ad1feffbec500f758f2f0a7a75fa9c5c3d782540f.
Fixed contract4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4.
