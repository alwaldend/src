"""Pure validation for deterministic Blender review-packet renders."""

import dataclasses
import hashlib
import json
import pathlib
import typing

_LOADED_SOURCE_SHA256 = hashlib.sha256(
    pathlib.Path(__file__).read_bytes()
).hexdigest()


class SpecError(ValueError):
    """Raised when a render-packet specification is invalid."""


RGB = typing.Tuple[int, int, int]


@dataclasses.dataclass(frozen=True)
class Resolution:
    """Pixel dimensions for every render in a packet."""

    width: int
    height: int


@dataclasses.dataclass(frozen=True)
class BeautyView:
    """One named beauty render."""

    name: str
    camera: str
    output: str


@dataclasses.dataclass(frozen=True)
class Fallback:
    """Treatment for renderable meshes absent from an ID mapping."""

    mode: str
    rgb: typing.Optional[RGB]


@dataclasses.dataclass(frozen=True)
class ComponentIDPass:
    """One named flat component-ID render."""

    name: str
    camera: str
    output: str
    objects: typing.Mapping[str, RGB]
    fallback: Fallback


@dataclasses.dataclass(frozen=True)
class RenderSpec:
    """Validated schema-v1 render-packet specification."""

    scene: str
    resolution: Resolution
    beauty_views: typing.Mapping[str, BeautyView]
    component_id_passes: typing.Mapping[str, ComponentIDPass]

    def output_names(self) -> typing.List[str]:
        """Return packet output paths in deterministic render order."""
        beauty = [view.output for view in self.beauty_views.values()]
        component_ids = [
            component_pass.output
            for component_pass in self.component_id_passes.values()
        ]
        return beauty + component_ids


def inverse_srgb_channel(channel: int) -> float:
    """Convert one requested 8-bit sRGB channel to scene-linear emission."""
    _validate_channel(channel, "channel")
    encoded = channel / 255.0
    if encoded <= 0.04045:
        return encoded / 12.92
    return ((encoded + 0.055) / 1.055) ** 2.4


def inverse_srgb(rgb: RGB) -> typing.Tuple[float, float, float]:
    """Convert an 8-bit sRGB triplet to scene-linear emission values."""
    validated = _rgb(list(rgb), "rgb")
    return (
        inverse_srgb_channel(validated[0]),
        inverse_srgb_channel(validated[1]),
        inverse_srgb_channel(validated[2]),
    )


def resolve_output_path(output_dir: pathlib.Path, output: str) -> pathlib.Path:
    """Resolve one validated output beneath an output directory.

    The second containment check catches pre-existing symlinks that would
    otherwise redirect a lexically safe packet path outside ``output_dir``.
    """
    normalized = _png_path(output, "output")
    root = output_dir.resolve()
    candidate = root.joinpath(
        *pathlib.PurePosixPath(normalized).parts
    ).resolve()
    try:
        candidate.relative_to(root)
    except ValueError as error:
        raise SpecError(
            f"output {output!r} resolves outside output directory"
        ) from error
    return candidate


def load_spec(path: pathlib.Path) -> RenderSpec:
    """Load and validate a schema-v1 JSON render specification."""
    try:
        with path.open("r", encoding="utf-8") as stream:
            value = json.load(stream, object_pairs_hook=_unique_object)
    except (OSError, UnicodeError, json.JSONDecodeError) as error:
        raise SpecError(f"cannot read spec {path}: {error}") from error
    return parse_spec(value)


def parse_spec(value: object) -> RenderSpec:
    """Validate a decoded schema-v1 render specification."""
    root = _object(value, "spec")
    _keys(
        root,
        required={"schema_version", "scene", "resolution", "beauty_views"},
        optional={"component_id_passes"},
        label="spec",
    )
    version = root["schema_version"]
    if isinstance(version, bool) or version != 1:
        raise SpecError("schema_version must be the integer 1")

    scene = _name(root["scene"], "scene")
    resolution = _resolution(root["resolution"])
    beauty_views = _beauty_views(root["beauty_views"])
    component_passes = _component_passes(root.get("component_id_passes", {}))
    _check_outputs(beauty_views, component_passes)
    return RenderSpec(
        scene=scene,
        resolution=resolution,
        beauty_views=beauty_views,
        component_id_passes=component_passes,
    )


def _unique_object(
    pairs: typing.List[typing.Tuple[str, object]],
) -> typing.Dict[str, object]:
    result: typing.Dict[str, object] = {}
    for key, value in pairs:
        if key in result:
            raise SpecError(f"duplicate JSON key {key!r}")
        result[key] = value
    return result


def _object(value: object, label: str) -> typing.Dict[str, object]:
    if not isinstance(value, dict):
        raise SpecError(f"{label} must be an object")
    result: typing.Dict[str, object] = {}
    for key, child in value.items():
        if not isinstance(key, str):
            raise SpecError(f"{label} keys must be strings")
        result[key] = child
    return result


def _keys(
    value: typing.Mapping[str, object],
    required: typing.Set[str],
    optional: typing.Set[str],
    label: str,
) -> None:
    missing = required - set(value)
    unknown = set(value) - required - optional
    if missing:
        raise SpecError(f"{label} is missing keys {sorted(missing)}")
    if unknown:
        raise SpecError(f"{label} has unknown keys {sorted(unknown)}")


