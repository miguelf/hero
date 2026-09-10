---
title: "Mark reference skills user-invocable: false; decide disable-model-invocation for workflows"
slug: hsu-user-invocable-flag
type: feature
status: planning
domain: engineering
parent: harness-skills-unification
size: small
priority: high
created: 2026-09-10
tags: [skills, frontmatter, copilot, claude, ux]
conflicts-with: [hsu-unreachable-skills-audit]
---

# `user-invocable` flag on reference skills

## Goal

On any harness where skills are `/`-invocable, the menu shows Hero's workflows
and nothing else.

## Kickoff

Deliver `hsu-user-invocable-flag` (parent `harness-skills-unification`). Add
`user-invocable: false` to the frontmatter of every canonical skill under
`core/skills` and `domains/*/skills` (source-level), add a test asserting the
invariant, and record a `/decide` on whether workflow skills rendered from
commands get `disable-model-invocation: true`. Verify Claude and Copilot honor
the flag; verify Codex/Grok/Cursor/OpenCode ignore unknown fields without error.

## Problem

All 57 reference skills are loaded by commands, agents, `stack-detection`, or
the AGENTS.md skills table — none is a standalone user action. Without the
flag they appear alongside the 29 workflows in the `/` menu.

## Changes

- `core/skills/*/SKILL.md`, `domains/*/skills/*/SKILL.md` — frontmatter.
- `internal/install/render.go` — `commandAsSkillRenderer` emits
  `user-invocable: true` (and `disable-model-invocation` per decision).
- Test: every canonical skill declares `user-invocable`; every rendered
  command-skill declares `user-invocable: true`.
- `.hero/knowledge/decisions/` — decision record.

## Acceptance Criteria

- AC-1 THE SYSTEM SHALL declare `user-invocable: false` in every canonical
  reference skill's frontmatter.
- AC-2 WHEN a command is rendered as a skill THE SYSTEM SHALL emit
  `user-invocable: true`.
- AC-3 WHEN installing for codex, grok, cursor, opencode, or generic THE SYSTEM
  SHALL install skills carrying the new field without error (smoke matrix).
- AC-4 THE SYSTEM SHALL record an accepted decision on
  `disable-model-invocation` for workflow skills, with rationale.
