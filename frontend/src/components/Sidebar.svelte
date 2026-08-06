<script lang="ts">
  import { RefreshCw, FolderPlus, FolderX, GitBranch, Settings } from 'lucide-svelte'
  import { api, pickFolder } from '../lib/api'
  import {
    ui,
    catalog,
    setCategory,
    setFolder,
    setTypeFilter,
    setLibraryOpen,
    setSettingsOpen
  } from '../lib/store.svelte.ts'
  import { typeChips } from '../lib/categories'

  let err = $state<string | null>(null)

  async function addFolder() {
    err = null
    const path = await pickFolder()
    if (!path) return
    try {
      await api.addFolder(path)
    } catch (e) {
      err = String(e)
    }
  }

  async function removeFolder(id: number) {
    if (!confirm('Remove this folder from the catalog? (files on disk are untouched)')) return
    await api.removeFolder(id)
    if (ui.selectedFolder === id) setFolder('all')
  }

  async function rescan() {
    await api.rescan()
  }

  const progressPct = $derived(
    ui.scanStatus && ui.scanStatus.total > 0
      ? Math.round((ui.scanStatus.current / ui.scanStatus.total) * 100)
      : 0
  )
</script>

<aside class="flex h-full w-48 shrink-0 flex-col border-r border-slate-800 bg-slate-900/60 sm:w-56 lg:w-64">
  <div class="flex h-14 items-center gap-2 border-b border-slate-800 px-4">
    <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-cyan-500 to-blue-600 font-bold text-white">
      IT
    </div>
    <div class="leading-tight">
      <div class="text-sm font-semibold text-slate-100">Solutions IT Toolkit</div>
      <div class="text-[10px] text-slate-500">Script console</div>
    </div>
  </div>

  <div class="flex-1 space-y-4 overflow-y-auto px-3 py-3">
    <div>
      <div class="mb-1.5 flex items-center justify-between px-1">
        <span class="text-[11px] font-semibold uppercase tracking-wider text-slate-500">Folders</span>
        <button
          class="rounded p-1 text-slate-400 hover:bg-slate-800 hover:text-slate-200"
          title="Add scripts folder"
          onclick={addFolder}
        >
          <FolderPlus size="14" />
        </button>
      </div>
      {#if catalog.folders.length === 0}
        <p class="px-1 text-xs text-slate-600">No folders added yet.</p>
      {/if}
      <ul class="space-y-0.5">
        {#each catalog.folders as f (f.id)}
          <li class="group flex items-center">
            <button
              class="flex min-w-0 flex-1 items-center gap-2 rounded-md px-2 py-1.5 text-left text-xs transition
                     {ui.selectedFolder === f.id
                ? 'bg-cyan-500/15 text-cyan-200'
                : 'text-slate-300 hover:bg-slate-800'}"
              onclick={() => setFolder(f.id === ui.selectedFolder ? 'all' : f.id)}
            >
              <span class="truncate">{f.displayName}</span>
              <span class="ml-auto shrink-0 rounded bg-slate-800 px-1.5 py-0.5 text-[10px] text-slate-400"
                >{f.scriptCount}</span
              >
            </button>
            <button
              class="hidden shrink-0 rounded p-1 text-slate-500 hover:text-rose-400 group-hover:block"
              title="Remove folder from catalog"
              onclick={() => removeFolder(f.id)}
            >
              <FolderX size="13" />
            </button>
          </li>
        {/each}
      </ul>
    </div>

    <div>
      <div class="mb-1.5 px-1 text-[11px] font-semibold uppercase tracking-wider text-slate-500">
        Categories
      </div>
      <ul class="space-y-0.5">
        <li>
          <button
            class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-xs transition
                   {ui.selectedCategory === 'all'
              ? 'bg-cyan-500/15 text-cyan-200'
              : 'text-slate-300 hover:bg-slate-800'}"
            onclick={() => setCategory('all')}
          >
            All repositories
            <span class="ml-auto rounded bg-slate-800 px-1.5 py-0.5 text-[10px] text-slate-400"
              >{catalog.repos.length}</span
            >
          </button>
        </li>
        {#each catalog.categories as c (c.name)}
          <li>
            <button
              class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-xs transition
                     {ui.selectedCategory === c.name
                ? 'bg-cyan-500/15 text-cyan-200'
                : 'text-slate-300 hover:bg-slate-800'}"
              onclick={() => setCategory(c.name)}
            >
              <span class="truncate">{c.name}</span>
              <span class="ml-auto rounded bg-slate-800 px-1.5 py-0.5 text-[10px] text-slate-400"
                >{c.count}</span
              >
            </button>
          </li>
        {/each}
      </ul>
    </div>

    <div>
      <div class="mb-1.5 px-1 text-[11px] font-semibold uppercase tracking-wider text-slate-500">
        File types
      </div>
      <div class="flex flex-wrap gap-1 px-1">
        {#each typeChips as t (t)}
          <button
            class="rounded-md border px-2 py-1 text-[11px] transition
                   {ui.typeFilter === t
              ? 'border-cyan-500/50 bg-cyan-500/15 text-cyan-200'
              : 'border-slate-700 bg-slate-800/60 text-slate-400 hover:text-slate-200'}"
            onclick={() => setTypeFilter(t)}
          >
            {t === 'all' ? 'All' : t}
          </button>
        {/each}
    </div>

    <div>
      <button
        class="flex w-full items-center gap-2 rounded-lg px-2 py-2 text-xs font-medium transition
               {ui.libraryOpen
          ? 'bg-gradient-to-r from-emerald-500/15 to-teal-500/15 text-emerald-300 border border-emerald-500/30'
          : 'text-slate-300 hover:bg-slate-800 border border-transparent'}"
        onclick={() => { setLibraryOpen(!ui.libraryOpen); if (!ui.libraryOpen) { ui.selectedCategory = 'all'; ui.typeFilter = 'all' } }}
      >
        <GitBranch size="14" />
        GitHub Library
      </button>
    </div>

    <div>
      <button
        class="flex w-full items-center gap-2 rounded-lg px-2 py-2 text-xs font-medium transition
               {ui.settingsOpen
          ? 'bg-gradient-to-r from-amber-500/15 to-orange-500/15 text-amber-300 border border-amber-500/30'
          : 'text-slate-300 hover:bg-slate-800 border border-transparent'}"
        onclick={() => { setSettingsOpen(!ui.settingsOpen); if (ui.settingsOpen) { setLibraryOpen(false) } }}
      >
        <Settings size="14" />
        Settings
      </button>
    </div>
  </div>
  </div>

  <div class="border-t border-slate-800 p-3">
    {#if err}<p class="mb-2 text-xs text-rose-400">{err}</p>{/if}
    <button
      class="flex w-full items-center justify-center gap-2 rounded-md border border-slate-700 bg-slate-800/80 px-3 py-2 text-xs font-medium text-slate-200 transition hover:border-slate-500 hover:bg-slate-700 disabled:opacity-50"
      onclick={rescan}
      disabled={ui.scanning}
    >
      <RefreshCw size="14" class={ui.scanning ? 'animate-spin' : ''} />
      {ui.scanning ? 'Scanning…' : 'Rescan folders'}
    </button>
    {#if ui.scanning && ui.scanStatus}
      <div class="mt-2">
        <div class="h-1 w-full overflow-hidden rounded bg-slate-800">
          <div
            class="h-full bg-cyan-500 transition-all"
            style="width: {progressPct}%"
          ></div>
        </div>
        <div class="mt-1 text-[10px] text-slate-500">
          {ui.scanStatus.phase} {ui.scanStatus.current}/{ui.scanStatus.total}
        </div>
      </div>
    {/if}
  </div>
</aside>
