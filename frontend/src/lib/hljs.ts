import hljs from 'highlight.js/lib/core'
import powershell from 'highlight.js/lib/languages/powershell'
import python from 'highlight.js/lib/languages/python'
import bash from 'highlight.js/lib/languages/bash'
import dos from 'highlight.js/lib/languages/dos'
import sql from 'highlight.js/lib/languages/sql'
import json from 'highlight.js/lib/languages/json'
import yaml from 'highlight.js/lib/languages/yaml'
import xml from 'highlight.js/lib/languages/xml'
import ini from 'highlight.js/lib/languages/ini'
import markdown from 'highlight.js/lib/languages/markdown'
import vbscript from 'highlight.js/lib/languages/vbscript'
import javascript from 'highlight.js/lib/languages/javascript'
import 'highlight.js/styles/github-dark.min.css'
import { escapeHtml } from './format'

hljs.registerLanguage('powershell', powershell)
hljs.registerLanguage('python', python)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('dos', dos)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('json', json)
hljs.registerLanguage('yaml', yaml)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('ini', ini)
hljs.registerLanguage('markdown', markdown)
hljs.registerLanguage('vbscript', vbscript)
hljs.registerLanguage('javascript', javascript)

const langMap: Record<string, string> = {
  PowerShell: 'powershell',
  Python: 'python',
  Shell: 'bash',
  Batch: 'dos',
  VBScript: 'vbscript',
  SQL: 'sql',
  JSON: 'json',
  YAML: 'yaml',
  XML: 'xml',
  INI: 'ini',
  Markdown: 'markdown',
  HTML: 'xml',
  Registry: 'dos',
  TOML: 'ini',
  Text: 'plaintext',
  CSV: 'plaintext'
}

export function highlight(lang: string, code: string): string {
  const l = langMap[lang] ?? 'plaintext'
  try {
    return hljs.highlight(code, { language: l }).value
  } catch {
    return escapeHtml(code)
  }
}
