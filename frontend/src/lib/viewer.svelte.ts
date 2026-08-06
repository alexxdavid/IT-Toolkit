import type { RepoDetail, ScriptFile } from '../lib/types'

export type ViewerState = {
  open: boolean
  script: ScriptFile | null
  repoDetail: RepoDetail | null
  tab: 'source' | 'readme'
}

export const viewerState = $state<ViewerState>({
  open: false,
  script: null,
  repoDetail: null,
  tab: 'source'
})

export function openScript(sc: ScriptFile) {
  viewerState.open = true
  viewerState.script = sc
  viewerState.repoDetail = null
  viewerState.tab = 'source'
}

export function openRepoReadme(d: RepoDetail) {
  viewerState.open = true
  viewerState.script = null
  viewerState.repoDetail = d
  viewerState.tab = 'readme'
}

export function closeViewer() {
  viewerState.open = false
  viewerState.script = null
  viewerState.repoDetail = null
}
