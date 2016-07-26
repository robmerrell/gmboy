package mmu

import (
	"github.com/robmerrell/gmboy/testhelpers"
	"testing"
)

func TestReadByte(t *testing.T) {
	m := NewMMU()
	m.WriteBytes([]byte{0x31}, 0)

	testhelpers.AssertByte(t, 0x31, m.ReadByte(0))
}

func TestReadWord(t *testing.T) {
	m := NewMMU()
	m.WriteBytes([]byte{0x31, 0x32}, 0)

	testhelpers.AssertWord(t, 0x3231, m.ReadWord(0))
}

func TestLoadBootRom(t *testing.T) {
	m := NewMMU()
	m.LoadBootRom("./testdata/testboot.bin")

	testhelpers.AssertWord(t, 0xFFFF, m.ReadWord(0))
	testhelpers.AssertByte(t, 0x00, m.ReadByte(0x100))
}
