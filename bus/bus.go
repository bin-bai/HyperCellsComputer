// bus
package bus

import (
	"github.com/bin-bai/HyperCellsComputer/types"
)

type Bus interface {
	Set(addr types.HCWORD, value types.HCWORD) error
	Load(addr types.HCWORD, dest []types.HCWORD) error

	GetFlag(index types.HCWORD) types.HCWORD
	OrFlag(index types.HCWORD, value types.HCWORD)
	AndFlag(index types.HCWORD, value types.HCWORD)

	// It's better to do lock in OS, without the supporting from HW
	// The lock result will be setted after tick
	Trylock(addr types.HCWORD, size types.HCWORD)
	Unlock(addr types.HCWORD, size types.HCWORD)
}

type Pender interface {
	GetPendingLock() (addr types.HCWORD, size types.HCWORD)
	SetLockResult(result types.HCWORD)
}
