<script lang="ts">
  import { onMount } from 'svelte'
  import { api } from '../lib/api'
  import { categoryColor, langBadge } from '../lib/categories'
  import { formatSize } from '../lib/format'
  import type { catalog } from '../../wailsjs/go/models'
  type ScriptFile = catalog.ScriptFile
  import { FileCode, Search, X, Plus, Upload, Trash2, RefreshCw, Star, Eye, Copy } from 'lucide-svelte'
  import ScriptEditor from './ScriptEditor.svelte'

  let scripts = $state<any[]>([])
  let categories = $state<string[]>([])
  let searchQuery = $state('')
  let activeCategory = $state('all')
  let loading = $state(false)

  // Editor state
  let editorOpen = $state(false)
  let editorId = $state(0)
  let editorName = $state('')

  // New script dialog
  let showNewScript = $state(false)
  let newScriptName = $state('')
  let newScriptCategory = $state('')
  let err = $state('')
  let showImportDialog = $state(false)
  let importCategory = $state('')
  let pendingImportSource = $state('')
  let scriptFavorites = $state<Record<string, boolean>>({})

  let filtered = $derived.by(() => {
    let list = scripts
    if (activeCategory !== 'all') list = list.filter((s: any) => s.category === activeCategory)
    if (searchQuery.trim()) {
      const q = searchQuery.trim().toLowerCase()
      list = list.filter((s: any) => s.name.toLowerCase().includes(q) || s.category.toLowerCase().includes(q))
    }
    return list
  })

  let categoryCounts = $derived.by(() => {
    const counts: Record<string, number> = { all: scripts.length }
    for (const s of scripts) {
      counts[s.category] = (counts[s.category] || 0) + 1
    }
    return counts
  })

  let sortedCategories = $derived(
    Object.entries(categoryCounts)
      .filter(([k]) => k !== 'all')
      .sort(([a], [b]) => a.localeCompare(b))
  )

  onMount(async () => {
    await load()
    const favs = await api.listFavorites('script') || []
    for (const f of favs) scriptFavorites[f] = true
  })

  async function load() {
    loading = true
    try {
      scripts = await api.listScriptsLib() || []
      categories = await api.libraryCategories() || []
    } catch { scripts = []; categories = [] }
    loading = false
  }

  function editScript(sc: any) {
    editorId = sc.id
    editorName = sc.name
    editorOpen = true
  }

  async function createNewScript() {
    if (!newScriptName.trim()) return
    err = ''
    try {
      const cat = newScriptCategory || 'Uncategorized'
      const id = await api.createScriptLib(newScriptName.trim(), cat)
      showNewScript = false
      newScriptName = ''
      newScriptCategory = ''
      await load()
      editorId = id
      editorName = newScriptName || 'script.ps1'
      editorOpen = true
    } catch (e) { err = String(e) }
  }

  async function importScript() {
    try {
      const source = await api.pickFile()
      if (!source) return
      pendingImportSource = source
      importCategory = activeCategory !== 'all' ? activeCategory : ''
      showImportDialog = true
    } catch (e) { err = String(e) }
  }

  async function doImport() {
    if (!pendingImportSource) return
    try {
      await api.importScriptToLib(pendingImportSource, importCategory || 'Uncategorized')
      showImportDialog = false
      pendingImportSource = ''
      await load()
    } catch (e) { err = String(e) }
  }

  async function deleteScript(sc: any) {
    if (!confirm(`Delete "${sc.name}" from the library?`)) return
    try {
      await api.deleteScriptLib(sc.id)
      await load()
    } catch (e) { err = String(e) }
  }

  async function toggleScriptFavorite(name: string) {
    const result = await api.toggleFavorite('script', name)
    scriptFavorites[name] = result
  }

  function langOf(name: string): string {
    const ext = name.split('.').pop()?.toLowerCase() || ''
    switch (ext) {
      case 'ps1': case 'psm1': case 'psd1': return 'PowerShell'
      case 'py': return 'Python'
      case 'bat': case 'cmd': return 'Batch'
      case 'vbs': return 'VBScript'
      case 'sh': return 'Shell'
      case 'sql': return 'SQL'
      default: return 'Text'
    }
  }
</script>

