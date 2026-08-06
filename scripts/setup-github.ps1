# IT Toolkit — GitHub Setup Script
# Run this to create the GitHub repo and update gist for auto-updates.
# Prerequisites: `gh auth login` first.

$ErrorActionPreference = 'Stop'

Write-Host "=== IT Toolkit GitHub Setup ===" -ForegroundColor Cyan

# 1. Create the GitHub repo
Write-Host "`n1. Creating GitHub repo..." -ForegroundColor Yellow
gh repo create IT-Toolkit --private --description "IT script catalog and tooling console" --source . --push 2>&1 | Write-Host

# 2. Create the update gist
Write-Host "`n2. Creating update manifest gist..." -ForegroundColor Yellow
$manifest = @{
    version = "1.0.0"
    build = 1
    installer_url = ""
    notes = "Initial release"
    force_update = $false
    previous_versions = @()
} | ConvertTo-Json -Depth 5

$gistId = ($manifest | gh gist create --desc "IT Toolkit Update Manifest" -f "update_manifest.json" --public 2>&1)
Write-Host "Gist created: $gistId" -ForegroundColor Green

# 3. Build the installer for upload
Write-Host "`n3. Building installer..." -ForegroundColor Yellow
.\build.ps1

# 4. Get the raw URL
$rawUrl = "https://gist.githubusercontent.com/$((gh api user).login)/$gistId/raw/update_manifest.json"
Write-Host "`n=== SETUP COMPLETE ===" -ForegroundColor Green
Write-Host "1. Repo: https://github.com/$((gh api user).login)/IT-Toolkit"
Write-Host "2. Gist raw URL: $rawUrl"
Write-Host "`nUpdate the gist URL in internal/update/update.go:"
Write-Host "  GistURL = `"$rawUrl`"" -ForegroundColor Yellow
Write-Host "`nTo release a new version:"
Write-Host "  1. Build the installer: .\build.ps1"
Write-Host "  2. Upload ITToolkit-Setup-X.X.X.exe to GitHub releases"
Write-Host "  3. Update the gist manifest with the download URL and new version"
