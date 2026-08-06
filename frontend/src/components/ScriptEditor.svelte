<script lang="ts">
  import { api } from '../lib/api'
  import { highlight } from '../lib/hljs'
  import { X, Save, Copy, FileCode } from 'lucide-svelte'

  let { open = $bindable(false), scriptId = 0, fileName = '' } = $props()

  let content = $state('')
  let originalContent = $state('')
  let loading = $state(false)
  let saving = $state(false)
  let saved = $state(false)
  let lang = $state('PowerShell')

  $effect(() => {
    if (open && scriptId) loadContent()
  })

  async function loadContent() {
    loading = true
    try {
      content = await api.getScriptLibContent(scriptId)
      originalContent = content
      const ext = fileName.split('.').pop()?.toLowerCase() || ''
      if (ext === 'py') lang = 'Python'
      else if (ext === 'bat' || ext === 'cmd') lang = 'Batch'
      else if (ext === 'sh') lang = 'Shell'
      else if (ext === 'sql') lang = 'SQL'
      else lang = 'PowerShell'
    } catch {
      content = ''
      originalContent = ''
    }
    loading = false
  }

  async function save() {
    saving = true
    try {
      await api.saveScriptLib(scriptId, content)
      originalContent = content
      saved = true
      setTimeout(() => saved = false, 2000)
    } catch (e) {
      alert('Save failed: ' + String(e))
    }
    saving = false
  }

  async function copyToClipboard() {
    await navigator.clipboard.writeText(content)
  }

  let edited = $derived(content !== originalContent)
  let hl = $derived(highlight(lang, content || ''))
</script>

{#if open}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    role="presentation"
    onclick={() => open = false}
    onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') open = false }}>
    <div class="flex h-[85vh] w-[min(900px,92vw)] flex-col overflow-hidden rounded-xl border border-slate-700 bg-slate-900 shadow-2xl"
      role="presentation"
      onclick={(e) => e.stopPropagation()}>
      <div class="flex items-center gap-3 border-b border-slate-800 px-5 py-3">
        <FileCode size="18" class="text-violet-400" />
        <span class="font-mono text-sm text-white">{fileName}</span>
        {#if edited}
          <span class="rounded bg-amber-500/15 px-2 py-0.5 text-xs font-medium text-amber-300">Modified</span>
        {/if}
        {#if saved}
          <span class="rounded bg-emerald-500/15 px-2 py-0.5 text-xs font-medium text-emerald-300">Saved</span>
        {/if}
        <div class="flex-1"></div>
        <button class="flex items-center gap-1.5 rounded-lg border border-white/10 bg-slate-800/60 px-3 py-1.5 text-sm text-slate-300 hover:text-white transition"
          onclick={copyToClipboard}>
          <Copy size="14" /> Copy
        </button>
        <button class="flex items-center gap-1.5 rounded-lg bg-emerald-500/15 px-3 py-1.5 text-sm font-medium text-emerald-300 ring-1 ring-emerald-500/40 transition hover:bg-emerald-500/25 disabled:opacity-50"
          onclick={save} disabled={saving || !edited}>
          <Save size="14" /> Save
        </button>
        <button class="rounded p-1.5 text-slate-400 hover:text-white" onclick={() => open = false}>
          <X size="16" />
        </button>
      </div>
      <div class="flex-1 overflow-auto">
        {#if loading}
          <div class="flex items-center justify-center py-16 text-slate-500">Loading…</div>
        {:else}
          <textarea
            bind:value={content}
            class="w-full h-full min-h-[50vh] resize-none bg-transparent p-5 font-mono text-sm leading-relaxed text-slate-200 outline-none"
            spellcheck="false"
            placeholder="Start typing your script…"
          ></textarea>
        {/if}
      </div>
      <div class="flex items-center justify-between border-t border-slate-800 px-5 py-2 text-xs text-slate-500">
        <span>Stored in library database</span>
        <span>{content.split('\n').length} lines · {content.length} chars</span>
      </div>
    </div>
  </div>
{/if}
