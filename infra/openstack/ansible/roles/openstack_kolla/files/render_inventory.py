#!/usr/bin/env python3
"""Replace host-bearing sections in Kolla's release-matched inventory."""

from __future__ import annotations

import argparse
import json
import re
from collections.abc import Iterable, Mapping, Sequence
from pathlib import Path

SECTION_RE = re.compile(r"^\[([^]]+)]\s*$")
REQUIRED_SECTIONS = frozenset(
    {
        "compute",
        "control",
        "deployment",
        "monitoring",
        "network",
        "storage",
    }
)


def source_sections(source: str) -> list[str]:
    """Return inventory section names and reject duplicate sections."""
    sections = [
        match.group(1)
        for line in source.splitlines()
        if (match := SECTION_RE.match(line))
    ]
    duplicates = sorted(
        section for section in set(sections) if sections.count(section) > 1
    )
    if duplicates:
        raise ValueError(
            "source inventory contains duplicate sections: "
            + ", ".join(duplicates)
        )
    return sections


def validate(replacements: Mapping[str, Sequence[str]]) -> None:
    """Validate the exact set and shape of replacement sections."""
    missing = REQUIRED_SECTIONS.difference(replacements)
    if missing:
        raise ValueError(
            "missing replacement sections: " + ", ".join(sorted(missing))
        )

    extra = set(replacements).difference(REQUIRED_SECTIONS)
    if extra:
        raise ValueError(
            "unexpected replacement sections: " + ", ".join(sorted(extra))
        )

    for section, lines in replacements.items():
        if not isinstance(lines, list):
            raise TypeError(f"section {section!r} must contain a list")
        if not all(isinstance(line, str) for line in lines):
            raise TypeError(
                f"section {section!r} must contain only strings"
            )
        if any(not line.strip() for line in lines):
            raise ValueError(f"section {section!r} contains an empty host entry")
        if any("\n" in line or "\r" in line for line in lines):
            raise ValueError(
                f"section {section!r} contains a multiline host entry"
            )


def render(source: str, replacements: Mapping[str, list[str]]) -> str:
    """Render a Kolla inventory with replaced host-bearing sections."""
    validate(replacements)
    sections = set(source_sections(source))
    missing = REQUIRED_SECTIONS.difference(sections)
    if missing:
        raise ValueError(
            "source inventory is missing sections: "
            + ", ".join(sorted(missing))
        )

    output: list[str] = []
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


def parse_args(argv: Iterable[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--source", type=Path, required=True)
    parser.add_argument("--sections", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    return parser.parse_args(argv)


def main() -> None:
    args = parse_args()
    replacements = json.loads(args.sections.read_text(encoding="utf-8"))
    rendered = render(
        args.source.read_text(encoding="utf-8"),
        replacements,
    )
    args.output.parent.mkdir(parents=True, exist_ok=True)
    if (
        args.output.exists()
        and args.output.read_text(encoding="utf-8") == rendered
    ):
        print("unchanged")
        return
    args.output.write_text(rendered, encoding="utf-8")
    print("changed")


if __name__ == "__main__":
    main()
