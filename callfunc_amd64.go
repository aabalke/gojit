package gojit

import (
	"reflect"
)

// to understand this page read here: https://aaronbalke.com/posts/calling-go-functions-from-jit-code/

// arg and result available rax, rbx, rcx, rdi
// saved over call r8, r9, r10, r11, rsi
// clobbers r12, r13

var (
	callPtr = getTaggedLabelAddr("func")
	exitPtr = getTaggedLabelAddr("exit")
)

func (a *Assembler) CallFunc(f any) {
	const offset = byte(4 + 3 + 10 + 10) // mov, movabs, movabs, jmp

	// lea r13, [rip+offset]
	a.byte(0x4D)
	a.byte(0x8D)
	a.byte(0x2D)
	a.byte(offset)
	a.byte(0)
	a.byte(0)
	a.byte(0)

	a.Mov(R13, Indirect{Rsp, 0, 64})

	a.MovAbs(uint64(getFuncPtr(f)), R12)
	a.MovAbs(uint64(callPtr), R13)

	// jmp r13
	a.byte(0x41)
	a.byte(0xff)
	a.byte(0xe5)
}

func getFuncPtr(f any) uintptr {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		panic("funcAddr: not a func")
	}
	return v.Pointer()
}
