// bubblesort
package program

import (
	"github.com/bin-bai/HyperCellsComputer/data"
	"github.com/bin-bai/HyperCellsComputer/logic"
	"github.com/bin-bai/HyperCellsComputer/types"
)

type BubbleSort struct {
	instructions []types.HCWORD
	datasize     types.HCWORD
}

func NewBubbleSort(args ...types.HCWORD) Program {
	s := new(BubbleSort)
	return s
}

func (s *BubbleSort) Init(da types.HCWORD, fi types.HCWORD, args ...types.HCWORD) {
	s.datasize = args[0]*2 + 1

	/*
		[0] [1] [2] [3] [4] [5] [6]
		 |   |   |   |   |   |   |
		 -----   -----   -----
		  C0      C1      C2
		[0] [1] [2] [3] [4] [5] [6]
		     |   |   |   |   |   |
		     -----   -----   -----
		      C0      C1      C2
	*/
	const (
		R_2ID = 0
		R_ONE = 1
		R_L   = 3
		R_R   = 4
		R_LV  = 5
		R_RV  = 6
	)
	s.instructions = []types.HCWORD{
		logic.ID, R_2ID, // R_2ID = ID
		logic.SHIFTL, R_2ID, 1, // R_2ID *= 2
		logic.SET, R_ONE, 1, // R_ONE = 1 const

		logic.SET, R_L, da, // R_L = dataaddr
		logic.ADD, R_L, R_L, R_2ID, // R_L = R_L + R_2ID
		logic.ADD, R_R, R_L, R_ONE, // R_R = R_L + R_ONE

		logic.REC, logic.R_USER_TOP, // R_USER_TOP = Next PC after REC

		logic.MARK, 4,
		logic.FORWARD, 16,
		logic.INC, R_L,
		logic.INC, R_R,

		logic.MARK, 4,
		logic.FORWARD, 8,
		logic.DEC, R_L,
		logic.DEC, R_R,

		logic.JMP, logic.R_USER_TOP, // JMP R_USER_TOP

		// Function which does comparation and swap
		logic.FUNC,
		logic.LOAD, R_L, 1, R_LV, // R_LV = *R_L
		logic.LOAD, R_R, 1, R_RV, // R_RV = *R_R
		// if R_LV <= R_RV{
		logic.CMP, R_LV, R_RV, // CMP R_LV, R_RV
		6, 6, 10,
		logic.NOP, // To synchronize without lock
		logic.NOP,
		logic.NOP,
		logic.RETURN,
		// } else {
		logic.SETMR, R_L, R_RV, // *R_L = R_RV
		logic.SETMR, R_R, R_LV, // *R_R = R_LV
		logic.ORFLAG, fi, R_ONE, // flag |= R_ONE
		logic.RETURN,
		// }
	}
}

func (s *BubbleSort) GetInstPerLoop() types.HCWORD {
	return types.HCWORD(25)
}

func (s *BubbleSort) GetInst() []types.HCWORD {
	return s.instructions
}

func (s *BubbleSort) GetDataSize() types.HCWORD {
	return s.datasize
}

func (s *BubbleSort) GetData(dest []types.HCWORD) {
	data.Random(dest[0:s.datasize], 10000)
}

func (s *BubbleSort) Finished() bool {
	return true
}
