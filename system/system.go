package system

import (
	"github.com/robmerrell/gmboy/system/cpu"
	"github.com/robmerrell/gmboy/system/mmu"
)

// System represents the Gameboy system as a whole
type System struct {
	cpu *cpu.CPU
	mmu *mmu.MMU
}

// NewSystem creates a new Gameboy system
func NewSystem() *System {
	m := mmu.NewMMU()
	c := cpu.NewCPU(m)
	return &System{cpu: c, mmu: m}
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
	}
}
