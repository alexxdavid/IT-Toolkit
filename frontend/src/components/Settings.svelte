<script lang="ts">
  import { onMount } from 'svelte'
  import { api } from '../lib/api'
  import { categoryColor } from '../lib/categories'
  import { CheckCircle2, Plus, Trash2, Pencil, RefreshCw, Download, Loader2, X } from 'lucide-svelte'
  import { setSettingsOpen, ui } from '../lib/store.svelte.ts'
  import { ArrowLeft } from 'lucide-svelte'

  let customCategories = $state<any[]>([])
  let customRepos = $state<any[]>([])
  let recCategories = $state<any[]>([])
  let newCatName = $state('')
  let newRepoName = $state('')
  let newRepoUrl = $state('')
  let newRepoCategory = $state('')
  let newRepoSummary = $state('')
  let newSwName = $state('')
  let newSwDownload = $state('')
  let newSwCategory = $state('')
  let newSwWinget = $state('')
  let newSwNotes = $state('')
  let customSoftware = $state<any[]>([])
  let scriptsPath = $state('')
  let editingCat = $state<number | null>(null)
  let editCatName = $state('')
  let err = $state<string | null>(null)
  let success = $state<string | null>(null)
  let version = $state('1.0.0')
  let build = $state(1)
  let updateInfo = $state<any>(null)
  let checking = $state(false)
  let downloading = $state(false)
  let downloadPct = $state(0)
  let pollTimer: ReturnType<typeof setInterval> | null = null

  const recCategoryNames = [
    'Microsoft 365', 'Windows Administration', 'Intune & Endpoint', 'Microsoft Defender',
    'Active Directory', 'Azure & Cloud', 'Network & Infrastructure', 'Security & Monitoring',
    'GRC & Compliance', 'Monitoring & SIEM', 'DevOps & Labs', 'Databases',
    'AI & Machine Learning', 'Utilities & Tools'
  ]

  onMount(async () => {
    await refresh()
    try { version = await api.getCurrentVersion() } catch {}
    try { build = await api.getCurrentBuild() } catch {}
    try { } catch {}
  })

  let allCategories = $derived.by((): string[] => {
    const cats = new Set<string>(recCategoryNames)
    customCategories.forEach((c: any) => cats.add(c.name))
    return Array.from(cats).sort()
  })

  async function refresh() {
    customCategories = await api.listCustomCategories() || []
    customRepos = await api.listCustomRepos() || []
    customSoftware = await api.listCustomSoftware() || []
    recCategories = await api.getRecommendedRepos() || []
  }

  function msg(t: string) { success = t; setTimeout(() => success = null, 3000) }
  function errMsg(t: string) { err = t; setTimeout(() => err = null, 5000) }

  async function addCategory() {
    if (!newCatName.trim()) return
    try {
      await api.createCustomCategory(newCatName.trim())
      newCatName = ''
      await refresh()
      msg('Category created')
    } catch (e) { errMsg(String(e)) }
  }

  async function renameCategory(id: number, name: string) {
    try {
      await api.renameCustomCategory(id, name)
      editingCat = null
      await refresh()
      msg('Category renamed')
    } catch (e) { errMsg(String(e)) }
  }

  async function deleteCategory(id: number) {
    if (!confirm('Delete this category? Repos in it will become uncategorized.')) return
    try {
      await api.deleteCustomCategory(id)
      await refresh()
      msg('Category deleted')
    } catch (e) { errMsg(String(e)) }
  }

  async function addRepo() {
    if (!newRepoName.trim() || !newRepoUrl.trim()) return
    const cat = newRepoCategory.trim() || 'Other'
    try {
      await api.addCustomRepo(newRepoName.trim(), newRepoUrl.trim(), cat, newRepoSummary.trim())
      newRepoName = ''; newRepoUrl = ''; newRepoCategory = ''; newRepoSummary = ''
      await refresh()
      msg('Repo added')
    } catch (e) { errMsg(String(e)) }
  }

  async function removeRepo(id: number) {
    if (!confirm('Remove this repo?')) return
    try {
      await api.removeCustomRepo(id)
      await refresh()
      msg('Repo removed')
    } catch (e) { errMsg(String(e)) }
  }

  async function addSoftware() {
    if (!newSwName.trim() || !newSwDownload.trim()) return
    try {
      await api.addCustomSoftware(newSwName.trim(), '', newSwCategory.trim() || 'Other', newSwDownload.trim(), newSwNotes.trim(), newSwWinget.trim())
      newSwName = ''; newSwDownload = ''; newSwCategory = ''; newSwWinget = ''; newSwNotes = ''
      await refresh()
      msg('Software added')
    } catch (e) { errMsg(String(e)) }
  }

  async function removeSoftware(id: number) {
    if (!confirm('Remove this software?')) return
    try {
      await api.removeCustomSoftware(id)
      await refresh()
      msg('Software removed')
    } catch (e) { errMsg(String(e)) }
  }

  async function checkUpdate() {
    checking = true
    try {
      const info = await api.checkForUpdate(true)
      updateInfo = info
      if (info?.available) msg(`Update available: v${info.version}`)
      else if (info?.stale) errMsg('Could not reach the update server — showing cached info')
      else msg('Already on latest version')
    } catch (e) {
      errMsg('Update check failed: ' + String(e))
    }
    checking = false
  }

  async function installUpdate() {
    if (!updateInfo) return
    downloading = true; downloadPct = 0
    pollTimer = setInterval(async () => {
      try { const p = await api.getUpdateProgress(); if (p?.percent > 0) downloadPct = Math.min(100, Math.round(p.percent)) } catch {}
    }, 300)
    try {
      const path = await api.downloadUpdate(updateInfo.installer_url, updateInfo.version)
      clearInterval(pollTimer!)
      if (path) { downloadPct = 100; await api.applyUpdate(path) }
    } catch { clearInterval(pollTimer!) }
    downloading = false
  }

  function getAllCategories(): string[] {
    const cats = new Set(recCategoryNames)
    customCategories.forEach((c: any) => cats.add(c.name))
    return Array.from(cats).sort()
  }
