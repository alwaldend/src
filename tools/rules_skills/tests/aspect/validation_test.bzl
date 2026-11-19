"""Execution test for the skill validation aspect."""

load("//skill:defs.bzl", "SkillInfo", "skill_validation_aspect")

def _skill_validation_test_impl(ctx):
    validation_files = ctx.attr.target[OutputGroupInfo].skill_validation
    executable = ctx.actions.declare_file(ctx.label.name + ".sh")
    checks = [
        "test -s \"${{TEST_SRCDIR}}/${{TEST_WORKSPACE}}/{}\"".format(
            validation_file.short_path,
        )
        for validation_file in validation_files.to_list()
    ]
    ctx.actions.write(
        output = executable,
        content = "#!/usr/bin/env bash\nset -euo pipefail\n{}\n".format(
            "\n".join(checks),
        ),
        is_executable = True,
    )
    return [
        DefaultInfo(
            executable = executable,
            runfiles = ctx.runfiles(transitive_files = validation_files),
        ),
    ]

skill_validation_test = rule(
    implementation = _skill_validation_test_impl,
    attrs = {
        "target": attr.label(
            aspects = [skill_validation_aspect],
            mandatory = True,
            providers = [SkillInfo],
        ),
    },
    test = True,
)
