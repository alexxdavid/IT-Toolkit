<script lang="ts">
  import { ChevronDown, ChevronRight, Copy, Play, FolderOpen, Eye, BookOpen, Loader2 } from 'lucide-svelte'
  import { api } from '../lib/api'
  import type { Repo, RepoDetail, ScriptFile } from '../lib/types'
  import { categoryColor, langBadge } from '../lib/categories'
  import { formatSize, formatDate } from '../lib/format'
  import { selection, toggleSelection, setDetail, getDetail } from '../lib/store.svelte.ts'
  import { openScript, openRepoReadme } from '../lib/viewer.svelte'

  let { repo }: { repo: Repo } = $props()

  let loading = $state(false)
  let err = $state<string | null>(null)
  let expanded = $state(false)

  const detail = $derived(getDetail(repo.id))

  async function toggle() {
    expanded = !expanded
    if (expanded && !detail) {
      loading = true
      err = null
      try {
        const d = await api.repoDetail(repo.id)
        setDetail(repo.id, d)
      } catch (e) {
        err = String(e)
      } finally {
        loading = false
      }
    }
  }

  function isSelected(p: string): boolean {
    return selection.has(p)
  }

  async function copyScript(sc: ScriptFile) {
    const text = await api.copyScript(sc.absPath)
    await navigator.clipboard.writeText(text)
  }

  async function runScript(sc: ScriptFile) {
    if (!confirm(`Run ${sc.name} in an elevated window?\n\nThis executes the script on this machine.`)) return
    await api.runScript(sc.absPath)
  }

  function toggleAll(scripts: ScriptFile[]) {
    const selected = scripts.filter((s) => isSelected(s.absPath)).length
    const select = selected !== scripts.length
    for (const s of scripts) {
      if (select && !isSelected(s.absPath)) selection.add(s.absPath)
      if (!select) selection.delete(s.absPath)
    }
  }
</script>

<div class="rounded-xl border border-slate-800 bg-slate-900/50 transition hover:border-slate-700">
  <div
    role="button"
    tabindex="0"
    class="flex w-full cursor-pointer items-center gap-3 px-4 py-3 text-left"
    onclick={toggle}
    onkeydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault()
        toggle()
      }
    }}
  >
    {#if expanded}
      <ChevronDown size="16" class="shrink-0 text-slate-500" />
    {:else}
      <ChevronRight size="16" class="shrink-0 text-slate-500" />
    {/if}
    <div class="min-w-0 flex-1">
      <div class="flex items-center gap-2">
        <span class="truncate font-mono text-sm font-semibold text-slate-100">{repo.name}</span>
        <span class="shrink-0 rounded border px-1.5 py-0.5 text-[10px] {categoryColor(repo.category)}"
          >{repo.category}</span
        >
      </div>
      <div class="mt-0.5 text-[11px] text-slate-500">
        {repo.scriptCount} scripts
      </div>
    </div>
    {#if repo.readmePath}
      <button
        class="flex shrink-0 items-center gap-1 rounded-md border border-slate-700 px-2 py-1 text-[11px] text-slate-300 transition hover:border-slate-500"
        title="View README"
        onclick={(e) => {
          e.stopPropagation()
          openRepoReadme({ repo, scripts: [] } as unknown as RepoDetail)
        }}
      >
        <BookOpen size="12" /> README
      </button>
    {/if}
    {#if detail && detail.scripts.length > 0}
      <button
        class="shrink-0 rounded-md border border-slate-700 px-2 py-1 text-[11px] text-slate-300 transition hover:border-slate-500"
        onclick={(e) => {
          e.stopPropagation()
          toggleAll(detail.scripts)
        }}
      >
        {detail.scripts.filter((s) => isSelected(s.absPath)).length > 0 ? 'Clear' : 'Select all'}
      </button>
    {/if}
  </div>

  {#if expanded}
    <div class="border-t border-slate-800">
      {#if loading}
        <div class="flex items-center gap-2 px-4 py-4 text-xs text-slate-500">
          <Loader2 size="14" class="animate-spin" /> Loading scripts…
        </div>
      {:else if err}
        <div class="px-4 py-4 text-xs text-rose-400">{err}</div>
      {:else if detail}
        {#if detail.scripts.length === 0}
          <div class="px-4 py-4 text-xs text-slate-500">No script files indexed in this repository.</div>
        {:else}
          <ul class="divide-y divide-slate-800/60">
            {#each detail.scripts as sc (sc.id)}
              <li class="group flex items-center gap-3 px-4 py-2 transition hover:bg-slate-800/40">
                <input
                  type="checkbox"
                  checked={isSelected(sc.absPath)}
                  onclick={() => toggleSelection(sc.absPath)}
                  class="h-3.5 w-3.5 shrink-0 accent-cyan-500"
                />
                <button class="min-w-0 flex-1 text-left" onclick={() => openScript(sc)}>
                  <div class="flex items-center gap-2">
                    <span class="truncate font-mono text-xs text-slate-200">{sc.relPath}</span>
                    {#if sc.isScript}
                      <span class="shrink-0 rounded border px-1.5 py-0.5 text-[9px] {langBadge(sc.lang)}"
                        >{sc.lang}</span
                      >
                    {:else}
                      <span class="shrink-0 rounded border px-1.5 py-0.5 text-[9px] {langBadge(sc.lang)}"
                        >doc</span
                      >
                    {/if}
                  </div>
                  <div class="mt-0.5 text-[10px] text-slate-600">
                    {formatSize(sc.size)} · {formatDate(sc.mtime)}
                  </div>
                </button>
                <div class="flex shrink-0 items-center gap-1 opacity-0 transition group-hover:opacity-100">
                  <button class="rounded p-1 text-slate-400 hover:bg-slate-700 hover:text-slate-100" title="View"
                    onclick={() => openScript(sc)}><Eye size="13" /></button>
                  <button class="rounded p-1 text-slate-400 hover:bg-slate-700 hover:text-slate-100" title="Copy"
                    onclick={() => copyScript(sc)}><Copy size="13" /></button>
                  {#if sc.isScript}
                    <button class="rounded p-1 text-emerald-400 hover:bg-slate-700" title="Run elevated"
                      onclick={() => runScript(sc)}><Play size="13" /></button>
                  {/if}
                  <button class="rounded p-1 text-slate-400 hover:bg-slate-700 hover:text-slate-100" title="Reveal in Explorer"
                    onclick={() => api.reveal(sc.absPath)}><FolderOpen size="13" /></button>
                </div>
              </li>
            {/each}
          </ul>
        {/if}
      {/if}
    </div>
  {/if}
</div>
