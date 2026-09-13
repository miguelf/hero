package install

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GitHub Copilot CLI/App reads project agents from .github/agents and
// Agent Skills from .github/skills. Commands are represented as
// user-invocable skills; reference skills are background instructions.
func runCopilot(opts Options) (*Result, error) {
	if opts.Mode != ModeProject || opts.TargetDir == "" {
		return nil, fmt.Errorf("copilot target only supports project mode")
	}
	info, statErr := os.Stat(opts.TargetDir)
	if statErr != nil || !info.IsDir() {
		return nil, fmt.Errorf("target directory does not exist: %s", opts.TargetDir)
	}

	result := &Result{}
	if err := cleanupCopilotLegacy(opts); err != nil {
		return nil, fmt.Errorf("cleanup legacy copilot outputs: %w", err)
	}

	agentsDest := filepath.Join(opts.TargetDir, ".github", "agents")
	if err := renderToFile(opts, result, "agents", agentsDest, renderCopilotAgentFile); err != nil {
		return nil, fmt.Errorf("render agents for copilot: %w", err)
	}

	skillDirs, err := renderCopilotSkills(opts, result)
	if err != nil {
		return nil, fmt.Errorf("render skills for copilot: %w", err)
	}
	skillsDest := filepath.Join(opts.TargetDir, ".github", "skills")
	if err := pruneStaleSkillDirs(opts, staleSkillPrune{dest: skillsDest, written: skillDirs}); err != nil {
		return nil, fmt.Errorf("prune stale copilot skills: %w", err)
	}
	result.skillDirs = skillDirs

	instructionsPath := filepath.Join(opts.TargetDir, ".github", "copilot-instructions.md")
	if err := installInstructionsMd(opts, result, instructionsPath, "copilot"); err != nil {
		return nil, fmt.Errorf("installing copilot-instructions.md: %w", err)
	}
	if err := installNativeInstructionFile(opts, result); err != nil {
		return nil, fmt.Errorf("installing AGENTS.md: %w", err)
	}
	return result, nil
}

// cleanupCopilotLegacy removes only files previously recorded as Hero output.
// Untracked prompt files and files in the old .github/copilot tree are left
// alone because those locations may contain user-authored material.
func cleanupCopilotLegacy(opts Options) error {
	// This tree was never a Copilot loader location, so its contents are
	// unambiguously dead Hero output.
	legacyBase := filepath.Join(opts.TargetDir, ".github", "copilot")
	for _, kind := range []string{"agents", "commands", "skills"} {
		if err := removeLegacyDir(opts, filepath.Join(legacyBase, kind)); err != nil {
			return err
		}
	}
	_ = os.Remove(legacyBase)
	st, err := ReadInstallState(opts.TargetDir)
	if err != nil || st == nil {
		st = nil
	}
	var prior TargetState
	if st != nil {
		prior = st.Targets[string(TargetCopilot)]
	}
	for _, rel := range prior.Files {
		rel = filepath.ToSlash(rel)
		if !strings.HasPrefix(rel, ".github/prompts/") && !strings.HasPrefix(rel, ".github/copilot/") {
			continue
		}
		full := filepath.Join(opts.TargetDir, filepath.FromSlash(rel))
		if _, err := os.Stat(full); err != nil {
			continue
		}
		if opts.DryRun {
			fmt.Fprintf(os.Stderr, "  cleanup %s (would remove legacy Copilot output)\n", full)
			continue
		}
		if err := os.Remove(full); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "  cleanup %s (removed legacy Copilot output)\n", full)
		_ = os.Remove(filepath.Dir(full))
	}
	return cleanupCopilotLegacyPrompts(opts)
}

