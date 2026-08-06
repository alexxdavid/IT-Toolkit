package catalog

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ProgressFunc reports scan progress to the UI layer.
type ProgressFunc func(folderID int64, folderPath, phase string, current, total int)

// Scanner walks watched folders and keeps the catalog index in sync.
type Scanner struct {
	store    *Store
	progress ProgressFunc
}

// NewScanner creates a Scanner bound to a store.
func NewScanner(store *Store, progress ProgressFunc) *Scanner {
	return &Scanner{store: store, progress: progress}
}

const maxIndexSize = 5 * 1024 * 1024

// maxContentIndex caps how much of each file is stored/searchable. The head of
// a script is enough for content search; pathological files (minified output,
// data dumps) would otherwise dominate indexing time.
const maxContentIndex = 256 * 1024

var skipDirs = map[string]bool{
	".git": true, "node_modules": true, ".github": true, ".vscode": true,
	".idea": true, ".svn": true, ".hg": true, ".cache": true, ".venv": true,
	"venv": true, "vendor": true, "dist": true, "build": true, "bin": true,
	"obj": true, "__pycache__": true,
}

var scriptExts = map[string]string{
	".ps1": "PowerShell", ".psm1": "PowerShell", ".psd1": "PowerShell",
	".py": "Python", ".bat": "Batch", ".cmd": "Batch",
	".vbs": "VBScript", ".sh": "Shell", ".sql": "SQL", ".reg": "Registry",
}

// contentExts are file types whose full text is read and indexed for content
// search. Python is deliberately excluded: most .py files in the catalog come
// from large library/AI repositories, and reading 400+ MB of library source
// makes the first index impractically slow. Python/SQL files are still listed,
// name-searchable and viewable (content is read from disk on demand).
var contentExts = map[string]bool{
	".ps1": true, ".psm1": true, ".psd1": true,
	".bat": true, ".cmd": true, ".vbs": true, ".sh": true, ".sql": true, ".reg": true,
}

func langForExt(ext string) string {
	if l, ok := scriptExts[ext]; ok {
		return l
	}
	return "text"
}

func isScriptExt(ext string) bool { _, ok := scriptExts[ext]; return ok }

// Only executable script files are indexed. READMEs and other documentation
// are intentionally not indexed to keep catalog scans fast.
func isIndexableExt(ext string) bool { return isScriptExt(ext) }

// ScanFolder indexes every repo directory under path for the given folder id.
// The whole folder is processed inside a single SQLite transaction for speed.
func (sc *Scanner) ScanFolder(folderID int64, path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	seenRepoNames := map[string]bool{}
	var repoDirs []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if skipDirs[strings.ToLower(name)] {
			continue
		}
		seenRepoNames[name] = true
		repoDirs = append(repoDirs, e)
	}
	sort.Slice(repoDirs, func(i, j int) bool { return repoDirs[i].Name() < repoDirs[j].Name() })

	tx, err := sc.store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Drop repos that no longer exist on disk.
	existing, err := listRepoIDsTx(tx, folderID)
	if err != nil {
		return err
	}
	for _, rid := range existing {
		var rp string
		if err := tx.QueryRow(`SELECT name FROM repos WHERE id = ?`, rid).Scan(&rp); err != nil {
			return err
		}
		if !seenRepoNames[rp] {
			if err := sc.store.deleteRepoData(tx, rid); err != nil {
				return err
			}
		}
	}

	total := len(repoDirs)
	for i, e := range repoDirs {
		name := e.Name()
		repoPath := filepath.Join(path, name)
		if sc.progress != nil {
			sc.progress(folderID, path, "repos", i, total)
		}
		if err := sc.scanRepo(tx, folderID, name, repoPath); err != nil {
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	if sc.progress != nil {
		sc.progress(folderID, path, "repos", total, total)
	}
	return sc.store.SetFolderScanned(folderID, time.Now().Unix())
}

func (sc *Scanner) scanRepo(tx *sql.Tx, folderID int64, name, repoPath string) error {
	category := Categorize(name)
	readmePath, readmeText := sc.findReadme(repoPath)

	repoID, err := upsertRepoTx(tx, folderID, name, repoPath, category, readmePath, readmeText)
	if err != nil {
		return err
	}

	seen := map[string]bool{}
	var files []struct {
		abs, rel, ext string
		size, mtime   int64
		content       bool
	}
	err = filepath.Walk(repoPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if p != repoPath && skipDirs[strings.ToLower(info.Name())] {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if !isIndexableExt(ext) {
			return nil
		}
		if info.Size() > maxIndexSize {
			return nil
		}
		rel, err := filepath.Rel(repoPath, p)
		if err != nil {
			return nil
		}
		seen[p] = true
		files = append(files, struct {
			abs, rel, ext string
			size, mtime   int64
			content       bool
		}{p, rel, ext, info.Size(), info.ModTime().Unix(), contentExts[ext]})
		return nil
	})
	if err != nil {
		return err
	}

	for _, f := range files {
		lang := langForExt(f.ext)
		snippet := ""
		content := ""
		if f.content {
			data, err := os.ReadFile(f.abs)
			if err != nil {
				continue
			}
			text := string(data)
			content = text
			if len(content) > maxContentIndex {
				content = content[:maxContentIndex]
			}
			snippet = strings.TrimSpace(text)
			if len(snippet) > 500 {
				snippet = snippet[:500]
			}
		}
		if _, err := upsertScriptTx(tx, repoID, f.rel, f.abs, filepath.Base(f.abs), f.ext, lang,
			f.size, f.mtime, true, snippet, content, f.content); err != nil {
			continue
		}
	}

	_, err = deleteScriptsNotInTx(tx, repoID, seen)
	return err
}

func (sc *Scanner) findReadme(repoPath string) (string, string) {
	preferred := []string{"readme.md", "readme.markdown", "readme.rst", "readme.txt", "readme"}
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return "", ""
	}
	found := ""
	lower := map[string]string{}
	for _, e := range entries {
		if !e.IsDir() {
			lower[strings.ToLower(e.Name())] = e.Name()
		}
	}
	for _, p := range preferred {
		if actual, ok := lower[p]; ok {
			found = filepath.Join(repoPath, actual)
			break
		}
	}
	if found == "" {
		if actual, ok := lower["docs"]; ok {
			for _, p := range preferred {
				cand := filepath.Join(repoPath, actual, p)
				if _, err := os.Stat(cand); err == nil {
					found = cand
					break
				}
			}
		}
	}
	if found == "" {
		return "", ""
	}
	info, err := os.Stat(found)
	if err != nil || info.Size() > 256*1024 {
		return found, ""
	}
	data, err := os.ReadFile(found)
	if err != nil {
		return found, ""
	}
	text := string(data)
	if len(text) > 65536 {
		text = text[:65536]
	}
	return found, text
}
