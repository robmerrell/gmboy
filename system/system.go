package system

import (
	"github.com/robmerrell/gmboy/system/cpu"
	"github.com/robmerrell/gmboy/system/debugger"
	"github.com/robmerrell/gmboy/system/display"
	"github.com/robmerrell/gmboy/system/mmu"
)

const (
	displayWidth  = 160
	displayHeight = 144
)

// System represents the Gameboy system as a whole
type System struct {
	cpu     *cpu.CPU
	mmu     *mmu.MMU
	display *display.Display
}

// NewSystem creates a new Gameboy system
func NewSystem() (*System, error) {
	m := mmu.NewMMU()
	c := cpu.NewCPU(m)

	d, err := display.NewDisplay(displayWidth, displayHeight, 1)
	if err != nil {
		return &System{}, err
	}

	return &System{cpu: c, mmu: m, display: d}, nil
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
		s.cpu.Step()
		s.display.PollOSEvents()
	}
}

// StartDebugger creates a new debugger and then attaches it to all of the relevant subsystems.
func (s *System) StartDebugger(file string) error {
	dbg := debugger.NewDebugger()
	err := dbg.LoadSourceFile(file)
	if err != nil {
		return err
	}

	s.cpu.AttachDebugger(dbg)
	s.mmu.AttachDebugger(dbg)
	return nil
}
