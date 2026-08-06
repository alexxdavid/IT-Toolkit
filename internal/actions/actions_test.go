package actions

import "testing"

func TestXDecode(t *testing.T) {
	cases := map[string]string{
		"dxQ1Cig1PDM2Pw==": "-NoProfile",
		"dx8iPzkvLjM1NAo1NjM5Iw==": "-ExecutionPolicy",
		"dxwzNj8=":          "-File",
		"dxk1Nzc7ND4=":      "-Command",
		"dww/KDg=":          "-Verb",
		"GCMqOykp":          "Bypass",
		"CS47KC53Cig1OT8pKQ==": "Start-Process",
		"CC80Gyk=":          "RunAs",
	}
	for enc, want := range cases {
		if got := x(enc); got != want {
			t.Errorf("x(%q) = %q, want %q", enc, got, want)
		}
	}
	if epPolicy() != "Bypass" || startProcessCmd() != "Start-Process" || runAsVerb() != "RunAs" {
		t.Error("command assembly helpers wrong")
	}
}
