package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expectMode cliMode
		expectPath string
		errorText  string
	}{
		{name: "no args enters repl", args: []string{}, expectMode: modeREPL},
		{name: "valid run command", args: []string{"run", "script.bas"}, expectMode: modeRunFile, expectPath: "script.bas"},
		{name: "missing run target", args: []string{"run"}, errorText: "missing script path for run command"},
		{name: "too many args", args: []string{"run", "a.bas", "extra"}, errorText: "run expects exactly one script path"},
		{name: "unknown command", args: []string{"debug", "script.bas"}, errorText: "unknown command \"debug\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, path, err := parseCommand(tt.args)
			if tt.errorText != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.errorText)
				}
				if err.Error() != tt.errorText {
					t.Fatalf("expected error %q, got %q", tt.errorText, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if mode != tt.expectMode {
				t.Fatalf("expected mode %v, got %v", tt.expectMode, mode)
			}
			if path != tt.expectPath {
				t.Fatalf("expected path %q, got %q", tt.expectPath, path)
			}
		})
	}
}

func TestRunNoArgsUsesREPLPath(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{}, strings.NewReader("quit\n"), &stdout, &stderr, "cart-basic")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "mb> ") {
		t.Fatalf("expected repl prompt in output, got %q", out)
	}
	if !strings.Contains(out, "Goodbye") {
		t.Fatalf("expected goodbye in output, got %q", out)
	}
}

func TestRunCommandFileErrors(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		errorText string
	}{
		{name: "missing target argument", args: []string{"run"}, errorText: "missing script path for run command"},
		{name: "missing file", args: []string{"run", filepath.Join(t.TempDir(), "missing.bas")}, errorText: "file error: cannot read"},
		{name: "unreadable target directory", args: []string{"run", t.TempDir()}, errorText: "file error: cannot read"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			exitCode := run(tt.args, strings.NewReader(""), &stdout, &stderr, "cart-basic")
			if exitCode == 0 {
				t.Fatalf("expected non-zero exit code")
			}
			if !strings.Contains(stderr.String(), tt.errorText) {
				t.Fatalf("expected stderr containing %q, got %q", tt.errorText, stderr.String())
			}
		})
	}
}

func TestRunCommandParseAndRuntimeFailures(t *testing.T) {
	tests := []struct {
		name      string
		script    string
		errorText string
	}{
		{name: "parse failure", script: "let x = ;", errorText: "parse failure"},
		{name: "runtime failure", script: "1 / 0;", errorText: "runtime failure"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptPath := filepath.Join(t.TempDir(), "script.bas")
			if err := os.WriteFile(scriptPath, []byte(tt.script), 0o644); err != nil {
				t.Fatalf("failed to write script: %v", err)
			}

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			exitCode := run([]string{"run", scriptPath}, strings.NewReader(""), &stdout, &stderr, "cart-basic")
			if exitCode == 0 {
				t.Fatalf("expected non-zero exit code")
			}

			if !strings.Contains(stderr.String(), tt.errorText) {
				t.Fatalf("expected stderr containing %q, got %q", tt.errorText, stderr.String())
			}
		})
	}
}

func TestRunCommandSuccess(t *testing.T) {
	scriptPath := filepath.Join(t.TempDir(), "ok.bas")
	if err := os.WriteFile(scriptPath, []byte("let x = 1 + 2; x;"), 0o644); err != nil {
		t.Fatalf("failed to write script: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{"run", scriptPath}, strings.NewReader(""), &stdout, &stderr, "cart-basic")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d with stderr %q", exitCode, stderr.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}
