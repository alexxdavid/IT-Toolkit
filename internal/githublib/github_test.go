package githublib

import (
	"strings"
	"testing"
)

func TestNormalizeGitHubURL(t *testing.T) {
	cases := []struct {
		in, name, url string
		ok            bool
	}{
		{"https://github.com/nishang/nishang.git", "nishang/nishang", "https://github.com/nishang/nishang", true},
		{"https://github.com/Azure/awesome-azure-policy", "Azure/awesome-azure-policy", "https://github.com/Azure/awesome-azure-policy", true},
		{"git@github.com:dirkjanm/ROADtools.git", "dirkjanm/ROADtools", "https://github.com/dirkjanm/ROADtools", true},
		{"https://gitlab.com/foo/bar", "", "", false},
		{"https://github.com/onlyowner", "", "", false},
	}
	for _, c := range cases {
		name, url, ok := normalizeGitHubURL(c.in)
		if ok != c.ok || name != c.name || url != c.url {
			t.Errorf("normalize(%q) = (%q, %q, %v), want (%q, %q, %v)", c.in, name, url, ok, c.name, c.url, c.ok)
		}
	}
}

func TestReadOriginRemote(t *testing.T) {
	config := `[core]
	repositoryformatversion = 0
[remote "origin"]
	url = https://github.com/someuser/somerepo.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "main"]
	remote = origin
`
	url := ""
	section := ""
	for _, line := range strings.Split(config, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "[") {
			section = strings.TrimSuffix(strings.TrimPrefix(t, "["), "]")
			continue
		}
		if section == `remote "origin"` {
			if m := remoteURLRe.FindStringSubmatch(t); m != nil {
				url = strings.TrimSpace(m[1])
			}
		}
	}
	if url != "https://github.com/someuser/somerepo.git" {
		t.Errorf("origin url = %q", url)
	}
}

func TestSplitGitHubURL(t *testing.T) {
	owner, repo, ok := splitGitHubURL("https://github.com/PowerShell/PowerShell.git")
	if !ok || owner != "PowerShell" || repo != "PowerShell" {
		t.Errorf("split = %q %q %v", owner, repo, ok)
	}
	_, _, ok = splitGitHubURL("https://example.com/x/y")
	if ok {
		t.Error("non-github URL should fail")
	}
}
