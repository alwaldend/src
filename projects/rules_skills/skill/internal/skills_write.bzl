"""Rules that write and validate a source-tree skill discovery directory."""

load(":skill_library.bzl", "SkillInfo")

SkillsWriteInfo = provider(
    doc = "Source-tree skill discovery entries derived from skill libraries.",
    fields = {
        "discovery_dir": "The normalized workspace-relative discovery directory.",
        "entries": "Deterministically ordered discovery entries by skill name.",
    },
)

_EntryInfo = provider(
    doc = "One resolved skill discovery entry.",
    fields = {
        "destination": "Entry path relative to the discovery directory.",
        "files": "Declared payload files for a written entry.",
        "kind": "Either symlink or written.",
        "name": "The logical skill name.",
        "symlink_target": "Relative symlink target for a symlink entry.",
    },
)

_VALID_NAME_CHARACTERS = "-0123456789abcdefghijklmnopqrstuvwxyz"

def _normalized_relative_path(path, description):
    if not path:
        fail("{} must be a non-empty workspace-relative path".format(description))
    if path.startswith("/"):
        fail("{} must be workspace-relative, got {}".format(description, path))
    parts = path.split("/")
    if any([part in ["", ".", ".."] for part in parts]):
        fail(
            (
                "{} must be normalized and must not contain empty, '.' or " +
                "'..' segments, got {}"
            ).format(description, path),
        )
    return "/".join(parts)

def _validate_name(name):
    invalid = not name
    for index in range(len(name)):
        if name[index] not in _VALID_NAME_CHARACTERS:
            invalid = True
            break
    if invalid:
        fail(
            (
                "skill discovery name must contain only lowercase ASCII " +
                "letters, digits, and hyphens, got {}"
            ).format(name),
        )

def _relative_path(from_directory, to_path):
    from_parts = from_directory.split("/")
    to_parts = to_path.split("/")
    common = 0
    for index in range(min(len(from_parts), len(to_parts))):
        if from_parts[index] != to_parts[index]:
            break
        common += 1
    relative_parts = [".."] * (len(from_parts) - common)
    relative_parts.extend(to_parts[common:])
    return "/".join(relative_parts) or "."

def _runfiles_path(workspace_name, file):
    short_path = file.short_path
    if short_path.startswith("../"):
        return short_path[3:]
    return "{}/{}".format(workspace_name, short_path)

def _symlink_entry(ctx, target):
    skill = target[SkillInfo]
    _validate_name(skill.name)
    if not skill.skill.is_source:
        fail(
            "symlink entries require a source SKILL.md in the same " +
            "repository; {} is generated or external".format(target.label),
        )
    if skill.skill.owner.repo_name != ctx.label.repo_name:
        fail(
            "symlink entries must live in the consuming repository: {} is " +
            "in repository {}".format(target.label, skill.skill.owner.repo_name),
        )
    root = _normalized_relative_path(
        skill.root,
        "root for skill {}".format(skill.name),
    )
    return _EntryInfo(
        destination = skill.name,
        files = [],
        kind = "symlink",
        name = skill.name,
        symlink_target = _relative_path(ctx.attr.discovery_dir, root),
    )

def _written_entry(ctx, target):
    skill = target[SkillInfo]
    _validate_name(skill.name)
    files = []
    skill_md = None
    for path in sorted(skill.files_by_path.keys()):
        file = skill.files_by_path[path]
        if path == "SKILL.md":
            skill_md = file
        files.append(struct(
            destination = path,
            file = file,
            path = _runfiles_path(ctx.workspace_name, file),
        ))
    if skill_md == None:
        fail("written entry for {} lacks SKILL.md".format(target.label))
    return _EntryInfo(
        destination = skill.name,
        files = files,
        kind = "written",
        name = skill.name,
        symlink_target = "",
    )

