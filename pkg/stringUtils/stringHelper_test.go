package stringutils

import (
	"slices"
	"testing"
)

func Test_IndexOfAny(t *testing.T) {
	const seg = byte(0x1E)
	const grp = byte(0x1D)

	tests := []struct {
		name      string
		str       string
		fromIndex int
		searchFor []byte
		want      int
	}{
		{"empty string", "", 0, []byte{seg}, -1},
		{"empty search set", "abc", 0, nil, -1},
		{"fromIndex at end", "abc", 3, []byte{'b'}, -1},
		{"fromIndex past end", "abc", 10, []byte{'b'}, -1},
		{"single byte match at start", "\x1Eabc", 0, []byte{seg}, 0},
		{"single byte match mid string", "ab\x1Ecd", 0, []byte{seg}, 2},
		{"single byte match at last byte", "abc\x1E", 0, []byte{seg}, 3},
		{"single byte match at fromIndex", "ab\x1Ec", 2, []byte{seg}, 2},
		{"single byte skips match before fromIndex", "\x1Eab\x1Ecd", 1, []byte{seg}, 3},
		{"single byte no match", "abcd", 0, []byte{seg}, -1},
		{"multi byte earliest match wins", "ab\x1Dc\x1Ed", 0, []byte{seg, grp}, 2},
		{"multi byte respects fromIndex", "ab\x1Dc\x1Ed", 3, []byte{seg, grp}, 4},
		{"multi byte no match", "abcd", 0, []byte{seg, grp}, -1},
		{"multi-byte rune content before separator", "ab€\x1Ec", 0, []byte{seg}, 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IndexOfAny(test.str, test.fromIndex, test.searchFor)
			if got != test.want {
				t.Errorf("IndexOfAny(%q, %v, %v) mismatch. Wanted: %v   Got: %v", test.str, test.fromIndex, test.searchFor, test.want, got)
			}
		})
	}
}

func Test_Substring(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		fromIndex int
		length    int
		want      string
	}{
		{"empty string", "", 0, 3, ""},
		{"fromIndex past end", "abc", 5, 1, ""},
		{"fromIndex at end", "abc", 3, -1, ""},
		{"negative length returns remainder", "abcdef", 2, -1, "cdef"},
		{"zero length returns empty", "abcdef", 2, 0, ""},
		{"length in range", "abcdef", 1, 3, "bcd"},
		{"length past end returns remainder", "abcdef", 4, 10, "ef"},
		{"whole string", "abcdef", 0, 6, "abcdef"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Substring(test.str, test.fromIndex, test.length)
			if got != test.want {
				t.Errorf("Substring(%q, %v, %v) mismatch. Wanted: %q   Got: %q", test.str, test.fromIndex, test.length, test.want, got)
			}
		})
	}
}

func Test_RightPadExact(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		padChar byte
		length  int
		want    string
	}{
		{"pads short string", "B1", ' ', 5, "B1   "},
		{"exact length unchanged", "B1", ' ', 2, "B1"},
		{"truncates long string", "BILLING", ' ', 4, "BILL"},
		{"pads empty string", "", '0', 3, "000"},
		{"zero length returns empty", "AB", ' ', 0, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := RightPadExact(test.str, test.padChar, test.length)
			if got != test.want {
				t.Errorf("RightPadExact(%q, %q, %v) mismatch. Wanted: %q   Got: %q", test.str, test.padChar, test.length, test.want, got)
			}
		})
	}
}

func Test_TrimAll(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		trimChars []byte
		want      string
	}{
		{"trims both ends", "  AB12  ", []byte{' '}, "AB12"},
		{"trims multiple characters", "0 0AB120 ", []byte{' ', '0'}, "AB12"},
		{"interior characters preserved", " A B ", []byte{' '}, "A B"},
		{"nothing to trim", "AB12", []byte{' '}, "AB12"},
		{"all trimmed", "   ", []byte{' '}, ""},
		{"empty string", "", []byte{' '}, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := TrimAll(test.str, test.trimChars)
			if got != test.want {
				t.Errorf("TrimAll(%q, %v) mismatch. Wanted: %q   Got: %q", test.str, test.trimChars, test.want, got)
			}
		})
	}
}

