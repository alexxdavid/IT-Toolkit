<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { GitBranch, Download, CheckCircle2, Loader2, FolderSearch, Package, Search, X, HardDrive, Settings } from 'lucide-svelte'
  import { api, pickFolder, onInstallProgress, onInstallDone } from '../lib/api'
  import type { InstallProgress } from '../lib/api'
  import { ui } from '../lib/store.svelte.ts'
  import { categoryColor } from '../lib/categories'
  import { setSettingsOpen, setLibraryOpen } from '../lib/store.svelte.ts'
  import SoftwareTab from './SoftwareTab.svelte'

  type Tab = 'recommended' | 'discovered' | 'software'
  type CategoryGroup = { name: string; repos: any[] }

  let activeTab = $state<Tab>('recommended')
  let selectedCategory = $state('all')
  let recommended = $state<any[]>([])
  let discoverFolder = $state('')
  let discovered = $state<any[]>([])
  let discoverLoading = $state(false)
  let discoverErr = $state<string | null>(null)
  let installDest = $state('')
  let selected = $state<Record<string, boolean>>({})
  let installing = $state(false)
  let installResults = $state<any[]>([])
  let installProgress = $state<InstallProgress | null>(null)
  let gitOK = $state(true)
  let searchQuery = $state('')
  let addedToCatalog = $state(false)
  let installedRepos = $state<Record<string, boolean>>({})
  let expandedRepo = $state<string | null>(null)
  let repoDetail = $state<any>(null)
  let repoDetailLoading = $state(false)
  let checking = $state(false)
  let updateInfo = $state<any>(null)
  let showUpdateNote = $state(false)
  let downloading = $state(false)
  let downloadPercent = $state(0)
  let offProgress: (() => void) | null = null
  let offDone: (() => void) | null = null

  let allRepos = $derived(activeTab === 'recommended' ? recommended : discovered)

  let categories = $derived.by((): CategoryGroup[] => {
    const q = searchQuery.trim().toLowerCase()
    const filtered = q ? allRepos.filter((r: any) => r.name.toLowerCase().includes(q) || r.category.toLowerCase().includes(q) || (r.summary && r.summary.toLowerCase().includes(q))) : allRepos
    const map = new Map<string, any[]>()
    for (const r of filtered) {
      const cat = r.category || 'Other'
      if (!map.has(cat)) map.set(cat, [])
      map.get(cat)!.push(r)
    }
    return Array.from(map.entries()).sort((a, b) => a[0].localeCompare(b[0])).map(([name, repos]) => ({ name, repos }))
  })

  let visibleRepos = $derived(
    searchQuery.trim()
      ? categories.flatMap((c: CategoryGroup) => c.repos)
      : selectedCategory === 'all'
        ? allRepos
        : categories.find((c: CategoryGroup) => c.name === selectedCategory)?.repos || []
  )

  let totalRepos = $derived(allRepos.length)

  onMount(async () => {
    recommended = await api.getRecommendedReposCombined() || []
    gitOK = await api.gitAvailable()
    const d = await api.defaultInstallDir()
    if (d) installDest = d
  })

  onDestroy(() => { offProgress?.(); offDone?.() })

  function toggleRepo(name: string) {
    selected[name] = !selected[name]
  }

  function toggleAll() {
    const all = visibleRepos.every((r: any) => selected[r.name])
    visibleRepos.forEach((r: any) => { selected[r.name] = !all })
  }

  function selectedCount(): number {
    return Object.values(selected).filter(Boolean).length
  }

  async function doInstall() {
    const count = selectedCount()
    if (count === 0 || !installDest) return
    installing = true; installResults = []; installProgress = null
    const names = Object.keys(selected).filter(k => selected[k])
    offProgress = onInstallProgress((p) => { installProgress = p })
    offDone = onInstallDone(() => { installing = false; offProgress?.(); offDone?.(); offProgress = null; offDone = null })
    try {
      installResults = await api.installRepos(names, installDest)
      if (installResults.some((r: any) => r.status === 'installed')) {
        try { await api.addFolder(installDest); addedToCatalog = true } catch { }
      }
      await checkInstalled()
    } catch (e) { installResults = [{ name: 'error', status: 'failed', message: String(e) }] }
    installing = false; offProgress?.(); offDone?.(); offProgress = null; offDone = null
  }

  async function discover() {
    if (!discoverFolder) return
    discoverLoading = true; discoverErr = null
    try { discovered = await api.discoverRepos(discoverFolder) || [] } catch (e) { discoverErr = String(e); discovered = [] }
    discoverLoading = false
  }

  async function browseFolder() { const p = await pickFolder(); if (p) discoverFolder = p }

  async function selectRepo(repo: any) {
    if (expandedRepo === repo.name) { expandedRepo = null; return }
    expandedRepo = repo.name
    repoDetail = null
    repoDetailLoading = true
    try {
      const slug = repo.url?.replace('https://github.com/', '') || repo.name
      const resp = await fetch(`https://api.github.com/repos/${slug}`, { signal: AbortSignal.timeout(5000) })
      if (resp.ok) {
        const data = await resp.json()
        repoDetail = {
          description: data.description || repo.summary || '',
          stars: data.stargazers_count || 0,
          language: data.language || '',
          topics: (data.topics || []).slice(0, 8),
          updated: data.updated_at || '',
          license: data.license?.spdx_id || '',
        }
      } else {
        repoDetail = { description: repo.summary || '', stars: 0, language: '', topics: [], updated: '', license: '' }
      }
    } catch {
      repoDetail = { description: repo.summary || '', stars: 0, language: '', topics: [], updated: '', license: '' }
    }
    repoDetailLoading = false
  }

  async function checkUpdate() {
    checking = true
    try {
      const info = await api.checkForUpdate(true)
      if (info?.available) { updateInfo = info; showUpdateNote = true }
      else { updateInfo = null }
    } catch { /* */ }
    checking = false
  }

  async function installUpdate() {
    if (!updateInfo) return
    downloading = true; downloadPercent = 0
    const poll = setInterval(async () => {
      try { const p = await api.getUpdateProgress(); if (p?.percent > 0) downloadPercent = Math.min(100, Math.round(p.percent)) } catch { }
    }, 300)
    try {
      const path = await api.downloadUpdate(updateInfo.installer_url, updateInfo.version)
      clearInterval(poll)
      if (path) { downloadPercent = 100; await api.applyUpdate(path) }
    } catch { clearInterval(poll) }
    downloading = false
  }

  async function checkInstalled() {
    const result: Record<string, boolean> = {}
    for (const r of allRepos) {
      try { result[r.name] = await api.isRepoInstalled(r.name, installDest) } catch { result[r.name] = false }
    }
    installedRepos = result
  }

  async function removeRepo(name: string) {
    if (!confirm(`Remove ${name.split('/').pop()} from ${installDest}?`)) return
    try {
      await api.removeRepo(name, installDest)
      installedRepos[name] = false
      selected[name] = false
    } catch (e) { alert(String(e)) }
  }

  $effect(() => {
    if (allRepos.length > 0 && installDest) checkInstalled()
  })
