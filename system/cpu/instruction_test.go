package cpu

import (
	"github.com/robmerrell/gmboy/system/mmu"
	"github.com/robmerrell/gmboy/testhelpers"
	"testing"
)

func mockCPU() *CPU {
	m := mmu.NewMMU()
	c := NewCPU(m)
	c.programCounter = 0x00

	return c
}

func assertSetFlags(t *testing.T, flagRegister byte, bitmask byte) {
	checkFlag := func(flag byte) {
		if flagRegister&flag == 0 {
			t.Errorf("Expected the flag %08b to be set, but it was not", flag)
		}
	}

	if bitmask&flagC != 0 {
		checkFlag(flagC)
	}

	if bitmask&flagH != 0 {
		checkFlag(flagH)
	}

	if bitmask&flagN != 0 {
		checkFlag(flagN)
	}

	if bitmask&flagZ != 0 {
		checkFlag(flagZ)
	}
}

func assertUnsetFlags(t *testing.T, flagRegister byte, bitmask byte) {
	checkFlag := func(flag byte) {
		if flagRegister&flag != 0 {
			t.Errorf("Expected the flag %08b to be reset, but it had a value", flag)
		}
	}

	if bitmask&flagC != 0 {
		checkFlag(flagC)
	}

	if bitmask&flagH != 0 {
		checkFlag(flagH)
	}

	if bitmask&flagN != 0 {
		checkFlag(flagN)
	}

	if bitmask&flagZ != 0 {
		checkFlag(flagZ)
	}
}

func Test0x21(t *testing.T) {
	c := mockCPU()
	c.mmu.WriteBytes([]byte{0x21, 0xFE, 0xFF}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0xFFFE, c.registers.HL.word())
}

func Test0x31(t *testing.T) {
	c := mockCPU()
	c.mmu.WriteBytes([]byte{0x31, 0xFE, 0xFF}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0xFFFE, c.stackPointer)
}

func Test0x32(t *testing.T) {
	c := mockCPU()
	c.registers.AF.setWord(0x0000)
	c.registers.HL.setWord(0x9FFF)
	c.mmu.WriteBytes([]byte{0x32}, 0)

	c.Step()

	testhelpers.AssertWord(t, 0x9FFE, c.registers.HL.word())
}

func Test0xAF(t *testing.T) {
	c := mockCPU()
	c.registers.AF.setWord(0xFFFE)
	c.mmu.WriteBytes([]byte{0xAF}, 0)

	c.Step()
	testhelpers.AssertByte(t, 0x00, c.registers.AF.low)
	assertSetFlags(t, *c.registers.flag(), flagZ)
	assertUnsetFlags(t, *c.registers.flag(), flagC|flagH|flagN)
}
