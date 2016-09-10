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
	assertFlagState(t, "--HC", r.flagToString())

	// test unsetting the flags
	r.resetFlag(flagC)
	assertFlagState(t, "--H-", r.flagToString())
}

func TestJumpOnCondition(t *testing.T) {
	// if a condition is true a jump should happen
	c := mockCPU()
	c.programCounter = 0x000A
	c.jumpOnCondition(0xFB, true)
	testhelpers.AssertWord(t, 0x0007, c.programCounter)

	// if a condition is false a jump shouldn't not happen and the pc should just be incremented
	c = mockCPU()
	c.programCounter = 0x000A
	c.jumpOnCondition(0xFB, false)
	testhelpers.AssertWord(t, 0x000C, c.programCounter)
}

func TestIncrementRegister(t *testing.T) {
	c := mockCPU()
	c.registers.BC.low = 0xF
	c.registers.BC.high = 0xFF
	c.registers.DE.low = 0x01

	c.incrementRegister(&c.registers.BC.low)
	testhelpers.AssertByte(t, 0x10, c.registers.BC.low)
	assertFlagState(t, "--H-", c.registers.flagToString())

	c.incrementRegister(&c.registers.BC.high)
	testhelpers.AssertByte(t, 0x0, c.registers.BC.high)
	assertFlagState(t, "Z-H-", c.registers.flagToString())

	c.incrementRegister(&c.registers.DE.low)
	testhelpers.AssertByte(t, 0x02, c.registers.DE.low)
	assertFlagState(t, "----", c.registers.flagToString())
}

func TestDecrementRegister(t *testing.T) {
	c := mockCPU()
	c.registers.BC.low = 0xF
	c.registers.BC.high = 0x10
	c.registers.DE.low = 0x01

	c.decrementRegister(&c.registers.BC.low)
	testhelpers.AssertByte(t, 0x0E, c.registers.BC.low)
	assertFlagState(t, "-N--", c.registers.flagToString())

	c.decrementRegister(&c.registers.BC.high)
	testhelpers.AssertByte(t, 0x0F, c.registers.BC.high)
	assertFlagState(t, "-NH-", c.registers.flagToString())

	c.decrementRegister(&c.registers.DE.low)
	testhelpers.AssertByte(t, 0x00, c.registers.DE.low)
	assertFlagState(t, "ZN--", c.registers.flagToString())
}

func TestLdIntoMemAndDec(t *testing.T) {
	c := mockCPU()
	c.registers.HL.setWord(0x1132)
	c.registers.BC.setWord(0x1030)
	c.registers.AF.low = 35

	c.ldIntoRegisterPairAddress(&c.registers.HL, c.registers.AF.low)
	testhelpers.AssertByte(t, 35, c.mmu.ReadByte(c.registers.HL.word()))

	prevPairValue := c.registers.BC.word()
	c.ldIntoRegisterPairAddressAndDec(&c.registers.BC, c.registers.AF.low)
	testhelpers.AssertByte(t, 35, c.mmu.ReadByte(prevPairValue))
	if c.registers.BC.word() != prevPairValue-1 {
		t.Error("Expected BC register pair to be decremented")
	}
}

func TestPushByteOntoStack(t *testing.T) {
	c := mockCPU()
	c.stackPointer = 0xFFFE
	c.pushByteOntoStack(0x04)

	testhelpers.AssertByte(t, 0x04, c.mmu.ReadByte(0xFFFD))
	testhelpers.AssertWord(t, 0xFFFD, c.stackPointer)
}

func TestPushWordOntoStack(t *testing.T) {
	c := mockCPU()
	c.stackPointer = 0xFFFE
	c.pushWordOntoStack(0x9432)

	testhelpers.AssertByte(t, 0x94, c.mmu.ReadByte(0xFFFD))
	testhelpers.AssertByte(t, 0x32, c.mmu.ReadByte(0xFFFC))
	testhelpers.AssertWord(t, 0xFFFC, c.stackPointer)
}

func TestPopStackIntoRegisterPair(t *testing.T) {
	c := mockCPU()
	c.stackPointer = 0xFFFE
	c.pushWordOntoStack(0x1234)
	c.popStackIntoRegisterPair(&c.registers.BC)

	testhelpers.AssertWord(t, 0x1234, c.registers.BC.word())
	testhelpers.AssertWord(t, 0xFFFE, c.stackPointer)
}

func TestCall(t *testing.T) {
	c := mockCPU()
	c.programCounter = 0x28
	c.stackPointer = 0xFFFE
	c.call(0x1234)

	testhelpers.AssertWord(t, c.stackPointer, 0xFFFC)
	testhelpers.AssertByte(t, 0x00, c.mmu.ReadByte(0xFFFD))
	testhelpers.AssertByte(t, 0x2b, c.mmu.ReadByte(0xFFFC))
	testhelpers.AssertWord(t, 0x1234, c.programCounter)
}

func TestRotateRegisterLeft(t *testing.T) {
	c := mockCPU()
	c.registers.setFlag(flagC)
	c.registers.BC.setWord(0x00d3)
	c.rotateRegisterLeft(&c.registers.BC.high)
	testhelpers.AssertByte(t, 0xA7, c.registers.BC.high)
	assertFlagState(t, "---C", c.registers.flagToString())

	c = mockCPU()
	c.registers.BC.setWord(0x0080)
	c.rotateRegisterLeft(&c.registers.BC.high)
	assertFlagState(t, "Z--C", c.registers.flagToString())

	c = mockCPU()
	c.registers.BC.setWord(0x0000)
	c.rotateRegisterLeft(&c.registers.BC.high)
	assertFlagState(t, "Z---", c.registers.flagToString())
}
