<script lang="ts">
  import { onMount } from 'svelte'
  import { api, pickFolder } from '../lib/api'
  import { categoryColor, langBadge } from '../lib/categories'
  import { formatSize } from '../lib/format'
  import type { catalog } from '../../wailsjs/go/models'
  type ScriptFile = catalog.ScriptFile
  import { Star, Download, Copy, CheckCircle2, FolderOpen, Loader2, X } from 'lucide-svelte'
  import { ui } from '../lib/store.svelte'

  let favoriteRepos = $state<any[]>([])
  let favoriteSoftware = $state<any[]>([])
  let favoriteScripts = $state<any[]>([])
  let loading = $state(false)
  let usbPath = $state('')
  let copyResults = $state<{name: string; status: string; message: string}[]>([])
  let copying = $state(false)

  let totalFavorites = $derived(favoriteRepos.length + favoriteSoftware.length + favoriteScripts.length)

  onMount(async () => {
    usbPath = await api.defaultInstallDir() + '\\Favorites'
    await load()
  })

  async function load() {
    loading = true
    try {
      // Get favorite names
      const repoNames = await api.listFavorites('repo') || []
      const swNames = await api.listFavorites('software') || []
      const scriptNames = await api.listFavorites('script') || []

      // Match to full data
      const allRepos = await api.getRecommendedReposCombined() || []
      const allSoftware = await api.getSoftwareCatalog() || []
      const allScripts = await api.listScriptsLib() || []

      favoriteRepos = allRepos.filter((r: any) => repoNames.includes(r.name))
      favoriteSoftware = allSoftware.filter((s: any) => swNames.includes(s.name))
      favoriteScripts = allScripts.filter((s: any) => scriptNames.includes(s.name))
    } catch {
      favoriteRepos = []
      favoriteSoftware = []
      favoriteScripts = []
    }
    loading = false
  }

  async function copyItemToUSB(item: any, kind: string) {
    if (!usbPath) { alert('Set a USB/favorites path first'); return }
    try {
      if (kind === 'repo') {
        // Repos are folders on disk — try the Repo directory
        const repoPath = (await api.defaultInstallDir()) + '\\' + item.name.split('/').pop()
        const exists = await api.isRepoInstalled(item.name, await api.defaultInstallDir())
        if (exists) {
          await api.copyToFavorites(repoPath, usbPath)
          alert(`Copied ${item.name} to USB`)
        } else {
          alert(`"${item.name}" must be installed first before copying to USB`)
        }
      } else if (kind === 'software') {
        // Software is in the Software directory — copy installer
        const swFiles = await api.listLocalSoftware() || []
        const match = swFiles.find((f: any) => f.name.toLowerCase().includes(item.name.split(' ')[0].toLowerCase()))
        if (match) {
          await api.copyToFavorites(match.path, usbPath)
          alert(`Copied ${item.name} installer to USB`)
        } else {
          alert(`"${item.name}" not downloaded yet — download first, then copy`)
        }
      } else if (kind === 'script') {
        const content = await api.getScriptLibContent(item.id)
        const scriptsDir = usbPath + '\\Scripts'
        await api.createScriptLib(item.name, item.category)
        copyResults = [...copyResults, { name: item.name, status: 'ok', message: 'copied to DB' }]
      }
    } catch (e) {
      alert('Copy failed: ' + String(e))
    }
  }

  async function copyAllToUSB() {
    if (!usbPath) { alert('Set a USB/favorites path first'); return }
    if (!confirm(`Copy all ${totalFavorites} favorited items to ${usbPath}?`)) return
    copying = true
    copyResults = []
    for (const r of favoriteRepos) {
      try {
        const exists = await api.isRepoInstalled(r.name, await api.defaultInstallDir())
        if (exists) {
          const path = (await api.defaultInstallDir()) + '\\' + (r.name.split('/').pop() || r.name)
          await api.copyToFavorites(path, usbPath + '\\Repos')
          copyResults = [...copyResults, { name: r.name, status: 'ok', message: 'copied' }]
        } else {
          copyResults = [...copyResults, { name: r.name, status: 'skip', message: 'not installed locally' }]
        }
      } catch (e) {
        copyResults = [...copyResults, { name: r.name, status: 'error', message: String(e) }]
      }
    }
    for (const s of favoriteSoftware) {
      try {
        const files = await api.listLocalSoftware() || []
        const match = files.find((f: any) => f.name.toLowerCase().includes(s.name.toLowerCase().split(' ')[0]))
        if (match) {
          await api.copyToFavorites(match.path, usbPath + '\\Software')
          copyResults = [...copyResults, { name: s.name, status: 'ok', message: 'copied' }]
        } else {
          copyResults = [...copyResults, { name: s.name, status: 'skip', message: 'not downloaded' }]
        }
      } catch (e) {
        copyResults = [...copyResults, { name: s.name, status: 'error', message: String(e) }]
      }
    }
    // Scripts: create in DB
    for (const sc of favoriteScripts) {
      try {
        await api.createScriptLib(sc.name, sc.category)
        copyResults = [...copyResults, { name: sc.name, status: 'ok', message: 'created' }]
      } catch (e) {
        copyResults = [...copyResults, { name: sc.name, status: 'error', message: String(e) }]
      }
    }
    copying = false
  }

  async function unstar(kind: string, name: string) {
    await api.toggleFavorite(kind, name)
    await load()
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

<div class="flex h-full min-h-0 flex-col overflow-hidden">
  <div class="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto p-5">
    <!-- Header -->
    <div class="flex items-center gap-3">
      <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-amber-500 to-orange-600 text-white shadow-lg shadow-amber-500/20">
        <Star size="20" />
      </div>
      <div class="flex-1">
        <h2 class="text-lg font-bold text-white">Favorites</h2>
        <p class="text-sm text-slate-400">
          {favoriteRepos.length} repos · {favoriteSoftware.length} software · {favoriteScripts.length} scripts
        </p>
      </div>
    </div>

    <!-- USB destination -->
    <div class="flex items-center gap-2">
      <span class="text-sm text-slate-400">Copy to:</span>
      <input bind:value={usbPath} placeholder="C:\USB\Favorites"
        class="flex-1 rounded-lg border border-white/10 bg-slate-800/60 py-2 px-3 text-sm text-slate-200 placeholder:text-slate-500 focus:border-amber-500/50 focus:outline-none" />
      <button class="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-sm text-slate-300 hover:text-white transition"
        onclick={async () => { const p = await pickFolder(); if (p) usbPath = p }}>
        <FolderOpen size="14" />
      </button>
    </div>

    {#if totalFavorites === 0}
      <div class="flex flex-col items-center gap-4 py-16 text-slate-500">
        <Star size="42" class="opacity-40" />
        <p class="text-sm">No favorites yet. Star repos, software, or scripts from their tabs to add them here.</p>
      </div>
    {:else}
      <button class="flex w-full items-center justify-center gap-2 rounded-lg bg-amber-500/15 px-4 py-2.5 text-sm font-bold text-amber-300 ring-1 ring-amber-500/40 transition hover:bg-amber-500/25 disabled:opacity-50"
        onclick={copyAllToUSB} disabled={copying || !usbPath}>
        {#if copying}
          <Loader2 size="14" class="animate-spin" /> Copying…
        {:else}
          <Copy size="14" /> Copy All ({totalFavorites} items) to USB
        {/if}
      </button>
    {/if}

    <!-- Copy results -->
    {#if copyResults.length > 0}
      <div class="rounded-xl border border-white/10 bg-slate-800/50 p-3 space-y-1">
        <div class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Copy Results</div>
        {#each copyResults as r}
          <div class="flex items-center gap-2 text-sm">
            {#if r.status === 'ok'}
              <CheckCircle2 size="12" class="text-emerald-400 shrink-0" />
              <span class="text-emerald-300">{r.name}</span>
            {:else if r.status === 'skip'}
              <span class="text-slate-400">{r.name} — {r.message}</span>
            {:else}
              <X size="12" class="text-red-400 shrink-0" />
              <span class="text-red-300">{r.name}</span>
              <span class="text-slate-500">{r.message}</span>
            {/if}
          </div>
        {/each}
      </div>
    {/if}

    {#if loading}
      <div class="flex items-center justify-center gap-3 py-16 text-slate-500">
        <div class="h-6 w-6 animate-spin rounded-full border-2 border-slate-700 border-t-amber-500"></div>
        <span class="text-sm">Loading favorites…</span>
      </div>
    {:else}
      <!-- Repos -->
      {#if favoriteRepos.length > 0}
        <div>
          <h3 class="mb-2 text-sm font-bold text-white">Favorite Repos</h3>
          <div class="space-y-1">
            {#each favoriteRepos as repo (repo.name)}
              <div class="flex items-center gap-3 rounded-xl border border-white/5 bg-slate-800/50 px-4 py-3">
                <Star size="14" class="shrink-0 text-amber-400 fill-amber-400" />
                <div class="min-w-0 flex-1">
                  <span class="font-mono text-sm font-semibold text-white">{repo.name.split('/').pop()}</span>
                  <span class="ml-2 rounded border px-2 py-0.5 text-[10px] {categoryColor(repo.category)}">{repo.category}</span>
                  {#if repo.summary}<div class="text-[11px] text-slate-500">{repo.summary}</div>{/if}
                </div>
                <button class="rounded p-1.5 text-slate-400 hover:text-amber-400" title="Unstar"
                  onclick={() => unstar('repo', repo.name)}><Star size="14" class="fill-current" /></button>
                <button class="rounded-lg border border-white/10 bg-slate-800/60 px-2 py-1 text-[11px] text-slate-300 hover:text-white transition"
                  onclick={() => copyItemToUSB(repo, 'repo')}><Copy size="12" /> Copy</button>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Software -->
      {#if favoriteSoftware.length > 0}
        <div>
          <h3 class="mb-2 text-sm font-bold text-white">Favorite Software</h3>
          <div class="space-y-1">
            {#each favoriteSoftware as sw (sw.name)}
              <div class="flex items-center gap-3 rounded-xl border border-white/5 bg-slate-800/50 px-4 py-3">
                <Star size="14" class="shrink-0 text-amber-400 fill-amber-400" />
                <div class="min-w-0 flex-1">
                  <span class="text-sm font-semibold text-white">{sw.name}</span>
                  <span class="ml-2 rounded border px-2 py-0.5 text-[10px] {categoryColor(sw.category)}">{sw.category}</span>
                </div>
                <button class="rounded p-1.5 text-slate-400 hover:text-amber-400" title="Unstar"
                  onclick={() => unstar('software', sw.name)}><Star size="14" class="fill-current" /></button>
                <button class="rounded-lg border border-white/10 bg-slate-800/60 px-2 py-1 text-[11px] text-slate-300 hover:text-white transition"
                  onclick={() => copyItemToUSB(sw, 'software')}><Copy size="12" /> Copy</button>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Scripts -->
      {#if favoriteScripts.length > 0}
        <div>
          <h3 class="mb-2 text-sm font-bold text-white">Favorite Scripts</h3>
          <div class="space-y-1">
            {#each favoriteScripts as sc (sc.id)}
              <div class="flex items-center gap-3 rounded-xl border border-white/5 bg-slate-800/50 px-4 py-3">
                <Star size="14" class="shrink-0 text-amber-400 fill-amber-400" />
                <div class="min-w-0 flex-1">
                  <span class="font-mono text-sm font-semibold text-white">{sc.name}</span>
                  <span class="ml-2 rounded border px-2 py-0.5 text-[10px] {categoryColor(sc.category)}">{sc.category}</span>
                  <span class="ml-1 rounded border px-2 py-0.5 text-[10px] {langBadge(langOf(sc.name))}">{langOf(sc.name)}</span>
                </div>
                <button class="rounded p-1.5 text-slate-400 hover:text-amber-400" title="Unstar"
                  onclick={() => unstar('script', sc.name)}><Star size="14" class="fill-current" /></button>
                <button class="rounded-lg border border-white/10 bg-slate-800/60 px-2 py-1 text-[11px] text-slate-300 hover:text-white transition"
                  onclick={() => copyItemToUSB(sc, 'script')}><Copy size="12" /> Copy</button>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>
