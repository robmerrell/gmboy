package cpu

import (
	"github.com/robmerrell/gmboy/testhelpers"
	"testing"
)

func TestRegisterWord(t *testing.T) {
	r := &register{0x32, 0x11}
	testhelpers.AssertWord(t, 0x3211, r.word())
}

func TestRegisterSetWord(t *testing.T) {
	r := &register{}
	r.setWord(0x3211)
	testhelpers.AssertByte(t, r.low, 0x32)
	testhelpers.AssertByte(t, r.high, 0x11)
}

func TestFlagRegisterIsF(t *testing.T) {
	r := &registers{}
	r.AF.high = 0x11

	testhelpers.AssertByte(t, 0x11, *r.flag())
}

func TestRegistersFlagSettingAndUnsetting(t *testing.T) {
	r := &registers{}
	r.setFlag(flagC)
	r.setFlag(flagH)

	if *r.flag()&flagC == 0 {
		t.Error("Expected C flag to be set, but it wasn't")
	}

	if *r.flag()&flagH == 0 {
		t.Error("Expected H flag to be set, but it wasn't")
	}

	if *r.flag()&flagZ != 0 {
		t.Error("Expected Z flag to not be set, but it was")
	}

	// test unsetting the flags
	r.resetFlag(flagC)

	if *r.flag()&flagC != 0 {
		t.Error("Expected C flag to not be set, but it was")
	}
}
