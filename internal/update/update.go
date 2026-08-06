package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	// GistURL is the raw URL for the Solutions IT Toolkit update manifest.
	GistURL = "https://gist.githubusercontent.com/alexxdavid/603ee193205abab222bee17188995f5a/raw/update_manifest.json"

	// CurrentVersion is the app version baked into the binary at build time.
	CurrentVersion = "236.12"
	CurrentBuild   = 5
)

// Info is the update manifest parsed from the GitHub Gist.
type Info struct {
	Available    bool   `json:"available"`
	Version      string `json:"version"`
	Build        int    `json:"build"`
	DownloadURL  string `json:"installer_url"`
	Notes        string `json:"notes"`
	ForceUpdate  bool   `json:"force_update"`
	SHA256       string `json:"sha256,omitempty"`
	Stale        bool   `json:"stale"`
}

// Manager handles update checking, downloading, and installing.
type Manager struct {
	downloading bool
	percent     int64
	total       int64
	mu          sync.Mutex
	cache       Info
	cacheTime   time.Time
}

// NewManager creates an update manager.
func NewManager() *Manager {
	return &Manager{}
}

func parseVersion(v string) [3]int {
	var out [3]int
	parts := strings.Split(strings.TrimSpace(v), ".")
	re := regexp.MustCompile(`^\d+`)
	for i := 0; i < len(parts) && i < 3; i++ {
		val, _ := strconv.Atoi(re.FindString(parts[i]))
		out[i] = val
	}
	return out
}

// Check fetches the update manifest and compares versions.
func (m *Manager) Check(force bool) (Info, error) {
	var info Info

	if !force {
		m.mu.Lock()
		if !m.cacheTime.IsZero() && time.Since(m.cacheTime) < 5*time.Minute {
			cached := m.cache
			m.mu.Unlock()
			return cached, nil
		}
		m.mu.Unlock()
	}

	cacheBust := time.Now().Unix() / 3600
	if force {
		cacheBust = time.Now().UnixNano()
	}
	url := fmt.Sprintf("%s?t=%d", GistURL, cacheBust)

	client := &http.Client{Timeout: 10 * time.Second}
	var resp *http.Response
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return info, err
		}
		req.Header.Set("User-Agent", "ITToolkit")
		resp, lastErr = client.Do(req)
		if lastErr != nil {
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = fmt.Errorf("rate limited")
			continue
		}
		break
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if lastErr != nil && !m.cacheTime.IsZero() {
		// Return stale info marked as stale so the UI can warn instead of
		// claiming everything is up to date while offline.
		stale := m.cache
		stale.Stale = true
		return stale, fmt.Errorf("could not reach update server: %v", lastErr)
	}
	if lastErr != nil {
		return info, lastErr
	}
	if resp == nil {
		return info, fmt.Errorf("update server unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return info, fmt.Errorf("update server returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return info, err
	}

	remoteV := parseVersion(info.Version)
	localV := parseVersion(CurrentVersion)
	hasUpdate := false
	if remoteV[0] > localV[0] {
		hasUpdate = true
	} else if remoteV[0] == localV[0] && remoteV[1] > localV[1] {
		hasUpdate = true
	} else if remoteV[0] == localV[0] && remoteV[1] == localV[1] && remoteV[2] > localV[2] {
		hasUpdate = true
	} else if remoteV[0] == localV[0] && remoteV[1] == localV[1] && remoteV[2] == localV[2] && info.Build > CurrentBuild {
		hasUpdate = true
	}

	info.Available = hasUpdate && info.DownloadURL != ""
	m.cache = info
	m.cacheTime = time.Now()
	return info, nil
}

