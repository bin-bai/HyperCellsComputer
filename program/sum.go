// sum
package program

import (
	"github.com/bin-bai/HyperCellsComputer/data"
	"github.com/bin-bai/HyperCellsComputer/logic"
	"github.com/bin-bai/HyperCellsComputer/types"
)

type Sum struct {
	instructions []types.HCWORD
	datasize     types.HCWORD
}

func NewSum(args ...types.HCWORD) Program {
	s := new(Sum)
	return s
}

func (s *Sum) Init(da types.HCWORD, fi types.HCWORD, args ...types.HCWORD) {
	s.datasize = args[0] * 2

	/*
		[0] [n/2] [1] [n/2+1] [2] [n/2+2]
		 |    |    |     |     |     |
		 ------    -------     -------
		   C0         C1          C2
		[0] [n/4] [1] [n/4+1] [2] [n/4+2]
		 |    |    |     |     |     |
		 ------    -------     -------
		   C0         C1          C2
		...
	*/
	s.instructions = []types.HCWORD{
		logic.ID, 0, // R0 = ID
		logic.SET, 1, 1, // R1 = 1 const
		logic.SET, 2, 0, // R2 = 0 const
		logic.SET, 3, s.datasize, // R3 = n

		logic.SET, 5, da, // R5 = dataaddr
		logic.ADD, 5, 5, 0, // R5 = R5 + R0

		logic.REC, logic.USER_REG_TOP, // R(USER_REG_TOP) = Next PC after REC

		logic.INC, 3,
		logic.SHIFTR, 3, 1, // R3 = R3 / 2

		// if id >= n then halt
		logic.CMP, 0, 3, // CMP R0, R3
		7, 6, 6,
		logic.HALT, // Finish in this cell

		logic.COPY, 6, 5, // R6 = R5
		logic.ADD, 6, 6, 3, // R6 = R6 + R3

		logic.LOAD, 5, 1, 7, // R7 = *R5
		logic.LOAD, 6, 1, 8, // R8 = *R6
		logic.ADD, 7, 7, 8, // R7 = R7 + R8
		logic.SETMR, 5, 7, // *R5 = R7
		logic.SETMD, 6, 0, // *R6 = 0

		logic.ORFLAG, fi, 1, // flag |= 1

		// if n <= 1 then halt
		logic.CMP, 3, 1, // CMP R3, R1
		6, 6, 7,
		logic.HALT, // Finish in this cell

		logic.JMP, logic.USER_REG_TOP, // JMP R(USER_REG_TOP)
	}
}

func (s *Sum) GetInstPerLoop() types.HCWORD {
	return types.HCWORD(13)
}

func (s *Sum) GetInst() []types.HCWORD {
	return s.instructions
}

func (s *Sum) GetDataSize() types.HCWORD {
	return s.datasize
}

func (s *Sum) GetData(dest []types.HCWORD) {
	data.Ordered(dest[0:s.datasize], 1)
}
