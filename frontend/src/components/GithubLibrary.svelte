<script lang="ts">
  import { onMount } from 'svelte'
  import { Download, Loader2, CheckCircle2, X } from 'lucide-svelte'
  import { api, pickFolder, onInstallProgress, onInstallDone } from '../lib/api'
  import type { InstallProgress } from '../lib/api'
  import { ui } from '../lib/store.svelte.ts'
  import SoftwareTab from './SoftwareTab.svelte'
  import ScriptsTab from './ScriptsTab.svelte'
  import LocalTab from './LocalTab.svelte'
  import FavoritesTab from './FavoritesTab.svelte'
  import { Star } from 'lucide-svelte'

  type CategoryGroup = { name: string; repos: any[] }

  let recommended = $state<any[]>([])
  let installDest = $state('')
  let selected = $state<Record<string, boolean>>({})
  let installing = $state(false)
  let installResults = $state<any[]>([])
  let installProgress = $state<InstallProgress | null>(null)
  let searchQuery = $state('')
  let selectedCategory = $state('all')
  let offProgress: (() => void) | null = null
  let offDone: (() => void) | null = null
  let repoFavorites = $state<Record<string, boolean>>({})

  let categories = $derived.by((): CategoryGroup[] => {
    const q = searchQuery.trim().toLowerCase()
    const filtered = q
      ? recommended.filter((r: any) => r.name.toLowerCase().includes(q) || r.category.toLowerCase().includes(q) || (r.summary && r.summary.toLowerCase().includes(q)))
      : recommended
    const map = new Map<string, any[]>()
    for (const r of filtered) {
      const cat = r.category || 'Other'
      if (!map.has(cat)) map.set(cat, [])
      map.get(cat)!.push(r)
    }
    return Array.from(map.entries()).sort((a, b) => a[0].localeCompare(b[0])).map(([name, repos]) => ({ name, repos }))
  })

  let visibleRepos = $derived(
    selectedCategory === 'all'
      ? categories.flatMap((c: CategoryGroup) => c.repos)
      : categories.find((c: CategoryGroup) => c.name === selectedCategory)?.repos || []
  )

  onMount(async () => {
    recommended = await api.getRecommendedReposCombined() || []
    try { installDest = await api.defaultInstallDir() } catch {}
    const favs = await api.listFavorites('repo') || []
    for (const f of favs) repoFavorites[f] = true
  })

  async function toggleFavorite(name: string) {
    const result = await api.toggleFavorite('repo', name)
    repoFavorites[name] = result
  }

  onDestroy(() => { offProgress?.(); offDone?.() })

  function onDestroy(fn: () => void) {
    onMount(fn)
  }

  function toggleRepo(name: string) { selected[name] = !selected[name] }

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
    offProgress = onInstallProgress((p: InstallProgress) => { installProgress = p })
    offDone = onInstallDone(() => { installing = false; offProgress?.(); offDone?.(); offProgress = null; offDone = null })
    try {
      installResults = await api.installRepos(names, installDest)
      if (installResults.some((r: any) => r.status === 'installed')) {
        try { await api.addFolder(installDest) } catch {}
      }
    } catch (e) { installResults = [{ name: 'error', status: 'failed', message: String(e) }] }
    installing = false; offProgress?.(); offDone?.(); offProgress = null; offDone = null
  }
</script>

