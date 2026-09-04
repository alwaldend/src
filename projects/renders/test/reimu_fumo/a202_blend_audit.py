"""Fail closed when the sanitized A202 Blender donor leaks local paths."""

import hashlib
import json
import math
import os
import re
import sys
from collections.abc import Mapping, Sequence

import bpy


EXPECTED_ARMATURE = "ReimuFumoRig"
EXPECTED_VIEW_CAMERAS = {
    "front": "Review_front_Camera",
    "rear": "Review_rear_Camera",
    "side": "Review_side_Camera",
    "three_quarter": "Review_three_quarter_Camera",
    "three_quarter_mirror": "Review_three_quarter_mirror_Camera",
}
EXPECTED_CAMERAS = set(EXPECTED_VIEW_CAMERAS.values())
FORBIDDEN_FRAGMENTS = ("/var/", "/home/", "simeon", "t3code-")
POSIX_ABSOLUTE = re.compile(
    r"(?<![A-Za-z0-9_./-])/(?:[A-Za-z0-9_.~-]+/)*[A-Za-z0-9_.~-]+"
)
WINDOWS_ABSOLUTE = re.compile(r"(?<![A-Za-z0-9_])[A-Za-z]:[\\/]")
UNC_ABSOLUTE = re.compile(r"(?<!\\)\\\\[^\\\s]+[\\/]?")
FILE_URI = re.compile(r"\bfile:(?://)?", re.IGNORECASE)
# These UI properties are populated from the test host after a file opens;
# they are not serialized asset values. All other reachable RNA strings remain
# in scope.
RUNTIME_ONLY_RNA_PROPERTIES = {
    "selected_studio_light",
    "system_bookmarks",
    "system_folders",
}

failures = []
privacy_findings = set()
visited_rna = set()


def fail(message):
    failures.append(message)


def fail_scan(location, error):
    fail(f"could not inspect {location}: {type(error).__name__}")


