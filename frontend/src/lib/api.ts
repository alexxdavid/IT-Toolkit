import * as App from '../../wailsjs/go/main/App'
import type { CatalogView, RepoDetail, SearchResult, ExportResult, RunResult, ScriptFile, RepoInfo, InstallResult } from '../lib/types'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export type ScanProgress = {
  folderId: number
  folder: string
  phase: string
  current: number
  total: number
}

export async function pickFolder(): Promise<string | null> {
  return App.PickFolder()
}

export function onScanProgress(cb: (p: ScanProgress) => void): () => void {
  return EventsOn('scan:progress', cb)
}

export function onScanDone(cb: () => void): () => void {
  return EventsOn('scan:done', cb)
}

export function onScanError(cb: (e: { folder: string; error: string }) => void): () => void {
  return EventsOn('scan:error', cb)
}

export type InstallProgress = {
  index: number
  total: number
  name: string
  status: string
  message: string
}

export function onInstallProgress(cb: (p: InstallProgress) => void): () => void {
  return EventsOn('install:progress', cb)
}

export function onInstallDone(cb: () => void): () => void {
  return EventsOn('install:done', cb)
}

export const api = {
  getCatalog: (): Promise<CatalogView> => App.GetCatalog(),
  rescan: (): Promise<string> => App.RescanAll(),
  addFolder: (path: string) => App.AddFolder(path),
  removeFolder: (id: number): Promise<void> => App.RemoveFolder(id),
  repoDetail: (id: number): Promise<RepoDetail> => App.GetRepoDetail(id),
  scriptContent: (path: string): Promise<string> => App.GetScriptContent(path),
  listScripts: (lang: string, q: string): Promise<ScriptFile[]> => App.ListScripts(lang, q),
  search: (query: string, inContent: boolean): Promise<SearchResult> => App.Search(query, inContent),
  runScript: (path: string): Promise<RunResult> => App.RunScript(path),
  copyScript: (path: string): Promise<string> => App.CopyScript(path),
  exportScripts: (paths: string[], dest: string): Promise<ExportResult> => App.ExportScripts(paths, dest),
  reveal: (path: string): Promise<void> => App.RevealInExplorer(path),
  exeDir: (): Promise<string> => App.GetExeDir(),
  defaultScriptsDir: (): Promise<string> => App.GetDefaultScriptsDir(),
  discoverRepos: (folder: string): Promise<RepoInfo[]> => App.DiscoverRepos(folder),
  gitAvailable: (): Promise<boolean> => App.GitAvailable(),
  defaultInstallDir: (): Promise<string> => App.GetDefaultInstallDir(),
  defaultSoftwareDir: (): Promise<string> => App.GetDefaultSoftwareDir(),
  getSoftwareCatalog: (): Promise<any[]> => App.GetSoftwareCatalog(),
  getSoftwareVersions: (): Promise<Record<string, any>> => App.GetSoftwareVersions(),
  invalidateSoftwareVersions: (): Promise<void> => App.InvalidateSoftwareVersions(),
  downloadSoftware: (name: string, url: string, dest: string): Promise<string> => App.DownloadSoftware(name, url, dest),
  getSoftwareProgress: (): Promise<any> => App.GetSoftwareProgress(),
  installRepos: (names: string[], dest: string): Promise<InstallResult[]> => App.InstallRepos(names, dest),
  getRecommendedRepos: (): Promise<any[]> => App.GetRecommendedRepos(),
  isRepoInstalled: (name: string, dest: string): Promise<boolean> => App.IsRepoInstalled(name, dest),
  removeRepo: (name: string, dest: string): Promise<void> => App.RemoveRepo(name, dest),
  checkForUpdate: (force: boolean): Promise<any> => App.CheckForUpdate(force),
  downloadUpdate: (url: string, version: string): Promise<string> => App.DownloadUpdate(url, version),
  getUpdateProgress: (): Promise<any> => App.GetUpdateProgress(),
  applyUpdate: (path: string): Promise<void> => App.ApplyUpdate(path),
  getCurrentVersion: (): Promise<string> => App.GetCurrentVersion(),
  listCustomCategories: (): Promise<any[]> => App.ListCustomCategories(),
  createCustomCategory: (name: string): Promise<any> => App.CreateCustomCategory(name),
  renameCustomCategory: (id: number, name: string): Promise<void> => App.RenameCustomCategory(id, name),
  deleteCustomCategory: (id: number): Promise<void> => App.DeleteCustomCategory(id),
  listCustomRepos: (): Promise<any[]> => App.ListCustomRepos(),
  addCustomRepo: (name: string, url: string, category: string, summary: string): Promise<any> => App.AddCustomRepo(name, url, category, summary),
  removeCustomRepo: (id: number): Promise<void> => App.RemoveCustomRepo(id),
  listCustomSoftware: (): Promise<any[]> => App.ListCustomSoftware(),
  addCustomSoftware: (name: string, version: string, category: string, download: string, notes: string, wingetId: string): Promise<any> => App.AddCustomSoftware(name, version, category, download, notes, wingetId),
  removeCustomSoftware: (id: number): Promise<void> => App.RemoveCustomSoftware(id),
  getRecommendedReposCombined: (): Promise<any[]> => App.GetRecommendedReposCombined()
}
