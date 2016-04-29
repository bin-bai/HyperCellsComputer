// chardev
package io

import (
	"github.com/binbai/HyperCellsComputer/types"
)

type CharDev interface {
	Read(addr types.HCWORD, dest []byte)
	Write(addr types.HCWORD, value []byte)
}
