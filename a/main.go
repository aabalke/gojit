package main

import (
	"fmt"
	"unsafe"

	"github.com/aabalke/gojit"
	"github.com/aabalke/gojit/amd64"
)

// integer args and results (params and returns) use RAX, RBX, RCX, RDI, RSI, R8, R9, R10, R11, then stack
// floats are X0 – X14
// https://go.dev/src/cmd/compile/abi-internal
// go tool asm


func CallBlk(asm *amd64.Assembler, framesize int32, f func()) {

    if framesize & 7 != 0 {
        panic("UNALIGNED FRAMESIZE")
    }

    // Requires manual sp adjustments

    asm.Sub(amd64.Imm{Val: framesize}, amd64.Rsp)

    f()

    asm.Add(amd64.Imm{Val: framesize}, amd64.Rsp)

}

func CallStack(off int32, bits uint8) amd64.Operand {
    return amd64.Indirect{Base: amd64.Rsp, Offset: off, Bits: bits}
}

func main() {

    // create

    asm, err := amd64.New(gojit.PageSize)
    if err != nil {
        panic(err)
    }

    asm.ABI = amd64.GoABI

    a := &adder{}

    ptr := uintptr(unsafe.Pointer(&a))

    fmt.Printf("PTR %X\n", ptr)



    //CallStack(asm, 0x0, 64)

    // emit

    // why min 0x10???
    CallBlk(asm, 0x18, func() {

        asm.Mov(amd64.Imm64{Val: int64(ptr)}, amd64.Rax)
        asm.Mov(amd64.Imm{Val: 0x69420}, amd64.Rbx)
        asm.Mov(amd64.Imm{Val: 0x7},     amd64.Rcx)
        //asm.CallFuncGo(add)
        asm.CallFuncGo((*adder).add)
    })


    asm.Ret() // make sure RAX is return value on amd64

    // build

    var fn func() int64
    asm.BuildTo(&fn)

    // call

    fmt.Printf("R %12X", fn())
}

type adder struct {
    v int64
}

func (a *adder) add(x, y int64) int64 {
    fmt.Printf("ADDER X %12X, Y %12X\n", x, y)

    a.v = x + y + 1

    return x + y
}

func add(x, y int64) int64 {
    fmt.Printf("X %12X, Y %12X\n", x, y)
    return x + y
}

