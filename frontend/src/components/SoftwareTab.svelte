<script lang="ts">
  import { api } from '../lib/api'
  import { categoryColor } from '../lib/categories'
  import { Download, FolderOpen, Search, X, ArrowUpCircle } from 'lucide-svelte'
  import { onMount } from 'svelte'

  let software = $state<any[]>([])
  let versions = $state<Record<string, any>>({})
  let searchQuery = $state('')
  let installDest = $state('')
  let downloading = $state<string | null>(null)
  let downloadPct = $state(0)
  let completed = $state<Record<string, string>>({})
  let pollTimer: ReturnType<typeof setInterval> | null = null
  let selectedCategory = $state('all')

  let filtered = $derived.by(() => {
    const q = searchQuery.trim().toLowerCase()
    let list = software
    if (q) list = list.filter((s: any) => s.name.toLowerCase().includes(q) || s.category.toLowerCase().includes(q) || (s.notes && s.notes.toLowerCase().includes(q)))
    if (selectedCategory !== 'all') list = list.filter((s: any) => s.category === selectedCategory)
    return list
  })

  let categories = $derived.by(() => {
    const cats = new Set(software.map((s: any) => s.category))
    return ['all', ...Array.from(cats).sort()]
  })

  let updateCount = $derived(Object.values(versions).filter((v: any) => v.updateAvailable).length)

  onMount(async () => {
    software = await api.getSoftwareCatalog() || []
    try { installDest = await api.defaultSoftwareDir() } catch {}
    versions = await api.getSoftwareVersions() || {}
  })

  async function downloadItem(item: any) {
    downloading = item.name
    downloadPct = 0
    pollTimer = setInterval(async () => {
      try { const p = await api.getSoftwareProgress(); if (p?.percent > 0) downloadPct = Math.min(100, Math.round(p.percent)) } catch {}
    }, 300)
    try {
      const path = await api.downloadSoftware(item.name, item.download, installDest)
      completed[item.name] = path
    } catch (e) {
      completed[item.name] = 'error: ' + String(e)
    }
    clearInterval(pollTimer!)
    downloading = null
    downloadPct = 0
    // Refresh versions so "On disk" / "Downloaded" state is current.
    try { await api.invalidateSoftwareVersions(); versions = await api.getSoftwareVersions() || {} } catch {}
  }

  function openFolder(path: string) {
    const dir = path.substring(0, path.lastIndexOf('\\'))
    api.reveal(dir)
  }
</script>

<div class="flex h-full flex-col overflow-hidden">
  <div class="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto p-5">
    <div class="flex items-center gap-3">
      <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-sky-500 to-blue-600 text-white shadow-lg shadow-sky-500/20">
        <Download size="20" />
      </div>
      <div>
        <h2 class="text-lg font-bold text-white">Software Catalog</h2>
        <p class="text-xs text-slate-400">
          {filtered.length} apps · {installDest || 'Software folder'}
          {#if updateCount > 0}
            <span class="ml-2 inline-flex items-center gap-1 rounded-full bg-amber-500/20 px-2 py-0.5 text-[10px] font-bold text-amber-300">
              <ArrowUpCircle size="10" /> {updateCount} update{updateCount === 1 ? '' : 's'} available
            </span>
          {/if}
        </p>
      </div>
    </div>

    <div class="flex gap-2">
      <div class="relative flex-1 max-w-md">
        <Search size="14" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
        <input bind:value={searchQuery} placeholder="Filter software…"
          class="w-full rounded-lg border border-white/10 bg-slate-800/60 py-2 pl-9 pr-8 text-xs text-slate-200 placeholder:text-slate-500 focus:border-sky-500/50 focus:outline-none" />
        {#if searchQuery}
          <button class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-200" onclick={() => searchQuery = ''}><X size="12" /></button>
        {/if}
      </div>
      <div class="flex gap-1 overflow-x-auto">
        {#each categories as cat (cat)}
          <button class="shrink-0 rounded-lg border px-2.5 py-1.5 text-[10px] font-medium transition {selectedCategory === cat ? 'border-sky-500/50 bg-sky-500/15 text-sky-200' : 'border-white/10 bg-slate-800/60 text-slate-400 hover:text-slate-200'}"
            onclick={() => selectedCategory = cat}>
            {cat === 'all' ? 'All' : cat}
          </button>
        {/each}
      </div>
    </div>

    <div class="space-y-1">
      {#each filtered as item (item.name)}
        {@const sv = versions[item.name]}
        <div class="flex items-center gap-3 rounded-xl border border-white/5 bg-slate-800/50 px-4 py-3 transition hover:border-white/10">
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <span class="text-xs font-bold text-white">{item.name}</span>
              <span class="shrink-0 rounded-lg border border-white/10 bg-slate-950/60 px-2 py-0.5 text-[9px] font-medium {categoryColor(item.category)}">{item.category}</span>
              {#if sv?.installedVersion}
                <span class="shrink-0 rounded bg-slate-700/50 px-1.5 py-0.5 text-[9px] font-mono text-slate-400">Installed: v{sv.installedVersion}</span>
              {/if}
              {#if sv?.latestVersion}
                <span class="shrink-0 rounded border border-sky-500/30 px-1.5 py-0.5 text-[9px] font-mono text-sky-300">Latest: v{sv.latestVersion}</span>
              {/if}
              {#if sv?.updateAvailable}
                <span class="shrink-0 rounded bg-amber-500/15 px-1.5 py-0.5 text-[9px] font-bold text-amber-300">Update available</span>
              {/if}
              {#if completed[item.name] && !completed[item.name].startsWith('error')}
                <span class="shrink-0 rounded bg-emerald-500/15 px-1.5 py-0.5 text-[9px] font-semibold text-emerald-300">Downloaded</span>
              {:else if sv?.hasDownload}
                <span class="shrink-0 rounded bg-emerald-500/10 px-1.5 py-0.5 text-[9px] text-emerald-400">On disk</span>
              {/if}
            </div>
            {#if item.notes}
              <div class="mt-0.5 text-[11px] text-slate-500">{item.notes}</div>
            {/if}
            {#if completed[item.name] && completed[item.name].startsWith('error')}
              <div class="text-[11px] text-rose-400">{completed[item.name]}</div>
            {/if}
          </div>
          {#if downloading === item.name}
            <div class="flex items-center gap-2 shrink-0">
              <div class="h-1.5 w-20 overflow-hidden rounded bg-slate-800">
                <div class="h-full rounded bg-sky-500 transition-all" style="width: {downloadPct}%"></div>
              </div>
              <span class="text-[10px] font-mono text-sky-400">{downloadPct}%</span>
            </div>
          {:else if completed[item.name] && !completed[item.name].startsWith('error')}
            <button class="flex shrink-0 items-center gap-1 rounded-lg border border-white/10 px-2.5 py-1.5 text-[10px] text-slate-300 hover:text-white transition"
              onclick={() => openFolder(completed[item.name])}>
              <FolderOpen size="12" /> Open
            </button>
          {:else}
            <button class="flex shrink-0 items-center gap-1.5 rounded-lg bg-sky-500/15 px-3 py-1.5 text-[10px] font-bold text-sky-300 ring-1 ring-sky-500/40 transition hover:bg-sky-500/25 disabled:opacity-50"
              onclick={() => downloadItem(item)} disabled={downloading !== null}>
              <Download size="12" /> {sv?.updateAvailable ? 'Update' : 'Download'}
            </button>
          {/if}
        </div>
      {/each}
    </div>
  </div>
</div>
