// pcmba
package bus

import (
	"errors"

	"github.com/bin-bai/HyperCellsComputer/mem"
	"github.com/bin-bai/HyperCellsComputer/types"
)

const (
	FLAGS_SIZE = 64
)

// Programmable CPU-Memory Bus Array
type PCMBA struct {
	Cells []Cell

	MemoryBank *mem.MemoryBank

	locked  []bool
	lockers []byte // Only use 2 bits, 0: No locker, 1: One locker, 2: >= 2 lockers, 3: reserved.

	flags []types.HCWORD
}

func (ba *PCMBA) Init(logicSize types.HCWORD, mb *mem.MemoryBank) {
	ba.Cells = make([]Cell, logicSize)

	ba.MemoryBank = mb
	ba.locked = make([]bool, mb.Size())
	ba.lockers = make([]byte, mb.Size())

	ba.flags = make([]types.HCWORD, FLAGS_SIZE)
}

func (ba *PCMBA) Load(addr types.HCWORD, dest []types.HCWORD) error {
	if !ba.MemoryBank.In(addr) || !ba.MemoryBank.In(addr+types.HCWORD(len(dest)-1)) {
		panic(1)
		return NewIllegalAddressError(addr)
	}

	ba.MemoryBank.ReadBatch(addr, dest)
	return nil
}

func (ba *PCMBA) Set(addr types.HCWORD, value types.HCWORD) error {
	if !ba.MemoryBank.In(addr) {
		panic(1)
		return NewIllegalAddressError(addr)
	}
	if ba.locked[addr] {
		panic(2)
		return errors.New("Attemp to write locked memory")
	}

	ba.MemoryBank.Write(addr, value)

	return nil
}

func (ba *PCMBA) SetBatch(addr types.HCWORD, values []types.HCWORD) error {
	if !ba.MemoryBank.In(addr) || !ba.MemoryBank.In(addr+types.HCWORD(len(values)-1)) {
		panic(1)
		return NewIllegalAddressError(addr)
	}

	for i := range values {
		if ba.locked[addr+types.HCWORD(i)] {
			panic(2)
			return errors.New("Attemp to write locked memory")
		}
	}

	ba.MemoryBank.WriteBatch(addr, values)

	return nil
}

func (ba *PCMBA) GetFlag(index types.HCWORD) types.HCWORD {
	return ba.flags[index]
}

func (ba *PCMBA) OrFlag(index types.HCWORD, value types.HCWORD) {
	ba.flags[index] |= value
}

func (ba *PCMBA) AndFlag(index types.HCWORD, value types.HCWORD) {
	ba.flags[index] &= value
}

// The lock result will be setted after tick
func (ba *PCMBA) Trylock(addr types.HCWORD, size types.HCWORD) {
	for i := 0; types.HCWORD(i) < size; i++ {
		if !ba.locked[addr+types.HCWORD(i)] {
			switch ba.lockers[addr+types.HCWORD(i)] {
			case 0:
				ba.lockers[addr+types.HCWORD(i)] = 1
			default:
				ba.lockers[addr+types.HCWORD(i)] = 2
			}
		}
	}
}

func (ba *PCMBA) Unlock(addr types.HCWORD, size types.HCWORD) {
	for i := 0; types.HCWORD(i) < size; i++ {
		ba.locked[addr+types.HCWORD(i)] = false
	}
}

func (ba *PCMBA) ToCell(id, addr types.HCWORD, values []types.HCWORD) error {
	if id < 0 || id >= types.HCWORD(len(ba.Cells)) {
		panic(10)
		return NewIllegalAddressError(id)
	}

	ba.Cells[id].Set(addr, values)
	return nil
}

func (ba *PCMBA) PreTick() {
	for i := range ba.lockers {
		ba.lockers[i] = 0
	}
}

func (ba *PCMBA) Tick() {
	// Handle pending requests
	ba.handleLocks()
}

// Simulate parallel lock
func (ba *PCMBA) handleLocks() {
	var addr, size types.HCWORD

	for i := range ba.Cells {
		addr, size = ba.Cells[i].GetPendingLock()
		if addr <= 0 || size <= 0 {
			continue
		}

		exclusive := true
		for j := 0; types.HCWORD(j) < size; j++ {
			if ba.lockers[addr+types.HCWORD(j)] != 1 {
				exclusive = false
				break
			}
		}

		if exclusive {
			for k := 0; types.HCWORD(k) < size; k++ {
				ba.locked[addr+types.HCWORD(k)] = true
			}
			ba.Cells[i].SetLockResult(1)
		} else {
			ba.Cells[i].SetLockResult(0)
		}
	}
}
