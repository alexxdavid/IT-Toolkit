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
	Name         string `json:"name"`
	Version      string `json:"version"`
	Category     string `json:"category"`
	Download     string `json:"download"`
	Notes        string `json:"notes"`
	WingetID     string `json:"wingetId"`
	Manufacturer string `json:"manufacturer"`
	Description  string `json:"description"`
	License      string `json:"license"`
	InstallerType string `json:"installerType"`
	Architecture string `json:"architecture"`
	OfficialSite string `json:"officialSite"`
	SilentArgs   string `json:"silentArgs"`
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
	{Name: "Google Chrome", Version: "", Category: "Browsers", Download: "https://dl.google.com/chrome/direct/googlechromestandaloneenterprise64_x64.msi", Notes: "64-bit enterprise installer", WingetID: "Google.Chrome"},
	{Name: "Mozilla Firefox", Version: "", Category: "Browsers", Download: "https://download.mozilla.org/?product=firefox-latest-ssl&os=win64&lang=en-US", Notes: "", WingetID: "Mozilla.Firefox"},
	{Name: "Microsoft Edge", Version: "", Category: "Browsers", Download: "https://www.microsoft.com/en-us/edge/download", Notes: "", WingetID: "Microsoft.Edge"},
	{Name: "Brave Browser", Version: "", Category: "Browsers", Download: "https://laptop-updates.brave.com/latest/standalone", Notes: "", WingetID: "Brave.Brave"},

	// Productivity
	{Name: "Microsoft Teams", Version: "", Category: "Productivity", Download: "https://go.microsoft.com/fwlink/p/?linkid=2112886", Notes: "New Teams client", WingetID: "Microsoft.Teams"},
	{Name: "Zoom", Version: "", Category: "Productivity", Download: "https://zoom.us/client/latest/ZoomInstallerFull.exe", Notes: "x64 on 64-bit systems", WingetID: "Zoom.Zoom"},
	{Name: "Slack", Version: "", Category: "Productivity", Download: "https://cdn.slack-edge.com/bundles/meting-19.8.0/15.8.0.17058/win64/SlackSetup.exe", Notes: "", WingetID: "SlackTechnologies.Slack"},
	{Name: "Notepad++", Version: "", Category: "Productivity", Download: "https://github.com/notepad-plus-plus/notepad-plus-plus/releases/download/v8.8.1/npp.8.8.1.Installer.x64.exe", Notes: "Use winget for always-current: winget install Notepad++.Notepad++", WingetID: "Notepad++.Notepad++"},
	{Name: "LibreOffice", Version: "", Category: "Productivity", Download: "https://download.documentfoundation.org/libreoffice/stable/25.2.4/win/x86_64/LibreOffice_25.2.4_Win_x86-64.msi", Notes: "64-bit MSI", WingetID: "TheDocumentFoundation.LibreOffice"},
	{Name: "VLC", Version: "", Category: "Productivity", Download: "https://get.videolan.org/vlc/3.0.21/win64/vlc-3.0.21-win64.exe", Notes: "Media player", WingetID: "VideoLAN.VLC"},
	{Name: "Adobe Acrobat Reader", Version: "", Category: "Productivity", Download: "https://ardownload2.adobe.com/pub/adobe/acrobat/win/AcrobatDC/2500720030/AcroRdrDCx64_2500720030_MUI.exe", Notes: "PDF reader", WingetID: "Adobe.Acrobat.Reader.64-bit"},

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

	// Networking
	{Name: "Nmap", Version: "", Category: "Networking", Download: "https://nmap.org/dist/nmap-7.97-setup.exe", Notes: "Network scanner", WingetID: "Insecure.Nmap"},

	// Security extras
	{Name: "1Password", Version: "", Category: "Security", Download: "https://downloads.1password.com/win/1PasswordSetup-x64.msi", Notes: "", WingetID: "AgileBits.1Password"},
	{Name: "OpenVPN", Version: "", Category: "Security", Download: "https://swupdate.openvpn.org/community/releases/OpenVPN-2.6.12-I001-amd64.msi", Notes: "", WingetID: "OpenVPNTechnologies.OpenVPN"},

	// Media & Office
	{Name: "SumatraPDF", Version: "", Category: "Utilities", Download: "https://github.com/sumatrapdfreader/sumatrapdf/releases/latest/download/SumatraPDF-3.5.2-64-install.exe", Notes: "Fast PDF viewer", WingetID: "SumatraPDF.SumatraPDF"},
	{Name: "Paint.NET", Version: "", Category: "Utilities", Download: "https://www.dotpdn.com/files/paint.net.5.0.14.install.x64.zip", Notes: "Image editor", WingetID: "dotpdn.paintdotnet"},
	{Name: "IrfanView", Version: "", Category: "Utilities", Download: "https://www.irfanview.com/iview470_x64_setup.exe", Notes: "Image viewer", WingetID: "IrfanView.IrfanView"},

	// Developer tools
	{Name: "Node.js LTS", Version: "", Category: "Development", Download: "https://nodejs.org/dist/v22.16.0/node-v22.16.0-x64.msi", Notes: "LTS version", WingetID: "OpenJS.NodeJS.LTS"},
	{Name: "Docker Desktop", Version: "", Category: "Development", Download: "https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe", Notes: "", WingetID: "Docker.DockerDesktop"},
	{Name: "GitHub Desktop", Version: "", Category: "Development", Download: "https://central.github.com/deployments/desktop/desktop/latest/win32", Notes: "", WingetID: "GitHub.GitHubDesktop"},
	{Name: "Postman", Version: "", Category: "Development", Download: "https://dl.pstmn.io/download/latest/win64", Notes: "API testing", WingetID: "Postman.Postman"},
	{Name: "WinMerge", Version: "", Category: "Development", Download: "https://github.com/WinMerge/winmerge/releases/latest/download/WinMerge-2.16.46-x64-Setup.exe", Notes: "File diff tool", WingetID: "WinMerge.WinMerge"},

	// RMM & Admin tools
	{Name: "DBeaver", Version: "", Category: "Utilities", Download: "https://dbeaver.io/files/dbeaver-ce-latest-win32.win32.x86_64.zip", Notes: "Database GUI", WingetID: "dbeaver.dbeaver"},
	{Name: "FileZilla", Version: "", Category: "Networking", Download: "https://download.filezilla-project.org/client/FileZilla_3.68.1_win64-setup.exe", Notes: "FTP client", WingetID: "FileZillaProject.FileZilla.Client"},
	{Name: "Sysinternals Suite", Version: "", Category: "Utilities", Download: "https://download.sysinternals.com/files/SysinternalsSuite.zip", Notes: "Essential IT admin tools", WingetID: "Microsoft.Sysinternals"},
	{Name: "WizTree", Version: "", Category: "Utilities", Download: "https://diskanalyzer.com/files/wiztree_4_25_setup.exe", Notes: "Disk space analyzer", WingetID: "AntibodySoftware.WizTree"},
	{Name: "CrystalDiskInfo", Version: "", Category: "Utilities", Download: "https://crystalmark.info/download/CrystalDiskInfo8_17_14.zip", Notes: "Disk health monitor", WingetID: "CrystalDewWorld.CrystalDiskInfo"},

	// Collaboration
	{Name: "Discord", Version: "", Category: "Productivity", Download: "https://discord.com/api/downloads/distributions/app/installers/win32/x86/DiscordSetup.exe", Notes: "", WingetID: "Discord.Discord"},
	{Name: "Telegram", Version: "", Category: "Productivity", Download: "https://telegram.org/dl/desktop/win64", Notes: "x64", WingetID: "Telegram.TelegramDesktop"},

	// Entertainment
	{Name: "Steam", Version: "", Category: "Utilities", Download: "https://cdn.cloudflare.steamstatic.com/client/installer/SteamSetup.exe", Notes: "Game store", WingetID: "Valve.Steam"},
	{Name: "OBS Studio", Version: "", Category: "Utilities", Download: "https://cdn-fastly.obsproject.com/downloads/OBS-Studio-31.1.2-Windows-Installer.exe", Notes: "Screen recording", WingetID: "OBSProject.OBSStudio"},

	// Productivity extras
	{Name: "AutoHotkey", Version: "", Category: "Development", Download: "https://www.autohotkey.com/download/ahk-install.exe", Notes: "Scripting for Windows automation", WingetID: "AutoHotkey.AutoHotkey"},
	{Name: "PowerToys", Version: "", Category: "Utilities", Download: "https://github.com/microsoft/PowerToys/releases/latest/download/PowerToysUserSetup-x64.exe", Notes: "Windows power user tools", WingetID: "Microsoft.PowerToys"},
	{Name: "O&O ShutUp10", Version: "", Category: "Utilities", Download: "https://dl5.oo-software.com/files/ooshutup10/OOSU10.exe", Notes: "Privacy tool", WingetID: "OO-Software.OOShutUp10"},
	{Name: "Revo Uninstaller", Version: "", Category: "Utilities", Download: "https://download.revouninstaller.com/download/revouninstallerpro-setup.exe", Notes: "Uninstaller", WingetID: "RevoUninstaller.RevoUninstaller"},

	// Backup extras
	{Name: "Clonezilla", Version: "", Category: "Backup", Download: "https://sourceforge.net/projects/clonezilla/files/clonezilla_live_stable/", Notes: "Disk cloning (USB boot)", WingetID: ""},
	{Name: "Macrium Reflect Free", Version: "", Category: "Backup", Download: "https://www.macrium.com/download/reflectionk", Notes: "Disk imaging", WingetID: "ParamountSoftwareUK.MacriumReflectFree"},
	{Name: "Duplicati", Version: "", Category: "Backup", Download: "https://github.com/duplicati/duplicati/releases/latest/download/duplicati-2.1.1.2_x64.msi", Notes: "Cloud backup", WingetID: "Duplicati.Duplicati"},
}

