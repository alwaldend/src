# Test a head-centered cylindrical surface chart

Root rejects035's single XZ graph prerequisite after the one permitted
ring still left two incompatible boundary depths. No additional ring
growth is authorized. The new question changes the surface coordinates,
not the model shape, camera or depth constraints.

The strongest objection to another chart is arbitrary coordinate shopping
to retain a bad boundary. Here the alternative has a physical basis:
the disputed side roots wrap around the head's vertical centerline, and
front projection discards precisely their differing depth. Test angle
around the model Z axis and height, with radius as the scalar unknown:
theta=atan2(X,-Y), height=Z, radius=sqrt(X*X+Y*Y). The axis is fixed at
X=Y=0 in the model frame; do not fit or sweep its location or direction.

Use only the exact438 recorded expanded-boundary points in
head_035_domain_probe.json, SHA256
ee3fb144c2eaf9cbe99c6b9160763063e31f53d618aa220da6826cad8666fa9f.
The chart can fail through angular wrapping, a pole, boundary crossings
or conflicting radius values. Test those once before any implementation.
This does not authorize a triangulation, radius solve, model creation or
save. Boundary simplicity alone cannot establish retained-face conditioning,
shell clearance, seam quality or reference fidelity.
