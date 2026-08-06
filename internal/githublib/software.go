package githublib

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// SoftwareItem is a downloadable IT tool (not a GitHub repo).
type SoftwareItem struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Category string `json:"category"`
	Download string `json:"download"`
	Notes    string `json:"notes"`
	WingetID string `json:"wingetId"`
}

var (
	swMu          sync.Mutex
	swDownloading = map[string]bool{}
	swPercent     = map[string]float64{}
	swCurrent     string
)

// SoftwareCatalog is the curated list of IT software shown in the Software tab.
// Download URLs are vendor "latest" endpoints or GitHub /releases/latest which
// resolve through redirects; the downloader captures the final filename.
var SoftwareCatalog = []*SoftwareItem{
	// Browsers
	{Name: "Google Chrome", Version: "", Category: "Browsers", Download: "https://dl.google.com/chrome/install/latest/chrome_installer.exe", Notes: "", WingetID: "Google.Chrome"},
	{Name: "Mozilla Firefox", Version: "", Category: "Browsers", Download: "https://download.mozilla.org/?product=firefox-latest-ssl&os=win64&lang=en-US", Notes: "", WingetID: "Mozilla.Firefox"},
	{Name: "Microsoft Edge", Version: "", Category: "Browsers", Download: "https://www.microsoft.com/en-us/edge/download", Notes: "", WingetID: "Microsoft.Edge"},
	{Name: "Brave Browser", Version: "", Category: "Browsers", Download: "https://laptop-updates.brave.com/latest/standalone", Notes: "", WingetID: "Brave.Brave"},

	// Productivity
	{Name: "Microsoft Teams", Version: "", Category: "Productivity", Download: "https://go.microsoft.com/fwlink/p/?linkid=2112886", Notes: "New Teams client", WingetID: "Microsoft.Teams"},
	{Name: "Zoom", Version: "", Category: "Productivity", Download: "https://zoom.us/client/latest/ZoomInstallerFull.exe", Notes: "", WingetID: "Zoom.Zoom"},
	{Name: "Slack", Version: "", Category: "Productivity", Download: "https://downloads.slack-edge.com/slack-installers/4.41.105/SlackSetup-x64.exe", Notes: "", WingetID: "SlackTechnologies.Slack"},
	{Name: "Notepad++", Version: "", Category: "Productivity", Download: "https://github.com/notepad-plus-plus/notepad-plus-plus/releases/latest/download/npp.installer.x64.exe", Notes: "", WingetID: "Notepad++.Notepad++"},

	// Development
	{Name: "Visual Studio Code", Version: "", Category: "Development", Download: "https://update.code.visualstudio.com/latest/win32-x64-user/stable", Notes: "", WingetID: "Microsoft.VisualStudioCode"},
	{Name: "Git for Windows", Version: "", Category: "Development", Download: "https://github.com/git-for-windows/git/releases/latest/download/Git-2.48.1-64-bit.exe", Notes: "Use winget for always-current: winget install Git.Git", WingetID: "Git.Git"},
	{Name: "PowerShell 7", Version: "", Category: "Development", Download: "https://github.com/PowerShell/PowerShell/releases/latest/download/PowerShell-7.5.0-win-x64.msi", Notes: "Use winget for always-current", WingetID: "Microsoft.PowerShell"},
	{Name: "Windows Terminal", Version: "", Category: "Development", Download: "https://github.com/microsoft/terminal/releases/latest/download/Microsoft.WindowsTerminal_1.21.3231.0_x64_8wekyb3d8bbwe.msixbundle", Notes: "Use winget for always-current", WingetID: "Microsoft.WindowsTerminal"},
	{Name: "Python 3", Version: "", Category: "Development", Download: "https://www.python.org/ftp/python/3.13.0/python-3.13.0-amd64.exe", Notes: "", WingetID: "Python.Python.3.13"},

	// Security
	{Name: "VeraCrypt", Version: "", Category: "Security", Download: "https://launchpad.net/veracrypt/trunk/1.26.15/+download/VeraCrypt_Setup_x64_1.26.15.exe", Notes: "Full disk encryption", WingetID: "IDRIX.VeraCrypt"},
	{Name: "KeePass", Version: "", Category: "Security", Download: "https://github.com/dlech/KeePass2.x/releases/latest/download/KeePass-2.58-Setup.exe", Notes: "Password manager", WingetID: "KeePassXTeam.KeePass"},
	{Name: "Bitwarden", Version: "", Category: "Security", Download: "https://github.com/bitwarden/clients/releases/latest/download/Bitwarden-Installer.exe", Notes: "", WingetID: "Bitwarden.Bitwarden"},

	// Utilities
	{Name: "7-Zip", Version: "", Category: "Utilities", Download: "https://www.7-zip.org/a/7z2501-x64.exe", Notes: "", WingetID: "7zip.7zip"},
	{Name: "PuTTY", Version: "", Category: "Utilities", Download: "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe", Notes: "", WingetID: "PuTTY.PuTTY"},
	{Name: "WinSCP", Version: "", Category: "Utilities", Download: "https://winscp.net/eng/downloads/winscp6005setup.exe", Notes: "SFTP/SCP client", WingetID: "WinSCP.WinSCP"},
	{Name: "Everything", Version: "", Category: "Utilities", Download: "https://www.voidtools.com/Everything-1.4.1.1026.x64-Setup.exe", Notes: "Instant file search", WingetID: "voidtools.Everything"},

	// Remote Access
	{Name: "RustDesk", Version: "", Category: "Remote Access", Download: "https://github.com/rustdesk/rustdesk/releases/latest/download/rustdesk-1.4.1-x86_64.exe", Notes: "Open source remote desktop", WingetID: "RustDesk.RustDesk"},
	{Name: "AnyDesk", Version: "", Category: "Remote Access", Download: "https://download.anydesk.com/AnyDesk.exe", Notes: "", WingetID: "AnyDesk.AnyDesk"},
	{Name: "TeamViewer", Version: "", Category: "Remote Access", Download: "https://download.teamviewer.com/download/TeamViewer_Setup_x64.exe", Notes: "", WingetID: "TeamViewer.TeamViewer"},
	{Name: "Tailscale", Version: "", Category: "Remote Access", Download: "https://pkgs.tailscale.com/stable/tailscale-setup-latest-amd64.exe", Notes: "Zero-config VPN", WingetID: "Tailscale.Tailscale"},

	// Networking
	{Name: "Wireshark", Version: "", Category: "Networking", Download: "https://www.wireshark.org/download/win64/Wireshark-win64-latest.exe", Notes: "", WingetID: "WiresharkFoundation.Wireshark"},
	{Name: "Advanced IP Scanner", Version: "", Category: "Networking", Download: "https://download.advanced-ip-scanner.com/advanced_ip_scanner.exe", Notes: "", WingetID: "Famatech.AdvancedIPScanner"},

	// Backup
	{Name: "Veeam Agent", Version: "", Category: "Backup", Download: "https://download4.veeam.com/VAW_agent_en_setup_x64.exe", Notes: "Free backup agent", WingetID: "Veeam.Agent.Windows"},
	{Name: "Google Drive", Version: "", Category: "Backup", Download: "https://dl.google.com/drive-file-stream/GoogleDriveSetup.exe", Notes: "File Stream", WingetID: "Google.Drive"},
}

