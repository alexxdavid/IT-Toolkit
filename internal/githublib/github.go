package githublib

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"ittoolkit/internal/catalog"
)

// RepoInfo describes a discovered GitHub repository.
type RepoInfo struct {
	Name      string `json:"name"` // owner/repo
	URL       string `json:"-"`
	LocalPath string `json:"localPath"`
	Category  string `json:"category"`
}

// InstallResult reports the outcome of one repository install.
type InstallResult struct {
	Name    string `json:"name"`
	URL     string `json:"-"`
	Status  string `json:"status"` // installed | skipped | failed
	Message string `json:"message"`
}

// ProgressFunc reports per-repo install progress.
type ProgressFunc func(index, total int, name, status, message string)

var remoteURLRe = regexp.MustCompile(`^\s*url\s*=\s*(.+?)\s*$`)

// urlRegistry maps lowercase repo name → full GitHub URL. Populated by
// Discover and by the recommended-repo list. InstallByName resolves names
// against this registry so the frontend never sees URLs.
var (
	urlRegistry   = map[string]string{}
	urlRegistryMu sync.RWMutex
)

// RegisterURL stores a name → URL mapping for install-time lookup.
func RegisterURL(name, url string) {
	urlRegistryMu.Lock()
	defer urlRegistryMu.Unlock()
	urlRegistry[strings.ToLower(name)] = url
}

func resolveURL(name string) (string, bool) {
	urlRegistryMu.RLock()
	defer urlRegistryMu.RUnlock()
	u, ok := urlRegistry[strings.ToLower(name)]
	return u, ok
}

// Discover scans a folder for cloned GitHub repositories by reading each
// subdirectory's .git/config remote.
func Discover(dir string) ([]*RepoInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []*RepoInfo
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		local := filepath.Join(dir, e.Name())
		url, ok := readOriginRemote(local)
		if !ok {
			continue
		}
		name, normURL, ok := normalizeGitHubURL(url)
		if !ok {
			continue
		}
		out = append(out, &RepoInfo{
			Name:      name,
			URL:       normURL,
			LocalPath: local,
			Category:  catalog.Categorize(e.Name()),
		})
		RegisterURL(name, normURL)
	}
	return out, nil
}

func readOriginRemote(repoDir string) (string, bool) {
	data, err := os.ReadFile(filepath.Join(repoDir, ".git", "config"))
	if err != nil {
		return "", false
	}
	section := ""
	for _, line := range strings.Split(string(data), "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "[") {
			section = strings.TrimSuffix(strings.TrimPrefix(t, "["), "]")
			continue
		}
		if section == `remote "origin"` {
			if m := remoteURLRe.FindStringSubmatch(t); m != nil {
				return strings.TrimSpace(m[1]), true
			}
		}
	}
	return "", false
}

