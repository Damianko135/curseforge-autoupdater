package version

import "testing"

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		v1, v2 string
		want   int
	}{
		{"1.2.3", "1.2.3", 0},
		{"1.2.4", "1.2.3", 1},
		{"1.2.2", "1.2.3", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.0.0-alpha", "1.0.0", -1},
	}
	for _, c := range cases {
		got, err := CompareVersions(c.v1, c.v2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != c.want {
			t.Errorf("CompareVersions(%q, %q) = %d, want %d", c.v1, c.v2, got, c.want)
		}
	}
}
