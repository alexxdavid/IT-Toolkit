; Solutions IT Toolkit — portable extract-only installer
; Compiled with Inno Setup 7 (ISCC)

#ifndef MyAppVersion
  #define MyAppVersion "236.12"
#endif

#define MyAppName "Solutions IT Toolkit"
#define MyAppExeName "ITToolkit.exe"
#define MyAppIcon "..\build\windows\icon.ico"

[Setup]
AppId={{9E7B0F4A-3C21-4D5E-8A9B-6F1C2D3E4A5B}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher=Solutions IT Toolkit
VersionInfoVersion={#MyAppVersion}
VersionInfoProductName={#MyAppName}
DefaultDirName={userdocs}\Solutions IT Toolkit
DisableProgramGroupPage=yes
DisableWelcomePage=yes
DisableReadyPage=yes
DisableStartupPrompt=yes
PrivilegesRequired=lowest
SetupIconFile={#MyAppIcon}
Compression=lzma2/ultra
SolidCompression=yes
WizardStyle=modern
Uninstallable=no
CreateUninstallRegKey=no
CreateAppDir=yes
OutputDir=..\build
OutputBaseFilename=ITToolkit-Setup-{#MyAppVersion}

[Dirs]
Name: "{app}\Repo"
Name: "{app}\Software"

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
; The application EXE plus the script library seed (imported into the DB on
; first launch so users can add/remove scripts freely).
Source: "..\build\bin\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\build\scripts_seed.json"; DestDir: "{app}"; Flags: ignoreversion

[Run]
; Nothing to launch after install by design — pure extraction.
