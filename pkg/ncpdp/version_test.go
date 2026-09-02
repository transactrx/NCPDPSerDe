package ncpdp

import "testing"

func TestVersionAtLeast(t *testing.T) {
	cases := []struct {
		version, minimum string
		want             bool
	}{
		{D0, D0, true},
		{F6, D0, true},
		{F6, F6, true},
		{D0, F6, false},
		// Unknown transmission versions rank oldest.
		{"ZZ", F6, false},
		{"", F6, false},
		{"", D0, false},
		// An unknown minimum can never be satisfied.
		{F6, "ZZ", false},
	}

	for _, c := range cases {
		if got := VersionAtLeast(c.version, c.minimum); got != c.want {
			t.Errorf("VersionAtLeast(%q, %q) = %v, want %v", c.version, c.minimum, got, c.want)
		}
	}
}

func TestOmitsGroupSeparator(t *testing.T) {
	if OmitsGroupSeparator(D0) {
		t.Error("D0 transmissions use group separators")
	}
	if !OmitsGroupSeparator(F6) {
		t.Error("F6 transmissions have no group separators")
	}
	if OmitsGroupSeparator("") {
		t.Error("unknown versions must default to D0-era separator behavior")
	}
}

func TestHeaderLeadsWithVersion(t *testing.T) {
	if HeaderLeadsWithVersion(D0) {
		t.Error("D0 request headers lead with the BIN, not the version")
	}
	if !HeaderLeadsWithVersion(F6) {
		t.Error("F6 headers lead with the version code")
	}
	if HeaderLeadsWithVersion("ZZ") {
		t.Error("unknown versions must default to D0-era header shape")
	}
}

func TestDetectTransmissionVersion(t *testing.T) {
	cases := []struct {
		name, raw, want string
	}{
		{"F6 request leads with version", "F6B100880151TEST      101...", F6},
		{"D0 response leads with version", "D0B11A011234567893     20210118", D0},
		{"D0 request carries version after BIN", "880151D0B1          101...", D0},
		// Version detection is exact-match against known versions, not an
		// F-prefix heuristic: an unmodeled F-family code is NOT treated as F6.
		{"unknown F-prefixed version is not F6", "FXB100880151TEST      101...", Empty},
		{"unrecognized raw", "GARBAGE-DATA", Empty},
		{"short raw", "D", Empty},
	}

	for _, c := range cases {
		if got := DetectTransmissionVersion(c.raw); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
