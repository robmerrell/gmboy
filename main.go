package main

import (
	"flag"
	"fmt"
	"github.com/robmerrell/gmboy/system"
)

func main() {
	debug := flag.String("debug", "", "")
	bootstrap := flag.String("bootstrap", "", "")
	flag.Usage = usage
	flag.Parse()

	romFile := flag.Arg(0)
	if romFile == "" {
		flag.Usage()
		return
	}

	sys := system.NewSystem()

	if *bootstrap != "" {
		err := sys.PerformBootstrap(*bootstrap)
		if err != nil {
			panic(err)
		}
	}

	if *debug != "" {
		err := sys.StartDebugger(*debug)
		if err != nil {
			fmt.Printf("Error loading %s\n", *debug)
			return
		}
	}

	sys.LoadRom(romFile)
	sys.Run()
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  gmbody file.gb\n")
	fmt.Println("  --bootstrap=file.bin   Run the bootstrap process using the specified file. Default is to not bootstrap.")
	fmt.Println("  --debug=file.js        Start the debugger and evaluate the specified file.")
	fmt.Println("  --help                 Show this help text.")
}
