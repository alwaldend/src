"""Analysis tests for skill_library and skill_validation_aspect."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")
load("//skill:defs.bzl", "SkillInfo", "skill_validation_aspect")

def _package_path(file):
    return "{}/{}".format(file.owner.package, file.owner.name)

def _package_paths(files):
    return sorted([_package_path(file) for file in files.to_list()])

def _skill_library_files_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    skill = target[SkillInfo]
    expected = [
        "tests/analysis/SKILL.md",
        "tests/analysis/agents/openai.yaml",
        "tests/analysis/reference.txt",
    ]

    asserts.equals(env, expected, _package_paths(skill.files))
    asserts.equals(env, expected, _package_paths(target[DefaultInfo].files))
    asserts.equals(env, "tests/analysis/SKILL.md", _package_path(skill.skill))
    asserts.equals(
        env,
        "tests/analysis/agents/openai.yaml",
        _package_path(skill.openai_yaml),
    )
    return analysistest.end(env)

skill_library_files_test = analysistest.make(
    _skill_library_files_test_impl,
)

def _skill_library_without_openai_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    skill = target[SkillInfo]

    asserts.equals(env, None, skill.openai_yaml)
    asserts.equals(
        env,
        [
            "tests/analysis/SKILL.md",
            "tests/analysis/reference.txt",
        ],
        _package_paths(skill.files),
    )
    return analysistest.end(env)

skill_library_without_openai_test = analysistest.make(
    _skill_library_without_openai_test_impl,
)

def _skill_library_generated_files_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    skill = target[SkillInfo]

    asserts.equals(env, 3, len(skill.files.to_list()))
    asserts.equals(env, "SKILL.md", skill.skill.basename)
    asserts.true(
        env,
        skill.openai_yaml.short_path.endswith("/agents/openai.yaml"),
    )
    return analysistest.end(env)

skill_library_generated_files_test = analysistest.make(
    _skill_library_generated_files_test_impl,
)

def _skill_library_failure_test_impl(ctx):
    env = analysistest.begin(ctx)
    asserts.expect_failure(env, ctx.attr.expected_message)
    return analysistest.end(env)

skill_library_failure_test = analysistest.make(
    _skill_library_failure_test_impl,
    attrs = {
        "expected_message": attr.string(mandatory = True),
    },
    expect_failure = True,
)

def _skill_library_aspect_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    outputs = target[OutputGroupInfo].skill_validation.to_list()

    asserts.equals(env, 1, len(outputs))
    asserts.equals(env, "library.skill_validation", outputs[0].basename)
    return analysistest.end(env)

skill_library_aspect_test = analysistest.make(
    _skill_library_aspect_test_impl,
    extra_target_under_test_aspects = [skill_validation_aspect],
)