<div class="flex-1 min-h-0 flex overflow-hidden">
  <!-- Left panel: categories -->
  <div class="flex w-56 shrink-0 flex-col border-r border-slate-800 bg-slate-900/60">
    <div class="border-b border-slate-800 p-3">
      <button class="flex w-full items-center justify-center gap-2 rounded-lg bg-violet-500/15 px-3 py-2.5 text-xs font-bold text-violet-300 ring-1 ring-violet-500/40 transition hover:bg-violet-500/25"
        onclick={() => { showNewScript = true; err = '' }}>
        <Plus size="14" /> New Script
      </button>
    </div>
    <div class="flex-1 min-h-0 overflow-y-auto p-2 space-y-0.5">
      <button class="flex w-full items-center justify-between rounded-lg px-2.5 py-2 text-left text-sm transition {activeCategory === 'all' ? 'bg-violet-500/15 text-violet-300 font-semibold' : 'text-slate-300 hover:bg-slate-800'}"
        onclick={() => activeCategory = 'all'}>
        <span>All scripts</span>
        <span class="rounded bg-slate-800/80 px-1.5 py-0.5 text-[10px] text-slate-400">{categoryCounts.all || 0}</span>
      </button>
      {#each sortedCategories as [cat, count] (cat)}
        <button class="flex w-full items-center gap-2 rounded-lg px-2.5 py-2 text-left text-sm transition {activeCategory === cat ? 'bg-violet-500/15 text-violet-300 font-semibold' : 'text-slate-300 hover:bg-slate-800'}"
          onclick={() => activeCategory = cat}>
          <span class="truncate">{cat}</span>
          <span class="ml-auto shrink-0 rounded bg-slate-800/80 px-1.5 py-0.5 text-[10px] text-slate-400">{count}</span>
        </button>
      {/each}
    </div>
    <div class="border-t border-slate-800 p-3">
      <button class="flex w-full items-center justify-center gap-2 rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-300 hover:text-white transition"
        onclick={importScript}>
        <Upload size="14" /> Import Script
      </button>
    </div>
  </div>

  <!-- Right panel: scripts table -->
  <div class="flex min-w-0 flex-1 min-h-0 flex-col">
    <!-- Header -->
    <div class="flex shrink-0 items-center gap-3 border-b border-slate-800 bg-slate-900/60 px-4 py-2.5">
      <div class="relative flex-1 max-w-md">
        <Search size="14" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
        <input bind:value={searchQuery} placeholder="Search scripts…"
          class="w-full rounded-lg border border-white/10 bg-slate-800/60 py-2 pl-9 pr-8 text-sm text-slate-200 placeholder:text-slate-500 focus:border-violet-500/50 focus:outline-none" />
        {#if searchQuery}
          <button class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-200" onclick={() => searchQuery = ''}>
            <X size="14" />
          </button>
        {/if}
      </div>
      <div class="flex-1"></div>
      <span class="text-sm text-slate-500">{filtered.length} scripts</span>
      <button class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-sm text-slate-300 hover:text-white transition"
        onclick={load}><RefreshCw size="15" /></button>
    </div>

    <!-- Table -->
    <div class="flex-1 min-h-0 overflow-y-auto">
      {#if loading}
        <div class="flex items-center justify-center gap-3 py-16 text-slate-500">
          <div class="h-6 w-6 animate-spin rounded-full border-2 border-slate-700 border-t-violet-500"></div>
          <span class="text-sm">Loading scripts…</span>
        </div>
      {:else if filtered.length === 0}
        <div class="flex flex-col items-center gap-4 py-16 text-slate-500">
          <FileCode size="42" class="opacity-40" />
          <p class="text-sm">No scripts found.</p>
          <div class="flex gap-3">
            <button class="flex items-center gap-1.5 rounded-lg bg-violet-500/15 px-4 py-2 text-sm text-violet-300 ring-1 ring-violet-500/40 transition hover:bg-violet-500/25"
              onclick={() => showNewScript = true}><Plus size="14" /> Create New Script</button>
            <button class="flex items-center gap-1.5 rounded-lg border border-white/10 px-4 py-2 text-sm text-slate-300 hover:text-white transition"
              onclick={importScript}><Upload size="14" /> Import Script</button>
          </div>
        </div>
      {:else}
        <!-- Table header -->
        <div class="grid grid-cols-[1fr_150px_100px] gap-2 px-4 py-2 text-[11px] font-bold uppercase tracking-wider text-slate-500 border-b border-slate-800/50">
          <span>Name</span>
          <span>Category</span>
          <span class="text-right">Size</span>
        </div>
        <!-- Script rows -->
            {#each filtered as sc (sc.id)}
          <div class="group grid grid-cols-[1fr_150px_100px] gap-2 items-center px-4 py-2.5 border-b border-slate-800/30 hover:bg-white/5 transition cursor-pointer"
            role="button" tabindex="0"
            onclick={() => editScript(sc)}
            onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); editScript(sc) } }}>
            <div class="flex items-center gap-2 min-w-0">
              <FileCode size="14" class="shrink-0 text-violet-400" />
              <span class="truncate font-mono text-sm font-medium text-white hover:text-violet-300 transition">{sc.name}</span>
            </div>
            <span class="shrink-0 rounded-lg border border-white/10 bg-slate-950/60 px-2 py-0.5 text-[10px] font-semibold {categoryColor(sc.category)}">{sc.category}</span>
            <div class="flex items-center justify-end gap-1">
              <span class="text-[10px] text-slate-500 w-16 text-right">{formatSize(sc.size)}</span>
              <div class="flex gap-0.5 opacity-0 group-hover:opacity-100 ml-2">
                <button class="rounded p-1 {scriptFavorites[sc.name] ? 'text-amber-400' : 'text-slate-400 hover:text-amber-400'}"
                  title={scriptFavorites[sc.name] ? 'Remove from favorites' : 'Add to favorites'}
                  onclick={(e: Event) => { e.stopPropagation(); toggleScriptFavorite(sc.name) }}>
                  <Star size="13" class={scriptFavorites[sc.name] ? 'fill-current' : ''} />
                </button>
                <button class="rounded p-1 text-slate-400 hover:bg-violet-500/20 hover:text-violet-300" title="Edit"
                  onclick={(e) => { e.stopPropagation(); editScript(sc) }}><Eye size="13" /></button>
                <button class="rounded p-1 text-slate-400 hover:bg-slate-700 hover:text-slate-100" title="Copy"
                  onclick={(e) => { e.stopPropagation(); navigator.clipboard.writeText(sc.name) }}><Copy size="13" /></button>
                <button class="rounded p-1 text-slate-400 hover:bg-rose-500/20 hover:text-rose-400" title="Delete"
                  onclick={(e) => { e.stopPropagation(); deleteScript(sc) }}><Trash2 size="13" /></button>
              </div>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>