// cleanupCopilotLegacyPrompts migrates prompt files from installs that
// predate the file manifest. Only names that exist in the current canonical
// content source are eligible; unrelated prompt files remain user-owned.
func cleanupCopilotLegacyPrompts(opts Options) error {
	srcFS := opts.sourceFS()
	if srcFS == nil {
		return nil
	}
	domain := opts.Domain
	if domain == "" {
		domain = "engineering"
	}
	agents, err := selectFlatContent(srcFS, "agents", domain)
	if err != nil {
		return err
	}
	commands, err := selectFlatContent(srcFS, "commands", domain)
	if err != nil {
		return err
	}
	known := make(map[string]bool, len(agents)+len(commands))
	for _, name := range append(agents, commands...) {
		known[name] = true
	}
	for _, kind := range []string{"agents", "commands"} {
		dir := filepath.Join(opts.TargetDir, ".github", "prompts", kind)
		entries, readErr := os.ReadDir(dir)
		if os.IsNotExist(readErr) {
			continue
		}
		if readErr != nil {
			return readErr
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".prompt.md") {
				continue
			}
			sourceName := strings.TrimSuffix(entry.Name(), ".prompt.md") + ".md"
			if !known[sourceName] {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			if opts.DryRun {
				fmt.Fprintf(os.Stderr, "  cleanup %s (would remove legacy Copilot output)\n", path)
				continue
			}
			if err := os.Remove(path); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "  cleanup %s (removed legacy Copilot output)\n", path)
		}
	}
	return nil
}

func renderCopilotSkills(opts Options, result *Result) ([]string, error) {
	srcFS := opts.sourceFS()
	if srcFS == nil {
		return nil, fmt.Errorf("no content source available")
	}
	domain := opts.Domain
	if domain == "" {
		domain = "engineering"
	}
	skills, err := selectSkillContent(srcFS)
	if err != nil {
		return nil, err
	}
	commands, err := selectFlatContent(srcFS, "commands", domain)
	if err != nil {
		return nil, err
	}
	commandNames := map[string]bool{}
	for _, file := range commands {
		commandNames[strings.TrimSuffix(file, ".md")] = true
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	sort.Strings(commands)

	dest := filepath.Join(opts.TargetDir, ".github", "skills")
	var dirs []string
	for _, skill := range skills {
		dirs = append(dirs, skill.Name)
		if commandNames[skill.Name] {
			continue // command wins the top-level collision
		}
		if err := writeCopilotSkill(opts, result, srcFS, skill.SourcePath, filepath.Join(dest, skill.Name, "SKILL.md"), false); err != nil {
			return nil, err
		}
	}
	for _, file := range commands {
		name := strings.TrimSuffix(file, ".md")
		dirs = append(dirs, name)
		if err := writeCopilotSkill(opts, result, srcFS, "commands/"+file, filepath.Join(dest, name, "SKILL.md"), true); err != nil {
			return nil, err
		}
		if !containsSkillName(skills, name) {
			continue
		}
		refPath := filepath.Join(dest, name, "references", name, "SKILL.md")
		refSource := copilotSkillSource(skills, name)
		if err := writeCopilotSkill(opts, result, srcFS, refSource, refPath, false); err != nil {
			return nil, err
		}
	}
	sort.Strings(dirs)
	return uniqueStrings(dirs), nil
}

func copilotSkillSource(skills []skillSource, name string) string {
	for _, skill := range skills {
		if skill.Name == name {
			return skill.SourcePath
		}
	}
	return ""
}

func containsSkillName(skills []skillSource, name string) bool {
	for _, skill := range skills {
		if skill.Name == name {
			return true
		}
	}
	return false
}

func uniqueStrings(values []string) []string {
	out := values[:0]
	for _, value := range values {
		if len(out) == 0 || out[len(out)-1] != value {
			out = append(out, value)
		}
	}
	return out
}

func writeCopilotSkill(opts Options, result *Result, srcFS fs.FS, srcPath, dst string, userInvocable bool) error {
	raw, err := fs.ReadFile(srcFS, srcPath)
	if err != nil {
		return err
	}
	entry := canonicalEntry{SourcePath: srcPath, Raw: raw}
	entry.Name = strings.TrimSuffix(filepath.Base(srcPath), ".md")
	entry.Frontmatter, entry.Body = parseSimpleFrontmatter(raw)
	rendered := renderCopilotSkill(entry, userInvocable)
	result.rendered = append(result.rendered, dst)
	if opts.DryRun {
		logRendered(opts, dst, "copilot skill")
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	if existing, err := os.ReadFile(dst); err == nil && bytes.Equal(existing, rendered) {
		return nil
	}
	if err := os.WriteFile(dst, rendered, 0o644); err != nil {
		return err
	}
	result.Copied = append(result.Copied, CopyAction{Source: srcPath, Dest: dst})
	logRendered(opts, dst, "copilot skill")
	return nil
}