def _name(value: object, label: str) -> str:
    if not isinstance(value, str) or not value.strip():
        raise SpecError(f"{label} must be a non-empty string")
    return value


def _positive_dimension(value: object, label: str) -> int:
    if isinstance(value, bool) or not isinstance(value, int):
        raise SpecError(f"{label} must be an integer")
    if not 1 <= value <= 65536:
        raise SpecError(f"{label} must be from 1 to 65536")
    return value


def _resolution(value: object) -> Resolution:
    raw = _object(value, "resolution")
    _keys(
        raw,
        required={"width", "height"},
        optional=set(),
        label="resolution",
    )
    return Resolution(
        width=_positive_dimension(raw["width"], "resolution.width"),
        height=_positive_dimension(raw["height"], "resolution.height"),
    )


def _png_path(value: object, label: str) -> str:
    raw = _name(value, label)
    if "\\" in raw:
        raise SpecError(f"{label} must use forward slashes")
    path = pathlib.PurePosixPath(raw)
    if path.is_absolute() or ".." in path.parts:
        raise SpecError(f"{label} must stay beneath output directory")
    normalized = path.as_posix()
    if normalized in {"", "."}:
        raise SpecError(f"{label} must name a file")
    if path.suffix != ".png":
        raise SpecError(f"{label} must end in .png")
    if normalized == "manifest.json":
        raise SpecError(f"{label} is reserved for the packet manifest")
    return normalized


def _beauty_views(value: object) -> typing.Dict[str, BeautyView]:
    raw_views = _object(value, "beauty_views")
    if not raw_views:
        raise SpecError("beauty_views must define at least one view")
    result: typing.Dict[str, BeautyView] = {}
    for raw_name in sorted(raw_views):
        name = _name(raw_name, "beauty view name")
        label = f"beauty_views.{name}"
        raw = _object(raw_views[raw_name], label)
        _keys(
            raw,
            required={"camera", "output"},
            optional=set(),
            label=label,
        )
        result[name] = BeautyView(
            name=name,
            camera=_name(raw["camera"], f"{label}.camera"),
            output=_png_path(raw["output"], f"{label}.output"),
        )
    return result


def _component_passes(value: object) -> typing.Dict[str, ComponentIDPass]:
    raw_passes = _object(value, "component_id_passes")
    result: typing.Dict[str, ComponentIDPass] = {}
    for raw_name in sorted(raw_passes):
        name = _name(raw_name, "component-ID pass name")
        label = f"component_id_passes.{name}"
        raw = _object(raw_passes[raw_name], label)
        _keys(
            raw,
            required={"camera", "output", "objects", "fallback"},
            optional=set(),
            label=label,
        )
        result[name] = ComponentIDPass(
            name=name,
            camera=_name(raw["camera"], f"{label}.camera"),
            output=_png_path(raw["output"], f"{label}.output"),
            objects=_component_objects(raw["objects"], label),
            fallback=_fallback(raw["fallback"], label),
        )
    return result


def _component_objects(
    value: object, pass_label: str
) -> typing.Dict[str, RGB]:
    label = f"{pass_label}.objects"
    raw_objects = _object(value, label)
    if not raw_objects:
        raise SpecError(f"{label} must map at least one mesh object")
    result: typing.Dict[str, RGB] = {}
    for raw_name in sorted(raw_objects):
        name = _name(raw_name, f"{label} object name")
        result[name] = _rgb(raw_objects[raw_name], f"{label}.{name}")
    return result


def _fallback(value: object, pass_label: str) -> Fallback:
    label = f"{pass_label}.fallback"
    raw = _object(value, label)
    mode = _name(raw.get("mode"), f"{label}.mode")
    if mode == "hide":
        _keys(raw, required={"mode"}, optional=set(), label=label)
        return Fallback(mode=mode, rgb=None)
    if mode == "color":
        _keys(raw, required={"mode", "rgb"}, optional=set(), label=label)
        return Fallback(mode=mode, rgb=_rgb(raw["rgb"], f"{label}.rgb"))
    raise SpecError(f"{label}.mode must be 'hide' or 'color'")


def _rgb(value: object, label: str) -> RGB:
    if not isinstance(value, list) or len(value) != 3:
        raise SpecError(f"{label} must contain exactly three RGB channels")
    channels = [
        _validate_channel(channel, f"{label}[{index}]")
        for index, channel in enumerate(value)
    ]
    return channels[0], channels[1], channels[2]


def _validate_channel(value: object, label: str) -> int:
    if isinstance(value, bool) or not isinstance(value, int):
        raise SpecError(f"{label} must be an integer from 0 to 255")
    if not 0 <= value <= 255:
        raise SpecError(f"{label} must be an integer from 0 to 255")
    return value


def _check_outputs(
    beauty_views: typing.Mapping[str, BeautyView],
    component_passes: typing.Mapping[str, ComponentIDPass],
) -> None:
    owners: typing.Dict[str, str] = {}
    items = [
        (view.output, f"beauty view {view.name!r}")
        for view in beauty_views.values()
    ] + [
        (component_pass.output, f"component-ID pass {component_pass.name!r}")
        for component_pass in component_passes.values()
    ]
    for output, owner in items:
        if output in owners:
            raise SpecError(
                f"duplicate output {output!r} for {owners[output]} and {owner}"
            )
        owners[output] = owner
