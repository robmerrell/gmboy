package cpu

import (
	"github.com/robmerrell/gmboy/testhelpers"
	"testing"
)

func TestRegisterWord(t *testing.T) {
	r := &register{0x32, 0x11}
	testhelpers.AssertWord(t, 0x3211, r.word())
}

func TestRegisterSetWord(t *testing.T) {
	r := &register{}
	r.setWord(0x3211)
	testhelpers.AssertByte(t, r.low, 0x32)
	testhelpers.AssertByte(t, r.high, 0x11)
}
