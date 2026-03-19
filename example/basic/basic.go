package main

import (
	"fmt"
	"unsafe"

	"github.com/aabalke/gojit"
)

// this example sets a = 0xBEEF and b = 0xDEAD from the JIT compiler

const PAGE_SIZE = 0x1000
var a, b uint32

func main() {

    // initialize a block of memory for jit instructions
	asm, err := gojit.New(PAGE_SIZE)
    if err != nil {
        panic(err)
    }

    // build jit instructions to be called

    // sets a = 0xBEEF
    // MOV Rax, 0xBEEF
    // MOV Rbx, *a
    // MOV [Rbx], Rax

    asm.Mov(gojit.Imm(0xBEEF), gojit.Rax)
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&a))), gojit.Rbx)
    asm.Mov(gojit.Eax, gojit.Indirect{Base: gojit.Rbx, Offset: 0, Bits: 32})

    // sets b = 0xDEAD
    // MOV Rax, 0xDEAD
    // MOV Rbx, *b
    // MOV [Rbx], Rax

    asm.Mov(gojit.Imm(0xDEAD), gojit.Rax)
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&b))), gojit.Rbx)
    asm.Mov(gojit.Eax, gojit.Indirect{Base: gojit.Rbx, Offset: 0, Bits: 32})

    gojit.ExitAssembler(asm)

    if err := asm.Error(); err != nil {
        panic(err)
    }

    fmt.Printf("BEFORE A %08X, B %08X\n", a, b)

    // call jit instructions
    gojit.CallJit(uintptr(unsafe.Pointer(&asm.Buf[0])))

    fmt.Printf("AFTER  A %08X, B %08X\n", a, b)

    // cleanup
    asm.Release()
}
