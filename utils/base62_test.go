package utils

import "testing"

func TestBase62Encode(t *testing.T) {

	tests := []struct {
		num int64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{62, "10"},
		{63, "11"},
		{12345, "3d7"},
	}

	for _, test := range tests {
		got := Base62Encode(test.num)

		if got != test.expected {
			t.Errorf("base62Encode(%d) = %s, want %s", test.num, got, test.expected)
		}
	}
}
