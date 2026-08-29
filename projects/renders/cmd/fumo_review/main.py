"""Command-line entry point for Fumo render-packet audits."""

import argparse
import os
import pathlib
import sys
import typing

from projects.renders.cmd.fumo_review.review import (
    ConfigError,
    audit_config,
    write_outputs,
)


def _workspace_path(path: pathlib.Path) -> pathlib.Path:
    if path.is_absolute():
        return path
    workspace = os.environ.get("BUILD_WORKSPACE_DIRECTORY")
    return pathlib.Path(workspace, path) if workspace else path


def main(args: typing.Optional[typing.Sequence[str]] = None) -> int:
    """Run a review audit and return a shell exit code."""
    parser = argparse.ArgumentParser(
        description="Audit fixed-view Fumo render packets."
    )
    parser.add_argument("--config", required=True, type=pathlib.Path)
    parser.add_argument("--output-dir", required=True, type=pathlib.Path)
    options = parser.parse_args(args)
    config_path = _workspace_path(options.config)
    output_dir = _workspace_path(options.output_dir)
    try:
        result = audit_config(config_path)
        results_path, html_path = write_outputs(result, output_dir)
    except (ConfigError, OSError) as error:
        parser.error(str(error))
    verdict = "PASS" if result["passed"] else "FAIL"
    print(f"{verdict}: {html_path}")
    print(f"results: {results_path}")
    return 0 if result["passed"] else 1


if __name__ == "__main__":
    sys.exit(main())
