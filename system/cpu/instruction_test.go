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

func Test0x31(t *testing.T) {
	c := mockCPU()
	c.mmu.WriteBytes([]byte{0x31, 0xFE, 0xFF}, 0)

	c.Step()
	testhelpers.AssertWord(t, 0xFFFE, c.stackPointer)
}

func Test0xAF(t *testing.T) {
	c := mockCPU()
	c.registers.AF.setWord(0xFFFE)
	c.mmu.WriteBytes([]byte{0xAF}, 0)

	c.Step()
	testhelpers.AssertByte(t, 0x00, c.registers.AF.low)
}
