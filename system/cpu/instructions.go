package cpu

import (
	"fmt"
)

type instruction struct {
	// The machine form of the instruction. This will be duplicated in the map key below, but is useful here
	// for the debugger.
	opcode byte

	// The human readable form of the instruction. Great for debugging and inspecting things.
	mnemonic string

	// Each kind of opcode consumes a different amount of CPU cycles. For example 8-bit register loads
	// take 8 cycles, while 16-bit register loads take 12 cycles.
	cycles int

	// The length the operands + instruction. Some instructions require more bytes for their parameters
	// than others. For example compare NoOp (0x00) with loading a value in the stack pointer (0x31). The
	// NoOp instruction takes 0 operands and so the entire instruction consumes 1 byte, just for the instruction itself.
	// On the other hand loading a 16-bit word into the stack pointer requires space for the instruction (0x31), plus
	// 2 more bytes for the 16-bit word.
	len uint16

	// Some operations change the program counter and in these cases the system shouldn't try to advance the program
	// counter in its usual manner. This flags those operations.
	changesProgramCounter bool

	// The function to execute for the instruction.
	fn func(*CPU)
}

func (i *instruction) Debug() map[string]interface{} {
	return map[string]interface{}{
		"opcode":    i.opcode,
		"opcodeHex": fmt.Sprintf("0x%2x", i.opcode),
		"mnemonic":  i.mnemonic,
		"cycles":    i.cycles,
		"len":       i.len,
	}

}

// I'm going off of this list for the info about the opcodes including the mnemonic: http://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html

var baseInstructions = map[byte]*instruction{
	0x00: &instruction{0x00, "NOP", 4, 1, false, func(c *CPU) {}},
	0x01: &instruction{0x01, "LD BC,d16", 12, 3, false, func(c *CPU) { c.registers.BC.setWord(c.operandWord()) }},
	0x05: &instruction{0x05, "DEC B", 4, 1, false, func(c *CPU) { c.decrementRegister(&c.registers.BC.low) }},
	0x06: &instruction{0x06, "LD B,d8", 8, 2, false, func(c *CPU) { c.registers.BC.low = c.operandByte() }},
	0x0C: &instruction{0x0C, "INC C", 4, 1, false, func(c *CPU) { c.incrementRegister(&c.registers.BC.high) }},
	0x0E: &instruction{0x0E, "LD C,d8", 8, 2, false, func(c *CPU) { c.registers.BC.high = c.operandByte() }},
	0x11: &instruction{0x11, "LD DE,d16", 12, 3, false, func(c *CPU) { c.registers.DE.setWord(c.operandWord()) }},
	0x13: &instruction{0x13, "INC DE", 8, 1, false, func(c *CPU) { c.registers.DE.setWord(c.registers.DE.word() + 1) }},
	0x1A: &instruction{0x1A, "LD A,(DE)", 8, 1, false, func(c *CPU) { c.registers.AF.low = c.mmu.ReadByte(c.registers.DE.word()) }},
	// TODO: I'm not entirely sure how to handle cycles for these jump calls. It looks like they can vary in some way...
	0x20: &instruction{0x20, "JR NZ,r8", 8, 2, true, func(c *CPU) {
		cond := *c.registers.flag()&flagZ == 0 // Z flag is unset
		c.jumpOnCondition(c.operandByte(), cond)
	}},
	0x21: &instruction{0x21, "LD HL,d16", 12, 3, false, func(c *CPU) { c.registers.HL.setWord(c.operandWord()) }},
	0x22: &instruction{0x22, "LD (HL+),A", 8, 1, false, func(c *CPU) { c.ldIntoRegisterPairAddressAndInc(&c.registers.HL, c.registers.AF.low) }},
	0x23: &instruction{0x23, "INC HL", 8, 1, false, func(c *CPU) { c.registers.HL.setWord(c.registers.HL.word() + 1) }},
	0x31: &instruction{0x31, "LD SP,d16", 12, 3, false, func(c *CPU) { c.stackPointer = c.operandWord() }},
	0x32: &instruction{0x32, "LD (HL-),A", 8, 1, false, func(c *CPU) { c.ldIntoRegisterPairAddressAndDec(&c.registers.HL, c.registers.AF.low) }},
	0x4F: &instruction{0x4F, "LD C,A", 4, 1, false, func(c *CPU) { c.registers.BC.high = c.registers.AF.low }},
	0x3E: &instruction{0x3E, "LD A,d8", 8, 2, false, func(c *CPU) { c.registers.AF.low = c.operandByte() }},
	0x77: &instruction{0x77, "LD (HL),A", 8, 1, false, func(c *CPU) { c.ldIntoRegisterPairAddress(&c.registers.HL, c.registers.AF.low) }},
	0xAF: &instruction{0xAF, "XOR A", 4, 1, false, func(c *CPU) { c.xorRegisters(&c.registers.AF.low, c.registers.AF.low) }},
	0xC1: &instruction{0xC1, "POP BC", 12, 1, false, func(c *CPU) { c.popStackIntoRegisterPair(&c.registers.BC) }},
	0xC5: &instruction{0xC5, "PUSH BC", 16, 1, false, func(c *CPU) { c.pushWordOntoStack(c.registers.BC.word()) }},
	0xC9: &instruction{0xC9, "RET", 16, 1, true, func(c *CPU) { c.ret() }},
	0xCD: &instruction{0xCD, "CALL a16", 24, 3, true, func(c *CPU) { c.call(c.operandWord()) }},
	0xE0: &instruction{0xE0, "LDH (a8),A", 12, 2, false, func(c *CPU) { c.mmu.WriteBytes([]byte{c.registers.AF.low}, 0xFF00+uint16(c.operandByte())) }},
	0xE2: &instruction{0xE2, "LD A,(C)", 8, 1, false, func(c *CPU) { c.mmu.WriteBytes([]byte{c.registers.AF.low}, 0xFF00+uint16(c.registers.BC.high)) }},
}

// The instruction length for the extended instructions is going to be what is in the above link-1. I believe that in the link above when they
// say the instruction length is 2 it's because they are counting both the 0xCB byte + the instruction.
var extendedInstructions = map[byte]*instruction{
	0x11: &instruction{0x11, "RL C", 8, 2, false, func(c *CPU) { c.rotateRegisterLeft(&c.registers.BC.high) }},
	0x7C: &instruction{0x7C, "BIT 7,H", 8, 1, false, func(c *CPU) { c.testRegisterBit(c.registers.HL.low, 7) }},
}
