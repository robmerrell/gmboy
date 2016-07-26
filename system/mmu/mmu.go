package mmu

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
)

/*
The gameboy memory map (from Pan Docs http://bgb.bircd.org/pandocs.htm)

  0000-3FFF   16KB ROM Bank 00     (in cartridge, fixed at bank 00)
  4000-7FFF   16KB ROM Bank 01..NN (in cartridge, switchable bank number)
  8000-9FFF   8KB Video RAM (VRAM) (switchable bank 0-1 in CGB Mode)
  A000-BFFF   8KB External RAM     (in cartridge, switchable bank, if any)
  C000-CFFF   4KB Work RAM Bank 0 (WRAM)
  D000-DFFF   4KB Work RAM Bank 1 (WRAM)  (switchable bank 1-7 in CGB Mode)
  E000-FDFF   Same as C000-DDFF (ECHO)    (typically not used)
  FE00-FE9F   Sprite Attribute Table (OAM)
  FEA0-FEFF   Not Usable
  FF00-FF7F   I/O Ports
  FF80-FFFE   High RAM (HRAM)
  FFFF        Interrupt Enable Register
*/

const memorySize = 0xFFFF

// MMU is the memory management unit for gmboy. The gameboy hardware doesn't have an MMU
// but we're creating one here to make accessing memory easier to deal with.
type MMU struct {
	memory []byte
}

// NewMMU creates a new MMU to manage loading, accessing and changing values in memory.
func NewMMU() *MMU {
	return &MMU{memory: make([]byte, memorySize)}
}

// LoadBootRom loads the given bootrom file into memory. When the system boots the bootrom is
// placed at 0x0000. The bootrom is expected to be 256 bytes, so it occupies up to 0x00FF.
func (m *MMU) LoadBootRom(romFile string) error {
	romContents, err := ioutil.ReadFile(romFile)
	if err != nil {
		return err
	}

	// make sure we don't exceed where we expect m
	if len(romContents) > 256 {
		return errors.New("The bootrom should not exceed 256 bytes in length.")
	}
	m.WriteBytes(romContents, 0)

	return nil
}

// ReadByte reads and returns a byte from memory at the given location.
func (m *MMU) ReadByte(location uint16) byte {
	return m.memory[location]
}

// ReadWord reads and returns a word from memory at the given location.
func (m *MMU) ReadWord(location uint16) uint16 {
	return binary.LittleEndian.Uint16(m.memory[location : location+2])
}

// WriteBytes write bytes into memory at the given location.
func (m *MMU) WriteBytes(content []byte, location uint16) {
	for _, b := range content {
		m.memory[location] = b
		location++
	}
}
