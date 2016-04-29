// LogicCell
package logic

import (
	"github.com/bin-bai/HyperCellsComputer/bus"
	"github.com/bin-bai/HyperCellsComputer/types"
)

// TODO Connect cells

const (
	// Instructions begin
	NOP = types.HCWORD(0)
	ID  = types.HCWORD(2)

	HALT    = types.HCWORD(11)
	REC     = types.HCWORD(12)
	JMP     = types.HCWORD(13)
	FORWARD = types.HCWORD(14)
	CMP     = types.HCWORD(16)

	SET  = types.HCWORD(21)
	COPY = types.HCWORD(22)

	LOAD  = types.HCWORD(31)
	SETMR = types.HCWORD(32)
	SETMD = types.HCWORD(33)

	ORFLAG  = types.HCWORD(41)
	ANDFLAG = types.HCWORD(42)

	LOCK   = types.HCWORD(51)
	UNLOCK = types.HCWORD(52)

	ADD = types.HCWORD(61)
	SUB = types.HCWORD(62)
	MUL = types.HCWORD(63)
	DIV = types.HCWORD(64)

	INC = types.HCWORD(65)
	DEC = types.HCWORD(66)

	SHIFTL = types.HCWORD(71)
	SHIFTR = types.HCWORD(72)
	// Instructions end

	// Instruction registers begin
	instReg = 0
	oprand1 = 1
	oprand2 = 2
	oprand3 = 3
	oprand4 = 4
	oprand5 = 5
	oprand6 = 6
	oprand7 = 7

	instRegSize = 8
	// Instruction registers end

	// User registers begin
	userRegSize  = 32
	USER_REG_TOP = userRegSize - 1
	// User registers end
)

var (
	busslot bus.Bus

	oprandCount = [...]types.HCWORD{
		NOP:     types.HCWORD(0),
		ID:      types.HCWORD(1),
		HALT:    types.HCWORD(0),
		REC:     types.HCWORD(1),
		JMP:     types.HCWORD(1),
		FORWARD: types.HCWORD(1),
		CMP:     types.HCWORD(5),
		SET:     types.HCWORD(2),
		COPY:    types.HCWORD(2),
		LOAD:    types.HCWORD(3),
		SETMR:   types.HCWORD(2),
		SETMD:   types.HCWORD(2),
		ORFLAG:  types.HCWORD(2),
		ANDFLAG: types.HCWORD(2),
		LOCK:    types.HCWORD(3),
		UNLOCK:  types.HCWORD(2),
		ADD:     types.HCWORD(3),
		SUB:     types.HCWORD(3),
		MUL:     types.HCWORD(3),
		DIV:     types.HCWORD(3),
		INC:     types.HCWORD(1),
		DEC:     types.HCWORD(1),
		SHIFTL:  types.HCWORD(2),
		SHIFTR:  types.HCWORD(2),
	}
)

func SetBus(b bus.Bus) {
	busslot = b
}

type LogicCell struct {
	id types.HCWORD
	PC types.HCWORD

	instruction []types.HCWORD
	registers   []types.HCWORD

	pendingaddr types.HCWORD
	pendingsize types.HCWORD
	pendingdest types.HCWORD
}

func (lc *LogicCell) Init(id types.HCWORD) {
	lc.id = id
	lc.PC = -1

	lc.instruction = make([]types.HCWORD, instRegSize)
	lc.registers = make([]types.HCWORD, userRegSize)
}

func (lc *LogicCell) Halt() {
	lc.PC = -1
}