// GetSoftwareCatalog returns the curated software list.
func GetSoftwareCatalog() []*SoftwareItem {
	return SoftwareCatalog
}

// SoftwareVersion holds the latest + installed version for a software item.
type SoftwareVersion struct {
	Name            string `json:"name"`
	LatestVersion   string `json:"latestVersion"`
	InstalledVersion string `json:"installedVersion"`
	HasDownload     bool   `json:"hasDownload"`
	UpdateAvailable bool   `json:"updateAvailable"`
}

var (
	versionCache     map[string]SoftwareVersion
	versionCacheTime time.Time
	versionMu        sync.Mutex
)

// GetSoftwareVersions queries winget for latest versions, checks installed
// status, and scans the download folder. Results cached for 1 hour.
func GetSoftwareVersions(items []*SoftwareItem, softwareDir string) map[string]SoftwareVersion {
	versionMu.Lock()
	if versionCache != nil && time.Since(versionCacheTime) < time.Hour {
		cached := versionCache
		versionMu.Unlock()
		return cached
	}
	versionMu.Unlock()

	// Query winget once (with timeout, async-safe).
	type wingetResult struct {
		installed map[string]string // key: winget Id
		available map[string]string // key: winget Id (newer version, if any)
	}
	wingetCh := make(chan wingetResult, 1)
	go func() {
		installed, available := queryWingetState()
		wingetCh <- wingetResult{installed, available}
	}()

	downloaded := scanDownloadFolder(softwareDir)

	var winget wingetResult
	select {
	case w := <-wingetCh:
		winget = w
	case <-time.After(35 * time.Second):
		winget = wingetResult{map[string]string{}, map[string]string{}}
	}

	result := map[string]SoftwareVersion{}
	for _, item := range items {
		id := item.WingetID
		inst := winget.installed[id]
		avail := winget.available[id]
		sv := SoftwareVersion{
			Name:            item.Name,
			InstalledVersion: inst,
			LatestVersion:   inst,
			HasDownload:     downloaded[strings.ToLower(item.Name)],
		}
		if avail != "" && avail != inst {
			sv.LatestVersion = avail
			sv.UpdateAvailable = true
		}
		result[item.Name] = sv
	}

	versionMu.Lock()
	versionCache = result
	versionCacheTime = time.Now()
	versionMu.Unlock()
	return result
}

