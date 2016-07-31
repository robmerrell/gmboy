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

	// The length the operands for the instruction. Some instructions require more bytes for their parameters
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
	0x0E: &instruction{0x0E, "LD C,d8", 8, 2, false, func(c *CPU) { c.registers.BC.high = c.operandByte() }},
	// TODO: I'm not entirely sure how to handle cycles for these jump calls. It looks like they can vary in some way...
	0x20: &instruction{0x20, "JR NZ,r8", 8, 2, true, func(c *CPU) {
		cond := *c.registers.flag()&flagZ == 0 // Z flag is unset
		c.jumpOnCondition(c.operandByte(), cond)
	}},
	0x21: &instruction{0x21, "LD HL,d16", 12, 3, false, func(c *CPU) { c.registers.HL.setWord(c.operandWord()) }},
	0x31: &instruction{0x31, "LD SP,d16", 12, 3, false, func(c *CPU) { c.stackPointer = c.operandWord() }},
	0x32: &instruction{0x32, "LD (HL-),A", 8, 1, false, func(c *CPU) { c.ldIntoMemAndDec(&c.registers.HL, c.registers.AF.low) }},
	0xAF: &instruction{0xAF, "XOR A", 4, 1, false, func(c *CPU) { c.xorRegisters(&c.registers.AF.low, c.registers.AF.low) }},
}

// The instruction length for thee extended instructions is going to be what is in the above link-1. I believe that in the link above when they
// say the instruction length is 2 it's because they are counting botht he 0xCB + the instruction.
var extendedInstructions = map[byte]*instruction{
	0x7c: &instruction{0x7C, "BIT 7,H", 8, 1, false, func(c *CPU) { c.testRegisterBit(c.registers.HL.low, 7) }},
}
