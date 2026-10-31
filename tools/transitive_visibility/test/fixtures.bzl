"""Fixtures distinguishing ordinary and execution-tool query edges."""

def _impl(ctx):
    return [DefaultInfo()]

edge_fixture = rule(
    implementation = _impl,
    attrs = {
        "deps": attr.label_list(),
        "tool": attr.label(cfg = "exec"),
    },
)
