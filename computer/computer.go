// computer
package computer

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/bin-bai/HyperCellsComputer/bus"
	"github.com/bin-bai/HyperCellsComputer/logic"
	"github.com/bin-bai/HyperCellsComputer/mem"
	"github.com/bin-bai/HyperCellsComputer/program"
	"github.com/bin-bai/HyperCellsComputer/types"
)

const (
	instaddr = types.HCWORD(1)
	dataaddr = types.HCWORD(1024)

	flagindex = types.HCWORD(1)
)

type Computer struct {
	Bus bus.PCMBA

	LogicCells []logic.LogicCell

	MemoryBank *mem.MemoryBank

	program program.Program

	tick   types.HCWORD
	status types.HCWORD

	// Try to accelerate execution by utilizing Go routines.
	// Well, can't see improvement on Windows 10.
	executers []*Executer
	exeWG     sync.WaitGroup
}

func NewComputer(logicSize, memorySize types.HCWORD) *Computer {
	c := new(Computer)
	c.status = 0

	c.MemoryBank = mem.NewMemoryBank(0, memorySize)

	c.Bus.Init(logicSize, c.MemoryBank)
	logic.SetBus(&c.Bus)

	c.LogicCells = make([]logic.LogicCell, logicSize)
	for i := range c.LogicCells {
		c.LogicCells[i].Init(types.HCWORD(i))
		c.LogicCells[i].Stack = make([]types.HCWORD, 0, 20)

		c.Bus.Cells[i] = &c.LogicCells[i]
	}

	c.executers = make([]*Executer, runtime.GOMAXPROCS(0))
	for i := range c.executers {
		c.executers[i] = NewExecuter(&c.exeWG)
	}

	return c
}

func (c *Computer) Init(newer program.Newer) {
	c.program = newer(0)
	c.program.Init(dataaddr, flagindex, types.HCWORD(len(c.LogicCells)))

	for i := range c.LogicCells {
		c.LogicCells[i].PC = instaddr
	}

	instructions := c.program.GetInst()
	c.MemoryBank.WriteBatch(instaddr, instructions)

	datasize := c.program.GetDataSize()
	c.program.GetData(c.MemoryBank.DMA(dataaddr, datasize))
}

func (c *Computer) Run() {
	c.tick = 0

	datasize := c.program.GetDataSize()
	c.Print(dataaddr, dataaddr+datasize-1)

	for i := range c.executers {
		c.executers[i].Start()
	}

	// Clear flag
	c.Bus.AndFlag(flagindex, types.HCWORD(0))

	ipl := c.program.GetInstPerLoop()
	looptick := types.HCWORD(0)
	for {
		c.Tick()

		// Check flag
		if c.Bus.GetFlag(flagindex) != 0 {
			// c.Print(dataaddr, dataaddr+datasize-1)
			c.Bus.AndFlag(flagindex, types.HCWORD(0))
			looptick = 0
		} else {
			looptick++
			if looptick > ipl*2 {
				// Run at least a whole loop but no more work to do
				break
			}
		}
	}

	for i := range c.executers {
		c.executers[i].Stop()
	}

	c.Print(dataaddr, dataaddr+datasize-1)

	fmt.Println()
	fmt.Printf("Tickcount is %d\n", c.tick)
	fmt.Printf("Run %d loops\n\n", c.tick/ipl)
}

func (c *Computer) Tick() {
	c.Bus.PreTick()

	/*
		for i := range c.LogicCells {
			c.LogicCells[i].Tick()
		}
	*/

	// Utilize Go routines begin
	nlc := len(c.LogicCells)
	nexe := len(c.executers)
	cpree := nlc / nexe

	c.exeWG.Add(nexe)

	for i := 0; i < nexe; i++ {
		c.executers[i].Exec(c.LogicCells[i*cpree : (i+1)*cpree])
	}

	for i := nexe * cpree; i < nlc; i++ {
		c.LogicCells[i].Tick()
	}

	c.exeWG.Wait()
	// Utilize Go routines end

	c.Bus.Tick()

	c.tick++
}

func (c *Computer) Print(from, to types.HCWORD) {
	c.MemoryBank.Print(from, to)
}

type Executer struct {
	exeWG *sync.WaitGroup

	worker chan []logic.LogicCell
	stop   chan bool
}

func NewExecuter(exeWG *sync.WaitGroup) *Executer {
	e := new(Executer)
	e.exeWG = exeWG
	e.worker = make(chan []logic.LogicCell)
	e.stop = make(chan bool)
	return e
}

func (e *Executer) Start() {
	go func() {
		for {
			select {
			case cells := <-e.worker:
				for i := range cells {
					cells[i].Tick()
				}
				e.exeWG.Done()

			case <-e.stop:
				return
			}
		}
	}()
}

func (e *Executer) Stop() {
	e.stop <- true
}

func (e *Executer) Exec(cells []logic.LogicCell) {
	e.worker <- cells
}
