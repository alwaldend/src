load("@rules_hugo//pkg/bzl:al_hugo_site.bzl", "al_hugo_site")
load("@rules_hugo//pkg/bzl:al_hugo_worker.bzl", "al_hugo_worker")
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

def _readme_impl(ctx):
    candidates = []
    for file in ctx.attr.src[DefaultInfo].files.to_list():
        if file.basename == "README.md" and file.short_path.startswith("projects/{}/".format(ctx.attr.project)):
            candidates.append(file)
    if not candidates:
        for file in ctx.attr.src[DefaultInfo].files.to_list():
            if file.basename == "README.md" and "/projects/{}/".format(ctx.attr.project) in file.short_path:
                candidates.append(file)
    if not candidates:
        candidates = [
            file
            for file in ctx.attr.src[DefaultInfo].files.to_list()
            if file.basename == "README.md"
        ]
    if not candidates:
        fail("could not find README.md for project {}".format(ctx.attr.project))
    readme = candidates[0]
    out = ctx.actions.declare_file("{}/_index.md".format(ctx.attr.name))
    ctx.actions.run_shell(
        inputs = [readme],
        outputs = [out],
        command = "sed 's/{{.*}}//' '{}' > '{}'".format(readme.path, out.path),
    )
    return [DefaultInfo(files = depset([out]))]

_readme = rule(
    implementation = _readme_impl,
    attrs = {
        "project": attr.string(mandatory = True),
        "src": attr.label(mandatory = True),
    },
)

def _write_config_impl(ctx):
    content = """\
[module.hugoVersion]
extended = true
min = "0.160.1"

disableKinds = ["taxonomy", "term"]
title = "{}"
baseURL = "/"
theme = ["github.com/google/docsy"]
enableRobotsTXT = true

[outputs]
home = ["html"]
page = ["html"]

[params]
github_repo = "{}"
github_branch = "master"

[[params.alwaldend.links]]
name = "Docs"
url = "{}"

[[params.alwaldend.links]]
name = "Github"
url = "{}"

[params.ui]
navbar_logo = false
showLightDarkModeMenu = true
sidebar_menu_compact = true

[markup.goldmark.renderer]
unsafe = true

[markup.goldmark.renderHooks.link]
useEmbedded = "always"

[markup.goldmark.parser.attribute]
block = true
title = true

[[module.mounts]]
source = "themes/github.com/twbs/bootstrap/scss"
target = "assets/vendor/bootstrap/scss"
[[module.mounts]]
source = "themes/github.com/twbs/bootstrap/js"
target = "assets/vendor/bootstrap/js"
[[module.mounts]]
source = "themes/github.com/twbs/bootstrap/dist"
target = "assets/vendor/bootstrap/dist"
[[module.mounts]]
source = "themes/github.com/FortAwesome/Font-Awesome/scss"
target = "assets/vendor/Font-Awesome/scss"
[[module.mounts]]
source = "themes/github.com/FortAwesome/Font-Awesome/webfonts"
target = "static/webfonts"
""".format(
        ctx.attr.title.replace('"', '\\"'),
        ctx.attr.repository_url.replace('"', '\\"'),
        ctx.attr.docs_url.replace('"', '\\"'),
        ctx.attr.repository_url.replace('"', '\\"'),
    )
    out = ctx.actions.declare_file("{}/hugo.toml".format(ctx.attr.name))
    ctx.actions.write(output = out, content = content)
    return [DefaultInfo(files = depset([out]))]

_write_config = rule(
    implementation = _write_config_impl,
    attrs = {
        "docs_url": attr.string(mandatory = True),
        "repository_url": attr.string(mandatory = True),
        "title": attr.string(mandatory = True),
    },
)

def al_hugo_landing(name, project, title, docs, docs_url, repository_url, **kwargs):
    """Creates a Hugo landing site archive from a project README."""
    _write_config(
        name = "{}_config".format(name),
        docs_url = docs_url,
        repository_url = repository_url,
        title = title,
    )

    _readme(
        name = "{}_readme".format(name),
        project = project,
        src = docs,
    )

    pkg_files(
        name = "{}_readme_file".format(name),
        srcs = [":{}_readme".format(name)],
        prefix = "content",
    )

    pkg_files(
        name = "{}_config_file".format(name),
        srcs = [":{}_config".format(name)],
        prefix = ".",
    )

    pkg_filegroup(
        name = "{}_source".format(name),
        srcs = [
            "//projects/hugo_landing/assets/scss:assets",
            "//projects/hugo_landing/layouts:layouts",
            ":{}_config_file".format(name),
            ":{}_readme_file".format(name),
            "@com_alwaldend_src_hugo_github_com_google_docsy",
            "@com_alwaldend_src_hugo_github_com_twbs_bootstrap",
            "@com_alwaldend_src_hugo_github_com_fortawesome_font_awesome",
        ],
    )

    pkg_tar(
        name = name,
        srcs = [":{}_source".format(name)],
        strip_prefix = ".",
        **kwargs
    )

def al_hugo_landing_site(name, project, title, docs, docs_url, repository_url, **kwargs):
    """Creates and builds a reusable Hugo landing site for a project."""
    al_hugo_landing(
        name = "{}_source".format(name),
        docs = docs,
        docs_url = docs_url,
        project = project,
        repository_url = repository_url,
        title = title,
    )
    al_hugo_site(
        name = "{}_site".format(name),
        postcss = "//tools/postcss",
        site = ":{}_source".format(name),
        tools = [
            "//tools/sass:dart-sass",
            "@com_alwaldend_src_hugo_github_com_fortawesome_font_awesome",
            "@com_alwaldend_src_hugo_github_com_google_docsy",
            "@com_alwaldend_src_hugo_github_com_twbs_bootstrap",
        ],
    )
    al_hugo_worker(
        name = name,
        arguments = [
            "build",
            "--logLevel",
            "debug",
            "--printPathWarnings",
        ],
        out_dir = "{}_site.destination".format(name),
        site = ":{}_site".format(name),
        worker = "@rules_hugo//cmd/hugo_worker",
        **kwargs
    )

def al_hugo_landing_sites(name, projects):
    """Creates rendered landing sites for every listed project."""
    targets = []
    for project in projects:
        if project.startswith("rules_"):
            docs = "@rules_{}//:docs".format(project.removeprefix("rules_"))
        else:
            docs = "//projects/{}:docs".format(project)
        al_hugo_landing_site(
            name = project,
            docs = docs,
            docs_url = "https://alwaldend.com/docs/projects/{}/".format(project),
            project = project,
            repository_url = "https://github.com/alwaldend/src",
            title = project.replace("_", " ").title(),
        )
        targets.append(":{}".format(project))

    native.filegroup(
        name = name,
        srcs = targets,
    )
