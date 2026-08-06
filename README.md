# Solutions IT Toolkit

A Windows desktop app for browsing, searching, viewing, copying, exporting and
running the IT script collection distributed alongside it. Built with
**Go + Wails 2 + Svelte 5 + Tailwind CSS 4**, packaged as a portable EXE and an
Inno Setup 7 extract-only installer.

## How it works

- Ships as an EXE that sits in the same folder as a `Scripts` folder
  (originally `C:\GitHub_IT_Scripts`). On first launch the app scans the
  `Scripts` folder next to the EXE, or any folder the user adds.
- The catalog is indexed into `%LOCALAPPDATA%\ITToolkit\index.db` (SQLite +
  FTS5). First index of a large collection (~250 repos / ~75k script files)
  takes a few minutes with a progress bar; later rescans are incremental and
  fast.
- Repositories are auto-categorised by folder name using a keyword map in
  `internal/catalog/categories.go`. Unknown repos land in *Uncategorized*.
- Content search indexes PowerShell/Batch/VBScript/Shell/SQL/Registry scripts.
  Python files (mostly large library/AI repos) are listed and name-searchable
  but not content-indexed to keep indexing fast; their contents are read from
  disk on demand.

## Features

- Category sidebar with live counts, folder filter, file-type filter
- Repo cards with expandable script lists, README viewer (rendered markdown)
- Script viewer with syntax highlighting, Copy / Run (elevated) / Reveal
- Multi-select export preserving `repo/subfolder/file` layout
- Content search across script file contents (FTS5)

## Build

Requires Go 1.25+, Node 20+, Wails CLI 2.12+, Inno Setup 7.

```powershell
# 1. Full app build (frontend type-check + EXE)
.\build.ps1

# 2. Installer (MUST be run as administrator, see Defender note below)
.\installer\build-installer.ps1
```

Outputs: `build\bin\ITToolkit.exe` and `build\ITToolkit-Setup-<ver>.exe`.

## Deploy

Give staff the EXE and the `Scripts` folder together (same folder), or run the
setup EXE which extracts the app to a folder of their choice. No shortcuts or
registry entries are created.

## Microsoft Defender false positive — important

The unsigned Go/Wails EXE is flagged by Microsoft Defender heuristics on some
machines (a well-known issue with Go desktop binaries that spawn processes and
scan files). Until the app is code-signed with a trusted certificate,
deploying to staff requires one of:

1. **Code-sign `ITToolkit.exe`** with an OV/EV Windows code-signing certificate
   (the durable fix — recommended before wide rollout), or
2. Add a Defender exclusion for the app/install folder via policy or
   `Set-MpPreference -ExclusionPath <dir>`, or
3. Distribute through an approved channel (e.g. Intune).

The elevated installer script (`installer\build-installer.ps1`) adds a
temporary exclusion for the build folder so the setup EXE can be compiled.
