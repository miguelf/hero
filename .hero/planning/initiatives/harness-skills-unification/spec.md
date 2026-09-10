---
title: "Harness Skills Unification — Hero workflows as native, user-invocable skills on every harness"
slug: harness-skills-unification
type: initiative
status: planning
domain: engineering
size: x-large
priority: high
horizon: now
created: 2026-09-10
tags: [install, harness, copilot, claude, codex, skills, commands, dx]
child:
  - hsu-rename-colliding-skills
  - hsu-copilot-native-surfaces
  - hsu-user-invocable-flag
  - hsu-skill-dir-supporting-files
  - hsu-unreachable-skills-audit
  - hsu-commands-as-skills-all-targets
---

# Harness Skills Unification

## Goal

A user who runs `hero install --target <any>` gets Hero's 29 workflows as
native `/`-callable entries in that harness — starting with the GitHub Copilot
CLI and App, where today they get **zero** — while the reference skills that
back those workflows stay out of the `/` menu. Along the way, collapse the
"commands vs skills" split that every target currently pays for in a different
way (VS Code-only prompts, `command-` prefixes, duplicate `/drive` entries).

## Kickoff

Continue the `harness-skills-unification` initiative. Read
`.hero/planning/initiatives/harness-skills-unification/spec.md`, check
`hero_blocked` and `hero_queue` for the next dependency-ready child, and
`/deliver` it. Wave 1 is `hsu-rename-colliding-skills`; Wave 2 children
(`hsu-copilot-native-surfaces`, `hsu-user-invocable-flag`,
`hsu-skill-dir-supporting-files`, `hsu-unreachable-skills-audit`) are
independent; Wave 3 (`hsu-commands-as-skills-all-targets`) needs all of Wave 2
except the audit. Every child is gated by the
`harness-changes-cover-all-targets` tripwire.

## Problem

Validated against a live install (morpheus-dev-utils, hero v0.34.2):

1. **Copilot CLI/App cannot call any Hero workflow.** `target_copilot.go`
   renders commands and agents only to `.github/prompts/{commands,agents}/*.prompt.md`
   — a VS Code-only surface. The Copilot CLI, App, and cloud agent read
   `.github/skills/<name>/SKILL.md` and `.github/agents/<name>.agent.md`; they
   never look in `.github/prompts/`. Skills already land correctly; commands
   and agents do not.

2. **Two canonical skills collide with commands by name.** `drive` and
   `roadmap-review` exist as both `commands/<name>.md` and
   `skills/<name>/SKILL.md` (deliberately — the skill is the command's
   doctrine). Codex/Grok dodge this with a `command-` prefix
   (`internal/install/render.go: commandSkillPrefix`), which is tolerable
   there because skills aren't slash-invocable, but on Copilot/VS Code/Claude
   skill `name` *is* the slash command, so `/command-design` would leak to
   users.

3. **Every target ships 29 + 57 = 86 skill-shaped items with no
   invocability distinction.** All 57 reference skills are loaded by a
   command, an agent, `stack-detection`, or the `AGENTS.md` skills table —
   none is a standalone user action. Without `user-invocable: false` they
   flood the `/` menu on any harness where skills are invocable.

4. **The installer cannot carry supporting files.** `installSkillsNested`
   copies only `<name>/SKILL.md`; `installSkillsFlat` (Cursor) flattens to one
   file. Zero canonical skills have a supporting file today, so the Agent
   Skills progressive-disclosure pattern (thin `SKILL.md` + linked doctrine)
   is unavailable.

5. **The command→skill graph is thinner than it looks.** Of 57 skills: 20 are
   shared by ≥2 commands (almost entirely via the two delivery-lead agents,
   i.e. agent dependencies), 14 are the private doctrine of exactly one
   command, and 23 are reachable from no command at all — only from the
   `AGENTS.md` table or from `stack-detection`.

## Outcome

- Copilot CLI/App/VS Code: `/design`, `/deliver`, … work; agents are
  selectable; `.github/prompts/` is gone.
- One `/drive`, one `/roadmap-review` on every harness.
- `/` menu on skill-invocable harnesses shows 29 workflows, not 86.
- Canonical skill directories may carry supporting files and they travel to
  every target (inlined where the target flattens).
- A per-target decision record for "commands as skills": adopted where skills
  are user-invocable (Claude, Copilot; Codex/Grok already are), commands dir
  retained where they aren't.
- The 23 command-unreachable skills are each either wired to a consumer or
  removed.

## Architecture reference

- Canonical content: `core/{commands,skills,agents}`,
  `domains/engineering/{commands,skills,agents}` — embedded via root
  `content.go`. 29 commands, 57 skills, 35 agents for an engineering install.
- Installers: `internal/install/target_{claude,codex,copilot,cursor,generic,grok,opencode}.go`;
  shared helpers `internal/install/content.go` (`installFlat`,
  `installSkillsNested`, `installSkillsFlat`), `internal/install/render.go`
  (`commandAsSkillRenderer`, `commandSkillPrefix`, `renderCopilotPromptFile`),
  `internal/install/agents_md.go` (routing table + per-harness sections).
- Harness surfaces (verified from vendor docs this design):
  - Copilot CLI/App/VS Code/cloud: `AGENTS.md`, `.github/copilot-instructions.md`,
    `.github/skills/<name>/SKILL.md`, `.github/agents/<name>.agent.md`. Skill
    frontmatter: `name` (must equal dir), `description`, `argument-hint`,
    `user-invocable` (default true), `disable-model-invocation` (default false).
  - Claude Code: `.claude/{agents,commands,skills}`, `CLAUDE.md`; skills are
    `/`-invocable and honor the same two flags.
  - Codex/Grok: skills only (`command-<name>` today); no user slash surface.
  - Cursor / OpenCode: commands dir + skills dir; skill invocability **unverified**
    (resolved in `hsu-commands-as-skills-all-targets`).

