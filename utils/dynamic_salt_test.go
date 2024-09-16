package utils

import (
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		length         int
		expectedLength int
	}{
		{10, 10},
		{0, 0},
		{5, 5},
	}

	for _, tt := range tests {
		result, err := GenerateRandomString(tt.length)
		if err != nil {
			t.Errorf("GenerateRandomString(%d) returned an error: %v", tt.length, err)
			continue
		}
		if len(result) != tt.expectedLength {
			t.Errorf("GenerateRandomString(%d) = %s; expected length %d", tt.length, result, tt.expectedLength)
		}
		for _, char := range result {
			if !isValidChar(char) {
				t.Errorf("GenerateRandomString(%d) = %s; contains invalid character %c", tt.length, result, char)
			}
		}
	}
}

func isValidChar(c rune) bool {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, validChar := range charset {
		if c == validChar {
			return true
		}
	}
	return false
}
