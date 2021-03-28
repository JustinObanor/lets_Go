package main

import (
	"testing"
)

type testCase struct {
	min, max    int
	testsAmount int
}

// checks if areacode is of correct length
func TestAreaCode(t *testing.T) {
	tc := testCase{
		min: 100, max: 1000,
		testsAmount: 10000,
	}

	for i := 0; i < tc.testsAmount; i++ {
		t.Run("", func(t *testing.T) {
			result := randomInt(tc.min, tc.max)

			if len(result) != 3 {
				t.Errorf("expected areacode with length of 3, got %s of length %d\n", result, len(result))
			}
		})
	}
}

// checks if number is of correct length
func TestNumber(t *testing.T) {
	tc := testCase{
		min: 1000000, max: 10000000,
		testsAmount: 10000,
	}

	for i := 0; i < tc.testsAmount; i++ {
		t.Run("", func(t *testing.T) {
			result := randomInt(tc.min, tc.max)

			if len(result) != 7 {
				t.Errorf("expected number with length of 7, got %s of length %d\n", result, len(result))
			}
		})
	}
}