## Specs

Six children in three waves.

### Wave 1 — unblock bare names

#### 1. Rename the two colliding reference skills
**Slug:** `hsu-rename-colliding-skills` | **Type:** chore | **Size:** small |
**Priority:** critical | **Deps:** none | **SHIPS FIRST**

Rename `skills/drive` → `skills/drive-protocol` and `skills/roadmap-review` →
`skills/roadmap-review-doctrine`; update every reference (commands, agents,
`AGENTS.md` skills tables, tests). With no name collision left, drop the
`command-` prefix from the Codex/Grok renderer so all targets use bare skill
names, and add a prune rule for the old `command-*` directories.

### Wave 2 — independent, parallel

#### 2. Copilot CLI/App native surfaces
**Slug:** `hsu-copilot-native-surfaces` | **Type:** feature | **Size:** large |
**Priority:** critical | **Deps:** `hsu-rename-colliding-skills`

Render commands to `.github/skills/<name>/SKILL.md` via
`commandAsSkillRenderer("GitHub Copilot")`, agents to
`.github/agents/<name>.agent.md`, write a thin `.github/copilot-instructions.md`
pointer, generalize the AGENTS.md "commands are skills" section (today
Codex-only), prune `.github/prompts/` as dead bytes, and update the ~8 test
files pinning `.github/prompts/` paths.

#### 3. `user-invocable` flag on reference skills
**Slug:** `hsu-user-invocable-flag` | **Type:** feature | **Size:** small |
**Priority:** high | **Deps:** none

Stamp `user-invocable: false` in the canonical frontmatter of all 57 reference
skills (source-level, so every target inherits it). Record a decision on
`disable-model-invocation: true` for workflow skills (keeps routing explicit
via AGENTS.md vs. lets the model auto-fire `/deliver`).

#### 4. Skill directories carry supporting files
**Slug:** `hsu-skill-dir-supporting-files` | **Type:** feature | **Size:** medium |
**Priority:** high | **Deps:** none

`installSkillsNested` copies the whole `<name>/` tree; `installSkillsFlat`
inlines files that `SKILL.md` links so flat targets lose nothing; prune logic
handles removed supporting files. Prerequisite for folding doctrine into
command skills.

#### 5. Audit the 23 command-unreachable skills
**Slug:** `hsu-unreachable-skills-audit` | **Type:** chore | **Size:** small |
**Priority:** low | **Deps:** none

For each skill reachable only from the `AGENTS.md` table (or `stack-detection`
for the 8 `*-stack` guides), decide: wire to a concrete consumer, keep as
model-auto-loaded reference with a stated trigger, or remove.

### Wave 3 — the merge

#### 6. Commands as skills on every target
**Slug:** `hsu-commands-as-skills-all-targets` | **Type:** feature | **Size:** large |
**Priority:** medium | **Deps:** `hsu-rename-colliding-skills`,
`hsu-user-invocable-flag`, `hsu-skill-dir-supporting-files`,
`hsu-copilot-native-surfaces`

Per-target matrix (verify Cursor/OpenCode skill invocability first). Where
skills are user-invocable, render commands as skills and retire the commands
dir with a prune rule; where they aren't, keep the commands dir. Fold the 14
single-consumer doctrine skills (`drive-protocol`, `roadmap-review-doctrine`,
`code-scrub`, `pr-review`, `cross-repo-peering`, `next-md`,
`html-mockup-generation`, `swiftui-mockup-renderer`, `stack-detection`,
`dependency-analysis`, `convention-writing`, `documentation-practices`,
`devops-and-operations`, `release-and-deployment`) into their command's skill
directory as linked supporting files.

## In-flight overlap watch

- `hsu-copilot-native-surfaces` ⇄ `hsu-skill-dir-supporting-files`: both edit
  `internal/install/content.go` and `render.go`.
- `hsu-rename-colliding-skills` ⇄ `hsu-unreachable-skills-audit`: both edit the
  canonical skill tree and `AGENTS.md` skills tables.
- `hsu-user-invocable-flag` ⇄ `hsu-unreachable-skills-audit`: both touch
  canonical skill frontmatter across the whole tree.

## Risks

- **Tripwire scope.** Every child is harness-facing; `harness-changes-cover-all-targets`
  applies to each. A Copilot-only fix in Wave 2 is acceptable only because
  the other targets already have a working command surface — the child must
  still run the full smoke matrix.
- **Cursor/OpenCode unknowns.** If their skills are not user-invocable, the
  merge removes their slash surface. Wave 3 must verify before removing any
  commands dir.
- **Existing installs.** Renames and dir moves leave dead bytes in user repos;
  every child that moves a path adds a prune rule and a migration test.
- **Agent-owned skill graph.** Most skill fan-in comes through
  `feature-delivery-lead` / `platform-delivery-lead`. Do not relocate shared
  skills into command dirs; only single-consumer doctrine moves.

## Open Questions

- `disable-model-invocation` default for workflow skills (decided in child 3).
- Whether hero-code (`hihcp-*`) consumes `user-invocable` — coordinate with
  `hero-in-hero-code-parity` if it reads canonical skill frontmatter.

## Notes

Design session 2026-09-10 validated: all 57 reference skills are
reference-only; fan-in = 20 shared / 14 single-consumer / 23 unreachable from
commands.
