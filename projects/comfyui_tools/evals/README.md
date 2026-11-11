---
title: Image-generation evaluation cases
---

Each case owns its prompt, supplied reference and review criteria. The operator
approved publishing both cases and their assets. Both reference files retain their original operator-supplied bytes.

The thorn case uses the final corpse-and-thorns prompt from the supplied
conversation. Its reference guides architectural form; the text requests a
painted digital-art style rather than copying the reference's drawing medium.
Compare text-only and reference-guided results with the same seeds.

The plush case requests visible petting motion and a seamless loop. Frame
prompts describe successive static hand poses, with identical first and last
poses. Its edit_region selects the hand area for masked comparisons. The matrix
compares chained feedback, original-reference sampling and native Qwen editing;
it does not imply a trained video model or guarantee smooth motion.

Execution, media validity, timings and cleanup can be checked automatically.
Review requested subject, style, reference adherence and artistic quality
separately on a 0–4 scale: absent, weak, partial, mostly achieved, consistently
achieved. Leave unreviewed scores null. Valid files are not quality acceptance.

Deployment configuration and credentials belong to infra/image_gen/evals.
Generated outputs and reports belong to the configured ignored output directory.
