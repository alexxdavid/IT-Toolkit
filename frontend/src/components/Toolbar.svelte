<script lang="ts">
  import { Search, X, Download, FolderOpen } from 'lucide-svelte'
  import { api } from '../lib/api'
  import {
    ui,
    selection,
    catalog,
    clearSelection,
    setQuery,
    setContentSearch,
    setExportOpen
  } from '../lib/store.svelte.ts'

  async function exportSelected() {
    setExportOpen(true)
  }

  async function openScriptsFolder() {
    const dir = await api.defaultScriptsDir()
    if (dir) await api.reveal(dir)
  }
</script>

<div class="flex min-h-14 shrink-0 flex-wrap items-center gap-2 border-b border-slate-800 bg-slate-900/60 px-4 py-2">
  <div class="relative flex-1 max-w-xl">
    <Search size="15" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
    <input
      value={ui.query}
      oninput={(e) => setQuery(e.currentTarget.value)}
      placeholder="Search repositories and scripts…"
      class="w-full rounded-md border border-slate-700 bg-slate-800/70 py-2 pl-9 pr-8 text-sm text-slate-200 placeholder:text-slate-500 focus:border-cyan-500/60 focus:outline-none"
    />
    {#if ui.query}
      <button
        class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-200"
        onclick={() => setQuery('')}
      >
        <X size="14" />
      </button>
    {/if}
  </div>

  <label class="flex cursor-pointer items-center gap-2 text-xs text-slate-400 select-none">
    <input
      type="checkbox"
      checked={ui.contentSearch}
      onchange={(e) => setContentSearch(e.currentTarget.checked)}
      class="h-3.5 w-3.5 accent-cyan-500"
    />
    Search contents
  </label>

  <button
    class="flex items-center gap-1.5 rounded-md border border-slate-700 bg-slate-800/70 px-2.5 py-1.5 text-xs text-slate-300 transition hover:border-slate-500 hover:text-slate-100"
    title="Open the Scripts folder next to this app"
    onclick={openScriptsFolder}
  >
    <FolderOpen size="14" />
  </button>

  <div class="flex-1"></div>

  <div class="text-xs text-slate-500">
    <span class="text-slate-300">{catalog.repos.length}</span> repos ·
    <span class="text-slate-300">{catalog.totalScripts}</span> scripts
  </div>

  {#if selection.size > 0}
    <button
      class="flex items-center gap-1.5 rounded-md bg-amber-500/15 px-3 py-1.5 text-xs font-medium text-amber-300 ring-1 ring-amber-500/40 transition hover:bg-amber-500/25"
      onclick={exportSelected}
    >
      <Download size="14" />
      Export {selection.size}
    </button>
    <button
      class="rounded-md border border-slate-700 px-2 py-1.5 text-xs text-slate-400 transition hover:text-slate-200"
      onclick={clearSelection}
    >
      Clear
    </button>
  {/if}
</div>