</script>

<div class="flex min-h-0 flex-1 flex-col overflow-y-auto p-6">
  <h2 class="mb-6 text-lg font-bold text-white">Settings</h2>

  {#if success}
    <div class="mb-4 flex items-center gap-2 rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-xs text-emerald-300"><CheckCircle2 size="14" /> {success}</div>
  {/if}
  {#if err}
    <div class="mb-4 flex items-center gap-2 rounded-lg border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-xs text-rose-300"><X size="14" /> {err}</div>
  {/if}

  <!-- Update Section -->
  <section class="mb-6 rounded-2xl border border-white/10 bg-gradient-to-br from-slate-800/80 to-slate-900/80 p-5 shadow-xl">
    <h3 class="mb-3 text-sm font-bold text-white">App Updates</h3>
    <div class="flex items-center gap-4">
      <span class="text-xs text-slate-400">Current version: <span class="font-mono text-white">v{version}</span> <span class="font-mono text-slate-500">(build {build})</span></span>
      {#if checking}
        <span class="text-xs text-slate-400"><Loader2 size="12" class="inline animate-spin" /> Checking…</span>
      {:else if updateInfo?.available}
        <span class="text-xs text-emerald-400">v{updateInfo.version} available</span>
        <button class="rounded-lg bg-emerald-500/15 px-3 py-1.5 text-xs font-bold text-emerald-300 ring-1 ring-emerald-500/40 hover:bg-emerald-500/25"
          onclick={installUpdate} disabled={downloading}>
          {#if downloading}<Loader2 size="12" class="inline animate-spin" /> {downloadPct}%{:else}<Download size="12" class="inline" /> Install{/if}
        </button>
      {/if}
      <button class="flex items-center gap-1.5 rounded-lg border border-white/10 bg-slate-800/60 px-3 py-1.5 text-xs text-slate-300 hover:text-white transition"
        onclick={checkUpdate} disabled={checking}>
        <RefreshCw size="12" class={checking ? 'animate-spin' : ''} /> Check for updates
      </button>
    </div>
  </section>

  <!-- Categories Section -->
  <section class="mb-6 rounded-2xl border border-white/10 bg-gradient-to-br from-slate-800/80 to-slate-900/80 p-5 shadow-xl">
    <h3 class="mb-3 text-sm font-bold text-white">Categories</h3>
    <div class="mb-4 flex gap-2">
      <input bind:value={newCatName} placeholder="New category name…" class="flex-1 rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none"
        onkeydown={(e) => { if (e.key === 'Enter') addCategory() }} />
      <button class="flex items-center gap-1.5 rounded-lg bg-emerald-500/15 px-3 py-2 text-xs font-bold text-emerald-300 ring-1 ring-emerald-500/40 hover:bg-emerald-500/25"
        onclick={addCategory}><Plus size="12" /> Add</button>
    </div>
    <div class="space-y-1.5">
      {#each customCategories as cat (cat.id)}
        <div class="flex items-center gap-2 rounded-lg border border-white/5 bg-slate-800/40 px-3 py-2">
          {#if editingCat === cat.id}
            <input bind:value={editCatName} class="flex-1 rounded border border-white/10 bg-slate-900/60 px-2 py-1 text-xs text-white focus:border-emerald-500/50 focus:outline-none"
              onkeydown={(e) => { if (e.key === 'Enter') renameCategory(cat.id, editCatName) }} />
            <button class="rounded bg-emerald-500/20 px-2 py-1 text-[10px] text-emerald-300" onclick={() => renameCategory(cat.id, editCatName)}>Save</button>
          {:else}
            <span class="shrink-0 rounded border px-2 py-0.5 text-[10px] font-semibold {categoryColor(cat.name)}">{cat.name}</span>
            <span class="flex-1 text-[10px] text-slate-500">{customRepos.filter((r: any) => r.category === cat.name).length} repos · {customSoftware.filter((s: any) => s.category === cat.name).length} software</span>
            <button class="rounded p-1 text-slate-400 hover:text-white" onclick={() => { editingCat = cat.id; editCatName = cat.name }}><Pencil size="12" /></button>
            <button class="rounded p-1 text-slate-400 hover:text-rose-400" onclick={() => deleteCategory(cat.id)}><Trash2 size="12" /></button>
          {/if}
        </div>
      {/each}
      {#if customCategories.length === 0}
        <p class="text-xs text-slate-600">No custom categories yet. Add one above.</p>
      {/if}
    </div>
  </section>

  <!-- Custom Repos Section -->
  <section class="mb-6 rounded-2xl border border-white/10 bg-gradient-to-br from-slate-800/80 to-slate-900/80 p-5 shadow-xl">
    <h3 class="mb-3 text-sm font-bold text-white">Add Custom Repo</h3>
    <div class="grid grid-cols-2 gap-2 mb-3">
      <input bind:value={newRepoName} placeholder="Repo name (e.g. MyTool)" class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none" />
      <input bind:value={newRepoUrl} placeholder="GitHub URL" class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none" />
    </div>
    <div class="grid grid-cols-2 gap-2 mb-3">
      <select bind:value={newRepoCategory} class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 focus:border-emerald-500/50 focus:outline-none">
        <option value="">Select category…</option>
        {#each getAllCategories() as cat (cat)}
          <option value={cat}>{cat}</option>
        {/each}
      </select>
      <input bind:value={newRepoSummary} placeholder="Short description" class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none" />
    </div>
    <button class="flex items-center gap-1.5 rounded-lg bg-emerald-500/15 px-4 py-2 text-xs font-bold text-emerald-300 ring-1 ring-emerald-500/40 hover:bg-emerald-500/25"
      onclick={addRepo}><Plus size="12" /> Add Repo</button>

    {#if customRepos.length > 0}
      <h4 class="mt-4 mb-2 text-xs font-semibold text-slate-400">Custom repos ({customRepos.length})</h4>
      <div class="space-y-1">
        {#each customRepos as repo (repo.id)}
          <div class="flex items-center gap-2 rounded-lg border border-white/5 bg-slate-800/40 px-3 py-2">
            <span class="min-w-0 truncate font-mono text-xs font-semibold text-white">{repo.name}</span>
            <span class="shrink-0 rounded border px-1.5 py-0.5 text-[9px] {categoryColor(repo.category)}">{repo.category}</span>
            <span class="min-w-0 flex-1 truncate text-[10px] text-slate-500">{repo.summary}</span>
            <button class="shrink-0 rounded p-1 text-slate-400 hover:text-rose-400" onclick={() => removeRepo(repo.id)}><Trash2 size="12" /></button>
          </div>
        {/each}
      </div>
    {/if}
  </section>

  <!-- Custom Software Section -->
  <section class="mb-6 rounded-2xl border border-white/10 bg-gradient-to-br from-slate-800/80 to-slate-900/80 p-5 shadow-xl">
    <h3 class="mb-3 text-sm font-bold text-white">Add Custom Software</h3>
    <div class="grid grid-cols-2 gap-2 mb-2">
      <input bind:value={newSwName} placeholder="Software name" class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 placeholder:text-slate-500 focus:border-sky-500/50 focus:outline-none" />
      <input bind:value={newSwDownload} placeholder="Download URL" class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 placeholder:text-slate-500 focus:border-sky-500/50 focus:outline-none" />
    </div>
    <div class="grid grid-cols-3 gap-2 mb-3">
      <select bind:value={newSwCategory} class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 focus:border-sky-500/50 focus:outline-none">
        <option value="">Category…</option>
        {#each allCategories as cat (cat)}
          <option value={cat}>{cat}</option>
        {/each}
      </select>
      <input bind:value={newSwWinget} placeholder="Winget ID (optional)" class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 placeholder:text-slate-500 focus:border-sky-500/50 focus:outline-none" />
      <input bind:value={newSwNotes} placeholder="Notes" class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-xs text-slate-200 placeholder:text-slate-500 focus:border-sky-500/50 focus:outline-none" />
    </div>
    <button class="flex items-center gap-1.5 rounded-lg bg-sky-500/15 px-4 py-2 text-xs font-bold text-sky-300 ring-1 ring-sky-500/40 hover:bg-sky-500/25"
      onclick={addSoftware}><Plus size="12" /> Add Software</button>

    {#if customSoftware.length > 0}
      <h4 class="mt-4 mb-2 text-xs font-semibold text-slate-400">Custom software ({customSoftware.length})</h4>
      <div class="space-y-1">
        {#each customSoftware as sw (sw.id)}
          <div class="flex items-center gap-2 rounded-lg border border-white/5 bg-slate-800/40 px-3 py-2">
            <span class="min-w-0 truncate text-xs font-semibold text-white">{sw.name}</span>
            <span class="shrink-0 rounded border px-1.5 py-0.5 text-[9px] {categoryColor(sw.category)}">{sw.category}</span>
            {#if sw.wingetId}
              <span class="shrink-0 rounded bg-slate-700/50 px-1.5 py-0.5 text-[9px] font-mono text-slate-400">winget:{sw.wingetId}</span>
            {/if}
            <span class="min-w-0 flex-1 truncate text-[10px] text-slate-500">{sw.download}</span>
            <button class="shrink-0 rounded p-1 text-slate-400 hover:text-rose-400" onclick={() => removeSoftware(sw.id)}><Trash2 size="12" /></button>
          </div>
        {/each}
      </div>
    {/if}
  </section>

</div>
