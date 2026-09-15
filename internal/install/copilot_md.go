package install

import (
	"fmt"
	"os"
	"path/filepath"
)

// resolveCopilotPaths returns the project-local Copilot instruction file.
// Copilot has no supported global project-instruction path in this installer.
func resolveCopilotPaths(opts Options) (string, error) {
	if opts.Mode != ModeProject || opts.TargetDir == "" {
		return "", fmt.Errorf("copilot target only supports project mode")
	}
	info, err := os.Stat(opts.TargetDir)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("target directory does not exist: %s", opts.TargetDir)
	}
	return filepath.Join(opts.TargetDir, ".github", "copilot-instructions.md"), nil
}

// installCopilotMd writes Hero's managed block into Copilot's native
// project instruction file, preserving user content outside the region.
func installCopilotMd(opts Options, result *Result, path string) error {
	return installManagedMarkdown(opts, result, installManagedSpec{
		Path:      path,
		Label:     ".github/copilot-instructions.md",
		DefaultH1: defaultInstructionH1(path),
		Sections:  defaultSections(opts, path),
	})
}
