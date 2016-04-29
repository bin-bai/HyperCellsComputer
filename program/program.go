// program
package program

import (
	"github.com/bin-bai/HyperCellsComputer/types"
)

type Program interface {
	// Do memory map
	Init(da types.HCWORD, fi types.HCWORD, args ...types.HCWORD)

	GetInstPerLoop() types.HCWORD
	GetInst() []types.HCWORD

	GetDataSize() types.HCWORD
	GetData(dest []types.HCWORD)
}

type Newer func(args ...types.HCWORD) Program
