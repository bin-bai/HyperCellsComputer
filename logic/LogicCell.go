// LogicCell
package logic

import (
	"github.com/bin-bai/HyperCellsComputer/bus"
	"github.com/bin-bai/HyperCellsComputer/types"
)

// TODO Connect cells

const (
	// Instructions begin
	NOP  = types.HCWORD(0)
	HALT = types.HCWORD(1)
	ID   = types.HCWORD(2)

	REC     = types.HCWORD(10)
	JMP     = types.HCWORD(11)
	FORWARD = types.HCWORD(12)

	CMP    = types.HCWORD(20)
	IFZERO = types.HCWORD(21)

	PUSH   = types.HCWORD(31)
	POP    = types.HCWORD(32)
	MARK   = types.HCWORD(33)
	FUNC   = types.HCWORD(34)
	RETURN = types.HCWORD(35)

	SET  = types.HCWORD(50)
	COPY = types.HCWORD(51)

	LOAD  = types.HCWORD(60)
	SETMR = types.HCWORD(62)
	SETMD = types.HCWORD(63)

	ADD = types.HCWORD(100)
	SUB = types.HCWORD(101)
	MUL = types.HCWORD(102)
	DIV = types.HCWORD(103)

	INC = types.HCWORD(110)
	DEC = types.HCWORD(111)

	SHIFTL = types.HCWORD(120)
	SHIFTR = types.HCWORD(121)

	GETFLAG = types.HCWORD(500)
	ORFLAG  = types.HCWORD(501)
	ANDFLAG = types.HCWORD(502)

	LOCK   = types.HCWORD(510)
	UNLOCK = types.HCWORD(511)

	TOCELL = types.HCWORD(600)
	IFIN   = types.HCWORD(601)
	GETIN  = types.HCWORD(602)
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
	UserRegSize = 32

	R_RET      = UserRegSize - 1
	R_USER_TOP = R_RET - 1
	// User registers end
)

var (
	busslot bus.Bus

	oprandCount = map[types.HCWORD]types.HCWORD{
		NOP:  types.HCWORD(0),
		HALT: types.HCWORD(0),
		ID:   types.HCWORD(1),

		REC:     types.HCWORD(1),
		JMP:     types.HCWORD(1),
		FORWARD: types.HCWORD(1),

		CMP:    types.HCWORD(5),
		IFZERO: types.HCWORD(3),

		PUSH:   types.HCWORD(1),
		POP:    types.HCWORD(1),
		MARK:   types.HCWORD(1),
		FUNC:   types.HCWORD(0),
		RETURN: types.HCWORD(0),

		SET:  types.HCWORD(2),
		COPY: types.HCWORD(2),

		LOAD:  types.HCWORD(3),
		SETMR: types.HCWORD(2),
		SETMD: types.HCWORD(2),

		ADD:    types.HCWORD(3),
		SUB:    types.HCWORD(3),
		MUL:    types.HCWORD(3),
		DIV:    types.HCWORD(3),
		INC:    types.HCWORD(1),
		DEC:    types.HCWORD(1),
		SHIFTL: types.HCWORD(2),
		SHIFTR: types.HCWORD(2),

		GETFLAG: types.HCWORD(2),
		ORFLAG:  types.HCWORD(2),
		ANDFLAG: types.HCWORD(2),
		LOCK:    types.HCWORD(3),
		UNLOCK:  types.HCWORD(2),

		TOCELL: types.HCWORD(4),
		IFIN:   types.HCWORD(2),
		GETIN:  types.HCWORD(2),
	}
)

func SetBus(b bus.Bus) {
	busslot = b
}

type LogicCell struct {
	id types.HCWORD

	PC types.HCWORD

	Stack []types.HCWORD

	instruction []types.HCWORD
	registers   []types.HCWORD

	pendingaddr types.HCWORD
	pendingsize types.HCWORD
	pendingdest types.HCWORD

	incomeaddr types.HCWORD
	incomesize types.HCWORD
}

func (lc *LogicCell) Init(id types.HCWORD) {
	lc.id = id
	lc.PC = -1

	lc.instruction = make([]types.HCWORD, instRegSize)
	lc.registers = make([]types.HCWORD, UserRegSize)

	lc.pendingaddr = -1
	lc.pendingsize = -1

	lc.incomeaddr = -1
	lc.incomesize = -1
}

func (lc *LogicCell) Halt() {
	lc.PC = -1
}

