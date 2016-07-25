package cpu

import (
	"encoding/binary"
	"github.com/robmerrell/gmboy/system/mmu"
	"log"
)

// register store the state of a register pair. The gameboy has 8 8-bit registers: A, B, C, D, E, F, H and L
// these registers are often accessed in 16-bit pairs (AF, BC, DE, HL). So by default we represent them as 16-bit
// pairs and if we need to access just a single 8-bit register we do so by accessing either the high or low byte of
// the word. For example: If the register pair is BC. B would be low and C high.
type register struct {
	low  byte
	high byte
}

// word returns the word represenation of the low and high registers together.
func (r *register) word() uint16 {
	return binary.LittleEndian.Uint16([]byte{r.high, r.low})
}

// setWord sets the low and high bits accordingly for the given word passed in.
func (r *register) setWord(word uint16) {
	r.low = byte((word & 0xFF00) >> 8)
	r.high = byte(word & 0x00FF)
}

// registers holds all of the CPU registers by name. The gameboy has 8 registers.
type registers struct {
	AF register // accumlator, flag and carry registers
	BC register
	DE register
	HL register
}

// CPU holds the current state of the CPU
type CPU struct {
	// registers are all of the working CPU registers
	registers *registers

	// stackPointer points to the current location on the call stack
	stackPointer uint16

	// programCounter keeps track of the next instruction to read from memory
	programCounter uint16

	// reference to the MMU
	mmu *mmu.MMU
}

// NewCPU returns a new CPU instance
func NewCPU(memoryManager *mmu.MMU) *CPU {
	return &CPU{registers: &registers{}, mmu: memoryManager}
}

// InitWithBoot initializes the CPU assuming we are executing the bootrom. The only absolute known of
// the CPU state when loading a bootrom is that the program counter should be at the start of the
// memory space. All other states are set by the bootrom itself.
func (c *CPU) InitWithBoot() {
	c.programCounter = 0x0000
}

// Step processes an instruction
func (c *CPU) Step() {
	// get the instruction of the opcode
	opcode := c.mmu.ReadByte(c.programCounter)
	instruction, exists := baseInstructions[opcode]
	if !exists {
		log.Printf("Opcode not yet implemented: 0x%2x\n", opcode)
		return
	}

	// execute the instruction
	instruction.fn(c)

	// advance the program counter
	if !instruction.changesProgramCounter {
		c.programCounter += instruction.len
	}
}

// operandWord reads and returns the current instructions operands as a word
func (c *CPU) operandWord() uint16 {
	return c.mmu.ReadWord(c.programCounter + 1)
}