// InvalidateVersionCache forces the next call to re-query winget.
func InvalidateVersionCache() {
	versionMu.Lock()
	versionCache = nil
	versionCacheTime = time.Time{}
	versionMu.Unlock()
}

// queryWingetState runs `winget list` once and parses the table into
// Id -> installed version and Id -> available (upgrade) version.
func queryWingetState() (installed map[string]string, available map[string]string) {
	installed = map[string]string{}
	available = map[string]string{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "winget", "list", "--accept-source-agreements", "--disable-interactivity")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // hidden window
	data, err := cmd.Output()
	if err != nil {
		return installed, available
	}

	lines := strings.Split(string(data), "\n")

	// Locate header row and its separator row.
	headerIdx := -1
	sepIdx := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Name") && strings.Contains(trimmed, "Id") && strings.Contains(trimmed, "Version") {
			headerIdx = i
			// separator is the next row
			for j := i + 1; j < len(lines); j++ {
				if strings.Contains(lines[j], "--") {
					sepIdx = j
					break
				}
			}
			break
		}
	}
	if headerIdx < 0 || sepIdx < 0 {
		return installed, available
	}

	header := lines[headerIdx]
	nameStart := strings.Index(header, "Name")
	idStart := strings.Index(header, "Id")
	verStart := strings.Index(header, "Version")
	availStart := strings.Index(header, "Available")
	if nameStart < 0 || idStart < 0 || verStart < 0 {
		return installed, available
	}
	if availStart < 0 {
		availStart = len(header) // available column may not exist in some versions
	}

	col := func(line string, start, end int) string {
		if start < 0 || start >= len(line) {
			return ""
		}
		if end < 0 || end > len(line) {
			end = len(line)
		}
		return strings.TrimSpace(line[start:end])
	}

	for i := sepIdx + 1; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.Contains(trimmed, "--") {
			continue
		}
		id := col(line, idStart, verStart)
		if id == "" {
			continue
		}
		version := col(line, verStart, availStart)
		avail := ""
		if availStart < len(header) {
			avail = col(line, availStart, len(line))
		}
		if version != "" {
			installed[id] = version
		}
		if avail != "" {
			available[id] = avail
		}
	}
	return installed, available
}

// scanDownloadFolder checks the software directory for matching installer files.
func scanDownloadFolder(dir string) map[string]bool {
	result := map[string]bool{}
	if dir == "" {
		return result
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return result
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.ToLower(e.Name())
		// Skip partial/cache downloads
		if strings.HasSuffix(name, ".part") || strings.HasSuffix(name, ".crdownload") || strings.HasSuffix(name, ".tmp") {
			continue
		}
		for _, item := range SoftwareCatalog {
			keywords := strings.Fields(strings.ToLower(item.Name))
			for _, kw := range keywords {
				if len(kw) > 3 && strings.Contains(name, kw) {
					result[strings.ToLower(item.Name)] = true
					break
				}
			}
		}
	}
	return result
}

