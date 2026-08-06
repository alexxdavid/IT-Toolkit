<script lang="ts">
  import { marked } from 'marked'
  import { onMount } from 'svelte'
  import { Copy, Play, FolderOpen, X, FileCode2, BookOpen } from 'lucide-svelte'
  import { api } from '../lib/api'
  import { viewerState, closeViewer } from '../lib/viewer.svelte'
  import { highlight } from '../lib/hljs'
  import { langBadge } from '../lib/categories'
  import { formatSize } from '../lib/format'

  let content = $state('')
  let currentPath = $state('')
  let loading = $state(false)
  let actionMsg = $state<string | null>(null)

  const detail = $derived(viewerState.repoDetail)
  const script = $derived(viewerState.script)
  const readmeHtml = $derived(
    detail?.repo?.readmeText ? marked.parse(detail.repo.readmeText) : ''
  )
  const hl = $derived(content ? highlight(script?.lang ?? 'Text', content) : '')

  $effect(() => {
    if (viewerState.open && viewerState.tab === 'source' && script && script.absPath !== currentPath) {
      currentPath = script.absPath
      content = ''
      loading = true
      api
        .scriptContent(script.absPath)
        .then((c) => {
          if (currentPath === script.absPath) {
            content = c
            loading = false
          }
        })
        .catch(() => {
          if (currentPath === script.absPath) loading = false
        })
    }
  })

  $effect(() => {
    if (!viewerState.open) {
      content = ''
      currentPath = ''
      actionMsg = null
    }
  })

  onMount(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') closeViewer()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  })

  async function copy() {
    if (!script) return
    try {
      const text = content || (await api.copyScript(script.absPath))
      await navigator.clipboard.writeText(text)
      actionMsg = 'Copied to clipboard'
    } catch {
      actionMsg = 'Copy failed'
    }
    setTimeout(() => (actionMsg = null), 2000)
  }

  async function run() {
    if (!script) return
    if (!confirm(`Run ${script.name} in an elevated window?\n\nThis executes the script on this machine.`)) return
    try {
      const res = await api.runScript(script.absPath)
      actionMsg = res.launched ? `Launched via ${res.message}` : res.message
    } catch (e) {
      actionMsg = String(e)
    }
    setTimeout(() => (actionMsg = null), 3000)
  }

  async function reveal() {
    if (!script) return
    await api.reveal(script.absPath)
  }

  const title = $derived(script?.relPath ?? detail?.repo?.name ?? '')
</script>

{#if viewerState.open}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    role="presentation"
    onclick={closeViewer}
    onkeydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        closeViewer()
      }
    }}
  >
    <div
      class="flex h-[85vh] w-[min(1100px,92vw)] flex-col overflow-hidden rounded-xl border border-slate-700 bg-slate-900 shadow-2xl"
      role="presentation"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
    >
      <div class="flex shrink-0 items-center gap-3 border-b border-slate-800 px-4 py-3">
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <span class="truncate font-mono text-sm text-slate-100">{title}</span>
            {#if script}
              <span class="shrink-0 rounded border px-1.5 py-0.5 text-[10px] {langBadge(script.lang)}"
                >{script.lang}</span
              >
            {/if}
          </div>
          <div class="mt-0.5 text-[11px] text-slate-500">
            {#if script}
              {script.repo} · {formatSize(script.size)}
            {:else}
              {detail?.repo?.category} · {detail?.repo?.scriptCount} scripts
            {/if}
          </div>        </div>
        {#if detail}
          <div class="flex rounded-lg border border-slate-700 bg-slate-800/60 p-0.5">
            <button
              class="flex items-center gap-1.5 rounded-md px-2.5 py-1 text-xs transition {viewerState.tab === 'source'
                ? 'bg-slate-700 text-slate-100'
                : 'text-slate-400 hover:text-slate-200'}"
              onclick={() => (viewerState.tab = 'source')}
            >
              <FileCode2 size="13" /> Scripts
            </button>
            <button
              class="flex items-center gap-1.5 rounded-md px-2.5 py-1 text-xs transition {viewerState.tab === 'readme'
                ? 'bg-slate-700 text-slate-100'
                : 'text-slate-400 hover:text-slate-200'}"
              onclick={() => (viewerState.tab = 'readme')}
            >
              <BookOpen size="13" /> README
            </button>
          </div>
        {/if}
        {#if script}
          <button
            class="flex items-center gap-1.5 rounded-md border border-slate-700 bg-slate-800/70 px-2.5 py-1.5 text-xs text-slate-300 transition hover:border-slate-500 hover:text-slate-100"
            onclick={reveal}
          >
            <FolderOpen size="14" /> Reveal
          </button>
          <button
            class="flex items-center gap-1.5 rounded-md border border-slate-700 bg-slate-800/70 px-2.5 py-1.5 text-xs text-slate-300 transition hover:border-slate-500 hover:text-slate-100"
            onclick={copy}
          >
            <Copy size="14" /> Copy
          </button>
          <button
            class="flex items-center gap-1.5 rounded-md bg-emerald-500/15 px-3 py-1.5 text-xs font-medium text-emerald-300 ring-1 ring-emerald-500/40 transition hover:bg-emerald-500/25"
            onclick={run}
          >
            <Play size="14" /> Run
          </button>
        {/if}
        <button class="rounded p-1.5 text-slate-400 hover:bg-slate-800 hover:text-slate-200" onclick={closeViewer}>
          <X size="16" />
        </button>
      </div>

      <div class="min-h-0 flex-1 overflow-auto bg-[#0d1424]">
        {#if viewerState.tab === 'source'}
          {#if script}
            {#if loading}
              <div class="p-6 text-sm text-slate-500">Loading script…</div>
            {:else}
              <pre class="p-4 font-mono text-[12.5px] leading-relaxed text-slate-200"><code class="hljs">{@html hl}</code></pre>
            {/if}
          {:else}
            <div class="p-6 text-sm text-slate-500">Select a script to view its source.</div>
          {/if}
        {:else if detail}
          {#if detail?.repo?.readmeText}
            <div class="markdown-body p-6">{@html readmeHtml}</div>
          {:else}
            <div class="p-6 text-sm text-slate-500">No README found for this repository.</div>
          {/if}
        {/if}
      </div>

      {#if actionMsg}
        <div class="shrink-0 border-t border-slate-800 bg-slate-900 px-4 py-2 text-xs text-cyan-300">
          {actionMsg}
        </div>
      {/if}
    </div>
  </div>
{/if}
