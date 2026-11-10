"""Build the editable army-spot remake with baked, script-free animation.

Run through the owning Blender target. The reference is used only for identity
validation; delivered geometry, materials, and animation do not depend on it.
"""

import argparse
import hashlib
import json
import math
import pathlib
import random
import sys

import bpy
import mathutils

_REFERENCE_SHA256: str = (
    "4aa8481103870d56131063ed62438f6999e3d0085a5826a647a03b6eb6e7917d"
)
"""Digest of the immutable user-supplied source clip."""


def _material(
    name: str, color: tuple[float, float, float], roughness: float = 0.7
) -> bpy.types.Material:
    """Create an opaque, deliberately matte surface material."""
    material: bpy.types.Material = bpy.data.materials.new(name)
    material.diffuse_color = (*color, 1)
    material.use_nodes = True
    shader: bpy.types.Node = material.node_tree.nodes.get("Principled BSDF")
    shader.inputs["Base Color"].default_value = (*color, 1)
    shader.inputs["Roughness"].default_value = roughness
    return material


def _camouflage(name: str, desert: bool) -> bpy.types.Material:
    """Construct deterministic small-scale digital camouflage without images."""
    palette: list[tuple[float, float, float]] = (
        [
            (0.42, 0.32, 0.20),
            (0.20, 0.13, 0.07),
            (0.58, 0.47, 0.32),
            (0.32, 0.23, 0.13),
        ]
        if desert
        else [
            (0.25, 0.26, 0.10),
            (0.07, 0.09, 0.035),
            (0.40, 0.33, 0.14),
            (0.15, 0.18, 0.065),
        ]
    )
    material: bpy.types.Material = _material(name, palette[0])
    nodes: bpy.types.Nodes = material.node_tree.nodes
    links: bpy.types.NodeLinks = material.node_tree.links
    coordinates: bpy.types.Node = nodes.new("ShaderNodeNewGeometry")
    scale: bpy.types.Node = nodes.new("ShaderNodeVectorMath")
    scale.operation = "SCALE"
    scale.inputs[3].default_value = 65
    quantize: bpy.types.Node = nodes.new("ShaderNodeVectorMath")
    quantize.operation = "FLOOR"
    cluster: bpy.types.Node = nodes.new("ShaderNodeVectorMath")
    cluster.operation = "SCALE"
    cluster.inputs[3].default_value = 0.16
    noise: bpy.types.Node = nodes.new("ShaderNodeTexNoise")
    noise.inputs["Scale"].default_value = 1.8
    noise.inputs["Detail"].default_value = 0
    ramp: bpy.types.Node = nodes.new("ShaderNodeValToRGB")
    ramp.color_ramp.interpolation = "CONSTANT"
    for index, color in enumerate(palette):
        element: bpy.types.ColorRampElement = (
            ramp.color_ramp.elements[index]
            if index < 2
            else ramp.color_ramp.elements.new(index / 4)
        )
        element.position = (0.0, 0.43, 0.52, 0.60)[index]
        element.color = (*color, 1)
    links.new(coordinates.outputs["Position"], scale.inputs[0])
    links.new(scale.outputs[0], quantize.inputs[0])
    links.new(quantize.outputs[0], cluster.inputs[0])
    links.new(cluster.outputs[0], noise.inputs["Vector"])
    links.new(noise.outputs["Fac"], ramp.inputs[0])
    links.new(
        ramp.outputs[0], nodes.get("Principled BSDF").inputs["Base Color"]
    )
    return material


def _attach(
    obj: bpy.types.Object, parent: bpy.types.Object | None
) -> bpy.types.Object:
    """Parent an object in its authored local coordinates."""
    obj.parent = parent
    return obj


def _empty(
    name: str, parent: bpy.types.Object | None = None
) -> bpy.types.Object:
    """Create an editable transform pivot for ordinary keyframe animation."""
    obj: bpy.types.Object = bpy.data.objects.new(name, None)
    bpy.context.collection.objects.link(obj)
    return _attach(obj, parent)


def _sphere(
    name: str,
    location: tuple[float, float, float],
    scale: tuple[float, float, float],
    material: bpy.types.Material,
    parent: bpy.types.Object | None = None,
) -> bpy.types.Object:
    """Create a smooth elliptical organic volume."""
    bpy.ops.mesh.primitive_uv_sphere_add(segments=32, ring_count=20)
    obj: bpy.types.Object = bpy.context.object
    obj.name = name
    obj.location = location
    obj.scale = scale
    obj.data.materials.append(material)
    for polygon in obj.data.polygons:
        polygon.use_smooth = True
    return _attach(obj, parent)


def _box(
    name: str,
    location: tuple[float, float, float],
    scale: tuple[float, float, float],
    material: bpy.types.Material,
    parent: bpy.types.Object | None = None,
    bevel: float = 0.025,
) -> bpy.types.Object:
    """Make a rounded rectangular garment or hard-surface part."""
    bpy.ops.mesh.primitive_cube_add(size=1)
    obj: bpy.types.Object = bpy.context.object
    obj.name = name
    obj.location = location
    obj.scale = scale
    bpy.ops.object.transform_apply(location=False, rotation=False, scale=True)
    obj.data.materials.append(material)
    if bevel:
        modifier: bpy.types.BevelModifier = obj.modifiers.new(
            "Rounded edges", "BEVEL"
        )
        modifier.width = bevel
        modifier.segments = 3
        obj.modifiers.new("Surface normals", "WEIGHTED_NORMAL")
    return _attach(obj, parent)


