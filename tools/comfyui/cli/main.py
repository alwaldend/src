"""Run the pinned upstream CLI with private repository settings and no telemetry."""

import os
import pathlib

import comfy_cli.constants


def _main() -> None:
    """Select writable task settings before upstream imports initialize them."""
    workspace: pathlib.Path = pathlib.Path(
        os.environ.get("BUILD_WORKSPACE_DIRECTORY", os.getcwd())
    )
    directory: pathlib.Path = pathlib.Path(
        os.environ.get(
            "COMFY_REPOSITORY_PROFILE",
            str(workspace / "out/comfyui/cli"),
        )
    ).resolve()
    directory.mkdir(parents=True, exist_ok=True, mode=0o700)
    platform: comfy_cli.constants.OS
    for platform in comfy_cli.constants.DEFAULT_CONFIG:
        comfy_cli.constants.DEFAULT_CONFIG[platform] = str(directory)
    os.environ["DO_NOT_TRACK"] = "1"
    os.environ["COMFY_NO_TELEMETRY"] = "1"
    os.environ["COMFY_CACHE_DIR"] = str(directory / "cache")

    # Upstream initializes ConfigManager during this import. Its pinned config
    # table has no environment override, so select the table before importing.
    import comfy_cli.__main__ as comfy_main

    comfy_main.main()


if __name__ == "__main__":
    _main()
