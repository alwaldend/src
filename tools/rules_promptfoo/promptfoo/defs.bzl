"""Public Promptfoo test rules."""

load(
    "//promptfoo/private:promptfoo_test.bzl",
    _promptfoo_test = "promptfoo_test",
    _promptfoo_validate_test = "promptfoo_validate_test",
)

promptfoo_test = _promptfoo_test
promptfoo_validate_test = _promptfoo_validate_test