def _mesh(
    name: str,
    vertices: list[tuple[float, float, float]],
    faces: list[tuple[int, ...]],
    material: bpy.types.Material,
    parent: bpy.types.Object | None = None,
) -> bpy.types.Object:
    """Create editable mesh topology from local positions and face indices."""
    mesh: bpy.types.Mesh = bpy.data.meshes.new(name)
    mesh.from_pydata(vertices, [], faces)
    mesh.update()
    obj: bpy.types.Object = bpy.data.objects.new(name, mesh)
    bpy.context.collection.objects.link(obj)
    obj.data.materials.append(material)
    for polygon in mesh.polygons:
        polygon.use_smooth = True
    return _attach(obj, parent)


def _curve(
    name: str,
    points: list[tuple[float, float, float]],
    radius: float,
    material: bpy.types.Material,
    parent: bpy.types.Object | None = None,
) -> bpy.types.Object:
    """Form a smooth narrow seam, eyebrow, or rim from a beveled path."""
    data: bpy.types.Curve = bpy.data.curves.new(name, "CURVE")
    data.dimensions = "3D"
    data.bevel_depth = radius
    data.bevel_resolution = 3
    spline: bpy.types.Spline = data.splines.new("BEZIER")
    spline.bezier_points.add(len(points) - 1)
    for point, coordinate in zip(spline.bezier_points, points):
        point.co = coordinate
        point.handle_left_type = "AUTO"
        point.handle_right_type = "AUTO"
    obj: bpy.types.Object = bpy.data.objects.new(name, data)
    bpy.context.collection.objects.link(obj)
    data.materials.append(material)
    return _attach(obj, parent)


def _segment(
    name: str,
    start: tuple[float, float, float],
    end: tuple[float, float, float],
    width: float,
    material: bpy.types.Material,
    parent: bpy.types.Object,
) -> bpy.types.Object:
    """Join two body landmarks with a tapered smooth limb volume."""
    vector: mathutils.Vector = mathutils.Vector(end) - mathutils.Vector(start)
    center: mathutils.Vector = (
        mathutils.Vector(start) + mathutils.Vector(end)
    ) / 2
    vertices: list[tuple[float, float, float]] = []
    faces: list[tuple[int, ...]] = []
    profile: tuple[float, ...] = (0.95, 1.0, 1.0, 0.98, 0.96, 0.94)
    for ring, radius in enumerate(profile):
        for column in range(24):
            angle: float = 2 * math.pi * column / 24
            vertices.append(
                (
                    width * radius * math.cos(angle),
                    width * radius * 0.9 * math.sin(angle),
                    vector.length * (ring / 5 - 0.5),
                )
            )
    for ring in range(5):
        for column in range(24):
            a: int = ring * 24 + column
            b: int = ring * 24 + (column + 1) % 24
            faces.append((a, b, b + 24, a + 24))
    faces.extend((tuple(reversed(range(24))), tuple(range(120, 144))))
    obj: bpy.types.Object = _mesh(name, vertices, faces, material, parent)
    obj.location = center
    obj.rotation_euler = vector.to_track_quat("Z", "Y").to_euler()
    return obj


def _head_surface(parent: bpy.types.Object, skin: bpy.types.Material) -> None:
    """Model a continuous tapered face with integrated nose and cheek volumes."""
    vertices: list[tuple[float, float, float]] = []
    faces: list[tuple[int, ...]] = []
    rings: int = 40
    columns: int = 80
    for row in range(rings + 1):
        latitude: float = -math.pi / 2 + math.pi * row / rings
        z: float = 0.38 * math.sin(latitude)
        jaw: float = 0.77 + 0.23 * (z + 0.38) / 0.76
        for column in range(columns):
            angle: float = 2 * math.pi * column / columns
            x: float = 0.255 * math.cos(latitude) * math.cos(angle) * jaw
            y: float = 0.205 * math.cos(latitude) * math.sin(angle)
            if y < 0:
                front: float = (-math.sin(angle)) ** 12
                nose: float = 0.17 * math.exp(
                    -((x / 0.046) ** 2) - (((z + 0.015) / 0.145) ** 4)
                )
                cheeks: float = 0.024 * math.exp(
                    -(((abs(x) - 0.13) / 0.06) ** 2)
                    - (((z + 0.09) / 0.13) ** 2)
                )
                y -= front * (nose + cheeks)
            vertices.append((x, y, z))
    for row in range(rings):
        for column in range(columns):
            a: int = row * columns + column
            b: int = row * columns + (column + 1) % columns
            faces.append((a, b, b + columns, a + columns))
    _mesh("Continuous face and nose", vertices, faces, skin, parent)


