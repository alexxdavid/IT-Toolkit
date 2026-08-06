package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"ittoolkit/internal/actions"
	"ittoolkit/internal/catalog"
	"ittoolkit/internal/githublib"
	"ittoolkit/internal/platform"
	"ittoolkit/internal/update"
)

// App is the Wails-bound application struct.
type App struct {
	ctx      context.Context
	store    *catalog.Store
	scanner  *catalog.Scanner
	scanMu   sync.Mutex
	scanning bool
	updateMgr *update.Manager
}

// NewApp creates the application.
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	debugLog("startup: begin — Solutions IT Toolkit")
	a.updateMgr = update.NewManager()
	dbPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "ITToolkit", "index.db")
	if dbPath == "\\ITToolkit\\index.db" {
		dbPath = filepath.Join(platform.HomeDir(), "AppData", "Local", "ITToolkit", "index.db")
	}
	debugLog("startup: db=" + dbPath)
	store, err := catalog.OpenStore(dbPath)
	if err != nil {
		debugLog("startup: ERROR " + err.Error())
		fmt.Println("Error opening catalog store:", err)
		return
	}
	a.store = store
	debugLog("startup: store ok")
	a.scanner = catalog.NewScanner(store, func(folderID int64, folderPath, phase string, current, total int) {
		runtime.EventsEmit(a.ctx, "scan:progress", map[string]any{
			"folderId": folderID, "folder": folderPath, "phase": phase,
			"current": current, "total": total,
		})
	})
	debugLog("startup: Solutions IT Toolkit done")
}

// Ping is a trivial connectivity check.
func (a *App) Ping() string { return "pong" }

// Log writes a line from the frontend to the debug log (for diagnosing UI issues).
func (a *App) Log(msg string) { debugLog("JS: " + msg) }

// PickFolder opens a native directory chooser dialog.
func (a *App) PickFolder() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a scripts folder",
	})
}

// GetCatalog returns the full catalog snapshot.
func (a *App) GetCatalog() (*catalog.CatalogView, error) {
	debugLog("GetCatalog: called")
	if a.store == nil {
		debugLog("GetCatalog: store is nil")
		return nil, fmt.Errorf("catalog store is not initialized")
	}
	v, err := a.store.CatalogView()
	if err != nil {
		debugLog("GetCatalog: ERROR " + err.Error())
		return nil, err
	}
	debugLog(fmt.Sprintf("GetCatalog: ok folders=%d repos=%d", len(v.Folders), len(v.Repos)))
	return v, nil
}

// debugLog appends a timestamped line to the app debug log.
func debugLog(msg string) {
	dir := filepath.Join(os.Getenv("LOCALAPPDATA"), "ITToolkit")
	_ = os.MkdirAll(dir, 0o755)
	f, err := os.OpenFile(filepath.Join(dir, "debug.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(time.Now().Format("15:04:05.000") + " " + msg + "\n")
}

// RescanAll rescans every registered folder in the background.
func (a *App) RescanAll() string {
	a.scanMu.Lock()
	if a.scanning {
		a.scanMu.Unlock()
		return "already scanning"
	}
	a.scanning = true
	a.scanMu.Unlock()
	go func() {
		defer func() {
			a.scanMu.Lock()
			a.scanning = false
			a.scanMu.Unlock()
			runtime.EventsEmit(a.ctx, "scan:done", map[string]any{})
		}()
		folders, err := a.store.ListFolders()
		if err != nil {
			return
		}
		for _, f := range folders {
			if err := a.scanner.ScanFolder(f.ID, f.Path); err != nil {
				runtime.EventsEmit(a.ctx, "scan:error", map[string]any{"folder": f.Path, "error": err.Error()})
			}
		}
	}()
	return "started"
}

// AddFolder registers a new watched folder and scans it.
func (a *App) AddFolder(path string) (*catalog.Folder, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %s", abs)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", abs)
	}
	f, err := a.store.AddFolder(abs, filepath.Base(abs), false)
	if err != nil {
		return nil, err
	}
	a.RescanAll()
	return f, nil
}

// RemoveFolder removes a watched folder and its indexed data.
func (a *App) RemoveFolder(id int64) error {
	return a.store.RemoveFolder(id)
}

