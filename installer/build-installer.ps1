# Builds the Inno Setup 7 installer. MUST be run from an elevated PowerShell
# because the unsigned ITToolkit.exe is a known false positive for Microsoft
# Defender heuristics; an exclusion for the build folder is added temporarily.
#
#   Right-click PowerShell > Run as administrator, then:
#   .\installer\build-installer.ps1
#
# The durable fix is to code-sign ITToolkit.exe with a trusted certificate;
# then this script is no longer needed.

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $root

$iscc = @(
    "$env:ProgramFiles(x86)\Inno Setup 7\ISCC.exe",
    "$env:ProgramFiles\Inno Setup 7\ISCC.exe",
    "$env:LOCALAPPDATA\Programs\Inno Setup 7\ISCC.exe"
) | Where-Object { Test-Path -LiteralPath $_ } | Select-Object -First 1
if (-not $iscc) {
    throw 'Inno Setup 7 (ISCC.exe) not found. Install with: winget install JRSoftware.InnoSetup.7'
}

$wails = Get-Content -LiteralPath 'wails.json' -Raw | ConvertFrom-Json
$version = if ($wails.info.productVersion) { $wails.info.productVersion } else { '1.0.0' }

# Temporary Defender exclusion for the build output (requires elevation).
$exclusionPath = Join-Path $root 'build'
Add-MpPreference -ExclusionPath $exclusionPath -ErrorAction SilentlyContinue
try {
    & $iscc "/DMyAppVersion=$version" 'installer\setup.iss'
    if ($LASTEXITCODE -ne 0) { throw 'ISCC failed' }
} finally {
    # Restore the previous exclusion state.
    $current = Get-MpPreference | Select-Object -ExpandProperty ExclusionPath
    if ($current -contains $exclusionPath) {
        Remove-MpPreference -ExclusionPath $exclusionPath -ErrorAction SilentlyContinue
    }
}

Write-Host ''
Write-Host "Installer built: build\ITToolkit-Setup-$version.exe" -ForegroundColor Green
Write-Host ''
Write-Host 'IMPORTANT: The app EXE is unsigned. Staff machines with Microsoft Defender'
Write-Host 'may block ITToolkit.exe. Recommended before wide rollout:' -ForegroundColor Yellow
Write-Host '  1. Sign ITToolkit.exe with a trusted code-signing certificate, or' -ForegroundColor Yellow
Write-Host '  2. Roll out a Defender exclusion for the install folder via policy, or' -ForegroundColor Yellow
Write-Host '  3. Distribute via Intune/approved-application channels.' -ForegroundColor Yellow
