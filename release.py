"""
Solutions IT Toolkit — Automated Release Script
Usage: python release.py <version> [--notes "release notes"]

Steps:
  1. Bumps version in setup.iss and internal/update/update.go
  2. Runs full build (frontend + Go + Inno Setup)
  3. Creates GitHub Release and uploads the installer
  4. Updates the update manifest gist with the new installer URL
  5. Commits and pushes all changes

Requires:
  - GITHUB_TOKEN env var or .secrets/github/token file
"""

import os, sys, re, json, shutil, subprocess, time, urllib.request, urllib.parse, ssl
from datetime import datetime

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
GITHUB_USER = "alexxdavid"
SOURCE_REPO = "IT-Toolkit"
RELEASES_REPO = "IT-Toolkit"
GIST_ID = "603ee193205abab222bee17188995f5a"
MANIFEST_FILE = "update_manifest.json"


def get_token():
    t = os.environ.get("GITHUB_TOKEN", "")
    if not t:
        for p in [os.path.join(BASE_DIR, ".secrets", "github", "token"), os.path.join(BASE_DIR, ".secrets", "github", "github_token"), os.path.join(BASE_DIR, ".github_token")]:
            if os.path.exists(p):
                t = open(p).read().strip()
                break
    if not t:
        print("[ERROR] No token. Set GITHUB_TOKEN or create .github_token file."); sys.exit(1)
    return t


def github_api(method, url, token, data=None):
    headers = {"Authorization": f"Bearer {token}", "Accept": "application/vnd.github+json", "User-Agent": "ITToolkitRelease"}
    body = json.dumps(data).encode("utf-8") if data else None
    if body:
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(url, data=body, headers=headers, method=method)
    if method == "DELETE":
        with urllib.request.urlopen(req, timeout=60) as resp:
            return None
    with urllib.request.urlopen(req, timeout=60) as resp:
        raw = resp.read().decode("utf-8").strip()
        return json.loads(raw) if raw else None


def github_upload(url, token, file_path):
    with open(file_path, "rb") as f:
        body = f.read()
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/octet-stream", "User-Agent": "ITToolkitRelease"}
    for attempt in range(1, 6):
        req = urllib.request.Request(url, data=body, headers=headers, method="POST")
        try:
            with urllib.request.urlopen(req, timeout=600) as resp:
                return json.loads(resp.read().decode("utf-8"))
        except Exception as e:
            print(f"  Upload attempt {attempt}/5 failed: {e}")
            if attempt < 5:
                time.sleep(min(2 ** attempt, 30))
            else:
                raise


def bump_version(version):
    build = 1
    files = {
        "installer/setup.iss": (r'(#define\s+MyAppVersion\s+")([^"]+)(")', version),
        "internal/update/update.go": (r'(CurrentVersion\s*=\s*")([^"]+)(")', version),
    }
    for filepath, (pattern, replacement) in files.items():
        full = os.path.join(BASE_DIR, filepath)
        if not os.path.exists(full):
            print(f"  [SKIP] {filepath}"); continue
        with open(full, "r", encoding="utf-8") as f:
            content = f.read()
        new_content, count = re.subn(pattern, rf"\g<1>{replacement}\g<3>", content, count=1)
        if count > 0:
            with open(full, "w", encoding="utf-8") as f:
                f.write(new_content)
            print(f"  [OK] {filepath} -> v{version}")

    # Bump build number
    go_file = os.path.join(BASE_DIR, "internal", "update", "update.go")
    if os.path.exists(go_file):
        with open(go_file, "r", encoding="utf-8") as f:
            c = f.read()
        m = re.search(r'(CurrentBuild\s*=\s*)(\d+)', c)
        if m:
            build = int(m.group(2)) + 1
        new_c, _ = re.subn(r'(CurrentBuild\s*=\s*)(\d+)', rf'\g<1>{build}', c, count=1)
        with open(go_file, "w", encoding="utf-8") as f:
            f.write(new_c)
        print(f"  [OK] build -> {build}")
    return build


def find_iscc():
    candidates = [
        os.path.join(os.environ.get("ProgramFiles(x86)", ""), "Inno Setup 7", "ISCC.exe"),
        os.path.join(os.environ.get("ProgramFiles", ""), "Inno Setup 7", "ISCC.exe"),
        os.path.join(os.environ.get("LOCALAPPDATA", ""), "Programs", "Inno Setup 7", "ISCC.exe"),
    ]
    for c in candidates:
        if os.path.isfile(c):
            return c
    return None


def run_build(version):
    print("\n[RUNNING] Building...")
    # Frontend (npm is a .cmd shim on Windows — use shell=True)
    subprocess.run("npm ci", cwd=os.path.join(BASE_DIR, "frontend"), check=True, shell=True)
    subprocess.run("npm run check", cwd=os.path.join(BASE_DIR, "frontend"), check=True, shell=True)
    # Wails build (wails is a .cmd shim on Windows — use shell=True)
    subprocess.run("wails build -clean -platform windows/amd64", cwd=BASE_DIR, check=True, shell=True)
    # Script library seed (bundled into the installer, imported into the DB on first run)
    subprocess.run("go run ./cmd/genseed", cwd=BASE_DIR, check=True, shell=True)
    # Installer
    iscc = find_iscc()
    if iscc:
        subprocess.run([iscc, f"/DMyAppVersion={version}", os.path.join(BASE_DIR, "installer", "setup.iss")], cwd=BASE_DIR, check=True)
    print("  [OK] Build complete")