// GetRepoDetail returns a repo and its indexed scripts.
func (a *App) GetRepoDetail(id int64) (*catalog.RepoDetail, error) {
	return a.store.RepoDetail(id)
}

// GetScriptContent returns full text of a script.
func (a *App) GetScriptContent(path string) (string, error) {
	return a.store.ScriptContent(path)
}

// Search queries repo/script names and optionally indexed content.
func (a *App) Search(query string, inContent bool) (*catalog.SearchResult, error) {
	res, err := a.store.SearchName(query)
	if err != nil {
		return nil, err
	}
	if inContent {
		scripts, err := a.store.SearchContent(query)
		if err == nil {
			res.Scripts = append(res.Scripts, scripts...)
		}
	}
	return res, nil
}

// ListScripts returns scripts filtered by language and/or name query.
func (a *App) ListScripts(lang, q string) ([]*catalog.ScriptFile, error) {
	return a.store.ListScripts(lang, q)
}

// RunScript launches a script in an elevated window.
func (a *App) RunScript(path string) (*catalog.RunResult, error) {
	return actions.RunScript(path)
}

// CopyScript returns script text for the clipboard.
func (a *App) CopyScript(path string) (string, error) {
	return actions.CopyScript(path)
}

// ExportScripts copies selected scripts into a destination folder.
func (a *App) ExportScripts(paths []string, dest string) (*catalog.ExportResult, error) {
	return actions.ExportScripts(a.store, paths, dest)
}

// RevealInExplorer opens/selects a path in Explorer.
func (a *App) RevealInExplorer(path string) error {
	return actions.RevealInExplorer(path)
}

// GetExeDir returns the directory containing the executable.
func (a *App) GetExeDir() string { return platform.ExeDir() }

// GetDefaultScriptsDir returns the default Scripts folder path.
func (a *App) GetDefaultScriptsDir() string { return platform.DefaultScriptsDir() }

// IsWebView2Installed reports the WebView2 runtime state.
func (a *App) IsWebView2Installed() bool { return platform.WebView2Installed() }

// DiscoverRepos lists cloned GitHub repositories inside a folder by reading
// each subdirectory's git remote.
func (a *App) DiscoverRepos(folder string) ([]*githublib.RepoInfo, error) {
	return githublib.Discover(folder)
}

// GitAvailable reports whether git is installed on this machine.
func (a *App) GitAvailable() bool { return githublib.GitAvailable() }

// GetDefaultInstallDir returns the default install destination (Repo subfolder).
func (a *App) GetDefaultInstallDir() string { return githublib.DefaultInstallDir() }

// GetDefaultSoftwareDir returns the default Software folder path.
func (a *App) GetDefaultSoftwareDir() string { return githublib.DefaultSoftwareDir() }

// GetRecommendedRepos returns the curated list of sysadmin GitHub repos.
func (a *App) GetRecommendedRepos() []*githublib.RecommendedRepo { return githublib.GetRecommendedRepos() }

// IsRepoInstalled checks if a repo is already downloaded.
func (a *App) IsRepoInstalled(name, dest string) bool { return githublib.IsRepoInstalled(name, dest) }

// RemoveRepo deletes an installed repo folder.
func (a *App) RemoveRepo(name, dest string) error { return githublib.RemoveRepo(name, dest) }

// CheckForUpdate checks the GitHub Gist for a new version.
func (a *App) CheckForUpdate(force bool) (update.Info, error) {
	return a.updateMgr.Check(force)
}

// DownloadUpdate downloads the installer from the manifest URL.
func (a *App) DownloadUpdate(url, version string) (string, error) {
	return a.updateMgr.Download(url, version)
}

// GetUpdateProgress returns download progress for the frontend polling loop.
func (a *App) GetUpdateProgress() map[string]interface{} {
	return a.updateMgr.Progress()
}

// ApplyUpdate launches the installer and exits this process.
func (a *App) ApplyUpdate(installerPath string) error {
	return update.Install(installerPath)
}

// GetCurrentVersion returns the version baked into this binary.
func (a *App) GetCurrentVersion() string { return update.CurrentVersion }

