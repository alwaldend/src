# Root cylindrical boundary diagnostic

No model was opened or changed. This is one mathematical transformation
of the exact035 recorded boundary, using a fixed physical model axis.
Angle is scaled by Wh/2 for metric units; segment intersection tests use
all nonadjacent pairs, including closing-edge pairs. Adjacent shared
endpoints are excluded. Nonparallel contacts include endpoints with1e-9
parameter tolerance; parallel collinear overlaps use1e-14 cross-product
and1e-8 overlap-parameter tolerances. No chart axis was fitted or swept.
The source's original3D boundary is unchanged. Retained-face conditioning
and inner-shell validity are separate prerequisites, not measured here.

```json
{
  "source_probe_sha256": "ee3fb144c2eaf9cbe99c6b9160763063e31f53d618aa220da6826cad8666fa9f",
  "method": "theta=atan2(X,-Y), Z; fixed model axis X=Y=0; angle scaled by Wh/2 only for numeric units",
  "boundary_count": 438,
  "theta_range_rad": [
    -1.1410863105814149,
    1.1410849468714916
  ],
  "radius_range_m": [
    0.02230328507721538,
    0.06581318620915973
  ],
  "min_chart_edge_m": 0.00016448605909558206,
  "nonadjacent_contacts": [],
  "geometry_changed": false,
  "limitations": [
    "Boundary test only; retained-face conditioning, inner chart and solver not tested.",
    "Uses piecewise-linear chart edges between transformed source vertices."
  ]
}
```
