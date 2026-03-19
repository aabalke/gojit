package main

import (
	"fmt"
	"unsafe"

	"github.com/aabalke/gojit"
)

// this example has PrintFunction called from jit code.
// for details on registers used for function/method arguments, and results,
// see https://go.dev/src/cmd/compile/abi-internal

const PAGE_SIZE = 0x1000

func main() {

    // initialize a block of memory for jit instructions
	asm, err := gojit.New(PAGE_SIZE)
    if err != nil {
        panic(err)
    }

    // build jit instructions to be called

    // calls PrintFunction with age as RAX (arg 1)
    // MOV Rax, 26
    // CALL PrintFunction

    asm.Mov(gojit.Imm(26), gojit.Rax)
    asm.InternalCallFunc(PrintFunction)
    gojit.ExitAssembler(asm)

    if err := asm.Error(); err != nil {
        panic(err)
    }

    // call jit instructions
    gojit.CallJit(uintptr(unsafe.Pointer(&asm.Buf[0])))

    // cleanup
    asm.Release()
}

//go:nosplit
func PrintFunction(years uint32) {
    fmt.Printf("I am %d years old.\n", years)
}
