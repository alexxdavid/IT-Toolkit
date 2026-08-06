<script lang="ts">
  import { onMount } from 'svelte'
  import { AlertTriangle, X } from 'lucide-svelte'
  import Sidebar from './components/Sidebar.svelte'
  import Toolbar from './components/Toolbar.svelte'
  import RepoGrid from './components/RepoGrid.svelte'
  import ScriptBrowser from './components/ScriptBrowser.svelte'
  import ScriptViewer from './components/ScriptViewer.svelte'
  import ExportDialog from './components/ExportDialog.svelte'
  import GithubLibrary from './components/GithubLibrary.svelte'
  import UpdateBanner from './components/UpdateBanner.svelte'
  import Settings from './components/Settings.svelte'
  import { api, onScanProgress, onScanDone, onScanError } from './lib/api'
  import type { SearchResult } from './lib/types'
  import {
    ui,
    setCatalog,
    catalog,
    scriptBrowser,
    setBrowserLoading,
    setScanning,
    setScanStatus,
    setCategory,
    setFolder,
    setQuery,
    setContentSearch,
    setTypeFilter
  } from './lib/store.svelte.ts'
  import { categoryColor } from './lib/categories'

  let offs: Array<() => void> = []
  let scanErr = $state<string | null>(null)
  let searchRes: SearchResult | null = $state(null)
  let searchLoading = $state(false)
  let searchTimer: ReturnType<typeof setTimeout> | undefined
  let browseTimer: ReturnType<typeof setTimeout> | undefined

  onMount(() => {
    offs.push(
      onScanProgress((p) => {
        setScanning(true)
        setScanStatus({ phase: p.phase, current: p.current, total: p.total })
      })
    )
    offs.push(
      onScanDone(() => {
        setScanning(false)
        setScanStatus(null)
        refresh()
      })
    )
    offs.push(
      onScanError((e) => {
        scanErr = e.error
        setTimeout(() => (scanErr = null), 6000)
      })
    )
    refresh()
    return () => {
      offs.forEach((f) => f())
      clearTimeout(searchTimer)
      clearTimeout(browseTimer)
    }
  })

  async function refresh() {
    try {
      const v = await api.getCatalog()
      ui.loadError = ''
      setCatalog(v)
    } catch (e) {
      ui.loadError = String(e)
    }
  }

  const mode = $derived(
    ui.contentSearch && ui.query.trim() ? 'search' : ui.typeFilter !== 'all' ? 'scripts' : 'browse'
  )

  $effect(() => {
    if (mode !== 'search') return
    const q = ui.query.trim()
    searchLoading = true
    searchTimer = setTimeout(async () => {
      try {
        searchRes = await api.search(q, true)
      } catch {
        searchRes = null
      } finally {
        searchLoading = false
      }
    }, 300)
    return () => clearTimeout(searchTimer)
  })

  $effect(() => {
    if (mode !== 'scripts') return
    const lang = ui.typeFilter
    const q = ui.query.trim()
    setBrowserLoading(true)
    scriptBrowser.lang = lang
    scriptBrowser.scripts = []
    browseTimer = setTimeout(async () => {
      try {
        scriptBrowser.scripts = await api.listScripts(lang, q)
      } catch {
        scriptBrowser.scripts = []
      } finally {
        setBrowserLoading(false)
      }
    }, 300)
    return () => clearTimeout(browseTimer)
  })

  function gotoRepo(r: { folderId: number; category: string }) {
    setCategory(r.category)
    setFolder(r.folderId)
    setQuery('')
    setContentSearch(false)
    setTypeFilter('all')
  }
</script>

<div class="flex h-screen flex-col bg-[#0b1120]">
  {#if !ui.loaded && ui.loadError}
    <div class="flex flex-1 flex-col items-center justify-center gap-3 text-slate-400">
      <AlertTriangle size="32" class="text-rose-400" />
      <p class="text-sm">Could not load the catalog: <span class="font-mono text-rose-300">{ui.loadError}</span></p>
      <button
        class="rounded-lg bg-cyan-500/15 px-4 py-2 text-sm font-medium text-cyan-300 ring-1 ring-cyan-500/40 transition hover:bg-cyan-500/25"
        onclick={refresh}
      >
        Retry
      </button>
    </div>
  {:else if !ui.loaded}
    <div class="flex flex-1 items-center justify-center gap-3 text-slate-500">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-slate-700 border-t-cyan-500"></div>
      <p class="text-sm">Loading catalog…</p>
    </div>
  {:else if ui.settingsOpen}
    <Settings />
  {:else if catalog.folders.length === 0 || ui.libraryOpen}
    <GithubLibrary />
  {:else}
    <div class="flex min-h-0 flex-1 overflow-hidden">
      <Sidebar />
      <div class="flex min-w-0 flex-1 flex-col">
        <Toolbar />
        {#if mode === 'browse'}
          <RepoGrid />
        {:else if mode === 'scripts'}
          <ScriptBrowser
            scripts={scriptBrowser.scripts}
            loading={ui.browserLoading}
            title={ui.typeFilter === 'all' ? 'All scripts' : `${ui.typeFilter} scripts`}
          />
        {:else}
          <div class="flex min-h-0 flex-1 flex-col">
            {#if searchRes && searchRes.repos.length > 0}
              <div class="shrink-0 border-b border-slate-800 px-4 py-3">
                <div class="mb-2 text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                  Matching repositories
                </div>
                <div class="flex flex-wrap gap-2">
                  {#each searchRes.repos as r (r.id)}
                    <button
                      class="flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/60 px-3 py-1.5 text-xs text-slate-200 transition hover:border-slate-500"
                      onclick={() => gotoRepo(r)}
                    >
                      <span class="font-mono">{r.name}</span>
                      <span class="rounded border px-1.5 py-0.5 text-[9px] {categoryColor(r.category)}"
                        >{r.category}</span
                      >
                    </button>
                  {/each}
                </div>
              </div>
            {/if}
            <ScriptBrowser
              scripts={searchRes?.scripts ?? []}
              loading={searchLoading}
              title="Matching scripts (file contents)"
            />
          </div>
        {/if}
      </div>
    </div>
  {/if}

  <ScriptViewer />
  <ExportDialog />
  <UpdateBanner />

  {#if scanErr}
    <div
      class="fixed bottom-4 right-4 z-50 flex max-w-md items-start gap-3 rounded-lg border border-rose-500/40 bg-rose-950/90 px-4 py-3 text-sm text-rose-200 shadow-xl"
    >
      <AlertTriangle size="16" class="mt-0.5 shrink-0" />
      <div class="min-w-0 flex-1">{scanErr}</div>
      <button class="shrink-0 text-rose-400 hover:text-rose-200" onclick={() => (scanErr = null)}>
        <X size="14" />
      </button>
    </div>
  {/if}
</div>
