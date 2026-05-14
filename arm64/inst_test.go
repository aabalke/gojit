package arm64

import (
	"arm64/cache"
	"fmt"
	"testing"
	"unsafe"
)

var v uint32

func TestInst(t *testing.T) {

	var tests = []struct {
		Want uint32
		buildFunc func(a *Assembler)
	}{

		{12, _testAdd},
		{0xBEEF << 1, _testLsl},
		{0xBEEF >> 1, _testLsr},

	}

	for i, tt := range tests {
		testname := fmt.Sprintf("Test %d", i)
		t.Run(testname, func(t *testing.T) {

			asm, err := New(PageSize)
			if err != nil {
				panic(err)
			}

			tt.buildFunc(asm)

			asm.Ret()

			if err := asm.Error(); err != nil {
				panic(err)
			}

			cache.ClearICache(asm.Buf)

			CallJit(uintptr(unsafe.Pointer(&asm.Buf[0])))

			asm.Release()

			if v != tt.Want {
				t.Errorf("got %X, want %X", v, tt.Want)
			}
		})
	}
}

func _testAdd(a *Assembler) {

	v = 8

	a.Mov64(R10, uint64(uintptr(unsafe.Pointer(&v))))
	a.LdrImm(R00, R10, 0, SIZE_WORD, false, true)
	a.ADDImm(R00, R00, 4, false, false ,false)
	a.StrImm(R00, R10, 0, SIZE_WORD, false, true)
}

func _testLsl(a *Assembler) {

	v = 0xBEEF

	a.Mov64(R10, uint64(uintptr(unsafe.Pointer(&v))))
	a.LdrImm(R00, R10, 0, SIZE_WORD, false, true)
	a.LslImm(R00, R00, 1, false)
	a.StrImm(R00, R10, 0, SIZE_WORD, false, true)
}

func _testLsr(a *Assembler) {

	v = 0xBEEF

	a.Mov64(R10, uint64(uintptr(unsafe.Pointer(&v))))
	a.LdrImm(R00, R10, 0, SIZE_WORD, false, true)
	a.LsrImm(R00, R00, 1, false)
	a.StrImm(R00, R10, 0, SIZE_WORD, false, true)
}
