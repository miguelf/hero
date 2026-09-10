---
title: "Audit the 23 reference skills unreachable from any command"
slug: hsu-unreachable-skills-audit
type: chore
status: planning
domain: engineering
parent: harness-skills-unification
size: small
priority: low
created: 2026-09-10
tags: [skills, scrub, knowledge]
conflicts-with: [hsu-rename-colliding-skills, hsu-user-invocable-flag]
---

# Audit the command-unreachable reference skills

## Goal

Every canonical skill has a named consumer (command, agent, skill, or an
explicit model-auto-load trigger) or is removed.

## Kickoff

Deliver `hsu-unreachable-skills-audit` (parent `harness-skills-unification`).
For each of the 23 skills below, determine whether it is loaded anywhere other
than the AGENTS.md skills table; then wire it to a concrete consumer, keep it
with a stated auto-load trigger in its `description`, or delete it. Record the
per-skill verdict in the delivery audit and update AGENTS.md tables.

## Problem

The 2026-09-10 fan-in mapping found these skills reachable from no command
(directly or via the agent the command delegates to):

`api-design-and-contracts`, `attention-lifecycle-awareness`, `database-stack`,
`deep-code-enrichment`, `deferred-work-suggestions`, `executive-report`,
`explainer-format`, `go-stack`, `greenfield-scaffolding`, `groovy-stack`,
`integration-boundaries`, `issue-list-report`, `java-stack`,
`javascript-stack`, `knowledge-flywheel`, `migration-safety`,
`nudge-awareness`, `performance-optimization`, `python-stack`, `react-stack`,
`root-cause-classification`, `rust-stack`, `test-strategy`.

The 8 `*-stack` skills are loaded by `stack-detection` and likely stay; the
rest need a verdict.

## Changes

- `core/skills/*`, `domains/engineering/skills/*` — deletions or `description`
  trigger phrasing.
- `domains/*/AGENTS.md`, `internal/install/agents_md.go` — skills tables.
- Possibly `domains/engineering/agents/*.md` — new explicit loads.

## Acceptance Criteria

- AC-1 THE SYSTEM SHALL have a recorded verdict (wire / keep-with-trigger /
  remove) for each of the 23 listed skills.
- AC-2 WHEN a skill is kept THE SYSTEM SHALL name at least one consumer or an
  explicit auto-load trigger in its `description`.
- AC-3 WHEN a skill is removed THE SYSTEM SHALL remove it from every AGENTS.md
  skills table and add a prune rule for installed copies.
- AC-4 THE SYSTEM SHALL pass the harness smoke matrix for all targets.
