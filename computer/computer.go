// computer
package computer

import (
	"fmt"

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

	// Clear flag
	c.Bus.AndFlag(flagindex, types.HCWORD(0))

	ipl := c.program.GetInstPerLoop()
	looptick := types.HCWORD(0)
	for {
		c.Tick()

		// Check flag
		if c.Bus.GetFlag(flagindex) != 0 {
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

	c.Print(dataaddr, dataaddr+datasize-1)

	fmt.Printf("Tickcount is %d\n", c.tick)
	fmt.Printf("Run %d loops\n\n", c.tick/ipl)
}

func (c *Computer) Tick() {
	c.Bus.PreTick()

	for i := range c.LogicCells {
		c.LogicCells[i].Tick()
	}

	c.Bus.Tick()

	c.tick++
}

func (c *Computer) Print(from, to types.HCWORD) {
	c.MemoryBank.Print(from, to)
	fmt.Println()
}
