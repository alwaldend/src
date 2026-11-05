# 092 cuff depth and actual-arm diagnosis

The inherited cuff depth is a measured collision cause. Its center is about
8.5 mm behind the actual arm section, and the saved cuff rim crosses that
section twice. Depth centering alone does not produce the requested high hand:
with the current XZ endpoints it moves the hand lower within the opening.

This is a read-only scalar analysis of the failed left outer and exported
evaluated arm. No Blender process, geometry candidate, model, goal, or Git
operation was performed by this diagnosis.

## Exact inputs

- Source087: `c171fea554811e825e2cd9d7f1d5a069c6b305e85888c6058d2390e1a20738e1`.
- Failed outer JSON: `6f1eaca20bb9ed6f7d7a0b94a8e12234b70b0f581053ef9a341224c897ecb714`.
- Actual left arm snapshot: `0ec5bd5181694f12b1b10bedb5d13bd8de203829c52016392f51a4943fd1bbfc`.
- Geometry092: `98ca7d1db97408e255fa6bb4e772e007caa8689d585c027d76d16fa395f44675`.

`sleeve_092_cuff_diagnosis.py` uses only Python scalar arithmetic on these
JSON arrays and prints the measurements. It does not execute any builder.

## Measured plane and collision

The declared top and bottom are (-72.5,-4,51) and (-58.5,-3,34.5) mm. Their
midpoint is C=(-65.5,-3.5,42.75) mm and half-vector is B=(-7,-0.5,8.25) mm.
The cuff is C+B*cos(theta)+N*h*sin(theta), with
N=(0.247383245,-0.968917711,0), front h=11/abs(Ny) mm and back
h=-10/abs(Ny) mm. It therefore lies in one oblique plane with normal
(0.742963259,0.189692747,0.641889599). The 96 saved cuff vertices agree with
that plane to 9.75 nm.

Actual arm vertices span signed distances -6.448 to +61.144 mm from this
plane: the plane cuts the arm. The cut is one closed polygon with 98 vertices
and area 99.044 mm². Its centroid is (-63.040628,-12.010370,42.418373) mm;
its Y bounds are [-17.211927,-6.793340] mm and Z bounds
[37.837131,47.007485] mm. The nominal arm-axis intersection is
(-63.353683,-12.090302,42.804345) mm, corroborating the depth discrepancy
but not substituting for the actual mesh section.

The saved polygonal cuff rim crosses the arm section at approximately:

- (-66.695295,-13.183830,46.995298) mm;
- (-59.143757,-12.146425,37.948098) mm.

46 of the 98 arm-section vertices lie outside the nominal asymmetric oval.
Its maximum normalized ellipse radius is 1.245442. This distinguishes actual
rim/arm clipping from a merely unappealing view or a liner-only problem.

## A measured depth equilibrium, not a new candidate

For diagnosis only, retain the exact092 XZ endpoints and shift both endpoint
Y values by the same scalar. Recompute the actual arm-plane section after
each scalar evaluation and solve for its centroid to sit midway between the
11/10 mm cuff lobes at that centroid's height. The bracket comes from the
actual arm Y bounds; 42 bisections solve this one condition. This is not an
offset sweep or a search for a visually acceptable candidate.

The equilibrium is ΔY=-8.436989 mm, giving C_y=-11.936989 mm (endpoint Y
values -12.436989 and -11.436989 mm). At that diagnostic plane the arm section
centroid is (-64.429598,-12.365000,41.637537) mm. Its area is 76.412 mm²;
no section vertex lies outside the oval, the maximum normalized radius is
0.623535, and the saved rim translated by that Y amount has zero planar
crossings with the arm section.

| Quantity | Failed092 | Depth equilibrium with unchanged XZ |
|---|---:|---:|
| Hand centroid above cuff bottom, fraction of cuff height | 47.990% | 43.258% |
| Hand centroid Z | 42.418 mm | 41.638 mm |
| Rim/arm-section crossings | 2 | 0 |

Thus this Y correction resolves the measured planar depth overlap while
moving the hand approximately 0.781 mm lower. It does not establish the
desired high-hand construction. The solved offset belongs to the incorrect
or disputed092 outline and must be recomputed for any corrected XZ trace.

## Reference cause and remaining scope

The092 plan explicitly used approximate canonical-front top/bottom trace
(714,684)/(670,736), and called inherited Y an assumption. The coordinator's
subsequent image inspection reports that (714,684) was likely a mid-cuff
location: the actual outer top is near (714,650), with bottom near (638,734).
Independent reference measurement is being handled separately. This diagnosis
did not remeasure those images and does not approve replacement endpoints.
The short, low-hand result therefore has two distinct causes to resolve:
the endpoint trace and the independently demonstrated rearward depth center.

These are finite polygonal plane-section measurements. Zero planar crossings
and an ellipse-containment result do not certify full three-dimensional arm,
liner, skirt, dash, or 0.6 mm clearance gates. "High inside" needs the corrected
reference landmarks; no target height fraction is invented here.
