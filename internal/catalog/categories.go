package catalog

import "strings"

type categoryRule struct {
	name     string
	keywords []string
}

// Rules are evaluated in order; the first rule whose keyword is a
// case-insensitive substring of the repo folder name wins.
var categoryRules = []categoryRule{
	{"Microsoft 365", []string{"office365", "office-365", "office-tool", "microsoft365", "microsoft-copilot", "onedrive", "m365", "cli-microsoft365"}},
	{"Intune & Endpoint", []string{"mem.zone", "mem-zone", "intune", "endpoint", "debloat", "windows toolbox", "win-debloat", "hyper-rev"}},
	{"OSINT", []string{"osint", "sherlock", "maigret", "osintgram", "spiderfoot", "web-check", "shannon"}},
	{"Security & Hardening", []string{"hardening", "baseline", "secure-host", "windows_hardening", "dsinternals", "fail2ban", "lynis", "adalanche"}},
	{"Security Tools", []string{"nishang", "payload", "hacking", "hacker", "imhex", "x64dbg", "rustscan", "gitleaks", "trufflehog", "trivy", "checkov", "dirsearch", "sn1per", "scanners", "hydra", "infosec", "strix", "pentest", "epv", "maigret"}},
	{"GRC & Compliance", []string{"grc", "ciso", "comply", "prowler", "pacbot", "verifywise", "probo", "claude-grc"}},
	{"Monitoring & SIEM", []string{"wazuh", "securityonion", "opensecurity", "fixinventory", "misp", "alert"}},
	{"Virtualization & Labs", []string{"hyperv", "hyper-v", "automatedlab", "detectionlab", "badlab", "vagrant", "machyperv", "osx-hyper", "vma", "windows-vagrant"}},
	{"Databases", []string{"sql-server", "sqlserver", "dbgate", "davmail", "fluentmigrator", "azuredatastudio", "beekeeper", "querybuilder"}},
	{"Automation & DevOps", []string{"ansible", "kpt", "kyverno", "opa", "cloudformation", "packer", "choco", "infisical", "azure-policy"}},
	{"AI & Reference", []string{"llms", "ai-for", "machine-learning", "minimind", "spacy", "flowise", "auto", "copilot", "cheatsheet", "awesome", "500-ai"}},
	{"PowerShell General", []string{"powershell"}},
}

// Uncategorized is the fallback category.
const Uncategorized = "Uncategorized"

// Categorize returns the category for a repo folder name.
func Categorize(folderName string) string {
	lower := strings.ToLower(folderName)
	for _, rule := range categoryRules {
		for _, kw := range rule.keywords {
			if strings.Contains(lower, kw) {
				return rule.name
			}
		}
	}
	return Uncategorized
}

// AllCategories returns category names in rule order (Uncategorized last).
func AllCategories() []string {
	names := make([]string, 0, len(categoryRules)+1)
	for _, rule := range categoryRules {
		names = append(names, rule.name)
	}
	names = append(names, Uncategorized)
	return names
}
