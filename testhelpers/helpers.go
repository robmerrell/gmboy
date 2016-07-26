package testhelpers

import (
	"testing"
)

// AssertWord asserts the value of a word
func AssertWord(t *testing.T, expected, actual uint16) {
	if expected != actual {
		t.Errorf("Expected value was 0x%4x, but was 0x%4x\n", expected, actual)
	}
}

// AssertByte asserts the value of a byte
func AssertByte(t *testing.T, expected, actual byte) {
	if expected != actual {
		t.Errorf("Expected value was 0x%2x, but was 0x%2x\n", expected, actual)
	}
}