def _helmet(
    parent: bpy.types.Object,
    material: bpy.types.Material,
    officer: bool,
    red: bpy.types.Material,
    gold: bpy.types.Material,
    dark: bpy.types.Material,
) -> None:
    """Build a hollow rolled-rim helmet or an officer's peaked service cap."""
    vertices: list[tuple[float, float, float]] = []
    faces: list[tuple[int, ...]] = []
    for row in range(17):
        latitude: float = (math.pi / 2) * row / 16
        for column in range(64):
            angle: float = 2 * math.pi * column / 64
            radius: float = math.cos(latitude)
            vertices.append(
                (
                    0.30 * radius * math.cos(angle),
                    0.255 * radius * math.sin(angle),
                    0.19 + 0.31 * math.sin(latitude),
                )
            )
    for row in range(16):
        for column in range(64):
            a: int = row * 64 + column
            b: int = row * 64 + (column + 1) % 64
            faces.append((a, b, b + 64, a + 64))
    shell: bpy.types.Object = _mesh(
        "Helmet shell" if not officer else "Peaked cap crown",
        vertices,
        faces,
        material,
        parent,
    )
    thickness: bpy.types.SolidifyModifier = shell.modifiers.new(
        "Shell thickness", "SOLIDIFY"
    )
    thickness.thickness = 0.018
    if not officer:
        _curve(
            "Rolled helmet rim",
            [
                (
                    0.303 * math.cos(a * math.pi / 16),
                    0.259 * math.sin(a * math.pi / 16),
                    0.19,
                )
                for a in range(33)
            ],
            0.012,
            material,
            parent,
        )
    if officer:
        bpy.data.objects.remove(shell, do_unlink=True)
        crown_vertices: list[tuple[float, float, float]] = []
        crown_faces: list[tuple[int, ...]] = []
        for z, radius in (
            (0.33, 0.28),
            (0.39, 0.35),
            (0.49, 0.35),
            (0.54, 0.28),
            (0.55, 0.02),
        ):
            for column in range(64):
                angle: float = 2 * math.pi * column / 64
                crown_vertices.append(
                    (
                        radius * math.cos(angle),
                        radius * 0.85 * math.sin(angle),
                        z,
                    )
                )
        for row in range(4):
            for column in range(64):
                a: int = row * 64 + column
                b: int = row * 64 + (column + 1) % 64
                crown_faces.append((a, b, b + 64, a + 64))
        _mesh(
            "Service cap crown", crown_vertices, crown_faces, material, parent
        )
        _sphere(
            "Red service cap band",
            (0, 0, 0.285),
            (0.284, 0.246, 0.092),
            red,
            parent,
        )
        _sphere(
            "Black polished visor",
            (0, -0.24, 0.21),
            (0.29, 0.16, 0.025),
            dark,
            parent,
        )
        for x in (-0.045, 0.045):
            vertices = [
                (
                    x
                    + (0.028 if i % 2 == 0 else 0.012)
                    * math.sin(i * math.pi / 5),
                    -0.249,
                    0.293
                    + (0.028 if i % 2 == 0 else 0.012)
                    * math.cos(i * math.pi / 5),
                )
                for i in range(10)
            ]
            _mesh("Cap star", vertices, [tuple(range(10))], gold, parent)
        for index in range(13):
            x: float = -0.22 + index * 0.036
            _sphere(
                "Visor gold braid",
                (x, -0.337 + 0.07 * (abs(x) / 0.22) ** 2, 0.23),
                (0.014, 0.009, 0.006),
                gold,
                parent,
            )


