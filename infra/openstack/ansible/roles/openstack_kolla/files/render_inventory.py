#!/usr/bin/env python3
"""Replace host-bearing sections in Kolla's release-matched inventory."""

from __future__ import annotations

import argparse
import json
import re
from pathlib import Path
from typing import Dict, Iterable, List

SECTION_RE = re.compile(r"^\[([^]]+)]\s*$", re.MULTILINE)


def render(source: str, replacements: Dict[str, List[str]]) -> str:
    source_sections = {match.group(1) for match in SECTION_RE.finditer(source)}
    missing_sections = set(replacements).difference(source_sections)
    if missing_sections:
        raise ValueError(
            "source inventory is missing sections: "
            + ", ".join(sorted(missing_sections))
        )

    output: List[str] = []
    skip_section_body = False

    for line in source.splitlines():
        match = SECTION_RE.match(line)
        if match:
            section = match.group(1)
            output.append(line)
            if section in replacements:
                output.extend(replacements[section])
                output.append("")
                skip_section_body = True
            else:
                skip_section_body = False
            continue

        if not skip_section_body:
            output.append(line)

    return "\n".join(output).rstrip() + "\n"


def validate(replacements: Dict[str, List[str]]) -> None:
    required = {"control", "network", "compute", "monitoring", "storage", "deployment"}
    missing = required.difference(replacements)
    if missing:
        raise ValueError(f"missing replacement sections: {', '.join(sorted(missing))}")

    for section, lines in replacements.items():
        if not isinstance(lines, list) or not all(isinstance(line, str) for line in lines):
            raise TypeError(f"section {section!r} must contain a list of strings")
        if any("\n" in line or "\r" in line for line in lines):
            raise ValueError(f"section {section!r} contains a multiline host entry")


def parse_args(argv: Iterable[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--source", type=Path, required=True)
    parser.add_argument("--sections", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    return parser.parse_args(argv)


def main() -> None:
    args = parse_args()
    replacements = json.loads(args.sections.read_text(encoding="utf-8"))
    validate(replacements)
    rendered = render(args.source.read_text(encoding="utf-8"), replacements)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    if args.output.exists() and args.output.read_text(encoding="utf-8") == rendered:
        print("unchanged")
        return
    args.output.write_text(rendered, encoding="utf-8")
    print("changed")


if __name__ == "__main__":
    main()
