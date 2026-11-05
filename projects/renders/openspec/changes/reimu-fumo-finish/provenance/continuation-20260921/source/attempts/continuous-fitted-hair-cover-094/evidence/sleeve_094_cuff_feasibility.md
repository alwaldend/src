# 094 corrected-cuff scalar feasibility

Correcting the093 cuff trace and recomputing depth centering admits the actual
arm section through the opening, but still leaves the hand low and projecting
beyond the cuff plane. This does not justify another candidate while the
reference-to-model registration and broader arm/root placement remain open.

## Inputs and method

The final `sleeve_093_reference_landmarks.md` SHA256 is
`4eea0441bad88253de71343b2dc2f5f2f4d755c06c127c46bad15744a6f10f88`.
Its corrected outer white cuff endpoints are (712,649) and (647,738), with
uncertainty ±4 and ±8 pixels respectively. Its attachment markers are not
exposed physical throat endpoints.

The scalar functions were reused from `sleeve_092_cuff_diagnosis.py`, verified
against SHA256 `aeb97bf2143fc86429af8b753d2c4c7c606ef0afbb096558f7f65d6eac5a524c`.
Only the prefix through the function definitions was evaluated, stopping
before the old092 equilibrium solve. The original input guards verified:

- Actual left arm snapshot: `0ec5bd5181694f12b1b10bedb5d13bd8de203829c52016392f51a4943fd1bbfc`.
- Source087: `c171fea554811e825e2cd9d7f1d5a069c6b305e85888c6058d2390e1a20738e1`.
- Failed092 outer: `6f1eaca20bb9ed6f7d7a0b94a8e12234b70b0f581053ef9a341224c897ecb714`.

The endpoint, center, half-vector, plane and analytic cuff-boundary bindings
were then replaced in memory. No092 rim or old092 Y offset was reused as the
new cuff. The analytic boundary uses the same 48 divisions per front/back arc
and 11/10 mm world-Y lobes, with no axial bow. No geometry was written.

## Corrected endpoints and measured depth equilibrium

Using the requested provisional map, scale=0.1165/368 m per pixel, center=485
and ground=845, then reflecting canonical right X onto the actual left arm:

| Endpoint | X | Initial Y assumption | Z | X/Z trace uncertainty |
|---|---:|---:|---:|---:|
| Upper | -71.862772 mm | -4 mm | 62.048913 mm | ±1.266304 mm |
| Lower | -51.285326 mm | -3 mm | 33.873641 mm | ±2.532609 mm |

The midpoint is (-61.574049,-3.5,47.961277) mm. The plane normal toward the
proximal arm is (0.787646127,0.201101139,0.582384676); its negative points
distally out of the opening. The depth direction remains
N=(0.247383245,-0.968917711,0).

Recomputing the actual arm-plane section and solving its centroid's depth
balance between the asymmetric cuff lobes gives **ΔY=-7.145523 mm**. The
balanced cuff center is (-61.574049,-10.645523,47.961277) mm. Corresponding
endpoint Y values would be -11.145523 and -10.145523 mm. This scalar result
belongs to this corrected XZ pair; it is not the old092 shift of -8.436989 mm.

The solve used the actual arm Y extrema as a bracket and 42 bisections of one
depth-balance condition. It was not a candidate or offset sweep.

| Actual-arm measurement | Corrected XZ, inherited Y | Corrected XZ, balanced Y |
|---|---:|---:|
| Arm-section centroid XYZ, mm | (-57.929961,-10.705518,45.520942) | (-59.170989,-11.022376,44.841389) |
| Centroid height above bottom / cuff height | 41.339% | 38.927% |
| Section area | 170.774 mm² | 155.056 mm² |
| Section vertices outside analytic oval | 38 / 98 | 0 / 96 |
| Maximum normalized ellipse radius | 1.289500 | 0.705610 |
| Polygonal rim/arm-section crossings | 2 | 0 |
| Maximum distal arm extension past cuff plane | 12.663157 mm | 11.226184 mm |
| Arm vertices distally beyond cuff plane | 300 / 1106 | 289 / 1106 |

At balanced Y, the actual section spans Z=38.785138–50.896496 mm within cuff
endpoints Z=33.873641–62.048913 mm. The centroid sits 10.967748 mm above the
bottom and 17.207524 mm below the top. This is a low section, not evidence of
the desired hand high inside the opening.

The distal extension is a perpendicular distance along the outward plane
normal, not world X or arm-axis distance. Its extreme actual vertex is
(-68.977877,-13.526263,39.693087) mm. This demonstrates actual arm protrusion
through the analytic opening even after its rim admits the planar section.

## Root-region discrepancy remains separate

Under the same provisional map, the093 upper attachment marker (579,588)
maps to |X|=29.758152 mm, Z=81.360054 mm. The lower marker (561,637) maps to
|X|=24.059783 mm, Z=65.847826 mm. These are visible loop/attachment witnesses,
not a fully observed throat ring; they must not be averaged into a new root.

The currently measured root center is |X|=38 mm, Z=56.829788 mm. Its actual
ring spans approximately |X|=32.315–43.685 mm and Z=46.989–66.670 mm. Relative
to the mapped marker region, the current root is substantially lower and
outboard. Thus a broader arm/root placement error is credible, in addition
to the now-corrected cuff-definition error.

That conclusion is conditional on the pattern mapping. The093 trace is
explicitly not an approved camera registration; projection, framing, landmark
roles, and occlusion must be reconciled before changing the arm or root.
This calculation neither diagnoses a camera error nor authorizes an arm
redesign.

## Limits

The scalar analysis used the actual evaluated left arm and an analytic
corrected cuff polygon. The normalized-radius and zero-crossing results
describe its planar section, not a completed sleeve. They do not certify
three-dimensional wall clearance, liner, skirt, dashes, or visual acceptance.
No Blender/native process, geometry candidate, model, helper, or goal change
was made. Only this new feasibility report was written.
