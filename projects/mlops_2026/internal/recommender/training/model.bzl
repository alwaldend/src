"""Build an upstream recommendation model with explicit data and output."""

def _course_model_impl(ctx):
    output = ctx.outputs.out
    args = ctx.actions.args()
    args.add("--trainer", ctx.executable.trainer)
    args.add("--dataset", ctx.file.dataset)
    args.add("--data", ctx.file.data)
    args.add("--output", output)
    ctx.actions.run(
        executable = ctx.executable._runner,
        arguments = [args],
        inputs = [ctx.file.dataset, ctx.file.data],
        tools = [ctx.attr.trainer[DefaultInfo].files_to_run],
        outputs = [output],
        mnemonic = "TrainCourseModel",
        progress_message = "Training course recommender model",
    )
    return [DefaultInfo(files = depset([output]))]

course_model = rule(
    implementation = _course_model_impl,
    attrs = {
        "data": attr.label(allow_single_file = True, mandatory = True),
        "dataset": attr.label(allow_single_file = True, mandatory = True),
        "out": attr.output(mandatory = True),
        "trainer": attr.label(executable = True, cfg = "exec", mandatory = True),
        "_runner": attr.label(
            default = Label("//projects/mlops_2026/internal/recommender/training:build_model"),
            executable = True,
            cfg = "exec",
        ),
    },
)