def _entries(ctx):
    discovery_dir = _normalized_relative_path(
        ctx.attr.discovery_dir,
        "discovery_dir",
    )
    if not ctx.attr.symlinks and not ctx.attr.archives:
        fail("skills_write requires at least one skill")

    entries = {}
    for target in ctx.attr.symlinks:
        entry = _symlink_entry(ctx, target)
        if entry.name in entries:
            fail("skills_write has duplicate skill name {} from {} and {}".format(
                entry.name,
                entries[entry.name].label,
                target.label,
            ))
        entries[entry.name] = struct(entry = entry, label = target.label)
    for target in ctx.attr.archives:
        entry = _written_entry(ctx, target)
        if entry.name in entries:
            fail("skills_write has duplicate skill name {} from {} and {}".format(
                entry.name,
                entries[entry.name].label,
                target.label,
            ))
        entries[entry.name] = struct(entry = entry, label = target.label)

    if discovery_dir.startswith(".."):
        fail("discovery_dir must be inside the consuming repository")
    for entry in entries.values():
        root = entry.entry.destination
        if root == discovery_dir or root.startswith(discovery_dir + "/"):
            fail("skill entry {} must not contain the discovery directory".format(root))

    return discovery_dir, [entries[name].entry for name in sorted(entries.keys())]

def _config_file(ctx, name, entries):
    """Writes the tab-separated entry and payload data read by the scripts."""
    lines = []
    for entry in entries:
        lines.append("\t".join([entry.kind, entry.name, entry.symlink_target]))
    lines.append("payloads")
    for index, entry in enumerate(entries):
        for payload in entry.files:
            lines.append("\t".join([
                str(index),
                payload.destination,
                payload.path,
            ]))
    config = ctx.actions.declare_file("{}.entries.txt".format(name))
    ctx.actions.write(output = config, content = "\n".join(lines) + "\n")
    return config

def _payload_files(entries):
    return depset([
        payload.file
        for entry in entries
        for payload in entry.files
    ])

def _skills_write_info(discovery_dir, entries):
    return SkillsWriteInfo(
        discovery_dir = discovery_dir,
        entries = entries,
    )

def _skills_write_updater_impl(ctx):
    discovery_dir, entries = _entries(ctx)
    config = _config_file(ctx, ctx.label.name, entries)
    executable = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.symlink(
        output = executable,
        target_file = ctx.executable._script,
        is_executable = True,
    )
    runfiles = ctx.runfiles(
        files = [config, ctx.executable._script],
        transitive_files = _payload_files(entries),
    )
    runfiles = runfiles.merge(
        ctx.attr._runfiles_library[DefaultInfo].default_runfiles,
    )
    runfiles = runfiles.merge(ctx.attr._script[DefaultInfo].default_runfiles)
    return [
        DefaultInfo(
            executable = executable,
            runfiles = runfiles,
        ),
        RunEnvironmentInfo(
            environment = {
                "SKILLS_WRITE_CONFIG": _runfiles_path(ctx.workspace_name, config),
                "SKILLS_WRITE_DISCOVERY_DIR": discovery_dir,
            },
        ),
        _skills_write_info(discovery_dir, entries),
    ]

skills_write_updater = rule(
    implementation = _skills_write_updater_impl,
    attrs = {
        "_runfiles_library": attr.label(
            default = Label("@bazel_tools//tools/bash/runfiles"),
        ),
        "_script": attr.label(
            cfg = "target",
            default = Label("//skill/internal/scripts:skills_write"),
            executable = True,
        ),
        "discovery_dir": attr.string(mandatory = True),
        "archives": attr.label_list(providers = [SkillInfo]),
        "symlinks": attr.label_list(providers = [SkillInfo]),
    },
    executable = True,
)