<div class="flex min-h-0 flex-1 flex-col overflow-hidden">
  <!-- Repos Tab -->
  {#if ui.localTab === 'repos'}
    <div class="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto p-5">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-bold text-white">Recommended Repos</h2>
        <span class="text-xs text-slate-400">{recommended.length} available</span>
      </div>

      <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
        <input bind:value={searchQuery} placeholder="Filter repos…"
          class="flex-1 rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-xs text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none" />
        <select bind:value={selectedCategory}
          class="rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-xs text-slate-200 focus:border-emerald-500/50 focus:outline-none">
          <option value="all">All categories</option>
          {#each categories as cat (cat.name)}
            <option value={cat.name}>{cat.name} ({cat.repos.length})</option>
          {/each}
        </select>
        <button class="rounded-lg bg-emerald-500/15 px-3 py-2 text-xs font-bold text-emerald-300 ring-1 ring-emerald-500/40 hover:bg-emerald-500/25"
          onclick={toggleAll}>
          {visibleRepos.every((r: any) => selected[r.name]) ? 'Deselect all' : 'Select all'}
        </button>
      </div>

      {#if visibleRepos.length === 0}
        <div class="py-16 text-center text-sm text-slate-500">No repos match.</div>
      {:else}
        <div class="space-y-1">
          {#each visibleRepos as repo (repo.name)}
            <div class="flex w-full items-center gap-3 rounded-xl border border-white/5 bg-slate-800/50 px-4 py-3 transition hover:border-white/10 {selected[repo.name] ? 'ring-1 ring-emerald-500/40 border-emerald-500/30 bg-emerald-500/5' : ''}">
              <input type="checkbox" checked={!!selected[repo.name]} class="h-3.5 w-3.5 shrink-0 accent-emerald-500" onclick={() => toggleRepo(repo.name)} />
              <button class="min-w-0 flex-1 text-left" onclick={() => toggleRepo(repo.name)}>
                <span class="font-mono text-xs font-semibold text-white">{repo.name.split('/').pop()}</span>
                <span class="ml-2 rounded border px-1.5 py-0.5 text-[9px] text-slate-400">{repo.category}</span>
                {#if repo.summary}<div class="text-[11px] text-slate-500">{repo.summary}</div>{/if}
              </button>
              <button class="shrink-0 rounded p-1.5 transition {repoFavorites[repo.name] ? 'text-amber-400' : 'text-slate-500 hover:text-amber-400'}"
                title={repoFavorites[repo.name] ? 'Remove from favorites' : 'Add to favorites'}
                onclick={() => toggleFavorite(repo.name)}>
                <Star size="14" class={repoFavorites[repo.name] ? 'fill-current' : ''} />
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>
    <!-- Install bar -->
    <div class="shrink-0 border-t border-slate-800 bg-slate-900/60 px-5 py-3 flex items-center gap-3">
      <button class="flex items-center justify-center gap-2 rounded-lg px-4 py-2 text-xs font-bold transition shadow-lg disabled:opacity-40 {selectedCount() > 0 ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-slate-950 shadow-emerald-500/20 hover:scale-105' : 'bg-slate-800 text-slate-400 cursor-default'}"
        onclick={doInstall} disabled={installing || selectedCount() === 0}>
        {#if installing}<Loader2 size="14" class="animate-spin" /> Installing…
        {:else if selectedCount() > 0}<Download size="14" /> Install {selectedCount()} repo{selectedCount() === 1 ? '' : 's'}
        {:else}Select repos to install
        {/if}
      </button>
      <input value={installDest} oninput={(e) => installDest = e.currentTarget.value} placeholder="Install to…"
        class="flex-1 rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-[11px] text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none" />
      {#if installing && installProgress}
        <div class="flex items-center gap-2"><Loader2 size="12" class="animate-spin text-emerald-400" /><span class="text-[10px] text-emerald-300">{installProgress.index + 1}/{installProgress.total}</span></div>
      {/if}
    </div>

    <!-- Install results -->
    {#if installResults.length > 0}
      <div class="shrink-0 border-t border-slate-800 bg-slate-900/60 px-5 py-2">
        {#each installResults.slice(-5) as r}
          <div class="flex items-center gap-2 text-[11px]">
            {#if r.status === 'installed'}<CheckCircle2 size="11" class="text-emerald-400 shrink-0" /><span class="text-emerald-300">{r.name}</span>
            {:else if r.status === 'skipped'}<span class="text-slate-500">{r.name} — skipped</span>
            {:else}<X size="11" class="text-red-400 shrink-0" /><span class="text-red-300">{r.name}</span><span class="text-slate-500">{r.message}</span>{/if}
          </div>
        {/each}
      </div>
    {/if}

    {:else if ui.localTab === 'software'}
      <SoftwareTab />

    {:else if ui.localTab === 'scripts'}
      <ScriptsTab />

    {:else if ui.localTab === 'local'}
    <LocalTab />
  {:else if ui.localTab === 'favorites'}
    <FavoritesTab />
  {/if}
</div>
