"""Rules for reconciling source-tree skill discovery symlinks."""

load(":skill_library.bzl", "SkillInfo")

SkillDiscoveryLinksInfo = provider(
    doc = "The deterministic source-tree links derived from skill libraries.",
    fields = {
        "discovery_dir": "The normalized workspace-relative discovery directory.",
        "names": "Skill names in deterministic order.",
        "roots": "Canonical workspace-relative skill roots by name order.",
        "targets": "Relative symlink targets by name order.",
    },
)

_SKILL_NAME_CHARACTERS = [
    "-",
    "0",
    "1",
    "2",
    "3",
    "4",
    "5",
    "6",
    "7",
    "8",
    "9",
    "a",
    "b",
    "c",
    "d",
    "e",
    "f",
    "g",
    "h",
    "i",
    "j",
    "k",
    "l",
    "m",
    "n",
    "o",
    "p",
    "q",
    "r",
    "s",
    "t",
    "u",
    "v",
    "w",
    "x",
    "y",
    "z",
]

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

def _validate_skill_name(name):
    invalid = not name
    for index in range(len(name)):
        if name[index] not in _SKILL_NAME_CHARACTERS:
            invalid = True
            break
    if invalid:
        fail(
            (
                "skill discovery link name must contain only lowercase ASCII " +
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

def _link_entries(ctx):
    discovery_dir = _normalized_relative_path(
        ctx.attr.discovery_dir,
        "discovery_dir",
    )
    if not ctx.attr.skills:
        fail("skill_discovery_links requires at least one skill")

    entries_by_name = {}
    for target in ctx.attr.skills:
        if target.label.repo_name != ctx.label.repo_name:
            fail(
                "skill discovery links can only target skills in the same " +
                "source repository: {} is in repository {}".format(
                    target.label,
                    target.label.repo_name,
                ),
            )

        skill = target[SkillInfo]
        _validate_skill_name(skill.name)
        root = _normalized_relative_path(
            skill.root,
            "root for skill {}".format(skill.name),
        )
        if not skill.skill.is_source:
            fail(
                "skill discovery links require a source SKILL.md; {} is " +
                "generated".format(target.label),
            )
        if (
            root == discovery_dir or
            root.startswith(discovery_dir + "/") or
            discovery_dir.startswith(root + "/")
        ):
            fail(
                "skill root {} and discovery directory {} must not contain " +
                "one another".format(root, discovery_dir),
            )
        if skill.name in entries_by_name:
            fail(
                "skill discovery links have duplicate name {} from {} and {}".format(
                    skill.name,
                    entries_by_name[skill.name].label,
                    target.label,
                ),
            )
        entries_by_name[skill.name] = struct(
            label = target.label,
            name = skill.name,
            root = root,
            skill_file = skill.skill,
            target = _relative_path(discovery_dir, root),
        )

    return (
        discovery_dir,
        [entries_by_name[name] for name in sorted(entries_by_name.keys())],
    )

def _shell_quote(value):
    return "'{}'".format(value.replace("'", "'\"'\"'"))

def _array_assignment(name, values):
    return "{}=({})".format(
        name,
        " ".join([_shell_quote(value) for value in values]),
    )

def _updater_script(discovery_dir, entries):
    lines = [
        "#!/usr/bin/env bash",
        "set -euo pipefail",
        "",
        "workspace=\"${BUILD_WORKSPACE_DIRECTORY:-}\"",
        "if [[ -z \"${workspace}\" || \"${workspace}\" != /* ]]; then",
        "    echo \"BUILD_WORKSPACE_DIRECTORY must be an absolute path\" >&2",
        "    exit 1",
        "fi",
        "discovery_relative={}".format(_shell_quote(discovery_dir)),
        "discovery=\"${workspace}/${discovery_relative}\"",
        _array_assignment("expected_names", [entry.name for entry in entries]),
        _array_assignment("expected_roots", [entry.root for entry in entries]),
        _array_assignment("expected_targets", [entry.target for entry in entries]),
        "",
        "expected_index() {",
        "    local candidate=\"$1\"",
        "    local index",
        "    for index in \"${!expected_names[@]}\"; do",
        "        if [[ \"${expected_names[$index]}\" == \"${candidate}\" ]]; then",
        "            printf '%s\\n' \"${index}\"",
        "            return 0",
        "        fi",
        "    done",
        "    return 1",
        "}",
        "",
        "for index in \"${!expected_roots[@]}\"; do",
        "    canonical=\"${workspace}/${expected_roots[$index]}\"",
        "    current=\"${workspace}\"",
        "    IFS='/' read -r -a root_parts <<< \"${expected_roots[$index]}\"",
        "    for root_part in \"${root_parts[@]}\"; do",
        "        current=\"${current}/${root_part}\"",
        "        if [[ -L \"${current}\" ]]; then",
        "            echo \"canonical skill path has a symlink ancestor: ${current}\" >&2",
        "            exit 1",
        "        fi",
        "    done",
        "    if [[ ! -d \"${canonical}\" || ! -f \"${canonical}/SKILL.md\" ]]; then",
        "        echo \"canonical skill root lacks SKILL.md: ${canonical}\" >&2",
        "        exit 1",
        "    fi",
        "done",
        "",
        "current=\"${workspace}\"",
        "IFS='/' read -r -a discovery_parts <<< \"${discovery_relative}\"",
        "for ((index = 0; index < ${#discovery_parts[@]} - 1; index++)); do",
        "    current=\"${current}/${discovery_parts[$index]}\"",
        "    if [[ -L \"${current}\" ]]; then",
        "        echo \"discovery parent must not be a symlink: ${current}\" >&2",
        "        exit 1",
        "    fi",
        "    if [[ -e \"${current}\" && ! -d \"${current}\" ]]; then",
        "        echo \"discovery parent must be a directory: ${current}\" >&2",
        "        exit 1",
        "    fi",
        "    if [[ ! -e \"${current}\" ]]; then",
        "        mkdir \"${current}\"",
        "    fi",
        "done",
        "",
        "lock=\"${discovery}.lock\"",
        "if ! mkdir \"${lock}\" 2>/dev/null; then",
        "    echo \"another updater holds the discovery lock: ${lock}\" >&2",
        "    exit 1",
        "fi",
        "trap 'rmdir \"${lock}\"' EXIT",
        "",
        "if [[ -L \"${discovery}\" ]]; then",
        "    unlink \"${discovery}\"",
        "elif [[ -e \"${discovery}\" && ! -d \"${discovery}\" ]]; then",
        "    echo \"discovery path must be a directory or symlink: ${discovery}\" >&2",
        "    exit 1",
        "fi",
        "if [[ ! -e \"${discovery}\" ]]; then",
        "    mkdir \"${discovery}\"",
        "fi",
        "",
        "shopt -s dotglob nullglob",
        "existing_entries=(\"${discovery}\"/*)",
        "shopt -u dotglob nullglob",
        "for entry in \"${existing_entries[@]}\"; do",
        "    entry_name=\"${entry##*/}\"",
        "    if [[ ! -L \"${entry}\" ]]; then",
        "        if expected_index \"${entry_name}\" >/dev/null; then",
        "            echo \"managed skill entry is not a symlink: ${entry}\" >&2",
        "        else",
        "            echo \"unexpected non-symlink in discovery directory: ${entry}\" >&2",
        "        fi",
        "        exit 1",
        "    fi",
        "done",
        "",
        "for entry in \"${existing_entries[@]}\"; do",
        "    entry_name=\"${entry##*/}\"",
        "    if ! expected_index \"${entry_name}\" >/dev/null; then",
        "        unlink \"${entry}\"",
        "    fi",
        "done",
        "",
        "for index in \"${!expected_names[@]}\"; do",
        "    entry=\"${discovery}/${expected_names[$index]}\"",
        "    expected=\"${expected_targets[$index]}\"",
        "    if [[ -L \"${entry}\" ]]; then",
        "        actual=\"$(readlink \"${entry}\")\"",
        "        if [[ \"${actual}\" == \"${expected}\" ]]; then",
        "            continue",
        "        fi",
        "        unlink \"${entry}\"",
        "    fi",
        "    ln -s \"${expected}\" \"${entry}\"",
        "done",
        "",
        "printf 'reconciled %s skill discovery links in %s\\n' \\",
        "    \"${#expected_names[@]}\" \"${discovery_relative}\"",
        "",
    ]
    return "\n".join(lines)

def _checker_script(discovery_dir, entries):
    marker = entries[0]
    lines = [
        "#!/usr/bin/env bash",
        "set -euo pipefail",
        "",
        "if [[ -n \"${RUNFILES_DIR:-}\" ]]; then",
        "    runfiles_library=\"${RUNFILES_DIR}/bazel_tools/tools/bash/runfiles/runfiles.bash\"",
        "elif [[ -n \"${RUNFILES_MANIFEST_FILE:-}\" ]]; then",
        "    runfiles_library=\"$(grep -m1 '^bazel_tools/tools/bash/runfiles/runfiles.bash ' \"${RUNFILES_MANIFEST_FILE}\" | cut -d ' ' -f 2-)\"",
        "else",
        "    echo 'Bazel runfiles are unavailable' >&2",
        "    exit 1",
        "fi",
        "if [[ ! -f \"${runfiles_library}\" ]]; then",
        "    echo \"Bazel runfiles library is unavailable: ${runfiles_library}\" >&2",
        "    exit 1",
        "fi",
        "source \"${runfiles_library}\"",
        "marker_relative={}".format(_shell_quote(marker.skill_file.short_path)),
        "marker=\"$(rlocation \"${TEST_WORKSPACE}/${marker_relative}\")\"",
        "if [[ -z \"${marker}\" ]]; then",
        "    echo \"source SKILL.md is absent from runfiles: ${marker_relative}\" >&2",
        "    exit 1",
        "fi",
        "resolve_path() {",
        "    local path=\"$1\"",
        "    local directory name target",
        "    local step",
        "    for ((step = 0; step < 64; step++)); do",
        "        directory=\"$(cd -P \"$(dirname \"${path}\")\" && pwd)\" || return 1",
        "        name=\"$(basename \"${path}\")\"",
        "        path=\"${directory}/${name}\"",
        "        if [[ ! -L \"${path}\" ]]; then",
        "            printf '%s\\n' \"${path}\"",
        "            return 0",
        "        fi",
        "        target=\"$(readlink \"${path}\")\" || return 1",
        "        if [[ \"${target}\" == /* ]]; then",
        "            path=\"${target}\"",
        "        else",
        "            path=\"${directory}/${target}\"",
        "        fi",
        "    done",
        "    echo \"too many source SKILL.md symlinks: $1\" >&2",
        "    return 1",
        "}",
        "if ! marker=\"$(resolve_path \"${marker}\")\"; then",
        "    echo \"cannot resolve source SKILL.md runfile: ${marker}\" >&2",
        "    exit 1",
        "fi",
        "marker_suffix={}".format(_shell_quote("/{}/SKILL.md".format(marker.root))),
        "if [[ \"${marker}\" != *\"${marker_suffix}\" ]]; then",
        "    echo \"could not locate source workspace from ${marker}\" >&2",
        "    exit 1",
        "fi",
        "workspace=\"${marker%\"${marker_suffix}\"}\"",
        "if [[ ! -f \"${workspace}/MODULE.bazel\" ]]; then",
        "    echo \"source SKILL.md did not resolve into a Bazel workspace: ${marker}\" >&2",
        "    exit 1",
        "fi",
        "discovery_relative={}".format(_shell_quote(discovery_dir)),
        "discovery=\"${workspace}/${discovery_relative}\"",
        _array_assignment("expected_names", [entry.name for entry in entries]),
        _array_assignment("expected_roots", [entry.root for entry in entries]),
        _array_assignment("expected_targets", [entry.target for entry in entries]),
        "",
        "expected_index() {",
        "    local candidate=\"$1\"",
        "    local index",
        "    for index in \"${!expected_names[@]}\"; do",
        "        if [[ \"${expected_names[$index]}\" == \"${candidate}\" ]]; then",
        "            printf '%s\\n' \"${index}\"",
        "            return 0",
        "        fi",
        "    done",
        "    return 1",
        "}",
        "",
        "current=\"${workspace}\"",
        "IFS='/' read -r -a discovery_parts <<< \"${discovery_relative}\"",
        "for index in \"${!discovery_parts[@]}\"; do",
        "    current=\"${current}/${discovery_parts[$index]}\"",
        "    if [[ -L \"${current}\" ]]; then",
        "        echo \"skill discovery path has a symlink ancestor: ${current}\" >&2",
        "        exit 1",
        "    fi",
        "done",
        "",
        "if [[ -L \"${discovery}\" || ! -d \"${discovery}\" ]]; then",
        "    echo \"skill discovery path is not a real directory: ${discovery}\" >&2",
        "    exit 1",
        "fi",
        "",
        "actual_count=0",
        "shopt -s dotglob nullglob",
        "existing_entries=(\"${discovery}\"/*)",
        "shopt -u dotglob nullglob",
        "for entry in \"${existing_entries[@]}\"; do",
        "    actual_count=$((actual_count + 1))",
        "    entry_name=\"${entry##*/}\"",
        "    if ! expected_index \"${entry_name}\" >/dev/null; then",
        "        echo \"unexpected skill discovery entry: ${entry_name}\" >&2",
        "        exit 1",
        "    fi",
        "done",
        "if [[ \"${actual_count}\" -ne \"${#expected_names[@]}\" ]]; then",
        "    echo \"skill discovery entry count is ${actual_count}, expected ${#expected_names[@]}\" >&2",
        "    exit 1",
        "fi",
        "",
        "for index in \"${!expected_names[@]}\"; do",
        "    entry=\"${discovery}/${expected_names[$index]}\"",
        "    if [[ ! -L \"${entry}\" ]]; then",
        "        echo \"skill discovery entry is not a symlink: ${entry}\" >&2",
        "        exit 1",
        "    fi",
        "    actual=\"$(readlink \"${entry}\")\"",
        "    expected=\"${expected_targets[$index]}\"",
        "    if [[ \"${actual}\" != \"${expected}\" ]]; then",
        "        echo \"${entry} points to ${actual}, expected ${expected}\" >&2",
        "        exit 1",
        "    fi",
        "    canonical=\"${workspace}/${expected_roots[$index]}\"",
        "    current=\"${workspace}\"",
        "    IFS='/' read -r -a root_parts <<< \"${expected_roots[$index]}\"",
        "    for root_part in \"${root_parts[@]}\"; do",
        "        current=\"${current}/${root_part}\"",
        "        if [[ -L \"${current}\" ]]; then",
        "            echo \"canonical skill path has a symlink ancestor: ${current}\" >&2",
        "            exit 1",
        "        fi",
        "    done",
        "    if [[ ! -d \"${canonical}\" || ! -f \"${canonical}/SKILL.md\" ]]; then",
        "        echo \"canonical skill root lacks SKILL.md: ${canonical}\" >&2",
        "        exit 1",
        "    fi",
        "done",
        "",
        "printf 'verified %s skill discovery links in %s\\n' \\",
        "    \"${#expected_names[@]}\" \"${discovery_relative}\"",
        "",
    ]
    return "\n".join(lines)

def _discovery_info(discovery_dir, entries):
    return SkillDiscoveryLinksInfo(
        discovery_dir = discovery_dir,
        names = [entry.name for entry in entries],
        roots = [entry.root for entry in entries],
        targets = [entry.target for entry in entries],
    )

def _skill_discovery_links_updater_impl(ctx):
    discovery_dir, entries = _link_entries(ctx)
    executable = ctx.actions.declare_file(ctx.label.name + ".sh")
    ctx.actions.write(
        output = executable,
        content = _updater_script(discovery_dir, entries),
        is_executable = True,
    )
    return [
        DefaultInfo(executable = executable),
        _discovery_info(discovery_dir, entries),
    ]

skill_discovery_links_updater = rule(
    implementation = _skill_discovery_links_updater_impl,
    attrs = {
        "discovery_dir": attr.string(mandatory = True),
        "skills": attr.label_list(
            mandatory = True,
            providers = [SkillInfo],
        ),
    },
    executable = True,
)

def _skill_discovery_links_test_impl(ctx):
    discovery_dir, entries = _link_entries(ctx)
    executable = ctx.actions.declare_file(ctx.label.name + ".sh")
    ctx.actions.write(
        output = executable,
        content = _checker_script(discovery_dir, entries),
        is_executable = True,
    )
    runfiles = ctx.runfiles(files = [entries[0].skill_file])
    runfiles = runfiles.merge(
        ctx.attr._runfiles_library[DefaultInfo].default_runfiles,
    )
    return [
        DefaultInfo(
            executable = executable,
            runfiles = runfiles,
        ),
        _discovery_info(discovery_dir, entries),
    ]

skill_discovery_links_test = rule(
    implementation = _skill_discovery_links_test_impl,
    attrs = {
        "_runfiles_library": attr.label(
            default = Label(
                "@bazel_tools//tools/bash/runfiles",
            ),
        ),
        "discovery_dir": attr.string(mandatory = True),
        "skills": attr.label_list(
            mandatory = True,
            providers = [SkillInfo],
        ),
    },
    test = True,
)

def skill_discovery_links(
        name,
        skills,
        discovery_dir = ".agents/skills",
        tags = None,
        testonly = False,
        visibility = None):
    """Declares a source-tree link updater and its exact-state test.

    Args:
        name: Name of the runnable updater target.
        skills: skill_library labels that are the complete discovery set.
        discovery_dir: Workspace-relative directory containing direct links.
        tags: Optional tags applied to both generated targets.
        testonly: Whether the generated targets are test-only.
        visibility: Optional visibility applied to both generated targets.
    """
    common = {
        "discovery_dir": discovery_dir,
        "skills": skills,
        "tags": list(tags or []),
    }
    if visibility != None:
        common["visibility"] = visibility

    updater_attributes = dict(common)
    updater_attributes["testonly"] = testonly
    skill_discovery_links_updater(
        name = name,
        **updater_attributes
    )
    test_attributes = dict(common)
    test_tags = list(common["tags"])
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
    skill_discovery_links_test(
        name = name + "_test",
        **test_attributes
    )
