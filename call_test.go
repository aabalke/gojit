package gojit

import (
	"runtime"
	"testing"
	"unsafe"
)

// note: variables called within go funcs have to be global

var called = false

func TestCall(t *testing.T) {

    pagesize := 512


	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.CallFunc(func() {
        called = true
	})

    exit(asm)

	callJIT(&asm.Buf[0])

	asm.Release()

    if !called {
        t.Errorf("Failed Test Call: called variable not set\n")
    }
}

var i = 1 << 16
func recursive() {
    if i > 0 {
        i--
        recursive()
    }
}

func TestCallRecursion(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.CallFunc(recursive)

    exit(asm)

	callJIT(&asm.Buf[0])

	asm.Release()
}

func TestCallGc(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.CallFunc(func ()  {
        runtime.GC()
	})

    exit(asm)

	callJIT(&asm.Buf[0])

	asm.Release()
}

var v uint64

func TestIndirect(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.CallFunc(func ()  {
        v = 0xBEEF
	})

	asm.Mov(Imm(0xDEAD), Rax)
	asm.MovAbs(uint64(uintptr(unsafe.Pointer(&v))), Rbx)
	asm.Mov(Rax, Indirect{Rbx, 0, 64})

	asm.CallFunc(func ()  {
        v = 0xABBA
	})

    exit(asm)

	callJIT(&asm.Buf[0])

	asm.Release()

    if v != 0xABBA {
        t.Errorf("Failed Test Call: v variable not set to 0xABBA\n")
    }
}


func TestCallArguments(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

    var o uint64

    asm.Mov(Imm(1), Rax)
    asm.Mov(Imm(1), Rbx)
    asm.Mov(Imm(1), Rcx)
    asm.Mov(Imm(1), Rdi)
    asm.Mov(Imm(1), Rsi)
    asm.Mov(Imm(1), R8)
    asm.Mov(Imm(1), R9)

	asm.CallFunc(func(a, b, c, d, e, f, g uint64) uint64  {
        return a+b+c+d+e+f+g
	})

    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&o))), Rbx)
    asm.Mov(Rax, Indirect{Rbx, 0, 64})
    exit(asm)

	callJIT(&asm.Buf[0])

	asm.Release()

    if o != 7 {
        t.Errorf("Failed Test Call Arguments: o variable not set to 0x7. Value: %X\n", o)
    }
}


func TestCallResults(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.CallFunc(func() (
        uint64, uint64, uint64, uint64,
        uint64, uint64, uint64, uint64)  {
        return 1, 2, 3, 4, 5, 6, 7, 8
	})

    var ta, tb, tc, td, te, tf, tg, th uint64
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&ta))), R11)
    asm.Mov(Rax, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&tb))), R11)
    asm.Mov(Rbx, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&tc))), R11)
    asm.Mov(Rcx, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&td))), R11)
    asm.Mov(Rdi, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&te))), R11)
    asm.Mov(Rsi, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&tf))), R11)
    asm.Mov(R8, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&tg))), R11)
    asm.Mov(R9, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&th))), R11)
    asm.Mov(R10, Indirect{R11, 0, 64})


    exit(asm)

	callJIT(&asm.Buf[0])

	asm.Release()

    if (
        ta != 1 ||
        tb != 2 ||
        tc != 3 ||
        td != 4 ||
        te != 5 ||
        tf != 6 ||
        tg != 7 ||
        th != 8) {
        t.Error("BAD")
    }
}
