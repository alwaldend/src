"""Declare the taxi-fare estimator and evaluation metadata as Bazel outputs."""

def _taxi_model_impl(ctx):
    model = ctx.actions.declare_file(ctx.label.name + "/model.pkl")
    metrics = ctx.actions.declare_file(ctx.label.name + "/metrics.json")
    args = ctx.actions.args()
    args.add("--dataset", ctx.file.dataset)
    args.add("--output-dir", model.dirname)
    args.add("--random-state", ctx.attr.random_state)
    args.add("--test-size", ctx.attr.test_size)
    ctx.actions.run(
        executable = ctx.executable.trainer,
        arguments = [args],
        inputs = [ctx.file.dataset],
        outputs = [model, metrics],
        env = {"OPENBLAS_NUM_THREADS": "1", "OMP_NUM_THREADS": "1", "MKL_NUM_THREADS": "1"},
        mnemonic = "TrainTaxiFareModel",
        progress_message = "Training taxi-fare model on the pinned upstream CSV",
    )
    return [DefaultInfo(files = depset([model, metrics]))]

taxi_model = rule(
    implementation = _taxi_model_impl,
    attrs = {
        "dataset": attr.label(allow_single_file = [".csv"], mandatory = True),
        "trainer": attr.label(executable = True, cfg = "exec", mandatory = True),
        "random_state": attr.int(default = 42),
        "test_size": attr.string(default = "0.2"),
    },
)
