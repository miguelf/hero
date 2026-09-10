---
title: "Rename colliding reference skills (drive, roadmap-review) and drop the command- skill prefix"
slug: hsu-rename-colliding-skills
type: chore
status: planning
domain: engineering
parent: harness-skills-unification
size: small
priority: critical
created: 2026-09-10
tags: [install, skills, codex, grok, rename]
conflicts-with: [hsu-unreachable-skills-audit]
---

# Rename colliding reference skills and drop the `command-` prefix

## Goal

No canonical skill shares a name with a canonical command, so every target can
render commands as bare-named skills without a prefix.

## Kickoff

Deliver `hsu-rename-colliding-skills` (parent `harness-skills-unification`).
Rename `domains/engineering/skills/drive` → `drive-protocol` and
`roadmap-review` → `roadmap-review-doctrine` (dir + `name:`), update all
references in commands/agents/AGENTS.md/tests, then remove `commandSkillPrefix`
from `internal/install/render.go` so Codex/Grok emit `<name>/SKILL.md`, with a
prune rule for legacy `command-*` dirs. Run the harness smoke matrix.

## Problem

`drive` and `roadmap-review` are both a command and a skill. Codex/Grok work
around it with `command-<name>`; Copilot/Claude cannot, because skill `name`
is the slash command.

## Changes

- `domains/engineering/skills/{drive,roadmap-review}/` → renamed dirs and
  `name:` frontmatter.
- `domains/engineering/commands/{drive,roadmap-review}.md`,
  `domains/engineering/agents/roadmap-reviewer.md`, `domains/*/AGENTS.md`
  skills tables, `internal/install/agents_md.go` — reference updates.
- `internal/install/render.go` — remove `commandSkillPrefix`/`commandSkillDir`.
- `internal/install/target_{codex,grok}.go`, `prune.go`/`uninstall.go` — prune
  `command-*` skill dirs; migration test.

## Acceptance Criteria

- AC-1 THE SYSTEM SHALL have zero name intersections between canonical
  `commands/*.md` and `skills/*/SKILL.md` (asserted by a test).
- AC-2 WHEN installing for codex or grok THE SYSTEM SHALL write commands to
  `<skills>/<name>/SKILL.md` with no `command-` prefix.
- AC-3 WHEN a prior install left `<skills>/command-*/` directories THE SYSTEM
  SHALL remove them on re-install and report them as pruned.
- AC-4 THE SYSTEM SHALL contain no remaining reference to the old skill names
  `drive`/`roadmap-review` as skills in commands, agents, or AGENTS.md.
- AC-5 THE SYSTEM SHALL pass `harness_smoke_test.go` for every target.