def find_installer(version):
    for name in [f"SITTOOLKIT-Setup-{version}.exe", "SITTOOLKIT-Setup-1.0.0.exe"]:
        p = os.path.join(BASE_DIR, "build", name)
        if os.path.exists(p):
            return p
    build_dir = os.path.join(BASE_DIR, "build")
    if os.path.isdir(build_dir):
        for f in os.listdir(build_dir):
            if f.endswith(".exe") and "Setup" in f:
                return os.path.join(build_dir, f)
    return None


def create_release(version, build, notes, installer_path, token):
    tag = f"v{version}"
    print(f"\n[RUNNING] Creating GitHub Release: {tag}")
    release = None
    try:
        existing = github_api("GET", f"https://api.github.com/repos/{GITHUB_USER}/{RELEASES_REPO}/releases/tags/{tag}", token)
        if existing.get("id"):
            release = existing
            github_api("PATCH", f"https://api.github.com/repos/{GITHUB_USER}/{RELEASES_REPO}/releases/{release['id']}", token,
                       {"name": f"Solutions IT Toolkit v{version}", "body": notes or f"Release v{version}"})
    except Exception:
        pass

    if release is None:
        release = github_api("POST", f"https://api.github.com/repos/{GITHUB_USER}/{RELEASES_REPO}/releases", token,
                             {"tag_name": tag, "name": f"Solutions IT Toolkit v{version}", "body": notes or f"Release v{version}", "draft": False, "prerelease": False})

    # Delete old installer assets for this version so only the newest remains.
    prefix = f"SITTOOLKIT-Setup-{version}"
    for a in release.get("assets", []):
        if a.get("name", "").startswith(prefix):
            try:
                github_api("DELETE", f"https://api.github.com/repos/{GITHUB_USER}/{RELEASES_REPO}/releases/assets/{a['id']}", token)
                print(f"  Deleted old asset {a['name']}")
            except Exception:
                pass

    upload_url = release["upload_url"].split("{")[0]
    asset_name = f"SITTOOLKIT-Setup-{version}.b{build}.exe"
    print(f"  Uploading {asset_name}...")
    github_upload(f"{upload_url}?name={urllib.parse.quote(asset_name)}", token, installer_path)

    release = github_api("GET", f"https://api.github.com/repos/{GITHUB_USER}/{RELEASES_REPO}/releases/tags/{tag}", token)
    print(f"  [OK] Release: {release.get('html_url', '')}")
    return release


def update_gist(version, build, installer_url, notes, token):
    if not GIST_ID:
        print("  [SKIP] No GIST_ID set — skipping manifest update."); return
    print("\n[RUNNING] Updating manifest gist...")
    manifest = {"version": version, "build": build, "installer_url": installer_url, "notes": notes or "", "force_update": False}
    github_api("PATCH", f"https://api.github.com/gists/{GIST_ID}", token, {"files": {MANIFEST_FILE: {"content": json.dumps(manifest, indent=2) + "\n"}}})
    print("  [OK] Manifest updated")


def git_commit_push(version, notes, token):
    print("\n[RUNNING] Git commit and push...")
    subprocess.run(["git", "add", "-A"], cwd=BASE_DIR, capture_output=True)
    subprocess.run(["git", "commit", "-m", f"release v{version}: {notes or 'update'}"], cwd=BASE_DIR, capture_output=True)
    origin = subprocess.run(["git", "remote", "get-url", "origin"], cwd=BASE_DIR, capture_output=True, text=True).stdout.strip()
    auth = origin.replace("https://github.com/", f"https://{token}@github.com/")
    subprocess.run(["git", "remote", "set-url", "origin", auth], cwd=BASE_DIR, capture_output=True)
    try:
        subprocess.run(["git", "push", "origin", "master"], cwd=BASE_DIR, capture_output=True)
        print("  [OK] Pushed")
    finally:
        subprocess.run(["git", "remote", "set-url", "origin", origin], cwd=BASE_DIR, capture_output=True)


def main():
    if len(sys.argv) < 2:
        print("Usage: python release.py <version> [--notes 'text']\n")
        print("  python release.py 1.1.0 --notes 'Category management, custom repos, update system'")
        sys.exit(1)

    version = sys.argv[1]
    notes = ""
    if "--notes" in sys.argv:
        idx = sys.argv.index("--notes")
        if idx + 1 < len(sys.argv):
            notes = sys.argv[idx + 1]

    token = get_token()
    print(f"{'='*60}\n  SOFTWARE SOLUTIONS IT TOOLKIT — RELEASE v{version}\n{'='*60}\n")

    build = bump_version(version)
    run_build(version)
    installer = find_installer(version)
    if not installer:
        print("[ERROR] Installer not found!"); sys.exit(1)

    release = create_release(version, build, notes, installer, token)
    current_asset = None
    for a in release.get("assets", []):
        if a.get("name", "").startswith(f"SITTOOLKIT-Setup-{version}"):
            current_asset = a; break
    # Use browser_download_url so the app can download the binary directly
    # (API asset URLs need an Accept header the update client may not send).
    installer_url = current_asset.get("browser_download_url", "") if current_asset else ""
    update_gist(version, build, installer_url, notes, token)
    git_commit_push(version, notes, token)

    print(f"\n{'='*60}\n  RELEASE v{version} DONE!\n  Release: {release.get('html_url', '')}\n{'='*60}\n")


if __name__ == "__main__":
    main()
