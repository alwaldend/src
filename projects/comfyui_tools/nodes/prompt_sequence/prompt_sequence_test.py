"""Verify the first-party ComfyUI node contract before GPU workflow acceptance."""

import collections.abc
import json
import os
import pathlib
import types
import unittest
from unittest import mock

from projects.comfyui_tools.nodes.prompt_sequence import prompt_sequence


class _PromptSequenceTest(unittest.TestCase):
    """Exercise paired prompt/prefix lists and rejection before sampling."""

    def test_unique_run_prefix(self) -> None:
        """Give repeated and cached-input runs distinct grouped panel prefixes."""
        node: prompt_sequence.PromptSequence = prompt_sequence.PromptSequence()
        first: tuple[list[str], list[str], str, object] = node.execute(
            "First\n\nSecond", "Context", "line", "comic", "Blur"
        )
        second: tuple[list[str], list[str], str, object] = node.execute(
            "First\n\nSecond", "Context", "line", "comic", "Blur"
        )
        self.assertEqual(first[0], ["Context\n\nFirst", "Context\n\nSecond"])
        self.assertEqual(first[2], "Blur")
        self.assertIsNone(first[3])
        self.assertRegex(first[1][0], r"^comic_[0-9a-f]{32}_001$")
        self.assertEqual(first[1][0][:-3], first[1][1][:-3])
        self.assertNotEqual(first[1][0], second[1][0])
        self.assertTrue(
            prompt_sequence.PromptSequence.IS_CHANGED()
            != prompt_sequence.PromptSequence.IS_CHANGED()
        )

    def test_native_panel_progress(self) -> None:
        """Count sampler invocations on a clone while preserving sampler results."""
        original: mock.Mock = mock.Mock()
        cloned: mock.Mock = original.clone.return_value
        extension: types.SimpleNamespace = types.SimpleNamespace(
            WrappersMP=types.SimpleNamespace(OUTER_SAMPLE="outer_sample")
        )
        server: mock.Mock = mock.Mock()
        with mock.patch.dict(
            "sys.modules",
            {
                "comfy": types.SimpleNamespace(patcher_extension=extension),
                "comfy.patcher_extension": extension,
                "server": server,
            },
        ):
            result: tuple[list[str], list[str], str, object] = (
                prompt_sequence.PromptSequence().execute(
                    "First\nSecond", "", "line", "comic", model=original
                )
            )
            wrapper: collections.abc.Callable[..., object] = (
                cloned.add_wrapper_with_key.call_args.args[2]
            )
            executor: mock.Mock = mock.Mock(return_value="sampled")
            self.assertEqual(wrapper(executor, "latent"), "sampled")
            self.assertEqual(wrapper(executor, "latent"), "sampled")
        self.assertIs(result[3], cloned)
        original.add_wrapper_with_key.assert_not_called()
        messages: list[str] = [
            call.args[0]
            for call in server.PromptServer.instance.send_progress_text.call_args_list
        ]
        self.assertEqual(messages, ["Panel 1 of 2", "Panel 2 of 2"])
        self.assertEqual(executor.call_count, 2)

    def test_order_and_shared_instructions(self) -> None:
        """Keep nonempty panels in input order with matching numbered prefixes."""
        node: prompt_sequence.PromptSequence = prompt_sequence.PromptSequence()
        prompts: list[str]
        prefixes: list[str]
        negative: str
        prompts, prefixes, negative = node.expand(
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
        self.assertEqual(negative, "")
        output: str = os.environ.get("TEST_UNDECLARED_OUTPUTS_DIR", "")
        if output:
            pathlib.Path(output, "prompt-sequence-result.json").write_text(
                json.dumps(
                    {
                        "prompts": prompts,
                        "prefixes": prefixes,
                        "negative": negative,
                    }
                )
            )

    def test_shared_negative_prompt(self) -> None:
        """Share one multiline negative without panel splitting or positive context."""
        node: prompt_sequence.PromptSequence = prompt_sequence.PromptSequence()
        prompts: list[str]
        prefixes: list[str]
        negative: str
        shared_negative: str = "blurry\n\nwatermark, text"
        prompts, prefixes, negative = node.expand(
            "A robot rides a bicycle.\n \nA robot reads a book.",
            "A watercolor comic panel.",
            "line",
            "comic",
            negative_prompt=shared_negative,
        )
        self.assertEqual(
            prompts,
            [
                "A watercolor comic panel.\n\nA robot rides a bicycle.",
                "A watercolor comic panel.\n\nA robot reads a book.",
            ],
        )
        self.assertEqual(prefixes, ["comic_001", "comic_002"])
        self.assertEqual(negative, shared_negative)
        output: str = os.environ.get("TEST_UNDECLARED_OUTPUTS_DIR", "")
        if output:
            pathlib.Path(output, "shared-negative-result.json").write_text(
                json.dumps(
                    {
                        "prompts": prompts,
                        "prefixes": prefixes,
                        "negative": negative,
                    }
                )
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
            (panels, ["comic_001", "comic_002"], ""),
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
