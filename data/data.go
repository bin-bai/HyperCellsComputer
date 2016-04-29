// data
package data

import (
	"math/rand"
	"time"

	"github.com/bin-bai/HyperCellsComputer/types"
)

func Random(dest []types.HCWORD, bound types.HCWORD) {
	rand.Seed(time.Now().UnixNano())
	for i := range dest {
		dest[i] = types.HCWORD(rand.Int63()) % bound
	}
}

func Ordered(dest []types.HCWORD, from types.HCWORD) {
	for i := range dest {
		dest[i] = from + types.HCWORD(i)
	}
}

func ReverseOrdered(dest []types.HCWORD, from types.HCWORD) {
	for i := range dest {
		dest[i] = from - types.HCWORD(i)
	}
}
