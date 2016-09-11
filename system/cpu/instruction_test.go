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

func Test0x01(t *testing.T) {
	c := mockCPU()
	c.mmu.WriteBytes([]byte{0x01, 0xFE, 0xFF}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0xFFFE, c.registers.BC.word())
}

func Test0x05(t *testing.T) {
	c := mockCPU()
	c.registers.BC.low = 0x10
	c.mmu.WriteBytes([]byte{0x05}, 0)

	c.Step()
	testhelpers.AssertByte(t, 0x0F, c.registers.BC.low)
	assertFlagState(t, "-NH-", c.registers.flagToString())
}

func Test0x06(t *testing.T) {
	c := mockCPU()
	c.mmu.WriteBytes([]byte{0x06, 0x04}, 0)
	c.Step()

	testhelpers.AssertByte(t, 0x04, c.registers.BC.low)
}

func Test0x0C(t *testing.T) {
	c := mockCPU()
	c.registers.BC.high = 0xF
	c.mmu.WriteBytes([]byte{0x0C}, 0)
	c.Step()

	testhelpers.AssertByte(t, 0x10, c.registers.BC.high)
	assertFlagState(t, "--H-", c.registers.flagToString())
}

func Test0x0E(t *testing.T) {
	c := mockCPU()
	c.mmu.WriteBytes([]byte{0x0E, 0x04}, 0)
	c.Step()

	testhelpers.AssertByte(t, 0x04, c.registers.BC.high)
}

func Test0x11(t *testing.T) {
	c := mockCPU()
	c.mmu.WriteBytes([]byte{0x11, 0xFE, 0xFF}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0xFFFE, c.registers.DE.word())
}

func Test0x13(t *testing.T) {
	c := mockCPU()
	c.registers.DE.setWord(0x9FFE)
	c.mmu.WriteBytes([]byte{0x13}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0x9FFF, c.registers.DE.word())
}

func Test0x1A(t *testing.T) {
	c := mockCPU()
	c.registers.DE.setWord(0xFEFF)
	c.mmu.WriteBytes([]byte{0x04}, 0xFEFF)
	c.mmu.WriteBytes([]byte{0x1A}, 0)

	c.Step()
	testhelpers.AssertByte(t, 0x04, c.registers.AF.low)
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

func Test0x22(t *testing.T) {
	c := mockCPU()
	c.registers.AF.setWord(0x1132)
	c.registers.HL.setWord(0x9FFE)
	c.mmu.WriteBytes([]byte{0x22}, 0)

	c.Step()

	testhelpers.AssertWord(t, 0x9FFF, c.registers.HL.word())
	testhelpers.AssertByte(t, 0x11, c.mmu.ReadByte(c.registers.HL.word()-1))
}

func Test0x23(t *testing.T) {
	c := mockCPU()
	c.registers.HL.setWord(0x9FFE)
	c.mmu.WriteBytes([]byte{0x23}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0x9FFF, c.registers.HL.word())
}

func Test0x31(t *testing.T) {
	c := mockCPU()
	c.mmu.WriteBytes([]byte{0x31, 0xFE, 0xFF}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0xFFFE, c.stackPointer)
}

func Test0x32(t *testing.T) {
	c := mockCPU()
	c.registers.AF.setWord(0x1132)
	c.registers.HL.setWord(0x9FFF)
	c.mmu.WriteBytes([]byte{0x32}, 0)

	c.Step()

	testhelpers.AssertWord(t, 0x9FFE, c.registers.HL.word())
	testhelpers.AssertByte(t, 0x11, c.mmu.ReadByte(c.registers.HL.word()+1))
}

func Test0x3E(t *testing.T) {
	c := mockCPU()
	c.mmu.WriteBytes([]byte{0x3E, 0x04}, 0)
	c.Step()

	testhelpers.AssertByte(t, 0x04, c.registers.AF.low)
}

func Test0x4F(t *testing.T) {
	c := mockCPU()
	c.registers.AF.low = 0x12
	c.registers.BC.high = 0x04
	c.mmu.WriteBytes([]byte{0x4f}, 0)
	c.Step()

	testhelpers.AssertByte(t, 0x12, c.registers.BC.high)
}

func Test0x77(t *testing.T) {
	c := mockCPU()
	c.registers.AF.setWord(0x1132)
	c.registers.HL.setWord(0x9FFF)
	c.mmu.WriteBytes([]byte{0x77}, 0)

	c.Step()

	testhelpers.AssertByte(t, 0x11, c.mmu.ReadByte(c.registers.HL.word()))
}

func Test0xAF(t *testing.T) {
	c := mockCPU()
	c.registers.AF.setWord(0xFFFE)
	c.mmu.WriteBytes([]byte{0xAF}, 0)

	c.Step()
	testhelpers.AssertByte(t, 0x00, c.registers.AF.low)
	assertFlagState(t, "Z---", c.registers.flagToString())
}

func Test0xC1(t *testing.T) {
	c := mockCPU()
	c.stackPointer = 0xFFFE
	c.pushWordOntoStack(0x1234)
	c.mmu.WriteBytes([]byte{0xC1}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0x1234, c.registers.BC.word())
}

func Test0xC5(t *testing.T) {
	c := mockCPU()
	c.stackPointer = 0xFFFE
	c.registers.BC.setWord(0x1234)
	c.mmu.WriteBytes([]byte{0xC5}, 0)

	c.Step()
	testhelpers.AssertByte(t, 0x12, c.mmu.ReadByte(0xFFFD))
	testhelpers.AssertByte(t, 0x34, c.mmu.ReadByte(0xFFFC))
}

func Test0xC9(t *testing.T) {
	c := mockCPU()
	c.stackPointer = 0xFFFE
	c.pushWordOntoStack(0x1123)
	c.mmu.WriteBytes([]byte{0xC9}, 0)

	c.Step()

	testhelpers.AssertWord(t, 0x1123, c.programCounter)
	testhelpers.AssertWord(t, 0xFFFE, c.stackPointer)
}

func Test0xCD(t *testing.T) {
	c := mockCPU()
	c.stackPointer = 0xFFFE
	c.mmu.WriteBytes([]byte{0xCD, 0x95, 0x00, 0x01}, 0x28)
	c.programCounter = 0x28

	c.Step()

	testhelpers.AssertWord(t, 0x95, c.programCounter)
	testhelpers.AssertByte(t, 0x0, c.mmu.ReadByte(0xFFFD))
	testhelpers.AssertByte(t, 0x2b, c.mmu.ReadByte(0xFFFC))
}

func Test0xE0(t *testing.T) {
	c := mockCPU()
	c.registers.AF.setWord(0x1203)
	c.mmu.WriteBytes([]byte{0xE0, 0x04}, 0)

	c.Step()
	testhelpers.AssertByte(t, c.registers.AF.low, c.mmu.ReadByte(0xFF04))
}

func Test0xE2(t *testing.T) {
	c := mockCPU()
	c.registers.AF.setWord(0x800a)
	c.registers.BC.setWord(0x0011)
	c.mmu.WriteBytes([]byte{0xE2}, 0)

	c.Step()
	testhelpers.AssertByte(t, c.registers.AF.low, c.mmu.ReadByte(0xFF11))
}

func Test0xCB11(t *testing.T) {
	c := mockCPU()
	c.registers.BC.setWord(0x0080)
	c.mmu.WriteBytes([]byte{0xCB, 0x11}, 0)

	c.Step()
	testhelpers.AssertByte(t, 0x00, c.registers.BC.high)
	assertFlagState(t, "Z--C", c.registers.flagToString())
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
