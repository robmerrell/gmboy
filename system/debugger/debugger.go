package debugger

import (
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
)

// This needs alot more documentation and examples, because it's a cool feature!

// Debugger is a javascript enabled debugger that allows a user to receive events emitted by the emulator,
// access and change CPU state, and access and change memory.
//
// The current events available to the javascript debugger are:
//   before_execute: fired before an opcode is executed. Passes the current instruction.
//   after_execute: fired after an opcode is executed. Passes the just executed instruction.
//   unimplemented_opcode: fired when an unimplemented opcode is encountered. Passes the opcode.
//
// Builtin functions:
//   dumpMemory() - returns an array of the system's memory
//   writeByte(location, byte) - write a byte at the given location in memory
//   cpuState() - returns an object with the current state of the CPU
//   ppSystem() - pretty prints the current system state
//   ppCPU() - pretty prints the current CPU state
//   ppInstruction(inst) - pretty prints an instruction
type Debugger struct {
	vm                *otto.Otto
	stopOnInstruction bool
	callbacks         map[string][]otto.Value

	// Flag for when we are actively stepping from a breakpoint
	BreakpointActive bool

	// Step is what we use to coordinate stepping through execution with the CPU when a breakpoint is encountered.
	Step chan bool

	// Cont signals that execution should continue and the breakpoint be removed.
	Cont chan bool
}

// NewDebugger returns a new debugger instance.
func NewDebugger() *Debugger {
	d := &Debugger{
		vm:        otto.New(),
		callbacks: make(map[string][]otto.Value),
		Step:      make(chan bool, 1),
		Cont:      make(chan bool, 1),
	}

	// Make the "on" function available to the js vm
	d.vm.Set("on", func(call otto.FunctionCall) otto.Value {
		name, _ := call.Argument(0).ToString()
		d.callbacks[name] = append(d.callbacks[name], call.Argument(1))
		return otto.Value{}
	})

	// add the pretty print functions
	d.vm.Run(prettPrintSrc)

	return d
}

// LoadSourceFile loads a javascript source file and executes it.
func (d *Debugger) LoadSourceFile(filename string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	d.vm.Run(contents)
	return nil
}

// AttachFunction allows all parts of the emulator to attach their own functions to the javascript vm. This is so
// that the inner parts of each sub-system (CPU, MMU, etc) don't have to be publically exposed to be able to debug.
// It does have the drawback of all javascript functions being defined from different places instead of one place.
func (d *Debugger) AttachFunction(name string, fn func(otto.FunctionCall) otto.Value) {
	d.vm.Set(name, fn)
}

// RunCallbacks runs all of a given type of callbacks. Passing arg as an argument to the javascript callback.
func (d *Debugger) RunCallbacks(name string, arg interface{}) {
	callbacks, exists := d.callbacks[name]
	if !exists {
		return
	}

	for _, callback := range callbacks {
		if callback.IsFunction() {
			_, err := callback.Call(otto.Value{}, arg)
			if err != nil {
				log.Println(err)
			}
		} else {
			val, _ := callback.ToString()
			log.Println("Expected a function for the callback, but got", val)
		}
	}
}

// Next sends a message to the system that we are ready to move to the next statement
func (d *Debugger) Next() {
	if d.BreakpointActive {
		d.Step <- true
	}
}

// Continue sends a message to the system that we are ready to remove the breakpoint and continue execution.
func (d *Debugger) Continue() {
	if d.BreakpointActive {
		d.Cont <- true
	}
}

var prettPrintSrc = `
ppInstruction = function(inst) {
  var mem = dumpMemory();
  var cpu = cpuState();

  var operands = mem.slice(cpu.programCounter+1, cpu.programCounter+inst.len);
  operands = operands.map(function(i) {
    return '0x'+i.toString(16);
  }).join(' ');

  console.log(inst.opcodeHex + ': ' + inst.mnemonic + '    operands: ' + operands);
};

padHex = function(value, mask) {
  value = value.toString(16);
  return (mask + value).slice(-mask.length);
};

ppCPU = function() {
  var state = cpuState();
  var output = 'PC:' + padHex(state.programCounter, '0000') + '  ' +
        'SP:' + padHex(state.stackPointer, '0000') + '  ' +
        'A:' + padHex(state.registers.A, '00') + '  ' +
        'B:' + padHex(state.registers.B, '00') + '  ' +
        'C:' + padHex(state.registers.C, '00') + '  ' +
        'D:' + padHex(state.registers.D, '00') + '  ' +
        'E:' + padHex(state.registers.E, '00') + '  ' +
        'F:' + padHex(state.registers.F, '00') + '  ' +
        'H:' + padHex(state.registers.H, '00') + '  ' +
        'L:' + padHex(state.registers.L, '00') + '  ' +
        'AF:' + padHex(state.registerPairs.AF, '0000') + '  ' +
        'BC:' + padHex(state.registerPairs.BC, '0000') + '  ' +
        'DE:' + padHex(state.registerPairs.DE, '0000') + '  ' +
        'HL:' + padHex(state.registerPairs.HL, '0000') + '  ' +
        'flags:' + state.flags;

  console.log(output);
};

ppSystem = function() {
  var state = cpuState();

  console.log('-----------------------------------------------------');
  console.log('                      CPU                            ');
  console.log('-----------------------------------------------------');

  console.log('  program counter: ' + state.programCounter.toString(16));
  console.log('  stack pointer:   ' + state.stackPointer.toString(16));
  console.log('  flags:           ' + state.flags);

  var keys = Object.keys(state.registerPairs).sort();
  for (var i in keys) {
    var key = keys[i];
    console.log('  ' + key + ': ' + state.registerPairs[key].toString(16));
  }

  console.log("");

  keys = Object.keys(state.registers).sort();
  for (var i in keys) {
    var key = keys[i];
    console.log('  ' + key + ': ' + state.registers[key].toString(16));
  }

  console.log('=====================================================\n');
};
`
