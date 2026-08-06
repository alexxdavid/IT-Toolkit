<script lang="ts">
  import { FolderOpen, Loader2, X, CheckCircle2 } from 'lucide-svelte'
  import { api, pickFolder } from '../lib/api'
  import type { ExportResult } from '../lib/types'
  import { selection, clearSelection, setExportOpen, ui } from '../lib/store.svelte.ts'

  let dest = $state('')
  let busy = $state(false)
  let result: ExportResult | null = $state(null)

  const count = $derived(selection.size)

  async function browse() {
    const d = await pickFolder()
    if (d) dest = d
  }

  async function doExport() {
    if (!dest || count === 0) return
    busy = true
    result = null
    try {
      result = await api.exportScripts(Array.from(selection), dest)
      if (result && result.errors.length === 0) {
        clearSelection()
      }
    } catch (e) {
      result = { copied: 0, skipped: 0, errors: [String(e)] }
    } finally {
      busy = false
    }
  }

  function close() {
    if (busy) return
    setExportOpen(false)
    result = null
    dest = ''
  }
</script>

{#if ui.exportOpen}
  <div
    class="fixed inset-0 z-40 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    role="presentation"
    onclick={close}
    onkeydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        close()
      }
    }}
  >
    <div
      class="w-[min(620px,92vw)] overflow-hidden rounded-xl border border-slate-700 bg-slate-900 shadow-2xl"
      role="presentation"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
    >
      <div class="flex items-center justify-between border-b border-slate-800 px-5 py-3">
        <h2 class="text-sm font-semibold text-slate-100">Export scripts</h2>
        <button class="rounded p-1 text-slate-400 hover:bg-slate-800 hover:text-slate-200" onclick={close}>
          <X size="16" />
        </button>
      </div>

      <div class="p-5">
        {#if result}
          <div class="space-y-3">
            <div class="flex items-center gap-2 text-sm text-emerald-300">
              <CheckCircle2 size="18" />
              Export finished: {result.copied} copied
              {#if result.skipped}· {result.skipped} skipped{/if}
            </div>
            {#if result.errors.length > 0}
              <div class="max-h-40 overflow-y-auto rounded-md border border-rose-500/30 bg-rose-500/10 p-3">
                {#each result.errors as e (e)}
                  <div class="font-mono text-[11px] text-rose-300">{e}</div>
                {/each}
              </div>
            {/if}
            <div class="text-xs text-slate-400">
              Files were copied to <code class="text-slate-200">{dest}</code> preserving repository folders.
            </div>
            <button
              class="rounded-md border border-slate-700 px-4 py-2 text-xs font-medium text-slate-200 transition hover:border-slate-500"
              onclick={close}
            >
              Done
            </button>
          </div>
        {:else}
          <p class="mb-4 text-sm text-slate-400">
            {count} script{count === 1 ? '' : 's'} selected. Choose a destination folder — repository
            subfolders will be preserved.
          </p>
          <div class="flex items-center gap-2">
            <input
              bind:value={dest}
              placeholder="Destination folder…"
              class="min-w-0 flex-1 rounded-md border border-slate-700 bg-slate-800/70 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500 focus:border-cyan-500/60 focus:outline-none"
            />
            <button
              class="flex shrink-0 items-center gap-1.5 rounded-md border border-slate-700 bg-slate-800/70 px-3 py-2 text-xs text-slate-300 transition hover:border-slate-500"
              onclick={browse}
            >
              <FolderOpen size="14" /> Browse
            </button>
          </div>
          <button
            class="mt-5 flex w-full items-center justify-center gap-2 rounded-md bg-amber-500/15 px-4 py-2.5 text-sm font-medium text-amber-300 ring-1 ring-amber-500/40 transition hover:bg-amber-500/25 disabled:opacity-50"
            disabled={!dest || busy}
            onclick={doExport}
          >
            {#if busy}<Loader2 size="15" class="animate-spin" /> Exporting…{:else}Export to folder{/if}
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}
