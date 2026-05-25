package gojit

import (
	"reflect"
)

// to understand this page read here: https://aaronbalke.com/posts/calling-go-functions-from-jit-code/

// r10, and r11 are used for pointers in conversion to abi0.
// do not use these registers for arguments or returns (use stack instead)

var callPtr = uint64(getTaggedLabelAddr(0x0))

func (a *Assembler) InternalCallFunc(f any) {
	// uses R10, R11 to give max cnt arg / result registers
	// see abi internal for amd64 registers
	// max available is 7 args, 8 returns

	ptr := funcAddr(f)

	a.MovAbs(uint64(callPtr), R11)

	offset := byte(4 + 3 + 10) // mov, movabs, jmp

	// lea r10, [rip+offset]
	a.byte(0x4C)
	a.byte(0x8D)
	a.byte(0x15)
	a.byte(offset)
	a.byte(0)
	a.byte(0)
	a.byte(0)

	a.Mov(R10, Indirect{Rsp, 0, 64})

	a.MovAbs(uint64(ptr), R10)

	// jmp r11
	a.byte(0x41)
	a.byte(0xff)
	a.byte(0xe3)
}

func funcAddr(f any) uintptr {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		panic("funcAddr: not a func")
	}
	return v.Pointer()
}
