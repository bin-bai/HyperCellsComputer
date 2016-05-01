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
	const (
		R_ID   = 0
		R_ONE  = 1
		R_ZERO = 2
		R_N    = 3
		R_L    = 5
		R_R    = 6
		R_LV   = 7
		R_RV   = 8
	)
	s.instructions = []types.HCWORD{
		logic.ID, R_ID, // R_ID = ID
		logic.SET, R_ONE, 1, // R_ONE = 1 const
		logic.SET, R_ZERO, 0, // R_ZERO = 0 const
		logic.SET, R_N, s.datasize, // R_N = n

		logic.SET, R_L, da, // R_L = dataaddr
		logic.ADD, R_L, R_L, R_ID, // R_L = R_L + R_ID

		logic.REC, logic.R_USER_TOP, // R_USER_TOP = Next PC after REC

		logic.INC, R_N,
		logic.SHIFTR, R_N, 1, // R_N /= 2

		// if id >= n then halt
		logic.CMP, R_ID, R_N, // CMP R_ID, R_N
		7, 6, 6,
		logic.HALT, // Finish in this cell

		logic.COPY, R_R, R_L, // R_R = R_L + R_N
		logic.ADD, R_R, R_R, R_N,

		logic.LOAD, R_L, 1, R_LV, // R_LV = *R_L
		logic.LOAD, R_R, 1, R_RV, // R_RV = *R_R
		logic.ADD, R_LV, R_LV, R_RV, // R_LV = R_LV + R_RV
		logic.SETMR, R_L, R_LV, // *R_L = R_LV
		logic.SETMD, R_R, 0, // *R_R = 0

		logic.ORFLAG, fi, R_ONE, // flag |= R_ONE

		// if n <= R_ONE then halt
		logic.CMP, R_N, R_ONE, // CMP R_N, R_ONE
		6, 6, 7,
		logic.HALT, // Finish in this cell

		logic.JMP, logic.R_USER_TOP, // JMP R_USER_TOP
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

func (s *Sum) Finished() bool {
	return true
}
