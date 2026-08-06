package main

import (
	"os"
	"path/filepath"
	"strings"

	"ittoolkit/internal/platform"
)

// LocalItem is a file or folder visible in the Local tab.
type LocalItem struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
	Kind  string `json:"kind"` // "repo" | "software" | "file"
}

// ListLocalRepos returns downloaded repos from the Repo folder.
func (a *App) ListLocalRepos() []*LocalItem {
	dir := filepath.Join(platform.ExeDir(), "Repo")
	return listDir(dir, "repo")
}

// ListLocalSoftware returns downloaded software from the Software folder.
func (a *App) ListLocalSoftware() []*LocalItem {
	dir := filepath.Join(platform.ExeDir(), "Software")
	return listDir(dir, "software")
}

// GetRepoContents returns files inside a downloaded repo folder.
func (a *App) GetRepoContents(repoPath string) []*LocalItem {
	return listDir(repoPath, "file")
}

func listDir(dir, kind string) []*LocalItem {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []*LocalItem{}
	}
	var out []*LocalItem
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, &LocalItem{
			Name:  name,
			Path:  filepath.Join(dir, name),
			IsDir: e.IsDir(),
			Size:  info.Size(),
			Kind:  kind,
		})
	}
	if out == nil {
		out = []*LocalItem{}
	}
	return out
}
