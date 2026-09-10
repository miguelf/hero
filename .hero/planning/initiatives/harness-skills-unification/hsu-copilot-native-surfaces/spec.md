---
title: "Copilot CLI/App native surfaces — commands as .github/skills, agents as .github/agents"
slug: hsu-copilot-native-surfaces
type: feature
status: planning
domain: engineering
parent: harness-skills-unification
depends-on: [hsu-rename-colliding-skills]
size: large
priority: critical
created: 2026-09-10
tags: [install, copilot, skills, agents, prompts]
conflicts-with: [hsu-skill-dir-supporting-files]
---

# Copilot CLI/App native surfaces

## Goal

`hero install --target copilot` produces Hero workflows and agents that the
GitHub Copilot CLI, App, cloud agent, and VS Code all load natively.

## Kickoff

Deliver `hsu-copilot-native-surfaces` (parent `harness-skills-unification`).
In `internal/install/target_copilot.go`, render commands via
`commandAsSkillRenderer("GitHub Copilot")` to `.github/skills/<name>/SKILL.md`
and agents to `.github/agents/<name>.agent.md`; write a thin
`.github/copilot-instructions.md` pointing at `AGENTS.md`; generalize the
Codex-only "commands are skills" AGENTS.md section in `agents_md.go`; prune
`.github/prompts/` as dead bytes; update tests pinning `.github/prompts/`.

## Problem

Commands and agents land only in `.github/prompts/*.prompt.md`, which the
Copilot CLI/App never read. Today a Copilot user can call zero Hero workflows.

## Changes

- `internal/install/target_copilot.go` — command and agent destinations;
  `.github/prompts/` moved to the legacy-cleanup list.
- `internal/install/render.go` — agent → `.agent.md` renderer (frontmatter:
  `name`, `description`, optional `tools`); retire `renderCopilotPromptFile`
  or keep behind an explicit VS Code-prompts opt-in (decide at delivery;
  default is retire).
- `internal/install/agents_md.go` — "Running Hero Workflows as skills" section
  parameterized by harness; Copilot routing table.
- New `.github/copilot-instructions.md` managed block.
- Tests: `harness_smoke_test.go`, `overlay_install_test.go`,
  `qa_overlay_install_test.go`, `attention_guidance_test.go`,
  `migration_test.go`, `prune_test.go`, `uninstall_test.go`.

## Acceptance Criteria

- AC-1 WHEN installing for copilot THE SYSTEM SHALL write every canonical
  command to `.github/skills/<name>/SKILL.md` with `name:` equal to the
  directory and the Hero execution preamble.
- AC-2 WHEN installing for copilot THE SYSTEM SHALL write every canonical agent
  to `.github/agents/<name>.agent.md` with `name` and `description`
  frontmatter.
- AC-3 WHEN installing for copilot THE SYSTEM SHALL write a managed
  `.github/copilot-instructions.md` that references `AGENTS.md` and does not
  duplicate the routing table.
- AC-4 IF a prior install left `.github/prompts/{commands,agents}/` THEN THE
  SYSTEM SHALL remove those files on re-install and report them as pruned.
- AC-5 WHEN `hero uninstall --target copilot` runs THE SYSTEM SHALL remove
  `.github/skills/<command>/`, `.github/agents/*.agent.md`, and the managed
  block in `.github/copilot-instructions.md`.
- AC-6 THE SYSTEM SHALL pass the harness smoke matrix for all targets.
- AC-7 WHEN a Copilot CLI session opens a repo installed by this build THE
  SYSTEM SHALL list `/design`, `/deliver`, `/diagnose` in the `/` menu (manual
  verification recorded in the delivery audit).
