# Foot033 selected machine receipts

Selected writer fields and full render receipt follow. Large control,
material-node and pose dictionaries are omitted from this view; their
successful comparisons are reported, not reproduced here. The raw writer
receipt remains ignored beside the candidate. No final acceptance.

Raw writer receipt SHA256
9ae2e7ec2440502f2bc627865a501fcf3f3d1892ffdad5b4d7ad4b19611b9eb9.
Render receipt SHA256
ad1d2d8e18cc1ba9c621f0b94f7944f9723a9b91a2d56bf7e87398c846475e0f.
Execution log SHA256
5cfa273711a28032cd8315563bd91b73d91523481c4b01a4d4f65b7a3c7be483.

## Selected writer receipt

```json
{
  "candidate": "foot_033_candidate.blend",
  "candidate_sha256": "98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48",
  "source": "head_032_candidate.blend",
  "source_sha256": "6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8",
  "version": "5.2.1 LTS",
  "build_hash": "9e2066aef7ef",
  "writer_sha256": "534e69c648a6097c76eb419717971e05df0492c4a9221711ed6f34d7ccf99660",
  "helper_sha256": "4b4ef4d8ef9a692382e9770f2503be7a2268b5b101bd964d77a8f5aa8b194a34",
  "target_names": [
    "Left black stuffed foot pod",
    "Right black stuffed foot pod"
  ],
  "control_count": 84,
  "controls_unchanged": true,
  "rig_pose_unchanged": true,
  "appearance_unchanged": true,
  "limitations": [
    "No visual, animation, whole-scene technical or final acceptance."
  ],
  "foot_construction": {
    "status": "Foot033 built and mask-tested in memory; no save, render, or stage pass",
    "decision": "front_cream_nonzero_only_requires_root_visual_review",
    "source": "/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-a13ca48d/out/reimu_fumo_finish/desktop_astra/head_032_candidate.blend",
    "source_sha256": "6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8",
    "writer_receipt_sha256": "9e98892445a0520e360276fd8811ea9415f992218a65fb7584506cb880fef4e4",
    "render_receipt_sha256": "1d5c9b6fc5d432376e8c78d06fb3b4a4f1eee81b05916346b5992a5c5d858013",
    "helper_sha256": "4b4ef4d8ef9a692382e9770f2503be7a2268b5b101bd964d77a8f5aa8b194a34",
    "target_names": [
      "Left black stuffed foot pod",
      "Right black stuffed foot pod"
    ],
    "created_names": [
      "Foot033 left conformal two-material stuffed pod",
      "Foot033 right conformal two-material stuffed pod"
    ],
    "created_mesh_names": [
      "Foot033 left conformal two-material stuffed pod mesh",
      "Foot033 right conformal two-material stuffed pod mesh"
    ],
    "created_collection_name": "Foot033 conformal two-material pods",
    "hidden_names": [
      "Left black stuffed foot pod",
      "Right black stuffed foot pod"
    ],
    "head_width_m": 0.11743925511837006,
    "formula": {
      "coordinates": "u=-side*(x-cx)/rx; p=(y-cy)/ry; zeta=(z-cz)/rz",
      "field": "F=zeta+0.60*u+0.60*p-0.64",
      "cream_rule": "F>0",
      "black_rule": "F<=0",
      "sweep_count": 0
    },
    "construction_metrics": {
      "Foot033 left conformal two-material stuffed pod": {
        "source_name": "Left black stuffed foot pod",
        "construction": {
          "side": "left",
          "source_bounds_m": [
            [
              -0.05128399282693863,
              -0.022316008806228638
            ],
            [
              -0.08391998708248138,
              -0.03607999533414841
            ],
            [
              -0.0006000008434057236,
              0.02484000101685524
            ]
          ],
          "source_half_extents_m": [
            0.014483992010354996,
            0.02391999587416649,
            0.012720000930130482
          ],
          "source_polygon_count": 2016,
          "source_triangle_count": 3920,
          "output_vertex_count": 2106,
          "output_triangle_count": 4208,
          "crossing_vertex_count": 144,
          "source_signed_volume_m3": 0.000019262366813912825,
          "output_signed_volume_m3": 0.000019262374863959585,
          "winding_reversed_by_guard": false,
          "donor_topology": {
            "vertices": 1962,
            "edges": 5880,
            "faces": 3920,
            "components": 1,
            "euler": 2,
            "genus": 0,
            "non_two_incident_edges": 0
          },
          "output_topology": {
            "vertices": 2106,
            "edges": 6312,
            "faces": 4208,
            "components": 1,
            "euler": 2,
            "genus": 0,
            "non_two_incident_edges": 0
          },
          "seam_edge_count": 144,
          "seam_vertex_count": 144,
          "seam_components": 1,
          "non_degree_two_seam_vertices": 0,
          "maximum_material_side_violation": 2.498834162922847e-7,
          "maximum_bounds_delta_m": 0,
          "maximum_bottom_20_percent_field": -0.6075732340289834,
          "distal_center_top_field": -0.24,
          "uv_layers_preserved": [
            "UVMap"
          ]
        },
        "replacement": {
          "bone": "Leg_L",
          "object_datablock_name": "Foot033 left conformal two-material stuffed pod",
          "mesh_datablock_name": "Foot033 left conformal two-material stuffed pod mesh",
          "material_names": [
            "Feet black velour.002",
            "Dress warm white cloth.002"
          ],
          "material_ids": [
            140052440699040,
            140052440702880
          ],
          "source_armature_modifier": {
            "name": "ReimuFumoRig",
            "object": "ReimuFumoRig",
            "settings": {
              "show_expanded": true,
              "show_in_editmode": false,
              "show_on_cage": false,
              "show_render": true,
              "show_viewport": true,
              "use_apply_on_spline": false,
              "use_bone_envelopes": false,
              "use_deform_preserve_volume": false,
              "use_multi_modifier": false,
              "use_pin_to_last": false,
              "use_vertex_groups": true,
              "vertex_group": "",
              "invert_vertex_group": false
            }
          },
          "uv_state": {
            "names": [
              "UVMap"
            ],
            "active_index": 0,
            "active_name": "UVMap",
            "active_render": [
              "UVMap"
            ]
          },
          "pose_roundtrip_max_error_m": 0,
          "evaluated_bounds_m": [
            [
              -0.05128399282693863,
              -0.022316008806228638
            ],
            [
              -0.08391998708248138,
              -0.03607999533414841
            ],
            [
              -0.0006000008434057236,
              0.02484000101685524
            ]
          ]
        }
      },
      "Foot033 right conformal two-material stuffed pod": {
        "source_name": "Right black stuffed foot pod",
        "construction": {
          "side": "right",
          "source_bounds_m": [
            [
              0.022112011909484863,
              0.0514879934489727
            ],
            [
              -0.08345998823642731,
              -0.03653999790549278
            ],
            [
              -0.0006000008434057236,
              0.02484000101685524
            ]
          ],
          "source_half_extents_m": [
            0.01468799076974392,
            0.023459995165467262,
            0.012720000930130482
          ],
          "source_polygon_count": 2016,
          "source_triangle_count": 3920,
          "output_vertex_count": 2106,
          "output_triangle_count": 4208,
          "crossing_vertex_count": 144,
          "source_signed_volume_m3": 0.000019158033665156317,
          "output_signed_volume_m3": 0.00001915804245217371,
          "winding_reversed_by_guard": false,
          "donor_topology": {
            "vertices": 1962,
            "edges": 5880,
            "faces": 3920,
            "components": 1,
            "euler": 2,
            "genus": 0,
            "non_two_incident_edges": 0
          },
          "output_topology": {
            "vertices": 2106,
            "edges": 6312,
            "faces": 4208,
            "components": 1,
            "euler": 2,
            "genus": 0,
            "non_two_incident_edges": 0
          },
          "seam_edge_count": 144,
          "seam_vertex_count": 144,
          "seam_components": 1,
          "non_degree_two_seam_vertices": 0,
          "maximum_material_side_violation": 2.948351521814274e-7,
          "maximum_bounds_delta_m": 0,
          "maximum_bottom_20_percent_field": -0.6075735850629711,
          "distal_center_top_field": -0.2400000760884935,
          "uv_layers_preserved": [
            "UVMap"
          ]
        },
        "replacement": {
          "bone": "Leg_R",
          "object_datablock_name": "Foot033 right conformal two-material stuffed pod",
          "mesh_datablock_name": "Foot033 right conformal two-material stuffed pod mesh",
          "material_names": [
            "Feet black velour.002",
            "Dress warm white cloth.002"
          ],
          "material_ids": [
            140052440699040,
            140052440702880
          ],
          "source_armature_modifier": {
            "name": "ReimuFumoRig",
            "object": "ReimuFumoRig",
            "settings": {
              "show_expanded": true,
              "show_in_editmode": false,
              "show_on_cage": false,
              "show_render": true,
              "show_viewport": true,
              "use_apply_on_spline": false,
              "use_bone_envelopes": false,
              "use_deform_preserve_volume": false,
              "use_multi_modifier": false,
              "use_pin_to_last": false,
              "use_vertex_groups": true,
              "vertex_group": "",
              "invert_vertex_group": false
            }
          },
          "uv_state": {
            "names": [
              "UVMap"
            ],
            "active_index": 0,
            "active_name": "UVMap",
            "active_render": [
              "UVMap"
            ]
          },
          "pose_roundtrip_max_error_m": 0,
          "evaluated_bounds_m": [
            [
              0.022112011909484863,
              0.0514879934489727
            ],
            [
              -0.08345998823642731,
              -0.03653999790549278
            ],
            [
              -0.0006000008434057236,
              0.02484000101685524
            ]
          ]
        }
      }
    },
    "first_hit_material_masks": {
      "front": {
        "Foot033 left conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            164,
            438,
            219,
            487
          ],
          "isolated_footprint_pixels": 1846,
          "isolated_footprint_bbox_px": [
            166,
            440,
            216,
            484
          ],
          "first_hit_counts": {
            "black": 1744,
            "cream": 102
          },
          "cream_bbox_px": [
            184,
            440,
            216,
            460
          ],
          "black_bbox_px": [
            166,
            442,
            216,
            484
          ],
          "cream_screen_horizontal_span_Wh": 0.16025571862174803,
          "black_screen_horizontal_span_fraction_of_footprint": 1,
          "cream_connected_components": 2,
          "cream_component_sizes_px": [
            101,
            1
          ],
          "cream_dominant_component_fraction": 0.9901960784313726,
          "cream_centroid_px": [
            201.87254901960785,
            445.6764705882353
          ]
        },
        "Foot033 right conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            292,
            438,
            349,
            487
          ],
          "isolated_footprint_pixels": 1868,
          "isolated_footprint_bbox_px": [
            295,
            440,
            345,
            484
          ],
          "first_hit_counts": {
            "black": 1766,
            "cream": 97,
            "occluder:Hem026 curled cotton strip": 5
          },
          "cream_bbox_px": [
            295,
            440,
            327,
            459
          ],
          "black_bbox_px": [
            295,
            442,
            345,
            484
          ],
          "cream_screen_horizontal_span_Wh": 0.16025571862174803,
          "black_screen_horizontal_span_fraction_of_footprint": 1,
          "cream_connected_components": 3,
          "cream_component_sizes_px": [
            93,
            2,
            2
          ],
          "cream_dominant_component_fraction": 0.9587628865979382,
          "cream_centroid_px": [
            309.0412371134021,
            445.8144329896907
          ]
        }
      },
      "side": {
        "Foot033 left conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            106,
            438,
            195,
            487
          ],
          "isolated_footprint_pixels": 3040,
          "isolated_footprint_bbox_px": [
            109,
            440,
            192,
            484
          ],
          "first_hit_counts": {
            "black": 43,
            "cream": 3,
            "occluder:Hem026 curled cotton strip": 51,
            "occluder:Right short hidden leg root": 106,
            "other_foot": 2837
          },
          "cream_bbox_px": [
            191,
            458,
            191,
            460
          ],
          "black_bbox_px": [
            109,
            444,
            192,
            483
          ],
          "cream_screen_horizontal_span_Wh": 0.004856233897628729,
          "black_screen_horizontal_span_fraction_of_footprint": 1,
          "cream_connected_components": 1,
          "cream_component_sizes_px": [
            3
          ],
          "cream_dominant_component_fraction": 1,
          "cream_centroid_px": [
            191,
            459
          ]
        },
        "Foot033 right conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            107,
            438,
            194,
            487
          ],
          "isolated_footprint_pixels": 2984,
          "isolated_footprint_bbox_px": [
            110,
            440,
            191,
            484
          ],
          "first_hit_counts": {
            "black": 2784,
            "cream": 53,
            "occluder:Hem026 curled cotton strip": 49,
            "occluder:Right short hidden leg root": 98
          },
          "cream_bbox_px": [
            138,
            440,
            190,
            458
          ],
          "black_bbox_px": [
            110,
            442,
            191,
            484
          ],
          "cream_screen_horizontal_span_Wh": 0.2573803965743226,
          "black_screen_horizontal_span_fraction_of_footprint": 1,
          "cream_connected_components": 4,
          "cream_component_sizes_px": [
            45,
            6,
            1,
            1
          ],
          "cream_dominant_component_fraction": 0.8490566037735849,
          "cream_centroid_px": [
            153.45283018867926,
            441.9811320754717
          ]
        }
      },
      "three_quarter": {
        "Foot033 left conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            99,
            438,
            173,
            488
          ],
          "isolated_footprint_pixels": 2503,
          "isolated_footprint_bbox_px": [
            101,
            441,
            170,
            484
          ],
          "first_hit_counts": {
            "black": 1871,
            "cream": 599,
            "occluder:Hem026 curled cotton strip": 24,
            "occluder:Left short hidden leg root": 9
          },
          "cream_bbox_px": [
            123,
            441,
            170,
            466
          ],
          "black_bbox_px": [
            101,
            443,
            169,
            484
          ],
          "cream_screen_horizontal_span_Wh": 0.23309922708617897,
          "black_screen_horizontal_span_fraction_of_footprint": 0.9857142857142858,
          "cream_connected_components": 1,
          "cream_component_sizes_px": [
            599
          ],
          "cream_dominant_component_fraction": 1,
          "cream_centroid_px": [
            149.3372287145242,
            451.9198664440735
          ]
        },
        "Foot033 right conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            190,
            439,
            264,
            488
          ],
          "isolated_footprint_pixels": 2493,
          "isolated_footprint_bbox_px": [
            193,
            441,
            261,
            485
          ],
          "first_hit_counts": {
            "black": 2431,
            "cream": 47,
            "occluder:Hem026 curled cotton strip": 6,
            "occluder:Right short hidden leg root": 9
          },
          "cream_bbox_px": [
            210,
            441,
            247,
            446
          ],
          "black_bbox_px": [
            193,
            443,
            261,
            485
          ],
          "cream_screen_horizontal_span_Wh": 0.1845368881098917,
          "black_screen_horizontal_span_fraction_of_footprint": 1,
          "cream_connected_components": 4,
          "cream_component_sizes_px": [
            43,
            2,
            1,
            1
          ],
          "cream_dominant_component_fraction": 0.9148936170212766,
          "cream_centroid_px": [
            228.72340425531914,
            442.1914893617021
          ]
        }
      },
      "rear": {
        "Foot033 left conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            293,
            438,
            348,
            487
          ],
          "isolated_footprint_pixels": 1847,
          "isolated_footprint_bbox_px": [
            295,
            440,
            345,
            484
          ],
          "first_hit_counts": {
            "black": 250,
            "occluder:Hem026 curled cotton strip": 449,
            "occluder:Skirt022 joined gathered panels": 1148
          },
          "cream_bbox_px": null,
          "black_bbox_px": [
            305,
            470,
            343,
            484
          ],
          "cream_screen_horizontal_span_Wh": 0,
          "black_screen_horizontal_span_fraction_of_footprint": 0.7647058823529411,
          "cream_connected_components": 0,
          "cream_component_sizes_px": [],
          "cream_dominant_component_fraction": 0,
          "cream_centroid_px": null
        },
        "Foot033 right conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            163,
            438,
            220,
            487
          ],
          "isolated_footprint_pixels": 1868,
          "isolated_footprint_bbox_px": [
            166,
            440,
            216,
            484
          ],
          "first_hit_counts": {
            "black": 478,
            "occluder:Hem026 curled cotton strip": 259,
            "occluder:Skirt022 joined gathered panels": 1131
          },
          "cream_bbox_px": null,
          "black_bbox_px": [
            166,
            465,
            212,
            484
          ],
          "cream_screen_horizontal_span_Wh": 0,
          "black_screen_horizontal_span_fraction_of_footprint": 0.9215686274509803,
          "cream_connected_components": 0,
          "cream_component_sizes_px": [],
          "cream_dominant_component_fraction": 0,
          "cream_centroid_px": null
        }
      },
      "three_quarter_mirror": {
        "Foot033 left conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            248,
            439,
            322,
            488
          ],
          "isolated_footprint_pixels": 2518,
          "isolated_footprint_bbox_px": [
            250,
            441,
            318,
            485
          ],
          "first_hit_counts": {
            "black": 2451,
            "cream": 54,
            "occluder:Left short hidden leg root": 13
          },
          "cream_bbox_px": [
            264,
            441,
            306,
            446
          ],
          "black_bbox_px": [
            250,
            443,
            318,
            485
          ],
          "cream_screen_horizontal_span_Wh": 0.20881805759803534,
          "black_screen_horizontal_span_fraction_of_footprint": 1,
          "cream_connected_components": 7,
          "cream_component_sizes_px": [
            45,
            3,
            2,
            1,
            1,
            1,
            1
          ],
          "cream_dominant_component_fraction": 0.8333333333333334,
          "cream_centroid_px": [
            282.35185185185185,
            442.48148148148147
          ]
        },
        "Foot033 right conformal two-material stuffed pod": {
          "projected_search_bbox_px": [
            339,
            438,
            413,
            488
          ],
          "isolated_footprint_pixels": 2476,
          "isolated_footprint_bbox_px": [
            342,
            441,
            409,
            484
          ],
          "first_hit_counts": {
            "black": 1857,
            "cream": 497,
            "occluder:Hem026 curled cotton strip": 117,
            "occluder:Right short hidden leg root": 5
          },
          "cream_bbox_px": [
            342,
            441,
            388,
            466
          ],
          "black_bbox_px": [
            342,
            443,
            409,
            484
          ],
          "cream_screen_horizontal_span_Wh": 0.22824299318855024,
          "black_screen_horizontal_span_fraction_of_footprint": 1,
          "cream_connected_components": 1,
          "cream_component_sizes_px": [
            497
          ],
          "cream_dominant_component_fraction": 1,
          "cream_centroid_px": [
            361.04426559356136,
            453.19315895372233
          ]
        }
      }
    },
    "paired_first_hit_summary": {
      "front_left_right": {
        "left": {
          "cream_first_hits": 102,
          "black_first_hits": 1744,
          "cream_fraction_of_visible_pod": 0.0552546045503792,
          "cream_components": 2,
          "dominant_component_fraction": 0.9901960784313726,
          "cream_centroid_px": [
            201.87254901960785,
            445.6764705882353
          ],
          "cream_screen_horizontal_span_Wh": 0.16025571862174803
        },
        "right": {
          "cream_first_hits": 97,
          "black_first_hits": 1766,
          "cream_fraction_of_visible_pod": 0.05206655931293613,
          "cream_components": 3,
          "dominant_component_fraction": 0.9587628865979382,
          "cream_centroid_px": [
            309.0412371134021,
            445.8144329896907
          ],
          "cream_screen_horizontal_span_Wh": 0.16025571862174803
        }
      },
      "near_foot_three_quarter_pair": {
        "camera_plus_x_right_foot": {
          "cream_first_hits": 47,
          "black_first_hits": 2431,
          "cream_fraction_of_visible_pod": 0.018966908797417272,
          "cream_components": 4,
          "dominant_component_fraction": 0.9148936170212766,
          "cream_centroid_px": [
            228.72340425531914,
            442.1914893617021
          ],
          "cream_screen_horizontal_span_Wh": 0.1845368881098917
        },
        "camera_minus_x_left_foot": {
          "cream_first_hits": 54,
          "black_first_hits": 2451,
          "cream_fraction_of_visible_pod": 0.02155688622754491,
          "cream_components": 7,
          "dominant_component_fraction": 0.8333333333333334,
          "cream_centroid_px": [
            282.35185185185185,
            442.48148148148147
          ],
          "cream_screen_horizontal_span_Wh": 0.20881805759803534
        }
      },
      "interpretation": "Raw paired evidence only; no automatic symmetry or visual pass."
    },
    "front_cream_first_hits": {
      "Foot033 left conformal two-material stuffed pod": 102,
      "Foot033 right conformal two-material stuffed pod": 97
    },
    "recorded_non_target_object_controls_unchanged": true,
    "rig_pose_unchanged": true,
    "recorded_scene_settings_unchanged": true,
    "persistent_material_ids_reused": true,
    "source_pods_preserved_and_hidden_in_active_view_layer_and_render": true,
    "active_view_layer": "ViewLayer",
    "tradeoff": "Black cap width and .06-.09 Wh cream reveal are reported, not simultaneous hard gates inside the fixed envelope.",
    "limitations": [
      "No file was saved and no image was rendered.",
      "Ray masks are pixel-center samples at the fixed 512 px views.",
      "Nonzero cream first hits do not constitute visual acceptance.",
      "The retained black ground band keeps the full pod-width mask.",
      "Animation beyond the frozen pose remains unverified.",
      "Material-node, light, and world appearance equality is not audited here; the root writer owns that guard."
    ]
  }
}
```

