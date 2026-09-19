---
title: Diagrams in the AC (After Clanker) era
linkTitle: Diagrams in the AC (After Clanker) era
date: 2026-09-18
description: Replacing hand-maintained Drawio diagrams with hermetically-built Mermaid diagrams.
tags:
  - mermaid
  - drawio
  - repo
draft: true
---

I have some diagrams in the repo, as one naturally does, because they are quite useful. They do not look particularly pretty, but they display information, which is the whole point. Here is an example of such a diagram:
![old_t3code_arch_diagram](./old_t3code_arch_diagram.svg).

In the BC (Before Clanker) era, I've manually made them using Drawio's desktop app. I've used this approach because it requires no setup and no effort - you just fire up the app and then click things in the UI. Everything is local, and the diagrams themselves are XML, so you can version control them. In the AC (After Clanker) era, I have more options, so I've decided to improve my diagrams.

## New tool

The first goal was, of course, to make them prettier - I wanted to see something pleasing to the eyes, not ugly boxes and straight arrows. Ideally, something like [Excalidraw](https://excalidraw.com/) - hand-drawn look, rounded corners, curved edges, all the rage.

The second goal was to simplify their source representation. Drawio diagrams are XML, so they can be version controlled and clanker-managed, but they are too complex and not readable. I wanted a simple DSL that allowed me to just describe the diagram itself - nodes, labels, edges, etc.

The third goal was theme support. While every diagram was supposed to be as simple as possible, they all should look pretty and use one style without being a part of the same project.

The forth goal was hermeticity. All these diagrams should be rendered locally and hermetically through Bazel without a CDN or any other funny business.

In the end, I've settled on [Mermaid](https://github.com/mermaid-js/mermaid) because it's the first big Open Source project I've found that fit my criteria and I did not care enough to research other options. After all, we are in the AC era - I can just rewrite things.

## Hermeticity

Nothing particularly interesting here. I've hermetically set up:

- [Chrome for Testing](https://developer.chrome.com/docs/automation-and-testing/chrome-for-testing)
- [Architects Daughter font](https://github.com/google/fonts/blob/main/ofl/architectsdaughter/METADATA.pb)
- [Mermaid](https://github.com/mermaid-js/mermaid)
- Some Bazel glue to make it all work together

And then I've had

## Final look

## Lesson

You can just do things, I guess.

## Links

- Excalidraw: https://excalidraw.com/
- Mermaid: https://github.com/mermaid-js/mermaid
- Actual blog post:
