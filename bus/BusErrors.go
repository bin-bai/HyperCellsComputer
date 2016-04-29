// BusErrors
package bus

import (
	"fmt"

	"github.com/bin-bai/HyperCellsComputer/types"
)

type IllegalAddressError struct {
	ErrorString string
	Address     types.HCWORD
}

func (ia *IllegalAddressError) Error() string {
	return fmt.Sprintf("%s: %d", ia.ErrorString, ia.Address)
}

func NewIllegalAddressError(addr types.HCWORD) *IllegalAddressError {
	ia := &IllegalAddressError{ErrorString: "Illegal Address", Address: addr}
	return ia
}
