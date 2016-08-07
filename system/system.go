package system

import (
	"github.com/robertkrimen/otto"
	"github.com/robmerrell/gmboy/system/cpu"
	"github.com/robmerrell/gmboy/system/debugger"
	"github.com/robmerrell/gmboy/system/mmu"
	"github.com/robmerrell/gmboy/system/ui"
	"log"
)

const (
	displayWidth  = 160
	displayHeight = 144
)

// System represents the Gameboy system as a whole
type System struct {
	cpu        *cpu.CPU
	mmu        *mmu.MMU
	display    *ui.Display
	inputState *ui.InputState
	debugger   *debugger.Debugger
}

// NewSystem creates a new Gameboy system
func NewSystem() (*System, error) {
	m := mmu.NewMMU()
	c := cpu.NewCPU(m)

	d, err := ui.NewDisplay(displayWidth, displayHeight, 1)
	if err != nil {
		return &System{}, err
	}

	i := ui.NewInput(d)

	return &System{cpu: c, mmu: m, display: d, inputState: i}, nil
}

// PerformBootstrap runs the given bootstrap rom on startup. I'm unclear on copyright issues with this, so
// to be safe you will need to provide your own when bootstrapping.
func (s *System) PerformBootstrap(romFile string) error {
	if err := s.mmu.LoadBootRom(romFile); err != nil {
		return err
	}

	s.cpu.InitWithBoot()
	return nil
}

// LoadRom loads the given rom file into memory
func (s *System) LoadRom(romFile string) {
}

// Run runs the system
func (s *System) Run() {
	for {
		if s.debugger.BreakpointActive {
			s.stepWithBreakpoint()
		} else {
			s.step()
		}
	}
}

// step executes an instruction
func (s *System) step() {
	s.cpu.Step()
	s.display.PollOSEvents()
}

// stepWithDebugger executes an instruction and waits for input from the debugger
func (s *System) stepWithBreakpoint() {
	cont := false
	for !cont {
		select {
		case <-s.debugger.Step:
			s.cpu.Step()
		case <-s.debugger.Cont:
			s.debugger.BreakpointActive = false
			cont = true
		default:
			s.display.PollOSEvents()
		}
	}
}

// StartDebugger creates a new debugger and then attaches it to all of the relevant subsystems.
func (s *System) StartDebugger(file string) error {
	dbg := debugger.NewDebugger()

	// create the breakpoint() function for the js debugger that will allow us to step through execution
	dbg.AttachFunction("breakpoint", func(call otto.FunctionCall) otto.Value {
		log.Println("Breakpoint reached")
		dbg.BreakpointActive = true
		return otto.Value{}
	})

	s.cpu.AttachDebugger(dbg)
	s.mmu.AttachDebugger(dbg)
	s.inputState.AttachDebugger(dbg)

	err := dbg.LoadSourceFile(file)
	if err != nil {
		return err
	}
	s.debugger = dbg

	return nil
}
