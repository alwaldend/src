## Context

The repository pins a Danbooru CSV with tag names, post counts and deprecation
flags. The pinned ComfyUI runtime provides Save Text, whose saved files are
shown in App mode and can be copied or downloaded.

## Goals / Non-Goals

Search the pinned vocabulary locally in a standalone graph or App. The search
runs on demand and needs no model, GPU, network call or browser extension.

## Decisions

Stream the CSV on each search and retain only the best requested matches.
Rank exact matches before prefix matches and prefix matches before other
substring matches. Break ties by descending post count and then tag name.
Treat spaces as underscores and compare without case. Hide deprecated tags by
default and allow the App user to include them.

Package the node and vocabulary together through a manual install target.
Generate the graph and App in the existing catalog. Send the search node's
comma-separated tag names to native Save Text for display and download.

## Risks / Trade-offs

Each new query scans the full vocabulary. Check that full-data latency is usable
and that packaged searches require no network. Exercise incorrect ordering or
normalization, result limits, deprecated tags, empty queries, unmatched queries,
missing data and the native text-output wiring.
