# Reproduce the actual-material026 comparison baseline

Use retained026 source, same pinned Blender, task-local TMPDIR/cache and
factory-startup/disabled-autoexec. First run render_026_existing_shaders.py
which produces the3q diagnostic and its receipt. Restore the embedded
node_palette.json under hem_026_eevee_diagnostic; it is structural metadata,
not a material mutation. Then run render_026_eevee_baseline.py for the other
four views and byte-for-byte reuse of3q. Scripts refuse existing outputdirs.
The resulting packet is a comparison baseline, not final-material acceptance.
Blender source regeneration may change serialization hashes; rebind any
regenerated source and images to their actual identity, not old claimed hashes.

## render_026_existing_shaders.py

```python
"""One fixed-camera Eevee diagnostic, existing shader nodes and lights."""

import hashlib
import json
import time
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "hem_026_candidate.blend"
EXPECTED = "56efb16739c746153c5a562195b221645865e0ae4a6c78a5f491783b2c700882"
OUT = ROOT / "hem_026_eevee_diagnostic"
CONTRACT = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/review_contract.json"
assert hashlib.sha256(SOURCE.read_bytes()).hexdigest() == EXPECTED
OUT.mkdir(exist_ok=False)
started = time.monotonic()
bpy.context.preferences.filepaths.temporary_directory = str(ROOT / "tmp")
bpy.ops.wm.open_mainfile(filepath=str(SOURCE), load_ui=False)
bpy.context.scene.frame_set(1)
bpy.context.view_layer.update()
scene = bpy.context.scene
report = {"candidate_sha256": EXPECTED,
          "runtime": {"version": bpy.app.version_string, "build_hash": bpy.app.build_hash.decode()},
          "saved_engine": scene.render.engine, "lights": [], "materials": {},
          "world": None, "view_settings": {key: getattr(scene.view_settings, key) for key in ("view_transform", "look", "exposure", "gamma")},
          "purpose": "Distinguish Workbench viewport colors from existing material-node shading; no material/lighting edits or geometry acceptance."}
for obj in scene.objects:
    if obj.type == "LIGHT":
        report["lights"].append({"name": obj.name, "type": obj.data.type, "energy": obj.data.energy,
            "color": list(obj.data.color), "location": list(obj.matrix_world.translation),
            "hide_render": obj.hide_render, "visible": obj.visible_get()})
if scene.world:
    report["world"] = {"name": scene.world.name, "color": list(scene.world.color),
        "background_nodes": [{"name": node.name, "color": list(node.inputs["Color"].default_value),
            "strength": node.inputs["Strength"].default_value} for node in scene.world.node_tree.nodes if node.type == "BACKGROUND"] if scene.world.node_tree else []}
for name in ("Face fabric clay", "Hair brown clay", "Dress red cloth.004", "Dress warm white cloth.002", "Feet black velour.002"):
    mat = bpy.data.materials.get(name)
    if mat is None:
        report["materials"][name] = {"missing": True}
        continue
    entry = {"viewport_diffuse_color": list(mat.diffuse_color), "base_color_ancestors": []}
    pending = []
    if mat.node_tree:
        for node in mat.node_tree.nodes:
            if node.type == "BSDF_PRINCIPLED":
                pending.extend(link.from_node for link in node.inputs["Base Color"].links)
    seen = set()
    while pending:
        node = pending.pop()
        if node.name in seen:
            continue
        seen.add(node.name)
        colors = {socket.name: list(socket.default_value) for socket in node.inputs
                  if socket.type == "RGBA" and not socket.is_linked}
        entry["base_color_ancestors"].append({"name": node.name, "type": node.type,
            "unlinked_color_inputs": colors,
            "linked_inputs": {socket.name: [link.from_node.name for link in socket.links] for socket in node.inputs if socket.is_linked}})
        pending.extend(link.from_node for socket in node.inputs for link in socket.links)
    report["materials"][name] = entry
report["missing_external_images"] = [img.name for img in bpy.data.images if img.source == "FILE" and not img.packed_file and not Path(bpy.path.abspath(img.filepath)).is_file()]
report["renderable_light_count"] = sum(light["energy"] > 0 and light["visible"] and not light["hide_render"] for light in report["lights"])
(OUT / "preflight.json").write_text(json.dumps(report, indent=2) + "\n")
assert report["renderable_light_count"] > 0, "No enabled positive-energy existing light; no lighting repair authorized."
assert not report["missing_external_images"], "Missing shader images; no material repair authorized."
contract = json.loads(CONTRACT.read_text())
spec = contract["fixed_views"]["three_quarter"]
camera_data = bpy.data.cameras.new("ExistingShaderDiagnosticCamera")
camera = bpy.data.objects.new("ExistingShaderDiagnosticCamera", camera_data)
scene.collection.objects.link(camera)
camera.location = spec["location_m"]
camera.rotation_euler = spec["rotation_euler_rad"]
camera_data.type = contract["camera"]["projection"]
camera_data.ortho_scale = contract["camera"]["ortho_scale_m"]
scene.camera = camera
scene.render.engine = "BLENDER_EEVEE"
scene.render.resolution_x, scene.render.resolution_y, scene.render.resolution_percentage = contract["camera"]["resolution"]
sample_setting = None
if hasattr(scene, "eevee") and hasattr(scene.eevee, "taa_render_samples"):
    sample_setting = {"original": scene.eevee.taa_render_samples, "diagnostic": 16}
    scene.eevee.taa_render_samples = 16
scene.render.image_settings.file_format = "PNG"
scene.render.filepath = str(OUT / "existing_nodes_three_quarter.png")
report["diagnostic_settings"] = {"engine": scene.render.engine, "sample_override": sample_setting,
    "camera": spec, "ortho_scale_m": camera_data.ortho_scale, "resolution": contract["camera"]["resolution"],
    "material_nodes_edited": False, "lights_edited": False, "world_edited": False, "view_settings_edited": False}
(OUT / "preflight.json").write_text(json.dumps(report, indent=2) + "\n")
bpy.ops.render.render(write_still=True)
report["elapsed_seconds"] = time.monotonic() - started
report["render_sha256"] = hashlib.sha256(Path(scene.render.filepath).read_bytes()).hexdigest()
report["candidate_preserved"] = hashlib.sha256(SOURCE.read_bytes()).hexdigest() == EXPECTED
(OUT / "render_receipt.json").write_text(json.dumps(report, indent=2) + "\n")
print(json.dumps(report, indent=2))

```

## render_026_eevee_baseline.py

```python
"""Complete real-material baseline: four new renders plus exact reused 3q."""

import hashlib
import json
import shutil
import time
from pathlib import Path

import bpy

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "hem_026_candidate.blend"
EXPECTED = "56efb16739c746153c5a562195b221645865e0ae4a6c78a5f491783b2c700882"
CONTRACT = ROOT.parents[2] / "projects/renders/assets/reimu_fumo/review_contract.json"
PRIOR = ROOT / "hem_026_eevee_diagnostic"
OUT = ROOT / "hem_026_eevee_review"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


assert sha(SOURCE) == EXPECTED
contract = json.loads(CONTRACT.read_text())
prior = json.loads((PRIOR / "render_receipt.json").read_text())
assert prior["candidate_sha256"] == EXPECTED
assert len(contract["fixed_views"]) == 5 and "three_quarter" in contract["fixed_views"]
assert prior["diagnostic_settings"]["camera"] == contract["fixed_views"]["three_quarter"]
prior_image = PRIOR / "existing_nodes_three_quarter.png"
assert sha(prior_image) == prior["render_sha256"] == "f204a37970f86502f222cd889dab070ad7587e78f55390e4696c17b2bc43c83b"
OUT.mkdir(exist_ok=False)
started = time.monotonic()
bpy.context.preferences.filepaths.temporary_directory = str(ROOT / "tmp")
bpy.ops.wm.open_mainfile(filepath=str(SOURCE), load_ui=False)
bpy.context.scene.frame_set(1)
bpy.context.view_layer.update()
scene = bpy.context.scene
scene.render.engine = "BLENDER_EEVEE"
scene.eevee.taa_render_samples = 16
scene.render.resolution_x, scene.render.resolution_y, scene.render.resolution_percentage = contract["camera"]["resolution"]
scene.render.image_settings.file_format = "PNG"
assert {key: getattr(scene.view_settings, key) for key in prior["view_settings"]} == prior["view_settings"]
manifest = {
    "purpose": "Frozen 026 actual-node five-view comparison baseline; no material or geometric acceptance",
    "candidate": str(SOURCE), "candidate_sha256": EXPECTED,
    "contract": str(CONTRACT), "contract_sha256": sha(CONTRACT),
    "runtime": {"version": bpy.app.version_string, "build_hash": bpy.app.build_hash.decode(), "background": bpy.app.background},
    "settings": {
        "engine": scene.render.engine, "taa_render_samples": scene.eevee.taa_render_samples,
        "frame": scene.frame_current, "resolution": contract["camera"]["resolution"],
        "projection": contract["camera"]["projection"], "ortho_scale_m": contract["camera"]["ortho_scale_m"],
        "view_settings": prior["view_settings"], "lights": prior["lights"], "world": prior["world"],
        "film_transparent": scene.render.film_transparent,
        "image_format": {key: getattr(scene.render.image_settings, key) for key in ("file_format", "color_mode", "color_depth", "compression")},
        "material_nodes_edited": False, "lights_edited": False, "world_edited": False,
    },
    "shader_inspection": str(PRIOR / "node_palette.json"),
    "shader_inspection_sha256": sha(PRIOR / "node_palette.json"),
    "renders": {},
}
copy_path = OUT / "candidate_three_quarter.png"
shutil.copyfile(prior_image, copy_path)
assert sha(copy_path) == prior["render_sha256"]
manifest["renders"]["three_quarter"] = {"path": copy_path.name, "sha256": sha(copy_path),
    "camera": contract["fixed_views"]["three_quarter"], "source": "exact copy of completed diagnostic, not rerendered", "copied_from": str(prior_image)}
for view, spec in contract["fixed_views"].items():
    if view == "three_quarter":
        continue
    camera_data = bpy.data.cameras.new("ExistingShaderBaseline_" + view)
    camera = bpy.data.objects.new("ExistingShaderBaseline_" + view, camera_data)
    scene.collection.objects.link(camera)
    camera.location = spec["location_m"]
    camera.rotation_euler = spec["rotation_euler_rad"]
    camera_data.type = contract["camera"]["projection"]
    camera_data.ortho_scale = contract["camera"]["ortho_scale_m"]
    scene.camera = camera
    path = OUT / ("candidate_" + view + ".png")
    scene.render.filepath = str(path)
    bpy.ops.render.render(write_still=True)
    manifest["renders"][view] = {"path": path.name, "sha256": sha(path), "camera": spec, "source": "new render with frozen existing nodes/lights"}
    (OUT / "render_receipt.json").write_text(json.dumps(manifest, indent=2) + "\n")
assert len(manifest["renders"]) == 5
manifest["elapsed_seconds_four_renders"] = time.monotonic() - started
manifest["candidate_preserved"] = sha(SOURCE) == EXPECTED
manifest["copied_three_quarter_preserved"] = sha(copy_path) == prior["render_sha256"]
(OUT / "render_receipt.json").write_text(json.dumps(manifest, indent=2) + "\n")
print(json.dumps(manifest, indent=2))

```

