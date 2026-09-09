"""Public API for packaging, writing, and validating Codex skills."""

load(
    "//skill/internal:skill_archive.bzl",
    _skill_archive = "skill_archive",
    _skill_archives = "skill_archives",
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
load(
    "//skill/internal:skills_write.bzl",
    _SkillsWriteInfo = "SkillsWriteInfo",
    _skills_write = "skills_write",
    _skills_write_check_test = "skills_write_check_test",
    _skills_write_updater = "skills_write_updater",
)

SkillInfo = _SkillInfo
SkillsWriteInfo = _SkillsWriteInfo
skill_archive = _skill_archive
skill_archives = _skill_archives
skill_library = _skill_library
skill_validation = _skill_validation
skill_validation_aspect = _skill_validation_aspect
skills_write = _skills_write
skills_write_check_test = _skills_write_check_test
skills_write_updater = _skills_write_updater
