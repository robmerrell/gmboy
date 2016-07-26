package cpu

type instruction struct {
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

// I'm going off of this list for the info about the opcodes including the mnemonic: http://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html

var baseInstructions = map[byte]*instruction{
	0x00: &instruction{"NOP", 4, 1, false, func(c *CPU) {}},
	0x31: &instruction{"LD SP,d16", 12, 3, false, func(c *CPU) { c.stackPointer = c.operandWord() }},
	0xAF: &instruction{"XOR A", 4, 1, false, func(c *CPU) { xorRegister(&c.registers.AF.low, c.registers.AF.low) }},
}