def _character(
    name: str,
    location: tuple[float, float, float],
    height: float,
    uniform: bpy.types.Material,
    officer: bool,
    materials: dict[str, bpy.types.Material],
) -> dict[str, bpy.types.Object]:
    """Assemble one complete character with connected local animation pivots."""
    root: bpy.types.Object = _empty(name)
    root.location = location
    root.scale = (height, height, height)
    legs: dict[str, bpy.types.Object] = {}
    torso: bpy.types.Object = _empty(name + " torso", root)
    jacket_vertices: list[tuple[float, float, float]] = []
    jacket_faces: list[tuple[int, ...]] = []
    for z, width, depth in (
        (0.94, 0.23, 0.14),
        (1.02, 0.235, 0.145),
        (1.18, 0.22, 0.145),
        (1.40, 0.245, 0.16),
        (1.54, 0.275, 0.14),
        (1.65, 0.12, 0.10),
    ):
        for column in range(40):
            angle: float = 2 * math.pi * column / 40
            jacket_vertices.append(
                (
                    width * math.cos(angle),
                    depth * math.sin(angle),
                    z - (0.20 * max(0, -math.sin(angle)) if z == 1.65 else 0),
                )
            )
    for row in range(5):
        for column in range(40):
            a: int = row * 40 + column
            b: int = row * 40 + (column + 1) % 40
            jacket_faces.append((a, b, b + 40, a + 40))
    jacket: bpy.types.Object = _mesh(
        "Continuous uniform shirt",
        jacket_vertices,
        jacket_faces,
        uniform,
        torso,
    )
    smooth: bpy.types.SubsurfModifier = jacket.modifiers.new(
        "Cloth silhouette", "SUBSURF"
    )
    smooth.levels = 2
    _box(
        "Front placket",
        (0, -0.143, 1.24),
        (0.025, 0.018, 0.55),
        uniform,
        torso,
        0.006,
    )
    for side in (-1, 1):
        leg: bpy.types.Object = _empty(
            name + (" right hip" if side == -1 else " left hip"), root
        )
        leg.location = (side * 0.12, 0, 1.02)
        legs["right_leg" if side == -1 else "left_leg"] = leg
        _segment(
            "Trouser leg",
            (0, 0, 0),
            (side * 0.025, 0.01, -0.86),
            0.103,
            materials["trousers"] if officer else uniform,
            leg,
        )
        _sphere(
            "Blue foot",
            (side * 0.025, -0.077, -0.945),
            (0.09, 0.18, 0.065),
            materials["skin"],
            leg,
        )
        _box(
            "Breast pocket",
            (side * 0.145, -0.135, 1.36),
            (0.13, 0.008, 0.14),
            uniform,
            torso,
            0.006,
        )
        collar: bpy.types.Object = _mesh(
            "Folded collar",
            [
                (side * 0.015, -0.16, 1.48),
                (side * 0.105, -0.135, 1.68),
                (side * 0.24, -0.075, 1.57),
                (side * 0.18, -0.18, 1.43),
            ],
            [(0, 1, 2, 3)],
            uniform,
            torso,
        )
        solid: bpy.types.SolidifyModifier = collar.modifiers.new(
            "Cloth edge", "SOLIDIFY"
        )
        solid.thickness = 0.006
        if officer:
            _box(
                "Shoulder epaulette",
                (side * 0.20, 0.01, 1.595),
                (0.105, 0.07, 0.018),
                uniform,
                torso,
                0.012,
            )
    _segment(
        "Blue neck", (0, 0, 1.48), (0, 0, 1.92), 0.11, materials["skin"], torso
    )
    head: bpy.types.Object = _empty(name + " head", torso)
    head.location = (0, -0.005, 1.99)
    _head_surface(head, materials["skin"])
    for side in (-1, 1):
        _sphere(
            "Ear",
            (side * 0.228, 0.0, -0.015),
            (0.049, 0.04, 0.093),
            materials["skin"],
            head,
        )
        _sphere(
            "Ear concha",
            (side * 0.254, -0.028, -0.013),
            (0.018, 0.012, 0.052),
            materials["skin_shadow"],
            head,
        )
        _sphere(
            "Eye white",
            (side * 0.101, -0.184, 0.063),
            (0.068, 0.035, 0.037),
            materials["eye"],
            head,
        )
        _sphere(
            "Pupil",
            (side * 0.103, -0.223, 0.06),
            (0.018, 0.008, 0.023),
            materials["dark"],
            head,
        )
        _sphere(
            "Eye glint",
            (side * 0.103 - 0.009, -0.234, 0.073),
            (0.006, 0.003, 0.007),
            materials["eye"],
            head,
        )
        _curve(
            "Upper eyelid",
            [
                (side * 0.043, -0.206, 0.082),
                (side * 0.101, -0.225, 0.093),
                (side * 0.164, -0.18, 0.084),
            ],
            0.013,
            materials["skin"],
            head,
        )
        _curve(
            "Lower eyelid",
            [
                (side * 0.043, -0.206, 0.047),
                (side * 0.101, -0.219, 0.023),
                (side * 0.164, -0.18, 0.05),
            ],
            0.008,
            materials["skin"],
            head,
        )
        _curve(
            "Expressive eyebrow",
            [
                (side * 0.045, -0.193, 0.137),
                (side * 0.102, -0.207, 0.149),
                (side * 0.167, -0.159, 0.137),
            ],
            0.011,
            materials["dark"],
            head,
        )
    mouth: bpy.types.Object = _sphere(
        "Animated mouth opening",
        (0, -0.191, -0.205),
        (0.106, 0.027, 0.035),
        materials["mouth"],
        head,
    )
    _curve(
        "Lower lip",
        [
            (-0.105, -0.198, -0.205),
            (0, -0.216, -0.247),
            (0.105, -0.198, -0.205),
        ],
        0.005,
        materials["skin"],
        head,
    )
    _helmet(
        head,
        materials["helmet"],
        officer,
        materials["red"],
        materials["gold"],
        materials["dark"],
    )
    arms: dict[str, bpy.types.Object] = {}
    for side, label in ((-1, "right"), (1, "left")):
        arm: bpy.types.Object = _empty(name + " " + label + " shoulder", torso)
        arm.location = (side * 0.25, 0, 1.55)
        _segment(
            "Upper sleeve",
            (0, 0, 0),
            (side * 0.065, 0, -0.30),
            0.10,
            uniform,
            arm,
        )
        elbow: bpy.types.Object = _empty(name + " " + label + " elbow", arm)
        elbow.location = (side * 0.065, 0, -0.30)
        _sphere(
            "Sleeve elbow fold",
            (0, 0, 0),
            (0.078, 0.075, 0.080),
            uniform,
            elbow,
        )
        _segment(
            "Forearm sleeve",
            (0, 0, 0),
            (side * 0.015, -0.02, -0.28),
            0.078,
            uniform,
            elbow,
        )
        _sphere(
            "Blue palm",
            (side * 0.015, -0.022, -0.345),
            (0.058, 0.035, 0.081),
            materials["skin"],
            elbow,
        )
        for finger in range(4):
            _segment(
                "Finger",
                (side * 0.015 - 0.038 + finger * 0.025, -0.03, -0.37),
                (
                    side * 0.015 - 0.038 + finger * 0.025,
                    -0.035,
                    -0.445 + 0.014 * abs(finger - 1),
                ),
                0.012,
                materials["skin"],
                elbow,
            )
        _segment(
            "Thumb",
            (side * 0.063, -0.018, -0.32),
            (side * 0.09, -0.05, -0.37),
            0.021,
            materials["skin"],
            elbow,
        )
        arms[label] = arm
        arms[label + "_elbow"] = elbow
    for z in (1.10, 1.26, 1.42):
        _sphere(
            "Jacket button",
            (0.02, -0.16, z),
            (0.014, 0.008, 0.014),
            materials["gold"] if officer else materials["helmet"],
            torso,
        )
    if officer:
        for row in range(3):
            for column in range(6):
                _box(
                    "Service ribbon",
                    (0.11 + 0.022 * column, -0.17, 1.37 - 0.028 * row),
                    (0.02, 0.012, 0.022),
                    materials[
                        ("red", "gold", "ribbon_blue")[(row + column) % 3]
                    ],
                    torso,
                    0.001,
                )
        for arm in (arms["right"], arms["left"]):
            arm.rotation_euler.x = 0.35
        for elbow in (arms["right_elbow"], arms["left_elbow"]):
            elbow.rotation_euler.x = 0.65
    return {
        "root": root,
        "torso": torso,
        "head": head,
        "mouth": mouth,
        **arms,
        **legs,
    }


