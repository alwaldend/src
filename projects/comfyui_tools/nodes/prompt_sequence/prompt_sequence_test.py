"""Verify the first-party ComfyUI node contract before GPU workflow acceptance."""

import json
import os
import pathlib
import unittest

from projects.comfyui_tools.nodes.prompt_sequence import prompt_sequence


class _PromptSequenceTest(unittest.TestCase):
    """Exercise paired prompt/prefix lists and rejection before sampling."""

    def test_order_and_shared_instructions(self) -> None:
        """Keep nonempty panels in input order with matching numbered prefixes."""
        node: prompt_sequence.PromptSequence = prompt_sequence.PromptSequence()
        prompts: list[str]
        prefixes: list[str]
        prompts, prefixes = node.expand(
            "A robot rides a bicycle.\n\nA robot reads a book.",
            "Use the shared reference character and style.",
            "line",
            "comic",
        )
        self.assertEqual(len(prompts), 2)
        self.assertTrue(prompts[0].endswith("rides a bicycle."))
        self.assertTrue(prompts[1].endswith("reads a book."))
        self.assertTrue(all(p.startswith("Use the shared") for p in prompts))
        self.assertEqual(prefixes, ["comic_001", "comic_002"])
        output: str = os.environ.get("TEST_UNDECLARED_OUTPUTS_DIR", "")
        if output:
            pathlib.Path(output, "prompt-sequence-result.json").write_text(
                json.dumps({"prompts": prompts, "prefixes": prefixes})
            )

    def test_paragraphs(self) -> None:
        """Preserve multiline descriptions when blank lines separate panels."""
        node: prompt_sequence.PromptSequence = prompt_sequence.PromptSequence()
        prompts: list[str] = node.expand(
            "A robot rides a bicycle.\nSnow falls.\n\nA robot reads a book.",
            "",
            "blank line",
            "comic",
        )[0]
        self.assertEqual(len(prompts), 2)
        self.assertIn("\nSnow falls.", prompts[0])

    def test_rejected_input(self) -> None:
        """Reject empty input and unsafe output prefixes without echoing prompts."""
        node: prompt_sequence.PromptSequence = prompt_sequence.PromptSequence()
        for prompts, prefix in [
            (" \n ", "comic"),
            ("A robot.", "../comic"),
            ("A robot.", ""),
        ]:
            with self.subTest(prefix=prefix):
                with self.assertRaises(ValueError):
                    node.expand(prompts, "", "line", prefix)

    def test_json_panels(self) -> None:
        """Preserve paragraphs, literal markup and panel order in editor data."""
        node: prompt_sequence.PromptSequence = prompt_sequence.PromptSequence()
        panels: list[str] = [
            "A robot.\n\nSnow falls.",
            "<robot> reads a book.",
        ]
        self.assertEqual(
            node.expand(json.dumps([*panels, "  "]), "", "JSON", "comic"),
            (panels, ["comic_001", "comic_002"]),
        )
        for invalid in [
            '{"private": "text"}',
            '["private", 2]',
            "private",
            "[]",
        ]:
            with self.subTest(value=invalid):
                with self.assertRaises(ValueError) as caught:
                    node.expand(invalid, "", "JSON", "comic")
                self.assertNotIn("private", str(caught.exception))


if __name__ == "__main__":
    unittest.main()
