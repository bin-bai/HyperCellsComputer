// sorting
package program

import (
	"github.com/bin-bai/HyperCellsComputer/data"
	"github.com/bin-bai/HyperCellsComputer/logic"
	"github.com/bin-bai/HyperCellsComputer/types"
)

type Sorting struct {
	instructions []types.HCWORD
	datasize     types.HCWORD
}

func NewSorting(args ...types.HCWORD) Program {
	s := new(Sorting)
	return s
}

func (s *Sorting) Init(da types.HCWORD, fi types.HCWORD, args ...types.HCWORD) {
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
	s.instructions = []types.HCWORD{
		logic.ID, 0, // R0 = ID
		logic.SET, 1, 1, // R1 = 1 const
		logic.SET, 2, 0, // R2 = 0
		logic.REC, logic.USER_REG_TOP, // R(USER_REG_TOP) = Next PC after REC

		// R3 = dataaddr + R2 + R0 * 2
		logic.SET, 3, da, // R3 = dataaddr
		logic.ADD, 3, 3, 2, // R3 = R3 + R2
		logic.ADD, 3, 3, 0, // R3 = R3 + R0
		logic.ADD, 3, 3, 0, // R3 = R3 + R0
		logic.COPY, 4, 3, // R4 = R3
		logic.INC, 4, // R4++
		logic.LOAD, 3, 2, 5, // R5 = *R3, R6 = *(R3 + 1)

		logic.CMP, 5, 6, // CMP R5, R6
		6, 6, 10,
		logic.NOP, // To synchronize without lock
		logic.NOP, // To synchronize without lock
		logic.FORWARD, 11,
		logic.SETMR, 3, 6, // *R3 = R6
		logic.SETMR, 4, 5, // *R4 = R5
		logic.ORFLAG, fi, 1, // flag |= 1

		logic.SUB, 2, 1, 2, // R2 = R1 - R2
		logic.JMP, logic.USER_REG_TOP, // JMP R(USER_REG_TOP)
	}
}

func (s *Sorting) GetInstPerLoop() types.HCWORD {
	return types.HCWORD(26)
}

func (s *Sorting) GetInst() []types.HCWORD {
	return s.instructions
}

func (s *Sorting) GetDataSize() types.HCWORD {
	return s.datasize
}

func (s *Sorting) GetData(dest []types.HCWORD) {
	data.Random(dest[0:s.datasize], 10000)
}
