package main

import (
	"github.com/robmerrell/gmboy/system"
)

func main() {
	sys := system.NewSystem()

	err := sys.PerformBootstrap("./roms/bootrom.bin")
	if err != nil {
		panic(err)
	}

	sys.LoadRom("~/tmp/tetris.gb")
	sys.Run()
}
