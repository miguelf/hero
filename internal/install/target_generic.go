package install

import (
	"fmt"
	"os"
	"path/filepath"
)

// runGeneric installs to a tool-agnostic layout using .ai/ for agents,
// commands, skills and AGENTS.md at the project root. `.ai/` is a Hero
// convention with no consuming loader — it's the catch-all for tools
// without a dedicated installer. Any MCP-capable tool can pick up Hero
// via .mcp.json (registered in Run()).

func runGeneric(opts Options) (*Result, error) {
	if opts.Mode != ModeProject || opts.TargetDir == "" {
		return nil, fmt.Errorf("generic target only supports project mode")
	}

	info, statErr := os.Stat(opts.TargetDir)
	if statErr != nil || !info.IsDir() {
		return nil, fmt.Errorf("target directory does not exist: %s", opts.TargetDir)
	}

	destBase := filepath.Join(opts.TargetDir, ".ai")
	result := &Result{}

	if err := installFlat(opts, result, "agents", filepath.Join(destBase, "agents")); err != nil {
		return nil, fmt.Errorf("installing agents: %w", err)
	}
	if err := installFlat(opts, result, "commands", filepath.Join(destBase, "commands")); err != nil {
		return nil, fmt.Errorf("installing commands: %w", err)
	}
	if err := installSkillsNested(opts, result, filepath.Join(destBase, "skills")); err != nil {
		return nil, fmt.Errorf("installing skills: %w", err)
	}
	if err := pruneNestedSkills(opts, result, filepath.Join(destBase, "skills")); err != nil {
		return nil, fmt.Errorf("prune stale skills: %w", err)
	}

	// Generic's root instruction file is AGENTS.md, written with the same
	// Hero-managed block as every other non-Claude target (harness-native
	// mapping) rather than the minimal legacy stub.
	if err := installNativeInstructionFile(opts, result); err != nil {
		return nil, fmt.Errorf("installing AGENTS.md: %w", err)
	}

	return result, nil
}
