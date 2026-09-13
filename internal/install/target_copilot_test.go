package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
