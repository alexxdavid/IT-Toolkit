<script lang="ts">
  import { onMount } from 'svelte'
  import { AlertTriangle, X } from 'lucide-svelte'
  import Sidebar from './components/Sidebar.svelte'
  import GithubLibrary from './components/GithubLibrary.svelte'
  import Settings from './components/Settings.svelte'
  import { api } from './lib/api'
  import { ui, setCatalog } from './lib/store.svelte'
  import logo from './lib/assets/logo.png'

  let scanErr = $state<string | null>(null)

  onMount(() => {
    refresh()
  })

  async function refresh() {
    try {
      const v = await api.getCatalog()
      ui.loadError = ''
      setCatalog(v)
    } catch (e) {
      ui.loadError = String(e)
    }
  }
</script>

<div class="flex h-screen flex-col bg-[#0b1120]">
  <!-- Logo top-right -->
  <div class="pointer-events-none absolute top-2 right-3 z-40">
    <img src={logo} alt="Solutions IT" class="h-10 w-auto object-contain" />
  </div>

  {#if ui.settingsOpen}
    <div class="flex min-h-0 flex-1 overflow-hidden">
      <Sidebar />
      <Settings />
    </div>
  {:else}
    <div class="flex min-h-0 flex-1 overflow-hidden">
      <Sidebar />
      <GithubLibrary />
    </div>
  {/if}

  {#if scanErr}
    <div class="fixed bottom-4 right-4 z-50 flex max-w-md items-start gap-3 rounded-lg border border-rose-500/40 bg-rose-950/90 px-4 py-3 text-sm text-rose-200 shadow-xl">
      <AlertTriangle size="16" class="mt-0.5 shrink-0" />
      <div class="min-w-0 flex-1">{scanErr}</div>
      <button class="shrink-0 text-rose-400 hover:text-rose-200" onclick={() => scanErr = null}><X size="14" /></button>
    </div>
  {/if}
</div>
