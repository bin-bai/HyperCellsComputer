// blockdev
package io

import (
	"github.com/bin-bai/HyperCellsComputer/types"
)

type BlockDev interface {
	Read(addr types.HCWORD, dest []byte)
	Write(addr types.HCWORD, value []byte)
}