// DownloadSoftware downloads a file to dest with progress tracking.
func DownloadSoftware(name, url, dest string) (string, error) {
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return "", err
	}
	safeName := strings.ReplaceAll(name, "/", "-")
	safeName = strings.ReplaceAll(safeName, " ", "_")

	swMu.Lock()
	if swCurrent != "" {
		swMu.Unlock()
		return "", fmt.Errorf("another download is already in progress")
	}
	swCurrent = name
	swDownloading[name] = true
	swPercent[name] = 0
	swMu.Unlock()

	defer func() {
		swMu.Lock()
		swDownloading[name] = false
		if swCurrent == name {
			swCurrent = ""
		}
		swMu.Unlock()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) SolutionsITToolkit/1.0")

	client := &http.Client{CheckRedirect: func(r *http.Request, via []*http.Request) error {
		r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) SolutionsITToolkit/1.0")
		return nil
	}}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("URL returned a web page, not an installer — check the link")
	}

	// Determine filename: Content-Disposition → final URL path → safe name.
	fileName := filenameFromContentDisposition(resp.Header.Get("Content-Disposition"))
	if fileName == "" {
		if resp.Request != nil && resp.Request.URL != nil {
			p := resp.Request.URL.Path
			base := filepath.Base(p)
			if base != "" && base != "." && base != "/" && !strings.Contains(base, "?") && strings.Contains(base, ".") {
				fileName = base
			}
		}
	}
	if fileName == "" || fileName == "." || fileName == "/" {
		fileName = safeName
	}
	if !strings.Contains(fileName, ".") {
		switch {
		case strings.Contains(contentType, "msi"):
			fileName += ".msi"
		case strings.Contains(contentType, "zip"):
			fileName += ".zip"
		case strings.Contains(contentType, "msix"):
			fileName += ".msixbundle"
		default:
			fileName += ".exe"
		}
	}

	// Download to .part then rename atomically.
	destPath := filepath.Join(dest, fileName)
	partPath := destPath + ".part"
	out, err := os.Create(partPath)
	if err != nil {
		return "", err
	}

	total := resp.ContentLength
	buf := make([]byte, 64*1024)
	var downloaded int64
	writeErr := error(nil)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, wErr := out.Write(buf[:n]); wErr != nil {
				writeErr = wErr
				break
			}
			downloaded += int64(n)
			swMu.Lock()
			if total > 0 {
				swPercent[name] = float64(downloaded) / float64(total) * 100
			}
			swMu.Unlock()
		}
		if readErr != nil {
			if readErr != io.EOF {
				writeErr = readErr
			}
			break
		}
	}
	out.Close()

	if writeErr != nil {
		os.Remove(partPath)
		return "", writeErr
	}
	if err := os.Rename(partPath, destPath); err != nil {
		os.Remove(partPath)
		return "", err
	}
	return destPath, nil
}

// filenameFromContentDisposition parses filename= from a Content-Disposition header.
func filenameFromContentDisposition(cd string) string {
	if cd == "" {
		return ""
	}
	for _, part := range strings.Split(cd, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToLower(part), "filename=") {
			v := strings.TrimPrefix(part, "filename=")
			v = strings.TrimPrefix(v, "filename*=UTF-8''")
			v = strings.Trim(v, `"`)
			if idx := strings.Index(v, "?"); idx >= 0 {
				v = v[:idx]
			}
			return v
		}
	}
	return ""
}

// SoftwareProgress returns current download state.
func SoftwareProgress() map[string]interface{} {
	swMu.Lock()
	defer swMu.Unlock()
	name := swCurrent
	pct := 0.0
	if name != "" {
		pct = swPercent[name]
	}
	return map[string]interface{}{
		"downloading": name != "",
		"name":        name,
		"percent":     pct,
	}
}
