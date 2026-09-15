package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hero-engine/hero/internal/managed"
)

func TestCopilotCommandSkillWinsReferenceCollision(t *testing.T) {
	source := t.TempDir()
	target := t.TempDir()
	write := func(name, body string) {
		path := filepath.Join(source, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("commands/design.md", "---\ndescription: command\n---\nRun the workflow.\n")
	write("skills/design/SKILL.md", "---\nname: design\ndescription: reference\n---\nReference guidance.\n")

	if _, err := Run(Options{SourceDir: source, Target: TargetCopilot, Mode: ModeProject, TargetDir: target, Force: true}); err != nil {
		t.Fatal(err)
	}
	commandPath := filepath.Join(target, ".github", "skills", "design", "SKILL.md")
	referencePath := filepath.Join(target, ".github", "skills", "design", "references", "design", "SKILL.md")
	command, err := os.ReadFile(commandPath)
	if err != nil {
		t.Fatal(err)
	}
	reference, err := os.ReadFile(referencePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(command), "user-invocable: true") || !strings.Contains(string(command), "Run the workflow.") {
		t.Fatalf("command skill did not win collision: %s", command)
	}
	if !strings.Contains(string(command), "read `references/design/SKILL.md`") {
		t.Fatalf("command skill does not load its preserved reference: %s", command)
	}
	if !strings.Contains(string(reference), "user-invocable: false") || !strings.Contains(string(reference), "Reference guidance.") {
		t.Fatalf("reference skill was not preserved as disabled nested instructions: %s", reference)
	}
}

func TestCopilotReferenceSkillPreservesMetadata(t *testing.T) {
	source := t.TempDir()
	target := t.TempDir()
	path := filepath.Join(source, "skills", "agent-reliability", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: agent-reliability\ndescription: Reliability rules.\nmetadata:\n  audience: all-agents\n  purpose: reliability-rules\n---\n\nReference guidance.\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(Options{SourceDir: source, Target: TargetCopilot, Mode: ModeProject, TargetDir: target, Force: true}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(target, ".github", "skills", "agent-reliability", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(got)
	for _, want := range []string{"user-invocable: false", "audience: all-agents", "purpose: reliability-rules"} {
		if !strings.Contains(text, want) {
			t.Errorf("rendered skill missing %q:\n%s", want, text)
		}
	}
}

// A freshly generated .github/copilot-instructions.md opens with a document
// heading, not a filename echo — parity with AGENTS.md's "# AGENTS.md" but
// readable as a title.
func TestCopilotInstructions_HasDocumentHeading(t *testing.T) {
	h := newInstallHarness(t)
	if err := os.MkdirAll(filepath.Join(h.TargetDir, ".hero"), 0o755); err != nil {
		t.Fatal(err)
	}
	h.Run(TargetCopilot, nil)

	path := filepath.Join(h.TargetDir, ".github", "copilot-instructions.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read copilot-instructions.md: %v", err)
	}
	content := string(data)
	first := strings.SplitN(strings.TrimSpace(content), "\n", 2)[0]
	if first != "# GitHub Copilot Instructions" {
		t.Errorf("first line = %q, want %q", first, "# GitHub Copilot Instructions")
	}
	if strings.Index(content, "# GitHub Copilot Instructions") > strings.Index(content, "hero:managed-start") {
		t.Errorf("heading must precede the managed region")
	}

	// The heading Hero writes itself must not make the file look
	// user-authored to the prune predicate.
	if !InstructionFileIsPrunable(path) {
		t.Errorf("a Hero-generated copilot-instructions.md must read as Hero-managed-only")
	}
}

// A file created by an earlier Hero (base-name heading) still reads as
// Hero-managed-only, so the heading change doesn't strand it as unprunable.
func TestCopilotInstructions_LegacyHeadingStillPrunable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "copilot-instructions.md")
	body := "# copilot-instructions.md\n\n" + managed.RenderManagedRegion("v0", "X") + "\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if !InstructionFileIsPrunable(path) {
		t.Errorf("legacy base-name heading must still count as Hero-managed-only")
	}
}

// User content outside the markers is never prunable, heading or not.
func TestCopilotInstructions_UserContentNotPrunable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "copilot-instructions.md")
	body := "# GitHub Copilot Instructions\n\n" + managed.RenderManagedRegion("v0", "X") + "\nUSER KEEP\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if InstructionFileIsPrunable(path) {
		t.Errorf("user content outside the markers must never be prunable")
	}
}