// normalizeGitHubURL converts any GitHub remote URL to https://github.com/owner/repo
func normalizeGitHubURL(raw string) (name, url string, ok bool) {
	u := strings.TrimSpace(raw)
	u = strings.TrimSuffix(u, "/")
	u = strings.TrimSuffix(u, ".git")
	switch {
	case strings.HasPrefix(u, "git@github.com:"):
		u = "https://github.com/" + strings.TrimPrefix(u, "git@github.com:")
	case strings.HasPrefix(u, "https://github.com/"):
		// keep
	case strings.HasPrefix(u, "http://github.com/"):
		u = "https://" + strings.TrimPrefix(u, "http://")
	default:
		return "", "", false
	}
	u = strings.TrimSuffix(u, "/")
	parts := strings.Split(strings.TrimPrefix(u, "https://github.com/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	name = parts[0] + "/" + parts[1]
	return name, u, true
}

// GitAvailable reports whether git is on PATH.
func GitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// DefaultInstallDir is the folder the app creates for installed repos.
func DefaultInstallDir() string { return filepath.Join(exeDir(), "Repo") }

// DefaultSoftwareDir is the folder the app creates for software downloads.
func DefaultSoftwareDir() string { return filepath.Join(exeDir(), "Software") }

// exeDir returns the directory containing the running executable.
func exeDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

// repoFolderName resolves the on-disk folder for a repo name. The folder is
// derived from the repo URL (owner/repo), so custom names like "MyTool"
// map to the correct installed folder. Falls back to the name itself.
func repoFolderName(name string) string {
	if u, ok := resolveURL(name); ok {
		if _, repo, ok2 := splitGitHubURL(u); ok2 {
			return repo
		}
	}
	repoName := name
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		repoName = name[idx+1:]
	}
	return repoName
}

// IsRepoInstalled checks whether a repo is already downloaded in dest.
func IsRepoInstalled(name, dest string) bool {
	info, err := os.Stat(filepath.Join(dest, repoFolderName(name)))
	return err == nil && info.IsDir()
}

// RemoveRepo deletes an installed repo folder from dest.
func RemoveRepo(name, dest string) error {
	folder := repoFolderName(name)
	if folder == "" || folder == "." || folder == ".." || strings.ContainsAny(folder, `/\`) {
		return fmt.Errorf("invalid repo name")
	}
	target := filepath.Join(dest, folder)
	// Safety: ensure target stays inside dest.
	rel, err := filepath.Rel(dest, target)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("refusing to remove path outside dest")
	}
	return os.RemoveAll(target)
}

// Install clones (or ZIP-downloads) the given repo URLs into dest.
func Install(urls []string, dest string, progress ProgressFunc) []*InstallResult {
	results := make([]*InstallResult, 0, len(urls))
	if err := os.MkdirAll(dest, 0o755); err != nil {
		for _, u := range urls {
			results = append(results, &InstallResult{Name: u, URL: u, Status: "failed", Message: err.Error()})
		}
		return results
	}
	useGit := GitAvailable()
	for i, u := range urls {
		owner, repo, ok := splitGitHubURL(u)
		if !ok {
			results = append(results, &InstallResult{Name: u, URL: u, Status: "failed", Message: "not a GitHub URL"})
			continue
		}
		name := owner + "/" + repo
		target := filepath.Join(dest, repo)
		if _, err := os.Stat(target); err == nil {
			if progress != nil {
				progress(i, len(urls), name, "skipped", "already exists")
			}
			results = append(results, &InstallResult{Name: name, URL: u, Status: "skipped", Message: "folder already exists"})
			continue
		}
		var installErr error
		if useGit {
			installErr = gitClone(u, target)
			if installErr != nil {
				// Fall back to ZIP download (works without git auth/proxy).
				os.RemoveAll(target)
				installErr = zipDownload(owner, repo, target)
			}
		} else {
			installErr = zipDownload(owner, repo, target)
		}
		if installErr != nil {
			os.RemoveAll(target)
			if progress != nil {
				progress(i, len(urls), name, "failed", installErr.Error())
			}
			results = append(results, &InstallResult{Name: name, URL: u, Status: "failed", Message: installErr.Error()})
			continue
		}
		if progress != nil {
			progress(i, len(urls), name, "installed", "")
		}
		results = append(results, &InstallResult{Name: name, URL: u, Status: "installed", Message: "installed to " + target})
	}
	return results
}

func splitGitHubURL(raw string) (owner, repo string, ok bool) {
	u := strings.TrimSpace(raw)
	u = strings.TrimSuffix(u, "/")
	u = strings.TrimSuffix(u, ".git")
	if strings.HasPrefix(u, "git@github.com:") {
		u = "https://github.com/" + strings.TrimPrefix(u, "git@github.com:")
	}
	if !strings.HasPrefix(u, "https://github.com/") {
		return "", "", false
	}
	parts := strings.Split(strings.TrimPrefix(u, "https://github.com/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func gitClone(url, target string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", url, target)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// zipDownload fetches the repo archive and extracts it (git not required).
func zipDownload(owner, repo, target string) error {
	client := &http.Client{Timeout: 20 * time.Minute}
	base := fmt.Sprintf("https://codeload.github.com/%s/%s/zip/refs/heads/", owner, repo)
	for _, branch := range []string{"main", "master"} {
		resp, err := client.Get(base + branch)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}
		tmp, err := os.CreateTemp("", "ittoolkit-dl-*.zip")
		if err != nil {
			resp.Body.Close()
			return err
		}
		_, copyErr := io.Copy(tmp, resp.Body)
		resp.Body.Close()
		tmp.Close()
		if copyErr != nil {
			os.Remove(tmp.Name())
			return copyErr
		}
		unzipErr := unzipStripFirst(tmp.Name(), target)
		os.Remove(tmp.Name())
		return unzipErr
	}
	return fmt.Errorf("could not download %s/%s (tried main/master branches)", owner, repo)
}

func init() {
	for _, r := range RecommendedRepos {
		RegisterURL(r.Name, r.URL)
	}
}

// InstallByName resolves repo names against the internal registry and installs
// them. The frontend never handles URLs — only names.
func InstallByName(names []string, dest string, progress ProgressFunc) []*InstallResult {
	var urls []string
	for _, n := range names {
		u, ok := resolveURL(n)
		if !ok {
			urls = append(urls, n) // will fail gracefully in Install
		} else {
			urls = append(urls, u)
		}
	}
	return Install(urls, dest, progress)
}

func unzipStripFirst(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		name := f.Name
		if idx := strings.IndexByte(name, '/'); idx >= 0 {
			name = name[idx+1:] // strip the top-level archive folder
		}
		if name == "" {
			continue
		}
		// Zip-slip guard: reject any path that escapes dest.
		cleaned := filepath.Clean(filepath.FromSlash(name))
		if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) || filepath.IsAbs(cleaned) {
			return fmt.Errorf("archive entry escapes target directory: %s", name)
		}
		target := filepath.Join(dest, cleaned)
		rel, err := filepath.Rel(dest, target)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return fmt.Errorf("archive entry escapes target directory: %s", name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(target)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
