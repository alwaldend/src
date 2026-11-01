"""Repository release pages using the site's shared Hugo presentation."""

load("@rules_al//:defs.bzl", "release")

_PAGE_TEMPLATE = """\
---
title: {release_name}
description: Release {release_name}
toc_hide: true
hide_summary: true
tags:
    - generated
    - release
---

{{{{< alwaldend/release standalone=true >}}}}
"""

def release_page(name, **kwargs):
    """Package a project release page with repository history and site markup.

    Args:
        name: Bazel target name.
        **kwargs: Remaining release rule arguments.
    """
    release(
        name = name,
        git_bundle = "@repository_git//:release_git_state",
        page_template = _PAGE_TEMPLATE,
        **kwargs
    )
