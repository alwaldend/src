"""Repository-internal Molecule scenarios with declared QEMU inputs."""

load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo")
load("@rules_python//python:py_info.bzl", "PyInfo")
load("//tools/ansible/main/bzl:al_ansible_scripts.bzl", "AL_ANSIBLE_SCRIPTS")

def _key(ctx, file):
    path = file.short_path
    return path[3:] if path.startswith("../") else ctx.workspace_name + "/" + path

def _impl(ctx):
    executable = ctx.actions.declare_file(ctx.label.name)
    configuration = ctx.actions.declare_file(ctx.label.name + ".json")
    ctx.actions.symlink(output = executable, target_file = ctx.executable._runner, is_executable = True)
    mappings = {}
    inputs = []
    for package in ctx.attr.collections:
        for info, _ in package[PackageFilegroupInfo].pkg_files:
            for destination, source in info.dest_src_map.items():
                if destination in mappings:
                    fail("duplicate package destination: " + destination)
                mappings[destination] = _key(ctx, source)
                inputs.append(source)
    scenario = {}
    prefix = ctx.label.package + "/"
    for source in ctx.files.scenario:
        path = source.short_path.removeprefix(prefix)
        scenario[path] = _key(ctx, source)
    tools = {}
    dependencies = ctx.attr._tools.values() + [ctx.attr._runner]
    for name, target in ctx.attr._tools.items():
        tools[name] = _key(ctx, target[DefaultInfo].files_to_run.executable)
    config = {
        "label": str(ctx.label),
        "scenario": ctx.attr.scenario_name,
        "files": scenario,
        "mappings": mappings,
        "tools": tools,
        "image": _key(ctx, ctx.file.image),
        "bios": _key(ctx, ctx.file._bios),
        "lifecycle": {f.basename: _key(ctx, f) for f in ctx.files._lifecycle},
        "cpus": ctx.attr.cpus,
        "memory_mb": ctx.attr.memory_mb,
        "disk_gb": ctx.attr.disk_gb,
        "data_disks_gb": ctx.attr.data_disks_gb,
        "ports": ctx.attr.ports,
        "accelerator": ctx.attr.accelerator,
        "boot_timeout": ctx.attr.boot_timeout,
        "tcg_boot_timeout": ctx.attr.tcg_boot_timeout,
        "ansible_args": ctx.attr.ansible_args,
        "python_imports": ctx.attr._tools["ansible-playbook"][PyInfo].imports.to_list(),
        "services": ctx.attr.services,
    }
    ctx.actions.write(configuration, json.encode_indent(config))
    runfiles = ctx.runfiles(files = inputs + ctx.files.scenario + ctx.files._lifecycle + [configuration, ctx.file.image, ctx.file._bios])
    runfiles = runfiles.merge_all([d[DefaultInfo].default_runfiles for d in dependencies])
    return [DefaultInfo(executable = executable, runfiles = runfiles)]

_molecule_test = rule(
    implementation = _impl,
    test = True,
    attrs = {
        "scenario": attr.label_list(allow_files = True, mandatory = True),
        "scenario_name": attr.string(default = "default"),
        "collections": attr.label_list(providers = [PackageFilegroupInfo]),
        "image": attr.label(default = "//third_party/org_fedora_cloud", allow_single_file = True),
        "cpus": attr.int(default = 2),
        "memory_mb": attr.int(default = 2048),
        "disk_gb": attr.int(default = 12),
        "data_disks_gb": attr.int_list(),
        "ports": attr.int_list(default = [22]),
        "accelerator": attr.string(default = "auto", values = ["auto", "kvm", "tcg"]),
        "boot_timeout": attr.int(default = 240),
        "tcg_boot_timeout": attr.int(default = 600),
        "ansible_args": attr.string_list(),
        "services": attr.string_list(default = ["sshd"]),
        "_runner": attr.label(default = "//tools/molecule/cmd/molecule_runner", executable = True, cfg = "exec"),
        "_bios": attr.label(default = "//third_party/com_github_hermeticbuild_qemu:bios", allow_single_file = True),
        "_lifecycle": attr.label(default = "//tools/molecule/internal:lifecycle"),
        "_tools": attr.string_keyed_label_dict(
            cfg = "exec",
            default = {
                "molecule": "//tools/ansible:molecule",
                "qemu": "//third_party/com_github_hermeticbuild_qemu:qemu_system_x86_64",
                "qemu-img": "//third_party/com_github_hermeticbuild_qemu:qemu_img",
                "seed": "//tools/molecule/internal:seed",
            } | {cli.replace("_", "-"): "//tools/ansible:" + cli for cli in AL_ANSIBLE_SCRIPTS},
        ),
    },
)

def molecule_test(name, tags = [], **kwargs):
    """Run a fresh scenario locally; package installation requires networking.

    Args:
        name: Test target name.
        tags: Additional test tags.
        **kwargs: Scenario, package, guest, and resource declarations.
    """
    _molecule_test(
        name = name,
        tags = ["manual", "external", "no-remote", "requires-network"] + tags,
        target_compatible_with = ["@platforms//os:linux", "@platforms//cpu:x86_64"],
        timeout = "eternal",
        **kwargs
    )
