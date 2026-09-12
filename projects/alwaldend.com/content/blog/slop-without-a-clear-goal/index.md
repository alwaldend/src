---
title: Slop without a clear goal
linkTitle: Slop without a clear goal
date: 2026-09-12
description: >-
  Imprecise goals without quality assurance lead to bad results, obviously.
---

A while back, I said to Jipity[^jipity]: "Jipity, how can I improve agent ergonomics in this monorepo?", and Jipity gave me some useful suggestions, but one of them was not like the rest. It said that context in the repo was too splintered — agents had to read many things in different directories to understand how things are done. In my infinite wisdom, I said, "Jipity, write me a comprehensive plan implementing a unified agent system with the goal of improving ergonomics in this monorepo". And it wrote me a plan, I did not really read it.

The plan consisted of several phases which gradually added more agent ergonomics improvements. Over the course of several days I made several PRs implementing the plan. In each PR, agents added things I did not understand according to the plan I did not read. After I prompted them to explain and still did not understand, I gave up and decided to just "trust the plan", and merged everything. After implementing the plan, my repo was littered with garbage manifests, I had a lot of new code, and nothing really changed.

I decided that maybe I did not trust the plan enough, and did a couple more phases of "ergonomics improvements", which also did nothing but increase line count and litter the repo with garbage. Overall, I got around crisp 30k lines of new code.

30k of code is not much in the grand scheme of things, but I don't like that my pristine monorepo had 30k lines of pointless slop, so I've deleted the "agent system".

The lesson here is that you need to have a clear idea of what you are trying to accomplish, I guess.

---

## Data

Jipity made me a table, seems about right.

| Date        | Timeline item                                          | Code files | Non-code files |                    Code LOC |                Non-code LOC |  Total LOC |
| ----------- | ------------------------------------------------------ | ---------: | -------------: | --------------------------: | --------------------------: | ---------: |
| 2026-08-31  | Phase 0 — System contract                              |          4 |             39 |             **21** (+21/−0) |      **6,908** (+6,825/−83) |  **6,929** |
| 2026-08-31  | Phase 1 — Shared semantics and controls                |         32 |             58 |     **2,661** (+2,456/−205) |      **2,779** (+2,736/−43) |  **5,440** |
| 2026-09-01  | Phases 2–3 — Context, recovery, planning, and evidence |        115 |             73 |   **18,018** (+17,850/−168) |     **3,121** (+2,980/−141) | **21,139** |
| 2026-09-02  | Phase 4 — Connect work, delivery, and release          |         26 |             19 |     **2,996** (+2,883/−113) |          **566** (+521/−45) |  **3,562** |
| 2026-09-02  | Phase 5 — Learning contracts                           |         21 |             16 |       **1,524** (+1,523/−1) |          **311** (+258/−53) |  **1,835** |
| 2026-09-03  | Phase 5 follow-up — Exercise the learning workflow     |         25 |             27 |          **914** (+883/−31) |      **1,043** (+1,007/−36) |  **1,957** |
| 2026-09-03  | Phase 6 — Reduce workflow friction                     |         28 |             69 |      **1,198** (+1,101/−97) |       **1,206** (+793/−413) |  **2,404** |
| 2026-09-05  | System simplification                                  |         41 |             49 |     **5,037** (+4,548/−489) |   **4,714** (+2,293/−2,421) |  **9,751** |
| **Overall** | **8 commits**                                          |    **292** |        **350** | **32,369** (+31,265/−1,104) | **20,648** (+17,413/−3,235) | **53,017** |

[^jipity]: GPT 5.6 Sol.