def file_sha256(filename):
    digest = hashlib.sha256()
    with open(filename, "rb") as source:
        for chunk in iter(lambda: source.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def classify_private_string(value):
    normalized = value.replace("\\", "/")
    lowered = normalized.casefold()
    reasons = []
    if any(fragment in lowered for fragment in FORBIDDEN_FRAGMENTS):
        reasons.append("forbidden local-path or identity marker")
    absolute = (
        value == "/"
        or POSIX_ABSOLUTE.search(normalized)
        or WINDOWS_ABSOLUTE.search(value)
        or UNC_ABSOLUTE.search(value)
        or FILE_URI.search(value)
    )
    if absolute:
        reasons.append("absolute path")
    return reasons


def inspect_string(location, value):
    if not isinstance(value, str) or not value:
        return
    value_digest = hashlib.sha256(value.encode("utf-8")).hexdigest()[:16]
    for reason in classify_private_string(value):
        privacy_findings.add(
            f"{location}: {reason}; value_sha256={value_digest}"
        )


def inspect_custom_value(location, value, depth=0):
    if isinstance(value, str):
        inspect_string(location, value)
        return
    if depth >= 32 or value is None:
        if depth >= 32:
            fail(f"custom-property nesting exceeds limit at {location}")
        return
    if isinstance(value, bpy.types.ID):
        return
    if isinstance(value, Mapping) or callable(getattr(value, "keys", None)):
        try:
            keys = list(value.keys())
        except Exception as error:
            fail_scan(location, error)
            return
        for index, key in enumerate(keys):
            inspect_string(f"{location}.key[{index}]", key)
            try:
                child = value[key]
            except Exception as error:
                fail_scan(f"{location}.value[{index}]", error)
                continue
            inspect_custom_value(
                f"{location}.value[{index}]", child, depth + 1
            )
        return
    if isinstance(value, Sequence) and not isinstance(
        value, (bytes, bytearray)
    ):
        try:
            for index, child in enumerate(value):
                inspect_custom_value(f"{location}[{index}]", child, depth + 1)
        except Exception as error:
            fail_scan(location, error)


def inspect_custom_properties(owner, location):
    keys_method = getattr(owner, "keys", None)
    if not callable(keys_method):
        return
    try:
        keys = list(keys_method())
    except TypeError:
        # Every bpy_struct exposes keys(), including types that cannot own
        # ID properties. Blender reports that case with TypeError.
        return
    except Exception as error:
        fail_scan(location, error)
        return
    for index, key in enumerate(keys):
        inspect_string(f"{location}.custom_key[{index}]", key)
        try:
            value = owner[key]
        except Exception as error:
            fail_scan(f"{location}.custom_value[{index}]", error)
            continue
        inspect_custom_value(f"{location}.custom_value[{index}]", value)
        metadata_method = getattr(owner, "id_properties_ui", None)
        if not callable(metadata_method):
            continue
        try:
            metadata = metadata_method(key).as_dict()
        except TypeError:
            # Groups and a few other IDProperty kinds have no UI metadata.
            continue
        except Exception as error:
            fail_scan(f"{location}.custom_ui[{index}]", error)
            continue
        inspect_custom_value(f"{location}.custom_ui[{index}]", metadata)


def rna_identity(owner):
    try:
        return owner.bl_rna.identifier, owner.as_pointer()
    except (AttributeError, ReferenceError, RuntimeError, TypeError):
        return type(owner).__name__, id(owner)


def inspect_rna(owner, location):
    if owner is None or not hasattr(owner, "bl_rna"):
        return
    identity = rna_identity(owner)
    if identity in visited_rna:
        return
    visited_rna.add(identity)
    inspect_custom_properties(owner, location)
    try:
        definitions = list(owner.bl_rna.properties)
    except Exception as error:
        fail_scan(location, error)
        return
    for definition in definitions:
        identifier = definition.identifier
        if identifier == "rna_type":
            continue
        if identifier in RUNTIME_ONLY_RNA_PROPERTIES:
            continue
        if definition.type not in {"COLLECTION", "POINTER", "STRING"}:
            continue
        child_location = f"{location}.{identifier}"
        try:
            value = getattr(owner, identifier)
        except Exception as error:
            fail_scan(child_location, error)
            continue
        if definition.type == "STRING":
            inspect_string(child_location, value)
        elif definition.type == "POINTER":
            inspect_rna(value, child_location)
        elif definition.type == "COLLECTION":
            try:
                for index, child in enumerate(value):
                    inspect_rna(child, f"{child_location}[{index}]")
            except Exception as error:
                fail_scan(child_location, error)
                continue


def data_blocks():
    blocks = []
    for definition in bpy.data.bl_rna.properties:
        if definition.type != "COLLECTION":
            continue
        try:
            collection = getattr(bpy.data, definition.identifier)
        except Exception as error:
            fail_scan(f"bpy.data.{definition.identifier}", error)
            continue
        try:
            items = list(collection)
        except Exception as error:
            fail_scan(f"bpy.data.{definition.identifier}", error)
            continue
        for index, block in enumerate(items):
            if isinstance(block, bpy.types.ID):
                location = f"bpy.data.{definition.identifier}[{index}]"
                blocks.append((location, block))
    return blocks


def inspect_known_path_surfaces():
    for scene_index, scene in enumerate(bpy.data.scenes):
        inspect_string(
            f"bpy.data.scenes[{scene_index}].render.filepath",
            scene.render.filepath,
        )

    path_collections = {
        "cache_files": ("filepath",),
        "fonts": ("filepath",),
        "images": ("filepath", "filepath_raw"),
        "libraries": ("filepath",),
        "movieclips": ("filepath",),
        "sounds": ("filepath",),
        "volumes": ("filepath",),
    }
    for collection_name, property_names in path_collections.items():
        collection = getattr(bpy.data, collection_name, ())
        for item_index, item in enumerate(collection):
            for property_name in property_names:
                value = getattr(item, property_name, "")
                inspect_string(
                    f"bpy.data.{collection_name}"
                    f"[{item_index}].{property_name}",
                    value,
                )

    for block_index, (_, block) in enumerate(data_blocks()):
        weak_reference = getattr(block, "library_weak_reference", None)
        if weak_reference is not None:
            inspect_string(
                f"datablock[{block_index}].library_weak_reference.filepath",
                getattr(weak_reference, "filepath", ""),
            )


def node_trees():
    trees = []
    seen = set()

    def append(tree):
        if tree is None:
            return
        identity = rna_identity(tree)
        if identity not in seen:
            seen.add(identity)
            trees.append(tree)

    for tree in bpy.data.node_groups:
        append(tree)
    for collection in (
        bpy.data.materials,
        bpy.data.worlds,
        bpy.data.lights,
        bpy.data.scenes,
    ):
        for owner in collection:
            append(getattr(owner, "node_tree", None))
    return trees


def inspect_file_output_nodes():
    for tree_index, tree in enumerate(node_trees()):
        for node_index, node in enumerate(tree.nodes):
            if (
                getattr(node, "type", "") != "OUTPUT_FILE"
                and getattr(node, "bl_idname", "")
                != "CompositorNodeOutputFile"
            ):
                continue
            prefix = f"node_trees[{tree_index}].nodes[{node_index}]"
            inspect_string(
                f"{prefix}.base_path", getattr(node, "base_path", "")
            )
            for slot_index, slot in enumerate(getattr(node, "file_slots", ())):
                inspect_string(
                    f"{prefix}.file_slots[{slot_index}].path",
                    getattr(slot, "path", ""),
                )


def inspect_texts():
    for index, text in enumerate(bpy.data.texts):
        inspect_string(f"bpy.data.texts[{index}].name", text.name)
        inspect_string(f"bpy.data.texts[{index}].body", text.as_string())


def inspect_external_dependencies():
    try:
        paths = list(
            bpy.utils.blend_paths(absolute=True, packed=False, local=False)
        )
    except (AttributeError, ReferenceError, RuntimeError, TypeError) as error:
        fail(
            "could not enumerate external dependencies: "
            f"{type(error).__name__}"
        )
        return
    missing = []
    for external_path in paths:
        if external_path and not os.path.exists(external_path):
            missing.append(
                hashlib.sha256(external_path.encode("utf-8")).hexdigest()[:16]
            )
    if missing:
        fail("missing external files: " + ", ".join(sorted(set(missing))))


def close_enough(actual, expected):
    return math.isclose(
        float(actual), float(expected), rel_tol=0.0, abs_tol=1.0e-6
    )


def inspect_review_contract(contract):
    if contract.get("schema_version") != 1:
        fail("review contract schema_version is not 1")
    fixed_views = contract.get("fixed_views")
    if not isinstance(fixed_views, dict):
        fail("review contract fixed_views is not an object")
        return
    actual_views = set(fixed_views)
    expected_views = set(EXPECTED_VIEW_CAMERAS)
    if actual_views != expected_views:
        fail("review contract does not declare the exact five fixed views")

    camera_contract = contract.get("camera")
    if not isinstance(camera_contract, dict):
        fail("review contract camera is not an object")
        return
    projection = camera_contract.get("projection")
    if projection != "ORTHO":
        fail("review contract projection is not ORTHO")
    expected_scale = camera_contract.get("ortho_scale_m")
    expected_resolution = camera_contract.get("resolution")
    if (
        not isinstance(expected_resolution, list)
        or len(expected_resolution) != 3
    ):
        fail("review contract resolution is not a three-item list")
    else:
        scene = bpy.context.scene
        actual_resolution = (
            scene.render.resolution_x,
            scene.render.resolution_y,
            scene.render.resolution_percentage,
        )
        if tuple(expected_resolution) != actual_resolution:
            fail("scene render resolution differs from review contract")

    contract_camera_names = {
        view.get("camera")
        for view in fixed_views.values()
        if isinstance(view, dict)
    }
    if contract_camera_names != EXPECTED_CAMERAS:
        fail("review contract does not declare the exact five camera names")

    for view_name, camera_name in EXPECTED_VIEW_CAMERAS.items():
        view_contract = fixed_views.get(view_name)
        if not isinstance(view_contract, dict):
            fail(f"review contract view is absent: {view_name}")
            continue
        if view_contract.get("camera") != camera_name:
            fail(f"review contract camera mismatch for view: {view_name}")
            continue
        camera = bpy.data.objects.get(camera_name)
        if camera is None or camera.type != "CAMERA" or camera.data is None:
            fail(f"required camera is absent: {camera_name}")
            continue
        if camera.data.type != projection:
            fail(f"camera projection differs for view: {view_name}")
        if expected_scale is None or not close_enough(
            camera.data.ortho_scale, expected_scale
        ):
            fail(f"camera ortho scale differs for view: {view_name}")
        for property_name, actual in (
            ("location_m", tuple(camera.location)),
            ("rotation_euler_rad", tuple(camera.rotation_euler)),
        ):
            expected = view_contract.get(property_name)
            if not isinstance(expected, list) or len(expected) != 3:
                fail(
                    f"review contract {property_name} is invalid for: "
                    f"{view_name}"
                )
                continue
            if any(
                not close_enough(component, expected_component)
                for component, expected_component in zip(actual, expected)
            ):
                fail(f"camera {property_name} differs for view: {view_name}")


def main():
    if "--" not in sys.argv:
        raise RuntimeError("expected model argument after --")
    arguments = sys.argv[sys.argv.index("--") + 1 :]
    if len(arguments) != 2:
        raise RuntimeError("expected model and review-contract arguments")
    expected_model = os.path.realpath(arguments[0])
    review_contract_path = os.path.realpath(arguments[1])
    loaded_model = os.path.realpath(bpy.data.filepath)
    if not os.path.isfile(expected_model):
        fail("expected model is not a regular file")
    elif not os.path.samefile(expected_model, loaded_model):
        fail("opened blend does not match the expected model")
    model_digest_before = file_sha256(expected_model)
    with open(review_contract_path, encoding="utf-8") as source:
        review_contract = json.load(source)

    if not bpy.app.background:
        fail("Blender is not running in background mode")
    if bpy.app.version[:2] != (5, 2):
        fail(f"Blender version is {bpy.app.version_string}, want 5.2.x")
    if not bpy.data.is_saved:
        fail("Blender did not open a saved blend")
    if bpy.data.is_dirty:
        fail("blend is dirty immediately after opening")
    if len(bpy.data.libraries) != 0:
        fail(f"linked library count is {len(bpy.data.libraries)}, want 0")

    armature_objects = [
        obj for obj in bpy.data.objects if obj.type == "ARMATURE"
    ]
    rig = bpy.data.objects.get(EXPECTED_ARMATURE)
    if len(armature_objects) != 1 or rig is None or rig.type != "ARMATURE":
        fail("expected sole ReimuFumoRig armature is absent")
    if len(bpy.data.armatures) != 1:
        fail(f"armature datablock count is {len(bpy.data.armatures)}, want 1")
    inspect_review_contract(review_contract)

    blocks = data_blocks()
    for location, block in blocks:
        inspect_rna(block, location)
    inspect_known_path_surfaces()
    inspect_file_output_nodes()
    inspect_texts()
    inspect_external_dependencies()

    if privacy_findings:
        findings = sorted(privacy_findings)
        failures.append(
            f"private path findings: {len(findings)}; "
            + " | ".join(findings[:50])
        )
    if bpy.data.is_dirty:
        fail("audit dirtied the blend")
    if file_sha256(expected_model) != model_digest_before:
        fail("audit changed the model bytes")
    if failures:
        for failure in failures:
            print(f"AUDIT FAILURE: {failure}", file=sys.stderr)
        raise RuntimeError(
            f"A202 blend audit failed with {len(failures)} errors"
        )

    print(
        json.dumps(
            {
                "armature": EXPECTED_ARMATURE,
                "background": True,
                "blender": bpy.app.version_string,
                "cameras_verified": len(EXPECTED_CAMERAS),
                "datablocks_scanned": len(blocks),
                "linked_libraries": 0,
                "missing_external_files": 0,
                "privacy_findings": 0,
                "saved": False,
            },
            sort_keys=True,
        )
    )


if __name__ == "__main__":
    main()
