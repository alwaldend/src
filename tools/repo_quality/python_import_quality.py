"""Run the pinned syntax linter against fixtures and tracked Python source."""

import argparse
import ast
import importlib.machinery
import json
import pathlib
import subprocess
import sys
import tempfile
import types
import typing

from python.runfiles import runfiles


def module_exists(name: str, modules: set[str]) -> bool:
    if name in modules:
        return True
    if isinstance(sys.modules.get(name), types.ModuleType):
        return True
    search: list[str] = list(sys.path)
    prefix: str = ""
    component: str
    for component in name.split("."):
        prefix = f"{prefix}.{component}" if prefix else component
        spec: importlib.machinery.ModuleSpec | None = (
            importlib.machinery.PathFinder.find_spec(prefix, search)
        )
        if spec is None:
            return False
        search = list(spec.submodule_search_locations or [])
    return True


def lint_source(
    source: str, modules: set[str], *, allow_exports: bool = False
) -> str:
    """Exclude verified module imports from the symbol-import syntax rule."""
    tree: ast.Module = ast.parse(source)
    exports: set[str] = set()
    statement: ast.stmt
    if allow_exports:
        for statement in tree.body:
            value: ast.expr | None = None
            if isinstance(statement, ast.Assign) and any(
                isinstance(target, ast.Name) and target.id == "__all__"
                for target in statement.targets
            ):
                value = statement.value
            elif (
                isinstance(statement, ast.AnnAssign)
                and isinstance(statement.target, ast.Name)
                and statement.target.id == "__all__"
            ):
                value = statement.value
            if isinstance(value, (ast.Tuple, ast.List)):
                exports = {
                    item.value
                    for item in value.elts
                    if isinstance(item, ast.Constant)
                    and isinstance(item.value, str)
                }
    lines: list[str] = source.splitlines(keepends=True)
    node: ast.AST
    for node in sorted(
        (item for item in ast.walk(tree) if isinstance(item, ast.ImportFrom)),
        key=lambda item: (item.lineno, item.col_offset),
        reverse=True,
    ):
        if isinstance(node, ast.ImportFrom):
            module_import: bool = bool(
                node.level == 0
                and node.module
                and all(
                    item.name != "*"
                    and module_exists(f"{node.module}.{item.name}", modules)
                    for item in node.names
                )
            )
            public_export: bool = allow_exports and all(
                item.name != "*" and (item.asname or item.name) in exports
                for item in node.names
            )
            if module_import or public_export:
                assert node.end_lineno is not None
                # Preserve other statements on the same line and source positions.
                start: int = node.lineno - 1
                end: int = node.end_lineno - 1
                before: str = lines[start][: node.col_offset]
                after: str = lines[end][node.end_col_offset :]
                lines[start] = before + "pass" + after
                index: int
                for index in range(start + 1, end + 1):
                    lines[index] = "\n"
    return "".join(lines)


def main() -> None:
    parser: argparse.ArgumentParser = argparse.ArgumentParser()
    parser.add_argument("--tool", required=True)
    parser.add_argument("--rule", required=True)
    parser.add_argument("--cases", required=True)
    parser.add_argument("--workspace", required=True)
    args: argparse.Namespace = parser.parse_args()
    resolver: runfiles.Runfiles | None = runfiles.Create()
    if resolver is None:
        raise RuntimeError("Missing test runfiles")

    def resolve(name: str) -> pathlib.Path:
        location: str | None = resolver.Rlocation(name)
        if location is None:
            raise FileNotFoundError(name)
        return pathlib.Path(location).resolve()

    tool: pathlib.Path = resolve(args.tool)
    workspace: pathlib.Path = resolve(args.workspace).parent
    scratch: pathlib.Path = workspace / "out/python-import-quality"
    scratch.mkdir(parents=True, exist_ok=True)
    tracked: subprocess.CompletedProcess[str] = subprocess.run(
        ["git", "ls-files", "-z", "--", "*.py"],
        cwd=workspace,
        capture_output=True,
        text=True,
        check=True,
        timeout=15,
    )
    sources: list[str] = [
        path
        for path in tracked.stdout.split("\0")
        if path and (workspace / path).is_file()
    ]
    modules: set[str] = set()
    path: str
    for path in sources:
        relative: str = path.removesuffix(".py")
        if relative.endswith("/__init__"):
            relative = relative.removesuffix("/__init__")
        modules.add(relative.replace("/", "."))
        # Legacy language source roots expose modules below main/py.
        if "/main/py/" in relative:
            modules.add(relative.split("/main/py/", 1)[1].replace("/", "."))
    rule: pathlib.Path = resolve(args.rule)
    cases: dict[str, list[str]] = typing.cast(
        dict[str, list[str]], json.loads(resolve(args.cases).read_text())
    )
    command: list[str] = [
        str(tool),
        "scan",
        "--rule",
        str(rule),
        "--no-ignore",
        "hidden",
        "--no-ignore",
        "vcs",
        "--color",
        "never",
    ]
    temporary: str
    with tempfile.TemporaryDirectory(dir=scratch) as temporary:
        source: pathlib.Path = pathlib.Path(temporary) / "fixture.py"
        category: str
        examples: list[str]
        for category, examples in cases.items():
            example: str
            for example in examples:
                source.write_text(
                    lint_source(
                        example + "\n",
                        modules,
                        allow_exports=category.endswith("_exports"),
                    )
                )
                result: subprocess.CompletedProcess[str] = subprocess.run(
                    [*command, str(source)],
                    capture_output=True,
                    text=True,
                    timeout=15,
                )
                expected: int = 0 if category.startswith("valid") else 1
                if result.returncode != expected:
                    raise AssertionError(
                        f"{category} fixture failed: {example}: {result.stderr}"
                    )
    if sources:
        with tempfile.TemporaryDirectory(dir=scratch) as temporary:
            directory: pathlib.Path = pathlib.Path(temporary)
            for path in sources:
                destination: pathlib.Path = directory / path
                destination.parent.mkdir(parents=True, exist_ok=True)
                destination.write_text(
                    lint_source(
                        (workspace / path).read_text(),
                        modules,
                        allow_exports=destination.name == "__init__.py",
                    )
                )
            subprocess.run([*command, str(directory)], check=True, timeout=60)


if __name__ == "__main__":
    main()
