package main

import (
	"fmt"
	"unsafe"

	a "github.com/aabalke/gojit/arm64"
	"github.com/aabalke/gojit/arm64/cache"
)

var c = uint64(0x80)

func main() {

	pagesize := 512

	asm, err := a.New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.Mov64(a.R01, uint64(uintptr(unsafe.Pointer(&c))))
	asm.LDRImm(a.R02, a.R01, 0, false, true, true)

	// r3 == 4
	asm.Movz(a.R03, 0xFF, a.HW_00, true)
	asm.Movz(a.R04, 0x1F, a.HW_00, true)

	asm.AndReg(a.R03, a.R04, a.R03, 0, 0, true, true)

	//asm.ADDReg(a.R03, a.R03, a.R02, 0, 0, false, true)

	//asm.ADDImm(a.R03, a.R03, 0x04, false, true)

	//asm.ADDImm(a.R02, a.R03, 0x00, false, true)
	asm.STRImm(a.R03, a.R01, 0, false, true, true)
	asm.Ret()

	cache.ClearICache(asm.Buf)

	fmt.Printf("C %08X\n", c)

	a.CallJit(uintptr(unsafe.Pointer(&asm.Buf[0])))

	fmt.Printf("C %08X\n", c)

	asm.Release()
}
