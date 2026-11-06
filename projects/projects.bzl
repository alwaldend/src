"""Projects that publish a landing page in the main site.

The alwaldend.com project is already served by the apex website and is excluded
from this list.
"""

PROJECTS = [
    "activitywatch_ingester_android",
    "al",
    "android_launcher",
    "anki_as_code",
    "anki_grpc",
    "ansible_collection",
    "autoscroll",
    "bazel_deps_analyzer",
    "bazel_gazelle_enforce_transitive_visibility",
    "ci_platform",
    "dotfiles",
    "icon_generator",
    "infinitime",
    "kustomization",
    "leetcode_downloader",
    "llama_responses_bridge",
    "mcp_cordis",
    "mermaid_3d",
    "nexus_security_plugin",
    "renders",
    "rules_trufflehog",
    "sri",
    "tf_modules",
    "un_no_excuses_campaign",
    "useless_qt_gui",
    "vastai_terraform_provider",
    "x_article_uploader",
]

# Standalone modules expose their landing content through their Bzlmod name.
PROJECT_MODULES = {
    "al": "rules_al",
    "bazel_gazelle_enforce_transitive_visibility": "bazel_gazelle_enforce_transitive_visibility",
    "rules_trufflehog": "rules_trufflehog",
}