def _skills_write_check_test_impl(ctx):
    discovery_dir, entries = _entries(ctx)
    marker = ctx.file.workspace_marker
    if marker == None or not marker.is_source:
        fail("skills_write workspace_marker must be a source file")
    if marker.owner.repo_name != ctx.label.repo_name:
        fail(
            "skills_write workspace_marker must belong to the consuming " +
            "repository; {} is in repository {}".format(
                marker.owner,
                marker.owner.repo_name,
            ),
        )
    config = _config_file(ctx, ctx.label.name, entries)
    transitive = depset(
        [marker] + [
            file
            for target in ctx.attr.symlinks + ctx.attr.archives
            for file in target[SkillInfo].files.to_list()
        ],
        transitive = [_payload_files(entries)],
    )
    runfiles = ctx.runfiles(
        files = [config, ctx.executable._script],
        transitive_files = transitive,
    )
    runfiles = runfiles.merge(
        ctx.attr._runfiles_library[DefaultInfo].default_runfiles,
    )
    runfiles = runfiles.merge(ctx.attr._script[DefaultInfo].default_runfiles)
    executable = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.symlink(
        output = executable,
        target_file = ctx.executable._script,
        is_executable = True,
    )
    marker_path = _runfiles_path(ctx.workspace_name, marker)
    return [
        DefaultInfo(
            executable = executable,
            runfiles = runfiles,
        ),
        RunEnvironmentInfo(
            environment = {
                "SKILLS_WRITE_CONFIG": _runfiles_path(ctx.workspace_name, config),
                "SKILLS_WRITE_DISCOVERY_DIR": discovery_dir,
                "SKILLS_WRITE_MARKER_SUFFIX": "/{}".format(marker.short_path),
                "SKILLS_WRITE_WORKSPACE_MARKER": marker_path,
            },
        ),
        _skills_write_info(discovery_dir, entries),
    ]

skills_write_check_test = rule(
    implementation = _skills_write_check_test_impl,
    attrs = {
        "_runfiles_library": attr.label(
            default = Label("@bazel_tools//tools/bash/runfiles"),
        ),
        "_script": attr.label(
            cfg = "target",
            default = Label("//skill/internal/scripts:skills_check"),
            executable = True,
        ),
        "discovery_dir": attr.string(mandatory = True),
        "archives": attr.label_list(providers = [SkillInfo]),
        "symlinks": attr.label_list(providers = [SkillInfo]),
        "workspace_marker": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
    },
    test = True,
)

def skills_write(
        name,
        workspace_marker,
        symlinks = [],
        archives = [],
        discovery_dir = ".agents/skills",
        tags = None,
        testonly = False,
        visibility = None):
    """Declares a source-tree skill writer and its exact-state check.

    Args:
        name: Name of the runnable writer target.
        workspace_marker: Source file in the consuming repository used to
            locate the workspace root from test runfiles.
        symlinks: Skill labels installed as relative symlinks to their roots.
        archives: Skill labels whose files are copied as regular files.
        discovery_dir: Workspace-relative directory receiving the entries.
        tags: Optional tags applied to both generated targets.
        testonly: Whether the generated targets are test-only.
        visibility: Optional visibility applied to both generated targets.
    """
    if not symlinks and not archives:
        fail("skills_write requires at least one skill")

    common = {
        "discovery_dir": discovery_dir,
        "archives": archives,
        "symlinks": symlinks,
    }
    if visibility != None:
        common["visibility"] = visibility

    updater_attributes = dict(common)
    updater_attributes["tags"] = list(tags or [])
    updater_attributes["testonly"] = testonly
    skills_write_updater(
        name = name,
        **updater_attributes
    )

    test_attributes = dict(common)
    test_attributes["workspace_marker"] = workspace_marker
    test_tags = list(tags or [])
    for tag in [
        "external",
        "local",
        "no-cache",
        "no-remote",
        "no-sandbox",
    ]:
        if tag not in test_tags:
            test_tags.append(tag)
    test_attributes["tags"] = test_tags
    skills_write_check_test(
        name = name + "_test",
        **test_attributes
    )
