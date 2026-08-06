<script lang="ts">
  import RepoCard from './RepoCard.svelte'
  import { SearchX, FolderTree } from 'lucide-svelte'
  import { ui, catalog, visibleRepos } from '../lib/store.svelte.ts'

  const repos = $derived(visibleRepos())
</script>

<div class="min-h-0 flex-1 overflow-y-auto p-4">
  {#if !ui.loaded}
    <div class="flex h-full flex-col items-center justify-center gap-3 text-slate-500">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-slate-700 border-t-cyan-500"></div>
      <p class="text-sm">Loading catalog…</p>
    </div>
  {:else if repos.length === 0}
    <div class="flex h-full flex-col items-center justify-center gap-3 text-slate-500">
      {#if ui.query.trim() || catalog.folders.length === 0}
        <SearchX size="36" />
        <p class="text-sm">No repositories match.</p>
      {:else}
        <FolderTree size="36" />
        <p class="text-sm">Nothing indexed yet — add a scripts folder from the sidebar.</p>
      {/if}
    </div>
  {:else}
    <div class="mx-auto max-w-4xl space-y-2">
      {#each repos as r (r.id)}
        <RepoCard repo={r} />
      {/each}
    </div>
  {/if}
</div>
