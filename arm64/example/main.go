package main

import (
	"unsafe"

	a "github.com/aabalke/gojit/arm64"
)

func testExit(asm *a.Assembler) {

	asm.Ret()

    if asm.Error() != nil {
        panic(asm.Error())
    }

	a.CallJit(uintptr(unsafe.Pointer(&asm.Buf[0])))

	asm.Release()
}

var called = false

func main() {

    pagesize := 512

	asm, err := a.New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.InternalCallFunc(func() {
        called = true
	})

    testExit(asm)

    if !called {
        panic("Failed Test Call: called variable not set\n")
    }
}


