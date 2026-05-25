package gojit

import (
	"reflect"
)

var callPtr = uint64(getTaggedLabelAddr(0x0))

func (a *Assembler) CallFunc(f any) {
	funcPtr := uint64(funcAddr(f))
	offset := 4 * 2

	a.Mov64(R23, funcPtr)
	a.Mov64(R24, callPtr)
	a.addInst(ADR(R25, int32(offset)))
	a.addInst(getBR(R24))
}

func ADR(rd Reg, imm int32) uint32 {
	// imm is signed byte offset, must fit in [-1MB, +1MB)
	u := uint32(imm)

	immlo := (u & 0x3) << 29
	immhi := ((u >> 2) & 0x7FFFF) << 5

	return 0x10000000 | immlo | immhi | (uint32(rd) & 0x1F)
}

func getBLR(rn Reg) uint32 {
	return (uint32(0b1101_0110_0011_1111) << 16) | uint32(rn<<5)
}

func getBR(rn Reg) uint32 {
	return (uint32(0b1101_0110_0001_1111) << 16) | uint32(rn<<5)
}

func funcAddr(f any) uintptr {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		panic("funcAddr: not a func")
	}
	return v.Pointer()
}
