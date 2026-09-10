---
title: "Skill directories carry supporting files through every installer"
slug: hsu-skill-dir-supporting-files
type: feature
status: planning
domain: engineering
parent: harness-skills-unification
size: medium
priority: high
created: 2026-09-10
tags: [install, skills, progressive-disclosure]
conflicts-with: [hsu-copilot-native-surfaces]
---

# Skill directories carry supporting files

## Goal

A canonical skill can be `<name>/SKILL.md` plus linked supporting files, and
every target receives the complete skill — copied as a tree where the harness
supports directories, inlined where it flattens to one file.

## Kickoff

Deliver `hsu-skill-dir-supporting-files` (parent `harness-skills-unification`).
In `internal/install/content.go`, make `installSkillsNested` copy the entire
`<name>/` tree (not only `SKILL.md`), make `installSkillsFlat` inline any file
that `SKILL.md` links via relative markdown link, and extend prune/uninstall to
remove supporting files that disappear from canonical content. Add a fixture
skill with a supporting file and cover it in the smoke matrix.

## Problem

`installSkillsNested` writes only `SKILL.md`; `installSkillsFlat` (Cursor)
produces one file. The Agent Skills progressive-disclosure pattern — thin
`SKILL.md` linking to doctrine — is impossible, which blocks folding
single-consumer doctrine into command skills.

## Changes

- `internal/install/content.go` — `installSkillsNested` tree copy;
  `installSkillsFlat` link-inlining.
- `internal/install/prune.go`, `uninstall.go` — supporting-file removal.
- `internal/install/*_test.go` — fixture skill with `SKILL.md` + `doctrine.md`.

## Acceptance Criteria

- AC-1 WHEN a canonical skill directory contains files beyond `SKILL.md` THE
  SYSTEM SHALL copy every file to `<dest>/<name>/` on nested-layout targets.
- AC-2 WHEN installing on a flat-layout target THE SYSTEM SHALL inline the
  content of each relatively-linked supporting file into `<dest>/<name>.md`.
- AC-3 IF a supporting file is removed from canonical content THEN THE SYSTEM
  SHALL prune it from the destination on re-install.
- AC-4 IF a destination skill directory contains a user-added file not in the
  manifest THEN THE SYSTEM SHALL leave it untouched.
- AC-5 THE SYSTEM SHALL pass the harness smoke matrix for all targets with the
  fixture skill present.
