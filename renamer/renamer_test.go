package main

import "testing"

func TestParser(t *testing.T) {
	testCases := []struct {
		input  string
		wanted string
	}{
		{"birthday_001.txt", "Birthday - 1 of 4.txt"},
		{"birthday_002.txt", "Birthday - 2 of 4.txt"},
		{"birthday_003.txt", "Birthday - 3 of 4.txt"},
		{"birthday_004.txt", "Birthday - 4 of 4.txt"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			gotten, _ := match(tc.input, 4)
			if gotten != tc.wanted {
				t.Errorf("wanted %s, but got %s", tc.wanted, gotten)
			}
		})
	}
}

func TestParser2(t *testing.T) {
	testCases := []struct {
		input  string
		wanted string
	}{
		{"n_008.txt", "N - 8 of 10.txt"},
		{"n_009.txt", "N - 9 of 10.txt"},
		{"n_010.txt", "N - 10 of 10.txt"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			gotten, _ := match(tc.input, 10)
			if gotten != tc.wanted {
				t.Errorf("wanted %s, but got %s", tc.wanted, gotten)
			}
		})
	}
}