## Saved-file render receipt

```json
{
  "candidate": "foot_033_candidate.blend",
  "candidate_sha256": "98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48",
  "contract_sha256": "4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4",
  "baseline_receipt_sha256": "1d5c9b6fc5d432376e8c78d06fb3b4a4f1eee81b05916346b5992a5c5d858013",
  "version": "5.2.1 LTS",
  "build_hash": "9e2066aef7ef",
  "settings": {
    "engine": "BLENDER_EEVEE",
    "taa_render_samples": 16,
    "frame": 1,
    "resolution": [
      512,
      512,
      100
    ],
    "projection": "ORTHO",
    "ortho_scale_m": 0.292,
    "view_settings": {
      "view_transform": "Standard",
      "look": "None",
      "exposure": -0.10000000149011612,
      "gamma": 1
    },
    "lights": [
      {
        "name": "Key",
        "type": "AREA",
        "energy": 10,
        "color": [
          1,
          1,
          1
        ],
        "location": [
          -0.3400000035762787,
          -0.4300000071525574,
          0.41999998688697815
        ],
        "hide_render": false,
        "visible": true
      },
      {
        "name": "Fill",
        "type": "AREA",
        "energy": 5.800000190734863,
        "color": [
          1,
          1,
          1
        ],
        "location": [
          0.3799999952316284,
          -0.3199999928474426,
          0.2800000011920929
        ],
        "hide_render": false,
        "visible": true
      },
      {
        "name": "Rim",
        "type": "AREA",
        "energy": 4,
        "color": [
          1,
          1,
          1
        ],
        "location": [
          0.25999999046325684,
          0.3400000035762787,
          0.4000000059604645
        ],
        "hide_render": false,
        "visible": true
      },
      {
        "name": "Top softbox",
        "type": "AREA",
        "energy": 3,
        "color": [
          1,
          1,
          1
        ],
        "location": [
          0,
          0.019999999552965164,
          0.6200000047683716
        ],
        "hide_render": false,
        "visible": true
      }
    ],
    "world": {
      "name": "Attempt11_World.002",
      "color": [
        0.05000000074505806,
        0.05000000074505806,
        0.05000000074505806
      ],
      "background_nodes": [
        {
          "name": "Background",
          "color": [
            0.07500000298023224,
            0.0949999988079071,
            0.12999999523162842,
            1
          ],
          "strength": 0.05999999865889549
        },
        {
          "name": "Background.001",
          "color": [
            0.6899999976158142,
            0.7300000190734863,
            0.7599999904632568,
            1
          ],
          "strength": 1
        }
      ]
    },
    "film_transparent": false,
    "image_format": {
      "file_format": "PNG",
      "color_mode": "RGBA",
      "color_depth": "8",
      "compression": 15
    },
    "material_nodes_edited": false,
    "lights_edited": false,
    "world_edited": false
  },
  "renders": {
    "front": {
      "path": "candidate_front.png",
      "sha256": "18c5cfafefebcb28f9f7fc184477cdb2d448323fdef5e2b5b44d2f06269ddf86",
      "camera": {
        "camera": "Review_front_Camera",
        "location_m": [
          0,
          -0.8,
          0.13
        ],
        "rotation_euler_rad": [
          1.570796,
          0,
          0
        ]
      }
    },
    "side": {
      "path": "candidate_side.png",
      "sha256": "2f7b72e525018b8d8b9fb0705b69429a7370245ab6e134d0e2f2af6cc90a12ff",
      "camera": {
        "camera": "Review_side_Camera",
        "location_m": [
          0.8,
          0,
          0.13
        ],
        "rotation_euler_rad": [
          1.570796,
          0,
          1.570796
        ]
      }
    },
    "three_quarter": {
      "path": "candidate_three_quarter.png",
      "sha256": "ffab5f531b9fbed2a93b401502a18cf5cc22759ef501a131a645440b3b5f67e0",
      "camera": {
        "camera": "Review_three_quarter_Camera",
        "location_m": [
          0.52,
          -0.52,
          0.135
        ],
        "rotation_euler_rad": [
          1.563997,
          0,
          0.785398
        ]
      }
    },
    "rear": {
      "path": "candidate_rear.png",
      "sha256": "081dbe865f4b1ddcc8ef1e780f7583ee6a324b92ecf259b0a46c449935b8963e",
      "camera": {
        "camera": "Review_rear_Camera",
        "location_m": [
          0,
          0.8,
          0.13
        ],
        "rotation_euler_rad": [
          -1.570796,
          3.141593,
          0
        ]
      }
    },
    "three_quarter_mirror": {
      "path": "candidate_three_quarter_mirror.png",
      "sha256": "c67917cd64021c1a17403ace54ad754f99a8e085f110bac9a95c955ed275fed7",
      "camera": {
        "camera": "Review_three_quarter_mirror_Camera",
        "location_m": [
          -0.52,
          -0.52,
          0.135
        ],
        "rotation_euler_rad": [
          1.563997,
          0,
          -0.785398
        ]
      }
    }
  },
  "purpose": "Same-settings actual-node construction comparison, no final acceptance",
  "candidate_preserved": true,
  "elapsed_seconds": 7.4879088470188435
}
```