func (lc *LogicCell) Fatal() {
	if lc.PC > 0 {
		lc.PC = -lc.PC
	}
}

func (lc *LogicCell) Tick() {
	if lc.PC < 0 {
		return
	}

	// Load instruction
	err := busslot.Load(lc.PC, lc.instruction[instReg:instReg+1])
	if err != nil {
		lc.Fatal()
		return
	}

	// Load oprand(s)
	oc, ok := oprandCount[lc.instruction[instReg]]
	if !ok {
		lc.Fatal()
		return
	}

	if oc > 0 {
		err := busslot.Load(lc.PC+1, lc.instruction[oprand1:oprand1+oc])
		if err != nil {
			lc.Fatal()
			return
		}
	}

	nextpc := lc.PC + oc + 1

	switch lc.instruction[instReg] {
	case NOP:
		// Pass
	case HALT:
		lc.Halt()
		return
	case ID:
		lc.registers[lc.instruction[oprand1]] = lc.id

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
	case IFZERO:
		if lc.registers[lc.instruction[oprand1]] == 0 {
			nextpc = lc.PC + lc.instruction[oprand2]
		} else {
			nextpc = lc.PC + lc.instruction[oprand3]
		}

	case PUSH:
		lc.Stack = append(lc.Stack, lc.registers[lc.instruction[oprand1]])
	case POP:
		top := len(lc.Stack) - 1
		if top < 0 {
			lc.Fatal()
			return
		}
		lc.registers[lc.instruction[oprand1]] = lc.Stack[top]
		lc.Stack = lc.Stack[0:top]
	case MARK:
		lc.Stack = append(lc.Stack, lc.PC+lc.instruction[oprand1])
	case FUNC:
		// Pass
	case RETURN:
		top := len(lc.Stack) - 1
		if top < 0 {
			lc.Fatal()
			return
		}
		nextpc = lc.Stack[top]
		lc.Stack = lc.Stack[0:top]

	case SET:
		lc.registers[lc.instruction[oprand1]] = lc.instruction[oprand2]
	case COPY:
		lc.registers[lc.instruction[oprand1]] = lc.registers[lc.instruction[oprand2]]

	case LOAD:
		size := lc.instruction[oprand2]
		dest := lc.instruction[oprand3]
		err := busslot.Load(lc.registers[lc.instruction[oprand1]], lc.registers[dest:dest+size])
		if err != nil {
			lc.Fatal()
			return
		}
	case SETMR:
		busslot.Set(lc.registers[lc.instruction[oprand1]], lc.registers[lc.instruction[oprand2]])
	case SETMD:
		busslot.Set(lc.registers[lc.instruction[oprand1]], lc.instruction[oprand2])

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

	case GETFLAG:
		lc.registers[lc.instruction[oprand2]] = busslot.GetFlag(lc.instruction[oprand1])
	case ORFLAG:
		busslot.OrFlag(lc.instruction[oprand1], lc.registers[lc.instruction[oprand2]])
	case ANDFLAG:
		busslot.AndFlag(lc.instruction[oprand1], lc.registers[lc.instruction[oprand2]])
	case LOCK:
		lc.DeferLock(lc.registers[lc.instruction[oprand1]], lc.instruction[oprand2], lc.instruction[oprand3])
	case UNLOCK:
		busslot.Unlock(lc.registers[lc.instruction[oprand1]], lc.instruction[oprand2])

	case TOCELL:
		src := lc.instruction[oprand3]
		size := lc.instruction[oprand4]
		values := lc.registers[src : src+size]
		busslot.ToCell(lc.registers[lc.instruction[oprand1]], lc.instruction[oprand2], values)
	case IFIN:
		if lc.incomeaddr >= 0 && lc.incomesize >= 0 {
			nextpc = lc.PC + lc.instruction[oprand1]
		} else {
			nextpc = lc.PC + lc.instruction[oprand2]
		}
	case GETIN:
		lc.registers[lc.instruction[oprand1]] = lc.incomeaddr
		lc.registers[lc.instruction[oprand2]] = lc.incomesize
		lc.incomeaddr = -1
		lc.incomesize = -1
	default:
		lc.Fatal()
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

	lc.pendingaddr = -1
	lc.pendingsize = -1
}

func (lc *LogicCell) Set(addr types.HCWORD, values []types.HCWORD) {
	copy(lc.registers[addr:], values)

	lc.incomeaddr = addr
	lc.incomesize = types.HCWORD(len(values))
}
