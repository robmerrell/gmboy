package cpu

import (
	"encoding/binary"
	"github.com/robertkrimen/otto"
	"github.com/robmerrell/gmboy/system/debugger"
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

	// debugger
	debugger       *debugger.Debugger
	debuggerActive bool
}

// NewCPU returns a new CPU instance
func NewCPU(memoryManager *mmu.MMU) *CPU {
	return &CPU{registers: &registers{}, mmu: memoryManager}
}

// AttachDebugger attaches a javascript debugger to the CPU
func (c *CPU) AttachDebugger(dbg *debugger.Debugger) {
	log.Println("Attaching debugger to CPU")
	c.debugger = dbg
	c.debuggerActive = true

	// create the cpuState() function for the js debugger that returns the current state of the CPU
	c.debugger.AttachFunction("cpuState", func(call otto.FunctionCall) otto.Value {
		registers := map[string]interface{}{
			"stackPointer":   c.stackPointer,
			"programCounter": c.programCounter,
			"registers": map[string]byte{
				"A": c.registers.AF.low,
				"B": c.registers.BC.low,
				"C": c.registers.BC.high,
				"D": c.registers.DE.low,
				"E": c.registers.DE.high,
				"F": c.registers.AF.high,
				"H": c.registers.HL.low,
				"L": c.registers.HL.high,
			},
			"registerPairs": map[string]uint16{
				"AF": c.registers.AF.word(),
				"BC": c.registers.BC.word(),
				"DE": c.registers.DE.word(),
				"HL": c.registers.HL.word(),
			},
		}
		val, _ := call.Otto.ToValue(registers)
		return val
	})
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
		if c.debuggerActive {
			c.debugger.RunCallbacks("unimplemented_opcode", opcode)
		}
		return
	}

	if c.debuggerActive {
		c.debugger.RunCallbacks("before_execute", instruction.Debug())
	}

	// execute the instruction
	instruction.fn(c)

	// advance the program counter
	if !instruction.changesProgramCounter {
		c.programCounter += instruction.len
	}

	if c.debuggerActive {
		c.debugger.RunCallbacks("after_execute", instruction.Debug())
	}
}

// operandWord reads and returns the current instructions operands as a word
func (c *CPU) operandWord() uint16 {
	return c.mmu.ReadWord(c.programCounter + 1)
}

// xorRegister xor's a source and operand register and saves it in the source register
func xorRegister(sourceRegister *byte, operandRegister byte) {
	*sourceRegister ^= operandRegister
}

// ldIntoMemAndDec loads the value of the copyRegister into the memory address stored in the ldRegister
// and then decrements the value stored in ldRegister
func (c *CPU) ldIntoMemAndDec(ldRegister *register, copyRegister byte) {
	address := ldRegister.word()
	c.mmu.WriteBytes([]byte{copyRegister}, address)
	ldRegister.setWord(ldRegister.word() - 1)
}
