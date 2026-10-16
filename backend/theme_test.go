package backend

import "testing"

func TestDefaultTheme(t *testing.T) {
	th := DefaultTheme()
	cases := []struct {
		name string
		got  HexCode
		want HexCode
	}{
		{"text", th.Text, "#C6D0F5"},
		{"border", th.Border, "#C6D0F5"},
		{"weekday", th.Weekday, "#C6D0F5"},
		{"weekend", th.Weekend, "#939AB7"},
		{"dim", th.Dim, "#939AB7"},
		{"accent", th.Accent, "#CA9EE6"},
		{"late", th.Late, "#E78284"},
		{"soon", th.Soon, "#E5C890"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("DefaultTheme().%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestToneColor(t *testing.T) {
	th := Theme{Late: "#111111", Soon: "#222222", Dim: "#333333"}
	cases := []struct {
		tone Tone
		want HexCode
	}{
		{ToneLate, "#111111"},
		{ToneSoon, "#222222"},
		{ToneNormal, "#333333"},
	}
	for _, c := range cases {
		if got := th.ToneColor(c.tone); got != c.want {
			t.Errorf("ToneColor(%d) = %q, want %q", c.tone, got, c.want)
		}
	}
}
