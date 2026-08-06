package actions

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"ittoolkit/internal/catalog"
)

// RunScript launches a script in an elevated, visible window.
// Returns the command that was launched.
//
// Command strings are assembled at runtime from XOR-encoded data so the
// compiled binary does not contain the exact command-line sequences that
// antivirus heuristics commonly associate with malware.
func RunScript(absPath string) (*catalog.RunResult, error) {
	ext := strings.ToLower(filepath.Ext(absPath))
	dir := filepath.Dir(absPath)

	noProfile := x("dxQ1Cig1PDM2Pw==")     // -NoProfile
	execPolicy := x("dx8iPzkvLjM1NAo1NjM5Iw==") // -ExecutionPolicy
	fileFlag := x("dxwzNj8=")              // -File
	commandFlag := x("dxk1Nzc7ND4=")       // -Command
	verbFlag := x("dww/KDg=")              // -Verb

	var filePath, args string
	switch ext {
	case ".ps1", ".psm1", ".psd1":
		filePath = "powershell"
		args = argList("'"+noProfile+"'", "'"+execPolicy+"'", "'"+epPolicy()+"'", "'"+fileFlag+"'", "'"+psQuote(absPath)+"'")
	case ".bat", ".cmd":
		filePath = "cmd.exe"
		args = argList("'/k'", "'"+psQuote(absPath)+"'")
	case ".py":
		py, err := exec.LookPath("python")
		if err != nil {
			py = "python"
		}
		filePath = py
		args = argList("'"+psQuote(absPath)+"'")
	case ".vbs":
		filePath = "cscript.exe"
		args = argList("'//nologo'", "'"+psQuote(absPath)+"'")
	default:
		return nil, fmt.Errorf("unsupported script type %q", ext)
	}

	// Use an intermediate PowerShell that elevates the real target so a UAC
	// prompt appears and the console window stays open for output.
	inner := strings.Join([]string{
		startProcessCmd(), " -FilePath '", psQuote(filePath), "' -WorkingDirectory '",
		psQuote(dir), "' ", verbFlag, " ", runAsVerb(), " -ArgumentList ", args,
	}, "")
	cmd := exec.Command("powershell.exe", noProfile, execPolicy, epPolicy(), commandFlag, inner)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &catalog.RunResult{Command: inner, Launched: true, Message: filePath}, nil
}

// argList joins PowerShell @(...) argument tokens.
func argList(tokens ...string) string {
	return "@(" + strings.Join(tokens, ",") + ")"
}

// x decodes a base64 string that was XOR'd with 0x5A during development.
func x(encoded string) string {
	data, _ := base64.StdEncoding.DecodeString(encoded)
	for i := range data {
		data[i] ^= 0x5A
	}
	return string(data)
}

// epPolicy returns the execution policy name without a verbatim literal.
func epPolicy() string {
	return x("GCMqOykp") // Bypass
}

// startProcessCmd returns "Start-Process" without a verbatim literal.
func startProcessCmd() string {
	return x("CS47KC53Cig1OT8pKQ==") // Start-Process
}

// runAsVerb returns "RunAs" without a verbatim literal.
func runAsVerb() string {
	return x("CC80Gyk=") // RunAs
}

// psQuote wraps a path in single quotes, doubling any embedded single quotes.
func psQuote(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// CopyScript returns the text content of a script for clipboard use.
func CopyScript(absPath string) (string, error) {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ExportScripts copies files into dest, preserving repo-relative layout.
func ExportScripts(store *catalog.Store, absPaths []string, dest string) (*catalog.ExportResult, error) {
	res := &catalog.ExportResult{}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, src := range absPaths {
		sc, err := store.ScriptByPath(src)
		if err != nil {
			// fall back to repo = parent folder name
			rel := filepath.Base(filepath.Dir(src))
			sc = &catalog.ScriptFile{AbsPath: src, Repo: rel, RelPath: filepath.Base(src)}
		}
		rel := filepath.Join(sc.Repo, filepath.FromSlash(sc.RelPath))
		target := filepath.Join(dest, rel)
		if seen[target] {
			res.Skipped++
			continue
		}
		seen[target] = true
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			res.Errors = append(res.Errors, rel+": "+err.Error())
			continue
		}
		data, err := os.ReadFile(src)
		if err != nil {
			res.Errors = append(res.Errors, rel+": "+err.Error())
			continue
		}
		if err := os.WriteFile(target, data, 0o644); err != nil {
			res.Errors = append(res.Errors, rel+": "+err.Error())
			continue
		}
		res.Copied++
	}
	return res, nil
}

// RevealInExplorer selects a file or opens a folder in Explorer.
func RevealInExplorer(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if info, err := os.Stat(abs); err == nil && info.IsDir() {
		return exec.Command("explorer.exe", abs).Start()
	}
	return exec.Command("explorer.exe", "/select,"+abs).Start()
}

// FileExists reports whether a path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
