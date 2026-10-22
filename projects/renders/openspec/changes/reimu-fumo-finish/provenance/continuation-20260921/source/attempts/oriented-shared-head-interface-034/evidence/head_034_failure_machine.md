# Selected034 failure diagnostics

Raw shared diagnostic SHA256
003720d883be055ed4bc6c62404b12ab10c1081914f4c07a0ab12996b9ff1af1.
Raw co-normal diagnostic SHA256
e2d100441a9bb65b8f35034456d64109f7b9d9b3dc0bf6bc56ba11c5fcd8369a.
The latter is approximately15MB and remains ignored; all arrays in its
compact view below are explicitly omitted with their lengths. Witnesses
and measured distributions are described in the paired diagnostic notes.

## Shared first state summary

```json
{
    "status": "root_review_required_before_writer",
    "source_sha256": "6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8",
    "evaluated_helper_sha256": "c9976b58dc310001b4478ddb399b6691bd801eb3a725b12358c6422000cb1c69",
    "entrypoint": "build_head_034_shared",
    "targets": ["Hair028 traced padded fringe", "Hair028 crown and back hood"],
    "created_object": "Hair034 shared crown fringe shell",
    "geometry_evaluations": 1,
    "blend_saved": false,
    "rendered": false,
    "source_preserved": true,
    "protected_non_target_controls_exact": 84,
    "retained_off_support_triangles": 72718,
    "off_support_coordinate_max_error_m": 0.0,
    "source_uv_layers": {"Hair028 traced padded fringe": [], "Hair028 crown and back hood": []},
    "boundary_supported_source_triangles": 366,
    "common_arc_knots": 432,
    "removed_original_fringe_rims": 306,
    "new_vertices": 44366,
    "new_faces": 81704,
    "nonmanifold_or_inconsistent_edges": 0,
    "strip_chord_m": {"min": 0.0008200246220152397, "median": 0.008154836206701137, "max": 0.010325792913552245},
    "paired_skin_separation_m": {"min": 0.0007105489915976281, "median": 0.0009124467630732082, "max": 0.001158310765263781},
    "ownership_flag": {"view": "three_quarter_mirror", "y": 175, "regions": [{"region": "retained_hood_outer", "x": [214, 236]}, {"region": "shared_outer_bridge", "x": [237, 248]}, {"region": "retained_fringe_outer", "x": [249, 249]}, {"region": "shared_outer_bridge", "x": [250, 250]}, {"region": "retained_fringe_outer", "x": [251, 279]}]},
    "flag_interpretation": "Repeated face-region ownership on one welded surface may be a projected cut-boundary reversal; not proof of competing overlapping sheets.",
    "unavailable_evidence": ["Exact per-hit candidate face/triangle IDs and local depth/normals at the reversal were not recorded.", "No strip intersection, envelope overshoot or junction tangent audit was run.", "Off-support shading normals were recomputed after triangulation and were not checked for equality."],
    "full_evidence": "head_034_shared_dryrun.json"
}
```

## Co-normal state scalar evidence

```json
{
  "source_sha256": "6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8",
  "helper_sha256": "68a035190a68d602813e93a501ad25208b68680483bdeb4fed02e3dff20c234d",
  "execution_succeeded": false,
  "error": "Traceback (most recent call last):\n  File \"/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/head_034c_diagnostic.py\", line 26, in <module>\n    construction=scope['build_head_034c']()\n  File \"/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/head_034c_draft.py\", line 438, in build_head_034c\n    _strip_preflight(builder)\n    ~~~~~~~~~~~~~~~~^^^^^^^^^\n  File \"/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/head_034c_draft.py\", line 173, in _strip_preflight\n    assert PREFLIGHT['bad_jacobian_count']==0 and PREFLIGHT['over90_junction_count']==0, (\n           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\nAssertionError: ('Fixed cut/correspondence still folds or collapses; no object created', 456, 6)\n",
  "preflight": {
    "phase": "local strip Jacobian and junction checks before object creation",
    "endpoint_directions": {
      "omitted_array_length": 1728
    },
    "derivative_vectors": {
      "omitted_array_length": 864
    },
    "invalid_endpoint_count": 0,
    "strip_jacobians": {
      "omitted_array_length": 6896
    },
    "precreation_junction_normals": {
      "omitted_array_length": 862
    },
    "bad_jacobian_count": 456,
    "over90_junction_count": 6,
    "junction_statistics": {
      "fringe_bridge": {
        "count": 431,
        "min": 0,
        "median": 1.7565789571219343,
        "max": 146.5000286321367
      },
      "bridge_hood": {
        "count": 431,
        "min": 0.12373852703836481,
        "median": 1.4609222578158507,
        "max": 167.18138989797092
      }
    }
  },
  "runtime": "5.2.1 LTS",
  "build_hash": "9e2066aef7ef",
  "created_object_present": false,
  "blend_saved": false,
  "rendered": false,
  "source_and_helper_preserved": true
}
```

Supplemental033 saved-file audit JSON SHA256
413174d450aa375727ab063d2de562dca4b2ca90e57e749a410c245541098119.
Its readable report is already canonical in closed033.