func Test_SplitBySeparator(t *testing.T) {
	const sep = byte(0x1C)

	tests := []struct {
		name string
		data string
		want []string
	}{
		{"no separator returns empty slice", "AM04", []string{}},
		{"empty string returns empty slice", "", []string{}},
		{"content before first separator is dropped", "AM04\x1Cf1\x1Cf2", []string{"f1", "f2"}},
		{"leading separator", "\x1Cf1\x1Cf2", []string{"f1", "f2"}},
		{"trailing separator yields empty element", "\x1Cf1\x1C", []string{"f1", ""}},
		{"consecutive separators yield empty element", "\x1Cf1\x1C\x1Cf2", []string{"f1", "", "f2"}},
		{"only separator", "\x1C", []string{""}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := SplitBySeparator(test.data, sep)
			if !slices.Equal(got, test.want) {
				t.Errorf("SplitBySeparator(%q) mismatch. Wanted: %q   Got: %q", test.data, test.want, got)
			}
		})
	}
}

func Test_Chunk(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		chunkSize int
		want      []string
	}{
		{"empty string returns nil", "", 2, nil},
		{"chunk size equal to length", "abcd", 4, []string{"abcd"}},
		{"chunk size larger than length", "abcd", 10, []string{"abcd"}},
		{"even chunks", "abcdef", 2, []string{"ab", "cd", "ef"}},
		{"uneven final chunk", "abcde", 2, []string{"ab", "cd", "e"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Chunk(test.str, test.chunkSize)
			if !slices.Equal(got, test.want) {
				t.Errorf("Chunk(%q, %v) mismatch. Wanted: %q   Got: %q", test.str, test.chunkSize, test.want, got)
			}
		})
	}
}

func Test_IndexOfAll(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		searchTerm string
		want       []int
	}{
		{"multiple matches", "a\x1Eb\x1Ec", "\x1E", []int{1, 3}},
		{"single match", "abc", "b", []int{1}},
		{"no match returns empty slice", "abc", "z", []int{}},
		{"overlapping matches", "aaa", "aa", []int{0, 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IndexOfAll(test.str, test.searchTerm)
			if !slices.Equal(got, test.want) {
				t.Errorf("IndexOfAll(%q, %q) mismatch. Wanted: %v   Got: %v", test.str, test.searchTerm, test.want, got)
			}
		})
	}
}

func Test_IndexAt(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		searchTerm string
		startIndex int
		want       int
	}{
		{"match after start index", "abcabc", "abc", 1, 3},
		{"match at start index", "abcabc", "abc", 3, 3},
		{"no match after start index", "abcdef", "abc", 1, -1},
		{"match at zero", "abc", "abc", 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IndexAt(test.str, test.searchTerm, test.startIndex)
			if got != test.want {
				t.Errorf("IndexAt(%q, %q, %v) mismatch. Wanted: %v   Got: %v", test.str, test.searchTerm, test.startIndex, test.want, got)
			}
		})
	}
}

func Test_SplitFirst(t *testing.T) {
	const sep = byte(0x1C)

	tests := []struct {
		name      string
		str       string
		fromIndex int
		want      string
	}{
		{"first field", "ab\x1Ccd\x1Cef", 0, "ab"},
		{"middle field", "ab\x1Ccd\x1Cef", 3, "cd"},
		{"last field returns remainder", "ab\x1Ccd\x1Cef", 6, "ef"},
		{"no separator returns remainder", "abcdef", 1, "bcdef"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := SplitFirst(test.str, test.fromIndex, []byte{sep})
			if got != test.want {
				t.Errorf("SplitFirst(%q, %v) mismatch. Wanted: %q   Got: %q", test.str, test.fromIndex, test.want, got)
			}
		})
	}
}