</script>

<div class="flex h-full overflow-hidden">
  <!-- Left panel: categories -->
  <div class="flex w-56 shrink-0 flex-col border-r border-slate-800 bg-slate-900/60">
    <div class="border-b border-slate-800 p-3">
      <div class="flex gap-1 rounded-lg bg-slate-800/60 p-0.5">
        <button class="flex-1 rounded-md px-2 py-1.5 text-[10px] font-bold transition {activeTab === 'recommended' ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-slate-950 shadow-lg' : 'text-slate-400 hover:text-slate-200'}"
          onclick={() => { activeTab = 'recommended'; selectedCategory = 'all' }}>
          <Package size="12" class="inline mr-1" />Repos
        </button>
        <button class="flex-1 rounded-md px-2 py-1.5 text-[10px] font-bold transition {activeTab === 'software' ? 'bg-gradient-to-r from-sky-500 to-blue-500 text-slate-950 shadow-lg' : 'text-slate-400 hover:text-slate-200'}"
          onclick={() => activeTab = 'software'}>
          <HardDrive size="12" class="inline mr-1" />Software
        </button>
        <button class="flex-1 rounded-md px-2 py-1.5 text-[10px] font-bold transition {activeTab === 'discovered' ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-slate-950 shadow-lg' : 'text-slate-400 hover:text-slate-200'}"
          onclick={() => { activeTab = 'discovered'; selectedCategory = 'all'; if (discovered.length === 0) discover() }}>
          <FolderSearch size="12" class="inline mr-1" />Local
        </button>
      </div>
    </div>
    <div class="flex-1 overflow-y-auto p-2 space-y-0.5">
      {#if activeTab !== 'software'}
        <button class="flex w-full items-center justify-between rounded-lg px-2.5 py-2 text-left text-xs transition {selectedCategory === 'all' ? 'bg-emerald-500/15 text-emerald-300 font-semibold' : 'text-slate-300 hover:bg-slate-800'}"
          onclick={() => selectedCategory = 'all'}>
          <span>All repos</span>
          <span class="rounded bg-slate-800/80 px-1.5 py-0.5 text-[10px] text-slate-400">{totalRepos}</span>
        </button>
        {#each categories as cat (cat.name)}
          <button class="flex w-full items-center gap-2 rounded-lg px-2.5 py-2 text-left text-xs transition {selectedCategory === cat.name ? 'bg-emerald-500/15 text-emerald-300 font-semibold' : 'text-slate-300 hover:bg-slate-800'}"
            onclick={() => selectedCategory = cat.name}>
            <span class="truncate">{cat.name}</span>
            <span class="ml-auto shrink-0 rounded bg-slate-800/80 px-1.5 py-0.5 text-[10px] text-slate-400">{cat.repos.length}</span>
          </button>
        {/each}
      {:else}
        <div class="px-2.5 py-4 text-center text-[11px] text-slate-600">
          Browse and download software directly from the list on the right.
        </div>
      {/if}
    </div>
    <div class="border-t border-slate-800 p-3 space-y-2">
      {#if activeTab !== 'software'}
        <button class="flex w-full items-center justify-center gap-2 rounded-lg px-3 py-2.5 text-xs font-bold transition shadow-lg disabled:opacity-40 {selectedCount() > 0 ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-slate-950 shadow-emerald-500/20 hover:scale-105' : 'bg-slate-800 text-slate-400 cursor-default'}"
          onclick={doInstall} disabled={installing || selectedCount() === 0}>
          {#if installing}
            <Loader2 size="14" class="animate-spin" /> Installing…
          {:else if selectedCount() > 0}
            <Download size="14" /> Install {selectedCount()} repo{selectedCount() === 1 ? '' : 's'}
          {:else}
            Tick repos on the right to install
          {/if}
        </button>
        <input value={installDest} oninput={(e) => installDest = e.currentTarget.value}
          placeholder="Install to…"
          class="w-full rounded-lg border border-white/10 bg-slate-800/60 py-1.5 px-2.5 text-[11px] text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none" />
      {/if}
      <button class="flex w-full items-center gap-2 rounded-lg px-2.5 py-2 text-xs text-slate-400 transition hover:bg-slate-800 hover:text-slate-200"
        onclick={() => { setSettingsOpen(true); setLibraryOpen(false) }}>
        <Settings size="13" /> Settings
      </button>
    </div>
  </div>

  <!-- Right panel: repos or software -->
  <div class="flex min-w-0 flex-1 flex-col">
    {#if activeTab === 'software'}
      <SoftwareTab />
    {:else}
    <div class="flex shrink-0 items-center gap-3 border-b border-slate-800 bg-slate-900/60 px-4 py-2.5">
      <div class="relative flex-1 max-w-md">
        <Search size="14" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
        <input value={searchQuery} oninput={(e) => searchQuery = e.currentTarget.value}
          placeholder="Filter repos…"
          class="w-full rounded-lg border border-white/10 bg-slate-800/60 py-2 pl-9 pr-8 text-xs text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none" />
        {#if searchQuery}
          <button class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-200" onclick={() => searchQuery = ''}><X size="12" /></button>
        {/if}
      </div>
      <div class="flex-1"></div>
      <span class="text-[11px] text-slate-500">{visibleRepos.length} repos</span>
      {#if visibleRepos.length > 0}
        <button class="rounded-md border border-white/10 px-2.5 py-1 text-[10px] font-medium text-slate-300 transition hover:border-white/20 hover:text-white"
          onclick={toggleAll}>
          {visibleRepos.every((r: any) => selected[r.name]) ? 'Deselect all' : 'Select all'}
        </button>
      {/if}
    </div>

    {#if activeTab === 'discovered' && discovered.length === 0}
      <div class="flex flex-1 flex-col items-center justify-center gap-3 p-6">
        <div class="flex w-full max-w-lg items-center gap-2">
          <input value={discoverFolder} oninput={(e) => discoverFolder = e.currentTarget.value}
            placeholder="Scripts folder path…"
            class="min-w-0 flex-1 rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-xs text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none" />
          <button class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-300 hover:text-white transition" onclick={browseFolder}>Browse</button>
          <button class="rounded-lg bg-gradient-to-r from-emerald-500 to-teal-500 px-4 py-2 text-xs font-bold text-slate-950 shadow-lg shadow-emerald-500/20 transition hover:scale-105 disabled:opacity-50"
            onclick={discover} disabled={discoverLoading}>
            {#if discoverLoading}<Loader2 size="14" class="animate-spin" />{:else}Discover{/if}
          </button>
        </div>
        {#if !gitOK}
          <div class="rounded-lg border border-amber-500/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-300">Git not installed. Install with: winget install Git.Git</div>
        {/if}
        {#if discoverErr}
          <div class="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-300">{discoverErr}</div>
        {/if}
      </div>
    {:else}
      <div class="flex-1 overflow-y-auto">
        {#if visibleRepos.length === 0}
          <div class="flex h-full items-center justify-center text-sm text-slate-500">No repos found.</div>
        {:else}
          <div class="divide-y divide-white/5">
            {#each visibleRepos as repo, idx (repo.name + '_' + idx)}
              {@const isInstalled = installedRepos[repo.name]}
              <div class="flex items-center gap-3 px-4 py-3 transition hover:bg-white/5 {selected[repo.name] ? 'bg-emerald-500/5' : ''}">
                {#if isInstalled}
                  <button class="h-3.5 w-3.5 shrink-0 rounded border border-emerald-500/50 bg-emerald-500/20 text-center text-[8px] leading-[14px] text-emerald-400" title="Installed — click to remove"
                    onclick={() => removeRepo(repo.name)}>✓</button>
                {:else}
                  <button class="h-3.5 w-3.5 shrink-0 rounded border border-slate-600 bg-transparent transition hover:border-slate-400"
                    onclick={() => toggleRepo(repo.name)}>
                    {#if selected[repo.name]}<span class="block h-full w-full rounded-sm bg-emerald-500"></span>{/if}
                  </button>
                {/if}
                <button class="min-w-0 flex-1 text-left" onclick={() => isInstalled ? null : toggleRepo(repo.name)}>
                  <div class="flex items-center gap-2">
                    <span class="font-mono text-xs font-semibold text-white">{repo.name.split('/').pop()}</span>
                    <span class="shrink-0 rounded border px-1.5 py-0.5 text-[9px] font-semibold {categoryColor(repo.category)}">{repo.category}</span>
                    {#if isInstalled}
                      <span class="shrink-0 rounded bg-emerald-500/15 px-1.5 py-0.5 text-[9px] font-semibold text-emerald-300">Installed</span>
                    {/if}
                  </div>
                  <div class="text-[11px] text-slate-500">{repo.summary || ''}</div>
                </button>
                {#if isInstalled}
                  <button class="shrink-0 rounded border border-rose-500/30 px-2 py-1 text-[10px] text-rose-400 transition hover:bg-rose-500/10 hover:text-rose-300"
                    onclick={() => removeRepo(repo.name)}>Remove</button>
                {/if}
              </div>
            {/each}
          </div>
        {/if}
      </div>
    {/if}

    {#if addedToCatalog}
      <div class="shrink-0 flex items-center gap-2 border-t border-slate-800 bg-emerald-500/10 px-4 py-2 text-xs text-emerald-300">
        <CheckCircle2 size="14" />
        Installed repos added to the catalog — open the sidebar to browse by category.
      </div>
    {/if}

    {#if installing && installProgress}
      <div class="shrink-0 border-t border-slate-800 bg-slate-900/60 px-4 py-2">
        <div class="flex items-center gap-2 text-xs text-emerald-300"><Loader2 size="14" class="animate-spin" /> Installing {installProgress.name} ({installProgress.index + 1}/{installProgress.total})…</div>
        <div class="mt-1.5 h-1.5 w-full overflow-hidden rounded bg-slate-800"><div class="h-full rounded bg-emerald-500 transition-all" style="width: {((installProgress.index + 1) / installProgress.total * 100)}%"></div></div>
      </div>
    {/if}

    {#if installResults.length > 0}
      <div class="shrink-0 space-y-1 border-t border-slate-800 bg-slate-900/60 px-4 py-2">
        {#each installResults.slice(-3) as r}
          <div class="flex items-center gap-2 text-[11px]">
            {#if r.status === 'installed'}<CheckCircle2 size="11" class="text-emerald-400 shrink-0" /><span class="text-emerald-300">{r.name}</span>
            {:else if r.status === 'skipped'}<span class="text-slate-400">{r.name} — skipped</span>
            {:else}<X size="11" class="text-red-400 shrink-0" /><span class="text-red-300">{r.name}</span>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
    {/if}
  </div>
</div>
