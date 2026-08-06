package catalog

// Folder is a watched root directory (e.g. the Scripts folder next to the EXE).
type Folder struct {
	ID          int64  `json:"id"`
	Path        string `json:"path"`
	DisplayName string `json:"displayName"`
	IsDefault   bool   `json:"isDefault"`
	LastScanned string `json:"lastScanned"`
	RepoCount   int    `json:"repoCount"`
	ScriptCount int    `json:"scriptCount"`
}

// Repo is a top-level directory inside a watched folder.
type Repo struct {
	ID          int64  `json:"id"`
	FolderID    int64  `json:"folderId"`
	FolderPath  string `json:"folderPath"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Category    string `json:"category"`
	ReadmePath  string `json:"readmePath"`
	ReadmeText  string `json:"readmeText"`
	FileCount   int    `json:"fileCount"`
	ScriptCount int    `json:"scriptCount"`
}

// ScriptFile is a single indexed file inside a repo.
type ScriptFile struct {
	ID       int64  `json:"id"`
	RepoID   int64  `json:"repoId"`
	Repo     string `json:"repo"`
	RelPath  string `json:"relPath"`
	AbsPath  string `json:"absPath"`
	Name     string `json:"name"`
	Ext      string `json:"ext"`
	Lang     string `json:"lang"`
	Size     int64  `json:"size"`
	MTime    int64  `json:"mtime"`
	IsScript bool   `json:"isScript"`
	Snippet  string `json:"snippet"`
}

// Category is a named group of repos with a count.
type Category struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// CatalogView is the full snapshot returned to the frontend.
type CatalogView struct {
	Folders      []*Folder   `json:"folders"`
	Categories   []*Category `json:"categories"`
	Repos        []*Repo     `json:"repos"`
	TotalScripts int         `json:"totalScripts"`
}

// RepoDetail is a repo plus its indexed scripts.
type RepoDetail struct {
	Repo    *Repo         `json:"repo"`
	Scripts []*ScriptFile `json:"scripts"`
}

// SearchResult contains matched repos and scripts.
type SearchResult struct {
	Repos   []*Repo       `json:"repos"`
	Scripts []*ScriptFile `json:"scripts"`
}

// ExportResult summarizes a multi-file export operation.
type ExportResult struct {
	Copied  int      `json:"copied"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
}

// RunResult describes a launch attempt of a script.
type RunResult struct {
	Command string `json:"command"`
	Launched bool  `json:"launched"`
	Message string `json:"message"`
}
