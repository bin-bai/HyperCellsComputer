// main
package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"
	"runtime/pprof"

	"github.com/bin-bai/HyperCellsComputer/computer"
	"github.com/bin-bai/HyperCellsComputer/program"
)

const (
	LogicSize  = 2048
	MemorySize = 32768
)

func main() {
	oninit()
	defer onexit()

	sortComputer := computer.NewComputer(LogicSize, MemorySize)

	sortComputer.Init(program.NewSum)
	sortComputer.Run()

	sortComputer.Init(program.NewBubbleSort)
	sortComputer.Run()
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to file")
var heapdump = flag.String("heapdump", "", "write heap dump to file")

func oninit() {
	flag.Parse()
	if len(*cpuprofile) > 0 {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

func onexit() {
	pprof.StopCPUProfile()

	if len(*memprofile) > 0 {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
	}

	if len(*heapdump) > 0 {
		f, err := os.Create(*heapdump)
		if err != nil {
			log.Fatal(err)
		}
		debug.WriteHeapDump(f.Fd())
	}
}