// Download fetches the installer to %TEMP% and returns the path.
func (m *Manager) Download(url, version string) (string, error) {
	tmpDir := os.TempDir()
	re := regexp.MustCompile(`[^0-9A-Za-z._-]`)
	safeVer := re.ReplaceAllString(version, "_")
	destPath := filepath.Join(tmpDir, fmt.Sprintf("ITToolkit-Setup-%s.exe", safeVer))

	// Prune old installers
	files, _ := os.ReadDir(tmpDir)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "ITToolkit-Setup-") && strings.HasSuffix(f.Name(), ".exe") {
			os.Remove(filepath.Join(tmpDir, f.Name()))
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "ITToolkit")
	// GitHub API asset URLs return JSON metadata unless we ask for the binary.
	req.Header.Set("Accept", "application/octet-stream")
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	m.mu.Lock()
	m.downloading = true
	m.percent = 0
	m.total = resp.ContentLength
	m.mu.Unlock()

	pw := &progressWriter{total: resp.ContentLength, mgr: m}
	_, err = io.Copy(out, io.TeeReader(resp.Body, pw))

	m.mu.Lock()
	m.downloading = false
	m.mu.Unlock()

	if err != nil {
		os.Remove(destPath)
		return "", err
	}
	out.Close()

	// Integrity check: the installer must be a Windows PE binary.
	if !isPEExecutable(destPath) {
		os.Remove(destPath)
		return "", fmt.Errorf("downloaded file is not a valid installer")
	}
	return destPath, nil
}

// isPEExecutable verifies the MZ/PE signature of a downloaded installer.
func isPEExecutable(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	header := make([]byte, 2)
	if _, err := f.ReadAt(header, 0); err != nil {
		return false
	}
	if header[0] != 'M' || header[1] != 'Z' {
		return false
	}
	// Check PE signature at the offset stored at 0x3C.
	var peOff [4]byte
	if _, err := f.ReadAt(peOff[:], 0x3C); err != nil {
		return false
	}
	off := int(peOff[0]) | int(peOff[1])<<8 | int(peOff[2])<<16 | int(peOff[3])<<24
	peSig := make([]byte, 4)
	if _, err := f.ReadAt(peSig, int64(off)); err != nil {
		return false
	}
	return peSig[0] == 'P' && peSig[1] == 'E' && peSig[2] == 0 && peSig[3] == 0
}

// Progress returns the current download state.
func (m *Manager) Progress() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	var pct float64
	if m.total > 0 {
		pct = float64(m.percent) / float64(m.total) * 100
	}
	return map[string]interface{}{
		"downloading": m.downloading,
		"percent":     pct,
	}
}

// Install launches the installer and exits the app.
func Install(installerPath string) error {
	shell32 := syscall.NewLazyDLL("shell32.dll")
	proc := shell32.NewProc("ShellExecuteW")
	opPtr, _ := syscall.UTF16PtrFromString("open")
	filePtr, _ := syscall.UTF16PtrFromString(installerPath)
	ret, _, err := proc.Call(0, uintptr(unsafe.Pointer(opPtr)), uintptr(unsafe.Pointer(filePtr)), 0, 0, 1)
	if ret <= 32 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return fmt.Errorf("ShellExecuteW failed: %d", ret)
	}
	os.Exit(0)
	return nil
}

// InstallWithArgs launches the installer with extra params and exits.
func InstallWithArgs(installerPath, params string) error {
	shell32 := syscall.NewLazyDLL("shell32.dll")
	proc := shell32.NewProc("ShellExecuteW")
	opPtr, _ := syscall.UTF16PtrFromString("open")
	filePtr, _ := syscall.UTF16PtrFromString(installerPath)
	paramsPtr, _ := syscall.UTF16PtrFromString(params)
	ret, _, err := proc.Call(0, uintptr(unsafe.Pointer(opPtr)), uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(paramsPtr)), 0, 1)
	if ret <= 32 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return fmt.Errorf("ShellExecuteW failed: %d", ret)
	}
	os.Exit(0)
	return nil
}

type progressWriter struct {
	total int64
	mgr   *Manager
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.mgr.mu.Lock()
	pw.mgr.percent += int64(n)
	pw.mgr.total = pw.total
	pw.mgr.mu.Unlock()
	return n, nil
}

// Helper to launch a program detached.
func HiddenCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	return cmd
}