// --- Custom Categories ---

func (a *App) ListCustomCategories() ([]*catalog.CustomCategory, error) {
	return a.store.ListCustomCategories()
}

func (a *App) CreateCustomCategory(name string) (*catalog.CustomCategory, error) {
	return a.store.CreateCustomCategory(name)
}

func (a *App) RenameCustomCategory(id int64, newName string) error {
	return a.store.RenameCustomCategory(id, newName)
}

func (a *App) DeleteCustomCategory(id int64) error {
	return a.store.DeleteCustomCategory(id)
}

// --- Custom Repos ---

func (a *App) ListCustomRepos() ([]*catalog.CustomRepo, error) {
	return a.store.ListCustomRepos()
}

func (a *App) AddCustomRepo(name, url, category, summary string) (*catalog.CustomRepo, error) {
	githublib.RegisterURL(name, url)
	return a.store.AddCustomRepo(name, url, category, summary)
}

func (a *App) RemoveCustomRepo(id int64) error {
	return a.store.RemoveCustomRepo(id)
}

// GetRecommendedReposCombined returns recommended + custom repos merged.
func (a *App) GetRecommendedReposCombined() ([]*githublib.RecommendedRepo, error) {
	rec := githublib.GetRecommendedRepos()
	custom, err := a.store.ListCustomRepos()
	if err != nil {
		return rec, err
	}
	for _, c := range custom {
		githublib.RegisterURL(c.Name, c.URL)
		rec = append(rec, &githublib.RecommendedRepo{
			Name:     c.Name,
			URL:      c.URL,
			Category: c.Category,
			Summary:  c.Summary,
		})
	}
	return rec, nil
}

// --- Software Catalog ---

func (a *App) GetSoftwareCatalog() []*githublib.SoftwareItem {
	builtin := githublib.GetSoftwareCatalog()
	custom, err := a.store.ListCustomSoftware()
	if err != nil {
		return builtin
	}
	for _, c := range custom {
		builtin = append(builtin, &githublib.SoftwareItem{
			Name:     c.Name,
			Version:  c.Version,
			Category: c.Category,
			Download: c.Download,
			Notes:    c.Notes,
			WingetID: c.WingetID,
		})
	}
	return builtin
}

// DownloadSoftware downloads a software installer to the specified folder.
func (a *App) DownloadSoftware(name, url, dest string) (string, error) {
	return githublib.DownloadSoftware(name, url, dest)
}

// SoftwareProgress reports download progress.
func (a *App) GetSoftwareProgress() map[string]interface{} {
	return githublib.SoftwareProgress()
}

// GetSoftwareVersions returns latest versions (via winget), installed status,
// and download status for each software item (built-in + custom).
func (a *App) GetSoftwareVersions() map[string]githublib.SoftwareVersion {
	softDir := filepath.Join(platform.ExeDir(), "Software")
	return githublib.GetSoftwareVersions(a.GetSoftwareCatalog(), softDir)
}

// InvalidateSoftwareVersions clears the version cache (e.g. after a download).
func (a *App) InvalidateSoftwareVersions() {
	githublib.InvalidateVersionCache()
}

// --- Custom Software ---

func (a *App) ListCustomSoftware() ([]*catalog.CustomSoftware, error) {
	return a.store.ListCustomSoftware()
}

func (a *App) AddCustomSoftware(name, version, category, download, notes, wingetID string) (*catalog.CustomSoftware, error) {
	return a.store.AddCustomSoftware(name, version, category, download, notes, wingetID)
}

func (a *App) RemoveCustomSoftware(id int64) error {
	return a.store.RemoveCustomSoftware(id)
}

// InstallRepos installs repos by name (URLs resolved server-side).
func (a *App) InstallRepos(names []string, dest string) []*githublib.InstallResult {
	results := githublib.InstallByName(names, dest, func(index, total int, name, status, message string) {
		runtime.EventsEmit(a.ctx, "install:progress", map[string]any{
			"index": index, "total": total, "name": name,
			"status": status, "message": message,
		})
	})
	runtime.EventsEmit(a.ctx, "install:done", map[string]any{"count": len(results)})
	return results
}
