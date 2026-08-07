export const categoryColors: Record<string, string> = {
  'Microsoft 365': 'border-sky-400/50 bg-sky-500/25 text-sky-200',
  'Intune & Endpoint': 'border-cyan-400/50 bg-cyan-500/25 text-cyan-200',
  'Microsoft Defender': 'border-purple-400/50 bg-purple-500/25 text-purple-200',
  'Windows Administration': 'border-blue-400/50 bg-blue-500/25 text-blue-200',
  'Windows Security': 'border-red-400/50 bg-red-500/25 text-red-200',
  'Active Directory': 'border-amber-400/50 bg-amber-500/25 text-amber-200',
  'Azure & Cloud': 'border-indigo-400/50 bg-indigo-500/25 text-indigo-200',
  'Network & Infrastructure': 'border-teal-400/50 bg-teal-500/25 text-teal-200',
  'Security (Defensive)': 'border-emerald-400/50 bg-emerald-500/25 text-emerald-200',
  'Security': 'border-emerald-400/50 bg-emerald-500/25 text-emerald-200',
  'Security Tools': 'border-rose-400/50 bg-rose-500/25 text-rose-200',
  'Security Tools (Offensive)': 'border-rose-400/50 bg-rose-500/25 text-rose-200',
  'Security & Monitoring': 'border-emerald-400/50 bg-emerald-500/25 text-emerald-200',
  'Security & Hardening': 'border-emerald-400/50 bg-emerald-500/25 text-emerald-200',
  'OSINT': 'border-fuchsia-400/50 bg-fuchsia-500/25 text-fuchsia-200',
  'GRC & Compliance': 'border-yellow-400/50 bg-yellow-500/25 text-yellow-200',
  'Monitoring & SIEM': 'border-violet-400/50 bg-violet-500/25 text-violet-200',
  'DevOps & Labs': 'border-orange-400/50 bg-orange-500/25 text-orange-200',
  'Databases': 'border-lime-400/50 bg-lime-500/25 text-lime-200',
  'AI & Machine Learning': 'border-pink-400/50 bg-pink-500/25 text-pink-200',
  'AI & Reference': 'border-lime-400/50 bg-lime-500/25 text-lime-200',
  'Utilities & Tools': 'border-zinc-300/50 bg-zinc-400/25 text-zinc-200',
  'Utilities': 'border-zinc-300/50 bg-zinc-400/25 text-zinc-200',
  'PowerShell General': 'border-blue-400/50 bg-blue-500/25 text-blue-200',
  'Uncategorized': 'border-slate-400/50 bg-slate-400/25 text-slate-300',
  'Other': 'border-slate-400/50 bg-slate-400/25 text-slate-300',
  'Virtualization & Labs': 'border-violet-400/50 bg-violet-500/25 text-violet-200',
  'Remote Access': 'border-indigo-400/50 bg-indigo-500/25 text-indigo-200',
  'Networking': 'border-teal-400/50 bg-teal-500/25 text-teal-200',
  'Backup': 'border-amber-400/50 bg-amber-500/25 text-amber-200',
  'Productivity': 'border-sky-400/50 bg-sky-500/25 text-sky-200',
  'Development': 'border-cyan-400/50 bg-cyan-500/25 text-cyan-200',
  'RMM': 'border-orange-400/50 bg-orange-500/25 text-orange-200',
  'Automated & DevOps': 'border-orange-400/50 bg-orange-500/25 text-orange-200',
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
