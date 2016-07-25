package cpu

import (
	"testing"
)

func assertWord(t *testing.T, expected, actual uint16) {
	if expected != actual {
		t.Errorf("Expected value was 0x%4x, but was 0x%4x\n", expected, actual)
	}
}

func assertByte(t *testing.T, expected, actual byte) {
	if expected != actual {
		t.Errorf("Expected value was 0x%2x, but was 0x%2x\n", expected, actual)
	}
}

func TestRegisterWord(t *testing.T) {
	r := &register{0x32, 0x11}
	assertWord(t, 0x3211, r.word())
}

func TestRegisterSetWord(t *testing.T) {
	r := &register{}
	r.setWord(0x3211)
	assertByte(t, r.low, 0x32)
	assertByte(t, r.high, 0x11)
}
