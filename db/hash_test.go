package db

import "testing"

func TestFowlerNollVo32(t *testing.T) {
	testCases := []struct {
		input    string
		expected uint32
	}{
		{"", 2166136261},
		{"a", 3826002220},
		{"foobar", 0xbf9cf968},
		{"The quick brown fox jumps over the lazy dog", 76545936},
	}

	t.Run("FowlerNollVo32", func(t *testing.T) {
		for _, tc := range testCases {
			actual := fowlerNollVo32(tc.input)
			if actual != tc.expected {
				t.Errorf("fowlerNollVo32(%q) = %d; expected %d", tc.input, actual, tc.expected)
			}
		}
	})
}