func (lc *LogicCell) Tick() {
	if lc.PC < 0 {
		return
	}

	lc.pendingaddr = 0
	lc.pendingsize = 0

	// Load instruction
	err := busslot.Load(lc.PC, lc.instruction[instReg:instReg+1])
	if err != nil {
		lc.Halt()
		return
	}

	// Load oprand(s)
	oc := oprandCount[lc.instruction[instReg]]
	if oc > 0 {
		err := busslot.Load(lc.PC+1, lc.instruction[oprand1:oprand1+oc])
		if err != nil {
			lc.Halt()
			return
		}
	}

	nextpc := lc.PC + oc + 1
	switch lc.instruction[instReg] {
	case NOP:
	// Pass
	case ID:
		lc.registers[lc.instruction[oprand1]] = lc.id
	case HALT:
		lc.Halt()
		return
	case REC:
		lc.registers[lc.instruction[oprand1]] = lc.PC + 2
	case JMP:
		nextpc = lc.registers[lc.instruction[oprand1]]
	case FORWARD:
		nextpc = lc.PC + lc.instruction[oprand1]
	case CMP:
		left := lc.registers[lc.instruction[oprand1]]
		right := lc.registers[lc.instruction[oprand2]]
		switch {
		case left < right:
			nextpc = lc.PC + lc.instruction[oprand3]
		case left == right:
			nextpc = lc.PC + lc.instruction[oprand4]
		default:
			nextpc = lc.PC + lc.instruction[oprand5]
		}
	case SET:
		lc.registers[lc.instruction[oprand1]] = lc.instruction[oprand2]
	case COPY:
		lc.registers[lc.instruction[oprand1]] = lc.registers[lc.instruction[oprand2]]
	case LOAD:
		size := lc.instruction[oprand2]
		dest := lc.instruction[oprand3]
		err := busslot.Load(lc.registers[lc.instruction[oprand1]], lc.registers[dest:dest+size])
		if err != nil {
			lc.Halt()
			return
		}
	case SETMR:
		busslot.Set(lc.registers[lc.instruction[oprand1]], lc.registers[lc.instruction[oprand2]])
	case SETMD:
		busslot.Set(lc.registers[lc.instruction[oprand1]], lc.instruction[oprand2])
	case ORFLAG:
		busslot.OrFlag(lc.instruction[oprand1], lc.registers[lc.instruction[oprand2]])
	case ANDFLAG:
		busslot.AndFlag(lc.instruction[oprand1], lc.registers[lc.instruction[oprand2]])
	case LOCK:
		lc.DeferLock(lc.registers[lc.instruction[oprand1]], lc.instruction[oprand2], lc.instruction[oprand3])
	case UNLOCK:
		busslot.Unlock(lc.registers[lc.instruction[oprand1]], lc.instruction[oprand2])
	case ADD:
		lc.registers[lc.instruction[oprand1]] = lc.registers[lc.instruction[oprand2]] + lc.registers[lc.instruction[oprand3]]
	case SUB:
		lc.registers[lc.instruction[oprand1]] = lc.registers[lc.instruction[oprand2]] - lc.registers[lc.instruction[oprand3]]
	case MUL:
		lc.registers[lc.instruction[oprand1]] = lc.registers[lc.instruction[oprand2]] * lc.registers[lc.instruction[oprand3]]
	case DIV:
		lc.registers[lc.instruction[oprand1]] = lc.registers[lc.instruction[oprand2]] / lc.registers[lc.instruction[oprand3]]
	case INC:
		lc.registers[lc.instruction[oprand1]]++
	case DEC:
		lc.registers[lc.instruction[oprand1]]--
	case SHIFTL:
		lc.registers[lc.instruction[oprand1]] <<= lc.instruction[oprand2].Abs()
	case SHIFTR:
		lc.registers[lc.instruction[oprand1]] >>= lc.instruction[oprand2].Abs()
	default:
		lc.Halt()
		return
	}

	lc.PC = nextpc
}

func (lc *LogicCell) DeferLock(addr, size, dest types.HCWORD) {
	lc.pendingaddr = addr
	lc.pendingsize = size
	lc.pendingdest = dest

	busslot.Trylock(addr, size)
}

func (lc *LogicCell) GetPendingLock() (addr types.HCWORD, size types.HCWORD) {
	return lc.pendingaddr, lc.pendingsize
}

func (lc *LogicCell) SetLockResult(result types.HCWORD) {
	lc.registers[lc.pendingdest] = result
}
