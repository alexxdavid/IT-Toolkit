<script lang="ts">
  import { onMount } from 'svelte'
  import { RefreshCw, Download, X, ExternalLink } from 'lucide-svelte'
  import { api } from '../lib/api'

  let updateInfo = $state<any>(null)
  let showBanner = $state(false)
  let downloading = $state(false)
  let downloadPercent = $state(0)
  let pollTimer: ReturnType<typeof setInterval> | null = null
  let checking = $state(false)

  const appVersion = '1.0.0'

  async function checkForUpdate(force: boolean = false) {
    checking = true
    try {
      const info = await api.checkForUpdate(force)
      if (info && info.available) {
        updateInfo = info
        showBanner = true
      }
    } catch { /* silent */ }
    checking = false
  }

  async function startDownload() {
    if (!updateInfo) return
    downloading = true
    downloadPercent = 0
    pollTimer = setInterval(async () => {
      try {
        const p = await api.getUpdateProgress()
        if (p?.percent > 0) downloadPercent = Math.min(100, Math.round(p.percent))
      } catch { /* */ }
    }, 300)
    try {
      const path = await api.downloadUpdate(updateInfo.installer_url, updateInfo.version)
      if (pollTimer) clearInterval(pollTimer)
      if (path) {
        downloadPercent = 100
        await api.applyUpdate(path)
      }
    } catch (e) {
      if (pollTimer) clearInterval(pollTimer)
      console.error('Update failed:', e)
    }
    downloading = false
  }

  onMount(() => {
    checkForUpdate(false)
  })
</script>

{#if showBanner}
  <div class="fixed bottom-0 left-0 right-0 z-50 border-t border-emerald-500/30 bg-slate-900/95 px-4 py-3 shadow-2xl backdrop-blur-md">
    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/20">
          <Download size="16" class="text-emerald-400" />
        </div>
        <div>
          <div class="text-xs font-bold text-white">Update available: v{updateInfo?.version}</div>
          {#if updateInfo?.notes}
            <div class="text-[10px] text-slate-400 max-w-md truncate">{updateInfo.notes}</div>
          {/if}
        </div>
      </div>
      <div class="flex-1"></div>
      {#if downloading}
        <div class="flex items-center gap-3">
          <div class="h-2 w-40 overflow-hidden rounded bg-slate-800">
            <div class="h-full rounded bg-emerald-500 transition-all" style="width: {downloadPercent}%"></div>
          </div>
          <span class="text-[10px] font-mono text-emerald-400">{downloadPercent}%</span>
        </div>
      {:else}
        <button class="flex items-center gap-2 rounded-lg bg-emerald-500 px-3 py-1.5 text-xs font-bold text-slate-950 transition hover:bg-emerald-400"
          onclick={startDownload}>
          <Download size="14" /> Install & Restart
        </button>
      {/if}
      <button class="rounded p-1 text-slate-400 hover:text-white" onclick={() => showBanner = false}>
        <X size="14" />
      </button>
    </div>
  </div>
{/if}
