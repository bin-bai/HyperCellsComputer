// main
package main

import (
	"flag"
	"log"
	"os"
	"runtime"
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

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	sortComputer := computer.NewComputer(LogicSize, MemorySize)

	sortComputer.Init(program.NewSum)
	sortComputer.Run()

	sortComputer.Init(program.NewBubbleSort)
	sortComputer.Run()
}

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to file")
	heapdump   = flag.String("heapdump", "", "write heap dump to file")
)

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
	if len(*cpuprofile) > 0 {
		pprof.StopCPUProfile()
	}

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
