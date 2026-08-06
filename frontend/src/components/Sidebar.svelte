<script lang="ts">
  import { Package, HardDrive, FolderSearch, Settings, FileCode, Star } from 'lucide-svelte'
  import { ui, setSettingsOpen, setLocalTab } from '../lib/store.svelte.ts'
</script>

<aside class="flex h-full w-48 shrink-0 flex-col border-r border-slate-800 bg-slate-900/60">
  <div class="flex flex-col gap-1 p-3">
    <button class="flex items-center gap-2 rounded-lg px-3 py-3 text-sm font-medium transition
      {ui.localTab === 'repos' ? 'bg-gradient-to-r from-emerald-500/15 to-teal-500/15 text-emerald-300 border border-emerald-500/30' : 'text-slate-300 hover:bg-slate-800 border border-transparent'}"
      onclick={() => { setLocalTab('repos'); setSettingsOpen(false) }}>
      <Package size="16" /> Repos
    </button>
    <button class="flex items-center gap-2 rounded-lg px-3 py-3 text-sm font-medium transition
      {ui.localTab === 'software' ? 'bg-gradient-to-r from-sky-500/15 to-blue-500/15 text-sky-300 border border-sky-500/30' : 'text-slate-300 hover:bg-slate-800 border border-transparent'}"
      onclick={() => { setLocalTab('software'); setSettingsOpen(false) }}>
      <HardDrive size="16" /> Software
    </button>
    <button class="flex items-center gap-2 rounded-lg px-3 py-3 text-sm font-medium transition
      {ui.localTab === 'scripts' ? 'bg-gradient-to-r from-violet-500/15 to-purple-500/15 text-violet-300 border border-violet-500/30' : 'text-slate-300 hover:bg-slate-800 border border-transparent'}"
      onclick={() => { setLocalTab('scripts'); setSettingsOpen(false) }}>
      <FileCode size="16" /> Scripts
    </button>
    <button class="flex items-center gap-2 rounded-lg px-3 py-3 text-sm font-medium transition
      {ui.localTab === 'local' ? 'bg-gradient-to-r from-emerald-500/15 to-teal-500/15 text-emerald-300 border border-emerald-500/30' : 'text-slate-300 hover:bg-slate-800 border border-transparent'}"
      onclick={() => { setLocalTab('local'); setSettingsOpen(false) }}>
      <FolderSearch size="16" /> Local
    </button>
    <button class="flex items-center gap-2 rounded-lg px-3 py-3 text-sm font-medium transition
      {ui.localTab === 'favorites' ? 'bg-gradient-to-r from-amber-500/15 to-orange-500/15 text-amber-300 border border-amber-500/30' : 'text-slate-300 hover:bg-slate-800 border border-transparent'}"
      onclick={() => { setLocalTab('favorites'); setSettingsOpen(false) }}>
      <Star size="16" /> Favorites
    </button>
  </div>

  <div class="flex-1"></div>

  <div class="p-3">
    <button
      class="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm transition
        {ui.settingsOpen
          ? 'bg-gradient-to-r from-amber-500/15 to-orange-500/15 text-amber-300 border border-amber-500/30'
          : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200 border border-transparent'}"
      onclick={() => setSettingsOpen(!ui.settingsOpen)}>
      <Settings size="14" /> Settings
    </button>
  </div>
</aside>
