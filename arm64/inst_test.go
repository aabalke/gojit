package arm64

import (
	"fmt"
	"testing"
	"unsafe"
)

var testV uint64

func TestInst(t *testing.T) {

	var tests = []struct {
		Want uint64
		buildFunc func(a *Assembler)
	}{
		{0xBEEF, _testB},
		{0xDEAD, _testBCondFail},
		{0x0, _testBCondPass},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("Test %d", i)
		t.Run(testname, func(t *testing.T) {

			asm, err := New(PageSize)
			if err != nil {
				panic(err)
			}

			tt.buildFunc(asm)

			asm.Exit()

			CallJit(uintptr(unsafe.Pointer(&asm.Buf[0])))

			asm.Release()

			if testV != tt.Want {
				t.Errorf("got %X, want %X", testV, tt.Want)
			}
		})
	}
}

func _testBCondPass(a *Assembler) {

	a.Mov64(R10, uint64(uintptr(unsafe.Pointer(&testV))))
	a.LdrImm(R00, R10, 0, SIZE_WORD, false, true)

	a.Movz(R00, 0x0, 0, false)
	a.TstReg(R00, R00, 0, 0, false)
	jump := a.BCond(Z)
	a.Movz(R00, 0xDEAD, 0, false)
	jump()

	a.StrImm(R00, R10, 0, SIZE_WORD, false, true)
}

func _testBCondFail(a *Assembler) {

	a.Mov64(R10, uint64(uintptr(unsafe.Pointer(&testV))))
	a.LdrImm(R00, R10, 0, SIZE_WORD, false, true)

	a.Movz(R00, 0x1, 0, false)
	a.TstReg(R00, R00, 0, 0, false)
	jump := a.BCond(Z)
	a.Movz(R00, 0xDEAD, 0, false)
	jump()

	a.StrImm(R00, R10, 0, SIZE_WORD, false, true)
}

func _testB(a *Assembler) {

	testV = 8

	a.Mov64(R10, uint64(uintptr(unsafe.Pointer(&testV))))
	a.LdrImm(R00, R10, 0, SIZE_WORD, false, true)

	a.Movz(R00, 0xBEEF, 0, false)

	jump := a.B()
	a.Movz(R00, 0xDEAD, 0, false)
	jump()

	a.StrImm(R00, R10, 0, SIZE_WORD, false, true)
}

//func _testAdd(a *Assembler) {
//  panic("v != testV")
//	v = 8
//
//	a.Mov64(R10, uint64(uintptr(unsafe.Pointer(&v))))
//	a.LdrImm(R00, R10, 0, SIZE_WORD, false, true)
//	a.ADDImm(R00, R00, 4, false, false ,false)
//	a.StrImm(R00, R10, 0, SIZE_WORD, false, true)
//}
//
//func _testLsl(a *Assembler) {
//
//  panic("v != testV")
//	v = 0xBEEF
//
//	a.Mov64(R10, uint64(uintptr(unsafe.Pointer(&v))))
//	a.LdrImm(R00, R10, 0, SIZE_WORD, false, true)
//	a.LslImm(R00, R00, 1, false)
//	a.StrImm(R00, R10, 0, SIZE_WORD, false, true)
//}
//
//func _testLsr(a *Assembler) {
//
//  panic("v != testV")
//	v = 0xBEEF
//
//	a.Mov64(R10, uint64(uintptr(unsafe.Pointer(&v))))
//	a.LdrImm(R00, R10, 0, SIZE_WORD, false, true)
//	a.LsrImm(R00, R00, 1, false)
//	a.StrImm(R00, R10, 0, SIZE_WORD, false, true)
//}