// GetSoftwareCatalog returns the curated software list.
func GetSoftwareCatalog() []*SoftwareItem {
	out := make([]*SoftwareItem, len(SoftwareCatalog))
	for i, item := range SoftwareCatalog {
		cp := *item
		if m, ok := softwareMeta[item.Name]; ok {
			cp.Manufacturer = m.Manufacturer
			cp.Description = m.Description
			cp.License = m.License
			cp.InstallerType = m.InstallerType
			cp.Architecture = m.Architecture
			cp.OfficialSite = m.OfficialSite
			cp.SilentArgs = m.SilentArgs
		}
		out[i] = &cp
	}
	return out
}

// softwareMeta holds rich metadata per software name, applied on top of the
// base catalog so entries stay readable.
var softwareMeta = map[string]struct {
	Manufacturer  string
	Description   string
	License       string
	InstallerType string
	Architecture  string
	OfficialSite  string
	SilentArgs    string
}{
	"Google Chrome":    {"Google", "Fast, secure web browser", "Proprietary/Free", "MSI", "x64", "https://www.google.com/chrome", "msiexec /i chrome.msi /quiet"},
	"Mozilla Firefox":  {"Mozilla", "Open-source web browser", "Open Source (MPL)", "EXE", "x64", "https://www.mozilla.org/firefox", "-ms"},
	"Microsoft Edge":   {"Microsoft", "Chromium-based web browser", "Proprietary/Free", "EXE", "x64", "https://www.microsoft.com/edge", "--silent-install"},
	"Brave Browser":    {"Brave", "Privacy-focused web browser", "Open Source (MPL-2.0)", "EXE", "x64", "https://brave.com", "--silent"},
	"Microsoft Teams":  {"Microsoft", "Team collaboration and meetings", "Proprietary/Free", "EXE", "x64", "https://www.microsoft.com/microsoft-teams", "--silent"},
	"Zoom":             {"Zoom Video", "Video conferencing and meetings", "Proprietary/Free", "EXE", "x64", "https://zoom.us", "/silent /norestart"},
	"Slack":            {"Slack", "Team messaging and collaboration", "Proprietary/Free", "EXE", "x64", "https://slack.com", "--silent"},
	"Notepad++":        {"Notepad++", "Advanced text and code editor", "Open Source (GPL-3.0)", "EXE", "x64", "https://notepad-plus-plus.org", "/S"},
	"LibreOffice":      {"The Document Foundation", "Open-source office suite", "Open Source (MPL-2.0)", "MSI", "x64", "https://www.libreoffice.org", "msiexec /i LibreOffice.msi /qn"},
	"VLC":              {"VideoLAN", "Open-source media player", "Open Source (GPL-2.0)", "EXE", "x64", "https://www.videolan.org/vlc", "/L=1033 /S"},
	"Adobe Acrobat Reader": {"Adobe", "PDF viewer and editor", "Proprietary/Free", "EXE", "x64", "https://get.adobe.com/reader", "/sPB"},
	"Visual Studio Code": {"Microsoft", "Cross-platform code editor", "Open Source (MIT)", "EXE", "x64", "https://code.visualstudio.com", "/VERYSILENT /NORESTART"},
	"Git for Windows":  {"The Git Project", "Distributed version control", "Open Source (GPL-2.0)", "EXE", "x64", "https://git-scm.com", "/VERYSILENT /NORESTART"},
	"PowerShell 7":     {"Microsoft", "Cross-platform task automation shell", "Open Source (MIT)", "MSI", "x64", "https://github.com/PowerShell/PowerShell", "msiexec /i PowerShell.msi /qn ADD_EXPLORER_CONTEXT_MENU_OPENPOWERSHELL=1"},
	"Windows Terminal": {"Microsoft", "Modern terminal for command-line tools", "Open Source (MIT)", "MSIX", "x64", "https://github.com/microsoft/terminal", "Add-AppxPackage"},
	"Python 3":         {"Python Software Foundation", "High-level programming language", "Open Source (PSF)", "EXE", "x64", "https://www.python.org", "/quiet InstallAllUsers=1 PrependPath=1"},
	"VeraCrypt":        {"IDRIX", "Free disk encryption software", "Open Source (Apache-2.0)", "EXE", "x64", "https://www.veracrypt.fr", "/S"},
	"KeePass":          {"KeePass", "Open-source password manager", "Open Source (GPL-2.0)", "EXE", "x64", "https://keepass.info", "/SILENT"},
	"Bitwarden":        {"Bitwarden", "Open-source password manager", "Open Source (GPL-3.0)", "EXE", "x64", "https://bitwarden.com", "--silent"},
	"1Password":        {"1Password", "Password manager for teams", "Proprietary", "MSI", "x64", "https://1password.com", "msiexec /i 1Password.msi /qn"},
	"OpenVPN":          {"OpenVPN", "Open-source VPN client", "Open Source (GPL-2.0)", "MSI", "x64", "https://openvpn.net", "msiexec /i OpenVPN.msi /qn"},
	"7-Zip":            {"Igor Pavlov", "Open-source file archiver", "Open Source (LGPL)", "EXE", "x64", "https://www.7-zip.org", "/S"},
	"PuTTY":            {"Simon Tatham", "SSH and telnet client", "Open Source (MIT)", "EXE", "x64", "https://www.putty.org", "/VERYSILENT"},
	"WinSCP":           {"WinSCP", "SFTP, SCP and FTP client", "Open Source (GPL-3.0)", "EXE", "x64", "https://winscp.net", "/VERYSILENT"},
	"Everything":       {"voidtools", "Instant file search utility", "Freeware", "EXE", "x64", "https://www.voidtools.com", "/S"},
	"RustDesk":         {"RustDesk", "Open-source remote desktop", "Open Source (AGPL-3.0)", "EXE", "x64", "https://rustdesk.com", "--silent-install"},
	"AnyDesk":          {"AnyDesk", "Remote desktop software", "Proprietary/Free", "EXE", "x64", "https://anydesk.com", "--silent"},
	"TeamViewer":       {"TeamViewer", "Remote access and support", "Proprietary/Free", "EXE", "x64", "https://teamviewer.com", "/S"},
	"Tailscale":        {"Tailscale", "Zero-config VPN based on WireGuard", "Proprietary/Free", "EXE", "x64", "https://tailscale.com", "/quiet"},
	"Wireshark":        {"Wireshark", "Network protocol analyzer", "Open Source (GPL-2.0)", "EXE", "x64", "https://www.wireshark.org", "/S"},
	"Advanced IP Scanner": {"Famatech", "Fast LAN network scanner", "Freeware", "EXE", "x64", "https://www.advanced-ip-scanner.com", "/VERYSILENT"},
	"Nmap":             {"Nmap", "Network discovery and security scanner", "Open Source (NPSL)", "EXE", "x64", "https://nmap.org", "/S"},
	"Veeam Agent":      {"Veeam", "Free backup agent for Windows", "Proprietary/Free", "EXE", "x64", "https://www.veeam.com", "/S"},
	"Google Drive":     {"Google", "Cloud file sync and storage", "Proprietary/Free", "EXE", "x64", "https://www.google.com/drive", "--silent"},
	"SumatraPDF":       {"SumatraPDF", "Lightweight PDF viewer", "Open Source (GPL-3.0)", "EXE", "x64", "https://www.sumatrapdfreader.org", "-s"},
	"Paint.NET":        {"dotPDN LLC", "Free image and photo editor", "Proprietary/Free", "ZIP", "x64", "https://www.getpaint.net", "n/a (portable ZIP)"},
	"IrfanView":        {"Irfan Skiljan", "Fast image viewer and converter", "Proprietary/Free", "EXE", "x64", "https://www.irfanview.com", "/silent"},
	"Node.js LTS":      {"OpenJS Foundation", "JavaScript runtime built on Chrome's V8", "Open Source (MIT)", "MSI", "x64", "https://nodejs.org", "msiexec /i node.msi /qn"},
	"Docker Desktop":   {"Docker", "Containerized app development", "Proprietary/Free", "EXE", "x64", "https://www.docker.com", "install --accept-license"},
	"GitHub Desktop":   {"GitHub", "Visual Git client", "Open Source (MIT)", "EXE", "x64", "https://desktop.github.com", "/VERYSILENT"},
	"Postman":          {"Postman", "API development and testing platform", "Proprietary/Free", "EXE", "x64", "https://www.postman.com", "/S"},
	"WinMerge":         {"WinMerge", "File and folder comparison tool", "Open Source (GPL-2.0)", "EXE", "x64", "https://winmerge.org", "/VERYSILENT"},
	"DBeaver":          {"DBeaver", "Universal database client", "Open Source (Apache-2.0)", "ZIP", "x64", "https://dbeaver.io", "n/a (portable ZIP)"},
	"FileZilla":        {"FileZilla Project", "Free FTP/SFTP client", "Open Source (GPL-2.0)", "EXE", "x64", "https://filezilla-project.org", "/S"},
	"Sysinternals Suite": {"Microsoft", "Advanced system utilities and troubleshooting", "Freeware", "ZIP", "x64", "https://learn.microsoft.com/sysinternals", "n/a (portable ZIP)"},
	"WizTree":          {"Antibody Software", "Fast disk space analyzer", "Freeware", "EXE", "x64", "https://wiztreefree.com", "/VERYSILENT"},
	"CrystalDiskInfo":  {"Crystal Dew World", "HDD/SSD health monitor", "Open Source (MIT)", "ZIP", "x64", "https://crystalmark.info", "n/a (portable ZIP)"},
	"Discord":          {"Discord", "Voice, video and text chat", "Proprietary/Free", "EXE", "x64", "https://discord.com", "/s"},
	"Telegram":         {"Telegram", "Cloud-based messaging", "Open Source (GPL-3.0)", "EXE", "x64", "https://telegram.org", "/VERYSILENT"},
	"Steam":            {"Valve", "Digital game distribution platform", "Proprietary/Free", "EXE", "x64", "https://store.steampowered.com", "/S"},
	"OBS Studio":       {"OBS Project", "Open-source screen recording and streaming", "Open Source (GPL-2.0)", "EXE", "x64", "https://obsproject.com", "/S"},
	"AutoHotkey":       {"AutoHotkey", "Windows automation scripting", "Open Source (GPL-2.0)", "EXE", "x64", "https://www.autohotkey.com", "/S"},
	"PowerToys":        {"Microsoft", "Windows power user utilities", "Open Source (MIT)", "EXE", "x64", "https://github.com/microsoft/PowerToys", "--silent"},
	"O&O ShutUp10":     {"O&O Software", "Windows 10/11 privacy and telemetry tool", "Freeware", "EXE", "x64", "https://www.oo-software.com", "/S"},
	"Revo Uninstaller": {"VS Revo Group", "Advanced uninstaller", "Proprietary/Free", "EXE", "x64", "https://www.revouninstaller.com", "/S"},
	"Macrium Reflect Free": {"Paramount Software", "Disk imaging and cloning", "Proprietary/Free", "EXE", "x64", "https://www.macrium.com", "/S"},
	"Duplicati":        {"Duplicati", "Open-source backup with cloud support", "Open Source (LGPL-2.1)", "MSI", "x64", "https://duplicati.com", "msiexec /i duplicati.msi /qn"},
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
			HasDownload:     downloaded[strings.ToLower(item.Name)],
		}
		if avail != "" && isVersion(avail) {
			sv.LatestVersion = avail
			sv.UpdateAvailable = (inst == "" || avail != inst)
		}
		result[item.Name] = sv
	}

	versionMu.Lock()
	versionCache = result
	versionCacheTime = time.Now()
	versionMu.Unlock()
	return result
}

// isVersion returns true if the string looks like a version number.
func isVersion(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || s == "winget" || s == "Version" {
		return false
	}
	return strings.ContainsAny(s, "0123456789")
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
