# Frozen030b winding inventory machine evidence

Complete parsed JSON from winding_inventory_030b.json, original SHA256
f4fb3d35ef2dfd1f4d2f1ddb2e0b5426c1dd1bd9fac888d2a75265362ccbce50.
Scope and limitations are in winding_inventory_030b.md. This is not an
individual disconnected-shell, curve, collision or animation proof.

```json
{"candidate":"bow_030b_candidate.blend","candidate_sha256":"d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6","runtime":{"version":"5.2.1 LTS","build_hash":"9e2066aef7ef","background":true},"scope":"MESH objects at frame1 with visible_get() true and hide_render false. Curves, hidden geometry, animation and contacts are excluded.","summary":{"visible_renderable_mesh_count":42,"evaluated_status_counts":{"closed_consistent_positive":41,"open_or_nonmanifold":1},"base_closed_count":26,"evaluated_closed_count":41,"nonfinite_evaluated_vertex_total":0,"degenerate_evaluated_face_total_threshold_1e_14_m2":0,"zero_evaluated_edge_total_threshold_1e_10_m":0,"inconsistent_manifold_evaluated_edge_total":0,"evaluated_boundary_edge_total":4,"evaluated_edges_with_more_than_two_faces_total":0,"evaluated_loose_edge_total":0},"closed_consistently_wound_negative_volume_objects":[],"method":"World-transformed bmesh topology and signed evaluated volume. Closed means every edge has exactly two incident faces; consistent means no inconsistent manifold edge winding. Signed volume classification uses +/-1e-12m3 tolerance. Degenerate thresholds are reported separately. No native flip was performed.","candidate_preserved":true,"limitations":["Closed consistent negative volume is an orientation witness, not a self-intersection or nesting audit.","Positive aggregate volume does not rule out a negative disconnected component; this inventory classifies objects, not individual shells.","Open cloth and curves are not newly required to be closed. Modifier risk is structural, not a repair authorization or an executed flip proof."]}
```
