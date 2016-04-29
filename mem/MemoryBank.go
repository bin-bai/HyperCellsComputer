// MemoryBank
package mem

import (
	"fmt"

	"github.com/bin-bai/HyperCellsComputer/types"
)

// TODO Check boundary

type MemoryBank struct {
	start types.HCWORD
	end   types.HCWORD

	memory []types.HCWORD
}

func NewMemoryBank(start types.HCWORD, size types.HCWORD) *MemoryBank {
	mb := new(MemoryBank)
	mb.start = start
	mb.end = start + size - 1
	mb.memory = make([]types.HCWORD, size)
	return mb
}

func (mb *MemoryBank) Size() types.HCWORD {
	return mb.end - mb.start + 1
}

func (mb *MemoryBank) In(addr types.HCWORD) bool {
	return mb.start <= addr && addr <= mb.end
}

func (mb *MemoryBank) Read(addr types.HCWORD) types.HCWORD {
	return mb.memory[addr-mb.start]
}

func (mb *MemoryBank) Write(addr types.HCWORD, value types.HCWORD) {
	mb.memory[addr-mb.start] = value
}

func (mb *MemoryBank) ReadBatch(addr types.HCWORD, dest []types.HCWORD) {
	copy(dest, mb.memory[addr-mb.start:])
}

func (mb *MemoryBank) WriteBatch(addr types.HCWORD, values []types.HCWORD) {
	copy(mb.memory[addr-mb.start:], values)
}

func (mb *MemoryBank) DMA(addr types.HCWORD, size types.HCWORD) []types.HCWORD {
	return mb.memory[addr-mb.start : addr-mb.start+size]
}

func (mb *MemoryBank) Print(from, to types.HCWORD) {
	fmt.Println(mb.memory[from : to+1])
}
