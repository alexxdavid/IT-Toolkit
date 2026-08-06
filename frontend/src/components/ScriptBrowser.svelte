<script lang="ts">
  import { Copy, Play, FolderOpen, Eye, Loader2, Inbox } from 'lucide-svelte'
  import { api } from '../lib/api'
  import type { ScriptFile } from '../lib/types'
  import { langBadge } from '../lib/categories'
  import { formatSize, formatDate } from '../lib/format'
  import { selection, toggleSelection } from '../lib/store.svelte.ts'
  import { openScript } from '../lib/viewer.svelte'

  let { scripts, loading = false, title = 'Scripts' }: { scripts: ScriptFile[]; loading?: boolean; title?: string } = $props()

  function isSelected(p: string): boolean {
    return selection.has(p)
  }

  function toggleAll() {
    const selected = scripts.filter((s) => isSelected(s.absPath)).length
    const select = selected !== scripts.length
    for (const s of scripts) {
      if (select && !isSelected(s.absPath)) selection.add(s.absPath)
      if (!select) selection.delete(s.absPath)
    }
  }

  async function copyScript(sc: ScriptFile) {
    const text = await api.copyScript(sc.absPath)
    await navigator.clipboard.writeText(text)
  }

  async function runScript(sc: ScriptFile) {
    if (!confirm(`Run ${sc.name} in an elevated window?\n\nThis executes the script on this machine.`)) return
    await api.runScript(sc.absPath)
  }
</script>

<div class="min-h-0 flex-1 overflow-y-auto p-4">
  <div class="mx-auto max-w-5xl">
    <div class="mb-3 flex items-center justify-between">
      <h2 class="text-sm font-semibold text-slate-300">{title}</h2>
      {#if scripts.length > 0}
        <button
          class="rounded-md border border-slate-700 px-2.5 py-1 text-[11px] text-slate-300 transition hover:border-slate-500"
          onclick={toggleAll}
        >
          {scripts.filter((s) => isSelected(s.absPath)).length > 0 ? 'Clear selection' : 'Select all'}
        </button>
      {/if}
    </div>

    {#if loading}
      <div class="flex items-center gap-2 py-10 text-sm text-slate-500">
        <Loader2 size="16" class="animate-spin" /> Searching…
      </div>
    {:else if scripts.length === 0}
      <div class="flex flex-col items-center gap-2 py-16 text-slate-600">
        <Inbox size="32" />
        <p class="text-sm">No matching scripts.</p>
      </div>
    {:else}
      <div class="overflow-hidden rounded-xl border border-slate-800">
        <table class="w-full text-left text-xs">
          <thead class="bg-slate-900 text-[11px] uppercase tracking-wider text-slate-500">
            <tr>
              <th class="w-8 px-3 py-2"></th>
              <th class="px-3 py-2">Name</th>
              <th class="px-3 py-2">Repository</th>
              <th class="px-3 py-2">Language</th>
              <th class="px-3 py-2">Size</th>
              <th class="px-3 py-2">Modified</th>
              <th class="w-28 px-3 py-2 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/60 bg-slate-900/40">
            {#each scripts as sc (sc.id)}
              <tr class="group transition hover:bg-slate-800/40">
                <td class="px-3 py-2">
                  <input
                    type="checkbox"
                    checked={isSelected(sc.absPath)}
                    onclick={() => toggleSelection(sc.absPath)}
                    class="h-3.5 w-3.5 accent-cyan-500"
                  />
                </td>
                <td class="max-w-[320px] px-3 py-2">
                  <button class="block w-full truncate text-left font-mono text-slate-200 hover:text-cyan-300" onclick={() => openScript(sc)}>
                    {sc.name}
                  </button>
                  <div class="truncate text-[10px] text-slate-600">{sc.relPath}</div>
                </td>
                <td class="px-3 py-2">
                  <span class="text-slate-400">{sc.repo}</span>
                </td>
                <td class="px-3 py-2">
                  <span class="rounded border px-1.5 py-0.5 text-[9px] {langBadge(sc.lang)}">{sc.lang}</span>
                </td>
                <td class="px-3 py-2 text-slate-400">{formatSize(sc.size)}</td>
                <td class="px-3 py-2 text-slate-500">{formatDate(sc.mtime)}</td>
                <td class="px-3 py-2">
                  <div class="flex justify-end gap-1 opacity-0 transition group-hover:opacity-100">
                    <button class="rounded p-1 text-slate-400 hover:bg-slate-700 hover:text-slate-100" title="View" onclick={() => openScript(sc)}><Eye size="13" /></button>
                    <button class="rounded p-1 text-slate-400 hover:bg-slate-700 hover:text-slate-100" title="Copy" onclick={() => copyScript(sc)}><Copy size="13" /></button>
                    {#if sc.isScript}
                      <button class="rounded p-1 text-emerald-400 hover:bg-slate-700" title="Run elevated" onclick={() => runScript(sc)}><Play size="13" /></button>
                    {/if}
                    <button class="rounded p-1 text-slate-400 hover:bg-slate-700 hover:text-slate-100" title="Reveal in Explorer" onclick={() => api.reveal(sc.absPath)}><FolderOpen size="13" /></button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
</div>
