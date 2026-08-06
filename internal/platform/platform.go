package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const webview2Key = `SOFTWARE\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`
const webview2KeyWoW64 = `SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`

// WebView2Installed checks the registry for the WebView2 runtime.
func WebView2Installed() bool {
	for _, root := range []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER} {
		for _, keyPath := range []string{webview2Key, webview2KeyWoW64} {
			k, err := registry.OpenKey(root, keyPath, registry.QUERY_VALUE)
			if err != nil {
				continue
			}
			pv, _, err := k.GetStringValue("pv")
			k.Close()
			if err == nil && pv != "" {
				return true
			}
		}
	}
	return false
}

// ExeDir returns the directory containing the running executable.
func ExeDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

// DefaultScriptsDir returns <exe dir>\Scripts when it exists.
func DefaultScriptsDir() string {
	return filepath.Join(ExeDir(), "Scripts")
}

// HomeDir returns the user home directory.
func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

// FileExists reports whether a path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ShowMessageBox displays a native modal message box.
func ShowMessageBox(title, text string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	msgbox := user32.NewProc("MessageBoxW")
	tp, _ := syscall.UTF16PtrFromString(title)
	txt, _ := syscall.UTF16PtrFromString(text)
	const mbIconExclamation = 0x00000030
	const mbSystemModal = 0x00001000
	msgbox.Call(0, uintptr(unsafe.Pointer(txt)), uintptr(unsafe.Pointer(tp)), mbIconExclamation|mbSystemModal)
}

// OpenURL opens a URL in the default browser.
func OpenURL(url string) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open url: %w", err)
	}
	return nil
}
