<script lang="ts">
  import { api } from '../lib/api'
  import { Folder, FileText, Download, ChevronRight, ChevronDown, ExternalLink } from 'lucide-svelte'
  import { onMount } from 'svelte'
  import { ui } from '../lib/store.svelte'

  let repos = $state<any[]>([])
  let software = $state<any[]>([])
  let expandedRepo = $state<string | null>(null)
  let repoFiles = $state<any[]>([])
  let activeTab = $state<'repos' | 'software'>('repos')
  let searchQuery = $state('')

  let filteredRepos = $derived(
    searchQuery.trim()
      ? repos.filter((r: any) => r.name.toLowerCase().includes(searchQuery.toLowerCase()))
      : repos
  )

  let filteredSoftware = $derived(
    searchQuery.trim()
      ? software.filter((s: any) => s.name.toLowerCase().includes(searchQuery.toLowerCase()))
      : software
  )

  async function load() {
    repos = await api.listLocalRepos() || []
    software = await api.listLocalSoftware() || []
    // Filter out non-repo dirs (like Desktop, etc.)
    repos = repos.filter((r: any) => r.isDir)
  }

  onMount(load)

  // Reload whenever the Local tab is opened (files may change on disk).
  $effect(() => {
    if (ui.localTab === 'local') load()
  })

  async function expandRepo(repo: any) {
    if (expandedRepo === repo.name) {
      expandedRepo = null
      repoFiles = []
      return
    }
    expandedRepo = repo.name
    repoFiles = await api.getRepoContents(repo.path) || []
  }

  function formatSize(bytes: number): string {
    if (!bytes) return '0 B'
    if (bytes >= 1048576) return (bytes / 1048576).toFixed(1) + ' MB'
    if (bytes >= 1024) return (bytes / 1024).toFixed(0) + ' KB'
    return bytes + ' B'
  }
</script>

<div class="flex h-full flex-col overflow-hidden">
  <div class="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto p-5">
    <div class="flex items-center gap-3">
      <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-emerald-500 to-teal-600 text-white shadow-lg shadow-emerald-500/20">
        <Folder size="20" />
      </div>
      <div class="flex-1">
        <h2 class="text-lg font-bold text-white">Downloaded</h2>
        <p class="text-xs text-slate-400">{repos.length} repos · {software.length} software installed</p>
      </div>
    </div>

    <!-- Toggle: Repos / Software -->
    <div class="flex gap-2 rounded-xl border border-white/10 bg-slate-800/60 p-1">
      <button class="flex-1 flex items-center justify-center gap-2 rounded-lg px-4 py-2 text-xs font-bold transition {activeTab === 'repos' ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-slate-950 shadow-lg' : 'text-slate-400 hover:text-slate-200'}"
        onclick={() => activeTab = 'repos'}>
        <Folder size="14" /> Repos ({filteredRepos.length})
      </button>
      <button class="flex-1 flex items-center justify-center gap-2 rounded-lg px-4 py-2 text-xs font-bold transition {activeTab === 'software' ? 'bg-gradient-to-r from-sky-500 to-blue-500 text-slate-950 shadow-lg' : 'text-slate-400 hover:text-slate-200'}"
        onclick={() => activeTab = 'software'}>
        <Download size="14" /> Software ({filteredSoftware.length})
      </button>
    </div>

    <input bind:value={searchQuery} placeholder="Filter…"
      class="rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-xs text-slate-200 placeholder:text-slate-500 focus:border-emerald-500/50 focus:outline-none" />

    {#if activeTab === 'repos'}
      {#if filteredRepos.length === 0}
        <div class="flex flex-col items-center gap-3 py-16 text-slate-600">
          <Folder size="36" class="opacity-40" />
          <p class="text-sm">No repos installed yet. Use the Repos tab to download some.</p>
        </div>
      {:else}
        <div class="space-y-1">
          {#each filteredRepos as repo (repo.name)}
            <div class="rounded-xl border border-white/5 bg-slate-800/50 overflow-hidden transition hover:border-white/10">
              <div class="flex items-center gap-3 px-4 py-3 cursor-pointer transition hover:bg-white/5"
                role="button" tabindex="0"
                onclick={() => expandRepo(repo)}
                onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); expandRepo(repo) } }}>
                {#if expandedRepo === repo.name}
                  <ChevronDown size="16" class="shrink-0 text-slate-400" />
                {:else}
                  <ChevronRight size="16" class="shrink-0 text-slate-400" />
                {/if}
                <Folder size="14" class="shrink-0 text-amber-400" />
                <span class="font-mono text-xs font-semibold text-white">{repo.name}</span>
                <div class="flex-1"></div>
                <button class="shrink-0 rounded-md border border-white/10 px-2 py-1 text-[10px] text-slate-400 hover:text-white transition"
                  onclick={(e: MouseEvent) => { e.stopPropagation(); api.reveal(repo.path) }}>
                  <ExternalLink size="10" /> Open
                </button>
              </div>
              {#if expandedRepo === repo.name}
                <div class="border-t border-white/5 divide-y divide-white/5">
                  {#if repoFiles.length === 0}
                    <div class="px-4 py-3 text-xs text-slate-500">Empty folder</div>
                  {:else}
                    {#each repoFiles as file (file.name)}
                      <div class="flex items-center gap-3 px-4 py-2 pl-10">
                        {#if file.isDir}
                          <Folder size="12" class="shrink-0 text-amber-400" />
                        {:else}
                          <FileText size="12" class="shrink-0 text-slate-500" />
                        {/if}
                        <span class="font-mono text-xs text-slate-200">{file.name}</span>
                        <span class="text-[10px] text-slate-500">{formatSize(file.size)}</span>
                      </div>
                    {/each}
                  {/if}
                </div>
              {/if}
            </div>
          {/each}
        </div>
      {/if}
    {:else}
      {#if filteredSoftware.length === 0}
        <div class="flex flex-col items-center gap-3 py-16 text-slate-600">
          <Download size="36" class="opacity-40" />
          <p class="text-sm">No software downloaded yet. Use the Software tab to download some.</p>
        </div>
      {:else}
        <div class="space-y-1">
          {#each filteredSoftware as sw (sw.name)}
            <div class="flex items-center gap-3 rounded-xl border border-white/5 bg-slate-800/50 px-4 py-3 transition hover:border-white/10">
              <FileText size="14" class="shrink-0 text-sky-400" />
              <span class="font-mono text-xs font-semibold text-white">{sw.name}</span>
              <span class="text-[10px] text-slate-500">{formatSize(sw.size)}</span>
              <div class="flex-1"></div>
              <button class="shrink-0 rounded-md border border-white/10 px-2 py-1 text-[10px] text-slate-400 hover:text-white transition"
                onclick={() => api.reveal(sw.path)}>
                <ExternalLink size="10" /> Open
              </button>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
</div>
