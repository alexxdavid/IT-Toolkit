export const categoryColors: Record<string, string> = {
  'Microsoft 365': 'border-sky-500/30 bg-sky-500/15 text-sky-300',
  'Intune & Endpoint': 'border-cyan-500/30 bg-cyan-500/15 text-cyan-300',
  'Microsoft Defender': 'border-purple-500/30 bg-purple-500/15 text-purple-300',
  'Windows Administration': 'border-blue-500/30 bg-blue-500/15 text-blue-300',
  'Windows Security': 'border-red-500/30 bg-red-500/15 text-red-300',
  'Active Directory': 'border-amber-500/30 bg-amber-500/15 text-amber-300',
  'Azure & Cloud': 'border-indigo-500/30 bg-indigo-500/15 text-indigo-300',
  'Network & Infrastructure': 'border-teal-500/30 bg-teal-500/15 text-teal-300',
  'Security (Defensive)': 'border-emerald-500/30 bg-emerald-500/15 text-emerald-300',
  'Security Tools': 'border-rose-500/30 bg-rose-500/15 text-rose-300',
  'OSINT': 'border-fuchsia-500/30 bg-fuchsia-500/15 text-fuchsia-300',
  'GRC & Compliance': 'border-yellow-500/30 bg-yellow-500/15 text-yellow-300',
  'Monitoring & SIEM': 'border-violet-500/30 bg-violet-500/15 text-violet-300',
  'DevOps & Labs': 'border-orange-500/30 bg-orange-500/15 text-orange-300',
  'Databases': 'border-lime-500/30 bg-lime-500/15 text-lime-300',
  'AI & Machine Learning': 'border-pink-500/30 bg-pink-500/15 text-pink-300',
  'Utilities & Tools': 'border-zinc-400/30 bg-zinc-400/15 text-zinc-300',
  'Uncategorized': 'border-slate-500/30 bg-slate-500/15 text-slate-400',
  'Other': 'border-slate-500/30 bg-slate-500/15 text-slate-400',
  'PowerShell General': 'border-blue-500/30 bg-blue-500/15 text-blue-300',
  'Security & Hardening': 'border-emerald-500/30 bg-emerald-500/15 text-emerald-300',
  'Automation & DevOps': 'border-orange-500/30 bg-orange-500/15 text-orange-300',
  'AI & Reference': 'border-lime-500/30 bg-lime-500/15 text-lime-300',
  'Security Tools (Offensive)': 'border-rose-500/30 bg-rose-500/15 text-rose-300',
  'Virtualization & Labs': 'border-violet-500/30 bg-violet-500/15 text-violet-300',
}

export const categoryColor = (c: string): string => categoryColors[c] ?? categoryColors['Uncategorized']

export const langBadges: Record<string, string> = {
  PowerShell: 'text-blue-300 bg-blue-500/10 border-blue-500/25',
  Python: 'text-yellow-300 bg-yellow-500/10 border-yellow-500/25',
  Batch: 'text-zinc-300 bg-zinc-500/10 border-zinc-500/25',
  VBScript: 'text-orange-300 bg-orange-500/10 border-orange-500/25',
  Shell: 'text-emerald-300 bg-emerald-500/10 border-emerald-500/25',
  SQL: 'text-violet-300 bg-violet-500/10 border-violet-500/25',
  Registry: 'text-cyan-300 bg-cyan-500/10 border-cyan-500/25',
  Markdown: 'text-slate-300 bg-slate-500/10 border-slate-500/25',
  JSON: 'text-amber-300 bg-amber-500/10 border-amber-500/25',
  YAML: 'text-rose-300 bg-rose-500/10 border-rose-500/25',
  XML: 'text-teal-300 bg-teal-500/10 border-teal-500/25',
  INI: 'text-lime-300 bg-lime-500/10 border-lime-500/25',
  Text: 'text-slate-400 bg-slate-500/10 border-slate-500/25',
  HTML: 'text-red-300 bg-red-500/10 border-red-500/25',
  CSV: 'text-indigo-300 bg-indigo-500/10 border-indigo-500/25'
}

export const langBadge = (l: string): string => langBadges[l] ?? langBadges['Text']

export const typeChips = ['all', 'PowerShell', 'Python', 'Batch', 'VBScript', 'Shell', 'SQL', 'Registry']