def _base(materials: dict[str, bpy.types.Material]) -> None:
    """Populate the sandy parade ground with arched huts, flag, and tanks."""
    _box(
        "Desert parade ground",
        (0, 0, -0.10),
        (400, 400, 0.2),
        materials["sand"],
        bevel=0,
    )
    for index in range(5):
        x: float = -21 + index * 2.7
        vertices: list[tuple[float, float, float]] = []
        faces: list[tuple[int, ...]] = []
        for y in (24, 29):
            for j in range(25):
                a: float = math.pi * j / 24
                vertices.append(
                    (x + 1.2 * math.cos(a), y, 0.7 + 1.2 * math.sin(a))
                )
        for j in range(24):
            faces.append((j, j + 1, j + 26, j + 25))
        _mesh("Arched barrack roof", vertices, faces, materials["hut"])
        _box(
            "Barrack front",
            (x, 24, 0.52),
            (2.4, 0.12, 1.04),
            materials["hut"],
            bevel=0.015,
        )
        _mesh(
            "Barrack arch front",
            [(x, 23.99, 0.7)]
            + [
                (
                    x + 1.2 * math.cos(j * math.pi / 24),
                    23.99,
                    0.7 + 1.2 * math.sin(j * math.pi / 24),
                )
                for j in range(25)
            ],
            [tuple(range(26))],
            materials["hut"],
        )
        _box(
            "Barrack entrance",
            (x, 23.90, 0.53),
            (0.6, 0.04, 1.06),
            materials["hut_door"],
            bevel=0.01,
        )
        for y in (24.03, 25.7, 27.3, 28.95):
            _curve(
                "Roof seam",
                [
                    (
                        x + 1.206 * math.cos(j * math.pi / 24),
                        y,
                        0.7 + 1.206 * math.sin(j * math.pi / 24),
                    )
                    for j in range(25)
                ],
                0.012,
                materials["hut_door"],
            )
    _segment(
        "Flag pole",
        (-9.8, 24, 0),
        (-9.8, 24, 3.2),
        0.022,
        materials["hut_door"],
        _empty("Flag"),
    )
    _mesh(
        "Yellow flag",
        [
            (-9.8, 24, 3.2),
            (-9.15, 24.08, 3.18),
            (-9.12, 24.14, 2.82),
            (-9.8, 24, 2.84),
        ],
        [(0, 1, 2, 3)],
        materials["flag"],
    )
    for index, x in enumerate((0, 5, 10, 15)):
        tank: bpy.types.Object = _empty("Background tank " + str(index))
        tank.location = (x, 26 + (index % 2) * 3, 0)
        tank.rotation_euler.z = -1.3
        _box(
            "Tank hull",
            (0, 0, 0.63),
            (2.3, 3.1, 0.65),
            materials["tank"],
            tank,
            0.2,
        )
        _box(
            "Turret",
            (0, -0.1, 1.18),
            (1.55, 1.6, 0.65),
            materials["tank"],
            tank,
            0.23,
        )
        _segment(
            "Tank barrel",
            (0, -0.7, 1.27),
            (0, -3.0, 1.27),
            0.09,
            materials["tank"],
            tank,
        )
        for side in (-1, 1):
            _box(
                "Continuous track",
                (side * 1.1, 0, 0.37),
                (0.38, 3.0, 0.65),
                materials["track"],
                tank,
                0.22,
            )
            for wheel in range(6):
                _sphere(
                    "Road wheel",
                    (side * 1.3, -1.15 + wheel * 0.46, 0.35),
                    (0.035, 0.23, 0.24),
                    materials["tank"],
                    tank,
                )
        for _patch in range(16):
            _sphere(
                "Tank camouflage patch",
                (random.uniform(-1, 1), random.uniform(-1.3, 1.3), 0.96),
                (0.24, 0.3, 0.015),
                materials["hut"],
                tank,
            )


def _key_pose(
    obj: bpy.types.Object, frame: int, rotation: tuple[float, float, float]
) -> None:
    """Insert a normal Euler animation key, requiring no runtime callbacks."""
    obj.rotation_euler = rotation
    obj.keyframe_insert(data_path="rotation_euler", frame=frame)


