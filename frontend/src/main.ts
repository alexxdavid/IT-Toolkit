import './app.css'
import { mount } from 'svelte'
import App from './App.svelte'
import * as AppBind from '../wailsjs/go/main/App'

window.addEventListener('error', (e) => {
  try {
    AppBind.Log('window.error: ' + (e.message || String(e.error)) + ' @ ' + (e.filename || '') + ':' + (e.lineno || ''))
  } catch {
    /* ignore */
  }
})

window.addEventListener('unhandledrejection', (e) => {
  try {
    const r = e.reason
    AppBind.Log('unhandledrejection: ' + (r && r.message ? r.message : String(r)))
  } catch {
    /* ignore */
  }
})

const app = mount(App, {
  target: document.getElementById('app')!
})

export default app
