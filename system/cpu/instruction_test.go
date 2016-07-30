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

// This is much easier and more human readable than checking the bits for our tests like I was doing before.
func assertFlagState(t *testing.T, expectedFlagString string, actualFlagString string) {
	if actualFlagString != expectedFlagString {
		t.Errorf("Expected the flag state to be %s, but is was %s", expectedFlagString, actualFlagString)
	}
}

func Test0x20(t *testing.T) {
	c := mockCPU()
	c.programCounter = 0x0000
	c.registers.resetFlag(flagZ)
	c.mmu.WriteBytes([]byte{0x20, 0x06}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0x0008, c.programCounter)
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
	assertFlagState(t, "Z---", c.registers.flagToString())
}

func Test0xCB7C(t *testing.T) {
	c := mockCPU()
	c.registers.HL.setWord(0xFFFF)
	c.mmu.WriteBytes([]byte{0xCB, 0x7C}, 0)
	c.Step()
	assertFlagState(t, "--H-", c.registers.flagToString())

	c = mockCPU()
	c.registers.HL.setWord(0x0000)
	c.mmu.WriteBytes([]byte{0xCB, 0x7C}, 0)
	c.Step()
	assertFlagState(t, "Z-H-", c.registers.flagToString())
}
