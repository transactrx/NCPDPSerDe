package stringutils

import "testing"

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
