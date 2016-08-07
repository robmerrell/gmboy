package cpu

import (
	"encoding/binary"
	"github.com/robertkrimen/otto"
	"github.com/robmerrell/gmboy/system/debugger"
	"github.com/robmerrell/gmboy/system/mmu"
	"log"
	"strings"
)

// These represent the flags for the flag register (register F). Some CPU operations will set or unset these
// flags on the flag register. And some will act differently when these flags are set.
const (
	flagC byte = 1 << 4
	flagH      = 1 << 5
	flagN      = 1 << 6
	flagZ      = 1 << 7
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

// flag is just a shortcut to get to the F register
func (r *registers) flag() *byte {
	return &r.AF.high
}

// resetFlag resets the given flag on the Flag register (F)
func (r *registers) resetFlag(flag byte) {
	r.AF.high &^= flag
}

//setFlag sets the given flag on the Flag register (F)
func (r *registers) setFlag(flag byte) {
	r.AF.high |= flag
}

// flagToString returns the flags in an easy to read format where the flags occupy one of
// four spaces in the string: ZNHC. If a flag is set, it's letter will be present, if not
// it will be zero.
func (r *registers) flagToString() string {
	flagStates := []string{"-", "-", "-", "-"}

	if *r.flag()&flagC != 0 {
		flagStates[3] = "C"
	}

	if *r.flag()&flagH != 0 {
		flagStates[2] = "H"
	}

	if *r.flag()&flagN != 0 {
		flagStates[1] = "N"
	}

	if *r.flag()&flagZ != 0 {
		flagStates[0] = "Z"
	}

	return strings.Join(flagStates, "")
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
			"flags": c.registers.flagToString(),
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

	var inst *instruction
	var exists bool
	if opcode == 0xCB { // extended instruction
		// since the extended opcode follows the 0xCB opcode increment the programCounter and read again
		c.programCounter++
		opcode = c.mmu.ReadByte(c.programCounter)

		inst, exists = extendedInstructions[opcode]
	} else {
		inst, exists = baseInstructions[opcode]
	}

	if !exists {
		if c.debuggerActive {
			c.debugger.RunCallbacks("unimplemented_opcode", opcode)
		}
		return
	}

	if c.debuggerActive {
		c.debugger.RunCallbacks("before_execute", inst.Debug())
	}

	// execute the instruction
	inst.fn(c)

	// advance the program counter
	if !inst.changesProgramCounter {
		c.programCounter += inst.len
	}

	if c.debuggerActive {
		c.debugger.RunCallbacks("after_execute", inst.Debug())
	}
}

// operandByte reads and returns the current instructions operand as a byte
func (c *CPU) operandByte() byte {
	return c.mmu.ReadByte(c.programCounter + 1)
}

// operandWord reads and returns the current instructions operands as a word
func (c *CPU) operandWord() uint16 {
	return c.mmu.ReadWord(c.programCounter + 1)
}

// xorRegisters XOR's a source and operand register and saves it in the source register. C,H,N flags are reset
// and the Z flag is set to 0 if the XOR results in a 0.
func (c *CPU) xorRegisters(sourceRegister *byte, operandRegister byte) {
	res := *sourceRegister ^ operandRegister

	c.registers.resetFlag(flagC)
	c.registers.resetFlag(flagH)
	c.registers.resetFlag(flagN)
	c.registers.resetFlag(flagZ)

	if res == 0 {
		c.registers.setFlag(flagZ)
	}

	*sourceRegister = res
}

// ldIntoMemAndDec loads the value of the copyRegister into the memory address stored in the ldRegister
// and then decrements the value stored in ldRegister
func (c *CPU) ldIntoMemAndDec(ldRegister *register, copyRegister byte) {
	address := ldRegister.word()
	c.mmu.WriteBytes([]byte{copyRegister}, address)
	ldRegister.setWord(ldRegister.word() - 1)
}

// testRegisterBit tests the bit in the register and sets the Z flag if that bit is 0
func (c *CPU) testRegisterBit(register byte, bitNum byte) {
	if register&(1<<bitNum) != 0 { // bit is set
		c.registers.resetFlag(flagZ)
	} else { // bit is not set
		c.registers.setFlag(flagZ)
	}

	c.registers.resetFlag(flagN)
	c.registers.setFlag(flagH)
}

// jumpOnCondition will jump if the condition is true and continue if not
func (c *CPU) jumpOnCondition(offset byte, condition bool) {
	// The jump assumes the program counter has alredy been incremented. And if no jump needs to happen
	// this needs to be incremented anyway to move onto the next instruction.
	c.programCounter += 2

	if condition {
		signedOffset := int8(offset)
		c.programCounter += uint16(signedOffset)
	}
}