def _animate(characters: list[dict[str, bpy.types.Object]]) -> None:
    """Bake conversational body motion, facial motion, salute, and reaction."""
    intervals: list[tuple[int, int, int]] = [
        (0, 1, 210),
        (3, 280, 335),
        (0, 336, 449),
        (3, 450, 643),
        (2, 644, 721),
        (0, 722, 780),
        (3, 781, 898),
    ]
    for index, character in enumerate(characters):
        character["root"].rotation_euler.z = (
            math.pi / 2 if index == 0 else -0.3
        )
        for frame in range(1, 899, 5):
            speaking: bool = any(
                speaker == index and begin <= frame <= end
                for speaker, begin, end in intervals
            )
            amount: float = 1 if speaking else 0.2
            _key_pose(
                character["head"],
                frame,
                (
                    amount * 0.04 * math.sin(frame * 0.064),
                    0.035 * math.sin(frame * 0.032 + index),
                    amount * 0.10 * math.sin(frame * 0.041)
                    - (0.06 if index == 3 else 0),
                ),
            )
            _key_pose(
                character["torso"],
                frame,
                (
                    0.012 * math.sin(frame * 0.041 + index),
                    0.018 * math.sin(frame * 0.035),
                    0,
                ),
            )
            character["mouth"].scale.z = 0.012 + (
                0.042 * abs(math.sin(frame * 0.71)) if speaking else 0
            )
            character["mouth"].keyframe_insert(data_path="scale", frame=frame)
        for eye in [
            obj
            for obj in bpy.data.objects
            if obj.parent == character["head"]
            and obj.name.startswith(("Eye white", "Pupil", "Eye glint"))
        ]:
            height: float = eye.scale.z
            for blink in (
                95 + index * 7,
                240 + index * 9,
                420 + index * 6,
                630 + index * 4,
                805 + index * 8,
            ):
                for frame, factor in (
                    (1, 1),
                    (blink - 3, 1),
                    (blink, 0.05),
                    (blink + 3, 1),
                    (898, 1),
                ):
                    eye.scale.z = height * factor
                    eye.keyframe_insert(data_path="scale", frame=frame)
        if index == 0:
            for frame in range(1, 102, 5):
                stride: float = (
                    0.26
                    * math.sin((frame - 1) * math.pi / 25)
                    * min(1, (101 - frame) / 15)
                )
                _key_pose(character["right_leg"], frame, (stride, 0, 0))
                _key_pose(character["left_leg"], frame, (-stride, 0, 0))
            for frame, x in ((1, -3.1), (75, -2.6), (100, -2.2), (898, -2.2)):
                character["root"].location.x = x
                character["root"].keyframe_insert(
                    data_path="location", frame=frame
                )
            _key_pose(character["left"], 720, (0.35, 0, 0))
            _key_pose(character["left"], 742, (-1.0, -0.2, -0.6))
            _key_pose(character["left"], 770, (-1.0, -0.2, -0.6))
            _key_pose(character["left"], 790, (0.35, 0, 0))
        else:
            for frame in (1, 430, 898):
                _key_pose(character["right"], frame, (0.04, 0, 0.06))
                _key_pose(character["left"], frame, (0.04, 0, -0.06))
    main: dict[str, bpy.types.Object] = characters[3]
    for frame, raised in (
        (1, False),
        (570, False),
        (596, True),
        (623, True),
        (646, False),
        (898, False),
    ):
        _key_pose(main["right"], frame, (0, 2.2 if raised else 0, 0.06))
        _key_pose(main["right_elbow"], frame, (0, 1.5 if raised else 0, 0))
    short: dict[str, bpy.types.Object] = characters[2]
    for frame, angle in (
        (1, 0),
        (640, 0),
        (664, -1.0),
        (710, -0.6),
        (735, 0),
        (898, 0),
    ):
        _key_pose(short["left"], frame, (angle, 0, -0.35 if angle else -0.06))
        _key_pose(short["left_elbow"], frame, (angle, 0, 0))


def _camera(
    name: str,
    position: tuple[float, float, float],
    target: tuple[float, float, float],
    lens: float,
) -> bpy.types.Object:
    """Create a calibrated perspective shot camera at an explicit target."""
    data: bpy.types.Camera = bpy.data.cameras.new(name)
    obj: bpy.types.Object = bpy.data.objects.new(name, data)
    bpy.context.collection.objects.link(obj)
    obj.location = position
    obj.rotation_euler = (
        (mathutils.Vector(target) - obj.location)
        .to_track_quat("-Z", "Y")
        .to_euler()
    )
    data.lens = lens
    data.clip_end = 500
    return obj


