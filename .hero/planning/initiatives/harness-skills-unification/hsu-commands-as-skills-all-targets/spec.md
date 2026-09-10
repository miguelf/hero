---
title: "Commands as skills on every target; fold single-consumer doctrine into command skills"
slug: hsu-commands-as-skills-all-targets
type: feature
status: planning
domain: engineering
parent: harness-skills-unification
depends-on: [hsu-rename-colliding-skills, hsu-user-invocable-flag, hsu-skill-dir-supporting-files, hsu-copilot-native-surfaces]
size: large
priority: medium
created: 2026-09-10
tags: [install, skills, commands, claude, cursor, opencode, codex, grok]
---

# Commands as skills on every target

## Goal

One artifact type — the skill — carries Hero's workflows on every harness
where skills are user-invocable, and the 14 doctrine skills consumed by exactly
one command live inside that command's skill directory as linked supporting
files.

## Kickoff

Deliver `hsu-commands-as-skills-all-targets` (parent
`harness-skills-unification`). First verify from vendor docs whether Cursor and
OpenCode skills are user-invocable and record the per-target matrix in the
spec. For targets where they are (Claude, Copilot; Codex/Grok already
skill-only), render commands via `commandAsSkillRenderer`, stop writing the
commands dir, and prune it. For targets where they are not, keep the commands
dir. Then move the 14 single-consumer doctrine skills into their command's
skill directory as `doctrine.md` linked from `SKILL.md`, and update references.

## Problem

Claude ships `.claude/commands/` and `.claude/skills/`; Copilot (after
`hsu-copilot-native-surfaces`) ships skills only; Codex/Grok ship skills only;
Cursor/OpenCode ship both. The split is per-target incidental, not designed.
Meanwhile 14 skills exist solely as the doctrine of one command (`drive`,
`roadmap-review`, `scrub`, `review`, `peer`, `handoff`, `mock` ×2, `check`,
`convention`, `docs`, `release` ×2) and are cheaper to read as a supporting
file of that command's skill.

## Changes

- Decision record: per-target matrix (skills invocable? → commands rendered as
  skills? → commands dir retained?).
- `internal/install/target_{claude,cursor,opencode,generic}.go` — per matrix.
- `internal/install/prune.go`, `uninstall.go`, migration tests — commands dir
  removal where retired.
- `domains/engineering/skills/{drive-protocol,roadmap-review-doctrine,code-scrub,
  pr-review,cross-repo-peering,html-mockup-generation,swiftui-mockup-renderer,
  stack-detection,dependency-analysis,convention-writing,documentation-practices,
  devops-and-operations,release-and-deployment}` and `core/skills/next-md` →
  `commands/<cmd>/doctrine.md` (or equivalent) under the command's skill dir;
  `SKILL.md` links it. Note `stack-detection` is also loaded by
  `/scrub`-adjacent agents — confirm single consumer before moving.
- Reference updates in agents and AGENTS.md skills tables.

## Acceptance Criteria

- AC-1 THE SYSTEM SHALL record an accepted per-target decision covering
  claude, copilot, codex, grok, cursor, opencode, generic.
- AC-2 WHERE a target's skills are user-invocable THE SYSTEM SHALL render every
  command as `<skills>/<name>/SKILL.md` with `user-invocable: true` and write
  no commands directory.
- AC-3 WHERE a target's skills are not user-invocable THE SYSTEM SHALL continue
  to write the commands directory unchanged.
- AC-4 IF a target retires its commands directory THEN THE SYSTEM SHALL prune
  the previously installed directory on re-install.
- AC-5 WHEN a command's skill directory contains linked doctrine THE SYSTEM
  SHALL install it on every target per `hsu-skill-dir-supporting-files`.
- AC-6 THE SYSTEM SHALL have no canonical skill whose only consumer is a single
  command still living as a standalone skill directory.
- AC-7 THE SYSTEM SHALL pass the harness smoke matrix for all targets.
