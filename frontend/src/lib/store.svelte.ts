import type { CatalogView, Repo, RepoDetail, ScriptFile } from '../lib/types'

export const ui = $state({
  loaded: false,
  scanning: false,
  scanStatus: null as { phase: string; current: number; total: number } | null,
  selectedCategory: 'all' as string,
  selectedFolder: 'all' as number | 'all',
  typeFilter: 'all' as string,
  query: '',
  contentSearch: false,
  browserLoading: false,
  exportOpen: false,
  libraryOpen: true,
  settingsOpen: false,
  localTab: 'local' as 'repos' | 'software' | 'scripts' | 'local' | 'favorites',
  scriptsFolder: '' as string,
  loadError: '' as string
})

export const catalog = $state<CatalogView>({
  folders: [],
  categories: [],
  repos: [],
  totalScripts: 0
} as unknown as CatalogView)

export const detailCache = $state(new Map<number, RepoDetail>())

export const scriptBrowser = $state<{ lang: string; scripts: ScriptFile[] }>({
  lang: 'all',
  scripts: []
})

export const selection = $state(new Set<string>())

export function toggleSelection(path: string) {
  if (selection.has(path)) {
    selection.delete(path)
  } else {
    selection.add(path)
  }
}

export function clearSelection() {
  selection.clear()
}

export function setScanning(v: boolean) {
  ui.scanning = v
}

export function setScanStatus(v: { phase: string; current: number; total: number } | null) {
  ui.scanStatus = v
}

export function setCategory(c: string) {
  ui.selectedCategory = c
}

export function setFolder(f: number | 'all') {
  ui.selectedFolder = f
}

export function setTypeFilter(t: string) {
  ui.typeFilter = t
}

export function setQuery(q: string) {
  ui.query = q
}

export function setContentSearch(v: boolean) {
  ui.contentSearch = v
}

export function setBrowserLoading(v: boolean) {
  ui.browserLoading = v
}

export function setExportOpen(v: boolean) {
  ui.exportOpen = v
}

export function setLocalTab(v: 'repos' | 'software' | 'scripts' | 'local' | 'favorites') {
  ui.localTab = v
}

export function setLibraryOpen(v: boolean) {
  ui.libraryOpen = v
}

export function setSettingsOpen(v: boolean) {
  ui.settingsOpen = v
}

export function setCatalog(v: CatalogView) {
  catalog.folders = v.folders
  catalog.categories = v.categories
  catalog.repos = v.repos
  catalog.totalScripts = v.totalScripts
  ui.loaded = true
}

export function getDetail(id: number): RepoDetail | undefined {
  return detailCache.get(id)
}

export function setDetail(id: number, d: RepoDetail) {
  detailCache.set(id, d)
}

export function visibleRepos(): Repo[] {
  let repos = catalog.repos
  if (ui.selectedFolder !== 'all') {
    repos = repos.filter((r) => r.folderId === ui.selectedFolder)
  }
  if (ui.selectedCategory !== 'all') {
    repos = repos.filter((r) => r.category === ui.selectedCategory)
  }
  if (ui.query.trim()) {
    const q = ui.query.trim().toLowerCase()
    repos = repos.filter(
      (r) => r.name.toLowerCase().includes(q) || r.category.toLowerCase().includes(q)
    )
  }
  return repos
}