## hem_026_eevee_diagnostic/node_palette.json

```json
{
  "candidate_sha256": "56efb16739c746153c5a562195b221645865e0ae4a6c78a5f491783b2c700882",
  "materials": {
    "Face fabric clay": [
      {
        "name": "Color Ramp",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.18000000715255737,
            "rgba": [
              0.25,
              0.17000000178813934,
              0.10499999672174454,
              1.0
            ]
          },
          {
            "position": 0.8199999928474426,
            "rgba": [
              0.30000001192092896,
              0.20499999821186066,
              0.12800000607967377,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Color Ramp.001",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.11999999731779099,
            "rgba": [
              0.9399999976158142,
              0.9399999976158142,
              0.9399999976158142,
              1.0
            ]
          },
          {
            "position": 0.8799999952316284,
            "rgba": [
              1.024999976158142,
              1.024999976158142,
              1.024999976158142,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Mix (Legacy)",
        "type": "MIX_RGB",
        "blend_type": "MULTIPLY",
        "factor_default": 1.0,
        "factor_linked": false
      }
    ],
    "Hair brown clay": [
      {
        "name": "Color Ramp",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.18000000715255737,
            "rgba": [
              0.040049999952316284,
              0.0106800002977252,
              0.004449999891221523,
              1.0
            ]
          },
          {
            "position": 0.8199999928474426,
            "rgba": [
              0.04859999939799309,
              0.012959999963641167,
              0.005400000140070915,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Color Ramp.001",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.11999999731779099,
            "rgba": [
              0.9399999976158142,
              0.9399999976158142,
              0.9399999976158142,
              1.0
            ]
          },
          {
            "position": 0.8799999952316284,
            "rgba": [
              1.024999976158142,
              1.024999976158142,
              1.024999976158142,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Mix (Legacy)",
        "type": "MIX_RGB",
        "blend_type": "MULTIPLY",
        "factor_default": 1.0,
        "factor_linked": false
      }
    ],
    "Dress red cloth.004": [
      {
        "name": "Color Ramp",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.18000000715255737,
            "rgba": [
              0.13795000314712524,
              0.0022249999456107616,
              0.004449999891221523,
              1.0
            ]
          },
          {
            "position": 0.8199999928474426,
            "rgba": [
              0.16740000247955322,
              0.0027000000700354576,
              0.005400000140070915,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Color Ramp.001",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.11999999731779099,
            "rgba": [
              0.9399999976158142,
              0.9399999976158142,
              0.9399999976158142,
              1.0
            ]
          },
          {
            "position": 0.8799999952316284,
            "rgba": [
              1.024999976158142,
              1.024999976158142,
              1.024999976158142,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Mix (Legacy)",
        "type": "MIX_RGB",
        "blend_type": "MULTIPLY",
        "factor_default": 1.0,
        "factor_linked": false
      },
      {
        "name": "Mix (Legacy).001",
        "type": "MIX_RGB",
        "blend_type": "MULTIPLY",
        "factor_default": 1.0,
        "factor_linked": false
      }
    ],
    "Dress warm white cloth.002": [
      {
        "name": "Color Ramp",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.18000000715255737,
            "rgba": [
              0.33000001311302185,
              0.3199999928474426,
              0.30000001192092896,
              1.0
            ]
          },
          {
            "position": 0.8199999928474426,
            "rgba": [
              0.4000000059604645,
              0.3880000114440918,
              0.36500000953674316,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Color Ramp.001",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.11999999731779099,
            "rgba": [
              0.9399999976158142,
              0.9399999976158142,
              0.9399999976158142,
              1.0
            ]
          },
          {
            "position": 0.8799999952316284,
            "rgba": [
              1.024999976158142,
              1.024999976158142,
              1.024999976158142,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Mix (Legacy)",
        "type": "MIX_RGB",
        "blend_type": "MULTIPLY",
        "factor_default": 1.0,
        "factor_linked": false
      },
      {
        "name": "Mix (Legacy).001",
        "type": "MIX_RGB",
        "blend_type": "MULTIPLY",
        "factor_default": 1.0,
        "factor_linked": false
      }
    ],
    "Feet black velour.002": [
      {
        "name": "Color Ramp",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.18000000715255737,
            "rgba": [
              0.0026700000744313,
              0.001601999974809587,
              0.001601999974809587,
              1.0
            ]
          },
          {
            "position": 0.8199999928474426,
            "rgba": [
              0.0032399999909102917,
              0.0019440000178292394,
              0.0019440000178292394,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Color Ramp.001",
        "type": "VALTORGB",
        "interpolation": "LINEAR",
        "elements": [
          {
            "position": 0.11999999731779099,
            "rgba": [
              0.9399999976158142,
              0.9399999976158142,
              0.9399999976158142,
              1.0
            ]
          },
          {
            "position": 0.8799999952316284,
            "rgba": [
              1.024999976158142,
              1.024999976158142,
              1.024999976158142,
              1.0
            ]
          }
        ]
      },
      {
        "name": "Mix (Legacy)",
        "type": "MIX_RGB",
        "blend_type": "MULTIPLY",
        "factor_default": 1.0,
        "factor_linked": false
      }
    ]
  },
  "candidate_preserved": true
}

```