<!-- New Script Dialog -->
{#if showNewScript}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    role="presentation"
    onclick={() => showNewScript = false}
    onkeydown={(e) => { if (e.key === 'Escape') showNewScript = false }}>
    <div class="w-[min(500px,92vw)] rounded-xl border border-slate-700 bg-slate-900 p-5 shadow-2xl space-y-4"
      role="presentation"
      onclick={(e) => e.stopPropagation()}>
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-bold text-white">New Script</h3>
        <button class="rounded p-1 text-slate-400 hover:text-white" onclick={() => showNewScript = false}><X size="14" /></button>
      </div>
      <input bind:value={newScriptName} placeholder="script-name.ps1"
        class="w-full rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-sm text-slate-200 placeholder:text-slate-500 focus:border-violet-500/50 focus:outline-none"
        onkeydown={(e) => { if (e.key === 'Enter') createNewScript() }} />
      <select bind:value={newScriptCategory}
        class="w-full rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-sm text-slate-200 focus:border-violet-500/50 focus:outline-none">
        <option value="">Category…</option>
        {#each categories as cat (cat)}
          <option value={cat}>{cat}</option>
        {/each}
      </select>
      {#if err}<p class="text-sm text-rose-400">{err}</p>{/if}
      <button class="w-full flex items-center justify-center gap-2 rounded-lg bg-violet-500/15 px-4 py-2.5 text-sm font-bold text-violet-300 ring-1 ring-violet-500/40 transition hover:bg-violet-500/25"
        onclick={createNewScript}>
        <Plus size="14" /> Create
      </button>
    </div>
  </div>
{/if}

<!-- Import Category Dialog -->
{#if showImportDialog}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    role="presentation"
    onclick={() => showImportDialog = false}
    onkeydown={(e) => { if (e.key === 'Escape') showImportDialog = false }}>
    <div class="w-[min(400px,92vw)] rounded-xl border border-slate-700 bg-slate-900 p-5 shadow-2xl space-y-4"
      role="presentation"
      onclick={(e) => e.stopPropagation()}>
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-bold text-white">Import Script</h3>
        <button class="rounded p-1 text-slate-400 hover:text-white" onclick={() => showImportDialog = false}><X size="14" /></button>
      </div>
      <p class="text-sm text-slate-400">Which category should this script go into?</p>
      <select bind:value={importCategory}
        class="w-full rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-sm text-slate-200 focus:border-violet-500/50 focus:outline-none">
        <option value="">Uncategorized</option>
        {#each categories as cat (cat)}
          <option value={cat}>{cat}</option>
        {/each}
        <option value="New category…">New category…</option>
      </select>
      {#if importCategory === 'New category…'}
        <input bind:value={importCategory} placeholder="Category name"
          class="w-full rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-sm text-slate-200 placeholder:text-slate-500 focus:border-violet-500/50 focus:outline-none" />
      {/if}
      <div class="flex gap-2">
        <button class="flex-1 flex items-center justify-center gap-2 rounded-lg bg-violet-500/15 px-4 py-2.5 text-sm font-bold text-violet-300 ring-1 ring-violet-500/40 transition hover:bg-violet-500/25"
          onclick={doImport}>
          <Upload size="14" /> Import
        </button>
        <button class="rounded-lg border border-white/10 px-4 py-2.5 text-sm text-slate-400 hover:text-white transition"
          onclick={() => showImportDialog = false}>Cancel</button>
      </div>
    </div>
  </div>
{/if}

<!-- Script Editor -->
<ScriptEditor bind:open={editorOpen} scriptId={editorId} fileName={editorName} />
