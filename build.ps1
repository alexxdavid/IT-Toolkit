# IT Toolkit — full build.
# Produces:
#   build\bin\ITToolkit.exe          (portable app)
#
# The Inno Setup 7 installer is built separately because the unsigned EXE is a
# known Microsoft Defender false positive; compiling the installer needs an
# elevated shell with a temporary Defender exclusion:
#
#   .\installer\build-installer.ps1   (run as administrator)
#
# The durable fix is to code-sign ITToolkit.exe.

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $root

# Version comes from wails.json (info.productVersion)
$wails = Get-Content -LiteralPath 'wails.json' -Raw | ConvertFrom-Json
$version = if ($wails.info.productVersion) { $wails.info.productVersion } else { '1.0.0' }
Write-Host "==> Building IT Toolkit v$version" -ForegroundColor Cyan

# 1. Frontend dependencies + type check
Write-Host '==> Frontend' -ForegroundColor Cyan
Push-Location frontend
if (Test-Path -LiteralPath 'node_modules') {
    npm ci
} else {
    npm install
}
npm run check
Pop-Location

# 2. Wails production build (Go backend + bundled frontend)
Write-Host '==> Wails build' -ForegroundColor Cyan
wails build -clean -platform windows/amd64
if ($LASTEXITCODE -ne 0) { throw 'wails build failed' }

Write-Host ''
Write-Host 'Build complete:' -ForegroundColor Green
Get-ChildItem -LiteralPath 'build\bin\ITToolkit.exe' | Select-Object FullName, @{n='MB';e={[math]::Round($_.Length/1MB,1)}}
Write-Host ''
Write-Host 'Next: run .\installer\build-installer.ps1 as administrator to produce the setup EXE.' -ForegroundColor Cyan
Write-Host 'Note: unsigned ITToolkit.exe is a Microsoft Defender false positive on some machines;' -ForegroundColor Yellow
Write-Host 'sign it with a code-signing certificate for clean rollout to staff.' -ForegroundColor Yellow
