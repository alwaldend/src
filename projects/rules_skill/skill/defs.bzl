"""Public API for packaging and validating Codex skills."""

load(
    "//skill/internal:skill_discovery_links.bzl",
    _skill_discovery_links = "skill_discovery_links",
)
load(
    "//skill/internal:skill_library.bzl",
    _SkillInfo = "SkillInfo",
    _skill_library = "skill_library",
)
load(
    "//skill/internal:skill_validation.bzl",
    _skill_validation = "skill_validation",
    _skill_validation_aspect = "skill_validation_aspect",
)

SkillInfo = _SkillInfo
skill_discovery_links = _skill_discovery_links
skill_library = _skill_library
skill_validation = _skill_validation
skill_validation_aspect = _skill_validation_aspect