def _main() -> None:
    """Validate the source and generate the scene, preview packet, or animation."""
    parser: argparse.ArgumentParser = argparse.ArgumentParser(
        description=__doc__
    )
    parser.add_argument("--reference", type=pathlib.Path, required=True)
    parser.add_argument("--soundtrack", type=pathlib.Path)
    parser.add_argument("--output", type=pathlib.Path, required=True)
    parser.add_argument("--preview", action="store_true")
    parser.add_argument("--animate", action="store_true")
    arguments: argparse.Namespace = parser.parse_args(
        sys.argv[sys.argv.index("--") + 1 :]
    )
    if (
        hashlib.sha256(arguments.reference.read_bytes()).hexdigest()
        != _REFERENCE_SHA256
    ):
        raise ValueError(
            "Reference SHA-256 does not match the controlling clip"
        )
    output: pathlib.Path = arguments.output.resolve()
    output.mkdir(parents=True, exist_ok=True)
    bpy.ops.object.select_all(action="SELECT")
    bpy.ops.object.delete(use_global=False)
    random.seed(11)
    materials: dict[str, bpy.types.Material] = {
        "skin": _material("Powder blue skin", (0.24, 0.48, 0.72), 0.6),
        "skin_shadow": _material("Ear recess blue", (0.24, 0.43, 0.54)),
        "eye": _material("Warm ivory eyes", (0.85, 0.84, 0.74), 0.4),
        "dark": _material("Brow and pupil charcoal", (0.022, 0.025, 0.018)),
        "mouth": _material("Mouth cavity", (0.025, 0.018, 0.014)),
        "helmet": _material("Olive drab helmet", (0.065, 0.075, 0.035)),
        "officer": _material("Officer olive jacket", (0.055, 0.06, 0.020)),
        "trousers": _material("Officer khaki trousers", (0.37, 0.31, 0.17)),
        "red": _material("Cap band crimson", (0.43, 0.065, 0.055)),
        "gold": _material("Brass insignia", (0.65, 0.49, 0.12), 0.4),
        "ribbon_blue": _material("Service ribbon blue", (0.11, 0.15, 0.42)),
        "sand": _material("Pale desert sand", (0.76, 0.68, 0.43)),
        "hut": _material("Canvas barracks", (0.59, 0.60, 0.47)),
        "hut_door": _material(
            "Barrack seams and entrance", (0.31, 0.33, 0.25)
        ),
        "tank": _material("Tank sand olive", (0.48, 0.48, 0.32)),
        "track": _material("Dusty tank track", (0.35, 0.35, 0.25)),
        "flag": _material("Signal flag yellow", (0.64, 0.52, 0.09)),
        "woodland": _camouflage("Digital woodland camouflage", False),
        "desert": _camouflage("Digital desert camouflage", True),
    }
    sky: bpy.types.Material = bpy.data.materials.new("Blue to ivory sky")
    sky.use_nodes = True
    sky.node_tree.nodes.clear()
    sky_geometry: bpy.types.Node = sky.node_tree.nodes.new(
        "ShaderNodeTexCoord"
    )
    sky_split: bpy.types.Node = sky.node_tree.nodes.new(
        "ShaderNodeSeparateXYZ"
    )
    sky_map: bpy.types.Node = sky.node_tree.nodes.new("ShaderNodeMapRange")
    sky_map.inputs["From Min"].default_value = 0.5
    sky_map.inputs["From Max"].default_value = 0.55
    sky_colors: bpy.types.Node = sky.node_tree.nodes.new("ShaderNodeValToRGB")
    sky_colors.color_ramp.elements[0].color = (0.86, 0.84, 0.64, 1)
    sky_colors.color_ramp.elements[1].color = (0.23, 0.36, 0.57, 1)
    sky_emission: bpy.types.Node = sky.node_tree.nodes.new(
        "ShaderNodeEmission"
    )
    sky_surface: bpy.types.Node = sky.node_tree.nodes.new(
        "ShaderNodeOutputMaterial"
    )
    sky.node_tree.links.new(
        sky_geometry.outputs["Generated"], sky_split.inputs[0]
    )
    sky.node_tree.links.new(sky_split.outputs["Z"], sky_map.inputs["Value"])
    sky.node_tree.links.new(sky_map.outputs[0], sky_colors.inputs[0])
    sky.node_tree.links.new(
        sky_colors.outputs[0], sky_emission.inputs["Color"]
    )
    sky.node_tree.links.new(
        sky_emission.outputs[0], sky_surface.inputs["Surface"]
    )
    dome: bpy.types.Object = _sphere(
        "Sky backdrop", (0, 0, 0), (150, 150, 150), sky
    )
    dome.visible_diffuse = False
    dome.visible_glossy = False
    dome.visible_shadow = False
    dome.visible_transmission = False
    _base(materials)
    characters: list[dict[str, bpy.types.Object]] = [
        _character(
            "Officer",
            (-2.2, -2.2, 0),
            1.0,
            materials["officer"],
            True,
            materials,
        ),
        _character(
            "Tall recruit",
            (-0.45, 0.7, 0),
            1.03,
            materials["woodland"],
            False,
            materials,
        ),
        _character(
            "Short recruit",
            (0.45, 0.55, 0),
            0.81,
            materials["woodland"],
            False,
            materials,
        ),
        _character(
            "Desert recruit",
            (2.3, 0.4, 0),
            1.0,
            materials["desert"],
            False,
            materials,
        ),
    ]
    _animate(characters)
    scene: bpy.types.Scene = bpy.context.scene
    scene.name = "Army spot remake"
    scene.frame_start = 1
    scene.frame_end = 898
    scene.render.fps = 30
    scene.render.fps_base = 1.001
    scene.render.engine = "CYCLES"
    scene.cycles.samples = 8
    scene.cycles.use_denoising = True
    scene.render.use_persistent_data = True
    scene.cycles.max_bounces = 3
    scene.render.resolution_x = 960
    scene.render.resolution_y = 540
    scene.render.resolution_percentage = 100
    scene.world.use_nodes = True
    sky_nodes: bpy.types.Nodes = scene.world.node_tree.nodes
    sky_coordinates: bpy.types.Node = sky_nodes.new("ShaderNodeTexCoord")
    sky_height: bpy.types.Node = sky_nodes.new("ShaderNodeSeparateXYZ")
    sky_range: bpy.types.Node = sky_nodes.new("ShaderNodeMapRange")
    sky_range.inputs["From Min"].default_value = 0.0
    sky_range.inputs["From Max"].default_value = 0.6
    sky_ramp: bpy.types.Node = sky_nodes.new("ShaderNodeValToRGB")
    sky_ramp.color_ramp.elements[0].color = (0.88, 0.86, 0.69, 1)
    sky_ramp.color_ramp.elements[1].color = (0.25, 0.38, 0.58, 1)
    sky_links: bpy.types.NodeLinks = scene.world.node_tree.links
    sky_links.new(sky_coordinates.outputs["Normal"], sky_height.inputs[0])
    sky_links.new(sky_height.outputs["Z"], sky_range.inputs["Value"])
    sky_links.new(sky_range.outputs[0], sky_ramp.inputs[0])
    sky_links.new(
        sky_ramp.outputs[0], sky_nodes.get("Background").inputs["Color"]
    )
    background: bpy.types.Node = sky_nodes.get("Background")
    for link in list(sky_links):
        if link.to_node == background:
            sky_links.remove(link)
    background.inputs["Color"].default_value = (0.65, 0.75, 0.95, 1)
    background.inputs["Strength"].default_value = 0.7
    bpy.ops.object.light_add(type="AREA", location=(-3, -4, 8))
    bpy.context.object.name = "Soft desert daylight"
    bpy.context.object.data.energy = 1700
    bpy.context.object.data.shape = "DISK"
    bpy.context.object.data.size = 6
    bpy.ops.object.light_add(type="SUN", location=(-8, -6, 10))
    bpy.context.object.name = "Low desert sun"
    bpy.context.object.rotation_euler = (0.35, -0.5, -0.4)
    bpy.context.object.data.energy = 2
    bpy.context.object.data.angle = 0.15
    scene.view_settings.view_transform = "AgX"
    cameras: dict[str, bpy.types.Object] = {
        "wide": _camera(
            "Wide parade ground", (3.2, -12.5, 0.95), (-0.3, 0.4, 1.35), 44
        ),
        "officer": _camera(
            "Officer closeup", (0.95, -2.70, 1.55), (-2.2, -2.2, 1.82), 48
        ),
        "pair": _camera(
            "Woodland recruits", (-0.3, -3.1, 2.0), (-0.08, 0.6, 1.67), 48
        ),
        "main": _camera(
            "Desert recruit closeup", (1.12, -2.8, 1.86), (2.3, 0.4, 1.95), 52
        ),
        "reverse": _camera(
            "Reverse group", (-4.8, -6.6, 2.5), (-0.1, 0.25, 1.25), 40
        ),
    }
    for frame, shot in (
        (1, "wide"),
        (108, "officer"),
        (217, "pair"),
        (281, "main"),
        (338, "reverse"),
        (451, "main"),
        (647, "pair"),
        (722, "wide"),
        (782, "main"),
    ):
        marker: bpy.types.TimelineMarker = scene.timeline_markers.new(
            shot, frame=frame
        )
        marker.camera = cameras[shot]
    scene.camera = cameras["wide"]
    if arguments.soundtrack:
        sound_strip: bpy.types.SoundStrip = (
            scene.sequence_editor_create().strips.new_sound(
                "Original soundtrack",
                str(arguments.soundtrack.resolve()),
                channel=1,
                frame_start=1,
            )
        )
        sound_strip.sound.pack()
        sound_strip.sound.filepath = "//army_spot_audio.m4a"
    scene.render.use_sequencer = False
    scene.render.image_settings.file_format = "PNG"
    scene.render.filepath = str(output / "frames" / "frame_")
    scene.render.use_file_extension = True
    scene.render.film_transparent = False
    bpy.context.preferences.filepaths.save_version = 0
    scene.frame_set(1)
    bpy.ops.wm.save_as_mainfile(filepath=str(output / "army_spot.blend"))
    manifest: dict[str, str | int] = {
        "reference_sha256": _REFERENCE_SHA256,
        "blender_version": bpy.app.version_string,
        "frame_count": 898,
        "fps_numerator": 30000,
        "fps_denominator": 1001,
        "width": 960,
        "height": 540,
        "objects": len(scene.objects),
    }
    (output / "scene_manifest.json").write_text(
        json.dumps(manifest, indent=2) + "\n"
    )
    if arguments.preview:
        for frame in (1, 181, 241, 301, 421, 601, 681, 841):
            scene.frame_set(frame)
            scene.render.filepath = str(output / f"preview_{frame:04}.png")
            bpy.ops.render.render(write_still=True)
    if arguments.animate:
        (output / "frames").mkdir(exist_ok=True)
        scene.render.filepath = str(output / "frames" / "frame_")
        bpy.ops.render.render(animation=True)


if __name__ == "__main__":
    _main()
